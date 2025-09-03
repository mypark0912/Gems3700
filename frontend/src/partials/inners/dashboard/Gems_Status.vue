<template>
    <div class="equipment-card">
      <!-- 헤더 -->
      <div class="card-header">
        <header class="header-content">
          <h2 class="card-title">{{ t('dashboard.channel.title') }}</h2>
          <div class="channel-info">
            <span class="channel-text">
              {{ computedChannel }} {{ t('dashboard.channel.channel') }}
            </span>
          </div>
        </header>
      </div>
  
      <!-- 데이터 섹션 -->  
      <div class="data-section">
        <!-- 설비 사양 - 한 줄 가로 배치 -->
        <div class="specification-section">
          <div class="section-header">
            <h4 class="section-title">{{ t('dashboard.channel.channelinfo') }}</h4>
          </div>
          <div class="spec-cards-wrapper">
            <div class="spec-cards-container">
              <div v-for="item in rawdata" :key="item.Name" class="spec-card">
                <div class="spec-value">
                  {{ item.Value }}{{ item.Unit }}
                </div>
                <div class="spec-label">
                  {{ t(`dashboard.transDiag.${item.Name}`) }}
                </div>
              </div>
            </div>
          </div>
        </div>
  
        <!-- 여기에 다른 컴포넌트들을 추가할 수 있습니다 -->
        <div class="additional-content">
            <StatusItem :channel="channel" :data="pqData" mode="pq" />
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
import StatusItem from './StatusItem_Trans_Claude.vue'
  
  export default {
    name: 'Gems_Status',
    components:{
        StatusItem
    },
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
      const pqData = ref({
       devName:'',
       devStatus: -2
     });
      // 채널 정보 계산
      const computedChannel = computed(() => {
        if (channel.value == 'Main' || channel.value == 'main')
          return 'Main'
        else
          return 'Sub'
      })
      
      const computedType = computed(()=> computedChannel.value == 'Main' ? AssetInfo.value.assetType_main: AssetInfo.value.assetType_sub)
      
  
  
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
      };

      const fetchPQData = async () => {
          if (!AssetInfo.value || (!AssetInfo.value.assetName_main && !AssetInfo.value.assetName_sub)) {
            //console.log("⏳ asset 준비 안됨. fetchData 대기중");
            //console.log(asset.value);
            return;
          }
         const chName = channel.value == 'main'? AssetInfo.value.assetName_main : AssetInfo.value.assetName_sub;
         const chType = channel.value == 'main'? AssetInfo.value.assetType_main : AssetInfo.value.assetType_sub;
         if(chName != ''){
           try {
             const response = await axios.get(`/api/getPQStatus/${chName}`);
             if (response.data.status >= 0) {
                pqData.value.devName = response.data.item;
                pqData.value.devStatus = response.data.status;
             }else{
               console.log('No Data');
             }
           }catch (error) {
             console.log("데이터 가져오기 실패:", error);
           } 
         }
     };
  
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
          }
          fetchAsset();
          fetchPQData();
        },
        { immediate: true }
      )
  
      return {
        // 데이터
        stData,
        transData,
        computedChannel,
        rawdata,
        LoadFactor,
        LoadRate,
        // 계산된 속성
        AssetInfo,
        channel,
        computedType,
        pqData,
        t,
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
  
  /* 데이터 섹션 */
  .data-section {
    @apply flex-1 p-4;
    @apply flex flex-col gap-4;
  }
  
  /* 설비 사양 섹션 */
  .specification-section {
    @apply bg-white dark:bg-gray-800 rounded-xl;
    @apply border border-gray-200 dark:border-gray-700;
    @apply shadow-sm;
    @apply overflow-hidden;
  }
  
  .section-header {
    @apply flex justify-between items-center;
    @apply px-4 py-3;
    @apply bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800;
    @apply border-b border-gray-200 dark:border-gray-600;
  }
  
  .section-title {
    @apply text-sm font-bold text-gray-700 dark:text-gray-200;
  }
  
  .channel-badge {
    @apply text-xs font-semibold;
    @apply px-2 py-1 rounded-full;
    @apply bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300;
  }
  
  /* 사양 카드 래퍼 - 가로 스크롤 */
  .spec-cards-wrapper {
    @apply p-3;
    @apply overflow-x-auto;
  }
  
  /* 사양 카드 컨테이너 - 한 줄 가로 배치 */
  .spec-cards-container {
    @apply flex gap-2;
    @apply min-w-fit;
  }
  
  /* 개별 사양 카드 */
  .spec-card {
    @apply flex-shrink-0;
    @apply bg-gradient-to-br from-gray-50 to-white dark:from-gray-700/50 dark:to-gray-800/50;
    @apply border border-gray-200 dark:border-gray-600;
    @apply rounded-lg px-4 py-3;
    @apply text-center;
    @apply transition-all duration-200;
    @apply hover:shadow-md hover:border-blue-300 dark:hover:border-blue-500;
    @apply cursor-default;
    @apply min-w-[100px];
  }
  
  .spec-value {
    @apply text-lg font-bold text-gray-900 dark:text-white;
    @apply mb-1;
  }
  
  .spec-label {
    @apply text-xs font-medium text-gray-500 dark:text-gray-400;
    @apply leading-tight;
    @apply whitespace-nowrap;
  }
  
  /* 추가 콘텐츠 영역 */
  .additional-content {
    @apply flex-1;
    /* 여기에 추가 콘텐츠 스타일 */
  }
  
  /* 스크롤바 스타일 */
  .spec-cards-wrapper::-webkit-scrollbar {
    @apply h-2;
  }
  
  .spec-cards-wrapper::-webkit-scrollbar-track {
    @apply bg-gray-100 dark:bg-gray-700;
    @apply rounded-full;
  }
  
  .spec-cards-wrapper::-webkit-scrollbar-thumb {
    @apply bg-gray-400 dark:bg-gray-500;
    @apply rounded-full;
    @apply hover:bg-gray-500 dark:hover:bg-gray-400;
  }
  
  /* 호버 효과 색상 */
  .spec-card:nth-child(4n+1):hover {
    @apply bg-gradient-to-br from-blue-50 to-white dark:from-blue-900/20 dark:to-gray-800/50;
  }
  
  .spec-card:nth-child(4n+2):hover {
    @apply bg-gradient-to-br from-green-50 to-white dark:from-green-900/20 dark:to-gray-800/50;
  }
  
  .spec-card:nth-child(4n+3):hover {
    @apply bg-gradient-to-br from-orange-50 to-white dark:from-orange-900/20 dark:to-gray-800/50;
  }
  
  .spec-card:nth-child(4n):hover {
    @apply bg-gradient-to-br from-purple-50 to-white dark:from-purple-900/20 dark:to-gray-800/50;
  }
  
  /* 반응형 디자인 */
  @media (max-width: 768px) {
    .equipment-info {
      @apply flex-col text-center gap-2;
    }
    
    .spec-cards-container {
      @apply gap-2;
    }
    
    .spec-card {
      @apply px-3 py-2;
      @apply min-w-[90px];
    }
    
    .spec-value {
      @apply text-base;
    }
  }
  
  @media (max-width: 480px) {
    .spec-card {
      @apply px-2 py-2;
      @apply min-w-[80px];
    }
    
    .spec-value {
      @apply text-sm;
    }
    
    .spec-label {
      @apply text-xs;
    }
  }
  
  /* 호버 효과 */
  .equipment-card:hover .avatar-image {
    @apply transform scale-105;
  }
  </style>