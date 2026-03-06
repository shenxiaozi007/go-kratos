<template>
  <el-container class="xsh-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">推广引流</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#1a1d24"
        text-color="#b0b8c4"
        active-text-color="#409eff"
      >
        <el-sub-menu index="product">
          <template #title>选品采集仓</template>
          <el-menu-item index="/xsh/product/parse">链接解析</el-menu-item>
          <el-menu-item index="/xsh/product/list">商品列表</el-menu-item>
          <el-menu-item index="/xsh/product/hot-ranks">爆款榜单</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="content">
          <template #title>AI 内容加工坊</template>
          <el-menu-item index="/xsh/content/templates">模板管理</el-menu-item>
          <el-menu-item index="/xsh/content/drafts">草稿列表</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="dispatch">
          <template #title>多账号分发</template>
          <el-menu-item v-if="isAdmin" index="/xsh/dispatch/accounts">账号绑定</el-menu-item>
          <el-menu-item index="/xsh/dispatch/schedules">定时规则</el-menu-item>
          <el-menu-item index="/xsh/dispatch/tasks">发布任务</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="inbox">
          <template #title>获客截流</template>
          <el-menu-item index="/xsh/inbox/comments">评论 Inbox</el-menu-item>
          <el-menu-item index="/xsh/inbox/actions">动作记录</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-main class="main">
      <h2 v-if="$route.meta.title" class="page-title">{{ $route.meta.title }}</h2>
      <router-view />
    </el-main>
  </el-container>
</template>

<script>
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProfile } from '../api/auth'

export default {
  name: 'XshLayout',
  setup () {
    const route = useRoute()
    const activeMenu = computed(() => route.path)
    const isAdmin = ref(true)
    onMounted(async () => {
      try {
        const p = await getProfile()
        isAdmin.value = p.role === 'admin'
      } catch (_) {
        isAdmin.value = false
      }
    })
    return { activeMenu, isAdmin }
  }
}
</script>

<style scoped>
.xsh-layout { height: 100vh; }
.sidebar { background: #1a1d24; }
.logo {
  height: 56px;
  line-height: 56px;
  text-align: center;
  color: #e4e7ed;
  font-weight: 600;
  border-bottom: 1px solid #2c2f36;
}
.main { background: #f0f2f5; padding: 16px 24px; overflow: auto; }
.page-title { margin: 0 0 16px; font-size: 18px; color: #303133; }
</style>
