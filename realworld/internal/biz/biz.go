package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
// 包含 xsh 选品/内容/分发/获客截流用例及内容加工坊占位实现（Copy/Cover/Video 未配置时使用 Noop）。
var ProviderSet = wire.NewSet(
	NewAuthUsecase,
	NewGreeterUsecase,
	NewArticleUseCase,
	NewTodoUsecase,
	NewGameUsecase,
	NewStatsUsecase,
	// xsh 选品
	NewXshProductUsecase,
	// xsh 内容加工坊（首期使用 Noop 文案/封面/视频，后续可替换为 Gemini/OSS 等）
	NewNoopCopyGenerator,
	NewNoopCoverGenerator,
	NewNoopVideoProcessor,
	NewXshContentUsecase,
	// xsh 分发
	NewXshDispatchUsecase,
	// xsh 获客截流
	NewXshInboxUsecase,
)
