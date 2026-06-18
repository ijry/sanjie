import { buildQuery, request } from './http'

export const getApprovals = params => request(`/api/approvals${buildQuery(params)}`)
export const getApproval = id => request(`/api/approvals/${id}`)
export const approveApproval = (id, data) => request(`/api/approvals/${id}/approve`, { method: 'POST', data })
export const rejectApproval = (id, data) => request(`/api/approvals/${id}/reject`, { method: 'POST', data })
export const transferApproval = (id, data) => request(`/api/approvals/${id}/transfer`, { method: 'POST', data })

