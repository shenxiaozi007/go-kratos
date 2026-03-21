package biz

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type WechatLoginStatus struct {
	LoggedIn bool
	SelfID   string
	SelfName string
}

// WechatDiagStatus 描述当前 Wechat 网关运行诊断快照。
type WechatDiagStatus struct {
	Mode           string
	HasToken       bool
	ScanReceived   bool
	StartupError   string
	QrcodeReady    bool
	LoggedIn       bool
	SelfID         string
	SelfName       string
}

// WechatGateway 定义 Wechat 基础能力抽象。
//
// 设计目的：
// - 隔离上层业务与具体微信协议实现（可切换 bridge、本地 SDK 或其他实现）；
// - 将发送与登录态读取统一到同一网关，便于 service 层保持稳定调用方式。
type WechatGateway interface {
	SendGroupText(ctx context.Context, roomID, roomTopic, text string) (string, error)
	SendGroupImage(ctx context.Context, roomID, roomTopic, imageURL string, imageBytes []byte, filename string) (string, error)
	SendGroupMentionText(ctx context.Context, roomID, roomTopic string, mentionIDs []string, text string) (string, error)
	GetLoginQrcode(ctx context.Context) (qrcode, qrcodeDataURL string, err error)
	GetLoginStatus(ctx context.Context) (*WechatLoginStatus, error)
	GetDiagStatus(ctx context.Context) (*WechatDiagStatus, error)
}

type WechatUsecase struct {
	gateway WechatGateway
	log     *log.Helper
}

// NewWechatUsecase 创建 Wechat 业务用例。
//
// 关键行为：
// - 注入网关依赖用于执行真实消息发送；
// - 注入日志组件用于记录发送链路异常，便于线上排障。
func NewWechatUsecase(gateway WechatGateway, logger log.Logger) *WechatUsecase {
	return &WechatUsecase{
		gateway: gateway,
		log:     log.NewHelper(logger),
	}
}

// SendGroupText 发送群文本消息。
//
// 边界约束：
// - room_id 与 room_topic 至少提供一个；
// - text 不能为空字符串。
func (uc *WechatUsecase) SendGroupText(ctx context.Context, roomID, roomTopic, text string) (string, error) {
	if err := validateRoomTarget(roomID, roomTopic); err != nil {
		return "", err
	}
	if strings.TrimSpace(text) == "" {
		return "", errors.New(400, "BAD_REQUEST", "text is required")
	}
	return uc.gateway.SendGroupText(ctx, roomID, roomTopic, text)
}

// SendGroupImage 发送群图片消息。
//
// 输入约束：
// - 图片来源二选一：image_url 或 image_bytes；
// - 两者都为空或同时存在都视为非法请求。
func (uc *WechatUsecase) SendGroupImage(ctx context.Context, roomID, roomTopic, imageURL string, imageBytes []byte, filename string) (string, error) {
	if err := validateRoomTarget(roomID, roomTopic); err != nil {
		return "", err
	}
	if strings.TrimSpace(imageURL) == "" && len(imageBytes) == 0 {
		return "", errors.New(400, "BAD_REQUEST", "image_url or image_bytes is required")
	}
	if strings.TrimSpace(imageURL) != "" && len(imageBytes) != 0 {
		return "", errors.New(400, "BAD_REQUEST", "provide only one image source")
	}
	return uc.gateway.SendGroupImage(ctx, roomID, roomTopic, imageURL, imageBytes, filename)
}

// SendGroupMentionText 发送群 @ 消息。
//
// 业务规则：
// - mention_ids 至少包含一个目标；
// - text 不能为空；
// - 群目标解析规则与其它发送接口保持一致。
func (uc *WechatUsecase) SendGroupMentionText(ctx context.Context, roomID, roomTopic string, mentionIDs []string, text string) (string, error) {
	if err := validateRoomTarget(roomID, roomTopic); err != nil {
		return "", err
	}
	if len(mentionIDs) == 0 {
		return "", errors.New(400, "BAD_REQUEST", "mention_ids is required")
	}
	if strings.TrimSpace(text) == "" {
		return "", errors.New(400, "BAD_REQUEST", "text is required")
	}
	return uc.gateway.SendGroupMentionText(ctx, roomID, roomTopic, mentionIDs, text)
}

// GetLoginQrcode 读取当前登录二维码快照。
func (uc *WechatUsecase) GetLoginQrcode(ctx context.Context) (qrcode, qrcodeDataURL string, err error) {
	return uc.gateway.GetLoginQrcode(ctx)
}

// GetLoginStatus 读取当前机器人登录状态。
func (uc *WechatUsecase) GetLoginStatus(ctx context.Context) (*WechatLoginStatus, error) {
	return uc.gateway.GetLoginStatus(ctx)
}

// GetDiagStatus 读取当前 Wechat 网关运行诊断信息。
func (uc *WechatUsecase) GetDiagStatus(ctx context.Context) (*WechatDiagStatus, error) {
	return uc.gateway.GetDiagStatus(ctx)
}

// validateRoomTarget 校验群定位参数。
//
// 兼容策略：
// - 支持 room_id 或 room_topic 任一方式定位；
// - 两者都缺失时返回统一 BAD_REQUEST，减少上层重复校验代码。
func validateRoomTarget(roomID, roomTopic string) error {
	if strings.TrimSpace(roomID) == "" && strings.TrimSpace(roomTopic) == "" {
		return errors.New(400, "BAD_REQUEST", "room_id or room_topic is required")
	}
	return nil
}
