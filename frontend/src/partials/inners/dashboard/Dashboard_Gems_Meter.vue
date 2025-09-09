<template>
  <div class="premium-dashboard-card">
    <!-- 헤더 섹션 -->
    <div class="card-header">
      <div class="header-content">
        <h2 class="card-title">{{t("dashboard.meter.title")}}</h2>
   
      </div>
    </div>

    <!-- 통합 요약 섹션 -->
    <div class="summary-section">
      <div class="summary-container">

        <!-- 전류 카드 -->
        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ data2.Itot || 0 }} <span class="summary-unit">A</span>
            </div>
            <div class="summary-label">
              {{ t("dashboard.meter.totcurrent") }}
            </div>
          </div>
        </div>

        <!-- 전압 카드 (다른 카드와 동일한 스타일) -->
        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ data2.Ig || 0 }} <span class="summary-unit">A</span>
            </div>
            <div class="summary-label">
              {{ t("dashboard.transDiag.Ig") }}
            </div>
          </div>
        </div>

        <!-- 주파수 카드 -->
        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ data2.Freq || 0 }} <span class="summary-unit">Hz</span>
            </div>
            <div class="summary-label">
              {{ t("dashboard.meter.frequency") }}
            </div>
          </div>
        </div>

        <!-- 역률 카드 -->
        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ data2.Temp || 0 }} <span class="summary-unit">℃</span>
            </div>
            <div class="summary-label">{{ t("dashboard.transDiag.Temperature") }}</div>
          </div>
        </div>

        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ data2.PF4 || 0 }} <span class="summary-unit">%</span>
            </div>
            <div class="summary-label">{{ t("dashboard.pq.powerfactor") }}</div>
          </div>
        </div>

        <!-- 유효전력 카드 -->
        <div class="summary-metric">
          <div class="summary-content">
            <div class="summary-value">
              {{ (data2.P4 / 1000).toFixed(2) || 0 }}
              <span class="summary-unit">kW</span>
            </div>
            <div class="summary-label">
              {{ t("dashboard.meter.activepower") }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 상세 정보 섹션 -->
    <div class="details-section">
      <div class="details-grid">
        <!-- 전압 상세 정보 -->
        <div class="detail-card voltage-detail">
          <div class="detail-header">
            <h3 class="detail-title">
              <span class="detail-icon">⚡</span>
              {{ t("dashboard.meter.voltage") }}
            </h3>
          </div>
          <div class="voltage-grid">
            <div
              class="voltage-item"
              v-for="(phase, index) in ['L1', 'L2', 'L3']"
              :key="phase"
            >
              <div class="phase-label-small" :class="`phase-${index + 1}`">
                {{ phase }}
              </div>
              <div class="phase-value">{{ data2[`U${index + 1}`] || 0 }}</div>
              <div class="phase-unit">V</div>
              <div class="phase-bar">
                <div
                  class="phase-fill"
                  :class="`phase-${index + 1}`"
                  :style="{
                    width:
                      Math.min(
                        ((data2[`U${index + 1}`] || 0) / 240) * 100,
                        100
                      ) + '%',
                  }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 불평형률 정보 -->
        <div class="detail-card unbalance-detail">
          <div class="detail-header">
            <h3 class="detail-title">
              <span class="detail-icon">⚖️</span>
              {{ t("dashboard.pq.unbalance") }}
            </h3>
          </div>
          <div class="unbalance-items">
            <div class="unbalance-item">
              <div class="unbalance-info">
                <span class="unbalance-type">{{
                  t("dashboard.pq.voltageunbalance")
                }}</span>
                <span
                  class="unbalance-value"
                  :class="getUnbalanceClass(data2.Ubal1)"
                >
                  {{ (data2.Ubal1 || 0).toFixed(1) }}%
                </span>
              </div>
              <div class="progress-container">
                <div class="progress-bar voltage">
                  <div
                    class="progress-fill"
                    :style="{ width: Math.min(data2.Ubal1 || 0, 100) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
            <div class="unbalance-item">
              <div class="unbalance-info">
                <span class="unbalance-type">{{
                  t("dashboard.pq.currentunbalance")
                }}</span>
                <span
                  class="unbalance-value"
                  :class="getUnbalanceClass(data2.Ibal1)"
                >
                  {{ (data2.Ibal1 || 0).toFixed(1) }}%
                </span>
              </div>
              <div class="progress-container">
                <div class="progress-bar current">
                  <div
                    class="progress-fill"
                    :style="{ width: Math.min(data2.Ibal1 || 0, 100) + '%' }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 고조파 차트 -->
        <div class="detail-card harmonics-detail">
          <div class="detail-header">
            <h3 class="detail-title">
              <span class="detail-icon">📊</span>
              {{ t("pq.tabs.harmonics") }}
            </h3>
          </div>
          <div class="harmonics-container">
            <DashboardCard_THD
              :data="data2"
              :height="120"
              @data-change="onDataChange"
            />
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import { watch, ref, computed, onMounted, onUnmounted, inject } from "vue";
import { useI18n } from "vue-i18n";
import DashboardCard_THD from "./DashboardCard_THD.vue";
import StatusItem from "./StatusItem_Trans.vue";
import { useSetupStore } from "@/store/setup";
import axios from "axios";

