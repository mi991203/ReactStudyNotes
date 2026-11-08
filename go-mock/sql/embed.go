// Package sql 内嵌建表与种子 SQL，供 dao 初始化数据库时使用。
package sql

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed seed.sql
var Seed string
