<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">{{ t('dashboard.pq.title') }}</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'main' ? t('dashboard.meter.subtitle_main') : t('dashboard.meter.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- 불평형률 & 고조파 통합 섹션 -->
    <div class="data-section">
      <!-- 불평형률 -->
      <div class="data-subsection">
        <h3 class="subsection-title">{{ t('dashboard.pq.unbalance') }}</h3>
        <div class="unbalance-vertical">
          <!-- 전압 불평형 -->
          <div class="compact-item">
            <div class="compact-header">
              <span class="compact-label">{{ t('dashboard.meter.voltage') }}</span>
              <span class="compact-value" :class="getUnbalanceClass(data2.Ubal1)">
                {{ data2.Ubal1?.toFixed(1) || 0 }}%
              </span>
            </div>
            <div class="mini-progress">
              <div 
                class="mini-fill bg-violet-400"
                :style="{ width: Math.min(data2.Ubal1 || 0, 100) + '%' }"
              ></div>
            </div>
          </div>

          <!-- 전류 불평형 -->
          <div class="compact-item">
            <div class="compact-header">
              <span class="compact-label">{{ t('dashboard.meter.current') }}</span>
              <span class="compact-value" :class="getUnbalanceClass(data2.Ibal1)">
                {{ data2.Ibal1?.toFixed(1) || 0 }}%
              </span>
            </div>
            <div class="mini-progress">
              <div 
                class="mini-fill bg-sky-400"
                :style="{ width: Math.min(data2.Ibal1 || 0, 100) + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 고조파 -->
      <div class="data-subsection">
        <h3 class="subsection-title">{{ t('pq.tabs.harmonics') }}</h3>
        
        <!-- CSS 바 차트 -->
        <div class="css-bar-chart">
          <div class="chart-bars">
            <div 
              v-for="(item, index) in chartItems" 
              :key="index"
              class="bar-item"
            >
              <!-- 바 컨테이너 -->
              <div class="bar-container">
                <div 
                  class="bar-fill"
                  :class="item.colorClass"
                  :style="{ height: item.height + '%' }"
                >
                  <!-- 값 라벨 -->
                  <div class="bar-label">{{ item.value }}%</div>
                </div>
              </div>
              
              <!-- 하단 라벨 -->
              <div class="bar-name">{{ item.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { watch, ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'DashboardCard_PowerQuality',
  props: {
    channel: String,
    data: Object,
  },
  emits: ['data-change'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const channel = ref(props.channel)
    const data2 = ref({})

    // 고조파 차트 데이터 설정
    const dataKeys = ['thdu total', 'thdi total', 'tddi total']
    const labels = ['THD-U', 'THD-I', 'TDD-I']
    const colors = ['purple', 'blue', 'green']

    // 차트 아이템 계산
    const chartItems = computed(() => {
      const values = dataKeys.map(key => parseFloat(data2.value[key] || 0))
      const maxValue = Math.max(...values, 10) // 최소 10으로 설정
      
      const items = labels.map((label, index) => {
        const value = values[index].toFixed(1)
        const height = (values[index] / maxValue) * 100
        
        return {
          label,
          value,
          height: Math.max(height, 5), // 최소 5% 높이
          colorClass: `bar-${colors[index]}`
        }
      })
      
      // 데이터 변경 이벤트 발생
      emit('data-change', { data: items })
      
      return items
    })

    // 불평형률 상태 클래스 - 다크모드 개선
    const getUnbalanceClass = (value) => {
      if (!value) return 'text-gray-500 dark:text-gray-400'
      if (value >= 3) return 'text-red-600 dark:text-red-400 font-bold'
      if (value >= 2) return 'text-orange-500 dark:text-orange-400 font-semibold'
      if (value >= 1) return 'text-yellow-600 dark:text-yellow-400 font-semibold'
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
      getUnbalanceClass,
      chartItems,
    }
  },
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-3;
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

.data-section {
  @apply px-5 py-3 space-y-6;
}

.data-subsection {
  @apply space-y-2;
}

.subsection-title {
  @apply text-sm font-semibold text-gray-700 dark:text-white;
  @apply flex items-center gap-2;
}

.subsection-title::before {
  content: '';
  @apply w-2 h-2 bg-blue-500 rounded-full;
}

.unbalance-vertical {
  @apply space-y-5;
}

.compact-item {
  @apply space-y-1;
}

.compact-header {
  @apply flex justify-between items-center;
}

.compact-label {
  @apply text-sm font-medium text-gray-600 dark:text-gray-300;
}

.compact-value {
  @apply text-sm font-bold;
}

.mini-progress {
  @apply w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden;
}

.mini-fill {
  @apply h-full rounded-full transition-all duration-500 ease-out;
}

/* CSS 바 차트 스타일 */
.css-bar-chart {
  @apply w-full p-4 pt-8;
}

.chart-bars {
  @apply flex justify-around items-end gap-4 h-24;
  @apply relative;
}

.bar-item {
  @apply flex-1 flex flex-col items-center;
  @apply max-w-[60px];
}

.bar-container {
  @apply relative w-full h-16 flex items-end justify-center;
}

.bar-fill {
  @apply w-full rounded-t transition-all duration-700 ease-out;
  @apply relative flex items-start justify-center;
  @apply min-h-[8px];
}

.bar-label {
  @apply text-xs font-bold;
  @apply absolute bottom-full left-1/2 transform -translate-x-1/2;
  @apply text-center whitespace-nowrap;
  @apply text-gray-800 dark:text-white;
  @apply mb-1;
}

.bar-name {
  @apply text-xs font-semibold text-gray-700 dark:text-white;
  @apply mt-2 text-center;
}

/* 바 색상 - 다크모드에서도 잘 보이도록 조정 */
.bar-purple {
  @apply bg-gradient-to-t from-purple-500 to-purple-400;
  @apply dark:from-purple-400 dark:to-purple-300;
}

.bar-blue {
  @apply bg-gradient-to-t from-blue-500 to-blue-400;
  @apply dark:from-blue-400 dark:to-blue-300;
}

.bar-green {
  @apply bg-gradient-to-t from-green-500 to-green-400;
  @apply dark:from-green-400 dark:to-green-300;
}

/* 반응형 개선 */
@media (max-width: 640px) {
  .unbalance-vertical {
    @apply space-y-4;
  }
  
  .data-section {
    @apply space-y-3;
  }
  
  .chart-bars {
    @apply gap-2 h-20;
  }
  
  .bar-container {
    @apply h-12;
  }
  
  .bar-label {
    @apply text-xs;
  }
  
  .bar-name {
    @apply text-xs;
  }
}

/* 추가 다크모드 개선 */
.mini-fill.bg-violet-400 {
  @apply dark:bg-violet-300;
}

.mini-fill.bg-sky-400 {
  @apply dark:bg-sky-300;
}
</style>