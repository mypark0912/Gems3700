<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl mt-4">
    
    <!-- 전력량 & 부하율 추이 카드 -->
    <div class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg mb-4">
      <div class="absolute top-0 left-0 right-0 h-0.5 bg-cyan-500" aria-hidden="true"></div>
      <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-cyan-500 mr-3">
            <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" fill="currentColor"/>
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            {{ t(`report.cardTitle.Energy2`)}} 
          </h3>

        </header>
      </div>
      <div class="px-4 py-3 space-y-2">
        <!-- 이중 Y축 차트 -->
        <div class="dual-axis-section">

          <!--차트 -->
          <div ref="dualAxisChart" class="dual-axis-chart"></div>

          <!--info card -->
          <div class="chart-info">
            <div class="info-card">
              <span class="info-label">{{ t(`report.cardContext.averageLoadRate`)}} </span>
              <span class="info-value">{{ averageLoadRate.toFixed(1) }}%</span>
            </div>
            <div class="info-card">
              <span class="info-label">{{ t(`report.cardContext.maxLoadRate`)}} </span>
              <span class="info-value">{{ maxLoadRate.toFixed(1) }}%</span>
            </div>
            <div class="info-card">
              <span class="info-label">{{ t(`report.cardContext.overloadCount`)}} </span>
              <span class="info-value" :class="{ 'text-red-600 dark:text-red-400': overloadCount > 0 }">
                {{ overloadCount }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 부하율 패턴 히트맵 & 분포 통계 카드 -->
    <div class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
      <div class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500" aria-hidden="true"></div>
      <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
            <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M3 3h18v4H3V3zm0 6h18v2H3V9zm0 4h18v2H3v-2zm0 4h18v4H3v-4z" fill="currentColor"/>
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            {{ t(`report.cardTitle.LoadPatternAnalysis`) }}
          </h3>
   
        </header>
      </div>
      <div class="px-4 py-3 space-y-4">

        <div v-if="heatmapLoading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          <p class="mt-2 text-gray-500">히트맵 데이터 로딩 중...</p>
        </div>


        <div v-else class="heatmap-section">
          <div ref="heatmapChart" class="heatmap-chart"></div>
        </div>


        <div v-if="!heatmapLoading" class="heatmap-stats-section">
          <h4 class="stats-title text-sm font-semibold text-gray-800 dark:text-gray-200 mb-3">{{ t(`report.cardTitle.weeklyLoadDistribution`) }}</h4>
          <div class="load-stats">
            <div class="stat-card light">
              <div class="stat-header">
                <h4>{{ t(`report.cardContext.lightLoadHours`) }}</h4>
                <span class="stat-icon">🟢</span>
              </div>
              <span class="stat-value">{{ apiHeatmapDistribution.light }}%</span>
              <span class="stat-desc">0-50% {{ t(`report.cardContext.loadRate`) }}</span>
            </div>
            <div class="stat-card medium">
              <div class="stat-header">
                <h4>{{ t(`report.cardContext.mediumLoadHours`) }}</h4>
                <span class="stat-icon">🟡</span>
              </div>
              <span class="stat-value">{{ apiHeatmapDistribution.medium }}%</span>
              <span class="stat-desc">50-80% {{ t(`report.cardContext.loadRate`) }}</span>
            </div>
            <div class="stat-card high">
              <div class="stat-header">
                <h4>{{ t(`report.cardContext.heavyLoadHours`) }}</h4>
                <span class="stat-icon">🟠</span>
              </div>
              <span class="stat-value">{{ apiHeatmapDistribution.high }}%</span>
              <span class="stat-desc">80-100% {{ t(`report.cardContext.loadRate`) }}</span>
            </div>
            <div class="stat-card overload">
              <div class="stat-header">
                <h4>{{ t(`report.cardContext.overloadHours`) }}</h4>
                <span class="stat-icon">🔴</span>
              </div>
              <span class="stat-value warning">{{ apiHeatmapDistribution.overload }}%</span>
              <span class="stat-desc">100% {{ t(`report.cardContext.overload`) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import * as echarts from 'echarts'
import { useI18n } from 'vue-i18n'
import { useReportData } from "@/composables/ReportDict";
import dayjs from 'dayjs'

export default {
  name: 'PowerLoadAnalysis',
  props:{
    channel:{
        type: String,
        default:''
      }
  },
  setup(props) {
    // 반응형 데이터
    const { t } = useI18n();
    const {
        reportData,
        getLoadFactorCalculated,
        getHeatmapLoadFactorData,
        loadEnergyHourlyData,
        } = useReportData();
    const selectedTimeRange = ref('24h')
    const channel = ref(props.channel);
    const dualAxisChart = ref(null)
    const heatmapChart = ref(null)
    const heatmapLoading = ref(false)

    // ✅ 실제 API 데이터 저장소
    const apiHeatmapData = ref([])
    const apiHeatmapStats = ref({})

    let dualAxisChartInstance = null
    let heatmapChartInstance = null

    // ✅ reportDict.js의 함수를 사용해서 히트맵 데이터 로드
    const loadHeatmapData = async () => {
      heatmapLoading.value = true;
      try {
        //console.log('히트맵 API 데이터 로딩 시작:', channel.value);
        
        const response = await getHeatmapLoadFactorData(channel.value, 4);
        
        //console.log('히트맵 API 응답:', response);
        
        if (response && response.success) {
          apiHeatmapData.value = response.heatmapData || [];
          apiHeatmapStats.value = response.distribution || {};
          
          //console.log('히트맵 데이터 개수:', apiHeatmapData.value.length);
          //console.log('히트맵 통계:', apiHeatmapStats.value);
          
          if (heatmapChartInstance && apiHeatmapData.value.length > 0) {
            createApiHeatmapChart();
          }
          
          setTimeout(() => {
            //console.log('3초 후 강제 히트맵 생성 시도');
            forceCreateHeatmap();
          }, 3000);
        } else {
          console.error('히트맵 API 실패:', response?.message);
        }
      } catch (error) {
        console.error('히트맵 데이터 로딩 실패:', error);
      } finally {
        heatmapLoading.value = false;
      }
    }

    const forceCreateHeatmap = () => {
      console.log('=== 강제 히트맵 생성 시도 ===');
      console.log('heatmapChart DOM 요소:', heatmapChart.value);
      
      if (!heatmapChart.value) {
        console.error('heatmapChart DOM 요소가 없음!');
        return;
      }
      
      if (heatmapChartInstance) {
        //console.log('기존 차트 인스턴스 해제');
        heatmapChartInstance.dispose();
      }
      
      //console.log('새로운 히트맵 차트 인스턴스 생성');
      heatmapChartInstance = echarts.init(heatmapChart.value);
      
      if (apiHeatmapData.value.length > 0) {
        //console.log('API 데이터로 차트 생성');
        createApiHeatmapChart();
      } else {
        //console.log('테스트 데이터로 차트 생성');
        createTestHeatmap();
      }
    }

    const createTestHeatmap = () => {
      //console.log('테스트 히트맵 생성');
      
      const days = ['월', '화', '수', '목', '금', '토', '일'];
      const hours = Array.from({length: 24}, (_, i) => i.toString().padStart(2, '0'));
      
      const testData = [];
      for (let d = 0; d < 7; d++) {
        for (let h = 0; h < 24; h++) {
          testData.push([h, d, Math.random() * 10]);
        }
      }
      
      const option = {
        title: { text: '테스트 히트맵', left: 'center' },
        tooltip: {
          position: 'top',
          formatter: function (params) {
            return `${days[params.data[1]]} ${hours[params.data[0]]}:00<br/>값: ${params.data[2].toFixed(1)}%`;
          }
        },
        grid: { height: '60%', top: '10%' },
        xAxis: {
          type: 'category',
          data: hours,
          splitArea: { show: true }
        },
        yAxis: { 
          type: 'category', 
          data: days, 
          splitArea: { show: true } 
        },
        visualMap: {
          min: 0,
          max: 10,
          calculable: true,
          orient: 'horizontal',
          left: 'center',
          bottom: '5%',
          inRange: {
            color: ['#f7fbff', '#08306b']
          }
        },
        series: [{
          name: '테스트',
          type: 'heatmap',
          data: testData,
          label: { show: false },
          itemStyle: {
            borderWidth: 1,
            borderColor: '#fff'
          }
        }]
      };
      
      try {
        heatmapChartInstance.setOption(option);
        //console.log('✅ 테스트 히트맵 생성 성공');
      } catch (error) {
        console.error('❌ 테스트 히트맵 생성 실패:', error);
      }
    }
    
    const apiHeatmapDistribution = computed(() => {
      if (!apiHeatmapStats.value || Object.keys(apiHeatmapStats.value).length === 0) {
        return { light: 0, medium: 0, high: 0, overload: 0 };
      }

      const total = apiHeatmapStats.value.light + apiHeatmapStats.value.medium + 
                   apiHeatmapStats.value.high + apiHeatmapStats.value.overload;
      
      if (total === 0) {
        return { light: 0, medium: 0, high: 0, overload: 0 };
      }

      return {
        light: Math.round((apiHeatmapStats.value.light / total) * 100),
        medium: Math.round((apiHeatmapStats.value.medium / total) * 100),
        high: Math.round((apiHeatmapStats.value.high / total) * 100),
        overload: Math.round((apiHeatmapStats.value.overload / total) * 100)
      };
    })

    const createApiHeatmapChart = () => {
      if (apiHeatmapData.value.length === 0) {
        console.warn('히트맵 데이터가 없습니다.');
        return;
      }
      
      //console.log('히트맵 차트 생성 시작, 데이터 길이:', apiHeatmapData.value.length);
      //console.log('차트용 데이터 샘플:', apiHeatmapData.value.slice(0, 5));
      
      const days = [
        t('report.cardContext.days.mon'),
        t('report.cardContext.days.tue'),
        t('report.cardContext.days.wed'),
        t('report.cardContext.days.thu'),
        t('report.cardContext.days.fri'),
        t('report.cardContext.days.sat'),
        t('report.cardContext.days.sun')
      ];
      
      const hours = Array.from({length: 24}, (_, i) => i.toString().padStart(2, '0'));
      
      let chartData = [];
      
      if (apiHeatmapData.value.length > 0 && Array.isArray(apiHeatmapData.value[0]) && apiHeatmapData.value[0].length === 3) {
        //console.log('데이터가 이미 [hour, day, value] 형식입니다.');
        chartData = apiHeatmapData.value;
      } else {
        //console.log('데이터 형식을 [hour, day, value]로 변환합니다.');
        chartData = apiHeatmapData.value.map(item => {
          if (Array.isArray(item)) {
            return item;
          } else if (typeof item === 'object') {
            return [item.hour || 0, item.day_of_week || 0, item.load_factor_percent || 0];
          } else {
            return [0, 0, 0];
          }
        });
      }
      
      //console.log('변환된 차트 데이터 샘플:', chartData.slice(0, 5));
      
      const values = chartData.map(item => item[2] || 0);
      const maxValue = Math.max(...values, 100);
      const minValue = Math.min(...values.filter(v => v > 0), 0);
      
      //console.log('값 범위:', { min: minValue, max: maxValue });
      
      const option = {
          grid: {
          left: '10%',     // 왼쪽 여백 증가 (5% → 10%)
          right: '10%',    // 오른쪽 여백도 균형있게 조정
          top: '12%',      
          bottom: '10%',   
          containLabel: true  // 라벨이 잘리지 않도록 자동 조정
        },
        tooltip: {
          position: 'top',
          formatter: function (params) {
            const hour = hours[params.data[0]] || '00';
            const day = days[params.data[1]] || '알 수 없음';
            const value = params.data[2] || 0;
            return `${day} ${hour}:00<br/>부하율: ${value.toFixed(1)}%`;
          }
        },
        grid: { height: '50%', top: '10%' },
        xAxis: {
          type: 'category',
          data: hours,
          splitArea: { show: true },
          axisLabel: { 
            formatter: function (value) { 
              const hourNum = parseInt(value);
              return hourNum % 2 === 0 ? value + ':00' : '';
            }
          }
        },
        yAxis: { 
          type: 'category', 
          data: days, 
          splitArea: { show: true } 
        },
        visualMap: {
          min: 0,
          max: Math.max(maxValue * 1.2, 10),
          calculable: true,
          orient: 'horizontal',
          left: 'center',
          bottom: '15%',
          inRange: {
            color: [
              '#22c55e', '#84cc16', '#a3e635', '#bef264', '#facc15', 
              '#fbbf24', '#f59e0b', '#f97316', '#ef4444', '#dc2626'
            ]
          },
          text: ['높음', '낮음'],
          textStyle: {
            color: '#333',
            fontSize: 12
          }
        },
        series: [{
          name: '부하율',
          type: 'heatmap',
          data: chartData,
          label: { 
            show: false 
          },
          itemStyle: {
            borderWidth: 0.5,
            borderColor: '#fff'
          },
          emphasis: { 
            itemStyle: { 
              shadowBlur: 10, 
              shadowColor: 'rgba(0, 0, 0, 0.5)',
              borderWidth: 2,
              borderColor: '#333'
            } 
          }
        }]
      }
      
     
      heatmapChartInstance.setOption(option, true);
    }

    // ✅ 더미 데이터 생성 (fallback용)
    const generateTimeSeriesData = (timeRange) => {
      const now = new Date()
      const dataPoints = timeRange === '24h' ? 24 : timeRange === '7d' ? 168 : 720
      const interval = timeRange === '24h' ? 1 : timeRange === '7d' ? 1 : 1
      
      const data = []
      for (let i = dataPoints - 1; i >= 0; i--) {
        const time = new Date(now.getTime() - i * interval * 60 * 60 * 1000)
        const hour = time.getHours()
        const dayOfWeek = time.getDay()
        
        data.push({
          time: time.toISOString(),
          loadRate: 0,
          powerConsumption: 0,
          hour: hour,
          dayOfWeek: dayOfWeek
        })
      }
      
      return data
    }

    const timeSeriesData = ref(generateTimeSeriesData('24h'))

    // ✅ 실제 API 데이터 기반 computed 속성들
    const averageLoadRate = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        const sum = reportData.loadrateData.reduce((acc, item) => acc + (item.load_factor_percent  || 0), 0);
        return sum / reportData.loadrateData.length;
      }
      
      const sum = timeSeriesData.value.reduce((acc, item) => acc + item.load_factor_percent , 0);
      return sum / timeSeriesData.value.length;
    });

    const maxLoadRate = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        return Math.max(...reportData.loadrateData.map(item => item.load_factor_percent  || 0));
      }
      
      return Math.max(...timeSeriesData.value.map(item => item.load_factor_percent ));
    });

    const overloadCount = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        return reportData.loadrateData.filter(item => (item.load_factor_percent  || 0) > 100).length;
      }
      
      return timeSeriesData.value.filter(item => item.load_factor_percent  > 100).length;
    });

