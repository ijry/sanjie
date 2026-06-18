import { request } from './http'

export const getAuditLogs = () => request('/api/audit-logs')

