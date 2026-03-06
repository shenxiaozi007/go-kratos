/**
 * Vite 配置：推广引流管理后台（front-xsh-admin）
 * 开发时 /api 代理到 Kratos 后端，与 xsh-docs 约定一致。
 */
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5174,
    // 与后端联通：所有 /api 请求转发到 Kratos HTTP 服务（默认 8000）
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8000',
        changeOrigin: true
      }
    }
  }
})
