#!/bin/bash

go run main.go -db=$1 -tb=$2 -dsn='root:EwqT42v1s2a78@tcp(127.0.0.1:3306)/'$1'?charset=utf8mb4&parseTime=True&loc=Local'
