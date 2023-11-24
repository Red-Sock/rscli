package environment

import (
	"github.com/spf13/cobra"

	"github.com/Red-Sock/rscli/internal/config"
	"github.com/Red-Sock/rscli/internal/io"
)

func newEnvInstallCmd(io io.IO, cfg *config.RsCliConfig) *cobra.Command {
	et := &envInstall{
		io:  io,
		cfg: cfg,
	}
	c := &cobra.Command{
		Use:   "tidy",
		Short: "Adds dependent binaries",

		RunE: et.RunInstall,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return c
}

type envInstall struct {
	io  io.IO
	cfg *config.RsCliConfig
}

func (e *envInstall) RunInstall(cmd *cobra.Command, arg []string) error {
	e.io.Println("Running rscli env install")
	// TODO install protoc
	// TODO install protoc-gen-go
	// TODO install docker desktop or anything else compatible
	return nil
}

func (e *envInstall) installMacOs() {

}
