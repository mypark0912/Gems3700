<template>
  <div class="card-wrap">
    <div class="card-header">
      <h3 class="card-title meter-accent-orange">{{ t('dashboard.pq.singletitle') }}</h3>
    </div>
    <div class="card-body">
      <!-- 상단: 불평형률 수평 바 -->
      <div class="section-group">
        <span class="section-title">{{ t('dashboard.pq.unbalance') }}</span>
        <div class="hbar-item">
          <span class="hbar-label">{{ t('dashboard.meter.voltage') }}</span>
          <div class="hbar-track">
            <div class="hbar-fill" :style="{ width: ubalVWidthPct + '%', backgroundColor: ubalVColor }"></div>
          </div>
          <span class="hbar-value" :style="{ color: ubalVColor }">{{ ubalVoltage.toFixed(1) }}%</span>
        </div>
        <div class="hbar-item">
          <span class="hbar-label">{{ t('dashboard.meter.current') }}</span>
          <div class="hbar-track">
            <div class="hbar-fill" :style="{ width: ubalIWidthPct + '%', backgroundColor: ubalIColor }"></div>
          </div>
          <span class="hbar-value" :style="{ color: ubalIColor }">{{ ubalCurrent.toFixed(1) }}%</span>
        </div>
      </div>

      <!-- 하단: 고조파 왜곡률 수직 바 차트 -->
      <div class="section-group mt-4">
        <span class="section-title">{{ t('pq.tabs.harmonics') }}</span>
        <div class="vbar-chart">
          <div v-for="item in thdItems" :key="item.label" class="vbar-item">
            <span class="vbar-value" :style="{ color: item.color }">{{ item.value }}%</span>
            <div class="vbar-track">
              <div class="vbar-fill" :style="{ height: item.heightPct + '%', backgroundColor: item.color }"></div>
            </div>
            <span class="vbar-label">{{ item.label }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRealtimeStore } from '@/store/realtime'

export default {
  name: 'DesignCard_PQ',
  setup() {
    const { t } = useI18n()
    const store = useRealtimeStore()
    const data2 = computed(() => store.getChannelData('Main') || {})

    const ubalVoltage = computed(() => parseFloat(data2.value.Ubal1) || 0)
    const ubalCurrent = computed(() => parseFloat(data2.value.Ibal1) || 0)

    const getUnbalColor = (v) => {
      if (v >= 3) return '#ef4444'
      if (v >= 2) return '#f59e0b'
      if (v >= 1) return '#eab308'
      return '#22c55e'
    }

    const ubalVColor = computed(() => getUnbalColor(ubalVoltage.value))
    const ubalIColor = computed(() => getUnbalColor(ubalCurrent.value))
    const ubalVWidthPct = computed(() => Math.min((ubalVoltage.value / 5) * 100, 100))
    const ubalIWidthPct = computed(() => Math.min((ubalCurrent.value / 5) * 100, 100))

    const thdItems = computed(() => {
      const items = [
        { label: 'THD-U', key: 'thdu total', color: '#8b5cf6' },
        { label: 'THD-I', key: 'thdi total', color: '#60a5fa' },
        { label: 'TDD-I', key: 'tddi total', color: '#2dd4bf' },
      ]
      const values = items.map(it => parseFloat(data2.value[it.key]) || 0)
      const maxVal = Math.max(...values, 10)
      return items.map((it, i) => ({
        label: it.label,
        value: values[i].toFixed(1),
        color: it.color,
        heightPct: Math.max((values[i] / maxVal) * 100, 5),
      }))
    })

    return {
      t, ubalVoltage, ubalCurrent, thdItems,
      ubalVColor, ubalIColor, ubalVWidthPct, ubalIWidthPct,
    }
  },
}
</script>

<style scoped>
.card-wrap {
  @apply col-span-full sm:col-span-6 xl:col-span-4;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply overflow-hidden;
}
.card-header {
  @apply flex justify-between items-center px-4 py-2.5;
}
.card-title {
  @apply text-base font-bold text-gray-800 dark:text-white flex items-center gap-2;
}
.card-title::before {
  content: '';
  @apply w-1 h-4 rounded-full inline-block flex-shrink-0;
}
.meter-accent-orange::before {
  @apply bg-orange-500;
}
.card-body {
  @apply px-4 py-3;
}

/* Section */
.section-group {
  @apply space-y-1.5;
}
.section-title {
  @apply text-sm font-bold text-gray-700 dark:text-gray-300 block mb-1.5;
}

/* Horizontal bar */
.hbar-item {
  @apply flex items-center gap-2;
}
.hbar-label {
  @apply text-sm text-gray-600 dark:text-gray-400 flex-shrink-0 whitespace-nowrap;
  min-width: 60px;
}
.hbar-track {
  @apply h-3 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden;
  width: 300px;
  flex-shrink: 0;
}
.hbar-fill {
  @apply h-full rounded-full transition-all duration-500;
  min-width: 4px;
}
.hbar-value {
  @apply text-sm font-bold tabular-nums text-right flex-shrink-0 whitespace-nowrap;
  min-width: 50px;
}

/* Vertical bar chart */
.vbar-chart {
  @apply flex items-end justify-around gap-3;
  min-height: 80px;
}
.vbar-item {
  @apply flex flex-col items-center gap-1 flex-1;
}
.vbar-value {
  @apply text-sm font-bold tabular-nums;
}
.vbar-track {
  @apply w-full flex items-end;
  height: 60px;
}
.vbar-fill {
  @apply w-full rounded-t-md transition-all duration-500;
  min-height: 3px;
}
.vbar-label {
  @apply text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap;
}
</style>
