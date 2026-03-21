// Package biz 萌宠之家（cat-admin）成就业务层
// 定义 Achievement、AchievementRepo、AchievementUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// Achievement 成就领域实体
type Achievement struct {
	ID             int64
	Code           string
	Name           string
	Description    string
	IconURL        string
	ConditionType  string
	ConditionValue int32
}

// AchievementRepo 成就仓储接口，由 data 层实现
type AchievementRepo interface {
	ListAll(ctx context.Context) ([]*Achievement, error)
	GetByID(ctx context.Context, id int64) (*Achievement, error)
	GetByCode(ctx context.Context, code string) (*Achievement, error)
	GetUnlockedByUser(ctx context.Context, userID uint64) (map[int64]int64, error)
	Unlock(ctx context.Context, userID uint64, achievementID int64) error
}

// AchievementUnlocker 成就解锁器：其他模块在满足条件时调用，按成就 code 解锁（幂等）
type AchievementUnlocker interface {
	UnlockByCode(ctx context.Context, userID uint64, code string) error
}

// achievementUnlocker 实现 AchievementUnlocker，供签到/领养/购买等触发
type achievementUnlocker struct {
	repo AchievementRepo
}

// NewAchievementUnlocker 创建成就解锁器
func NewAchievementUnlocker(repo AchievementRepo) AchievementUnlocker {
	return &achievementUnlocker{repo: repo}
}

// UnlockByCode 按成就编码解锁，成就不存在或已解锁则静默忽略
func (u *achievementUnlocker) UnlockByCode(ctx context.Context, userID uint64, code string) error {
	if code == "" {
		return nil
	}
	a, err := u.repo.GetByCode(ctx, code)
	if err != nil || a == nil {
		return nil
	}
	return u.repo.Unlock(ctx, userID, a.ID)
}

// AchievementUsecase 成就用例：列表（含解锁状态）、详情（含获得时间）
type AchievementUsecase struct {
	repo AchievementRepo
	log  *log.Helper
}

// NewAchievementUsecase 创建成就用例
func NewAchievementUsecase(repo AchievementRepo, logger log.Logger) *AchievementUsecase {
	return &AchievementUsecase{repo: repo, log: log.NewHelper(logger)}
}

// ListAchievements 返回全部成就及当前用户各成就是否已解锁、解锁时间
func (uc *AchievementUsecase) ListAchievements(ctx context.Context, userID uint64) ([]*Achievement, map[int64]int64, error) {
	list, err := uc.repo.ListAll(ctx)
	if err != nil {
		return nil, nil, err
	}
	unlocked, err := uc.repo.GetUnlockedByUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	return list, unlocked, nil
}

var ErrAchievementNotFound = errors.New("achievement not found")

// GetAchievement 成就详情，含当前用户是否已解锁及解锁时间
func (uc *AchievementUsecase) GetAchievement(ctx context.Context, userID uint64, id int64) (*Achievement, bool, int64, error) {
	a, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, false, 0, err
	}
	if a == nil {
		return nil, false, 0, ErrAchievementNotFound
	}
	unlocked, err := uc.repo.GetUnlockedByUser(ctx, userID)
	if err != nil {
		return nil, false, 0, err
	}
	at, ok := unlocked[id]
	if !ok {
		return a, false, 0, nil
	}
	return a, true, at, nil
}
