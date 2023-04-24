package patterns

var (
	MigrationsUtilityPrefix = []byte(`
#==============
# migrations
#==============
`)
	MigrationsUtility = []byte(`
GOOSE_VERSION=$(shell goose -version)
MIG_DIR="migrations/"
goose-dep:
ifeq ("$(GOOSE_VERSION)", "")
	@echo "installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
else
	@echo "goose is installed!"
endif
`)
)

var SectionSeparator = []byte(`
#==============
`)
