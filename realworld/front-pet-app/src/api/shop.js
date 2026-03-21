/**
 * 商店与背包 API
 */

import { api } from './api'

/** GET /api/v1/shop/items?category= */
export function listShopItems(category = '') {
  const q = category ? `?category=${encodeURIComponent(category)}` : ''
  return api.get(`/api/v1/shop/items${q}`)
}

/** POST /api/v1/shop/buy  Body: { item_id, quantity } */
export function buy(itemId, quantity = 1) {
  return api.post('/api/v1/shop/buy', { item_id: itemId, quantity })
}

/** GET /api/v1/inventory */
export function getInventory() {
  return api.get('/api/v1/inventory')
}

/** POST /api/v1/inventory/use  Body: { item_id, pet_id? } */
export function useItem(itemId, petId = 0) {
  return api.post('/api/v1/inventory/use', { item_id: itemId, pet_id: petId })
}
