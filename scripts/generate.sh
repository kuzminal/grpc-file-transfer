#!/bin/bash

PROJECT_ROOT=$(pwd)
API_DIR="${PROJECT_ROOT}/api/v1"

# Генерация для Go
protoc --go_out=${API_DIR} \
       --go-grpc_out=${API_DIR} \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       -I ${API_DIR} ${API_DIR}/*.proto

echo "Code generated successfully"