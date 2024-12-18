package impl_gen

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	errors "github.com/Red-Sock/trace-errors"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/internal/utils/cases"
	"github.com/Red-Sock/rscli/plugins/project"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
)

var (
	//go:embed template/impl.go.pattern
	implPattern  string
	implTemplate = template.Must(template.New("grpc_impl").Parse(implPattern))

	ErrNoGoPackageOption = errors.New("no Go package option")
)

type genArgs struct {
	FullProjPath string
	GrpcPackage  string
	ServiceName  string
}

func GenerateImpl(cfg *rscliconfig.RsCliConfig, proj project.IProject) ([]*folder.Folder, error) {
	grpcFolder := proj.GetFolder().GetByPath(cfg.Env.PathToServerDefinition, patterns.GRPCServer)
	if grpcFolder == nil {
		return nil, nil
	}

	out := make([]*folder.Folder, 0, 1)
	for _, f := range grpcFolder.Inner {
		if !strings.HasSuffix(f.Name, ".proto") {
			continue
		}

		stub, err := generateImpl(proj, f.Content)
		if err != nil {
			return nil, errors.Wrap(err, "error generating stub")
		}

		if stub != nil {
			out = append(out, stub)
		}
	}

	return out, nil
}

func generateImpl(proj project.IProject, protoContract []byte) (*folder.Folder, error) {
	out := bytes.Buffer{}

	args := genArgs{
		FullProjPath: proj.GetName(),
		ServiceName:  cases.ToPascal(extractServiceName(protoContract)),
	}

	var err error

	args.GrpcPackage, err = extractGoGrpcPackage(protoContract)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if proj.GetFolder().GetByPath(patterns.TransportFolder, args.GrpcPackage, "impl.go") != nil {
		return nil, nil
	}

	err = implTemplate.Execute(&out, args)
	if err != nil {
		return nil, errors.Wrap(err, "error executing template for grpc impl")
	}

	outF := &folder.Folder{
		Name: args.GrpcPackage + "_impl",
		Inner: []*folder.Folder{
			{
				Name:    "impl.go",
				Content: out.Bytes(),
			},
		},
	}
	return outF, nil
}

func extractGoGrpcPackage(contract []byte) (string, error) {
	const patternToFind = "option go_package ="
	idxStart := bytes.Index(contract, []byte(patternToFind))
	if idxStart == -1 {
		return "", errors.Wrap(ErrNoGoPackageOption)
	}

	endIdx := idxStart + bytes.Index(contract[idxStart:], []byte("\n"))

	goPackage := string(contract[idxStart+len(patternToFind) : endIdx-1])

	goPackage = strings.NewReplacer("\"", "", " ", "").Replace(goPackage)
	if goPackage[0] == '/' {
		goPackage = goPackage[1:]
	}
	return goPackage, nil
}

func extractServiceName(contract []byte) string {
	const patternToFind = "\nservice "
	startIdx := bytes.Index(contract, []byte(patternToFind))
	if startIdx == -1 {
		return ""
	}

	endIdx := startIdx + bytes.Index(contract[startIdx:], []byte("{"))
	serviceName := contract[startIdx+len(patternToFind) : endIdx]
	serviceName = bytes.TrimSpace(serviceName)
	return string(serviceName)
}
