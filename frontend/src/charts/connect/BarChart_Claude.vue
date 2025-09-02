<template>
  <div class="modern-chart-container">
    <!-- 헤더 섹션 -->
    <div class="chart-header">
      <div class="header-content">
        <h3 class="chart-title">진단 세부항목 상태</h3>
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
        <span class="loading-text">데이터 로딩 중...</span>
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
          t('diagnosis.tabContext.st0'),
          t('diagnosis.tabContext.st1'),
          t('diagnosis.tabContext.st2'),
          t('diagnosis.tabContext.st3'),
          t('diagnosis.tabContext.st4')
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
                    t('diagnosis.tabContext.st0'),
                    t('diagnosis.tabContext.st1'),
                    t('diagnosis.tabContext.st2'),
                    t('diagnosis.tabContext.st3'),
                    t('diagnosis.tabContext.st4')
                  ]
                  const actualValue = value === 0.2 ? 0 : value
                  return `상태: ${statusLabels[Math.floor(actualValue)]}`
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

/* 차트 래퍼 */
.chart-wrapper {
  @apply relative px-3 py-2;
  min-height: 450px;
}

.chart-background {
  @apply relative;
  @apply bg-gradient-to-b from-transparent to-gray-50/30 dark:to-gray-800/30;
  @apply rounded-xl;
  @apply backdrop-blur-sm;
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

/* 반응형 */
@media (max-width: 640px) {
  .chart-header,
  .chart-wrapper {
    @apply p-4;
  }
}
</style>