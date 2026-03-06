/**
 * 推广引流管理后台 - Vue CLI 配置
 * 开发时 /api 代理到 Kratos 后端，与 xsh-docs 约定一致（接口前缀 /api/xsh/v1）
 */
const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 5174,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8000',
        changeOrigin: true
      }
    }
  }
})
