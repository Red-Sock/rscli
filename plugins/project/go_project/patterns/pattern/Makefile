include rscli.mk

# generates folders and installs dependencies
warmup:
	make .prepare-grpc-folders
	make .deps-grpc
	PROTOPACKPATH=proto_deps protopack mod download
# generates code on warm project
codegen:
	PROTOPACKPATH=proto_deps protopack generate
	#cd {{ .NPM_PACKAGE_PATH }} && npm run build

lint:
	golangci-lint run ./...