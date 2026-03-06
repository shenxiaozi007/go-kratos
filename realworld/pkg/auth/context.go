// Package auth 提供 v2 鉴权相关：JWT 签发/校验、context 中用户信息存取。
package auth

import "context"

// contextKey 用于在 context 中存储当前用户信息，避免与业务 key 冲突。
type contextKey struct{ name string }

var userContextKey = &contextKey{"xsh_user"}

// UserInfo 当前请求中的用户信息，由 JWT 中间件解析后写入 context。
type UserInfo struct {
	UserID   uint64
	Username string
	Role     string // operator / admin
}

// WithUser 将用户信息写入 context，供后续业务与权限校验使用。
func WithUser(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// FromContext 从 context 读取当前用户；未登录或未经过 JWT 中间件时返回 nil。
func FromContext(ctx context.Context) *UserInfo {
	v := ctx.Value(userContextKey)
	if v == nil {
		return nil
	}
	u, _ := v.(*UserInfo)
	return u
}
