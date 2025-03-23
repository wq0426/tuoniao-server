#!/bin/bash

REPLACE_IP=127.0.0.1:8289
echo "Running on Swag..."
swag init -g cmd/server/main.go
