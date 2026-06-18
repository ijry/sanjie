import { buildQuery, request } from './http'

export const getReincarnations = params => request(`/api/reincarnations${buildQuery(params)}`)
export const getReincarnation = id => request(`/api/reincarnations/${id}`)
export const approveReincarnation = (id, data) => request(`/api/reincarnations/${id}/approve`, { method: 'POST', data })
export const rejectReincarnation = (id, data) => request(`/api/reincarnations/${id}/reject`, { method: 'POST', data })
export const reviewReincarnation = (id, data) => request(`/api/reincarnations/${id}/review`, { method: 'POST', data })
export const assignReincarnationQuota = (id, data) => request(`/api/reincarnations/${id}/assign-quota`, { method: 'POST', data })
export const issueReincarnationSoup = (id, data) => request(`/api/reincarnations/${id}/soup`, { method: 'POST', data })
export const completeReincarnation = (id, data) => request(`/api/reincarnations/${id}/complete`, { method: 'POST', data })

