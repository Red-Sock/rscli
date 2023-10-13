package makefile

import (
	"bytes"
	"os"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/rscli/cmd/environment/project/compose/env"
)

var (
	ErrParsingMakefile     = errors.New("error parsing makefile")
	ErrMarshallingMakefile = errors.New("error marshalling makefile")
)

const phony = ".PHONY"

type Rule struct {
	Name      []byte
	PhonyName []byte
	Commands  [][]byte
	isInline  bool // flag showing that this part is for calling multiple other make rules
}

type Makefile struct {
	variables *env.Container
	rules     []Rule
}

func ReadMakeFile(pth string) (*Makefile, error) {
	makeFile, err := os.ReadFile(pth)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, errors.Wrap(err, "error reading makefile")
	}

	return NewMakeFile(makeFile)
}

func NewMakeFile(in []byte) (*Makefile, error) {
	m := &Makefile{
		variables: &env.Container{},
	}

	lines := bytes.Split(in, []byte{'\n'})

	for idx := 0; idx < len(lines); idx++ {
		if len(lines[idx]) == 0 {
			continue
		}

		l := lines[idx]

		if index := bytes.Index(l, []byte{'='}); index != -1 {
			m.variables.Append(parseVariable(index, l))
			continue
		}

		if index := bytes.Index(l, []byte{':'}); index != -1 {
			rule, offset, err := parseRule(lines[idx:])
			if err != nil {
				return nil, errors.Wrap(err, "error parsing rule")
			}
			m.rules = append(m.rules, rule)
			idx += offset
		}
	}
	return m, nil
}

func MewEmptyMakefile() *Makefile {
	return &Makefile{
		variables: &env.Container{},
	}
}

func (m *Makefile) GetRules() []Rule {
	return m.rules
}

func (m *Makefile) GetRuleByName(name string) *Rule {

	for _, item := range m.rules {
		if string(item.Name) == name {
			return &item
		}
	}

	return nil
}
func (m *Makefile) GetVars() *env.Container {
	return m.variables
}

func (m *Makefile) ContainsRule(name string) bool {
	for _, rule := range m.rules {
		if string(rule.Name) == name {
			return true
		}
	}

	return false
}

func (m *Makefile) AppendRule(rule Rule) {
	m.rules = append(m.rules, rule)
}

func (m *Makefile) Merge(external *Makefile) {
	for _, item := range external.GetVars().Content() {
		if !m.GetVars().ContainsByName(item.Name) {
			m.GetVars().Append(item)
		}
	}

	for _, item := range external.GetRules() {
		if !m.ContainsRule(string(item.Name)) {
			m.AppendRule(item)
		}
	}
}

func (m *Makefile) Marshal() ([]byte, error) {
	sb := bytes.Buffer{}

	sb.Write(m.variables.MarshalEnv())

	if sb.Len() != 0 {
		sb.WriteByte('\n')
		sb.WriteByte('\n')
	}

	for _, rule := range m.rules {
		if len(rule.PhonyName) != 0 {
			sb.WriteString(phony)
			sb.WriteByte(':')
			if rule.PhonyName[0] != ' ' {
				sb.WriteByte(' ')
			}

			sb.Write(rule.PhonyName)
			sb.WriteByte('\n')
		}

		if len(rule.Name) == 0 {
			return nil, errors.Wrap(ErrMarshallingMakefile, "no name provided for a rule")
		}

		sb.Write(rule.Name)
		sb.WriteByte(':')
		sb.WriteByte('\n')

		for _, r := range rule.Commands {
			if len(r) == 0 {
				return nil, errors.Wrap(ErrMarshallingMakefile, "empty command rule in "+string(rule.Name))
			}
			if r[0] != '\t' {
				sb.WriteByte('\t')
			}

			sb.Write(r)
			sb.WriteByte('\n')
		}

		sb.WriteByte('\n')
	}

	return sb.Bytes(), nil
}

func parseVariable(equalIndex int, b []byte) env.Variable {
	return env.Variable{
		Name:  string(b[:equalIndex]),
		Value: string(b[equalIndex+1:]),
	}
}

func parseRule(b [][]byte) (rule Rule, idx int, err error) {
	delimeterIdx := bytes.Index(b[idx], []byte{':'})
	if delimeterIdx == -1 {
		return rule, idx, errors.Wrap(err, "no \":\" symbol at the beginning of a make rule")
	}

	rule.Name = b[idx][:delimeterIdx]
	if string(rule.Name) == phony {
		rule.PhonyName = b[idx][delimeterIdx+1:]

		idx++
		delimeterIdx = bytes.Index(b[idx], []byte{':'})
		if delimeterIdx == -1 {
			return rule, idx, errors.Wrap(err, "rule must contain a name. "+
				"A proper format is \"rule-name:\", but \""+string(b[idx])+"\" is given")
		}
		rule.Name = b[idx][:delimeterIdx]
	}

	idx++

	for ; idx < len(b); idx++ {
		if !bytes.HasPrefix(b[idx], []byte("\t")) {
			return rule, idx, nil
		}
		rule.Commands = append(rule.Commands, b[idx])
	}

	return rule, idx, nil
}
