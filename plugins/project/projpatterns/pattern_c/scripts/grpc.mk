gen: deps gen-server

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

local-link:
	rm -rf api/google
	rm -rf api/validate
	ln -sf $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google api/google
	ln -sf $(GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.0.2/validate api/validate

gen-server: .pre-gen-server local-link .gen-server

.pre-gen-server:
	mkdir -p pkg/

.gen-server:
	protoc \
    	-I=./api \
    	-I $(GOPATH)/bin \
    	--grpc-gateway_out=logtostderr=true:./pkg/ \
    	--openapiv2_out ./api \
		--descriptor_set_out=./pkg/velez_api/api_descriptor.pb \
    	--go_out=./pkg/ \
    	--go-grpc_out=./pkg/ \
	    --validate_out="lang=go:./pkg" \
    	./api/grpc/*.proto
