FROM golang:latest

RUN apt update && apt install -y --assume-yes protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY ../protos /protos

WORKDIR /stats
RUN protoc --proto_path=/protos --go_out=. --go-grpc_out=. /protos/statistics.proto

COPY stats/go.* .
RUN go mod download -x

COPY stats/ .
RUN go build .

ENTRYPOINT ["./stats"]
