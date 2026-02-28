package server

import (
	"context"
	v2 "realworld/api/blog/v1"
	v1 "realworld/api/realworld/v1"
	todov1 "realworld/api/todo/v1"
	"realworld/internal/conf"
	"realworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
// 该函数负责根据配置创建一个 HTTP 服务器实例，并注册所有 HTTP 路由。
// 我们在此函数中额外注册 Todo 模块相关的 RESTful API。
func NewHTTPServer(
	c *conf.Server,
	greeter *service.RealWorldService,
	logger log.Logger,
	articleService *service.ArticleService,
	todoService *service.TodoService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			//middeware.JWTAuth(),
			CorsMiddleware(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// 注册 CORS 处理：在所有路由之前注册一个全局的 OPTIONS 处理器，
	// 用于处理浏览器的 CORS 预检请求（preflight request）。
	// 这样即使路由不存在，OPTIONS 请求也能被正确处理。
	r := srv.Route("/")
	r.OPTIONS("/api/{path:.*}", handleCORS)

	v1.RegisterRealworldHTTPServer(srv, greeter)

	v2.RegisterArticleHTTPServer(srv, articleService)

	// 使用 proto 生成的 HTTP 注册函数来注册 Todo 模块路由，
	// 完全遵循 Kratos 的 proto 驱动方式。
	// 具体路由定义见 api/todo/v1/todo.proto 中的 google.api.http 注解。
	todov1.RegisterTodoServiceHTTPServer(srv, todoService)

	return srv
}

// handleCORS 处理 CORS 预检请求（OPTIONS）。
// 当浏览器发送跨域请求时，会先发送一个 OPTIONS 请求来检查服务器是否允许该请求。
func handleCORS(ctx http.Context) error {
	// 设置 CORS 响应头
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	ctx.Response().Header().Set("Access-Control-Max-Age", "3600")
	// 返回 200 状态码，表示允许跨域请求
	return ctx.Result(200, map[string]string{"message": "OK"})
}

// CorsMiddleware 是一个 HTTP 中间件，用于在所有响应中添加 CORS 头部。
// 这样即使不是 OPTIONS 请求，正常的 GET/POST 等请求也会携带 CORS 头部。
func CorsMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从 context 中获取 transport 信息
			tr, ok := transport.FromServerContext(ctx)
			if ok {
				// 设置 CORS 响应头
				tr.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
				tr.ReplyHeader().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
				tr.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			}

			// 调用下一个中间件或处理器
			return handler(ctx, req)
		}
	}
}
