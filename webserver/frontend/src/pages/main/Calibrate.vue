<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        
        <!-- Site header -->
        <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />
  
        <main class="grow flex flex-col">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full flex-1">
          <!-- Page header -->
          <div class="sm:flex sm:justify-between sm:items-center mb-4">
            <!-- Left: Title -->
            <div class="mb-4 sm:mb-0">
              <h2
                class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold mb-3"
              >

              <template
                v-if="
                  ![                   
                    'Maintenance'
                  ].includes(formattedChannel)
                "
              >
                
              </template>
                {{ t('sidebar.config') }}  > 

                
                {{ t(`config.sitemap.${formattedChannel}`) }}
              </h2>
            </div>
            </div>
            <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl mb-8 flex-1 flex flex-col">
              <ServicePanel v-if="channel == 'Service'"/>
                <!--div v-else-if="channel == 'Calibrate'" class="flex flex-col md:flex-row md:-mr-px">
                  <CaliPanel :items="items" :channel="channel"/>
                  <CaliSidebar :channel="channel" :commands="commands"  @startPolling="startPolling" @stopPolling="stopPolling" />
                </div-->
                <div v-else-if="channel == 'Calibrate'" class="flex flex-col md:flex-row md:-mr-px gap-4 flex-1">
                  <div class="md:w-3/4">
                    <!-- 또는 md:w-3/5, md:flex-[2] 등 -->
                    <CaliPanel :items="items" :channel="channel"/>
                  </div>
                  <div class="md:w-1/4 md:flex">
                    <!-- 또는 md:w-2/5, md:flex-[1] 등 -->
                    <CaliSidebar class="flex-1 w-full" :channel="channel" :commands="commands" @startPolling="startPolling" @stopPolling="stopPolling" />
                  </div>
                </div>
                <System v-else-if="channel == 'Command'" />
                <Maintenance v-else-if="channel == 'Maintenance'" />
                <LogView v-else />
            </div>
        </div>
        </main>
        <Footer />
      </div> 
  
    </div>
  </template>
  
  <script>
  import { useI18n } from 'vue-i18n' 
  import { ref , computed, watch,provide, reactive} from 'vue'
  import Sidebar from '../common/SideBar.vue'
  import Header from '../common/Header.vue'
  //import UsersTabsCard from '../../partials/community/CaliCard.vue'
  import CaliSidebar from '../../partials/inners/config/CaliSideBar.vue'
  import CaliPanel from '../../partials/inners/config/CaliPanel.vue'
  import ServicePanel from '../../partials/inners/config/ServicePanel.vue'
  import Maintenance from '../../partials/inners/config/MaintenancePanel.vue'
  import LogView from '../../partials/inners/config/LogPanel.vue' 
  import System from '../../partials/inners/config/CommandPanel.vue'
  //import { useAuthStore } from "@/store"; // ✅ Pinia Store 사용
   import axios from 'axios';
  export default {
    name: 'Calibrate',
    props:['channel'],
    components: {
      Sidebar,
      Header,
      //UsersTabsCard,
      CaliSidebar,
      CaliPanel,
      ServicePanel,
      Maintenance,
      System,
      LogView
    },
    setup(props) {
      //const authStore = useAuthStore();
      const { t } = useI18n();
      //const langset = computed(() => authStore.getLang);
      //const installed = computed(() => authStore.getInstalled);
      const sidebarOpen = ref(false)
      const channel = ref(props.channel);
      let updateInterval = null;

      watch(() => props.channel, (newChannel) => {
      if (newChannel !== channel.value) {
        channel.value = newChannel;
      }
    }, { immediate: true });

      const commands = ref([
      {
        "name":"Start", "Label":"Start"
      },
      {
        "name":"End", "Label":"End"
      },
      {
        "name":"CAL_VGAIN","number":1, "Label":"Set U", "Param":"U", "Type":"SET"
      },
      {
        "name":"CLR_VGAIN","number":7, "Label":"Clear U", "Type":"CLEAR"
      },
      {
        "name":"CAL_VPPGAIN","number":2, "Label":"Set Upp","Param":"U", "Type":"SET"
      },
      {
        "name":"CLR_VPPGAIN","number":8, "Label":"Clear Upp", "Type":"CLEAR"
      },
      {
        "name":"CAL_IGAIN","number":3, "Label":"Set I", "Param":"I", "Type":"SET"
      },
      {
        "name":"CLR_IGAIN","number":9, "Label":"Clear I", "Type":"CLEAR"
      },
      {
        "name":"CAL_INGAIN","number":4, "Label":"Set Ig", "Param":"In", "Type":"SET"
      },
      {
        "name":"CLR_INGAIN","number":12, "Label":"Clear Ig", "Type":"CLEAR"
      },
      {
        "name":"CAL_WGAIN", "number":5,"Label":"Set Power", "Param":"U,I","Type":"SET"
      },
      {
        "name":"CLR_WGAIN", "number":11,"Label":"Clear Power", "Type":"CLEAR"
      },
      {
        "name":"CAL_PHGAIN","number":6, "Label":"Set Phase", "Param":"P", "Type":"SET"
      },
      {
        "name":"CLR_PHGAIN", "number":10,"Label":"Clear Phase", "Type":"CLEAR"
      }
      // {
      //   "name":"CALDCos", "Label":"Set DCos", "Type":"SET"
      // },
      // {
      //   "name":"CLRDCos", "Label":"Clear DCos", "Type":"CLEAR"
      // }
    ]);
      const items = ref([
        {
          id: 0,
          name: 'Phase Voltage',
          value:[230,230,230],
          error:[0,0,0],
          unit:"V"
        },
        {
          id: 1,
          name: 'Phase Current',
          value:[30,10.5,0.8],
          error:[0,0,0],
          unit:"A"
        },
        {
          id: 2,
          name: 'Power Angle',
          value:[0,0,0],
          error:[0,0,0],
          unit:"deg"
        },
        {
          id: 3,
          name: 'Active Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kW"
        },
        {
          id: 4,
          name: 'Reactive Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kvar"
        },
        {
          id: 5,
          name: 'Apparent Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kVA"
        },     
      ])
      const formattedChannel = computed(() => {
        console.log("channel:", channel.value);
        if (channel.value === "Service") {
          return "Service";
        } else if (channel.value === "Calibrate") {
          return "Calibrate";
        }else if (channel.value === "Command") {
          return "Command";
        } else {
          return "Maintenance";
        }
    
      });
      const items2 = ref([
        {
          id: 0,
          name: 'Phase Voltage',
          value:[230,230,230],
          error:[0,0,0],
          unit:"V"
        },
        {
          id: 1,
          name: 'Phase Current',
          value:[30,10.5,0.8],
          error:[0,0,0],
          unit:"A"
        },
        {
          id: 2,
          name: 'Power Angle',
          value:[0,0,0],
          error:[0,0,0],
          unit:"deg"
        },
        {
          id: 3,
          name: 'Active Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kW"
        },
        {
          id: 4,
          name: 'Reactive Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kvar"
        },
        {
          id: 5,
          name: 'Apparent Power',
          value:[0,0,0],
          error:[0,0,0],
          unit:"kVA"
        },     
      ])
      const channels = reactive({
      "main":{       
        "Phase Voltage": {
            "view":[
              { subTitle: "U_A", value: 0, error: 0, limit:0 },
              { subTitle: "U_B", value: 0, error: 0, limit:0 },
              { subTitle: "U_C", value: 0, error: 0, limit:0 }
            ],
          },
          "Phase Current":{
              "view": [
              { subTitle: "I_A", value: 0, error: 0, limit:0 },
              { subTitle: "I_B", value: 0, error: 0, limit:0 },
              { subTitle: "I_C", value: 0, error: 0, limit:0 },
              { subTitle: "Ig", value: 0, error: 0, limit:0 }
            ],
          },
        "Power Angle": {
            "view":[
              { subTitle: "Angle_A", value: 0, error: 0, limit:0 },
              { subTitle: "Angle_B", value: 0, error: 0, limit:0 },
              { subTitle: "Angle_C", value: 0, error: 0, limit:0 }
            ],
          },
          "Active Power":{
            "view" :[
              { subTitle: "Watt_A", value: 0, error: 0, limit:0 },
              { subTitle: "Watt_B", value: 0, error: 0, limit:0 },
              { subTitle: "Watt_C", value: 0, error: 0, limit:0 }
            ]
          },
          "Reactive Power":{
            "view":[
              { subTitle: "Var_A", value: 0, error: 0, limit:0 },
              { subTitle: "Var_B", value: 0, error: 0, limit:0 },
              { subTitle: "Var_C", value: 0, error: 0, limit:0 }
            ]
          },
          "Apparent Power":{
            "view":[
              { subTitle: "VA_A", value: 0, error: 0, limit:0 },
              { subTitle: "VA_B", value: 0, error: 0, limit:0 },
              { subTitle: "VA_C", value: 0, error: 0, limit:0 }
            ]
          },
      },
      "sub":{       
        "Phase Voltage": {
            "view":[
              { subTitle: "U_A", value: 0, error: 0, limit:0 },
              { subTitle: "U_B", value: 0, error: 0, limit:0 },
              { subTitle: "U_C", value: 0, error: 0, limit:0 }
            ],
          },
          "Phase Current":{
              "view": [
              { subTitle: "I_A", value: 0, error: 0, limit:0 },
              { subTitle: "I_B", value: 0, error: 0, limit:0 },
              { subTitle: "I_C", value: 0, error: 0, limit:0 },
              { subTitle: "Ig", value: 0, error: 0, limit:0 }
            ],
          },
        "Power Angle": {
            "view":[
              { subTitle: "Angle_A", value: 0, error: 0, limit:0 },
              { subTitle: "Angle_B", value: 0, error: 0, limit:0 },
              { subTitle: "Angle_C", value: 0, error: 0, limit:0 }
            ],
          },
          "Active Power":{
            "view" :[
              { subTitle: "Watt_A", value: 0, error: 0, limit:0 },
              { subTitle: "Watt_B", value: 0, error: 0, limit:0 },
              { subTitle: "Watt_C", value: 0, error: 0, limit:0 }
            ]
          },
          "Reactive Power":{
            "view":[
              { subTitle: "Var_A", value: 0, error: 0, limit:0 },
              { subTitle: "Var_B", value: 0, error: 0, limit:0 },
              { subTitle: "Var_C", value: 0, error: 0, limit:0 }
            ]
          },
          "Apparent Power":{
            "view":[
              { subTitle: "VA_A", value: 0, error: 0, limit:0 },
              { subTitle: "VA_B", value: 0, error: 0, limit:0 },
              { subTitle: "VA_C", value: 0, error: 0, limit:0 }
            ]
          },
      },
    });
    provide('channels', channels);
const dataMapping = [
  { index: 0, category: "Phase Voltage", dataCount: 3, refKey: "U" },
  { index: 1, category: "Phase Current", dataCount: 4, refKey: "I" },
  { index: 2, category: "Power Angle", dataCount: 3, refKey: "P" },
  { index: 3, category: "Active Power", dataCount: 3, refKey: null },
  { index: 4, category: "Reactive Power", dataCount: 3, refKey: null },
  { index: 5, category: "Apparent Power", dataCount: 3, refKey: null }
];


function updateChannelData(response) {
  // main과 sub 채널 동시 처리
  ['main', 'sub'].forEach(channelType => {
    const sourceData = channelType === 'main' ? response.mainData : response.subData;
    const refData = channelType === 'main' ? response.mainRef : response.subRef; // refData 가져오기
    
    if (!sourceData) return;
    const THRESHOLD = refData?.Error || 0;
   
    // 매핑 테이블에 따라 자동 업데이트
    dataMapping.forEach(({ index, category, dataCount, refKey }) => {
      if (sourceData[index]?.data) {
        const targetView = channels[channelType][category].view;
        
        // 데이터 배열의 각 값을 view 배열에 할당
        for (let i = 0; i < dataCount && i < targetView.length; i++) {
          if (sourceData[index].data[i] !== undefined) {
            const value = sourceData[index].data[i]["value"];
            let error = 0;
            let limit = 0;
            // refData가 있고 해당 category의 refKey가 있으면 오차 계산
            if (refData && refKey && refData[refKey] !== undefined && refData[refKey] !== 0) {
              let refValue;
              if(index == 1 && i == 3){
                refValue = refData["In"];
              }else{
                refValue = refData[refKey];
              }
              //const refValue = refData[refKey];
              error = (Math.abs(value - refValue) / refValue * 100).toFixed(3); // 백분율로 계산
              limit = parseFloat(error) > THRESHOLD ? 1 : 0; //limit = error > limit ? 1 : 0;
            }
            
            channels[channelType][category].view[i] = {
              ...channels[channelType][category].view[i],
              value: value,
              error: error,
              limit: limit,
            };
          }
        }
      }
    });
  });
}
    const startPolling = () => {
      if (updateInterval) clearInterval(updateInterval);
      updateInterval = setInterval(() => fetchData(), 1000);
    };

    const stopPolling = () => {
      if (updateInterval) {
        clearInterval(updateInterval);
        updateInterval = null;
      }
    };
    const fetchData = async (ch) => {
      try {
        const response = await axios.get(`/config/calibrateNow`);
        updateChannelData(response.data);
        
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    
  
      return {
        sidebarOpen,
        items,
        //langset,
        items2,
        channel,
        formattedChannel,
        commands,
        t,
        channels,
        fetchData,
        startPolling,
        stopPolling,
        //installed,
      }  
    }
    
  }
  </script>