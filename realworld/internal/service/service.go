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
	NewWechatService,
	// 萌宠之家（cat-admin）宠物
	NewPetService,
	// 萌宠之家（cat-admin）签到
	NewCheckinService,
	// 萌宠之家（cat-admin）用户资料
	NewUserService,
	// 萌宠之家（cat-admin）商店与背包
	NewShopService,
	NewInventoryService,
	// 萌宠之家（cat-admin）品种
	NewBreedService,
	// 萌宠之家（cat-admin）成就
	NewAchievementService,
	// 萌宠之家（cat-admin）社交
	NewFriendService,
)
