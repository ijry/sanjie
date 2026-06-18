export const roleText = {
  jade_emperor: '玉帝',
  yanwang: '阎王',
  judge: '判官',
  impermanence: '黑白无常',
  mengpo: '孟婆',
  heavenly_general: '天兵统领',
  city_god: '城隍',
  ghost_clerk: '鬼差',
  auditor: '审计神仙'
}

export function roleName(role) {
  return roleText[role] || role || '-'
}

