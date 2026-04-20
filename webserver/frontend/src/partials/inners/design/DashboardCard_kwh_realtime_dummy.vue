<template>
  <div class="card-wrap">
    <div class="card-header">
      <h3 class="card-title meter-accent-indigo">{{ t('dashboard.kwh') }}</h3>
    </div>

    <div class="card-body">
      <!-- Summary Section -->
      <div class="summary-grid">
        <div class="summary-item today">
          <div class="summary-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
            </svg>
          </div>
          <div class="summary-content">
            <span class="summary-label">{{ t('dashboard.kwh_realtime.today') }}</span>
            <div class="summary-value">61.24 <span class="summary-unit">kWh</span></div>
            <div class="summary-compare">
              <span class="compare-label">{{ t('dashboard.kwh_realtime.comparetoyesterday') }}</span>
              <span class="compare-value up">+12%</span>
            </div>
          </div>
        </div>

        <div class="summary-item week">
          <div class="summary-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <div class="summary-content">
            <span class="summary-label">{{ t('dashboard.kwh_realtime.thisweek') }}</span>
            <div class="summary-value">387.50 <span class="summary-unit">kWh</span></div>
            <div class="summary-compare">
              <span class="compare-label">{{ t('dashboard.kwh_realtime.comparetolastweek') }}</span>
              <span class="compare-value up">+5%</span>
            </div>
          </div>
        </div>

        <div class="summary-item month">
          <div class="summary-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
            </svg>
          </div>
          <div class="summary-content">
            <span class="summary-label">{{ t('dashboard.kwh_realtime.thismonth') }}</span>
            <div class="summary-value">1,542.80 <span class="summary-unit">kWh</span></div>
            <div class="summary-compare">
              <span class="compare-label">{{ t('dashboard.kwh_realtime.comparetolastmonth') }}</span>
              <span class="compare-value down">-3%</span>
            </div>
          </div>
        </div>

        <div class="summary-item year">
          <div class="summary-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
            </svg>
          </div>
          <div class="summary-content">
            <span class="summary-label">{{ t('dashboard.kwh_realtime.thisyear') }}</span>
            <div class="summary-value">18,204.30 <span class="summary-unit">kWh</span></div>
            <div class="summary-compare">
              <span class="compare-label">{{ t('dashboard.kwh_realtime.comparetolastyyear') }}</span>
              <span class="compare-value down">-8%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Chart -->
      <div class="chart-area">
        <RealtimeChart :data="dummyChartData" width="100%" height="200" />
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { chartAreaGradient } from '../../../charts/ChartjsConfig'
import RealtimeChart from '../../../charts/connect/RealtimeChart2.vue'
import { tailwindConfig, hexToRGB } from '../../../utils/Utils'

export default {
  name: 'DesignCard_Energy',
  components: {
    RealtimeChart,
  },
  setup() {
    const { t } = useI18n()
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

    return { t, dummyChartData };
  }
}
</script>

<style scoped>
.card-wrap {
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply overflow-hidden;
}
.card-header {
  @apply flex justify-between items-center px-4 py-2.5;
}
.card-title {
  @apply text-base font-bold text-gray-800 dark:text-white flex items-center gap-2;
}
.card-title::before {
  content: '';
  @apply w-1 h-4 rounded-full inline-block flex-shrink-0;
}
.meter-accent-indigo::before {
  @apply bg-indigo-500;
}
.card-body {
  @apply px-4 py-3;
}

/* Summary Grid */
.summary-grid {
  @apply grid grid-cols-2 lg:grid-cols-4 gap-3 mb-4;
}
.summary-item {
  @apply flex items-start gap-3 p-3 rounded-lg;
  @apply border border-gray-100 dark:border-gray-700;
  @apply transition-all duration-200 hover:shadow-md;
}
.summary-item.today {
  @apply bg-blue-50/60 dark:bg-blue-900/20;
}
.summary-item.today .summary-icon {
  @apply text-blue-500 dark:text-blue-400;
}
.summary-item.week {
  @apply bg-emerald-50/60 dark:bg-emerald-900/20;
}
.summary-item.week .summary-icon {
  @apply text-emerald-500 dark:text-emerald-400;
}
.summary-item.month {
  @apply bg-violet-50/60 dark:bg-violet-900/20;
}
.summary-item.month .summary-icon {
  @apply text-violet-500 dark:text-violet-400;
}
.summary-item.year {
  @apply bg-amber-50/60 dark:bg-amber-900/20;
}
.summary-item.year .summary-icon {
  @apply text-amber-500 dark:text-amber-400;
}
.summary-icon {
  @apply flex-shrink-0 mt-0.5;
}
.summary-content {
  @apply flex-1 min-w-0;
}
.summary-label {
  @apply text-sm text-gray-500 dark:text-gray-400 block;
}
.summary-value {
  @apply text-base font-bold text-gray-800 dark:text-white tabular-nums leading-tight;
}
.summary-unit {
  @apply text-xs font-medium text-gray-600 dark:text-gray-400;
}
.summary-compare {
  @apply flex items-center gap-1 mt-0.5;
}
.compare-label {
  @apply text-xs text-gray-600 dark:text-gray-400;
}
.compare-value {
  @apply text-xs font-semibold;
}
.compare-value.up {
  @apply text-green-600 dark:text-green-400;
}
.compare-value.down {
  @apply text-red-500 dark:text-red-400;
}

/* Chart */
.chart-area {
  @apply bg-gray-50/50 dark:bg-gray-800/50 rounded-lg p-2;
  @apply border border-gray-100 dark:border-gray-700;
}

@media (max-width: 640px) {
  .summary-grid {
    @apply grid-cols-2 gap-2;
  }
  .summary-value {
    @apply text-lg;
  }
}
</style>
