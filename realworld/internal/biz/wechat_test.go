package biz

import (
	"context"
	"io"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

type mockWechatGateway struct{}

func (m *mockWechatGateway) SendGroupText(context.Context, string, string, string) (string, error) {
	return "msg-text-1", nil
}

func (m *mockWechatGateway) SendGroupImage(context.Context, string, string, string, []byte, string) (string, error) {
	return "msg-image-1", nil
}

func (m *mockWechatGateway) SendGroupMentionText(context.Context, string, string, []string, string) (string, error) {
	return "msg-mention-1", nil
}

func (m *mockWechatGateway) GetLoginQrcode(context.Context) (string, string, error) {
	return "qrcode", "data:image/png;base64,xxx", nil
}

func (m *mockWechatGateway) GetLoginStatus(context.Context) (*WechatLoginStatus, error) {
	return &WechatLoginStatus{LoggedIn: true, SelfID: "wxid_1", SelfName: "bot"}, nil
}

func TestWechatUsecase_SendGroupText_Validate(t *testing.T) {
	uc := NewWechatUsecase(&mockWechatGateway{}, log.NewStdLogger(io.Discard))
	if _, err := uc.SendGroupText(context.Background(), "", "", "hello"); err == nil {
		t.Fatalf("expected room target validation error")
	}
	if _, err := uc.SendGroupText(context.Background(), "room1", "", ""); err == nil {
		t.Fatalf("expected empty text validation error")
	}
}

func TestWechatUsecase_SendGroupImage_Validate(t *testing.T) {
	uc := NewWechatUsecase(&mockWechatGateway{}, log.NewStdLogger(io.Discard))
	if _, err := uc.SendGroupImage(context.Background(), "room1", "", "", nil, ""); err == nil {
		t.Fatalf("expected image source validation error")
	}
	if _, err := uc.SendGroupImage(context.Background(), "room1", "", "http://img", []byte("x"), "a.png"); err == nil {
		t.Fatalf("expected dual image source validation error")
	}
}

func TestWechatUsecase_SendGroupMentionText_Validate(t *testing.T) {
	uc := NewWechatUsecase(&mockWechatGateway{}, log.NewStdLogger(io.Discard))
	if _, err := uc.SendGroupMentionText(context.Background(), "room1", "", nil, "hi"); err == nil {
		t.Fatalf("expected mention ids validation error")
	}
	if _, err := uc.SendGroupMentionText(context.Background(), "room1", "", []string{"wxid_2"}, ""); err == nil {
		t.Fatalf("expected text validation error")
	}
}
