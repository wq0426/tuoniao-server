#!/bin/sh

cd /go/src && export PATH=$PATH:/usr/local/go/bin && export GOPATH=/root/go && go build -o admin .
if pgrep admin > /dev/null 2>&1
then
    pkill admin
else
    echo "Process admin is not running"
fi
/go/src/admin -c /go/src/config.docker.yaml