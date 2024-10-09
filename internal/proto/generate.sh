#!/bin/bash

printf "removing old generated rpc...\n"
rm -v *.pb.go

printf "\ngenerating users rpc...\n"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./users.proto

printf "\ngenerating products rpc...\n"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./products.proto

printf "\ngenerating orders rpc...\n"
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./orders.proto
