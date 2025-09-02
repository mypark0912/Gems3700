<template>
    <div class="modern-chart-container">
      <!-- 헤더 섹션 -->
      <div class="chart-header">
        <div class="header-content">
          <h3 class="chart-title">진단 세부항목 상태</h3>
        </div>
      </div>
  
      <!-- 히트맵 컨테이너 -->
      <div class="heatmap-wrapper">
        <div class="heatmap-grid" :style="gridStyle">
          <div
            v-for="(item, index) in heatmapData"
            :key="index"
            :class="['heatmap-cell', getStatusClass(item.value)]"
            @mouseenter="showTooltip($event, item, index)"
            @mouseleave="hideTooltip"
          >
            <div class="cell-title">{{ item.title }}</div>
            <div class="cell-status">{{ getStatusText(item.value) }}</div>
          </div>
        </div>
        
        <!-- 툴팁 -->
        <div 
          v-if="tooltip.show" 
          class="tooltip"
          :style="tooltipStyle"
        >
          <div class="tooltip-title">{{ tooltip.title }}</div>
          <div class="tooltip-status">상태: {{ tooltip.status }}</div>
        </div>
        
        <!-- 로딩 오버레이 -->
        <div v-if="isLoading" class="loading-overlay">
          <div class="loading-spinner"></div>
          <span class="loading-text">데이터 로딩 중...</span>
        </div>
      </div>
  
      <!-- 범례 -->
      <div class="legend">
        <div class="legend-title">상태 범례</div>
        <div class="legend-items">
          <div class="legend-item">
            <div class="legend-color status-disabled"></div>
            <span>{{ t('diagnosis.tabContext.st0') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color status-normal"></div>
            <span>{{ t('diagnosis.tabContext.st1') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color status-warning"></div>
            <span>{{ t('diagnosis.tabContext.st2') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color status-alert"></div>
            <span>{{ t('diagnosis.tabContext.st3') }}</span>
          </div>
          <div class="legend-item">
            <div class="legend-color status-error"></div>
            <span>{{ t('diagnosis.tabContext.st4') }}</span>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted } from 'vue'
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
  
      // 히트맵 데이터 변환
      const heatmapData = computed(() => {
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
  
      // 그리드 스타일 계산 (가변 열 개수)
      const gridStyle = computed(() => {
        const itemCount = heatmapData.value.length
        let columns = 5 // 기본 4열
        
        // 항목 개수에 따른 최적 열 개수 계산
        if (itemCount <= 4) columns = 3
        else if (itemCount <= 9) columns = 4
        else if (itemCount <= 16) columns = 5
        else columns = 6
        
        return {
          'grid-template-columns': `repeat(${columns}, 1fr)`,
          'row-gap': '12px',
          'column-gap': '8px'
        }
      })
  
      // 상태별 CSS 클래스 반환
      const getStatusClass = (value) => {
        const statusMap = {
          0: 'status-disabled',
          1: 'status-normal',
          2: 'status-warning', 
          3: 'status-alert',
          4: 'status-error'
        }
        return statusMap[Math.floor(value)] || 'status-disabled'
      }
  
      // 상태 텍스트 반환
      const getStatusText = (value) => {
        const statusLabels = [
          t('diagnosis.tabContext.st0'),
          t('diagnosis.tabContext.st1'),
          t('diagnosis.tabContext.st2'),
          t('diagnosis.tabContext.st3'),
          t('diagnosis.tabContext.st4')
        ]
        return statusLabels[Math.floor(value)] || statusLabels[0]
      }
  
      // 툴팁 표시
      const showTooltip = (event, item, index) => {
        const rect = event.target.getBoundingClientRect()
        tooltip.value = {
          show: true,
          title: item.title,
          status: getStatusText(item.value),
          x: rect.left + rect.width / 2,
          y: rect.top - 10
        }
      }
  
      // 툴팁 숨기기
      const hideTooltip = () => {
        tooltip.value.show = false
      }
  
      // 툴팁 위치 스타일
      const tooltipStyle = computed(() => ({
        left: `${tooltip.value.x}px`,
        top: `${tooltip.value.y}px`,
        transform: 'translate(-50%, -100%)'
      }))
  
      onMounted(() => {
        // 애니메이션 효과를 위한 초기화
        isLoading.value = false
      })
  
      return {
        isLoading,
        heatmapData,
        gridStyle,
        tooltip,
        tooltipStyle,
        getStatusClass,
        getStatusText,
        showTooltip,
        hideTooltip,
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
    @apply p-4 pb-2;
    @apply bg-gradient-to-r from-blue-50/50 via-purple-50/30 to-pink-50/50;
    @apply dark:from-blue-900/20 dark:via-purple-900/15 dark:to-pink-900/20;
    @apply border-b border-gray-200/50 dark:border-gray-700/50;
  }
  
  .header-content {
    @apply flex items-center;
  }
  
  .chart-title {
    @apply text-xl font-bold;
    @apply bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent;
    @apply dark:from-gray-100 dark:via-blue-200 dark:to-purple-200;
  }
  
  /* 히트맵 래퍼 */
  .heatmap-wrapper {
    @apply relative px-4 pt-4 pb-2;
    min-height: 300px;
  }
  
  .heatmap-grid {
    @apply grid w-full;
    @apply gap-x-2 gap-y-3;
  }
  
  /* 히트맵 셀 */
  .heatmap-cell {
    @apply relative p-3 rounded-xl border-2 cursor-pointer;
    @apply transition-all duration-300 ease-in-out;
    @apply flex flex-col items-center justify-center text-center;
    @apply min-h-[80px] max-w-[280px] w-full;
    @apply transform hover:scale-105 hover:shadow-lg;
    animation: fadeInUp 0.6s ease-out forwards;
    opacity: 0;
    transform: translateY(20px);
  }
  
  .heatmap-cell:nth-child(1) { animation-delay: 0.1s; }
  .heatmap-cell:nth-child(2) { animation-delay: 0.2s; }
  .heatmap-cell:nth-child(3) { animation-delay: 0.3s; }
  .heatmap-cell:nth-child(4) { animation-delay: 0.4s; }
  .heatmap-cell:nth-child(5) { animation-delay: 0.5s; }
  .heatmap-cell:nth-child(6) { animation-delay: 0.6s; }
  .heatmap-cell:nth-child(n+7) { animation-delay: 0.7s; }
  
  .cell-title {
    @apply text-sm font-semibold mb-1;
    @apply text-white drop-shadow-sm;
  }
  
  .cell-status {
    @apply text-xs font-medium;
    @apply text-white/90;
  }
  
  /* 상태별 스타일 */
  .status-disabled {
    @apply bg-gradient-to-br from-slate-400 to-slate-600;
    @apply border-slate-500 shadow-slate-200 dark:shadow-slate-800;
  }
  
  .status-normal {
    @apply bg-gradient-to-br from-emerald-400 to-emerald-600;
    @apply border-emerald-500 shadow-emerald-200 dark:shadow-emerald-800;
  }
  
  .status-warning {
    @apply bg-gradient-to-br from-amber-400 to-amber-600;
    @apply border-amber-500 shadow-amber-200 dark:shadow-amber-800;
  }
  
  .status-alert {
    @apply bg-gradient-to-br from-orange-400 to-orange-600;
    @apply border-orange-500 shadow-orange-200 dark:shadow-orange-800;
  }
  
  .status-error {
    @apply bg-gradient-to-br from-red-400 to-red-600;
    @apply border-red-500 shadow-red-200 dark:shadow-red-800;
  }
  
  /* 툴팁 */
  .tooltip {
    @apply fixed z-50 px-3 py-2 rounded-lg shadow-lg;
    @apply bg-gray-900 dark:bg-gray-800 text-white;
    @apply text-sm font-medium;
    @apply border border-gray-700 dark:border-gray-600;
    @apply pointer-events-none;
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
  @media (max-width: 640px) {
    .heatmap-grid {
      @apply gap-2;
      grid-template-columns: repeat(2, 1fr) !important;
    }
    
    .heatmap-cell {
      @apply min-h-[70px] p-2;
    }
    
    .cell-title {
      @apply text-xs;
    }
    
    .cell-status {
      @apply text-xs;
    }
    
    .legend-items {
      @apply gap-2;
    }
  }
  
  @media (max-width: 480px) {
    .heatmap-grid {
      grid-template-columns: repeat(1, 1fr) !important;
    }
  }
  </style>