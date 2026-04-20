<template>
    <div class="en50160-harmonics bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <!-- 헤더 -->
      <div class="flex items-center justify-between mb-4">
        <h4 class="text-base font-semibold text-gray-800 dark:text-gray-100">
          Voltage Harmonics
        </h4>
        <div class="flex items-center gap-2">
          <!-- 페이즈 선택 탭 -->
          <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
            <button
              v-for="phase in phases"
              :key="phase"
              @click="selectedPhase = phase"
              :class="[
                'px-3 py-1 text-sm font-medium rounded-md transition-colors',
                selectedPhase === phase
                  ? 'bg-white dark:bg-gray-600 text-gray-800 dark:text-white shadow-sm'
                  : 'text-gray-600 dark:text-gray-100 hover:text-gray-800 dark:hover:text-white'
              ]"
            >
              {{ phase }}
            </button>
          </div>
        </div>
      </div>
  
      <!-- 전체 결과 요약 -->
      <div v-if="overallResult" class="mb-4 flex items-center gap-4">
        <span class="text-sm text-gray-600 dark:text-gray-100">Overall Result:</span>
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
        <span class="text-sm text-gray-500 dark:text-gray-100">
          ({{ totalSamples }} samples)
        </span>
      </div>
  
      <!-- 바 차트 -->
      <div class="mb-4">
        <div ref="chartRef" style="height: 250px; width: 100%;"></div>
      </div>
  
      <!-- 테이블 -->
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="bg-gray-50 dark:bg-gray-700">
            <tr>
              <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                Harmonic
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                Limit (%)
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                Max (%)
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                Avg (%)
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                95th (%)
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                In Range
              </th>
              <th class="px-3 py-2 text-center text-xs font-medium text-gray-500 dark:text-gray-100 uppercase">
                Result
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
            <tr
              v-for="(item, key) in currentPhaseData"
              :key="key"
              :class="item.result === 'FAIL' ? 'bg-red-50 dark:bg-red-900/20' : ''"
            >
              <td class="px-3 py-2 text-gray-700 dark:text-gray-100 font-medium">
                {{ formatHarmonicName(key) }}
              </td>
              <td class="px-3 py-2 text-center text-gray-600 dark:text-gray-100">
                {{ item.limit }}
              </td>
              <td class="px-3 py-2 text-center text-gray-700 dark:text-gray-100">
                {{ item.max }}
              </td>
              <td class="px-3 py-2 text-center text-gray-700 dark:text-gray-100">
                {{ item.avg }}
              </td>
              <td class="px-3 py-2 text-center text-gray-700 dark:text-gray-100">
                {{ item.percentile_95 }}
              </td>
              <td class="px-3 py-2 text-center text-gray-700 dark:text-gray-100">
                {{ item.in_range_95_percent }}%
              </td>
              <td class="px-3 py-2 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 text-xs font-medium rounded',
                    item.result === 'PASS'
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                      : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                  ]"
                >
                  {{ item.result }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
  import * as echarts from 'echarts'
  
  export default {
    name: 'EN50160Harmonics',
    props: {
      harmonicsData: {
        type: Object,
        default: null
      }
    },
    setup(props) {
      const chartRef = ref(null)
      const selectedPhase = ref('L1')
      let chartInstance = null
  
      const phases = ['L1', 'L2', 'L3']
  
      // 현재 선택된 페이즈 데이터
      const currentPhaseData = computed(() => {
        if (!props.harmonicsData?.phases) return {}
        return props.harmonicsData.phases[selectedPhase.value] || {}
      })
  
      // 전체 결과 (하나라도 FAIL이면 FAIL)
      const overallResult = computed(() => {
        if (!props.harmonicsData?.phases) return null
        
        for (const phase of Object.values(props.harmonicsData.phases)) {
          for (const harmonic of Object.values(phase)) {
            if (harmonic.result === 'FAIL') return 'FAIL'
          }
        }
        return 'PASS'
      })
  
      // 총 샘플 수
      const totalSamples = computed(() => {
        return props.harmonicsData?.total_samples || 0
      })
  
      // 고조파 이름 포맷
      const formatHarmonicName = (key) => {
        // h2 -> H2, h3 -> H3 ...
        return key.toUpperCase()
      }
  
      // 차트 초기화
      const initChart = () => {
        if (!chartRef.value || !currentPhaseData.value) return
  
        if (chartInstance) {
          chartInstance.dispose()
        }
  
        chartInstance = echarts.init(chartRef.value)
  
        const harmonics = Object.entries(currentPhaseData.value)
        const categories = harmonics.map(([key]) => formatHarmonicName(key))
        const maxValues = harmonics.map(([, data]) => data.max)
        const limitValues = harmonics.map(([, data]) => data.limit)
        const percentile95 = harmonics.map(([, data]) => data.percentile_95)
  
        // 색상 결정 (limit 초과시 빨간색)
        const barColors = harmonics.map(([, data]) => 
          data.percentile_95 > data.limit ? '#ef4444' : '#3b82f6'
        )
  
        const option = {
          tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' },
            formatter: (params) => {
              const idx = params[0].dataIndex
              const [key, data] = harmonics[idx]
              return `
                <div style="font-size: 12px;">
                  <div style="font-weight: bold; margin-bottom: 4px;">${formatHarmonicName(key)}</div>
                  <div>Max: <b>${data.max}%</b></div>
                  <div>95th: <b>${data.percentile_95}%</b></div>
                  <div>Limit: <b>${data.limit}%</b></div>
                  <div>Result: <b>${data.result}</b></div>
                </div>
              `
            }
          },
          legend: {
            data: ['95th Percentile', 'Limit'],
            top: 0,
            textStyle: { fontSize: 11, color: '#6b7280' }
          },
          grid: {
            left: '5%',
            right: '5%',
            bottom: '15%',
            top: '15%',
            containLabel: true
          },
          xAxis: {
            type: 'category',
            data: categories,
            axisLabel: {
              fontSize: 10,
              color: '#6b7280',
              rotate: 45
            },
            axisLine: { lineStyle: { color: '#e5e7eb' } }
          },
          yAxis: {
            type: 'value',
            name: '[%]',
            nameTextStyle: { fontSize: 11, color: '#6b7280' },
            axisLabel: { fontSize: 10, color: '#6b7280' },
            axisLine: { show: true, lineStyle: { color: '#e5e7eb' } },
            splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } }
          },
          series: [
            {
              name: '95th Percentile',
              type: 'bar',
              data: percentile95.map((value, idx) => ({
                value,
                itemStyle: { color: barColors[idx] }
              })),
              barWidth: '60%'
            },
            {
              name: 'Limit',
              type: 'line',
              data: limitValues,
              symbol: 'none',
              lineStyle: {
                color: '#f44336',
                type: 'dashed',
                width: 2
              }
            }
          ]
        }
  
        chartInstance.setOption(option)
      }
  
      // 리사이즈 핸들러
      const handleResize = () => {
        chartInstance?.resize()
      }
  
      onMounted(() => {
        initChart()
        window.addEventListener('resize', handleResize)
      })
  
      watch(() => props.harmonicsData, () => {
        initChart()
      }, { deep: true })
  
      watch(selectedPhase, () => {
        initChart()
      })
  
      onBeforeUnmount(() => {
        window.removeEventListener('resize', handleResize)
        if (chartInstance) {
          chartInstance.dispose()
          chartInstance = null
        }
      })
  
      return {
        chartRef,
        selectedPhase,
        phases,
        currentPhaseData,
        overallResult,
        totalSamples,
        formatHarmonicName
      }
    }
  }
  </script>
  
  <style scoped>
  .en50160-harmonics {
    transition: box-shadow 0.2s ease;
  }
  
  .en50160-harmonics:hover {
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  }
  </style>