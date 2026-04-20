<template>
  <div 
    class="modern-chart-container"
    :class="isPdfMode ? 'bg-white border-gray-200' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700'"
  >
    <!-- 헤더 섹션 -->
    <div 
      class="chart-header"
      :class="isPdfMode ? 'bg-blue-50 border-blue-100' : 'bg-blue-50 dark:bg-gray-900 border-blue-100 dark:border-gray-700'"
    >
      <div class="header-content">
        <h3 
          class="chart-title"
          :class="isPdfMode ? 'text-gray-800' : 'text-gray-800 dark:text-gray-100'"
        >
          {{ chartTitle }}
        </h3>
        <span class="update-time"> {{data_state }} / {{ data_recordtime }}</span>
      </div>
    </div>

    <!-- 차트 컨테이너 -->
    <div class="chart-wrapper">
      <div 
        class="chart-background"
        :class="isPdfMode ? 'bg-gray-50' : 'bg-gray-50 dark:bg-gray-900'"
      >
        <canvas ref="canvas" :data="data" :width="width" :height="height"></canvas>
      </div>
      
      <!-- 로딩 오버레이 -->
      <div v-if="isLoading" class="loading-overlay">
        <div class="loading-spinner"></div>
        <span class="loading-text">loading...</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted, onUnmounted, computed, inject } from 'vue'  // ✅ inject 추가
import { useDark } from '@vueuse/core'
import { chartColors } from '../ChartjsConfig'
import { useI18n } from 'vue-i18n'

import {
  Chart, BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend, CategoryScale
} from 'chart.js'
import 'chartjs-adapter-moment'

// Import utilities
import { tailwindConfig, formatValue } from '../../utils/Utils'

Chart.register(BarController, BarElement, LinearScale, CategoryScale, TimeScale, Tooltip, Legend)

