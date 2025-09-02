<template>
  <tbody v-if="mode=='Event'" class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
    <!-- Data rows -->
    <tr v-for="(item, index) in eventData" :key="index" class="h-14">
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div>{{ item.Type }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="font-medium text-gray-800 dark:text-gray-100">{{ item.StartTime }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left font-medium text-gray-600">{{ item.Duration }} ms</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left font-medium text-gray-600">{{ item.EndTime }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left">{{ item.Phase }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left">{{ item.Level }}</div>
      </td>
    </tr>
  </tbody>
  
  <tbody v-else class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
    <!-- Data rows -->
    <tr v-for="(item, index) in eventData" :key="index" class="h-14">
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="font-medium text-gray-800 dark:text-gray-100">{{ item.alarm_ts }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left font-medium text-green-600">{{ item.condition_str }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left">{{ item.status }}</div>
      </td>
      <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
        <div class="text-left">{{ item.value.toFixed(2) }}</div>
      </td>
    </tr>
  </tbody>
</template>

<script>
import { ref, inject } from 'vue'
import axios from 'axios'
import { useI18n } from "vue-i18n";

export default {
  name: 'EventTableItem',
  props: ['channel', 'mode'],
  setup(props) {
    const { t } = useI18n();
    const sidebarOpen = ref(false)
    const channel = ref(props.channel);
    const mode = ref(props.mode);
    const eventData = inject('eventData');

    return {
      t,
      sidebarOpen,
      eventData,
      channel,
      mode
    }
  }
}
</script>

<style scoped>
/* 테이블 컨테이너에 최소 높이 설정 제거 */
.table-container {
  /* min-height: 700px; 제거 */
  position: relative;
  display: flex;
  flex-direction: column;
}

/* 테이블 자체 최소 높이 제거 */
.table-container table {
  /* min-height: 350px; 제거 */
  width: 100%;
}

/* 테이블이 늘어나지 않도록 설정 */
.table-container table tbody {
  height: auto;
}

/* 빈 상태 메시지 스타일 - 테이블 아래에 표시 */
.empty-state {
  text-align: center;
  padding: 2rem;
  margin-top: 2rem;
}

/* 다크 모드에서도 잘 보이도록 */
.dark .empty-state {
  color: #9ca3af;
}
</style>