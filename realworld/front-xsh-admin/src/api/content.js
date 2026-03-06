/**
 * AI 内容加工坊 API
 */
import request from './request'

/** 模板列表 GET /templates */
export function listTemplates (params = {}) {
  return request.get('/templates', { params })
}

/** 创建模板 POST /templates */
export function createTemplate (data) {
  return request.post('/templates', data)
}

/** 更新模板 PUT /templates/:id */
export function updateTemplate (id, data) {
  return request.put(`/templates/${id}`, { id, ...data })
}

/** 删除模板 DELETE /templates/:id */
export function deleteTemplate (id) {
  return request.delete(`/templates/${id}`)
}

/** 创建草稿 POST /drafts */
export function createDraft (data) {
  return request.post('/drafts', data)
}

/** 草稿分页列表 GET /drafts */
export function listDrafts (params = {}) {
  return request.get('/drafts', { params: { page: 1, page_size: 20, ...params } })
}

/** 草稿详情 GET /drafts/:id */
export function getDraft (id) {
  return request.get(`/drafts/${id}`)
}

/** 更新草稿 PUT /drafts/:id */
export function updateDraft (id, data) {
  return request.put(`/drafts/${id}`, { id, ...data })
}

/** 删除草稿 DELETE /drafts/:id */
export function deleteDraft (id) {
  return request.delete(`/drafts/${id}`)
}

/** 生成文案 POST /drafts/:draft_id/generate-copy（body: style，对应模板风格） */
export function generateCopy (draftId, style) {
  return request.post(`/drafts/${draftId}/generate-copy`, { style: style || 'default' })
}

/** 生成封面 POST /drafts/:draft_id/generate-cover */
export function generateCover (draftId) {
  return request.post(`/drafts/${draftId}/generate-cover`, {})
}

/** 视频去重 POST /drafts/:draft_id/process-video */
export function processVideo (draftId, options = '') {
  return request.post(`/drafts/${draftId}/process-video`, { draft_id: draftId, options })
}
