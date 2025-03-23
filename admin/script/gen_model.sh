#!/bin/bash

cd cmd/server/gorm_model && sh gen.sh flashbear $1
cat models/$1.go
rm -f models/$1.go