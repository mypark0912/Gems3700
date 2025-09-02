<template>
  <div class="modern-chart-container">
    <!-- 헤더 섹션 -->
    <div class="chart-header">
      <div class="header-content">
        <h3 class="chart-title">{{t('diagnosis.tabTitle.detailTitle')}}</h3>
      </div>
    </div>

    <!-- 배터리 컨테이너 -->
    <div class="battery-wrapper">
      <div class="battery-grid" :style="gridStyle">
        <div
          v-for="(item, index) in batteryData"
          :key="index"
          class="battery-item"
          @mouseenter="showTooltip($event, item, index)"
          @mouseleave="hideTooltip"
          @mousemove="updateTooltipPosition($event)"
        >
          <!-- 항목 타이틀 -->
          <div class="item-title">{{ item.title }}</div>
          
          <!-- 배터리 컨테이너 -->
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
                      'active': segment <= item.value && item.value > 0,
                      'inactive': segment > item.value || item.value === 0
                    }
                  ]"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 툴팁 -->
      <div 
        v-if="tooltip.show" 
        class="custom-tooltip"
        :style="tooltipStyle"
      >
        <div class="tooltip-content">
          <div class="tooltip-title">{{ tooltip.title }}</div>
          <div class="tooltip-status">{{ tooltip.status }}</div>
        </div>
      </div>
      
      <!-- 로딩 오버레이 -->
      <div v-if="isLoading" class="loading-overlay">
        <div class="loading-spinner"></div>
        <span class="loading-text">Loading...</span>
      </div>
    </div>

    <!-- 범례 -->
    <div class="legend">
      <div class="legend-title">{{t('diagnosis.tabContext.legend')}}</div>
      <div class="legend-items">
        <div class="legend-item">
          <div class="legend-color status-stopped"></div>
          <span>{{t('diagnosis.tabContext.st0')}}</span>
        </div>
        <div class="legend-item">
          <div class="legend-color status-normal"></div>
          <span>{{t('diagnosis.tabContext.st1')}}</span>
        </div>
        <div class="legend-item">
          <div class="legend-color status-warning"></div>
          <span>{{t('diagnosis.tabContext.st2')}}</span>
        </div>
        <div class="legend-item">
          <div class="legend-color status-alert"></div>
          <span>{{t('diagnosis.tabContext.st3')}}</span>
        </div>
        <div class="legend-item">
          <div class="legend-color status-danger"></div>
          <span>{{t('diagnosis.tabContext.st4')}}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ModernHeatmapChart',
  props: ['data', 'width', 'height'],
  setup(props) {
    const { t, locale } = useI18n()
    const isLoading = ref(false)
    const darkMode = useDark()
    
    // 툴팁 상태
    const tooltip = ref({
      show: false,
      title: '',
      status: '',
      x: 0,
      y: 0
    })

    // 배터리 데이터 변환
    const batteryData = computed(() => {
      if (!props.data?.labels || !props.data?.datasets?.[0]?.data) return []
      
      return props.data.labels.map((label, index) => {
        const value = props.data.datasets[0].data[index]
        const title = props.data.titles?.[index]?.[locale.value] || label
        
        return {
          title: title,
          value: value === 0.2 ? 0 : value, // 0.2는 실제로 0으로 처리
          originalIndex: index
        }
      })
    })

    // 그리드 스타일 계산
    const gridStyle = computed(() => {
      const itemCount = batteryData.value.length
      let columns = 4 // 기본 4열
      
      if (itemCount <= 4) columns = 4
      else if (itemCount <= 8) columns = 4
      else if (itemCount <= 12) columns = 4
      else columns = 5
      
      return {
        'grid-template-columns': `repeat(${columns}, 1fr)`,
        'gap': '20px'
      }
    })

    // 상태 텍스트 반환
    const getStatusText = (value) => {
      const statusLabels = [
        t('diagnosis.tabContext.st0'), 
        t('diagnosis.tabContext.st1'), 
        t('diagnosis.tabContext.st2'), 
        t('diagnosis.tabContext.st3'), 
        t('diagnosis.tabContext.st4')
      ]
      return statusLabels[Math.floor(value)] || t('diagnosis.tabContext.st0')
    }

    // 툴팁 표시
    const showTooltip = (event, item, index) => {
      const batteryItem = event.currentTarget
      const containerRect = batteryItem.closest('.battery-wrapper').getBoundingClientRect()
      const itemRect = batteryItem.getBoundingClientRect()
      
      // 컨테이너 기준 상대 위치 계산
      tooltip.value = {
        show: true,
        title: item.title,
        status: getStatusText(item.value),
        x: itemRect.left - containerRect.left + itemRect.width / 2,
        y: itemRect.top - containerRect.top - 10
      }
    }

    // 마우스 이동 시 툴팁 위치 업데이트
    const updateTooltipPosition = (event) => {
      if (tooltip.value.show) {
        const batteryItem = event.currentTarget
        const containerRect = batteryItem.closest('.battery-wrapper').getBoundingClientRect()
        
        // 컨테이너 기준 마우스 위치
        tooltip.value.x = event.clientX - containerRect.left
        tooltip.value.y = event.clientY - containerRect.top - 40
      }
    }

    // 툴팁 숨기기
    const hideTooltip = () => {
      tooltip.value.show = false
    }

    // 툴팁 위치 스타일
    const tooltipStyle = computed(() => {
      return {
        left: `${tooltip.value.x}px`,
        top: `${tooltip.value.y}px`,
        transform: 'translate(-50%, -100%)'
      }
    })

    // 스크롤 이벤트 핸들러
    const handleScroll = () => {
      if (tooltip.value.show) {
        hideTooltip()
      }
    }

    onMounted(() => {
      // 애니메이션 효과를 위한 초기화
      isLoading.value = false
      
      // 스크롤 시 툴팁 숨기기
      window.addEventListener('scroll', handleScroll, true)
    })

    onUnmounted(() => {
      window.removeEventListener('scroll', handleScroll, true)
    })

    return {
      isLoading,
      batteryData,
      gridStyle,
      tooltip,
      tooltipStyle,
      getStatusText,
      showTooltip,
      hideTooltip,
      updateTooltipPosition,
      t,
      locale,
    }
  }
}
</script>

