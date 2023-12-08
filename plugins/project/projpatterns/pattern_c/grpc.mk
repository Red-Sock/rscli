gen: deps gen-server

deps:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

gen-server: .pre-gen-server .gen-server
.pre-gen-server:
	mkdir -p pkg/

.gen-server:
	protoc --go_out=./pkg/ --go-grpc_out=./pkg/ \
	--proto_path=. \
	./api/*.proto


### TODO вынести в кодген
.generate-certs:
	openssl req -x509 \
	-newkey rsa:4096 \
	-keyout key.pem \
	-out cert.pem \
	-sha256 \
	-days 3650 \
	-nodes \
	-subj "/C=RE/ST=RSCLI_EXAMPLE/L=RSCLI_EXAMPLE/O=RSCLI_EXAMPLE/OU=RSCLI_EXAMPLE/CN=RSCLI_EXAMPLE"

.generate-server:
	openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=RE/ST=RSCLI_EXAMPLE/L=RSCLI_EXAMPLE/O=RSCLI_EXAMPLE/OU=RSCLI_EXAMPLE/CN=RSCLI_EXAMPLE/emailAddress=RSCLI_EXAMPLE"
