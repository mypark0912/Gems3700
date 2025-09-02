<template>
    <div class="flex h-[100dvh] overflow-hidden">
  
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = true" :channel="channel"/>
  
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
                <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">  {{ channel == 'Main' ? t('diagnosis.sitemap.main'):t('diagnosis.sitemap.sub') }} > 전력품질 진단 </h2>
              </div>
  
  
            </div>
  
            <!-- Cards -->
            <div class="grid grid-cols-12 gap-6 overflow-visible">
              <div class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full">
                <!-- Tab Navigation -->
                <div class="px-4">
                    <ul class="text-sm font-medium flex flex-nowrap overflow-x-auto no-scrollbar w-full">
                        <li v-for="(tab, index) in tabs" :key="index" class="mr-4 last:mr-0 relative">
                        <button
                            @click="changeTab(tab.name)"
                            class="px-4 py-2 whitespace-nowrap transition duration-200 ease-in-out"
                            :class="activeTab === tab.name
                            ? 'text-violet-500 border-b-2 border-violet-500'
                            : 'text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'">
                            {{ t(`diagnosis.tabTitle.${tab.label }`)}}
                        </button>
                        </li>
                    </ul>
                </div>
                <!-- Tab Content -->
                <div v-for="(tab, index) in tabs" :key="index">
                    <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-white text-left pt-3 px-4">
                        <!-- 차트 컨테이너 -->
                        <div class="flex flex-col space-y-2"> 
                            <DiagnosisTab v-if="activeTab === 'PowerQuality'" :channel="channel" :asset="asset" :mode="activeTab" />     
                            <DiagnosisTab v-else-if="activeTab === 'Fault'" :channel="channel" :asset="asset" :mode="activeTab" />  
                            <DiagnosisTab v-else :channel="channel" :asset="asset" :mode="activeTab" />                    
                            <!--BarChart v-if="activeTab === 'Harmonics' && tbdata !== null && btnOptions === 'chart'" :data="chartData" :width="600" :height="250" :mode="mode1"/>
                            <BarChart v-if="activeTab === 'Harmonics' && tbdata !== null && btnOptions === 'chart'" :data="chartData2" :width="600" :height="250" :mode="mode2"/-->
                            <!--Report_table v-if="activeTab === 'EN50160'" :data = "tbdata" /-->
                        </div>
                    </div>
                </div>
              </div>
            </div>
  
          </div>
        </main>
        <Footer />
      </div> 
  
    </div>
</template>
  
<script>
import { ref, watch, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSetupStore } from '@/store/setup'; // Pinia Store 
import axios from 'axios'
import Sidebar from '../common/SideBar3.vue'
import Header from '../common/Header.vue'
import Footer from "../common/Footer.vue";
//import DiagnosisTab from '../../partials/dashboard/DiagnosisTab_Trans.vue'
import DiagnosisTab from '../../partials/inners/diagnosis/DiagnosisTab.vue'
//import Diagnosis_Info from '../../partials/dashboard/Diagnosis_Info_Trans.vue'
import Diagnosis_Info from '../../partials/inners/diagnosis/Diagnosis_Info.vue'
import { useI18n } from 'vue-i18n'  

export default {
  name: 'Diagnosis2',
  props:['channel'],
  components: {
    Sidebar,
    Header,
    Footer,
    DiagnosisTab,
    Diagnosis_Info,
  },
  setup(props) {
    const { t } = useI18n();
    const route = useRoute()
    const sidebarOpen = ref(true)
    const channel = ref(props.channel)
    const setupStore = useSetupStore();
    //const langset = computed(() => authStore.getLang);
    const asset = computed(() => setupStore.getAssetConfig);
    const rawdata = ref([]);

    const activeTab = ref('PowerQuality');
    const tabs = ref([
      { name: 'PowerQuality', label: 'PowerQuality' },
      { name: 'Event', label: 'Event' },
      { name: 'Fault', label: 'Fault' },
    ]);

    const changeTab = (tabName) => {
        activeTab.value = tabName;
    };

    // 채널 변경 시 Status 탭으로 리셋하는 함수
    const resetToStatusTab = () => {
        activeTab.value = 'Status';
    };

    const fetchAsset = async () => {
      if(!asset.value)
        await setupStore.checkSetting(); 
      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      //const chType = channel.value == 'Main'? asset.value.assetType_main : asset.value.assetType_sub;
      if(chName != ''){
          try {
          //const ch = 'Fan';
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

    const fetchRealtimeAsset = async () => {
      if(!asset.value)
        await setupStore.checkSetting(); 
      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      const chType = channel.value == 'Main'? asset.value.assetType_main : asset.value.assetType_sub;
      if(chType != 'Transformer'){
          try {
          //const ch = 'Fan';
          const response = await axios.get(`/api/getRealTime/${chType}/${chName}`);
          if (response.data.success) {

            for(let i = 0 ; i < response.data.data.length ; i++)
              rawdata.value.push(response.data.data[i]);
          }else{
            console.log('No Data');
          }
        }catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    onMounted(async () => {
      await fetchAsset();
      //await fetchRealtimeAsset();
    });

    // 라우트 채널 변경 감지 - Status 탭으로 리셋
    watch(() => route.params.channel, async(newChannel) => {
      channel.value = newChannel
      //console.log('Updated Channel:', channel.value)
      resetToStatusTab(); // 🎯 채널 변경 시 Status 탭으로 이동
      await fetchAsset();
    });

    // Props 채널 변경 감지 - Status 탭으로 리셋
    watch(() => props.channel, async(newChannel) => {
      if (newChannel !== channel.value) {
        channel.value = newChannel;
        resetToStatusTab(); // 🎯 채널 변경 시 Status 탭으로 이동
        await fetchAsset();
      }
    });

    return {
      sidebarOpen,
      channel,
      //langset,
      asset,
      fetchAsset,
      activeTab,
      changeTab,
      rawdata,
      tabs,
      t,
    }  
  }
}
</script>