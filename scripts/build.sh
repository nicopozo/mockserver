#!/bin/bash

cd ..

docker build -t nicopozo/mock-service:3.4.0 .
docker build -t nicopozo/mock-service:latest .
docker push nicopozo/mock-service:3.4.0
docker push nicopozo/mock-service:latest