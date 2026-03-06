/**
 * 封装 axios，统一 baseURL、错误与响应格式
 * v2：请求头自动带 Authorization: Bearer <token>；401 时清除登录态并跳转登录页
 */
import axios from 'axios'

const TOKEN_KEY = 'xsh_access_token'

const request = axios.create({
  baseURL: '/api/xsh/v1',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem('xsh_expires_at')
      if (!window.location.pathname.startsWith('/login')) {
        window.location.pathname = '/login'
      }
      return Promise.reject(new Error('登录已过期，请重新登录'))
    }
    const data = err.response?.data
    const msg = data?.message || err.message || '请求失败'
    return Promise.reject(new Error(msg))
  }
)

export default request
export { TOKEN_KEY }
