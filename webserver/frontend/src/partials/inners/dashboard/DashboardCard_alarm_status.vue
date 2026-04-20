<template>
  <div class="premium-dashboard-card alarm-card">
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">알람 / 이벤트 상태</h2>
      </header>
    </div>

    <div class="tab-bar">
      <div class="tab-group">
        <button
          class="tab-btn"
          :class="{ 'tab-active': activeTab === 'alarm' }"
          @click="activeTab = 'alarm'"
        >
          알람 상태
          <span class="tab-count" :class="activeTab === 'alarm' ? 'tab-count-active' : ''">{{ alarmData.length }}</span>
        </button>
        <button
          class="tab-btn"
          :class="{ 'tab-active': activeTab === 'event' }"
          @click="activeTab = 'event'"
        >
          이벤트 상태
          <span class="tab-count" :class="activeTab === 'event' ? 'tab-count-active' : ''">{{ eventData.length }}</span>
        </button>
      </div>
    </div>

    <div v-show="activeTab === 'alarm'" class="table-container">
      <div class="table-wrapper mt-2">
        <table class="alarm-table">
          <thead>
            <tr>
              <th>알람 기준</th>
              <th>상태</th>
              <th class="text-right">알람 카운트</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in alarmData" :key="'alarm-' + index">
              <td class="cell-condition">{{ item.condition_str }}</td>
              <td>
                <span class="status-badge" :class="item.status === 'OCCURRED' ? 'status-occurred' : 'status-cleared'">
                  {{ item.status }}
                </span>
              </td>
              <td class="cell-value">{{ item.count }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-show="activeTab === 'event'" class="table-container">
      <div class="table-wrapper mt-2">
        <table class="alarm-table">
          <thead>
            <tr>
              <th>유형</th>
              <th>발생 시간</th>
              <th>지속 시간</th>
              <th>상</th>
              <th>레벨</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in eventData" :key="'event-' + index">
              <td>
                <span class="event-type-badge" :class="getEventTypeClass(item.Type)">{{ item.Type }}</span>
              </td>
              <td class="cell-time">{{ item.StartTime }}</td>
              <td class="cell-duration">{{ item.Duration }} ms</td>
              <td>{{ item.Phase }}</td>
              <td class="cell-level">{{ item.Level }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'DashboardCard_alarm_status',
  setup() {
    const activeTab = ref('alarm')

    const alarmData = ref([
      { condition_str: 'Temperature OVER 75.0', status: 'OCCURRED', count: 12 },
      { condition_str: 'Frequency UNDER 59.5', status: 'CLEARED', count: 5 },
      { condition_str: 'Phase Voltage L1 UNDER 190.0', status: 'OCCURRED', count: 28 },
      { condition_str: 'Active Power Total OVER 50000', status: 'OCCURRED', count: 124 },
      { condition_str: 'Power Factor Total UNDER 0.85', status: 'CLEARED', count: 56 },
    ])

    const eventData = ref([
      { Type: 'SAG', StartTime: '2025-09-09 12:53:36.550', Duration: 43, Phase: 'L3', Level: 'L1:65.36, L2:60.37, L3:55.94' },
      { Type: 'SAG', StartTime: '2025-09-09 12:52:27.575', Duration: 7345, Phase: 'L1 L2 L3', Level: 'L1:47.97, L2:48.30, L3:47.88' },
      { Type: 'LONG INTERRUPT', StartTime: '2025-09-08 21:30:25.129', Duration: 9005, Phase: 'L1 L2 L3', Level: 'L1:5.03, L2:4.75, L3:5.00' },
    ])

    const getEventTypeClass = (type) => {
      const map = {
        'SAG': 'type-sag',
        'SWELL': 'type-swell',
        'LONG INTERRUPT': 'type-long-interrupt',
        'SHORT INTERRUPT': 'type-short-interrupt',
      }
      return map[type] || ''
    }

    return { activeTab, alarmData, eventData, getEventTypeClass }
  }
}
</script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm transition-all duration-300 hover:shadow-xl;
  height:  100%;
  
}

.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50 rounded-t-xl;
  @apply bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20;
}

.header-content { @apply flex justify-between items-center; }

.card-title {
  @apply text-lg font-bold;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 bg-clip-text text-transparent;
}

.tab-bar {
  @apply px-5 pt-2 pb-0 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-white dark:bg-gray-800;
}

.tab-group { @apply flex gap-1; }
.tab-btn {
  @apply px-3 py-1.5 text-sm font-medium rounded-t-lg transition-all duration-200;
  @apply text-gray-500 dark:text-gray-400 flex items-center gap-1.5;
  @apply border border-transparent border-b-0;
}
.tab-active {
  @apply text-gray-900 dark:text-white bg-white dark:bg-gray-700;
  @apply border-gray-200 dark:border-gray-600;
  @apply shadow-sm;
}
.tab-count {
  @apply text-xs px-1.5 py-0.5 rounded-full bg-gray-200 dark:bg-gray-600 text-gray-500 dark:text-gray-400;
}
.tab-count-active {
  @apply bg-blue-100 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400;
}

.table-container { @apply flex-1 overflow-auto; min-height: 0; }

/* kwh 카드처럼 좌우 여백 */
.table-wrapper { @apply px-5 py-2; }

.alarm-table { @apply w-full text-sm; }
.alarm-table thead { @apply sticky top-0 z-10; }
.alarm-table thead tr {
  @apply bg-gray-50 dark:bg-gray-700/80 border-b border-gray-200 dark:border-gray-600;
}
.alarm-table th {
  @apply py-2.5 text-left text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider whitespace-nowrap;
}
.alarm-table th:last-child { @apply text-right; }
.alarm-table tbody tr {
  @apply border-b border-gray-100 dark:border-gray-700/50 transition-colors duration-150 hover:bg-gray-50 dark:hover:bg-gray-700/30;
}
.alarm-table td {
  @apply py-2.5 text-sm text-gray-700 dark:text-gray-300 whitespace-nowrap;
}
.alarm-table td + td { @apply pl-3; }
.alarm-table th:first-child,
.alarm-table td:first-child {
  @apply pl-4;
}
.alarm-table th:last-child,
.alarm-table td:last-child {
  @apply pr-4;
}
.cell-condition { @apply font-medium text-emerald-600 dark:text-emerald-400; }
.cell-value { @apply font-mono text-right font-bold text-blue-600 dark:text-blue-400; }
.cell-time { @apply font-medium text-gray-800 dark:text-gray-100; }

.status-badge { @apply inline-flex items-center px-2 py-0.5 text-xs font-semibold rounded-full; }
.status-occurred { @apply bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400; }
.status-cleared { @apply bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400; }

.event-type-badge { @apply inline-flex items-center px-2 py-0.5 text-xs font-bold rounded; }
.type-sag { @apply bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-400; }
.type-swell { @apply bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-400; }
.type-long-interrupt { @apply bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-400; }
.type-short-interrupt { @apply bg-orange-100 text-orange-800 dark:bg-orange-900/40 dark:text-orange-400; }
</style>