export default {
  name: 'ModernBarChart',
  props: ['data', 'width', 'height','mode'],
  setup(props) {
    // ✅ PDF 모드 inject
    const isPdfMode = inject('isPdfMode', false)

    const data_recordtime_raw = inject('data_recordtime', ref(''));
    const data_state_raw = inject('data_state', ref(''));

    const data_state = computed(() => data_state_raw.value || '-');
    const data_recordtime = computed(() => data_recordtime_raw.value || '-');


    
    const { t, locale } = useI18n()
    const canvas = ref(null)
    const legend = ref(null)
    const isLoading = ref(false)
    let chart = null
    const darkMode = useDark()

    const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors

    // 모드에 따른 타이틀 계산
    const chartTitle = computed(() => {
      const titleMap = {
        'PowerQuality': t('diagnosis.tabTitle.detailTitle_pq'),
        'DiagnosisDetail': t('diagnosis.tabTitle.detailTitle'),
      }
      const result = titleMap[props.mode] || t('diagnosis.tabTitle.detailTitle_pq')
      console.log(props.mode,'Chart title:', result)
      
      return result
    })

    // 업데이트 시간 포맷팅 (국제 공통 숫자 형식)
    const getUpdateTime = () => {
      // props.data에 updateTime이나 timestamp 필드가 없으면 빈 문자열 반환
      const time = props.data?.updateTime || props.data?.timestamp
      if (!time) return '-'
      
      const date = time instanceof Date ? time : new Date(time)
      
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      const seconds = String(date.getSeconds()).padStart(2, '0')
      
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    }

    // 동적 바 굵기 계산
    const dynamicBarThickness = computed(() => {
      const itemCount = props.data?.labels?.length || 1
      
      if (itemCount <= 5) return 80
      else if (itemCount <= 10) return 45
      else if (itemCount <= 15) return 35
      else return 30
    })

    const getBarColors = (value) => {
      const colorMap = {
        0: {
          bg: 'rgba(100, 116, 139, 0.8)',
          border: 'rgba(100, 116, 139, 1)',
        },
        1: {
          bg: 'rgba(16, 185, 129, 0.8)',
          border: 'rgba(16, 185, 129, 1)',
        },
        2: {
          bg: 'rgba(245, 158, 11, 0.8)',
          border: 'rgba(245, 158, 11, 1)',
        },
        3: {
          bg: 'rgba(249, 115, 22, 0.8)',
          border: 'rgba(249, 115, 22, 1)',
        },
        4: {
          bg: 'rgba(239, 68, 68, 0.8)',
          border: 'rgba(239, 68, 68, 1)',
        }
      }
      return colorMap[Math.floor(value)] || colorMap[0]
    }

    watch(locale, () => {
      if (!chart) return
      chart.options.scales.y.ticks.callback = (value) => {
        const textstr = [
          t('diagnosis.tabContext.pqfe0'),
          t('diagnosis.tabContext.pqfe1'),
          t('diagnosis.tabContext.pqfe2'),
          t('diagnosis.tabContext.pqfe3'),
          t('diagnosis.tabContext.pqfe4')
        ]
        return textstr[Math.floor(value)]
      }

      chart.options.scales.x.ticks.callback = (value, index) => {
        const title = props.data.titles?.[index]?.[locale.value]
        const label = (title && title.trim() !== '') ? title : chart.data.labels?.[index]
        return label
      }

      chart.update('none')
    })
    
    onMounted(() => {
      const ctx = canvas.value
      chart = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: props.data.labels,
          datasets: [{
            label: 'Diagnosis',
            data: props.data.datasets[0].data.map(value => value === 0 ? 0.2 : value),
            backgroundColor: props.data.datasets[0].data.map(value => getBarColors(value).bg),
            borderColor: props.data.datasets[0].data.map(value => getBarColors(value).border),
            borderWidth: 2,
            borderRadius: {
              topLeft: 8,
              topRight: 8,
              bottomLeft: 0,
              bottomRight: 0,
            },
            borderSkipped: false,
            barThickness: dynamicBarThickness.value,
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          layout: {
            padding: {
              top: 30,
              bottom: 8,
              left: 10,
              right: 10,
            },
          },
          scales: {
            y: {
              max: 4,
              border: {
                display: false,
              },
              grid: {
                display: true,
                color: (context) => {
                  const colors = [
                    'rgba(100, 116, 139, 0.1)',
                    'rgba(16, 185, 129, 0.1)',
                    'rgba(245, 158, 11, 0.1)',
                    'rgba(249, 115, 22, 0.1)',
                    'rgba(239, 68, 68, 0.1)'
                  ]
                  return colors[context.tick.value] || 'rgba(156, 163, 175, 0.1)'
                },
                lineWidth: 1,
                drawBorder: false,
              },
              ticks: {
                stepSize: 1,
                padding: 8,
                font: {
                  size: 12,
                  weight: '600',
                },
                callback: (value) => {
                  const textstr = [
                    t('diagnosis.tabContext.pqfe0'),
                    t('diagnosis.tabContext.pqfe1'),
                    t('diagnosis.tabContext.pqfe2'),
                    t('diagnosis.tabContext.pqfe3'),
                    t('diagnosis.tabContext.pqfe4')
                  ]
                  return textstr[Math.floor(value)]
                }, 
                color: (context) => {
                  const colors = [
                    '#64748B', '#10B981', '#F59E0B', '#F97316', '#EF4444'
                  ]
                  return colors[context.tick.value] || '#6B7280'
                },
              },
            },
            x: {
              border: {
                display: false,
              },
              grid: {
                display: false,
              },
              ticks: {
                padding: 8,
                font: {
                  size: 12,
                  weight: '600',
                },
                // ✅ PDF 모드일 때는 다크모드 무시
                color: isPdfMode ? '#6B7280' : (darkMode.value ? '#D1D5DB' : '#6B7280'),
                maxRotation: 45,
                minRotation: 0,
                callback: (value, index) => {
                  const title = props.data.titles?.[index]?.[locale.value]
                  const label = (title && title.trim() !== '') ? title : props.data.labels?.[index]
                  return label
                },
              },
            },
          },
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              enabled: true,
              // ✅ PDF 모드일 때는 라이트 테마 강제
              backgroundColor: isPdfMode ? 'rgba(255, 255, 255, 0.95)' : (darkMode.value ? 'rgba(17, 24, 39, 0.95)' : 'rgba(255, 255, 255, 0.95)'),
              titleColor: isPdfMode ? '#111827' : (darkMode.value ? '#F9FAFB' : '#111827'),
              bodyColor: isPdfMode ? '#374151' : (darkMode.value ? '#E5E7EB' : '#374151'),
              borderColor: isPdfMode ? '#E5E7EB' : (darkMode.value ? '#4B5563' : '#E5E7EB'),
              borderWidth: 1,
              cornerRadius: 12,
              padding: 12,
              displayColors: false,
              titleFont: {
                size: 13,
                weight: '600',
              },
              bodyFont: {
                size: 12,
                weight: '500',
              },
              callbacks: {
                title: (context) => {
                  return context[0].label
                },
                label: (context) => {
                  const value = context.parsed.y
                  const statusLabels = [
                    t('diagnosis.tabContext.pqfe0'),
                    t('diagnosis.tabContext.pqfe1'),
                    t('diagnosis.tabContext.pqfe2'),
                    t('diagnosis.tabContext.pqfe3'),
                    t('diagnosis.tabContext.pqfe4')
                  ]
                  const actualValue = value === 0.2 ? 0 : value
                  return `${statusLabels[Math.floor(actualValue)]}`
                },
              },
            },
          },
          interaction: {
            intersect: false,
            mode: 'index',
          },
          animation: {
            duration: 1200,
            easing: 'easeInOutQuart',
            delay: (context) => {
              return context.dataIndex * 100
            }
          },
          hover: {
            animationDuration: 300,
          },
          onHover: (event, elements) => {
            event.native.target.style.cursor = elements.length > 0 ? 'pointer' : 'default'
          }
        },
      })
    })

    onUnmounted(() => chart?.destroy())

    watch(
      () => darkMode.value,
      () => {
        if (!chart || isPdfMode) return  // ✅ PDF 모드일 때는 다크모드 변경 무시
        
        if (darkMode.value) {
          chart.options.scales.x.ticks.color = '#D1D5DB'
          chart.options.plugins.tooltip.backgroundColor = 'rgba(17, 24, 39, 0.95)'
          chart.options.plugins.tooltip.titleColor = '#F9FAFB'
          chart.options.plugins.tooltip.bodyColor = '#E5E7EB'
          chart.options.plugins.tooltip.borderColor = '#4B5563'
        } else {
          chart.options.scales.x.ticks.color = '#6B7280'
          chart.options.plugins.tooltip.backgroundColor = 'rgba(255, 255, 255, 0.95)'
          chart.options.plugins.tooltip.titleColor = '#111827'
          chart.options.plugins.tooltip.bodyColor = '#374151'
          chart.options.plugins.tooltip.borderColor = '#E5E7EB'
        }
        chart.update('none')
      }
    )    

    return {
      canvas,
      legend,
      isLoading,
      dynamicBarThickness,
      isPdfMode,  // ✅ 반환 목록에 추가
      getUpdateTime,
      t,
      locale,
      chartTitle,
      data_state,
      data_recordtime,
    }
  }
}
</script>

