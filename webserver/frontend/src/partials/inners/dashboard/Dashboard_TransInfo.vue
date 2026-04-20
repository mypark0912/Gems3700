<template>
  <div class="equipment-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">{{ t('dashboard.singleinfo.title') }}</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ channel == 'main' ? t('dashboard.diagnosis.subtitle_main') : t('dashboard.diagnosis.subtitle_sub') }}
          </span>
        </div>
      </header>
    </div>

    <!-- 설비 정보 -->
    <div class="equipment-info">
      <div class="equipment-avatar">
        <img class="avatar-image" :src="motorImageSrc" alt="설비 이미지" />
      </div>
      <div class="equipment-details">
      </div>
      <div class="equipment-status">
        <span 
          class="status-badge"
          :class="isRunning ? 'status-running' : 'status-stopped'"
        >
          <span class="status-indicator"></span>
          <span class="status-text">
            {{ isRunning ? t('dashboard.singleinfo.running') : t('dashboard.singleinfo.stopped') }}
          </span>
        </span>
      </div>
    </div>

    <!-- 데이터 섹션 -->
    <div class="data-section">
      <!-- 상태 카드들 - 수평 정렬 -->
      <div class="status-cards">
        <!-- 온도 카드 - 데이터 있을 때만 표시 -->
        <div v-if="displayData.Temp && displayData.Temp.length > 0" class="status-card temperature-card">
          <div class="status-value">
            <div class="temp-grid">
              <div v-for="(temp, index) in displayData.Temp" :key="index" class="temp-row">
                <span class="phase-label">{{ ['R', 'S', 'T'][index] }}</span>
                <span class="value-number">{{ Number(temp) < -900 ? '-': Number(temp).toFixed(2) }}</span>
                <span class="value-unit">℃</span>
              </div>
            </div>
          </div>
          <div class="status-label">{{ t('dashboard.transDiag.Temperature') }}</div>
        </div>

        <!-- 부하율 카드 -->
        <div class="status-card load-card">
          <div class="status-value" :class="getLoadRateClass(LoadRate)">
            <span class="value-number">{{ LoadRate }}</span>
            <span class="value-unit">%</span>
          </div>
          <div class="status-label">{{ t('dashboard.transDiag.LoadFactor') }}</div>
        </div>

        <!-- 역률 카드 -->
        <div class="status-card power-card">
          <div class="status-value">
            <span class="value-number">{{ displayData.PF }}</span>
            <span class="value-unit">%</span>
          </div>
          <div class="status-label">{{ t('dashboard.transDiag.PowerFactor') }}</div>
        </div>
      </div>

      <!-- 설비 사양 -->
      <div class="data-table">
        <div class="table-header">
          <h4 class="table-title">{{ t('dashboard.singleinfo.EquipmentInformation') }}</h4>
        </div>
        <div class="table-content">
          <div v-for="item in rawdata" :key="item.Name" class="data-row">
            <span class="data-label">{{ t(`dashboard.transDiag.${item.Name}`) }}</span>
            <span class="data-value">
              {{ item.Value }} <span class="data-unit">{{ item.Unit }}</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import sampleImg from '@/images/trans.png'
import { useSetupStore } from '@/store/setup'
import { useRealtimeStore } from '@/store/realtime'
import axios from 'axios'

