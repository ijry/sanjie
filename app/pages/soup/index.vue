<template>
  <view class="page stack dense-page">
    <view class="module-hero">
      <text class="module-hero__eyebrow">SOUP INVENTORY</text>
      <text class="module-hero__title">孟婆汤库存</text>
      <text class="module-hero__desc">汤剂库存、预警线和发放记录统一查看，避免投胎窗口断供。</text>
      <view class="module-hero__stats">
        <view class="module-hero__stat">
          <text>{{ inventory.length }}</text>
          <text>汤剂</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ lowStockCount }}</text>
          <text>低库存</text>
        </view>
        <view class="module-hero__stat">
          <text>{{ records.length }}</text>
          <text>发放</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <SectionTitle title="库存" />
      <view class="card-list">
        <view v-for="item in inventory" :key="item.id" class="work-card">
          <view class="work-card__head">
            <view class="work-card__title">
              <text>{{ item.name }}</text>
              <text class="work-card__subtitle">{{ item.formulaNote }}</text>
            </view>
            <up-button size="mini" type="primary" text="+10" @click="adjust(item, 10)" />
          </view>
          <view class="work-card__grid">
            <view class="work-card__field">
              <text class="work-card__field-label">库存</text>
              <text class="work-card__field-value">{{ item.stock }}</text>
            </view>
            <view class="work-card__field">
              <text class="work-card__field-label">预警线</text>
              <text class="work-card__field-value">{{ item.warningStock }}</text>
            </view>
          </view>
          <view v-if="item.stock <= item.warningStock" class="work-card__note">库存低于预警线，建议补货。</view>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <SectionTitle title="发放记录" />
      <EmptyState v-if="records.length === 0" text="暂无发放记录" />
      <view v-else class="mini-list">
        <view v-for="item in records" :key="item.id" class="mini-row">
          <view class="mini-row__main">
            <text class="mini-row__title">{{ item.soulName }}</text>
            <text class="mini-row__meta">{{ item.inventoryName }} · 剂量 {{ item.dose }} · 残留 {{ item.memoryAfter }}%</text>
          </view>
          <text class="record-id">#{{ item.id }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import SectionTitle from '../../components/SectionTitle.vue'
import EmptyState from '../../components/EmptyState.vue'
import { adjustSoupInventory, getSoupInventory, getSoupRecords } from '../../api/soup'
import { showError, showSuccess } from '../../utils/format'

const inventory = ref([])
const records = ref([])
const lowStockCount = computed(() => inventory.value.filter(item => Number(item.stock) <= Number(item.warningStock)).length)

async function load() {
  try {
    inventory.value = await getSoupInventory()
    records.value = await getSoupRecords()
  } catch (error) {
    showError(error)
  }
}

async function adjust(item, delta) {
  try {
    await adjustSoupInventory(item.id, { delta, note: '移动端库存调整' })
    showSuccess('已调整')
    await load()
  } catch (error) {
    showError(error)
  }
}

onShow(load)
</script>

<style scoped>
.record-id {
  color: #a93226;
  font-size: 22rpx;
  font-weight: 900;
}
</style>
