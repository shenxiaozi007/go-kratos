// Package service 萌宠之家（cat-admin）宠物 HTTP 服务层
// 实现 api/pet/v1.PetServiceHTTPServer，从 context 取当前用户 ID，调用 biz.PetUsecase 并转换为 proto 响应，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	petv1 "realworld/api/pet/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// 宠物模块 HTTP 错误码，与 Kratos 约定一致
var (
	ErrUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")
	ErrPetNotFound  = errors.New(404, "PET_NOT_FOUND", "宠物不存在或无权操作")
	ErrBadRequest   = errors.New(400, "BAD_REQUEST", "参数错误")
)

// PetService 实现 pet.v1.PetServiceHTTPServer，所有接口需 JWT 鉴权
// 注：protoc-gen-go-http 未生成 Unimplemented 嵌入，仅实现接口方法即可
type PetService struct {
	uc  *biz.PetUsecase
	app *biz.AppearanceUsecase
	log *log.Helper
}

// NewPetService 创建宠物 HTTP 服务，uc、app 由 wire 注入
func NewPetService(uc *biz.PetUsecase, app *biz.AppearanceUsecase, logger log.Logger) *PetService {
	return &PetService{uc: uc, app: app, log: log.NewHelper(logger)}
}

// getCurrentUserID 从 context 取出当前用户 ID（由 JWT 中间件写入），未登录返回 0 与 err
func (s *PetService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrUnauthorized
	}
	return u.UserID, nil
}

// bizToProto 将 biz.Pet 转为 proto Pet，供 HTTP 返回
func (s *PetService) bizToProto(p *biz.Pet) *petv1.Pet {
	if p == nil {
		return nil
	}
	return &petv1.Pet{
		Id:            p.ID,
		UserId:        p.UserID,
		Name:          p.Name,
		Species:       s.speciesToProto(p.Species),
		BreedId:       p.BreedID,
		AvatarUrl:     p.AvatarURL,
		BackgroundUrl: p.BackgroundURL,
		Mood:          p.Mood,
		Affection:     p.Affection,
		Fullness:      p.Fullness,
		Happiness:     p.Happiness,
		Cleanliness:   p.Cleanliness,
	}
}

func (s *PetService) speciesToProto(species string) petv1.Species {
	switch species {
	case "CAT":
		return petv1.Species_SPECIES_CAT
	case "DOG":
		return petv1.Species_SPECIES_DOG
	default:
		return petv1.Species_SPECIES_UNKNOWN
	}
}

// speciesFromProto 将 proto Species 转为 biz 使用的字符串
func (s *PetService) speciesFromProto(v petv1.Species) string {
	switch v {
	case petv1.Species_SPECIES_CAT:
		return "CAT"
	case petv1.Species_SPECIES_DOG:
		return "DOG"
	default:
		return "CAT"
	}
}

// actionToString 将 proto InteractionAction 转为 biz 使用的字符串
func (s *PetService) actionToString(v petv1.InteractionAction) string {
	switch v {
	case petv1.InteractionAction_INTERACTION_ACTION_STROKE:
		return "STROKE"
	case petv1.InteractionAction_INTERACTION_ACTION_TEASE:
		return "TEASE"
	case petv1.InteractionAction_INTERACTION_ACTION_TREAT:
		return "TREAT"
	case petv1.InteractionAction_INTERACTION_ACTION_BATH:
		return "BATH"
	default:
		return "STROKE"
	}
}

// GetMyPet 获取当前用户主宠物，GET /api/v1/pet/me
func (s *PetService) GetMyPet(ctx context.Context, _ *petv1.GetMyPetRequest) (*petv1.GetMyPetReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.uc.GetMyPet(ctx, userID)
	if err != nil {
		if err == biz.ErrPetNotFound {
			return nil, ErrPetNotFound
		}
		s.log.Errorf("GetMyPet failed: %v", err)
		return nil, err
	}
	return &petv1.GetMyPetReply{Pet: s.bizToProto(p)}, nil
}

