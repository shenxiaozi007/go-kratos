package data

import (
	"context"

	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// TodoPO 是 Todo 在数据层（Persistence Object）中的结构定义。
// 它与数据库表结构一一对应，并带有 GORM 的字段标签，用于 ORM 映射。
type TodoPO struct {
	// ID 是主键，自增。
	ID uint64 `gorm:"primaryKey;autoIncrement"`
	// Value 是待办事项的具体内容，限制长度为 255 字符，不能为空。
	Value string `gorm:"type:varchar(255);not null"`
	// IsCompleted 表示是否完成，使用 tinyint(1) 存储布尔值。
	IsCompleted bool `gorm:"type:tinyint(1);not null;default:0"`
}

// TableName 指定该实体所对应的数据库表名为 "list"。
// 这样可以满足题目中“从数据库的 'list' 集合中查询”的要求。
func (TodoPO) TableName() string {
	return "list"
}

// todoRepo 是 TodoRepo 接口的具体实现，负责与数据库交互。
type todoRepo struct {
	// db 封装了底层的 *gorm.DB 对象。
	db *gorm.DB
	// log 用于记录数据访问层的关键日志。
	log *log.Helper
}

// NewTodoRepo 创建一个新的 Todo 仓储实现。
// 参数：
//   - data: Data 结构体，内部持有全局的 *gorm.DB。
//   - logger: 全局日志记录器，用于包装成 log.Helper。
func NewTodoRepo(data *Data, logger log.Logger) biz.TodoRepo {
	return &todoRepo{
		db:  data.DB,
		log: log.NewHelper(logger),
	}
}

// poToBiz 将数据层的 TodoPO 转换为业务层的 biz.Todo。
func poToBiz(po *TodoPO) *biz.Todo {
	if po == nil {
		return nil
	}
	return &biz.Todo{
		ID:          po.ID,
		Value:       po.Value,
		IsCompleted: po.IsCompleted,
	}
}

// bizToPO 将业务层的 biz.Todo 转换为数据层的 TodoPO。
func bizToPO(todo *biz.Todo) *TodoPO {
	if todo == nil {
		return nil
	}
	return &TodoPO{
		ID:          todo.ID,
		Value:       todo.Value,
		IsCompleted: todo.IsCompleted,
	}
}

// ListTodos 查询并返回所有待办事项。
func (r *todoRepo) ListTodos(ctx context.Context) ([]*biz.Todo, error) {
	r.log.Infof("query all todos from database")

	var pos []TodoPO
	// 使用 WithContext 将上层传入的 context 透传给 GORM，
	// 以便支持超时控制、链路追踪、取消等功能。
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		r.log.Errorf("query todos failed: %v", err)
		return nil, err
	}

	todos := make([]*biz.Todo, 0, len(pos))
	for _, po := range pos {
		todos = append(todos, poToBiz(&po))
	}
	return todos, nil
}

// CreateTodo 在数据库中插入一条新的待办事项记录。
func (r *todoRepo) CreateTodo(ctx context.Context, todo *biz.Todo) (*biz.Todo, error) {
	r.log.Infof("create todo in database, value=%s, isCompleted=%v", todo.Value, todo.IsCompleted)

	po := bizToPO(todo)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		r.log.Errorf("create todo failed: %v", err)
		return nil, err
	}

	// 数据库插入完成后，自增 ID 会回写到 po.ID 中，这里再转换回业务实体返回。
	return poToBiz(po), nil
}

// ToggleTodoCompleted 根据 ID 查询待办事项，并将其完成状态取反后写回数据库。
func (r *todoRepo) ToggleTodoCompleted(ctx context.Context, id uint64) (*biz.Todo, error) {
	r.log.Infof("toggle todo completed status in database, id=%d", id)

	var po TodoPO
	// 先根据主键查询待办事项，如果不存在则返回错误。
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Warnf("toggle todo failed: record not found, id=%d", id)
		} else {
			r.log.Errorf("toggle todo query failed: %v", err)
		}
		return nil, err
	}

	// 将完成状态取反。
	po.IsCompleted = !po.IsCompleted

	// 只更新 is_completed 字段，避免无意义的字段写入。
	if err := r.db.WithContext(ctx).Model(&TodoPO{}).
		Where("id = ?", po.ID).
		Update("is_completed", po.IsCompleted).Error; err != nil {
		r.log.Errorf("toggle todo update failed: %v", err)
		return nil, err
	}

	return poToBiz(&po), nil
}

// DeleteTodo 根据 ID 删除一条待办事项记录。
func (r *todoRepo) DeleteTodo(ctx context.Context, id uint64) error {
	r.log.Infof("delete todo from database, id=%d", id)

	if err := r.db.WithContext(ctx).Delete(&TodoPO{}, id).Error; err != nil {
		r.log.Errorf("delete todo failed: %v", err)
		return err
	}
	return nil
}
