export function showError(error) {
  uni.showToast({
    title: error?.message || '操作失败',
    icon: 'none'
  })
}

export function showSuccess(title = '操作成功') {
  uni.showToast({ title, icon: 'success' })
}

export function displayTime(value) {
  if (!value) return '-'
  return String(value).replace('T', ' ').replace(/\+.*/, '')
}

