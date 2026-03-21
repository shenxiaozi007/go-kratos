// Package service 萌宠之家（cat-admin）品种 HTTP 服务层
// 实现 api/breed/v1.BreedServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md
// 品种列表与详情为只读，可不强制 JWT（与商店列表一致），若全局中间件已鉴权则仍会校验

package service

import (
	"context"

	breedv1 "realworld/api/breed/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrBreedNotFoundHTTP 品种不存在时返回 404
var ErrBreedNotFoundHTTP = errors.New(404, "BREED_NOT_FOUND", "品种不存在")

// BreedService 实现 breed.v1.BreedServiceHTTPServer
type BreedService struct {
	uc  *biz.BreedUsecase
	log *log.Helper
}

// NewBreedService 创建品种服务
func NewBreedService(uc *biz.BreedUsecase, logger log.Logger) *BreedService {
	return &BreedService{uc: uc, log: log.NewHelper(logger)}
}

// ListBreeds GET /api/v1/breeds，query: species 可选
func (s *BreedService) ListBreeds(ctx context.Context, req *breedv1.ListBreedsRequest) (*breedv1.ListBreedsReply, error) {
	list, err := s.uc.ListBreeds(ctx, req.GetSpecies())
	if err != nil {
		s.log.Errorf("ListBreeds failed: %v", err)
		return nil, err
	}
	breeds := make([]*breedv1.Breed, 0, len(list))
	for _, b := range list {
		breeds = append(breeds, &breedv1.Breed{
			Id:       b.ID,
			Species:  b.Species,
			NameCn:   b.NameCn,
			NameEn:   b.NameEn,
			SizeTag:  b.SizeTag,
			LifeSpan: b.LifeSpan,
			ImageUrl: b.ImageURL,
		})
	}
	return &breedv1.ListBreedsReply{Breeds: breeds}, nil
}

// GetBreed GET /api/v1/breeds/{id}
func (s *BreedService) GetBreed(ctx context.Context, req *breedv1.GetBreedRequest) (*breedv1.GetBreedReply, error) {
	b, err := s.uc.GetBreed(ctx, req.GetId())
	if err != nil {
		if err == biz.ErrBreedNotFound {
			return nil, ErrBreedNotFoundHTTP
		}
		s.log.Errorf("GetBreed failed: %v", err)
		return nil, err
	}
	return &breedv1.GetBreedReply{
		Breed: &breedv1.Breed{
			Id:       b.ID,
			Species:  b.Species,
			NameCn:   b.NameCn,
			NameEn:   b.NameEn,
			SizeTag:  b.SizeTag,
			LifeSpan: b.LifeSpan,
			ImageUrl: b.ImageURL,
		},
	}, nil
}
