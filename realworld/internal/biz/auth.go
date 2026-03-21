// Package biz 提供 v2 鉴权业务：注册、登录、获取当前用户资料。
package biz

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 业务错误
var (
	ErrUserExists   = errors.New("username already exists")
	// ErrAuthUserNotFound 避免与其他模块的 ErrUserNotFound 命名冲突。
	ErrAuthUserNotFound = errors.New("user not found")
	ErrBadPassword      = errors.New("invalid password")
)

// User 鉴权用用户实体（与 data.UserPO 解耦，避免 biz 依赖 data）。
// 萌宠之家扩展：Nickname、AvatarURL、Signature、Level、Coins、HeartPoints 用于个人资料与签到奖励。
type User struct {
	ID           uint64
	Username     string
	PasswordHash string
	Role         string
	Nickname     string
	AvatarURL    string
	Signature    string
	Level        int32
	Coins        int32
	HeartPoints  int32
}

// UserRepo 用户仓储接口，由 data 层实现。
// 萌宠之家：UpdateProfile 更新资料；AddCoinsAndHeart 增加金币/爱心（如签到奖励）。
type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uint64) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
	Count(ctx context.Context) (int64, error)
	UpdateProfile(ctx context.Context, userID uint64, nickname, avatarURL, signature string) error
	AddCoinsAndHeart(ctx context.Context, userID uint64, coins, heart int32) error
	DeductCoins(ctx context.Context, userID uint64, amount int32) error // 商店购买扣金币，余额不足返回错误
	// 社交：搜索用户（username/nickname 模糊）、排行榜（按 heart_points 或 level 排序）
	SearchByKeyword(ctx context.Context, keyword string, limit int) ([]*User, error)
	ListByRank(ctx context.Context, rankType string, limit int) ([]*User, error)
}

// AuthUsecase 鉴权用例：注册、登录、按 ID 查用户（供 GetProfile 使用）。
type AuthUsecase struct {
	userRepo UserRepo
}

// NewAuthUsecase 创建鉴权用例。
func NewAuthUsecase(userRepo UserRepo) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

// Register 注册新用户。密码由调用方哈希后传入；角色默认 operator，若当前无用户则设为 admin。
func (uc *AuthUsecase) Register(ctx context.Context, username, passwordHash string) (*User, error) {
	_, err := uc.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	n, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	role := "operator"
	if n == 0 {
		role = "admin"
	}
	u := &User{
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
	}
	return uc.userRepo.Create(ctx, u)
}

// Login 根据用户名与明文密码校验登录（内部使用 bcrypt 与存储哈希比对）；成功返回用户，失败返回 ErrUserNotFound 或 ErrBadPassword。
func (uc *AuthUsecase) Login(ctx context.Context, username, plainPassword string) (*User, error) {
	u, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthUserNotFound
		}
		return nil, err
	}
	if !uc.checkPassword(u.PasswordHash, plainPassword) {
		return nil, ErrBadPassword
	}
	return u, nil
}

// GetByID 根据用户 ID 查询用户，供 GetProfile 使用。
func (uc *AuthUsecase) GetByID(ctx context.Context, userID uint64) (*User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

// checkPassword 使用 bcrypt 校验明文密码与存储哈希是否一致。
func (uc *AuthUsecase) checkPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
