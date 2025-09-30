<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">역률/전력</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'main' ? t('dashboard.meter.subtitle_main') : t('dashboard.meter.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- 데이터 섹션 - 2열 레이아웃 -->
    <div class="data-section-grid">
      <!-- 왼쪽: 역률 게이지 -->
      <div class="gauge-column">
        <h3 class="subsection-title">역률</h3>
        
        <!-- Chart.js 게이지 -->
        <div class="gauge-container">
          <canvas ref="gaugeChart" width="160" height="160"></canvas>
          <div class="gauge-center">
            <div class="gauge-value">{{ (data2.powerFactor || 0).toFixed(2) }}</div>
            <div class="gauge-label">PF</div>
          </div>
        </div>
        
        <!-- 역률 상태 -->
        <div class="gauge-status">
          <span class="status-text">상태:</span>
          <span class="status-badge" :class="getPowerFactorStatusClass(data2.powerFactor)">
            {{ getPowerFactorStatus(data2.powerFactor) }}
          </span>
        </div>
        <div class="status-description">
          {{ getPowerFactorDescription(data2.powerFactor) }}
        </div>
      </div>

      <!-- 오른쪽: 전력 측정값 -->
      <div class="power-column">
        <h3 class="subsection-title">전력</h3>
        
        <div class="power-list">
          <!-- 유효 전력 -->
          <div class="power-item">
            <div class="power-info">
              <span class="power-label">유효 전력</span>
              <span class="power-value text-green-600 dark:text-green-400">
                {{ formatPower(data2.activePower) }}
                <span class="power-unit">kW</span>
              </span>
            </div>
            <div class="power-bar-container">
              <div class="power-bar bg-green-500/20 dark:bg-green-400/20">
                <div 
                  class="power-bar-fill bg-green-500 dark:bg-green-400"
                  :style="{ width: getPowerPercentage(data2.activePower, 2000) + '%' }"
                ></div>
              </div>
              <div class="power-meta">
                <span class="power-max">Max: 2,000 kW</span>
                <span class="power-percent">{{ getPowerPercentage(data2.activePower, 2000) }}%</span>
              </div>
            </div>
          </div>

          <!-- 무효 전력 -->
          <div class="power-item">
            <div class="power-info">
              <span class="power-label">무효 전력</span>
              <span class="power-value text-yellow-600 dark:text-yellow-400">
                {{ formatPower(data2.reactivePower) }}
                <span class="power-unit">kVAR</span>
              </span>
            </div>
            <div class="power-bar-container">
              <div class="power-bar bg-yellow-500/20 dark:bg-yellow-400/20">
                <div 
                  class="power-bar-fill bg-yellow-500 dark:bg-yellow-400"
                  :style="{ width: getPowerPercentage(data2.reactivePower, 500) + '%' }"
                ></div>
              </div>
              <div class="power-meta">
                <span class="power-max">Max: 500 kVAR</span>
                <span class="power-percent">{{ getPowerPercentage(data2.reactivePower, 500) }}%</span>
              </div>
            </div>
          </div>

          <!-- 피상 전력 -->
          <div class="power-item">
            <div class="power-info">
              <span class="power-label">피상 전력</span>
              <span class="power-value text-blue-600 dark:text-blue-400">
                {{ formatPower(data2.apparentPower) }}
                <span class="power-unit">kVA</span>
              </span>
            </div>
            <div class="power-bar-container">
              <div class="power-bar bg-blue-500/20 dark:bg-blue-400/20">
                <div 
                  class="power-bar-fill bg-blue-500 dark:bg-blue-400"
                  :style="{ width: getPowerPercentage(data2.apparentPower, 2000) + '%' }"
                ></div>
              </div>
              <div class="power-meta">
                <span class="power-max">Max: 2,000 kVA</span>
                <span class="power-percent">{{ getPowerPercentage(data2.apparentPower, 2000) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { watch, ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

export default {
  name: 'DashboardCard_PowerMetrics',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const { t } = useI18n()
    const channel = ref(props.channel)
    const gaugeChart = ref(null)
    let chartInstance = null
    
    // 더미 데이터
    const data2 = ref({
      powerFactor: 0.92,
      activePower: 1250.5,
      reactivePower: 320.8,
      apparentPower: 1290.3,
    })

    // Chart.js 게이지 생성
    const createGaugeChart = () => {
      if (!gaugeChart.value) return
      
      const ctx = gaugeChart.value.getContext('2d')
      const pf = data2.value.powerFactor || 0
      
      // 색상 결정
      let gaugeColor = '#3b82f6' // blue (default for good)
      if (pf < 0.85) gaugeColor = '#ef4444' // red
      else if (pf < 0.90) gaugeColor = '#f59e0b' // yellow
      else if (pf >= 0.95) gaugeColor = '#10b981' // green
      
      chartInstance = new Chart(ctx, {
        type: 'doughnut',
        data: {
          datasets: [{
            data: [pf, 1 - pf],
            backgroundColor: [gaugeColor, '#e5e7eb'],
            borderWidth: 0,
          }]
        },
        options: {
          responsive: false,
          maintainAspectRatio: false,
          rotation: -90,
          circumference: 360,
          cutout: '75%',
          plugins: {
            legend: { display: false },
            tooltip: { enabled: false }
          }
        }
      })
    }

    // 차트 업데이트
    const updateChart = () => {
      if (!chartInstance) return
      
      const pf = data2.value.powerFactor || 0
      let gaugeColor = '#3b82f6' // blue (default for good)
      if (pf < 0.85) gaugeColor = '#ef4444' // red
      else if (pf < 0.90) gaugeColor = '#f59e0b' // yellow
      else if (pf >= 0.95) gaugeColor = '#10b981' // green
      
      chartInstance.data.datasets[0].data = [pf, 1 - pf]
      chartInstance.data.datasets[0].backgroundColor = [gaugeColor, '#e5e7eb']
      chartInstance.update('none')
    }

    // 역률 상태 관련 함수들
    const getPowerFactorStatus = (pf) => {
      if (!pf) return '측정중'
      if (pf >= 0.95) return '우수'
      if (pf >= 0.90) return '양호'
      if (pf >= 0.85) return '보통'
      return '개선필요'
    }

    const getPowerFactorStatusClass = (pf) => {
      if (!pf) return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
      if (pf >= 0.95) return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
      if (pf >= 0.90) return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
      if (pf >= 0.85) return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
    }

    const getPowerFactorDescription = (pf) => {
      if (!pf) return '데이터 수집중'
      if (pf >= 0.95) return '최적 효율 운영'
      if (pf >= 0.90) return '정상 범위 내 운영'
      if (pf >= 0.85) return '개선 검토 필요'
      return '개선 시급'
    }

    // 전력 관련 함수들
    const formatPower = (value) => {
      if (!value) return '0'
      return value.toLocaleString('ko-KR', { maximumFractionDigits: 1 })
    }

    const getPowerPercentage = (value, max) => {
      if (!value || !max) return 0
      return Math.min(Math.round((value / max) * 100), 100)
    }

    // Lifecycle
    onMounted(() => {
      setTimeout(() => {
        createGaugeChart()
      }, 100)
    })

    onUnmounted(() => {
      if (chartInstance) {
        chartInstance.destroy()
      }
    })

    // props.data 감시
    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          data2.value = { ...data2.value, ...newData }
          updateChart()
        }
      },
      { immediate: true }
    )

    // data2 변경 감시
    watch(
      () => data2.value.powerFactor,
      () => {
        updateChart()
      }
    )

    return {
      channel,
      data2,
      t,
      gaugeChart,
      getPowerFactorStatus,
      getPowerFactorStatusClass,
      getPowerFactorDescription,
      formatPower,
      getPowerPercentage,
    }
  },
}
</script>

