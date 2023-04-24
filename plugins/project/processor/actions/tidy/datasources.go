package tidy

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	_consts "github.com/Red-Sock/rscli/plugins/config/pkg/const"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

func DataSources(p interfaces.Project) error {
	dataSources, err := p.GetConfig().GetDataSourceOptions()
	if err != nil {
		return errors.Wrap(err, "error extracting data sources from config")
	}

	makeFile := p.GetFolder().GetByPath(patterns.RsCliMkFileName)
	if makeFile == nil {
		return ErrNoMakeFile
	}

	sectionStartIdx := bytes.Index(makeFile.Content, patterns.MigrationsUtilityPrefix)
	sectionEndIdx := sectionStartIdx
	if sectionStartIdx != -1 {
		sectionEndIdx = bytes.Index(makeFile.Content[sectionStartIdx+len(patterns.MigrationsUtilityPrefix):], patterns.SectionSeparator)
		if sectionEndIdx == -1 {
			sectionEndIdx = len(makeFile.Content)
		} else {
			sectionEndIdx += sectionStartIdx + len(patterns.MigrationsUtilityPrefix) + len(patterns.SectionSeparator)
		}
	}

	migUpSection := `
mig-up:
`
	// migration up call
	for _, dsn := range dataSources {
		if dsn.Type == _consts.SourceNamePostgres {
			migUpSection += fmt.Sprintf(`
	@echo "applying migration on %s"
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=%s goose up
`, dsn.Name, dsn.ConnectionString)
		}
	}

	migUpSecBytes := bytes.Join([][]byte{
		patterns.MigrationsUtilityPrefix,
		patterns.MigrationsUtility,
		[]byte(migUpSection),
		patterns.SectionSeparator,
	}, []byte{})

	if sectionStartIdx == sectionEndIdx {
		makeFile.Content = append(makeFile.Content, migUpSecBytes...)
	} else {
		prev := makeFile.Content[:sectionStartIdx]
		post := makeFile.Content[sectionEndIdx:]

		makeFile.Content = bytes.Join([][]byte{
			prev,
			migUpSecBytes,
			post,
		}, []byte{})

	}

	return nil
}
