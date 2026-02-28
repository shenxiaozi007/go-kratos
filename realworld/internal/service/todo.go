package service

import (
	"context"

	todov1 "realworld/api/todo/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TodoService 负责实现由 proto（api/todo/v1/todo.proto）定义的 Todo gRPC/HTTP 服务。
// 通过实现 todov1.TodoServer 接口，并在 HTTP Server 中通过 RegisterTodoHTTPServer
// 注册，即可完全按照 Kratos 官方推荐的 proto 驱动方式来提供 RESTful API。
type TodoService struct {
	// 为了兼容后续可能新增的 RPC 方法，嵌入 UnimplementedTodoServer。
	todov1.UnimplementedTodoServiceServer

	// uc 是与 Todo 相关的业务用例。
	uc *biz.TodoUsecase
	// log 用于记录 Service 层的关键日志。
	log *log.Helper
}

// NewTodoService 创建一个新的 TodoService 实例。
// 参数：
//   - uc: TodoUsecase，封装了核心业务逻辑。
//   - logger: 全局日志记录器，用于包装成 log.Helper。
func NewTodoService(uc *biz.TodoUsecase, logger log.Logger) *TodoService {
	return &TodoService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// GetTodoList 实现 proto 中的 GetTodoList RPC。
// 通过 google.api.http 注解，该方法对应的 HTTP 接口为：
//
//	GET /api/get-todo
//
// 请求参数：无（使用 google.protobuf.Empty）
// 返回结果：GetTodoListReply，JSON 结构为：
//
//	{
//	  "items": [
//	    { "id": 1, "value": "示例", "isCompleted": false },
//	    ...
//	  ]
//	}
func (s *TodoService) GetTodoList(ctx context.Context, _ *emptypb.Empty) (*todov1.GetTodoListReply, error) {
	// 调用业务用例获取所有待办事项。
	todos, err := s.uc.ListTodos(ctx)
	if err != nil {
		s.log.Errorf("list todos failed: %v", err)
		return nil, err
	}

	// 将业务实体转换为 proto 定义的 Todo。
	reply := &todov1.GetTodoListReply{
		Items: make([]*todov1.Todo, 0, len(todos)),
	}
	for _, t := range todos {
		reply.Items = append(reply.Items, &todov1.Todo{
			Id:          t.ID,
			Value:       t.Value,
			IsCompleted: t.IsCompleted,
		})
	}
	return reply, nil
}

// AddTodo 实现 proto 中的 AddTodo RPC。
// 通过 google.api.http 注解，该方法对应的 HTTP 接口为：
//
//	POST /api/add-todo
//	Body:
//	{
//	  "value": "待办事项内容",
//	  "isCompleted": false
//	}
//
// 返回结果：TodoReply，JSON 结构为：
//
//	{
//	  "todo": { "id": 1, "value": "待办事项内容", "isCompleted": false }
//	}
func (s *TodoService) AddTodo(ctx context.Context, req *todov1.AddTodoRequest) (*todov1.TodoReply, error) {
	// 调用业务用例创建待办事项。
	todo, err := s.uc.CreateTodo(ctx, req.GetValue(), req.GetIsCompleted())
	if err != nil {
		s.log.Errorf("create todo failed: %v", err)
		return nil, err
	}

	return &todov1.TodoReply{
		Todo: &todov1.Todo{
			Id:          todo.ID,
			Value:       todo.Value,
			IsCompleted: todo.IsCompleted,
		},
	}, nil
}

// UpdateTodo 实现 proto 中的 UpdateTodo RPC。
// 通过 google.api.http 注解，该方法对应的 HTTP 接口为：
//
//	POST /api/update-todo/{id}
//
// 路径参数：
//   - id: 待办事项的唯一标识，例如 /api/update-todo/1
//
// 返回结果：TodoReply，包含更新后的待办事项。
func (s *TodoService) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.TodoReply, error) {
	id := req.GetId()

	// 调用业务用例反转完成状态。
	todo, err := s.uc.ToggleTodoCompleted(ctx, id)
	if err != nil {
		s.log.Errorf("toggle todo failed: %v", err)
		return nil, err
	}

	return &todov1.TodoReply{
		Todo: &todov1.Todo{
			Id:          todo.ID,
			Value:       todo.Value,
			IsCompleted: todo.IsCompleted,
		},
	}, nil
}

// DeleteTodo 实现 proto 中的 DeleteTodo RPC。
// 通过 google.api.http 注解，该方法对应的 HTTP 接口为：
//
//	POST /api/del-todo/{id}
//
// 路径参数：
//   - id: 待办事项的唯一标识，例如 /api/del-todo/1
//
// 返回结果：DeleteTodoReply，JSON 结构为：
//
//	{
//	  "success": true
//	}
func (s *TodoService) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoReply, error) {
	id := req.GetId()

	// 调用业务用例执行删除操作。
	if err := s.uc.DeleteTodo(ctx, id); err != nil {
		s.log.Errorf("delete todo failed: %v", err)
		return nil, err
	}

	return &todov1.DeleteTodoReply{
		Success: true,
	}, nil
}
