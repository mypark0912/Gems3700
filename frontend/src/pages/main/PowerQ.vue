<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = true" :channel="channel" />

    <!-- Content area -->
    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <!-- Site header -->
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-4 sm:px-6 lg:px-8 py-5 w-full max-w-full">
          <!-- Dashboard actions -->
          <div class="sm:flex sm:justify-between sm:items-center mb-5">
            <!-- Left: Title -->
            <div class="mb-3 sm:mb-0">
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">
                {{
                  channel == "Main" ? t("pq.sitemap.main") : t("pq.sitemap.sub")
                }} {{ assetName }}
                > {{ t("pq.sitemap.title") }}
              </h2>
            </div>
          </div>

          <!-- Cards -->
          <div class="grid grid-cols-1 md:grid-cols-12 gap-6">
            <!-- Card (탭이 포함될 카드 섹션) -->
            <div class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full">
              <!-- Tab Navigation -->
              <div class="px-4">
                <ul class="text-sm font-medium flex flex-nowrap overflow-x-auto no-scrollbar w-full">
                  <li v-for="(tab, index) in tabs" :key="index" class="mr-4 last:mr-0 relative">
                    <button @click="changeTab(tab.name)"
                      class="px-4 py-2 whitespace-nowrap transition duration-200 ease-in-out" :class="activeTab === tab.name
                        ? 'text-violet-500 border-b-2 border-violet-500'
                        : 'text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                        ">
                      {{ tab.label }}
                    </button>
                  </li>
                </ul>
              </div>

              <!-- Tab Content -->
              <div v-for="(tab, index) in tabs" :key="index">
                <div v-if="activeTab === tab.name" class="text-gray-700 dark:text-white text-left pt-3 px-4">
                  <!-- 콤보박스 -->
                  <div class="mt-2 mb-4 flex items-center gap-x-2">
                    <label class="sr-only" :for="tab.name + '-select'">
                      Select Option
                    </label>
                    <select :id="tab.name + '-select'" :value="selectedOptions[tab.name]"
                      class="form-select w-64 flex-shrink-0 bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-violet-500 focus:border-violet-500"
                      @change="selectChanged(tab.name, $event.target.value)">
                      <option v-for="(option, i) in tab.options" :key="i" :value="option.key">
                        {{ option.label }}
                      </option>
                    </select>
                    <div v-if="activeTab === 'Harmonics'" class="flex flex-wrap -space-x-px">
                      <button v-for="option in options" :key="option.value" :value="option.value"
                        @click.prevent="setbtnOption(option.value)" :class="[
                          'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg',
                          btnOptions === option.value
                            ? 'bg-violet-500 text-white border-violet-500'
                            : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                        ]">
                        {{ option.label }}
                      </button>
                    </div>
                    <!--button v-if="activeTab === 'Waveform'"  @click.prevent="startWave()"
                    class="btn border px-4 py-2 transition-colors duration-200 rounded-lg bg-violet-500 text-white border-violet-500"> Show Chart</button-->
                  </div>

                  <!-- 차트 컨테이너 -->
                  <div class="flex flex-col space-y-4">
                    <BarChart v-if="
                      activeTab === 'Harmonics' &&
                      tbdataH !== null &&
                      btnOptions === 'chart'
                    " :data="chartData" :width="600" :height="250" :mode="mode1"
                      :key="`harmonics-chart1-${chartUpdateKey}`" />
                    <BarChart v-if="
                      activeTab === 'Harmonics' &&
                      tbdataH !== null &&
                      btnOptions === 'chart'
                    " :data="chartData2" :width="600" :height="250" :mode="mode2"
                      :key="`harmonics-chart2-${chartUpdateKey}`" />
                    <BarChart v-if="
                      activeTab === 'Harmonics' &&
                      tbdataH !== null &&
                      btnOptions === 'chart'
                    " :data="chartData3" :width="600" :height="250" :mode="mode3"
                      :key="`harmonics-chart3-${chartUpdateKey}`" />
                    <PowerQ_Table v-if="
                      activeTab === 'Harmonics' &&
                      tbdataH !== null &&
                      btnOptions === 'table'
                    " :data="tbdataH" :option="selectedOptions.Harmonics" 
                      :key="`harmonics-table-${selectedOptions.Harmonics}`" />
                    <LineChart v-if="activeTab === 'Waveform'" :data="linechartData" width="495" height="348"
                      :mode="'Voltage'" />
                    <LineChart v-if="activeTab === 'Waveform'" :data="linechartData2" width="495" height="348"
                      :mode="'Current'" />
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
import { ref, watch, computed, nextTick, onMounted, onUnmounted, onBeforeUnmount, reactive } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";
import Sidebar from "../common/SideBar3.vue";
import Header from "../common/Header.vue";
import BarChart from "../../charts/connect/BarChart05.vue";
import LineChart from "../../charts/connect/LineChart02.vue";
import Report_table from "../../partials/inners/power/ReportTable2.vue";
import PowerQ_Table from "../../partials/inners/power/PowerQ_Table.vue";
import { tailwindConfig } from "../../utils/Utils";
import { useI18n } from "vue-i18n";
import { useSetupStore } from '@/store/setup'

