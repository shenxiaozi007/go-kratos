# 可访问性说明（Accessibility）

## 图片与 alt 文本

- **当前状态**：项目中未使用 `<img>` 等需要 alt 的图片元素，无缺失 alt 问题。
- **若后续添加图片**：请为所有 `<img>` 提供有意义的 `alt` 属性；纯装饰性图片使用 `alt=""` 或 `role="presentation"` 并配合 `aria-hidden="true"`，避免读屏朗读。

## 已做的可访问性改进

- **颜色对比度**：提高副标题、提示文字在深色背景上的对比度（如 `#94a3b8` → `#b8c5d6`，`#6b7280` → `#4b5563`）。
- **语义化 HTML**：使用 `<main>`、`<header>`、`<nav>`、`<section>`、`<aside>`；品牌标题改为 `<p>`，每页保留单一 `<h1>`；表格表头使用 `scope="col"`。
- **ARIA**：主导航 `aria-label="主导航"`；加载区域 `role="status"`、`aria-live="polite"`、`aria-busy`；错误信息 `role="alert"`；按钮组 `role="group"` 与 `aria-label`；切换按钮 `aria-pressed`；游戏画布 `role="img"` 与 `aria-label`；分页使用 `<nav aria-label="历史战绩分页">`。
- **键盘导航**：全局「跳到主内容」链接（Tab 首项聚焦时显示）；所有按钮与链接具备 `:focus-visible` 轮廓（2px 实线 + offset），便于键盘用户识别焦点。

## 键盘操作建议

- 使用 **Tab** / **Shift+Tab** 在链接与按钮间移动。
- 游戏页：画布为装饰性内容，焦点在「返回首页」等控件上；游戏内操作依赖 WASD / 空格 / E（由 GameLoop 处理）。
