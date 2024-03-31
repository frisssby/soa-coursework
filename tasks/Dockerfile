FROM golang:1.22-alpine

WORKDIR /tasks
COPY go.* .

RUN go mod download -x

COPY . .

RUN go build .

ENTRYPOINT ["./tasks"]
