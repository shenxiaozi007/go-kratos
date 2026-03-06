<template>
  <div class="product-list">
    <el-card>
      <el-form :inline="true" class="filter-bar">
        <el-form-item label="神价">
          <el-select v-model="filterIsHot" placeholder="全部" clearable style="width:100px">
            <el-option label="是" :value="true" />
            <el-option label="否" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源">
          <el-input v-model="filterRankSource" placeholder="rank_source" clearable style="width:140px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load">查询</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="items" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="original_price" label="原价" width="90" />
        <el-table-column prop="coupon_amount" label="券额" width="80" />
        <el-table-column prop="is_hot" label="神价" width="70">
          <template #default="{ row }">
            <el-tag v-if="row.is_hot" type="danger" size="small">神价</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="rank_source" label="来源" width="100" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goEdit(row.id)">编辑</el-button>
            <el-button link type="primary" @click="goCreateDraft(row.id)">创建草稿</el-button>
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
    <el-dialog v-model="editVisible" title="编辑商品" width="500">
      <el-form v-if="editForm" :model="editForm" label-width="100px">
        <el-form-item label="标题"><el-input v-model="editForm.title" /></el-form-item>
        <el-form-item label="原价"><el-input-number v-model="editForm.original_price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="券额"><el-input-number v-model="editForm.coupon_amount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="神价"><el-switch v-model="editForm.is_hot" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saveLoading" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listProducts, updateProduct, deleteProduct } from '@/api/product'

export default {
  name: 'ProductList',
  setup () {
    const router = useRouter()
    const loading = ref(false)
    const items = ref([])
    const total = ref(0)
    const page = ref(1)
    const pageSize = ref(20)
    const filterIsHot = ref(undefined)
    const filterRankSource = ref('')
    const editVisible = ref(false)
    const editForm = ref(null)
    const saveLoading = ref(false)

    async function load () {
      loading.value = true
      try {
        const params = { page: page.value, page_size: pageSize.value }
        if (filterIsHot.value !== undefined && filterIsHot.value !== '') params.is_hot = filterIsHot.value
        if (filterRankSource.value) params.rank_source = filterRankSource.value
        const res = await listProducts(params)
        items.value = res.items || []
        total.value = res.total || 0
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    function goEdit (id) {
      const row = items.value.find(r => r.id === id)
      if (!row) return
      editForm.value = {
        id: row.id,
        title: row.title,
        original_price: row.original_price,
        coupon_amount: row.coupon_amount,
        is_hot: row.is_hot
      }
      editVisible.value = true
    }

    async function saveEdit () {
      if (!editForm.value) return
      saveLoading.value = true
      try {
        await updateProduct(editForm.value.id, editForm.value)
        ElMessage.success('已更新')
        editVisible.value = false
        load()
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        saveLoading.value = false
      }
    }

    function onDelete (row) {
      ElMessageBox.confirm('确定删除该商品？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await deleteProduct(row.id)
          ElMessage.success('已删除')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    function goCreateDraft (productId) {
      router.push({ path: '/xsh/content/draft/edit', query: { product_id: productId } })
    }

    onMounted(load)
    return {
      loading, items, total, page, pageSize, filterIsHot, filterRankSource,
      editVisible, editForm, saveLoading, load, goEdit, saveEdit, onDelete, goCreateDraft
    }
  }
}
</script>

<style scoped>
.filter-bar { margin-bottom: 12px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>
