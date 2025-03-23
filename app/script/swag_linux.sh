#!/bin/bash

REPLACE_IP=47.106.112.30:8086
echo "Running on Swag..."
swag init -g cmd/server/main.go && sed -i "" "s/127.0.0.1:8086/$REPLACE_IP/g" docs/swagger.yaml &&  sed -i "" "s/127.0.0.1:8086/$REPLACE_IP/g" docs/swagger.json &&  sed -i "" "s/127.0.0.1:8086/$REPLACE_IP/g" docs/docs.go
