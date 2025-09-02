<template>
  <div v-show="open" role="alert">
    <div class="inline-flex flex-col w-full max-w-lg px-4 py-2 rounded-lg text-sm bg-white dark:bg-gray-800 shadow-sm border border-gray-200 dark:border-gray-700/60 text-gray-600 dark:text-gray-400">
      <div class="flex w-full justify-between items-start">
        <div class="flex items-center">
          <!-- 상태 아이콘 -->
          <svg v-if="statusType === 'warning'" class="shrink-0 fill-current text-yellow-500 mt-[3px] mr-3" width="16" height="16" viewBox="0 0 16 16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-8zm0 12c-.6 0-1-.4-1-1s.4-1 1-1-.4 1-1 1zm1-3H7V4h2v5z" />
          </svg>
          <svg v-else-if="statusType === 'repair'" class="shrink-0 fill-current text-red-500 mt-[3px] mr-3" width="16" height="16" viewBox="0 0 16 16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-8zm3.5 10.1l-1.4 1.4L8 9.4l-2.1 2.1-1.4-1.4L6.6 8 4.5 5.9l1.4-1.4L8 6.6l2.1-2.1 1.4 1.4L9.4 8l2.1 2.1z" />
          </svg>
          <svg v-else-if="statusType === 'ok'" class="shrink-0 fill-current text-green-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-8zm-1 11.4L3.6 8 5 6.6l2 2 4-4L12.4 6 7 11.4z" />
          </svg>
          <svg v-else class="shrink-0 fill-current text-gray-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-8zm1 12H7V7h2v5zM8 6c-.6 0-1-.4-1-1s.4-1 1-1 .4 1 1 1-.4 1-1 1z" />
          </svg>

          <!-- 장비 이름 -->
          <span class="font-semibold ml-2">{{ equipmentName }}</span>
        </div>
      </div>

      <!-- 내부 슬롯 영역 (모터 아이콘 + 상태) -->
      <div class="flex items-center mt-2">
        <!-- 장비 이미지 -->
        <img src="../images/motor_pump.png" alt="장비 이미지" class="w-10 h-10 object-cover rounded-lg mr-3" />

        <!-- 장비 상태 -->
        <span class="text-sm font-medium">장비 상태 :</span>
        <div class="inline-flex bg-green-500/20 text-green-700 rounded-full px-3 py-1 ml-2">
          {{ equipmentStatus }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Notification',
  props: {
    open: Boolean,
    equipmentName: String, // 장비 이름
    equipmentStatus: String, // 장비 상태
  },
  computed: {
    statusType() {
      // 장비 상태에 따른 알림 타입 지정
      const status = this.equipmentStatus.toLowerCase();
      if (status === 'warning') return 'warning';
      if (status === 'repair') return 'repair';
      if (status === 'ok') return 'ok';
      return 'info'; // 기본값
    },
  },
};
</script>