// ListPets 获取当前用户宠物列表，GET /api/v1/pet/list
func (s *PetService) ListPets(ctx context.Context, _ *petv1.ListPetsRequest) (*petv1.ListPetsReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.uc.ListPets(ctx, userID)
	if err != nil {
		s.log.Errorf("ListPets failed: %v", err)
		return nil, err
	}
	pets := make([]*petv1.Pet, 0, len(list))
	for _, p := range list {
		pets = append(pets, s.bizToProto(p))
	}
	return &petv1.ListPetsReply{Pets: pets}, nil
}

// CreatePet 创建/领养宠物，POST /api/v1/pet
func (s *PetService) CreatePet(ctx context.Context, req *petv1.CreatePetRequest) (*petv1.CreatePetReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetName() == "" {
		return nil, ErrBadRequest
	}
	species := s.speciesFromProto(req.GetSpecies())
	p, err := s.uc.CreatePet(ctx, userID, req.GetName(), species, req.GetAvatarUrl(), req.GetBackgroundUrl(), req.GetBreedId())
	if err != nil {
		s.log.Errorf("CreatePet failed: %v", err)
		return nil, err
	}
	return &petv1.CreatePetReply{Pet: s.bizToProto(p)}, nil
}

// UpdatePet 更新宠物基础信息，PUT /api/v1/pet/{id}
func (s *PetService) UpdatePet(ctx context.Context, req *petv1.UpdatePetRequest) (*petv1.UpdatePetReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	species := s.speciesFromProto(req.GetSpecies())
	p, err := s.uc.UpdatePet(ctx, userID, req.GetId(), req.GetName(), species, req.GetAvatarUrl(), req.GetBackgroundUrl())
	if err != nil {
		if err == biz.ErrPetNotFound {
			return nil, ErrPetNotFound
		}
		s.log.Errorf("UpdatePet failed: %v", err)
		return nil, err
	}
	return &petv1.UpdatePetReply{Pet: s.bizToProto(p)}, nil
}

// Interact 与宠物互动，POST /api/v1/pet/interact
func (s *PetService) Interact(ctx context.Context, req *petv1.InteractRequest) (*petv1.InteractReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	action := s.actionToString(req.GetAction())
	p, err := s.uc.Interact(ctx, userID, action)
	if err != nil {
		if err == biz.ErrPetNotFound {
			return nil, ErrPetNotFound
		}
		s.log.Errorf("Interact failed: %v", err)
		return nil, err
	}
	return &petv1.InteractReply{Pet: s.bizToProto(p)}, nil
}

// GetPetAppearance 获取宠物当前装扮，GET /api/v1/pet/{id}/appearance
func (s *PetService) GetPetAppearance(ctx context.Context, req *petv1.GetPetAppearanceRequest) (*petv1.GetPetAppearanceReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	head, body, neck, err := s.app.GetAppearance(ctx, userID, req.GetId())
	if err != nil {
		if err == biz.ErrPetNotFound {
			return nil, ErrPetNotFound
		}
		s.log.Errorf("GetPetAppearance failed: %v", err)
		return nil, err
	}
	return &petv1.GetPetAppearanceReply{
		HeadItemId: head,
		BodyItemId: body,
		NeckItemId: neck,
	}, nil
}

// UpdatePetAppearance 保存宠物装扮，PUT /api/v1/pet/{id}/appearance
func (s *PetService) UpdatePetAppearance(ctx context.Context, req *petv1.UpdatePetAppearanceRequest) (*petv1.UpdatePetAppearanceReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.app.UpdateAppearance(ctx, userID, req.GetId(), req.GetHeadItemId(), req.GetBodyItemId(), req.GetNeckItemId())
	if err != nil {
		if err == biz.ErrPetNotFound {
			return nil, ErrPetNotFound
		}
		s.log.Errorf("UpdatePetAppearance failed: %v", err)
		return nil, err
	}
	return &petv1.UpdatePetAppearanceReply{
		HeadItemId: req.GetHeadItemId(),
		BodyItemId: req.GetBodyItemId(),
		NeckItemId: req.GetNeckItemId(),
	}, nil
}
