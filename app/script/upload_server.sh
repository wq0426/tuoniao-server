#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tuoliao cmd/server/main.go
chmod +x tuoliao && scp tuoliao root@8.130.22.115:/root/
