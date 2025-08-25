# Generate Go code into pkg/gen/
PROTO_FILES := $(shell find proto -name "*.proto")

generate-proto:
	protoc --proto_path=proto --go_out=pkg/gen --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/gen --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)
