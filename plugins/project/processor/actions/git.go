package actions

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/Red-Sock/rscli/pkg/cmd"
	"github.com/Red-Sock/rscli/plugins/project/processor/interfaces"
)

func InitGit(p interfaces.Project) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"init"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return errors.Wrap(err, "error initiating git repository")
	}

	err = GitCommit(p.GetProjectPath(), "initializing project via rscli")
	if err != nil {
		return errors.Wrap(err, "error committing changes")
	}

	return nil
}

func GitCommit(pth, msg string) error {
	_, err := cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"add", "."},
		WorkDir: pth,
	})
	if err != nil {
		return errors.Wrap(err, "error adding files to git repository")
	}

	_, err = cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"commit", "-m", "\"" + msg + "\""},
		WorkDir: pth,
	})
	if err != nil {
		return errors.Wrap(err, "error committing files to git repository")
	}

	return nil
}

type gitChangesType int

const (
	GitChangesTypeInvalid = iota
	GitChangesTypeNotStaged
	GitChangesTypeNotCommitted
)

func (g gitChangesType) Msg() string {
	switch g {
	case GitChangesTypeNotStaged:
		return "Not Staged changes:"
	case GitChangesTypeNotCommitted:
		return "Not committed changes:"
	default:
		return "Unknown git changes type!"
	}
}

type Status []GitChanges

func (s Status) GetFilesListed() string {
	sb := strings.Builder{}
	for _, item := range s {
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
			sb.WriteString(strings.Join(v, "\n"))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
func (s Status) String() string {
	sb := strings.Builder{}
	for _, item := range s {
		sb.WriteString(item.Type.Msg())
		sb.WriteString("\n\t")
		sb.WriteString(strings.Join(item.Changelist, "\n\t"))
	}

	return sb.String()
}

type GitChanges struct {
	Type       gitChangesType
	Changelist []string
}

func GitStatus(p interfaces.Project) (uncommitted Status, err error) {
	executeOut, err := cmd.Execute(cmd.Request{
		Tool:    "git",
		Args:    []string{"status"},
		WorkDir: p.GetProjectPath(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting git status")
	}

	out := make([]GitChanges, 0, 2)
	{
		const messageForUntrackedFiles = "Untracked files"

		startIdx := strings.Index(executeOut, messageForUntrackedFiles)
		if startIdx != -1 {

			startIdx += len(messageForUntrackedFiles)
			changeList := strings.Split(executeOut[startIdx:], "\n")

			gitChanges := GitChanges{
				Type:       GitChangesTypeNotStaged,
				Changelist: make([]string, 0, len(changeList)),
			}
			for _, item := range changeList {
				if len(item) == 0 {
					continue
				}
				if strings.HasPrefix(item, "\n") {

					gitChanges.Changelist = append(gitChanges.Changelist, item)
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

				gitChanges := GitChanges{
					Type:       GitChangesTypeNotCommitted,
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
