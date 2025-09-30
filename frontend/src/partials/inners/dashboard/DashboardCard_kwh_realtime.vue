<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">{{ t('dashboard.kwh') }}</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'Main' ? t('dashboard.meter.subtitle_main') : t('dashboard.meter.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- Summary Section -->
    <div class="summary-section">
      <div class="summary-grid">
        <!-- 금일 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{t('dashboard.kwh_realtime.today')}}</div>
            <div class="summary-value-container">
              <div class="summary-value">24.7 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{ t('dashboard.kwh_realtime.comparetoyesterday') }}</span>
              <span class="change-value positive">+49%</span>
            </div>
          </div>
        </div>

        <!-- 금주 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{t('dashboard.kwh_realtime.thisweek')}}</div>
            <div class="summary-value-container">
              <div class="summary-value">56.9 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{t('dashboard.kwh_realtime.comparetolastweek')}}</span>
              <span class="change-value positive">+47%</span>
            </div>
          </div>
        </div>

        <!-- 금월 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{t('dashboard.kwh_realtime.thismonth')}}</div>
            <div class="summary-value-container">
              <div class="summary-value">80.9 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{t('dashboard.kwh_realtime.comparetolastmonth')}}</span>
              <span class="change-value negative">-7%</span>
            </div>
          </div>
        </div>

        <!-- 연간 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{t('dashboard.kwh_realtime.thisyear')}}</div>
            <div class="summary-value-container">
              <div class="summary-value">11,340 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{t('dashboard.kwh_realtime.comparetolastyyear')}}</span>
              <span class="change-value negative">-17%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <RealtimeChart :data="chartData" width="595" height="248" />
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
//import Tooltip from '../../components/Tooltip.vue'
import { chartAreaGradient } from '../../../charts/ChartjsConfig'
import RealtimeChart from '../../../charts/connect/RealtimeChart2.vue'

// Import utilities
import { tailwindConfig, hexToRGB } from '../../../utils/Utils'
import { useI18n } from 'vue-i18n'  // ✅ 추가

export default {
  name: 'DashboardCard05',
  components: {
    //Tooltip,
    RealtimeChart,
  },
  props: {
    channel: {
      type: String,
      default: ''
    },
  },
  setup(props) {
    const { t } = useI18n();
    let updateInterval = null;
    const channel = computed(() => props.channel == 'main' ? 'Main' : 'Sub');
    
    // IMPORTANT:
    // Code below is for demo purpose only, and it's not covered by support.
    // If you need to replace dummy data with real data,
    // refer to Chart.js documentation: https://www.chartjs.org/docs/latest    

    const counter = ref(0)
    const range = ref(35)

    // Dummy data to be looped
    const sampleData = [
      57.81, 57.75, 55.48, 54.28, 53.14, 52.25, 51.04, 52.49, 55.49, 56.87,
      53.73, 56.42, 58.06, 55.62, 58.16, 55.22, 58.67, 60.18, 61.31, 63.25,
      65.91, 64.44, 65.97, 62.27, 60.96, 59.34, 55.07, 59.85, 53.79, 51.92,
      50.95, 49.65, 48.09, 49.81, 47.85, 49.52, 50.21, 52.22, 54.42, 53.42,
      50.91, 58.52, 53.37, 57.58, 59.09, 59.36, 58.71, 59.42, 55.93, 57.71,
      50.62, 56.28, 57.37, 53.08, 55.94, 55.82, 53.94, 52.65, 50.25,
    ]
    
    const slicedData = ref(sampleData.slice(0, range.value))

    // Generate fake dates from now to back in time
    const generateDates = () => {
      const now = new Date()
      const dates = []
      sampleData.forEach((v, i) => {
        dates.push(new Date(now - 2000 - i * 2000))
      })
      return dates
    }
    
    const slicedLabels = ref(generateDates().slice(0, range.value).reverse())

    // Fake update every 2 seconds
    const interval = ref(null)
    onMounted(() => {
      interval.value = setInterval(() => {
        counter.value++
      }, 2000)
    })
    onUnmounted(() => {
      clearInterval(interval)
    })

    // Loop through data array and update
    watch(counter, () => {
      range.value++;
      if (range.value >= sampleData.length) {
        range.value = 0;
      }
      slicedData.value.shift();
      slicedData.value.push(sampleData[range.value]);      
      slicedLabels.value.shift()
      slicedLabels.value.push(new Date())
    })

    const chartData = computed(() => {
      return {
        labels: slicedLabels.value,
        datasets: [
          {
            data: [...slicedData.value],
            fill: true,
            backgroundColor: function(context) {
              const chart = context.chart;
              const {ctx, chartArea} = chart;
              return chartAreaGradient(ctx, chartArea, [
                { stop: 0, color: `rgba(${hexToRGB(tailwindConfig().theme.colors.violet[500])}, 0)` },
                { stop: 1, color: `rgba(${hexToRGB(tailwindConfig().theme.colors.violet[500])}, 0.2)` }
              ]);
            },
            borderColor: tailwindConfig().theme.colors.violet[500],
            borderWidth: 2,
            pointRadius: 0,
            pointHoverRadius: 3,
            pointBackgroundColor: tailwindConfig().theme.colors.violet[500],
            clip: 20,
            tension: 0.2,
          },
        ],
      }
    })

    return {
      counter,
      range,
      slicedData,
      slicedLabels,
      interval,
      chartData,
      t,
      channel,
    }  
  }  
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-6;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-50/50 to-purple-50/50 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-white;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 bg-clip-text text-transparent;
}

.channel-info {
  @apply flex items-center;
}

.channel-text {
  @apply text-xs font-semibold text-gray-400 dark:text-gray-300 uppercase;
}

.summary-section {
  @apply px-5 py-4;
  @apply bg-white dark:bg-gray-800;
}

.summary-grid {
  @apply grid grid-cols-4 gap-4;
}

.summary-item {
  @apply flex-1;
  @apply p-3 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-md hover:bg-gray-100 dark:hover:bg-gray-700;
}

.summary-content {
  @apply w-full text-left;
}

.summary-label {
  @apply text-sm font-medium text-gray-500 dark:text-gray-300 mb-1;
}

.summary-value-container {
  @apply flex items-start mb-1;
}

.summary-value {
  @apply text-2xl font-bold text-gray-800 dark:text-white;
}

.summary-unit {
  @apply text-lg font-semibold text-gray-600 dark:text-gray-300 ml-1;
}

.summary-change {
  @apply flex items-start gap-1;
}

.change-label {
  @apply text-xs text-gray-500 dark:text-gray-300;
}

.change-value {
  @apply text-xs font-medium;
}

.change-value.positive {
  @apply text-green-600 dark:text-green-400;
}

.change-value.negative {
  @apply text-red-500 dark:text-red-400;
}

/* 반응형 개선 */
@media (max-width: 768px) {
  .summary-grid {
    @apply grid-cols-2 gap-3;
  }
  
  .summary-value {
    @apply text-xl;
  }
  
  .summary-unit {
    @apply text-base;
  }
}

@media (max-width: 1024px) and (min-width: 769px) {
  .summary-grid {
    @apply grid-cols-4 gap-3;
  }
  
  .summary-item {
    @apply p-3;
  }
  
  .summary-value {
    @apply text-xl;
  }
}
</style>