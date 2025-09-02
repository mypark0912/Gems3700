<template>
  <div class="equipment-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">{{ t('dashboard.singleinfo.title') }}</h2>
        <div class="channel-info">
          <span class="channel-text">
            {{ computedChannel }} channel
          </span>
        </div>
      </header>
    </div>

    <!-- 설비 정보 -->
    <div class="equipment-info">
      <div class="equipment-avatar">
        <img class="avatar-image" :src="motorImageSrc" :alt="stData.devType" />
      </div>
      <div class="equipment-details">
        <h3 class="equipment-name" @click="goToEquipmentDetail">
          {{ computedChannel == 'Main' ? AssetInfo.assetNickname_main : AssetInfo.assetNickname_sub }}
        </h3>
        <p class="equipment-type">
          {{ computedChannel == 'Main' ? AssetInfo.assetType_main : AssetInfo.assetType_sub }}
        </p>
      </div>
    </div>

    <!-- 데이터 테이블 -->  
    <div class="data-section">
      <div class="data-grid">
        <!-- 운용 현황 카드들 - 동적 개수 (운전시간 + transData 배열 길이) -->
        <div class="status-cards">
          <!-- 운전시간 카드 (항상 첫 번째) -->

          <template v-if="computedType == 'Transformer'">
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
          </template>
          <template v-else>
            <div class="status-card hours-card">
            <div class="status-value">
              <span class="value-number">132</span>
              <span class="value-unit">hrs</span>
            </div>
            <div class="status-label">{{ t('dashboard.singleinfo.Operatingtime') }}</div>
          </div>
          <!-- transData 배열 기반 동적 카드들 -->
          <div 
            v-for="(item, index) in getVisibleTransData()" 
            :key="index"
            class="status-card"
            :class="getCardClass(index)"
          >
            <div class="status-value">
              <span class="value-number">{{ isNaN(item.Value) ? 0 : (item.Value?.toFixed(1) || 0) }}</span>
              <span class="value-unit">{{ item.Unit || '' }}</span>
            </div>
            <div v-if="item.Title.length < 10" class="status-label" :title="`${item.Assembly} : ${item.Title}`">{{ item.Assembly }} : {{ item.Title }}</div>
            <div v-else class="status-label" :title="`${item.Assembly} : ${item.Title}`">{{ item.Title }}</div>
          </div>
          </template>
        </div>

        <!-- 설비 사양 (우측 테이블은 필요시 추가) -->
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
import { ref, watch, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSetupStore } from '@/store/setup'
import { useRouter } from 'vue-router'
import axios from 'axios'
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'

