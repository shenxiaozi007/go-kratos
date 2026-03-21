// Package biz 萌宠之家（cat-admin）宠物业务层
// 定义领域模型 Pet、仓储接口 PetRepo 及用例 PetUsecase，与 api/pet/v1、data 层解耦，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// ErrPetNotFound 当前用户没有宠物或指定 ID 的宠物不存在
var ErrPetNotFound = errors.New("pet not found")

// Pet 萌宠之家领域层宠物实体，与 data.PetPO、api/pet/v1.Pet 解耦
// 用于 biz 与 data 之间传递，由 service 层转换为 proto 返回
type Pet struct {
	ID            int64
	UserID        uint64
	Name          string
	Species       string   // CAT / DOG
	BreedID       int64
	AvatarURL     string
	BackgroundURL string
	Mood          string   // 心情文案
	Affection     int32    // 亲密度 0-1000
	Fullness      int32    // 饱食度 0-100
	Happiness     int32    // 心情值 0-100
	Cleanliness   int32    // 清洁度 0-100
	AgeYears      int32
	HealthStatus  string
}

// PetRepo 宠物仓储接口，由 data 层实现，供 PetUsecase 调用
type PetRepo interface {
	GetByUserID(ctx context.Context, userID uint64) (*Pet, error)
	ListByUserID(ctx context.Context, userID uint64) ([]*Pet, error)
	GetByID(ctx context.Context, id int64) (*Pet, error)
	Save(ctx context.Context, p *Pet) error
	Update(ctx context.Context, p *Pet) error
	UpdateAffectionAndMood(ctx context.Context, petID int64, affection, fullness, happiness, cleanliness int32, mood string) error
}

// PetUsecase 萌宠之家宠物用例：获取主宠、列表、创建、更新、互动；首次领养触发成就
type PetUsecase struct {
	repo        PetRepo
	achUnlocker AchievementUnlocker
	log         *log.Helper
}

// NewPetUsecase 创建宠物用例，repo、achUnlocker 由 wire 注入
func NewPetUsecase(repo PetRepo, achUnlocker AchievementUnlocker, logger log.Logger) *PetUsecase {
	return &PetUsecase{repo: repo, achUnlocker: achUnlocker, log: log.NewHelper(logger)}
}

// GetMyPet 获取当前用户的主宠物（列表首只），无则返回 nil, ErrPetNotFound
func (uc *PetUsecase) GetMyPet(ctx context.Context, userID uint64) (*Pet, error) {
	p, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrPetNotFound
	}
	return p, nil
}

// ListPets 获取当前用户的全部宠物列表
func (uc *PetUsecase) ListPets(ctx context.Context, userID uint64) ([]*Pet, error) {
	return uc.repo.ListByUserID(ctx, userID)
}

// CreatePet 创建/领养一只新宠物，归属当前用户，默认成长数据见下方
func (uc *PetUsecase) CreatePet(ctx context.Context, userID uint64, name, species, avatarURL, backgroundURL string, breedID int64) (*Pet, error) {
	p := &Pet{
		UserID:        userID,
		Name:          name,
		Species:       species,
		BreedID:       breedID,
		AvatarURL:     avatarURL,
		BackgroundURL: backgroundURL,
		Mood:          "很开心",
		Affection:     0,
		Fullness:      80,
		Happiness:     80,
		Cleanliness:   80,
	}
	if err := uc.repo.Save(ctx, p); err != nil {
		uc.log.Errorf("CreatePet Save failed: %v", err)
		return nil, err
	}
	if uc.achUnlocker != nil {
		_ = uc.achUnlocker.UnlockByCode(ctx, userID, "first_pet")
	}
	return p, nil
}

// UpdatePet 更新宠物基础信息（名称、种类、头像、背景），需校验归属当前用户
func (uc *PetUsecase) UpdatePet(ctx context.Context, userID uint64, id int64, name, species, avatarURL, backgroundURL string) (*Pet, error) {
	p, err := uc.repo.GetByID(ctx, id)
	if err != nil || p == nil {
		return nil, ErrPetNotFound
	}
	if p.UserID != userID {
		return nil, ErrPetNotFound
	}
	p.Name = name
	p.Species = species
	p.AvatarURL = avatarURL
	p.BackgroundURL = backgroundURL
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Interact 根据动作类型增加亲密度并更新心情等，返回更新后的宠物；无宠物或非归属用户返回 ErrPetNotFound
// action 与 proto InteractionAction 对应：STROKE/TEASE/TREAT/BATH
func (uc *PetUsecase) Interact(ctx context.Context, userID uint64, action string) (*Pet, error) {
	p, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil || p == nil {
		return nil, ErrPetNotFound
	}

	// 根据动作增加亲密度（上限 1000），并微调饱食度/心情值/清洁度
	affection := p.Affection + interactionAffectionDelta(action)
	if affection > 1000 {
		affection = 1000
	}
	fullness := p.Fullness
	happiness := p.Happiness
	cleanliness := p.Cleanliness
	switch action {
	case "STROKE", "INTERACTION_ACTION_STROKE":
		happiness = clampInt32(p.Happiness+2, 0, 100)
	case "TEASE", "INTERACTION_ACTION_TEASE":
		happiness = clampInt32(p.Happiness+3, 0, 100)
	case "TREAT", "INTERACTION_ACTION_TREAT":
		fullness = clampInt32(p.Fullness+5, 0, 100)
	case "BATH", "INTERACTION_ACTION_BATH":
		cleanliness = clampInt32(p.Cleanliness+8, 0, 100)
	}
	mood := moodFromAffection(affection)

	if err := uc.repo.UpdateAffectionAndMood(ctx, p.ID, affection, fullness, happiness, cleanliness, mood); err != nil {
		return nil, err
	}
	p.Affection = affection
	p.Fullness = fullness
	p.Happiness = happiness
	p.Cleanliness = cleanliness
	p.Mood = mood
	return p, nil
}

// interactionAffectionDelta 各动作对应的亲密度增量
func interactionAffectionDelta(action string) int32 {
	switch action {
	case "STROKE", "INTERACTION_ACTION_STROKE":
		return 5
	case "TEASE", "INTERACTION_ACTION_TEASE":
		return 4
	case "TREAT", "INTERACTION_ACTION_TREAT":
		return 6
	case "BATH", "INTERACTION_ACTION_BATH":
		return 5
	default:
		return 2
	}
}

// moodFromAffection 根据亲密度返回心情文案，供前端展示
func moodFromAffection(affection int32) string {
	switch {
	case affection >= 800:
		return "非常开心"
	case affection >= 500:
		return "很开心"
	case affection >= 200:
		return "开心"
	default:
		return "慢慢熟悉中"
	}
}

func clampInt32(v, lo, hi int32) int32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
