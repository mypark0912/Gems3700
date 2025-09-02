<template>
  <div class="status-card">
    <div class="card-content">
      <!-- 12 컬럼 그리드로 변경 -->
      <div class="grid grid-cols-12 gap-4 items-center">
        <!-- 장비 정보 섹션 - 6 컬럼 -->
        <div class="col-span-6">
          <div class="flex items-center gap-3">
            <div class="equipment-avatar">
              <img class="avatar-image" :src="motorImageSrc" :alt="stData.devType" />
            </div>
            <div class="equipment-info">
              <h3 class="equipment-name">{{ stData.devNickname }}</h3>
              <!-- <div class="equipment-type">{{ stData.devType }}</div> -->
            </div>
          </div>
        </div>
 
        <!-- 메트릭 섹션 - 6 컬럼 -->
        <div class="col-span-6">
          <div class="metrics-section">
            <!-- Transformer -->
            <template v-if="stData.devType.includes('Transformer')">
              <div class="metric-box temperature">
                <div class="metric-label">{{ t('dashboard.transDiag.Temperature') }}</div>
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Temp?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">℃</span>
                </div>
              </div>
              <div class="metric-box load">
                <div class="metric-label">{{ t('dashboard.transDiag.LoadFactor') }}</div>
                <div class="metric-main">
                  <span class="metric-value">{{ LoadRate }}</span>
                  <span class="metric-unit">%</span>
                </div>
              </div>
              <div class="metric-box current">
                <div class="metric-label">{{ t('dashboard.transDiag.Ig') }}</div>
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Ig?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">A</span>
                </div>
              </div>
            </template>
 
            <!-- VFD -->
            <template v-else-if="stData.devType=='VFD'">
              <div class="metric-box temperature">
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Temp?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">℃</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.Temperature') }}</div>
              </div>
              <div class="metric-box load">
                <div class="metric-main">
                  <span class="metric-value">56.9</span>
                  <span class="metric-unit">%</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.LoadFactor') }}</div>
              </div>
              <div class="metric-box current">
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Ig?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">A</span>
                </div>
                <div class="metric-label">{{ t('dashboard.meter.current') }}</div>
              </div>
            </template>
 
            <!-- 기타 장비 -->
            <template v-else>
              <div v-if="true" class="metric-box hours">
                <div class="metric-main">
                  <span class="metric-value">132</span>
                  <span class="metric-unit">hrs</span>
                </div>
                <div class="metric-label">{{ t('dashboard.singleinfo.Operatingtime') }}</div>
              </div>
              <div class="metric-box hours">
                <div class="metric-main">
                  <span class="metric-value">{{stData.Ig}}</span>
                  <span class="metric-unit">A</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.Ig') }}</div>
              </div>
              <div class="metric-box primary">
                <div class="metric-main">
                  <span class="metric-value">{{ transData[0]?.Value?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">{{ transData[0]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[0]["Title"] }}</div>
              </div>
              <div v-if="transData.length > 2" class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ transData[1]?.Value?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">{{ transData[1]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[1]["Title"] }}</div>
              </div>
              <div v-if="transData.length > 2" class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ transData[2]?.Value?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">{{ transData[2]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[2]["Title"] }}</div>
              </div>
              <div v-else class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ transData[1]?.Value?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">{{ transData[1]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[1]["Title"] }}</div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
 </template>
 
 <script>
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
 import { ref, watch, computed, onMounted } from 'vue';
 import { useI18n } from 'vue-i18n'
 import { useSetupStore } from '@/store/setup'
 import { useRouter } from 'vue-router'
 
 export default {
  name: 'CompactStatusCard',
  props: {
    data: Object,
    channel: String,
    transData: Object,
  },
  setup(props) {
    const { t } = useI18n();
    const setupStore = useSetupStore();
    const route = useRouter();
    
    const channel = ref(props.channel);
    const computedChannel = computed(() => {
      if (channel.value == 'Main' || channel.value == 'main')
        return 'Main';
      else
        return 'Sub';
    })
    
    const stData = ref(props.data);
    const transData = ref(props.transData);
    const LoadRate = ref(0);
 
    const strList = computed(() => {
      return [
        { "text": t('dashboard.diagnosis.st0'), "css": 'bg-gray-500/20 text-gray-700 font-semibold' },
        { "text": t('dashboard.diagnosis.st1'), "css": 'bg-green-500/20 text-green-700 font-semibold' },
        { "text": t('dashboard.diagnosis.st2'), "css": 'bg-yellow-500/20 text-yellow-700 font-semibold' },
        { "text": t('dashboard.diagnosis.st3'), "css": 'bg-orange-500/20 text-orange-700 font-semibold' },
        { "text": t('dashboard.diagnosis.st4'), "css": 'bg-red-500/20 text-red-700 font-semibold' }
      ]
    })
    
    const motorImageSrc = computed(() => {
      switch (stData.value.devType) {
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
    });
 
    const goToEquipmentDetail = () => {
      route.push(`/diagnosis/${computedChannel.value}`)
    }
 
    const statusClass = computed(() => {
      switch (stData.value.devStatus) {
        case 1:
          return 'status-ok';
        case 2:
          return 'status-warning';
        case 3:
          return 'status-inspect';
        case 4:
          return 'status-repair';
        default:
          return 'status-no-data';
      }
    });
 
    const LoadFactor = computed(() => {
      let kva = -1;
      if (computedChannel.value == 'Main' && stData.value.devType == 'Transformer') {
        kva = setupStore.getMkva;
      }
      if (computedChannel.value == 'sub' && stData.value.devType == 'Transformer') {
        kva = setupStore.getSkva;
      }
      return kva;
    })
 
    watch(
      () => LoadFactor.value,
      (newVal) => {
        if (newVal > 0 && transData.value?.Stotal) {
          LoadRate.value = ((transData.value.Stotal / newVal) * 100).toFixed(1);
        }
      },
      { immediate: true }
    );
 
    const statusStr = computed(() => {
      switch (stData.value.devStatus) {
        case 1:
          return t('dashboard.diagnosis.st1');
        case 2:
          return t('dashboard.diagnosis.st2');
        case 3:
          return t('dashboard.diagnosis.st3');
        case 4:
          return t('dashboard.diagnosis.st4');
        default:
          return t('dashboard.diagnosis.st0');
      }
    });
 
    return {
      stData,
      motorImageSrc,
      statusStr,
      statusClass,
      transData,
      t,
      strList,
      LoadFactor,
      channel,
      LoadRate,
      goToEquipmentDetail,
      computedChannel,
    }
  }
 }
 </script>
 
 <style scoped>
 .status-card {
   @apply col-span-full sm:col-span-12 xl:col-span-12;
   @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
   @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
   @apply backdrop-blur-sm;
   @apply transition-all duration-300 hover:shadow-xl;
 }
 
 .card-content {
   @apply p-3;
 }
 
 /* 장비 정보 스타일 */
 .equipment-avatar {
   @apply relative flex-shrink-0;
 }
 
 .avatar-image {
   @apply w-10 h-10 rounded-lg object-cover;
   @apply shadow-sm border border-gray-200 dark:border-gray-600;
 }
 
 .equipment-info {
   @apply flex flex-col justify-center;
 }
 
 .equipment-name {
   @apply text-base font-bold text-gray-800 dark:text-gray-100;
   @apply cursor-pointer hover:text-blue-600 dark:hover:text-blue-400;
   @apply transition-colors duration-200 leading-tight;
   @apply whitespace-nowrap;
 }
 
 .equipment-type {
   @apply text-xs font-semibold text-blue-500 dark:text-blue-400;
   @apply bg-blue-100 dark:bg-blue-700 px-2 py-0.5 rounded;
   @apply inline-block w-fit mt-0.5;
 }
 
 /* 메트릭 섹션 - 6 컬럼 내에서 flex 사용 */
 .metrics-section {
   @apply flex items-center gap-3;
   @apply w-full justify-end;
 }
 
 .metric-box {
   @apply flex flex-row items-center justify-between px-4 py-2 rounded-lg; /* flex-col을 flex-row로, items-center를 justify-between으로 변경 */
   @apply bg-gray-50 dark:bg-gray-700/50;
   @apply border border-gray-200 dark:border-gray-600;
   @apply flex-1;
   @apply h-14;
   @apply transition-all duration-200 hover:shadow-sm;
 }
 
 /* 메트릭 박스 컬러 */
 .metric-box.temperature {
   @apply border-violet-300 bg-violet-50/50 dark:bg-violet-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-box.load {
   @apply border-blue-300 bg-blue-50/50 dark:bg-blue-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-box.current {
   @apply border-green-300 bg-green-50/50 dark:bg-green-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-box.hours {
   @apply border-purple-300 bg-purple-50/50 dark:bg-purple-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-box.primary {
   @apply border-indigo-300 bg-indigo-50/50 dark:bg-indigo-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-box.secondary {
   @apply border-pink-300 bg-pink-50/50 dark:bg-pink-900/20;
   @apply hover:shadow-sm hover:scale-105;
 }
 
 .metric-main {
   @apply flex items-baseline gap-1; /* justify-center 제거 */
   @apply leading-none;
 }
 
 .metric-value {
   @apply text-sm font-bold text-gray-900 dark:text-white;
   @apply leading-none;
 }
 
 .metric-unit {
   @apply text-xs font-medium text-gray-500 dark:text-white;
   @apply leading-none;
 }
 
 .metric-label {
   @apply text-sm font-semibold text-gray-600 dark:text-white;
   @apply leading-tight whitespace-nowrap;
   /* mt-1 제거 */
 }
 
 /* 반응형 */
 @media (max-width: 1024px) {
   .grid {
     @apply gap-3;
   }
   
   .metric-box {
     @apply px-3;
   }
 }
 
 @media (max-width: 768px) {
   /* 태블릿에서는 5:7 비율 */
   .grid > div:first-child {
     @apply col-span-5;
   }
   
   .grid > div:last-child {
     @apply col-span-7;
   }
   
   .metric-box {
     @apply h-12 px-2 py-1; /* max-w-[100px] 제거 */
   }
   
   .metric-value {
     @apply text-xs;
   }
   
   .metric-label {
     @apply text-xs;
   }
 }
 
 @media (max-width: 640px) {
   /* 모바일에서는 세로 배치 */
   .grid {
     @apply grid-cols-1;
   }
   
   .grid > div {
     @apply col-span-full;
   }
   
   .metrics-section {
     @apply justify-start overflow-x-auto;
     scrollbar-width: thin;
   }
   
   .metrics-section::-webkit-scrollbar {
     height: 4px;
   }
   
   .metrics-section::-webkit-scrollbar-thumb {
     @apply bg-gray-300 dark:bg-gray-600 rounded;
   }
   
   .metric-box {
     @apply flex-shrink-0 min-w-[80px];
     @apply flex-col items-center justify-center; /* 모바일에서는 다시 세로 정렬 */
   }
   
   .metric-label {
     @apply mt-1; /* 모바일에서만 상단 마진 추가 */
   }
 }
 </style>