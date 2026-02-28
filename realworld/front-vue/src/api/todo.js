/**
 * Todo API 服务
 * 封装所有与后端交互的 API 调用
 */

const API_BASE_URL = 'http://127.0.0.1:8000'

/**
 * 通用请求函数
 */
async function request(url, options = {}) {
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      ...options
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`API Error: ${response.status} ${errorText}`)
    }

    const data = await response.json()
    return data
  } catch (error) {
    console.error('API Request failed:', error)
    throw error
  }
}

/**
 * 添加待办事项
 * POST /api/add-todo
 * Body: { "value": "待办内容", "isCompleted": false }
 */
export async function addTodo(value, isCompleted = false) {
  return request('/api/add-todo', {
    method: 'POST',
    body: JSON.stringify({
      value,
      isCompleted
    })
  })
}

/**
 * 删除待办事项
 * POST /api/del-todo/{id}
 */
export async function deleteTodo(id) {
  return request(`/api/del-todo/${id}`, {
    method: 'POST',
    body: JSON.stringify({})
  })
}

/**
 * 查询所有待办事项
 * GET /api/get-todo
 */
export async function getTodoList() {
  return request('/api/get-todo', {
    method: 'GET'
  })
}

/**
 * 更新待办事项状态（将 isCompleted 取反）
 * POST /api/update-todo/{id}
 */
export async function updateTodo(id) {
  return request(`/api/update-todo/${id}`, {
    method: 'POST',
    body: JSON.stringify({})
  })
}
