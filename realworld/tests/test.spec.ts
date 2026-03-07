import { test, expect, type Page } from '@playwright/test'

// 账号密码：与 tests/README.md 一致，可在此修改或通过环境变量覆盖
const USERNAME = process.env.E2E_USERNAME ?? 'hxc111'
const PASSWORD = process.env.E2E_PASSWORD ?? 'hxc111'

/** 登录后等待跳转到后台的超时（接口慢时可适当加大） */
const LOGIN_NAV_TIMEOUT = 15000
/** 页面内容可见性断言超时 */
const VISIBLE_TIMEOUT = 10000

/**
 * 执行登录并等待进入后台（商品列表页）。
 * 超时未跳转请检查：后端已启动、前端代理 /api、账号存在。
 */
async function doLogin(page: Page) {
  await page.getByPlaceholder('请输入用户名').fill(USERNAME)
  await page.getByPlaceholder('请输入密码').fill(PASSWORD)
  await page.getByRole('button', { name: '登录' }).click()
  await expect(page).toHaveURL(/\/xsh\/product\/list/, { timeout: LOGIN_NAV_TIMEOUT })
  await expect(page.getByRole('heading', { name: '商品列表' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
}

test.describe('推广引流管理后台 E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test.describe('未登录', () => {
    test('访问首页应重定向到登录页', async ({ page }) => {
      await expect(page).toHaveURL(/\/login/)
      await expect(page.getByText('推广引流管理后台')).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '登录' })).toBeVisible()
    })

    test('登录页空表单提交应显示校验', async ({ page }) => {
      await page.getByRole('button', { name: '登录' }).click()
      await expect(page.getByText('请输入用户名')).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByText('请输入密码')).toBeVisible()
    })

    test('使用正确账号密码可登录成功', async ({ page }) => {
      await doLogin(page)
    })

    test('可进入注册页并返回登录', async ({ page }) => {
      await page.getByRole('button', { name: '去注册' }).click()
      await expect(page).toHaveURL(/\/register/)
      await expect(page.getByPlaceholder('请再次输入密码')).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '注册并登录' })).toBeVisible()
      await page.getByRole('button', { name: '返回登录' }).click()
      await expect(page).toHaveURL(/\/login/)
    })
  })

  test.describe('已登录 - 选品采集仓', () => {
    test.beforeEach(async ({ page }) => {
      await doLogin(page)
    })

    test('商品列表页正常展示', async ({ page }) => {
      await expect(page.getByRole('heading', { name: '商品列表' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '查询' })).toBeVisible()
    })

    test('链接解析页正常展示', async ({ page }) => {
      await page.goto('/xsh/product/parse')
      await expect(page.getByText('粘贴链接解析')).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '解析' })).toBeVisible()
      await expect(page.getByPlaceholder('粘贴淘宝/猫超链接')).toBeVisible()
    })

    test('爆款榜单页正常展示', async ({ page }) => {
      await page.goto('/xsh/product/hot-ranks')
      await expect(page.getByRole('heading', { name: '爆款榜单' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })
  })

  test.describe('已登录 - AI 内容加工坊', () => {
    test.beforeEach(async ({ page }) => {
      await doLogin(page)
    })

    test('模板管理页正常展示', async ({ page }) => {
      await page.goto('/xsh/content/templates')
      await expect(page.getByRole('heading', { name: '模板管理' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '新增模板' })).toBeVisible()
    })

    test('草稿列表页正常展示', async ({ page }) => {
      await page.goto('/xsh/content/drafts')
      await expect(page.getByRole('heading', { name: '草稿列表' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })
  })

  test.describe('已登录 - 多账号分发', () => {
    test.beforeEach(async ({ page }) => {
      await doLogin(page)
    })

    test('账号绑定页正常展示', async ({ page }) => {
      await page.goto('/xsh/dispatch/accounts')
      await expect(page.getByRole('heading', { name: '账号绑定' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
      await expect(page.getByRole('button', { name: '查询' })).toBeVisible()
    })

    test('定时规则页正常展示', async ({ page }) => {
      await page.goto('/xsh/dispatch/schedules')
      await expect(page.getByRole('heading', { name: '定时规则' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })

    test('发布任务页正常展示', async ({ page }) => {
      await page.goto('/xsh/dispatch/tasks')
      await expect(page.getByRole('heading', { name: '发布任务' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })
  })

  test.describe('已登录 - 获客截流', () => {
    test.beforeEach(async ({ page }) => {
      await doLogin(page)
    })

    test('评论 Inbox 页正常展示', async ({ page }) => {
      await page.goto('/xsh/inbox/comments')
      await expect(page.getByRole('heading', { name: '评论 Inbox' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })

    test('动作记录页正常展示', async ({ page }) => {
      await page.goto('/xsh/inbox/actions')
      await expect(page.getByRole('heading', { name: '动作记录' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })
  })

  test.describe('侧边栏导航', () => {
    test.beforeEach(async ({ page }) => {
      await doLogin(page)
    })

    test('通过菜单可切换到各模块', async ({ page }) => {
      const sidebar = page.locator('.sidebar')

      await expect(page.getByRole('heading', { name: '商品列表' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })

      // 子菜单可能渲染在 .sidebar 内或 body 的 popper 里，用 page 查找更稳；展开后等子项可见再点
      const clickSubItem = async (submenuTitle: string, itemName: string) => {
        await sidebar.getByText(submenuTitle).first().click()
        const item = page.getByRole('link', { name: itemName }).or(page.getByText(itemName).first())
        await item.waitFor({ state: 'visible', timeout: 5000 })
        await item.click()
      }

      await clickSubItem('选品采集仓', '链接解析')
      await expect(page).toHaveURL(/\/xsh\/product\/parse/)
      await expect(page.getByText('粘贴链接解析')).toBeVisible({ timeout: VISIBLE_TIMEOUT })

      await clickSubItem('AI 内容加工坊', '模板管理')
      await expect(page).toHaveURL(/\/xsh\/content\/templates/)
      await expect(page.getByRole('heading', { name: '模板管理' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })

      await clickSubItem('多账号分发', '发布任务')
      await expect(page).toHaveURL(/\/xsh\/dispatch\/tasks/)
      await expect(page.getByRole('heading', { name: '发布任务' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })

      await clickSubItem('获客截流', '评论 Inbox')
      await expect(page).toHaveURL(/\/xsh\/inbox\/comments/)
      await expect(page.getByRole('heading', { name: '评论 Inbox' })).toBeVisible({ timeout: VISIBLE_TIMEOUT })
    })
  })
})
