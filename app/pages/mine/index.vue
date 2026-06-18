<template>
  <view class="page stack dense-page mine-page">
    <view class="passport-card">
      <view class="passport-card__seal">{{ avatarText }}</view>
      <view class="passport-card__main">
        <text class="passport-card__eyebrow">IMMORTAL PASSPORT</text>
        <text class="passport-card__name">{{ user.current?.nickname || '未登录' }}</text>
        <text class="passport-card__role">{{ roleName(user.current?.role) }}</text>
      </view>
      <view class="passport-card__chip">
        <text>{{ user.current?.id ? `#${user.current.id}` : '-' }}</text>
      </view>
    </view>

    <view class="identity-panel">
      <view class="identity-panel__head">
        <view>
          <text class="identity-panel__kicker">ACCESS TOKEN</text>
          <text class="identity-panel__title">仙籍令牌</text>
        </view>
        <up-button size="mini" type="primary" text="切换身份" @click="go('/pages/login/index')" />
      </view>
      <view class="token-strip">
        <view v-for="item in tokenLines" :key="item.label" class="token-line">
          <text class="token-line__label">{{ item.label }}</text>
          <text class="token-line__value">{{ item.value }}</text>
        </view>
      </view>
    </view>

    <view class="demo-panel">
      <view class="demo-panel__head">
        <view>
          <text class="demo-panel__kicker">MOCK SESSION</text>
          <text class="demo-panel__title">演示状态</text>
        </view>
        <up-button v-if="mockEnabled" size="mini" type="warning" text="重置演示" @click="confirmReset" />
      </view>
      <view class="demo-metrics">
        <view v-for="item in demoMetrics" :key="item.label" class="demo-metric">
          <text class="demo-metric__value">{{ item.value }}</text>
          <text class="demo-metric__label">{{ item.label }}</text>
        </view>
      </view>
      <text class="demo-panel__note">{{ demoNote }}</text>
    </view>

    <view class="quick-panel">
      <view class="quick-panel__head">
        <text>值守快捷入口</text>
        <text>{{ dutyHint }}</text>
      </view>
      <view class="quick-grid">
        <view v-for="item in roleShortcuts" :key="item.url" class="quick-tile" :class="`quick-tile--${item.tone}`" @click="go(item.url)">
          <text class="quick-tile__mark">{{ item.mark }}</text>
          <text class="quick-tile__name">{{ item.name }}</text>
          <text class="quick-tile__desc">{{ item.desc }}</text>
        </view>
      </view>
    </view>

    <view class="permission-panel">
      <view class="permission-panel__head">
        <text>权限模拟</text>
        <text>当前身份推荐演示动作</text>
      </view>
      <view class="permission-list">
        <view v-for="item in permissionCards" :key="item.name" class="permission-card" :class="{ 'permission-card--active': item.active }">
          <text class="permission-card__name">{{ item.name }}</text>
          <text class="permission-card__desc">{{ item.desc }}</text>
        </view>
      </view>
    </view>

    <view class="system-panel">
      <view class="system-row" @click="go('/pages/audit/index')">
        <text>天眼审计日志</text>
        <text>查看敏感操作</text>
      </view>
      <view class="system-row" @click="showApiInfo">
        <text>接口模式</text>
        <text>{{ apiMode }}</text>
      </view>
      <view class="system-row" @click="showDocsInfo">
        <text>项目文档</text>
        <text>docs-site / H5 demo</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { API_BASE_URL, MOCK_ENABLED } from '../../api/http'
import { getMockSession, resetMockSession } from '../../api/mock'
import { useUserStore } from '../../stores/user'
import { roleName } from '../../utils/auth'
import { displayTime, showError, showSuccess } from '../../utils/format'

const user = useUserStore()
const mockSession = ref(null)
const avatarText = computed(() => (user.current?.nickname || '三界').slice(0, 1))
const apiMode = computed(() => (MOCK_ENABLED ? '前端 Mock' : 'Go + SQLite'))
const mockEnabled = computed(() => MOCK_ENABLED)

const tokenLines = computed(() => [
  { label: '当前角色', value: roleName(user.current?.role) },
  { label: '账号标识', value: user.current?.username || '-' },
  { label: '接口地址', value: MOCK_ENABLED ? 'mock://sanjie' : API_BASE_URL },
  { label: '演示模式', value: MOCK_ENABLED ? '本地持久化，可一键重置' : '真实后端，SQLite 持久化' }
])

