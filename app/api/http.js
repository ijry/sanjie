import { mockRequest } from '../mock/server'

export const API_BASE_URL = import.meta.env.VITE_SANJIE_API_BASE_URL ?? 'http://localhost:8080'
export const MOCK_ENABLED = import.meta.env.VITE_SANJIE_MOCK === 'true'

export function buildQuery(params = {}) {
  const search = Object.entries(params)
    .filter(([, value]) => value !== undefined && value !== null && value !== '' && value !== 'all')
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
    .join('&')
  return search ? `?${search}` : ''
}

export function request(path, options = {}) {
  if (MOCK_ENABLED) {
    return mockRequest(path, options)
  }
  const userId = uni.getStorageSync('sanjie_user_id') || 1
  return new Promise((resolve, reject) => {
    uni.request({
      url: `${API_BASE_URL}${path}`,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'X-User-ID': String(userId),
        ...(options.header || {})
      },
      success(res) {
        const body = res.data
        if (res.statusCode >= 200 && res.statusCode < 300 && body && body.code === 0) {
          resolve(body.data)
          return
        }
        reject(new Error((body && body.message) || '请求失败'))
      },
      fail(err) {
        reject(new Error(err.errMsg || '网络异常'))
      }
    })
  })
}
