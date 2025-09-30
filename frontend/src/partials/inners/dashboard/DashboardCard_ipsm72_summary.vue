<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 -->
    <div class="card-header">
      <header class="header-content">
        <h2 class="card-title">실시간 모듈 상태</h2>
        <div class="channel-info">
          <span class="channel-text">
            보조 채널
          </span>
        </div>
      </header>
    </div>

    <!-- 데이터 섹션 -->
    <div class="data-section">
      <!-- 전체 상태 요약 -->
      <div class="summary-section">
        <h3 class="subsection-title">상태 요약</h3>
        
        <div class="status-grid">
          <div class="status-card" v-for="status in statusSummary" :key="status.type">
            <div class="status-left">
              <div class="status-indicator" :class="status.colorClass"></div>
              <div class="status-label">{{ status.label }}</div>
            </div>
            <div class="status-count">{{ status.count }}</div>
          </div>
        </div>
      </div>

      <!-- 2열 레이아웃: 채널 분포 & 최근 이벤트 -->
      <div class="bottom-grid">
        <!-- 채널별 상태 분포 -->
        <div class="distribution-section">
          <h3 class="subsection-title">모듈별 분포</h3>
          
          <div class="channel-stats">
            <!-- 메인 채널 -->
            <div class="channel-item">
              <div class="channel-header">
                <span class="channel-name">IPSM #1</span>
                <span class="channel-total">72 개</span>
              </div>
              <div class="channel-bar">
                <div 
                  v-for="(segment, index) in mainChannelSegments" 
                  :key="index"
                  class="bar-segment"
                  :class="segment.colorClass"
                  :style="{ width: segment.percentage + '%' }"
                  :title="`${segment.label}: ${segment.count}개`"
                >
                </div>
              </div>
            </div>

            <!-- 서브 채널 -->
            <div class="channel-item">
              <div class="channel-header">
                <span class="channel-name">IPSM #2</span>
                <span class="channel-total">72 개</span>
              </div>
              <div class="channel-bar">
                <div 
                  v-for="(segment, index) in subChannelSegments" 
                  :key="index"
                  class="bar-segment"
                  :class="segment.colorClass"
                  :style="{ width: segment.percentage + '%' }"
                  :title="`${segment.label}: ${segment.count}개`"
                >
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 최근 이벤트 -->
        <div class="events-section">
          <h3 class="subsection-title">최근 이벤트</h3>
          
          <div class="event-list">
            <div v-for="event in recentEvents" :key="event.id" class="event-item">
              <div class="event-indicator" :class="getStatusColorClass(event.status)"></div>
              <div class="event-details">
                <div class="event-channel">{{ event.channel }}</div>
                <div class="event-status">{{ event.label }}</div>
              </div>
              <div class="event-time">{{ event.time }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ChannelStatusSummary',
  props: {
    channel: String,
    data: Object,
  },
  setup(props) {
    const { t } = useI18n()
    
    // 더미 데이터 - 상태별 카운트
    const statusData = ref({
      OFF: 15,
      ON: 42,
      TRIP: 3,
      OPR: 8,
      OC: 2,
      OCG: 1,
      'OCG+': 0,
      BAD: 1
    })

    // 상태 정의
    const statusDefinitions = {
      OFF: { label: 'OFF', colorClass: 'bg-gray-200 dark:bg-gray-600' },
      ON: { label: '정상', colorClass: 'bg-green-500 dark:bg-green-400' },
      TRIP: { label: '차단', colorClass: 'bg-red-500 dark:bg-red-400' },
      OPR: { label: '운전', colorClass: 'bg-blue-500 dark:bg-blue-400' },
      OC: { label: '과전류', colorClass: 'bg-orange-500 dark:bg-orange-400' },
      OCG: { label: '지락', colorClass: 'bg-purple-500 dark:bg-purple-400' },
      'OCG+': { label: '지락+', colorClass: 'bg-amber-700 dark:bg-amber-600' },
      BAD: { label: '불량', colorClass: 'bg-gray-800 dark:bg-gray-700' }
    }

    // 상태 요약 계산
    const statusSummary = computed(() => {
      return Object.entries(statusData.value).map(([type, count]) => ({
        type,
        count,
        ...statusDefinitions[type]
      })).filter(status => status.count > 0)
    })

    // 메인 채널 분포 (더미 데이터)
    const mainChannelData = ref({
      OFF: 8,
      ON: 25,
      TRIP: 1,
      OPR: 4,
      OC: 1,
      OCG: 0,
      'OCG+': 0,
      BAD: 1
    })

    const mainChannelTotal = computed(() => {
      return Object.values(mainChannelData.value).reduce((sum, val) => sum + val, 0)
    })

    const mainChannelSegments = computed(() => {
      const total = mainChannelTotal.value
      return Object.entries(mainChannelData.value)
        .filter(([_, count]) => count > 0)
        .map(([type, count]) => ({
          ...statusDefinitions[type],
          count,
          percentage: (count / total * 100).toFixed(1)
        }))
    })

    // 서브 채널 분포 (더미 데이터)
    const subChannelData = ref({
      OFF: 7,
      ON: 17,
      TRIP: 2,
      OPR: 4,
      OC: 1,
      OCG: 1,
      'OCG+': 0,
      BAD: 0
    })

    const subChannelTotal = computed(() => {
      return Object.values(subChannelData.value).reduce((sum, val) => sum + val, 0)
    })

    const subChannelSegments = computed(() => {
      const total = subChannelTotal.value
      return Object.entries(subChannelData.value)
        .filter(([_, count]) => count > 0)
        .map(([type, count]) => ({
          ...statusDefinitions[type],
          count,
          percentage: (count / total * 100).toFixed(1)
        }))
    })

    // 최근 이벤트 (더미 데이터)
    const recentEvents = ref([
      { id: 1, channel: 'CT1-3', status: 'TRIP', label: '차단 발생', time: '2분 전' },
      { id: 2, channel: 'CT15', status: 'OC', label: '과전류 감지', time: '15분 전' },
      { id: 3, channel: 'CT27', status: 'ON', label: '정상 복귀', time: '1시간 전' },
      { id: 4, channel: 'CT42', status: 'OCG', label: '지락 발생', time: '2시간 전' }
    ])

    const getStatusColorClass = (status) => {
      return statusDefinitions[status]?.colorClass || 'bg-gray-400'
    }

    return {
      t,
      statusSummary,
      mainChannelTotal,
      mainChannelSegments,
      subChannelTotal,
      subChannelSegments,
      recentEvents,
      getStatusColorClass
    }
  }
}
</script>

