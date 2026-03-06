// Package biz 获客截流（Inbox）业务层：评论汇总、标记已处理、同步、回复/私信、动作记录的领域实体与用例。

package biz

import "context"

// XshComment 评论领域实体，对应表 xsh_comments。WantLink 由 AI 或规则标记。
type XshComment struct {
	ID                  uint64
	AccountID           uint64
	Platform            string
	PlatformContentID   string
	PlatformCommentID   string
	AuthorNickname      string
	Content             string
	WantLink            bool
	Status              string // pending | replied | ignored
	CreatedAt           int64
	UpdatedAt           int64
}

// XshInboxAction 私信/回复动作领域实体，对应表 xsh_inbox_actions。
type XshInboxAction struct {
	ID            uint64
	CommentID    uint64
	ActionType   string
	Content      string
	SentAt       int64
	ResultStatus string
	ErrorMessage string
	CreatedAt    int64
}

// XshCommentRepo 评论仓储接口。
type XshCommentRepo interface {
	List(ctx context.Context, accountID uint64, wantLink *bool, status string, page, pageSize int32) ([]*XshComment, int64, error)
	GetByID(ctx context.Context, id uint64) (*XshComment, error)
	Update(ctx context.Context, c *XshComment) error
	UpsertByPlatformCommentID(ctx context.Context, c *XshComment) error
}

// XshInboxActionRepo 动作记录仓储接口。
type XshInboxActionRepo interface {
	Create(ctx context.Context, a *XshInboxAction) (*XshInboxAction, error)
	ListByCommentID(ctx context.Context, commentID uint64, page, pageSize int32) ([]*XshInboxAction, int64, error)
}

// XshInboxUsecase 获客截流用例。
type XshInboxUsecase struct {
	commentRepo     XshCommentRepo
	actionRepo      XshInboxActionRepo
}

// NewXshInboxUsecase 创建获客截流用例。
func NewXshInboxUsecase(commentRepo XshCommentRepo, actionRepo XshInboxActionRepo) *XshInboxUsecase {
	return &XshInboxUsecase{
		commentRepo: commentRepo,
		actionRepo:  actionRepo,
	}
}

// ListComments 分页列表评论，支持 account_id、want_link、status 筛选。
func (uc *XshInboxUsecase) ListComments(ctx context.Context, accountID uint64, wantLink *bool, status string, page, pageSize int32) ([]*XshComment, int64, error) {
	return uc.commentRepo.List(ctx, accountID, wantLink, status, page, pageSize)
}

// GetComment 按 ID 查询评论。
func (uc *XshInboxUsecase) GetComment(ctx context.Context, id uint64) (*XshComment, error) {
	return uc.commentRepo.GetByID(ctx, id)
}

// MarkCommentHandled 标记评论已处理（replied/ignored）。
func (uc *XshInboxUsecase) MarkCommentHandled(ctx context.Context, id uint64, status string) error {
	c, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil || c == nil {
		return err
	}
	c.Status = status
	return uc.commentRepo.Update(ctx, c)
}

// SyncComments 同步评论；首期可由定时任务或外部 Worker 调用，此处仅返回占位成功。
func (uc *XshInboxUsecase) SyncComments(ctx context.Context, accountID uint64) (int32, string, error) {
	// 首期不实现真实拉取，返回 0 与提示；后续由 Playwright Worker 或定时任务拉取并写库。
	return 0, "首期未接入评论拉取", nil
}

// SendReply 记录一条回复动作；实际发送由 Playwright Worker 或外部服务完成。
func (uc *XshInboxUsecase) SendReply(ctx context.Context, commentID uint64, content string) (actionID uint64, err error) {
	a := &XshInboxAction{
		CommentID:   commentID,
		ActionType:  "reply",
		Content:     content,
		ResultStatus: "pending",
	}
	created, err := uc.actionRepo.Create(ctx, a)
	if err != nil {
		return 0, err
	}
	return created.ID, nil
}

// SendPrivateMessage 记录一条私信动作；实际发送由外部完成。
func (uc *XshInboxUsecase) SendPrivateMessage(ctx context.Context, commentID uint64, content string) (actionID uint64, err error) {
	a := &XshInboxAction{
		CommentID:   commentID,
		ActionType:  "private_message",
		Content:     content,
		ResultStatus: "pending",
	}
	created, err := uc.actionRepo.Create(ctx, a)
	if err != nil {
		return 0, err
	}
	return created.ID, nil
}

// ListInboxActions 按 comment_id 分页查询动作记录。
func (uc *XshInboxUsecase) ListInboxActions(ctx context.Context, commentID uint64, page, pageSize int32) ([]*XshInboxAction, int64, error) {
	return uc.actionRepo.ListByCommentID(ctx, commentID, page, pageSize)
}
