<template>
    <div class="flex flex-col col-span-full sm:col-span-6 xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
      <div class="px-5 pt-5">
        <header class="flex justify-between items-start mb-2">
          <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ channel }} channel</h2>
        </header>
        <div class="grid grid-cols-12 gap-4 py-3">
          <div
              v-for="(key, index) in Object.keys(data)"
              :key="key"
              class="col-span-12 sm:col-span-2 xl:col-span-2"
            >
            <SubDash v-if="Object.keys(data[key]).length > 0" :id="key" :data="JSON.parse(data[key])" :channel="channel" />
          </div>
        </div>
        </div>
      </div>
</template>
 
 <script>
 import { ref, computed,watchEffect, onMounted, onUnmounted } from 'vue'
 //import Notification from '../../components/Notification2.vue'
import SubDash from '../job/SubDash.vue'
 import { useAuthStore } from '@/store'; // ✅ Pinia Store 사용
 import axios from 'axios'
 
 
 export default {
   name: 'DashboardCard04',
   props: {
     channel: String,
     data:Object,
   },
   components: {
    SubDash,
   },
   setup(props) {
     const channel = ref(props.channel);
     const data = ref(props.data);
    //  const chartdata =ref(null);
     // const stData = ref({
     //   devName:'',
     //   devType:'',
     //   devStatus: -2
     // });
    //  const DiagEnable = ref(false);
    //  const authStore = useAuthStore();
    //  const channelStatus = computed(() => authStore.getChannelSetting);
    //  const asset = computed(() => authStore.getAssetConfig);
    //  let updateInterval = null;
 
    //  const fetchData = async () => {
    //      const chName = channel.value == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
    //      if(chName != ''){
    //        try {
    //          const response = await axios.get(`/api/getDiagnosis/${chName}`);
    //          if (response.data.success) {
    //            let itemlist = [], valuelist=[];
    //            for(let i = 0; i < response.data.data.length;i++){
    //              itemlist.push(response.data.data[i]["Name"]);
    //              valuelist.push(response.data.data[i]["Status"]);
    //            }
    //            chartdata.value = {"Names" : itemlist, "Values" : valuelist};
    //          }else{
    //            console.log('No Data');
    //          }
    //        }catch (error) {
    //          console.log("데이터 가져오기 실패:", error);
    //        } 
    //      }
    //  };
 
 
    //  onUnmounted(() => {
    //    if (updateInterval) {
    //      clearInterval(updateInterval);
    //    }
    //  });
 
    //  watchEffect(() => {
    //    console.log("🔄 channelStatus 값 변경 감지됨:", channelStatus.value);
    //    if(channel.value == 'main')
    //      DiagEnable.value = channelStatus.value.MainDiagnosis
    //    else
    //      DiagEnable.value = channelStatus.value.SubDiagnosis
    //    console.log(channel.value + " : 🔄 channelStatus 값 변경 감지됨:", DiagEnable.value);
    //  });
 
     return {
       channel,
       data,
     }    
   }
 }
 </script>