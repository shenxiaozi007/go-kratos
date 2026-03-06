// Package service 获客截流（Inbox）HTTP 服务：实现 api/xsh/v1 InboxService。
// 评论列表与标记、同步占位、回复/私信动作记录；实际拉评与发送由 Worker 或外部完成。

package service

import (
	"context"

	xshv1 "realworld/api/xsh/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// XshInboxService 实现 xsh.v1.InboxService 的 HTTP 服务。
type XshInboxService struct {
	uc  *biz.XshInboxUsecase
	log *log.Helper
}

// NewXshInboxService 创建获客截流 HTTP 服务实例。
func NewXshInboxService(uc *biz.XshInboxUsecase, logger log.Logger) *XshInboxService {
	return &XshInboxService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListComments 分页列表评论，支持 account_id、want_link、status 筛选。
func (s *XshInboxService) ListComments(ctx context.Context, req *xshv1.ListCommentsRequest) (*xshv1.ListCommentsReply, error) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	var wantLink *bool
	if req.GetWantLink() {
		v := true
		wantLink = &v
	}
	items, total, err := s.uc.ListComments(ctx, req.GetAccountId(), wantLink, req.GetStatus(), page, pageSize)
	if err != nil {
		s.log.Errorf("ListComments failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.Comment, 0, len(items))
	for _, c := range items {
		list = append(list, &xshv1.Comment{
			Id:                  c.ID,
			AccountId:           c.AccountID,
			Platform:            c.Platform,
			PlatformContentId:   c.PlatformContentID,
			PlatformCommentId:   c.PlatformCommentID,
			AuthorNickname:      c.AuthorNickname,
			Content:             c.Content,
			WantLink:            c.WantLink,
			Status:              c.Status,
			CreatedAt:           c.CreatedAt,
			UpdatedAt:           c.UpdatedAt,
		})
	}
	return &xshv1.ListCommentsReply{Items: list, Total: total}, nil
}

// MarkCommentHandled 标记评论已处理（replied/ignored）。
func (s *XshInboxService) MarkCommentHandled(ctx context.Context, req *xshv1.MarkCommentHandledRequest) (*xshv1.MarkCommentHandledReply, error) {
	if err := s.uc.MarkCommentHandled(ctx, req.GetId(), req.GetStatus()); err != nil {
		s.log.Errorf("MarkCommentHandled failed: %v", err)
		return nil, err
	}
	return &xshv1.MarkCommentHandledReply{Success: true}, nil
}

// SyncComments 同步评论；首期返回占位成功与提示。
func (s *XshInboxService) SyncComments(ctx context.Context, req *xshv1.SyncCommentsRequest) (*xshv1.SyncCommentsReply, error) {
	count, msg, err := s.uc.SyncComments(ctx, req.GetAccountId())
	if err != nil {
		s.log.Errorf("SyncComments failed: %v", err)
		return nil, err
	}
	return &xshv1.SyncCommentsReply{Success: true, SyncedCount: count, Message: msg}, nil
}

// SendReply 记录一条回复动作；实际发送由 Worker 完成。
func (s *XshInboxService) SendReply(ctx context.Context, req *xshv1.SendReplyRequest) (*xshv1.SendReplyReply, error) {
	actionID, err := s.uc.SendReply(ctx, req.GetCommentId(), req.GetContent())
	if err != nil {
		s.log.Errorf("SendReply failed: %v", err)
		return nil, err
	}
	return &xshv1.SendReplyReply{Success: true, ActionId: actionID}, nil
}

// SendPrivateMessage 记录一条私信动作；实际发送由 Worker 完成。
func (s *XshInboxService) SendPrivateMessage(ctx context.Context, req *xshv1.SendPrivateMessageRequest) (*xshv1.SendPrivateMessageReply, error) {
	actionID, err := s.uc.SendPrivateMessage(ctx, req.GetCommentId(), req.GetContent())
	if err != nil {
		s.log.Errorf("SendPrivateMessage failed: %v", err)
		return nil, err
	}
	return &xshv1.SendPrivateMessageReply{Success: true, ActionId: actionID}, nil
}

// ListInboxActions 按 comment_id 分页查询动作记录。
func (s *XshInboxService) ListInboxActions(ctx context.Context, req *xshv1.ListInboxActionsRequest) (*xshv1.ListInboxActionsReply, error) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	items, total, err := s.uc.ListInboxActions(ctx, req.GetCommentId(), page, pageSize)
	if err != nil {
		s.log.Errorf("ListInboxActions failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.InboxAction, 0, len(items))
	for _, a := range items {
		list = append(list, &xshv1.InboxAction{
			Id:            a.ID,
			CommentId:    a.CommentID,
			ActionType:   a.ActionType,
			Content:      a.Content,
			SentAt:       a.SentAt,
			ResultStatus: a.ResultStatus,
			ErrorMessage: a.ErrorMessage,
			CreatedAt:     a.CreatedAt,
		})
	}
	return &xshv1.ListInboxActionsReply{Items: list, Total: total}, nil
}
