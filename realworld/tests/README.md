# E2E 自动化测试 - 问题排查与修复

跑测试前请确保：
- 前端已启动：`cd front-xsh-admin && npm run serve`（默认 5174）
- 后端已启动且前端已配置代理 `/api` 到后端
- 存在测试账号：用户名 `hxc111`，密码 `hxc111`

**环境变量（可选）：**
- `PLAYWRIGHT_BASE_URL`：前端地址，默认 `http://192.168.31.115:5174`。例：`PLAYWRIGHT_BASE_URL=http://localhost:5174 npm run test:e2e`
- `E2E_USERNAME` / `E2E_PASSWORD`：登录账号，默认 `hxc111` / `hxc111`

---

## 一、运行与查看失败

```bash
# 项目根目录执行
npm run test:e2e

# 只看某个用例（用例名用正则）
npx playwright test -g "登录成功"

# 调试模式：有界面、慢速、保留 trace
npx playwright test --headed --slow-mo=500
```

失败时终端会标出**失败用例名**和**报错信息**，并生成 HTML 报告。

---

## 二、查看失败报告（必做）

测试失败后会提示报告路径，例如：

```bash
# 打开上次运行的 HTML 报告
npx playwright show-report
```

在报告里可以：
- 看**失败时的截图**
- 点开用例看**每一步的 trace**（需先开启 trace，见下文）
- 看报错堆栈和断言失败原因

---

## 三、常见问题与修复

### 1. 登录后没跳转 / 一直停在登录页

**可能原因：** 后端未启动、代理未配置、或账号不存在。

**排查：**
- 浏览器直接访问 `http://192.168.31.115:5174`，手动用 `hxc111` / `hxc111` 登录，看能否进「商品列表」。
- 打开开发者工具 Network，看 `POST .../api/v1/auth/login` 是否 200、响应里是否有 `access_token`。

**修复：**
- 启动后端服务。
- 在 `front-xsh-admin/vue.config.js`（或当前 dev 配置）里确认 `/api` 代理到后端地址。
- 若账号不对，在测试里改 `tests/test.spec.ts` 顶部的 `USERNAME` / `PASSWORD`，或在系统里创建对应用户。

---

### 2. 找不到元素（Timeout / locator not found）

**可能原因：** 文案或结构改了、元素在弹层/折叠里、页面未加载完。

**修复：**

- **改选择器**  
  打开 `tests/test.spec.ts`，找到报错里提到的用例，把对应的 `getByText(...)` / `getByRole(...)` / `getByPlaceholder(...)` 改成当前页面上真实存在的文案或角色。  
  优先用：`getByRole('button', { name: '登录' })`、`getByPlaceholder('请输入用户名')` 等语义化方式。

- **先等再点**  
  若元素是异步出来的，在操作前加等待，例如：
  ```ts
  await expect(page.getByRole('heading', { name: '商品列表' })).toBeVisible({ timeout: 10000 })
  ```

- **侧栏/子菜单**  
  若「链接解析」「模板管理」等在折叠菜单里，需要先点开父级菜单再点子项；若菜单结构改了，在用例里改成先点新的父级再点新的子项。

---

### 3. 超时（Test timeout / Expect timeout）

**可能原因：** 接口慢、页面加载慢、或等待时间设得太短。

**修复：**

- **单次等待加 timeout**  
  例如：`await expect(page).toHaveURL(/\/xsh\/product\/list/, { timeout: 20000 })`

- **全局加时**  
  在 `playwright.config.ts` 里调大：
  ```ts
  timeout: 60000,        // 单用例总超时
  expect: { timeout: 15000 }
  ```

- **先确认环境**  
  同一操作在浏览器里是否也要等很久；若是，先优化接口或前端加载，再适当加大超时。

---

### 4. 端口或地址不对

当前配置里前端地址在 `playwright.config.ts` 的 `baseURL`（默认 `http://192.168.31.115:5174`）。

**修复：**

- 用环境变量覆盖（推荐，无需改代码）：
  ```bash
  PLAYWRIGHT_BASE_URL=http://localhost:5174 npm run test:e2e
  ```
- 或直接改 `playwright.config.ts` 里 `baseURL`。测试脚本使用相对路径（如 `page.goto('/')`），会自动用该 baseURL。

---

### 5. 想看浏览器操作过程（调试用）

```bash
# 有界面运行，方便看点击/输入是否对
npx playwright test --headed

# 配合慢速
npx playwright test --headed --slow-mo=800

# UI 模式：选用例、看 trace、逐步执行
npm run test:e2e:ui
```

当前配置已开启失败时保留 trace（`trace: 'retain-on-failure'`）和首次重试录屏（`video: 'on-first-retry'`）。无需再改。

然后：

```bash
npx playwright show-report
# 在报告里点开失败用例 → 查看 Trace
```

---

## 四、修改测试代码的推荐方式

| 问题           | 改哪里                     | 示例 |
|----------------|----------------------------|------|
| 账号/密码      | 环境变量或 `tests/test.spec.ts` 顶部 | `E2E_USERNAME=xx E2E_PASSWORD=xx npm run test:e2e` 或改 `USERNAME`、`PASSWORD` |
| 前端地址/端口  | 环境变量或 `playwright.config.ts` | `PLAYWRIGHT_BASE_URL=http://localhost:5174` 或改 `baseURL` |
| 元素找不到     | 对应用例里的 locator       | 改 `getByText` / `getByRole` 等 |
| 等待不够       | `test.spec.ts` 里 `LOGIN_NAV_TIMEOUT` / `VISIBLE_TIMEOUT` 或单句 `{ timeout }` | 加大常量或 `expect(..., { timeout: 15000 })` |
| 全局超时       | `playwright.config.ts`     | 改 `timeout`、`expect.timeout` |
| 侧栏/菜单步骤  | `test.spec.ts` 侧边栏导航用例 | 先点父菜单（如「选品采集仓」）再点子项 |

改完后保存，再执行：

```bash
npm run test:e2e
```

确认通过后，可把本次修改（含 `tests/test.spec.ts`、`playwright.config.ts`、本 README）一并提交。
