<template>
  <div class="dispatch-accounts">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="平台">
          <el-select v-model="filterPlatform" placeholder="全部" clearable style="width:120px">
            <el-option label="抖音" value="douyin" />
            <el-option label="小红书" value="xiaohongshu" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
        <el-form-item><el-button type="success" @click="openAdd">绑定账号</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="platform" label="平台" width="100" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="proxy_url" label="代理" width="140">
          <template #default="{ row }">{{ maskProxy(row.proxy_url) }}</template>
        </el-table-column>
        <el-table-column prop="bind_at" label="绑定时间" width="160">
          <template #default="{ row }">{{ formatTime(row.bind_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="onUnbind(row)">解绑</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        class="pagination"
        @current-change="load"
      />
    </el-card>
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑账号' : '绑定账号'" width="500">
      <el-form :model="form" label-width="120px">
        <el-form-item label="平台" required>
          <el-select v-model="form.platform" placeholder="选择平台" style="width:100%" :disabled="isEdit">
            <el-option label="抖音" value="douyin" />
            <el-option label="小红书" value="xiaohongshu" />
          </el-select>
        </el-form-item>
        <el-form-item label="昵称"><el-input v-model="form.nickname" placeholder="备注名" /></el-form-item>
        <el-form-item label="代理 URL"><el-input v-model="form.proxy_url" placeholder="http://user:pass@host:port" /></el-form-item>
        <el-form-item label="Cookie 存储路径"><el-input v-model="form.cookie_storage_path" placeholder="登录态路径" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saveLoading" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPlatformAccounts, bindAccount, updateAccount, unbindAccount } from '@/api/dispatch'

function formatTime (ts) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString()
}
function maskProxy (s) {
  if (!s) return '-'
  if (s.length <= 20) return s.replace(/:.*@/, ':***@')
  return s.slice(0, 12) + '***' + s.slice(-8)
}

export default {
  name: 'DispatchAccounts',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterPlatform = ref('')
    const dialogVisible = ref(false)
    const isEdit = ref(false)
    const saveLoading = ref(false)
    const form = reactive({ id: null, platform: '', nickname: '', proxy_url: '', cookie_storage_path: '' })

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterPlatform.value) params.platform = filterPlatform.value
        const res = await listPlatformAccounts(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    function openAdd () {
      isEdit.value = false
      form.id = null
      form.platform = ''
      form.nickname = ''
      form.proxy_url = ''
      form.cookie_storage_path = ''
      dialogVisible.value = true
    }

    function openEdit (row) {
      isEdit.value = true
      form.id = row.id
      form.platform = row.platform
      form.nickname = row.nickname
      form.proxy_url = row.proxy_url
      form.cookie_storage_path = row.cookie_storage_path || ''
      dialogVisible.value = true
    }

    async function save () {
      if (!form.platform) {
        ElMessage.warning('请选择平台')
        return
      }
      saveLoading.value = true
      try {
        if (isEdit.value) {
          await updateAccount(form.id, { nickname: form.nickname, proxy_url: form.proxy_url, cookie_storage_path: form.cookie_storage_path })
        } else {
          await bindAccount({ platform: form.platform, nickname: form.nickname, proxy_url: form.proxy_url, cookie_storage_path: form.cookie_storage_path })
        }
        ElMessage.success('已保存')
        dialogVisible.value = false
        load()
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        saveLoading.value = false
      }
    }

    function onUnbind (row) {
      ElMessageBox.confirm('确定解绑该账号？', '提示', { type: 'warning' }).then(async () => {
        try {
          await unbindAccount(row.id)
          ElMessage.success('已解绑')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    onMounted(load)
    return { loading, items, total, page, pageSize, filterPlatform, dialogVisible, isEdit, form, saveLoading, load, openAdd, openEdit, save, onUnbind, formatTime, maskProxy }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; }
</style>
