<template>
  <div v-show="open" role="alert">
    <div class="inline-flex flex-col w-full max-w-lg px-4 py-2 rounded-lg text-sm bg-white shadow-sm border border-gray-200 text-gray-600">
      <div class="flex w-full justify-between items-start">
        <div class="flex">
          <!-- ✅ `notiType.value` 대신 `notiType` 사용 (Vue가 감지 가능하도록) -->
          <svg v-if="notiType === 'warning'" class="shrink-0 fill-current text-yellow-500 mt-[3px] mr-3" width="16" height="16" viewBox="0 0 16 16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm0 12c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1zm1-3H7V4h2v5z" />
          </svg>
          <svg v-else-if="notiType === 'error'" class="shrink-0 fill-current text-red-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm3.5 10.1l-1.4 1.4L8 9.4l-2.1 2.1-1.4-1.4L6.6 8 4.5 5.9l1.4-1.4L8 6.6l2.1-2.1 1.4 1.4L9.4 8l2.1 2.1z" />
          </svg>
          <svg v-else-if="notiType === 'success'" class="shrink-0 fill-current text-green-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zM7 11.4L3.6 8 5 6.6l2 2 4-4L12.4 6 7 11.4z" />
          </svg>
          <svg v-else-if="notiType === 'inspect'" class="shrink-0 fill-current text-orange-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm1 12H7V7h2v5zM8 6c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1z" />
          </svg>
          <svg v-else class="shrink-0 fill-current text-gray-500 mt-[3px] mr-3" width="16" height="16">
            <path d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm1 12H7V7h2v5zM8 6c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1z" />
          </svg>

          <div>
            <div class="font-medium text-gray-800 mb-1">{{ stData.devName }}</div>
            <div class="flex items-center mt-2">
              <img :src="motorImageSrc" alt="장비 이미지" class="w-20 h-20 object-cover rounded-lg mr-3" />
              <div class="flex flex-col">
                <span class="text-base">장비 종류 : {{ stData.devType }}</span>
                <span class="text-base mt-2">장비 상태 : 
                  <div class="inline-block text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center" :class="statusClass">
                    {{ statusStr }}
                  </div>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>  
</template>

<script>
import { ref, watch, computed, onMounted } from 'vue';
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
export default {
  name: 'Notification',
  props: {
    data: Object,
    open: Boolean,
  },
  setup(props) {
    const stData = ref(props.data);

    const notiType = ref('none');

    const updateNotiType = (status) => {
      switch (status) {
        case 1:
          notiType.value = 'success';
          break;
        case 2:
          notiType.value = 'warning';
          break;
        case 3:
          notiType.value = 'inspect';
          break;
        case 4:
          notiType.value = 'error';
          break;
        default:
          notiType.value = 'none';
      }
    };

    const statusClass = computed(() => {
      switch (stData.value.devStatus) {
        case 1:
          return 'bg-green-500/20 text-green-700 font-semibold';
        case 2:
        return 'bg-yellow-500/20 text-yellow-700 font-semibold';
        case 3:
          return 'bg-orange-500/20 text-orange-700 font-semibold';
        case 4:
          return 'bg-red-500/20 text-red-700 font-semibold';
        default:
          return 'bg-gray-500/20 text-gray-700 font-semibold';
      }
    });

    const statusStr =  computed(() => {
      switch (stData.value.devStatus) {
        case 1:
          return 'OK';
        case 2:
        return 'Warning';
        case 3:
          return 'Inspect';
        case 4:
          return 'Repair';
        default:
          return 'No Data';
      }
    });

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

    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          console.log("🔄 새 데이터 업데이트:", newData);
          stData.value = { ...newData };  
          updateNotiType(stData.value.devStatus);
          console.log("✅ 업데이트된 notiType:", notiType.value);
        }
      },
      { immediate: true, deep: true }
    );

    return {
      stData,
      notiType,
      statusClass,
      motorImageSrc,
      statusStr,
    };
  },
};
</script>
