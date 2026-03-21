// Package biz 萌宠之家（cat-admin）宠物装扮业务层
// 定义 AppearanceRepo、AppearanceUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// AppearanceRepo 装扮仓储接口，由 data 层实现
// 按宠物 ID 查询/保存头、身、颈三个槽位的 item_id，0 表示未装备
type AppearanceRepo interface {
	GetByPetID(ctx context.Context, petID int64) (headItemID, bodyItemID, neckItemID int64, err error)
	Save(ctx context.Context, petID int64, headItemID, bodyItemID, neckItemID int64) error
}

// AppearanceUsecase 装扮用例：获取/保存宠物装扮，需校验宠物归属当前用户
type AppearanceUsecase struct {
	repo    AppearanceRepo
	petRepo PetRepo
	log     *log.Helper
}

// NewAppearanceUsecase 创建装扮用例
func NewAppearanceUsecase(repo AppearanceRepo, petRepo PetRepo, logger log.Logger) *AppearanceUsecase {
	return &AppearanceUsecase{repo: repo, petRepo: petRepo, log: log.NewHelper(logger)}
}

// GetAppearance 获取宠物当前装扮，需校验 pet 归属 userID
func (uc *AppearanceUsecase) GetAppearance(ctx context.Context, userID uint64, petID int64) (head, body, neck int64, err error) {
	pet, err := uc.petRepo.GetByID(ctx, petID)
	if err != nil || pet == nil {
		return 0, 0, 0, ErrPetNotFound
	}
	if pet.UserID != userID {
		return 0, 0, 0, ErrPetNotFound
	}
	return uc.repo.GetByPetID(ctx, petID)
}

// UpdateAppearance 保存宠物装扮，需校验宠物归属；0 表示该槽位不装备
func (uc *AppearanceUsecase) UpdateAppearance(ctx context.Context, userID uint64, petID int64, head, body, neck int64) error {
	pet, err := uc.petRepo.GetByID(ctx, petID)
	if err != nil || pet == nil {
		return ErrPetNotFound
	}
	if pet.UserID != userID {
		return ErrPetNotFound
	}
	return uc.repo.Save(ctx, petID, head, body, neck)
}
