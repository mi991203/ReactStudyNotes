// Package dao 封装 SQLite 访问。
package dao

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	mocksql "go-mock/sql"
)

// Open 打开（必要时创建）SQLite 数据库，并执行建表与种子数据（均为幂等）。
func Open(path string) (*sql.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	// SQLite 单写者，mock 场景用单连接避免 database is locked
	db.SetMaxOpenConns(1)

	if err := Init(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// Init 执行建表语句，并在 comment 表为空时导入种子数据。
func Init(db *sql.DB) error {
	if _, err := db.Exec(mocksql.Schema); err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM comment`).Scan(&n); err != nil {
		return fmt.Errorf("count comment: %w", err)
	}
	// 只在空库时导入种子，避免用户删除的数据在重启后复活
	if n == 0 {
		if _, err := db.Exec(mocksql.Seed); err != nil {
			return fmt.Errorf("exec seed: %w", err)
		}
	}
	return nil
}
