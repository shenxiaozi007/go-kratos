package server

import (
	"context"
	"io"

	achievementv1 "realworld/api/achievement/v1"
	authv1 "realworld/api/auth/v1"
	v2 "realworld/api/blog/v1"
	breedv1 "realworld/api/breed/v1"
	checkinv1 "realworld/api/checkin/v1"
	friendv1 "realworld/api/friend/v1"
	gamev1 "realworld/api/game/v1"
	inventoryv1 "realworld/api/inventory/v1"
	petv1 "realworld/api/pet/v1"
	v1 "realworld/api/realworld/v1"
	shopv1 "realworld/api/shop/v1"
	statsv1 "realworld/api/stats/v1"
	todov1 "realworld/api/todo/v1"
	userv1 "realworld/api/user/v1"
	wechatv1 "realworld/api/wechat/v1"
	xshv1 "realworld/api/xsh/v1"
	"realworld/internal/conf"
	"realworld/internal/service"
	"realworld/pkg/middeware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
// 该函数负责根据配置创建一个 HTTP 服务器实例，并注册所有 HTTP 路由。
// 我们在此函数中额外注册 Todo 模块相关的 RESTful API。
// NewHTTPServer 创建 HTTP 服务并注册所有路由，包括 xsh 选品/内容/分发/获客截流 API（/api/xsh/v1/...）。
func NewHTTPServer(
	c *conf.Server,
	greeter *service.RealWorldService,
	logger log.Logger,
	authService *service.AuthService,
	articleService *service.ArticleService,
	todoService *service.TodoService,
	gameService *service.GameService,
	statsService *service.StatsService,
	petService *service.PetService,
	checkinService *service.CheckinService,
	userService *service.UserService,
	shopService *service.ShopService,
	inventoryService *service.InventoryService,
	breedService *service.BreedService,
	achievementService *service.AchievementService,
	friendService *service.FriendService,
	xshProductService *service.XshProductService,
	xshContentService *service.XshContentService,
	xshDispatchService *service.XshDispatchService,
	xshInboxService *service.XshInboxService,
	wechatService *service.WechatService,
) *http.Server {
	// v2：若配置了 JWT 密钥，则启用鉴权中间件；/api/v1/auth/login、register 白名单放行，其余需 Bearer Token。
	jwtSecret := ""
	if c.GetAuth() != nil {
		jwtSecret = c.GetAuth().GetJwtSecret()
	}
	middlewares := []middleware.Middleware{recovery.Recovery(), CorsMiddleware()}
	if jwtSecret != "" {
		middlewares = append(middlewares, middeware.JWTAuth(jwtSecret))
	}
	var opts = []http.ServerOption{
		http.Middleware(middlewares...),
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
	r.POST("/api/company/contact", handleCompanyContact(logger))
	r.POST("/api/wechat/v1/group/image/upload", handleWechatImageUpload(wechatService, logger))
	r.GET("/api/wechat/v1/diag", handleWechatDiag(wechatService, logger))

	v1.RegisterRealworldHTTPServer(srv, greeter)

	// v2 鉴权：注册、登录、获取当前用户（GetProfile 需 JWT）
	authv1.RegisterAuthServiceHTTPServer(srv, authService)

	v2.RegisterArticleHTTPServer(srv, articleService)

	// 使用 proto 生成的 HTTP 注册函数来注册 Todo 模块路由，
	// 完全遵循 Kratos 的 proto 驱动方式。
	// 具体路由定义见 api/todo/v1/todo.proto 中的 google.api.http 注解。
	todov1.RegisterTodoServiceHTTPServer(srv, todoService)

	// 注册 GameService 的 HTTP 路由，具体映射见 api/game/v1/game.proto 中的注解。
	gamev1.RegisterGameServiceHTTPServer(srv, gameService)
	statsv1.RegisterStatsServiceHTTPServer(srv, statsService)

	// 萌宠之家（cat-admin）宠物 API，见 docs-all/cat-admin
	petv1.RegisterPetServiceHTTPServer(srv, petService)
	// 萌宠之家（cat-admin）签到 API
	checkinv1.RegisterCheckinServiceHTTPServer(srv, checkinService)
	// 萌宠之家（cat-admin）用户资料 API
	userv1.RegisterUserServiceHTTPServer(srv, userService)
	// 萌宠之家（cat-admin）商店与背包 API
	shopv1.RegisterShopServiceHTTPServer(srv, shopService)
	inventoryv1.RegisterInventoryServiceHTTPServer(srv, inventoryService)
	// 萌宠之家（cat-admin）品种 API
	breedv1.RegisterBreedServiceHTTPServer(srv, breedService)
	// 萌宠之家（cat-admin）成就 API
	achievementv1.RegisterAchievementServiceHTTPServer(srv, achievementService)
	// 萌宠之家（cat-admin）社交 API
	friendv1.RegisterFriendServiceHTTPServer(srv, friendService)

	// xsh 选品/内容/分发/获客截流（见 xsh-docs）
	xshv1.RegisterProductServiceHTTPServer(srv, xshProductService)
	xshv1.RegisterContentServiceHTTPServer(srv, xshContentService)
	xshv1.RegisterDispatchServiceHTTPServer(srv, xshDispatchService)
	xshv1.RegisterInboxServiceHTTPServer(srv, xshInboxService)
	wechatv1.RegisterWechatServiceHTTPServer(srv, wechatService)

	return srv
}

// contactRequest 官网联系表单请求体。
type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// handleCompanyContact 处理官网联系表单：记录日志并返回 200。邮件发送可后续接 SMTP。
func handleCompanyContact(logger log.Logger) func(http.Context) error {
	helper := log.NewHelper(logger)
	return func(ctx http.Context) error {
		var req contactRequest
		if err := ctx.Bind(&req); err != nil {
			helper.Errorf("company contact bind error: %v", err)
			return ctx.Result(400, map[string]string{"error": "invalid request"})
		}
		helper.Infof("company contact form submission: name=%q email=%q message=%q", req.Name, req.Email, req.Message)
		return ctx.Result(200, map[string]interface{}{"ok": true, "message": "received"})
	}
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

type wechatMultipartImageReply struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
}

// handleWechatImageUpload 兼容 multipart/form-data 图片上传，调用同一业务能力发送群图。
func handleWechatImageUpload(wechatService *service.WechatService, logger log.Logger) func(http.Context) error {
	helper := log.NewHelper(logger)
	return func(ctx http.Context) error {
		req := ctx.Request()
		if err := req.ParseMultipartForm(10 << 20); err != nil {
			return ctx.Result(400, map[string]string{"error": "invalid multipart form"})
		}
		roomID := req.FormValue("room_id")
		roomTopic := req.FormValue("room_topic")
		file, fileHeader, err := req.FormFile("file")
		if err != nil {
			return ctx.Result(400, map[string]string{"error": "file is required"})
		}
		defer file.Close()
		imageBytes, err := io.ReadAll(file)
		if err != nil {
			helper.Errorf("read upload image failed: %v", err)
			return ctx.Result(500, map[string]string{"error": "read file failed"})
		}
		reply, err := wechatService.SendGroupImage(ctx, &wechatv1.SendGroupImageRequest{
			RoomId:     roomID,
			RoomTopic:  roomTopic,
			ImageBytes: imageBytes,
			Filename:   fileHeader.Filename,
		})
		if err != nil {
			return err
		}
		return ctx.Result(200, &wechatMultipartImageReply{
			Success:   reply.GetSuccess(),
			MessageID: reply.GetMessageId(),
		})
	}
}

// handleWechatDiag 提供 Wechat 运行态诊断信息，便于快速排查二维码/登录异常。
func handleWechatDiag(wechatService *service.WechatService, logger log.Logger) func(http.Context) error {
	helper := log.NewHelper(logger)
	return func(ctx http.Context) error {
		diag, err := wechatService.DiagStatus(ctx)
		if err != nil {
			helper.Errorf("wechat diag failed: %v", err)
			return err
		}
		return ctx.Result(200, map[string]interface{}{
			"mode":          diag.Mode,
			"has_token":     diag.HasToken,
			"scan_received": diag.ScanReceived,
			"qrcode_ready":  diag.QrcodeReady,
			"logged_in":     diag.LoggedIn,
			"self_id":       diag.SelfID,
			"self_name":     diag.SelfName,
			"startup_error": diag.StartupError,
		})
	}
}
