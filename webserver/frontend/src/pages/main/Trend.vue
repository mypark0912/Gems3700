<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <!-- Sidebar -->
    <Sidebar
      :sidebarOpen="sidebarOpen"
      @close-sidebar="sidebarOpen = true"
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
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">{{ t('trend.sitemap.title') }}</h2>
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
                          : 'text-gray-500 dark:text-gray-100 hover:text-gray-600 dark:hover:text-gray-300'
                      "
                    >
                    {{ t(`trend.tabTitle.${tab.label }`)}}               
                    </button>
                  </li>
                </ul>
              </div>

              <!-- Tab Content -->
              <div v-for="(tab, index) in tabs" :key="index">
                    <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-gray-100 text-left pt-3 px-4">
                        <!-- 차트 컨테이너 -->
                        <div class="flex flex-col space-y-2">                                         
                            <TrendTab v-if="activeTab === 'Meter'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>       
                            <TrendTab v-if="activeTab === 'Energy'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>  
                            <TrendTab v-if="activeTab === 'PowerQuality'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>   
                            <TrendTab v-if="activeTab === 'Diagnosis'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>
                            <TrendTab v-if="activeTab === 'Demand'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>
                            <TrendTab v-if="activeTab === 'Parameters'" :key="`${activeTab}-${channel}`" :channel="channel" :startdate="startDate" :enddate="endDate" :tap="activeTab"/>
                                             
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
import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
import axios from "axios";
import Sidebar from "../common/SideBar.vue";
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
    const sidebarOpen = ref(true);
    const channel = ref('Main');
    const setupStore = useSetupStore();
    const activeTab = ref("Meter"); // 기본 활성 탭 설정
    //const langset = computed(() => authStore.getLang);
    const DigEnable = computed(() => false);

   // Meter 파라미터 설정 여부 확인
   const hasMeterConfig = ref(false);
   
   const checkMeterConfig = async () => {
      try {
        const response = await axios.get(`/api/getTrendParameters/${channel.value}`);
        if (response.data.success) {
          const enabledParams = response.data.data.params;
          // Energy 관련 파라미터와 "None"을 제외한 나머지가 있는지 확인
          const filteredParams = enabledParams.filter(param =>
            param !== "None" && !["Energy", "Active Energy", "Reactive Energy", "Apparent Energy"].includes(param)
          );
          hasMeterConfig.value = filteredParams.length > 0;
        } else {
          hasMeterConfig.value = false;
        }
      } catch (error) {
        console.log("Meter 설정 확인 실패:", error);
        hasMeterConfig.value = false;
      }
   };

   const tabs = computed(()=>{
      const tabList = [];
      
      // Meter 탭은 파라미터 설정이 있을 때만 추가
      if (hasMeterConfig.value) {
        tabList.push({ name: "Meter", label: "Meter" });
      }
      
      // Energy 탭은 항상 표시
      tabList.push({ name: "Energy", label: "Energy" });
      tabList.push({ name: "Demand", label: "Demand" });
      // Diagnosis가 활성화된 경우 추가 탭들
      if(DigEnable.value){
        tabList.push(
          { name: "PowerQuality", label: "Power Quality" },
          { name: "Diagnosis", label: "Diagnosis" },
          { name: "Parameters", label: "Parameters" }
        );
      }
      
      return tabList;
   });

    const channelComputed = computed(() => props.channel);

    // ✅ Meter 탭이 숨겨졌을 때 기본 탭을 첫 번째 탭으로 변경
    watch([hasMeterConfig, tabs], ([hasMeter, currentTabs]) => {
      if (!hasMeter && activeTab.value === 'Meter' && currentTabs.length > 0) {
        activeTab.value = currentTabs[0].name;
      }
    });

    // ✅ 컴포넌트 마운트 시 Meter 설정 확인 및 기본 탭 설정
    onMounted(async () => {
      await checkMeterConfig();
      // Meter 설정이 없으면 첫 번째 탭으로 전환
      if (!hasMeterConfig.value && tabs.value.length > 0) {
        activeTab.value = tabs.value[0].name;
      }
    });

    // const tabs = ref([
    //   { name: "Meter", label: "Meter" },
    //   { name: "PowerQuality", label: "Power Quality" },
    //   { name: "Diagnosis", label: "Diagnosis" },
    // ]);

    // ✅ 탭 변경 시 차트 업데이트
    const changeTab = (tabName) => {
      activeTab.value = tabName;
    };
    
    // ✅ Meter 탭이 숨겨졌을 때 기본 탭을 Energy로 변경
    watch([hasMeterConfig, tabs], ([hasMeter, currentTabs]) => {
      if (!hasMeter && activeTab.value === 'Meter' && currentTabs.length > 0) {
        activeTab.value = currentTabs[0].name;
      }
    });
    
    const nowDate = new Date();
    const twoDaysAgo = new Date(nowDate.setDate(nowDate.getDate() - 2));
    const startDate = ref(twoDaysAgo);
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
    };
  },
};
</script>

<style>
/* 전역 다크모드 스타일 - select 드롭다운에 강제 적용 */
.dark select.flatpickr-monthDropdown-months {
  background-color: #1f2937 !important;
  color: #f3f4f6 !important;
  border: 1px solid #4b5563 !important;
  color-scheme: dark;
}

