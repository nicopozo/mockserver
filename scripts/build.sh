#!/bin/bash

cd ..

docker build -t nicopozo/mock-service:3.5.1 .
docker build -t nicopozo/mock-service:latest .
docker push nicopozo/mock-service:3.5.1
docker push nicopozo/mock-service:latest