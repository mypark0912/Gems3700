<template>
  <div class="card-wrap">
    <div class="card-header">
      <h3 class="card-title meter-accent-teal">{{ t('dashboard.alarmLog') }}</h3>
    </div>
    <div class="card-body">
      <div class="log-table-wrap">
        <table class="log-table">
          <thead>
            <tr>
              <th>{{ t('dashboard.alarmTable.time') }}</th>
              <th>{{ t('dashboard.alarmTable.criteria') }}</th>
              <th>{{ t('dashboard.alarmTable.status') }}</th>
              <th>{{ t('dashboard.alarmTable.value') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, idx) in logItems" :key="idx">
              <td class="col-time">{{ item.ts_formatted }}</td>
              <td class="col-msg">{{ item.criteria }}</td>
              <td>
                <span class="status-tag" :class="item.status === 'OCCURRED' ? 'tag-active' : 'tag-resolved'">{{ item.status }}</span>
              </td>
              <td class="col-value">{{ item.value.toFixed(2) }}</td>
            </tr>
            <tr v-if="logItems.length === 0">
              <td colspan="4" class="text-center text-gray-400 py-4">No data</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

export default {
  name: 'DesignCard_EventLog',
  setup() {
    const { t } = useI18n()
    const logItems = ref([])
    let timer = null

    const fetchRecentAlarmLog = async () => {
      try {
        const response = await axios.get('/api/getRecentAlarmLog')
        if (response.data.success) {
          logItems.value = response.data.data
        }
      } catch (error) {
        console.error('알람로그 가져오기 실패:', error)
      }
    }

    onMounted(() => {
      fetchRecentAlarmLog()
      timer = setInterval(fetchRecentAlarmLog, 3600000) // 1시간
    })

    onUnmounted(() => {
      if (timer) clearInterval(timer)
    })

    return { t, logItems }
  },
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
.meter-accent-teal::before {
  @apply bg-teal-500;
}
.card-body {
  @apply px-4 py-3;
}

/* Table */
.log-table-wrap {
  @apply overflow-y-auto;
}
.log-table {
  @apply w-full text-left;
  border-collapse: collapse;
}
.log-table thead th {
  @apply text-xs font-semibold text-gray-500 dark:text-gray-400 pb-2;
  @apply border-b border-gray-200 dark:border-gray-700;
  @apply sticky top-0 bg-white dark:bg-gray-800;
}
.log-table tbody tr {
  @apply border-b border-gray-100 dark:border-gray-700/50;
}
.log-table tbody td {
  @apply py-2 text-sm text-gray-600 dark:text-gray-300;
}
.col-time {
  @apply tabular-nums text-gray-600 dark:text-gray-400 whitespace-nowrap pr-2 text-sm;
}
.col-msg {
  @apply truncate max-w-[160px];
}
.col-value {
  @apply tabular-nums text-sm;
}

/* Status tag */
.status-tag {
  @apply text-sm font-medium;
}
.tag-active {
  @apply text-red-500 dark:text-red-400;
}
.tag-resolved {
  @apply text-green-500 dark:text-green-400;
}
</style>
