// Package biz 选品采集仓（Input）业务层：商品、爆款榜单、链接解析的领域实体与用例。
// 不依赖具体存储实现，通过 XshProductRepo / XshHotRankRepo 接口与 data 层交互。

package biz

import "context"

// XshProduct 商品领域实体，对应表 xsh_products。
// 用于在 biz 与 data 之间传递，不暴露 GORM 等实现细节。
type XshProduct struct {
	ID             uint64
	SourceURL      string
	Title          string
	OriginalPrice  float64
	CouponAmount   float64
	CommissionRate float64
	ImageURL       string
	IsHot          bool
	RankSource     string
	CreatedAt      int64
	UpdatedAt      int64
}

// XshHotRank 爆款榜单快照领域实体，对应表 xsh_hot_ranks。
type XshHotRank struct {
	ID           uint64
	RankType     string
	SyncedAt     int64
	SnapshotJSON string
	CreatedAt    int64
}

// XshProductRepo 选品仓储接口：商品 CRUD 与链接解析。
// 由 data 层实现，biz 层仅依赖此接口。
type XshProductRepo interface {
	GetByID(ctx context.Context, id uint64) (*XshProduct, error)
	List(ctx context.Context, page, pageSize int32, isHot *bool, rankSource string) ([]*XshProduct, int64, error)
	Create(ctx context.Context, p *XshProduct) (*XshProduct, error)
	Update(ctx context.Context, p *XshProduct) error
	Delete(ctx context.Context, id uint64) error
	// ParseProductLink 解析链接；首期可不接联盟，返回 (nil, errorMessage, nil) 表示未解析到商品。
	ParseProductLink(ctx context.Context, url string) (*XshProduct, string, error)
}

// XshHotRankRepo 爆款榜单仓储接口。
type XshHotRankRepo interface {
	List(ctx context.Context, rankType string, page, pageSize int32) ([]*XshHotRank, int64, error)
	Sync(ctx context.Context, rankType string) (bool, string, error)
}

// XshProductUsecase 选品用例：封装商品列表、创建、更新、删除、解析链接、榜单同步等业务逻辑。
type XshProductUsecase struct {
	productRepo  XshProductRepo
	hotRankRepo  XshHotRankRepo
}

// NewXshProductUsecase 创建选品用例实例。
func NewXshProductUsecase(productRepo XshProductRepo, hotRankRepo XshHotRankRepo) *XshProductUsecase {
	return &XshProductUsecase{
		productRepo: productRepo,
		hotRankRepo: hotRankRepo,
	}
}

// ParseProductLink 解析链接；首期可不接联盟，返回 (nil, errorMessage, nil) 表示未解析到商品，由调用方提示用户手动录入。
func (uc *XshProductUsecase) ParseProductLink(ctx context.Context, url string) (*XshProduct, string, error) {
	return uc.productRepo.ParseProductLink(ctx, url)
}

// ListProducts 分页列表，支持 is_hot、rank_source 筛选。
func (uc *XshProductUsecase) ListProducts(ctx context.Context, page, pageSize int32, isHot *bool, rankSource string) ([]*XshProduct, int64, error) {
	return uc.productRepo.List(ctx, page, pageSize, isHot, rankSource)
}

// GetProduct 按 ID 查询商品。
func (uc *XshProductUsecase) GetProduct(ctx context.Context, id uint64) (*XshProduct, error) {
	return uc.productRepo.GetByID(ctx, id)
}

// CreateProduct 创建商品（手动录入或解析后入库）。
func (uc *XshProductUsecase) CreateProduct(ctx context.Context, p *XshProduct) (*XshProduct, error) {
	return uc.productRepo.Create(ctx, p)
}

// UpdateProduct 更新商品。
func (uc *XshProductUsecase) UpdateProduct(ctx context.Context, p *XshProduct) error {
	return uc.productRepo.Update(ctx, p)
}

// DeleteProduct 删除商品。
func (uc *XshProductUsecase) DeleteProduct(ctx context.Context, id uint64) error {
	return uc.productRepo.Delete(ctx, id)
}

// ListHotRanks 分页查询爆款榜单快照。
func (uc *XshProductUsecase) ListHotRanks(ctx context.Context, rankType string, page, pageSize int32) ([]*XshHotRank, int64, error) {
	return uc.hotRankRepo.List(ctx, rankType, page, pageSize)
}

// SyncHotRank 同步爆款榜单；首期可仅记录一次快照。
func (uc *XshProductUsecase) SyncHotRank(ctx context.Context, rankType string) (bool, string, error) {
	return uc.hotRankRepo.Sync(ctx, rankType)
}
