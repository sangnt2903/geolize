FROM golang:latest

WORKDIR /

RUN apt-get update && \
    apt-get install -y protobuf-compiler

RUN GO111MODULE=on go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN GO111MODULE=on go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN GO111MODULE=on go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN GO111MODULE=on go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Add the installed binary to the PATH
ENV PATH="${PATH}:/go/bin"