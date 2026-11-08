# go-mock

通用 Mock 服务：protobuf + gRPC + grpc-gateway，SQLite 存储。
目前提供 `comment-list` 评论案例所需的 `http://localhost:3004/list` 接口。

## 目录结构

```
protocols/        .proto 定义与 protoc 生成的桩代码（同目录，comment/v1/）
api/              接口实现：func (s *Server) ListComments 等，负责 pb 转换与错误映射
cmd/server/       入口：启动 gRPC(:9090) + HTTP gateway(:3004)，含 CORS
logic/            业务逻辑（校验、排序规则），不依赖 pb
dao/              SQLite 访问与数据库初始化
model/            表结构对应的 Go 结构体
sql/              schema.sql / seed.sql（go:embed）
docs/db-init.md   数据库初始化文档
third_party/      google/api/*.proto（http 注解）
```

## 快速开始

```bash
make run                       # 启动，首次运行自动建库并导入 12 条初始评论
curl localhost:3004/list       # 评论列表（裸数组，默认按时间倒序）
curl 'localhost:3004/list?type=hot'
curl -X POST localhost:3004/list -H 'Content-Type: application/json' \
     -d '{"uid":"30009257","content":"hello"}'
curl -X DELETE localhost:3004/list/1
```

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/list?type=hot\|time` | 评论列表，`hot` 按点赞、`time`(默认) 按时间 |
| POST | `/list` | 发布评论 `{uid, content}` |
| DELETE | `/list/{rpid}` | 删除评论（软删除） |

响应字段：`rpid, user{uid,avatar,uname}, content, ctime, like`。

## 修改 proto 后重新生成

`.proto` 在 `protocols/comment/v1/`，生成的 `*.pb.go` 也输出到同一目录。需要 `protoc`，首次先安装插件（装到 `./bin`）：

```bash
make tools
make proto
```

## 参数

```bash
go run ./cmd/server -db data/mock.db -http :3004 -grpc 127.0.0.1:9090 -cors-origin '*'
```

数据库初始化、重置详见 [docs/db-init.md](docs/db-init.md)。
