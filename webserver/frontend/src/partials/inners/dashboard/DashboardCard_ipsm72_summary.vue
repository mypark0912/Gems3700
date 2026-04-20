<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">IPSM 모듈 상태</h2>    
      </header>
    </div>

    <!-- 데이터 섹션 -->
    <div class="data-section">

      <!-- 전체 상태 요약 -->
      <div class="summary-section">
        <h3 class="subsection-title">상태 요약</h3>
        
        <div class="status-grid-wrapper">
          <div class="status-grid">
            <div class="status-card" v-for="status in statusSummary.slice(0, 4)" :key="status.type">
              <div class="status-left">
                <div class="status-indicator" :class="status.colorClass"></div>
                <div class="status-label">{{ status.label }}</div>
              </div>
              <div class="status-count">{{ status.count }}</div>
            </div>
          </div>
          <div class="status-grid" v-if="statusSummary.length > 4">
            <div class="status-card" v-for="status in statusSummary.slice(4)" :key="status.type">
              <div class="status-left">
                <div class="status-indicator" :class="status.colorClass"></div>
                <div class="status-label">{{ status.label }}</div>
              </div>
              <div class="status-count">{{ status.count }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 모듈별 분포 섹션 -->
      <div class="modules-section">
        <h3 class="subsection-title">모듈별 채널 상태 (CH1~72)</h3>
        <div class="chart-section">               
          <div class="module-list">
            <div v-for="module in modules" :key="module.id" class="module-item">
              <div class="module-header">
                <span class="module-name">{{ module.name }}</span>
                <span class="module-count">{{ module.total }}개</span>
              </div>
              <div class="module-bar-wrapper">
                <div class="module-bar">
                  <!-- 채널 1~72를 순서대로 표시 -->
                  <div 
                    v-for="channel in module.channels" 
                    :key="channel.id"
                    class="bar-segment" 
                    :class="getStatusClass(channel.status)"
                    :style="{ width: getChannelWidth() + '%' }"
                    :title="`CH${channel.id}: ${getStatusLabel(channel.status)}`"
                  ></div>
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
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ChannelStatusSummary',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const { t } = useI18n()
    
    // 상태 정의
    const statusDefinitions = {
      OFF: { label: 'OFF', colorClass: 'bg-gray-400' },
      ON: { label: '정상', colorClass: 'bg-green-500' },
      TRIP: { label: '차단', colorClass: 'bg-red-500' },
      OPR: { label: '운전', colorClass: 'bg-blue-500' },
      OC: { label: '과전류', colorClass: 'bg-orange-500' },
      OCG: { label: '지락', colorClass: 'bg-purple-500' },
      'OCG+': { label: '지락+', colorClass: 'bg-amber-700' },
      BAD: { label: '불량', colorClass: 'bg-gray-800 dark:bg-gray-600' }
    }

    // 채널 생성 함수 (1~72, 순서대로)
    const generateChannels = (moduleId) => {
      const channels = []
      
      for (let i = 1; i <= 72; i++) {
        let status = 'ON' // 기본값
        
        // 테스트 데이터 (실제로는 props.data에서 가져와야 함)
        if (moduleId === 1) {
          if (i === 2 || i === 11 || i === 23 || i === 35 || i === 47 || i === 59 || i === 68) status = 'OFF'
          else if ( i === 54) status = 'TRIP'
          else if (i === 15 ) status = 'OPR'
          else if (i === 19 || i === 51) status = 'OC'
          else if (i === 33 || i === 65) status = 'OCG'
          else if (i === 44) status = 'OCG+'
          else if (i === 71) status = 'BAD'
          else status = 'ON'
        } 
        // IPSM #2 - 다른 패턴으로 골고루 배치
        else {
          if (i === 4 || i === 14 || i === 26 || i === 38 || i === 50 || i === 61 || i === 70) status = 'OFF'
          else if (i === 10 ) status = 'TRIP'
          else if (i === 17 ) status = 'OPR'
          else if (i === 22 || i === 48) status = 'OC'
          else if (i === 36 || i === 63) status = 'OCG'
          else if (i === 53) status = 'OCG+'
          else if (i === 69) status = 'BAD'
          else status = 'ON'
        }
        
        
        channels.push({
          id: i,
          status: status
        })
      }
      return channels
    }
    
    // 모듈 데이터
    const modules = ref([
      {
        id: 1,
        name: 'IPSM #1',
        total: 72,
        channels: generateChannels(1)
      },
      {
        id: 2,
        name: 'IPSM #2',
        total: 72,
        channels: generateChannels(2)
      }
    ])

    // 전체 통합 상태 계산
    const totalStatus = computed(() => {
      const total = {
        OFF: 0,
        ON: 0,
        TRIP: 0,
        OPR: 0,
        OC: 0,
        OCG: 0,
        'OCG+': 0,
        BAD: 0
      }
      
      modules.value.forEach(module => {
        module.channels.forEach(channel => {
          total[channel.status]++
        })
      })
      
      return total
    })

    // 상태 요약 (카운트가 0보다 큰 것만)
    const statusSummary = computed(() => {
      return Object.entries(totalStatus.value)
        .filter(([_, count]) => count > 0)
        .map(([type, count]) => ({
          type,
          count,
          label: statusDefinitions[type].label,
          colorClass: statusDefinitions[type].colorClass
        }))
    })

    // 상태별 색상 클래스
    const getStatusClass = (status) => {
      return statusDefinitions[status]?.colorClass || 'bg-gray-400'
    }

    // 상태 레이블
    const getStatusLabel = (status) => {
      return statusDefinitions[status]?.label || 'Unknown'
    }

    // 각 채널의 너비 (72개로 나눔)
    const getChannelWidth = () => {
      return (100 / 72).toFixed(4)
    }

    return {
      t,
      modules,
      statusSummary,
      getStatusClass,
      getStatusLabel,
      getChannelWidth
    }
  }
}
</script>

