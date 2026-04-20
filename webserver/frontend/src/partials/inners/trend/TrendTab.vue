<template>
  <div class="grid grid-cols-12 gap-6">
    <div class="col-span-4 flex flex-col gap-6 h-auto">
      <div v-show="true" role="alert">
        <div
          class="flex flex-col w-full px-4 py-2 rounded-lg text-sm bg-white dark:bg-gray-800 shadow-sm border border-gray-200 dark:border-gray-700/60 text-gray-600 dark:text-gray-100"
        >
          <div class="flex flex-col w-full gap-3">
            <div class="flex items-start">
              <svg
                class="shrink-0 fill-current text-violet-500 mt-[3px] mr-3"
                width="16"
                height="16"
                viewBox="0 0 16 16"
              >
                <path
                  d="M8 0C3.6 0 0 3.6 0 8s3.6 8 8 8 8-3.6 8-8-3.6-8-8-8zm1 12H7V7h2v5zM8 6c-.6 0-1-.4-1-1s.4-1 1-1 1 .4 1 1-.4 1-1 1z"
                />
              </svg>
              <div><slot />{{ t("trend.TrendTab.slot") }}.</div>
            </div>
          </div>
          <div class="text-right mt-1">
            <a v-if="tap == `Meter`"
              class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400"
              href="#0"
              :class="{ 'opacity-50 pointer-events-none': isLoading }"
              @click.prevent="drawMeterChart"
            >
              {{
                isLoading
                  ? t("trend.TrendTab.loading")
                  : t("trend.TrendTab.Plot")
              }}
            </a>
            <a
              v-else-if="tap == `Energy`"
              class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400"
              href="#0"
              :class="{ 'opacity-50 pointer-events-none': isLoading }"
              @click.prevent="drawEnergyChart"
            >
              {{
                isLoading
                  ? t("trend.TrendTab.loading")
                  : t("trend.TrendTab.Plot")
              }}
            </a>
            <a
              v-else-if="tap == `Demand`"
              class="font-medium text-violet-500 hover:text-violet-600 dark:hover:text-violet-400"
              href="#0"
              :class="{ 'opacity-50 pointer-events-none': isLoading }"
              @click.prevent="drawDemandChart"
            >
              {{
                isLoading
                  ? t("trend.TrendTab.loading")
                  : t("trend.TrendTab.Plot")
              }}
            </a>
          </div>
        </div>
      </div>

      <Trend_treetable
        v-if="items.length > 0"
        :channel="channel"
        :data="items"
        :expanded="false"
        :checked-ids="checkedIds"
        :checked-names="checkedNames"
        @check-change="onCheckChange"
      />
    </div>

    <div class="col-span-8 relative">
      <div
        v-if="isLoading"
        class="absolute inset-0 bg-white/80 dark:bg-gray-900/80 flex items-center justify-center z-10 rounded-lg"
      >
        <div class="flex flex-col items-center gap-3">
          <svg
            class="animate-spin h-12 w-12 text-violet-500"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <p class="text-lg font-medium text-gray-700 dark:text-gray-300">
            {{ t("trend.TrendTab.loading") }}
          </p>
        </div>
      </div>
      <LineChart
        v-if="tap == `Diagnosis`"
        :chart-data="option.lineData"
        :chart-labels="option.lineLabels"
        :chartLastDate="lastDate"
      />
      <LineChart
        v-if="tap == `Meter`"
        :chart-data="meterOption.lineData"
        :chart-labels="meterOption.lineLabels"
        :chartLastDate="lastDate"
      />
      <LineChart
        v-if="tap == `PowerQuality`"
        :chart-data="option.lineData"
        :chart-labels="option.lineLabels"
        :chartLastDate="lastDate"
      />
      <LineChart
        v-if="tap == `Parameters`"
        :chart-data="option.lineData"
        :chart-labels="option.lineLabels"
        :chartLastDate="lastDate"
      />
      <LineChart
        v-if="tap == `Energy`"
        :chart-data="energyOption.lineData"
        :chart-labels="energyOption.lineLabels"
        :chartLastDate="lastDate"
      />
      <LineChart
        v-if="tap == `Demand`"
        :chart-data="demandOption.lineData"
        :chart-labels="demandOption.lineLabels"
        :chartLastDate="lastDate"
      />
    </div>
  </div>
