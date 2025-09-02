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
        <h3 class="equipment-name">
          {{ channel == 'main' ? AssetInfo.assetNickname_main : AssetInfo.assetNickname_sub }}
        </h3>
        <p class="equipment-type">
          {{ channel == 'main' ? AssetInfo.assetType_main : AssetInfo.assetType_sub }}
        </p>
      </div>
    </div>

    <!-- 데이터 테이블 -->
    <div class="data-section">
      <div class="data-grid">
        <!-- 운용 현황 카드들 -->
        <div class="status-cards">
          <!-- 온도 카드 -->
          <div class="status-card temperature-card">
            <div class="status-value">
              <span class="value-number">{{ transData.Temp }}</span>
              <span class="value-unit">℃</span>
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
              <span class="value-number">{{ transData.PF }}</span>
              <span class="value-unit">%</span>
            </div>
            <div class="status-label">{{ t('dashboard.transDiag.PowerFactor') }}</div>
          </div>

          <!-- 전류 카드 -->
          <div class="status-card current-card">
            <div class="status-value">
              <span class="value-number">{{ transData.Ig }}</span>
              <span class="value-unit">A</span>
            </div>
            <div class="status-label">{{ t('dashboard.transDiag.Ig') }}</div>
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
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import sampleImg from '@/images/transformer_m.png'
import { useSetupStore } from '@/store/setup'
import axios from 'axios'

export default {
  name: 'DashboardCard_TransInfo',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const { t } = useI18n()
    const setupStore = useSetupStore()
    
    // 반응형 데이터
    const channel = ref(props.channel)
    const transData = ref({})
    const LoadRate = ref(0)
    const LoadFactor = ref(-1)
    const rawdata = ref([])

    // 계산된 속성
    const AssetInfo = computed(() => setupStore.getAssetConfig)
    const motorImageSrc = computed(() => sampleImg)

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

    const fetchRealData = async (chType, chName) => {
      transData.value = {};
      try {
          const response = await axios.get(`/api/getRealTime/${chType}/${chName}`);

          if (response.data.success) {

            transData.value = response.data.data;
          }else{
            console.log('No Data');
          }
        }catch (error) {
          console.log("데이터 가져오기 실패:", error);
        } 
     };

    // API 호출
    const fetchAsset = async () => {
      if (!AssetInfo.value) {
        await setupStore.checkSetting()
      }

      const chName = channel.value.toLowerCase() === 'main' 
        ? AssetInfo.value.assetName_main 
        : AssetInfo.value.assetName_sub

      if (chName === '') {
        alert('등록된 설비가 없습니다.')
        return
      }

      const chType = channel.value.toLowerCase() === 'main' 
        ? AssetInfo.value.assetType_main 
        : AssetInfo.value.assetType_sub

      try {
        const response = await axios.get(`/api/getAsset/${chName}`)
        
        if (response.data.success) {
          rawdata.value = response.data.data
          
          // RatedKVA 찾기
          const ratedKVAItem = rawdata.value.find(item => item.Name === "RatedKVA")
          if (ratedKVAItem) {
            LoadFactor.value = parseFloat(ratedKVAItem.Value)
          }
          
          // 부하율 계산
          if (LoadFactor.value > 0 && transData.value?.Stotal) {
            LoadRate.value = ((transData.value.Stotal / LoadFactor.value) * 100).toFixed(2)
          }
        } else {
          console.log('No Data')
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error)
      }
      if(chType != 'Transformer'){
        fetchRealData(chType, chName);
      }
    };

    // props.data 감시
    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          const chType = channel.value.toLowerCase() === 'main' ? AssetInfo.value.assetType_main : AssetInfo.value.assetType_sub;

          if(chType == 'Transformer'){
            transData.value = {
              Temp: newData.Temp,
              Ig: newData.Ig,
              Stotal: newData.S4,
              PF: newData.PF4,
            }
          }
          fetchAsset();
        }
      },
      { immediate: true }
    )

    return {
      // 데이터
      motorImageSrc,
      channel,
      AssetInfo,
      transData,
      LoadRate,
      LoadFactor,
      rawdata,
      
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

/* 헤더 섹션 - 두 번째 컴포넌트와 동일한 스타일 */
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

.data-section {
  @apply flex-1 p-4;
}

.data-grid {
  @apply grid grid-cols-1 lg:grid-cols-2 gap-4;
  @apply items-start;
}

/* 상태 카드들 */
.status-cards {
  @apply grid grid-cols-2 gap-3;
  @apply h-full;
}

.status-card {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-md;
  @apply text-center;
  @apply flex flex-col justify-center;
  @apply min-h-[5rem];
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

/* 개별 카드 스타일 */
.temperature-card:hover {
  @apply bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-700;
}

.load-card:hover {
  @apply bg-orange-50 dark:bg-orange-900/20 border-orange-200 dark:border-orange-700;
}

.power-card:hover {
  @apply bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-700;
}

.current-card:hover {
  @apply bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-700;
}

.data-table {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg overflow-hidden;
  @apply border border-gray-200 dark:border-gray-600;
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

/* 반응형 개선 */
@media (max-width: 768px) {
  .equipment-info {
    @apply flex-col text-center gap-2;
  }
  
  .data-grid {
    @apply grid-cols-1;
  }
  
  .status-cards {
    @apply grid-cols-2 gap-2;
  }
  
  .status-card {
    @apply p-3;
    @apply h-20;
    @apply min-h-[5rem];
  }
  
  .value-number {
    @apply text-xl;
  }
}

@media (max-width: 480px) {
  .status-cards {
    @apply grid-cols-1;
  }
}

/* 호버 효과 */
.equipment-card:hover .avatar-image {
  @apply transform scale-105;
}
</style>