<template>
  <div class="grid grid-cols-12 gap-x-2 gap-y-6">
    <!-- 메인 채널 카드 - 6칸 -->
    <Dashboard_MeterCard 
      v-if="channelState.MainEnable" 
      :channel="'main'" 
      class="col-span-6"
    />

    <!-- 진단 현황 카드 - 6칸, 세로 2줄 차지 -->
    <DashboardCard_Meter_Integrated v-if="channelState.SubEnable"
      class="col-span-6 row-span-2"
      :channel="'sub'"
      :diagData="diagData_sub"
    />

    <!-- 서브 채널 카드 - 6칸 -->
    <Dashboard_MeterCard 
      v-if="channelState.SubEnable" 
      :channel="'sub'" 
      class="col-span-6"
    />
  </div>
</template>

<script>
  import { ref } from 'vue';
import Dashboard_MeterCard from "../../partials/inners/dashboard/DashboardCard_Meter_Ch3.vue";
import DashboardCard04 from "../../partials/inners/dashboard/DashboardCard_New.vue";
import DashboardCard_Meter_Integrated from '../../partials/inners/dashboard/DashboardCard_Meter_Integrated_total.vue'
import { useI18n } from 'vue-i18n'
export default {
  name: "DualChannelLayout",
  components: {
    Dashboard_MeterCard,
    DashboardCard04,
    DashboardCard_Meter_Integrated,
  },
  props: {
    channelState: {
      type: Object,
      required: true,
    },
  },
    setup(props){
    const { t } = useI18n()

    const diagData_main = ref(null);
    const diagData_sub = ref(null);
    const mainDataReady = ref(false);
    const subDataReady = ref(false);
    const isLoading = ref(false);

    return {
      diagData_main,
      diagData_sub,
      mainDataReady,
      subDataReady,
      isLoading,
      t,
    }
  }
};
</script>
