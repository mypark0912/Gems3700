<template>
  <div class="status-card">
    <div class="card-content">
      <div class="grid grid-cols-12 gap-4 items-center">
        <!-- 장비 정보 섹션 - 왼쪽 고정 -->
        <div class="col-span-3">
          <div class="flex items-center gap-3">
            <div class="equipment-avatar">
              <img class="avatar-image" :src="motorImageSrc" :alt="stData.devType" />
            </div>
            <div class="equipment-info">
              <h3 class="equipment-name">{{ stData.devNickname }}</h3>
            </div>
          </div>
        </div>

        <!-- 메트릭 섹션 - 오른쪽 정렬 (가동중 포함) -->
        <div class="col-span-9">
          <div class="metrics-section">
            <!-- 가동중 배지 -->
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

            <!-- Transformer -->
            <template v-if="stData.devType.includes('Transformer')">
              <template v-if="hasTempData">
                <div v-for="(label, index) in ['R', 'S', 'T']" :key="index" class="metric-box temperature">
                  <div class="metric-main">
                    <span class="metric-value">{{ transData.Temp?.[index]?.toFixed(2) ?? '-' }}</span>
                    <span class="metric-unit">℃</span>
                  </div>
                  <div class="metric-label">Temp {{ label }}</div>
                </div>
              </template>
              <div class="metric-box load" :class="{ 'max-w-[160px]': !hasTempData }">
                <div class="metric-main">
                  <span class="metric-value">{{ LoadRate }}</span>
                  <span class="metric-unit">%</span>
                </div>
                <div class="metric-label">{{ t('dashboard.transDiag.LoadFactor') }}</div>
              </div>
            </template>

            <!-- VFD -->
            <template v-else-if="stData.devType=='VFD'">
              <template v-if="hasTempData">
                <div v-for="(label, index) in ['R', 'S', 'T']" :key="index" class="metric-box temperature">
                  <div class="metric-main">
                    <span class="metric-value">{{ transData.Temp?.[index]?.toFixed(2) ?? '-' }}</span>
                    <span class="metric-unit">℃</span>
                  </div>
                  <div class="metric-label">Temp {{ label }}</div>
                </div>
              </template>
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
    status: {
      type: Boolean,
      default: true
    }
  },
  setup(props) {
    const { t } = useI18n();
    const setupStore = useSetupStore();
    const route = useRouter();
    
    const channel = ref(props.channel);
    const isRunning = ref(props.status);
    
    const computedChannel = computed(() => {
      if (channel.value == 'Main' || channel.value == 'main')
        return 'Main';
      else
        return 'Sub';
    })
    
    const stData = ref(props.data);
    const transData = ref(props.transData);
    const LoadRate = ref(0);

    // Temp 데이터 유무 체크
    const hasTempData = computed(() => {
      return transData.value?.Temp?.length > 0
    })

    // Transformer 타입 체크
    const isTransformer = computed(() => {
      return stData.value.devType?.includes('Transformer') || false;
    });

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
          return pumpImg;
        case 'Fan':
          return fanImg;
        case 'Compressor':
          return compImg;
        case 'PSupply':
          return powerImg;
        case 'PowerSupply':
          return powerImg;
        case 'PrimaryTransformer':
          return transImg;
        case 'Transformer':
          return transImg;
        default:
          return defaultImg;
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
      if (computedChannel.value == 'Main' && stData.value.devType?.includes('Transformer')) {
        kva = setupStore.getMkva;
      }
      if (computedChannel.value == 'Sub' && stData.value.devType?.includes('Transformer')) {
        kva = setupStore.getSkva;
      }
      return kva;
    })

    watch(
      () => LoadFactor.value,
      (newVal) => {
        if (newVal > 0 && transData.value?.Stotal) {
          LoadRate.value = (((transData.value.Stotal/1000) / newVal) * 100).toFixed(2);
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
      isRunning,
      isTransformer,
      hasTempData,
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

/* 상태 배지 */
.equipment-status {
  @apply flex items-center flex-shrink-0;
}

.status-badge {
  @apply flex items-center gap-1.5 px-2.5 py-1;
  @apply rounded-full font-semibold text-xs;
  @apply border-2 transition-all duration-300;
  @apply shadow-sm whitespace-nowrap;
}

.status-indicator {
  @apply w-1.5 h-1.5 rounded-full animate-pulse;
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
  @apply bg-gray-400 animate-none;
}

/* 메트릭 섹션 */
.metrics-section {
  @apply flex items-center gap-3 w-full justify-end;
}

.metric-box {
  @apply flex flex-col items-center justify-center px-4 py-2 rounded-lg;
  @apply bg-gray-50 dark:bg-gray-700/50;
  @apply border border-gray-200 dark:border-gray-600;
  @apply flex-1 max-w-[160px] h-14;
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
  @apply flex items-baseline gap-1 leading-none;
}

.metric-value {
  @apply text-sm font-bold text-gray-900 dark:text-white leading-none;
}

.metric-unit {
  @apply text-xs font-medium text-gray-500 dark:text-gray-400 leading-none;
}

.metric-label {
  @apply text-xs font-medium text-gray-500 dark:text-gray-400;
  @apply leading-tight whitespace-nowrap mt-1;
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
  .grid > div:first-child {
    @apply col-span-4;
  }
  
  .grid > div:last-child {
    @apply col-span-8;
  }

  .status-badge {
    @apply px-2 py-0.5 text-[10px] gap-1;
  }

  .status-indicator {
    @apply w-1 h-1;
  }
  
  .metric-box {
    @apply h-12 px-2 py-1;
  }
  
  .metric-value {
    @apply text-xs;
  }
  
  .metric-label {
    @apply text-xs;
  }
}

@media (max-width: 640px) {
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
  }
}
</style>