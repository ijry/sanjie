import { buildQuery, request } from './http'

export const getHellFloors = () => request('/api/hell/floors')
export const getHellSentences = params => request(`/api/hell/sentences${buildQuery(params)}`)
export const getHellSentence = id => request(`/api/hell/sentences/${id}`)
export const reviewHellSentence = (id, data) => request(`/api/hell/sentences/${id}/review`, { method: 'POST', data })
export const dispatchHellFloor = (id, data) => request(`/api/hell/floors/${id}/dispatch`, { method: 'POST', data })

