<template>
  <EventTable :channel="channel" :mode="mode" />

    <!-- Pagination -->
    <div class="mt-8">
        <PaginationClassic v-if="totalPages > 1" :channel="channel" :mode="mode" />
    </div>
</template>
  <script>
import EventTable from './EventTable.vue'
import PaginationClassic from '../../../pages/common/PaginationClassic.vue'
import { ref, provide, onMounted, watch } from 'vue';
import axios  from 'axios'
export default {
  name: 'EventTab',
  props:['channel','mode'],
  components: {
    EventTable,
    PaginationClassic,
  },
  setup(props){
    const channel = ref(props.channel);
    const mode = ref(props.mode);
    const eventData = ref([]);
    const alarmData = ref([]);
    const totalRecord = ref(0);
    const totalPages = ref(1);
    const curPage = ref(1);
    const param = ref(-1);
    const start = ref('');
    const end = ref('');
    const search = ref(false);
    const period = ref(false);
    const parameter = ref(false);

    const eventType = [
          { id: 0, label: "All" },
          { id: 1, label: "SAG" },
          { id: 2, label: "SWELL" },
          { id: 3, label: "SHORT INTERRUPT" },
          { id: 4, label: "LONG INTERRUPT" },
          { id: 5, label: "OVER CURRENT" },
          { id: 6, label: "UNDER CURRENT" },
          { id: 7, label: "VOLTAGE TRANSIENT" },
          { id: 8, label: "CURRENT TRANSIENT" },
        ];

    const parseMask = (mask) => {
        const phases = [];
        if (mask & 0b001) phases.push("L1");
        if (mask & 0b010) phases.push("L2");
        if (mask & 0b100) phases.push("L3");
        return phases;
    }

    const fetchEventData = async (ch, pg) => {
      try {
        const response = await axios.get(`/api/getEventlist/${ch}/${pg}`);
        if(response.data.success){
          const tmpList = response.data.data;
          eventData.value = [];
          for (let i = 0 ; i < tmpList.length;i++){
            const phase = parseMask(tmpList[i]["mask"]).join(" ");
            const levels = [tmpList[i]["level_l1"],tmpList[i]["level_l2"],tmpList[i]["level_l3"]]
            const levelstr = [
                `L1:${levels[0].toFixed(2)}`,
                `L2:${levels[1].toFixed(2)}`,
                `L3:${levels[2].toFixed(2)}`
            ].join(", ");
            //const typeIdx = parseInt(tmpList[i]["type"]);
            const endtime = tmpList[i]["end_ts"]
            eventData.value.push({"Type":tmpList[i]["event_type"] ,"StartTime":tmpList[i]["start_ts"], "Duration":tmpList[i]["duration"],"EndTime":endtime,"Phase":phase, "Level":levelstr});
          }
          totalRecord.value = response.data.totalRecord;
          totalPages.value = response.data.totalPages;
        }else{
          //alert('No Data');
        }
      } catch (error) {
        console.error("데이터를 가져오는 중 오류 발생:", error);
      }
    };

    function formatTimestampWithMs(timestamp) {
    const date = new Date(timestamp);
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    const milliseconds = String(date.getMilliseconds()).padStart(3, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
}

    const searchEventData = async (ch, pg) => {
      const searchRequest = {
        param: param.value.toString(),
        StartDate: start.value,
        EndDate: end.value,
      };
      //console.log(searchRequest);
      try {
        const response = await axios.post(`/api/getEventSearch/${ch}/${pg}`, searchRequest, {
          headers: { "Content-Type": "application/json" },
        });
        if(response.data.success){
          const tmpList = response.data.data;
          eventData.value = [];
          for (let i = 0 ; i < tmpList.length;i++){
            const phase = parseMask(tmpList[i]["mask"]).join(" ");
            const levels = [tmpList[i]["level_l1"],tmpList[i]["level_l2"],tmpList[i]["level_l3"]]
            const levelstr = [
                `L1:${levels[0].toFixed(2)}`,
                `L2:${levels[1].toFixed(2)}`,
                `L3:${levels[2].toFixed(2)}`
            ].join(", ");
            //const typeIdx = parseInt(tmpList[i]["type"]);
            const endtime = tmpList[i]["end_ts"];
            eventData.value.push({"Type":tmpList[i]["event_type"] ,"StartTime":tmpList[i]["start_ts"], "Duration":tmpList[i]["duration"],"EndTime":endtime,"Phase":phase, "Level":levelstr});
          }
          totalRecord.value = response.data.totalRecord;
          totalPages.value = response.data.totalPages;
          curPage.value = pg;
        }else{
          alert('No Data');
        }
      } catch (error) {
        console.error("데이터를 가져오는 중 오류 발생:", error);
      }
    };

    const fetchAlarmData = async (ch, pg) => {
      try {
        //const response = await axios.get(`/api/getAlarmLog/${ch}`);
        const response = await axios.get(`/api/getAlarms/${ch}/${pg}`);
        if(response.data.success){
          eventData.value = response.data.data;
          totalRecord.value = response.data.totalRecord;
          totalPages.value = response.data.totalPages;
        }else{
          //alert('No Data');
        }
      } catch (error) {
        console.error("데이터를 가져오는 중 오류 발생:", error);
      }
    };

    const searchAlarmData = async (ch, pg) => {
      const searchRequest = {
        param: param.value,
        StartDate: start.value,
        EndDate: end.value,
      };
      try {
        const response = await axios.post(`/api/getAlarmSearch/${ch}/${pg}`, searchRequest, {
          headers: { "Content-Type": "application/json" },
        });
        if(response.data.success){
          eventData.value = response.data.data;
          totalRecord.value = response.data.totalRecord;
          totalPages.value = response.data.totalPages;
          curPage.value = pg;
        }else{
          alert('No Data');
        }
      } catch (error) {
        console.error("데이터를 가져오는 중 오류 발생:", error);
      }
    };

    onMounted(()=>{
      if(mode.value == 'Event')
        fetchEventData(channel.value, curPage.value);
      else
        fetchAlarmData(channel.value, curPage.value);
    });

    provide('eventData',eventData);
    provide('curPage',curPage);
    provide('param',param)
    provide('start',start)
    provide('end',end)
    provide('search',search)
    provide('totalPages',totalPages);
    provide('totalRecord',totalRecord);
    provide('period',period);
    provide('parameter',parameter);

    watch(() => curPage.value, (newPage) => {
      if(mode.value != 'Event'){
        if(!period.value && !parameter.value){
          fetchAlarmData(channel.value, newPage)
        }else{
          searchAlarmData(channel.value, newPage);
        }
      }else{
        if(!period.value && !parameter.value){
          fetchEventData(channel.value, newPage)
        }else{
          searchEventData(channel.value, newPage);
        }
      }
    });

    watch(() => search.value, (val) => {
      //fetchAlarmData(channel.value, newPage)
      if(val){
        //console.log('SEARCH:' ,start.value , end.value, param.value);
        if(mode.value != 'Event'){
          if(!period.value && !parameter.value){
            fetchAlarmData(channel.value, 1);
            curPage.value = 1;
          }
          else{
            searchAlarmData(channel.value, 1);
            curPage.value = 1;
          }
          search.value = false;   
        }else{
          if(!period.value && !parameter.value){
            fetchEventData(channel.value, 1);
            curPage.value = 1;
          }
          else{
            searchEventData(channel.value, 1);
            curPage.value = 1;
          }
          search.value = false;   
        }    
      }
    });

    return{
      channel,
      eventData,
      alarmData,
      mode,
      totalRecord,
      totalPages,
      curPage,
      param,
      start,
      end,
      search,
      period,
      parameter,
    }
  }
}
</script>