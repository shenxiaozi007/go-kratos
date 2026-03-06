<template>
  <div class="content-drafts">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="状态">
          <el-select v-model="filterStatus" placeholder="全部" clearable style="width:120px">
            <el-option label="draft" value="draft" />
            <el-option label="ready" value="ready" />
            <el-option label="published" value="published" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
        <el-form-item><el-button @click="$router.push('/xsh/content/draft/edit')">新建草稿</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_id" label="商品ID" width="90" />
        <el-table-column prop="template_id" label="模板ID" width="90" />
        <el-table-column prop="copy_text" label="文案" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/xsh/content/draft/edit/${row.id}`)">编辑</el-button>
            <el-button link type="primary" @click="goCreateTask(row)">创建发布任务</el-button>
            <el-button link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        class="pagination"
        @current-change="load"
        @size-change="load"
      />
    </el-card>
    <el-dialog v-model="taskDialogVisible" title="创建发布任务" width="480">
      <el-form label-width="100px">
        <el-form-item label="选择账号">
          <el-select v-model="taskAccountIds" multiple placeholder="多选账号" style="width:100%">
            <el-option v-for="a in accounts" :key="a.id" :label="a.nickname || a.platform + '-' + a.id" :value="a.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="计划时间">
          <el-date-picker v-model="taskScheduledAt" type="datetime" value-format="x" placeholder="选择时间" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taskDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="taskLoading" @click="submitTask">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listDrafts, deleteDraft } from '@/api/content'
import { listPlatformAccounts } from '@/api/dispatch'
import { createPublishTask } from '@/api/dispatch'

export default {
  name: 'ContentDrafts',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterStatus = ref('')
    const taskDialogVisible = ref(false)
    const taskDraftId = ref(null)
    const taskAccountIds = ref([])
    const taskScheduledAt = ref('')
    const accounts = ref([])
    const taskLoading = ref(false)

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterStatus.value) params.status = filterStatus.value
        const res = await listDrafts(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    async function goCreateTask (row) {
      if (row.status !== 'ready') {
        ElMessage.warning('请先将草稿状态改为 ready 再创建发布任务')
        return
      }
      taskDraftId.value = row.id
      taskAccountIds.value = []
      taskScheduledAt.value = Date.now() + 3600000
      try {
        const res = await listPlatformAccounts({ page_size: 100 })
        accounts.value = res.items || []
      } catch (e) {
        ElMessage.error(e.message)
        return
      }
      taskDialogVisible.value = true
    }

    async function submitTask () {
      if (!taskDraftId.value || !taskAccountIds.value.length) {
        ElMessage.warning('请选择至少一个账号')
        return
      }
      taskLoading.value = true
      try {
        await createPublishTask({
          draft_id: taskDraftId.value,
          account_ids: taskAccountIds.value,
          scheduled_at: Math.floor((taskScheduledAt.value || Date.now()) / 1000)
        })
        ElMessage.success('已创建发布任务')
        taskDialogVisible.value = false
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        taskLoading.value = false
      }
    }

    function onDelete (row) {
      ElMessageBox.confirm('确定删除该草稿？', '提示', { type: 'warning' }).then(async () => {
        try {
          await deleteDraft(row.id)
          ElMessage.success('已删除')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    onMounted(load)
    return {
      loading, items, total, page, pageSize, filterStatus,
      taskDialogVisible, taskAccountIds, taskScheduledAt, accounts, taskLoading,
      load, goCreateTask, submitTask, onDelete
    }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; }
</style>
