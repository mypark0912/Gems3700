<template>
  <div class="flex h-[100dvh] overflow-hidden">

    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" :channel="channel" />

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
                  <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">{{ t('event.sitemap.title') }} > {{ channel == 'Main' ? t('event.sitemap.main'):t('event.sitemap.sub') }} {{ assetName }}</h2>
                </div>
              </div>

              <!-- Content: Flex container for EventTables -->
              <div class="grid grid-cols-1 md:grid-cols-12 gap-6">
            
            <!-- Card (탭이 포함될 카드 섹션) -->
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
                      {{ tab.label }}
                    </button>
                  </li>
                </ul>
              </div>

              <!-- Tab Content -->
              <div v-for="(tab, index) in tabs" :key="index">
                <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-white text-left pt-3 px-4">
                
                  <!-- 차트 컨테이너 -->
                  <div class="flex flex-col space-y-2">
                    <AlarmStatusTable v-if="activeTab === 'Status'" :mode="mode3" :channel="channelComputed" :key="`status-${channelComputed}`"/>
                    <EventTab v-if="activeTab === 'AlarmLog'" :mode="mode2" :channel="channelComputed" :key="`alarmlog-${channelComputed}`"/>
                    <EventTab v-if="activeTab === 'EventLog'" :mode="mode1" :channel="channelComputed" :key="`eventlog-${channelComputed}`"/>
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
import { ref, computed, nextTick } from 'vue';
import { useRoute } from 'vue-router'
import Sidebar from "../common/SideBar3.vue";
import Header from '../common/Header.vue';
//import { useAuthStore } from '@/store'; // ✅ Pinia Store 사용
import EventTab from '../../partials/inners/event/EventTab.vue'
import AlarmStatusTable from '../../partials/inners/event/AlarmStatusTable.vue'
import Footer from "../common/Footer.vue";
import { useI18n } from "vue-i18n";
import { useSetupStore } from '@/store/setup'
export default {
  name: 'Event',
  components: {
    Sidebar,
    Header,
    Footer,
    EventTab,
    AlarmStatusTable,
  },
  props:['channel'],
  setup(props) {
    const { t } = useI18n();
    const sidebarOpen = ref(false)
    const mode1=ref('Event');
    const mode2=ref('Alarm Log');
    const mode3=ref('Alarm Status');
    //const authStore = useAuthStore();
    const route = useRoute()
    const setupStore = useSetupStore();
    //const langset = computed(() => authStore.getLang);
    const channelComputed = computed(() => props.channel || route.params.channel || 'Default')
    const activeTab = ref('Status');
    const assetName = computed(()=> {
      const mainName = setupStore.getAssetConfig.assetNickname_main;
      const subName = setupStore.getAssetConfig.assetNickname_sub;
      if(channelComputed.value == 'Main' || channelComputed.value == 'main'){
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
    const tabs = computed(() =>[
      { name: 'Status', label: t("event.tabs.alarmStatus")},
      { name: 'AlarmLog', label:t("event.tabs.alarmLog")},
      { name: 'EventLog', label:t("event.tabs.eventLog")},
      // { name: 'EN50160', label: 'EN50160', options: ['Voltage Variations', 'Flicker', 'Harmonic Distortion'] },
      // { name: 'ITIC', label: 'ITIC', options: ['Voltage Sag', 'Overvoltage', 'Short Interruptions'] },
    ]);

    const changeTab = (tabName) => {
      activeTab.value = tabName;
      nextTick(() => {
      });
    };

    return {
      t,
      sidebarOpen,
      mode1,
      mode2,
      mode3,
      //langset,
      channelComputed,
      tabs,
      activeTab,
      changeTab,
      assetName,
    }  
  }
}
</script>