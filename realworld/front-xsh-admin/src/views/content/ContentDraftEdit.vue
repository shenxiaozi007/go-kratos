<template>
  <div class="draft-edit">
    <el-card v-loading="loading">
      <el-form label-width="100px">
        <el-form-item label="关联商品">
          <el-select v-model="form.product_id" placeholder="选择商品" filterable style="width:100%" @focus="loadProducts">
            <el-option v-for="p in products" :key="p.id" :label="p.title" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="模板">
          <el-select v-model="form.template_id" placeholder="选择模板" clearable style="width:100%">
            <el-option v-for="t in templates" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="文案">
          <el-input v-model="form.copy_text" type="textarea" :rows="5" />
          <el-button type="primary" size="small" :loading="genCopyLoading" class="mt-1" @click="onGenerateCopy">一键生成文案（Gemini）</el-button>
        </el-form-item>
        <el-form-item label="封面图URL">
          <el-input v-model="form.cover_image_url" placeholder="封面图地址" />
          <el-button type="primary" size="small" :loading="genCoverLoading" class="mt-1" @click="onGenerateCover">一键合成封面</el-button>
        </el-form-item>
        <el-form-item label="视频URL">
          <el-input v-model="form.video_url" placeholder="原片地址" />
          <el-button type="primary" size="small" :loading="processVideoLoading" class="mt-1" @click="onProcessVideo">视频去重</el-button>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width:120px">
            <el-option label="draft" value="draft" />
            <el-option label="ready" value="ready" />
            <el-option label="published" value="published" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saveLoading" @click="save">保存</el-button>
          <el-button @click="$router.push('/xsh/content/drafts')">返回列表</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getDraft, updateDraft, createDraft, listTemplates, generateCopy, generateCover, processVideo } from '@/api/content'
import { listProducts } from '@/api/product'

export default {
  name: 'ContentDraftEdit',
  setup () {
    const route = useRoute()
    const id = computed(() => route.params.id)
    const loading = ref(false)
    const saveLoading = ref(false)
    const genCopyLoading = ref(false)
    const genCoverLoading = ref(false)
    const processVideoLoading = ref(false)
    const products = ref([])
    const templates = ref([])
    const form = reactive({
      product_id: null,
      template_id: null,
      copy_text: '',
      cover_image_url: '',
      video_url: '',
      status: 'draft'
    })

    async function loadProducts () {
      if (products.value.length) return
      try {
        const res = await listProducts({ page_size: 200 })
        products.value = res.items || []
      } catch (e) {
        ElMessage.error(e.message)
      }
    }

    async function loadTemplates () {
      try {
        const res = await listTemplates({})
        templates.value = res.items || []
      } catch (e) {
        ElMessage.error(e.message)
      }
    }

    async function loadDraft () {
      if (!id.value) {
        const q = route.query.product_id
        if (q) form.product_id = Number(q)
        return
      }
      loading.value = true
      try {
        const res = await getDraft(id.value)
        const d = res.draft || res
        if (d) {
          form.product_id = d.product_id
          form.template_id = d.template_id || null
          form.copy_text = d.copy_text || ''
          form.cover_image_url = d.cover_image_url || ''
          form.video_url = d.video_url || ''
          form.status = d.status || 'draft'
        }
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    async function save () {
      saveLoading.value = true
      try {
        if (id.value) {
          await updateDraft(Number(id.value), form)
          ElMessage.success('已更新')
        } else {
          const res = await createDraft({ product_id: form.product_id, template_id: form.template_id || 0 })
          ElMessage.success('已创建')
          if (res.draft && res.draft.id) {
            form.copy_text = res.draft.copy_text || form.copy_text
            form.status = res.draft.status || form.status
            window.history.replaceState(null, '', `/xsh/content/draft/edit/${res.draft.id}`)
          }
        }
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        saveLoading.value = false
      }
    }

    async function onGenerateCopy () {
      const draftId = id.value
      if (!draftId) {
        ElMessage.warning('请先保存草稿后再生成文案')
        return
      }
      genCopyLoading.value = true
      try {
        const t = templates.value.find(x => x.id === form.template_id)
        const res = await generateCopy(Number(draftId), t ? t.style_type : 'default')
        if (res.copy_text) form.copy_text = res.copy_text
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        genCopyLoading.value = false
      }
    }

    async function onGenerateCover () {
      const draftId = id.value
      if (!draftId) {
        ElMessage.warning('请先保存草稿后再合成封面')
        return
      }
      genCoverLoading.value = true
      try {
        const res = await generateCover(Number(draftId))
        if (res.cover_image_url) form.cover_image_url = res.cover_image_url
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        genCoverLoading.value = false
      }
    }

    async function onProcessVideo () {
      const draftId = id.value
      if (!draftId) {
        ElMessage.warning('请先保存草稿后再处理视频')
        return
      }
      processVideoLoading.value = true
      try {
        const res = await processVideo(Number(draftId), '')
        if (res.video_url) form.video_url = res.video_url
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        processVideoLoading.value = false
      }
    }

    onMounted(async () => {
      await loadTemplates()
      await loadDraft()
    })

    return {
      loading, saveLoading, genCopyLoading, genCoverLoading, processVideoLoading,
      form, products, templates,
      loadProducts, save, onGenerateCopy, onGenerateCover, onProcessVideo
    }
  }
}
</script>

<style scoped>
.mt-1 { margin-top: 8px; }
</style>
