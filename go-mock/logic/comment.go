// Package logic 业务逻辑层：参数校验、排序规则等，不依赖 protobuf / gRPC。
package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-mock/dao"
	"go-mock/model"
)

// 业务错误，由 api 层映射为 gRPC status。
var (
	ErrInvalidType     = errors.New("type must be hot or time")
	ErrEmptyContent    = errors.New("content is required")
	ErrEmptyUID        = errors.New("uid is required")
	ErrUserNotFound    = errors.New("user not found")
	ErrCommentNotFound = errors.New("comment not found")
)

// CommentLogic 评论业务
type CommentLogic struct {
	dao *dao.CommentDAO
}

func NewCommentLogic(d *dao.CommentDAO) *CommentLogic {
	return &CommentLogic{dao: d}
}

// List 查询评论列表。typ 取值 ""(等同 time) / time / hot。
func (l *CommentLogic) List(ctx context.Context, typ string) ([]*model.Comment, error) {
	order := dao.OrderByTime
	switch typ {
	case "", "time":
	case "hot":
		order = dao.OrderByHot
	default:
		return nil, ErrInvalidType
	}
	return l.dao.List(ctx, order)
}

// Create 发布评论并返回完整评论（含用户信息）。
func (l *CommentLogic) Create(ctx context.Context, uid, content string) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrEmptyContent
	}
	if uid == "" {
		return nil, ErrEmptyUID
	}
	ok, err := l.dao.UserExists(ctx, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUserNotFound
	}

	id, err := l.dao.Create(ctx, uid, content, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return l.dao.Get(ctx, id)
}

// Delete 软删除评论。
func (l *CommentLogic) Delete(ctx context.Context, rpid int32) error {
	err := l.dao.Delete(ctx, rpid)
	if errors.Is(err, dao.ErrNotFound) {
		return ErrCommentNotFound
	}
	return err
}
