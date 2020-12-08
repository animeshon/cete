#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

GATEWAY=$(go list -m -f "{{.Dir}}" github.com/grpc-ecosystem/grpc-gateway)
PROTO=./protobuf
CLIENT=./client

# This is the server.
protoc \
    --proto_path=. \
    --proto_path=${PROTO} \
    --proto_path=${GATEWAY} \
    --proto_path=${GATEWAY}/third_party/googleapis \
    --go_out=. \
    --go_opt=paths=source_relative \
        ${PROTO}/*.proto

# This is the server gateway.
protoc \
    --proto_path=. \
    --proto_path=${PROTO} \
    --proto_path=${GATEWAY} \
    --proto_path=${GATEWAY}/third_party/googleapis \
    --grpc-gateway_out=. \
    --grpc-gateway_opt=logtostderr=true \
    --grpc-gateway_opt=paths=source_relative \
    --grpc-gateway_opt=generate_unbound_methods=true \
    --grpc-gateway_opt=allow_delete_body=true \
        ${PROTO}/*.proto

# This is the client.
protoc \
    --proto_path=. \
    --proto_path=${PROTO} \
    --proto_path=${GATEWAY} \
    --proto_path=${GATEWAY}/third_party/googleapis \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
        ${PROTO}/*.proto

exit 0