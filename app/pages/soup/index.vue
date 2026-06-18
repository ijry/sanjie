<template>
  <view class="page stack dense-page soup-page">
    <view class="soup-hero">
      <view class="soup-cauldron">
        <view class="soup-cauldron__steam soup-cauldron__steam--one"></view>
        <view class="soup-cauldron__steam soup-cauldron__steam--two"></view>
        <view class="soup-cauldron__pot">
          <text>{{ lowStockCount }}</text>
          <text>低库存</text>
        </view>
      </view>
      <view class="soup-hero__main">
        <text class="soup-hero__eyebrow">MENGPO BREWERY</text>
        <text class="soup-hero__title">汤坊库存台</text>
        <text class="soup-hero__desc">库存、预警、补货和发放记录放在同一张值班台，避免奈何桥断汤。</text>
      </view>
      <view class="soup-hero__stats">
        <view class="soup-stat">
          <text>{{ inventory.length }}</text>
          <text>汤剂种类</text>
        </view>
        <view class="soup-stat">
          <text>{{ totalStock }}</text>
          <text>总库存</text>
        </view>
        <view class="soup-stat">
          <text>{{ records.length }}</text>
          <text>发放记录</text>
        </view>
      </view>
    </view>

    <view v-if="selectedItem" class="brew-console">
      <view class="brew-console__head">
        <view>
          <text class="brew-console__kicker">BREW CONTROL</text>
          <text class="brew-console__title">{{ selectedItem.name }}</text>
        </view>
        <text class="brew-console__stock">库存 {{ selectedItem.stock }}</text>
      </view>
      <view class="brew-console__form">
        <view class="brew-console__field">
          <text>调整数量</text>
          <up-input v-model="deltaText" type="number" placeholder="正数补货，负数扣减" border="none" />
        </view>
        <view class="brew-console__field">
          <text>备注</text>
          <up-input v-model="noteText" placeholder="例如：桥头窗口加急补货" border="none" />
        </view>
      </view>
      <view class="brew-console__actions">
        <up-button size="small" plain text="-5" @click="deltaText = '-5'" />
        <up-button size="small" plain text="+10" @click="deltaText = '10'" />
        <up-button size="small" type="primary" text="提交调整" @click="applyAdjust" />
      </view>
    </view>

    <view class="stock-board">
      <view class="stock-board__head">
        <view>
          <text class="stock-board__kicker">INVENTORY</text>
          <text class="stock-board__title">库存雷达</text>
        </view>
        <text class="stock-board__hint">{{ dangerNames }}</text>
      </view>

      <view class="stock-list">
        <view
          v-for="item in inventory"
          :key="item.id"
          class="soup-stock-card"
          :class="{ 'soup-stock-card--low': item.stock <= item.warningStock }"
          @click="selectAdjust(item)"
        >
          <view class="soup-stock-card__head">
            <view>
              <text class="soup-stock-card__name">{{ item.name }}</text>
              <text class="soup-stock-card__formula">{{ item.formulaNote }}</text>
            </view>
            <text class="soup-stock-card__badge">{{ item.stock <= item.warningStock ? '补' : '稳' }}</text>
          </view>
          <view class="soup-stock-card__gauge">
            <view class="soup-stock-card__gauge-inner" :style="{ width: stockPercent(item) }"></view>
          </view>
          <view class="soup-stock-card__meta">
            <text>库存 {{ item.stock }}</text>
            <text>预警 {{ item.warningStock }}</text>
            <text>{{ item.stock <= item.warningStock ? '需要补货' : '供应稳定' }}</text>
          </view>
          <view class="soup-stock-card__actions">
            <up-button size="mini" type="primary" text="调汤" @click.stop="selectAdjust(item)" />
            <up-button size="mini" plain text="+10" @click.stop="quickAdjust(item, 10)" />
          </view>
        </view>
      </view>
    </view>

    <view class="record-board">
      <view class="record-board__head">
        <view>
          <text class="record-board__kicker">ISSUE LEDGER</text>
          <text class="record-board__title">发放流水</text>
        </view>
        <text class="record-board__hint">记忆残留追踪</text>
      </view>
      <EmptyState v-if="records.length === 0" text="暂无发放记录" />
      <view v-else class="soup-ledger">
        <view v-for="item in records" :key="item.id" class="soup-ledger__row">
          <text class="soup-ledger__mark">#{{ item.id }}</text>
          <view class="soup-ledger__main">
            <text class="soup-ledger__title">{{ item.soulName }}</text>
            <text class="soup-ledger__meta">{{ item.inventoryName }} · 剂量 {{ item.dose }} · 残留 {{ item.memoryAfter }}%</text>
          </view>
          <text class="soup-ledger__time">{{ displayTime(item.createdAt) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { adjustSoupInventory, getSoupInventory, getSoupRecords } from '../../api/soup'
import { displayTime, showError, showSuccess } from '../../utils/format'

const inventory = ref([])
const records = ref([])
const selectedItem = ref(null)
const deltaText = ref('10')
const noteText = ref('移动端汤坊库存调整')

const lowStockCount = computed(() => inventory.value.filter(item => Number(item.stock) <= Number(item.warningStock)).length)
const totalStock = computed(() => inventory.value.reduce((total, item) => total + Number(item.stock || 0), 0))
const dangerNames = computed(() => {
  const names = inventory.value.filter(item => Number(item.stock) <= Number(item.warningStock)).map(item => item.name)
  return names.length ? `预警：${names.join('、')}` : '库存稳定'
})

async function load() {
  try {
    inventory.value = await getSoupInventory()
    records.value = await getSoupRecords()
    if (selectedItem.value) {
      selectedItem.value = inventory.value.find(item => item.id === selectedItem.value.id) || null
    } else {
      selectedItem.value = inventory.value.find(item => item.stock <= item.warningStock) || inventory.value[0] || null
    }
  } catch (error) {
    showError(error)
  }
}

function selectAdjust(item) {
  selectedItem.value = item
  deltaText.value = item.stock <= item.warningStock ? '10' : '-5'
  noteText.value = item.stock <= item.warningStock ? '库存低于预警线，移动端补货' : '移动端抽检扣减'
}

async function quickAdjust(item, delta) {
  selectedItem.value = item
  deltaText.value = String(delta)
  await applyAdjust()
}

async function applyAdjust() {
  if (!selectedItem.value) return
  const delta = Number(deltaText.value)
  if (!Number.isFinite(delta) || delta === 0) {
    showError(new Error('请输入非 0 调整数量'))
    return
  }
  try {
    await adjustSoupInventory(selectedItem.value.id, { delta, note: noteText.value || '移动端库存调整' })
    showSuccess(delta > 0 ? '已补货' : '已扣减')
    await load()
  } catch (error) {
    showError(error)
  }
}

function stockPercent(item) {
  const base = Math.max(Number(item.warningStock || 0) * 3, Number(item.stock || 0), 1)
  return `${Math.min(100, Math.round((Number(item.stock || 0) / base) * 100))}%`
}

onShow(load)
</script>

<style scoped>
.soup-page {
  background:
    radial-gradient(circle at 50% 0%, rgba(20, 184, 166, 0.16), transparent 25%),
    linear-gradient(180deg, #10201d 0%, #1f3a2f 30%, #f7f1e3 30%, #f7f1e3 100%);
}

.soup-hero {
  display: grid;
  grid-template-columns: 150rpx minmax(0, 1fr);
  gap: 18rpx;
  padding: 26rpx;
  color: #ecfdf5;
  border: 1rpx solid rgba(45, 212, 191, 0.22);
  border-radius: 32rpx;
  background:
    radial-gradient(circle at 82% 0%, rgba(45, 212, 191, 0.2), transparent 28%),
    linear-gradient(135deg, #052e2b, #164e3f 58%, #3f2a12);
  box-shadow: 0 20rpx 46rpx rgba(15, 78, 63, 0.2);
}

.soup-cauldron {
  position: relative;
  width: 140rpx;
  height: 142rpx;
}

.soup-cauldron__steam {
  position: absolute;
  top: 0;
  width: 20rpx;
  height: 58rpx;
  border-radius: 999rpx;
  background: rgba(204, 251, 241, 0.22);
}

.soup-cauldron__steam--one {
  left: 38rpx;
}

.soup-cauldron__steam--two {
  right: 38rpx;
  top: 14rpx;
}

.soup-cauldron__pot {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 96rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 5rpx solid rgba(204, 251, 241, 0.28);
  border-radius: 28rpx 28rpx 42rpx 42rpx;
  background: linear-gradient(180deg, #14b8a6, #0f766e);
}

.soup-cauldron__pot text:first-child {
  font-size: 36rpx;
  font-weight: 900;
  line-height: 1;
}

.soup-cauldron__pot text:last-child {
  margin-top: 5rpx;
  font-size: 20rpx;
  font-weight: 900;
}

.soup-hero__main {
  min-width: 0;
}

.soup-hero__eyebrow,
.soup-hero__title,
.soup-hero__desc {
  display: block;
}

.soup-hero__eyebrow {
  color: #5eead4;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.soup-hero__title {
  margin-top: 6rpx;
  font-size: 44rpx;
  font-weight: 900;
}

.soup-hero__desc {
  margin-top: 10rpx;
  color: rgba(236, 253, 245, 0.72);
  font-size: 24rpx;
  line-height: 1.45;
}

.soup-hero__stats {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10rpx;
}

.soup-stat {
  padding: 14rpx 12rpx;
  border: 1rpx solid rgba(204, 251, 241, 0.16);
  border-radius: 16rpx;
  background: rgba(4, 47, 46, 0.42);
}

.soup-stat text:first-child {
  display: block;
  font-size: 30rpx;
  font-weight: 900;
}

.soup-stat text:last-child {
  display: block;
  margin-top: 4rpx;
  color: rgba(236, 253, 245, 0.66);
  font-size: 20rpx;
}

.brew-console,
.stock-board,
.record-board {
  padding: 22rpx;
  border: 1rpx solid rgba(15, 118, 110, 0.18);
  border-radius: 24rpx;
  background: rgba(255, 250, 240, 0.94);
  box-shadow: 0 12rpx 30rpx rgba(15, 78, 63, 0.08);
}

.brew-console {
  background:
    linear-gradient(90deg, rgba(15, 118, 110, 0.08) 1rpx, transparent 1rpx),
    linear-gradient(135deg, #f0fdfa, #fffaf0);
  background-size: 30rpx 30rpx;
}

.brew-console__head,
.stock-board__head,
.record-board__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14rpx;
}

.brew-console__kicker,
.stock-board__kicker,
.record-board__kicker {
  display: block;
  color: #0f766e;
  font-size: 19rpx;
  font-weight: 900;
  letter-spacing: 2rpx;
}

.brew-console__title,
.stock-board__title,
.record-board__title {
  display: block;
  margin-top: 4rpx;
  color: #134e4a;
  font-size: 32rpx;
  font-weight: 900;
}

.brew-console__stock,
.stock-board__hint,
.record-board__hint {
  color: #64748b;
  font-size: 22rpx;
  font-weight: 800;
}

.brew-console__form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
  margin-top: 18rpx;
}

.brew-console__field {
  min-width: 0;
  padding: 14rpx;
  border: 1rpx solid rgba(15, 118, 110, 0.16);
  border-radius: 16rpx;
  background: rgba(255, 250, 240, 0.74);
}

.brew-console__field text {
  display: block;
  color: #0f766e;
  font-size: 21rpx;
  font-weight: 900;
}

.brew-console__field :deep(.u-input) {
  margin-top: 8rpx;
  padding: 0 !important;
  background: transparent !important;
}

.brew-console__field :deep(.uni-input-input),
.brew-console__field :deep(.u-input__content__field-wrapper__field) {
  color: #134e4a;
  font-size: 25rpx;
  font-weight: 800;
}

.brew-console__actions {
  display: grid;
  grid-template-columns: 1fr 1fr 2fr;
  gap: 10rpx;
  margin-top: 16rpx;
}

.stock-list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  margin-top: 18rpx;
}

.soup-stock-card {
  padding: 20rpx;
  border: 1rpx solid rgba(15, 118, 110, 0.16);
  border-radius: 22rpx;
  background: linear-gradient(135deg, #f0fdfa, #fffaf0);
}

.soup-stock-card--low {
  border-color: rgba(185, 28, 28, 0.28);
  background: linear-gradient(135deg, #fef2f2, #fffaf0);
}

.soup-stock-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14rpx;
}

.soup-stock-card__name,
.soup-stock-card__formula {
  display: block;
}

.soup-stock-card__name {
  color: #134e4a;
  font-size: 31rpx;
  font-weight: 900;
}

.soup-stock-card__formula {
  margin-top: 6rpx;
  color: #64748b;
  font-size: 22rpx;
  line-height: 1.4;
}

.soup-stock-card__badge {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 58rpx;
  height: 58rpx;
  color: #fffaf0;
  border-radius: 50%;
  background: #0f766e;
  font-size: 26rpx;
  font-weight: 900;
}

.soup-stock-card--low .soup-stock-card__badge {
  background: #a93226;
}

.soup-stock-card__gauge {
  height: 18rpx;
  margin-top: 18rpx;
  overflow: hidden;
  border-radius: 999rpx;
  background: rgba(15, 118, 110, 0.12);
}

.soup-stock-card__gauge-inner {
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #14b8a6, #0f766e);
}

.soup-stock-card--low .soup-stock-card__gauge-inner {
  background: linear-gradient(90deg, #f97316, #a93226);
}

.soup-stock-card__meta {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8rpx;
  margin-top: 14rpx;
}

.soup-stock-card__meta text {
  overflow: hidden;
  padding: 10rpx;
  border-radius: 12rpx;
  background: rgba(255, 250, 240, 0.72);
  color: #0f766e;
  font-size: 20rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: center;
}

.soup-stock-card__actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10rpx;
  margin-top: 14rpx;
}

.soup-ledger {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-top: 16rpx;
}

.soup-ledger__row {
  display: grid;
  grid-template-columns: 58rpx minmax(0, 1fr) 130rpx;
  gap: 12rpx;
  align-items: center;
  padding: 16rpx;
  border: 1rpx dashed rgba(15, 118, 110, 0.2);
  border-radius: 16rpx;
  background: #f0fdfa;
}

.soup-ledger__mark {
  color: #0f766e;
  font-size: 22rpx;
  font-weight: 900;
}

.soup-ledger__main {
  min-width: 0;
}

.soup-ledger__title,
.soup-ledger__meta,
.soup-ledger__time {
  display: block;
}

.soup-ledger__title {
  overflow: hidden;
  color: #134e4a;
  font-size: 26rpx;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.soup-ledger__meta,
.soup-ledger__time {
  overflow: hidden;
  margin-top: 4rpx;
  color: #64748b;
  font-size: 20rpx;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