<style scoped>
.modern-chart-container {
 @apply rounded-xl border mt-2;
 @apply overflow-hidden shadow-lg;
}

/* 헤더 */
.chart-header {
 @apply px-5 py-4;
 @apply border-b;
}

.header-content {
 @apply flex items-center justify-between gap-3;
}

.header-icon {
 @apply w-10 h-10 bg-blue-500 rounded-lg;
 @apply flex items-center justify-center;
 @apply text-white shadow-md;
}

.chart-title {
 @apply text-lg font-bold;
}

/* 업데이트 시간 스타일 */
.update-time {
  @apply text-xs font-medium;
  @apply tabular-nums;
}

/* 차트 래퍼 */
.chart-wrapper {
 @apply relative p-4;
 min-height: 450px;
}

.chart-background {
 @apply relative;
 @apply rounded-lg p-2;
}

/* 로딩 오버레이 */
.loading-overlay {
 @apply absolute inset-0 flex flex-col items-center justify-center;
 @apply bg-white/90 dark:bg-gray-900/90;
 @apply rounded-lg;
}

.loading-spinner {
 @apply w-8 h-8 border-3 border-blue-200 border-t-blue-500 rounded-full;
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

/* 호버 효과 */
.modern-chart-container:hover {
 @apply shadow-xl transition-shadow duration-300;
}

.header-icon {
 @apply transition-transform duration-300;
}

.modern-chart-container:hover .header-icon {
 @apply transform scale-110;
}

/* 반응형 */
@media (max-width: 640px) {
 .chart-header {
   @apply px-4 py-3;
 }
 
 .chart-title {
   @apply text-base;
 }
 
 .header-icon {
   @apply w-9 h-9;
 }
 
 .update-time {
   @apply text-[10px];
 }
}
</style>