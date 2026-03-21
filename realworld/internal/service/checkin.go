// Package service 萌宠之家（cat-admin）签到 HTTP 服务层
// 实现 api/checkin/v1.CheckinServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	checkinv1 "realworld/api/checkin/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrCheckinUnauthorized 未登录无法签到
var ErrCheckinUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// CheckinService 实现 checkin.v1.CheckinServiceHTTPServer
type CheckinService struct {
	uc  *biz.CheckinUsecase
	log *log.Helper
}

// NewCheckinService 创建签到服务
func NewCheckinService(uc *biz.CheckinUsecase, logger log.Logger) *CheckinService {
	return &CheckinService{uc: uc, log: log.NewHelper(logger)}
}

func (s *CheckinService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrCheckinUnauthorized
	}
	return u.UserID, nil
}

// DoCheckin POST /api/v1/checkin 每日签到
func (s *CheckinService) DoCheckin(ctx context.Context, _ *checkinv1.DoCheckinRequest) (*checkinv1.DoCheckinReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	already, continuous, coins, heart, err := s.uc.DoCheckin(ctx, userID)
	if err != nil {
		s.log.Errorf("DoCheckin failed: %v", err)
		return nil, err
	}
	return &checkinv1.DoCheckinReply{
		AlreadyCheckedToday: already,
		ContinuousDays:      continuous,
		CoinsReward:         coins,
		HeartReward:         heart,
	}, nil
}

// GetCheckinStatus GET /api/v1/checkin/status
func (s *CheckinService) GetCheckinStatus(ctx context.Context, _ *checkinv1.GetCheckinStatusRequest) (*checkinv1.GetCheckinStatusReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	checkedToday, continuous, last7, err := s.uc.GetCheckinStatus(ctx, userID)
	if err != nil {
		s.log.Errorf("GetCheckinStatus failed: %v", err)
		return nil, err
	}
	return &checkinv1.GetCheckinStatusReply{
		CheckedToday:   checkedToday,
		ContinuousDays: continuous,
		Last_7Days:     last7,
	}, nil
}
