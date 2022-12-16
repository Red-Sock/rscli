package patterns

import (
	_ "embed"
)

// Build and deploy
var (
	//go:embed pattern_c/Dockerfile
	Dockerfile []byte
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

	//go:embed pattern_c/examples/http-client.env.json
	HttpEnvironment []byte
)
