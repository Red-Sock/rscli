RSCLI_VERSION=v0.0.19-alpha
rscli-version:
	@echo $(RSCLI_VERSION)

#==============
# migrations
#==============

GOOSE_VERSION=$(shell goose -version)
MIG_DIR="migrations/"
goose-dep:
ifeq ("$(GOOSE_VERSION)", "")
	@echo "installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
else
	@echo "goose is installed!"
endif

#==============