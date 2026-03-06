// Package biz 多账号分发矩阵（Dispatch）业务层：平台账号、定时规则、发布任务的领域实体与用例。
// 实际发布执行由 Playwright-go Worker 独立进程消费任务，本层仅管理任务与账号配置。

package biz

import "context"

// XshPlatformAccount 平台账号领域实体，对应表 xsh_platform_accounts。一账号一 IP（ProxyURL）。
type XshPlatformAccount struct {
	ID                 uint64
	Platform           string
	Nickname           string
	ProxyURL           string
	CookieStoragePath  string
	BindAt             int64
	CreatedAt          int64
	UpdatedAt          int64
}

// XshPublishTask 发布任务领域实体，对应表 xsh_publish_tasks。
// Status: pending | running | success | failed | cancelled。
type XshPublishTask struct {
	ID                  uint64
	DraftID             uint64
	AccountID           uint64
	ScheduledAt         int64
	PublishedAt        int64
	Status              string
	PlatformContentID   string
	ErrorMessage        string
	CreatedAt           int64
	UpdatedAt           int64
}

// XshPublishSchedule 定时规则领域实体，对应表 xsh_publish_schedules。
type XshPublishSchedule struct {
	ID          uint64
	Name        string
	CronExpr    string
	AccountIDs  string
	IsEnabled   bool
	CreatedAt   int64
	UpdatedAt   int64
}

// XshPlatformAccountRepo 平台账号仓储接口。
type XshPlatformAccountRepo interface {
	List(ctx context.Context, platform string, page, pageSize int32) ([]*XshPlatformAccount, int64, error)
	GetByID(ctx context.Context, id uint64) (*XshPlatformAccount, error)
	Create(ctx context.Context, a *XshPlatformAccount) (*XshPlatformAccount, error)
	Update(ctx context.Context, a *XshPlatformAccount) error
	Delete(ctx context.Context, id uint64) error
}

// XshPublishTaskRepo 发布任务仓储接口。
type XshPublishTaskRepo interface {
	GetByID(ctx context.Context, id uint64) (*XshPublishTask, error)
	List(ctx context.Context, page, pageSize int32, accountID uint64, status string) ([]*XshPublishTask, int64, error)
	Create(ctx context.Context, t *XshPublishTask) (*XshPublishTask, error)
	Update(ctx context.Context, t *XshPublishTask) error
	SetStatusRetry(ctx context.Context, id uint64) error
	SetStatusCancel(ctx context.Context, id uint64) error
}

// XshScheduleRepo 定时规则仓储接口。
type XshScheduleRepo interface {
	List(ctx context.Context) ([]*XshPublishSchedule, error)
	GetByID(ctx context.Context, id uint64) (*XshPublishSchedule, error)
	Create(ctx context.Context, s *XshPublishSchedule) (*XshPublishSchedule, error)
	Update(ctx context.Context, s *XshPublishSchedule) error
	Delete(ctx context.Context, id uint64) error
}

// XshDispatchUsecase 多账号分发用例。
type XshDispatchUsecase struct {
	accountRepo XshPlatformAccountRepo
	taskRepo    XshPublishTaskRepo
	scheduleRepo XshScheduleRepo
	draftRepo   XshDraftRepo
}

// NewXshDispatchUsecase 创建分发用例。
func NewXshDispatchUsecase(
	accountRepo XshPlatformAccountRepo,
	taskRepo XshPublishTaskRepo,
	scheduleRepo XshScheduleRepo,
	draftRepo XshDraftRepo,
) *XshDispatchUsecase {
	return &XshDispatchUsecase{
		accountRepo:  accountRepo,
		taskRepo:     taskRepo,
		scheduleRepo: scheduleRepo,
		draftRepo:    draftRepo,
	}
}

// ListPlatformAccounts 分页列表平台账号，可按 platform 筛选。
func (uc *XshDispatchUsecase) ListPlatformAccounts(ctx context.Context, platform string, page, pageSize int32) ([]*XshPlatformAccount, int64, error) {
	return uc.accountRepo.List(ctx, platform, page, pageSize)
}

