<template>
  <div class="inbox-actions">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="评论ID"><el-input v-model="filterCommentId" placeholder="comment_id" clearable style="width:140px" /></el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="comment_id" label="评论ID" width="90" />
        <el-table-column prop="action_type" label="类型" width="120" />
        <el-table-column prop="content" label="内容" min-width="180" show-overflow-tooltip />
        <el-table-column prop="sent_at" label="发送时间" width="160">
          <template #default="{ row }">{{ formatTime(row.sent_at) }}</template>
        </el-table-column>
        <el-table-column prop="result_status" label="结果" width="90" />
        <el-table-column prop="error_message" label="错误" min-width="140" show-overflow-tooltip />
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
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listInboxActions } from '@/api/inbox'

function formatTime (ts) {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString()
}

export default {
  name: 'InboxActions',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterCommentId = ref('')

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterCommentId.value) params.comment_id = Number(filterCommentId.value)
        const res = await listInboxActions(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    onMounted(load)
    return { loading, items, total, page, pageSize, filterCommentId, load, formatTime }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; }
</style>
