<template>
    <div class="col-span-full xl:col-span-6 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
      <header class="px-5 py-4 border-b border-gray-100 dark:border-gray-700/60">
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">ITIC</h2>
      </header>
      <div class="p-3">
        <div class="flex flex-wrap justify-start items-center space-x-6 mt-2">
          <div class="flex flex-wrap items-center -m-1.5">
              <div class="m-1.5">
                <!-- Start -->
                  <div class="mt-2 mb-4 flex items-center gap-x-2">
                        <div class="flex flex-wrap -space-x-px">
                          <button
                            v-for="option in options"
                            :key="option.value"
                            :value="option.value"
                            @click.prevent="setSelectedOption(option.value)"
                            :class="[
                                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg',
                                    selectedOption === option.value
                                      ? 'bg-violet-500 text-white border-violet-500'  // ✅ 선택된 버튼: 보라색 배경 + 흰색 텍스트
                                      : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900' // ✅ 선택되지 않은 버튼: 흰 배경 + 보라 텍스트
                                  ]">
                            {{ option.label }}
                          </button>
                        </div>
                  </div>
                <!-- End -->
              </div>
            </div>
        </div>
        <br/>
        <!-- Table -->
        <div class="overflow-x-auto">
            <LineChart :data="linechartData" width="595" height="248" :mode="mode3" />
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue' // onMounted 추가
  import LineChart from '../../../charts/connect/LineChart02.vue'
  import { tailwindConfig } from '../../../utils/Utils'
  export default {
    name: 'Report_table',
    props:{
      channel:{
        type: String,
        default:''
      }
    },
    components:{
        LineChart,
    },
    setup(props) {
      const selectedOption = ref('voltage_sag') // ✅ 초기값 설정
      const channel = ref(props.channel);
      const mode3=ref('L3');
        const options = ref([
          { label: "Voltage Sag", value: "voltage_sag" },
          { label: "Over Voltage", value: "over_voltage" },
          { label: "Short Interruptions", value: "short_interruptions" },
        ]);
  
      const setSelectedOption = (value) => {
        selectedOption.value = value;
        // console.log("선택된 옵션:", value);
        // getData(value); // 데이터 가져오기 실행
      };
  
      const linechartData = ref({
      labels: Array.from({ length: 31 }, (_, i) => i + 1), // ✅ X축: 1~31일
      datasets: [
        // Power Quality Data → Current 라인
        {
          label: "Current",
          data: [
            0.04, 0.48, 0.01, 1.35, 0.01, 0.74, 0, 0.48, 0, 0.21, 
            0, 0.27, 0, 0.12, 0, 0.19, 0, 0.11, 0, 0.24, 
            0, 0.43, 0.01, 0.29, 0, 0.35, 0.01, 0.71, 0, 0.7, 
            0.01
          ],
          borderColor: tailwindConfig().theme.colors.violet[500], // ✅ 보라색 라인
          fill: false,
          borderWidth: 2,
          pointRadius: 0,
          pointHoverRadius: 3,
          pointBackgroundColor: tailwindConfig().theme.colors.violet[500],
          clip: 20,
          tension: 0.2, // ✅ 부드러운 곡선 적용
        }
      ],
    });
      return {
        selectedOption,
        setSelectedOption,
        options,
        channel,
        linechartData,
        mode3,
      }
    }
  }
  </script>
  