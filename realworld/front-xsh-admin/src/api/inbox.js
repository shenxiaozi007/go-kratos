/**
 * 获客截流 API
 */
import request from './request'

/** 评论列表 GET /comments */
export function listComments (params = {}) {
  return request.get('/comments', { params: { page: 1, page_size: 20, ...params } })
}

/** 标记评论已处理 POST /comments/:id/handled */
export function markCommentHandled (id, status) {
  return request.post(`/comments/${id}/handled`, { id, status })
}

/** 同步评论 POST /comments/sync */
export function syncComments (accountId) {
  return request.post('/comments/sync', { account_id: accountId })
}

/** 回复评论 POST /comments/:comment_id/reply */
export function sendReply (commentId, content) {
  return request.post(`/comments/${commentId}/reply`, { comment_id: commentId, content })
}

/** 私信 POST /comments/:comment_id/private-message */
export function sendPrivateMessage (commentId, content) {
  return request.post(`/comments/${commentId}/private-message`, { comment_id: commentId, content })
}

/** 动作记录列表 GET /inbox-actions */
export function listInboxActions (params = {}) {
  return request.get('/inbox-actions', { params: { page: 1, page_size: 20, ...params } })
}
