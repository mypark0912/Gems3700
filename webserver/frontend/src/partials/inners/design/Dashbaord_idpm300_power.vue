<template>
  <div class="card-wrap">
    <div class="card-header">
      <h3 class="card-title meter-accent-violet">{{ t('dashboard.powerCard') }}</h3>
    </div>
    <div class="card-body">
      <div class="pf-row">
        <!-- 역률 링 게이지 -->
        <div class="gauge-section">
          <div class="ring-gauge">
            <svg :width="ringSize" :height="ringSize" class="-rotate-90">
              <circle :cx="ringSize/2" :cy="ringSize/2" :r="ringR"
                fill="none" stroke="#f1f5f9" :stroke-width="ringStroke"
                class="dark:stroke-gray-700" />
              <circle :cx="ringSize/2" :cy="ringSize/2" :r="ringR"
                fill="none" :stroke="pfColor" :stroke-width="ringStroke"
                :stroke-dasharray="ringCirc" :stroke-dashoffset="ringOffset"
                stroke-linecap="round" class="transition-all duration-700" />
            </svg>
            <div class="ring-center">
              <span class="ring-value" :style="{ color: pfColor }">{{ pfDisplay }}</span>
              <span class="ring-unit">%</span>
            </div>
          </div>
          <span class="gauge-title">{{ t('dashboard.pq.powerfactor') }}</span>
          <span class="pf-status" :class="pfStatusClass">{{ pfStatusText }}</span>
        </div>

        <!-- 전력 수치 -->
        <div class="power-section">
          <div class="power-item">
            <span class="power-label">{{ t('dashboard.activePower') }}</span>
            <div class="power-value">{{ data2.P4?.toFixed(2) || 0 }} <span class="power-unit">kW</span></div>
          </div>
          <div class="power-item">
            <span class="power-label">{{ t('dashboard.reactivePower') }}</span>
            <div class="power-value">{{ data2.Q4?.toFixed(2) || 0 }} <span class="power-unit">kVar</span></div>
          </div>
          <div class="power-item">
            <span class="power-label">{{ t('dashboard.apparentPower') }}</span>
            <div class="power-value">{{ data2.S4?.toFixed(2) || 0 }} <span class="power-unit">kVA</span></div>
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
  name: 'DesignCard_Power',
  setup() {
    const { t } = useI18n()
    const store = useRealtimeStore()
    const data2 = computed(() => store.getChannelData('Main') || {})

    const pfValue = computed(() => parseFloat(data2.value.PF4) || 0)
    const pfDisplay = computed(() => pfValue.value.toFixed(2))

    const ringSize = 96
    const ringStroke = 8
    const ringR = (ringSize - ringStroke) / 2
    const ringCirc = 2 * Math.PI * ringR
    const ringOffset = computed(() => {
      const pct = Math.min(pfValue.value / 100, 1)
      return ringCirc * (1 - pct)
    })

    const pfColor = computed(() => {
      const v = pfValue.value
      if (v >= 95) return '#22c55e'
      if (v >= 85) return '#f59e0b'
      return '#ef4444'
    })

    const pfStatusClass = computed(() => {
      const v = pfValue.value
      if (v >= 95) return 'text-green-500'
      if (v >= 85) return 'text-amber-500'
      return 'text-red-500'
    })

    const pfStatusText = computed(() => {
      const v = pfValue.value
      if (v >= 95) return t('dashboard.pfGood')
      if (v >= 85) return t('dashboard.pfWarn')
      return t('dashboard.pfBad')
    })

    return {
      t, data2, pfDisplay, pfColor,
      ringSize, ringStroke, ringR, ringCirc, ringOffset,
      pfStatusClass, pfStatusText,
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
.meter-accent-violet::before {
  @apply bg-violet-500;
}
.card-body {
  @apply px-4 py-3;
}

.pf-row {
  @apply grid items-center gap-3;
  grid-template-columns: 4fr 6fr;
}

/* Ring Gauge */
.gauge-section {
  @apply flex flex-col items-center;
}
.ring-gauge {
  @apply relative;
}
.ring-center {
  @apply absolute inset-0 flex items-center justify-center flex-col;
}
.ring-value {
  @apply text-xl font-extrabold tabular-nums leading-none;
}
.ring-unit {
  @apply text-xs text-gray-600 dark:text-gray-400;
}
.gauge-title {
  @apply text-sm text-gray-600 dark:text-gray-400 mt-1;
}
.pf-status {
  @apply text-sm font-semibold;
}

/* Power */
.power-section {
  @apply space-y-1.5;
}
.power-item {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg px-2.5 py-1.5;
}
.power-label {
  @apply text-sm text-gray-600 dark:text-gray-400;
}
.power-value {
  @apply text-lg font-bold text-gray-800 dark:text-white tabular-nums;
}
.power-unit {
  @apply text-sm font-medium text-gray-600 dark:text-gray-400;
}
</style>
