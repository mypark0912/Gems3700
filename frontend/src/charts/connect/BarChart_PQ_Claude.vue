<template>
  <div class="modern-chart-container">
    <!-- 헤더 섹션 -->
    <div class="chart-header">
      <div class="header-content">
        <h3 class="chart-title">{{t('diagnosis.tabTitle.detailTitle_pq')}}</h3>
      </div>
    </div>

    <!-- 차트 컨테이너 -->
    <div class="chart-wrapper">
      <div class="chart-background">
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
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
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
  props: ['data', 'width', 'height'],
  setup(props) {
    const { t, locale } = useI18n()
    const canvas = ref(null)
    const legend = ref(null)
    const isLoading = ref(false)
    let chart = null
    const darkMode = useDark()
    const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors

    // 동적 바 굵기 계산
    const dynamicBarThickness = computed(() => {
      const itemCount = props.data?.labels?.length || 1
      
      // 항목 개수에 따른 바 굵기 조정
      if (itemCount <= 5) return 80      // 적은 항목일 때 굵게
      else if (itemCount <= 10) return 45  // 중간 개수
      else if (itemCount <= 15) return 35  // 많은 항목
      else return 30                       // 매우 많은 항목일 때 얇게
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
            barThickness: dynamicBarThickness.value, // 동적 바 굵기
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
            y: { // Y축이 값 축
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
                    t('diagnosis.tabContext.st0'),
                    t('diagnosis.tabContext.st1'),
                    t('diagnosis.tabContext.st2'),
                    t('diagnosis.tabContext.st3'),
                    t('diagnosis.tabContext.st4')
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
            x: { // X축이 라벨 축
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
                color: darkMode.value ? '#D1D5DB' : '#6B7280',
                maxRotation: 45,
                minRotation: 0,
                callback: (value, index) => {
                  const title = props.data.titles?.[index]?.[locale.value]
                  const label = (title && title.trim() !== '') ? title : props.data.labels?.[index]
                  // 말줄임표 없이 전체 라벨 표시 (긴 경우 회전됨)
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
              backgroundColor: darkMode.value ? 'rgba(17, 24, 39, 0.95)' : 'rgba(255, 255, 255, 0.95)',
              titleColor: darkMode.value ? '#F9FAFB' : '#111827',
              bodyColor: darkMode.value ? '#E5E7EB' : '#374151',
              borderColor: darkMode.value ? '#4B5563' : '#E5E7EB',
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
        if (!chart) return
        
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
      t,
      locale,
    }
  }
}
</script>

<style scoped>
.modern-chart-container {
 @apply bg-white dark:bg-gray-800;
 @apply rounded-xl border border-gray-200 dark:border-gray-700 mt-2;
 @apply overflow-hidden shadow-lg;
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

/* 차트 래퍼 */
.chart-wrapper {
 @apply relative p-4;
 min-height: 450px;
}

.chart-background {
 @apply relative;
 @apply bg-gray-50 dark:bg-gray-900;
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
}
</style>