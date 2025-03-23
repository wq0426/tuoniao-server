#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o liar_task cmd/task/main.go
chmod +x liar_task && scp liar_task root@8.130.22.115:/root/liar/app
