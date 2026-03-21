// Package biz 萌宠之家（cat-admin）商店业务层
// 定义 ShopItem、ShopRepo、ShopUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// ShopItem 商品领域实体
type ShopItem struct {
	ID          int64
	Name        string
	Category    string
	PriceCoins  int32
	ImageURL    string
	Description string
}

// ShopRepo 商店仓储接口，由 data 层实现
type ShopRepo interface {
	List(ctx context.Context, category string) ([]*ShopItem, error)
	GetByID(ctx context.Context, itemID int64) (*ShopItem, error)
}

// ShopUsecase 商店用例：商品列表、购买（扣金币+加背包）；首次购买触发成就
type ShopUsecase struct {
	repo        ShopRepo
	userRepo    UserRepo
	invRepo     InventoryRepo
	achUnlocker AchievementUnlocker
	log         *log.Helper
}

// NewShopUsecase 创建商店用例
func NewShopUsecase(repo ShopRepo, userRepo UserRepo, invRepo InventoryRepo, achUnlocker AchievementUnlocker, logger log.Logger) *ShopUsecase {
	return &ShopUsecase{repo: repo, userRepo: userRepo, invRepo: invRepo, achUnlocker: achUnlocker, log: log.NewHelper(logger)}
}

// ListItems 商品列表，category 为空则全部
func (uc *ShopUsecase) ListItems(ctx context.Context, category string) ([]*ShopItem, error) {
	return uc.repo.List(ctx, category)
}

var (
	ErrShopItemNotFound = errors.New("shop item not found")
	ErrInsufficientCoins = errors.New("insufficient coins")
)

// Buy 购买：校验商品存在、扣金币、增加背包数量；金币不足或商品不存在返回错误
func (uc *ShopUsecase) Buy(ctx context.Context, userID uint64, itemID int64, quantity int32) error {
	if quantity <= 0 {
		return ErrShopItemNotFound
	}
	item, err := uc.repo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrShopItemNotFound
	}
	total := item.PriceCoins * quantity
	if total <= 0 {
		return ErrShopItemNotFound
	}
	if err := uc.userRepo.DeductCoins(ctx, userID, total); err != nil {
		// 余额不足或用户不存在时 data 层返回错误，统一视为金币不足
		return ErrInsufficientCoins
	}
	if err := uc.invRepo.AddOrIncrement(ctx, userID, itemID, quantity); err != nil {
		return err
	}
	if uc.achUnlocker != nil {
		_ = uc.achUnlocker.UnlockByCode(ctx, userID, "first_purchase")
	}
	return nil
}
