// Package biz 萌宠之家（cat-admin）社交业务层
// 定义 FriendRepo、FriendRequestRepo、FriendUsecase，见 docs-all/cat-admin/02-后端开发文档.md

package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// FriendRequest 好友申请领域实体
type FriendRequest struct {
	ID         int64
	FromUserID uint64
	ToUserID   uint64
	Message    string
	Status     string // PENDING / ACCEPTED / REJECTED
	CreatedAt  time.Time
}

// FriendRepo 好友关系仓储
type FriendRepo interface {
	ListFriendIDs(ctx context.Context, userID uint64) ([]uint64, error)
	AddFriendship(ctx context.Context, userID, friendID uint64) error
}

// FriendRequestRepo 好友申请仓储
type FriendRequestRepo interface {
	Create(ctx context.Context, fromUserID, toUserID uint64, message string) (*FriendRequest, error)
	GetByID(ctx context.Context, id int64) (*FriendRequest, error)
	ListReceivedByToUser(ctx context.Context, toUserID uint64) ([]*FriendRequest, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

var (
	ErrFriendRequestNotFound = errors.New("friend request not found")
	ErrCannotRequestSelf     = errors.New("cannot send request to self")
	ErrTargetUserNotFound    = errors.New("target user not found")
	ErrDuplicateRequest     = errors.New("already sent or already friends")
)

// FriendUsecase 社交用例：好友列表、搜索、申请、同意/拒绝、排行榜
type FriendUsecase struct {
	friendRepo        FriendRepo
	friendRequestRepo FriendRequestRepo
	userRepo          UserRepo
	log               *log.Helper
}

// NewFriendUsecase 创建社交用例
func NewFriendUsecase(friendRepo FriendRepo, friendRequestRepo FriendRequestRepo, userRepo UserRepo, logger log.Logger) *FriendUsecase {
	return &FriendUsecase{
		friendRepo:        friendRepo,
		friendRequestRepo: friendRequestRepo,
		userRepo:          userRepo,
		log:               log.NewHelper(logger),
	}
}

// ListFriends 返回当前用户的好友列表（*User），需根据 friendIDs 查 UserRepo.GetByID
func (uc *FriendUsecase) ListFriends(ctx context.Context, userID uint64) ([]*User, error) {
	ids, err := uc.friendRepo.ListFriendIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	list := make([]*User, 0, len(ids))
	for _, id := range ids {
		u, err := uc.userRepo.GetByID(ctx, id)
		if err != nil {
			continue
		}
		if u != nil {
			list = append(list, u)
		}
	}
	return list, nil
}

// SearchUsers 搜索用户（委托 UserRepo）
func (uc *FriendUsecase) SearchUsers(ctx context.Context, keyword string, limit int) ([]*User, error) {
	return uc.userRepo.SearchByKeyword(ctx, keyword, limit)
}

// SendFriendRequest 发送好友申请：不能给自己发、对方需存在、可选：检查是否已存在申请或已是好友
func (uc *FriendUsecase) SendFriendRequest(ctx context.Context, fromUserID, toUserID uint64, message string) error {
	if fromUserID == toUserID {
		return ErrCannotRequestSelf
	}
	_, err := uc.userRepo.GetByID(ctx, toUserID)
	if err != nil {
		return err
	}
	// 简化：不检查是否已有 PENDING 或已是好友，直接创建（可后续加唯一约束或业务校验）
	_, err = uc.friendRequestRepo.Create(ctx, fromUserID, toUserID, message)
	return err
}

// ListFriendRequests 收到的好友申请列表
func (uc *FriendUsecase) ListFriendRequests(ctx context.Context, toUserID uint64) ([]*FriendRequest, error) {
	return uc.friendRequestRepo.ListReceivedByToUser(ctx, toUserID)
}

// AcceptFriendRequest 同意申请：校验 to_user_id 为当前用户、状态为 PENDING，更新状态并建立双向好友
func (uc *FriendUsecase) AcceptFriendRequest(ctx context.Context, toUserID uint64, requestID int64) error {
	req, err := uc.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil || req == nil {
		return ErrFriendRequestNotFound
	}
	if req.ToUserID != toUserID {
		return ErrFriendRequestNotFound
	}
	if req.Status != "PENDING" {
		return ErrFriendRequestNotFound
	}
	if err := uc.friendRequestRepo.UpdateStatus(ctx, requestID, "ACCEPTED"); err != nil {
		return err
	}
	return uc.friendRepo.AddFriendship(ctx, req.FromUserID, req.ToUserID)
}

// RejectFriendRequest 拒绝申请
func (uc *FriendUsecase) RejectFriendRequest(ctx context.Context, toUserID uint64, requestID int64) error {
	req, err := uc.friendRequestRepo.GetByID(ctx, requestID)
	if err != nil || req == nil {
		return ErrFriendRequestNotFound
	}
	if req.ToUserID != toUserID {
		return ErrFriendRequestNotFound
	}
	if req.Status != "PENDING" {
		return ErrFriendRequestNotFound
	}
	return uc.friendRequestRepo.UpdateStatus(ctx, requestID, "REJECTED")
}

// GetRanking 排行榜（委托 UserRepo）
func (uc *FriendUsecase) GetRanking(ctx context.Context, rankType string, limit int) ([]*User, error) {
	return uc.userRepo.ListByRank(ctx, rankType, limit)
}
