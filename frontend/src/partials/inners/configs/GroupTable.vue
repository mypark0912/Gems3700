<template>
  <div class="mt-1 mb-2">
    <table class="table-auto w-full border-collapse text-sm">
      <thead class="bg-gray-100 dark:bg-gray-800 text-left">
        <tr>
          <th class="border border-gray-300 dark:border-gray-600 px-4 py-2 w-40 text-gray-900 dark:text-gray-100">Title</th>
          <th class="border border-gray-300 dark:border-gray-600 px-4 py-2 w-40 text-gray-900 dark:text-gray-100">SubTitle</th>
          <th class="border border-gray-300 dark:border-gray-600 px-4 py-2 w-40 text-center text-gray-900 dark:text-gray-100">Value</th>
          <th class="border border-gray-300 dark:border-gray-600 px-4 py-2 w-40 text-center text-gray-900 dark:text-gray-100">Error</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="(group, titleKey) in channels" :key="titleKey">
          <template v-if="group.view && Array.isArray(group.view)">
            <tr v-for="(row, rowIndex) in group.view" :key="row.subTitle + rowIndex">
              <!-- Title (ex: Phase Voltage) -->
              <td
                v-if="rowIndex === 0"
                :rowspan="group.view.length"
                class="border border-gray-300 dark:border-gray-600 px-4 py-2 align-top bg-gray-50 dark:bg-gray-700/50"
              >
                <div class="font-semibold mb-2 text-gray-900 dark:text-gray-100">{{ titleKey }}</div>
                <FormBlock v-if="group.forms" :forms="group.forms" class="mt-2" />
              </td>

              <!-- SubTitle (ex: U_A, U_B, etc.) -->
              <td class="border border-gray-300 dark:border-gray-600 px-4 py-2 text-gray-800 dark:text-gray-200 bg-white dark:bg-gray-800">
                {{ row.subTitle }}
              </td>
              <td class="border border-gray-300 dark:border-gray-600 px-4 py-2 text-center text-gray-800 dark:text-gray-200 bg-white dark:bg-gray-800">
                {{ row.value }}
              </td>
              <td class="border border-gray-300 dark:border-gray-600 px-4 py-2 text-center text-gray-800 dark:text-gray-200 bg-white dark:bg-gray-800">
                {{ row.error }}
              </td>
            </tr>
          </template>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import FormBlock from './FormArea.vue'

const props = defineProps({
  channels: Object,
})
</script>

<style scoped>
/* 다크모드 테이블 스타일 */
@media (prefers-color-scheme: dark) {
  table {
    border-color: rgb(75 85 99); /* gray-600 */
  }
  
  /* hover 효과 - Title 셀을 제외한 나머지 셀만 적용 */
  tbody tr:hover td:not([rowspan]) {
    background-color: rgb(55 65 81) !important; /* gray-700 */
    color: rgb(229 231 235) !important; /* gray-200 */
  }
}

/* 라이트모드 테이블 스타일 */
@media (prefers-color-scheme: light) {
  table {
    border-color: #ccc;
  }
  
  /* hover 효과 - Title 셀을 제외한 나머지 셀만 적용 */
  tbody tr:hover td:not([rowspan]) {
    background-color: rgb(249 250 251) !important; /* gray-50 */
    color: rgb(31 41 55) !important; /* gray-800 */
  }
}

/* 테이블 전체 스타일 개선 */
table {
  border-spacing: 0;
  border-radius: 0.375rem;
  overflow: hidden;
}

thead tr:first-child th:first-child {
  border-top-left-radius: 0.375rem;
}

thead tr:first-child th:last-child {
  border-top-right-radius: 0.375rem;
}

tbody tr:last-child td:first-child {
  border-bottom-left-radius: 0.375rem;
}

tbody tr:last-child td:last-child {
  border-bottom-right-radius: 0.375rem;
}
</style>