import { request } from './http'

export const getMe = () => request('/api/users/me')
export const switchUser = userId => request('/api/users/switch', { method: 'POST', data: { userId } })