const createDualAxisChart = () => {
  //console.log('=== createDualAxisChart 실행 ===');
  //console.log('reportData.energyHourlyData:', reportData.energyHourlyData);
  //console.log('reportData.loadrateData!!!!!!!!!!:', reportData.loadrateData);

  let times = [];
  let powerData = [];
  let loadData = [];
  
  // ✅ 실제 전력량 데이터 처리
  if (reportData.energyHourlyData && Array.isArray(reportData.energyHourlyData) && reportData.energyHourlyData.length > 0) {
    //console.log('실제 전력량 데이터 사용');
    
    const now = dayjs();
    const currentHour = now.hour();
    
    const powerDataMap = {};
    reportData.energyHourlyData.forEach((item, index) => {
 
      
      if (item?.hour !== undefined && item?.value !== undefined) {
        const hourLabel = `${item.hour.toString().padStart(2, '0')}:00`;
        const parsedValue = parseFloat(item.value);
        powerDataMap[hourLabel] = parsedValue;
        
        //console.log(`매핑: ${hourLabel} -> ${parsedValue} (원본: ${item.value})`);
      }
    });
    
    for (let i = 0; i <= currentHour; i++) {
      const hourLabel = `${i.toString().padStart(2, '0')}:00`;
      times.push(hourLabel);
      const value = powerDataMap[hourLabel] || 0;
      powerData.push(value);
      

    }
    
    //console.log('실제 전력량 times:', times);
    //console.log('실제 전력량 data:', powerData);
  } else {
    //console.log('전력량 데이터가 없어서 더미 데이터 사용');
    times = timeSeriesData.value.map(item => {
      const time = new Date(item.time);
      return time.toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit' });
    });
    powerData = timeSeriesData.value.map(() => 0); // ✅ 모든 값을 0으로 변경
  }

  // ✅ 실제 부하율 데이터 처리
  if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
    //console.log('실제 부하율 데이터 사용');
    
    const loadDataMap = {};
    reportData.loadrateData.forEach(item => {
      if (item?.hour !== undefined && item?.load_factor_percent  !== undefined) {
        const hourLabel = `${item.hour.toString().padStart(2, '0')}:00`;
        loadDataMap[hourLabel] = parseFloat(item.load_factor_percent ) || 0;
      }
    });
    
    loadData = times.map(time => loadDataMap[time] || 0);
    
    console.log('실제 부하율 data:', loadData);
  } else {
    //console.log('부하율 데이터가 없어서 더미 데이터 사용');
    if (reportData.energyHourlyData && Array.isArray(reportData.energyHourlyData)) {
      loadData = times.map(() => 0); // ✅ 모든 값을 0으로 변경
    } else {
      loadData = timeSeriesData.value.map(() => 0); // ✅ 모든 값을 0으로 변경
    }
  }

  //console.log('최종 차트 데이터:');
  //console.log('times:', times);
  //console.log('powerData:', powerData);
  //console.log('loadData:', loadData);

  // ✅ 실제 데이터 값 범위 계산 (수정된 버전)
  const powerValues = powerData.filter(val => !isNaN(val) && val > 0);
  const loadValues = loadData.filter(val => !isNaN(val) && val > 0);
  
  const maxPowerValue = powerValues.length > 0 ? Math.max(...powerValues) : 0.01;
  const maxLoadValue = loadValues.length > 0 ? Math.max(...loadValues) : 10;
  
  // ✅ Y축 범위를 깔끔한 값으로 계산
  const powerAxisMax = maxPowerValue > 0 ? 
    Math.ceil(maxPowerValue * 1.3 * 1000) / 1000 : // 소수점 3자리까지 올림
    0.01;
  const loadAxisMax = Math.ceil(Math.max(maxLoadValue * 1.3, 10));
  
  //console.log('전력량 최대값:', maxPowerValue, '-> Y축 최대값:', powerAxisMax);
  //console.log('부하율 최대값:', maxLoadValue, '-> Y축 최대값:', loadAxisMax);

  const option = {
    grid: {
      left: '1%',     // 왼쪽 여백 증가 (5% → 10%)
      right: '1%',    // 오른쪽 여백도 균형있게 조정
      top: '12%',      
      bottom: '10%',   
      containLabel: true  // 라벨이 잘리지 않도록 자동 조정
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross', crossStyle: { color: '#999' } },
      formatter: function (params) {
        let result = params[0].axisValueLabel + '<br/>';
        params.forEach(function (item) {
          if (item.seriesName === '전력량') {
            result += item.marker + item.seriesName + ': ' + item.value + ' kWh<br/>';
          } else {
            result += item.marker + item.seriesName + ': ' + item.value + '%<br/>';
          }
        });
        return result;
      }
    },
    legend: { 
      data: [t('report.cardTitle.Energy'), t('report.cardContext.loadRate')],
      top: '2%',  // ✅ 범례를 위로 올려서 공간 절약
      textStyle: { fontSize: 12 }
    },
    xAxis: [{ 
      type: 'category', 
      data: times,
      boundaryGap: false, // ✅ 카테고리 양쪽 여백 제거
      axisLabel: {
        formatter: function (value) {
          if (times.length > 20) {
            const hourNum = parseInt(value.split(':')[0]);
            return hourNum % 2 === 0 ? value : '';
          }
          return value;
        },
        rotate: times.length > 20 ? 45 : 0
      }
    }],
    yAxis: [
      {
        type: 'value',
        name: t('report.cardTitle.Energy')+' (kWh)',
        position: 'left',
        axisLabel: { 
          formatter: function(value) {
            // 0이면 0, 아니면 소수점 3자리까지 표시
            return value === 0 ? '0' : value.toFixed(3);
          },
          margin: 55
        },
        splitLine: { show: true, lineStyle: { type: 'dashed', color: '#e0e0e0' } },
        min: 0,
        max: powerAxisMax,
        interval: powerAxisMax / 5 // 5개 구간으로 나누기
      },
      {
        type: 'value',
        name: t('report.cardContext.loadRate') + '(%)',
        position: 'right',
        axisLabel: { 
          formatter: function(value) {
            // 0이면 0, 아니면 소수점 3자리까지 표시
            return value === 0 ? '0' : value.toFixed(3);
          },
          fontSize: 10,  // ✅ 폰트 크기 줄여서 공간 절약
          margin: 55      // ✅ 축과 라벨 사이 간격 줄임
        },
        min: 0,
        max: loadAxisMax,
        splitLine: { show: false }
      }
    ],
    series: [
      {
        name: t('report.cardTitle.Energy'),
        type: 'bar',
        yAxisIndex: 0,
        data: powerData,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        }
      },
      {
        name: t('report.cardContext.loadRate'),
        type: 'line',
        yAxisIndex: 1,
        data: loadData,
        lineStyle: { width: 3, color: '#ff6b6b' },
        itemStyle: { color: '#ff6b6b' },
        areaStyle: { opacity: 0.1, color: '#ff6b6b' },
        smooth: true
      }
    ],
    tooltip: {
  trigger: 'axis',
  valueFormatter: (value) => value.toFixed(2)
}
  };
  
  //console.log('ECharts 옵션 설정 완료');
  dualAxisChartInstance.setOption(option);
};
    const refreshData = () => {
      loadHeatmapData();
    }

    // ✅ 라이프사이클 훅 - 실제 데이터 로드 후 차트 생성
    onMounted(async () => {
      //console.log('PowerLoadAnalysis 컴포넌트 마운트됨, channel:', channel.value);

      if (heatmapChart.value.hasAttribute('_echarts_instance_')) {
        const oldInstance = echarts.getInstanceByDom(heatmapChart.value);
        if (oldInstance) {
          oldInstance.dispose();
          //console.log('기존 히트맵 인스턴스 제거됨');
        }
      }
      
      dualAxisChartInstance = echarts.init(dualAxisChart.value);
      heatmapChartInstance = echarts.init(heatmapChart.value);
      
      try {
        //console.log('API 데이터 로딩 시작...');
        
        await Promise.all([
          loadEnergyHourlyData(channel.value),
          getLoadFactorCalculated(channel.value)
        ]);
        
        //console.log('API 데이터 로딩 완료');
        //console.log('energyHourlyData:', reportData.energyHourlyData);
        //console.log('loadrateData:', reportData.loadrateData);
        
        createDualAxisChart();
        
      } catch (error) {
        console.error('API 데이터 로딩 실패:', error);
        createDualAxisChart();
      }
      
      await loadHeatmapData();
      
      window.addEventListener('resize', () => {
        dualAxisChartInstance?.resize();
        heatmapChartInstance?.resize();
      });
    });

    onUnmounted(() => {
      dualAxisChartInstance?.dispose()
      heatmapChartInstance?.dispose()
      window.removeEventListener('resize', () => {})
    })

    return {
      selectedTimeRange,
      dualAxisChart,
      heatmapChart,
      heatmapLoading,
      averageLoadRate,
      maxLoadRate,
      overloadCount,
      apiHeatmapDistribution,
      refreshData,
      t,
      getLoadFactorCalculated,
      channel,
    }
  }
}
</script>

