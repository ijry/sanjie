import { buildQuery, request } from './http'

export const getSouls = params => request(`/api/souls${buildQuery(params)}`)
export const getSoul = id => request(`/api/souls/${id}`)

