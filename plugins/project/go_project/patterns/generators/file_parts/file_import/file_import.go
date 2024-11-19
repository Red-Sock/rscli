package file_import

import (
	"path"
)

type FileImports interface {
	Build() map[string]string
	AddImport(importPath string) (alias string)
	AddWithAlias(alias, importPath string) (newAlias string)
}

type fileImports struct {
	importsToAliases map[string]string
	aliasToImports   map[string]string
}

func New() FileImports {
	return &fileImports{
		importsToAliases: make(map[string]string),
		aliasToImports:   make(map[string]string),
	}
}

func (f *fileImports) AddImport(importPath string) (alias string) {
	alias, ok := f.importsToAliases[importPath]
	if ok {
		return alias
	}
	alias = f.extractAlias(importPath)

	f.importsToAliases[importPath] = alias
	f.aliasToImports[alias] = importPath

	return alias
}

func (f *fileImports) AddWithAlias(alias, importPath string) (newAlias string) {
	oldAlias, ok := f.importsToAliases[importPath]
	if ok {
		return oldAlias
	}

	f.importsToAliases[importPath] = alias
	f.aliasToImports[alias] = importPath
	return alias
}

func (f *fileImports) Build() map[string]string {
	return f.importsToAliases
}

func (f *fileImports) extractAlias(importPath string) string {
	base := path.Base(importPath)
	return base
}