<style scoped>
/* 차트 영역 */
.dual-axis-chart {
  @apply h-80 w-full;
}

.heatmap-chart {
  @apply h-56 w-full;
}

/* 히트맵 통계 섹션 */
.heatmap-stats-section {
  @apply pt-3 border-t border-gray-200 dark:border-gray-700;
}

.stats-title {
  @apply text-sm font-semibold text-gray-800 dark:text-gray-200 mb-3;
}

.chart-info {
  @apply flex justify-around mt-3 pt-3 border-t border-gray-200 dark:border-gray-700;
}

.info-card {
  @apply text-center;
}

.info-label {
  @apply block text-sm text-gray-500 dark:text-gray-400 mb-1;
}

.info-value {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
}

/* 통계 카드들 - 히트맵 색상 매칭 */
.load-stats {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3;
}

.stat-card {
  @apply bg-gray-50 dark:bg-gray-700/50 p-3 rounded-lg text-center;
  @apply border-l-4;
}

.stat-card.light {
  @apply border-l-8;
  border-left-color: #22c55e;
  @apply bg-gradient-to-br from-green-50 to-emerald-50 dark:from-green-900/20 dark:to-emerald-900/20;
}

.stat-card.medium {
  @apply border-l-8;
  border-left-color: #eab308;
  @apply bg-gradient-to-br from-yellow-50 to-amber-50 dark:from-yellow-900/20 dark:to-amber-900/20;
}

