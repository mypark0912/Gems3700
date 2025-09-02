<template>
  <div class="grid grid-cols-12 gap-6">
    <!-- 왼쪽 박스 (Diagnosis_Info + Diagnosis_Barchart) -->
    <div class="col-span-4 flex flex-col gap-6 h-auto">
      <!-- ✅ flex-col로 수직 정렬 강제 -->

 <div v-show="true" role="alert">
  <div
    class="flex flex-col w-full px-4 py-2 rounded-lg text-sm bg-white dark:bg-gray-800 shadow-sm border border-gray-200 dark:border-gray-700/60 text-gray-600 dark:text-gray-400">
    <div class="flex w-full justify-between items-start">
      <div class="flex">
        <svg class="shrink-0 fill-current text-violet-500 mt-[3px] mr-3" width="16" height="16"
          viewBox="0 0 16 16">
          <path
            d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm1 12H7V7h2v5zM8 6c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1z" />
        </svg>
        <div>
          <slot />{{ t("trend.TrendTab.slot") }}.
        </div>
      </div>
    </div>
    <div class="text-right mt-1">
      <a v-if="tap == `Meter`"
        class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400" href="#0"
        @click.prevent="drawMeterChart">{{ t("trend.TrendTab.Plot") }}</a>
      <a v-else-if="tap == `Energy`"
        class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400" href="#0"
        @click.prevent="drawEnergyChart">{{ t("trend.TrendTab.Plot") }}</a>
      <a v-else class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400" href="#0"
        @click.prevent="drawDiagnosisChart">{{ t("trend.TrendTab.Plot") }}</a>
    </div>
  </div>
</div>

      <Trend_treetable v-if="items.length > 0" :channel="channel" :data="items" :expanded="false"
        :checked-ids="checkedIds" :checked-names="checkedNames" @check-change="onCheckChange" />
    </div>

    <!-- 오른쪽 박스 (Diagnosis_TreeTable) - 높이 제한 없음 -->
    <div class="col-span-8">
      <LineChart v-if="tap == `Diagnosis`" :chart-data="option.lineData" :chart-labels="option.lineLabels" />
      <LineChart v-if="tap == `Meter`" :chart-data="meterOption.lineData" :chart-labels="meterOption.lineLabels" />
      <LineChart v-if="tap == `PowerQuality`" :chart-data="option.lineData" :chart-labels="option.lineLabels" />
      <LineChart v-if="tap == `Parameters`" :chart-data="option.lineData" :chart-labels="option.lineLabels" />
      <LineChart v-if="tap == `Energy`" :chart-data="energyOption.lineData" :chart-labels="energyOption.lineLabels" />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from "vue";
