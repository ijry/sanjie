import { buildQuery, request } from './http'

export const getCaptureTasks = params => request(`/api/capture-tasks${buildQuery(params)}`)
export const getCaptureTask = id => request(`/api/capture-tasks/${id}`)
export const startCaptureTask = id => request(`/api/capture-tasks/${id}/start`, { method: 'POST' })
export const completeCaptureTask = (id, data) => request(`/api/capture-tasks/${id}/complete`, { method: 'POST', data })
export const reportCaptureException = (id, data) => request(`/api/capture-tasks/${id}/exception`, { method: 'POST', data })
export const cancelCaptureTask = id => request(`/api/capture-tasks/${id}/cancel`, { method: 'POST' })

