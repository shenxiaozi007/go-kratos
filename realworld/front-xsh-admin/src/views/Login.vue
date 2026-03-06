<template>
  <div class="login-page">
    <el-card class="login-card" shadow="always">
      <template #header>
        <span>推广引流管理后台</span>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px" @submit.prevent="onSubmit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password clearable @keyup.enter="onSubmit" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" @click="onSubmit">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="tip">
        v2：需先注册或使用已有账号登录
        <el-button type="primary" link @click="$router.replace('/register')">去注册</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { login } from '../api/auth'
import { TOKEN_KEY } from '../api/request'

export default {
  name: 'LoginPage',
  data() {
    return {
      form: { username: '', password: '' },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      loading: false
    }
  },
  methods: {
    onSubmit() {
      this.$refs.formRef.validate(async (valid) => {
        if (!valid) return
        this.loading = true
        try {
          const res = await login(this.form.username, this.form.password)
          // 兼容后端 protojson 的 camelCase 与 snake_case
          const token = res.access_token ?? res.accessToken
          const expiresIn = res.expires_in ?? res.expiresIn ?? 86400
          if (!token) {
            this.$message.error('登录响应缺少 token')
            return
          }
          localStorage.setItem(TOKEN_KEY, token)
          localStorage.setItem('xsh_expires_at', String(Date.now() + expiresIn * 1000))
          this.$message.success('登录成功')
          const redirect = this.$route.query.redirect || '/xsh/product/list'
          this.$router.replace(redirect)
        } catch (e) {
          this.$message.error(e.message || '登录失败')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}
.login-card {
  width: 400px;
}
.tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}
</style>
