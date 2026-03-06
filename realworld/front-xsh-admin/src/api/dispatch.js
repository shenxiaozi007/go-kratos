/**
 * 多账号分发矩阵 API
 */
import request from './request'

/** 平台账号列表 GET /accounts */
export function listPlatformAccounts (params = {}) {
  return request.get('/accounts', { params: { page: 1, page_size: 20, ...params } })
}

/** 绑定账号 POST /accounts（一账号一 IP：platform, nickname, proxy_url, cookie_storage_path） */
export function bindAccount (data) {
  return request.post('/accounts', {
    platform: data.platform,
    nickname: data.nickname,
    proxy_url: data.proxy_url,
    cookie_storage_path: data.cookie_storage_path
  })
}

/** 更新账号 PUT /accounts/:id */
export function updateAccount (id, data) {
  return request.put(`/accounts/${id}`, { id, ...data })
}

/** 解绑账号 DELETE /accounts/:id */
export function unbindAccount (id) {
  return request.delete(`/accounts/${id}`)
}

/** 定时规则列表 GET /schedules */
export function listSchedules () {
  return request.get('/schedules')
}

/** 创建定时规则 POST /schedules */
export function createSchedule (data) {
  return request.post('/schedules', data)
}

/** 更新定时规则 PUT /schedules/:id */
export function updateSchedule (id, data) {
  return request.put(`/schedules/${id}`, { id, ...data })
}

/** 删除定时规则 DELETE /schedules/:id */
export function deleteSchedule (id) {
  return request.delete(`/schedules/${id}`)
}

/** 创建发布任务 POST /publish-tasks（draft_id, account_ids[], scheduled_at） */
export function createPublishTask (data) {
  return request.post('/publish-tasks', data)
}

/** 发布任务列表 GET /publish-tasks */
export function listPublishTasks (params = {}) {
  return request.get('/publish-tasks', { params: { page: 1, page_size: 20, ...params } })
}

/** 重试任务 POST /publish-tasks/:id/retry */
export function retryPublishTask (id) {
  return request.post(`/publish-tasks/${id}/retry`, { id })
}

/** 取消任务 POST /publish-tasks/:id/cancel */
export function cancelPublishTask (id) {
  return request.post(`/publish-tasks/${id}/cancel`, { id })
}
