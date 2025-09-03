<template>
  <div class="flex h-[100dvh] overflow-hidden">

    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = true" :channel="channel" />

    <!-- Content area -->
    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      
      <!-- Site header -->
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          
          <!-- Dashboard actions -->
          <div class="sm:flex sm:justify-between sm:items-center mb-4">

            <!-- Left: Title -->
            <div class="mb-4 sm:mb-0">
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold ml-1">{{ channel == 'Main' ? t('meter.sitemap.main'):t('meter.sitemap.sub') }} > {{ t('meter.sitemap.title') }}</h2>
            </div>

          </div>

          <!-- Cards -->
          <div class="grid grid-cols-12 gap-4">
            
            <MeterDetail2 v-if="Object.keys(meterDatas).length > 0" :data="meterDatas.meterData" :channel="channel" :title="'Meter'" />
              <div class="col-span-5 md:col-span-5 flex flex-col gap-4">
                <div class="col-span-5">
                  <CanvasAngle2 v-if="Object.keys(phaseDict).length > 0"
                    :degree="phaseDict.degree"
                    :magnitude="phaseDict.magnitude"
                    :texts="phaseDict.texts"
                    :maxlist="phaseDict.max"
                    :channel="channel"
                  />
                </div>
                <div class="flex gap-4">
                <!--div class="basis-1/2">
                  <MeterKwh
                    v-if="Object.keys(meterDatas).length > 0"
                    :data="energyData"
                    :channel="channel"
                    :title="'Energy'"
                    :mode="'import'"
                  />
                </div>
                <div class="basis-1/2">
                  <MeterKwh
                    v-if="Object.keys(energyData).length > 0"
                    :data="energyData"
                    :channel="channel"
                    :title="'Energy'"
                    :mode="'export'"
                  />
                </div-->
              </div>
              </div>

            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.thdData" :data="powerThd.thdData" :channel="channel" :title="'THD'" />
            <!--MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.demandDataP" :data="powerThd.demandDataP" :channel="channel" :title="'Demand'" /-->
            
            <MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.powerData" :data="powerThd.powerData" :channel="channel" :title="'Power'" />
            <!--MeterDetail2 v-if="Object.keys(powerThd).length > 0 && powerThd.demandDataI" :data="powerThd.demandDataI" :channel="channel" :title="'Demand I'" /-->

          </div>

        </div>
      </main>
      <Footer />
    </div> 

  </div>
</template>

