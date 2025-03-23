#!/bin/sh

cd web \
		&& tar -zcf dist.tar.gz dist \
		&& scp ./dist.tar.gz root@47.106.112.30:/root/countrybattle/admin/dist.tar.gz \
		&& rm -rf dist.tar.gz dist