.dark select.flatpickr-monthDropdown-months option {
  background-color: #1f2937 !important;
  color: #f3f4f6 !important;
}
</style>

<style scoped>
/* 다크모드에서 Flatpickr 스타일 개선 */
:deep(.dark .flatpickr-calendar) {
  background: #1f2937 !important;
  border-color: #374151 !important;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3) !important;
}

:deep(.dark .flatpickr-months) {
  background: #1f2937 !important;
}

:deep(.dark .flatpickr-month) {
  background: #1f2937 !important;
  color: #f3f4f6 !important;
}

/* 월 선택 드롭다운 - select 요소 */
:deep(.dark .flatpickr-current-month .flatpickr-monthDropdown-months) {
  background-color: #1f2937 !important;
  color: #f3f4f6 !important;
  border: 1px solid #4b5563 !important;
}

:deep(.dark .flatpickr-current-month .flatpickr-monthDropdown-months:hover) {
  background-color: #1f2937 !important;
}

:deep(.dark .flatpickr-current-month .flatpickr-monthDropdown-months:focus) {
  background-color: #1f2937 !important;
  outline: none !important;
}

/* 드롭다운 옵션들 */
:deep(.dark .flatpickr-current-month select.flatpickr-monthDropdown-months option) {
  background-color: #1f2937 !important;
  color: #f3f4f6 !important;
  padding: 8px 12px !important;
}

:deep(.dark .flatpickr-current-month select.flatpickr-monthDropdown-months option:hover) {
  background-color: #374151 !important;
  color: #ffffff !important;
}

:deep(.dark .flatpickr-current-month select.flatpickr-monthDropdown-months option:checked),
:deep(.dark .flatpickr-current-month select.flatpickr-monthDropdown-months option:focus) {
  background-color: #8b5cf6 !important;
  color: #ffffff !important;
}

/* Static 월 선택기 스타일 (monthSelectorType: 'static' 사용 시) */
:deep(.dark .flatpickr-monthSelect-months) {
  background: #1f2937 !important;
  border: 1px solid #374151 !important;
  border-radius: 8px !important;
  padding: 8px !important;
}

:deep(.dark .flatpickr-monthSelect-month) {
  background: #374151 !important;
  color: #f3f4f6 !important;
  border-radius: 4px !important;
  margin: 4px !important;
  padding: 8px !important;
  cursor: pointer !important;
}

:deep(.dark .flatpickr-monthSelect-month:hover) {
  background: #4b5563 !important;
  color: #ffffff !important;
}

:deep(.dark .flatpickr-monthSelect-month.selected) {
  background: #8b5cf6 !important;
  color: #ffffff !important;
}


/* 년도 입력 */
:deep(.dark .flatpickr-current-month input.cur-year) {
  color: #f3f4f6 !important;
  background: #374151 !important;
}

:deep(.dark .flatpickr-current-month input.cur-year:hover) {
  background: #4b5563 !important;
}

/* 요일 헤더 */
:deep(.dark .flatpickr-weekdays) {
  background: #1f2937 !important;
}

:deep(.dark .flatpickr-weekday) {
  color: #f3f4f6 !important;
}

/* 날짜 영역 */
:deep(.dark .flatpickr-days) {
  background: #1f2937 !important;
}

:deep(.dark .flatpickr-day) {
  color: #f3f4f6 !important;
}

:deep(.dark .flatpickr-day:hover) {
  background: #374151 !important;
  border-color: #4b5563 !important;
}

:deep(.dark .flatpickr-day.selected) {
  background: #8b5cf6 !important;
  border-color: #8b5cf6 !important;
  color: #ffffff !important;
}

:deep(.dark .flatpickr-day.today) {
  border-color: #8b5cf6 !important;
}

:deep(.dark .flatpickr-day.disabled) {
  color: #6b7280 !important;
}

/* 월 선택 그리드 뷰 */
:deep(.dark .flatpickr-monthSelect-months) {
  background: #1f2937 !important;
}

:deep(.dark .flatpickr-monthSelect-month) {
  background: #374151 !important;
  color: #f3f4f6 !important;
  border-radius: 4px !important;
  margin: 2px !important;
}

:deep(.dark .flatpickr-monthSelect-month:hover) {
  background: #4b5563 !important;
}

:deep(.dark .flatpickr-monthSelect-month.selected) {
  background: #8b5cf6 !important;
  color: #ffffff !important;
}

/* 숫자 입력 화살표 */
:deep(.dark .numInputWrapper:hover) {
  background: #374151 !important;
}

:deep(.dark .arrowUp),
:deep(.dark .arrowDown) {
  border-bottom-color: #f3f4f6 !important;
}

:deep(.dark .arrowDown) {
  border-top-color: #f3f4f6 !important;
}

/* 이전/다음 달 버튼 */
:deep(.dark .flatpickr-prev-month),
:deep(.dark .flatpickr-next-month) {
  fill: #f3f4f6 !important;
}

:deep(.dark .flatpickr-prev-month:hover),
:deep(.dark .flatpickr-next-month:hover) {
  fill: #ffffff !important;
}

:deep(.dark .flatpickr-prev-month:hover svg),
:deep(.dark .flatpickr-next-month:hover svg) {
  fill: #ffffff !important;
}
</style>