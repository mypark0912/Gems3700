<template>
  <div class="en50160-section bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 mb-4">
    <!-- 섹션 헤더 -->
    <div class="flex items-center justify-between mb-4">
      <h4 class="text-base font-semibold text-gray-800 dark:text-gray-100">
        {{ title }}
      </h4>
      <div v-if="statistics" class="flex items-center gap-2">
        <span 
          :class="[
            'px-2 py-1 text-xs font-medium rounded',
            overallResult === 'PASS' 
              ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' 
              : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
          ]"
        >
          {{ overallResult }}
        </span>
      </div>
    </div>

    <!-- 라인 차트 -->
    <div class="mb-4">
      <div ref="lineChartRef" :style="{ height: chartHeight, width: '100%' }"></div>
    </div>

    <!-- 범례 -->
    <div v-if="legendItems && legendItems.length > 0" class="flex items-center gap-4 mb-4 px-2">
      <div v-for="item in legendItems" :key="item.name" class="flex items-center gap-1.5">
        <span 
          class="w-3 h-3 rounded-sm" 
          :style="{ backgroundColor: item.color }"
        ></span>
        <span class="text-xs text-gray-600 dark:text-gray-300">{{ item.name }}</span>
      </div>
    </div>

    <!-- 히스토그램 -->
    <div v-if="showHistogram && histogramData" class="mb-4 border-t border-gray-100 dark:border-gray-700 pt-4">
      <EN50160Histogram
        :title="`${title} Distribution`"
        :histogramData="histogramData"
        :unit="histogramUnit"
        :barColor="histogramColor"
        :limitMin="histogramLimitMin"
        :limitMax="histogramLimitMax"
        :limitValue="histogramLimitValue"
        :chartHeight="histogramHeight"
      />
    </div>

    <!-- 3상 히스토그램 (Voltage, THD 등) -->
    <div v-if="showPhaseHistograms && phaseHistograms" class="mb-4 border-t border-gray-100 dark:border-gray-700 pt-4">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <EN50160Histogram
          v-for="(phaseData, phaseKey) in phaseHistograms"
          :key="phaseKey"
          :title="`${phaseKey} Distribution`"
          :histogramData="phaseData"
          :unit="histogramUnit"
          :barColor="phaseColors[phaseKey] || '#5470c6'"
          :limitMin="histogramLimitMin"
          :limitMax="histogramLimitMax"
          :limitValue="histogramLimitValue"
          chartHeight="180px"
        />
      </div>
    </div>

    <!-- 테이블 -->
    <div v-if="tableData && tableData.length > 0" class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th 
              v-for="col in tableColumns" 
              :key="col.key"
              class="px-3 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-300"
            >
              {{ col.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
          <tr 
            v-for="(row, idx) in tableData" 
            :key="idx"
            :class="row.rowClass || ''"
          >
            <td 
              v-for="col in tableColumns" 
              :key="col.key"
              class="px-3 py-2 text-gray-700 dark:text-gray-200"
            >
              <template v-if="col.key === 'result'">
                <span 
                  :class="[
                    'px-2 py-0.5 text-xs font-medium rounded',
                    row[col.key] === 'PASS' 
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' 
                      : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                  ]"
                >
                  {{ row[col.key] }}
                </span>
              </template>
              <template v-else>
                {{ row[col.key] }}
              </template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch, onBeforeUnmount, computed } from 'vue'
import * as echarts from 'echarts'
import EN50160Histogram from './EN50160Histogram.vue'

