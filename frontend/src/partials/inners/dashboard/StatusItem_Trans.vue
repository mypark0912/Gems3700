<template>
    <div class="detail-card">
      <!-- 헤더 섹션 -->
      <div class="card-header">
        <h3 class="card-title">
          {{ mode === 'pq' ? t('dashboard.diagnosis.pq') : t('dashboard.diagnosis.diagnostic') }}
        </h3>
      </div>
  
      <!-- 상태 인디케이터 섹션 -->
      <div class="status-section">
        <div class="status-indicators">
          <template v-for="(item, idx) in displayedStatuses" :key="idx">
            <div 
              class="status-item"
              :class="{ 'status-active': stData.devStatus === idx + 1 }"
            >
              <!-- 상태 아이콘 -->
              <div class="status-icon-wrapper">
                <div 
                  class="status-icon"
                  :class="[
                    item.fill,
                    stData.devStatus === idx + 1 ? 'status-icon-active' : 'status-icon-inactive'
                  ]"
                >
                  <!-- 활성 상태일 때 펄스 효과 -->
                  <div 
                    v-if="(stData.devStatus === idx + 1) && stData.devStatus > 1" 
                    class="status-pulse"
                    :class="item.fill"
                  ></div>
                </div>
              </div>
              
              <!-- 상태 텍스트 -->
              <span 
                class="status-text"
                :class="stData.devStatus === idx + 1 ? 'status-text-active' : 'status-text-inactive'"
              >
                {{ item.text }}
              </span>
            </div>
          </template>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { useI18n } from 'vue-i18n'
  
  const props = defineProps({
    data: Object,
    channel: String,
    mode: String,
  })
  
  const { t } = useI18n()
  const stData = ref(props.data)
  const currentTime = ref(new Date())
  
  // 시간 업데이트
  let timeInterval = null
  onMounted(() => {
    timeInterval = setInterval(() => {
      currentTime.value = new Date()
    }, 1000)
  })
  
  onUnmounted(() => {
    if (timeInterval) {
      clearInterval(timeInterval)
    }
  })
  
  // 표시할 상태들 (정지 제외)
  const displayedStatuses = computed(() => {
    const fills = [
      'bg-emerald-500',    // 정상 - 더 선명한 초록
      'bg-amber-400',      // 경고 - 더 부드러운 노랑
      'bg-orange-500',     // 점검 - 더 선명한 주황
      'bg-red-500',        // 수리 - 기존 빨강 유지
    ]
  
    const texts = props.mode === 'diagnosis'
      ? [
          t('dashboard.diagnosis.st1'),
          t('dashboard.diagnosis.st2'),
          t('dashboard.diagnosis.st3'),
          t('dashboard.diagnosis.st4'),
        ]
      : [
          t('dashboard.diagnosis.pqfe1'),
          t('dashboard.diagnosis.pqfe2'),
          t('dashboard.diagnosis.pqfe3'),
          t('dashboard.diagnosis.pqfe4'),
        ]
  
    return texts.map((text, i) => ({
      text,
      fill: fills[i],
    }))
  })
  
  // 현재 상태 텍스트 가져오기
  const getCurrentStatusText = () => {
    if (stData.value.devStatus === 0) {
      return props.mode === 'diagnosis' 
        ? t('dashboard.diagnosis.st0') 
        : t('dashboard.diagnosis.pqfe0')
    }
    
    const statusIndex = stData.value.devStatus - 1
    return displayedStatuses.value[statusIndex]?.text || '알 수 없음'
  }
  
  // 현재 시간 포맷팅
  const getCurrentTime = () => {
    return currentTime.value.toLocaleTimeString('ko-KR', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }
  </script>
  
  <style scoped>
  /* detail-card 스타일과 동일하게 - mb-2 제거하고 flex 추가 */
  .detail-card {
    @apply bg-white dark:bg-gray-800 rounded-lg;
    @apply border border-gray-200 dark:border-gray-700;
    @apply shadow-sm hover:shadow-md transition-all duration-200;
    @apply overflow-hidden;
    @apply flex flex-col; /* flex 컨테이너로 설정 */
    height: 100%; /* 부모 높이에 맞춤 */
  }
  
  /* detail-header 스타일과 동일하게 */
  .card-header {
    @apply p-3 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-600;
    @apply flex justify-between items-center;
    @apply flex-shrink-0; /* 헤더 크기 고정 */
  }
  
  .header-content {
    @apply flex justify-between items-center;
  }
  
  /* detail-title 스타일과 동일하게 */
  .card-title {
    @apply text-sm font-semibold text-gray-900 dark:text-gray-100;
    @apply flex items-center gap-2;
  }
  
  .channel-info {
    @apply flex items-center;
  }
  
  .channel-text {
    @apply text-xs font-semibold text-gray-400 dark:text-white uppercase;
  }
  
  /* 상태 섹션 - flex-grow로 남은 공간 차지하고 중앙 정렬 */
  .status-section {
    @apply px-2 py-4;
    @apply flex-grow flex items-center justify-center; /* 중앙 정렬 */
  }
  
  .status-indicators {
    @apply grid grid-cols-2 md:grid-cols-4 gap-4;
    @apply w-full; /* 전체 너비 사용 */
  }
  
  .status-item {
    @apply flex flex-col items-center justify-center gap-3;
    @apply p-4 rounded-lg transition-all duration-300;
  }
  
  .status-active {
    @apply transform scale-105;
  }
  
  .status-icon-wrapper {
    @apply relative;
  }
  
  .status-icon {
    @apply w-7 h-7 rounded-full transition-all duration-300;
    @apply flex items-center justify-center;
    @apply shadow-md;
  }
  
  .status-icon-active {
    @apply shadow-lg transform scale-110;
  }
  
  .status-icon-inactive {
    @apply bg-gray-300 dark:bg-gray-600;
    @apply opacity-60;
  }
  
  .status-pulse {
    @apply absolute inset-0 rounded-full;
    @apply animate-ping opacity-75;
  }
  
  .status-text {
    @apply text-sm font-medium text-center;
    @apply transition-all duration-200;
    @apply leading-tight;
  }
  
  .status-text-active {
    @apply text-gray-900 dark:text-white font-semibold;
  }
  
  .status-text-inactive {
    @apply text-gray-500 dark:text-white;
  }
  
  /* 반응형 개선 */
  @media (max-width: 1024px) {
    .status-indicators {
      @apply grid-cols-2 gap-3;
    }
    
    .detail-card {
      height: auto; /* 태블릿에서는 자동 높이 */
      min-height: 200px; /* 최소 높이 보장 */
    }
  }
  
  @media (max-width: 640px) {
    .detail-card {
      @apply col-span-full;
      height: auto; /* 모바일에서는 자동 높이 */
      min-height: 180px; /* 모바일 최소 높이 */
    }
    
    .card-header {
      @apply p-2;
    }
    
    .status-section {
      @apply p-3;
    }
    
    .status-indicators {
      @apply grid-cols-2 gap-3;
    }
    
    .status-item {
      @apply p-3 gap-2;
    }
    
    .status-icon {
      @apply w-6 h-6;
    }
    
    .status-text {
      @apply text-xs;
    }
  }
  
  @media (prefers-color-scheme: dark) {
    .detail-card {
      @apply bg-gray-800;
    }
  }
  </style>