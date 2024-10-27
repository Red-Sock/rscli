package dependencies

import (
	"bytes"
	stderrs "errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"sort"
	"strings"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka"
	"github.com/godverv/matreshka/resources"

	"github.com/Red-Sock/rscli/internal/cmd"
	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/actions/go_actions/dependencies/grpc_discovery"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators"
	"github.com/Red-Sock/rscli/plugins/project/go_project/projpatterns/generators/config_generators"
)

var (
	modPkg = os.Getenv("GOPATH") + "/pkg/mod/"
)

type GrpcClient struct {
	Modules []string

	Cfg *rscliconfig.RsCliConfig
	Io  io.IO
}

func (g GrpcClient) GetFolderName() string {
	return "grpc"
}

func (g GrpcClient) AppendToProject(proj project.Project) error {
	if len(g.Modules) == 0 {
		return nil
	}

	succeeded, failed := g.getPackages(g.Modules)
	if len(failed) != 0 {
		g.Io.Error("some packages weren't found: " + strings.Join(failed, ","))
	}

	var errs error
	for _, item := range succeeded {
		err := g.applyLink(proj, item)
		if err != nil {
			errs = stderrs.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}

	grpcClientConnFilePath := path.Join(g.Cfg.Env.PathsToClients[0], projpatterns.GRPCServer, projpatterns.ConnFileName)
	if proj.GetFolder().GetByPath(grpcClientConnFilePath) == nil {
		proj.GetFolder().Add(projpatterns.GrpcClientConnFile.CopyWithNewName(grpcClientConnFilePath))
	}

	return nil
}

func (g GrpcClient) getPackages(args []string) (succeeded, failed []string) {
	for _, packageName := range args {
		if g.getPackage(packageName) {
			succeeded = append(succeeded, packageName)
		} else {
			failed = append(failed, packageName)
		}
	}

	return succeeded, failed
}

func (g GrpcClient) getPackage(packageName string) (ok bool) {
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

// Users/alexbukov/go/pkg/mod/github.com/godverv/hello_world/pkg
func (g GrpcClient) applyLink(proj project.Project, packageName string) error {
	discovery := grpc_discovery.GrpcDiscovery{Cfg: g.Cfg}

	pkg, err := discovery.DiscoverPackage(packageName)
	if err != nil {
		return errors.Wrap(err, "error discovering package")
	}

	if pkg == nil {
		return nil
	}

	if idx := strings.Index(packageName, "@"); idx > -1 {
		packageName = packageName[:idx]
	}
	grpcPkgPath := path.Join(g.Cfg.Env.PathsToClients[0], projpatterns.GRPCServer)

	grpcClientsFolder := proj.GetFolder().GetByPath(grpcPkgPath)
	if grpcClientsFolder == nil {
		grpcClientsFolder = &folder.Folder{
			Name: grpcPkgPath,
		}

		proj.GetFolder().Add(grpcClientsFolder)
	}

	resourceName := resources.GrpcResourceName + "_" + generators.NormalizeResourceName(path.Base(packageName))

	grpcResource, err := proj.GetConfig().DataSources.GRPC(resourceName)
	if err != nil {
		if !errors.Is(err, matreshka.ErrNotFound) {
			return errors.Wrap(err, "error getting grpc resource from config")
		}

		grpcResource = &resources.GRPC{
			Name:             resources.Name(resourceName),
			Module:           packageName,
			ConnectionString: "0.0.0.0:50051",
		}
		proj.GetConfig().DataSources = append(proj.GetConfig().DataSources, grpcResource)
	}

	grpcClientFile, err := config_generators.GenerateGRPCClient(*pkg)
	if err != nil {
		return errors.Wrap(err, "error generating grpc client")
	}

	grpcClientsFolder.Add(
		&folder.Folder{
			Name:    path.Base(packageName) + ".go",
			Content: grpcClientFile,
		})

	return nil
}

func (g GrpcClient) getPathToModule(packageName string) (pathToModule string, err error) {
	packageName = g.filterPackageName(packageName)
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
			// TODO sorting by name is wrong. need to sort by version
			return potentialDirs[i].Name() < potentialDirs[i].Name()
		})

		packagePath = path.Join(root, potentialDirs[0].Name())

	}

	return packagePath, nil
}

func (g GrpcClient) getCompiledGRPCContractFromPackage(packagePath string) ([]*grpcPackage, error) {
	var errs error
	var packages []*grpcPackage

	for _, compiledClientPath := range g.Cfg.Env.PathsToCompiledClients {
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

func (g GrpcClient) filterPackageName(packageName string) string {
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
