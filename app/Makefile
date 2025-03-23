all: swag run
all_linux: swag_linux build-linux upload upload_task

.PHONY: wire
wire:
	nunu wire

.PHONY: swag
swag:
	/bin/bash ./script/swag.sh

.PHONY: swag_linux
swag_linux:
	/bin/bash ./script/swag_linux.sh

.PHONY: run
run:
	make swag && nunu run

.PHONY: run_task
run_task:
	make swag && go build -o countrybattle_task ./cmd/task/main.go && ./countrybattle_task
.PHONY: upload
upload:
	make swag_linux && /bin/bash ./script/upload_server.sh

.PHONY: upload_task
upload_task:
	make swag && /bin/bash ./script/upload_task.sh

.PHONY: proto
proto:
	protoc --go_out=./internal/grpc --go-grpc_out=./internal/grpc ./internal/grpc/push_message.proto

.PHONY: clean
clean:
	rm -rf countrybattle