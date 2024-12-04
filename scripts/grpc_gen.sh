protoc --go_out=./internal/daemon \
    --go-grpc_out=./internal/daemon \
    ./api/pm0.proto