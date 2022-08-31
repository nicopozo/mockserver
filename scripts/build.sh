#!/bin/bash

cd ..

docker build -t nicopozo/mock-service:2.1.0 .
docker build -t nicopozo/mock-service:latest .
docker push nicopozo/mock-service:2.1.0
docker push nicopozo/mock-service:latest