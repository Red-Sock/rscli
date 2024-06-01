gen-server-grpc: .pre-gen-server-grpc .deps-grpc .gen-server-grpc

.deps-grpc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

	rm -rf api/google
	rm -rf api/validate
	ln -sf $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google api/google
	ln -sf $(GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.0.2/validate api/validate
.pre-gen-server-grpc:
	mkdir -p pkg/

.gen-server-grpc:
	protoc \
        	-I=./api \
        	--grpc-gateway_out=logtostderr=true:./pkg/ \
        	--swagger_out=allow_merge=true,merge_file_name=api:./api \
    		--descriptor_set_out=./pkg/api_discriptor.pb \
        	--go_out=./pkg/. \
        	--go-grpc_out=./pkg/. \
        	./api/grpc/*.proto
