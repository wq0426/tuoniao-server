#!/bin/sh

cd /go/src && export PATH=$PATH:/usr/local/go/bin && export GOPATH=/root/go && go build -o admin . && /go/src/admin -c config.docker.local.yaml
