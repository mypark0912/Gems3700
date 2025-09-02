<template>
  <div class="px-5 py-3">
    <div class="flex flex-wrap justify-between items-end gap-y-2 gap-x-4">
      <div class="flex items-start">
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100 mr-2">{{ title }}</div>
      </div>
      <div class="grow mb-1">
        <ul ref="legend" class="flex flex-wrap gap-x-4 sm:justify-end"></ul>
      </div>
    </div>
  </div>
  <!-- Chart built with Chart.js 3 -->
  <div class="grow">
    <canvas ref="canvas" :data="data" :width="width" :height="height"></canvas>
  </div>  
</template>

<script>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useDark } from '@vueuse/core'
import { chartColors } from '../ChartjsConfig'

import {
  Chart, LineController, LineElement, Filler, PointElement, LinearScale, TimeScale, Tooltip,
} from 'chart.js'
import 'chartjs-adapter-moment'

// Import utilities
import { tailwindConfig, formatValue } from '../../utils/Utils'

Chart.register(LineController, LineElement, Filler, PointElement, LinearScale, TimeScale, Tooltip)

export default {
  name: 'LineChart02',
  props: ['data', 'width', 'height', 'mode'],
  setup(props) {
    const canvas = ref(null)
    const legend = ref(null)
    let chart = null
    const title = ref(props.mode);
    const darkMode = useDark()
    const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors

    const renderChart = () => {
      if (!canvas.value) return
      const ctx = canvas.value.getContext('2d')
      if (chart) chart.destroy()
      chart = new Chart(ctx, {
        type: 'line',
        data: props.data,
        options: {
          layout: { padding: 20 },
          animation: {
            duration: 800,
            easing: 'easeOutQuart',
          },
          elements: {
            point: { radius: 0, hoverRadius: 3 },
          },
          spanGaps: true,
          scales: {
            y: {
              border: { display: false },
              beginAtZero: true,
              ticks: {
                maxTicksLimit: 5,
                callback: (value) => value,
                color: darkMode.value ? textColor.dark : textColor.light,
              },
              grid: {
                color: darkMode.value ? gridColor.dark : gridColor.light,
              },
            },
            x: {
              type: 'linear',
              border: { display: false },
              grid: { display: false },
              ticks: {
                autoSkipPadding: 48,
                maxRotation: 0,
                callback: (value) => value,
                color: darkMode.value ? textColor.dark : textColor.light,
              },
            },
          },
          plugins: {
            legend: { display: true,
              onClick: false
             },
            tooltip: {
              callbacks: {
                title: () => false,
                label: (context) => context.parsed.y.toFixed(3),
              },
              bodyColor: darkMode.value ? tooltipBodyColor.dark : tooltipBodyColor.light,
              backgroundColor: darkMode.value ? tooltipBgColor.dark : tooltipBgColor.light,
              borderColor: darkMode.value ? tooltipBorderColor.dark : tooltipBorderColor.light,
            },
            datalabels: {
              display: false  // ✅ 표시 비활성화
            },
          },
          interaction: {
            intersect: false,
            mode: 'nearest',
          },
          maintainAspectRatio: false,
          resizeDelay: 200,
        },
      })
    }

    onMounted(renderChart)

    watch(
      () => props.data,
      (newData) => {
        if (chart) {
          chart.data = newData
          chart.update()
        }
      },
      { deep: true } // 🔥 중요: 데이터 객체 내부까지 감지
    )

    onUnmounted(() => {
      if (chart) chart.destroy()
    })

    watch(() => props.data, () => {
      renderChart()
    }, { deep: true })

    watch(() => darkMode.value, () => {
      renderChart()
    })

    return {
      canvas,
      legend,
      //mode: props.mode,
      title,
    }
  }
}
</script>