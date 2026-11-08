GOBIN := $(CURDIR)/bin
export PATH := $(GOBIN):$(PATH)

PROTO_DIR := protocols
OUT_DIR   := .
DB_PATH   ?= data/mock.db

.PHONY: tools proto run build init-db reset-db tidy

# 安装 protoc 插件到 ./bin
tools:
	GOBIN=$(GOBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(GOBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(GOBIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# 生成 pb / grpc / gateway 桩代码到 protocols/comment/v1（与 .proto 同目录）
proto:
	protoc -I $(PROTO_DIR) -I third_party \
		--go_out=$(OUT_DIR) --go_opt=module=go-mock \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=module=go-mock \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=module=go-mock \
		$(PROTO_DIR)/comment/v1/comment.proto

tidy:
	go mod tidy

build:
	go build -o bin/go-mock ./cmd/server

run:
	go run ./cmd/server -db $(DB_PATH)

# 用 sqlite3 命令行手动初始化（服务启动时也会自动初始化）
init-db:
	mkdir -p $(dir $(DB_PATH))
	sqlite3 $(DB_PATH) < sql/schema.sql
	sqlite3 $(DB_PATH) < sql/seed.sql

# 删除数据库并重新初始化
reset-db:
	rm -f $(DB_PATH) $(DB_PATH)-shm $(DB_PATH)-wal
	$(MAKE) init-db
