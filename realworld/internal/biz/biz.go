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
	// wechat 群发与登录
	NewWechatUsecase,
	// 萌宠之家（cat-admin）宠物
	NewPetUsecase,
	// 萌宠之家（cat-admin）签到
	NewCheckinUsecase,
	// 萌宠之家（cat-admin）用户资料
	NewUserUsecase,
	// 萌宠之家（cat-admin）商店与背包
	NewShopUsecase,
	NewInventoryUsecase,
	// 萌宠之家（cat-admin）品种
	NewBreedUsecase,
	// 萌宠之家（cat-admin）宠物装扮
	NewAppearanceUsecase,
	// 萌宠之家（cat-admin）成就与解锁器（签到/领养/购买触发）
	NewAchievementUsecase,
	NewAchievementUnlocker,
	// 萌宠之家（cat-admin）社交
	NewFriendUsecase,
)
