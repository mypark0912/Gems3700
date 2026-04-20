<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">{{ t('dashboard.kwh') }}</h2>
      </header>
    </div>

    <!-- Summary Section -->
    <div class="summary-section">
      <div class="summary-container">
        <!-- 금월 전력량 -->
        <div class="summary-item">
          <div class="icon-wrapper current-month">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
          </div>
          <div class="summary-info">
            <div class="summary-label">금월 전력량</div>
            <div class="summary-value">32,000 <span class="unit">kWh</span></div>
          </div>
        </div>

        <!-- 전월 전력량 -->
        <div class="summary-item">
          <div class="icon-wrapper previous-month">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
          </div>
          <div class="summary-info">
            <div class="summary-label">전월 전력량</div>
            <div class="summary-value">31,002 <span class="unit">kWh</span></div>
          </div>
        </div>

        <!-- 누적 전력량 -->
        <div class="summary-item">
          <div class="icon-wrapper total">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
            </svg>
          </div>
          <div class="summary-info">
            <div class="summary-label">누적 전력량</div>
            <div class="summary-value">12,300 <span class="unit">kWh</span></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'DashboardCard05',
  components: {},
  props: {
    channel: {
      type: String,
      default: ''
    },
  },
  setup(props) {
    const { t } = useI18n();
    const channel = computed(() => props.channel == 'main' ? 'Main' : 'Sub');

    return {
      t,
      channel,
    }  
  }  
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-2 xl:col-span-2;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

/* Summary Section - 여백 최소화 */
.summary-section {
  @apply p-4;
}

.summary-container {
  @apply space-y-4;
}

.summary-item {
  @apply flex items-center gap-3;
  @apply p-3 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply border border-gray-100 dark:border-gray-600;
  @apply transition-all duration-200;
  @apply hover:shadow-md hover:border-gray-200 dark:hover:border-gray-500;
}

.icon-wrapper {
  @apply w-11 h-11 rounded-lg flex items-center justify-center flex-shrink-0;
  @apply shadow-sm;
}

.icon-wrapper.current-month {
  @apply bg-blue-500 dark:bg-blue-600;
}

.icon-wrapper.previous-month {
  @apply bg-purple-500 dark:bg-purple-600;
}

.icon-wrapper.total {
  @apply bg-green-500 dark:bg-green-600;
}

.icon {
  @apply w-6 h-6 text-white;
}

.summary-info {
  @apply flex-1 min-w-0;
}

.summary-label {
  @apply text-xs font-semibold text-gray-500 dark:text-gray-400 mb-1;
}

.summary-value {
  @apply text-2xl font-bold text-gray-900 dark:text-white;
  @apply leading-none;
}

.unit {
  @apply text-sm font-medium text-gray-400 dark:text-gray-500 ml-1;
}

/* 반응형 */
@media (max-width: 768px) {
  .summary-section {
    @apply p-3;
  }
  
  .summary-container {
    @apply space-y-3;
  }
  
  .summary-item {
    @apply p-2.5 gap-2.5;
  }
  
  .icon-wrapper {
    @apply w-10 h-10;
  }

  .icon {
    @apply w-5 h-5;
  }
  
  .summary-value {
    @apply text-xl;
  }
  
  .unit {
    @apply text-xs;
  }
}
</style>