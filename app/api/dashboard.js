import { request } from './http'

export const getDashboard = () => request('/api/dashboard')
export const getAlerts = () => request('/api/alerts')
export const handleAlert = id => request(`/api/alerts/${id}/handle`, { method: 'POST' })