<style scoped>
/* 기존 카드 스타일 - 높이 맞춤 */
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-4 xl:col-span-4;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply p-3 border-b border-gray-200/50 dark:border-gray-700/50;
  @apply bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

/* 데이터 섹션 */
.data-section {
  @apply p-4 space-y-3;
}

.subsection-title {
  @apply text-sm font-semibold text-gray-700 dark:text-white mb-2;
  @apply flex items-center gap-2;
}

.subsection-title::before {
  content: '';
  @apply w-2 h-2 bg-blue-500 rounded-full;
}

/* 상태 요약 그리드 */
.summary-section {
  flex-shrink: 0;
}

.status-grid-wrapper {
  @apply space-y-2;
}

.status-grid {
  @apply grid grid-cols-4 gap-2;
}

.status-card {
  @apply flex items-center justify-between px-2.5 py-1.5 bg-gray-50 dark:bg-gray-800/50 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
}

.status-left {
  @apply flex items-center gap-1.5;
}

.status-indicator {
  @apply w-2.5 h-2.5 rounded-full flex-shrink-0;
}

.status-label {
  @apply text-xs text-gray-600 dark:text-gray-400;
}

.status-count {
  @apply text-base font-bold text-gray-800 dark:text-white;
}

/* 모듈별 분포 섹션 */
.modules-section {
  margin-top: 0.5rem;
}

.chart-section {
  @apply bg-white dark:bg-gray-800/30 rounded-lg p-3;
  @apply border border-gray-200 dark:border-gray-700;
}

.module-list {
  @apply space-y-2.5;
}

.module-item {
  @apply space-y-1;
}

.module-header {
  @apply flex justify-between items-center;
}

.module-name {
  @apply text-sm font-semibold text-gray-700 dark:text-gray-300;
}

.module-count {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.module-bar-wrapper {
  @apply bg-gray-100 dark:bg-gray-700/50 rounded-lg overflow-hidden;
  @apply shadow-inner;
  height: 26px;
}

.module-bar {
  @apply flex h-full w-full;
}

.bar-segment {
  @apply h-full transition-all duration-200 cursor-pointer;
  position: relative;
}

.bar-segment:hover {
  @apply opacity-80;
  filter: brightness(1.1);
}

/* 반응형 */
@media (max-width: 640px) {
  .status-grid {
    @apply grid-cols-2;
  }
}
</style>