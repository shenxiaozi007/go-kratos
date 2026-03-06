# Gemini 3 Flash 文案生成接入

本文档描述 v2 中将内容加工坊「一键生成文案」从 Noop 换为真实调用 Gemini 3 Flash（或当前可用等价模型）的设计，包括配置、实现、错误与限流及 Wire 注入。

---

## 一、目标

- 实现 [CopyGenerator](internal/biz/xsh_content.go) 接口的 Gemini 版本，供 `XshContentUsecase.GenerateCopy` 调用。
- 入参：`draftID`、`style`、`productTitle`、`templateContent`；出参：生成文案文本（string）。
- 在 Wire 中用该实现替换 `NoopCopyGenerator`，并注入 xsh.ai 配置。

---

## 二、配置

- **来源**：沿用 [xsh-docs/03-deployment](../xsh-docs/03-deployment.md) 的 xsh.ai 配置块。
- **字段**：
  - `xsh.ai.api_key`：Gemini API Key（必填，无 key 时不注册 Gemini 实现或回退 Noop）。
  - `xsh.ai.model`：模型 ID，如 `gemini-2.0-flash` 或官方当前等价型号；可选，缺省时使用 SDK/官方默认。
  - `xsh.ai.endpoint`（可选）：自定义 API 端点，若使用代理或自托管。
- **读取与注入**：在 `internal/conf` 中增加 xsh.ai 的配置结构体与解析；通过 Wire 将 conf 注入到 Gemini CopyGenerator 实现（及可选 Biz 层）。

---

## 三、实现位置与接口

- **接口**：`CopyGenerator` 定义在项目 [internal/biz/xsh_content.go](../../internal/biz/xsh_content.go)：
  - `GenerateCopy(ctx, draftID, style, productTitle, templateContent) (string, error)`
- **实现**：新增实现类，例如：
  - `internal/biz/xsh_gemini_copy.go`：实现 `CopyGenerator`，内部调用 Gemini API；或
  - `pkg/gemini/client.go`：封装 Gemini 调用，biz 层包装为 `CopyGenerator`。
- **请求格式**：按 Google Gemini API（REST 或 Go SDK）约定，构造 prompt：将 `style`、`productTitle`、`templateContent` 组合为系统/用户消息，请求生成一段小红书风格文案；只返回纯文本内容，便于直接写入草稿 `copy_text`。

---

## 四、错误与限流

- **API 失败**：网络错误、4xx/5xx、响应解析失败时返回明确 error，业务层可映射为统一错误码（如 503 或 500），前端提示「文案生成服务暂时不可用，请稍后重试」。
- **超时**：建议设置请求超时（如 30s），超时返回可区分提示「请求超时」。
- **限流（429）**：返回可区分错误或重试标识，前端提示「请求过于频繁，请稍后重试」；可选在实现中加入指数退避重试（如最多重试 2 次）。
- **熔断（可选）**：连续失败 N 次后短暂不再请求 Gemini，直接返回「服务暂不可用」；可在后续迭代实现。

---

## 五、Wire 与依赖

- 在 [cmd/realworld/wire_gen.go](../../cmd/realworld/wire_gen.go)（或 wire 输入）中：
  - 当 xsh.ai.api_key 已配置时，提供 Gemini 版 `CopyGenerator` 的构造函数，并注入 xsh.ai 配置；
  - 将该实现注入 `XshContentUsecase`，替代原来的 `NoopCopyGenerator`。
- 当未配置 api_key 时，可继续使用 NoopCopyGenerator，避免启动报错；业务层已能处理 `ErrCopyGeneratorNotAvailable`，前端可提示「未配置文案生成服务」。

---

## 六、与现有 Noop 的衔接

- 保留项目 [internal/biz/xsh_content_noop.go](../../internal/biz/xsh_content_noop.go) 中的 Noop 实现，用于未配置或降级场景。
- `XshContentUsecase` 仅依赖 `CopyGenerator` 接口，通过 Wire 在运行时注入 Gemini 实现或 Noop，无需改业务逻辑代码。
