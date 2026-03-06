<template>
  <div class="dispatch-schedules">
    <el-card>
      <el-button type="primary" @click="openAdd">新增定时规则</el-button>
      <el-table v-loading="loading" :data="items" border stripe class="table">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="cron_expr" label="Cron" width="140" />
        <el-table-column prop="account_ids" label="关联账号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="is_enabled" label="启用" width="80">
          <template #default="{ row }"><el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">{{ row.is_enabled ? '是' : '否' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑规则' : '新增规则'" width="500">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="Cron 表达式"><el-input v-model="form.cron_expr" placeholder="如 0 10,12,20 * * *" /></el-form-item>
        <el-form-item label="账号 IDs"><el-input v-model="form.account_ids" type="textarea" :rows="2" placeholder="逗号分隔或 JSON 数组" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.is_enabled" /></el-form-item>
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
import { listSchedules, createSchedule, updateSchedule, deleteSchedule } from '@/api/dispatch'

export default {
  name: 'DispatchSchedules',
  setup () {
    const loading = ref(false)
    const items = ref([])
    const dialogVisible = ref(false)
    const isEdit = ref(false)
    const saveLoading = ref(false)
    const form = reactive({ id: null, name: '', cron_expr: '', account_ids: '', is_enabled: true })

    async function load () {
      loading.value = true
      try {
        const res = await listSchedules()
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
      form.cron_expr = ''
      form.account_ids = ''
      form.is_enabled = true
      dialogVisible.value = true
    }

    function openEdit (row) {
      isEdit.value = true
      form.id = row.id
      form.name = row.name
      form.cron_expr = row.cron_expr
      form.account_ids = row.account_ids || ''
      form.is_enabled = row.is_enabled !== false
      dialogVisible.value = true
    }

    async function save () {
      saveLoading.value = true
      try {
        if (isEdit.value) {
          await updateSchedule(form.id, { name: form.name, cron_expr: form.cron_expr, account_ids: form.account_ids, is_enabled: form.is_enabled })
        } else {
          await createSchedule({ name: form.name, cron_expr: form.cron_expr, account_ids: form.account_ids, is_enabled: form.is_enabled })
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
      ElMessageBox.confirm('确定删除该规则？', '提示', { type: 'warning' }).then(async () => {
        try {
          await deleteSchedule(row.id)
          ElMessage.success('已删除')
          load()
        } catch (e) {
          ElMessage.error(e.message)
        }
      }).catch(() => {})
    }

    onMounted(load)
    return { loading, items, dialogVisible, isEdit, form, saveLoading, load, openAdd, openEdit, save, onDelete }
  }
}
</script>

<style scoped>
.table { margin-top: 12px; }
</style>
