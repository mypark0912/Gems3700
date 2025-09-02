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
    Chart, BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend,
  } from 'chart.js'
  import 'chartjs-adapter-moment'
  
  // Import utilities
  import { tailwindConfig, formatValue } from '../../utils/Utils'
  
  Chart.register(BarController, BarElement, LinearScale, TimeScale, Tooltip, Legend)
  
  export default {
    name: 'BarChart01',
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

        // Y축 텍스트 다시 지정
        chart.options.scales.x.ticks.callback = (value) => {
          let key = chart.data.labels[value];
          return t(`${key}`);
          //return t(`report.cardContext.energyInfo.${key}`);
        };

        // htmlLegend 강제 재실행을 위해 update
        chart.update();
      });
      
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
              x:{
                ticks: {
                  callback: (value, index) => {
                    return props.data.labels[index]; 
                    //return t(`report.cardContext.energyInfo.${props.data.labels[index]}`);
                  },
                  color: darkMode.value ? textColor.dark : textColor.light,
                }
              },
              y: {
                border: {
                  display: false,
                },
                ticks: {
                  maxTicksLimit: 5,
                  color: darkMode.value ? textColor.dark : textColor.light,
                  callback: (value) => Number(value), //formatValue(value),                
                },
                grid: {
                  color: darkMode.value ? gridColor.dark : gridColor.light,
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
                  label: (context) => parseFloat(context.parsed.y).toFixed(2)//formatValue(context.parsed.y),
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
          plugins: [{
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
                label.classList.add('text-gray-500', 'dark:text-gray-200')
                label.style.fontSize = tailwindConfig().theme.fontSize.sm[0]
                label.style.lineHeight = tailwindConfig().theme.fontSize.sm[1].lineHeight
                
                //const theValue = c.data.datasets[item.datasetIndex].data.reduce((a, b) => a + b, 0)
                //const valueText = document.createTextNode(formatValue(theValue))
                const labelText = document.createTextNode(item.text);
                //const labelText = document.createTextNode( t(`report.cardContext.energyInfo.${item.text}`))
                //value.appendChild(valueText)
                label.appendChild(labelText)
                li.appendChild(button)
                button.appendChild(box)
                button.appendChild(labelContainer)
                labelContainer.appendChild(value)
                labelContainer.appendChild(label)
                ul.appendChild(li)
              })
            },
          }],
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
        
        watch(
          () => props.data,
          (newData) => {
            if (chart) {
              chart.data = newData;
              chart.update();
            }
          },
          { deep: true }
        );
  
      return {
        canvas,
        legend,
      }
    }
  }
  </script>