<!--template>
    <div class="grid grid-cols-12 gap-6">
      
      <Dashboard_Gems_Meter v-if="channelState.MainEnable" 
      :channel="'main'" 
      :data="mainData" />
      <Dashboard_TransInfo_final v-if="opMode === 'device2' && channelState.MainDiagnosis"
        :channel="channel" 
        :data="mainData" 
      />
      
      <Gems_Status v-else-if="opMode !== 'device2' && channelState.MainDiagnosis"
        :channel="channel" 
        :data="mainData" 
      />
      
      <DashboardCard_kwh v-if="Object.keys(mainData).length > 0 && channelState.MainEnable"
        :channel="channel" 
      />
      
      <Gems_Event v-if="channelState.MainDiagnosis"
        :channel="channel" 
      />
    </div>
  </template-->
<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 왼쪽 메인 영역을 감싸는 div -->
    <div class="col-span-7">
      <div class="grid grid-cols-7 gap-6">
        <Dashboard_Gems_Meter v-if="channelState.MainEnable" 
          :channel="'main'" 
          :data="mainData" />
        
        <DashboardCard_kwh v-if="Object.keys(mainData).length > 0 && channelState.MainEnable"
          :channel="channel" 
        />
      </div>
    </div>
    
    <!-- 오른쪽 사이드 패널 -->
    <div class="col-span-5">
      <Gems_Status v-if="channelState.MainDiagnosis"
        :channel="channel" 
        :data="mainData" 
      />
    </div>
  </div>
</template>
  <script>
  import DashboardCard_Meter_Claude from '../../partials/inners/dashboard/DashboardCard_Meter_Ch1.vue'
  import DashboardCard_PQ_Claude from '../../partials/inners/dashboard/DashboardCard_PQ_Claude2.vue'
  import Dashboard_TransInfo_final from '../../partials/inners/dashboard/Dashboard_TransInfo_final.vue'
  import Dashboard_Single_Info from '../../partials/inners/dashboard/Dashboard_Single_Info.vue'
  import Gems_Status from '../../partials/inners/dashboard/Gems_EventStatus.vue'
  import Gems_Event from '../../partials/inners/dashboard/Gems_Event.vue'
  import DashboardCard_kwh from '../../partials/inners/dashboard/DashboardCard_kwh_realtime.vue'
  import DashboardCard_Diagnosis from '../../partials/inners/dashboard/DashboardCard_Diagnosis.vue'
  import Dashboard_Gems_Meter from '../../partials/inners/dashboard/Dashboard_Gems_Meter.vue'
  import { useAuthStore } from '@/store/auth'
  import { computed } from 'vue'
  
  export default {
    name: 'GemsLayout2',
    components: {
      DashboardCard_Meter_Claude,
      DashboardCard_PQ_Claude,
      Dashboard_TransInfo_final,
      DashboardCard_kwh,
      DashboardCard_Diagnosis,
      Dashboard_Single_Info,
      Gems_Status,
      Gems_Event,
      Dashboard_Gems_Meter
    },
    props: {
      mainData: {
        type: Object,
        required: true
      },
      subData: {
        type: Object,
        default: () => ({})
      },
      channelState: {
        type: Object,
        required: true
      },
      channel: {
        type: String,
        required: true
      }
    },
    setup() {
      const authStore = useAuthStore()
      const opMode = computed(() => authStore.getOpMode);
      
      return {
        opMode
      }
    }
  }
  </script>