<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title"> 전압 / 전류 </h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'main' ? t('dashboard.meter.subtitle_main') : t('dashboard.meter.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- 주요 지표 요약 -->
    <div class="summary-section">
      <div class="summary-grid">
        <!-- 평균 전압 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-value">{{ data2.U4 || 0 }} <span class="summary-unit">V</span></div>
            <div class="summary-label">{{ t('dashboard.meter.avgvoltage') }}</div>
          </div>
        </div>

        <!-- 총 전류 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-value">{{ data2.Itot || 0 }} <span class="summary-unit">A</span></div>
            <div class="summary-label">{{ t('dashboard.meter.totcurrent') }}</div>
          </div>
        </div>

        <!-- 주파수 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-value">{{ data2.Ig || 0 }} <span class="summary-unit">A</span></div>
            <div class="summary-label"> 누설 전류 </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 상세 측정값 섹션 -->
    <div class="details-section">
      <div class="details-grid">
        <!-- 전압 상세 -->
        <div class="data-table">
          <div class="detail-header">
            <h4 class="detail-title">{{ t('dashboard.meter.voltage') }}</h4>
          </div>
          <div class="detail-values">
            <div class="value-row">
              <span class="phase-label">L1</span>
              <div class="phase-value" :class="getVoltageClass(data2.U1)">
                {{ data2.U1 || 0 }} 
                <span class="value-unit">V</span>
              </div>
            </div>
            <div class="value-row">
              <span class="phase-label">L2</span>
              <div class="phase-value" :class="getVoltageClass(data2.U2)">
                {{ data2.U2 || 0 }} 
                <span class="value-unit">V</span>
              </div>
            </div>
            <div class="value-row">
              <span class="phase-label">L3</span>
              <div class="phase-value" :class="getVoltageClass(data2.U3)">
                {{ data2.U3 || 0 }} 
                <span class="value-unit">V</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 전류 상세 -->
        <div class="data-table">
          <div class="detail-header">
            <span class="detail-title">{{ t('dashboard.meter.current') }}</span>
          </div>
          <div class="detail-values">
            <div class="value-row">
              <span class="phase-label">L1</span>
              <div class="phase-value" :class="getCurrentClass(data2.I1)">
                {{ data2.I1 || 0 }} 
                <span class="value-unit">A</span>
              </div>
            </div>
            <div class="value-row">
              <span class="phase-label">L2</span>
              <div class="phase-value" :class="getCurrentClass(data2.I2)">
                {{ data2.I2 || 0 }} 
                <span class="value-unit">A</span>
              </div>
            </div>
            <div class="value-row">
              <span class="phase-label">L3</span>
              <div class="phase-value" :class="getCurrentClass(data2.I3)">
                {{ data2.I3 || 0 }} 
                <span class="value-unit">A</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { watch, ref } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'DashboardCard_Meter_Claude',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const { t } = useI18n()
    const channel = ref(props.channel)
    const data2 = ref({})

    // 전압 상태 클래스 - 다크모드 개선
    const getVoltageClass = (value) => {
      if (!value) return 'text-gray-500 dark:text-gray-400'
      if (value < 200 || value > 240) return 'text-red-600 dark:text-red-400 font-bold'
      if (value < 210 || value > 230) return 'text-orange-500 dark:text-orange-400 font-semibold'
      return 'text-green-600 dark:text-green-400 font-medium'
    }

    // 전류 상태 클래스 - 다크모드 개선
    const getCurrentClass = (value) => {
      if (!value) return 'text-gray-500 dark:text-gray-400'
      if (value > 100) return 'text-red-600 dark:text-red-400 font-bold'
      if (value > 80) return 'text-orange-500 dark:text-orange-400 font-semibold'
      if (value > 60) return 'text-yellow-600 dark:text-yellow-400 font-semibold'
      return 'text-green-600 dark:text-green-400 font-medium'
    }

    // 역률 상태 클래스 - 다크모드 개선
    const getPowerFactorClass = (value) => {
      if (!value) return 'text-gray-500 dark:text-gray-400'
      if (value < 80) return 'text-red-600 dark:text-red-400 font-bold'
      if (value < 90) return 'text-orange-500 dark:text-orange-400 font-semibold'
      if (value < 95) return 'text-yellow-600 dark:text-yellow-400 font-semibold'
      return 'text-green-600 dark:text-green-400 font-medium'
    }

    // props.data 감시
    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          data2.value = newData
        }
      },
      { immediate: true }
    )

    return {
      channel,
      data2,
      t,
      getVoltageClass,
      getCurrentClass,
      getPowerFactorClass,
    }
  },
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply col-span-full sm:col-span-4 xl:col-span-4;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-50/50 to-purple-50/50 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-white;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 bg-clip-text text-transparent;
}

