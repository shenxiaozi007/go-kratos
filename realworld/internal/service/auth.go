// Package service 实现 auth.v1 的 HTTP 服务：注册、登录、获取当前用户资料（v2 鉴权）。
package service

import (
	"context"
	stderrors "errors"

	authv1 "realworld/api/auth/v1"
	"realworld/internal/biz"
	"realworld/internal/conf"
	"realworld/pkg/auth"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func errBadRequest(msg string) error   { return kerrors.New(400, "BAD_REQUEST", msg) }
func errUnauthorized(msg string) error { return kerrors.New(401, "UNAUTHORIZED", msg) }

// AuthService 实现 api/auth/v1 AuthServiceHTTPServer，提供注册、登录、GetProfile。
type AuthService struct {
	uc     *biz.AuthUsecase
	conf   *conf.Server
	logger *log.Helper
}

// NewAuthService 创建鉴权 HTTP 服务。需传入鉴权用例与服务器配置（用于 JWT 签发与过期时间）。
func NewAuthService(uc *biz.AuthUsecase, c *conf.Server, logger log.Logger) *AuthService {
	return &AuthService{uc: uc, conf: c, logger: log.NewHelper(logger)}
}

// Register 注册新用户：校验用户名唯一、bcrypt 哈希密码、写入 users 表；首个用户角色为 admin。
func (s *AuthService) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterReply, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errBadRequest("username and password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u, err := s.uc.Register(ctx, req.Username, string(hash))
	if err != nil {
		if stderrors.Is(err, biz.ErrUserExists) {
			return nil, errBadRequest("username already exists")
		}
		return nil, err
	}
	return &authv1.RegisterReply{UserId: int64(u.ID), Username: u.Username}, nil
}

// Login 登录：校验用户名与密码，签发 JWT 并返回 access_token、expires_in。
func (s *AuthService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginReply, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errBadRequest("username and password required")
	}
	ac := s.conf.GetAuth()
	if ac == nil {
		return nil, errBadRequest("auth not configured")
	}
	secret := ac.GetJwtSecret()
	if secret == "" {
		return nil, errBadRequest("auth not configured")
	}
	expireSeconds := ac.GetJwtExpireSeconds()
	if expireSeconds <= 0 {
		expireSeconds = 86400
	}

	u, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		if stderrors.Is(err, biz.ErrAuthUserNotFound) || stderrors.Is(err, biz.ErrBadPassword) {
			return nil, errUnauthorized("invalid username or password")
		}
		return nil, err
	}

	token, _, err := auth.Sign(secret, u.ID, u.Username, u.Role, expireSeconds)
	if err != nil {
		return nil, err
	}
	return &authv1.LoginReply{
		AccessToken: token,
		ExpiresIn:   int64(expireSeconds),
	}, nil
}

// GetProfile 返回当前登录用户信息（含 role）；依赖 JWT 中间件已将用户写入 context。
func (s *AuthService) GetProfile(ctx context.Context, _ *emptypb.Empty) (*authv1.GetProfileReply, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return nil, errUnauthorized("invalid or expired token")
	}
	po, err := s.uc.GetByID(ctx, u.UserID)
	if err != nil {
		return nil, err
	}
	return &authv1.GetProfileReply{
		UserId:   int64(po.ID),
		Username: po.Username,
		Role:     po.Role,
	}, nil
}
