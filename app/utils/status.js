export const statusText = {
  pending: '待处理',
  running: '执行中',
  completed: '已完成',
  exception: '异常',
  canceled: '已撤销',
  queued: '排队中',
  approved: '已通过',
  rejected: '已驳回',
  reviewing: '复核中',
  pending_soup: '待喝汤',
  issued: '已发放',
  failed: '失败',
  critical: '严重',
  warning: '预警',
  info: '提示',
  none: '未复核',
  routed: '已转派',
  transferred: '已转交',
  normal: '正常'
}

export const statusType = {
  pending: 'warning',
  running: 'primary',
  completed: 'success',
  exception: 'error',
  canceled: 'info',
  queued: 'warning',
  approved: 'success',
  rejected: 'error',
  reviewing: 'primary',
  pending_soup: 'warning',
  issued: 'success',
  failed: 'error',
  critical: 'error',
  warning: 'warning',
  info: 'info',
  none: 'info',
  routed: 'primary',
  transferred: 'primary',
  normal: 'success'
}

export function labelOf(value) {
  return statusText[value] || value || '-'
}

export function typeOf(value) {
  return statusType[value] || 'info'
}

