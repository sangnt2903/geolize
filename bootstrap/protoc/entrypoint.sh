#!/bin/bash

mkdir -p generated

protoc --proto_path=defs \
 defs/"$SERVICE"/*.proto \
 --go_out=:generated \
 --go-grpc_out=require_unimplemented_servers=false:generated \
 --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
 --plugin=$(go env GOPATH)/bin/protoc-gen-openapiv2 \
 --plugin=protoc-gen-grpc-gateway=$(go env GOPATH)/bin/protoc-gen-grpc-gateway \
 --grpc-gateway_out generated \
 --grpc-gateway_opt generate_unbound_methods=true \
 --openapiv2_out generated