<template>
  <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl relative alarm-wrapper">
    <div class="alarm-content">
      <!-- Table -->
      <div class="overflow-x-auto table-scroll-container">
        <!-- 테이블 컨테이너에 최소 높이 적용 -->
        <div class="table-container">
          <table class="table-auto w-full dark:text-white divide-y divide-gray-100 dark:divide-gray-700/60">
            <!-- Table header -->
            <thead class="text-xs uppercase text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-900/20 border-t border-gray-100 dark:border-gray-700/60 sticky top-0 z-10">
              <tr>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">{{ t("event.table.index") }}</div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">{{ t("event.table.alarmChannel") }}</div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">{{ t("event.table.state") }}</div>
                </th>
                <th class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="font-semibold text-left">{{ t("event.table.count") }}</div>
                </th>
              </tr>
            </thead>
            <!-- Table body -->
            <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
              <!-- Data rows -->
              <tr v-for="(item, index) in statuslist" :key="index" class="h-14">
                <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="text-left">{{ index+1 }}</div>
                </td>
                <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="text-left">
                    {{ item.chan_text }} {{ item.condition }} {{ item.level }} 
                  </div>
                </td>
                <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="text-left">{{ item.status_text }}</div>
                </td>
                <td class="px-2 first:pl-5 last:pr-5 py-3 whitespace-nowrap">
                  <div class="text-left"> {{  item.count }}</div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from "vue-i18n";

export default {
  name: 'AlarmStatusTable',
  props: ['channel'],
  setup(props) {
    const { t } = useI18n();
    const statuslist = ref([]);

    let interval = null;

    watch(() => props.channel, async (newVal, oldVal) => {
      if (newVal !== oldVal && oldVal !== undefined) {
        if (interval) {
          clearInterval(interval);
        }
        
        await fetchData(newVal);
        
        interval = setInterval(async () => {
          await fetchData(newVal);
        }, 60000);
      }
    });

    // 초기 실행은 onMounted에서
    onMounted(async () => {
      if (props.channel) {
        await fetchData(props.channel);
        
        interval = setInterval(async () => {
          await fetchData(props.channel);
        }, 60000);
      }
    });

    onUnmounted(() => {
      if (interval) {
        clearInterval(interval)
      }
    })

    const formatTimestamp = (timestamp) => {
      return new Date(timestamp * 1000).toLocaleString('ko-KR', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    };

    const fetchData = async (ch) => {
      try {
        const response = await axios.get("/api/getAlarmStatus/" + ch);
        if (response.data.success) {
          const resultDict = response.data.data;
         
          // 한 번에 처리
          statuslist.value = Object.entries(resultDict)
              .map(([key, value]) => {
                const parsed = JSON.parse(value);
                return {
                  index: parseInt(key) + 1,
                  ...parsed,
                  last_update_formatted: formatTimestamp(parsed.last_update)
                };
              })
              .sort((a, b) => a.index - b.index);
          
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    return {
      t,
      statuslist,
    }
  }
}
</script>

<style scoped>
/* 전체 래퍼에 최소 높이 설정 */
.alarm-wrapper {
  min-height: 500px;
  max-height: 600px; /* 최대 높이 설정 */
  display: flex;
  flex-direction: column;
}

/* 내부 컨텐츠 영역 */
.alarm-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 스크롤 컨테이너 */
.table-scroll-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: auto;
  position: relative;
}

/* 테이블 컨테이너 */
.table-container {
  display: flex;
  flex-direction: column;
}

/* 테이블 */
.table-container table {
  width: 100%;
}

/* 테이블 행 높이 고정 */
tr {
  height: 56px;
  max-height: 56px;
}

/* 테이블 셀 정렬 */
td {
  vertical-align: middle;
}

/* 헤더 고정 스타일 */
thead {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: rgb(249 250 251);
}

/* 다크모드 헤더 배경 */
.dark thead {
  background-color: rgba(17, 24, 39, 0.2);
}

/* 스크롤바 스타일 (선택사항) */
.table-scroll-container::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.table-scroll-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.table-scroll-container::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.table-scroll-container::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* 다크모드 스크롤바 */
.dark .table-scroll-container::-webkit-scrollbar-track {
  background: #374151;
}

.dark .table-scroll-container::-webkit-scrollbar-thumb {
  background: #6b7280;
}

.dark .table-scroll-container::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>