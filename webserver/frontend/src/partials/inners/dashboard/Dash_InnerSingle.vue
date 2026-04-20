<template>
    <div class="data-table">
      <div class="table-header">
        <h4 class="table-title">{{ t('dashboard.singleinfo.OperationStatus') || '운용 현황' }}</h4>
      </div>
      <div class="table-content">
        <!-- 로딩 상태 -->
        <div v-if="isLoading" class="loading-state">
          <div class="spinner"></div>
          <span>{{ t('common.loading') || 'Loading...' }}</span>
        </div>
        
        <!-- 데이터 없음 상태 -->
        <div v-else-if="!hasData" class="empty-state">
          <span>{{ t('common.noData') || 'No data available' }}</span>
        </div>
        
        <!-- 데이터 표시 -->
        <template v-else>
          <!-- 운전시간 -->
          <div class="data-row">
            <span class="data-label">{{ t('dashboard.singleinfo.Operatingtime') }}</span>
            <span class="data-value">
              {{ runhours }} <span class="data-unit">hrs</span>
            </span>
          </div>
          
          <!-- 실시간 데이터 -->
          <div v-for="(item, index) in transData" :key="index" class="data-row">
            <span class="data-label">
              <template v-if="item.Title && item.Title.length < 10">
                {{ item.Assembly }} : {{ item.Title }}
              </template>
              <template v-else>
                {{ item.Title || '' }}
              </template>
            </span>
            <span class="data-value">
              {{ isNaN(item.Value) ? 0 : (item.Value?.toFixed(1) || 0) }}
              <span class="data-unit">{{ item.Unit || '' }}</span>
            </span>
          </div>
        </template>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { useI18n } from 'vue-i18n'
  import axios from 'axios'
  
  export default {
    name: 'Dash_InnerSingle',
    props: {
      channel: String,   // 채널 정보
    },
    setup(props) {
      const { t } = useI18n()
      let updateInterval = null;
      
      // 반응형 데이터
      const channel = ref(props.channel)
      const transData = ref([])
      const runhours = ref(0)
      const isLoading = ref(false)
      
      // 채널 정보 계산
      const computedChannel = computed(() => {
        if (channel.value == 'Main' || channel.value == 'main')
          return 'Main'
        else
          return 'Sub'
      })
      
      // 데이터 존재 여부
      const hasData = computed(() => {
        return runhours.value > 0 || transData.value.length > 0
      })
  
      const fetchRealData = async () => {
        isLoading.value = true
        try {
          const response = await axios.get(`/api/getRealTimeSingle/${computedChannel.value}`);

          if (response.data.success) {
            transData.value = response.data.data || [];
            runhours.value = response.data.runhours || 0;
          }else{
            console.log('No Data');
          }
        }catch (error) {
          console.log("데이터 가져오기 실패:", error);
        } finally {
          isLoading.value = false
        }
      };
  
      onMounted(async()=>{
        await fetchRealData();
        
        // 기존 타이머 정리
        if (updateInterval) {
          clearInterval(updateInterval);
          updateInterval = null;
        }
        
        // 5분마다 자동 갱신
        updateInterval = setInterval(async () => {
          await fetchRealData();
        }, 300000);  // 5분
      });
  
      onUnmounted(() => {
        if (updateInterval) {
          clearInterval(updateInterval);
          updateInterval = null;
        }
      });
  
      return {
        // 데이터
        transData,
        computedChannel,
        runhours,
        isLoading,
        hasData,
        // 계산된 속성
        channel,
        // 함수
        t,
        fetchRealData,
      }
    }
  }
  </script>
  
  <style scoped>
  .data-table {
    @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
    @apply border border-gray-200 dark:border-gray-600;
    @apply h-full;
  }
  
  .table-header {
    @apply bg-gray-100 dark:bg-gray-600 px-3 py-2 border-b border-gray-200 dark:border-gray-500;
    @apply flex justify-between items-center;
  }
  
  .table-title {
    @apply text-sm font-bold text-gray-700 dark:text-gray-200;
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
  
  /* 반응형 */
  @media (max-width: 768px) {
    .table-content {
      @apply p-2;
    }
    
    .data-row {
      @apply text-xs;
    }
    
    .data-label,
    .data-value {
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