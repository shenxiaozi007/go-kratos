// Package service 选品采集仓（Input）HTTP 服务：实现 api/xsh/v1 ProductService 的 HTTP 接口。
// 将 proto 请求/响应与 biz 层转换，委托 XshProductUsecase 执行业务逻辑。

package service

import (
	"context"

	xshv1 "realworld/api/xsh/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// XshProductService 实现 xsh.v1.ProductService 的 HTTP 服务。
type XshProductService struct {
	uc  *biz.XshProductUsecase
	log *log.Helper
}

// NewXshProductService 创建选品 HTTP 服务实例。
func NewXshProductService(uc *biz.XshProductUsecase, logger log.Logger) *XshProductService {
	return &XshProductService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ParseProductLink 解析链接；首期未接联盟时返回空商品与提示信息。
func (s *XshProductService) ParseProductLink(ctx context.Context, req *xshv1.ParseProductLinkRequest) (*xshv1.ParseProductLinkReply, error) {
	product, errMsg, err := s.uc.ParseProductLink(ctx, req.GetUrl())
	if err != nil {
		s.log.Errorf("ParseProductLink failed: %v", err)
		return nil, err
	}
	reply := &xshv1.ParseProductLinkReply{ErrorMessage: errMsg}
	if product != nil {
		reply.Product = bizProductToProto(product)
	}
	return reply, nil
}

// ListProducts 分页列表商品，支持 is_hot、rank_source 筛选。
func (s *XshProductService) ListProducts(ctx context.Context, req *xshv1.ListProductsRequest) (*xshv1.ListProductsReply, error) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	// 仅当客户端显式传 is_hot=true 时筛选神价商品
	var isHot *bool
	if req.GetIsHot() {
		v := true
		isHot = &v
	}
	items, total, err := s.uc.ListProducts(ctx, page, pageSize, isHot, req.GetRankSource())
	if err != nil {
		s.log.Errorf("ListProducts failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.Product, 0, len(items))
	for _, p := range items {
		list = append(list, bizProductToProto(p))
	}
	return &xshv1.ListProductsReply{Items: list, Total: total}, nil
}

// GetProduct 按 ID 查询商品。
func (s *XshProductService) GetProduct(ctx context.Context, req *xshv1.GetProductRequest) (*xshv1.GetProductReply, error) {
	p, err := s.uc.GetProduct(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return &xshv1.GetProductReply{}, nil
	}
	return &xshv1.GetProductReply{Product: bizProductToProto(p)}, nil
}

// CreateProduct 创建商品（手动录入）。
func (s *XshProductService) CreateProduct(ctx context.Context, req *xshv1.CreateProductRequest) (*xshv1.CreateProductReply, error) {
	p := &biz.XshProduct{
		SourceURL:      req.GetSourceUrl(),
		Title:          req.GetTitle(),
		OriginalPrice:  req.GetOriginalPrice(),
		CouponAmount:   req.GetCouponAmount(),
		CommissionRate: req.GetCommissionRate(),
		ImageURL:       req.GetImageUrl(),
		IsHot:          req.GetIsHot(),
		RankSource:     req.GetRankSource(),
	}
	created, err := s.uc.CreateProduct(ctx, p)
	if err != nil {
		s.log.Errorf("CreateProduct failed: %v", err)
		return nil, err
	}
	return &xshv1.CreateProductReply{Product: bizProductToProto(created)}, nil
}

// UpdateProduct 更新商品；请求中非空字段覆盖已有值。
func (s *XshProductService) UpdateProduct(ctx context.Context, req *xshv1.UpdateProductRequest) (*xshv1.UpdateProductReply, error) {
	existing, err := s.uc.GetProduct(ctx, req.GetId())
	if err != nil || existing == nil {
		return nil, err
	}
	// 仅用请求中非零/非空字段覆盖
	if req.GetSourceUrl() != "" {
		existing.SourceURL = req.GetSourceUrl()
	}
	if req.GetTitle() != "" {
		existing.Title = req.GetTitle()
	}
	if req.GetOriginalPrice() != 0 {
		existing.OriginalPrice = req.GetOriginalPrice()
	}
	if req.GetCouponAmount() != 0 {
		existing.CouponAmount = req.GetCouponAmount()
	}
	if req.GetCommissionRate() != 0 {
		existing.CommissionRate = req.GetCommissionRate()
	}
	if req.GetImageUrl() != "" {
		existing.ImageURL = req.GetImageUrl()
	}
	existing.IsHot = req.GetIsHot()
	if req.GetRankSource() != "" {
		existing.RankSource = req.GetRankSource()
	}
	if err := s.uc.UpdateProduct(ctx, existing); err != nil {
		s.log.Errorf("UpdateProduct failed: %v", err)
		return nil, err
	}
	return &xshv1.UpdateProductReply{Product: bizProductToProto(existing)}, nil
}

// DeleteProduct 删除商品。
func (s *XshProductService) DeleteProduct(ctx context.Context, req *xshv1.DeleteProductRequest) (*xshv1.DeleteProductReply, error) {
	if err := s.uc.DeleteProduct(ctx, req.GetId()); err != nil {
		s.log.Errorf("DeleteProduct failed: %v", err)
		return nil, err
	}
	return &xshv1.DeleteProductReply{Success: true}, nil
}

// ListHotRanks 分页查询爆款榜单快照。
func (s *XshProductService) ListHotRanks(ctx context.Context, req *xshv1.ListHotRanksRequest) (*xshv1.ListHotRanksReply, error) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	items, total, err := s.uc.ListHotRanks(ctx, req.GetRankType(), page, pageSize)
	if err != nil {
		s.log.Errorf("ListHotRanks failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.HotRank, 0, len(items))
	for _, r := range items {
		list = append(list, &xshv1.HotRank{
			Id:           r.ID,
			RankType:     r.RankType,
			SyncedAt:     r.SyncedAt,
			SnapshotJson: r.SnapshotJSON,
			CreatedAt:    r.CreatedAt,
		})
	}
	return &xshv1.ListHotRanksReply{Items: list, Total: total}, nil
}

// SyncHotRank 同步爆款榜单；首期仅记录快照。
func (s *XshProductService) SyncHotRank(ctx context.Context, req *xshv1.SyncHotRankRequest) (*xshv1.SyncHotRankReply, error) {
	ok, msg, err := s.uc.SyncHotRank(ctx, req.GetRankType())
	if err != nil {
		s.log.Errorf("SyncHotRank failed: %v", err)
		return nil, err
	}
	return &xshv1.SyncHotRankReply{Success: ok, Message: msg}, nil
}

func bizProductToProto(p *biz.XshProduct) *xshv1.Product {
	if p == nil {
		return nil
	}
	return &xshv1.Product{
		Id:             p.ID,
		SourceUrl:      p.SourceURL,
		Title:          p.Title,
		OriginalPrice:  p.OriginalPrice,
		CouponAmount:   p.CouponAmount,
		CommissionRate: p.CommissionRate,
		ImageUrl:       p.ImageURL,
		IsHot:          p.IsHot,
		RankSource:     p.RankSource,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
