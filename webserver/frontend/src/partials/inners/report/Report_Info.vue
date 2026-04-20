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
  <div class="flex gap-6 items-start">  <!-- items-center -> items-start -->

  <div class="min-w-[120px] flex flex-col space-y-2">
    <span class="text-xs font-bold text-gray-500 dark:text-white uppercase">
      {{ t('diagnosis.info.drivetype') }}
    </span>
    <span class="text-lg font-bold text-gray-800 dark:text-gray-100">
      {{ drType == 'DOL'? t('diagnosis.info.dr1') : t('diagnosis.info.dr2') }}
    </span>
  </div>

  <div v-if="devLocation != ''" class="min-w-[120px] flex flex-col space-y-2">
    <span class="text-xs font-bold text-gray-500 dark:text-gray-100 uppercase whitespace-nowrap">
      {{ t('report.cardTitle.installation') }}
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
      <span class="text-xs font-bold text-gray-500 dark:text-gray-100 uppercase whitespace-nowrap">
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
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'  // ✅ 추가
import { useSetupStore } from '@/store/setup'
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
    const rawdata = ref([]);
    const drType = ref('');
    //const { loadInfoData } = useReportData()

    const devLocation = computed(()=>{
      return setupStore.getDevLocation;
    });



    const fetchAsset = async () => {
      await setupStore.checkSetting();
      const chName = '';

      if(chName != ''){
        try {

            const response = await axios.get(`/api/getAsset/${chName}`);
            if (response.data.success) {
              rawdata.value = response.data.data;
              drType.value = response.data.driveType;
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

    return {
      channel,
      rawdata,
      t,
      fetchAsset,
      devLocation,
      drType,
    }
  }
}
</script>
<style>

@import '../../../css/card-styles.css';
</style>