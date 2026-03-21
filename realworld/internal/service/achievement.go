// Package service 萌宠之家（cat-admin）成就 HTTP 服务层
// 实现 api/achievement/v1.AchievementServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	achievementv1 "realworld/api/achievement/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrAchievementUnauthorized 未登录无法查看成就
var ErrAchievementUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// ErrAchievementNotFoundHTTP 成就不存在时返回 404
var ErrAchievementNotFoundHTTP = errors.New(404, "ACHIEVEMENT_NOT_FOUND", "成就不存在")

// AchievementService 实现 achievement.v1.AchievementServiceHTTPServer
type AchievementService struct {
	uc  *biz.AchievementUsecase
	log *log.Helper
}

// NewAchievementService 创建成就服务
func NewAchievementService(uc *biz.AchievementUsecase, logger log.Logger) *AchievementService {
	return &AchievementService{uc: uc, log: log.NewHelper(logger)}
}

func (s *AchievementService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrAchievementUnauthorized
	}
	return u.UserID, nil
}

func bizAchievementToProto(a *biz.Achievement) *achievementv1.Achievement {
	if a == nil {
		return nil
	}
	return &achievementv1.Achievement{
		Id:             a.ID,
		Code:           a.Code,
		Name:           a.Name,
		Description:    a.Description,
		IconUrl:        a.IconURL,
		ConditionType:  a.ConditionType,
		ConditionValue: a.ConditionValue,
	}
}

// ListAchievements GET /api/v1/achievements 成就列表（含是否已解锁）
func (s *AchievementService) ListAchievements(ctx context.Context, _ *achievementv1.ListAchievementsRequest) (*achievementv1.ListAchievementsReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	list, unlocked, err := s.uc.ListAchievements(ctx, userID)
	if err != nil {
		s.log.Errorf("ListAchievements failed: %v", err)
		return nil, err
	}
	items := make([]*achievementv1.AchievementItem, 0, len(list))
	for _, a := range list {
		at := unlocked[a.ID]
		items = append(items, &achievementv1.AchievementItem{
			Achievement: bizAchievementToProto(a),
			Unlocked:    at > 0,
			UnlockedAt:  at,
		})
	}
	return &achievementv1.ListAchievementsReply{Items: items}, nil
}

// GetAchievement GET /api/v1/achievements/{id} 成就详情（描述、获得时间）
func (s *AchievementService) GetAchievement(ctx context.Context, req *achievementv1.GetAchievementRequest) (*achievementv1.GetAchievementReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	a, unlocked, at, err := s.uc.GetAchievement(ctx, userID, req.GetId())
	if err != nil {
		if err == biz.ErrAchievementNotFound {
			return nil, ErrAchievementNotFoundHTTP
		}
		s.log.Errorf("GetAchievement failed: %v", err)
		return nil, err
	}
	return &achievementv1.GetAchievementReply{
		Achievement: bizAchievementToProto(a),
		Unlocked:    unlocked,
		UnlockedAt:  at,
	}, nil
}
