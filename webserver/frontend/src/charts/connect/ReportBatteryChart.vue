<template>
  <div 
    class="battery-chart-container"
    :class="isPdfMode ? 'bg-white border-gray-200' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700'"
  >
    <!-- 헤더 섹션 -->
    <div 
      class="chart-header"
      :class="isPdfMode ? 'bg-blue-50 border-blue-100' : 'bg-blue-50 dark:bg-gray-900 border-blue-100 dark:border-gray-700'"
    >
      <div class="flex justify-between items-center">
        <h3 
          class="chart-title"
          :class="isPdfMode ? 'text-gray-800' : 'text-gray-800 dark:text-gray-100'"
        >
          {{ chartTitle }}
        </h3>
        <!-- 타임스탬프 표시 -->
        <div v-if="timestamp" class="flex items-center gap-2 px-3 py-1.5 bg-indigo-50 dark:bg-indigo-900/30 rounded-md">
          <span class="text-sm text-indigo-600 dark:text-indigo-400 font-medium">📌 {{ t('report.lastSaved') || 'Last saved datetime' }}:</span>
          <span class="text-sm text-indigo-700 dark:text-indigo-300 font-semibold">{{ formatTimestamp(timestamp) }}</span>
        </div>
      </div>
    </div>

    <!-- 배터리 그리드 -->
    <div class="battery-wrapper">
      <div class="battery-grid" :style="gridStyle">
        <div
          v-for="(item, index) in batteryData"
          :key="index"
          class="battery-item"
        >
          <!-- 항목 타이틀 (위에) -->
          <div 
            class="item-title"
            :class="isPdfMode ? 'text-gray-700' : 'text-gray-700 dark:text-gray-300'"
          >
            {{ item.title }}
          </div>
          
          <!-- 배터리 바디 (세그먼트 방식) -->
          <div class="battery-container">
            <div class="battery-body">
              <div class="battery-segments">
                <div 
                  v-for="segment in 4" 
                  :key="segment"
                  class="battery-segment"
                  :class="[
                    `segment-${segment}`,
                    { 
                      'active': segment <= item.status && item.status > 0,
                      'inactive': segment > item.status || item.status === 0
                    }
                  ]"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 범례 -->
    <div class="chart-legend">
      <div class="legend-title">{{ t('diagnosis.tabContext.legend') || 'Legend' }}</div>
      <div class="legend-items">
        <div v-for="(item, idx) in legendItems" :key="idx" class="legend-item">
          <div class="legend-color" :class="item.colorClass"></div>
          <span 
            class="legend-label"
            :class="isPdfMode ? 'text-gray-600' : 'text-gray-600 dark:text-gray-400'"
          >
            {{ item.label }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, inject } from 'vue'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ReportBatteryChart',
  props: {
    data: {
      type: Object,
      required: true
    },
    mode: {
      type: String,
      default: 'DiagnosisDetail'
    },
    height: {
      type: Number,
      default: 350
    },
    timestamp: {
      type: String,
      default: null
    }
  },
  setup(props) {
    const isPdfMode = inject('isPdfMode', false)
    const { t, locale } = useI18n()
    const darkMode = useDark()

    // 타임스탬프 포맷
    const formatTimestamp = (timestamp) => {
      if (!timestamp) return '';
      const date = new Date(timestamp);
      return date.toLocaleString('ko-KR', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    };

    // 모드에 따른 타이틀
    const chartTitle = computed(() => {
      const titleMap = {
        'PowerQuality': t('diagnosis.tabTitle.detailTitle_pq'),
        'DiagnosisDetail': t('diagnosis.tabTitle.detailTitle'),
      }
      return titleMap[props.mode] || t('diagnosis.tabTitle.detailTitle')
    })

    // 배터리 데이터 변환
    const batteryData = computed(() => {
      if (!props.data?.Names || !props.data?.Values) return []

      return props.data.Names.map((name, index) => {
        const status = props.data.Values[index]
        const title = props.data.Titles?.[index]?.[locale.value] || name
        return {
          title: title,
          status: status === 0.2 ? 0 : Math.floor(status),
          originalIndex: index
        }
      })
    })

    // 그리드 스타일 계산
    const gridStyle = computed(() => {
      const itemCount = batteryData.value.length
      let columns = 6
      
      if (itemCount <= 4) columns = 4
      else if (itemCount <= 6) columns = 6
      else if (itemCount <= 9) columns = 6
      else if (itemCount <= 12) columns = 6
      else columns = 8
      
      return {
        'grid-template-columns': `repeat(${columns}, 1fr)`,
        'gap': '16px'
      }
    })

    // 범례 아이템
    const legendItems = computed(() => {
      if (props.mode === 'PowerQuality') {
        return [
          { label: t('diagnosis.tabContext.pqfe0'), colorClass: 'status-stopped' },
          { label: t('diagnosis.tabContext.pqfe1'), colorClass: 'status-normal' },
          { label: t('diagnosis.tabContext.pqfe2'), colorClass: 'status-warning' },
          { label: t('diagnosis.tabContext.pqfe3'), colorClass: 'status-alert' },
          { label: t('diagnosis.tabContext.pqfe4'), colorClass: 'status-danger' },
        ]
      } else {
        return [
          { label: t('diagnosis.tabContext.st0'), colorClass: 'status-stopped' },
          { label: t('diagnosis.tabContext.st1'), colorClass: 'status-normal' },
          { label: t('diagnosis.tabContext.st2'), colorClass: 'status-warning' },
          { label: t('diagnosis.tabContext.st3'), colorClass: 'status-alert' },
          { label: t('diagnosis.tabContext.st4'), colorClass: 'status-danger' },
        ]
      }
    })

    // timestamp computed로 변경 (reactive)
    const timestampValue = computed(() => props.timestamp)

    return {
      isPdfMode,
      t,
      chartTitle,
      batteryData,
      gridStyle,
      legendItems,
      formatTimestamp,
      timestamp: timestampValue,
    }
  }
}
</script>

