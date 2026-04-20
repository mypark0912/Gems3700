<template>
  <div class="meter-card p-5">
    <!-- 헤더 -->
    <div class="meter-card-header !px-0 !py-0 mb-4">
      <h3 class="meter-card-title meter-accent-blue">System / Voltage</h3>
      <span class="text-xs font-mono text-gray-500 dark:text-gray-400">{{ formattedTimestamp }}</span>
    </div>

    <!-- System 인라인 -->
    <div class="flex flex-wrap items-center gap-x-8 gap-y-2 pb-4 mb-5 border-b border-gray-200 dark:border-gray-700/60">
      <div class="flex items-baseline gap-2">
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">Temperature</span>
        <span class="text-2xl font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums">
          {{ systemData ? systemData.temperature.toFixed(1) : '-' }}
          <span class="text-sm font-semibold text-gray-400 dark:text-gray-500">°C</span>
        </span>
      </div>
      <div class="flex items-baseline gap-2">
        <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">Frequency</span>
        <span class="text-2xl font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums">
          {{ systemData ? systemData.frequency.toFixed(2) : '-' }}
          <span class="text-sm font-semibold text-gray-400 dark:text-gray-500">Hz</span>
        </span>
      </div>
    </div>

    <!-- Voltage (상전압 · 선간전압 수평 배치) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 상전압 -->
      <div class="voltage-section">
        <div class="flex items-center gap-2 mb-3">
          <div class="w-1 h-4 bg-indigo-500 rounded-full"></div>
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 uppercase tracking-wider">상전압 (Phase)</span>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div v-for="(lbl, i) in ['L1','L2','L3']" :key="'p'+i" class="voltage-cell">
            <span class="voltage-label">{{ lbl }}</span>
            <div class="flex items-baseline gap-1.5">
              <span class="voltage-value">{{ fmt(voltageData?.u?.[i]) }}</span>
              <span class="voltage-unit">V</span>
            </div>
            <div class="voltage-thd">
              <span class="thd-label">THD</span>
              <span class="thd-value">{{ fmt(voltageData?.thd_u?.[i]) }}<span class="thd-unit">%</span></span>
            </div>
          </div>
        </div>

        <div class="voltage-summary">
          <div class="summary-item">
            <span class="summary-label">정상분</span>
            <span class="summary-value">{{ fmt(voltageData?.u_pos) }}<span class="summary-unit">V</span></span>
          </div>
          <div class="summary-divider"></div>
          <div class="summary-item">
            <span class="summary-label">불평형률</span>
            <span class="summary-value">{{ fmt(voltageData?.u_unbal) }}<span class="summary-unit">%</span></span>
          </div>
        </div>
      </div>

      <!-- 선간전압 -->
      <div class="voltage-section">
        <div class="flex items-center gap-2 mb-3">
          <div class="w-1 h-4 bg-emerald-500 rounded-full"></div>
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 uppercase tracking-wider">선간전압 (Line)</span>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div v-for="(lbl, i) in ['L1-L2','L2-L3','L3-L1']" :key="'l'+i" class="voltage-cell">
            <span class="voltage-label">{{ lbl }}</span>
            <div class="flex items-baseline gap-1.5">
              <span class="voltage-value">{{ fmt(voltageData?.ull?.[i]) }}</span>
              <span class="voltage-unit">V</span>
            </div>
            <div class="voltage-thd">
              <span class="thd-label">THD</span>
              <span class="thd-value">{{ fmt(voltageData?.thd_ull?.[i]) }}<span class="thd-unit">%</span></span>
            </div>
          </div>
        </div>

        <div class="voltage-summary">
          <div class="summary-item">
            <span class="summary-label">평균</span>
            <span class="summary-value">{{ fmt(voltageData?.ull_avg) }}<span class="summary-unit">V</span></span>
          </div>
          <div class="summary-divider"></div>
          <div class="summary-item">
            <span class="summary-label">불평형률</span>
            <span class="summary-value">{{ fmt(voltageData?.ull_unbal) }}<span class="summary-unit">%</span></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  systemData: { type: Object, default: null },
  voltageData: { type: Object, default: null },
})

const fmt = (v) => v != null ? v.toFixed(2) : '-'

const formattedTimestamp = computed(() => {
  if (!props.systemData?.timestamp) return '-'
  return new Date(props.systemData.timestamp * 1000).toLocaleString()
})
</script>

<style scoped>
@import '../../css/meter-card.css';

/* Voltage 섹션 래퍼 */
.voltage-section {
  @apply rounded-xl p-4 bg-gray-50/60 dark:bg-gray-700/20 border border-gray-200/60 dark:border-gray-700/60;
}

/* Voltage 요약 (정상분/평균 · 불평형률) */
.voltage-summary {
  @apply mt-3 pt-3 flex items-center justify-around gap-3 border-t border-gray-200 dark:border-gray-700/70;
}
.summary-item {
  @apply flex flex-col items-center flex-1 gap-0.5;
}
.summary-label {
  @apply text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider;
}
.summary-value {
  @apply text-lg font-mono font-bold text-gray-800 dark:text-gray-100 tabular-nums;
}
.summary-unit {
  @apply text-xs font-semibold text-gray-400 dark:text-gray-500 ml-1;
}
.summary-divider {
  @apply w-px h-8 bg-gray-200 dark:bg-gray-700/70;
}

/* Voltage 셀 */
.voltage-cell {
  @apply relative flex flex-col items-center py-4 px-3 rounded-lg;
  @apply bg-gradient-to-b from-gray-50 to-white dark:from-gray-700/30 dark:to-gray-800/30;
  @apply border border-gray-200/60 dark:border-gray-700/60;
}
.voltage-label {
  @apply text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1;
}
.voltage-value {
  @apply text-xl md:text-2xl font-mono font-bold text-gray-900 dark:text-gray-100 tabular-nums leading-tight;
}
.voltage-unit {
  @apply text-sm font-semibold text-gray-400 dark:text-gray-500;
}
.voltage-thd {
  @apply flex items-center gap-2 mt-3 pt-3 w-full justify-center border-t border-gray-200/70 dark:border-gray-700/70;
}
.thd-label {
  @apply text-xs font-bold text-gray-400 dark:text-gray-500 uppercase tracking-wider;
}
.thd-value {
  @apply text-xl md:text-2xl font-mono font-bold text-indigo-600 dark:text-indigo-400 tabular-nums;
}
.thd-unit {
  @apply text-sm text-gray-400 dark:text-gray-500 ml-0.5 font-semibold;
}
</style>
