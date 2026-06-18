export function createMockState() {
  const now = '2026-06-18T02:00:00+08:00'
  return {
    currentUserId: 1,
    users: [
      { id: 1, username: 'jade', nickname: '玉帝', role: 'jade_emperor', avatar: '' },
      { id: 2, username: 'yanwang', nickname: '阎王', role: 'yanwang', avatar: '' },
      { id: 3, username: 'judge', nickname: '判官', role: 'judge', avatar: '' },
      { id: 4, username: 'black', nickname: '黑无常', role: 'impermanence', avatar: '' },
      { id: 5, username: 'mengpo', nickname: '孟婆', role: 'mengpo', avatar: '' }
    ],
    alerts: [
      { id: 1, level: 'critical', title: '生死簿批量划名风险', content: '检测到同一时辰 3 次寿命变更申请', handled: false, createdAt: now },
      { id: 2, level: 'warning', title: '奈何桥排队拥堵', content: '投胎队列已超过阈值，建议增开普通胎位', handled: false, createdAt: now },
      { id: 3, level: 'info', title: '孟婆汤库存提醒', content: '低糖无梦款接近预警库存', handled: false, createdAt: now }
    ],
    souls: [
      { id: 1, name: '张三丰外卖版', birthInfo: '甲子年三月初三', deathInfo: '夜跑误入鬼门关', status: 'pending_capture', meritScore: 72, karmaScore: 18, memoryResidue: 12, relationshipRisk: 'low' },
      { id: 2, name: '李小卷', birthInfo: '乙丑年七月十五', deathInfo: '熬夜写奏折', status: 'reincarnation_waiting', meritScore: 91, karmaScore: 6, memoryResidue: 25, relationshipRisk: 'medium' },
      { id: 3, name: '王大胆', birthInfo: '丙寅年正月初九', deathInfo: '十八层体验券过期', status: 'hell', meritScore: 21, karmaScore: 88, memoryResidue: 4, relationshipRisk: 'high' }
    ],
    captureTasks: [
      { id: 1, soulId: 1, soulName: '张三丰外卖版', assigneeId: 4, assigneeName: '黑无常', status: 'pending', location: '人间西市烧烤摊', scheduledTime: now, actualTime: '', exceptionType: '', exceptionNote: '' },
      { id: 2, soulId: 3, soulName: '王大胆', assigneeId: 4, assigneeName: '黑无常', status: 'running', location: '忘川河边自拍点', scheduledTime: now, actualTime: '', exceptionType: '', exceptionNote: '' },
      { id: 3, soulId: 2, soulName: '李小卷', assigneeId: 4, assigneeName: '黑无常', status: 'exception', location: '天庭加班区', scheduledTime: now, actualTime: '', exceptionType: 'not_found', exceptionNote: '疑似临时飞升' }
    ],
    lifeBookRecords: [
      { id: 1, soulId: 1, soulName: '张三丰外卖版', expectedLifespan: 88, actualLifespan: 73, deathReason: '夜跑误入鬼门关', fateLevel: '小富即安', calamityCount: 3, locked: false, riskFlag: 'warning', updatedAt: now },
      { id: 2, soulId: 2, soulName: '李小卷', expectedLifespan: 99, actualLifespan: 99, deathReason: '寿终正寝', fateLevel: '天选打工魂', calamityCount: 1, locked: false, riskFlag: 'normal', updatedAt: now },
      { id: 3, soulId: 3, soulName: '王大胆', expectedLifespan: 66, actualLifespan: 45, deathReason: '主动试胆', fateLevel: '困难模式', calamityCount: 9, locked: true, riskFlag: 'critical', updatedAt: now }
    ],
    approvals: [
      { id: 1, type: 'lifespan_change', targetId: 1, title: '寿命变更申请', status: 'pending', applicantId: 3, approverId: null, reason: '救人一命，建议补阳寿', resultNote: '', createdAt: now, updatedAt: now },
      { id: 2, type: 'hell_review', targetId: 1, title: '地狱刑期复核', status: 'transferred', applicantId: 2, approverId: null, reason: '疑似判官手滑，多判三百年', resultNote: '', createdAt: now, updatedAt: now }
    ],
    reincarnations: [
      { id: 1, soulId: 2, soulName: '李小卷', targetRealm: '人间', priority: 8, soupStatus: 'pending', quotaType: '', status: 'queued', reviewNote: '', createdAt: now, updatedAt: now },
      { id: 2, soulId: 1, soulName: '张三丰外卖版', targetRealm: '猫界', priority: 6, soupStatus: 'issued', quotaType: '普通胎', status: 'pending_soup', reviewNote: '建议保留方向感', createdAt: now, updatedAt: now }
    ],
    hellFloors: [
      { id: 1, floorNo: 1, name: '拔舌地狱', capacity: 300, occupancy: 286, equipmentStatus: '舌钳紧张', loadLevel: 'critical', updatedAt: now },
      { id: 2, floorNo: 2, name: '剪刀地狱', capacity: 260, occupancy: 120, equipmentStatus: '剪刀已磨', loadLevel: 'normal', updatedAt: now },
      { id: 3, floorNo: 18, name: '终极加班地狱', capacity: 80, occupancy: 79, equipmentStatus: '工位不足', loadLevel: 'warning', updatedAt: now }
    ],
    hellSentences: [
      { id: 1, soulId: 3, soulName: '王大胆', floorId: 1, floorName: '拔舌地狱', crimeType: '造谣地府 KPI', sentenceDays: 365, painLevel: 7, reviewStatus: 'none', equipmentId: 'BT-001', createdAt: now, updatedAt: now },
      { id: 2, soulId: 1, soulName: '张三丰外卖版', floorId: 3, floorName: '终极加班地狱', crimeType: '迟到早退', sentenceDays: 18, painLevel: 4, reviewStatus: 'reviewing', equipmentId: 'OT-996', createdAt: now, updatedAt: now }
    ],
    wishes: [
      { id: 1, title: '一夜暴富', wishType: 'wealth', requesterName: '赵小财', incenseAmount: 12, meritScore: 3, status: 'pending', assignedDeity: '', resultNote: '', createdAt: now, updatedAt: now },
      { id: 2, title: '考试全会', wishType: 'study', requesterName: '孙卷王', incenseAmount: 6, meritScore: 18, status: 'routed', assignedDeity: '文昌帝君', resultNote: '', createdAt: now, updatedAt: now }
    ],
    soupInventory: [
      { id: 1, name: '经典忘忧汤', stock: 120, warningStock: 30, formulaNote: '标准版，适合大多数投胎场景', updatedAt: now },
      { id: 2, name: '低糖无梦汤', stock: 18, warningStock: 20, formulaNote: '控糖魂魄专用', updatedAt: now }
    ],
    soupRecords: [
      { id: 1, soulId: 2, soulName: '李小卷', orderId: 2, inventoryId: 1, inventoryName: '经典忘忧汤', dose: 1, memoryAfter: 0, operatorId: 5, createdAt: now }
    ],
    auditLogs: [
      { id: 1, actorId: 3, action: 'life_book.freeze', targetType: 'life_book', targetId: 3, note: '高危寿命记录冻结', createdAt: now },
      { id: 2, actorId: 5, action: 'soup.issue', targetType: 'reincarnation', targetId: 2, note: '孟婆汤已发放', createdAt: now }
    ],
    nextApprovalId: 3,
    nextAuditLogId: 3,
    nextSoupRecordId: 2
  }
}
