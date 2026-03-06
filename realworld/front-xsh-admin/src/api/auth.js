/**
 * v2 鉴权 API：登录、获取当前用户资料（与 backend api/auth/v1 一致）
 */
import axios from 'axios'

const authRequest = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' }
})

/**
 * 登录，返回 { access_token, expires_in }
 */
export function login(username, password) {
  return authRequest.post('/auth/login', { username, password }).then((res) => res.data)
}

/**
 * 注册，返回 { user_id, username }
 * 注意：注册不会返回 token，需要再调用 login 或跳转登录页。
 */
export function register(username, password) {
  return authRequest.post('/auth/register', { username, password }).then((res) => res.data)
}

/**
 * 获取当前用户信息（需已登录，Header 带 Authorization: Bearer <token>）
 * 返回 { user_id, username, role, ... }
 */
export function getProfile() {
  const token = localStorage.getItem('xsh_access_token')
  return authRequest.get('/auth/profile', {
    headers: token ? { Authorization: `Bearer ${token}` } : {}
  }).then((res) => res.data)
}

export default authRequest
