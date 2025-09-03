<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" :channel="channel" />
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
        
        <!-- Site header -->
        <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" :langset="langset" />
  
        <main class="grow">
          <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
            
            <!-- Dashboard actions -->
            <div class="sm:flex sm:justify-between sm:items-center mb-4">

              <!-- Left: Title -->
              <div class="mb-4 sm:mb-0">
                <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">Sub Channel</h2>
              </div>


            </div>
  
            <!-- Cards -->
            <div class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full">
              <GemsSubCard v-if="Object.keys(subDataA).length > 0" :data="subDataA" :channel="'A'" />
              </div>
            <!--div class="grid grid-cols-12 gap-6">

              <GemsDashCard v-if="Object.keys(subDataA).length > 0" :data="subDataA" :channel="'A'"/>
              <GemsDashCard v-if="Object.keys(subDataB).length > 0" :data="subDataB" :channel="'B'"/>
              <GemsDashCard_sub v-if="Object.keys(subData).length > 0" :data="subData" :channel="'B'"/>
            </div-->
  
          </div>
        </main>
        <Footer />
      </div> 
  
    </div>
  </template>
  
  <script>
  import { ref, watch, onMounted, computed, onUnmounted ,nextTick} from 'vue'
  //import Sidebar from '../partials/Sidebar.vue'
  import Sidebar from '../common/SideBarGems.vue'
  import Header from '../common/Header.vue'
  import Footer from "../common/Footer.vue";
  import axios from 'axios'
  import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
  import GemsDashCard from '../../partials/inners/dashboard/GemsDashCard.vue'
  import GemsSubCard from '../../partials/inners/dashboard/GemsSubCard.vue'
  
  export default {
    name: 'GemsDashboard',
    props:['user', 'id'],
    components: {
      Sidebar,
      Header,
      Footer,
      GemsDashCard,
      GemsSubCard,
      //GemsDashCard_sub,
    },
    setup(props) {
      const sidebarOpen = ref(false)
      const authStore = useAuthStore();
      const subDataA = ref({});
      const subDataB = ref({});
      const subDataC = ref({});
      const langset = computed(() => authStore.getLang);
      const opMode = computed(() => authStore.getOpMode);

      // const activeTab = ref('A');
      // const tabs = ref([
      //   { name: 'A', label: '1-12'},
      //   { name: 'C', label: '13-24' },
      //   { name: 'B', label: '25-36' },
      //   // { name: 'EN50160', label: 'EN50160', options: ['Voltage Variations', 'Flicker', 'Harmonic Distortion'] },
      //   // { name: 'ITIC', label: 'ITIC', options: ['Voltage Sag', 'Overvoltage', 'Short Interruptions'] },
      // ]);

      // const changeTab = (tabName) => {
      // activeTab.value = tabName;
      // nextTick(() => {
      //   fetchData(activeTab.value);
      // });
      //};

    onMounted(()=>{
      // activeTab.value = 'A';
      fetchData('A');
    })
      //const subEnable = computed(() => authStore.getSubEnable);
  
      //let updateInterval = null;
  
      const fetchData = async (ch) => {
        try {

            const response = await axios.get(`/api/getSubData/${ch}`);
            if (response.data.success) {
              if(ch == 'A')
                subDataA.value = { ...subDataA.value, ...response.data.data };
              else if(ch == 'B')
                subDataB.value = { ...subDataB.value, ...response.data.data };
              else
                subDataC.value = { ...subDataC.value, ...response.data.data };
            }
          } catch (error) {
            console.log("데이터 가져오기 실패:", error);
          }
      };
  
      // onMounted(() => {
      //   fetchData('A');
      //   fetchData('B');
      //   // updateInterval = setInterval(() => {
      //   //   fetchData('Main');
      //   // }, 60000);
      // });
  
      // onUnmounted(() => {
      //   if (updateInterval) {
      //     clearInterval(updateInterval);
      //   }
      // });
  
      return {
        sidebarOpen,
        langset,
        opMode,
        subDataA,
        subDataB,
        subDataC,
        //changeTab,
        //activeTab,
        //tabs,
        //subEnable,
        //user,
        // stdata,
        // stdata2,
      }  
    }
  }
  </script>