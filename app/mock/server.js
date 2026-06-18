import { createMockState } from './data'

const MOCK_STATE_KEY = 'sanjie_mock_state_v1'
const MOCK_STARTED_KEY = 'sanjie_mock_started_at'

let state = loadMockState()

export function mockRequest(path, options = {}) {
  const method = (options.method || 'GET').toUpperCase()
  const data = options.data || {}
  const userId = Number(uni.getStorageSync('sanjie_user_id') || state.currentUserId || 1)
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      try {
        resolve(route(method, path, data, userId))
      } catch (error) {
        reject(error)
      }
    }, 120)
  })
}

function route(method, rawPath, data, userId) {
  const { pathname, searchParams } = parsePath(rawPath)
  if (method === 'GET' && pathname === '/api/health') return { status: 'ok', mock: true }
  if (method === 'GET' && pathname === '/api/mock/session') return mockSession()
  if (method === 'POST' && pathname === '/api/mock/reset') return resetMockState()
  if (method === 'GET' && pathname === '/api/users/me') return currentUser(userId)
  if (method === 'POST' && pathname === '/api/users/switch') return switchUser(data.userId)
  if (method === 'GET' && pathname === '/api/dashboard') return dashboard()
  if (method === 'GET' && pathname === '/api/alerts') return state.alerts
  if (method === 'POST' && match(pathname, '/api/alerts/:id/handle')) return handleAlert(idOf(pathname))
  if (method === 'GET' && pathname === '/api/audit-logs') return state.auditLogs
  if (method === 'GET' && pathname === '/api/souls') return filterKeyword(state.souls, searchParams, item => item.name)
  if (method === 'GET' && match(pathname, '/api/souls/:id')) return findByID(state.souls, idOf(pathname))
  if (method === 'GET' && pathname === '/api/capture-tasks') return filterStatus(state.captureTasks, searchParams)
  if (method === 'GET' && match(pathname, '/api/capture-tasks/:id')) return findByID(state.captureTasks, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/capture-tasks/:id/start')) return updateCapture(idOf(pathname), { status: 'running' })
  if (method === 'POST' && match(pathname, '/api/capture-tasks/:id/complete')) return updateCapture(idOf(pathname), { status: 'completed', actualTime: now() })
  if (method === 'POST' && match(pathname, '/api/capture-tasks/:id/exception')) return updateCapture(idOf(pathname), { status: 'exception', exceptionType: data.exceptionType || data.type || 'wrong_soul', exceptionNote: data.note || '' })
  if (method === 'POST' && match(pathname, '/api/capture-tasks/:id/cancel')) return updateCapture(idOf(pathname), { status: 'canceled' })
  if (method === 'GET' && pathname === '/api/life-book') return filterLifeBook(searchParams)
  if (method === 'GET' && match(pathname, '/api/life-book/:id')) return findByID(state.lifeBookRecords, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/life-book/:id/freeze')) return freezeLifeBook(idOf(pathname), userId)
  if (method === 'POST' && match(pathname, '/api/life-book/:id/change-lifespan')) return createLifespanApproval(idOf(pathname), data, userId)
  if (method === 'GET' && match(pathname, '/api/life-book/:id/audit-logs')) return state.auditLogs.filter(item => item.targetType === 'life_book' && item.targetId === idOf(pathname))
  if (method === 'GET' && pathname === '/api/approvals') return filterStatus(state.approvals, searchParams)
  if (method === 'GET' && match(pathname, '/api/approvals/:id')) return findByID(state.approvals, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/approvals/:id/approve')) return updateApproval(idOf(pathname), 'approved', data, userId)
  if (method === 'POST' && match(pathname, '/api/approvals/:id/reject')) return updateApproval(idOf(pathname), 'rejected', data, userId)
  if (method === 'POST' && match(pathname, '/api/approvals/:id/transfer')) return updateApproval(idOf(pathname), 'transferred', data, userId)
  if (method === 'GET' && pathname === '/api/reincarnations') return filterStatus(state.reincarnations, searchParams)
  if (method === 'GET' && match(pathname, '/api/reincarnations/:id')) return findByID(state.reincarnations, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/approve')) return updateReincarnation(idOf(pathname), { status: 'approved', reviewNote: data.note || '' })
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/reject')) return updateReincarnation(idOf(pathname), { status: 'rejected', reviewNote: data.note || '' })
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/review')) return updateReincarnation(idOf(pathname), { status: 'reviewing', reviewNote: data.note || '' })
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/assign-quota')) return updateReincarnation(idOf(pathname), { status: 'pending_soup', quotaType: data.quotaType || '普通胎' })
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/soup')) return issueSoup(idOf(pathname), data, userId)
  if (method === 'POST' && match(pathname, '/api/reincarnations/:id/complete')) return updateReincarnation(idOf(pathname), { status: 'completed' })
  if (method === 'GET' && pathname === '/api/hell/floors') return state.hellFloors
  if (method === 'GET' && pathname === '/api/hell/sentences') return state.hellSentences
  if (method === 'GET' && match(pathname, '/api/hell/sentences/:id')) return findByID(state.hellSentences, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/hell/sentences/:id/review')) return reviewHellSentence(idOf(pathname), data, userId)
  if (method === 'POST' && match(pathname, '/api/hell/floors/:id/dispatch')) return dispatchHellFloor(idOf(pathname), data, userId)
  if (method === 'GET' && pathname === '/api/wishes') return filterStatus(state.wishes, searchParams)
  if (method === 'GET' && match(pathname, '/api/wishes/:id')) return findByID(state.wishes, idOf(pathname))
  if (method === 'POST' && match(pathname, '/api/wishes/:id/route')) return updateWish(idOf(pathname), { status: 'routed', assignedDeity: data.assignedDeity || '财神', resultNote: data.note || '' }, userId, 'wish.route')
  if (method === 'POST' && match(pathname, '/api/wishes/:id/resolve')) return updateWish(idOf(pathname), { status: 'resolved', resultNote: data.note || '' }, userId, 'wish.resolve')
  if (method === 'POST' && match(pathname, '/api/wishes/:id/reject')) return updateWish(idOf(pathname), { status: 'rejected', resultNote: data.note || '' }, userId, 'wish.reject')
  if (method === 'GET' && pathname === '/api/soup/inventory') return state.soupInventory
  if (method === 'POST' && match(pathname, '/api/soup/inventory/:id/adjust')) return adjustSoup(idOf(pathname), data, userId)
  if (method === 'GET' && pathname === '/api/soup/records') return state.soupRecords
  throw new Error(`mock endpoint not implemented: ${method} ${pathname}`)
}

function currentUser(userId) {
  const user = state.users.find(item => item.id === Number(userId)) || state.users[0]
  return { ...user, roles: [user.role] }
}

function switchUser(userId) {
  state.currentUserId = Number(userId) || 1
  persistMockState()
  return { userId: state.currentUserId }
}

function mockSession() {
  return {
    persisted: true,
    storageKey: MOCK_STATE_KEY,
    startedAt: uni.getStorageSync(MOCK_STARTED_KEY) || '',
    currentUserId: state.currentUserId,
    auditLogCount: state.auditLogs.length,
    approvalCount: state.approvals.length,
    soupRecordCount: state.soupRecords.length,
    handledAlertCount: state.alerts.filter(item => item.handled).length
  }
}

function resetMockState() {
  state = createMockState()
  uni.setStorageSync('sanjie_user_id', 1)
  persistMockState(true)
  return mockSession()
}

function dashboard() {
  return {
    pendingCaptureCount: state.captureTasks.filter(item => item.status === 'pending').length,
    pendingApprovalCount: state.approvals.filter(item => ['pending', 'transferred'].includes(item.status)).length,
    reincarnationQueueCount: state.reincarnations.filter(item => ['queued', 'reviewing', 'pending_soup'].includes(item.status)).length,
    highRiskLifeBookCount: state.lifeBookRecords.filter(item => item.riskFlag !== 'normal').length,
    unhandledAlertCount: state.alerts.filter(item => !item.handled).length
  }
}

function handleAlert(id) {
  const item = findByID(state.alerts, id)
  item.handled = true
  persistMockState()
  return { id, handled: true }
}

function updateCapture(id, patch) {
  const item = findByID(state.captureTasks, id)
  Object.assign(item, patch)
  persistMockState()
  return { id, status: item.status }
}

function freezeLifeBook(id, userId) {
  const item = findByID(state.lifeBookRecords, id)
  item.locked = true
  item.riskFlag = 'critical'
  item.updatedAt = now()
  addAudit(userId, 'life_book.freeze', 'life_book', id, 'mock 冻结生死簿')
  persistMockState()
  return { id, locked: true }
}

function createLifespanApproval(id, data, userId) {
  const approval = {
    id: state.nextApprovalId++,
    type: 'lifespan_change',
    targetId: id,
    title: '寿命变更申请',
    status: 'pending',
    applicantId: userId,
    approverId: null,
    reason: data.reason || `调整到 ${data.newLifespan}`,
    resultNote: '',
    createdAt: now(),
    updatedAt: now()
  }
  state.approvals.unshift(approval)
  addAudit(userId, 'approval.create', 'life_book', id, approval.reason)
  persistMockState()
  return { approvalId: approval.id }
}

function updateApproval(id, status, data, userId) {
  const item = findByID(state.approvals, id)
  item.status = status
  item.approverId = userId
  item.resultNote = data.note || ''
  item.updatedAt = now()
  addAudit(userId, `approval.${status}`, 'approval', id, item.resultNote)
  persistMockState()
  return { id, status }
}

function updateReincarnation(id, patch) {
  const item = findByID(state.reincarnations, id)
  Object.assign(item, patch, { updatedAt: now() })
  persistMockState()
  return { id, status: item.status, soupStatus: item.soupStatus, quotaType: item.quotaType }
}

function issueSoup(id, data, userId) {
  const order = findByID(state.reincarnations, id)
  const inventory = findByID(state.soupInventory, Number(data.inventoryId) || 1)
  inventory.stock = Math.max(0, inventory.stock - (Number(data.dose) || 1))
  order.soupStatus = 'issued'
  state.soupRecords.unshift({
    id: state.nextSoupRecordId++,
    soulId: order.soulId,
    soulName: order.soulName,
    orderId: order.id,
    inventoryId: inventory.id,
    inventoryName: inventory.name,
    dose: Number(data.dose) || 1,
    memoryAfter: Number(data.memoryAfter) || 0,
    operatorId: userId,
    createdAt: now()
  })
  addAudit(userId, 'soup.issue', 'reincarnation', id, 'mock 发放孟婆汤')
  persistMockState()
  return { id, soupStatus: 'issued' }
}

function reviewHellSentence(id, data, userId) {
  const item = findByID(state.hellSentences, id)
  item.reviewStatus = 'reviewing'
  item.updatedAt = now()
  const approval = {
    id: state.nextApprovalId++,
    type: 'hell_review',
    targetId: id,
    title: '地狱刑期复核',
    status: 'pending',
    applicantId: userId,
    approverId: null,
    reason: data.note || '移动端发起刑期复核',
    resultNote: '',
    createdAt: now(),
    updatedAt: now()
  }
  state.approvals.unshift(approval)
  addAudit(userId, 'hell.review', 'hell_sentence', id, data.note || '')
  persistMockState()
  return { id, reviewStatus: 'reviewing' }
}

function dispatchHellFloor(id, data, userId) {
  const source = findByID(state.hellFloors, id)
  const target = findByID(state.hellFloors, Number(data.targetFloorId))
  const amount = Number(data.amount) || 0
  source.occupancy = Math.max(0, source.occupancy - amount)
  target.occupancy = Math.min(target.capacity, target.occupancy + amount)
  source.updatedAt = now()
  target.updatedAt = now()
  addAudit(userId, 'hell.dispatch', 'hell_floor', id, `向 ${target.name} 分流 ${amount} 魂`)
  persistMockState()
  return { sourceFloorId: source.id, targetFloorId: target.id, amount }
}

function updateWish(id, patch, userId, action) {
  const item = findByID(state.wishes, id)
  Object.assign(item, patch, { updatedAt: now() })
  addAudit(userId, action, 'wish', id, patch.resultNote || '')
  persistMockState()
  return { id, status: item.status }
}

function adjustSoup(id, data, userId) {
  const item = findByID(state.soupInventory, id)
  item.stock = Math.max(0, item.stock + (Number(data.delta) || 0))
  item.updatedAt = now()
  addAudit(userId, Number(data.delta) >= 0 ? 'soup.restock' : 'soup.adjust', 'soup_inventory', id, data.note || 'mock 调整孟婆汤库存')
  persistMockState()
  return { id, stock: item.stock }
}

function addAudit(actorId, action, targetType, targetId, note) {
  state.auditLogs.unshift({ id: state.nextAuditLogId++, actorId, action, targetType, targetId, note, createdAt: now() })
}

function filterStatus(items, searchParams) {
  const status = searchParams.get('status')
  if (!status || status === 'all') return items
  return items.filter(item => item.status === status)
}

function filterLifeBook(searchParams) {
  const keyword = searchParams.get('keyword') || ''
  const riskFlag = searchParams.get('riskFlag') || searchParams.get('risk') || 'all'
  return state.lifeBookRecords.filter(item => {
    const keywordOK = !keyword || item.soulName.includes(keyword)
    const riskOK = !riskFlag || riskFlag === 'all' || item.riskFlag === riskFlag
    return keywordOK && riskOK
  })
}

function filterKeyword(items, searchParams, getter) {
  const keyword = searchParams.get('keyword') || ''
  if (!keyword) return items
  return items.filter(item => getter(item).includes(keyword))
}

function findByID(items, id) {
  const item = items.find(entry => entry.id === Number(id))
  if (!item) throw new Error(`mock record not found: ${id}`)
  return item
}

function parsePath(rawPath) {
  const url = new URL(rawPath, 'http://mock.local')
  return { pathname: url.pathname, searchParams: url.searchParams }
}

function match(pathname, pattern) {
  const pathParts = pathname.split('/').filter(Boolean)
  const patternParts = pattern.split('/').filter(Boolean)
  if (pathParts.length !== patternParts.length) return false
  return patternParts.every((part, index) => part.startsWith(':') || part === pathParts[index])
}

function idOf(pathname) {
  const id = pathname.split('/').filter(Boolean).findLast(part => /^\d+$/.test(part))
  return Number(id)
}

function now() {
  return new Date().toISOString()
}

function loadMockState() {
  try {
    const raw = uni.getStorageSync(MOCK_STATE_KEY)
    if (!raw) {
      const fresh = createMockState()
      uni.setStorageSync(MOCK_STARTED_KEY, now())
      uni.setStorageSync(MOCK_STATE_KEY, JSON.stringify(fresh))
      return fresh
    }
    return JSON.parse(raw)
  } catch (error) {
    const fresh = createMockState()
    uni.setStorageSync(MOCK_STARTED_KEY, now())
    uni.setStorageSync(MOCK_STATE_KEY, JSON.stringify(fresh))
    return fresh
  }
}

function persistMockState(resetStartedAt = false) {
  if (resetStartedAt || !uni.getStorageSync(MOCK_STARTED_KEY)) {
    uni.setStorageSync(MOCK_STARTED_KEY, now())
  }
  uni.setStorageSync(MOCK_STATE_KEY, JSON.stringify(state))
}
