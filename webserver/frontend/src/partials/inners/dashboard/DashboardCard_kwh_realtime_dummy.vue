<template>
  <div class="premium-dashboard-card">
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">전력량 (kWh)</h2>
        <div class="channel-info">
          <!-- <span class="channel-text">메인 계측</span> -->
        </div>
      </header>
    </div>

    <div class="summary-section">
      <div class="summary-grid">
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">금일</div>
            <div class="summary-value-container">
              <div class="summary-value">61.24 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">전일 대비</span>
              <span class="change-value positive">+12%</span>
            </div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">금주</div>
            <div class="summary-value-container">
              <div class="summary-value">387.50 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">전주 대비</span>
              <span class="change-value positive">+5%</span>
            </div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">금월</div>
            <div class="summary-value-container">
              <div class="summary-value">1,542.80 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">전월 대비</span>
              <span class="change-value negative">-3%</span>
            </div>
          </div>
        </div>
        <div class="summary-item">
          <div class="summary-content">
            <div class="summary-label">연간</div>
            <div class="summary-value-container">
              <div class="summary-value">18,204.30 <span class="summary-unit">kWh</span></div>
            </div>
            <div class="summary-change">
              <span class="change-label">전년 대비</span>
              <span class="change-value negative">-8%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="chart-area">
      <RealtimeChart :data="dummyChartData" width="595" height="248" />
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { chartAreaGradient } from '../../../charts/ChartjsConfig'
import RealtimeChart from '../../../charts/connect/RealtimeChart2.vue'
import { tailwindConfig, hexToRGB } from '../../../utils/Utils'

export default {
  name: 'DashboardCard_kwh_realtime_dummy',
  components: {
    RealtimeChart,
  },
  setup() {
    const dummyChartData = computed(() => {
      const labels = [];
      const data = [];
      const now = new Date();
      const currentHour = now.getHours();
      const hourlyValues = [
        1.2, 1.0, 0.8, 0.7, 0.9, 1.5,
        3.2, 5.8, 6.4, 5.9, 5.2, 4.8,
        4.5, 5.1, 5.6, 6.0, 6.3, 5.5,
        4.2, 3.8, 3.1, 2.5, 1.8, 1.3
      ];

      for (let hour = 0; hour <= currentHour; hour++) {
        const hourValue = hourlyValues[hour] || 0;
        const nextHourValue = hourlyValues[hour + 1] || hourValue;
        const lastMinute = (hour === currentHour) ? now.getMinutes() : 59;

        for (let minute = 0; minute <= lastMinute; minute += 10) {
          const t = new Date();
          t.setHours(hour, minute, 0, 0);
          const progress = minute / 60;
          const val = hourValue + (nextHourValue - hourValue) * progress;
          labels.push(t);
          data.push(Math.max(0, val + (Math.random() - 0.5) * 0.3));
        }
      }

      return {
        labels,
        datasets: [{
          data,
          fill: true,
          backgroundColor: function(context) {
            const chart = context.chart;
            const { ctx, chartArea } = chart;
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
    });

    return { dummyChartData };
  }
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm transition-all duration-300 hover:shadow-xl;
  height: 100%;
}

.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content { @apply flex justify-between items-center; }

.card-title {
  @apply text-lg font-bold;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 bg-clip-text text-transparent;
}

.channel-info { @apply flex items-center; }
.channel-text { @apply text-xs font-semibold text-gray-400 dark:text-gray-300 uppercase; }

.summary-section { @apply px-5 py-4 bg-white dark:bg-gray-800; }
.summary-grid { @apply grid grid-cols-4 gap-4; }

.summary-item {
  @apply p-3 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-md hover:bg-gray-100 dark:hover:bg-gray-700;
}

.summary-content { @apply w-full text-left; }
.summary-label { @apply text-sm font-medium text-gray-500 dark:text-gray-300 mb-1; }
.summary-value-container { @apply flex items-start mb-1; }
.summary-value { @apply text-2xl font-bold text-gray-800 dark:text-white; }
.summary-unit { @apply text-lg font-semibold text-gray-600 dark:text-gray-300 ml-1; }
.summary-change { @apply flex items-start gap-1; }
.change-label { @apply text-xs text-gray-500 dark:text-gray-300; }
.change-value { @apply text-xs font-medium; }
.change-value.positive { @apply text-green-600 dark:text-green-400; }
.change-value.negative { @apply text-red-500 dark:text-red-400; }

.chart-area { @apply px-2 pb-3; }

@media (max-width: 768px) {
  .summary-grid { @apply grid-cols-2 gap-3; }
  .summary-value { @apply text-xl; }
  .summary-unit { @apply text-base; }
}
</style>
