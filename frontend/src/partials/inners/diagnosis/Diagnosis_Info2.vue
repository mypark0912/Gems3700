<template>
    <div class="flex flex-col col-span-full sm:col-span-6 xl:col-span-7 bg-white dark:bg-gray-800">
      <div class="px-5 py-4"> <!-- ✅ 위아래 여백 조정 -->
                <div class="flex items-center gap-5"> <!-- ✅ 이미지와 텍스트 간 간격 조정 -->
          <!-- 장비 이미지 -->
          <img :src="equipImageSrc" alt="장비 이미지"
            class="w-20 h-20 object-cover rounded-lg shadow-md border border-gray-300 dark:border-gray-600" />
  
          <!-- 장비 이름 & 타입 (수직 정렬) -->
          <div class="flex flex-col space-y-2"> <!-- ✅ 내부 간격 조정 -->
            <!-- 이름 -->
            <div class="flex flex-col">
              <div class="flex flex-row space-x-6">
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-gray-500 dark:text-white uppercase">{{ t('diagnosis.info.name') }}</span>
                  <span class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ thisAsset.name }}</span>
                </div>
                <div class="flex flex-col">
                  <span class="text-xs font-bold text-gray-500 dark:text-white uppercase">{{ t('diagnosis.info.type') }}</span>
                  <span class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ thisAsset.type }}</span>
                </div>
              </div>
            </div> 
            <!-- 타입 -->
            <div class="flex flex-col">
              <div class="flex items-center gap-5 py-3">
                <div v-for="item in AssetInfo" :key="item.Title" class="flex flex-col space-y-3">
                  <div class="flex flex-col">
                    <span class="text-xs font-bold text-gray-500 dark:text-white uppercase">{{ t(`diagnosis.info.${item.Title}`) }}</span>
                    <span class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ item.Value }} {{ item.Unit }}</span>
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
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
import { useI18n } from 'vue-i18n'  // ✅ 추가
export default {
  name: 'Diagnosis_Info2',
  props: {
    asset:{
      type: Object,
      default: () => {}
    },
    data:{
      type: Object,
      default: () => []
    },
    channel:{
      type: String,
      default: ''
    }
  },
  setup(props){
    const { t } = useI18n();
    const asset = ref(props.asset);
    const channel = ref(props.channel);
    const thisAsset = computed(()=> {
      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      const chType = channel.value == 'Main'? asset.value.assetType_main : asset.value.assetType_sub;
      return ({"name":chName, "type":chType});
    });
    const AssetInfo = ref(props.data);
    const equipImageSrc = computed(() => {
      switch (thisAsset.value.type) {
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

    return {
      channel,
      asset,
      thisAsset,
      AssetInfo,
      equipImageSrc,
      t,
    }
  }
}
</script>