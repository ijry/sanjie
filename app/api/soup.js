import { request } from './http'

export const getSoupInventory = () => request('/api/soup/inventory')
export const adjustSoupInventory = (id, data) => request(`/api/soup/inventory/${id}/adjust`, { method: 'POST', data })
export const getSoupRecords = () => request('/api/soup/records')

