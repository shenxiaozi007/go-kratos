// Package middeware 提供 v2 鉴权中间件：JWT 校验、将当前用户写入 context。
// 白名单路径（登录/注册）不校验；其余路径需携带有效 Bearer Token，否则返回 401。
package middeware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"

	"realworld/pkg/auth"
)

// 需要跳过 JWT 校验的路径（小写，无尾斜杠）；用于登录、注册等公开接口。
var skipAuthPaths = map[string]bool{
	"/api/v1/auth/login":    true,
	"/api/v1/auth/register": true,
}

// JWTAuth 返回 JWT 鉴权中间件。secret 为空时不校验（放行所有请求，用于未配置鉴权时兼容）。
// 非白名单路径需在 Header 中携带 Authorization: Bearer <token>，否则返回 401。
func JWTAuth(secret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if secret == "" {
				return handler(ctx, req)
			}

			// 仅 HTTP 传输可获取 Request 路径；非 HTTP 直接放行（如 gRPC 可另行鉴权）。
			reqObj, ok := httptransport.RequestFromServerContext(ctx)
			if ok && reqObj != nil {
				path := reqObj.URL.Path
				if len(path) > 0 && path[len(path)-1] == '/' {
					path = path[:len(path)-1]
				}
				if skipAuthPaths[path] {
					return handler(ctx, req)
				}
			}

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New(500, "INVALID_CTX", "missing transport")
			}
			authHeader := tr.RequestHeader().Get("Authorization")
			tokenString := ""
			if strings.HasPrefix(strings.TrimSpace(authHeader), "Bearer ") {
				tokenString = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(authHeader), "Bearer "))
			}
			if tokenString == "" {
				return nil, errors.New(401, "UNAUTHORIZED", "invalid or expired token")
			}

			claims, err := auth.Parse(secret, tokenString)
			if err != nil {
				return nil, errors.New(401, "UNAUTHORIZED", "invalid or expired token")
			}

			ctx = auth.WithUser(ctx, &auth.UserInfo{
				UserID:   claims.UserID,
				Username: claims.Username,
				Role:     claims.Role,
			})
			return handler(ctx, req)
		}
	}
}
