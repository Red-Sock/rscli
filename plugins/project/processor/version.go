package processor

import (
	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func GetCurrentVersion() *Version {
	return &Version{
		major:      0,
		minor:      0,
		negligible: 8,
		additional: tagVersionAlpha,
	}
}

type Version struct {
	major, minor, negligible int
	additional               tagVersion
}

func (v *Version) IsOlderThan(v2 *Version) bool {
	if v2.major != v.major {
		return v2.major > v.major
	}
	if v.minor != v2.minor {
		return v2.minor > v.minor
	}
	if v.negligible != v2.negligible {
		return v2.negligible > v.negligible
	}

	if v.additional != v2.additional {
		return v2.additional > v.additional
	}

	return false
}

func (v *Version) String() string {
	var tag string
	switch v.additional {
	case tagVersionAlpha:
		tag = "-alpha"
	case tagVersionBeta:
		tag = "-beta"
	}

	return "V" + strconv.Itoa(v.major) + "." + strconv.Itoa(v.minor) + "." + strconv.Itoa(v.negligible) + tag
}

type tagVersion int

const (
	tagVersionAlpha tagVersion = iota
	tagVersionBeta
	tagVersionMain
)

func LoadProjectVersion(p *Project) error {
	out, err := cmd.Execute(cmd.Request{
		Tool:    "make",
		Args:    []string{"-f", patterns.RsCliMkFileName, "rscli-version"},
		WorkDir: p.ProjectPath,
	})
	if err != nil {
		return errors.Wrap(err, "error executing make rscli-version")
	}

	out = strings.NewReplacer(
		"\r", "",
		"\n", "",
	).Replace(out)

	return nil
}
