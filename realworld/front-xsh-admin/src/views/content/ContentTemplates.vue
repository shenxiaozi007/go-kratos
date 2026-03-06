<template>
  <div class="content-templates">
    <el-card>
      <el-button type="primary" @click="openAdd">新增模板</el-button>
      <el-table v-loading="loading" :data="items" border stripe class="table">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="style_type" label="风格类型" width="120" />
        <el-table-column prop="content_template" label="模板内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑模板' : '新增模板'" width="560">
      <el-form ref="formRef" :model="form" label-width="100px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="风格类型"><el-input v-model="form.style_type" placeholder="如 暴躁省钱风" /></el-form-item>
        <el-form-item label="模板内容"><el-input v-model="form.content_template" type="textarea" :rows="4" placeholder="占位符 {{title}} {{price}}" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
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
import { listTemplates, createTemplate, updateTemplate, deleteTemplate } from '@/api/content'

export default {
  name: 'ContentTemplates',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const dialogVisible = ref(false)
    const isEdit = ref(false)
    const saveLoading = ref(false)
    const formRef = ref(null)
    const form = reactive({ id: null, name: '', style_type: '', content_template: '', sort_order: 0 })

    async function load () {
      loading.value = true
      try {
        const res = await listTemplates({})
        items.value = res.items || []
      } catch (e) {
        ElMessage.error(e.message)
      } finally {
        loading.value = false
      }
    }

    function openAdd () {
      isEdit.value = false
      form.id = null
      form.name = ''
      form.style_type = ''
      form.content_template = ''
      form.sort_order = 0
      dialogVisible.value = true
    }

    function openEdit (row) {
      isEdit.value = true
      form.id = row.id
      form.name = row.name
      form.style_type = row.style_type
      form.content_template = row.content_template
      form.sort_order = row.sort_order ?? 0
      dialogVisible.value = true
    }

    async function save () {
      saveLoading.value = true
      try {
        if (isEdit.value) {
          await updateTemplate(form.id, { name: form.name, style_type: form.style_type, content_template: form.content_template, sort_order: form.sort_order })
        } else {
          await createTemplate({ name: form.name, style_type: form.style_type, content_template: form.content_template, sort_order: form.sort_order })
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

    function onDelete (row) {
      ElMessageBox.confirm('确定删除该模板？', '提示', { type: 'warning' }).then(async () => {
        try {
          await deleteTemplate(row.id)
          ElMessage.success('已删除')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    onMounted(load)
    return { loading, items, dialogVisible, isEdit, formRef, form, saveLoading, load, openAdd, openEdit, save, onDelete }
  }
}
</script>

<style scoped>
.table { margin-top: 12px; }
</style>
