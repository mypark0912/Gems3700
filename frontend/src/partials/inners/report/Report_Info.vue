<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
    <div class="card">
      <div class="premium-card-header">
        <div class="header-content">
          <div class="header-left">
            <h2 class="card-title">
              {{ t('report.cardTitle.Channel') }}
            </h2>    
          </div>  
        </div>
      </div> 
    </div>   
    <div class="flex flex-col gap-4 p-4 ml-2">
  <!-- 첫 번째 줄: 장비 정보 -->
  <div class="flex gap-6 items-center">
    <img :src="equipImageSrc" alt="장비 이미지"
         class="w-14 h-14 object-cover rounded-lg shadow-md border border-gray-300 dark:border-gray-600" />
    
    <div class="min-w-[120px] flex flex-col space-y-2">
      <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
        {{ channel == 'Main' ? AssetInfo.assetType_main : AssetInfo.assetType_sub }}
      </span>
      <span class="text-lg font-bold text-gray-800 dark:text-gray-100">
        {{ channel == 'Main' ? AssetInfo.assetNickname_main : AssetInfo.assetNickname_sub }}
      </span>
    </div>
    
    <div v-if="devLocation != ''" class="min-w-[120px] flex flex-col space-y-2">
      <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
        설치장소
      </span>
      <span class="text-lg font-bold text-gray-800 dark:text-gray-100">
        {{ devLocation }}
      </span>
    </div>
  </div>

  <!-- 두 번째 줄: 측정 데이터들 -->
  <div class="flex gap-6 overflow-x-auto">
    <div
      v-for="item in rawdata"
      :key="item.Name"
      class="min-w-[120px] flex flex-col space-y-2 flex-shrink-0"
    >
      <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase whitespace-nowrap">
        {{ t(`dashboard.transDiag.${item.Name}`) }}
      </span>
      <span class="text-lg font-bold text-gray-800 dark:text-gray-100">
        {{ item.Value }} {{ item.Unit }}
      </span>
    </div>
  </div>
</div>
  </div>
</template>
  

<script>
import { ref, onMounted, watch, computed } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'  // ✅ 추가
import { useSetupStore } from '@/store/setup'
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
//import { useReportData } from '@/composables/reportDict'

export default {
  name: 'Diagnosis_Info',
  props: {
    channel:{
      type:String,
      default: ''
    }
  },
  setup(props){
    const { t } = useI18n();
    const channel = ref(props.channel);
    const setupStore = useSetupStore();
    const route = useRoute()
    const rawdata = ref([]);
    //const { loadInfoData } = useReportData()

    const AssetInfo = computed(()=>{
      return setupStore.getAssetConfig;
    });

    const devLocation = computed(()=>{
      return setupStore.getDevLocation;
    });

    const equipImageSrc = computed(() => {
      const eqType = channel.value.toLowerCase() == 'main'? AssetInfo.value.assetType_main : AssetInfo.value.assetType_sub;
      switch (eqType) {
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


    const fetchAsset = async () => {
      if(!AssetInfo.value)
        await setupStore.checkSetting(); 
      const chName = channel.value.toLowerCase() == 'main'? AssetInfo.value.assetName_main : AssetInfo.value.assetName_sub;

      if(chName != ''){
        try {

            const response = await axios.get(`/api/getAsset/${chName}`);
            if (response.data.success) {
              rawdata.value = response.data.data;
            }else{
              console.log('No Data');
            }
        }catch (error) {
            console.log("데이터 가져오기 실패:", error);
        }
      }else{
        alert('There are no registered Asset.');
      }
    };

    onMounted(()=>{
      fetchAsset();
    });

    watch(() => route.params.channel, async (newChannel) => {
        channel.value = newChannel
        //console.log('Updated Channel:', channel.value)
        await fetchAsset();
      });

    return {
      channel,
      rawdata,
      t,
      AssetInfo,
      equipImageSrc,
      fetchAsset,
      devLocation,
    }
  }
}
</script>
<style>

@import '../../../css/card-styles.css';
</style>