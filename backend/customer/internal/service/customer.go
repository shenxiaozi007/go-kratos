package service

import (
	"context"

	pb "customer/api/customer"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) GetCaptcha(ctx context.Context, req *pb.GetCaptchaRequest) (*pb.GetCaptchaReply, error) {
	return &pb.GetCaptchaReply{}, nil
}

func getVerifyCode() {

}
