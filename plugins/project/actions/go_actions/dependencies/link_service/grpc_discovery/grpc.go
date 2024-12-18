package grpc_discovery

import (
	"bytes"
	stderrs "errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"go.redsock.ru/rerrors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
)

type GrpcDiscovery struct {
	Cfg *rscliconfig.RsCliConfig
}

var defaultDiscoverer GrpcDiscovery

func DiscoverPackage(packageName string) (*GrpcPackage, error) {
	if defaultDiscoverer.Cfg == nil {
		defaultDiscoverer.Cfg = rscliconfig.GetConfig()
	}

	return defaultDiscoverer.DiscoverPackage(packageName)
}

func (g GrpcDiscovery) DiscoverPackage(packageName string) (*GrpcPackage, error) {
	pathToPackageInMod, err := GetPathToGlobalModule(packageName)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting path to package in global module")
	}

	grpcPackage, err := g.getGrpcPackageFromMod(pathToPackageInMod)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting grpc package from mod")
	}

	return grpcPackage, nil
}

type GrpcPackage struct {
	ImportPath  string
	Constructor string
	ClientName  string
}

func (g GrpcDiscovery) getGrpcPackageFromMod(packagePath string) (*GrpcPackage, error) {
	var errs error

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
			var pkg *GrpcPackage
			pkg, err = readGrpcPackageFromPackageClientPath(packagePath, path.Join(compiledClientPath, file.Name()))
			if err != nil {
				return nil, rerrors.Wrap(err, "error reading package from path")
			}

			if pkg != nil {
				packagePath = packagePath[len(modFolderPath):]
				atIdx := strings.Index(packagePath, "@")
				if atIdx != -1 {
					packagePath = packagePath[:atIdx]
				}

				pkg.ImportPath = path.Join(packagePath, pkg.ImportPath)
				return pkg, nil
			}
		}
	}

	return nil, nil
}

func readGrpcPackageFromPackageClientPath(projectPath, apiContractPath string) (*GrpcPackage, error) {
	packagePath := path.Join(projectPath, apiContractPath)
	files, err := os.ReadDir(packagePath)
	if err != nil {
		return nil, rerrors.Wrap(err, "error reading contracts dir")
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
		return nil, rerrors.Wrap(err, "error reading client contract files")
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path.Base(clientContractPath), clientContractB, 0)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing go contract file")
	}

	out := &GrpcPackage{}

	// TODO делать в горутине?
	ast.Inspect(f, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.GenDecl:
			if out.ClientName != "" {
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
				out.ClientName = spec.Name.Name
			}
		case *ast.FuncDecl:
			if out.Constructor != "" {
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
				out.Constructor = fn.Name.Name
			}
		}

		return out.ClientName == "" || out.Constructor == ""
	})

	if out.Constructor == "" {
		return nil, nil
	}

	if !f.Package.IsValid() {
		return nil, nil
	}
	packageSubSet := clientContractB[f.Package:]
	packageSubSet = packageSubSet[:bytes.IndexByte(packageSubSet, '\n')]
	packageSubSet = packageSubSet[bytes.IndexByte(packageSubSet, ' ')+1:]
	packageName := string(packageSubSet)
	out.ImportPath = path.Join(path.Dir(apiContractPath), packageName)

	return out, nil
}
