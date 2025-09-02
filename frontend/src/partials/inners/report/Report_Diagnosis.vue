<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl mt-4">
              <div class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
            >
              <div
                class="absolute top-0 left-0 right-0 h-0.5 bg-sky-500"
                aria-hidden="true"
              ></div>
              <div
                class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
              >
                <header class="flex items-center mb-2">
                  <div class="w-6 h-6 rounded-full shrink-0 bg-sky-500 mr-3 flex items-center justify-center">
                    <svg 
                      xmlns="http://www.w3.org/2000/svg" 
                      class="w-5 h-5 text-white"
                      viewBox="0 0 24 24" 
                      fill="none"
                      stroke="currentColor"
                      stroke-linecap="round" 
                      stroke-linejoin="round"
                      stroke-width="2"
                    > 
                      <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14z"/> 
                      <path d="M3 10h18"/> 
                      <path d="M10 3v18"/> 
                    </svg>
                  </div>

                  <h3
                    class="text-lg text-gray-800 dark:text-gray-100 font-semibold"
                  >
                  {{ t('report.cardContext.info.status') }}
                  </h3>
                </header>
              </div>
              <div class="grid grid-cols-12 p-2">
                  <Diagnosis_Barchart v-if="chartdata !== null" :channel="channel" :data="chartdata" :mode="'Status'" class="h-auto" />
              </div>
              <div class="flex flex-col gap-4 p-2 pl-4 pr-4">
                <template v-for="item in items">
                  <StatusReport :data="item" />
                </template>
              </div>
              <div class="grid grid-cols-12 gap-4 p-2 pl-4 pr-4">
                <div v-for="option in chartOption" class="col-span-6">
                  <ReportTrend :data="option" />
                </div>
                <div class="col-span-6">
                  <ReportEventFault v-if="eventData.length > 0" :data="eventData" />
                </div>
              </div>
              
      </div>
  </div>
</template>

<script>
import { ref, watch, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSetupStore } from '@/store/setup'; // ✅ Pinia Store 사용
import axios from 'axios'
import DashboardCard_table3 from './DashboardCard_table4.vue'
import Diagnosis_Barchart from '../diagnosis/Diagnosis_BarChart2.vue'
import Diagnosis_Info from '../diagnosis/Diagnosis_Info2.vue'
import StatusReport from './StatusReport.vue'
import ReportTrend from './ReportTrend.vue'
import { useI18n } from 'vue-i18n'  // ✅ 추가
//import Diagnosis_TreeTable from './Diagnosis_TreeTable2.vue'
import ReportEventFault from './ReportEventFault.vue'
import { useInputDict } from "@/composables/useInputDict";
//import TreeRow from './TreeRow.vue';

