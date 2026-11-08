# 数据库初始化文档

go-mock 使用本地 SQLite 文件作为存储，默认路径 `data/mock.db`（已被 `.gitignore` 忽略）。

## 1. 相关文件

| 文件 | 作用 |
| --- | --- |
| `sql/schema.sql` | 建表语句（`CREATE TABLE/INDEX IF NOT EXISTS`，可重复执行） |
| `sql/seed.sql` | 初始 mock 数据（`INSERT OR IGNORE`，可重复执行） |
| `sql/embed.go` | 用 `go:embed` 把上面两个 SQL 编进二进制 |
| `dao/db.go` | `Open` / `Init`：打开数据库并执行初始化 |

## 2. 表结构

### user（用户）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `uid` | TEXT PK | 用户 id |
| `uname` | TEXT NOT NULL | 昵称 |
| `avatar` | TEXT NOT NULL DEFAULT '' | 头像地址 |

### comment（评论）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `rpid` | INTEGER PK AUTOINCREMENT | 评论 id |
| `uid` | TEXT NOT NULL → `user.uid` | 发布者 |
| `content` | TEXT NOT NULL | 评论内容 |
| `like_count` | INTEGER NOT NULL DEFAULT 0 | 点赞数 |
| `ctime` | INTEGER NOT NULL | 发布时间，unix 秒；接口输出时格式化为 `MM-DD HH:mm` |
| `is_deleted` | INTEGER NOT NULL DEFAULT 0 | 软删除标记 |

索引：`idx_comment_ctime(ctime DESC)`、`idx_comment_like(like_count DESC)`，分别服务于 `type=time` 与 `type=hot` 排序。

### 建表 SQL

```sql
CREATE TABLE IF NOT EXISTS user (
  uid    TEXT PRIMARY KEY,
  uname  TEXT NOT NULL,
  avatar TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS comment (
  rpid       INTEGER PRIMARY KEY AUTOINCREMENT,
  uid        TEXT    NOT NULL REFERENCES user(uid),
  content    TEXT    NOT NULL,
  like_count INTEGER NOT NULL DEFAULT 0,
  ctime      INTEGER NOT NULL,
  is_deleted INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_comment_ctime ON comment(ctime DESC);
CREATE INDEX IF NOT EXISTS idx_comment_like  ON comment(like_count DESC);
```

## 3. 自动初始化（推荐）

直接启动服务即可：

```bash
make run          # 等价于 go run ./cmd/server -db data/mock.db
```

启动时 `dao.Open` 会依次：

1. 创建 `data/` 目录（不存在时）；
2. 打开（或创建）SQLite 文件，开启外键约束；
3. 执行 `schema.sql` 建表；
4. **仅当 `comment` 表为空时**导入 `seed.sql`。

> 第 4 步只在空库时执行，避免你通过 `DELETE` 接口删掉的数据在重启后又被种子数据"复活"。

## 4. 手动初始化（sqlite3 命令行）

不启动服务也可以直接建库：

```bash
make init-db      # 建表 + 导入种子，等价于下面两行
# sqlite3 data/mock.db < sql/schema.sql
# sqlite3 data/mock.db < sql/seed.sql
```

## 5. 重置数据库

```bash
make reset-db     # 删除 data/mock.db(-shm/-wal) 后重新 init-db
```

或者只删库文件，下次启动服务时会自动重建：

```bash
rm -f data/mock.db data/mock.db-shm data/mock.db-wal
```

## 6. 验证

```bash
sqlite3 data/mock.db '.tables'                                  # comment  user
sqlite3 data/mock.db 'SELECT COUNT(*) FROM comment;'            # 12
sqlite3 -header -column data/mock.db \
  "SELECT c.rpid, u.uname, c.content, c.like_count,
          datetime(c.ctime,'unixepoch','localtime') AS ctime
   FROM comment c JOIN user u ON u.uid = c.uid
   WHERE c.is_deleted = 0 ORDER BY c.ctime DESC;"
```

## 7. 初始数据

`sql/seed.sql` 初始化 5 个用户和 12 条评论（2023-11-11 ~ 2023-11-13），点赞数各不相同，
可直接验证 `type=hot`（点赞倒序）与 `type=time`（时间倒序）两种排序。

## 8. 修改 mock 数据

- 改初始数据：编辑 `sql/seed.sql`，然后 `make reset-db`（或删库后重启服务）。
- 改表结构：编辑 `sql/schema.sql`。由于使用 `IF NOT EXISTS`，对已存在的库不会生效，需要先重置数据库。
