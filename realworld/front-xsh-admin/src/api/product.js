/**
 * 选品采集仓 API（与 backend api-catalog 一致）
 */
import request from './request'

/** 解析链接 POST /products/parse */
export function parseProductLink (url) {
  return request.post('/products/parse', { url })
}

/** 商品分页列表 GET /products */
export function listProducts (params = {}) {
  return request.get('/products', { params: { page: 1, page_size: 20, ...params } })
}

/** 商品详情 GET /products/:id */
export function getProduct (id) {
  return request.get(`/products/${id}`)
}

/** 创建商品 POST /products */
export function createProduct (data) {
  return request.post('/products', data)
}

/** 更新商品 PUT /products/:id */
export function updateProduct (id, data) {
  return request.put(`/products/${id}`, { id, ...data })
}

/** 删除商品 DELETE /products/:id */
export function deleteProduct (id) {
  return request.delete(`/products/${id}`)
}

/** 爆款榜单列表 GET /hot-ranks */
export function listHotRanks (params = {}) {
  return request.get('/hot-ranks', { params: { page: 1, page_size: 20, ...params } })
}

/** 同步爆款榜单 POST /hot-ranks/sync */
export function syncHotRank (rankType) {
  return request.post('/hot-ranks/sync', { rank_type: rankType })
}
