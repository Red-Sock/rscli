package v0_0_18_alpha

import (
	"bytes"
	"errors"
	"fmt"

	interfaces2 "github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
	"github.com/Red-Sock/rscli/plugins/project/processor/patterns"
)

var Version = interfaces2.Version{
	Major:      0,
	Minor:      0,
	Negligible: 18,
	Additional: interfaces2.TagVersionAlpha,
}

func Do(p interfaces2.Project) (err error) {
	defer func() {
		if err != nil {
			return
		}

		updErr := Version.UpdateProjectVersion(p)
		if updErr == nil {
			return
		}

		if err == nil {
			err = updErr
			return
		}

		err = errors.Join(err, updErr)
	}()
	// update migration file
	{
		migrationsSection := []byte("\n#==============\n# migrations\n#==============")
		idxStart := bytes.Index(patterns.RscliMK, migrationsSection)
		if idxStart == -1 {
			return fmt.Errorf("no migration section in rscli.mk source file")
		}
		sectionEndBytes := []byte("#==============")

		idxEnd := bytes.Index(patterns.RscliMK[idxStart+len(migrationsSection):], sectionEndBytes)
		if idxEnd == -1 {
			idxEnd = len(patterns.RscliMK)
		} else {
			idxEnd += idxStart + +len(migrationsSection) + len(sectionEndBytes)
		}
		mkFile := p.GetFolder().GetByPath(patterns.RsCliMkFileName)
		if mkFile == nil {
			return fmt.Errorf("no file %s", patterns.RsCliMkFileName)
		}

		mkFile.Content = append(mkFile.Content, patterns.RscliMK[idxStart:idxEnd]...)

	}

	return p.GetFolder().Build()
}