<style scoped>
/* 기존 카드 스타일 유지 */
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-6;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

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

/* 데이터 섹션 */
.data-section {
  @apply p-4 space-y-4;
}

.subsection-title {
  @apply text-sm font-semibold text-gray-700 dark:text-white mb-3;
  @apply flex items-center gap-2;
}

.subsection-title::before {
  content: '';
  @apply w-2 h-2 bg-blue-500 rounded-full;
}

/* 상태 요약 그리드 */
.summary-section {
  @apply mb-4;
}

.status-grid {
  @apply grid grid-cols-4 gap-2;
}

.status-card {
  @apply flex items-center justify-between px-2 py-1.5 bg-gray-50 dark:bg-gray-800/50 rounded-lg;
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

/* 2열 그리드 레이아웃 */
.bottom-grid {
  @apply grid grid-cols-2 gap-4;
}

.distribution-section,
.events-section {
  @apply flex flex-col;
}

/* 채널별 분포 */
.channel-stats {
  @apply space-y-3;
}

.channel-item {
  @apply space-y-2;
}

.channel-header {
  @apply flex justify-between items-center text-sm;
}

.channel-name {
  @apply font-medium text-gray-700 dark:text-gray-300;
}

.channel-total {
  @apply text-gray-500 dark:text-gray-400;
}

.channel-bar {
  @apply flex w-full h-5 rounded-lg overflow-hidden bg-gray-200 dark:bg-gray-700;
}

.bar-segment {
  @apply h-full transition-all duration-500;
  @apply first:rounded-l-lg last:rounded-r-lg;
}

/* 최근 이벤트 */
.event-list {
  @apply space-y-2 max-h-32 overflow-y-auto;
}

.event-item {
  @apply flex items-center gap-2 p-2 bg-gray-50 dark:bg-gray-800/50 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
}

.event-indicator {
  @apply w-2 h-2 rounded-full flex-shrink-0;
}

.event-details {
  @apply flex-1 min-w-0;
}

.event-channel {
  @apply text-xs font-medium text-gray-800 dark:text-white;
}

.event-status {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.event-time {
  @apply text-xs text-gray-400 dark:text-gray-500 flex-shrink-0;
}

/* 반응형 */
@media (max-width: 640px) {
  .status-grid {
    @apply grid-cols-2;
  }
  
  .bottom-grid {
    @apply grid-cols-1;
  }
}
</style>