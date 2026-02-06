package service

import (
	"context"
	v1 "realworld/api/realworld/v1"
	"realworld/internal/biz"
)

// RealWorldService is a greeter service.
type RealWorldService struct {
	v1.UnimplementedRealworldServer

	uc *biz.GreeterUsecase
}

// NewRealWorldService new a greeter service.
func NewRealWorldService(uc *biz.GreeterUsecase) *RealWorldService {
	return &RealWorldService{uc: uc}
}

func (ws *RealWorldService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.UserReply, error) {
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Email: "1212",
		},
	}, nil
}
