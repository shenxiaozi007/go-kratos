/**
 * 品种 API
 */

import { api } from './api'

/** GET /api/v1/breeds?species=CAT|DOG */
export function listBreeds(species = '') {
  const q = species ? `?species=${encodeURIComponent(species)}` : ''
  return api.get(`/api/v1/breeds${q}`)
}

/** GET /api/v1/breeds/:id */
export function getBreed(id) {
  return api.get(`/api/v1/breeds/${id}`)
}
