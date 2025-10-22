<template>
  <div class="sidebar-panel">
    <!-- 전력 품질 진단 섹션 -->
    <div class="panel-section">
      <div class="section-header">
        <h3 class="section-title">{{ t('dashboard.channel.title') }}</h3>
        <span class="channel-badge">
          Main channel
        </span>
      </div>
      
      <div class="section-content">
        <StatusItem :channel="channel" :data="pqData" mode="pq" />
      </div>
    </div>

    <!-- 이벤트 및 장애 상태 섹션 -->
    <div class="panel-section">
      <div class="section-content">
        <!-- 2컬럼 그리드로 이벤트와 장애 상태 배치 -->
        <div class="cards-grid">
          <!-- 이벤트 상태 카드 -->
          <div class="event-card">
            <div class="event-card-header">
              <h4 class="event-card-title">{{ t('dashboard.event.eventstatus') }}</h4>
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

          <!-- 장애 상태 카드 -->
          <div class="event-card">
            <div class="event-card-header">
              <h4 class="event-card-title">{{ t('dashboard.event.faultstatus') }}</h4>
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
  </div>
</template>

<script>
import { ref, computed, watchEffect, onMounted, onUnmounted, watch } from 'vue'
import StatusItem from './StatusItem_Trans_Claude.vue'
import { useSetupStore } from '@/store/setup'
import axios from 'axios'
import { useI18n } from 'vue-i18n'

export default {
  name: 'IntegratedSidebar',
  props: {
    channel: String,
    data: Object
  },
  components: {
    StatusItem,
  },
  setup(props) {
    const { t } = useI18n()
    const setupStore = useSetupStore()
    const AssetInfo = computed(() => setupStore.getAssetConfig)

    // 기본 데이터
    const channel = ref(props.channel)
    const stData = ref(props.data)
    
    // PQ 데이터
    const pqData = ref({
      devName: '',
      devStatus: -2
    })

    // 이벤트/장애 데이터
    const eventData = ref([])
    const faultData = ref([])
    
    let updateInterval = null

    // 계산된 속성
    const computedChannel = computed(() => {
      if (channel.value == 'Main' || channel.value == 'main')
        return 'Main'
      else
        return 'Sub'
    })

    // 상태 관련 함수들
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

    // PQ 데이터 가져오기
    const fetchPQData = async () => {
      if (!AssetInfo.value || (!AssetInfo.value.assetName_main && !AssetInfo.value.assetName_sub)) {
        return
      }
      
      const chName = channel.value == 'main' ? AssetInfo.value.assetName_main : AssetInfo.value.assetName_sub
      
      if (chName != '') {
        try {
          const response = await axios.get(`/api/getPQStatus/${chName}`)
          if (response.data.status >= 0) {
            pqData.value.devName = response.data.item
            pqData.value.devStatus = response.data.status
          }
        } catch (error) {
          console.log("PQ 데이터 가져오기 실패:", error)
        }
      }
    }

    // 이벤트 데이터 가져오기
    const fetchEventData = async () => {
      if (!AssetInfo.value || (!AssetInfo.value.assetName_main && !AssetInfo.value.assetName_sub)) {
        return
      }
      
      const chName = channel.value == "main" 
        ? AssetInfo.value.assetName_main 
        : AssetInfo.value.assetName_sub

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getEvents/${chName}`)
          if (response.data.success) {
            eventData.value = response.data.data_status || []
          }
        } catch (error) {
          console.log("이벤트 데이터 가져오기 실패:", error)
          eventData.value = []
        }
      }
    }

    // 장애 데이터 가져오기
    const fetchFaultData = async () => {
      if (!AssetInfo.value || (!AssetInfo.value.assetName_main && !AssetInfo.value.assetName_sub)) {
        return
      }
      
      const chName = channel.value == "main" 
        ? AssetInfo.value.assetName_main 
        : AssetInfo.value.assetName_sub

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getFaults/${chName}`)
          if (response.data.success) {
            faultData.value = response.data.data_status || []
          }
        } catch (error) {
          console.log("장애 데이터 가져오기 실패:", error)
          faultData.value = []
        }
      }
    }

    // 모든 데이터 가져오기
    const fetchAllData = async () => {
      await Promise.all([
        fetchPQData(),
        fetchEventData(),
        fetchFaultData()
      ])
    }

    // props.data 변경 감지
    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          fetchPQData()
        }
      },
      { immediate: true }
    )

    // AssetInfo 변경 감지
    watch(AssetInfo, (newVal) => {
      if (newVal) {
        fetchAllData()
        
        if (updateInterval) {
          clearInterval(updateInterval)
        }
        
        // 5분마다 업데이트
        updateInterval = setInterval(fetchAllData, 300000)
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

    return {
      channel,
      computedChannel,
      pqData,
      eventData,
      faultData,
      getStatusColor2,
      getStatusCText,
      getStatusText,
      t,
    }
  }
}
</script>

<style scoped>
.sidebar-panel {
  @apply h-full flex flex-col;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

.panel-section {
  @apply flex-1;
}

.panel-section:first-child {
  flex: 2; /* 상단 섹션을 2배 크기로 */
}

.panel-section:last-child {
  flex: 1; /* 하단 섹션은 기본 크기 */
}

.section-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-50/50 to-purple-50/50 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
  @apply flex justify-between items-center;
}

.section-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

.channel-badge {
  @apply text-xs font-semibold text-gray-400 dark:text-white uppercase;
}

.section-content {
  @apply p-4;
}

/* 카드 그리드 - 2컬럼 배치 */
.cards-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-3;
}

/* 이벤트 카드 스타일 */
.event-card {
  @apply bg-white dark:bg-gray-800;
  @apply rounded-lg border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-shadow;
  @apply overflow-hidden;
  @apply mb-3 last:mb-0;
}

.event-card-header {
  @apply px-3 py-2;
  @apply bg-gray-50 dark:bg-gray-700;
  @apply border-b border-gray-200 dark:border-gray-600;
}

.event-card-title {
  @apply text-xs font-bold text-gray-700 dark:text-gray-200;
}

.event-card-body {
  @apply p-3;
}

.event-list {
  @apply space-y-2;
}

.event-item {
  @apply flex items-center justify-between;
  @apply p-2 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply hover:bg-gray-100 dark:hover:bg-gray-700;
  @apply transition-colors;
}

.event-label {
  @apply text-xs font-medium text-gray-700 dark:text-gray-300;
  @apply flex-1 mr-2;
  @apply truncate;
}

.event-status-badge {
  @apply px-2 py-1 rounded-full;
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

/* 데이터 없을 때 */
.event-list:empty::after {
  content: "No data available";
  @apply block text-center text-gray-500 dark:text-gray-400;
  @apply py-4 text-xs;
}
</style>