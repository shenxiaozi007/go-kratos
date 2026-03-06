# 关键组件与复用约定

以下为各模块可复用的组件与约定，便于与 [views-routes.md](views-routes.md) 配合实现。

---

## 通用组件

| 组件 | 用途 | 说明 |
|------|------|------|
| XshTable | 通用表格 | 支持分页、排序、筛选列、加载态；接收 columns、data、total、page、pageSize、@change 等 |
| XshFilterBar | 筛选栏 | 表单项（关键词、下拉、日期范围等）+ 查询/重置按钮 |
| XshEmpty | 空状态 | 无数据时的占位与操作引导 |
| XshConfirm | 二次确认 | 删除、取消任务等操作的弹窗确认 |

---

## 选品采集仓

| 组件 | 用途 |
|------|------|
| ProductLinkInput | 链接输入框 + 解析按钮，展示解析结果卡片（标题、原价、券额等），支持编辑后入库 |
| ProductTable | 商品列表表格，神价标签、来源标签、操作列（编辑、删除、去创建草稿） |
| HotRankPanel | 爆款榜单区块：类型切换、同步按钮、列表（首期可为占位） |

---

## AI 内容加工坊

| 组件 | 用途 |
|------|------|
| TemplateSelect | 模板选择器（下拉或卡片），展示风格类型与名称 |
| DraftForm | 草稿表单：商品选择器、模板选择、文案输入、封面图预览、视频上传/去重入口 |
| CopyPreview | 文案预览（可带占位符渲染结果） |
| CoverPreview | 封面图预览，支持一键生成后展示 |

---

## 多账号分发矩阵

| 组件 | 用途 |
|------|------|
| AccountForm | 账号表单：平台、昵称、proxy_url、cookie_storage_path 等；代理可脱敏展示 |
| ScheduleForm | 定时规则表单：名称、cron 或时间点配置、关联账号多选 |
| TaskTable | 发布任务表格：草稿摘要、账号、计划时间、状态、操作（重试、取消） |
| TimeRangePicker | 定时发布用时间选择（如每天 10:00、12:00、20:00 多选） |

---

## 获客截流

| 组件 | 用途 |
|------|------|
| CommentList | 评论列表：账号 Tab、求链接高亮、内容、时间、操作（回复、私信、标记已处理） |
| ReplyModal | 回复/私信弹窗：输入内容、发送按钮 |
| ActionRecordTable | 动作记录表格：评论摘要、动作类型、发送内容、结果、时间 |

---

## 与现有 front-vue 的复用

- 若与 [front-vue](../../front-vue/) 同仓库：可抽离公共 UI 组件（如按钮、输入框、表格）到 `shared/` 或通过 npm 包复用；业务组件（如 ProductLinkInput、DraftForm）放在 front-xsh-admin 内，不依赖游戏相关逻辑。
- 若完全独立 front-xsh-admin：不依赖 front-vue，仅参考其 API 封装与路由风格即可。
