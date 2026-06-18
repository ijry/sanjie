<template>
  <view class="page stack">
    <view class="panel">
      <SectionTitle title="库存" />
      <up-cell-group :border="false">
        <up-cell v-for="item in inventory" :key="item.id" :title="item.name" :label="item.formulaNote" :value="`${item.stock}/${item.warningStock}`">
          <template #right-icon>
            <up-button size="mini" type="primary" text="+10" @click="adjust(item, 10)" />
          </template>
        </up-cell>
      </up-cell-group>
    </view>

    <view class="panel">
      <SectionTitle title="发放记录" />
      <EmptyState v-if="records.length === 0" text="暂无发放记录" />
      <up-cell v-for="item in records" :key="item.id" :title="item.soulName" :label="`${item.inventoryName} · 剂量 ${item.dose} · 残留 ${item.memoryAfter}%`" />
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import SectionTitle from '../../components/SectionTitle.vue'
import EmptyState from '../../components/EmptyState.vue'
import { adjustSoupInventory, getSoupInventory, getSoupRecords } from '../../api/soup'
import { showError, showSuccess } from '../../utils/format'

const inventory = ref([])
const records = ref([])

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

