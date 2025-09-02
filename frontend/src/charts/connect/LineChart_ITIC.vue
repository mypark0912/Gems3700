<template>
    <div class="px-5 py-3">
      <div class="flex flex-wrap justify-between items-end gap-y-2 gap-x-4">
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
    Chart, LineController, LineElement, Filler, PointElement, LinearScale, TimeScale, Tooltip, LogarithmicScale,ScatterController
  } from 'chart.js'
  import 'chartjs-adapter-moment'
  
  // Import utilities
  import { tailwindConfig, formatValue } from '../../utils/Utils'
  
  Chart.register(LineController, LineElement, Filler, PointElement, LinearScale, TimeScale, Tooltip, LogarithmicScale,ScatterController)
  
  export default {
    name: 'LineChart_ITIC',
    props: ['data', 'width', 'height'],
    setup(props) {
  
      const canvas = ref(null)
      const legend = ref(null)
      let chart = null
      const darkMode = useDark()
      const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors
      
      onMounted(() => {
        const ctx = canvas.value
        chart = new Chart(ctx, {
          type: 'line',
          data: props.data,
          options: {
            layout: {
              padding: 20,
            },
            scales: {
              x: {
                type: 'logarithmic',        // 🔍 로그 스케일
                min: 0.0001,
                max: 100,
                ticks: {
                  color:darkMode.value ? textColor.dark : textColor.light,
                  callback: val => val,
                },
                grid: {
                  color: darkMode.value ? gridColor.dark : gridColor.light,
                },
              },
              y: {
                min: 0,
                max: 500,
                ticks:{
                  color:darkMode.value ? textColor.dark : textColor.light,
                },
                grid: {
                  color: darkMode.value ? gridColor.dark : gridColor.light,
                },
              }
            },
            plugins: {
              legend: {
                display: false,
              },
              tooltip: {
                callbacks: {
                  title: () => false, // Disable tooltip title
                  label: (context) => context.parsed.y, // ✅ 툴팁도 숫자로 출력
                },
                bodyColor: darkMode.value ? tooltipBodyColor.dark : tooltipBodyColor.light,
                backgroundColor: darkMode.value ? tooltipBgColor.dark : tooltipBgColor.light,
                borderColor: darkMode.value ? tooltipBorderColor.dark : tooltipBorderColor.light,
                textColor:darkMode.value ? textColor.dark : textColor.light,
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
      })
  
  
      onUnmounted(() => chart.destroy())

      watch(
        () => props.data,
        (newData) => {
          if (chart) {
            chart.data = newData
            chart.update()
          }
        },
        { deep: true }
      )
  
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
      }
    }
  }
  </script>