<style scoped>
.battery-chart-container {
  @apply rounded-xl border overflow-hidden shadow-lg;
}

.chart-header {
  @apply px-5 py-4 border-b;
}

.chart-title {
  @apply text-lg font-bold;
}

.battery-wrapper {
  @apply relative p-5;
  min-height: 200px;
}

.battery-grid {
  @apply grid w-full;
}

/* 배터리 아이템 */
.battery-item {
  @apply flex flex-col items-center gap-2 p-3 rounded-lg;
  @apply bg-gray-50/50 dark:bg-gray-700/30;
  @apply border border-gray-200/30 dark:border-gray-700/30;
  animation: fadeInUp 0.6s ease-out forwards;
  opacity: 0;
  transform: translateY(20px);
}

.battery-item:nth-child(1) { animation-delay: 0.1s; }
.battery-item:nth-child(2) { animation-delay: 0.2s; }
.battery-item:nth-child(3) { animation-delay: 0.3s; }
.battery-item:nth-child(4) { animation-delay: 0.4s; }
.battery-item:nth-child(5) { animation-delay: 0.5s; }
.battery-item:nth-child(6) { animation-delay: 0.6s; }
.battery-item:nth-child(n+7) { animation-delay: 0.7s; }

/* 아이템 타이틀 */
.item-title {
  @apply text-xs font-medium text-center leading-tight;
  @apply max-w-full truncate;
}

/* 배터리 컨테이너 */
.battery-container {
  @apply relative;
}

.battery-body {
  @apply w-20 h-10 rounded-lg border-2 border-gray-300 dark:border-gray-600;
  @apply bg-gray-100 dark:bg-gray-700;
  @apply p-1;
  @apply transition-all duration-300;
}

.battery-segments {
  @apply flex gap-1 h-full;
}

.battery-segment {
  @apply flex-1 rounded-sm;
  @apply transition-all duration-500 ease-out;
}

/* 각 세그먼트별 고정 색상 */
.battery-segment.segment-1 {
  @apply bg-emerald-500;
}

.battery-segment.segment-2 {
  @apply bg-yellow-500;
}

.battery-segment.segment-3 {
  @apply bg-orange-500;
}

.battery-segment.segment-4 {
  @apply bg-red-500;
}

/* 활성 상태 */
.battery-segment.active {
  @apply opacity-100 scale-100;
  @apply shadow-inner;
  filter: none;
}

/* 비활성 상태 - 회색 블러 처리 */
.battery-segment.inactive {
  @apply opacity-40;
  @apply bg-gray-400 dark:bg-gray-500;
  filter: blur(1px) grayscale(80%);
  transform: scale(0.95);
}

/* 범례 */
.chart-legend {
  @apply px-4 py-3 border-t border-gray-100 dark:border-gray-700;
  @apply bg-gray-50/30 dark:bg-gray-800/30;
}

.legend-title {
  @apply text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2;
}

.legend-items {
  @apply flex flex-wrap justify-center gap-4;
}

.legend-item {
  @apply flex items-center gap-2;
}

.legend-color {
  @apply w-3 h-3 rounded-full border;
}

.status-stopped { @apply bg-gray-400 border-gray-500; }
.status-normal { @apply bg-emerald-500 border-emerald-600; }
.status-warning { @apply bg-yellow-500 border-yellow-600; }
.status-alert { @apply bg-orange-500 border-orange-600; }
.status-danger { @apply bg-red-500 border-red-600; }

.legend-label {
  @apply text-xs font-medium;
}

/* 호버 효과 */
.battery-chart-container:hover {
  @apply shadow-xl transition-shadow duration-300;
}

/* 애니메이션 */
@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 반응형 */
@media (max-width: 768px) {
  .battery-grid {
    grid-template-columns: repeat(4, 1fr) !important;
    gap: 12px !important;
  }
}

@media (max-width: 640px) {
  .battery-grid {
    grid-template-columns: repeat(3, 1fr) !important;
    gap: 10px !important;
  }
  
  .battery-body {
    @apply w-16 h-8;
  }
  
  .chart-legend {
    @apply gap-2;
  }
  
  .legend-label {
    @apply text-[10px];
  }
}
</style>