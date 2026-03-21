// Package biz 萌宠之家（cat-admin）签到业务层
// 定义 CheckIn 领域、CheckinRepo 接口及 CheckinUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// CheckIn 签到记录领域实体
type CheckIn struct {
	ID             int64
	UserID         uint64
	Date           time.Time
	ContinuousDays int32
}

// CheckinRepo 签到仓储接口，由 data 层实现
type CheckinRepo interface {
	GetByUserAndDate(ctx context.Context, userID uint64, date time.Time) (*CheckIn, error)
	Save(ctx context.Context, c *CheckIn) error
	ListByUserLastNDays(ctx context.Context, userID uint64, n int) ([]*CheckIn, error)
}

// CheckinUsecase 签到用例：执行签到、获取状态与 7 日日历；签到奖励写入用户金币/爱心（依赖 UserRepo）；首次签到触发成就
type CheckinUsecase struct {
	repo        CheckinRepo
	userRepo    UserRepo
	achUnlocker AchievementUnlocker
	log         *log.Helper
}

// NewCheckinUsecase 创建签到用例
func NewCheckinUsecase(repo CheckinRepo, userRepo UserRepo, achUnlocker AchievementUnlocker, logger log.Logger) *CheckinUsecase {
	return &CheckinUsecase{repo: repo, userRepo: userRepo, achUnlocker: achUnlocker, log: log.NewHelper(logger)}
}

// DoCheckin 执行每日签到。若今日已签返回 already=true 与当前连续天数；否则写入记录并计算连续天数与奖励
func (uc *CheckinUsecase) DoCheckin(ctx context.Context, userID uint64) (already bool, continuousDays int32, coinsReward, heartReward int32, err error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	exist, err := uc.repo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		return false, 0, 0, 0, err
	}
	if exist != nil {
		return true, exist.ContinuousDays, 0, 0, nil
	}

	// 计算连续天数：查昨天是否签到
	yesterday := today.AddDate(0, 0, -1)
	prev, _ := uc.repo.GetByUserAndDate(ctx, userID, yesterday)
	continuous := int32(1)
	if prev != nil {
		continuous = prev.ContinuousDays + 1
	}

	c := &CheckIn{UserID: userID, Date: today, ContinuousDays: continuous}
	if err := uc.repo.Save(ctx, c); err != nil {
		return false, 0, 0, 0, err
	}

	// 简单奖励规则：连续 3 天多送，其余固定
	coinsReward = 10
	heartReward = 5
	if continuous >= 3 {
		coinsReward = 15
		heartReward = 8
	}
	// 将奖励写入用户表
	if coinsReward > 0 || heartReward > 0 {
		_ = uc.userRepo.AddCoinsAndHeart(ctx, userID, coinsReward, heartReward)
	}
	// 首次签到解锁成就（连续天数为 1 表示今日是第一次签到）
	if continuous == 1 && uc.achUnlocker != nil {
		_ = uc.achUnlocker.UnlockByCode(ctx, userID, "first_checkin")
	}
	return false, continuous, coinsReward, heartReward, nil
}

// GetCheckinStatus 获取今日是否已签、连续天数、最近 7 天签到情况（用于前端日历）
func (uc *CheckinUsecase) GetCheckinStatus(ctx context.Context, userID uint64) (checkedToday bool, continuousDays int32, last7 []bool, err error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	exist, err := uc.repo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		return false, 0, nil, err
	}
	if exist != nil {
		checkedToday = true
		continuousDays = exist.ContinuousDays
	}

	list, err := uc.repo.ListByUserLastNDays(ctx, userID, 7)
	if err != nil {
		return checkedToday, continuousDays, nil, err
	}
	// 构造最近 7 天（从 6 天前到今天）的签到布尔数组
	last7 = make([]bool, 7)
	for i := 0; i < 7; i++ {
		d := today.AddDate(0, 0, -6+i)
		for _, c := range list {
			cd := time.Date(c.Date.Year(), c.Date.Month(), c.Date.Day(), 0, 0, 0, 0, c.Date.Location())
			if cd.Equal(d) {
				last7[i] = true
				break
			}
		}
	}
	return checkedToday, continuousDays, last7, nil
}