export default {
  name: 'SingleChannelTransformerView',
  props: {
    data: Object,      // stData 
    channel: String,   // 채널 정보
  },
  setup(props) {
    const { t } = useI18n()
    const setupStore = useSetupStore()
    const router = useRouter()
    const AssetInfo = computed(() => setupStore.getAssetConfig)
    const rawdata = ref([]);
    // 반응형 데이터
    const channel = ref(props.channel)
    const stData = ref(props.data);
    const transData = ref({})
    const LoadRate = ref(0)
    const LoadFactor = ref(-1)
    // 채널 정보 계산
    const computedChannel = computed(() => {
      if (channel.value == 'Main' || channel.value == 'main')
        return 'Main'
      else
        return 'Sub'
    })
    const computedType = computed(()=> computedChannel.value == 'Main' ? AssetInfo.value.assetType_main: AssetInfo.value.assetType_sub)
    // 설비 이미지 결정
    const motorImageSrc = computed(() => {
      //let assetType = computedChannel.value == 'Main' ? AssetInfo.value.assetType_main : AssetInfo.value.assetType_sub
      switch (computedType.value) {
        case 'Motor':
          return motorImg;
        case 'MotorFeed':
          return motorImg;          
        case 'Pump':
          return pumpImg;//'/images/motor_pump.png';
        case 'Fan':
          return fanImg;
        case 'Compressor':
          return compImg;//'/images/fan_m.png';
        case 'PSupply':
          return powerImg;//'/images/power.png';
        case 'PowerSupply':
          return powerImg;//'/images/power.png';  
        case 'PrimaryTransformer':
          return transImg;//'/images/trans.png';       
        case 'Transformer':
          return transImg;//'/images/trans.png';                  
        default:
          return defaultImg;//'@/images/cleaned_logo.png'
      }
    })

    // 표시할 transData 배열 반환 (최대 3개까지, 운전시간 제외)
    const getVisibleTransData = () => {
      if (Array.isArray(transData.value)) {
        // 최대 3개까지만 표시 (운전시간 + 3개 = 총 4개)
        return transData.value.slice(0, 3)
      }
      return []
    }

    // 카드별 스타일 클래스 결정
    const getCardClass = (index) => {
      // index는 transData의 인덱스 (0, 1, 2)
      switch(index) {
        case 0: return 'primary-card'
        case 1: return 'secondary-card'
        case 2: return 'secondary-card'
        default: return 'secondary-card'
      }
    }

    // 설비 상세로 이동
    const goToEquipmentDetail = () => {
      router.push(`/diagnosis/${computedChannel.value}`)
    }

    const fetchRealData = async (chName, chType) => {

         if(chName != ''){

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
         }
     };

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
          
          if(chType == 'Transformer'){
            const ratedKVAItem = rawdata.value.find(item => item.Name === "RatedKVA")
            if (ratedKVAItem) {
              LoadFactor.value = parseFloat(ratedKVAItem.Value)
            }
            // 부하율 계산
            if (LoadFactor.value > 0 && transData.value?.Stotal) {
              LoadRate.value = ((transData.value.Stotal / LoadFactor.value) * 100).toFixed(2)
            }
          }

        } else {
          console.log('No Data')
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error)
      }
      if(chType != 'Transformer'){
        fetchRealData(chName, chType);
      }
    };

    // watch(
    //   () => props.data,
    //   (newVal) => {
    //     console.log(newVal);
    //     if (Object.keys(newVal).length > 0 ) {
    //       stData.value = newVal;
    //       fetchAsset();
    //     }
    //   },
    //   { immediate: true }
    // )
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
            //console.log(transData.value);
          }

        }else{
          if (Object.keys(newVal).length > 0 ) {
            stData.value = newVal;
          }
        }
        fetchAsset();
      },
      { immediate: true }
    )

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

    return {
      // 데이터
      stData,
      transData,
      computedChannel,
      rawdata,
      LoadFactor,
      LoadRate,
      // 계산된 속성
      motorImageSrc,
      AssetInfo,
      channel,
      computedType,
      // 함수
      t,
      getVisibleTransData,
      getCardClass,
      goToEquipmentDetail,
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
  @apply cursor-pointer hover:text-blue-600 dark:hover:text-blue-400;
  @apply transition-colors duration-200 leading-tight;
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

/* 상태 카드들 - 동적 그리드 */
.status-cards {
  @apply gap-3 h-fit;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  max-width: 100%;
}

/* 카드가 1개일 때 */
.status-cards:has(.status-card:only-child) {
  grid-template-columns: 1fr;
}

/* 카드가 2개일 때 */
.status-cards:has(.status-card:nth-child(2):last-child) {
  grid-template-columns: 1fr 1fr;
}

/* 카드가 3개일 때 */
.status-cards:has(.status-card:nth-child(3):last-child) {
  grid-template-columns: 1fr 1fr 1fr;
}

/* 카드가 4개일 때 */
.status-cards:has(.status-card:nth-child(4)) {
  grid-template-columns: 1fr 1fr;
}

.status-card {
  @apply bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-600;
  @apply transition-all duration-200 hover:shadow-md;
  @apply text-center;
  @apply flex flex-col justify-center;
  @apply min-h-[4rem];
}

.status-value {
  @apply flex items-baseline justify-center gap-1 mb-1;
}

.value-number {
  @apply text-xl font-bold text-gray-800 dark:text-gray-100;
}

.value-unit {
  @apply text-xs font-semibold text-gray-500 dark:text-white;
}

.status-label {
  @apply text-xs font-medium text-gray-600 dark:text-white;
  @apply leading-tight;
  @apply line-clamp-2;
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

.hours-card:hover {
  @apply bg-purple-50 dark:bg-purple-900/20 border-purple-200 dark:border-purple-700;
}

.primary-card:hover {
  @apply bg-indigo-50 dark:bg-indigo-900/20 border-indigo-200 dark:border-indigo-700;
}

.secondary-card:hover {
  @apply bg-pink-50 dark:bg-pink-900/20 border-pink-200 dark:border-pink-700;
}

.empty-card {
  @apply opacity-50;
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
    @apply p-2;
    @apply min-h-[3.5rem];
  }
  
  .value-number {
    @apply text-lg;
  }
  
  .status-label {
    @apply text-xs;
  }
}

@media (max-width: 480px) {
  .status-cards {
    @apply grid-cols-2 gap-2;
  }
  
  .status-card {
    @apply p-2;
    @apply min-h-[3rem];
  }
  
  .value-number {
    @apply text-base;
  }
  
  .status-label {
    @apply text-xs;
  }
}

/* 호버 효과 */
.equipment-card:hover .avatar-image {
  @apply transform scale-105;
}
</style>