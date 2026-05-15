#!/bin/bash

cd ..

docker build -t nicopozo/mock-service:3.5.2 .
docker build -t nicopozo/mock-service:latest .
docker push nicopozo/mock-service:3.5.2
docker push nicopozo/mock-service:latest