const demoMetrics = computed(() => [
  { label: '审计日志', value: mockSession.value?.auditLogCount ?? '-' },
  { label: '审批单', value: mockSession.value?.approvalCount ?? '-' },
  { label: '发汤记录', value: mockSession.value?.soupRecordCount ?? '-' },
  { label: '已处理预警', value: mockSession.value?.handledAlertCount ?? '-' }
])

const demoNote = computed(() => {
  if (!MOCK_ENABLED) return '当前连接真实 Go + SQLite 后端，数据由后端数据库保存。'
  if (!mockSession.value?.startedAt) return 'mock 状态正在读取。'
  return `当前 mock 会话始于 ${displayTime(mockSession.value.startedAt)}，刷新页面不会丢失演示操作。`
})

const dutyHint = computed(() => {
  const hints = {
    jade_emperor: '总览审批与高危命簿',
    yanwang: '看地狱负载和刑期复核',
    judge: '处理审批和命簿风险',
    impermanence: '执行勾魂任务',
    mengpo: '盯汤坊库存和投胎单',
    heavenly_general: '巡看工作台告警',
    city_god: '协助勾魂和投胎',
    auditor: '追查审计留痕'
  }
  return hints[user.current?.role] || '按推荐路线演示'
})

const roleShortcuts = computed(() => {
  const common = [
    { name: '工作台', desc: '三界总控大屏', mark: '台', tone: 'ink', url: '/pages/index/index' },
    { name: '审计', desc: '天眼留痕', mark: '眼', tone: 'blue', url: '/pages/audit/index' }
  ]
  const byRole = {
    jade_emperor: [
      { name: '审批中心', desc: '处理高危申请', mark: '审', tone: 'amber', url: '/pages/approval/index' },
      { name: '生死簿', desc: '检查关系户阳寿', mark: '命', tone: 'paper', url: '/pages/life-book/search' }
    ],
    yanwang: [
      { name: '十八层', desc: '楼层负载调度', mark: '狱', tone: 'red', url: '/pages/hell/index' },
      { name: '审批中心', desc: '刑期复核', mark: '审', tone: 'amber', url: '/pages/approval/index' }
    ],
    judge: [
      { name: '审批中心', desc: '判官案牍', mark: '审', tone: 'amber', url: '/pages/approval/index' },
      { name: '生死簿', desc: '改寿留痕', mark: '命', tone: 'paper', url: '/pages/life-book/search' }
    ],
    impermanence: [
      { name: '勾魂任务', desc: '夜巡任务池', mark: '魂', tone: 'cyan', url: '/pages/capture/index' },
      { name: '投胎队列', desc: '移交流程', mark: '桥', tone: 'teal', url: '/pages/reincarnation/index' }
    ],
    mengpo: [
      { name: '孟婆汤', desc: '汤坊库存台', mark: '汤', tone: 'teal', url: '/pages/soup/index' },
      { name: '投胎队列', desc: '发汤归档', mark: '桥', tone: 'teal', url: '/pages/reincarnation/index' }
    ],
    auditor: [
      { name: '审计日志', desc: '追溯风险操作', mark: '眼', tone: 'blue', url: '/pages/audit/index' },
      { name: '生死簿', desc: '冻结高危记录', mark: '命', tone: 'paper', url: '/pages/life-book/search' }
    ]
  }
  return [...(byRole[user.current?.role] || byRole.jade_emperor), ...common].slice(0, 4)
})

const permissionCards = computed(() => {
  const role = user.current?.role
  return [
    { name: '勾魂外勤', desc: '开始、完成、异常上报', active: ['impermanence', 'city_god', 'jade_emperor'].includes(role) },
    { name: '判官审批', desc: '同意、驳回、转交', active: ['judge', 'yanwang', 'jade_emperor'].includes(role) },
    { name: '汤坊调度', desc: '补货、扣减、发汤', active: ['mengpo', 'jade_emperor'].includes(role) },
    { name: '审计追踪', desc: '查看敏感操作链路', active: ['auditor', 'jade_emperor', 'judge'].includes(role) }
  ]
})

function go(url) {
  if (['/pages/index/index', '/pages/capture/index', '/pages/reincarnation/index', '/pages/approval/index', '/pages/mine/index'].includes(url)) {
    uni.switchTab({ url })
    return
  }
  uni.navigateTo({ url })
}

