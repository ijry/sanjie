import { buildQuery, request } from './http'

export const getLifeBookRecords = params => request(`/api/life-book${buildQuery(params)}`)
export const getLifeBookRecord = id => request(`/api/life-book/${id}`)
export const changeLifespan = (id, data) => request(`/api/life-book/${id}/change-lifespan`, { method: 'POST', data })
export const freezeLifeBookRecord = id => request(`/api/life-book/${id}/freeze`, { method: 'POST' })
export const getLifeBookAuditLogs = id => request(`/api/life-book/${id}/audit-logs`)

