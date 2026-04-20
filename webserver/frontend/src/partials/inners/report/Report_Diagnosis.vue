<template>
    <div class="col-span-full xl:col-span-12 space-y-4">
  
      <!-- 로딩 표시 -->
      <div v-if="isLoading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
      </div>
  
      <!-- 컨텐츠 -->
      <div v-else class="space-y-4">
        
        <!-- 차트 -->
        <div class="grid grid-cols-12">
          <Diagnosis_Barchart 
            v-if="chartData !== null" 
            :channel="channel" 
            :data="chartData" 
            :mode="mode === 'diagnosis' ? 'DiagnosisDetail' : 'PowerQuality'" 
            :height="300"
            :timestamp="timestamp"
            class="col-span-12" 
          />
        </div>
        
        <!-- 상태 리포트 -->
        <div class="flex flex-col gap-4">
          <template v-for="item in items" :key="item.Item.Name">
            <StatusReport :data="item" />
          </template>
        </div>
        
        <!-- 트렌드 차트 -->
        <div class="grid grid-cols-12 gap-4">
          <div v-for="(option, idx) in chartOptions" :key="idx" class="col-span-6">
            <ReportTrend :data="option" />
          </div>
        </div>
        
        <!-- 데이터 없음 -->
        <div v-if="chartData === null && !isLoading" class="text-gray-500 text-center py-8">
          {{ t('report.noData') }}
        </div>
        
      </div>
  
    </div>
  </template>
  
  <script>
  import { ref, watch, computed, toRef } from 'vue'
  import axios from 'axios'
  import Diagnosis_Barchart from '../../../charts/connect/ReportBatteryChart.vue'
  import StatusReport from './StatusReport.vue'
  import ReportTrend from './ReportTrend.vue'
  import { useI18n } from 'vue-i18n'
  
  export default {
    name: 'Report_Diagnosis',
    props: {
      channel: { type: String, default: '' },
      mode: { type: String, default: 'diagnosis' }, // 'diagnosis' or 'powerquality'
      reportData: { type: Object, default: () => ({ main: [], detail: [], timestamp: null }) }
    },
    components: {
      Diagnosis_Barchart,
      StatusReport,
      ReportTrend,
    },
    setup(props) {
      const { t, locale } = useI18n();

      // === props를 reactive하게 참조 ===
      const timestamp = computed(() => props.reportData?.timestamp || null);
  
      // === 상태 ===
      const isLoading = ref(false);
      const chartData = ref(null);
      const items = ref([]);
      const chartOptions = ref([]);
  
// === 트렌드 차트 데이터 조회 ===
const setChartData = async (itemName, title) => {
  let option = {};

  try {
    const response = await axios.get(`/report/status_trend/${props.mode}/${itemName}`);

    if (response.data.success) {
      const trendData = response.data.data.trend;
      
      if (trendData && trendData.length > 0) {
        const labels = trendData.map(item => item.timestamp);
        const data = trendData.map(item => item.status);

        const datasets = [{
          name: title,
          data: data,
          isThreshold: false,
        }];

        option = { 
          lineLabels: labels, 
          lineData: datasets, 
          lineTitle: title,
          mode: props.mode === 'diagnosis' ? 'DiagnosisDetail' : 'PowerQuality'  // ✅ mode 추가
        };
      } else {
        option = { 
          lineLabels: [], 
          lineData: [], 
          lineTitle: title,
          mode: props.mode === 'diagnosis' ? 'DiagnosisDetail' : 'PowerQuality'  // ✅ mode 추가
        };
      }
    } else {
      option = { 
        lineLabels: [], 
        lineData: [], 
        lineTitle: title,
        mode: props.mode === 'diagnosis' ? 'DiagnosisDetail' : 'PowerQuality'  // ✅ mode 추가
      };
    }
  } catch (error) {
    console.error("트렌드 데이터 요청 실패:", error);
    option = { 
      lineLabels: [], 
      lineData: [], 
      lineTitle: title,
      mode: props.mode === 'diagnosis' ? 'DiagnosisDetail' : 'PowerQuality'  // ✅ mode 추가
    };
  }
  return option;
};
  
      // === 트렌드 차트 로드 ===
      const loadTrendCharts = async (chartList) => {
        const options = [];
        
        // 중복 제거 (item_name 기준)
        const uniqueItems = [];
        const seenNames = new Set();
        
        for (const item of chartList) {
          if (!seenNames.has(item.Name)) {
            seenNames.add(item.Name);
            uniqueItems.push(item);
          }
        }
  
        for (const item of uniqueItems) {
          const titleName = item.Title || item.Name;
          const chartValue = await setChartData(item.Name, titleName);
          if (chartValue.lineLabels.length > 0) {
            options.push(chartValue);
          }
        }
  
        return options;
      };
  
      // === InfluxDB 데이터 변환 ===
      const transformInfluxData = (main, detail) => {
        let itemlist = [], valuelist = [], datalist = [];
        for (let i = 0; i < main.length; i++) {
          const titleKey = `title_${locale.value}`;
          const displayName = main[i][titleKey] || main[i].title || main[i].item_name;
          itemlist.push(displayName);
          valuelist.push(main[i].status);
          datalist.push({
            en: main[i].title_en || main[i].item_name,
            ko: main[i].title_ko || main[i].item_name,
            ja: main[i].title_ja || main[i].item_name
          });
        }
        const chartDataResult = { "Names": itemlist, "Values": valuelist, "Titles": datalist };
  
        const mainByName = {};
        for (let i = 0; i < main.length; i++) {
          const key = main[i].item_name.replace(/\s/g, '');
          mainByName[key] = main[i];
        }
  
        const groupedByParent = {};
        for (let i = 0; i < detail.length; i++) {
          const parentName = detail[i].parent_name;
          const parentKey = parentName.replace(/\s/g, '');
          
          if (!groupedByParent[parentName]) {
            const parentInfo = mainByName[parentKey] || {};
            groupedByParent[parentName] = {
              item_name: parentName,
              status: 0,
              Titles: {
                en: parentInfo.title_en || parentName,
                ko: parentInfo.title_ko || parentName,
                ja: parentInfo.title_ja || parentName
              },
              Descriptions: {
                en: parentInfo.description_en || '',
                ko: parentInfo.description_ko || '',
                ja: parentInfo.description_ja || ''
              },
              children: []
            };
          }
          if (detail[i].status > groupedByParent[parentName].status) {
            groupedByParent[parentName].status = detail[i].status;
          }
          groupedByParent[parentName].children.push(detail[i]);
        }
  
        const itemsResult = [], chartList = [];
        for (const parentName in groupedByParent) {
          const parent = groupedByParent[parentName];
          if (parent.status > 1) {
            const childDict = [];
            for (let j = 0; j < parent.children.length; j++) {
              const child = parent.children[j];
              if (child.status > 1) {
                const childTitle = child[`title_${locale.value}`] || child.title || child.item_name;
                childDict.push({
                  Title: childTitle,
                  Assembly: child.assembly_id,
                  Value: child.value !== undefined ? child.value : 'NaN'
                });
              }
            }
            if (childDict.length > 0) {
              // main의 item_name으로 트렌드 조회 (띄어쓰기 제거)
              chartList.push({
                Name: parentName.replace(/\s/g, ''),
                Title: parent.Titles[locale.value],
                Status: parent.status
              });
  
              itemsResult.push({
                Item: {
                  Name: parentName,
                  Title: parent.Titles[locale.value],
                  Titles: parent.Titles,
                  Descriptions: parent.Descriptions,
                  Status: parent.status
                },
                Child: childDict
              });
            }
          }
        }
  
        return { chartData: chartDataResult, items: itemsResult, chartList };
      };
  
      // === 데이터 처리 ===
      const processReportData = async () => {
        const { main, detail } = props.reportData;
        
        if (!main || main.length === 0) {
          chartData.value = null;
          items.value = [];
          chartOptions.value = [];
          return;
        }
  
        isLoading.value = true;
  
        const { chartData: cd, items: it, chartList } = transformInfluxData(main, detail);
        
        chartData.value = cd;
        items.value = it;
        chartOptions.value = await loadTrendCharts(chartList);
  
        isLoading.value = false;
      };
  
      // === Watch reportData ===
      watch(() => props.reportData, async (newVal) => {
        if (newVal && newVal.main && newVal.main.length > 0) {
          await processReportData();
        } else {
          chartData.value = null;
          items.value = [];
          chartOptions.value = [];
        }
      }, { deep: true, immediate: true });
  
      return {
        t,
        isLoading,
        chartData,
        items,
        chartOptions,
        timestamp,
        channel: props.channel,
        mode: props.mode,
      };
    }
  }
  </script>
  