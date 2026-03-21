// Package service 萌宠之家（cat-admin）社交 HTTP 服务层
// 实现 api/friend/v1.FriendServiceHTTPServer，见 docs-all/cat-admin/02-后端开发文档.md

package service

import (
	"context"

	friendv1 "realworld/api/friend/v1"
	"realworld/internal/biz"
	"realworld/pkg/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// ErrFriendUnauthorized 未登录
var ErrFriendUnauthorized = errors.New(401, "UNAUTHORIZED", "请先登录")

// FriendService 实现 friend.v1.FriendServiceHTTPServer
type FriendService struct {
	uc  *biz.FriendUsecase
	userRepo biz.UserRepo
	log *log.Helper
}

// NewFriendService 创建社交服务
func NewFriendService(uc *biz.FriendUsecase, userRepo biz.UserRepo, logger log.Logger) *FriendService {
	return &FriendService{uc: uc, userRepo: userRepo, log: log.NewHelper(logger)}
}

func (s *FriendService) getCurrentUserID(ctx context.Context) (uint64, error) {
	u := auth.FromContext(ctx)
	if u == nil {
		return 0, ErrFriendUnauthorized
	}
	return u.UserID, nil
}

func userToBrief(u *biz.User) *friendv1.UserBrief {
	if u == nil {
		return nil
	}
	return &friendv1.UserBrief{
		Id:          u.ID,
		Username:    u.Username,
		Nickname:    u.Nickname,
		AvatarUrl:   u.AvatarURL,
		Level:       u.Level,
		HeartPoints: u.HeartPoints,
	}
}

// ListFriends GET /api/v1/friends
func (s *FriendService) ListFriends(ctx context.Context, _ *friendv1.ListFriendsRequest) (*friendv1.ListFriendsReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	list, err := s.uc.ListFriends(ctx, userID)
	if err != nil {
		s.log.Errorf("ListFriends failed: %v", err)
		return nil, err
	}
	friends := make([]*friendv1.UserBrief, 0, len(list))
	for _, u := range list {
		friends = append(friends, userToBrief(u))
	}
	return &friendv1.ListFriendsReply{Friends: friends}, nil
}

// SearchUsers GET /api/v1/friends/search
func (s *FriendService) SearchUsers(ctx context.Context, req *friendv1.SearchUsersRequest) (*friendv1.SearchUsersReply, error) {
	_, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	limit := 20
	if req.GetKeyword() == "" {
		return &friendv1.SearchUsersReply{Users: nil}, nil
	}
	list, err := s.uc.SearchUsers(ctx, req.GetKeyword(), limit)
	if err != nil {
		s.log.Errorf("SearchUsers failed: %v", err)
		return nil, err
	}
	users := make([]*friendv1.UserBrief, 0, len(list))
	for _, u := range list {
		users = append(users, userToBrief(u))
	}
	return &friendv1.SearchUsersReply{Users: users}, nil
}

// SendFriendRequest POST /api/v1/friends/request
func (s *FriendService) SendFriendRequest(ctx context.Context, req *friendv1.SendFriendRequestRequest) (*friendv1.SendFriendRequestReply, error) {
	fromID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.uc.SendFriendRequest(ctx, fromID, req.GetToUserId(), req.GetMessage())
	if err != nil {
		if err == biz.ErrCannotRequestSelf {
			return &friendv1.SendFriendRequestReply{Success: false, Message: "不能添加自己为好友"}, nil
		}
		if err == biz.ErrTargetUserNotFound {
			return &friendv1.SendFriendRequestReply{Success: false, Message: "用户不存在"}, nil
		}
		s.log.Errorf("SendFriendRequest failed: %v", err)
		return nil, err
	}
	return &friendv1.SendFriendRequestReply{Success: true, Message: "已发送申请"}, nil
}

// ListFriendRequests GET /api/v1/friends/requests
func (s *FriendService) ListFriendRequests(ctx context.Context, _ *friendv1.ListFriendRequestsRequest) (*friendv1.ListFriendRequestsReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	reqs, err := s.uc.ListFriendRequests(ctx, userID)
	if err != nil {
		s.log.Errorf("ListFriendRequests failed: %v", err)
		return nil, err
	}
	items := make([]*friendv1.FriendRequestItem, 0, len(reqs))
	for _, r := range reqs {
		fromUser, _ := s.userRepo.GetByID(ctx, r.FromUserID)
		item := &friendv1.FriendRequestItem{
			Id:         r.ID,
			FromUserId: r.FromUserID,
			FromUser:   userToBrief(fromUser),
			Message:    r.Message,
			Status:     r.Status,
			CreatedAt:  r.CreatedAt.Unix(),
		}
		items = append(items, item)
	}
	return &friendv1.ListFriendRequestsReply{Requests: items}, nil
}

// AcceptFriendRequest POST /api/v1/friends/requests/{id}/accept
func (s *FriendService) AcceptFriendRequest(ctx context.Context, req *friendv1.AcceptFriendRequestRequest) (*friendv1.AcceptFriendRequestReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.uc.AcceptFriendRequest(ctx, userID, req.GetId())
	if err != nil {
		if err == biz.ErrFriendRequestNotFound {
			return &friendv1.AcceptFriendRequestReply{Success: false, Message: "申请不存在或已处理"}, nil
		}
		s.log.Errorf("AcceptFriendRequest failed: %v", err)
		return nil, err
	}
	return &friendv1.AcceptFriendRequestReply{Success: true, Message: "已同意"}, nil
}

// RejectFriendRequest POST /api/v1/friends/requests/{id}/reject
func (s *FriendService) RejectFriendRequest(ctx context.Context, req *friendv1.RejectFriendRequestRequest) (*friendv1.RejectFriendRequestReply, error) {
	userID, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.uc.RejectFriendRequest(ctx, userID, req.GetId())
	if err != nil {
		if err == biz.ErrFriendRequestNotFound {
			return &friendv1.RejectFriendRequestReply{Success: false, Message: "申请不存在或已处理"}, nil
		}
		s.log.Errorf("RejectFriendRequest failed: %v", err)
		return nil, err
	}
	return &friendv1.RejectFriendRequestReply{Success: true, Message: "已拒绝"}, nil
}

// GetRanking GET /api/v1/ranking
func (s *FriendService) GetRanking(ctx context.Context, req *friendv1.GetRankingRequest) (*friendv1.GetRankingReply, error) {
	_, err := s.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	rankType := req.GetType()
	if rankType == "" {
		rankType = "heart_points"
	}
	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 50
	}
	list, err := s.uc.GetRanking(ctx, rankType, limit)
	if err != nil {
		s.log.Errorf("GetRanking failed: %v", err)
		return nil, err
	}
	ranking := make([]*friendv1.UserBrief, 0, len(list))
	for _, u := range list {
		ranking = append(ranking, userToBrief(u))
	}
	return &friendv1.GetRankingReply{Ranking: ranking}, nil
}