export default {
  name: 'DashboardCard_TransInfo',
  props: {
    channel: String,
  },
  setup(props) {
    const { t } = useI18n()
    const setupStore = useSetupStore()
    const realtimeStore = useRealtimeStore()
    let updateInterval = null;
    // 반응형 데이터
    const LoadFactor = ref(-1)
    const rawdata = ref([])
    const isRunning = ref(0);
    // 계산된 속성
    const motorImageSrc = computed(() => sampleImg)

    // Store에서 실시간 데이터 가져오기
    const channelName = computed(() => 
      props.channel?.toLowerCase() === 'main' ? 'Main' : 'Sub'
    )

    const realtimeData = computed(() => 
      realtimeStore.getChannelData(channelName.value) || {}
    )

    // 화면 표시용 데이터 (설비 타입에 따라 매핑)
    const displayData = computed(() => {
      const data = realtimeData.value;
      console.log(data);
      if (!data || Object.keys(data).length === 0) {
        return { Temp: [], Ig: '-', Stotal: 0, PF: '-' }
      }

      return {
        Temp: data.Temp2 ?? [],
        Ig: data.Ig ?? '-',
        Stotal: data.Stotal ?? 0,
        PF: data.PF ?? '-',
      }
    })

    // 부하율 계산
    const LoadRate = computed(() => {
      if (LoadFactor.value > 0 && displayData.value.Stotal > 0) {
        return (((displayData.value.Stotal/1000) / LoadFactor.value) * 100).toFixed(2)
      }
      return 0
    })

    // 상태 클래스 함수들
    const getLoadRateClass = (rate) => {
      if (rate >= 90) return 'text-red-600 font-bold'
      if (rate >= 80) return 'text-orange-500 font-semibold'
      if (rate >= 60) return 'text-yellow-600 font-semibold'
      return 'text-green-600'
    }

    const getTemperatureStatus = (temp) => {
      if (temp >= 80) return 'status-critical'
      if (temp >= 60) return 'status-warning'
      return 'status-normal'
    }

    const getLoadStatus = (load) => {
      if (load >= 90) return 'status-critical'
      if (load >= 80) return 'status-warning'
      return 'status-normal'
    }

    const getPowerFactorStatus = (pf) => {
      if (pf < 80) return 'status-warning'
      if (pf < 90) return 'status-caution'
      return 'status-normal'
    }

    // 설비 정보 가져오기 (한 번만 호출)
    const fetchAsset = async () => {
      await setupStore.checkSetting()

      try {
        const chName = props.channel?.toLowerCase() === 'main' ? 'Main' : 'Sub'
        const response = await axios.get(`/api/getAsset/${chName}`)

        if (response.data.success) {
          rawdata.value = response.data.data

          // RatedKVA 찾기
          const ratedKVAItem = rawdata.value.find(item => item.Name === 'RatedKVA')
          if (ratedKVAItem) {
            LoadFactor.value = parseFloat(ratedKVAItem.Value)
          }
        } else {
          console.log('No Data')
        }
      } catch (error) {
        console.log('데이터 가져오기 실패:', error)
      }
    }

    const getStatus = async() => {
      const chName = props.channel?.toLowerCase() === 'main' ? 'Main' : 'Sub'
      const response = await axios.get(`/api/getEquipStatus/${chName}`);
      if (response.data.success) {
        isRunning.value = response.data.status;
      }
    }

    // 마운트 시 설비 정보 로드
    onMounted(async() => {
      await fetchAsset()
      await getStatus();
      updateInterval = setInterval(async () => {
        await fetchAsset()
        await getStatus();
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
      motorImageSrc,
      displayData,
      LoadRate,
      LoadFactor,
      rawdata,
      isRunning,
      // 함수
      t,
      getLoadRateClass,
      getTemperatureStatus,
      getLoadStatus,
      getPowerFactorStatus,
    }
  }
}
</script>

<style scoped>
.equipment-card {
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

.equipment-info {
  @apply flex items-center gap-3 px-5 py-3;
  @apply dark:from-gray-700/50 dark:to-gray-800;
}

.equipment-avatar {
  @apply relative;
}

.avatar-image {
  @apply w-12 h-12 rounded-xl object-cover;
  @apply shadow-sm border border-white dark:border-gray-600;
}

.equipment-details {
  @apply flex flex-col gap-1;
}

.equipment-name {
  @apply text-sm font-bold text-gray-800 dark:text-gray-100;
  @apply leading-tight;
}

.equipment-type {
  @apply text-xs font-medium text-gray-600 dark:text-gray-300;
  @apply bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300;
  @apply px-2 py-1 rounded-md inline-block;
}

.equipment-status {
  @apply flex items-center;
  @apply flex-shrink-0 ml-auto;
}

.status-badge {
  @apply flex items-center gap-2 px-3 py-1.5;
  @apply rounded-full font-semibold text-xs;
  @apply border-2 transition-all duration-300;
  @apply shadow-sm;
  @apply whitespace-nowrap;
}

.status-indicator {
  @apply w-2 h-2 rounded-full;
  @apply animate-pulse;
}

.status-running {
  @apply bg-green-50 dark:bg-green-900/30;
  @apply border-green-500 dark:border-green-600;
  @apply text-green-700 dark:text-green-300;
}

.status-running .status-indicator {
  @apply bg-green-500;
}

.status-stopped {
  @apply bg-gray-50 dark:bg-gray-700/30;
  @apply border-gray-400 dark:border-gray-500;
  @apply text-gray-700 dark:text-gray-300;
}

.status-stopped .status-indicator {
  @apply bg-gray-400;
  @apply animate-none;
}

/* 데이터 섹션 - 수평 배치 */
.data-section {
  @apply flex-1 p-4;
  @apply flex flex-row gap-4 items-stretch;
}

/* 상태 카드들 - 수평 정렬 */
.status-cards {
  @apply flex flex-wrap gap-3;
  @apply flex-1;
}

.status-card {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-md;
  @apply text-center;
  @apply flex flex-col justify-center;
  @apply flex-1 min-w-[80px];
}

.status-value {
  @apply flex items-baseline justify-center gap-1 mb-1;
}

.value-number {
  @apply text-2xl font-bold text-gray-800 dark:text-gray-100;
}

.value-unit {
  @apply text-sm font-semibold text-gray-500 dark:text-white;
}

.status-label {
  @apply text-sm font-medium text-gray-600 dark:text-white;
  @apply leading-tight;
}

/* 온도 그리드 (세로 배치) */
.temp-grid {
  @apply flex flex-col gap-1;
}

.temp-row {
  @apply flex items-center justify-center gap-1;
}

.phase-label {
  @apply text-xs font-semibold text-gray-500 dark:text-gray-400;
  @apply min-w-[16px];
}

.temp-row .value-number {
  @apply text-lg;
}

.temp-row .value-unit {
  @apply text-xs;
}

/* 개별 카드 호버 스타일 */
.temperature-card:hover {
  @apply bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700;
}

.load-card:hover {
  @apply bg-orange-50 dark:bg-orange-900/20 border-orange-200 dark:border-orange-700;
}

.power-card:hover {
  @apply bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-700;
}

/* 데이터 테이블 - 가로 너비만 축소 */
.data-table {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
  @apply border border-gray-200 dark:border-gray-600;
  @apply w-[240px] flex-shrink-0;
}

.table-header {
  @apply bg-gray-100 dark:bg-gray-600 px-3 py-2 border-b border-gray-200 dark:border-gray-500;
}

.table-title {
  @apply text-sm font-bold text-gray-700 dark:text-gray-200;
}

.table-content {
  @apply p-3;
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
  .equipment-info {
    @apply flex-col text-center gap-2;
  }
  
  .data-section {
    @apply flex-col;
  }
  
  .status-cards {
    @apply w-full;
  }
  
  .data-table {
    @apply w-full;
  }
}

/* 호버 효과 */
.equipment-card:hover .avatar-image {
  @apply transform scale-105;
}
</style>
