package service

import (
	"context"
	"customer/api/verifyCode"
	"regexp"
	"time"

	pb "customer/api/customer"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) GetCaptcha(ctx context.Context, req *pb.GetCaptchaRequest) (*pb.GetCaptchaReply, error) {
	// 校验手机号
	if req.Telephone == "" {
		return &pb.GetCaptchaReply{
			Code:    1,
			Message: "手机号不能为空",
		}, nil
	}

	phoneRegex := `^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`

	// 正则校验
	if !regexp.MustCompile(phoneRegex).MatchString(req.Telephone) {
		return &pb.GetCaptchaReply{
			Code:    1,
			Message: "手机号格式不正确",
		}, nil
	}

	// 连接grpc服务器 用dialInsecure
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)

	defer func() {
		_ = conn.Close()
	}()

	if err != nil {
		return &pb.GetCaptchaReply{
			Code:    1,
			Message: "连接服务器失败",
		}, nil
	}

	// 发送验证码请求
	client := verifyCode.NewVerifyCodeClient(conn)
	
	reply, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   1,
	})

	return &pb.GetCaptchaReply{
		Code:             0,
		Message:          "验证码发送成功",
		VerifyCode:       reply.Code,
		VerifyCodeExpire: int32(time.Now().Add(time.Minute * 5).Unix()),
		VerifyCodeLife:   60,
	}, nil
}
