import { buildQuery, request } from './http'

export const getWishes = params => request(`/api/wishes${buildQuery(params)}`)
export const getWish = id => request(`/api/wishes/${id}`)
export const routeWish = (id, data) => request(`/api/wishes/${id}/route`, { method: 'POST', data })
export const resolveWish = (id, data) => request(`/api/wishes/${id}/resolve`, { method: 'POST', data })
export const rejectWish = (id, data) => request(`/api/wishes/${id}/reject`, { method: 'POST', data })