<style scoped>
.modern-chart-container {
  @apply bg-gradient-to-br from-white via-gray-50 to-white dark:from-gray-800 dark:via-gray-900 dark:to-gray-800;
  @apply rounded-2xl border border-gray-200/50 dark:border-gray-700/50 mt-2;
  @apply overflow-hidden backdrop-blur-sm;
}

/* 헤더 */
.chart-header {
 @apply px-5 py-4;
 @apply bg-blue-50 dark:bg-gray-900;
 @apply border-b border-blue-100 dark:border-gray-700;
}

.header-content {
 @apply flex items-center gap-3;
}

.header-icon {
 @apply w-10 h-10 bg-blue-500 rounded-lg;
 @apply flex items-center justify-center;
 @apply text-white shadow-md;
}

.chart-title {
 @apply text-lg font-bold;
 @apply text-gray-800 dark:text-gray-100;
}

/* 배터리 래퍼 */
.battery-wrapper {
  @apply relative px-4 pt-4 pb-2;
  min-height: 300px;
}

.battery-grid {
  @apply grid w-full;
}

/* 배터리 아이템 */
.battery-item {
  @apply flex flex-col items-center p-4 rounded-lg cursor-pointer;
  @apply transition-all duration-300 ease-in-out;
  @apply transform hover:scale-105 hover:shadow-md;
  @apply bg-white/60 dark:bg-gray-800/60 backdrop-blur-sm;
  @apply border border-gray-200/30 dark:border-gray-700/30;
  @apply min-h-[100px];
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

/* 아이템 타이틀 - 포인터 이벤트 무시 */
.item-title {
  @apply text-sm font-semibold text-gray-800 dark:text-gray-200;
  @apply text-center mb-3;
  @apply leading-tight;
  pointer-events: none;
}

/* 배터리 - 포인터 이벤트 무시 */
.battery-container {
  @apply relative;
  pointer-events: none;
}

.battery-body {
  @apply w-24 h-12 rounded-lg border-2 border-gray-300 dark:border-gray-600;
  @apply bg-gray-100 dark:bg-gray-700;
  @apply p-1.5;
  @apply transition-all duration-300;
}

.battery-segments {
  @apply flex gap-1.5 h-full;
}

.battery-segment {
  @apply flex-1 rounded-sm;
  @apply transition-all duration-500 ease-out;
  @apply relative;
}

/* 각 세그먼트별 고정 색상 */
.battery-segment.segment-1 {
  @apply bg-green-500; /* 정상 - 녹색 */
}

.battery-segment.segment-2 {
  @apply bg-yellow-500; /* 주의 - 노란색 */
}

.battery-segment.segment-3 {
  @apply bg-orange-500; /* 경고 - 오렌지색 */
}

.battery-segment.segment-4 {
  @apply bg-red-500; /* 위험 - 빨간색 */
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

/* 툴팁 */
.custom-tooltip {
  @apply absolute z-[9999] px-3 py-2 rounded-lg shadow-lg;
  @apply bg-gray-900 dark:bg-gray-800 text-white;
  @apply text-sm font-medium;
  @apply border border-gray-700 dark:border-gray-600;
  @apply pointer-events-none;
  @apply whitespace-nowrap;
  max-width: 200px;
  margin-top: -10px;
}

.custom-tooltip::after {
  content: '';
  @apply absolute top-full left-1/2 transform -translate-x-1/2;
  @apply border-4 border-transparent border-t-gray-900 dark:border-t-gray-800;
}

.tooltip-content {
  @apply relative z-10;
}

.tooltip-title {
  @apply font-semibold;
}

.tooltip-status {
  @apply text-gray-300 text-xs mt-1;
}

/* 범례 */
.legend {
  @apply px-4 pb-3 pt-2 border-t border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gray-50/30 dark:bg-gray-800/30;
}

.legend-title {
  @apply text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2;
}

.legend-items {
  @apply flex flex-wrap gap-4;
}

.legend-item {
  @apply flex items-center gap-2;
}

.legend-color {
  @apply w-3 h-3 rounded-full border;
}

.status-stopped { @apply bg-gray-400 border-gray-500; }
.status-normal { @apply bg-green-500 border-green-600; }
.status-warning { @apply bg-yellow-500 border-yellow-600; }
.status-alert { @apply bg-orange-500 border-orange-600; }
.status-danger { @apply bg-red-500 border-red-600; }

.legend-item span {
  @apply text-xs text-gray-600 dark:text-gray-400;
}

/* 로딩 오버레이 */
.loading-overlay {
  @apply absolute inset-0 flex flex-col items-center justify-center;
  @apply bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm;
  @apply rounded-xl;
}

.loading-spinner {
  @apply w-8 h-8 border-3 border-blue-200 border-t-blue-600 rounded-full;
  animation: spin 1s linear infinite;
}

.loading-text {
  @apply mt-3 text-sm font-medium text-gray-600 dark:text-gray-400;
}

/* 애니메이션 */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 반응형 */
@media (max-width: 768px) {
  .battery-grid {
    grid-template-columns: repeat(3, 1fr) !important;
    @apply gap-4;
  }
}

@media (max-width: 640px) {
  .battery-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    @apply gap-3;
  }
  
  .battery-item {
    @apply p-3;
  }
  
  .legend-items {
    @apply gap-2;
  }
}

@media (max-width: 480px) {
  .battery-grid {
    grid-template-columns: repeat(1, 1fr) !important;
  }
  
  .battery-wrapper {
    @apply px-3;
  }
  
  .chart-header {
    @apply p-3;
  }
}
</style>