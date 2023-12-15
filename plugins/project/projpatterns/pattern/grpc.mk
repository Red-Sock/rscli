gen: deps link-gw gen-server

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest

link-gw:
	ln -s $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google api/google
	ln -s $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway

gen-server: .pre-gen-server .gen-server

.pre-gen-server:
	mkdir -p pkg/

.gen-server:
	protoc \
        	-I=./api \
        	--grpc-gateway_out=logtostderr=true:./pkg/ \
        	--swagger_out=allow_merge=true,merge_file_name=api:./api \
    		--descriptor_set_out=./pkg/velez_api/api_descriptor.pb \
        	--go_out=./pkg/ \
        	--go-grpc_out=./pkg/ \
        	./api/grpc/*.proto

