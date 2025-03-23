#!/bin/sh

if pgrep admin > /dev/null 2>&1
then
    pkill admin
else
    echo "Process admin is not running"
fi
/go/src/admin -c /go/src/config.docker.yaml