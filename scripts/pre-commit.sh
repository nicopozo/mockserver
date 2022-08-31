#!/bin/bash

swag init -g cmd/mocks/main.go -o ./api

go mod tidy

echo "###################################"
echo "Executing tests"
echo "###################################"
go test ./...

if [ $? -ne 0 ]
then
    echo 'There are tests in failure'
    exit 1
fi

echo ""
echo "###################################"
echo "Executing Lint"
echo "###################################"
golangci-lint run

if [ $? -ne 0 ]
then
    exit 1
fi

echo ""
echo "Build Success!!! You can commit your changes!!"