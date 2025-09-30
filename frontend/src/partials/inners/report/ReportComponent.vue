<template>
      <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl mt-4">
            <!-- Tab 1 -->
            <Report_PQ v-if="mode && devMode=='device0'" :channel="channel" />
            <div
              class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg mb-6 mt-6"
            >
              <div
                class="absolute top-0 left-0 right-0 h-0.5 bg-blue-500"
                aria-hidden="true"
              ></div>
              <div
                class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
              >
                <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-blue-500 mr-3">
                    <svg
                      class="w-6 h-6 fill-current text-white"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M3 20h18M5 16l4-4 4 4 6-6"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </div>
                  <h3
                    class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                  >
                    EN 50160
                  </h3>
                </header>
              </div>
              <div class="px-4 py-4 space-y-3">
                <!-- No Data 표시 또는 테이블 표시 -->
                <div v-if="!tbdata || Object.keys(tbdata).length === 0" 
                     class="flex flex-col items-center justify-center py-12">
                  <svg 
                    class="w-16 h-16 text-gray-400 dark:text-gray-600 mb-4" 
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24" 
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path 
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="1.5" 
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                  <p class="text-lg font-medium text-gray-500 dark:text-gray-400">No Data Available</p>
                  <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">EN50160 data will appear here when available</p>
                </div>
                
                <!-- 데이터가 있을 때만 테이블 표시 -->
                <table v-else class="table-auto w-full dark:text-white">
                  <!-- Table header -->
                  <thead class="text-xs uppercase text-gray-400 bg-gray-50 dark:bg-gray-400 dark:text-gray-200 dark:bg-opacity-50 rounded-sm">
                    <tr>
                      <th class="p-2">
                        <div class="font-bold text-left">{{ t('report.cardContext.pqInfo.parameter') }}</div>
                      </th>
                      <th class="p-2">
                        <div class="font-bold text-center">L1</div>
                      </th>
                      <th class="p-2">
                        <div class="font-bold text-center">L2</div>
                      </th>
                      <th class="p-2">
                        <div class="font-bold text-center">L3</div>
                      </th>
                      <th class="p-2">
                        <div class="font-bold text-center">{{ t('report.cardContext.pqInfo.compliance') }}</div>
                      </th>
                      <th class="p-2">
                        <div class="font-bold text-center">{{ t('report.cardContext.pqInfo.RequiredQuality') }}</div>
                      </th>
                    </tr>
                  </thead>
                  <!-- Table body -->
                  <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
                    <tr v-for="(row, index) in data" :key="index">
                      <td class="p-2 text-left">{{ row.parameter }}</td>
                      <td class="p-2 text-center">
                        {{ tbdata[makeKey(row.parameter, 'L1')] ?? '-' }}
                      </td>
                      <td class="p-2 text-center">
                        {{ tbdata[makeKey(row.parameter, 'L2')] ?? '-' }}
                      </td>
                      <td class="p-2 text-center">
                        {{ tbdata[makeKey(row.parameter, 'L3')] ?? '-' }}
                      </td>
                      <td class="p-2 text-center">
                        {{ getComp(row.parameter) ?? '-' }}
                      </td>
                      <td class="p-2 text-center">{{ row.required }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
            <!-- ITIC Section -->
            <div
              class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
            >
              <div
                class="absolute top-0 left-0 right-0 h-0.5 bg-green-500"
                aria-hidden="true"
              ></div>
              <div
                class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
              >
                <header class="flex items-center mb-2">
                    <div class="w-6 h-6 rounded-full shrink-0 bg-green-500 mr-3">
                    <svg
                      class="w-6 h-6 fill-current text-white"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2"
                        stroke="currentColor"
                        stroke-width="2"
                        fill="none"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <circle cx="3" cy="12" r="1.5" fill="currentColor" />
                      <circle cx="21" cy="12" r="1.5" fill="currentColor" />
                      <path
                        d="M5 12v7m14-7v7"
                        stroke="currentColor"
                        stroke-width="2"
                        fill="none"
                        stroke-linecap="round"
                      />
                    </svg>
                  </div>
                  <h3
                    class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                  >
                    ITIC
                  </h3>
                </header>
              </div>
              <div class="px-4 py-4 space-y-3">
                <div class="flex flex-wrap -space-x-px">
                    <!--<button
                    v-for="option in options2"
                    :key="option.value"
                    :value="option.value"
                    @click.prevent="setSelectedOption2(option.value)"
                    :class="[
                            'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg',
                            selectedOption2 === option.value
                              ? 'bg-violet-500 text-white border-violet-500'
                              : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900'
                          ]">
                    {{ option.label }}
                  </button>  -->
                  </div>
                <div class="flex items-center gap-x-6 mb-4">
                  <div class="flex items-center gap-x-2">
                    <span class="w-3 h-3 rounded-full" style="background-color: orange"></span>
                    <span class="text-sm text-gray-600 dark:text-gray-400">sag/interruption</span>
                  </div>
                  <div class="flex items-center gap-x-2">
                    <span class="w-3 h-3 rounded-full" style="background-color: green"></span>
                    <span class="text-sm text-gray-600 dark:text-gray-400">swell</span>
                  </div>
                </div>
                <LineChart :data="linechartData" width="595" height="248" />
              </div>
            </div>
            <!-- <Report_PQ_detail /> -->
        </div>

  </template>
  
  <script>
  import { ref, onMounted, computed } from 'vue'
  import LineChart from '../../../charts/connect/LineChart_ITIC.vue'
  import Report_PQ from './Report_PQ.vue'
  import Report_PQ_detail from './Report_PQ_detail.vue'
  import axios from 'axios'
  import { useI18n } from 'vue-i18n'
  import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
  export default {
    name: 'ReportComponent',
    components:{
        LineChart,
        Report_PQ,
        Report_PQ_detail,
    },
    props: {
      data: {
        type: Array,
        default: () => []
      },
      channel:{
        type: String,
        default:''
      },
      mode: {
        type: Boolean,
        default: false
      },
    },
    setup(props) {
      const { t } = useI18n();
      const channel = ref(props.channel);
      const mode = computed(()=>props.mode );
      const selectedX = ref(null);
      const iticDataList = ref([]);
      const selectedOption2 = ref('voltage_sag')
      const authStore = useAuthStore();
      const options2 = ref([
        { label: "Voltage Sag", value: "voltage_sag" },
        { label: "Over Voltage", value: "over_voltage" },
        { label: "Short Interruptions", value: "short_interruptions" },
      ]);
      const devMode = computed(()=> authStore.getOpMode);
      const tbdata = ref(null);
      
      const COMPLIANCE_BITS = {
        "Frequency Variation 1": [0],
        "Frequency Variation 2": [1],
        "Voltage Variation 1": [2, 3, 4],
        "Voltage Variation 2": [5, 6, 7],
        "Voltage Unbalance": [8],
        "THD": [9, 10, 11],
        "Harmonics": [12, 13, 14],
        "Pst": [15, 16, 17],
        "Plt": [18, 19, 20],
        "Signal Vol.": [21, 22, 23]
      }

      onMounted(()=>{
        fetchData();
        fetchITICData();
        setSelectedOption2(selectedOption2.value); 
      });

      const fetchData = async () => {
        try {
          const response = await axios.get(`/api/getEn50160/${channel.value}`);
          if (response.data.success) {
              tbdata.value = response.data.data
          } else {
              console.warn("서버 응답이 success: false 입니다.");
              tbdata.value = {}; // 빈 객체로 설정
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
          tbdata.value = {}; // 에러 시에도 빈 객체로 설정
        }
      };

      const parseMask = (mask) => {
          const phases = [];
          if (mask & 0b001) phases.push("L1");
          if (mask & 0b010) phases.push("L2");
          if (mask & 0b100) phases.push("L3");
          return phases;
      }

      const getfinValue = (phaselist, levellist, mode) =>{
        let finalValue = 0;
        if (phaselist.length === 1) {
            // phaselist가 하나일 때 해당하는 Level 값 선택
            switch (phaselist[0]) {
              case "L1":
                finalValue = levellist[0];
                break;
              case "L2":
                finalValue = levellist[1];
                break;
              case "L3":
                finalValue = levellist[2];
                break;
              default:
                finalValue = null; // 예상 외 값일 경우 처리
            }
          } else if (phaselist.length > 1) {
            // 2개 이상일 때 가장 큰 값 선택
            if (mode == 'SWELL')
              finalValue = Math.max(...levellist);
            else
            finalValue = Math.min(...levellist);
          } else {
            finalValue = null; // phaselist가 비어있는 경우
          }
          return finalValue;
      }

      const fetchITICData = async () => {
        try {
          const response = await axios.get(`/api/getITIC/${channel.value}`);
          if (response.data.success) {
              const data = response.data.data;
              const ratedV = parseInt(response.data.ratedV);
              if(iticDataList.value.length > 0)
                iticDataList.value = [];
              for (let i = 0; i < data.length;i++){
                let phaselist = parseMask(data[i]["mask"]);
                  let levellist = [];
                  levellist.push(data[i]["level_l1"]);
                  levellist.push(data[i]["level_l2"]);
                  levellist.push(data[i]["level_l3"]);
                  let Yvalue = (getfinValue(phaselist, levellist, data[i].event_type)*100)/ratedV;
                  let Xvalue = data[i]["duration"]/1000.0;
                  iticDataList.value.push({
                    label: 'Selected Point',
                    type: 'scatter',
                    data: [{ x: Xvalue, y: Yvalue.toFixed(2) }],
                    backgroundColor: data[i].event_type == 'SWELL'? 'green' : 'orange',
                    pointRadius: 3,
                    pointHoverRadius: 8,
                    showLine: false,
                  });
              }
          } else {
              console.warn("서버 응답이 success: false 입니다.");
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      };
  
      const setSelectedOption2 = (value) => {
        selectedOption2.value = value;
        switch(selectedOption2.value){
          case "voltage_sag":
            selectedX.value = 0.02;
            break;
          case "over_voltage":
            selectedX.value=0.5;
            break;
          case "short_interruptions":
            selectedX.value=10;
            break;
        }
      };

      const baseChart = {
        datasets: [
          {
            label: 'Series 0',
            data: [
              { x: 0.0001, y: 500 },
              { x: 0.001, y: 200 },
              { x: 0.003, y: 140 },
              { x: 0.003, y: 120 },
              { x: 0.02, y: 120 },
              { x: 0.5, y: 120 },
              { x: 0.5, y: 110 },
              { x: 10, y: 110 },
              { x: 100, y: 110 },
            ],
            borderColor: 'blue',
            backgroundColor: 'transparent',
            borderWidth: 2,
            pointRadius: 0,
            tension: 0,
          },
          {
            label: 'Series 1',
            data: [
              { x: 0.02, y: 0 },
              { x: 0.02, y: 70 },
              { x: 0.5, y: 70 },
              { x: 0.5, y: 80 },
              { x: 10, y: 80 },
              { x: 10, y: 90 },
              { x: 100, y: 90 },
            ],
            borderColor: 'red',
            backgroundColor: 'transparent',
            borderWidth: 2,
            pointRadius: 0,
            tension: 0,
          }
        ]
      }

      const linechartData = computed(() => {
        const base = JSON.parse(JSON.stringify(baseChart));
        if(iticDataList.value.length > 0){
          for(let i = 0 ; i < iticDataList.value.length;i++){
            base.datasets.push(iticDataList.value[i]);
          }
        }
        // if (selectedX.value != null) {
        //   base.datasets.push({
        //       label: 'Selected Point',
        //       type: 'scatter',
        //       data: [{ x: selectedX.value, y: 125 }],
        //       backgroundColor: 'orange',
        //       pointRadius: 6,
        //       pointHoverRadius: 8,
        //       showLine: false,
        //     },
        //   {
        //       label: 'Selected Point2',
        //       type: 'scatter',
        //       data: [{ x: selectedX.value+20, y: 225 }],
        //       backgroundColor: 'green',
        //       pointRadius: 6,
        //       pointHoverRadius: 8,
        //       showLine: false,
        //     })
        // }

        return base
      })

      const makeKey = (param, phase) => {
        const suffixMap = {
          L1: "L1",
          L2: "L2",
          L3: "L3",
          "Multi Phase": "Multi Phase"
        }

        // 예외 처리
        if (param === "Frequency Variation 1") {
          return phase === 'L1' ? "Frequency Variation 1" : undefined
        }
        if (param === "Frequency Variation 2") {
          return phase === 'L1' ? "Frequency Variation 2" : undefined
        }
        if (param === "Voltage Unbalance") {
          return phase === 'L1' ? "Voltage Unbalance" : undefined
        }
        if (param === "Voltage Variation 1") return `Voltage Variation 1 ${suffixMap[phase]}(%)`
        if (param === "Voltage Variation 2") return `Voltage Variation 2 ${suffixMap[phase]}(%)`
        if (param === "THD") return `THDs Variation ${suffixMap[phase]}(%)`
        if (param === "Harmonics") return `Harmonics Variatiopn ${suffixMap[phase]}(%)`
        if (param === "Pst") return `Flickers Pst ${suffixMap[phase]}(%)`
        if (param === "Plt") return `Flickers Plt ${suffixMap[phase]}(%)`
        if (param === "Signal Vol.") return `Signaling Voltage ${suffixMap[phase]}(%)`
        if (param === "Voltage Sag") return `Voltage Dips ${suffixMap[phase]}`
        if (param === "Voltage Swell") return `Voltage Swells ${suffixMap[phase]}`
        if (param === "Short Interruption") return `Short Interruptions ${suffixMap[phase]}`
        if (param === "Long Interruption") return `Long Interruptions ${suffixMap[phase]}`
        if (param === "Signaling Volt.") return `Signaling Voltage ${suffixMap[phase]}(%)`

        // 기본값
        return `${param} ${suffixMap[phase]}`
      }
      
      const getComp = (param) => {
        if (!tbdata.value || !tbdata.value["status&compliance"]) return "-"
        
        const comp = tbdata.value["status&compliance"]
        const bits = COMPLIANCE_BITS[param]
        if (!bits) return "-"
        
        const isAnyOK = bits.some(bit => (comp & (1 << bit)) !== 0)
        return isAnyOK ? "Failed" : "OK"
      }
  
      return {
        setSelectedOption2,
        channel,
        linechartData,
        options2,
        selectedOption2,
        mode,
        tbdata,
        makeKey,
        getComp,
        selectedX,
        t,
        iticDataList,
        devMode,
      }
    }
  }
  </script>