SRC_DIR = ./fiber-proto
DST_DIR = ./protobuf
PACKAGES = eth,types,api

CLIENT_DIR = ./protobuf

.PHONY:
proto:
	mkdir -p $(CLIENT_DIR)/{api,eth,types}
	protoc -I=$(SRC_DIR) --go_opt=paths=source_relative \
		--go_opt=Mtypes.proto=github.com/chainbound/fiber-go/protobuf/types \
		--go_out=$(CLIENT_DIR)/types \
		types.proto
	protoc -I=$(SRC_DIR) --go_opt=paths=source_relative \
		--go_opt=Mtypes.proto=github.com/chainbound/fiber-go/protobuf/types \
		--go_opt=Meth.proto=github.com/chainbound/fiber-go/protobuf/eth \
		--go_out=$(CLIENT_DIR)/eth \
		eth.proto
	protoc -I=$(SRC_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--go_opt=Mtypes.proto=github.com/chainbound/fiber-go/protobuf/types \
		--go_opt=Meth.proto=github.com/chainbound/fiber-go/protobuf/eth \
		--go_opt=Mapi.proto=github.com/chainbound/fiber-go/protobuf/api \
		--go-grpc_opt=Mtypes.proto=github.com/chainbound/fiber-go/protobuf/types \
		--go-grpc_opt=Meth.proto=github.com/chainbound/fiber-go/protobuf/eth \
		--go-grpc_opt=Mapi.proto=github.com/chainbound/fiber-go/protobuf/api \
		--go_out=$(CLIENT_DIR)/api \
		--go-grpc_out=$(CLIENT_DIR)/api \
		api.proto


# protoc -I=$(SRC_DIR) --go_opt=Mapi.proto=github.com/chainbound/fiber/protobuf/api --go_out=./protobuf/api --go-grpc_out=./protobuf/api api.proto
.PHONY:
clean:
	rm -rf ./protobuf/{$(PACKAGES)}