package patterns

import (
	_ "embed"
)

// Build and deploy
var (
	//go:embed pattern_c/Dockerfile
	Dockerfile []byte

	//go:embed pattern_c/.gitignore
	GitIgnore []byte

	//go:embed pattern_c/.golangci.yaml
	Linter []byte
)

// Documentation
var (
	//go:embed pattern_c/README.md
	Readme []byte
)

// Testing files
var (
	//go:embed pattern_c/examples/api.http
	ApiHTTP []byte
)
