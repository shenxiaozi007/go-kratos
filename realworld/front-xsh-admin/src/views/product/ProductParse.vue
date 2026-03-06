<template>
  <div class="product-parse">
    <el-card>
      <template #header>粘贴链接解析（首期可人工补全后入库）</template>
      <el-input
        v-model="url"
        type="textarea"
        :rows="2"
        placeholder="粘贴淘宝/猫超链接"
      />
      <div style="margin-top:12px">
        <el-button type="primary" :loading="loading" @click="onParse">解析</el-button>
        <el-button v-if="product" type="success" :loading="saveLoading" @click="onSave">一键入库</el-button>
      </div>
    </el-card>
    <el-card v-if="product || parseError" class="result-card">
      <template #header>解析结果</template>
      <el-alert v-if="parseError" type="warning" :title="parseError" show-icon />
      <el-descriptions v-else-if="product" :column="1" border>
        <el-descriptions-item label="标题">{{ product.title || '-' }}</el-descriptions-item>
        <el-descriptions-item label="原价">{{ product.original_price }}</el-descriptions-item>
        <el-descriptions-item label="券额">{{ product.coupon_amount }}</el-descriptions-item>
        <el-descriptions-item label="佣金比例">{{ product.commission_rate }}%</el-descriptions-item>
        <el-descriptions-item label="主图"><img v-if="product.image_url" :src="product.image_url" style="max-width:120px;max-height:80px" /></el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { parseProductLink, createProduct } from '@/api/product'

export default {
  name: 'ProductParse',
  setup () {
    const router = useRouter()
    const url = ref('')
    const loading = ref(false)
    const saveLoading = ref(false)
    const product = ref(null)
    const parseError = ref('')

    async function onParse () {
      if (!url.value.trim()) {
        ElMessage.warning('请输入链接')
        return
      }
      loading.value = true
      product.value = null
      parseError.value = ''
      try {
        const res = await parseProductLink(url.value.trim())
        if (res.product) {
          product.value = res.product
        } else {
          parseError.value = res.error_message || '未解析到商品，可手动创建'
        }
      } catch (e) {
        parseError.value = e.message
      } finally {
        loading.value = false
      }
    }

    async function onSave () {
      if (!product.value) return
      saveLoading.value = true
      try {
        await createProduct({
          source_url: product.value.source_url || url.value,
          title: product.value.title || '未命名',
          original_price: product.value.original_price || 0,
          coupon_amount: product.value.coupon_amount || 0,
          commission_rate: product.value.commission_rate || 0,
          image_url: product.value.image_url || '',
          is_hot: false,
          rank_source: ''
        })
        ElMessage.success('已入库')
        router.push('/xsh/product/list')
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        saveLoading.value = false
      }
    }

    return { url, loading, saveLoading, product, parseError, onParse, onSave }
  }
}
</script>

<style scoped>
.product-parse { max-width: 720px; }
.result-card { margin-top: 16px; }
</style>