</template>
<script>
import { ref, onMounted, computed } from "vue";
import axios from "axios";
import Trend_treetable from "./Trend_treetable.vue";
import LineChart from "../../../charts/connect/LineChart01_Echart.vue";
import Notification from "./Notification.vue";
import {
  trendTreeData,
  energyTreeData,
  demandTreeData,
  meterParamMap,
  meterChartParamMap,
  demandParamMap,
  buildNameToTitleMap,
} from "@/composables/useTrendTree";
import { useI18n } from "vue-i18n";
import { useAuthStore } from '@/store/auth';

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
    Notification,
  },
  setup(props) {
    const { t } = useI18n();
    const channel = ref(props.channel);
    const tap = ref(props.tap);
    const items = ref([]);
    const authStore = useAuthStore();
    const saveCsv = ref(0);
    const withFreq = ref(false);
    const withCurrent = ref(false);
    const option = ref({
      lineLabels: [],
      lineData: [],
    });
    const meterOption = ref({ lineLabels: [], lineData: [] });
    const energyOption = ref({ lineLabels: [], lineData: [] });
    const demandOption = ref({ lineLabels: [], lineData: [] });
    const pqOption = ref({ lineLabels: [], lineData: [] });
    const checkedIds = ref([]);
    const checkedNames = ref([]);
    const isLoading = ref(false);
    const lastDate = ref("");

    const isNtek = computed(() => {
      const userName = authStore.getUser;
      if (userName == 'ntek')
        return true;
      else
        return false;
    });

    const expandEnabledNames = (enabledParams, paramMap) => {
      const paramNames = enabledParams.flatMap(
        (param) => paramMap[param] || []
      );

      const expanded = new Set(paramNames);

      const treeData = tap.value === "Energy" ? energyTreeData : tap.value === "Demand" ? demandTreeData : trendTreeData;

      // 3단계 트리 지원: 재귀적으로 자식 노드 확장
      const expandChildren = (nodes) => {
        nodes.forEach((node) => {
          if (paramNames.includes(node.Name)) {
            if (node.children) {
              node.children.forEach((child) => {
                expanded.add(child.Name);
                // 손자 노드도 확장
                if (child.children) {
                  child.children.forEach((grandchild) => {
                    expanded.add(grandchild.Name);
                  });
                }
              });
            }
          }
          // 재귀적으로 자식 노드 처리
          if (node.children) {
            expandChildren(node.children);
          }
        });
      };

      expandChildren(treeData);

      return Array.from(expanded);
    };

    const fetchData = async () => {
      try {
        if (tap.value == "Diagnosis") {
          if (props.asset != "") {
            const response = await axios.get(`/api/getDiagnosisDetail/${props.asset}`);
            if (response.data.success) {
              items.value = response.data.data_tree;
            }
          }
        } else if (tap.value == "PowerQuality") {
          if (props.asset != "") {
            const response = await axios.get(`/api/getDiagPQ/${props.asset}`);
            if (response.data.success) {
              items.value = response.data.data_tree;
            }
          }
        } else if (tap.value == "Parameters") {
          if (props.asset != "") {
            const response = await axios.get(
              `/api/getTrendParameters/${props.asset}/parameter`
            );
            if (response.data.success) {
              items.value = response.data.superlist;
            } else {
              // No Data
            }
          }
        } else if (tap.value === "Meter") {
          const response = await axios.get(
            `/api/getTrendParameters/${channel.value}`
          );
          if (response.data.success) {
            const enabledParams = response.data.data.params;
            // enabledParams loaded
            const filteredParams = enabledParams.filter(
              (param) =>
                ![
                  "Energy",
                  "Active Energy",
                  "Reactive Energy",
                  "Apparent Energy",
                ].includes(param)
            );

            const enabledNames = expandEnabledNames(
              filteredParams,
              meterParamMap
            );

            // 3단계 트리를 지원하는 deepFilterTree
            const deepFilterTree = (nodes) => {
              return nodes
                .map((node) => {
                  const isNodeEnabled = enabledNames.includes(node.Name);

                  if (node.children) {
                    // 재귀적으로 자식 노드 필터링 (3단계 지원)
                    const filteredChildren = deepFilterTree(node.children);

                    if (filteredChildren.length > 0 || isNodeEnabled) {
                      return {
                        ...node,
                        children: filteredChildren.length > 0 ? filteredChildren : undefined,
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
          items.value = energyTreeData;
        } else if (tap.value === "Demand") {
          items.value = demandTreeData;
        }
        // items loaded
      } catch (error) {
        // 데이터 가져오기 실패
      }
    };

    // 3단계 트리 지원: 모든 자손 ID와 이름 수집 헬퍼 함수
    function getAllDescendantIds(node) {
      const ids = [];
      if (node.children) {
        node.children.forEach((child) => {
          ids.push(child.ID);
          ids.push(...getAllDescendantIds(child));
        });
      }
      return ids;
    }

    function getAllDescendantNames(node) {
      const names = [];
      if (node.children) {
        node.children.forEach((child) => {
          names.push(child.Name);
          names.push(...getAllDescendantNames(child));
        });
      }
      return names;
    }

    // 리프 노드만 수집 (실제 데이터를 가진 노드)
    function getLeafDescendants(node, maxCount = 4) {
      const leaves = [];

      const collectLeaves = (n) => {
        if (leaves.length >= maxCount) return;

        if (!n.children || n.children.length === 0) {
          leaves.push({ id: n.ID, name: n.Name });
        } else {
          for (const child of n.children) {
            if (leaves.length >= maxCount) break;
            collectLeaves(child);
          }
        }
      };

      collectLeaves(node);
      return leaves;
    }

    function onCheckChange({
      id,
      checked,
      name,
      children = [],
      childrenNames = [],
    }) {
      // 현재 노드 찾기 (3단계 지원)
      const findNode = (nodes, targetId) => {
        for (const node of nodes) {
          if (node.ID === targetId) return node;
          if (node.children) {
            const found = findNode(node.children, targetId);
            if (found) return found;
          }
        }
        return null;
      };

      const currentNode = findNode(items.value, id);

      if (checked) {
        if (!checkedIds.value.includes(id)) checkedIds.value.push(id);
        if (!checkedNames.value.includes(name)) checkedNames.value.push(name);

        if (currentNode) {
          // 리프 노드만 최대 3개 체크 (3단계 트리에서 실제 데이터 노드)
          const leafNodes = getLeafDescendants(currentNode, 4);
          leafNodes.forEach(({ id: leafId, name: leafName }) => {
            if (!checkedIds.value.includes(leafId)) checkedIds.value.push(leafId);
            if (!checkedNames.value.includes(leafName)) checkedNames.value.push(leafName);
          });

          // 중간 노드들도 체크 (부모-자식 관계 유지)
          const allDescendantIds = getAllDescendantIds(currentNode);
          const allDescendantNames = getAllDescendantNames(currentNode);

          // 리프 노드의 부모들도 체크 상태 유지
          allDescendantIds.forEach((descId) => {
            const descNode = findNode(items.value, descId);
            if (descNode && descNode.children) {
              // 중간 노드: 자식 중 하나라도 체크되어 있으면 체크
              const hasCheckedChild = descNode.children.some((child) =>
                checkedIds.value.includes(child.ID)
              );
              if (hasCheckedChild && !checkedIds.value.includes(descId)) {
                checkedIds.value.push(descId);
              }
            }
          });
        } else {
          // 기존 로직 유지 (children, childrenNames가 전달된 경우)
          children.slice(0, 4).forEach((cid) => {
            if (!checkedIds.value.includes(cid)) checkedIds.value.push(cid);
          });
          childrenNames.slice(0, 4).forEach((cname) => {
            if (!checkedNames.value.includes(cname)) checkedNames.value.push(cname);
          });
        }
      } else {
        // 체크 해제: 현재 노드와 모든 자손 노드 체크 해제
        const allDescendantIds = currentNode ? getAllDescendantIds(currentNode) : children;
        const allDescendantNames = currentNode ? getAllDescendantNames(currentNode) : childrenNames;

        checkedIds.value = checkedIds.value.filter(
          (item) => item !== id && !allDescendantIds.includes(item)
        );
        checkedNames.value = checkedNames.value.filter(
          (item) => item !== name && !allDescendantNames.includes(item)
        );
      }

      syncParentCheck();
    }

    // 3단계 트리 지원: 부모 체크 상태 동기화
    function syncParentCheck() {
      const traverse = (node) => {
        if (node.children) {
          // 먼저 자식 노드들을 재귀적으로 처리
          node.children.forEach(traverse);

          // 모든 자식이 체크되었는지 확인
          const allChildrenChecked = node.children.every((child) =>
            checkedIds.value.includes(child.ID)
          );

          // 자식 중 일부라도 체크되었는지 확인
          const someChildrenChecked = node.children.some((child) =>
            checkedIds.value.includes(child.ID)
          );

          if (allChildrenChecked) {
            if (!checkedIds.value.includes(node.ID)) {
              checkedIds.value.push(node.ID);
            }
            if (!checkedNames.value.includes(node.Name)) {
              checkedNames.value.push(node.Name);
            }
          } else if (!someChildrenChecked) {
            // 자식이 하나도 체크되지 않으면 부모도 체크 해제
            checkedIds.value = checkedIds.value.filter((id) => id !== node.ID);
            checkedNames.value = checkedNames.value.filter((name) => name !== node.Name);
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
      const month = pad(date.getMonth() + 1);
      const day = pad(date.getDate());
      let hours, minutes, seconds;

      if (soe === 0) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
      } else if (soe === 1) {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
      } else if (soe === 2) {
        hours = pad(0);
        minutes = pad(0);
        seconds = pad(0);
      } else {
        hours = pad(23);
        minutes = pad(59);
        seconds = pad(59);
      }

      const timezoneOffset = -date.getTimezoneOffset();
      const offsetSign = timezoneOffset >= 0 ? "+" : "-";
      const offsetHours = pad(Math.abs(Math.floor(timezoneOffset / 60)));
      const offsetMinutes = pad(Math.abs(timezoneOffset % 60));

      const milliseconds = String(date.getMilliseconds()).padStart(3, "0");
      return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}.${milliseconds}${offsetSign}${offsetHours}:${offsetMinutes}`;
    };

    const drawMeterChart = async () => {
      if (!checkedNames.value || checkedNames.value.length === 0) {
        meterOption.value = { lineLabels: [], lineData: [] };
        return;
      }

      // 선택된 필드들을 하나의 배열로 수집
      let selectedFields = [];
      let actualParamCount = 0;

      checkedNames.value.forEach((param) => {
        const keys = meterChartParamMap[param];
        if (keys) {
          actualParamCount += keys.length;
          selectedFields = selectedFields.concat(keys);
        }
      });

      if (actualParamCount > 4) {
        alert(t("trend.Linechart.metercount"));
        return;
      }

      isLoading.value = true;

      try {
        const url = `/api/getMeterTrend/${channel.value}/${saveCsv.value}`;

        const response = await axios.post(url, {
          startDate: formatToISOString(props.startdate, 2),
          endDate: formatToISOString(props.enddate, 3),
          fields: selectedFields,
        });

        const responseData = response.data.data;
        lastDate.value = response.data.date;
        const datasets = [];
        const labels = [];
        const selectedParams = checkedNames.value;

        labels.push(...responseData.map((row) => row._time));

        selectedParams.forEach((param) => {
          const keys = meterChartParamMap[param];
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

        meterOption.value = {
          lineLabels: labels,
          lineData: datasets,
        };
      } catch (error) {
        // 데이터 가져오기 실패
        if (error.response?.data?.error) {
          alert(`Error: ${error.response.data.error}`);
        }
      } finally {
        isLoading.value = false;
      }
    };

    const drawEnergyChart = async () => {
      if (!checkedNames.value || checkedNames.value.length === 0) {
        energyOption.value = { lineLabels: [], lineData: [] };
        return;
      }

      isLoading.value = true;

      try {
        let url = `/api/getEnergyTrend/${channel.value}`;
        const response = await axios.get(url, {
          params: {
            startDate: formatToISOString(props.startdate, 2),
            endDate: formatToISOString(props.enddate, 3),
          },
        });

        if (response.data.result) {
          const responseData = response.data.data;
          lastDate.value = response.data.date;
          if (!Array.isArray(responseData)) {
            // Energy 응답 데이터가 배열이 아님
            energyOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          if (responseData.length === 0) {
            // Energy 데이터가 비어있음
            energyOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          const datasets = [];
          const labels = [];
          const selectedParams = checkedNames.value;

          try {
            labels.push(...responseData.map((row) => row._time));
          } catch (timeError) {
            // 시간 라벨 생성 실패
            energyOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          const energyParamMap = {
            kwh_import_consumption: ["kwh_import_consumption"],
            kvarh_import_consumption: ["kvarh_import_consumption"],
            kvah_import_consumption: ["kvah_import_consumption"],
          };

          selectedParams.forEach((param) => {
            const keys = energyParamMap[param];
            if (!keys) {
              // Energy 파라미터 매핑을 찾을 수 없음
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
              } catch (dataError) {
                // Energy 데이터 생성 실패
              }
            });
          });

          energyOption.value = {
            lineLabels: labels,
            lineData: datasets,
          };
        } else {
          // Energy API 응답 실패
          energyOption.value = { lineLabels: [], lineData: [] };
        }
      } catch (error) {
        // Energy 데이터 가져오기 실패
        energyOption.value = { lineLabels: [], lineData: [] };
      } finally {
        isLoading.value = false;
      }
    };

    // Demand 차트 그리기
    const drawDemandChart = async () => {
      if (!checkedNames.value || checkedNames.value.length === 0) {
        demandOption.value = { lineLabels: [], lineData: [] };
        return;
      }

      isLoading.value = true;

      try {
        let url = `/api/getDemandTrend/${channel.value}`;
        const response = await axios.get(url, {
          params: {
            startDate: formatToISOString(props.startdate, 2),
            endDate: formatToISOString(props.enddate, 3),
          },
        });

        if (response.data.result) {
          const responseData = response.data.data;
          lastDate.value = response.data.date;
          if (!Array.isArray(responseData)) {
            demandOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          if (responseData.length === 0) {
            demandOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          const datasets = [];
          const labels = [];
          const selectedParams = checkedNames.value;

          try {
            labels.push(...responseData.map((row) => row._time));
          } catch (timeError) {
            demandOption.value = { lineLabels: [], lineData: [] };
            return;
          }

          selectedParams.forEach((param) => {
            const keys = demandParamMap[param];
            if (!keys) {
              return;
            }

            keys.forEach((key) => {
              try {
                const data = responseData.map((row) => row[key] != null ? row[key] / 1000 : row[key]);
                datasets.push({
                  name: nameToTitleMap[key] || key,
                  data: data,
                  isThreshold: false,
                });
              } catch (dataError) {
                // Demand 데이터 생성 실패
              }
            });
          });

          demandOption.value = {
            lineLabels: labels,
            lineData: datasets,
          };
        } else {
          demandOption.value = { lineLabels: [], lineData: [] };
        }
      } catch (error) {
        demandOption.value = { lineLabels: [], lineData: [] };
      } finally {
        isLoading.value = false;
      }
    };

    // 3단계 트리 지원: 노드가 부모인지 확인 (자식이 있는지)
    function isParentNode(node) {
      return node && node.children && node.children.length > 0;
    }

    function isParentId(id) {
      const findNode = (nodes) => {
        for (const node of nodes) {
          if (node.ID === id) {
            return isParentNode(node);
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

    // 3단계 트리 지원: 리프 노드의 ID만 수집
    function getEffectiveParameterIds() {
      const result = [];

      const collectLeafIds = (nodes) => {
        nodes.forEach((node) => {
          if (node.children && node.children.length > 0) {
            // 자식이 있으면 재귀적으로 처리
            collectLeafIds(node.children);
          } else {
            // 리프 노드이고 체크되어 있으면 추가
            if (checkedIds.value.includes(node.ID)) {
              result.push(node.ID);
            }
          }
        });
      };

      collectLeafIds(items.value);
      return result;
    }

    // 3단계 트리 지원: 리프 노드의 AssemblyID와 Name 수집
    function getEffectiveParameters() {
      const result = [];

      const collectLeafParams = (nodes) => {
        nodes.forEach((node) => {
          if (node.children && node.children.length > 0) {
            // 자식이 있으면 재귀적으로 처리
            collectLeafParams(node.children);
          } else {
            // 리프 노드이고 체크되어 있으면 추가
            if (checkedIds.value.includes(node.ID)) {
              result.push({
                AssemblyID: node.AssemblyID || "",
                Name: node.Name || ""
              });
            }
          }
        });
      };

      collectLeafParams(items.value);
      return result;
    }

    const drawDiagnosisChart = async () => {
      const effectiveIds = getEffectiveParameterIds();

      if (!effectiveIds || effectiveIds.length === 0) {
        option.value = {
          lineLabels: [],
          lineData: [],
        };
        return;
      }

      if (effectiveIds.length > 4) {
        alert(t("trend.Linechart.parametercount"));
        return;
      }

      isLoading.value = true;

      const trendDataRequest = {
        ParametersIds: effectiveIds,
        StartDate: formatToISOString(props.startdate, 0),
        EndDate: formatToISOString(props.enddate, 1),
      };
      // 트렌드 데이터 요청
      try {
        const url = `/api/getTrendData`;
        const response = await axios.post(url, trendDataRequest, {
          headers: { "Content-Type": "application/json" },
        });

        if (response.data.success) {
          const resData = response.data.data;
          lastDate.value = response.data.date;
          let datasets = [];
          let labels = [];

          Object.keys(resData).forEach((key) => {
            if (key !== "Thresholds") {
              const dataPoints = resData[key].data;
              if (dataPoints && dataPoints.length > 0) {
                if (labels.length === 0) {
                  labels = dataPoints.map((point) => point.XAxis);
                }
                datasets.push({
                  name: resData[key].Title,
                  data: dataPoints.map((point) => point.YAxis),
                  isThreshold: false,
                });
              }
            }
          });

          if (
            resData.Thresholds &&
            resData.Thresholds.length > 0 &&
            labels.length > 0
          ) {
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

            if (resData.Thresholds[0].Thresholds != null) {
              const thresholdCount = resData.Thresholds[0].Thresholds.length;

              for (let idx = 0; idx < thresholdCount; idx++) {
                const hasValidValue = resData.Thresholds.some((t) => {
                  const value = t.Thresholds[idx];
                  return (
                    value !== "NaN" &&
                    value !== null &&
                    value !== undefined &&
                    typeof value === "number"
                  );
                });

                if (!hasValidValue) continue;

                const timeList = resData.Thresholds.filter((t) => {
                  const value = t.Thresholds[idx];
                  return (
                    value !== "NaN" &&
                    value !== null &&
                    value !== undefined &&
                    typeof value === "number"
                  );
                })
                  .map((t) => ({
                    time: new Date(t.XAxis),
                    value: t.Thresholds[idx],
                  }))
                  .sort((a, b) => a.time - b.time);

                if (timeList.length === 0) continue;

                const thresholdData = labels.map((lbl) => {
                  const labelTime = new Date(lbl);

                  let applicableThreshold = timeList[0].value;

                  for (let i = 0; i < timeList.length; i++) {
                    if (labelTime >= timeList[i].time) {
                      applicableThreshold = timeList[i].value;
                    } else {
                      break;
                    }
                  }

                  return applicableThreshold;
                });

                datasets.push({
                  name: ThresholdString[idx],
                  data: thresholdData,
                  isThreshold: true,
                });
              }
            }
          }

          option.value = {
            lineLabels: labels,
            lineData: datasets,
          };
        } else {
          // 서버 오류
          option.value = {
            lineLabels: [],
            lineData: [],
          };
        }
      } catch (error) {
        // 요청 실패
      } finally {
        isLoading.value = false;
      }
    };

    /**
     * drawDiagnosisChartByName - API 형식에 맞춰 수정
     * 요청 형식:
     * {
     *   "AssetName": "fan",
     *   "StartDate": "2025-12-01T21:28:41.9779368-05:00",
     *   "EndDate": "2025-12-24T21:28:41.9779445-04:00",
     *   "Parameters": [
     *     { "AssemblyID": "PWS", "Name": "Load" },
     *     { "AssemblyID": "PWS", "Name": "LoadFLA" }
     *   ]
     * }
     */
    const drawDiagnosisChartByName = async () => {
      // 리프 노드에서 AssemblyID와 Name 수집
      const effectiveParams = getEffectiveParameters();

      if (withFreq.value) {
        effectiveParams.push({
          AssemblyID: "VOL",
          Name: "LF"
        });
      }
      // 전류 체크박스가 선택되었을 때
      if (withCurrent.value) {
        effectiveParams.push({
          AssemblyID: "PWS",
          Name: "rmsiA"
        });
        effectiveParams.push({
          AssemblyID: "PWS",
          Name: "rmsiB"
        });
        effectiveParams.push({
          AssemblyID: "PWS",
          Name: "rmsiC"
        });
      }

      // 선택된 파라미터
      // 현재 items
      // 체크된 IDs

      if (!effectiveParams || effectiveParams.length === 0) {
        option.value = {
          lineLabels: [],
          lineData: [],
        };
        return;
      }

      if (effectiveParams.length > 8) {
        alert(t("trend.Linechart.parametercount"));
        return;
      }

      isLoading.value = true;

      // API 요청 형식에 맞게 구성
      const trendDataRequest = {
        AssetName: props.asset,
        StartDate: formatToISOString(props.startdate, 0),
        EndDate: formatToISOString(props.enddate, 1),
        Parameters: effectiveParams,
        ParametersIds: []  // 기존 TrendData 모델 호환을 위한 더미 필드
      };

      // 트렌드 데이터 요청 ByName

      try {
        const url = `/api/getTrendbyName`;
        const response = await axios.post(url, trendDataRequest, {
          headers: { "Content-Type": "application/json" },
        });

        if (response.data.success) {
          // 서버 응답 수신
          const resData = response.data.data;
          
          // lastDate 형식 변환 (ISO 형식 → YYYY-MM-DD HH:MM:SS)
          if (response.data.date) {
            const dateStr = response.data.date;
            try {
              const dateObj = new Date(dateStr);
              if (!isNaN(dateObj.getTime())) {
                const year = dateObj.getFullYear();
                const month = String(dateObj.getMonth() + 1).padStart(2, '0');
                const day = String(dateObj.getDate()).padStart(2, '0');
                const hours = String(dateObj.getHours()).padStart(2, '0');
                const minutes = String(dateObj.getMinutes()).padStart(2, '0');
                const seconds = String(dateObj.getSeconds()).padStart(2, '0');
                lastDate.value = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
              } else {
                lastDate.value = dateStr;
              }
            } catch (e) {
              lastDate.value = dateStr;
            }
          } else {
            lastDate.value = "";
          }
          
          let datasets = [];
          let labels = [];

          Object.keys(resData).forEach((key) => {
            if (key !== "Thresholds") {
              const dataPoints = resData[key].data;
              if (dataPoints && dataPoints.length > 0) {
                if (labels.length === 0) {
                  labels = dataPoints.map((point) => point.XAxis);
                }
                datasets.push({
                  name: resData[key].Title,
                  data: dataPoints.map((point) => point.YAxis),
                  isThreshold: false,
                });
              }
            }
          });

          // Threshold 처리
          if (
            resData.Thresholds &&
            resData.Thresholds.length > 0 &&
            labels.length > 0
          ) {
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

            if (resData.Thresholds[0].Thresholds != null) {
              const thresholdCount = resData.Thresholds[0].Thresholds.length;

              for (let idx = 0; idx < thresholdCount; idx++) {
                const hasValidValue = resData.Thresholds.some((t) => {
                  const value = t.Thresholds[idx];
                  return (
                    value !== "NaN" &&
                    value !== null &&
                    value !== undefined &&
                    typeof value === "number"
                  );
                });

                if (!hasValidValue) continue;

                const timeList = resData.Thresholds.filter((t) => {
                  const value = t.Thresholds[idx];
                  return (
                    value !== "NaN" &&
                    value !== null &&
                    value !== undefined &&
                    typeof value === "number"
                  );
                })
                  .map((t) => ({
                    time: new Date(t.XAxis),
                    value: t.Thresholds[idx],
                  }))
                  .sort((a, b) => a.time - b.time);

                if (timeList.length === 0) continue;

                const thresholdData = labels.map((lbl) => {
                  const labelTime = new Date(lbl);
                  let applicableThreshold = timeList[0].value;

                  for (let i = 0; i < timeList.length; i++) {
                    if (labelTime >= timeList[i].time) {
                      applicableThreshold = timeList[i].value;
                    } else {
                      break;
                    }
                  }

                  return applicableThreshold;
                });

                datasets.push({
                  name: ThresholdString[idx],
                  data: thresholdData,
                  isThreshold: true,
                });
              }
            }
          }

          option.value = {
            lineLabels: labels,
            lineData: datasets,
          };
        } else {
          // 서버 오류
          option.value = {
            lineLabels: [],
            lineData: [],
          };
        }
      } catch (error) {
        // 요청 실패
        // 에러 상세
        // 에러 상태
        option.value = {
          lineLabels: [],
          lineData: [],
        };
      } finally {
        isLoading.value = false;
      }
    };

    const nameToTitleMap = buildNameToTitleMap(trendTreeData, energyTreeData, demandTreeData);

    onMounted(() => {
      fetchData();
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
      isLoading,
      onCheckChange,
      drawMeterChart,
      drawEnergyChart,
      drawDemandChart,
      demandOption,
      drawDiagnosisChart,
      drawDiagnosisChartByName,
      trendTreeData,
      lastDate,
      saveCsv,
      isNtek,
      withFreq,
      withCurrent,
    };
  },
};
</script>