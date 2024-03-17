protoc --go_out=./internal/daemon \
    --go-grpc_out=./internal/daemon \
    ./.proto/pm0.proto