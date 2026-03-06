<template>
  <div class="hot-ranks">
    <el-card>
      <template #header>
        <span>爆款榜单</span>
        <el-button type="primary" size="small" :loading="syncLoading" style="float:right" @click="onSync">同步</el-button>
      </template>
      <el-form :inline="true">
        <el-form-item label="榜单类型">
          <el-input v-model="rankType" placeholder="如 high_commission" clearable style="width:160px" />
        </el-form-item>
        <el-form-item><el-button type="primary" @click="load">查询</el-button></el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="rank_type" label="类型" width="140" />
        <el-table-column prop="synced_at" label="同步时间" width="180">
          <template #default="{ row }">{{ formatTime(row.synced_at) }}</template>
        </el-table-column>
        <el-table-column prop="snapshot_json" label="快照" min-width="200" show-overflow-tooltip />
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
import { listHotRanks, syncHotRank } from '@/api/product'

function formatTime (ts) {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  return d.toLocaleString()
}

export default {
  name: 'ProductHotRanks',
  setup () {
    const loading = ref(false)
    const syncLoading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const rankType = ref('')

    async function load () {
      loading.value = true
      try {
        const res = await listHotRanks({ rank_type: rankType.value, page: page.value, page_size: pageSize.value })
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
        const res = await syncHotRank(rankType.value || 'default')
        ElMessage.success(res.message || '已同步')
        load()
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        syncLoading.value = false
      }
    }

    onMounted(load)
    return { loading, syncLoading, items, total, page, pageSize, rankType, load, onSync, formatTime }
  }
}
</script>

<style scoped>
.pagination { margin-top: 16px; }
</style>