.channel-info {
  @apply flex items-center;
}

.channel-text {
  @apply text-xs font-semibold text-gray-400 dark:text-gray-300 uppercase;
}

.summary-section {
  @apply px-5 py-3 mt-1;
  @apply bg-white dark:bg-gray-800;
}

.summary-grid {
  @apply grid grid-cols-3 gap-4;
}

.summary-item {
  @apply flex items-center gap-3;
  @apply p-3 rounded-lg bg-white dark:bg-gray-800/50;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-sm;
}

.summary-icon {
  @apply w-10 h-10 rounded-lg flex items-center justify-center;
  @apply shadow-sm;
}

.voltage-icon {
  @apply bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400;
}

.current-icon {
  @apply bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400;
}

.frequency-icon {
  @apply bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400;
}

.summary-content {
  @apply flex-1;
}

.summary-value {
  @apply text-2xl font-bold text-gray-800 dark:text-white;
  @apply flex items-baseline gap-1;
}

.summary-unit {
  @apply text-lg font-semibold text-gray-500 dark:text-gray-300;
}

.summary-label {
  @apply text-sm font-medium text-gray-600 dark:text-gray-300;
  @apply mt-1;
}

.data-table {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
  @apply border border-gray-200 dark:border-gray-600;
}

.details-section {
  @apply px-5 py-4;
}

.details-title {
  @apply text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3;
  @apply flex items-center gap-2;
}

.details-title::before {
  content: '';
  @apply w-2 h-2 bg-blue-500 rounded-full;
}

.details-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-3;
}

.detail-card {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-sm;
}

.detail-header {
  @apply bg-gray-100 dark:bg-gray-600 px-3 py-2 border-b border-gray-200 dark:border-gray-500;
}

.detail-icon {
  @apply w-6 h-6 rounded flex items-center justify-center;
}

.detail-title {
  @apply text-sm font-bold text-gray-700 dark:text-white;
}

.detail-values {
  @apply py-1 px-3 space-y-2;
}

.value-row {
  @apply flex justify-between items-center;
  @apply py-1;
}

.phase-label {
  @apply text-xs font-medium text-gray-600 dark:text-gray-300;
}

.phase-value {
  @apply text-lg font-bold text-gray-800 dark:text-white;
}

/* 단위 텍스트 스타일 개선 */
.value-unit {
  @apply text-sm font-semibold text-gray-500 dark:text-gray-300;
}

.power-summary {
  @apply px-5 py-3 space-y-3;
  @apply bg-gray-50 dark:bg-gray-700/20;
}

.power-factor-card {
  @apply bg-white dark:bg-gray-800 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-600;
}

.pf-header {
  @apply flex justify-between items-center mb-2;
}

.pf-label {
  @apply text-sm font-medium text-gray-600 dark:text-gray-300;
}

.pf-value {
  @apply text-lg font-bold text-gray-800 dark:text-white;
}

.pf-progress {
  @apply w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden;
}

.pf-fill {
  @apply h-full bg-gradient-to-r from-green-400 to-green-500 rounded-full;
  @apply transition-all duration-500 ease-out;
}

.energy-cards {
  @apply grid grid-cols-2 gap-3;
}

.energy-card {
  @apply bg-white dark:bg-gray-800 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-600;
  @apply flex items-center gap-3;
  @apply transition-all duration-200 hover:shadow-sm;
}

.energy-icon {
  @apply w-8 h-8 rounded-lg flex items-center justify-center;
  @apply bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400;
}

.import-card .energy-icon {
  @apply bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400;
}

.export-card .energy-icon {
  @apply bg-orange-100 text-orange-600 dark:bg-orange-900/30 dark:text-orange-400;
}

.energy-content {
  @apply flex flex-col;
}

.energy-value {
  @apply text-lg font-bold text-gray-800 dark:text-white;
}

.energy-label {
  @apply text-xs font-medium text-gray-600 dark:text-gray-300;
}

/* 반응형 개선 */
@media (max-width: 768px) {
  .summary-grid {
    @apply grid-cols-1 gap-3;
  }
  
  .summary-item {
    @apply p-4;
  }
  
  .details-grid {
    @apply grid-cols-1 gap-4;
  }
  
  .energy-cards {
    @apply grid-cols-1;
  }
}

@media (max-width: 1024px) {
  .summary-grid {
    @apply grid-cols-1 gap-3;
  }
  
  .details-grid {
    @apply grid-cols-1 gap-3;
  }
}

/* 호버 효과 */
.summary-item:hover,
.detail-card:hover,
.energy-card:hover {
  @apply transform scale-105;
}

.pf-fill {
  position: relative;
}

.pf-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255,255,255,0.3) 50%, transparent 100%);
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
</style>