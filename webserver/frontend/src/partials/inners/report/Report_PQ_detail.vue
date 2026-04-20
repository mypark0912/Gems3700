<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl mt-4">
    
    <!-- 시계열 추이 카드 -->
    <div class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
      <div class="absolute top-0 left-0 right-0 h-0.5 bg-blue-500" aria-hidden="true"></div>
      <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-blue-500 mr-3">
            <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M3 20h18M5 16l4-4 4 4 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            {{ t(`report.cardTitle.parameter`)}}
          </h3>
          <div class="ml-auto flex items-center gap-3">
            <select v-model="selectedTimeRange" @change="updateData" class="time-selector">
              <option value="24h">{{ t(`report.cardContext.pqInfo.hour`)}}</option>
              <option value="7d">{{ t(`report.cardContext.pqInfo.day`)}}</option>
              <option value="30d">{{ t(`report.cardContext.pqInfo.month`)}}</option>
            </select>
            <button class="header-action-btn" @click="refreshData">
              🔄
            </button>
          </div>
        </header>
      </div>
      <div class="px-4 py-4 space-y-3">
        <!-- 시계열 차트 -->
        <div class="timeseries-section">
          <div ref="timeseriesChart" class="timeseries-chart"></div>
          <div class="legend-info">
            <div class="legend-item">
              <span class="legend-color" style="background-color: #ff6b6b;"></span>
              <span>{{ t(`dashboard.pq.unbalance`)}} (%)</span>
            </div>
            <div class="legend-item">
              <span class="legend-color" style="background-color: #4ecdc4;"></span>
              <span>{{ t(`report.cardTitle.THD1`)}} (%)</span>
            </div>
            <div class="legend-item">
              <span class="legend-color" style="background-color: #45b7d1;"></span>
              <span>{{ t(`report.cardTitle.THD2`)}} (%)</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from 'vue-i18n'  

