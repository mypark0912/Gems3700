<template>
    <div class="equipment-card">
      <!-- 헤더 -->
      <div class="card-header">
        <header class="header-content">
          <h2 class="card-title">{{ t('dashboard.singleinfo.title') }}</h2>
          <div class="channel-info">
            <span class="channel-text">
              {{ computedChannel == 'Main' ? t('dashboard.meter.subtitle_main') : t('dashboard.meter.subtitle_sub') }}
            </span>
          </div>
        </header>
      </div>
  
      <!-- 설비 정보 -->
      <div class="equipment-info">
        <div class="equipment-avatar">
          <img class="avatar-image" :src="motorImageSrc" />
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
      <!-- 데이터 테이블 -->  
      <div class="data-section">
        <div class="data-grid">
            <Dash_InnerSingle :channel="computedChannel"/> 
            <Dash_InnerEquipment :channel="computedChannel" />
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { useI18n } from 'vue-i18n'
  import motorImg from '@/images/motor_m.png'
  import fanImg from '@/images/fan_m.png'
  import pumpImg from '@/images/pump_m.png'
  import compImg from '@/images/comp_m.png'
  import powerImg from '@/images/power_m.png'
  import defaultImg from '@/images/cleaned_logo.png'
  import transImg from '@/images/trans.png'
  import Dash_InnerSingle from '../dashboard/Dash_InnerSingle.vue'
  import Dash_InnerEquipment from '../dashboard/Dash_InnerEquipment.vue'
  import axios from 'axios'
  export default {
    name: 'SingleChannelTransformerView',
    components:{
        Dash_InnerSingle,
        Dash_InnerEquipment,
    },
    props: {
      channel: String,   // 채널 정보
    },
    setup(props) {
      const { t } = useI18n()
      //const rawdata = ref([]);
      // 반응형 데이터
      const channel = ref(props.channel)
      const isRunning = ref(false)
      let updateInterval = null;
      // 채널 정보 계산
      const computedChannel = computed(() => {
        if (channel.value == 'Main' || channel.value == 'main')
          return 'Main'
        else
          return 'Sub'
      })

      const getStatus = async() =>{
        const response = await axios.get(`/api/getEquipStatus/${computedChannel.value}`);
        if (response.data.success){
          //console.log(response.data);
          isRunning.value = response.data.status;
        }
      }

      onMounted(async()=>{
        await getStatus();
        updateInterval = setInterval(async () => {
          await getStatus();
        }, 300000);  // 5분
      })

      onUnmounted(() => {
        if (updateInterval) {
          clearInterval(updateInterval);
          updateInterval = null;
        }
      });
      
      // 설비 이미지 결정
      const motorImageSrc = computed(() => {
        return defaultImg;
      })

      return {
        // 데이터
        computedChannel,
        isRunning,
        //rawdata,
        // 계산된 속성
        motorImageSrc,
        channel,
        // 함수
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
    @apply justify-between;
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
    @apply flex-shrink-0;
  }
  
  .equipment-name {
    @apply text-sm font-bold text-gray-800 dark:text-gray-100;
    @apply cursor-pointer hover:text-blue-600 dark:hover:text-blue-400;
    @apply transition-colors duration-200 leading-tight;
  }
  
  .equipment-type {
    @apply text-xs font-medium text-gray-600 dark:text-gray-300;
    @apply bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300;
    @apply px-2 py-1 rounded-md;
    @apply inline-block w-auto;
  }

  /* 상태 배지 */
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
      @apply flex-wrap;
    }

    .equipment-status {
      @apply w-full justify-end;
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

    .status-badge {
      @apply px-2 py-1 text-[10px];
    }

    .status-indicator {
      @apply w-1.5 h-1.5;
    }
  }
  
  /* 호버 효과 */
  .equipment-card:hover .avatar-image {
    @apply transform scale-105;
  }
  </style>