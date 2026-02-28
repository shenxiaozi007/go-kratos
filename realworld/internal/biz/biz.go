package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
// 虽然当前项目的依赖注入是通过 wire_gen.go 手动生成好的，
// 但仍然在 ProviderSet 中补充 TodoUsecase，方便后续如果重新生成 wire 代码时自动加入。
var ProviderSet = wire.NewSet(
	NewGreeterUsecase,
	NewArticleUseCase,
	NewTodoUsecase,
)
