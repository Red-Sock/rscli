package interfaces

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
	"github.com/Red-Sock/rscli/plugins/project/projpatterns"
)

const (
	TagVersionInvalid tagVersion = iota
	TagVersionAlpha
	TagVersionBeta
	TagVersionMain
)

const (
	tagVersionAlphaStr = "alpha"
	TagVersionBetaStr  = "beta"
	TagVersionMainStr  = ""
)

var mapStrVersionToInt = map[string]tagVersion{
	tagVersionAlphaStr: TagVersionAlpha,
	TagVersionBetaStr:  TagVersionBeta,
	TagVersionMainStr:  TagVersionMain,
}

var (
	ErrNoVersionVarInMkFile = errors.New("Cannot find rscli version env variable in rscli.mk file")
	ErrInvalidVersion       = errors.New("invalid project version, suppose to be like v0.1.2[-optional-version]")
)

type Version struct {
	Major, Minor, Negligible int
	Additional               tagVersion
}

func (v *Version) NeedUpdate(v2 Version) bool {
	if v2.Major != v.Major {
		return v2.Major < v.Major
	}

	if v.Minor != v2.Minor {
		return v2.Minor < v.Minor
	}

	if v.Negligible != v2.Negligible {
		return v2.Negligible < v.Negligible
	}

	if v.Additional != v2.Additional {
		return v2.Additional < v.Additional
	}

	return false
}

func (v *Version) String() string {
	var tag string
	switch v.Additional {
	case TagVersionAlpha:
		tag = "-alpha"
	case TagVersionBeta:
		tag = "-beta"
	}

	return "V" + strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Negligible) + tag
}

func (v *Version) UpdateProjectVersion(p Project) error {

	mkFile := p.GetFolder().GetByPath(projpatterns.RsCliMkFileName)

	rvBytes := []byte("RSCLI_VERSION=")
	startIdx := bytes.Index(mkFile.Content, rvBytes)
	if startIdx == -1 {
		return ErrNoVersionVarInMkFile
	}

	startIdx += len(rvBytes)

	endIdx := bytes.IndexByte(mkFile.Content[startIdx:], '\n')
	if endIdx == -1 {
		endIdx = len(mkFile.Content)
	} else {
		endIdx += startIdx
	}

	newVersion := []byte(v.String())

	out := make([]byte, len(mkFile.Content[:startIdx])+len(newVersion)+len(mkFile.Content[endIdx:]))
	copy(out[:startIdx], mkFile.Content[:startIdx])
	copy(out[startIdx:endIdx], newVersion)
	copy(out[endIdx:], mkFile.Content[endIdx:])

	mkFile.Content = out

	p.SetVersion(*v)

	return nil
}

type tagVersion int

func LoadProjectVersion(p Project) error {
	out, err := cmd.Execute(cmd.Request{
		Tool:    "make",
		Args:    []string{"-f", projpatterns.RsCliMkFileName, "rscli-version"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error executing make rscli-version")
	}

	out = strings.NewReplacer(
		"\r", "",
		"\n", "",
		"v", "",
		"V", "",
	).Replace(out)

	versionStr := strings.Split(out, ".")
	if len(versionStr) != 3 {
		return ErrInvalidVersion
	}

	var v Version
	v.Major, err = strconv.Atoi(versionStr[0])
	if err != nil {
		return err
	}

	v.Minor, err = strconv.Atoi(versionStr[1])
	if err != nil {
		return err
	}

	neglibleAndAdditional := strings.Split(versionStr[2], "-")
	if len(neglibleAndAdditional) == 2 {
		var ok bool
		v.Additional, ok = mapStrVersionToInt[neglibleAndAdditional[1]]
		if !ok {
			return errors.New("unknown additional version " + neglibleAndAdditional[1])
		}
	}

	v.Negligible, err = strconv.Atoi(neglibleAndAdditional[0])
	if err != nil {
		return err
	}

	p.SetVersion(v)

	return nil
}
