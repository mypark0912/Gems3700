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
            <div class="summary-label">{{ t('dashboard.kwh_realtime.today') }}</div>
            <div class="summary-value-container">
              <div v-if="energyData && energyData.length > 0" class="summary-value"> 61 <span class="summary-unit">kwh</span></div>
              <div v-else class="summary-value"> {{ parseFloat(consumptionData.daily).toFixed(2) }} <span class="summary-unit">kwh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{ t('dashboard.kwh_realtime.comparetoyesterday') }}</span>
              <span class="change-value positive"> {{ parseInt(consumptionData.daily_comparison)}}%</span>
            </div>
          </div>
        </div>

        <!-- 금주 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{ t('dashboard.kwh_realtime.thisweek') }}</div>
            <div class="summary-value-container">
              <div class="summary-value">{{ parseFloat(consumptionData.weekly).toFixed(2) }} <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{ t('dashboard.kwh_realtime.comparetolastweek') }}</span>
              <span class="change-value positive">{{ parseInt(consumptionData.weekly_comparison) }}%</span>
            </div>
          </div>
        </div>

        <!-- 금월 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{ t('dashboard.kwh_realtime.thismonth') }}</div>
            <div class="summary-value-container">
              <div v-if="energyData && energyData.length > 0" class="summary-value">{{energyData[0].data[1].value.toFixed(2) }} <span class="summary-unit">{{ energyData[0].data[1].unit }}</span></div>
              <div v-else class="summary-value">{{ parseFloat(consumptionData.monthly).toFixed(2) }} <span class="summary-unit">kwh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{ t('dashboard.kwh_realtime.comparetolastmonth') }}</span>
              <span class="change-value negative">{{ parseInt(consumptionData.monthly_comparison)}} %</span>
            </div>
          </div>
        </div>

        <!-- 연간 -->
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">{{ t('dashboard.kwh_realtime.thisyear') }}</div>
            <div class="summary-value-container">
              <div class="summary-value">{{ parseFloat(consumptionData.yearly).toFixed(2) }} <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">{{ t('dashboard.kwh_realtime.comparetolastyyear') }}</span>
              <span class="change-value negative">{{ parseInt(consumptionData.yearly_comparison) }}%</span>
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
import { chartAreaGradient } from '../../../charts/ChartjsConfig'
import RealtimeChart from '../../../charts/connect/RealtimeChart2.vue'
import axios from 'axios'
import { tailwindConfig, hexToRGB } from '../../../utils/Utils'
import { useI18n } from 'vue-i18n'