<script>
import { reactive,ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

import axios from 'axios'
import Sidebar from '../common/SideBar3.vue'
import Header from '../common/Header.vue'
import Footer from "../common/Footer.vue";

import CanvasAngle2 from '../../partials/inners/meter/CanvasAngle4.vue'

import MeterKwh from '../../partials/inners/meter/MeterKwh.vue'
import MeterDetail2 from '../../partials/inners/meter/MeterDetail2.vue'
import { useSetupStore } from '@/store/setup'
import { useI18n } from 'vue-i18n'

export default {
  name: 'Meter',
  props:['channel'],
  components: {
    Sidebar,
    Header,
    Footer,
    CanvasAngle2,
    MeterKwh,
    MeterDetail2,
  },
  setup(props) {
    const { t } = useI18n();
    const route = useRoute();
    const sidebarOpen = ref(true)
    const channel = ref(props.channel)
    const meterDatas = ref({});
    const powerThd = ref({});
    const phaseDict = ref({});
    const energyData = ref({});

    // const demandData = ref({
    //   demandData: [
    //     {
    //       subTitle: "Active",
    //       data: [
    //         { id: 0, subTitle: "Imp.", value: "0.000", unit: "W", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Imp. Max", value: "0.000", unit: "W", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 2, subTitle: "Exp.", value: "0.000", unit: "W", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 3, subTitle: "Exp. Max", value: "0.000", unit: "W", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     },
    //     {
    //       subTitle: "Reactive",
    //       data: [
    //         { id: 0, subTitle: "Imp.", value: "0.000", unit: "VAR", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Imp. Max", value: "0.000", unit: "VAR", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 2, subTitle: "Exp.", value: "0.000", unit: "VAR", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 3, subTitle: "Exp. Max", value: "0.000", unit: "VAR", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     },
    //     {
    //       subTitle: "Apparent",
    //       data: [
    //         { id: 0, subTitle: "Imp.", value: "0.000", unit: "VA", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Imp. Max", value: "0.000", unit: "VA", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 2, subTitle: "Exp.", value: "0.000", unit: "VA", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 3, subTitle: "Exp. Max", value: "0.000", unit: "VA", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     }
    //   ],
    //   demandIData: [
    //     {
    //       subTitle: "L1",
    //       data: [
    //         { id: 0, subTitle: "", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Max", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     },
    //     {
    //       subTitle: "L2",
    //       data: [
    //         { id: 0, subTitle: "", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Max", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     },
    //     {
    //       subTitle: "L3",
    //       data: [
    //         { id: 0, subTitle: "", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" },
    //         { id: 1, subTitle: "Max", value: "0.000", unit: "A", min: "0.000", max: "0.000", minTime: "1970-01-01 00:00:00", maxTime: "1970-01-01 00:00:00" }
    //       ]
    //     }
    //   ]
    // }); // 새로 추가된 demand 데이터
    const setupStore = useSetupStore();
    let updateInterval_ones = null;
    let updateInterval_onem = null;
    let updateInterval_fifthm = null;
    let updateInterval_oneh = null;
    
    const assetName = computed(()=> {
      const mainName = setupStore.getAssetConfig.assetNickname_main;
      const subName = setupStore.getAssetConfig.assetNickname_sub;
      if(channel.value == 'Main' || channel.value == 'main'){
        if (mainName != ''){
            return  "("+ mainName+")";
        }else{
          return "";
        }
      }else{
        if (subName != ''){
            return  "("+ subName+")";
        }else{
          return "";
        }
      }
    })

    const setup = computed(() => setupStore.getSetup);

    const unbal = computed(()=> setupStore.getUnbalance);
    
    const ReleaseInterval = () =>{
      if (updateInterval_ones) {
        clearInterval(updateInterval_ones);
      }
      if (updateInterval_onem) {
        clearInterval(updateInterval_onem);
      }
      if (updateInterval_fifthm) {
        clearInterval(updateInterval_fifthm);
      }
      if (updateInterval_oneh) {
        clearInterval(updateInterval_oneh);
      }
    };
    
    onUnmounted(() => {
      ReleaseInterval();
    });

    watch(() => props.channel, (newChannel) => {      
      if (newChannel !== channel.value) {
        console.log("Channel changed to:", newChannel);
        channel.value = newChannel;
        meterDatas.value = {};      
        powerThd.value = {}; 
        phaseDict.value = {}; 
        energyData.value ={}; 
        demandData.value = { 
          demandData: [
            { subTitle: "Active", data: [] },
            { subTitle: "Reactive", data: [] },
            { subTitle: "Apparent", data: [] }
          ], 
          demandIData: [
            { subTitle: "L1", data: [] },
            { subTitle: "L2", data: [] },
            { subTitle: "L3", data: [] }
          ] 
        }; // 새로 추가
        ReleaseInterval();
        startFetching();
      }
    },{ immediate: true });

    watch(() => route.params.channel, (newChannel) => {
      if (newChannel !== channel.value) {
        channel.value = newChannel;
        meterDatas.value = {};                
        powerThd.value = {}; 
        phaseDict.value ={}; 
        energyData.value ={}; 




        ReleaseInterval();
        startFetching();
      }
    }, { immediate: true });

    const fetchRedisOnesData = async (ch) => {
      try {
        const response = await axios.get(`/api/getOnesfromRedis/${ch}/${unbal.value}`);
        if (response.data.success) {
          Object.assign(meterDatas.value, response.data.retData.retData);

        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisOnemData = async (ch) => {
      console.log('Channel:', ch);
      try {
        const response = await axios.get(`/api/getonemfromRedis/${ch}`);
        console.log(response.status);
        if (response.data.success) {
          phaseDict.value = response.data.retData.retData.angleData;
          
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisfifthmData = async (ch) => {
      try {
        const response = await axios.get(`/api/getFifthMfromRedis/${ch}`);
        if (response.data.success) {
          powerThd.value = response.data.retData.retData;
          powerThd.value.powerData.forEach(section => {
            section.data.forEach(item => {
              item.value = (item.value / 1000).toFixed(2);
            });
          });
          
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchRedisOnehData = async (ch) => {
      try {
        const response = await axios.get(`/api/getOnehfromRedis/${ch}`);

        if (response.data.success) {
          energyData.value = response.data.retData.retData.energyData;
        }else{
          if ('retData' in response.data) {
            energyData.value = response.data.retData.retData.energyData;
          } 
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    // 새로 추가된 Demand 데이터 가져오기 함수
    // const fetchDemandData = async (ch) => {
    //   try {
    //     // 실제 API 엔드포인트에 맞게 수정 필요
    //     // const response = await axios.get(`/api/getDemandData/${ch}`);
    //     // if (response.data.success && response.data.retData && response.data.retData.retData) {
    //     //   demandData.value = response.data.retData.retData;
    //     // }
        
    //     // 더미 데이터로 설정 (초기값과 동일하게)
    //     demandData.value = {
    //       demandData: [
    //         {
    //           subTitle: "Active",
    //           data: [
    //             { id: 0, subTitle: "Imp.", value: "0.000", unit: "W", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" },
    //             { id: 1, subTitle: "Exp.", value: "0.000", unit: "W", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         },
    //         {
    //           subTitle: "Reactive", 
    //           data: [
    //             { id: 0, subTitle: "Imp.", value: "0.000", unit: "VAR", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" },
    //             { id: 1, subTitle: "Exp.", value: "0.000", unit: "VAR", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         },
    //         {
    //           subTitle: "Apparent",
    //           data: [
    //             { id: 0, subTitle: "Imp.", value: "0.000", unit: "VA", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" },
    //             { id: 1, subTitle: "Exp.", value: "0.000", unit: "VA", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         }
    //       ],
    //       demandIData: [
    //         {
    //           subTitle: "L1",
    //           data: [
    //             { id: 0, subTitle: "", value: "0.000", unit: "A", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         },
    //         {
    //           subTitle: "L2", 
    //           data: [
    //             { id: 0, subTitle: "", value: "0.000", unit: "A", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         },
    //         {
    //           subTitle: "L3",
    //           data: [
    //             { id: 0, subTitle: "", value: "0.000", unit: "A", min: "-", max: "0.000", minTime: "-", maxTime: "1970-01-01 00:00:00" }
    //           ]
    //         }
    //       ]
    //     };
    //   } catch (error) {
    //     console.log("Demand 데이터 가져오기 실패:", error);
    //     // 에러 시에도 동일한 구조 유지
    //     demandData.value = {
    //       demandData: [
    //         { subTitle: "Active", data: [] },
    //         { subTitle: "Reactive", data: [] },
    //         { subTitle: "Apparent", data: [] }
    //       ],
    //       demandIData: [
    //         { subTitle: "L1", data: [] },
    //         { subTitle: "L2", data: [] },
    //         { subTitle: "L3", data: [] }
    //       ]
    //     };
    //   }
    // };

    const startFetching = () => {
      if (updateInterval_ones) clearInterval(updateInterval_ones);
      fetchRedisOnesData(channel.value);
        updateInterval_ones = setInterval(() => {
          fetchRedisOnesData(channel.value);
        }, 1000);

      if (updateInterval_onem) clearInterval(updateInterval_onem);
      fetchRedisOnemData(channel.value);
        updateInterval_onem = setInterval(() => {
          fetchRedisOnemData(channel.value);
        }, 60*1000);

      if (updateInterval_fifthm) clearInterval(updateInterval_fifthm);
      fetchRedisfifthmData(channel.value);
        updateInterval_fifthm = setInterval(() => {
          fetchRedisfifthmData(channel.value);
        }, 15*60*1000);

      // if (updateInterval_oneh) clearInterval(updateInterval_oneh);
      // fetchRedisOnehData(channel.value);
      //   updateInterval_oneh = setInterval(() => {
      //     fetchRedisOnehData(channel.value);
      //   }, 60*60*1000);

      // Demand 데이터도 주기적으로 가져오기
      // fetchDemandData(channel.value);
      // setInterval(() => {
      //   fetchDemandData(channel.value);
      // }, 60*1000); // 1분마다 업데이트
    };

    watch(() => setup.value, (newChannel) => {
      if (newChannel) {
        startFetching();
      }
    },{ immediate: true });

    return {
      sidebarOpen,
      channel,
      meterDatas,
      phaseDict,
      t,
      assetName,
      powerThd,
      energyData,
      //demandData, // 새로 추가
      setup,
      unbal,
    }  
  }
}
</script>