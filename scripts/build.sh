#!/bin/bash

cd ..
swag init -g cmd/mocks/main.go
go mod tidy
cd cmd/mocks
go build