package embeded

import (
	"fmt"
	"github.com/Red-Sock/rscli/internal/utils/shared"
	"github.com/Red-Sock/rscli/pkg/commands"
	"net/url"
	"os"
	"path"
	"strings"
)

type DeletePlugin struct{}

func (d *DeletePlugin) Run(flgs map[string][]string) error {

	if len(flgs) != 1 {
		return fmt.Errorf("invalid amount of agruments for %s plugins. Expected %d got %d", commands.GetUtil, 1, len(flgs))
	}

	var repoURL string
	for k := range flgs {
		repoURL = k
	}

	allPluginsDir := shared.GetPluginsDir(flgs)

	u, err := url.Parse(repoURL)
	if err != nil {
		return err
	}

	err = os.RemoveAll(strings.ReplaceAll(path.Join(allPluginsDir, u.Host, u.Path), "*", ""))
	if err != nil {
		return err
	}

	return nil
}

func (d *DeletePlugin) GetName() string {
	return commands.Delete
}
