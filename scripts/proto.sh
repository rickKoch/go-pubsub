#!/bin/bash

set -e

protoc \
  --proto_path=api/protobuf "api/protobuf/pubsub.proto" \
  "--go_out=genproto/pubsub" --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=genproto/pubsub" --go-grpc_opt=paths=source_relative
