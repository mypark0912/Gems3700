<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <!-- Sidebar -->
    <Sidebar
      :sidebarOpen="sidebarOpen"
      @close-sidebar="sidebarOpen = true"
      :channel="channel"
    />

    <!-- Content area -->
    <div
      class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden"
    >
      <!-- Site header -->
      <Header
        :sidebarOpen="sidebarOpen"
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      />

      <main class="grow">
        <div class="px-4 sm:px-6 lg:px-8 py-5 w-full max-w-full">
          <!-- Dashboard actions -->
          <div class="sm:flex sm:justify-between sm:items-center mb-5">
            <!-- Left: Title -->
            <div class="mb-3 sm:mb-0">
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">{{ t('trend.sitemap.title') }} > {{ channel == 'Main' ? t('trend.sitemap.main'):t('trend.sitemap.sub') }} {{ assetName }}</h2>
            </div>
          </div>
          <!-- 날짜 범위 선택 및 Apply 버튼 -->
          <div class="mb-4 flex gap-4 items-center">
            <label class="flex items-center">
              <span class="mr-2" style="white-space: nowrap">{{ t('trend.date.StartDate') }}:</span>

              <flat-pickr
                v-model="startDate"
                :config="dateConfig"
                class="form-input w-full p-2 border border-gray-300 rounded-md text-gray-700"
              />

              <!-- <input
                type="date"
                v-model="startDate"
                class="border p-1 rounded"
              /> -->
            </label>
            <label class="flex items-center">
              <span class="mr-2" style="white-space: nowrap">{{ t('trend.date.EndDate') }}:</span>

              <flat-pickr
                v-model="endDate"
                :config="dateConfig"
                class="form-input w-full p-2 border border-gray-300 rounded-md text-gray-700"
              />
            </label>
          </div>
          <!-- Cards -->
          <div class="grid grid-cols-1 md:grid-cols-12 gap-6">
            <!-- Card (탭이 포함될 카드 섹션) -->
            <div
              class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full"
            >
              <!-- Tab Navigation -->
              <div class="px-4">
                <ul
                  class="text-sm font-medium flex flex-nowrap overflow-x-auto no-scrollbar w-full"
                >
                  <li
                    v-for="(tab, index) in tabs"
                    :key="index"
                    class="mr-4 last:mr-0 relative"
                  >
                    <button
                      @click="changeTab(tab.name)"
                      class="px-4 py-2 whitespace-nowrap transition duration-200 ease-in-out"
                      :class="
                        activeTab === tab.name
                          ? 'text-violet-500 border-b-2 border-violet-500'
                          : 'text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                      "
                    >
                    {{ t(`trend.tabTitle.${tab.label }`)}}               
                    </button>
                  </li>
                </ul>
              </div>

              <!-- Tab Content -->
              <div v-for="(tab, index) in tabs" :key="index">
                    <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-white text-left pt-3 px-4">
                        <!-- 차트 컨테이너 -->
                        <div class="flex flex-col space-y-2">                                         
                            <TrendTab v-if="activeTab === 'Meter'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab" :asset="asset"/>       
                            <TrendTab v-if="activeTab === 'Energy'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab" :asset="asset"/>  
                            <TrendTab v-if="activeTab === 'PowerQuality'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab" :asset="asset"/>   
                            <TrendTab v-if="activeTab === 'Diagnosis'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab" :asset="asset"/>
                            <TrendTab v-if="activeTab === 'Parameters'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab" :asset="asset"/>
                                             
                        </div>
                    </div>
                </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { useI18n } from "vue-i18n";
import { ref, watch, computed, nextTick, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
import axios from "axios";
import Sidebar from "../common/SideBar3.vue";
import Header from "../common/Header.vue";
import TrendTab from '../../partials/inners/trend/TrendTab.vue'
//import TrendTab_trans from '../../partials/TrendTab_trans.vue'

import { tailwindConfig } from "../../utils/Utils";
import flatPickr from "vue-flatpickr-component";
import "flatpickr/dist/flatpickr.css";
import * as echarts from "echarts";

export default {
  name: "Trend",
  props: ["channel"],
  components: {
    Sidebar,
    Header,
    TrendTab,
    flatPickr,
    //TrendTab_trans,
  },
  setup(props) {
    const { t } = useI18n();
    const route = useRoute();
    const sidebarOpen = ref(true);
    const channel = ref(props.channel);
    const setupStore = useSetupStore();
    const activeTab = ref("Meter"); // 기본 활성 탭 설정
    //const langset = computed(() => authStore.getLang);
    const DigEnable = computed(()=> {
      if(channel.value == 'Main')
        return setupStore.getMainDiagnosis
      else
        return setupStore.getSubDiagnosis
   });
   const assetName = ref('');
   const asset = computed(()=> {

      const assetConfig = setupStore.getAssetConfig;
      console.log("assetConfig",assetConfig);
      if(channel.value == 'Main'){
        if (assetConfig.assetNickname_main !== '')
          assetName.value = "("+assetConfig.assetNickname_main+")";
        else
          assetName.value = ''
        return assetConfig.assetName_main;
      }
      else{
        if (assetConfig.assetNickname_sub !== '')
          assetName.value = "("+assetConfig.assetNickname_sub+")";
        else
          assetName.value = ''
        return assetConfig.assetName_sub;
      }         
    });

   const tabs = computed(()=>{
      if(DigEnable.value){
        return [
          { name: "Meter", label: "Meter" },
          { name: "Energy", label: "Energy" },
          { name: "PowerQuality", label: "Power Quality" },
          { name: "Diagnosis", label: "Diagnosis" },
          { name: "Parameters", label: "Parameters" },

        ];   
      }else{
        return [
          { name: "Meter", label: "Meter" },
          { name: "Energy", label: "Energy" },
        ];  
      }
   })

    const channelComputed = computed(
      () => props.channel || route.params.channel || "Default"
    );

    watch(
      () => route.params.channel,
      async (newChannel) => {
        //console.log("Updated Channel:", newChannel);
        channel.value = newChannel;
      },
      { immediate: true }
    ); // ✅ 컴포넌트가 처음 마운트될 때 실행

    // const tabs = ref([
    //   { name: "Meter", label: "Meter" },
    //   { name: "PowerQuality", label: "Power Quality" },
    //   { name: "Diagnosis", label: "Diagnosis" },
    // ]);

    // ✅ 탭 변경 시 차트 업데이트
    const changeTab = (tabName) => {
      activeTab.value = tabName;
    };

    const startDate = ref(new Date());
    const endDate = ref(new Date());
    const chartContainer = ref(null);
    let chartInstance = null;
    const rawChartData = ref([]);
    const applychart = () => {
      fetchData(channel.value, startDate.value, endDate.value);
    };

    const dateConfig = {
      dateFormat: "Y-m-d", // YYYY-MM-DD 포맷
      defaultDate: new Date(), // 기본값
      onChange: (selectedDates, dateStr) => {
        //console.log("Selected Date:", dateStr);
      },
    };
    return {
      t,
      sidebarOpen,
      channel,
      channelComputed,
      activeTab,
      tabs,
      startDate,
      endDate,
      changeTab,
      //langset,
      applychart,
      dateConfig,
      chartContainer,
      asset,
      assetName,
    };
  },
};
</script>