gen: deps pre-gen gen-client gen-server

deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

pre-gen:
	mkdir -p pkg/gen

gen-client:
	protoc --go_out=./pkg/gen \
	-I /Users/alexbukov/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
	--proto_path=. \
	./pkg/proto/clients/*.proto

gen-server:
	protoc --go-grpc_out=./pkg/gen \
	-I /Users/alexbukov/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
	--proto_path=. \
	./pkg/proto/*.proto