// Package service 萌宠之家（cat-admin）背包 HTTP 服务层
// 实现 api/inventory/v1.InventoryServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	inventoryv1 "realworld/api/inventory/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrInventoryUnauthorized 未登录无法访问背包
var ErrInventoryUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// InventoryService 实现 inventory.v1.InventoryServiceHTTPServer
type InventoryService struct {
	uc  *biz.InventoryUsecase
	log *log.Helper
}

// NewInventoryService 创建背包服务
func NewInventoryService(uc *biz.InventoryUsecase, logger log.Logger) *InventoryService {
	return &InventoryService{uc: uc, log: log.NewHelper(logger)}
}

func (s *InventoryService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrInventoryUnauthorized
	}
	return u.UserID, nil
}

// List GET /api/v1/inventory 我的背包
func (s *InventoryService) List(ctx context.Context, _ *inventoryv1.ListRequest) (*inventoryv1.ListReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.uc.List(ctx, userID)
	if err != nil {
		s.log.Errorf("Inventory List failed: %v", err)
		return nil, err
	}
	items := make([]*inventoryv1.InventoryItem, 0, len(list))
	for _, e := range list {
		items = append(items, &inventoryv1.InventoryItem{
			ItemId:   e.ItemID,
			Name:     e.Name,
			Category: e.Category,
			Quantity: e.Quantity,
			ImageUrl: e.ImageURL,
		})
	}
	return &inventoryv1.ListReply{Items: items}, nil
}

// UseItem POST /api/v1/inventory/use 使用道具，喂食等需传 pet_id
func (s *InventoryService) UseItem(ctx context.Context, req *inventoryv1.UseItemRequest) (*inventoryv1.UseItemReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	petID := req.GetPetId()
	err = s.uc.UseItem(ctx, userID, req.GetItemId(), petID)
	if err != nil {
		if err == biz.ErrInventoryItemNotFound {
			return &inventoryv1.UseItemReply{Success: false, Message: "物品不存在或数量不足"}, nil
		}
		if err == biz.ErrPetNotFound {
			return &inventoryv1.UseItemReply{Success: false, Message: "宠物不存在"}, nil
		}
		s.log.Errorf("UseItem failed: %v", err)
		return nil, err
	}
	return &inventoryv1.UseItemReply{Success: true, Message: "使用成功"}, nil
}