export default {
  name: 'Report_Diagnosis',
  props:['channel'],
  components: {
    DashboardCard_table3,
    Diagnosis_Barchart,
    Diagnosis_Info,
    StatusReport,
    ReportTrend,
    //Diagnosis_TreeTable,
    ReportEventFault,
    //TreeRow,
  },
  setup(props) {
    const { t, locale } = useI18n();
    const route = useRoute()
    const channel = ref(props.channel)
    const setupStore = useSetupStore();
    //const langset = computed(() => authStore.getLang);
    const assetdata = ref([]);
    const rawdata = ref([]);
    const chartdata = ref(null);
    const items = ref([]);
    const asset = computed(() => setupStore.getAssetConfig);
    const chartOption = ref([]);
    const eventData = ref([]);
    const today = ref(new Date());
    const weekAgo = computed(() => {
      const date = new Date()
      date.setDate(date.getDate() - 7)
      return date
    });
    const {
     formatToISOString,
    } = useInputDict();

    const setParamIds = async(asset, datalist) =>{
      const idList = [];
      try{
        const response = await axios.get(`/api/getParameters/${asset}/diagnostic`);
        if (response.data.success) {
          const paramData = response.data.data;
          //console.log('datalist',datalist);
          //console.log('paramData', paramData);
          for(let i = 0 ; i < paramData.length ; i++){
            for(let j = 0; j < datalist.length; j++){
              if(paramData[i]["Name"] == datalist[j].Name && paramData[i]["AssemblyID"] == datalist[j].AssemblyID ){
                idList.push({idx:paramData[i]["ID"], Assembly:paramData[i]["AssemblyID"], title:paramData[i]["Title"]});
              }
            }
          }
        }
      }catch (error){
        console.error(error);
      }
      return idList;
    }

    const setChartData = async (effectiveIds, Title) => {
      let option = {};

      const trendDataRequest = {
        ParametersIds: [effectiveIds],
        StartDate: formatToISOString(weekAgo.value, 0),
        EndDate: formatToISOString(today.value, 1),
      };

      try {
        const url = `/api/getTrendData`; // FastAPI 엔드포인트
        const response = await axios.post(url, trendDataRequest, {
          headers: { "Content-Type": "application/json" },
        });

        if (response.data.success) {
          const resData = response.data.data;
          let datasets = [];
          let labels = [];

          // 1. 메인 데이터 처리 (Thresholds 제외)
          Object.keys(resData).forEach((key) => {
            if (key !== "Thresholds") {
              const dataPoints = resData[key].data;
              if (dataPoints && dataPoints.length > 0) {
                if (labels.length === 0) {
                  labels = dataPoints.map((point) => point.XAxis);
                }
                datasets.push({
                  name: resData[key].Title, // 예: "Current RMS A"
                  data: dataPoints.map((point) => point.YAxis),
                  isThreshold: false, // 일반 데이터
                });
              }
            }
          });

          // 2. Thresholds 처리 (Thresholds 배열이 2개일 때)
          if (
            resData.Thresholds &&
            resData.Thresholds.length == 2 &&
            labels.length > 0
          ) {
            // t1: 첫 번째 Threshold의 XAxis, t2: 두 번째 Threshold의 XAxis
            let timeList = []
            for(let i = 0 ; i < resData.Thresholds.length;i++){
              timeList.push(new Date(resData.Thresholds[i].XAxis));
            }
            const t1 = new Date(resData.Thresholds[0].XAxis);
            const t2 = new Date(resData.Thresholds[1].XAxis);

            // 첫 번째 Threshold 리스트의 값 (예: 인덱스 4, 5, 6 등)
            resData.Thresholds[0].Thresholds.forEach((value, idx) => {
              if (value !== "NaN" && value !== null && value !== undefined) {
                // 두 번째 Threshold의 값
                const secondValue = resData.Thresholds[1].Thresholds[idx];
                if (
                  secondValue === "NaN" ||
                  secondValue === null ||
                  secondValue === undefined
                ) {
                  // 두 번째 값이 유효하지 않으면 스킵
                  return;
                }

                const thresholdData = labels.map((lbl) => {
                  const dt = new Date(lbl);
                  return dt < t1 ? value : secondValue;
                });
                const ThresholdString = [
                  "Out of Range(Down side)",
                  "Repair",
                  "Inspect",
                  "Warning",
                  "Warning",
                  "Inspect",
                  "Repair",
                  "Out of Range(Upper side)",
                ];
                datasets.push({
                  name: ThresholdString[idx], // 예: "4", "5", "6" 등 (해당 인덱스)
                  data: thresholdData,
                  isThreshold: true, // threshold series임을 표시 (LineChart 컴포넌트에서 점선 처리)
                });
              }
            });
          }

          // 3. 차트 업데이트
          option = {
            lineLabels: labels,
            lineData: datasets,
            lineTitle : Title,
          };
        } else {
          console.error("서버 오류:", response.data.error);
          option = {
            lineLabels: [],
            lineData: [],
            lineTitle : '',
          };
        }
      } catch (error) {
        console.error("요청 실패:", error);
      }
      return option;
    };

    const fetchDetailData = async () => {
      items.value = [];
      chartOption.value = [];
      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      try {

        const response = await axios.get(`/api/getDiagnosisDetail/${chName}`);

        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
       
            itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"])
          }
          chartdata.value = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};

          const treeItems = response.data.data_tree;
          const chartList = [];
          for(let i = 0 ; i< treeItems.length;i++){
            if(treeItems[i]["Status"] > 1 ){
              //items.value.push();
              const childDict = [];
              if(treeItems[i]["children"].length > 0){
                for(let j = 0 ; j < treeItems[i]["children"].length ; j++ ){
                  if(treeItems[i]["children"][j]["isSub"]){
                    for(let k = 0 ; k < treeItems[i]["children"][j]["children"].length ; k++ ){
                      if( treeItems[i]["children"][j]["children"][k]["Status"] > 1){
                        chartList.push(treeItems[i]["children"][j]["children"][k]);
                        childDict.push({Title:treeItems[i]["children"][j]["children"][k]["Title"], Assembly:treeItems[i]["children"][j]["AssemblyID"], Value:treeItems[i]["children"][j]["children"][k]["Value"]})
                      }
                    }
                  }else{
                    if(treeItems[i]["children"][j]["Status"] > 1){
                      chartList.push(treeItems[i]["children"][j]);
                      childDict.push({Title:treeItems[i]["children"][j]["Title"],Assembly:treeItems[i]["children"][j]["AssemblyID"], Value:treeItems[i]["children"][j]["Value"]})
                    }
                  }
                }
              }else{
                childDict = [{Title:treeItems[i]["Title"],Assembly:treeItems[i]["children"][j]["AssemblyID"], Value:treeItems[i]["Value"]} ]
              }
              items.value.push({Item:treeItems[i], Child:childDict});            
            }
          }
          // for(let i = 0 ; i< treeItems.length;i++){
          //   if(treeItems[i]["Status"] > 1 ){
          //     //items.value.push();
          //     const childDict = [];
          //     for(let j = 0 ; j < treeItems[i]["children"].length ; j++ ){
          //       if(treeItems[i]["children"][j]["Status"] > 1){
          //         chartList.push(treeItems[i]["children"][j]);
          //         childDict.push({Title:treeItems[i]["children"][j]["Title"],Assembly:treeItems[i]["children"][j]["AssemblyID"] , Value:treeItems[i]["children"][j]["Value"]})
          //       }
          //     }
          //     items.value.push({Item:treeItems[i], Child:childDict});            
          //   }
          // }
          if(chartList.length > 0){
            let chartValue=null;
            const effectiveIds = await setParamIds(chName, chartList)
            let idxList=[],idList=[];
            for (let i = 0 ; i < effectiveIds.length;i++){
              if (idxList.includes(effectiveIds[i].idx)){
                continue;
              }else{
                idList.push(effectiveIds[i])
                idxList.push(effectiveIds[i].idx)
              }
            }
            //console.log(idxList);
            for (let i = 0 ; i < idList.length;i++){
                const titleName = '['+idList[i].Assembly + ']' + idList[i].title;
                chartValue = await setChartData(idList[i].idx, titleName);
                chartOption.value.push(chartValue);
            }
            //console.log('Diagnosis',idxList);
          }

        }else{
          console.log('No Data');
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const fetchEventData = async () => {

      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      try {
        //const ch = 'Fan';
        const response = await axios.get(`/api/getEvents/${chName}`);
        if (response.data.success) {
          const eventList = response.data.data_tree;
          for(let i = 0; i < eventList.length ; i++){
            if(eventList[i]["Status"] > 1){
              eventData.value.push(eventList[i]);
            }
          }
        }else{
          console.log('No Data');
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }

      try {

        const response = await axios.get(`/api/getFaults/${chName}`);
        if (response.data.success) {
          const faultList = response.data.data_tree;
          for(let i = 0; i < faultList.length ; i++){
            if(faultList[i]["Status"] > 1){
              eventData.value.push(faultList[i]);
            }
          }
        }else{
          console.log('No Data');
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    onMounted(async () => {

      await fetchDetailData();
      await fetchEventData();
    });

    watch(() => route.params.channel, async(newChannel) => {
      channel.value = newChannel
      eventData.value = [];
      await fetchDetailData();
      await fetchEventData();

    });

//     const items = ref([
//   {
//     id: 1,
//     name: '상위 항목 1',
//     status: 'OK',
//     children: [
//       { id: 2, name: '하위 항목 1', status: 'Warning' },
//       {
//         id: 3,
//         name: '하위 항목 2',
//         status: 'Repair',
//         children: [{ id: 4, name: '세부 항목', status: 'Inspect' }],
//       },
//     ],
//   },
//   { id: 5, name: '상위 항목 2', status: 'NoData' },
// ])

    return {
      channel,
      //langset,
      assetdata,
      rawdata,
      asset,
      items,
      chartdata,
      fetchDetailData,
      setChartData,
      locale,
      formatToISOString,
      t,
      chartOption,
      setParamIds,
      eventData,
      today,
      weekAgo,
    }  
  }
}
</script>