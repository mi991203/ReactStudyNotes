// Package model 定义与数据库表一一对应的结构体。
package model

// User 对应 user 表
type User struct {
	UID    string
	UName  string
	Avatar string
}

// Comment 对应 comment 表（连表查询时带上发布者信息）
type Comment struct {
	RPID      int32
	UID       string
	Content   string
	LikeCount int32
	CTime     int64 // unix 秒
	IsDeleted bool

	// 连表 user 得到的发布者信息
	User User
}
