<template>
  <div class="px-5 py-3">
    <div class="flex flex-wrap justify-between items-end gap-y-2 gap-x-4">
      <div class="flex items-center">
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100 mr-2">{{mode}}</div>
      </div>
      <div class="grow mb-1">
        <ul ref="legend" class="flex flex-wrap gap-x-4 sm:justify-end"></ul>
      </div>
    </div>
  </div>
  <div class="grow">
    <canvas ref="canvas" :data="data" :width="width" :height="height"></canvas>
  </div>  
</template>

<script>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useDark } from '@vueuse/core'
import { chartColors } from '../ChartjsConfig'

import {
  Chart, BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend,
} from 'chart.js'
import 'chartjs-adapter-moment'

// Import utilities
import { tailwindConfig, formatValue } from '../../utils/Utils'

Chart.register(BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend)

export default {
  name: 'BarChart05',
  props: ['data', 'width', 'height', 'mode'],
  setup(props) {

    const canvas = ref(null)
    const legend = ref(null)
    const mode = ref(props.mode)
    //console.log(mode.value);
    let chart = null
    const darkMode = useDark()
    const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors
    
    onMounted(() => {
      const ctx = canvas.value
      chart = new Chart(ctx, {
        type: 'bar',
        data: props.data,
        options: {
          layout: {
            padding: {
              top: 12,
              bottom: 16,
              left: 20,
              right: 20,
            },
          },
          scales: {
            y: {
              border: {
                display: false,
              },
              ticks: {
                maxTicksLimit: 5,
                callback: (value) => Number(value), // ✅ 숫자만 출력되도록 보장
                color: darkMode.value ? textColor.dark : textColor.light,
              },
              grid: {
                color: darkMode.value ? gridColor.dark : gridColor.light,
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
                color: darkMode.value ? textColor.dark : textColor.light,
              },
            },
          },
          plugins: {
            legend: {
              display: false, // ✅ 범례 완전 제거
            },
            title: {
              display: false, // ✅ 제목 숨기기
              text:'',
            },
            tooltip: {
              callbacks: {
                title: () => '', // ✅ 툴팁 제목 완전 제거
                label: (context) => Number(context.parsed.y), // ✅ 숫자만 표시
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
          animation: {
            duration: 500,
          },
          maintainAspectRatio: false,
          resizeDelay: 200,
        },
        plugins: [] // ✅ 플러그인 제거 → 범례(htmlLegend) 관련 코드 삭제
      });
    });



    onUnmounted(() => chart.destroy())

    watch(
      () => darkMode.value,
      () => {
        if (darkMode.value) {
          chart.options.scales.x.ticks.color = textColor.dark
          chart.options.scales.y.ticks.color = textColor.dark
          chart.options.scales.y.grid.color = gridColor.dark
          chart.options.plugins.tooltip.bodyColor = tooltipBodyColor.dark
          chart.options.plugins.tooltip.backgroundColor = tooltipBgColor.dark
          chart.options.plugins.tooltip.borderColor = tooltipBorderColor.dark
        } else {
          chart.options.scales.x.ticks.color = textColor.light
          chart.options.scales.y.ticks.color = textColor.light
          chart.options.scales.y.grid.color = gridColor.light
          chart.options.plugins.tooltip.bodyColor = tooltipBodyColor.light
          chart.options.plugins.tooltip.backgroundColor = tooltipBgColor.light
          chart.options.plugins.tooltip.borderColor = tooltipBorderColor.light
        }
        chart.update('none')
      })     

    return {
      canvas,
      legend,
      mode,
    }
  }
}
</script>