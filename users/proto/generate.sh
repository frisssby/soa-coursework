#!/bin/bash

protoc --go_out=./users --go-grpc_out=./users ./protos/tasks.proto
