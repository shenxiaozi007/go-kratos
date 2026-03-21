// Package biz 萌宠之家（cat-admin）用户资料业务层
// 提供 GetProfile、UpdateProfile 用例，复用 UserRepo（与鉴权共用），见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// UserUsecase 用户资料用例：获取/更新当前用户资料（昵称、头像、个性签名、等级、金币、爱心）
type UserUsecase struct {
	userRepo UserRepo
	log      *log.Helper
}

// NewUserUsecase 创建用户资料用例
func NewUserUsecase(userRepo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, log: log.NewHelper(logger)}
}

// GetProfile 获取用户资料，供 GET /api/v1/user/profile
func (uc *UserUsecase) GetProfile(ctx context.Context, userID uint64) (*User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

// UpdateProfile 更新昵称、头像、个性签名；空字符串表示不更新该字段
func (uc *UserUsecase) UpdateProfile(ctx context.Context, userID uint64, nickname, avatarURL, signature string) (*User, error) {
	if err := uc.userRepo.UpdateProfile(ctx, userID, nickname, avatarURL, signature); err != nil {
		uc.log.Errorf("UpdateProfile failed: %v", err)
		return nil, err
	}
	return uc.userRepo.GetByID(ctx, userID)
}
