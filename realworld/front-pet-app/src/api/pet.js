/**
 * 宠物、互动、装扮 API
 */

import { api } from './api'

/** GET /api/v1/pet/me 主宠物 */
export function getMyPet() {
  return api.get('/api/v1/pet/me')
}

/** GET /api/v1/pet/list 宠物列表 */
export function listPets() {
  return api.get('/api/v1/pet/list')
}

/** POST /api/v1/pet 领养 */
export function createPet(body) {
  return api.post('/api/v1/pet', body)
}

/** PUT /api/v1/pet/:id 更新 */
export function updatePet(id, body) {
  return api.put(`/api/v1/pet/${id}`, body)
}

/** POST /api/v1/pet/interact  Body: { action: "STROKE"|"TEASE"|"TREAT"|"BATH" } */
export function interact(action) {
  return api.post('/api/v1/pet/interact', { action })
}

/** GET /api/v1/pet/:id/appearance */
export function getPetAppearance(id) {
  return api.get(`/api/v1/pet/${id}/appearance`)
}

/** PUT /api/v1/pet/:id/appearance */
export function updatePetAppearance(id, body) {
  return api.put(`/api/v1/pet/${id}/appearance`, body)
}