export default {
  name: 'DashboardCard05',
  components: {
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
  const updateInterval = ref(null);
  const channel = computed(() => (props.channel == 'main' || props.channel == 'Main') ? 'Main' : 'Sub');
  
  const energyData = ref([]);
  const consumptionData = ref({
    today: 0,
    daily:0,
    weekly: 0,
    monthly: 0,
    yearly: 0,
    daily_comparison: 0,
    weekly_comparison: 0,
    monthly_comparison: 0,
    yearly_comparison: 0,
  });
  
  // InfluxDB에서 가져온 오늘의 시간별 데이터
  const todayHourlyData = ref([]);
  
  // InfluxDB에서 시간별 데이터 가져오기
  const fetchTodayHourlyData = async (ch) => {
    try {
      //console.log("Fetching today's hourly data for channel:", ch);
      
      const response = await axios.get(`/api/getHourlyEnergyData/${ch}`, {
        params: {
          date: new Date().toISOString().split('T')[0],
        }
      });

      //console.log("API Response:", response.data);

      if (response.data.success && response.data.hourlyData) {
        todayHourlyData.value = response.data.hourlyData;
        //console.log("Received hourly data:", todayHourlyData.value);
      }
    } catch (error) {
      console.error("시간별 데이터 가져오기 실패:", error);
    }
  };

  // Redis 데이터 가져오기
  const fetchRedisOnehData = async (ch) => {
    try {
      const response = await axios.get(`/api/getEnergyRedis/${ch}`);
      
      if (response.data.success) {
        consumptionData.value = {
          today: response.data.data.today || 0,
          daily: response.data.data.daily || 0,
          weekly: response.data.data.weekly || 0,
          monthly: response.data.data.monthly || 0,
          yearly: response.data.data.yearly || 0,
          daily_comparison: response.data.data.daily_comparison || 0,
          weekly_comparison: response.data.data.weekly_comparison || 0,
          monthly_comparison: response.data.data.monthly_comparison || 0,
          yearly_comparison: response.data.data.yearly_comparison || 0,
        };
      }
    } catch (error) {
      console.error("데이터 가져오기 실패:", error);
    }
  };

  // Chart.js용 데이터 포맷 (실제 시간별 데이터 + 보간)
  const chartData = computed(() => {
    const now = new Date();
    const currentHour = now.getHours();
    const currentMinute = now.getMinutes();
    
    const labels = [];
    const data = [];
    
    // 데이터가 없으면 빈 차트
    if (!todayHourlyData.value || todayHourlyData.value.length === 0) {
      return {
        labels: [],
        datasets: [{
          data: [],
          fill: true,
          backgroundColor: function(context) {
            const chart = context.chart;
            const {ctx, chartArea} = chart;
            if (!chartArea) return null;
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
          tension: 0.4,
        }]
      };
    }
    
    // 시간별 데이터를 Map으로 변환
    const hourDataMap = new Map();
    todayHourlyData.value.forEach(item => {
      hourDataMap.set(item.hour, item.value);
    });
    
    // 10분 간격으로 보간된 데이터 생성
    const interval = 10; // 분 단위
    
    for (let hour = 0; hour <= currentHour; hour++) {
      const hourValue = hourDataMap.get(hour) || 0;
      const nextHourValue = hourDataMap.get(hour + 1) || hourValue;
      
      const lastMinute = (hour === currentHour) ? currentMinute : 59;
      
      for (let minute = 0; minute <= lastMinute; minute += interval) {
        // 오늘 날짜의 특정 시간으로 Date 객체 생성
        const today = new Date();
        today.setHours(hour);
        today.setMinutes(minute);
        today.setSeconds(0);
        today.setMilliseconds(0);
        
        // 선형 보간
        const progress = minute / 60;
        const interpolatedValue = hourValue + (nextHourValue - hourValue) * progress;
        
        labels.push(today);
        data.push(Math.max(0, interpolatedValue));
      }
    }
    
    return {
      labels: labels,
      datasets: [
        {
          data: data,
          fill: true,
          backgroundColor: function(context) {
            const chart = context.chart;
            const {ctx, chartArea} = chart;
            if (!chartArea) return null;
            
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
          tension: 0.4,
        },
      ],
    };
  });

  // 초기화 함수
  const initializeData = async () => {
    // 기존 인터벌 정리
    if (updateInterval.value) {
      clearInterval(updateInterval.value);
      updateInterval.value = null;
    }
    
    // 데이터 가져오기
    await fetchRedisOnehData(channel.value);
    await fetchTodayHourlyData(channel.value);
    
    // 1분마다 차트 업데이트 (보간 데이터 재생성)
    updateInterval.value = setInterval(() => {
      // chartData는 computed이므로 자동 업데이트됨
      //console.log("Chart auto-updated");
    }, 60000);
    
    // 매시간 정각에 새 데이터 가져오기
    scheduleHourlyUpdate();
  };

  // 정각 업데이트 스케줄링
  const scheduleHourlyUpdate = () => {
    const now = new Date();
    const msUntilNextHour = (60 - now.getMinutes()) * 60000 - now.getSeconds() * 1000;
    
    setTimeout(() => {
      fetchTodayHourlyData(channel.value);
      fetchRedisOnehData(channel.value);
      
      // 이후 매시간 반복
      setInterval(() => {
        fetchTodayHourlyData(channel.value);
        fetchRedisOnehData(channel.value);
      }, 3600000);
    }, msUntilNextHour);
  };

  // 컴포넌트 마운트 시 실행
  onMounted(async () => {
    await initializeData();
  });

  // props.channel 변경 감지
  watch(
    () => props.channel,
    async (newChannel, oldChannel) => {
      if (newChannel !== oldChannel) {
        await initializeData();
      }
    }
  );

  // 정리
  onUnmounted(() => {
    if (updateInterval.value) {
      clearInterval(updateInterval.value);
      updateInterval.value = null;
    }
  });

  return {
    chartData,  // computed 함수 직접 반환
    t,
    channel,
    energyData,
    consumptionData,
    todayHourlyData,
  };
} 
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-7;
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