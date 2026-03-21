# 萌宠之家（cat-admin）文档中心

萌宠之家 - 云养宠伴侣，基于 code.html 原型的全栈应用文档，位于 realworld 仓库内。

## 文档索引

| 文档 | 说明 |
|------|------|
| [01-项目文档.md](01-项目文档.md) | 项目概述、目标与范围、技术架构、路由与模块总览 |
| [02-后端开发文档.md](02-后端开发文档.md) | 后端 API 设计、鉴权、错误约定、Proto/Biz/Data/Service 实现要点 |
| [03-前端开发文档.md](03-前端开发文档.md) | 前端页面、组件、路由、设计规范、与后端 API 对接说明 |
| [database/](database/) | 数据库表结构、字段说明、迁移与首期/二期拆分 |

## 本地与依赖

- **后端**：realworld（Go + Kratos），`go mod tidy` 后 `go build ./cmd/realworld`。
- **前端**：front-pet-app（规划中），安装最新依赖请在该目录下执行 `npm install`。
- **文档**：本目录为纯 Markdown，无需安装；可用任意 Markdown 预览打开。

## 相关资源

- 原型参考：code.html（萌宠之家 - 云养宠伴侣）
- 后端仓库路径：realworld（与坦克/选品等模块并存）
- 前端规划路径：front-pet-app（与 front-vue、front-xsh-admin 并列）
