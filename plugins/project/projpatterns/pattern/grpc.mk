gen: deps link-gw gen-server

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest

link-gw:
	ln -sfn $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google api/google

gen-server: .pre-gen-server .gen-server

.pre-gen-server:
	mkdir -p pkg/api

.gen-server:
	protoc \
        	-I=./api \
        	--grpc-gateway_out=logtostderr=true:./pkg/api/ \
        	--swagger_out=allow_merge=true,merge_file_name=api:./api \
    		--descriptor_set_out=./pkg/api/api_discriptor.pb \
        	--go_out=./pkg/api/. \
        	--go-grpc_out=./pkg/api/. \
        	./api/grpc/*.proto

