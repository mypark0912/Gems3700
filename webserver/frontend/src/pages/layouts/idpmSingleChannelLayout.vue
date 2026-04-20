<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 상단 카드들 -->
    <DashboardCard_idpm300_meter />
    <Dashboard_idpm300_power />
    <DashboardCard_PQ_Claude />
    <div class="col-span-full grid grid-cols-12 gap-6 items-stretch">
      <!-- 하단 좌측: 실시간 kWh 차트 (더미) -->
      <div class="col-span-full lg:col-span-7 flex">
        <DashboardCard_kwh_realtime_dummy class="w-full" />
      </div>

      <!-- 하단 우측: 실시간 알람/이벤트 -->
      <div class="col-span-full lg:col-span-5 flex">
        <DashboardCard_alarm_status class="w-full" />
      </div>
    </div>
  </div>
</template>

<script>
import DashboardCard_idpm300_meter from "../../partials/inners/dashboard/DashboardCard_idpm300_meter.vue";
import Dashboard_idpm300_power from "../../partials/inners/dashboard/Dashbaord_idpm300_power.vue";
import DashboardCard_PQ_Claude from "../../partials/inners/dashboard/DashboardCard_PQ_Claude2.vue";
import DashboardCard_kwh_realtime_dummy from "../../partials/inners/dashboard/DashboardCard_kwh_realtime_dummy.vue";
import DashboardCard_alarm_status from "../../partials/inners/dashboard/DashboardCard_alarm_status.vue";

export default {
  name: "SingleChannelLayout",
  components: {
    DashboardCard_idpm300_meter,
    Dashboard_idpm300_power,
    DashboardCard_PQ_Claude,
    DashboardCard_kwh_realtime_dummy,
    DashboardCard_alarm_status,
  },
};
</script>
