<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title"> Event & Fault</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'main' ? t('dashboard.diagnosis.subtitle_main') : t('dashboard.diagnosis.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- 콘텐츠 섹션 -->
    <div class="content-section">
      <div class="cards-grid">
        <!-- 첫 번째 이벤트 카드 -->
        <div class="event-card">
          <div class="event-card-header">
            <h3 class="event-card-title">Event Status</h3>
          </div>
          <div class="event-card-body">
            <div class="event-list">
              <div v-for="(item, index) in eventData" 
                   :key="`event1-${item.Name}`" 
                   class="event-item">
                <span class="event-label">{{ item.Title }}</span>
                <div class="event-status-badge" 
                     :class="getStatusColor2(getStatusText(item.Status))">
                  {{ getStatusCText(getStatusText(item.Status)) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 두 번째 이벤트 카드 -->
        <div class="event-card">
          <div class="event-card-header">
            <h3 class="event-card-title"> Fault Status</h3>
          </div>
          <div class="event-card-body">
            <div class="event-list">
              <div v-for="(item, index) in faultData" 
                   :key="`event2-${item.Name}`" 
                   class="event-item">
                <span class="event-label">{{ item.Title }}</span>
                <div class="event-status-badge" 
                     :class="getStatusColor2(getStatusText(item.Status))">
                  {{ getStatusCText(getStatusText(item.Status)) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watchEffect, onMounted, onUnmounted, watch } from 'vue'
import StatusItem from './StatusItem_Trans_Claude.vue'
import StatusItem2 from './StatusItem2.vue'
import { useSetupStore } from '@/store/setup'
import axios from 'axios'
import { useI18n } from 'vue-i18n'

export default {
  name: 'DashboardCard04',
  props: {
    channel: String
  },
  components: {
    StatusItem,
    StatusItem2,
  },
  setup(props) {
    const { t } = useI18n()
    const channel = ref(props.channel)
    const eventData = ref([])
    const faultData = ref([])
    const DiagEnable = ref(false)
    const setupStore = useSetupStore()
    const channelStatus = computed(() => setupStore.getChannelSetting)
    const asset = computed(() => setupStore.getAssetConfig)
    const assetTypes = ref('')
    const status = ref('Normal')

    let updateInterval = null

    const getStatusColor2 = (status) => {
      switch (status) {
        case 'OK': return 'status-ok'
        case 'Low': return 'status-low'
        case 'Medium': return 'status-medium'
        case 'High': return 'status-high'
        default: return 'status-default'
      }
    }

    const getStatusCText = (status) => {
      switch (status) {
        case 'OK': return t('diagnosis.tabContext.pqfe1')
        case 'Low': return t('diagnosis.tabContext.pqfe2')
        case 'Medium': return t('diagnosis.tabContext.pqfe3')
        case 'High': return t('diagnosis.tabContext.pqfe4')
        default: return t('diagnosis.tabContext.pqfe0')
      }
    }

    const getStatusText = (status) => {
      switch (status) {
        case 1: return 'OK'
        case 2: return 'Low'
        case 3: return 'Medium'
        case 4: return 'High'
        default: return 'No Data'
      }
    }

    const fetchEventData = async () => {
      if (!asset.value || (!asset.value.assetName_main && !asset.value.assetName_sub)) {
        return
      }
      
      const chName = channel.value == "main" 
        ? asset.value.assetName_main 
        : asset.value.assetName_sub

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getEvents/${chName}`)
          if (response.data.success) {
            eventData.value = response.data.data_status || []
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error)
          eventData.value = []
        }
      }
    }

    const fetchFaultData = async () => {
      if (!asset.value || (!asset.value.assetName_main && !asset.value.assetName_sub)) {
        return
      }
      
      const chName = channel.value == "main" 
        ? asset.value.assetName_main 
        : asset.value.assetName_sub

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getFaults/${chName}`)
          if (response.data.success) {
            faultData.value = response.data.data_status || []
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error)
          faultData.value = []
        }
      }
    }

    watch(asset, (newVal) => {
      if (newVal) {
        if(channel.value == 'main')
          assetTypes.value = newVal.assetType_main
        else
          assetTypes.value = newVal.assetType_sub
        
        fetchEventData()
        fetchFaultData();
        if (updateInterval) {
          clearInterval(updateInterval)
        }
        
        updateInterval = setInterval(async () => {
          await fetchEventData()
          await fetchFaultData();
        }, 300000) // 5분
      }
    }, { immediate: true })

    onMounted(async () => {
      await setupStore.checkSetting()
    })

    onUnmounted(() => {
      if (updateInterval) {
        clearInterval(updateInterval)
      }
    })

    watchEffect(() => {
      if(channel.value == 'main')
        DiagEnable.value = channelStatus.value.MainDiagnosis
      else
        DiagEnable.value = channelStatus.value.SubDiagnosis
    })

    return {
      channel,
      channelStatus,
      DiagEnable,
      fetchEventData,
      fetchFaultData,
      asset,
      status,
      assetTypes,
      getStatusColor2,
      getStatusCText,
      getStatusText,
      t,
      eventData,
      faultData,
    }
  }
}
</script>

<style scoped>
/* 메인 카드 컨테이너 */
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-5;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-50/50 to-purple-50/50 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

.channel-info {
  @apply flex items-center;
}

.channel-text {
  @apply text-xs font-semibold text-gray-400 dark:text-white uppercase;
}

/* 콘텐츠 섹션 */
.content-section {
  @apply p-4;
}

/* 카드 그리드 */
.cards-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-4;
}

/* 이벤트 카드 */
.event-card {
  @apply bg-white dark:bg-gray-800;
  @apply rounded-lg border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-shadow;
  @apply overflow-hidden;
}

.event-card-header {
  @apply px-4 py-3;
  @apply bg-gray-50 dark:bg-gray-700;
  @apply border-b border-gray-200 dark:border-gray-600;
}

.event-card-title {
  @apply text-sm font-bold text-gray-700 dark:text-gray-200;
}

.event-card-body {
  @apply p-4;
}

/* 이벤트 리스트 */
.event-list {
  @apply space-y-3;
}

.event-item {
  @apply flex items-center justify-between;
  @apply p-2 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply hover:bg-gray-100 dark:hover:bg-gray-700;
  @apply transition-colors;
}

.event-label {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
  @apply flex-1 mr-3;
}

/* 상태 배지 */
.event-status-badge {
  @apply px-3 py-1 rounded-full;
  @apply text-xs font-semibold;
  @apply flex-shrink-0;
  @apply transition-all duration-200;
}

/* 상태별 색상 */
.status-ok {
  @apply bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300;
}

.status-low {
  @apply bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-300;
}

.status-medium {
  @apply bg-orange-100 text-orange-700 dark:bg-orange-900/40 dark:text-orange-300;
}

.status-high {
  @apply bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300;
}

.status-default {
  @apply bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300;
}

/* 호버 효과 */
.event-status-badge:hover {
  @apply scale-105;
}

/* 반응형 디자인 */
@media (max-width: 768px) {
  .cards-grid {
    @apply grid-cols-1;
  }
  
  .event-item {
    @apply flex-col items-start gap-2;
  }
  
  .event-label {
    @apply mb-1;
  }
  
  .event-status-badge {
    @apply w-full text-center;
  }
}

/* 데이터 없을 때 */
.event-list:empty::after {
  content: "No event data available";
  @apply block text-center text-gray-500 dark:text-gray-400;
  @apply py-8 text-sm;
}
</style>