export default {
  name: "Dashboard_Gems_Meter",
  props: {
    channel: String,
    data: Object,
  },
  components: {
    DashboardCard_THD,
    StatusItem,
  },
  setup(props) {
    const { t } = useI18n();

    const channel = ref(props.channel);
    const data2 = ref({});

    // StatusItem 관련 추가 데이터
    const stData = ref({
      devName: "",
      devType: "",
      devStatus: -2,
      devNickname: "",
      Ig: 0,
    });

    const pqData = ref({
      devName: "",
      devStatus: -2,
    });
    const eventData = ref({});

    const transData = ref({});
    const rawdata = ref([]); // rawdata 추가
    const DiagEnable = ref(false);
    const setupStore = useSetupStore();
    const channelStatus = computed(() => setupStore.getChannelSetting);
    const asset = computed(() => setupStore.getAssetConfig);
    const assetTypes = ref("");
    const meterDictMain = inject("meterDictMain");
    const meterDictSub = inject("meterDictSub");
    let updateInterval = null;

    // 전체 시스템 상태 판정
    const getOverallStatus = () => {
      const voltage = data2.value.U4 || 0;
      const current = data2.value.Itot || 0;
      const freq = data2.value.Freq || 0;
      const pf = data2.value.PF4 || 0;

      // 임계값 체크
      const voltageOk = voltage >= 200 && voltage <= 240;
      const freqOk = freq >= 49 && freq <= 51;
      const pfOk = pf >= 80;

      if (!voltageOk || !freqOk || !pfOk) return "critical";
      if (pf < 90 || voltage < 210) return "warning";
      return "good";
    };

    const getOverallStatusText = () => {
      const status = getOverallStatus();
      switch (status) {
        case "critical":
          return "주의 필요";
        case "warning":
          return "경고";
        case "good":
          return "정상";
        default:
          return "알 수 없음";
      }
    };

    // 불평형률 상태
    const getUnbalanceClass = (value) => {
      if (!value) return "status-unknown";
      if (value >= 3) return "status-critical";
      if (value >= 2) return "status-warning";
      if (value >= 1) return "status-caution";
      return "status-good";
    };

    const getUnbalanceStatusClass = () => {
      const ubal = data2.value.Ubal1 || 0;
      const ibal = data2.value.Ibal1 || 0;
      const maxUnbalance = Math.max(ubal, ibal);

      if (maxUnbalance >= 3) return "critical";
      if (maxUnbalance >= 2) return "warning";
      if (maxUnbalance >= 1) return "caution";
      return "good";
    };

    const getUnbalanceStatusText = () => {
      const status = getUnbalanceStatusClass();
      switch (status) {
        case "critical":
          return "위험";
        case "warning":
          return "경고";
        case "caution":
          return "주의";
        case "good":
          return "양호";
        default:
          return "정상";
      }
    };

    const getStatusColor2 = (status) => {
        switch (status) {
          case 'OK': return 'bg-green-500/20 text-green-700 font-semibold dark:bg-green-600/40 dark:text-green-300';
          case 'Low': return 'bg-yellow-500/20 text-yellow-700 font-semibold dark:bg-yellow-600/40 dark:text-yellow-300';
          case 'Medium': return 'bg-orange-500/20 text-orange-700 font-semibold dark:bg-orange-600/40 dark:text-orange-300';
          case 'High': return 'bg-red-500/20 text-red-700 font-semibold dark:bg-red-600/40 dark:text-red-300';
          default: return 'bg-gray-500/20 text-gray-700 font-semibold dark:bg-gray-600/40 dark:text-gray-300';
        }
      }; 

      const getStatusCText = (status) => {
        switch (status) {
          case 'OK': return t('diagnosis.tabContext.pqfe1');
          case 'Low': return t('diagnosis.tabContext.pqfe2');
          case 'Medium': return t('diagnosis.tabContext.pqfe3');
          case 'High': return t('diagnosis.tabContext.pqfe4');
          default: return t('diagnosis.tabContext.pqfe0');
        }
      }; 

      const getStatusText = (status) => {
        switch (status) {
          case 1: return 'OK';
          case 2: return 'Low';
          case 3: return 'Medium';
          case 4: return 'High';
          default: return 'No Data';
        }
      }

    // StatusItem 관련 함수들
    const fetchData = async () => {
      if (
        !asset.value ||
        (!asset.value.assetName_main && !asset.value.assetName_sub)
      ) {
        console.log("⏳ asset 준비 안됨. fetchData 대기중");
        return;
      }
      const chName =
        channel.value == "main"
          ? asset.value.assetName_main
          : asset.value.assetName_sub;
      const chType =
        channel.value == "main"
          ? asset.value.assetType_main
          : asset.value.assetType_sub;
      const chNick =
        channel.value == "main"
          ? asset.value.assetNickname_main
          : asset.value.assetNickname_sub;

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getStatus/${chName}`);
          if (response.data.status >= 0) {
            stData.value.devName = chName;
            stData.value.devType = chType;
            stData.value.devStatus = response.data.status;
            stData.value.devNickname = chNick;
            if (assetTypes.value == "Transformer") {
              if (channel.value === "main") {
                transData.value = {
                  Temp: meterDictMain.value.Temp,
                  Ig: meterDictMain.value.Ig,
                  Stotal: meterDictMain.value.S4,
                };
              } else {
                transData.value = {
                  Temp: meterDictSub.value.Temp,
                  Ig: meterDictSub.value.Ig,
                  Stotal: meterDictMain.value.S4,
                };
              }
            } else {
              stData.value.Ig =
                channel.value === "main"
                  ? meterDictMain.value.Ig
                  : meterDictSub.value.Ig;
            }
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    const fetchRealData = async () => {
      if (
        !asset.value ||
        (!asset.value.assetName_main && !asset.value.assetName_sub)
      ) {
        return;
      }
      const chName =
        channel.value == "main"
          ? asset.value.assetName_main
          : asset.value.assetName_sub;
      const chType =
        channel.value == "main"
          ? asset.value.assetType_main
          : asset.value.assetType_sub;

      if (chName != "") {
        try {
          const response = await axios.get(
            `/api/getRealTime/${chType}/${chName}`
          );
          if (response.data.success) {
            transData.value = response.data.data;
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    const fetchPQData = async () => {
      if (
        !asset.value ||
        (!asset.value.assetName_main && !asset.value.assetName_sub)
      ) {
        return;
      }
      const chName =
        channel.value == "main"
          ? asset.value.assetName_main
          : asset.value.assetName_sub;

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getPQStatus/${chName}`);
          if (response.data.status >= 0) {
            pqData.value.devName = response.data.item;
            pqData.value.devStatus = response.data.status;
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    const fetchEventData = async () => {
      if (
        !asset.value ||
        (!asset.value.assetName_main && !asset.value.assetName_sub)
      ) {
        return;
      }
      const chName =
        channel.value == "main"
          ? asset.value.assetName_main
          : asset.value.assetName_sub;

      if (chName != "") {
        try {
          const response = await axios.get(`/api/getEvents/${chName}`);
          if (response.data.success) {
            eventData.value = response.data.data_status;
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
      }
    };

    // 이벤트 핸들러
    const onChartReady = (chartInstance) => {
      console.log("고조파 차트 준비 완료:", chartInstance);
    };

    const onDataChange = (chartInfo) => {
      //console.log("차트 데이터 변경:", chartInfo);
    };

    // 데이터 감시
    watch(
      () => props.data,
      (newData) => {
        if (newData && Object.keys(newData).length > 0) {
          data2.value = {
            U4: newData.U4 || 0,
            U1: newData.U1 || 0,
            U2: newData.U2 || 0,
            U3: newData.U3 || 0,
            Itot: newData.Itot || 0,
            Freq: newData.Freq || 0,
            PF4: newData.PF4 || 0,
            P4: newData.P4 || 0,
            Ubal1: newData.Ubal1 || 0,
            Ibal1: newData.Ibal1 || 0,
            ...newData,
          };
        }
      },
      { immediate: true }
    );

    // asset 감시
    // watch(
    //   asset,
    //   (newVal) => {
    //     if (newVal) {
    //       if (channel.value == "main") assetTypes.value = newVal.assetType_main;
    //       else assetTypes.value = newVal.assetType_sub;

    //       fetchData();
    //       if (assetTypes.value != "Transformer") {
    //         fetchRealData();
    //       }
    //       fetchPQData();
    //       fetchEventData();

    //       if (updateInterval) {
    //         clearInterval(updateInterval);
    //       }
    //       if (!updateInterval) {
    //         updateInterval = setInterval(async () => {
    //           await fetchData();
    //           if (assetTypes.value != "Transformer") {
    //             await fetchRealData();
    //           }
    //           await fetchPQData();
    //           await fetchEventData();
    //         }, 300000); // 5분
    //       }
    //     }
    //   },
    //   { immediate: true }
    // );

    // DiagEnable 감시
    watch(
      () => channelStatus.value,
      () => {
        if (channel.value == "main")
          DiagEnable.value = channelStatus.value.MainDiagnosis;
        else DiagEnable.value = channelStatus.value.SubDiagnosis;
      },
      { immediate: true }
    );

    onMounted(async () => {
      await setupStore.checkSetting();
    });

    onUnmounted(() => {
      if (updateInterval) {
        clearInterval(updateInterval);
      }
    });

    return {
      channel,
      data2,
      stData,
      pqData,
      transData,
      rawdata,
      DiagEnable,
      t,
      getOverallStatus,
      getOverallStatusText,
      getUnbalanceClass,
      getUnbalanceStatusClass,
      getUnbalanceStatusText,
      onChartReady,
      onDataChange,
      channelStatus,
      asset,
      assetTypes,
      eventData,
      getStatusText,
      getStatusColor2,
      getStatusCText,
    };
  },
};
</script>

<style scoped>
.premium-dashboard-card {
  @apply col-span-full sm:col-span-7 xl:col-span-7;
  @apply bg-gradient-to-br from-white to-gray-50 dark:from-gray-800 dark:to-gray-900;
  @apply shadow-lg rounded-xl border border-gray-200/50 dark:border-gray-700/50;
  @apply backdrop-blur-sm;
  @apply transition-all duration-300 hover:shadow-xl;
}

/* 헤더 섹션 */
.card-header {
  @apply px-5 py-4 bg-gradient-to-r from-blue-50/50 to-purple-50/50 dark:from-blue-900/20 dark:to-purple-900/20;
  @apply rounded-t-xl;
}

.header-content {
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-lg font-bold text-gray-900 dark:text-gray-100;
  @apply bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent;
}

.channel-info {
  @apply flex items-center;
}

.channel-text {
  @apply text-xs font-semibold text-gray-400 dark:text-gray-300 uppercase;
}

/* StatusItem 섹션 */
.status-item-section {
  @apply px-5 py-3;
  @apply bg-gradient-to-r from-gray-50/30 to-gray-100/30 dark:from-gray-800/30 dark:to-gray-900/30;
}

/* 통합 요약 섹션 */
.summary-section {
  @apply px-5 py-3 mt-2;
}

.summary-container {
  @apply grid grid-cols-6 gap-4;
}

.summary-metric {
  @apply p-4 rounded-lg bg-white dark:bg-gray-800;
  @apply border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-all duration-200;
  @apply flex flex-col items-center text-center;
  @apply min-h-[80px];
}

.summary-content {
  @apply w-full;
}

.summary-value {
  @apply text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2;
}

.summary-unit {
  @apply text-base font-medium text-gray-500 dark:text-gray-400;
}

.summary-label {
  @apply text-sm text-gray-500 dark:text-gray-400 leading-tight;
}

/* 상세 정보 섹션 */
.details-section {
  @apply px-5 py-2;
}

.details-grid {
  @apply grid grid-cols-1 lg:grid-cols-3 gap-4;
}

.detail-card {
  @apply bg-white dark:bg-gray-800 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-all duration-200;
  @apply overflow-hidden;
  @apply min-h-[160px];
}

.detail-header {
  @apply p-3 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-600;
  @apply flex justify-between items-center;
}

.detail-title {
  @apply text-sm font-semibold text-gray-900 dark:text-gray-100;
  @apply flex items-center gap-2;
}

.detail-icon {
  @apply text-base;
}

.detail-badge {
  @apply px-2 py-1 text-xs font-medium rounded-full;
  @apply bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400;
}

.detail-badge.good {
  @apply bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400;
}

.detail-badge.warning {
  @apply bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400;
}

.detail-badge.critical {
  @apply bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400;
}

.detail-badge.chart-badge {
  @apply bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400;
}

/* 전압 상세 - 패딩과 간격 증가 */
.voltage-grid {
  @apply p-4 space-y-3;
}

.voltage-item {
  @apply flex items-center gap-4;
}

.phase-label-small {
  @apply w-7 h-7 rounded-full flex items-center justify-center;
  @apply text-sm font-bold text-white;
}

.phase-label-small.phase-1 {
  @apply bg-red-500;
}

.phase-label-small.phase-2 {
  @apply bg-yellow-500;
}

.phase-label-small.phase-3 {
  @apply bg-blue-500;
}

.phase-value {
  @apply text-xl font-bold text-gray-900 dark:text-gray-100 min-w-[70px];
}

.phase-unit {
  @apply text-base text-gray-500 dark:text-gray-400 min-w-[25px];
}

.phase-bar {
  @apply flex-1 h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden;
}

.phase-fill {
  @apply h-full rounded-full transition-all duration-1000;
}

.phase-fill.phase-1 {
  @apply bg-gradient-to-r from-red-400 to-red-600;
}

.phase-fill.phase-2 {
  @apply bg-gradient-to-r from-yellow-400 to-yellow-600;
}

.phase-fill.phase-3 {
  @apply bg-gradient-to-r from-blue-400 to-blue-600;
}

/* 불평형률 - 패딩과 간격 증가 */
.unbalance-items {
  @apply p-4 space-y-4;
}

.unbalance-item {
  @apply space-y-3;
}

.unbalance-info {
  @apply flex justify-between items-center;
}

.unbalance-type {
  @apply text-sm text-gray-600 dark:text-gray-400;
}

.unbalance-value {
  @apply text-sm font-bold;
}

.unbalance-value.status-good {
  @apply text-green-600 dark:text-green-400;
}

.unbalance-value.status-warning {
  @apply text-yellow-600 dark:text-yellow-400;
}

.unbalance-value.status-critical {
  @apply text-red-600 dark:text-red-400;
}

.progress-container {
  @apply w-full;
}

.progress-bar {
  @apply w-full h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden;
}

.progress-bar.voltage .progress-fill {
  @apply bg-gradient-to-r from-violet-400 to-violet-600;
}

.progress-bar.current .progress-fill {
  @apply bg-gradient-to-r from-sky-400 to-sky-600;
}

.progress-fill {
  @apply h-full rounded-full transition-all duration-1000;
}

/* 고조파 차트 - 패딩 조정 */
.harmonics-container {
  @apply p-1;
}

/* 데이터 테이블 */
.data-table {
  @apply bg-white dark:bg-gray-800 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-all duration-200;
  /* @apply h-full flex flex-col;   */
}

.table-header {
  @apply p-3 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-600;
  @apply flex justify-between items-center;
}

.table-title {
  @apply text-sm font-semibold text-gray-900 dark:text-gray-100;
  @apply flex items-center gap-2;
}

.table-content {
  @apply px-3 pt-1;
  /* @apply flex-1 overflow-y-auto;  */
}

.data-row {
  @apply flex justify-between items-center;
  @apply border-b border-gray-200/50 dark:border-gray-600/50 last:border-b-0;
  @apply transition-colors duration-200 hover:bg-gray-100/50 dark:hover:bg-gray-600/30;
  @apply rounded px-3;
  @apply min-h-[2.8rem];
}

.data-label {
  @apply text-sm font-medium text-gray-600 dark:text-white;
  @apply flex-1;
}

.data-value {
  @apply text-sm font-bold text-gray-800 dark:text-gray-100;
  @apply flex items-baseline gap-1;
}

.data-unit {
  @apply text-xs font-semibold text-gray-500 dark:text-white;
}

/* 반응형 */
@media (max-width: 1024px) {
  .summary-container {
    @apply grid-cols-3 gap-3;
  }

  .details-grid {
    @apply grid-cols-1 gap-3;
  }

  .summary-metric {
    @apply min-h-[70px] p-3;
  }

  .summary-value {
    @apply text-xl;
  }
}

@media (max-width: 640px) {
  .premium-dashboard-card {
    @apply col-span-full;
  }

  .card-header {
    @apply px-3 py-2;
  }

  .summary-section {
    @apply px-3 py-2;
  }

  .details-section {
    @apply px-3 py-2;
  }

  .status-item-section {
    @apply px-3 py-2;
  }

  .summary-container {
    @apply grid-cols-2 gap-2;
  }

  .summary-value {
    @apply text-lg;
  }

  .summary-metric {
    @apply min-h-[60px] p-2;
  }

  .detail-card {
    @apply min-h-[160px];
  }
}
.detail-card {
  @apply bg-white dark:bg-gray-800 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-all duration-200;
  @apply overflow-hidden;
  @apply mb-2;
}

.card-header {
  @apply p-3 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-600;
  @apply flex justify-between items-center;
}

.card-title {
  @apply text-sm font-semibold text-gray-900 dark:text-gray-100;
  @apply flex items-center gap-2;
}

/* StatusItem의 status-section 스타일 적용 */
.status-section {
  @apply px-2 py-4;
}

/* Event Status용 data-row 스타일 */
.data-row {
  @apply flex items-center justify-between;
  @apply px-3 py-2;
  @apply rounded-lg transition-all duration-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-700/30;
}

/* 마지막 행 마진 조정 */
.data-row:not(.last-row) {
  @apply mb-2;
}

/* data-label - StatusItem의 status-text와 유사한 스타일 */
.data-label {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
  @apply transition-all duration-200;
}

/* status-badge - 통일된 뱃지 스타일 */
.status-badge {
  @apply text-sm rounded-full px-3 py-1;
  @apply min-w-[100px] text-center;
  @apply font-medium;
  @apply transition-all duration-200;
  @apply shadow-sm;
  @apply inline-block;
}

/* 반응형 디자인 - StatusItem과 통일 */
@media (max-width: 640px) {
  .detail-card {
    @apply col-span-full;
  }
  
  .card-header {
    @apply p-2;
  }
  
  .status-section {
    @apply p-3;
  }
  
  .data-row {
    @apply px-2 py-2;
  }
  
  .data-label {
    @apply text-xs;
  }
  
  .status-badge {
    @apply text-xs px-2 py-0.5 min-w-[80px];
  }
}

/* 다크모드 대응 */
@media (prefers-color-scheme: dark) {
  .detail-card {
    @apply bg-gray-800;
  }
}
/* StatusItem 섹션 - 높이 조정 */
.status-item-section {
  @apply px-5 py-3;
  @apply bg-gradient-to-r from-gray-50/30 to-gray-100/30 dark:from-gray-800/30 dark:to-gray-900/30;
}

/* 이벤트 카드 높이 조정 */
.status-item-section .detail-card {
  @apply bg-white dark:bg-gray-800 rounded-lg;
  @apply border border-gray-200 dark:border-gray-700;
  @apply shadow-sm hover:shadow-md transition-all duration-200;
  @apply overflow-hidden;
  @apply mb-0; /* 기존 mb-2 제거 */
  height: 100%; /* 부모 높이에 맞춤 */
}

/* 이벤트 상태 섹션 패딩 조정 */
.status-section {
  @apply px-4 py-4;
  display: flex;
  flex-direction: column;
}

/* 데이터 행 간격 조정 */
.data-row {
  @apply flex items-center justify-between;
  @apply px-3 py-2.5; /* 패딩 약간 증가 */
  @apply rounded-lg transition-all duration-200;
  @apply hover:bg-gray-50 dark:hover:bg-gray-700/30;
}

/* 이벤트 항목 간격 */
.status-item-section .space-y-2 > div:not(:last-child) {
  @apply mb-2;
}

/* 반응형 - 모바일에서는 세로 정렬 유지 */
@media (max-width: 1024px) {
  .status-item-section .h-full {
    height: auto;
  }
  
  .status-item-section .detail-card {
    height: auto;
    min-height: 200px;
  }
}
</style>