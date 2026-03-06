package service

import "github.com/google/wire"

// ProviderSet is service providers.
// 包含 auth（v2 鉴权）与 xsh 选品/内容/分发/获客截流 HTTP 服务。
var ProviderSet = wire.NewSet(
	NewAuthService,
	NewRealWorldService,
	NewArticleService,
	NewTodoService,
	NewGameService,
	NewStatsService,
	NewXshProductService,
	NewXshContentService,
	NewXshDispatchService,
	NewXshInboxService,
)
