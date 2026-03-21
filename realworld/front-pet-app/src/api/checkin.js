/**
 * 签到 API
 */

import { api } from './api'

/** POST /api/v1/checkin 每日签到 */
export function doCheckin() {
  return api.post('/api/v1/checkin', {})
}

/** GET /api/v1/checkin/status 签到状态与 7 日日历 */
export function getCheckinStatus() {
  return api.get('/api/v1/checkin/status')
}
