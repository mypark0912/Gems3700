<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl mt-4">
    <div
      class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
    >
      <div
        class="absolute top-0 left-0 right-0 h-0.5 bg-lime-500"
        aria-hidden="true"
      ></div>
      <div
        class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
      >
        <header class="flex items-center mb-2">
            <div class="w-6 h-6 rounded-full shrink-0 bg-lime-500 mr-3">
              <svg
                class="w-6 h-6 fill-current text-white"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M13 2L3 14h7l-1 8 8-12h-6l2-8z"
                  stroke="currentColor"
                  stroke-width="2"
                  fill="none"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
          </div>
          <h3
            class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
          >
            {{ t('report.cardContext.info.energy') }}
          </h3>
        </header>
      </div>
      <div class="px-4 py-4 space-y-3">
    <!-- 토글 버튼 -->
      <div class="flex flex-wrap -space-x-px">
        <button
          v-for="option in options"
          :key="option.value"
          :value="option.value"
          @click.prevent="setSelectedOption(option.value)"
          :class="[
            'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg',
            selectedOption === option.value
              ? 'bg-violet-500 text-white border-violet-500'
              : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900'
          ]"
        >
          {{ option.label }}
        </button>
      </div>
      <BarChart
        v-if="selectedChart"
        :data="selectedChart"
        width="595"
        height="248"
      />

    </div>
  </div> 
      <Report_kwh_detail v-if="assettypes == 'Transformer'" :channel="channel"/>
        <!-- <BarChart v-if="Object.keys(chartData).length > 0" :data="chartData" width="595" height="248" /> -->
      </div>
  </template>
  
  <script>
  import { ref, onMounted, computed } from 'vue' // ✅ onMounted 추가
  //import LineChart from '../charts/LineChart_ITIC.vue'
  import BarChart from '../../../charts/connect/BarChart01_Energy.vue'
  import Report_kwh_detail from './Report_kwh_detail.vue'
  import { tailwindConfig } from '../../../utils/Utils'
  import { useReportData } from "@/composables/ReportDict";
  import { useSetupStore } from '@/store/setup'
  import axios from 'axios'
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  import dayjs from 'dayjs'
  export default {
    name: 'ReportComponent',
    components:{
        //LineChart,
        BarChart,
        Report_kwh_detail,
    },
    props: {
      data: {
        type: Array,
        default: () => [] // ✅ required 제거
      },
      channel:{
        type: String,
        default:''
      },
      mode: {
        type: Boolean,
        default: false // ✅ required 제거
      },
    },
    setup(props) {
      const { t } = useI18n();
      const setupStore = useSetupStore();
      const {
        reportData,
        loadEnergyHourlyData,
        loadEnergyDailyData, // ✅ 함수명 수정
        loadEnergyMonthlyData,
        } = useReportData();
      const selectedOption = ref('hourly') // ✅ 초기값 설정
      const channel = ref(props.channel);
      const assettypes = computed(()=> {
        const assetInfo = setupStore.getAssetConfig;
        if(channel.value == 'Main')
          return assetInfo.assetType_main;
        else
          return assetInfo.assetType_sub;
      })
      const mode = computed(()=>props.mode );
      const isLoading = ref(false); // ✅ 로딩 상태 추가
      const options = ref([
        { label: 'Hourly', value: 'hourly' },
        { label: 'Daily', value: 'daily' },
        { label: 'Monthly', value: 'monthly' },
      ])
      
      const generateHourlyLabels = () => {
        return Array.from({ length: 24 }, (_, i) => `${i.toString().padStart(2, '0')}:00`)
      }

      const generateDailyLabels = () => {
        const today = dayjs()
        const daysInMonth = today.daysInMonth()
        return Array.from({ length: daysInMonth }, (_, i) =>
          today.date(i + 1).format('MM-DD')
        )
      }

      const generateMonthlyLabels = () => {
        return ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
      }

      // ✅ 더미데이터 (fallback용)
      const fakeChartData = {
        hourly: {
          labels: generateHourlyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: 24 }, () => 0),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        },
        daily: {
          labels: generateDailyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: generateDailyLabels().length }, () => 0),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        },
        monthly: {
          labels: generateMonthlyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: 12 }, () => 0),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        },
      }

      // ✅ API 데이터를 차트 형식으로 변환하는 함수들 (디버깅 포함)
      const transformHourlyData = (hourlyData) => {
        //console.log('transformHourlyData 호출됨, 입력 데이터:', hourlyData);
        
        if (!hourlyData || !Array.isArray(hourlyData) || hourlyData.length === 0) {
          console.log('hourly 데이터가 없어서 fakeChartData 사용');
          return fakeChartData.hourly
        }

        //console.log('실제 hourly 데이터 사용, 길이:', hourlyData.length);
        
        const labels = hourlyData.map(item => {
          const hour = item?.hour?.toString().padStart(2, '0') || '00'
          return `${hour}:00`
        })

        const data = hourlyData.map(item => item?.value || 0)

        //console.log('변환된 labels:', labels);
        //console.log('변환된 data:', data);

        return {
          labels,
          datasets: [
            {
              label: 'Import (kWh)',
              data,
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      const transformDailyData = (dailyData) => {
        if (!dailyData || !Array.isArray(dailyData) || dailyData.length === 0) {
          return fakeChartData.daily
        }

        const labels = dailyData.map(item => {
          const date = new Date(item?.date || new Date())
          return dayjs(date).format('MM-DD')
        })

        const data = dailyData.map(item => item?.value || 0)

        return {
          labels,
          datasets: [
            {
              label: 'Import (kWh)',
              data,
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      const transformMonthlyData = (monthlyData) => {
        if (!monthlyData || !Array.isArray(monthlyData) || monthlyData.length === 0) {
          return fakeChartData.monthly
        }

        const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
        const labels = monthlyData.map(item => monthNames[(item?.month || 1) - 1])
        const data = monthlyData.map(item => item?.value || 0)

        return {
          labels,
          datasets: [
            {
              label: 'Import (kWh)',
              data,
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      // ✅ 실제 데이터 또는 더미데이터를 반환하는 computed (디버깅 포함)
      const selectedChart = computed(() => {
        //console.log('selectedChart computed 실행됨');
        //console.log('selectedOption.value:', selectedOption.value);
        //console.log('isLoading.value:', isLoading.value);
        //console.log('reportData:', reportData);
        
        if (isLoading.value) {
          console.log('로딩 중이므로 null 반환');
          return null // 로딩 중에는 차트를 보여주지 않음
        }

        // ✅ 옵셔널 체이닝 사용으로 안전한 접근
        switch (selectedOption.value) {
          case 'hourly':
            //console.log('hourly 데이터:', reportData?.energyHourlyData);
            const hourlyResult = transformHourlyData(reportData?.energyHourlyData);
            //console.log('hourly 변환 결과:', hourlyResult);
            return hourlyResult;
          case 'daily':
            //console.log('daily 데이터:', reportData?.energyDailyData);
            const dailyResult = transformDailyData(reportData?.energyDailyData);
            //console.log('daily 변환 결과:', dailyResult);
            return dailyResult;
          case 'monthly':
            //console.log('monthly 데이터:', reportData?.energyMonthlyData);
            const monthlyResult = transformMonthlyData(reportData?.energyMonthlyData);
            //console.log('monthly 변환 결과:', monthlyResult);
            return monthlyResult;
          default:
            //console.log('기본 fakeChartData 사용');
            return fakeChartData[selectedOption.value];
        }
      })

      const chartData = ref({});
      const tbdata = ref({});

      // ✅ 버튼 클릭 시 데이터 로드하는 함수
      const setSelectedOption = async (value) => {
        selectedOption.value = value;
        isLoading.value = true;

        try {
          console.log("선택된 옵션:", value);
          
          switch (value) {
            case 'hourly':
              await loadEnergyHourlyData(channel.value);
              break;
            case 'daily':
              await loadEnergyDailyData(channel.value);
              break;
            case 'monthly':
              await loadEnergyMonthlyData(channel.value);
              break;
          }
        } catch (error) {
          console.error(`${value} 데이터 로딩 실패:`, error);
        } finally {
          isLoading.value = false;
        }
      };

      const fetchData = async () => {
        try {
          //console.log("Data Get");
          const response = await axios.get(`/api/getEnergyRedis/${channel.value}`);
          if (response.data.success) {
              tbdata.value = response.data.data;
              chartData.value = {
                labels: [
                  'Import', 'Export'
                ],
                datasets: [
                  // Light blue bars
                  {
                    label: 'This Month',
                    data: [
                      Number(tbdata.value["import kwh this month"]), Number(tbdata.value["export kwh this month"])
                    ],
                    backgroundColor: tailwindConfig().theme.colors.sky[500],
                    hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
                    barPercentage: 0.7,
                    categoryPercentage: 0.7,
                    borderRadius: 4,
                  },
                  // Blue bars
                  {
                    label: 'Last Month',
                    data: [
                      Number(tbdata.value["import kwh last month"]), Number(tbdata.value["export kwh last month"])
                    ],
                    backgroundColor: tailwindConfig().theme.colors.violet[500],
                    hoverBackgroundColor: tailwindConfig().theme.colors.violet[600],
                    barPercentage: 0.7,
                    categoryPercentage: 0.7,
                    borderRadius: 4,
                  },
                ],
              };
          } else {
              console.warn("서버 응답이 success: false 입니다.");
              tbdata.value = {}; // 기본값 설정
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      };

      // ✅ 초기 로드: Hourly 데이터만 로드 (디버깅 포함)
      onMounted(async () => {
        //console.log('컴포넌트 마운트됨, channel:', channel.value);
        isLoading.value = true;
        try {
          //console.log('loadEnergyHourlyData 호출 전');
          const result = await loadEnergyHourlyData(channel.value);
          //console.log('loadEnergyHourlyData 결과:', result);
          //console.log('reportData 상태:', reportData);
          //console.log('reportData.energyHourlyData:', reportData?.energyHourlyData);
        } catch (error) {
          console.error("초기 데이터 로딩 실패:", error);
        } finally {
          isLoading.value = false;
          console.log('로딩 완료, isLoading:', isLoading.value);
        }
        // fetchData(); // 실제 API 호출 
      });

      //   const chartData = ref({
      //   labels: [
      //     '02-01-2023', '03-01-2023', '04-01-2023', '05-01-2023',
      //   ],
      //   datasets: [
      //     // Blue bars
      //     {
      //       label: 'New Visitors',
      //       data: [
      //         8000, 3800, 5350, 7800,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.violet[500],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.violet[600],
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //     // Light blue bars
      //     {
      //       label: 'Returning Visitors',
      //       data: [
      //         4000, 6500, 2200, 5800,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.sky[500],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //   ],
      // })

      //   const chartData = ref({
      //   labels: [
      //     '12-01-2022', '01-01-2023', '02-01-2023',
      //     '03-01-2023', '04-01-2023', '05-01-2023',
      //   ],
      //   datasets: [
      //     // Stack
      //     {
      //       label: 'Direct',
      //       data: [
      //         5000, 4000, 4000, 3800, 5200, 5100,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.violet[700],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.violet[800],
      //       barPercentage: 0.7,
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //     // Stack
      //     {
      //       label: 'Referral',
      //       data: [
      //         2500, 2600, 4000, 4000, 4800, 3500,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.violet[500],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.violet[600],
      //       barPercentage: 0.7,
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //     // Stack
      //     {
      //       label: 'Organic Search',
      //       data: [
      //         2300, 2000, 3100, 2700, 1300, 2600,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.violet[300],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.violet[400],
      //       barPercentage: 0.7,
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //     // Stack
      //     {
      //       label: 'Social',
      //       data: [
      //         4800, 4200, 4800, 1800, 3300, 3500,
      //       ],
      //       backgroundColor: tailwindConfig().theme.colors.violet[100],
      //       hoverBackgroundColor: tailwindConfig().theme.colors.violet[200],
      //       barPercentage: 0.7,
      //       categoryPercentage: 0.7,
      //       borderRadius: 4,
      //     },
      //   ],
      // })

  
      // onMounted(() => {
      //   selectedOption.value = 'this_week' // ✅ 마운트 후 강제 설정
      //   console.log('selectOption' + selectedOption.value)
      // });
  
      return {
        selectedOption,
        setSelectedOption,
        selectedChart,
        options,
        channel,
        chartData,
        tbdata,
        t,
        mode,
        isLoading, // ✅ 로딩 상태 노출
        loadEnergyHourlyData,
        loadEnergyDailyData, // ✅ 함수명 수정
        loadEnergyMonthlyData,
        assettypes,
      }
    }
  }
  </script>