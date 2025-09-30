<template>
    <div class="premium-dashboard-card">
     <!-- 헤더 -->
     <div class="card-header">
       <header class="header-content">
         <h2 class="card-title">{{ t('dashboard.diagnosis.title') }}</h2>
         <div class="channel-info">
           <span class="channel-text">
             {{ channel == 'main' ? t('dashboard.diagnosis.subtitle_main') : t('dashboard.diagnosis.subtitle_sub') }}
           </span>
         </div>
       </header>
     </div>

     <div class="content-section">
       <div class="grid grid-cols-12 gap-6">
         <!-- 상태 아이템 -->
         <StatusItem v-if="DiagEnable" :channel="channel" :data="stData" mode="diagnosis" />
         <StatusItem v-if="DiagEnable" :channel="channel" :data="pqData" mode="pq" />
       </div>
       <br/>
     </div>
   </div>
 </template>
 
 <script>
 import { ref, computed,watchEffect, onMounted, onUnmounted, watch, inject } from 'vue'
 //import Notification from '../../components/Notification.vue'
 import StatusItem from './StatusItem_Trans_Claude.vue'
 import StatusItem2 from './StatusItem2.vue'
  //import StatusItem_Trans from './StatusItem_Trans.vue'
  //import AlarmItem from '../job/AlarmItem.vue'
 //import DashboardCard_NotFound from './DashboardCard_NotFound.vue'
 import { useSetupStore } from '@/store/setup'; // ✅ Pinia Store 사용
 import axios from 'axios'
 import { useI18n } from 'vue-i18n'  // ✅ 추가

 
 export default {
   name: 'DashboardCard04',
   props: {
     channel: String
   },
   components: {
     //Notification,
     StatusItem,
     StatusItem2,
     //StatusItem_Trans,
     //AlarmItem,
     //DashboardCard_NotFound,
   },
   setup(props) {
     const { t } = useI18n();
     const channel = ref(props.channel);
     const stData = ref({
       devName:'',
       devType:'',
       devStatus: -2
     });
     const pqData = ref({
       devName:'',
       devStatus: -2
     });
     const DiagEnable = ref(false);
     const setupStore = useSetupStore();
     const channelStatus = computed(() => setupStore.getChannelSetting);
     const asset = computed(() => setupStore.getAssetConfig);
     const assetTypes = ref('');
     const status = ref('Normal');
     const alarmContext = ref('');
     const data = ref([]);
     const meterDictMain = inject('meterDictMain');
     const meterDictSub = inject('meterDictSub');
     //const alarmEnable = ref(false);
     let updateInterval = null;

     const transData = ref({});
       const LoadRate = ref(0);

    const LoadFactor = computed(()=>{
      let kva = -1;
      if(channel.value.toLowerCase() == 'main' && stData.value.devType=='Transformer'){
        kva = setupStore.getMkva;
      }
      if(channel.value.toLowerCase() == 'sub' && stData.value.devType=='Transformer'){
        kva = setupStore.getSkva;
      }
      return kva;
    })

    watch(
      () => LoadFactor.value, // ✅ 이렇게 해야 실제 값이 감지됨
      (newVal) => {
        if (newVal > 0 && transData.value?.Stotal) {
          LoadRate.value = ((transData.value.Stotal / newVal) * 100).toFixed(2);
        }
      },
      { immediate: true }
    );
 
     const fetchData = async () => {
          if (!asset.value || (!asset.value.assetName_main && !asset.value.assetName_sub)) {
            //console.log("⏳ asset 준비 안됨. fetchData 대기중");
            //console.log(asset.value);
            return;
          }
         const chName = channel.value == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
         const chType = channel.value == 'main'? asset.value.assetType_main : asset.value.assetType_sub;
         if(chName != ''){

           try {
             const response = await axios.get(`/api/getStatus/${chName}`);
              //console.log(response.data.status);
             if (response.data.status >= 0) {
                stData.value.devName = chName;
                stData.value.devType = chType;
                stData.value.devStatus = response.data.status;
                if(assetTypes.value == 'Transformer'){                 
                    if (channel.value === 'main') {
                      transData.value = { Temp: meterDictMain.value.Temp, Ig: meterDictMain.value.Ig, Stotal:meterDictMain.value.S4 }
                    } else {
                      transData.value = { Temp: meterDictSub.value.Temp, Ig: meterDictSub.value.Ig, Stotal:meterDictSub.value.S4 }
                    }
                }
             }else{
               console.log('No Data');
             }
           }catch (error) {
             console.log("데이터 가져오기 실패:", error);
           } 
         }
     };

     
     const fetchRealData = async () => {
          if (!asset.value || (!asset.value.assetName_main && !asset.value.assetName_sub)) {
            //console.log("⏳ asset 준비 안됨. fetchData 대기중");
            //console.log(asset.value);
            return;
          }
         const chName = channel.value == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
         const chType = channel.value == 'main'? asset.value.assetType_main : asset.value.assetType_sub;
         if(chName != ''){

           try {
             const response = await axios.get(`/api/getRealTime/${chType}/${chName}`);

             if (response.data.success) {
                transData.value = response.data.data;
             }else{
               console.log('No Data');
             }
           }catch (error) {
             console.log("데이터 가져오기 실패:", error);
           } 
         }
     };

    const fetchPQData = async () => {
          if (!asset.value || (!asset.value.assetName_main && !asset.value.assetName_sub)) {
            //console.log("⏳ asset 준비 안됨. fetchData 대기중");
            //console.log(asset.value);
            return;
          }
         const chName = channel.value == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
         const chType = channel.value == 'main'? asset.value.assetType_main : asset.value.assetType_sub;
         if(chName != ''){
           try {
             const response = await axios.get(`/api/getPQStatus/${chName}`);
             if (response.data.status >= 0) {
                pqData.value.devName = response.data.item;
                pqData.value.devStatus = response.data.status;
             }else{
               console.log('No Data');
             }
           }catch (error) {
             console.log("데이터 가져오기 실패:", error);
           } 
         }
     };

     
     const fetchAlarmData = async () => {
      try {
        const chName = channel.value == 'main'?'Main':'Sub';
        const response = await axios.get(`/api/getAlarmLast/${chName}`);
        if (response.data.success) {
            data.value = response.data.data;
            if (data.value.length > 0){
              status.value = 'Alarm';
              alarmContext.value = `${response.data.last.AlarmChannel}, ${response.data.last.TimeStamp}`;
              alarmEnable.value = true;
            }
            else{
              status.value = 'Normal'
              alarmContext.value = 'No Alarm'
              alarmEnable.value = true;
            }
        }else{
          alarmEnable.value = false;
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
        alarmEnable.value = false;
      }
    };

    watch(asset, (newVal) => {
      if (newVal) {
        if(channel.value == 'main')
          assetTypes.value = newVal.assetType_main;
        else
          assetTypes.value = newVal.assetType_sub;
        fetchData();
        if(assetTypes.value != 'Transformer'){
          fetchRealData();
        }
        fetchPQData();
        //fetchAlarmData();

        // 그리고 나서 주기 업데이트 걸어
        if (updateInterval) {
          clearInterval(updateInterval);  // ✅ 혹시 이전에 걸린 거 있으면 지우고
        }
        if (!updateInterval) {
          updateInterval = setInterval(async () => {
            await fetchData();
             if(assetTypes.value != 'Transformer'){
                await fetchRealData();
              }
            await fetchPQData();
            //await fetchAlarmData();
          }, 300000);  //5 x 60 × 1000 : 5min
        }
      }
    }, { immediate: true }); // <-- 바로 실행 시도
 
    onMounted(async () => {
        await setupStore.checkSetting();   // ✅ setupStore에서 서버 데이터 다시 가져오기
      });

 
     onUnmounted(() => {
       if (updateInterval) {
         clearInterval(updateInterval);
       }
     });
 
     watchEffect(() => {
       if(channel.value == 'main')
         DiagEnable.value = channelStatus.value.MainDiagnosis
       else
         DiagEnable.value = channelStatus.value.SubDiagnosis
     });
 
     return {
       channel,
       stData,
       channelStatus,
       DiagEnable,
       fetchData,
       asset,
       status,
       data,
       pqData,
       alarmContext,
       assetTypes,
       transData,
       fetchRealData,
       LoadFactor,
       t,
       LoadRate,
     }    
   }
 }
 </script>

<style scoped>
.premium-dashboard-card {
  @apply flex flex-col col-span-full sm:col-span-6 xl:col-span-6;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
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

.content-section {
  @apply px-5 pt-5;
}
</style>