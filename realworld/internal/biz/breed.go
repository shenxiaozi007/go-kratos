// Package biz 萌宠之家（cat-admin）品种业务层
// 定义 Breed、BreedRepo、BreedUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// Breed 品种领域实体
type Breed struct {
	ID       int64
	Species  string // CAT / DOG
	NameCn   string
	NameEn   string
	SizeTag  string
	LifeSpan string
	ImageURL string
}

// BreedRepo 品种仓储接口，由 data 层实现
type BreedRepo interface {
	List(ctx context.Context, species string) ([]*Breed, error)
	GetByID(ctx context.Context, id int64) (*Breed, error)
}

// BreedUsecase 品种用例：列表、详情（领养时前端可先选品种再调 CreatePet）
type BreedUsecase struct {
	repo BreedRepo
	log  *log.Helper
}

// NewBreedUsecase 创建品种用例
func NewBreedUsecase(repo BreedRepo, logger log.Logger) *BreedUsecase {
	return &BreedUsecase{repo: repo, log: log.NewHelper(logger)}
}

// ListBreeds 品种列表，species 为空则全部
func (uc *BreedUsecase) ListBreeds(ctx context.Context, species string) ([]*Breed, error) {
	return uc.repo.List(ctx, species)
}

var ErrBreedNotFound = errors.New("breed not found")

// GetBreed 品种详情，不存在返回 ErrBreedNotFound
func (uc *BreedUsecase) GetBreed(ctx context.Context, id int64) (*Breed, error) {
	b, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBreedNotFound
	}
	return b, nil
}