function showApiInfo() {
  uni.showModal({
    title: '接口模式',
    content: MOCK_ENABLED ? '当前使用前端 mock 数据，页面刷新后恢复初始状态。' : `当前连接后端：${API_BASE_URL}`,
    showCancel: false
  })
}

function showDocsInfo() {
  uni.showModal({
    title: '文档站',
    content: '项目文档和嵌入式 H5 演示由 docs-site 发布，可从 README 中打开 GitHub Pages。',
    showCancel: false
  })
}

async function loadMockSession() {
  if (!MOCK_ENABLED) return
  try {
    mockSession.value = await getMockSession()
  } catch (error) {
    showError(error)
  }
}

function confirmReset() {
  uni.showModal({
    title: '重置演示数据',
    content: '将恢复 mock 初始数据、清空本轮操作状态，并切回玉帝身份。',
    success: async res => {
      if (!res.confirm) return
      try {
        await resetMockSession()
        uni.setStorageSync('sanjie_demo_guide_seen', '')
        await user.fetchMe()
        await loadMockSession()
        showSuccess('演示已重置')
      } catch (error) {
        showError(error)
      }
    }
  })
}

onShow(async () => {
  await user.fetchMe()
  await loadMockSession()
})
</script>

<style scoped>
.mine-page {
  background:
    radial-gradient(circle at 14% 0%, rgba(169, 50, 38, 0.16), transparent 25%),
    radial-gradient(circle at 88% 8%, rgba(217, 119, 6, 0.18), transparent 25%),
    linear-gradient(180deg, #fff7ed 0%, #f7f1e3 42%, #f7f1e3 100%);
}

.passport-card {
  position: relative;
  overflow: hidden;
  display: grid;
  grid-template-columns: 116rpx minmax(0, 1fr) 86rpx;
  gap: 18rpx;
  align-items: center;
  padding: 28rpx;
  color: #fffaf0;
  border: 1rpx solid rgba(251, 191, 36, 0.26);
  border-radius: 30rpx;
  background:
    linear-gradient(90deg, rgba(251, 191, 36, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #111827, #7f1d1d 62%, #b45309);
  background-size: 34rpx 34rpx;
  box-shadow: 0 18rpx 42rpx rgba(49, 28, 13, 0.2);
}

.passport-card::after {
  content: "";
  position: absolute;
  right: -60rpx;
  top: -70rpx;
  width: 210rpx;
  height: 210rpx;
  border: 2rpx dashed rgba(255, 250, 240, 0.18);
  border-radius: 50%;
}

.passport-card__seal {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 108rpx;
  height: 108rpx;
  border: 6rpx double rgba(255, 250, 240, 0.42);
  border-radius: 50%;
  color: #fef3c7;
  font-size: 48rpx;
  font-weight: 900;
}

.passport-card__main {
  min-width: 0;
}

.passport-card__eyebrow,
.passport-card__name,
.passport-card__role {
  display: block;
}

.passport-card__eyebrow {
  color: #fbbf24;
  font-size: 19rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.passport-card__name {
  margin-top: 6rpx;
  font-size: 46rpx;
  font-weight: 900;
  line-height: 1.12;
}

.passport-card__role {
  margin-top: 8rpx;
  color: rgba(255, 250, 240, 0.72);
  font-size: 24rpx;
}

.passport-card__chip {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 66rpx;
  border: 1rpx solid rgba(255, 250, 240, 0.18);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.1);
  font-size: 22rpx;
  font-weight: 900;
}

.identity-panel,
.demo-panel,
.quick-panel,
.permission-panel,
.system-panel {
  padding: 22rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.16);
  border-radius: 24rpx;
  background: rgba(255, 250, 240, 0.94);
  box-shadow: 0 12rpx 30rpx rgba(66, 41, 20, 0.08);
}

.demo-panel {
  background:
    linear-gradient(90deg, rgba(37, 99, 235, 0.06) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #eff6ff, #fffaf0);
  background-size: 30rpx 30rpx;
}

.identity-panel__head,
.demo-panel__head,
.quick-panel__head,
.permission-panel__head {
  display: flex;
  justify-content: space-between;
  gap: 14rpx;
}

.identity-panel__kicker {
  display: block;
  color: #a93226;
  font-size: 19rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.demo-panel__kicker {
  display: block;
  color: #1d4ed8;
  font-size: 19rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.identity-panel__title,
.demo-panel__title,
.quick-panel__head text:first-child,
.permission-panel__head text:first-child {
  display: block;
  margin-top: 4rpx;
  color: #25211b;
  font-size: 32rpx;
  font-weight: 900;
}

.demo-metrics {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10rpx;
  margin-top: 18rpx;
}

.demo-metric {
  min-width: 0;
  padding: 14rpx 8rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.12);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.78);
}

.demo-metric__value,
.demo-metric__label {
  display: block;
  text-align: center;
}

.demo-metric__value {
  color: #1e3a8a;
  font-size: 28rpx;
  font-weight: 900;
}

.demo-metric__label {
  overflow: hidden;
  margin-top: 5rpx;
  color: #475569;
  font-size: 19rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.demo-panel__note {
  display: block;
  margin-top: 12rpx;
  color: #1e3a8a;
  font-size: 22rpx;
  font-weight: 800;
  line-height: 1.42;
}

.quick-panel__head text:last-child,
.permission-panel__head text:last-child {
  color: #8c806d;
  font-size: 22rpx;
  font-weight: 800;
}

.token-strip {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-top: 18rpx;
}

.token-line {
  display: grid;
  grid-template-columns: 150rpx minmax(0, 1fr);
  gap: 12rpx;
  padding: 14rpx;
  border: 1rpx dashed rgba(169, 50, 38, 0.18);
  border-radius: 16rpx;
  background: #fff7ed;
}

.token-line__label {
  color: #9a3412;
  font-size: 22rpx;
  font-weight: 900;
}

.token-line__value {
  overflow: hidden;
  color: #25211b;
  font-size: 22rpx;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.quick-grid,
.permission-list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
  margin-top: 18rpx;
}

.quick-tile {
  min-width: 0;
  padding: 18rpx;
  border-radius: 20rpx;
  color: #fffaf0;
  background: #1f2937;
}

.quick-tile--amber { background: linear-gradient(135deg, #b45309, #78350f); }
.quick-tile--paper { background: linear-gradient(135deg, #92400e, #3f2a12); }
.quick-tile--red { background: linear-gradient(135deg, #a93226, #450a0a); }
.quick-tile--cyan { background: linear-gradient(135deg, #0284c7, #0f172a); }
.quick-tile--teal { background: linear-gradient(135deg, #0f766e, #134e4a); }
.quick-tile--blue { background: linear-gradient(135deg, #1d4ed8, #0f172a); }
.quick-tile--ink { background: linear-gradient(135deg, #25211b, #0f0d0a); }

.quick-tile__mark {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 54rpx;
  height: 54rpx;
  border: 1rpx solid rgba(255, 250, 240, 0.24);
  border-radius: 16rpx;
  font-weight: 900;
}

.quick-tile__name,
.quick-tile__desc {
  display: block;
}

.quick-tile__name {
  margin-top: 16rpx;
  font-size: 28rpx;
  font-weight: 900;
}

.quick-tile__desc {
  margin-top: 6rpx;
  color: rgba(255, 250, 240, 0.7);
  font-size: 21rpx;
}

.permission-card {
  padding: 16rpx;
  border: 1rpx solid rgba(120, 53, 15, 0.14);
  border-radius: 18rpx;
  background: #f8fafc;
  opacity: 0.58;
}

.permission-card--active {
  background: #fff7ed;
  opacity: 1;
}

.permission-card__name,
.permission-card__desc {
  display: block;
}

.permission-card__name {
  color: #25211b;
  font-size: 25rpx;
  font-weight: 900;
}

.permission-card__desc {
  margin-top: 6rpx;
  color: #756f63;
  font-size: 21rpx;
  line-height: 1.35;
}

.system-panel {
  padding: 0;
  overflow: hidden;
}

.system-row {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  padding: 20rpx 22rpx;
  border-bottom: 1rpx solid rgba(120, 53, 15, 0.12);
}

.system-row:last-child {
  border-bottom: 0;
}

.system-row text:first-child {
  color: #25211b;
  font-size: 26rpx;
  font-weight: 900;
}

.system-row text:last-child {
  overflow: hidden;
  color: #756f63;
  font-size: 22rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
