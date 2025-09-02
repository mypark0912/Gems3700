<template>
    <div class="px-5 py-3">
      <ul ref="legend" class="flex flex-wrap gap-x-4"></ul>
    </div>
    <div class="grow">
      <canvas ref="canvas" :data="data" :width="width" :height="height"></canvas>
    </div>
  </template>
  
  <script>
  import { ref, watch, onMounted, onUnmounted } from 'vue'
  import { useDark } from '@vueuse/core'
  import { chartColors } from '../ChartjsConfig'
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  import {
    Chart, BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend, CategoryScale
  } from 'chart.js'
  import 'chartjs-adapter-moment'
  
  // Import utilities
  import { tailwindConfig, formatValue } from '../../utils/Utils'
  
  Chart.register(BarController, BarElement, LinearScale, CategoryScale, TimeScale, Tooltip, Legend)
  
  export default {
    name: 'BarChart_PQ',
    props: ['data', 'width', 'height'],
    setup(props) {
      const { t, locale } = useI18n();
      const canvas = ref(null)
      const legend = ref(null)
      let chart = null
      const darkMode = useDark()
      const { textColor, gridColor, tooltipBodyColor, tooltipBgColor, tooltipBorderColor } = chartColors

      watch(locale, () => {
          if (!chart) return;
          chart.options.scales.y.ticks.callback = (value) => {
            const textstr = [
              t('diagnosis.tabContext.pqfe0'),
              t('diagnosis.tabContext.pqfe1'),
              t('diagnosis.tabContext.pqfe2'),
              t('diagnosis.tabContext.pqfe3'),
              t('diagnosis.tabContext.pqfe4')
            ];
            return textstr[Math.floor(value)];
          };
          chart.options.scales.x.ticks.callback = (value, index) => {
            const title = props.data.titles?.[index]?.[locale.value];
            return (title && title.trim() !== '') ? title : chart.data.labels?.[index];
          };
          chart.update('none');
        })
      
      onMounted(() => {
        const ctx = canvas.value
        chart = new Chart(ctx, {
          type: 'bar',
          data: {
            labels: props.data.labels,
            datasets: [{
              label: 'Diagnosis',
              data: props.data.datasets[0].data.map(value => value === 0 ? 0.2 : value), // 0이면 0.1로 변경
              backgroundColor: props.data.datasets[0].data.map(value => {
                // 값(0~4)에 따라 색상 지정
                const colors = [
                  'rgba(128, 128, 128, 1)', // 회색 (0: No Data)
                  'rgba(0, 128, 0, 1)',     // 녹색 (1: OK)
                  'rgba(255, 206, 86, 1)',  // 노랑 (2: Warning)
                  'rgba(255, 165, 0, 1)',   // 오렌지 (3: Inspect)
                  'rgba(255, 0, 0, 1)'      // 빨강 (4: Repair)
                ];
                return colors[Math.floor(value)] || 'rgba(200, 200, 200, 0.5)';
              }),
              borderWidth: 1
            }]
          },
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
                max: 4,
                border: {
                  display: true,
                },
                ticks: {
                  stepSize: 1,
                  callback: (value) => {
                    //const textstr = ['No Data' , 'OK' , 'Low', 'Medium', 'High'];
                    const textstr = [ t('diagnosis.tabContext.pqfe0') , t('diagnosis.tabContext.pqfe1') , t('diagnosis.tabContext.pqfe2'), t('diagnosis.tabContext.pqfe3'), t('diagnosis.tabContext.pqfe4')];
                    const index = Math.floor(value);
                    return textstr[index];
                  }, 
                  color: (context) => {
                      const value = context.tick.value;
                      const colors = [
                        'rgba(128, 128, 128, 1)', // 회색 (0~20)
                        'rgba(0, 128, 0, 1)',     // 녹색 (20~40)
                        'rgba(255, 206, 86, 1)',  // 노랑 (40~60)
                        'rgba(255, 165, 0, 1)',   // 오렌지 (60~80)
                        'rgba(255, 0, 0, 1)'      // 빨강 (80~100)
                      ];
                      const index = Math.floor(value);
                      return colors[index];
                    },
                },
                grid: {
                  drawBorder: false,
                  color: (context) => {
                    const value = context.tick.value;
                        const colors = [
                      'rgba(128, 128, 128, 1)',   // 회색 (0~20)
                      'rgba(0, 128, 0, 1)',       // 녹색 (20~40)
                      'rgba(255, 206, 86, 1)',    // 노란색 (40~60)
                      'rgba(255, 165, 0, 1)',     // 오렌지색 (60~80)
                      'rgba(255, 0, 0, 1)'        // 빨간색 (80~100)
                    ];
                    //const index = Math.floor(value / 20);
                    return colors[value] || 'rgba(200, 200, 200, 0.5)';
                  },
                  lineWidth: (context) => {
                    return context.tick.value % 1 === 0 ? 2 : 1; 
                  },
                },              
              },
              x: {
                type: 'category',
                labels: props.data.labels, 
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
                display: false,
              },
              tooltip: {
                callbacks: {
                  title: () => false, // Disable tooltip title
                  label: function(context) {
                    const value = context.parsed.y;
                    return value === 0.2 ? '0' : Number(value);
                  },
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
          /*plugins: [{
            id: 'htmlLegend',
            afterUpdate(c, args, options) {
              const ul = legend.value
              if (!ul) return
              // Remove old legend items
              while (ul.firstChild) {
                ul.firstChild.remove()
              }
              // Reuse the built-in legendItems generator
              const items = c.options.plugins.legend.labels.generateLabels(c)
              items.forEach((item) => {
                const li = document.createElement('li')
                // Button element
                const button = document.createElement('button')
                button.style.display = 'inline-flex'
                button.style.alignItems = 'center'
                button.style.opacity = item.hidden ? '.3' : ''
                button.onclick = () => {
                  c.setDatasetVisibility(item.datasetIndex, !c.isDatasetVisible(item.datasetIndex))
                  c.update()
                }
                // Color box
                const box = document.createElement('span')
                box.style.display = 'block'
                box.style.width = tailwindConfig().theme.width[3]
                box.style.height = tailwindConfig().theme.height[3]
                box.style.borderRadius = tailwindConfig().theme.borderRadius.full
                box.style.marginRight = tailwindConfig().theme.margin[2]
                box.style.borderWidth = '3px'
                box.style.borderColor = item.fillStyle
                box.style.pointerEvents = 'none'
                // Label
                const labelContainer = document.createElement('span')
                labelContainer.style.display = 'flex'
                labelContainer.style.alignItems = 'center'
                const value = document.createElement('span')
                value.classList.add('text-gray-800', 'dark:text-gray-100')
                value.style.fontSize = tailwindConfig().theme.fontSize['3xl'][0]
                value.style.lineHeight = tailwindConfig().theme.fontSize['3xl'][1].lineHeight
                value.style.fontWeight = tailwindConfig().theme.fontWeight.bold
                value.style.marginRight = tailwindConfig().theme.margin[2]
                value.style.pointerEvents = 'none'
                const label = document.createElement('span')
                label.classList.add('text-gray-500', 'dark:text-gray-400')
                label.style.fontSize = tailwindConfig().theme.fontSize.sm[0]
                label.style.lineHeight = tailwindConfig().theme.fontSize.sm[1].lineHeight
                const theValue = c.data.datasets[item.datasetIndex].data.reduce((a, b) => a + b, 0)
                const valueText = document.createTextNode(formatValue(theValue))
                const labelText = document.createTextNode(item.text)
                value.appendChild(valueText)
                label.appendChild(labelText)
                li.appendChild(button)
                button.appendChild(box)
                button.appendChild(labelContainer)
                labelContainer.appendChild(value)
                labelContainer.appendChild(label)
                ul.appendChild(li)
              })
            },
          }], */
        })
      })
  
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
      }
    }
  }
  </script>