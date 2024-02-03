protoc --go_out=bufs --go_opt=paths=source_relative \
    --go-grpc_out=bufs --go-grpc_opt=paths=source_relative \
    server.proto