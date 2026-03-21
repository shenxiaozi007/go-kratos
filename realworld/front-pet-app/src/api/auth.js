/**
 * 认证：登录、注册、获取资料、登出
 * 文档：docs-all/cat-admin/02-后端开发文档.md
 */

import { api } from './api'

const TOKEN_KEY = 'pet_app_token'

export function getStoredToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export function logout() {
  setToken(null)
  window.location.href = '/'
}

/** POST /api/v1/auth/register  Body: { username, password } */
export function register(username, password) {
  return api.post('/api/v1/auth/register', { username, password })
}

/** POST /api/v1/auth/login  Body: { username, password } 返回含 token */
export function login(username, password) {
  return api.post('/api/v1/auth/login', { username, password })
}

/** GET /api/v1/auth/profile 当前用户资料（需 token） */
export function getProfile() {
  return api.get('/api/v1/auth/profile')
}
