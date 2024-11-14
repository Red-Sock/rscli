package grpc_discovery

import (
	"strings"

	"github.com/Red-Sock/rscli/internal/cmd"
)

func FetchPackage(packageName string) (ok bool) {
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
