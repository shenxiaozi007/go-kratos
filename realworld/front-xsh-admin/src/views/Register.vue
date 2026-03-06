<template>
  <div class="register-page">
    <el-card class="register-card" shadow="always">
      <template #header>
        <span>注册账号</span>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px" @submit.prevent="onSubmit">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input
            v-model="form.confirm_password"
            type="password"
            placeholder="请再次输入密码"
            show-password
            clearable
            @keyup.enter="onSubmit"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width: 100%" @click="onSubmit">注册并登录</el-button>
        </el-form-item>
      </el-form>

      <div class="tip">
        已有账号？
        <el-button type="primary" link @click="$router.replace('/login')">返回登录</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { register, login } from '../api/auth'
import { TOKEN_KEY } from '../api/request'

export default {
  name: 'RegisterPage',
  data() {
    return {
      form: { username: '', password: '', confirm_password: '' },
      loading: false,
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
        confirm_password: [
          { required: true, message: '请再次输入密码', trigger: 'blur' },
          {
            validator: (rule, value, callback) => {
              // 自定义校验：确认密码需与密码一致
              if (value !== this.form.password) return callback(new Error('两次输入的密码不一致'))
              callback()
            },
            trigger: 'blur'
          }
        ]
      }
    }
  },
  methods: {
    onSubmit() {
      this.$refs.formRef.validate(async (valid) => {
        if (!valid) return
        this.loading = true
        try {
          // 先注册，再自动登录，减少一次人工跳转
          await register(this.form.username, this.form.password)
          const res = await login(this.form.username, this.form.password)
          localStorage.setItem(TOKEN_KEY, res.access_token)
          const expiresIn = res.expires_in || 86400
          localStorage.setItem('xsh_expires_at', String(Date.now() + expiresIn * 1000))
          this.$message.success('注册成功，已自动登录')
          this.$router.replace('/xsh/product/list')
        } catch (e) {
          this.$message.error(e.message || '注册失败')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}
.register-card {
  width: 420px;
}
.tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
  text-align: right;
}
</style>