export default {
  name: 'EN50160ChartWithHistogram',
  components: {
    EN50160Histogram
  },
  props: {
    // 기본 정보
    title: { type: String, default: '' },
    
    // 라인 차트 데이터
    timeseries: {
      type: Object,
      default: () => ({ labels: [], datasets: [] })
    },
    yAxis: {
      type: Object,
      default: () => ({ name: '', min: null, max: null })
    },
    limitLines: {
      type: Array,
      default: () => []
    },
    markAreas: {
      type: Array,
      default: () => []
    },
    legendItems: {
      type: Array,
      default: () => []
    },
    chartHeight: {
      type: String,
      default: '280px'
    },
    
    // 히스토그램 관련
    showHistogram: {
      type: Boolean,
      default: false
    },
    histogramData: {
      type: Object,
      default: null
    },
    showPhaseHistograms: {
      type: Boolean,
      default: false
    },
    phaseHistograms: {
      type: Object,
      default: null
    },
    histogramUnit: {
      type: String,
      default: ''
    },
    histogramColor: {
      type: String,
      default: '#5470c6'
    },
    histogramLimitMin: {
      type: Number,
      default: null
    },
    histogramLimitMax: {
      type: Number,
      default: null
    },
    histogramLimitValue: {
      type: Number,
      default: null
    },
    histogramHeight: {
      type: String,
      default: '200px'
    },
    
    // 테이블 데이터
    tableColumns: {
      type: Array,
      default: () => []
    },
    tableData: {
      type: Array,
      default: () => []
    },
    
    // 통계 (결과 판정용)
    statistics: {
      type: Object,
      default: null
    }
  },
  setup(props) {
    const lineChartRef = ref(null)
    let chartInstance = null
    let darkModeObserver = null

    // 3상 색상
    const phaseColors = {
      'L1': '#e53935',
      'L2': '#43a047',
      'L3': '#1e88e5'
    }

    // 다크모드 감지
    const isDarkMode = () => {
      return document.documentElement.classList.contains('dark')
    }

    // 테마별 색상 정의
    const getThemeColors = () => {
      const dark = isDarkMode()
      return {
        textColor: dark ? '#e5e7eb' : '#374151',
        textColorSecondary: dark ? '#f3f4f6' : '#6b7280',
        axisLineColor: dark ? '#4b5563' : '#e5e7eb',
        splitLineColor: dark ? '#374151' : '#f3f4f6',
        tooltipBg: dark ? '#1f2937' : '#fff',
        tooltipBorder: dark ? '#374151' : '#e5e7eb'
      }
    }

    // 전체 결과 계산
    const overallResult = computed(() => {
      if (!props.statistics) return null
      
      // 단일 result
      if (props.statistics.result) return props.statistics.result
      
      // 99.5% / 100% 결과
      if (props.statistics.result_99_5 === 'FAIL' || props.statistics.result_100 === 'FAIL') {
        return 'FAIL'
      }
      if (props.statistics.result_95 === 'FAIL') {
        return 'FAIL'
      }
      
      return 'PASS'
    })

    const initLineChart = () => {
      if (!lineChartRef.value) return

      const { labels, datasets } = props.timeseries || {}
      
      if (!labels || labels.length === 0) {
        if (chartInstance) {
          chartInstance.dispose()
          chartInstance = null
        }
        return
      }

      if (chartInstance) {
        chartInstance.dispose()
      }

      chartInstance = echarts.init(lineChartRef.value)
      
      const colors = getThemeColors()

      // 시리즈 생성
      const series = (datasets || []).map(ds => ({
        name: ds.name,
        type: 'line',
        data: ds.data,
        symbol: 'none',
        lineStyle: {
          color: ds.color,
          width: 1.5
        },
        itemStyle: {
          color: ds.color
        },
        areaStyle: ds.areaStyle ? {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: ds.color + '40' },
            { offset: 1, color: ds.color + '05' }
          ])
        } : undefined
      }))

      // 한계선 추가
      const markLineData = (props.limitLines || []).map(line => ({
        yAxis: line.value,
        lineStyle: {
          color: line.color || '#f44336',
          type: 'dashed',
          width: 1.5
        },
        label: {
          show: false
        }
      }))

      if (markLineData.length > 0 && series.length > 0) {
        series[0].markLine = {
          silent: true,
          symbol: 'none',
          data: markLineData
        }
      }

      // 영역 표시
      if (props.markAreas && props.markAreas.length > 0 && series.length > 0) {
        series[0].markArea = {
          silent: true,
          data: props.markAreas.map(area => [
            {
              yAxis: area.from,
              itemStyle: { color: area.color }
            },
            {
              yAxis: area.to
            }
          ])
        }
      }

      const option = {
        grid: {
          left: '8%',
          right: '3%',
          bottom: '12%',
          top: '8%',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            lineStyle: {
              color: colors.axisLineColor
            }
          },
          backgroundColor: colors.tooltipBg,
          borderColor: colors.tooltipBorder,
          textStyle: {
            color: colors.textColor
          },
          formatter: (params) => {
            let html = `<div style="font-size: 12px;"><b>${params[0].axisValue}</b><br/>`
            params.forEach(p => {
              html += `<span style="color:${p.color}">●</span> ${p.seriesName}: <b>${p.value}</b><br/>`
            })
            html += '</div>'
            return html
          }
        },
        xAxis: {
          type: 'category',
          data: labels,
          axisLabel: {
            fontSize: 10,
            color: colors.textColorSecondary,
            rotate: 0,
            interval: Math.floor(labels.length / 7)
          },
          axisLine: {
            lineStyle: { color: colors.axisLineColor }
          },
          axisTick: {
            lineStyle: { color: colors.axisLineColor }
          }
        },
        yAxis: {
          type: 'value',
          name: props.yAxis?.name || '',
          nameTextStyle: {
            fontSize: 11,
            color: colors.textColorSecondary
          },
          min: props.yAxis?.min,
          max: props.yAxis?.max,
          axisLabel: {
            fontSize: 10,
            color: colors.textColorSecondary,
            formatter: props.yAxis?.formatter
          },
          axisLine: {
            show: true,
            lineStyle: { color: colors.axisLineColor }
          },
          splitLine: {
            lineStyle: {
              color: colors.splitLineColor,
              type: 'dashed'
            }
          }
        },
        series: series
      }

      chartInstance.setOption(option)
    }

    // 다크모드 변경 감지
    const setupDarkModeObserver = () => {
      darkModeObserver = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          if (mutation.attributeName === 'class') {
            initLineChart()
          }
        })
      })
      
      darkModeObserver.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['class']
      })
    }

    const handleResize = () => {
      chartInstance?.resize()
    }

    onMounted(() => {
      initLineChart()
      setupDarkModeObserver()
      window.addEventListener('resize', handleResize)
    })

    watch(
      () => props.timeseries,
      () => initLineChart(),
      { deep: true }
    )

    watch(
      () => [props.limitLines, props.markAreas],
      () => initLineChart(),
      { deep: true }
    )

    onBeforeUnmount(() => {
      window.removeEventListener('resize', handleResize)
      if (darkModeObserver) {
        darkModeObserver.disconnect()
      }
      if (chartInstance) {
        chartInstance.dispose()
        chartInstance = null
      }
    })

    return {
      lineChartRef,
      phaseColors,
      overallResult
    }
  }
}
</script>

<style scoped>
.en50160-section {
  transition: box-shadow 0.2s ease;
}

.en50160-section:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style>