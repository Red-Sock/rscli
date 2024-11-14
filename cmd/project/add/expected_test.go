package add

var (
	grpcExpectedProtoFile = []byte(`
syntax = "proto3";

package grpc_api;

import "google/protobuf/timestamp.proto";
import "google/api/annotations.proto";

option go_package = "/example_api";

service grpcAPI {
  rpc Version(PingRequest) returns (PingResponse) {
    option (google.api.http) = {
      post: "/api/version"
      body: "*"
    };
  };
}

message PingRequest {
  google.protobuf.Timestamp client_timestamp = 1;
}

message PingResponse {
   uint32 took = 1;
}
`)[1:]
	grpcMatreshkaConfigExpected = []byte(`
app_info:
    name: test_project
    version: v0.0.1
    startup_duration: 10s
servers:
    80:
        /{GRPC}:
            module: Test_AddDependency/GRPC
            gateway: /api
`)[1:]
)
