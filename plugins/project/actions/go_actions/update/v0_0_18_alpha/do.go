package v0_0_18_alpha

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/Red-Sock/rscli/plugins/project/interfaces"
	projpatterns "github.com/Red-Sock/rscli/plugins/project/patterns"
)

var Version = interfaces.Version{
	Major:      0,
	Minor:      0,
	Negligible: 18,
	Additional: interfaces.TagVersionAlpha,
}

func Do(p interfaces.Project) (err error) {
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
		idxStart := bytes.Index(projpatterns.RscliMK, migrationsSection)
		if idxStart == -1 {
			return fmt.Errorf("no migration section in rscli.mk source file")
		}
		sectionEndBytes := []byte("#==============")

		idxEnd := bytes.Index(projpatterns.RscliMK[idxStart+len(migrationsSection):], sectionEndBytes)
		if idxEnd == -1 {
			idxEnd = len(projpatterns.RscliMK)
		} else {
			idxEnd += idxStart + +len(migrationsSection) + len(sectionEndBytes)
		}
		mkFile := p.GetFolder().GetByPath(projpatterns.RsCliMkFileName)
		if mkFile == nil {
			return fmt.Errorf("no file %s", projpatterns.RsCliMkFileName)
		}

		mkFile.Content = append(mkFile.Content, projpatterns.RscliMK[idxStart:idxEnd]...)

	}

	return p.GetFolder().Build()
}
