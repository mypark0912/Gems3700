<template>
    <div class="status-card">
      <div class="card-content">
        <!-- 첫 번째 줄: 장비 정보 -->
        <div class="first-row">
          <div class="equipment-section">
            <div class="equipment-avatar">
              <img class="avatar-image" :src="motorImageSrc" :alt="stData.devType" />
            </div>
            <div class="equipment-info">
              <h3 class="equipment-name" @click="goToEquipmentDetail">{{ stData.devNickname }}</h3>
              <div class="equipment-type">{{ stData.devType }}</div>
            </div>
          </div>
        </div>
  
        <!-- 두 번째 줄: 메트릭 데이터들 -->
        <div class="second-row">
          <div class="metrics-section">
            <!-- Transformer -->
            <template v-if="stData.devType.includes('Transformer')">
              <div class="metric-box temperature">
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Temp?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">℃</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.Temperature') }}</div>
              </div>
              <div class="metric-box load">
                <div class="metric-main">
                  <span class="metric-value">{{ LoadRate }}</span>
                  <span class="metric-unit">%</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.LoadFactor') }}</div>
              </div>
              <div class="metric-box current">
                <div class="metric-main">
                  <span class="metric-value">{{ transData.Ig?.toFixed(1) || 0 }}</span>
                  <span class="metric-unit">A</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.Ig') }}</div>
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
                  <span class="metric-value">{{ stData.runhour }}</span>
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
                  <span class="metric-value">{{ (parseFloat(transData[0]?.Value) || 0).toFixed(1) }}</span>
                  <span class="metric-unit">{{ transData[0]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[0]["Assembly"] }} : {{ transData[0]["Title"] }}</div>
              </div>
              <div v-if="transData.length > 2" class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ (parseFloat(transData[1]?.Value) || 0).toFixed(1) }}</span>
                  <span class="metric-unit">{{ transData[1]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[1]["Assembly"] }} : {{ transData[1]["Title"] }}</div>
              </div>
              <div v-if="transData.length > 2" class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ (parseFloat(transData[2]?.Value) || 0).toFixed(1)}}</span>
                  <span class="metric-unit">{{ transData[2]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[2]["Assembly"] }} : {{ transData[2]["Title"] }}</div>
              </div>
              <div v-else class="metric-box secondary">
                <div class="metric-main">
                  <span class="metric-value">{{ (parseFloat(transData[1]?.Value) || 0).toFixed(1) }}</span>
                  <span class="metric-unit">{{ transData[1]?.Unit || '' }}</span>
                </div>
                <div class="metric-label">{{ transData[1]["Assembly"] }} : {{ transData[1]["Title"] }}</div>
              </div>
            </template>
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
  
      // 세그먼트 활성화 계산 함수 제거 (불필요)
  
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
    @apply bg-white dark:bg-gray-800 shadow-sm rounded-xl;
    @apply border border-gray-200 dark:border-gray-700;
    @apply transition-all duration-300 hover:shadow-md;
  }
  
  .card-content {
    @apply flex flex-col p-3 gap-3;
  }
  
  /* 첫 번째 줄 */
  .first-row {
    @apply flex items-center;
    @apply min-h-[48px];
  }
  
  /* 장비 섹션 */
  .equipment-section {
    @apply flex items-center gap-3 w-full;
  }
  
  .equipment-avatar {
    @apply relative flex-shrink-0;
  }
  
  .avatar-image {
    @apply w-10 h-10 rounded-lg object-cover;
    @apply shadow-sm border border-gray-200 dark:border-gray-600;
  }
  
  .equipment-info {
    @apply flex flex-col min-w-0 flex-1;
  }
  
  .equipment-name {
    @apply text-base font-bold text-gray-800 dark:text-gray-100;
    @apply cursor-pointer hover:text-blue-600 dark:hover:text-blue-400;
    @apply transition-colors duration-200 leading-tight;
    @apply truncate;
  }
  
  .equipment-type {
    @apply text-xs font-semibold text-blue-500 dark:text-blue-400;
    @apply bg-blue-100 dark:bg-blue-700 px-2 py-1 rounded-lg;
    @apply inline-block w-fit mt-1;
  }
  
  /* 두 번째 줄 */
  .second-row {
    @apply flex items-center;
    @apply min-h-[70px];
  }
  
  /* 메트릭 섹션 */
  .metrics-section {
    @apply flex items-center gap-3 w-full;
  }
  
  .metric-box {
    @apply flex flex-col items-center justify-center text-center p-2 rounded-lg;
    @apply bg-gray-50 dark:bg-gray-700/50;
    @apply border border-gray-200 dark:border-gray-600;
    @apply flex-1 h-14;
    @apply transition-all duration-200 hover:shadow-sm;
  }
  
  .metric-box.temperature {
    @apply hover:border-red-300 hover:bg-red-50/50 dark:hover:bg-red-900/20;
  }
  
  .metric-box.load {
    @apply hover:border-blue-300 hover:bg-blue-50/50 dark:hover:bg-blue-900/20;
  }
  
  .metric-box.current {
    @apply hover:border-green-300 hover:bg-green-50/50 dark:hover:bg-green-900/20;
  }
  
  .metric-box.hours {
    @apply hover:border-purple-300 hover:bg-purple-50/50 dark:hover:bg-purple-900/20;
  }
  
  .metric-box.primary {
    @apply hover:border-indigo-300 hover:bg-indigo-50/50 dark:hover:bg-indigo-900/20;
  }
  
  .metric-box.secondary {
    @apply hover:border-pink-300 hover:bg-pink-50/50 dark:hover:bg-pink-900/20;
  }
  
  .metric-main {
    @apply flex items-baseline justify-center gap-1;
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
    @apply text-xs text-gray-600 dark:text-white;
    @apply mt-1 leading-tight;
  }
  
  /* 반응형 */
  @media (max-width: 1024px) {
    .metrics-section {
      @apply gap-2;
    }
    
    .metric-box {
      @apply h-12 p-1.5;
    }
    
    .metric-value {
      @apply text-xs;
    }
  }
  
  @media (max-width: 768px) {
    .card-content {
      @apply p-2 gap-2;
    }
    
    .equipment-section {
      @apply gap-2;
    }
    
    .avatar-image {
      @apply w-8 h-8;
    }
    
    .equipment-name {
      @apply text-sm;
    }
    
    .equipment-type {
      @apply text-xs px-1.5 py-0.5;
    }
    
    .metrics-section {
      @apply gap-2;
    }
    
    .metric-box {
      @apply h-10 p-1;
    }
    
    .metric-value {
      @apply text-xs;
    }
    
    .metric-label {
      @apply text-xs;
    }
  }
  
  @media (max-width: 640px) {
    .metrics-section {
      @apply flex-wrap;
    }
    
    .metric-box {
      @apply min-w-[80px];
    }
  }
  </style>