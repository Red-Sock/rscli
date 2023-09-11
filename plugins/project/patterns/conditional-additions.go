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

var (
	GRPCSection = []byte(`
#==============
# grpc
#==============
`)
	GRPCUtilityInstallGoProtocHeader = []byte(`
protoc-dep: .install-protoc .install-protoc-gen-go .get-grpc-gateway
`)
	GRPCInstallProtocViaGolangEnvOSBased = []byte(`
GOOS_NAME=$(shell go env go env GOOS)
.install-protoc:
ifeq ("$(GOOS_NAME)", "darwin")
	brew install protobuf
endif

ifeq ("$(GOOS_NAME)", "windows")
	choco install protoc -y
endif

ifeq ("$(GOOS_NAME)", "linux")
ifeq ("$(RSCLI_LINUX_CONFIRM_INSTALL)","true")
	PROTOC_VERSION=$(curl -s "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
	curl -Lo protoc.zip "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
	sudo unzip -q protoc.zip bin/protoc -d /usr/local
	sudo chmod a+x /usr/local/bin/protoc
	rm -rf protoc.zip
else
	@echo "this code has never been properly tested on linux. In order to execute run with environment variable RSCLI_LINUX_CONFIRM_INSTALL set to true"
endif
	@echo "unknown OS"
endif
`)

	GRPCInstallGolangProtoc = []byte(`
.install-protoc-gen-go:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
`)

	GRPCGenerateGoCode = []byte(`
.generate-proto:
	protoc --go_out=./pkg --go_opt=paths=source_relative --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative -I pkg/proto/ pkg/proto/*/*.proto
`)

	GRPCGenerateGoCodeWithDependencies = []byte(`
generate-proto: protoc-dep .generate-proto 
`)

	GRPCGatewayDependency = []byte(`
.get-grpc-gateway:	
	git clone -b v1 https://github.com/grpc-ecosystem/grpc-gateway && \
    mv grpc-gateway/third_party/googleapis/google pkg/proto/google/ && \
    rm -rf grpc-gateway
`)
)
