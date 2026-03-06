<template>
  <div class="inbox-comments">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="账号ID"><el-input v-model="filterAccountId" placeholder="account_id" clearable style="width:120px" /></el-form-item>
        <el-form-item label="求链接">
          <el-select v-model="filterWantLink" placeholder="全部" clearable style="width:100px">
            <el-option label="是" :value="true" />
            <el-option label="否" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filterStatus" placeholder="全部" clearable style="width:100px">
            <el-option label="pending" value="pending" />
            <el-option label="replied" value="replied" />
            <el-option label="ignored" value="ignored" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
        <el-form-item><el-button :loading="syncLoading" @click="onSync">同步评论</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="account_id" label="账号ID" width="90" />
        <el-table-column prop="author_nickname" label="作者" width="100" />
        <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span :class="{ 'want-link': row.want_link }">{{ row.content }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="want_link" label="求链接" width="80">
          <template #default="{ row }"><el-tag v-if="row.want_link" type="success" size="small">是</el-tag><span v-else>-</span></template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openReply(row)">回复</el-button>
            <el-button link type="primary" @click="openPrivate(row)">私信</el-button>
            <el-button link @click="markHandled(row, 'replied')">标已回复</el-button>
            <el-button link @click="markHandled(row, 'ignored')">标忽略</el-button>
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
    <el-dialog v-model="replyVisible" :title="replyTitle" width="480">
      <el-input v-model="replyContent" type="textarea" :rows="3" placeholder="输入回复/私信内容" />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" :loading="replyLoading" @click="submitReply">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listComments, markCommentHandled, syncComments, sendReply, sendPrivateMessage } from '@/api/inbox'

export default {
  name: 'InboxComments',
  setup () {
    const loading = ref(false)
    const syncLoading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterAccountId = ref('')
    const filterWantLink = ref(undefined)
    const filterStatus = ref('')
    const replyVisible = ref(false)
    const replyTitle = ref('回复')
    const replyContent = ref('')
    const replyLoading = ref(false)
    const currentCommentId = ref(null)
    const isPrivate = ref(false)

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterAccountId.value) params.account_id = Number(filterAccountId.value)
        if (filterWantLink.value !== undefined) params.want_link = filterWantLink.value
        if (filterStatus.value) params.status = filterStatus.value
        const res = await listComments(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    async function onSync () {
      syncLoading.value = true
      try {
        const res = await syncComments(filterAccountId.value ? Number(filterAccountId.value) : 0)
        ElMessage.success(res.message || '已触发同步')
        load()
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        syncLoading.value = false
      }
    }

    function openReply (row) {
      currentCommentId.value = row.id
      isPrivate.value = false
      replyTitle.value = '回复评论'
      replyContent.value = ''
      replyVisible.value = true
    }

    function openPrivate (row) {
      currentCommentId.value = row.id
      isPrivate.value = true
      replyTitle.value = '私信引导'
      replyContent.value = ''
      replyVisible.value = true
    }

    async function submitReply () {
      if (!replyContent.value.trim()) {
        ElMessage.warning('请输入内容')
        return
      }
      replyLoading.value = true
      try {
        if (isPrivate.value) {
          await sendPrivateMessage(currentCommentId.value, replyContent.value)
        } else {
          await sendReply(currentCommentId.value, replyContent.value)
        }
        ElMessage.success('已发送')
        replyVisible.value = false
        load()
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        replyLoading.value = false
      }
    }

    async function markHandled (row, status) {
      try {
        await markCommentHandled(row.id, status)
        ElMessage.success('已标记')
        load()
      } catch (e) {
        ElMessage.error(e.message)
      }
    }

    onMounted(load)
    return {
      loading, syncLoading, items, total, page, pageSize,
      filterAccountId, filterWantLink, filterStatus,
      replyVisible, replyTitle, replyContent, replyLoading,
      load, onSync, openReply, openPrivate, submitReply, markHandled
    }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; }
.want-link { color: #e6a23c; font-weight: 500; }
</style>
