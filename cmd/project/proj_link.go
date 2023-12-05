package project

import (
	"bytes"
	stderrs "errors"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path"
	"sort"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/resources"
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/rw"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

var (
	modPkg = os.Getenv("GOPATH") + "/pkg/mod/"
)

type projectLink struct {
	io     io.IO
	config *config.RsCliConfig

	proj *project.Project
	path string
}

func newLinkCmd(pl projectLink) *cobra.Command {
	c := &cobra.Command{
		Use:   "link",
		Short: "Links other projects",
		Long:  `Can be used to link another project's contracts`,

		RunE: pl.run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	c.Flags().StringP(pathFlag, pathFlag[:1], "", `path to folder with project`)

	return c
}

func (p *projectLink) run(_ *cobra.Command, args []string) (err error) {
	if p.proj == nil {
		p.proj, err = project.LoadProject(p.path, p.config)
		if err != nil {
			return errors.Wrap(err, "error fetching project")
		}
	}
	succeeded, failed := p.getPackages(args)
	if len(failed) != 0 {
		p.io.Error("some packages weren't found: " + strings.Join(failed, ","))
	}

	var errs error
	for _, item := range succeeded {
		err := p.applyLink(item)
		if err != nil {
			errs = stderrs.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}

	err = go_actions.TidyAction{}.Do(p.proj)
	if err != nil {
		return errors.Wrap(err, "error tiding project")
	}

	return nil
}

func (p *projectLink) getPackages(args []string) (succeeded, failed []string) {
	for _, packageName := range args {
		if p.getPackage(packageName) {
			succeeded = append(succeeded, packageName)
		} else {
			failed = append(failed, packageName)
		}
	}

	return succeeded, failed
}

func (p *projectLink) getPackage(packageName string) (ok bool) {
	if !strings.Contains(packageName, "@") {
		packageName += "@latest"
	}

	_, err := cmd.Execute(cmd.Request{
		Tool: "go",
		Args: []string{"get", packageName},
	})
	if err != nil {
		return false
	}

	return true
}

func (p *projectLink) applyLink(packageName string) error {
	packagePath, err := p.getPathToModule(packageName)
	if err != nil {
		return errors.Wrapf(err, "error getting path to module: %s", packageName)
	}

	compiledPackages, err := p.getCompiledGRPCContractFromPackage(packagePath)
	if err != nil {
		return errors.Wrap(err, "error getting compiled grpc contracts")
	}

	const grpcConstructorPattern = `
package grpc
	 
import (
	"context"

	{{.PackageAlias}} "{{.ApiPath}}"
	errors "github.com/Red-Sock/trace-errors"

	"{{.FullProjectName}}/internal/config"
)

func {{.Constructor}}(ctx context.Context, cfg config.Config) ({{.PackageAlias}}.{{.ClientName}}, error) {
	connCfg, err := cfg.Resources().GRPC(config.{{.ConfigKey}})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find key"+ config.{{.ConfigKey}}+ " grpc connection in config")
	}

	conn, err := connect(ctx, connCfg)
	if err != nil {
		return nil, errors.Wrap(err, "error connection to "+connCfg.PackageName)
	}

	return {{.PackageAlias}}.{{.Constructor}}(conn), nil
}
`

	tmplt, err := template.New("grpcConstructor").Parse(grpcConstructorPattern)
	if err != nil {
		return errors.Wrap(err, "error parsing grpc constructor pattern")
	}

	grpcClientsFolder := p.proj.GetFolder().GetByPath(projpatterns.InternalFolder, p.config.Env.PathsToCompiledClients[0], projpatterns.GRPCServer)
	if grpcClientsFolder == nil {
		grpcClientsFolder = &folder.Folder{
			Name: path.Join(p.config.Env.PathsToClients[0], projpatterns.GRPCServer),
		}

		p.proj.GetFolder().Add(grpcClientsFolder)
	}
	sb := &rw.RW{}
	var errs error
	for _, c := range compiledPackages {
		resourceName := resources.GrpcResourceName + "_" + path.Base(packageName)

		_, err = p.proj.GetConfig().Resources.GRPC(resourceName)
		if err == nil {
			continue
		}

		if !errors.Is(err, matreshka.ErrResourceNotFound) {
			return errors.Wrap(err, "error getting grpc resource from config")
		}

		p.proj.GetConfig().Resources = append(p.proj.GetConfig().Resources, &resources.GRPC{
			Name:             resources.Name(resourceName),
			Module:           packageName,
			ConnectionString: "http://0.0.0.0:50051",
		})

		err = tmplt.Execute(sb, map[string]string{
			"PackageAlias":    "pb",
			"FullProjectName": p.proj.GetName(),
			"Constructor":     c.constructor,
			"ClientName":      c.clientName,
			"ConfigKey":       "ResourceGRPC" + strings.ReplaceAll(path.Base(packageName), "-", "_"), //TODO генерация на стороне матрёшки как  "ResourceGRPC" + base от пакета
			"ApiPath":         path.Join(packageName, c.importPath),
		})
		if err != nil {
			errs = stderrs.Join(errs, err)
			continue
		}

		grpcClientsFolder.Add(&folder.Folder{
			Name:    path.Base(packageName) + ".go",
			Content: sb.Bytes(),
		})
	}

	return nil
}

func (p *projectLink) getPathToModule(packageName string) (pathToModule string, err error) {
	packageName = p.filterPackageName(packageName)
	packagePath := path.Join(modPkg, packageName)

	if !strings.Contains(path.Base(packagePath), "@") {
		root := path.Dir(packagePath)
		potentialDirs, err := os.ReadDir(root)
		if err != nil {
			return "", errors.Wrap(err, "error reading potential packages paths")
		}
		baseName := path.Base(packageName)
		moveIdx := 0
		for idx := range potentialDirs {
			if !strings.HasPrefix(potentialDirs[idx].Name(), baseName) {
				potentialDirs[moveIdx], potentialDirs[idx] = potentialDirs[idx], potentialDirs[moveIdx]
				moveIdx++
			}

		}
		potentialDirs = potentialDirs[moveIdx:]

		sort.Slice(potentialDirs, func(i, j int) bool {
			return potentialDirs[i].Name() < potentialDirs[i].Name()
		})

		packagePath = path.Join(root, potentialDirs[0].Name())

	}

	return packagePath, nil
}

func (p *projectLink) getCompiledGRPCContractFromPackage(packagePath string) ([]*grpcPackage, error) {
	var errs error
	var packages []*grpcPackage

	for _, compiledClientPath := range p.config.Env.PathsToCompiledClients {
		apiPath := path.Join(packagePath, compiledClientPath)
		files, err := os.ReadDir(apiPath)
		if err != nil {
			errs = stderrs.Join(errs, err)
		}

		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			pkg, err := readGrpcPackage(packagePath, path.Join(compiledClientPath, file.Name()))
			if err != nil {
				errs = stderrs.Join(errs, err)
			}
			if pkg != nil {
				packages = append(packages, pkg)
			}

		}

	}

	return packages, errs
}

func (p *projectLink) filterPackageName(packageName string) string {
	packageNameB := []rune(packageName)
	out := make([]rune, 0, len(packageNameB))

	for idx := range packageNameB {
		if packageNameB[idx] >= 'A' && packageNameB[idx] <= 'Z' {
			packageNameB[idx] = packageNameB[idx] + 32
			out = append(out, '!')
		}

		out = append(out, packageNameB[idx])
	}

	return string(out)
}

type grpcPackage struct {
	importPath  string
	constructor string
	clientName  string
}

func readGrpcPackage(projectPath, apiContractPath string) (*grpcPackage, error) {
	packagePath := path.Join(projectPath, apiContractPath)
	files, err := os.ReadDir(packagePath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading contracts dir")
	}

	var clientContractPath string

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
		if strings.HasSuffix(fileName, "_grpc.pb.go") {
			clientContractPath = path.Join(packagePath, fileName)
			break
		}
	}

	if clientContractPath == "" {
		return nil, nil
	}

	clientContractB, err := os.ReadFile(clientContractPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading client contract files")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Base(clientContractPath), clientContractB, 0)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing go contract file")
	}

	out := &grpcPackage{}

	// TODO делать в горутине?
	ast.Inspect(f, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.GenDecl:
			if out.clientName != "" {
				break
			}
			if fn.Tok != token.TYPE {
				break
			}

			if len(fn.Specs) == 0 {
				break
			}
			spec, ok := fn.Specs[0].(*ast.TypeSpec)
			if !ok {
				break
			}

			if strings.HasSuffix(spec.Name.Name, "Client") {
				out.clientName = spec.Name.Name
			}
		case *ast.FuncDecl:
			if out.constructor != "" {
				break
			}

			if !strings.HasPrefix(fn.Name.Name, "New") || !strings.HasSuffix(fn.Name.Name, "Client") {
				return true
			}

			if fn.Type == nil || fn.Type.Params == nil || len(fn.Type.Params.List) != 1 || fn.Type.Params.List[0] == nil {

				return true
			}

			startIdx, endIdx := int(fn.Type.Params.List[0].Pos()), int(fn.Type.Params.List[0].End())
			if bytes.Contains(clientContractB[startIdx:endIdx], []byte("grpc.ClientConnInterface")) {
				out.constructor = fn.Name.Name
			}
		}

		return out.clientName == "" || out.constructor == ""
	})

	if out.constructor == "" {
		return nil, nil
	}

	{
		if !f.Package.IsValid() {
			return nil, nil
		}
		packageSubSet := clientContractB[f.Package:]
		packageSubSet = packageSubSet[:bytes.IndexByte(packageSubSet, '\n')]
		packageSubSet = packageSubSet[bytes.IndexByte(packageSubSet, ' ')+1:]
		packageName := string(packageSubSet)
		out.importPath = path.Join(path.Dir(apiContractPath), packageName)
	}

	return out, nil
}
