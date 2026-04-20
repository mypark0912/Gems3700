<template>
<div v-if="mode != 'Event'" class="grid grid-cols-12 gap-6">
    <!-- 왼쪽 박스 (Diagnosis_Info + Diagnosis_Barchart) -->
     <div class="col-span-12 flex flex-col gap-6 h-auto"> 
        <Diagnosis_Barchart v-if="chartdata !== null" :channel="channel" :data="chartdata" :mode="mode" class="h-auto" />
      </div>

          <!-- 오른쪽 박스 (Diagnosis_TreeTable) - 높이 제한 없음 -->
      <div class="col-span-12">
        <Diagnosis_TreeTable v-if="items.length > 0" :channel="channel" :data="items" :mode="mode" />
      </div>
</div>
<div v-else class="grid grid-cols-12 gap-6">
    <!-- 왼쪽 박스 (Diagnosis_Info + Diagnosis_Barchart) -->
     <div class="col-span-8 flex flex-col gap-6 h-auto"> 
        <Diagnosis_Barchart v-if="chartdata !== null" :channel="channel" :data="chartdata" :mode="mode" class="h-auto" />
      </div>

          <!-- 오른쪽 박스 (Diagnosis_TreeTable) - 높이 제한 없음 -->
      <div class="col-span-4 pt-4">
        <Diagnosis_TreeTable v-if="items.length > 0" :channel="channel" :data="items" :mode="mode" />
      </div>
</div>
</template>
<script>
import { ref, watch, computed, onMounted, provide } from 'vue'
import { useI18n } from 'vue-i18n'  // ✅ 추가
//import { useAuthStore } from '@/store'; // ✅ Pinia Store 사용
import axios from 'axios'
import Diagnosis_TreeTable from './Diagnosis_TreeTable.vue'
import Diagnosis_Barchart from './Diagnosis_BarChart.vue'
import Transformer from './Transformer.vue'

export default {
  name: 'Diagnosis',
  props:{
    asset:{
      type: Object,
      default: () => {}
    },
    channel:{
      type:String,
      default: ''
    },
    mode:{
      type:String,
      default: ''
    }
  },
  components: {
    Diagnosis_TreeTable,
    Diagnosis_Barchart,
    Transformer,
  },
  setup(props) {
    const { t, locale } = useI18n();
    const channel = ref(props.channel)
    //const authStore = useAuthStore();
    const asset = ref(props.asset);
    const chartdata = ref(null);
    const data_recordtime = ref('');
    const data_state = ref('');
    const items = ref([]);
    const mode=ref(props.mode);

    const fetchDetailData = async () => {

      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      try {
        //const ch = 'Fan';
        //console.log('locale...............=',locale.value);
        //const response = await axios.get(`/api/getDiagnosisDetail/${chName}`);
        const response = await axios.get(`/api/getDiagnosisDetail/${chName}`);
        //console.log('Title...........='+response.data.data_status[0]["Titles"][locale.value])
        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
       
            itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"])
            // datalist.push({"Titles" : response.data.data_status[i]["Titles"][locale.value], "Status" : response.data.data_status[i]["Status"], 
            // "Description" : response.data.data_status[i]["Description"], "Descriptions" : response.data.data_status[i]["Descriptions"]})
          }
          chartdata.value = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};
          items.value = response.data.data_tree;
          data_recordtime.value = response.data.data_recordtime;
          data_state.value = response.data.data_state;
        }else{
          console.log('No Data');
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

const fetchPQData = async () => {
  const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
  
  // 원하는 순서 정의 (Name 기준)
  const pqDisplayOrder = [
    // 전압 관련
    'VoltagePhaseAngle',
    'VoltageRMS',
    'DC',
    // 전류 관련
    'CurrentRMS',
    'CurrentPhaseAngle',
    'CrestFactor',
    // 파형, 왜곡 관련
    'Unbalance',
    'Harmonics',
    'ZeroSequence',
    'NegativeSequence',
    // 전력, 효율 관련
    'Power',
    'PowerFactor',
    'TotalDemandDistortion',
    // 기타
    'PhaseAngle',
    'Events',
  ];

  try {
    const response = await axios.get(`/api/getDiagPQ/${chName}`);

    if (response.data.success) {
      // data_status 정렬
      const sortedStatus = [...response.data.data_status].sort((a, b) => {
        const indexA = pqDisplayOrder.indexOf(a.Name);
        const indexB = pqDisplayOrder.indexOf(b.Name);
        if (indexA === -1) return 1;
        if (indexB === -1) return -1;
        return indexA - indexB;
      });

      // data_tree도 같은 순서로 정렬
      const sortedTree = [...response.data.data_tree].sort((a, b) => {
        const indexA = pqDisplayOrder.indexOf(a.Name);
        const indexB = pqDisplayOrder.indexOf(b.Name);
        if (indexA === -1) return 1;
        if (indexB === -1) return -1;
        return indexA - indexB;
      });

      let itemlist = [], valuelist = [], datalist = [];
      for (let i = 0; i < sortedStatus.length; i++) {
        itemlist.push(sortedStatus[i]["Titles"][locale.value]);
        valuelist.push(sortedStatus[i]["Status"]);
        datalist.push(sortedStatus[i]["Titles"]);
      }
      
      chartdata.value = { "Names": itemlist, "Values": valuelist, "Titles": datalist };
      items.value = sortedTree;  // 정렬된 tree 사용
      data_recordtime.value = response.data.data_recordtime;
      data_state.value = response.data.data_state;
    } else {
      console.log('No Data');
    }
  } catch (error) {
    console.log("데이터 가져오기 실패:", error);
  }
};
    const fetchFaultData = async () => {

      const chName = channel.value == 'Main'? asset.value.assetName_main : asset.value.assetName_sub;
      try {
        //const ch = 'Fan';
        //const response = await axios.get(`/api/getFaults/${chName}`);
        const response = await axios.get(`/api/getFaults/${chName}`);
        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
            //itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            itemlist.push(response.data.data_status[i]["Title"]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"]);
            // datalist.push({"Item" : response.data.data_status[i]["Title"], "Status" : response.data.data_status[i]["Status"], 
            // "Description" : response.data.data_status[i]["Description"], "Descriptions" : response.data.data_status[i]["Descriptions"]})
          }
          chartdata.value = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};
          items.value = response.data.data_tree;
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
        //const response = await axios.get(`/api/getEvents/${chName}`);
        const response = await axios.get(`/api/getEvents/${chName}`);
        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
            itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"]);
            // datalist.push({"Item" : response.data.data_status[i]["Titles"][locale.value], "Status" : response.data.data_status[i]["Status"], "Description" : response.data.data_status[i]["Description"]})
          }
          chartdata.value = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};
          items.value = response.data.data_tree;
        }else{
          console.log('No Data');
        }
      }catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };


    onMounted(async () => {
     
        if(mode.value == 'Status')
          await fetchDetailData();
        else if(mode.value == 'PowerQuality')
          await fetchPQData();
        else if(mode.value == 'Fault')
          await fetchFaultData();
        else
          await fetchEventData();
        
    });

    watch(() => props.asset, async(newChannel) => {
      if (newChannel !== asset.value) {
        asset.value = newChannel;
        await fetchDetailData();
      }
    });

    provide('data_recordtime', computed(() => data_recordtime.value));
    provide('data_state', computed(() => data_state.value));

    return {
      channel,
      //langset,
      chartdata,
      items,
      asset,
      fetchDetailData,
      mode,
      data_recordtime,
      data_state,
    }  
  }



}
</script>