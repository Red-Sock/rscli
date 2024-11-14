package git

import (
	"strings"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/internal/cmd"
)

func Status(pth string) (uncommitted StatusDiff, err error) {
	executeOut, err := cmd.Execute(cmd.Request{
		Tool:    bin,
		Args:    []string{"status"},
		WorkDir: pth,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting git status")
	}

	out := make([]Changes, 0, 2)
	{
		const messageForUntrackedFiles = "Untracked files"

		startIdx := strings.Index(executeOut, messageForUntrackedFiles)
		if startIdx != -1 {

			startIdx += len(messageForUntrackedFiles)
			changeList := strings.Split(executeOut[startIdx:], "\n")

			gitChanges := Changes{
				Type:       ChangesTypeNotStaged,
				Changelist: make([]string, 0, len(changeList)),
			}
			for _, item := range changeList {
				if len(item) == 0 {
					continue
				}
				if strings.HasPrefix(item, "\t") {

					gitChanges.Changelist = append(gitChanges.Changelist, item[1:])
				}

			}

			out = append(out, gitChanges)
		}
	}

	{
		var keyWords = []string{"deleted", "modified", "new file"}

		for _, message := range []string{"Changes to be committed", "Changes not staged for commit"} {

			startIdx := strings.Index(executeOut, message)
			if startIdx != -1 {

				startIdx += len(message)
				changeList := strings.Split(executeOut[startIdx:], "\n")

				gitChanges := Changes{
					Type:       ChangesTypeNotCommitted,
					Changelist: make([]string, 0, len(changeList)),
				}

				for _, item := range changeList {

					if len(item) == 0 {
						continue
					}

					item = strings.ReplaceAll(item, "\t", "")
					for _, keyWord := range keyWords {
						if strings.HasPrefix(item, keyWord) {
							gitChanges.Changelist = append(gitChanges.Changelist, item)
							break
						}
					}

				}
				out = append(out, gitChanges)
			}
		}

	}

	return out, nil
}

type Changes struct {
	Type       gitChangesType
	Changelist []string
}

type StatusDiff []Changes

func (s StatusDiff) GetFilesListed() string {
	sb := strings.Builder{}
	for _, item := range s {
		if len(item.Changelist) == 0 {
			continue
		}
		sb.WriteString(item.Type.Msg())
		sb.WriteString("\n")

		changeTypeToFile := map[string][]string{}
		for _, line := range item.Changelist {
			splited := strings.Split(line, ":")
			if len(splited) != 2 {
				continue
			}
			changeTypeToFile[splited[0]] = append(changeTypeToFile[splited[0]], splited[1])
		}

		for k, v := range changeTypeToFile {
			sb.WriteString(k + ": \n\t")
			sb.WriteString(strings.Join(v, "; "))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
func (s StatusDiff) String() string {
	sb := strings.Builder{}
	for _, item := range s {
		if len(item.Changelist) == 0 {
			continue
		}
		sb.WriteString(item.Type.Msg())
		sb.WriteString("\n\t")
		sb.WriteString(strings.Join(item.Changelist, "\n\t"))
	}

	return sb.String()
}

type gitChangesType int

func (g gitChangesType) Msg() string {
	switch g {
	case ChangesTypeNotStaged:
		return "Not Staged changes:"
	case ChangesTypeNotCommitted:
		return "Not committed changes:"
	default:
		return "Unknown git changes type!"
	}
}