// GetPlatformAccount 按 ID 查询账号。
func (uc *XshDispatchUsecase) GetPlatformAccount(ctx context.Context, id uint64) (*XshPlatformAccount, error) {
	return uc.accountRepo.GetByID(ctx, id)
}

// BindAccount 绑定新平台账号。
func (uc *XshDispatchUsecase) BindAccount(ctx context.Context, a *XshPlatformAccount) (*XshPlatformAccount, error) {
	return uc.accountRepo.Create(ctx, a)
}

// UpdateAccount 更新账号信息。
func (uc *XshDispatchUsecase) UpdateAccount(ctx context.Context, a *XshPlatformAccount) error {
	return uc.accountRepo.Update(ctx, a)
}

// UnbindAccount 解绑账号（删除记录）。
func (uc *XshDispatchUsecase) UnbindAccount(ctx context.Context, id uint64) error {
	return uc.accountRepo.Delete(ctx, id)
}

// ListSchedules 列表定时规则。
func (uc *XshDispatchUsecase) ListSchedules(ctx context.Context) ([]*XshPublishSchedule, error) {
	return uc.scheduleRepo.List(ctx)
}

// GetSchedule 按 ID 查询定时规则。
func (uc *XshDispatchUsecase) GetSchedule(ctx context.Context, id uint64) (*XshPublishSchedule, error) {
	return uc.scheduleRepo.GetByID(ctx, id)
}

// CreateSchedule 创建定时规则。
func (uc *XshDispatchUsecase) CreateSchedule(ctx context.Context, s *XshPublishSchedule) (*XshPublishSchedule, error) {
	return uc.scheduleRepo.Create(ctx, s)
}

// UpdateSchedule 更新定时规则。
func (uc *XshDispatchUsecase) UpdateSchedule(ctx context.Context, s *XshPublishSchedule) error {
	return uc.scheduleRepo.Update(ctx, s)
}

// DeleteSchedule 删除定时规则。
func (uc *XshDispatchUsecase) DeleteSchedule(ctx context.Context, id uint64) error {
	return uc.scheduleRepo.Delete(ctx, id)
}

// CreatePublishTasks 为 draft_id 与每个 account_id 创建一条发布任务，返回全部新建任务。
func (uc *XshDispatchUsecase) CreatePublishTasks(ctx context.Context, draftID uint64, accountIDs []uint64, scheduledAt int64) ([]*XshPublishTask, error) {
	var out []*XshPublishTask
	for _, accountID := range accountIDs {
		t := &XshPublishTask{
			DraftID:     draftID,
			AccountID:   accountID,
			ScheduledAt: scheduledAt,
			Status:      "pending",
		}
		created, err := uc.taskRepo.Create(ctx, t)
		if err != nil {
			return nil, err
		}
		out = append(out, created)
	}
	return out, nil
}

// ListPublishTasks 分页列表发布任务，支持 account_id、status 筛选。
func (uc *XshDispatchUsecase) ListPublishTasks(ctx context.Context, page, pageSize int32, accountID uint64, status string) ([]*XshPublishTask, int64, error) {
	return uc.taskRepo.List(ctx, page, pageSize, accountID, status)
}

// GetPublishTask 按 ID 查询发布任务。
func (uc *XshDispatchUsecase) GetPublishTask(ctx context.Context, id uint64) (*XshPublishTask, error) {
	return uc.taskRepo.GetByID(ctx, id)
}

// RetryPublishTask 将 status=failed 的任务置回 pending。
func (uc *XshDispatchUsecase) RetryPublishTask(ctx context.Context, id uint64) error {
	return uc.taskRepo.SetStatusRetry(ctx, id)
}

// CancelPublishTask 取消任务（置为 cancelled）。
func (uc *XshDispatchUsecase) CancelPublishTask(ctx context.Context, id uint64) error {
	return uc.taskRepo.SetStatusCancel(ctx, id)
}
