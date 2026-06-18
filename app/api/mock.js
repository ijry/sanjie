import { request } from './http'

export const getMockSession = () => request('/api/mock/session')
export const resetMockSession = () => request('/api/mock/reset', { method: 'POST' })
