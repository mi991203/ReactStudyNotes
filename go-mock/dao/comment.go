package dao

import (
	"context"
	"database/sql"
	"errors"

	"go-mock/model"
)

// ErrNotFound 记录不存在
var ErrNotFound = errors.New("record not found")

// CommentDAO 评论数据访问
type CommentDAO struct {
	db *sql.DB
}

func NewCommentDAO(db *sql.DB) *CommentDAO {
	return &CommentDAO{db: db}
}

const selectComment = `
SELECT c.rpid, c.uid, c.content, c.like_count, c.ctime, c.is_deleted,
       u.uid, u.uname, u.avatar
FROM comment c
JOIN user u ON u.uid = c.uid`

// OrderBy 列表排序方式
type OrderBy int

const (
	OrderByTime OrderBy = iota // 发布时间倒序
	OrderByHot                 // 点赞数倒序
)

// List 查询未删除的评论
func (d *CommentDAO) List(ctx context.Context, order OrderBy) ([]*model.Comment, error) {
	q := selectComment + ` WHERE c.is_deleted = 0 `
	if order == OrderByHot {
		q += `ORDER BY c.like_count DESC, c.ctime DESC, c.rpid DESC`
	} else {
		q += `ORDER BY c.ctime DESC, c.rpid DESC`
	}

	rows, err := d.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// Get 按 rpid 查询未删除的评论
func (d *CommentDAO) Get(ctx context.Context, rpid int32) (*model.Comment, error) {
	row := d.db.QueryRowContext(ctx, selectComment+` WHERE c.rpid = ? AND c.is_deleted = 0`, rpid)
	c, err := scanComment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

// Create 新增评论，返回新评论 id
func (d *CommentDAO) Create(ctx context.Context, uid, content string, ctime int64) (int32, error) {
	res, err := d.db.ExecContext(ctx,
		`INSERT INTO comment (uid, content, like_count, ctime) VALUES (?, ?, 0, ?)`,
		uid, content, ctime)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int32(id), err
}

// Delete 软删除评论
func (d *CommentDAO) Delete(ctx context.Context, rpid int32) error {
	res, err := d.db.ExecContext(ctx,
		`UPDATE comment SET is_deleted = 1 WHERE rpid = ? AND is_deleted = 0`, rpid)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// UserExists 判断用户是否存在
func (d *CommentDAO) UserExists(ctx context.Context, uid string) (bool, error) {
	var n int
	err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM user WHERE uid = ?`, uid).Scan(&n)
	return n > 0, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanComment(s scanner) (*model.Comment, error) {
	var c model.Comment
	var deleted int
	if err := s.Scan(&c.RPID, &c.UID, &c.Content, &c.LikeCount, &c.CTime, &deleted,
		&c.User.UID, &c.User.UName, &c.User.Avatar); err != nil {
		return nil, err
	}
	c.IsDeleted = deleted != 0
	return &c, nil
}
