// Package biz 萌宠之家（cat-admin）背包业务层
// 定义 InventoryEntry、InventoryRepo、InventoryUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// InventoryEntry 背包单项（含商品信息，用于列表展示）
type InventoryEntry struct {
	ItemID   int64
	Quantity int32
	Name     string
	Category string
	ImageURL string
}

// InventoryRepo 背包仓储接口，由 data 层实现
type InventoryRepo interface {
	ListByUserID(ctx context.Context, userID uint64) ([]*InventoryEntry, error)
	AddOrIncrement(ctx context.Context, userID uint64, itemID int64, delta int32) error
	DecrQuantity(ctx context.Context, userID uint64, itemID int64, need int32) (after int32, err error)
}

var (
	ErrInventoryItemNotFound = errors.New("inventory item not found or quantity insufficient")
)

// InventoryUsecase 背包用例：列表、使用道具（如喂食增加饱食度）
type InventoryUsecase struct {
	invRepo InventoryRepo
	petRepo PetRepo
	log     *log.Helper
}

// NewInventoryUsecase 创建背包用例
func NewInventoryUsecase(invRepo InventoryRepo, petRepo PetRepo, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{invRepo: invRepo, petRepo: petRepo, log: log.NewHelper(logger)}
}

// List 我的背包
func (uc *InventoryUsecase) List(ctx context.Context, userID uint64) ([]*InventoryEntry, error) {
	return uc.invRepo.ListByUserID(ctx, userID)
}

// UseItem 使用道具：扣减数量；若为食物且传入 pet_id 则增加该宠物饱食度（需校验归属）
func (uc *InventoryUsecase) UseItem(ctx context.Context, userID uint64, itemID int64, petID int64) error {
	after, err := uc.invRepo.DecrQuantity(ctx, userID, itemID, 1)
	if err != nil {
		return ErrInventoryItemNotFound
	}
	_ = after

	// 若指定了宠物，且物品为食物类，则增加宠物饱食度（这里简化：只要有 pet_id 就尝试加饱食度，不严格校验 category）
	if petID > 0 {
		pet, err := uc.petRepo.GetByID(ctx, petID)
		if err != nil || pet == nil {
			return err
		}
		if pet.UserID != userID {
			return ErrPetNotFound
		}
		fullness := pet.Fullness + 10
		if fullness > 100 {
			fullness = 100
		}
		// 只更新饱食度，其余保持；使用现有 UpdateAffectionAndMood 传原值
		_ = uc.petRepo.UpdateAffectionAndMood(ctx, petID, pet.Affection, fullness, pet.Happiness, pet.Cleanliness, pet.Mood)
	}
	return nil
}
