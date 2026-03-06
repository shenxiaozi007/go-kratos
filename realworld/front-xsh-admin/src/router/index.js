/**
 * 推广引流管理后台路由（与 xsh-docs/frontend/views-routes 一致）
 * v2：/login 为登录页；/xsh 下所有路由需登录后访问
 */
import { createRouter, createWebHistory } from 'vue-router'
import XshLayout from '../views/XshLayout.vue'

const TOKEN_KEY = 'xsh_access_token'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue'), meta: { title: '登录', public: true } },
  {
    path: '/xsh',
    component: XshLayout,
    redirect: '/xsh/product/list',
    meta: { requiresAuth: true },
    children: [
      // 选品采集仓
      { path: 'product/parse', name: 'ProductParse', component: () => import('../views/product/ProductParse.vue'), meta: { title: '链接解析' } },
      { path: 'product/list', name: 'ProductList', component: () => import('../views/product/ProductList.vue'), meta: { title: '商品列表' } },
      { path: 'product/hot-ranks', name: 'ProductHotRanks', component: () => import('../views/product/ProductHotRanks.vue'), meta: { title: '爆款榜单' } },
      // AI 内容加工坊
      { path: 'content/templates', name: 'ContentTemplates', component: () => import('../views/content/ContentTemplates.vue'), meta: { title: '模板管理' } },
      { path: 'content/drafts', name: 'ContentDrafts', component: () => import('../views/content/ContentDrafts.vue'), meta: { title: '草稿列表' } },
      { path: 'content/draft/edit/:id?', name: 'ContentDraftEdit', component: () => import('../views/content/ContentDraftEdit.vue'), meta: { title: '草稿编辑' } },
      // 多账号分发
      { path: 'dispatch/accounts', name: 'DispatchAccounts', component: () => import('../views/dispatch/DispatchAccounts.vue'), meta: { title: '账号绑定' } },
      { path: 'dispatch/schedules', name: 'DispatchSchedules', component: () => import('../views/dispatch/DispatchSchedules.vue'), meta: { title: '定时规则' } },
      { path: 'dispatch/tasks', name: 'DispatchTasks', component: () => import('../views/dispatch/DispatchTasks.vue'), meta: { title: '发布任务' } },
      // 获客截流
      { path: 'inbox/comments', name: 'InboxComments', component: () => import('../views/inbox/InboxComments.vue'), meta: { title: '评论 Inbox' } },
      { path: 'inbox/actions', name: 'InboxActions', component: () => import('../views/inbox/InboxActions.vue'), meta: { title: '动作记录' } }
    ]
  },
  { path: '/', redirect: '/xsh/product/list' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (to.path === '/login') {
    if (token) return next('/xsh/product/list')
    return next()
  }
  if (to.path.startsWith('/xsh') && !token) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }
  next()
})

export default router
