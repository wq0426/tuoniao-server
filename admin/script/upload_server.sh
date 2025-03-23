#!/bin/sh

cd server \
  && scp admin root@47.106.112.30:/root/countrybattle/admin/server \
  && scp config.docker.yaml root@47.106.112.30:/root/countrybattle/admin/server