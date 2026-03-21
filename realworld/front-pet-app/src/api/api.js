/**
 * 萌宠之家 API 封装
 * 文档：docs-all/cat-admin/03-前端开发文档.md
 * 统一 baseURL、Authorization: Bearer <token>
 */

const getBaseURL = () =>
  typeof window !== 'undefined' && window.APP_CONFIG?.apiBaseUrl
    ? window.APP_CONFIG.apiBaseUrl
    : ''

const getToken = () => localStorage.getItem('pet_app_token')

export async function request(path, options = {}) {
  const url = path.startsWith('http') ? path : `${getBaseURL()}${path}`
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(url, { ...options, headers })
  const text = await res.text()
  let data = null
  try {
    data = text ? JSON.parse(text) : null
  } catch (_) {}

  if (!res.ok) {
    const err = new Error(data?.message || res.statusText || '请求失败')
    err.status = res.status
    err.data = data
    throw err
  }
  return data
}

export const api = {
  get: (path) => request(path, { method: 'GET' }),
  post: (path, body) => request(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  put: (path, body) => request(path, { method: 'PUT', body: body ? JSON.stringify(body) : undefined }),
}
