// Package service 多账号分发矩阵（Dispatch）HTTP 服务：实现 api/xsh/v1 DispatchService。
// 平台账号、定时规则、发布任务 CRUD；实际发布由 Playwright-go Worker 消费任务执行。

package service

import (
	"context"

	xshv1 "realworld/api/xsh/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// XshDispatchService 实现 xsh.v1.DispatchService 的 HTTP 服务。
type XshDispatchService struct {
	uc  *biz.XshDispatchUsecase
	log *log.Helper
}

// NewXshDispatchService 创建分发 HTTP 服务实例。
func NewXshDispatchService(uc *biz.XshDispatchUsecase, logger log.Logger) *XshDispatchService {
	return &XshDispatchService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListPlatformAccounts 分页列表平台账号。
func (s *XshDispatchService) ListPlatformAccounts(ctx context.Context, req *xshv1.ListPlatformAccountsRequest) (*xshv1.ListPlatformAccountsReply, error) {
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
	items, total, err := s.uc.ListPlatformAccounts(ctx, req.GetPlatform(), page, pageSize)
	if err != nil {
		s.log.Errorf("ListPlatformAccounts failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.PlatformAccount, 0, len(items))
	for _, a := range items {
		list = append(list, accountToProto(a))
	}
	return &xshv1.ListPlatformAccountsReply{Items: list, Total: total}, nil
}

// BindAccount 绑定新平台账号；仅 admin 可操作（v2 权限）。
func (s *XshDispatchService) BindAccount(ctx context.Context, req *xshv1.BindAccountRequest) (*xshv1.BindAccountReply, error) {
	if u := auth.FromContext(ctx); u == nil || u.Role != "admin" {
		return nil, errors.New(403, "FORBIDDEN", "admin only")
	}
	a := &biz.XshPlatformAccount{
		Platform:          req.GetPlatform(),
		Nickname:          req.GetNickname(),
		ProxyURL:          req.GetProxyUrl(),
		CookieStoragePath: req.GetCookieStoragePath(),
	}
	created, err := s.uc.BindAccount(ctx, a)
	if err != nil {
		s.log.Errorf("BindAccount failed: %v", err)
		return nil, err
	}
	return &xshv1.BindAccountReply{Account: accountToProto(created)}, nil
}

// UpdateAccount 更新账号；请求中非空字段覆盖已有值。
func (s *XshDispatchService) UpdateAccount(ctx context.Context, req *xshv1.UpdateAccountRequest) (*xshv1.UpdateAccountReply, error) {
	existing, err := s.uc.GetPlatformAccount(ctx, req.GetId())
	if err != nil || existing == nil {
		return nil, err
	}
	if req.GetNickname() != "" {
		existing.Nickname = req.GetNickname()
	}
	if req.GetProxyUrl() != "" {
		existing.ProxyURL = req.GetProxyUrl()
	}
	if req.GetCookieStoragePath() != "" {
		existing.CookieStoragePath = req.GetCookieStoragePath()
	}
	if err := s.uc.UpdateAccount(ctx, existing); err != nil {
		s.log.Errorf("UpdateAccount failed: %v", err)
		return nil, err
	}
	return &xshv1.UpdateAccountReply{Account: accountToProto(existing)}, nil
}

// UnbindAccount 解绑账号（删除记录）；仅 admin 可操作（v2 权限）。
func (s *XshDispatchService) UnbindAccount(ctx context.Context, req *xshv1.UnbindAccountRequest) (*xshv1.UnbindAccountReply, error) {
	if u := auth.FromContext(ctx); u == nil || u.Role != "admin" {
		return nil, errors.New(403, "FORBIDDEN", "admin only")
	}
	if err := s.uc.UnbindAccount(ctx, req.GetId()); err != nil {
		s.log.Errorf("UnbindAccount failed: %v", err)
		return nil, err
	}
	return &xshv1.UnbindAccountReply{Success: true}, nil
}

// ListSchedules 列表定时规则。
func (s *XshDispatchService) ListSchedules(ctx context.Context, _ *xshv1.ListSchedulesRequest) (*xshv1.ListSchedulesReply, error) {
	list, err := s.uc.ListSchedules(ctx)
	if err != nil {
		s.log.Errorf("ListSchedules failed: %v", err)
		return nil, err
	}
	items := make([]*xshv1.PublishSchedule, 0, len(list))
	for _, v := range list {
		items = append(items, &xshv1.PublishSchedule{
			Id:         v.ID,
			Name:       v.Name,
			CronExpr:   v.CronExpr,
			AccountIds: v.AccountIDs,
			IsEnabled:  v.IsEnabled,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}
	return &xshv1.ListSchedulesReply{Items: items}, nil
}

// CreateSchedule 创建定时规则。
func (s *XshDispatchService) CreateSchedule(ctx context.Context, req *xshv1.CreateScheduleRequest) (*xshv1.CreateScheduleReply, error) {
	sch := &biz.XshPublishSchedule{
		Name:       req.GetName(),
		CronExpr:   req.GetCronExpr(),
		AccountIDs: req.GetAccountIds(),
		IsEnabled:  req.GetIsEnabled(),
	}
	created, err := s.uc.CreateSchedule(ctx, sch)
	if err != nil {
		s.log.Errorf("CreateSchedule failed: %v", err)
		return nil, err
	}
	return &xshv1.CreateScheduleReply{Schedule: &xshv1.PublishSchedule{
		Id:         created.ID,
		Name:       created.Name,
		CronExpr:   created.CronExpr,
		AccountIds: created.AccountIDs,
		IsEnabled:  created.IsEnabled,
		CreatedAt:  created.CreatedAt,
		UpdatedAt:  created.UpdatedAt,
	}}, nil
}

// UpdateSchedule 更新定时规则。
func (s *XshDispatchService) UpdateSchedule(ctx context.Context, req *xshv1.UpdateScheduleRequest) (*xshv1.UpdateScheduleReply, error) {
	existing, err := s.uc.GetSchedule(ctx, req.GetId())
	if err != nil || existing == nil {
		return nil, err
	}
	if req.GetName() != "" {
		existing.Name = req.GetName()
	}
	if req.GetCronExpr() != "" {
		existing.CronExpr = req.GetCronExpr()
	}
	if req.GetAccountIds() != "" {
		existing.AccountIDs = req.GetAccountIds()
	}
	existing.IsEnabled = req.GetIsEnabled()
	if err := s.uc.UpdateSchedule(ctx, existing); err != nil {
		s.log.Errorf("UpdateSchedule failed: %v", err)
		return nil, err
	}
	return &xshv1.UpdateScheduleReply{Schedule: &xshv1.PublishSchedule{
		Id:         existing.ID,
		Name:       existing.Name,
		CronExpr:   existing.CronExpr,
		AccountIds: existing.AccountIDs,
		IsEnabled:  existing.IsEnabled,
		CreatedAt:  existing.CreatedAt,
		UpdatedAt:  existing.UpdatedAt,
	}}, nil
}

// DeleteSchedule 删除定时规则。
func (s *XshDispatchService) DeleteSchedule(ctx context.Context, req *xshv1.DeleteScheduleRequest) (*xshv1.DeleteScheduleReply, error) {
	if err := s.uc.DeleteSchedule(ctx, req.GetId()); err != nil {
		s.log.Errorf("DeleteSchedule failed: %v", err)
		return nil, err
	}
	return &xshv1.DeleteScheduleReply{Success: true}, nil
}

// CreatePublishTask 为 draft_id 与每个 account_id 创建一条发布任务，返回全部新建任务。
func (s *XshDispatchService) CreatePublishTask(ctx context.Context, req *xshv1.CreatePublishTaskRequest) (*xshv1.CreatePublishTaskReply, error) {
	accountIDs := req.GetAccountIds()
	if len(accountIDs) == 0 {
		return &xshv1.CreatePublishTaskReply{Tasks: nil}, nil
	}
	tasks, err := s.uc.CreatePublishTasks(ctx, req.GetDraftId(), accountIDs, req.GetScheduledAt())
	if err != nil {
		s.log.Errorf("CreatePublishTask failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.PublishTask, 0, len(tasks))
	for _, t := range tasks {
		list = append(list, taskToProto(t))
	}
	return &xshv1.CreatePublishTaskReply{Tasks: list}, nil
}

// ListPublishTasks 分页列表发布任务。
func (s *XshDispatchService) ListPublishTasks(ctx context.Context, req *xshv1.ListPublishTasksRequest) (*xshv1.ListPublishTasksReply, error) {
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
	items, total, err := s.uc.ListPublishTasks(ctx, page, pageSize, req.GetAccountId(), req.GetStatus())
	if err != nil {
		s.log.Errorf("ListPublishTasks failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.PublishTask, 0, len(items))
	for _, t := range items {
		list = append(list, taskToProto(t))
	}
	return &xshv1.ListPublishTasksReply{Items: list, Total: total}, nil
}

// RetryPublishTask 将 failed 任务置回 pending。
func (s *XshDispatchService) RetryPublishTask(ctx context.Context, req *xshv1.RetryPublishTaskRequest) (*xshv1.RetryPublishTaskReply, error) {
	if err := s.uc.RetryPublishTask(ctx, req.GetId()); err != nil {
		s.log.Errorf("RetryPublishTask failed: %v", err)
		return nil, err
	}
	return &xshv1.RetryPublishTaskReply{Success: true}, nil
}

// CancelPublishTask 取消任务。
func (s *XshDispatchService) CancelPublishTask(ctx context.Context, req *xshv1.CancelPublishTaskRequest) (*xshv1.CancelPublishTaskReply, error) {
	if err := s.uc.CancelPublishTask(ctx, req.GetId()); err != nil {
		s.log.Errorf("CancelPublishTask failed: %v", err)
		return nil, err
	}
	return &xshv1.CancelPublishTaskReply{Success: true}, nil
}

func accountToProto(a *biz.XshPlatformAccount) *xshv1.PlatformAccount {
	if a == nil {
		return nil
	}
	return &xshv1.PlatformAccount{
		Id:                 a.ID,
		Platform:           a.Platform,
		Nickname:           a.Nickname,
		ProxyUrl:           a.ProxyURL,
		CookieStoragePath:  a.CookieStoragePath,
		BindAt:             a.BindAt,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func taskToProto(t *biz.XshPublishTask) *xshv1.PublishTask {
	if t == nil {
		return nil
	}
	return &xshv1.PublishTask{
		Id:                  t.ID,
		DraftId:             t.DraftID,
		AccountId:          t.AccountID,
		ScheduledAt:         t.ScheduledAt,
		PublishedAt:         t.PublishedAt,
		Status:              t.Status,
		PlatformContentId:   t.PlatformContentID,
		ErrorMessage:        t.ErrorMessage,
		CreatedAt:           t.CreatedAt,
		UpdatedAt:           t.UpdatedAt,
	}
}
