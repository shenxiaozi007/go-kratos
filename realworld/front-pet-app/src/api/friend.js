/**
 * 社交、排行榜 API
 */

import { api } from './api'

/** GET /api/v1/friends */
export function listFriends() {
  return api.get('/api/v1/friends')
}

/** GET /api/v1/friends/search?keyword= */
export function searchUsers(keyword) {
  return api.get(`/api/v1/friends/search?keyword=${encodeURIComponent(keyword)}`)
}

/** POST /api/v1/friends/request  Body: { to_user_id, message } */
export function sendFriendRequest(toUserId, message = '') {
  return api.post('/api/v1/friends/request', { to_user_id: toUserId, message })
}

/** GET /api/v1/friends/requests */
export function listFriendRequests() {
  return api.get('/api/v1/friends/requests')
}

/** POST /api/v1/friends/requests/:id/accept */
export function acceptFriendRequest(id) {
  return api.post(`/api/v1/friends/requests/${id}/accept`, { id })
}

/** POST /api/v1/friends/requests/:id/reject */
export function rejectFriendRequest(id) {
  return api.post(`/api/v1/friends/requests/${id}/reject`, { id })
}

/** GET /api/v1/ranking?type=heart_points|level&limit= */
export function getRanking(type = 'heart_points', limit = 50) {
  return api.get(`/api/v1/ranking?type=${type}&limit=${limit}`)
}
