<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" :channel="channels" />
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        
        <!-- Site header -->
        <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />
  
        <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <!-- Page header -->
          <div class="sm:flex sm:justify-between sm:items-center mb-4">
            <!-- Left: Title -->
            <div class="mb-4 sm:mb-0">
              <h2
                class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold mb-3"
              >
                {{ t('sidebar.setup') }}  > 
                {{ t(`config.sitemap.${formattedChannel}`) }}
              </h2>
            </div>
            </div>
            <div class="bg-white dark:bg-gray-800 shadow-sm rounded-xl mb-8">
              <ServicePanel v-if="channels == 'Service'"/>
                <div v-else-if="channels == 'Calibrate'" class="flex flex-col md:flex-row md:-mr-px">
                  <CaliPanel :items="items" :channel="channels"/>
                  <CaliSidebar :channel="channels" :commands="commands" />
                </div>
                <Maintenance v-else />
            </div>
        </div>
        </main>
        <Footer />
      </div> 
  
    </div>
  </template>
  
  <script>
  import { useI18n } from 'vue-i18n' 
  import { ref , computed, watch} from 'vue'
  import Sidebar from '../common/SideBar3.vue'
  import Header from '../common/Header.vue'
  //import UsersTabsCard from '../../partials/community/CaliCard.vue'
  import CaliSidebar from '../../partials/inners/config/CaliSideBar2.vue'
  import CaliPanel from '../../partials/inners/config/CaliPanel2.vue'
  import ServicePanel from '../../partials/inners/config/ServicePanel.vue'
  import Maintenance from '../../partials/inners/config/MaintenancePanel.vue'
  //import { useAuthStore } from "@/store"; // ✅ Pinia Store 사용
   import { useRoute } from 'vue-router'
  
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
    },
    setup(props) {
      //const authStore = useAuthStore();
      const route = useRoute();
      const { t } = useI18n();
      //const langset = computed(() => authStore.getLang);
      //const installed = computed(() => authStore.getInstalled);
      const sidebarOpen = ref(false)
      const channels = ref(props.channel);

      watch(() => props.channel, (newChannel) => {
      if (newChannel !== channels.value) {
        channels.value = newChannel;
      }
    }, { immediate: true });

  // route.params.channel 변경 감시
    watch(() => route.params.channel, (newChannel) => {
      if (newChannel !== channels.value) {
        channels.value = newChannel;
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
        "name":"CALV", "Label":"Set U", "Param":"U", "Type":"SET"
      },
      {
        "name":"CLRV", "Label":"Clear U", "Type":"CLEAR"
      },
      {
        "name":"CALVPP", "Label":"Set Upp", "Type":"SET"
      },
      {
        "name":"CLRVPP", "Label":"Clear Upp", "Type":"CLEAR"
      },
      {
        "name":"CALA", "Label":"Set I", "Param":"I", "Type":"SET"
      },
      {
        "name":"CLRA", "Label":"Clear I", "Type":"CLEAR"
      },
      {
        "name":"CALIN", "Label":"Set In", "Param":"In", "Type":"SET"
      },
      {
        "name":"CLRIN", "Label":"Clear In", "Type":"CLEAR"
      },
      {
        "name":"CALW", "Label":"Set Power", "Type":"SET"
      },
      {
        "name":"CLRW", "Label":"Clear Power", "Type":"CLEAR"
      },
      {
        "name":"CALP", "Label":"Set Phase", "Param":"P", "Type":"SET"
      },
      {
        "name":"CLRP", "Label":"Clear Phase", "Type":"CLEAR"
      },
      {
        "name":"CALDCos", "Label":"Set DCos", "Type":"SET"
      },
      {
        "name":"CLRDCos", "Label":"Clear DCos", "Type":"CLEAR"
      }]);
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

        if (channels.value === "Service") {
          return "Service";
        } else if (channels.value === "Calibrate") {
          return "Calibrate";
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
  
      return {
        sidebarOpen,
        items,
        //langset,
        items2,
        channels,
        formattedChannel,
        commands,
        t,
        //installed,
      }  
    }
    
  }
  </script>