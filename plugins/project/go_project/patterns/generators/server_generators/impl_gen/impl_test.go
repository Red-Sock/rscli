package impl_gen

import (
	"testing"

	"github.com/stretchr/testify/require"

	rscliconfig "github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io/folder"
	"github.com/Red-Sock/rscli/plugins/project/go_project/patterns"
	"github.com/Red-Sock/rscli/tests/mocks"
)

func TestGenerateImpl(t *testing.T) {
	projMock := mocks.NewIProjectMock(t)
	projMock.GetNameMock.Expect().Return("test_Proj")

	folders := &folder.Folder{
		Name: "",
		Inner: []*folder.Folder{
			{
				Name: "api",
				Inner: []*folder.Folder{
					{
						Name: patterns.GRPCServer,
						Inner: []*folder.Folder{
							patterns.ProtoContract.Copy(),
						},
					},
				},
			},
		},
	}

	projMock.GetFolderMock.
		Expect().
		Return(folders)

	cfg := &rscliconfig.RsCliConfig{
		Env: rscliconfig.Project{PathToServerDefinition: "api"},
	}

	out, err := GenerateImpl(cfg, projMock)
	require.NoError(t, err)
	_ = out
}
