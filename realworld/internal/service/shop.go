// Package service 萌宠之家（cat-admin）商店 HTTP 服务层
// 实现 api/shop/v1.ShopServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	shopv1 "realworld/api/shop/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrShopUnauthorized 未登录无法访问商店
var ErrShopUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// ShopService 实现 shop.v1.ShopServiceHTTPServer
type ShopService struct {
	uc  *biz.ShopUsecase
	log *log.Helper
}

// NewShopService 创建商店服务
func NewShopService(uc *biz.ShopUsecase, logger log.Logger) *ShopService {
	return &ShopService{uc: uc, log: log.NewHelper(logger)}
}

func (s *ShopService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrShopUnauthorized
	}
	return u.UserID, nil
}

// ListItems GET /api/v1/shop/items，category 由 query 绑定到 ListItemsRequest.Category
func (s *ShopService) ListItems(ctx context.Context, req *shopv1.ListItemsRequest) (*shopv1.ListItemsReply, error) {
	items, err := s.uc.ListItems(ctx, req.GetCategory())
	if err != nil {
		s.log.Errorf("ListItems failed: %v", err)
		return nil, err
	}
	out := make([]*shopv1.ShopItem, 0, len(items))
	for _, it := range items {
		out = append(out, &shopv1.ShopItem{
			Id:          it.ID,
			Name:        it.Name,
			Category:    it.Category,
			PriceCoins:  it.PriceCoins,
			ImageUrl:    it.ImageURL,
			Description: it.Description,
		})
	}
	return &shopv1.ListItemsReply{Items: out}, nil
}

// Buy POST /api/v1/shop/buy 需登录，扣金币并写入背包
func (s *ShopService) Buy(ctx context.Context, req *shopv1.BuyRequest) (*shopv1.BuyReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	quantity := req.GetQuantity()
	if quantity <= 0 {
		quantity = 1
	}
	err = s.uc.Buy(ctx, userID, req.GetItemId(), quantity)
	if err != nil {
		if err == biz.ErrInsufficientCoins {
			return &shopv1.BuyReply{Success: false, Message: "金币不足"}, nil
		}
		if err == biz.ErrShopItemNotFound {
			return &shopv1.BuyReply{Success: false, Message: "商品不存在"}, nil
		}
		s.log.Errorf("Buy failed: %v", err)
		return nil, err
	}
	return &shopv1.BuyReply{Success: true, Message: "购买成功"}, nil
}