export default {
  name: 'PowerQualityAnalysis',
  setup() {
    // 반응형 데이터
    const { t } = useI18n();
    const selectedTimeRange = ref('24h')
    const timeseriesChart = ref(null)

    let timeseriesChartInstance = null

    // 더미 데이터 생성 함수
    const generateDummyData = (timeRange) => {
      const now = new Date()
      const dataPoints = timeRange === '24h' ? 24 : timeRange === '7d' ? 168 : 720
      const interval = timeRange === '24h' ? 1 : timeRange === '7d' ? 1 : 1 // 시간 단위
      
      const data = []
      for (let i = dataPoints - 1; i >= 0; i--) {
        const time = new Date(now.getTime() - i * interval * 60 * 60 * 1000)
        
        // 시간대별 패턴 시뮬레이션
        const hour = time.getHours()
        let baseVUF = 0.8 + Math.random() * 0.4
        let baseTHDV = 2.0 + Math.random() * 1.0
        let baseTHDI = 3.0 + Math.random() * 1.5
        
        // 특정 시간대에 품질 저하
        if ((hour >= 9 && hour <= 11) || (hour >= 14 && hour <= 16)) {
          baseVUF += Math.random() * 0.8
          baseTHDV += Math.random() * 1.5
          baseTHDI += Math.random() * 2.0
        }
        
        // 가끔 스파이크 발생
        if (Math.random() < 0.05) {
          baseVUF *= 1.5
          baseTHDV *= 1.3
          baseTHDI *= 1.4
        }
        
        data.push({
          time: time.toISOString(),
          vuf: Math.min(baseVUF, 3.0),
          thdv: Math.min(baseTHDV, 8.0),
          thdi: Math.min(baseTHDI, 10.0)
        })
      }
      
      return data
    }

    const dummyData = ref(generateDummyData('24h'))

    // 계산된 속성들
    const averageVUF = computed(() => {
      const sum = dummyData.value.reduce((acc, item) => acc + item.vuf, 0)
      return sum / dummyData.value.length
    })

    const averageTHDV = computed(() => {
      const sum = dummyData.value.reduce((acc, item) => acc + item.thdv, 0)
      return sum / dummyData.value.length
    })

    const averageTHDI = computed(() => {
      const sum = dummyData.value.reduce((acc, item) => acc + item.thdi, 0)
      return sum / dummyData.value.length
    })

    // 시계열 차트 생성
    const createTimeseriesChart = () => {
      const times = dummyData.value.map(item => new Date(item.time))
      const vufData = dummyData.value.map(item => item.vuf.toFixed(2))
      const thdvData = dummyData.value.map(item => item.thdv.toFixed(2))
      const thdiData = dummyData.value.map(item => item.thdi.toFixed(2))
      
      const option = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          },
          formatter: function (params) {
            let result = params[0].axisValueLabel + '<br/>'
            params.forEach(function (item) {
              result += item.marker + item.seriesName + ': ' + item.value + '%<br/>'
            })
            return result
          }
        },
        legend: {
              data: [
                      t('dashboard.pq.unbalance'),
                      t('report.cardTitle.THD1'),
                      t('report.cardTitle.THD2')
                    ]
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            boundaryGap: false,
            data: times.map(time => time.toLocaleTimeString('ko-KR', {
              hour: '2-digit',
              minute: '2-digit'
            }))
          }
        ],
        yAxis: [
          {
            type: 'value',
            name: t('report.cardTitle.distortionRate') + ' (%)',
            min: 0,
            max: 10,
            axisLabel: {
              formatter: '{value}%'
            }
          }
        ],
        dataZoom: [
          {
            type: 'inside',
            start: 70,
            end: 100
          },
          {
            start: 70,
            end: 100
          }
        ],
        series: [
          {
            name: t('dashboard.pq.unbalance'),
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: vufData,
            color: '#ff6b6b',
            markLine: {
              data: [{ yAxis: 2, name: '불평형률 기준선' }]
            }
          },
          {
            name: t(`report.cardTitle.THD1`),
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: thdvData,
            color: '#4ecdc4',
            markLine: {
              data: [{ yAxis: 5, name: 'THDV 기준선' }]
            }
          },
          {
            name: t(`report.cardTitle.THD2`),
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: thdiData,
            color: '#45b7d1',
            markLine: {
              data: [{ yAxis: 5, name: 'THDI 기준선' }]
            }
          }
        ]
      }
      
      timeseriesChartInstance.setOption(option)
    }

    // 데이터 업데이트
    const updateData = () => {
      dummyData.value = generateDummyData(selectedTimeRange.value)
      
      if (timeseriesChartInstance) {
        createTimeseriesChart()
      }
    }

    const refreshData = () => {
      updateData()
    }

    // 라이프사이클 훅
    onMounted(() => {
      timeseriesChartInstance = echarts.init(timeseriesChart.value)
      
      createTimeseriesChart()
      
      window.addEventListener('resize', () => {
        timeseriesChartInstance?.resize()
      })
    })

    onUnmounted(() => {
      timeseriesChartInstance?.dispose()
      window.removeEventListener('resize', () => {})
    })

    return {
      selectedTimeRange,
      timeseriesChart,
      averageVUF,
      averageTHDV,
      averageTHDI,
      updateData,
      refreshData,
      t
    }
  }
}
</script>

<style scoped>
/* 시간 선택기 */
.time-selector {
  @apply px-3 py-1 text-xs bg-gray-100 dark:bg-gray-700 rounded border-none;
  @apply text-gray-700 dark:text-gray-300;
}

/* 헤더 액션 버튼 */
.header-action-btn {
  @apply w-8 h-8 flex items-center justify-center;
  @apply text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200;
  @apply rounded-full hover:bg-gray-100 dark:hover:bg-gray-700;
  @apply transition-all duration-200;
}

/* 요약 통계 그리드 */
.summary-stats-grid {
  @apply grid grid-cols-1 gap-4;
}

.stat-card {
  @apply bg-gray-50 dark:bg-gray-700/50 p-4 rounded-lg text-center;
}

.stat-card h4 {
  @apply text-gray-600 dark:text-gray-400 text-sm font-medium mb-2;
}

.stat-value {
  @apply text-gray-900 dark:text-gray-100 text-xl font-bold;
}

/* 시계열 섹션 */
.timeseries-chart {
  @apply h-96 w-full;
}

.legend-info {
  @apply flex justify-center gap-6 mt-4 pt-4 border-t border-gray-200 dark:border-gray-700;
}

.legend-item {
  @apply flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400;
}

.legend-color {
  @apply w-3 h-3 rounded-sm;
}

/* 반응형 */
@media (max-width: 1024px) {
  .timeseries-chart {
    @apply h-80;
  }
}

@media (max-width: 640px) {
  .legend-info {
    @apply flex-col gap-2;
  }
  
  .summary-stats-grid {
    @apply grid-cols-1;
  }
}
</style>