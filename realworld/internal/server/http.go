package server

import (
	v2 "realworld/api/blog/v1"
	v1 "realworld/api/realworld/v1"
	todov1 "realworld/api/todo/v1"
	"realworld/internal/conf"
	"realworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
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
	v1.RegisterRealworldHTTPServer(srv, greeter)

	v2.RegisterArticleHTTPServer(srv, articleService)

	// 使用 proto 生成的 HTTP 注册函数来注册 Todo 模块路由，
	// 完全遵循 Kratos 的 proto 驱动方式。
	// 具体路由定义见 api/todo/v1/todo.proto 中的 google.api.http 注解。
	todov1.RegisterTodoServiceHTTPServer(srv, todoService)

	return srv
}
