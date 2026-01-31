package service

import (
	"context"
	"math/rand"

	pb "verifyCode/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode(int(req.Length), req.Type),
	}, nil
}

func RandCode(l int, t pb.TYPE) string {
	var chars string
	switch t {
	case pb.TYPE_DIGIT:
		chars = randCode("1234567890", l)
	case pb.TYPE_LETTER:
		chars = randCode("ABCDEFGHIJKLMNOPQRSTUVWXYZ", l)
	case pb.TYPE_MIXED:
		chars = randCode("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", l)
	}
	//909912 816218

	return chars
}

func randCode(s string, l int) string {
	// 传进来的字符串s 随机取len长度的字符串

	randString := make([]byte, l)
	lenS := len(s)
	for i := 0; i < l; i++ {
		randString[i] = s[rand.Intn(lenS)]
	}

	return string(randString)
}
