<template>
        <div class="data-table">
        <div class="table-header">
          <h4 class="table-title">{{ t('dashboard.singleinfo.EquipmentInformation') }}</h4>
        </div>
        <div class="table-content">
          <!-- 로딩 상태 -->
          <div v-if="isLoading" class="loading-state">
            <div class="spinner"></div>
            <span>Loading...</span>
          </div>
          
          <!-- 데이터 없음 상태 -->
          <div v-else-if="rawdata.length == 0" class="empty-state">
            <span>No Data</span>
          </div>
          
          <!-- 데이터 표시 -->
          <template v-else>
            <div v-for="item in rawdata" :key="item.Name" class="data-row">
              <span class="data-label">{{ t(`dashboard.transDiag.${item.Name}`) }}</span>
              <span class="data-value">
                {{ item.Value }} <span class="data-unit">{{ item.Unit }}</span>
              </span>
            </div>
          </template>
        </div>
      </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSetupStore } from '@/store/setup'
import axios from 'axios'

export default {
  name: 'Dash_InnerEquipment',
  props: {
    channel: {
      type: String,
      required: true
    },
  },
  setup(props) {
    const { t } = useI18n()
    const setupStore = useSetupStore()
    const isLoading = ref(false);
    // 반응형 데이터
    const rawdata = ref([])
    
    // 채널 정보 계산
    const computedChannel = computed(() => {
      if (props.channel == 'Main' || props.channel == 'main')
        return 'Main'
      else
        return 'Sub'
    })
    
    

    // 설비 정보 가져오기
    const fetchAsset = async () => {
      if (isLoading.value) return // 중복 호출 방지

      await setupStore.checkSetting()

      isLoading.value = true

      try {
        const response = await axios.get(`/api/getAsset/${computedChannel.value}`)

        if (response.data.success) {
          rawdata.value = response.data.data
          console.log('rawdata', rawdata.value);
        } else {
          console.error('No Data')
        }
      } catch (error) {
        console.error("데이터 가져오기 실패:", error)
      } finally {
        isLoading.value = false
      }
    }
    
    // 부하율 계산
    // const calculateLoadRate = () => {
    //   if (LoadFactor.value > 0 && props.loadData?.stotal) {
    //     LoadRate.value = (props.loadData.stotal / LoadFactor.value) * 100
    //   }
    // }
    
    // 자동 새로고침 설정
    // const setupAutoRefresh = () => {
    //   if (props.refreshInterval > 0) {
    //     clearInterval(refreshTimer)
    //     refreshTimer = setInterval(() => {
    //       fetchAsset()
    //     }, props.refreshInterval)
    //   }
    // }
    
    // // 자동 새로고침 해제
    // const clearAutoRefresh = () => {
    //   if (refreshTimer) {
    //     clearInterval(refreshTimer)
    //     refreshTimer = null
    //   }
    // }
    
    // loadData prop 변경 감지
    // watch(
    //   () => props.loadData,
    //   (newData) => {
    //     if (newData?.stotal && LoadFactor.value > 0) {
    //       calculateLoadRate()
    //     }
    //   },
    //   { deep: true }
    // )
    
    // refreshTrigger 변경 감지 (수동 새로고침)
  
    
    // 컴포넌트 마운트
    onMounted(async() => {
      await fetchAsset()
    })
    
    
    return {
      // 데이터
      rawdata,
      isLoading,
      // 계산된 속성
      computedChannel,
      // 함수
      t,
      fetchAsset,
    }
  }
}
</script>

<style scoped>
.specification-section {
  @apply p-4 border-t border-gray-200/50 dark:border-gray-700/50;
}

.data-grid {
  @apply grid grid-cols-1 lg:grid-cols-2 gap-4;
  @apply items-start;
}

.data-table {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200;
}

.table-header {
  @apply bg-gray-100 dark:bg-gray-600 px-3 py-2 border-b border-gray-200 dark:border-gray-500;
  @apply flex justify-between items-center;
}

.table-title {
  @apply text-sm font-bold text-gray-700 dark:text-gray-200;
}

.refresh-button {
  @apply p-1 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700;
  @apply transition-all duration-200;
  @apply focus:outline-none focus:ring-2 focus:ring-blue-500;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.refresh-icon {
  @apply w-4 h-4 text-gray-600 dark:text-gray-300;
}

.refresh-button.loading .refresh-icon {
  @apply animate-spin;
}

.table-content {
  @apply p-3;
  @apply min-h-[100px];
}

.loading-state,
.empty-state {
  @apply flex flex-col items-center justify-center;
  @apply text-gray-500 dark:text-gray-400;
  @apply py-8;
}

.spinner {
  @apply w-8 h-8 border-2 border-gray-300 dark:border-gray-600;
  @apply border-t-blue-500 rounded-full animate-spin;
  @apply mb-2;
}

.data-row {
  @apply flex justify-between items-center;
  @apply py-1 border-b border-gray-200/50 dark:border-gray-600/50 last:border-b-0;
  @apply transition-colors duration-200 hover:bg-gray-100/50 dark:hover:bg-gray-600/30;
  @apply rounded px-2 -mx-2;
  @apply min-h-[2rem];
}

.load-rate-row {
  @apply mt-2 pt-2 border-t-2 border-gray-300 dark:border-gray-500;
  @apply font-semibold;
}

.data-label {
  @apply text-sm font-medium text-gray-600 dark:text-white;
  @apply flex-1;
}

.data-value {
  @apply text-sm font-bold text-gray-800 dark:text-gray-100;
  @apply flex items-baseline gap-1;
}

.data-unit {
  @apply text-xs font-semibold text-gray-500 dark:text-white;
}

/* 부하율 색상 클래스 - 직접 스타일 정의 */
.load-rate-critical {
  color: #dc2626; /* red-600 */
}
.dark .load-rate-critical {
  color: #f87171; /* red-400 */
}

.load-rate-warning {
  color: #f97316; /* orange-500 */
}
.dark .load-rate-warning {
  color: #fb923c; /* orange-400 */
}

.load-rate-caution {
  color: #ca8a04; /* yellow-600 */
}
.dark .load-rate-caution {
  color: #facc15; /* yellow-400 */
}

.load-rate-normal {
  color: #16a34a; /* green-600 */
}
.dark .load-rate-normal {
  color: #4ade80; /* green-400 */
}

.additional-info {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4;
  @apply border border-gray-200 dark:border-gray-600;
}

/* 반응형 */
@media (max-width: 768px) {
  .data-grid {
    @apply grid-cols-1;
  }
  
  .table-content {
    @apply p-2;
  }
  
  .data-row {
    @apply text-xs;
  }
}

/* 애니메이션 */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>