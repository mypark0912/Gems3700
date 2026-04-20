<template>
  <div class="card-wrap">
    <div class="card-header">
      <h3 class="card-title meter-accent-blue">{{ t('dashboard.meter.singletitle') }}</h3>
    </div>

    <!-- 주요 지표 요약 -->
    <div class="card-body">
      <div class="summary-row">
        <div class="summary-item">
          <span class="summary-label">{{ t('dashboard.meter.avgvoltage') }}</span>
          <div class="summary-value">
            <span class="status-dot" :class="getVoltageDotClass(realtimeData?.U4)"></span>
            {{ realtimeData?.U4?.toFixed(1) }} <span class="summary-unit">V</span>
          </div>
        </div>
        <div class="summary-item">
          <span class="summary-label">{{ t('dashboard.meter.totcurrent') }}</span>
          <div class="summary-value">
            <span class="status-dot" :class="getCurrentDotClass(realtimeData?.I4)"></span>
            {{ realtimeData?.I4?.toFixed(2) }} <span class="summary-unit">A</span>
          </div>
        </div>
        <div class="summary-item">
          <span class="summary-label">{{ t('dashboard.meter.frequency') }}</span>
          <div class="summary-value">
            <span class="status-dot" :class="getFreqDotClass(realtimeData?.Freq)"></span>
            {{ realtimeData?.Freq?.toFixed(2) }} <span class="summary-unit">Hz</span>
          </div>
        </div>
      </div>

      <!-- 상세 측정값 -->
      <div class="detail-grid">
        <!-- 전압 -->
        <div class="detail-block">
          <span class="detail-block-title">{{ t('dashboard.meter.voltage') }}</span>
          <div class="phase-section">
            <div class="phase-row">
              <span class="phase-label">L1</span>
              <span class="phase-value" :class="getVoltageTextClass(realtimeData?.U1)">{{ realtimeData?.U1?.toFixed(1) }} <span class="phase-unit">V</span></span>
            </div>
            <div class="phase-row">
              <span class="phase-label">L2</span>
              <span class="phase-value" :class="getVoltageTextClass(realtimeData?.U2)">{{ realtimeData?.U2?.toFixed(1) }} <span class="phase-unit">V</span></span>
            </div>
            <div class="phase-row">
              <span class="phase-label">L3</span>
              <span class="phase-value" :class="getVoltageTextClass(realtimeData?.U3)">{{ realtimeData?.U3?.toFixed(1) }} <span class="phase-unit">V</span></span>
            </div>
          </div>
        </div>
        <!-- 전류 -->
        <div class="detail-block">
          <span class="detail-block-title">{{ t('dashboard.meter.current') }}</span>
          <div class="phase-section">
            <div class="phase-row">
              <span class="phase-label">L1</span>
              <span class="phase-value" :class="getCurrentTextClass(realtimeData?.I1)">{{ realtimeData?.I1?.toFixed(2) }} <span class="phase-unit">A</span></span>
            </div>
            <div class="phase-row">
              <span class="phase-label">L2</span>
              <span class="phase-value" :class="getCurrentTextClass(realtimeData?.I2)">{{ realtimeData?.I2?.toFixed(2) }} <span class="phase-unit">A</span></span>
            </div>
            <div class="phase-row">
              <span class="phase-label">L3</span>
              <span class="phase-value" :class="getCurrentTextClass(realtimeData?.I3)">{{ realtimeData?.I3?.toFixed(2) }} <span class="phase-unit">A</span></span>
            </div>
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
  name: 'DesignCard_Meter',
  setup() {
    const { t } = useI18n()
    const store = useRealtimeStore()
    const realtimeData = computed(() => store.getChannelData('Main') || {})

    const getVoltageDotClass = (v) => {
      if (!v) return 'dot-unknown'
      if (v < 200 || v > 240) return 'dot-danger'
      if (v < 210 || v > 230) return 'dot-warn'
      return 'dot-good'
    }

    const getCurrentDotClass = (v) => {
      if (!v) return 'dot-unknown'
      if (v > 100) return 'dot-danger'
      if (v > 80) return 'dot-warn'
      return 'dot-good'
    }

    const getFreqDotClass = (v) => {
      if (!v) return 'dot-unknown'
      if (v < 59.5 || v > 60.5) return 'dot-danger'
      if (v < 59.8 || v > 60.2) return 'dot-warn'
      return 'dot-good'
    }

    const getVoltageTextClass = (v) => {
      if (!v) return 'text-gray-400'
      if (v < 200 || v > 240) return 'text-red-500 dark:text-red-400'
      if (v < 210 || v > 230) return 'text-amber-500 dark:text-amber-400'
      return 'text-green-600 dark:text-green-400'
    }

    const getCurrentTextClass = (v) => {
      if (!v) return 'text-gray-400'
      if (v > 100) return 'text-red-500 dark:text-red-400'
      if (v > 80) return 'text-amber-500 dark:text-amber-400'
      return 'text-green-600 dark:text-green-400'
    }

    return {
      t, realtimeData,
      getVoltageDotClass, getCurrentDotClass, getFreqDotClass,
      getVoltageTextClass, getCurrentTextClass,
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
.meter-accent-blue::before {
  @apply bg-blue-500;
}
.card-body {
  @apply px-4 py-3;
}

/* Summary row */
.summary-row {
  @apply grid grid-cols-3 gap-3 mb-3;
}
.summary-item {
  @apply flex flex-col items-center text-center;
}
.summary-label {
  @apply text-sm text-gray-600 dark:text-gray-400 mb-1;
}
.summary-value {
  @apply text-2xl font-extrabold text-gray-800 dark:text-white tabular-nums flex items-center justify-center gap-1.5;
}
.summary-unit {
  @apply text-sm font-medium text-gray-600 dark:text-gray-400;
}

/* Status dot */
.status-dot {
  @apply w-2 h-2 rounded-full flex-shrink-0;
}
.dot-good { @apply bg-green-500; }
.dot-warn { @apply bg-amber-500; }
.dot-danger { @apply bg-red-500; }
.dot-unknown { @apply bg-gray-400; }

/* Detail grid */
.detail-grid {
  @apply grid grid-cols-2 gap-3;
}
.detail-block {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
}
.detail-block-title {
  @apply block text-sm font-bold text-gray-700 dark:text-white px-3 py-1.5;
  @apply bg-gray-100 dark:bg-gray-600 border-b border-gray-200 dark:border-gray-500;
}
.phase-section {
  @apply px-3 py-2 space-y-1;
}
.phase-row {
  @apply flex justify-between items-center;
}
.phase-label {
  @apply text-sm font-semibold text-gray-600 dark:text-gray-400;
}
.phase-value {
  @apply text-base font-bold tabular-nums;
}
.phase-unit {
  @apply text-sm font-medium text-gray-600 dark:text-gray-400;
}
</style>
