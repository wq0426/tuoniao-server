all: swag run
all_linux: swag_linux build-linux upload

.PHONY: swag
swag:
	/bin/bash ./script/swag.sh

.PHONY: swag_linux
swag_linux:
	/bin/bash ./script/swag_linux.sh

.PHONY: restart-web
restart-web:
	source ~/.bashrc && stopweb && upweb

.PHONY: restart-server
restart-server:
	source ~/.bashrc && stopcountryserver && upcountryserver

.PHONY: build-web
build-web:
	#docker-compose stop gva-web && docker-compose rm -f gva-web
	cd web && cnpm run build

.PHONY: build-server
build-server:
	cd server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o admin .

.PHONY: upload-web
upload-web:
	/bin/bash ./script/upload_web.sh

.PHONY: upload-server
upload-server:
	/bin/bash ./script/upload_server.sh

.PHONY: clean
clean:
	rm -rf flashbear