// Package service 萌宠之家（cat-admin）用户资料 HTTP 服务层
// 实现 api/user/v1.UserServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	userv1 "realworld/api/user/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrUserUnauthorized 未登录
var ErrUserUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// UserService 实现 user.v1.UserServiceHTTPServer
type UserService struct {
	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService 创建用户资料服务
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

func (s *UserService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrUserUnauthorized
	}
	return u.UserID, nil
}

func (s *UserService) userToProto(u *biz.User) *userv1.GetProfileReply {
	if u == nil {
		return nil
	}
	return &userv1.GetProfileReply{
		UserId:      int64(u.ID),
		Username:    u.Username,
		Nickname:    u.Nickname,
		AvatarUrl:   u.AvatarURL,
		Signature:   u.Signature,
		Level:       u.Level,
		Coins:       u.Coins,
		HeartPoints: u.HeartPoints,
	}
}

// GetProfile GET /api/v1/user/profile
func (s *UserService) GetProfile(ctx context.Context, _ *userv1.GetProfileRequest) (*userv1.GetProfileReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	u, err := s.uc.GetProfile(ctx, userID)
	if err != nil {
		s.log.Errorf("GetProfile failed: %v", err)
		return nil, err
	}
	return s.userToProto(u), nil
}

// UpdateProfile PUT /api/v1/user/profile
func (s *UserService) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	u, err := s.uc.UpdateProfile(ctx, userID, req.GetNickname(), req.GetAvatarUrl(), req.GetSignature())
	if err != nil {
		s.log.Errorf("UpdateProfile failed: %v", err)
		return nil, err
	}
	return &userv1.UpdateProfileReply{Profile: s.userToProto(u)}, nil
}
