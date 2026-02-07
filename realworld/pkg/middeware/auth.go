package middeware

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// 中间件
func JWTAuth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New(500, "invalid authorization", "")
			}

			auth := tr.RequestHeader().Get("Authorization")
			spew.Dump(auth)
			if auth == "" {
				return nil, errors.New(401, "invalid authorization", "")
			}

			return handler(ctx, req)
		}
	}
}
