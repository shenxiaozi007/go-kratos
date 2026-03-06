<template>
  <div class="dispatch-tasks">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="账号ID"><el-input v-model="filterAccountId" placeholder="account_id" clearable style="width:120px" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterStatus" placeholder="全部" clearable style="width:120px">
            <el-option label="pending" value="pending" />
            <el-option label="running" value="running" />
            <el-option label="success" value="success" />
            <el-option label="failed" value="failed" />
            <el-option label="cancelled" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="draft_id" label="草稿ID" width="90" />
        <el-table-column prop="account_id" label="账号ID" width="90" />
        <el-table-column prop="scheduled_at" label="计划时间" width="160">
          <template #default="{ row }">{{ formatTime(row.scheduled_at) }}</template>
        </el-table-column>
        <el-table-column prop="published_at" label="发布时间" width="160">
          <template #default="{ row }">{{ formatTime(row.published_at) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error_message" label="错误信息" min-width="140" show-overflow-tooltip />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'failed'" link type="primary" @click="onRetry(row)">重试</el-button>
            <el-button v-if="['pending', 'running'].includes(row.status)" link type="warning" @click="onCancel(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @current-change="load"
        @size-change="load"
      />
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPublishTasks, retryPublishTask, cancelPublishTask } from '@/api/dispatch'

function formatTime (ts) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString()
}
function statusType (s) {
  const map = { success: 'success', failed: 'danger', pending: 'info', running: 'warning', cancelled: 'info' }
  return map[s] || 'info'
}

export default {
  name: 'DispatchTasks',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterAccountId = ref('')
    const filterStatus = ref('')

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterAccountId.value) params.account_id = Number(filterAccountId.value)
        if (filterStatus.value) params.status = filterStatus.value
        const res = await listPublishTasks(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    async function onRetry (row) {
      try {
        await retryPublishTask(row.id)
        ElMessage.success('已提交重试')
        load()
      } catch (e) {
        ElMessage.error(e.message)
      }
    }

    function onCancel (row) {
      ElMessageBox.confirm('确定取消该任务？', '提示', { type: 'warning' }).then(async () => {
        try {
          await cancelPublishTask(row.id)
          ElMessage.success('已取消')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    onMounted(load)
    return { loading, items, total, page, pageSize, filterAccountId, filterStatus, load, onRetry, onCancel, formatTime, statusType }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; }
</style>
