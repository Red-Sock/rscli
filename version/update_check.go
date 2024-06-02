package version

import (
	"encoding/json"
	"net/http"

	v "github.com/hashicorp/go-version"
)

func CanUpdate() (string, bool) {
	resp, err := http.Get("https://api.github.com/repos/Red-Sock/rscli/releases/latest")
	if err != nil {
		return "", false
	}
	if resp.StatusCode != 200 {
		return "", false
	}
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return "", false
	}

	tag := m["tag_name"]
	if tag == "" {
		return "", false
	}
	tagStr, _ := tag.(string)
	if tagStr == "" {
		return "", false
	}

	originVersion, err := v.NewVersion(tagStr)
	if err != nil {
		return "", false
	}

	localVersion, err := v.NewVersion(version)
	if err != nil {
		return "", false
	}

	return tagStr, originVersion.GreaterThan(localVersion)
}