export default {
  name: "PowerQ",
  props: ["channel"],
  components: {
    Sidebar,
    Header,
    BarChart,
    LineChart,
    Report_table,
    PowerQ_Table,
  },
  setup(props) {
    const { t } = useI18n();
    const route = useRoute();
    const sidebarOpen = ref(true);
    const channel = ref(props.channel);
    const activeTab = ref("Harmonics");
    const mode1 = ref("L1");
    const mode2 = ref("L2");
    const mode3 = ref("L3");
    const tbdata = ref(null);
    const tbdataH = ref(null);
    const btnOptions = ref("chart");
    const chartUpdateKey = ref(0); // 차트 강제 업데이트를 위한 키
    
    const options = computed(() => [
      { label: t("pq.Button.table"), value: "table" },
      { label: t("pq.Button.chart"), value: "chart" },
    ]);

    const chartData = ref({});
    const chartData2 = ref({});
    const chartData3 = ref({});
    const linechartData = ref({});
    const setupStore = useSetupStore();
    const linechartData2 = ref({});
    const waveUpdate = ref(false);
    const waveInterval = ref(null);  // Waveform용 interval 추가
    const channelComputed = computed(
      () => props.channel || route.params.channel || "Default"
    );

    const assetName = computed(() => {
      const mainName = setupStore.getAssetConfig.assetName_main;
      const subName = setupStore.getAssetConfig.assetName_sub;
      if (channelComputed.value == 'Main' || channelComputed.value == 'main') {
        if (mainName != '') {
          return "(" + mainName + ")";
        } else {
          return "";
        }
      } else {
        if (subName != '') {
          return "(" + subName + ")";
        } else {
          return "";
        }
      }
    })

    const scaledDataV = (waveform) => {
      const vscale = tbdata.value["vscale"];
      if (!Array.isArray(waveform) || typeof vscale !== "number") return [];
      const scaleFactor = Number(vscale.toFixed(5));
      return waveform.map((y) => y * scaleFactor);
    };

    const scaledDataI = (waveform) => {
      const iscale = tbdata.value["iscale"];
      if (!Array.isArray(waveform) || typeof iscale !== "number") return [];
      const scaleFactor = Number(iscale.toFixed(5));
      return waveform.map((y) => y * scaleFactor);
    };

    const waveMap = {
      phaseVoltage: ["Wave Form V1", "Wave Form V2", "Wave Form V3"],
      lineVoltage: ["Wave Form V12", "Wave Form V23", "Wave Form V31"],
      current: ["Wave Form I1", "Wave Form I2", "Wave Form I3"],
    };
    
    const harmonicsMap = {
      phaseVoltage: ["U1", "U2", "U3"],
      lineVoltage: ["Upp1", "Upp2", "Upp3"],
      current: ["I1", "I2", "I3"],
    };

    const selectChanged = (tabName, selectedValue) => {
      //console.log('콤보박스 변경:', tabName, selectedValue);
      selectedOptions[tabName] = selectedValue;
      // 이미 가져온 데이터로 차트만 업데이트 (새로 fetch하지 않음)
      updateChartByOption(tabName, selectedValue);
    };

    const updateChartByOption = async(tab, option) => {
      //console.log('Updating chart for:', tab, 'with option:', option);

      if (tab === "Waveform") {
        if (!tbdata.value) return;

        const [k1, k2, k3] = waveMap[option] || [];
        const scaleFactorV = typeof tbdata.value["vscale"] === "number" ? Number(tbdata.value["vscale"]) : 1;
        const scaleFactorI = typeof tbdata.value["iscale"] === "number" ? Number(tbdata.value["iscale"]) : 1;

        const scale = (arr, isCurrent = false) =>
          Array.isArray(arr) ? arr.map((y) => y * (isCurrent ? scaleFactorI : scaleFactorV)) : [];

        linechartData.value = {
          labels: Array.from({ length: 160 }, (_, i) => i + 1),
          datasets: [
            {
              label: "L1",
              data: scale(tbdata.value[k1], false),
              borderColor: tailwindConfig().theme.colors.violet[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
            {
              label: "L2",
              data: scale(tbdata.value[k2], false),
              borderColor: tailwindConfig().theme.colors.sky[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
            {
              label: "L3",
              data: scale(tbdata.value[k3], false),
              borderColor: tailwindConfig().theme.colors.lime[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
          ],
        };

        const [k11, k12, k13] = waveMap['current'] || [];
        linechartData2.value = {
          labels: Array.from({ length: 160 }, (_, i) => i + 1),
          datasets: [
            {
              label: "L1",
              data: scale(tbdata.value[k11], true),
              borderColor: tailwindConfig().theme.colors.violet[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
            {
              label: "L2",
              data: scale(tbdata.value[k12], true),
              borderColor: tailwindConfig().theme.colors.sky[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
            {
              label: "L3",
              data: scale(tbdata.value[k13], true),
              borderColor: tailwindConfig().theme.colors.lime[500],
              backgroundColor: "transparent",
              borderWidth: 2,
              tension: 0.2,
            },
          ],
        };
        await nextTick();
      }
      else if (tab === "Harmonics") {
        if (!tbdataH.value) return;
        //tbdataH.value = tbdataH.value.slice(2);
        const [k1, k2, k3] = harmonicsMap[option] || [];
        //console.log('Harmonics keys:', k1, k2, k3);
        //console.log('Data available:', !!tbdataH.value[k1], !!tbdataH.value[k2], !!tbdataH.value[k3]);

        chartData.value = {
          labels: Array.from({ length: 62 }, (_, i) => i + 2),
          datasets: [
            {
              label: "L1",
              data: tbdataH.value[k1].slice(2) || [],
              backgroundColor: tailwindConfig().theme.colors.violet[500],
              borderRadius: 4,
            },
          ],
        };
        
        chartData2.value = {
          labels: Array.from({ length: 62 }, (_, i) => i + 2),
          datasets: [
            {
              label: "L2",
              data: tbdataH.value[k2].slice(2) || [],
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              borderRadius: 4,
            },
          ],
        };
        
        chartData3.value = {
          labels: Array.from({ length: 62 }, (_, i) => i + 2),
          datasets: [
            {
              label: "L3",
              data: tbdataH.value[k3].slice(2) || [],
              backgroundColor: tailwindConfig().theme.colors.lime[500],
              borderRadius: 4,
            },
          ],
        };

        // 차트 강제 업데이트를 위해 key 증가
        chartUpdateKey.value++;
        await nextTick();
      }
    };

    const fetchWave = async (ch) => {
      try {
        const response = await axios.get(`/api/getWave/${ch}`);
        if (response.data.success) {
          tbdata.value = { ...response.data.data };
          
          linechartData.value = {
            labels: Array.from({ length: 160 }, (_, i) => i + 1),
            datasets: [
              {
                label: "L1",
                data: scaledDataV(tbdata.value["Wave Form V1"]),
                borderColor: tailwindConfig().theme.colors.violet[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
              {
                label: "L2",
                data: scaledDataV(tbdata.value["Wave Form V2"]),
                borderColor: tailwindConfig().theme.colors.sky[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
              {
                label: "L3",
                data: scaledDataV(tbdata.value["Wave Form V3"]),
                borderColor: tailwindConfig().theme.colors.lime[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
            ],
          };
          
          linechartData2.value = {
            labels: Array.from({ length: 160 }, (_, i) => i + 1),
            datasets: [
              {
                label: "L1",
                data: scaledDataI(tbdata.value["Wave Form I1"]),
                borderColor: tailwindConfig().theme.colors.violet[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
              {
                label: "L2",
                data: scaledDataI(tbdata.value["Wave Form I2"]),
                borderColor: tailwindConfig().theme.colors.sky[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
              {
                label: "L3",
                data: scaledDataI(tbdata.value["Wave Form I3"]),
                borderColor: tailwindConfig().theme.colors.lime[500],
                backgroundColor: "transparent",
                borderWidth: 2,
                tension: 0.2,
              },
            ],
          };
        } else {
          console.warn("서버 응답이 success: false 입니다.");
          tbdata.value = null;
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchData = async (ch) => {
      try {
        // 컴포넌트가 비활성화되었으면 실행하지 않음
        if (!isComponentActive.value) {
          
          return;
        }

        //console.log("Fetching harmonics data for channel:", ch);
        const response = await axios.get(`/api/getHarmonics/${ch}`);
        
        if (response.data.success) {
          // 다시 한번 컴포넌트 상태 확인
          if (!isComponentActive.value) {
            //console.log("데이터 수신 후 컴포넌트 비활성화 확인됨, 처리 중단");
            return;
          }

          // 모든 harmonics 데이터 저장 (U1, U2, U3, Upp1, Upp2, Upp3, I1, I2, I3)
          tbdataH.value = { ...response.data.data };
          //console.log("All harmonics data received and stored");

          // 현재 선택된 옵션에 맞는 차트 업데이트
          updateChartByOption('Harmonics', selectedOptions.Harmonics);

          //console.log("Charts updated successfully");
        } else {
          console.warn("서버 응답이 success: false 입니다.");
          tbdataH.value = null;
        }
      } catch (error) {
        if (isComponentActive.value) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    // 상태 관리를 위한 ref들
    const dataInterval = ref(null);
    const isComponentActive = ref(false);

    // 데이터 갱신 함수
    const refreshData = async () => {
      if (!isComponentActive.value) {
        //console.log("컴포넌트 비활성화로 refreshData 중단");
        return;
      }

      if (channel.value && activeTab.value === 'Harmonics') {
        await fetchData(channel.value);
      }
    };

    // interval 시작 함수
    const startInterval = () => {
      // 기존 interval이 있다면 먼저 정리
      stopInterval();
      
      if (activeTab.value === 'Harmonics' && isComponentActive.value) {
        //console.log("Harmonics 탭 활성화 - interval 시작");
        refreshData(); // 초기 데이터 로드
        dataInterval.value = setInterval(refreshData, 5000);
        //console.log("Interval 시작됨:", dataInterval.value);
      }
    };

    // interval 정리 함수
    const stopInterval = () => {
      if (dataInterval.value) {
        //console.log("Clearing interval:", dataInterval.value);
        clearInterval(dataInterval.value);
        dataInterval.value = null;
        //console.log("Interval 정리됨");
      }
    };

    const stopAllIntervals = () => {
      stopInterval();      // Harmonics interval 정리
      stopWaveInterval();  // Waveform interval 정리
    };

    const refreshWaveData = async () => {
      if (!isComponentActive.value || activeTab.value !== 'Waveform') {
        //console.log("Waveform 갱신 조건 미충족");
        return;
      }

      try {
        const response = await axios.get(`/api/getWave/${channel.value}`);
        if (response.data.success) {
          tbdata.value = { ...response.data.data };
          
          // 현재 선택된 옵션으로 차트 업데이트
          updateChartByOption('Waveform', selectedOptions.Waveform);
          //console.log("Waveform 데이터 갱신됨");
        }
      } catch (error) {
        console.error("Waveform 데이터 갱신 실패:", error);
      }
    };

// Waveform interval 시작 함수
    const startWaveInterval = () => {
      // 기존 interval 정리
      stopWaveInterval();
      
      if (activeTab.value === 'Waveform' && isComponentActive.value && waveUpdate.value) {
        //console.log("Waveform interval 시작");
        waveInterval.value = setInterval(refreshWaveData, 5000);  // 5초마다 갱신
      }
    };

    // Waveform interval 정리 함수
    const stopWaveInterval = () => {
      if (waveInterval.value) {
        //console.log("Clearing wave interval:", waveInterval.value);
        clearInterval(waveInterval.value);
        waveInterval.value = null;
      }
    };

    // startWave 함수 수정
    const startWave = async () => {
      try {
        const response = await axios.get(`/api/setWave/${channel.value}/1`);
        
        if (response.data.success) {
          //console.log('Wave 업데이트 시작 요청 성공');
          waveUpdate.value = true;
          
          // 2초 후에 첫 데이터 가져오기 및 interval 시작
          setTimeout(async () => {
            if (activeTab.value === 'Waveform' && isComponentActive.value) {
              await fetchWave(channel.value);  // 첫 데이터 로드
              startWaveInterval();  // interval 시작
            }
          }, 2000);
        }
      } catch (error) {
        console.error('Wave 업데이트 시작 실패:', error);
      }
    };

    // endWave 함수 수정
    const endWave = async () => {
      try {
        // interval 먼저 정리
        stopWaveInterval();
        
        const response = await axios.get(`/api/setWave/${channel.value}/0`);
        
        if (response.data.success) {
          waveUpdate.value = false;
          //console.log('Wave 업데이트 종료됨');
        }
      } catch (error) {
        console.error('Wave 업데이트 종료 실패:', error);
      }
    };

    // 컴포넌트 마운트 시
    onMounted(() => {
      //console.log("PowerQ 컴포넌트 마운트됨");
      isComponentActive.value = true;
      startInterval();
    });

    // 컴포넌트 언마운트 시
    onBeforeUnmount(() => {
      //console.log("PowerQ 컴포넌트 언마운트 준비됨");
      isComponentActive.value = false;

  
      // 차트 데이터를 빈 객체로 초기화 (차트 컴포넌트가 안전하게 정리되도록)
      chartData.value = { labels: [], datasets: [] };
      chartData2.value = { labels: [], datasets: [] };
      chartData3.value = { labels: [], datasets: [] };
      linechartData.value = { labels: [], datasets: [] };
      linechartData2.value = { labels: [], datasets: [] };
      // stopInterval();
      stopAllIntervals();  // 모든 interval 정리
  
      // Wave 업데이트가 진행 중이었다면 종료
      if (waveUpdate.value) {
        endWave();
      }
    });

    onUnmounted(() => {
      //console.log("PowerQ 컴포넌트 완전히 언마운트됨");
      // 한번 더 확실하게 정리
      isComponentActive.value = false;
      stopInterval();
      if(waveUpdate.value)
        endWave();
    });

    // 채널 변경 감지
    // watch(
    //   () => route.params.channel,
    //   async (newChannel) => {
    //     console.log("Updated Channel:", newChannel);
        
    //     // 채널 변경 시에만 interval 재시작
    //     channel.value = newChannel;
        
    //     // 차트 데이터 초기화
    //     chartData.value = {};
    //     chartData2.value = {};
    //     chartData3.value = {};
    //     linechartData.value = {};
    //     linechartData2.value = {};
    //     tbdataH.value = null; // harmonics 데이터도 초기화
        
    //     // interval 재시작 (새 채널의 데이터를 가져오기 위해)
    //     startInterval();
    //   },
    //   { immediate: true }
    // );
    watch(
      () => route.params.channel,
      async (newChannel) => {
        //console.log("Updated Channel:", newChannel);
        
        // 모든 interval 정리
        stopAllIntervals();
        
        // Wave 업데이트가 진행 중이었다면 종료
        if (waveUpdate.value) {
          await endWave();
        }
        
        channel.value = newChannel;
        
        // 차트 데이터 초기화
        chartData.value = {};
        chartData2.value = {};
        chartData3.value = {};
        linechartData.value = {};
        linechartData2.value = {};
        tbdata.value = null;
        tbdataH.value = null;
        
        // 현재 탭에 따라 적절한 interval 시작
        if (activeTab.value === 'Harmonics') {
          startInterval();
        } else if (activeTab.value === 'Waveform') {
          startWave();
        }
      },
      { immediate: true }
    );

    // 탭 변경 시
    watch(activeTab, async (newTab, oldTab) => {
      //console.log(`탭 변경: ${oldTab} -> ${newTab}`);
      await nextTick();

      if (oldTab === 'Waveform' && waveUpdate.value) {
        // Waveform 탭에서 나갈 때 업데이트 종료
        await endWave();
      }

      if (newTab === 'Harmonics') {
        startInterval();  // Harmonics interval 시작
      } else {
        stopInterval();  // Harmonics interval 정지
      }
      
      if (newTab === 'Waveform') {
        // Waveform 탭 진입 시 Wave 업데이트 시작
        startWave();
      }
    });

    // watch(activeTab, (newTab, oldTab) => {
    //   console.log(`탭 변경: ${oldTab} -> ${newTab}`);

    //   if (newTab === 'Harmonics') {
    //     endWave();
    //     startInterval();
    //   } else {
    //     stopInterval();
    //   }
      
    //   if (newTab === 'Waveform' && !linechartData.value.labels) {
    //     //fetchWave(channel.value);
    //     startWave();
    //   }
    // });

    const tabs = computed(() => [
      {
        name: "Harmonics",
        label: t("pq.tabs.harmonics"),
        options: [
          { key: "phaseVoltage", label: t("pq.options.phaseVoltage") },
          { key: "lineVoltage", label: t("pq.options.lineVoltage") },
          { key: "current", label: t("pq.options.current") },
        ],
      },
      {
        name: "Waveform",
        label: t("pq.tabs.waveform"),
        options: [
          { key: "phaseVoltage", label: t("pq.options.phaseVoltage") },
          { key: "lineVoltage", label: t("pq.options.lineVoltage") },
        ],
      },
    ]);

    const setbtnOption = (value) => {
      //console.log("테이블/차트 버튼 변경:", value);
      btnOptions.value = value;
      
      // 차트/테이블 전환 시 현재 선택된 옵션으로 다시 렌더링
      if (value === 'chart' && tbdataH.value) {
        // 차트로 전환 시 현재 선택된 옵션으로 차트 업데이트
        updateChartByOption('Harmonics', selectedOptions.Harmonics);
      }
      // 테이블은 이미 :option prop으로 selectedOptions.Harmonics를 받고 있으므로 자동 업데이트됨
    };

    // 선택된 옵션 저장
    const selectedOptions = reactive({
      Harmonics: "phaseVoltage",
      Waveform: "phaseVoltage",
    });

    // 탭 변경 함수
    const changeTab = (tabName) => {
      activeTab.value = tabName;
    };

    return {
      t,
      sidebarOpen,
      channel,
      channelComputed,
      activeTab,
      tabs,
      selectedOptions,
      chartData,
      chartData2,
      chartData3,
      chartUpdateKey,
      linechartData,
      linechartData2,
      changeTab,
      mode1,
      mode2,
      mode3,
      tbdata,
      tbdataH,
      options,
      setbtnOption,
      btnOptions,
      fetchWave,
      selectChanged,
      assetName,
      startWave,
      endWave,
      waveUpdate,
      waveInterval,
      startWaveInterval,  // 추가
      stopWaveInterval,  // 추가
      stopAllIntervals,  // 추가
    };
  },
};
</script>