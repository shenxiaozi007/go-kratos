package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

// Todo 代表业务层中的待办事项实体。
// 这里不直接暴露底层数据库的实现细节（例如 GORM 的字段标签），
// 以便在将来需要更换存储介质时，只需要修改 data 层即可。
type Todo struct {
	// ID 是待办事项的唯一标识，由数据库自增生成。
	ID uint64
	// Value 是待办事项的具体内容。
	Value string
	// IsCompleted 表示该待办事项是否已经完成。
	IsCompleted bool
}

// TodoRepo 定义了 Todo 领域所需的数据访问接口。
// biz 层只依赖该接口，而不依赖具体的数据实现（如 MySQL/GORM），
// 这样可以保持领域逻辑的纯净。
type TodoRepo interface {
	// ListTodos 返回所有待办事项列表。
	ListTodos(ctx context.Context) ([]*Todo, error)
	// CreateTodo 创建一个新的待办事项，并返回包含自增 ID 的实体。
	CreateTodo(ctx context.Context, todo *Todo) (*Todo, error)
	// ToggleTodoCompleted 根据 ID 反转待办事项的完成状态，并返回更新后的实体。
	ToggleTodoCompleted(ctx context.Context, id uint64) (*Todo, error)
	// DeleteTodo 根据 ID 删除待办事项。
	DeleteTodo(ctx context.Context, id uint64) error
}

// TodoUsecase 封装了与 Todo 相关的业务用例。
// 它依赖 TodoRepo 接口而不是具体实现，并通过 log.Helper 记录关键日志。
type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

// NewTodoUsecase 创建一个新的 TodoUsecase 实例。
// 参数：
//   - repo: 具体的 TodoRepo 实现，由 data 层提供。
//   - logger: 全局日志记录器，用于构建 log.Helper。
func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// ListTodos 获取所有待办事项。
// 这里可以添加权限控制、审计日志等跨领域逻辑。
func (uc *TodoUsecase) ListTodos(ctx context.Context) ([]*Todo, error) {
	uc.log.Infof("list all todos")
	return uc.repo.ListTodos(ctx)
}

// CreateTodo 创建新的待办事项。
// 在进入数据层之前，这里会做基础的业务校验，例如内容不能为空。
func (uc *TodoUsecase) CreateTodo(ctx context.Context, value string, isCompleted bool) (*Todo, error) {
	uc.log.Infof("create todo, value=%s, isCompleted=%v", value, isCompleted)

	// 简单的业务校验：待办内容不能为空。
	if value == "" {
		uc.log.Warnf("create todo failed: empty value")
		return nil, ErrInvalidTodoContent
	}

	todo := &Todo{
		Value:       value,
		IsCompleted: isCompleted,
	}
	return uc.repo.CreateTodo(ctx, todo)
}

// ToggleTodoCompleted 根据 ID 反转待办事项的完成状态。
func (uc *TodoUsecase) ToggleTodoCompleted(ctx context.Context, id uint64) (*Todo, error) {
	uc.log.Infof("toggle todo status, id=%d", id)
	if id == 0 {
		uc.log.Warnf("toggle todo failed: invalid id=0")
		return nil, ErrInvalidTodoID
	}
	return uc.repo.ToggleTodoCompleted(ctx, id)
}

// DeleteTodo 根据 ID 删除待办事项。
func (uc *TodoUsecase) DeleteTodo(ctx context.Context, id uint64) error {
	uc.log.Infof("delete todo, id=%d", id)
	if id == 0 {
		uc.log.Warnf("delete todo failed: invalid id=0")
		return ErrInvalidTodoID
	}
	return uc.repo.DeleteTodo(ctx, id)
}

// 下面定义与 Todo 相关的业务错误。

// ErrInvalidTodoContent 表示创建或更新待办事项时内容不合法（例如为空）。
var ErrInvalidTodoContent = errors.New("待办内容不能为空")

// ErrInvalidTodoID 表示传入的待办事项 ID 不合法。
var ErrInvalidTodoID = errors.New("待办事项 ID 不合法")

