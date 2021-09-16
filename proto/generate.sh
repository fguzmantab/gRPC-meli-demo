#Generate GO code
protoc --go_out=./ ./proto/item.proto

#Generate Java code
protoc --java_out=./proto-generated/item-proto ./proto/item.proto

#Generate gRPC GO Code
protoc --go_out=./ --go-grpc_out=./ ./proto/item.proto

#Generate REST gateway
protoc -I . --grpc-gateway_out ./gen/go \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    ./proto/item.proto