.stat-card.high {
  @apply border-l-8;
  border-left-color: #f97316;
  @apply bg-gradient-to-br from-orange-50 to-red-50 dark:from-orange-900/20 dark:to-red-900/20;
}

.stat-card.overload {
  @apply border-l-8;
  border-left-color: #dc2626;
  @apply bg-gradient-to-br from-red-50 to-pink-50 dark:from-red-900/20 dark:to-pink-900/20;
}

.stat-header {
  @apply flex justify-between items-center mb-2;
}

.stat-card h4 {
  @apply text-gray-700 dark:text-gray-300 text-sm font-semibold;
}

.stat-icon {
  @apply text-xl;
}

.stat-card.light .stat-icon {
  color: #22c55e;
}

.stat-card.medium .stat-icon {
  color: #eab308;
}

.stat-card.high .stat-icon {
  color: #f97316;
}

.stat-card.overload .stat-icon {
  color: #dc2626;
}

.stat-value {
  @apply block text-3xl font-bold mb-2;
  @apply transition-all duration-300;
}

.stat-card.light .stat-value {
  color: #22c55e;
}

.stat-card.medium .stat-value {
  color: #eab308;
}

.stat-card.high .stat-value {
  color: #f97316;
}

.stat-card.overload .stat-value {
  color: #dc2626;
}

.stat-desc {
  @apply text-sm text-gray-500 dark:text-gray-400 font-medium;
}

/* 반응형 */
@media (max-width: 1024px) {
  .dual-axis-chart {
    @apply h-80;
  }
  
  .heatmap-chart {
    @apply h-64;
  }
}

@media (max-width: 640px) {
  .chart-info {
    @apply flex-col gap-3;
  }
  
  .load-stats {
    @apply grid-cols-1;
  }
}
</style>