import axios from "axios";
import Trend_treetable from "./Trend_treetable.vue";
import LineChart from "../../../charts/connect/LineChart01_Echart.vue";
import Notification_origin from "./Notification_origin.vue";
import { useI18n } from "vue-i18n";
export default {
  name: "TrendTab",
  props: {
    startdate: {
      type: Object,
      default: null,
    },
    channel: {
      type: String,
      default: "",
    },
    tap: {
      type: String,
      default: "",
    },
    enddate: {
      type: Object,
      default: null,
    },
    asset: {
      type: String,
      default: "",
    },
  },
  components: {
    LineChart,
    Trend_treetable,
    Notification_origin,
  },
  setup(props) {
    const { t } = useI18n();
    const channel = ref(props.channel);
    const tap = ref(props.tap);
    const items = ref([]);
    const option = ref({
      lineLabels: [],
      lineData: [],
    });
    const meterOption = ref({ lineLabels: [], lineData: [] });
    const energyOption = ref({ lineLabels: [], lineData: [] });
    const pqOption = ref({ lineLabels: [], lineData: [] });
    const checkedIds = ref([]);
    const checkedNames = ref([]);

    // Meter 탭용 트리 데이터 (Energy 제외)
    const trendTreeData = [
      {
        ID: 1,
        Name: "Temp",
        Title: "Temperature",
        isParent: true,
      },
      { ID: 2, Name: "Freq", Title: "Frequency", isParent: true },
      {
        ID: 3, Name: "PF", Title: "Power Factor", isParent: true,
        children: [
          { ID: 4, Name: "PF1", Title: "Power Factor L1" },
          { ID: 5, Name: "PF2", Title: "Power Factor L2" },
          { ID: 6, Name: "PF3", Title: "Power Factor L3" },
        ],
      },
      {
        ID: 7,
        Name: "linevoltage",
        Title: "Line Voltage",
        isParent: true,
        children: [
          { ID: 8, Name: "Upp1", Title: "Line Voltage L12" },
          { ID: 9, Name: "Upp2", Title: "Line Voltage L23" },
          { ID: 10, Name: "Upp3", Title: "Line Voltage L31" },
          { ID: 11, Name: "Upp4", Title: "Line Voltage Average" },
        ],
      },
      {
        ID: 12,
        Name: "voltage",
        Title: "Phase Voltage",
        isParent: true,
        children: [
          { ID: 13, Name: "U1", Title: "Voltage L1" },
          { ID: 14, Name: "U2", Title: "Voltage L2" },
          { ID: 15, Name: "U3", Title: "Voltage L3" },
        ],
      },
      {
        ID: 16,
        Name: "current",
        Title: "Current",
        isParent: true,
        children: [
          { ID: 17, Name: "I1", Title: "Current L1" },
          { ID: 18, Name: "I2", Title: "Current L2" },
          { ID: 19, Name: "I3", Title: "Current L3" },
        ],
      },
      {
        ID: 20,
        Name: "unbalance",
        Title: "Unbalance",
        isParent: true,
        children: [
          { ID: 21, Name: "Ibal1", Title: "Current Unbalance" },
          { ID: 22, Name: "Ubal1", Title: "Voltage Unbalance" },
        ],
      },
      {
        ID: 23,
        Name: "Power",
        Title: "Power",
        isParent: true,
        children: [
          { ID: 24, Name: "P4", Title: "Active Power" },
          { ID: 25, Name: "Q4", Title: "Reactive Power" },
          { ID: 26, Name: "S4", Title: "Apparent Power" },
        ],
      },
      {
        ID: 31,
        Name: "THD",
        Title: "THD",
        isParent: true,
        children: [
          { ID: 32, Name: "THD_U1", Title: "THD Voltage L1" },
          { ID: 33, Name: "THD_U2", Title: "THD Voltage L2" },
          { ID: 34, Name: "THD_U3", Title: "THD Voltage L3" },
          { ID: 35, Name: "THD_Upp1", Title: "THD Voltage L12" },
          { ID: 36, Name: "THD_Upp2", Title: "THD Voltage L23" },
          { ID: 37, Name: "THD_Upp3", Title: "THD Voltage L31" },
          { ID: 38, Name: "THD_I1", Title: "THD Current L1" },
          { ID: 39, Name: "THD_I2", Title: "THD Current L2" },
          { ID: 40, Name: "THD_I3", Title: "THD Current L3" },
        ],
      },
      {
        ID: 41,
        Name: "TDD",
        Title: "TDD",
        isParent: true,
        children: [
          { ID: 42, Name: "TDD_I1", Title: "TDD Current L1" },
          { ID: 43, Name: "TDD_I2", Title: "TDD Current L2" },
          { ID: 44, Name: "TDD_I3", Title: "TDD Current L3" },
          //{ ID: 45, Name: "tdd_i_avg", Title: "TDD Current Average" },
        ],
      },
    ];

    // Energy 탭 전용 트리 데이터
    const energyTreeData = [
      {
        ID: 27,
        Name: "Energy",
        Title: "Energy",
        isParent: true,
        children: [
          { ID: 28, Name: "kwh_import_consumption", Title: "Active Energy" },
          { ID: 29, Name: "kvarh_import_consumption", Title: "Reactive Energy" },
          { ID: 30, Name: "kvah_import_consumption", Title: "Apparent Energy" },
        ],
      },
    ];

    const meterParamMap = {
      Temperature: ["Temp"],
      Frequency: ["Freq"],
      "Phase Voltage": ["voltage"],
      "Line Voltage": ["linevoltage"],
      "Line Voltage L12": ["Upp1"],
      "Line Voltage L23": ["Upp2"],
      "Line Voltage L31": ["Upp3"],
      "Line Voltage Average": ["Upp4"],
      Current: ["current"],
      "Current L1": ["I1"],
      "Current L2": ["I2"],
      "Current L3": ["I3"], 
      Power: ["Power"],
      "Active Power": ["P4"],
      "Reactive Power": ["Q4"],
      "Apparent Power": ["S4"],
      PF: ["PF"],
      Unbalance: ["unbalance"],
      "Current Unbalance": ["Ibal1"],
      "Voltage Unbalance": ["Ubal1"],
      THD: ["THD"],
      "THD Voltage L1": ["THD_U1"],
      "THD Voltage L2": ["THD_U2"],
      "THD Voltage L3": ["THD_U3"],
      "THD Voltage L12": ["THD_Upp1"],
      "THD Voltage L23": ["THD_Upp2"],
      "THD Voltage L31": ["THD_Upp3"],
      "THD Current L1": ["THD_I1"],
      "THD Current L2": ["THD_I2"],
      "THD Current L3": ["THD_I3"],
      TDD: ["TDD"],
      "TDD I1": ["TDD_I1"],
      "TDD I2": ["TDD_I2"],
      "TDD I3": ["TDD_I3"],
      "TDD Iavg": ["tdd_i_avg"],
    };

 

    const expandEnabledNames = (enabledParams, paramMap) => {
      const paramNames = enabledParams.flatMap(
        (param) => paramMap[param] || []
      );

      const expanded = new Set(paramNames);

      // 사용할 트리 데이터 결정
      const treeData = tap.value === "Energy" ? energyTreeData : trendTreeData;

      treeData.forEach((node) => {
        if (node.children && paramNames.includes(node.Name)) {
          node.children.forEach((child) => {
            expanded.add(child.Name);
          });
        }
      });

      return Array.from(expanded);
    };

    const fetchData = async () => {
      try {
        if (tap.value == "Diagnosis") {
          if(props.asset != ""){
            const response = await axios.get(
            `/api/getTrendParameters/${props.asset}/Diagnostic`
            );
            if (response.data.success) {
              items.value = response.data.superlist;
            } else {
              console.log("No Data");
            }
          }

        } else if (tap.value == "PowerQuality") {
          if(props.asset != ""){
            const response = await axios.get(
            `/api/getTrendParameters/${props.asset}/PowerQuality`
            );
            if (response.data.success) {
              items.value = response.data.superlist;
            } else {
              console.log("No Data");
            }
          }
        }
        else if (tap.value == "Parameters") {
          if(props.asset != ""){
            const response = await axios.get(
              `/api/getTrendParameters/${props.asset}/parameter`
            );
            if (response.data.success) {
              items.value = response.data.superlist;
            } else {
              console.log("No Data");
            }
          }
        } else if (tap.value === "Meter") {
          const response = await axios.get(
            `/api/getTrendParameters/${channel.value}`
          );
          if (response.data.success) {
            const enabledParams = response.data.data.params;
            console.log("enabledParams",enabledParams)
            // Energy 관련 파라미터 제외
            const filteredParams = enabledParams.filter(param =>
              !["Energy", "Active Energy", "Reactive Energy", "Apparent Energy"].includes(param)
            );

            const enabledNames = expandEnabledNames(filteredParams, meterParamMap);

            const deepFilterTree = (nodes) => {
              return nodes
                .map((node) => {
                  const isNodeEnabled = enabledNames.includes(node.Name);

                  if (node.children) {
                    const filteredChildren = node.children.filter((child) =>
                      enabledNames.includes(child.Name)
                    );
                    if (filteredChildren.length > 0 || isNodeEnabled) {
                      return {
                        ...node,
                        children: filteredChildren,
                      };
                    } else {
                      return null;
                    }
                  } else {
                    return isNodeEnabled ? node : null;
                  }
                })
                .filter(Boolean);
            };

            items.value = deepFilterTree(trendTreeData);
          }
        } else if (tap.value === "Energy") {
          // Energy 탭은 항상 모든 Energy 관련 항목을 표시
          items.value = energyTreeData;
        }
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    function onCheckChange({
      id,
      checked,
      name,
      children = [],
      childrenNames = [],
    }) {
      if (checked) {
        if (!checkedIds.value.includes(id)) checkedIds.value.push(id);
        children.forEach((cid) => {
          if (!checkedIds.value.includes(cid)) checkedIds.value.push(cid);
        });
        if (!checkedNames.value.includes(name)) checkedNames.value.push(name);
        childrenNames.forEach((cname) => {
          if (!checkedNames.value.includes(cname))
            checkedNames.value.push(cname);
        });
      } else {
        checkedIds.value = checkedIds.value.filter(
          (item) => item !== id && !children.includes(item)
        );
        checkedNames.value = checkedNames.value.filter(
          (item) => item !== name && !childrenNames.includes(item)
        );
      }
      syncParentCheck();
      //console.log("checkedIds", checkedIds);
      //console.log("checkedNames", checkedNames);
    }

    function syncParentCheck() {
      const traverse = (node) => {
        if (node.children) {
          node.children.forEach(traverse);
          const allChildrenChecked = node.children.every((child) =>
            checkedIds.value.includes(child.ID)
          );
          if (allChildrenChecked) {
            if (!checkedIds.value.includes(node.ID)) {
              checkedIds.value.push(node.ID);
            }
          } else {
            checkedIds.value = checkedIds.value.filter((id) => id !== node.ID);
          }
        }
      };
      items.value.forEach(traverse);
    }

    const formatToISOString = (date, soe) => {
      if (typeof date === "string") {
        date = new Date(date);
      }
      if (!(date instanceof Date) || isNaN(date)) {
        throw new Error("Invalid date");
      }

      const pad = (num, size = 2) => String(num).padStart(size, "0");

      const year = date.getFullYear();
      const month = pad(date.getMonth() + 1); // 월은 0부터 시작
      const day = pad(date.getDate());
      let hours, minutes, seconds, milliseconds;

      if (soe === 0) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
        milliseconds = pad(1, 7);
      } else if (soe === 1) {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
        milliseconds = pad(999, 7);
      } else if (soe === 2) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
        milliseconds = pad(1, 2);
      } else {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
        milliseconds = pad(99, 2);
      }
      // 타임존 오프셋 계산
      const timezoneOffset = -date.getTimezoneOffset();
      const offsetSign = timezoneOffset >= 0 ? "+" : "-";
      const offsetHours = pad(Math.abs(Math.floor(timezoneOffset / 60)));
      const offsetMinutes = pad(Math.abs(timezoneOffset % 60));
      if (soe === 2 || soe === 3) {
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}Z`;
      } else {
        return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}${offsetSign}${offsetHours}:${offsetMinutes}`;
      }
    };

    const drawMeterChart = async () => {
      // checkedNames가 빈 배열이면 API 요청하지 않음
      if (!checkedNames.value || checkedNames.value.length === 0) {       
        meterOption.value = { lineLabels: [], lineData: [] };
        return;
      }
      try {
        //console.log("Selected Parameters:", checkedNames.value);

        let url = `/api/getMeterTrend/${channel.value}`;
        const StartDate = formatToISOString(props.startdate, 2);
        const EndDate = formatToISOString(props.enddate, 3);
        url += `?startDate=${StartDate}&endDate=${EndDate}`;

        const response = await axios.get(url);

        //console.log("dblenght:", response.data.data.length);
        //console.log(response.data);

        const responseData = response.data.data;
        const datasets = [];
        const labels = [];
        const selectedParams = checkedNames.value;

        // 1. 시간 라벨 (_time)
        labels.push(...responseData.map((row) => row._time));

        // 2. 파라미터 그룹 정의 (Energy 제외)
        const paramMap = {
          Temp: ["Temp"],
          U1: ["U1"],
          U2: ["U2"],
          U3: ["U3"],
          U4: ["U4"],
          Upp1: ["Upp1"],
          Upp2: ["Upp2"],
          Upp3: ["Upp3"],
          Upp4: ["Upp4"],          
          I1: ["I1"],
          I2: ["I2"],
          I3: ["I3"],
          P4: ["P4"],
          Q4: ["Q4"],
          S4: ["S4"],
          Freq: ["Freq"],      
          PF1: ["PF1"],
          PF2: ["PF2"],
          PF3: ["PF3"],
          vunbal: ["Ubal1"],
          curunbal: ["Ibal1"],
          THD_U1: ["THD_U1"],
          THD_U2: ["THD_U2"],
          THD_U3: ["THD_U3"],
          THD_Upp1: ["THD_Upp1"],
          THD_Upp2: ["THD_Upp2"],
          THD_Upp3: ["THD_Upp3"],
          THD_I1: ["THD_I1"],
          THD_I2: ["THD_I2"],
          THD_I3: ["THD_I3"],
          TDD_I1: ["TDD_I1"],
          TDD_I2: ["TDD_I2"],
          TDD_I3: ["TDD_I3"],
          Ubal1: ["Ubal1"],
          Ibal1: ["Ibal1"],
        };
        console.log("selectedParams",selectedParams)
        // 3. 선택된 항목만 datasets 생성
        selectedParams.forEach((param) => {
          const keys = paramMap[param];
          if (!keys) return;

          keys.forEach((key) => {
            const data = responseData.map((row) => row[key]);
            datasets.push({
              name: nameToTitleMap[key] || key,
              data: data,
              isThreshold: false,
            });
          });
        });

        // 4. 차트 옵션에 데이터 적용
        meterOption.value = {
          lineLabels: labels,
          lineData: datasets,
        };
      } catch (error) {
        console.log("데이터 가져오기 실패:", error);
      }
    };

    const drawEnergyChart = async () => {
      try {
        // checkedNames가 빈 배열이면 API 요청하지 않음
        if (!checkedNames.value || checkedNames.value.length === 0) {       
          energyOption.value = { lineLabels: [], lineData: [] };
          return;
        }

        let url = `/api/getEnergyTrend/${channel.value}`;
        const StartDate = formatToISOString(props.startdate, 2);
        const EndDate = formatToISOString(props.enddate, 3);
        url += `?startDate=${StartDate}&endDate=${EndDate}`;

        // console.log("Energy API 호출:", url);
        const response = await axios.get(url);

        if (response.data.result){
          const responseData = response.data.data;
     
        if (!Array.isArray(responseData)) {
          console.error("Energy 응답 데이터가 배열이 아닙니다:", responseData);
          energyOption.value = { lineLabels: [], lineData: [] };
          return;
        }

        // 빈 배열 체크
        if (responseData.length === 0) {
          console.warn("Energy 데이터가 비어있습니다");
          energyOption.value = { lineLabels: [], lineData: [] };
          return;
        }

        const datasets = [];
        const labels = [];
        const selectedParams = checkedNames.value;

        // 1. 시간 라벨 (_time) - 안전하게 처리
        try {
          labels.push(...responseData.map((row) => row._time));
          //console.log("Energy 시간 라벨 생성됨:", labels.length, "개");
        } catch (timeError) {
          console.error("시간 라벨 생성 실패:", timeError);
          energyOption.value = { lineLabels: [], lineData: [] };
          return;
        }

        // 2. Energy 파라미터 그룹 정의
        const energyParamMap = {
          kwh_import_consumption: ["kwh_import_consumption"],
          kvarh_import_consumption: ["kvarh_import_consumption"],
          kvah_import_consumption: ["kvah_import_consumption"],
        };

        // 3. 선택된 Energy 항목만 datasets 생성
        selectedParams.forEach((param) => {
          const keys = energyParamMap[param];
          if (!keys) {
            console.warn("Energy 파라미터 매핑을 찾을 수 없음:", param);
            return;
          }

          keys.forEach((key) => {
            try {
              const data = responseData.map((row) => row[key]);
              datasets.push({
                name: nameToTitleMap[key] || key,
                data: data,
                isThreshold: false,
              });
              //console.log(`Energy 데이터셋 생성됨: ${key}, 데이터 개수: ${data.length}`);
            } catch (dataError) {
              console.error(`Energy 데이터 생성 실패 (${key}):`, dataError);
            }
          });
        });

        // 4. 차트 옵션에 데이터 적용
        energyOption.value = {
          lineLabels: labels,
          lineData: datasets,
        };

        //console.log("Energy 차트 옵션 업데이트 완료:", energyOption.value);
        }else{
          console.error("Energy API 응답 실패:", response.data);
          energyOption.value = { lineLabels: [], lineData: [] };
        }
      } catch (error) {
        console.error("Energy 데이터 가져오기 실패:", error);
        energyOption.value = { lineLabels: [], lineData: [] };
      }
    };

    function isParentId(id) {
      const findNode = (nodes) => {
        for (const node of nodes) {
          if (node.ID === id) {
            return node.isParent === true; // 부모라면 true 반환
          }
          if (node.children) {
            const found = findNode(node.children);
            if (found !== undefined) return found;
          }
        }
        return undefined;
      };
      return findNode(items.value) || false;
    }

    function getEffectiveParameterIds() {
      const result = [];

      items.value.forEach((node) => {
        if (node.children && node.children.length > 0) {
          // 자식 있는 부모
          const childIds = node.children.map((c) => c.ID);
          const allChildrenChecked = childIds.every((cid) =>
            checkedIds.value.includes(cid)
          );
          if (allChildrenChecked) {
            result.push(...childIds); // 자식만 보냄
          } else {
            node.children.forEach((c) => {
              if (checkedIds.value.includes(c.ID)) {
                result.push(c.ID);
              }
            });
          }
        } else {
          // 리프 노드 (부모든 아니든)
          if (checkedIds.value.includes(node.ID)) {
            result.push(node.ID);
          }
        }
      });

      return result;
    }

    const drawDiagnosisChart = async () => {
      const effectiveIds = getEffectiveParameterIds();
      // effectiveIds가 빈 배열이면 API 요청하지 않음
      if (!effectiveIds || effectiveIds.length === 0) {

        option.value = {
          lineLabels: [],
          lineData: [],
        };
        return;
      }
      const trendDataRequest = {
        ParametersIds: effectiveIds,
        StartDate: formatToISOString(props.startdate, 0),
        EndDate: formatToISOString(props.enddate, 1),
      };

      try {
        const url = `/api/getTrendData`; // FastAPI 엔드포인트
        const response = await axios.post(url, trendDataRequest, {
          headers: { "Content-Type": "application/json" },
        });

        if (response.data.success) {
          console.log("서버 응답 데이터:", response.data.data);
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
            for (let i = 0; i < resData.Thresholds.length; i++) {
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
                // threshold 시리즈 데이터 생성:
                // for each label in labels:
                // - if labelTime < t1: use value (첫 번째 threshold 값)
                // - if labelTime >= t1: use secondValue (두 번째 threshold 값)
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
          option.value = {
            lineLabels: labels,
            lineData: datasets,
          };
        } else {
          console.error("서버 오류:", response.data.error);
          option.value = {
            lineLabels: [],
            lineData: [],
          };
        }
      } catch (error) {
        console.error("요청 실패:", error);
      }
    };

    const nameToTitleMap = {};

    function buildNameToTitleMap(tree) {
      tree.forEach((node) => {
        if (node.children) {
          node.children.forEach((child) => {
            nameToTitleMap[child.Name] = child.Title;
          });
        } else {
          nameToTitleMap[node.Name] = node.Title;
        }
      });
    }

    onMounted(() => {
      fetchData();
      buildNameToTitleMap(trendTreeData);
      buildNameToTitleMap(energyTreeData);
    });

    return {
      t,
      option,
      meterOption,
      energyOption,
      pqOption,
      channel,
      tap,
      items,
      checkedIds,
      checkedNames,
      onCheckChange,
      drawMeterChart,
      drawEnergyChart,
      drawDiagnosisChart,
      trendTreeData,
    };
  },
};
</script>