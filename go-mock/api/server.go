// Package api 实现 protocols 中定义的 gRPC 接口（func (s *Server) Xxx）。
// 本层只做协议转换：pb 请求 → logic 调用 → pb 响应，并把业务错误映射为 gRPC status。
package api

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-mock/logic"
	"go-mock/model"
	commentv1 "go-mock/protocols/comment/v1"
)

// 前端展示格式：MM-DD HH:mm
const ctimeLayout = "01-02 15:04"

// Server 实现 commentv1.CommentServiceServer
type Server struct {
	commentv1.UnimplementedCommentServiceServer
	comment *logic.CommentLogic
}

func NewServer(comment *logic.CommentLogic) *Server {
	return &Server{comment: comment}
}

func (s *Server) ListComments(ctx context.Context, req *commentv1.ListCommentsRequest) (*commentv1.ListCommentsResponse, error) {
	items, err := s.comment.List(ctx, req.GetType())
	if err != nil {
		return nil, toStatus(err)
	}
	resp := &commentv1.ListCommentsResponse{List: make([]*commentv1.Comment, 0, len(items))}
	for _, it := range items {
		resp.List = append(resp.List, toPB(it))
	}
	return resp, nil
}

func (s *Server) CreateComment(ctx context.Context, req *commentv1.CreateCommentRequest) (*commentv1.Comment, error) {
	c, err := s.comment.Create(ctx, req.GetUid(), req.GetContent())
	if err != nil {
		return nil, toStatus(err)
	}
	return toPB(c), nil
}

func (s *Server) DeleteComment(ctx context.Context, req *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	if err := s.comment.Delete(ctx, req.GetRpid()); err != nil {
		return nil, toStatus(err)
	}
	return &commentv1.DeleteCommentResponse{}, nil
}

// toStatus 把业务错误映射为 gRPC status（gateway 会再转成 HTTP 状态码）。
func toStatus(err error) error {
	switch {
	case errors.Is(err, logic.ErrInvalidType),
		errors.Is(err, logic.ErrEmptyContent),
		errors.Is(err, logic.ErrEmptyUID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, logic.ErrUserNotFound),
		errors.Is(err, logic.ErrCommentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func toPB(c *model.Comment) *commentv1.Comment {
	return &commentv1.Comment{
		Rpid: c.RPID,
		User: &commentv1.User{
			Uid:    c.User.UID,
			Avatar: c.User.Avatar,
			Uname:  c.User.UName,
		},
		Content: c.Content,
		Ctime:   time.Unix(c.CTime, 0).Format(ctimeLayout),
		Like:    c.LikeCount,
	}
}
