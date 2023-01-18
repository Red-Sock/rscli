package embeded

import (
	"github.com/Red-Sock/rscli/pkg/commands"
	"github.com/Red-Sock/rscli/pkg/flag"
	"github.com/Red-Sock/rscli/pkg/flag/flags"
	"github.com/pkg/errors"
	"os"
	"path"
)

type FixupPlugin struct{}

func (f *FixupPlugin) Run(flgs map[string][]string) error {
	pluginDir, err := flag.ExtractOneValueFromFlags(flgs, flags.PluginsDirFlag)
	if err != nil {
		return errors.Wrapf(err, "error extracting %s", flags.PluginsDirFlag)
	}

	if pluginDir != "" {
		return f.fix(pluginDir)
	}

	pluginDir, _ = os.LookupEnv(flags.PluginsDirEnv)
	if pluginDir != "" {
		return f.fix(pluginDir)
	}

	var exePath string
	exePath, err = os.Executable()
	if err == nil {
		pluginDir, _ = path.Split(exePath)
		pluginDir = path.Join(pluginDir, "rscli_plugins")

		return f.fix(path.Join(pluginDir))
	}

	return err
}

func (f *FixupPlugin) GetName() string {
	return commands.FixUtil
}

func (f *FixupPlugin) fix(pth string) (err error) {
	err = os.Setenv(flags.PluginsDirEnv, pth)
	if err != nil {
		return errors.Wrapf(err, "error setting environment variable %s to value %s", flags.PluginsDirEnv, pth)
	}

	if _, err = os.ReadDir(pth); err == nil {
		return nil
	}

	err = os.MkdirAll(pth, 0755)
	if err != nil {
		return errors.Wrap(err, "error creating directory for plugins")
	}

	return nil
}
