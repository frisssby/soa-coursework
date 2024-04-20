FROM golang:latest

COPY ../protos /protos

WORKDIR /tasks

RUN apt update && apt install -y --assume-yes protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN protoc --proto_path=/protos --go_out=. --go-grpc_out=. /protos/tasks.proto

COPY tasks/go.* .
RUN go mod download -x

COPY tasks/ .
RUN go build .

ENTRYPOINT ["./tasks"]
