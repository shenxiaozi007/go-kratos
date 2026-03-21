package service

import (
	"context"

	wechatv1 "realworld/api/wechat/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WechatService struct {
	uc  *biz.WechatUsecase
	log *log.Helper
}

// NewWechatService 创建 Wechat 服务实例。
//
// 作用：
// - 作为 HTTP/gRPC 入口承接接口请求；
// - 将参数转发到 biz 层并返回统一响应结构。
func NewWechatService(uc *biz.WechatUsecase, logger log.Logger) *WechatService {
	return &WechatService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// SendGroupText 处理群文本发送请求。
//
// 返回约定：
// - 成功时 success=true；
// - 失败时返回 Kratos 标准错误结构（由上层中间件统一序列化）。
func (s *WechatService) SendGroupText(ctx context.Context, req *wechatv1.SendGroupTextRequest) (*wechatv1.SendGroupTextReply, error) {
	messageID, err := s.uc.SendGroupText(ctx, req.GetRoomId(), req.GetRoomTopic(), req.GetText())
	if err != nil {
		return nil, err
	}
	return &wechatv1.SendGroupTextReply{
		Success:   true,
		MessageId: messageID,
	}, nil
}

// SendGroupImage 处理群图片发送请求。
//
// 兼容性说明：
// - 同时支持 URL 与字节流两种输入，具体二选一校验在 biz 层完成；
// - service 层仅负责参数透传与返回结构封装。
func (s *WechatService) SendGroupImage(ctx context.Context, req *wechatv1.SendGroupImageRequest) (*wechatv1.SendGroupImageReply, error) {
	messageID, err := s.uc.SendGroupImage(ctx, req.GetRoomId(), req.GetRoomTopic(), req.GetImageUrl(), req.GetImageBytes(), req.GetFilename())
	if err != nil {
		return nil, err
	}
	return &wechatv1.SendGroupImageReply{
		Success:   true,
		MessageId: messageID,
	}, nil
}

// SendGroupMentionText 处理群 @ 文本发送请求。
func (s *WechatService) SendGroupMentionText(ctx context.Context, req *wechatv1.SendGroupMentionTextRequest) (*wechatv1.SendGroupMentionTextReply, error) {
	messageID, err := s.uc.SendGroupMentionText(ctx, req.GetRoomId(), req.GetRoomTopic(), req.GetMentionIds(), req.GetText())
	if err != nil {
		return nil, err
	}
	return &wechatv1.SendGroupMentionTextReply{
		Success:   true,
		MessageId: messageID,
	}, nil
}

// GetLoginQrcode 查询当前可用登录二维码。
//
// 注意：
// - 若尚未触发扫码事件，qrcode 可能为空字符串；
// - 前端应结合登录状态接口一起判断展示逻辑。
func (s *WechatService) GetLoginQrcode(ctx context.Context, _ *emptypb.Empty) (*wechatv1.GetLoginQrcodeReply, error) {
	qrcode, qrcodeDataURL, err := s.uc.GetLoginQrcode(ctx)
	if err != nil {
		return nil, err
	}
	return &wechatv1.GetLoginQrcodeReply{
		Qrcode:        qrcode,
		QrcodeDataUrl: qrcodeDataURL,
	}, nil
}

// GetLoginStatus 查询机器人当前登录状态。
func (s *WechatService) GetLoginStatus(ctx context.Context, _ *emptypb.Empty) (*wechatv1.GetLoginStatusReply, error) {
	status, err := s.uc.GetLoginStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &wechatv1.GetLoginStatusReply{
		LoggedIn: status.LoggedIn,
		SelfId:   status.SelfID,
		SelfName: status.SelfName,
	}, nil
}

// DiagStatus 读取 Wechat 网关诊断状态，供 HTTP 自定义路由透传使用。
func (s *WechatService) DiagStatus(ctx context.Context) (*biz.WechatDiagStatus, error) {
	return s.uc.GetDiagStatus(ctx)
}