<style scoped>
/* 기존 코드와 일치하는 카드 스타일 */
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-4;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 - 기존과 동일 */
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

/* 2열 그리드 레이아웃 */
.data-section-grid {
  @apply grid grid-cols-5 gap-4 p-4;
}

.gauge-column {
  @apply col-span-2 flex flex-col;
}

.power-column {
  @apply col-span-3;
}

.subsection-title {
  @apply text-sm font-semibold text-gray-700 dark:text-white mb-3;
  @apply flex items-center gap-2;
}

.subsection-title::before {
  content: '';
  @apply w-2 h-2 bg-blue-500 rounded-full;
}

.subsection-title-left {
  @apply self-start;
}

/* 게이지 스타일 */
.gauge-container {
  @apply relative mx-auto;
}

.gauge-center {
  @apply absolute inset-0 flex flex-col items-center justify-center;
}

.gauge-value {
  @apply text-xl font-bold text-gray-800 dark:text-white;
}

.gauge-label {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.gauge-status {
  @apply flex items-center gap-2 mt-3;
}

.status-text {
  @apply text-xs text-gray-600 dark:text-gray-400;
}

.status-badge {
  @apply px-2 py-0.5 text-xs font-semibold rounded-full;
}

.status-description {
  @apply text-xs text-gray-500 dark:text-gray-400 text-center mt-2;
}

/* 전력 리스트 */
.power-list {
  @apply space-y-4;
}

.power-item {
  @apply space-y-2;
}

.power-info {
  @apply flex justify-between items-baseline;
}

.power-label {
  @apply text-xs font-medium text-gray-600 dark:text-gray-300;
}

.power-value {
  @apply text-base font-bold;
}

.power-unit {
  @apply text-xs font-normal text-gray-500 dark:text-gray-400 ml-1;
}

.power-bar-container {
  @apply space-y-1;
}

.power-bar {
  @apply w-full h-6 rounded-md overflow-hidden relative;
}

.power-bar-fill {
  @apply h-full rounded-md transition-all duration-700 ease-out;
}

.power-meta {
  @apply flex justify-between text-xs text-gray-500 dark:text-gray-400;
}

.power-max {
  @apply text-xs;
}

.power-percent {
  @apply font-semibold;
}

/* 반응형 */
@media (max-width: 640px) {
  .data-section-grid {
    @apply grid-cols-1;
  }
  
  .gauge-column {
    @apply col-span-1;
  }
  
  .power-column {
    @apply col-span-1;
  }
}
</style>