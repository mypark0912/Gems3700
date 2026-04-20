<template>
  <div
    class="relative col-span-full xl:col-span-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5"
      :class="type === 'diagnostic' ? 'bg-purple-500' : 'bg-purple-500'"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center mb-2">
        <div 
          class="w-6 h-6 rounded-full shrink-0 mr-3"
          :class="type === 'diagnostic' ? 'bg-purple-500' : 'bg-purple-500'"
        >
          <!-- Diagnostic/PQ 아이콘 -->
          <svg
            v-if="false"
            class="w-6 h-6 fill-current text-white"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <circle
              cx="17"
              cy="17"
              r="4"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
            />
            <path
              d="M20 20l2 2"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          <!-- Faults/Events 아이콘 -->
          <svg
            
            class="w-6 h-6 fill-current text-white"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </div>
        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
          {{ type === 'diagnostic' ? 'Equipment / PQ Status' : 'Faults / Events Status' }}
        </h3>
      </header>
    </div>

    <div class="px-4 py-3 space-y-4">
      <!-- 탭 네비게이션 -->
      <div class="flex border-b border-gray-200 dark:border-gray-700/60">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors duration-200',
            activeTab === tab.id
              ? 'text-violet-600 dark:text-violet-400 border-b-2 border-violet-600 dark:border-violet-400'
              : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
          ]"
        >
          {{ tab.label }}
        </button>
      </div>

      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-3 border border-gray-200 dark:border-gray-700/60"
      >
        <div class="mt-6 pt-2">
          <div class="flex flex-col">
            <!-- 헤더 -->
            <div class="grid grid-cols-[minmax(180px,1fr)_80px_120px] gap-4 text-xs font-semibold text-gray-500 dark:text-gray-400 px-1 mb-4 border-b pb-2">
              <div class="text-left">Parameter</div>
              <div class="text-center">Enable</div>
              <div class="text-center">Select Level</div>
            </div>

            <!-- 데이터가 없을 때 -->
            <div v-if="currentTabData.length === 0" class="text-center py-8 text-gray-400">
              <p class="text-sm">No items available for this asset type</p>
            </div>

            <!-- 스크롤 가능한 데이터 목록 -->
            <div v-else class="overflow-y-auto max-h-[400px] space-y-2 pr-2">
              <div
                v-for="(item, idx) in currentTabData"
                :key="idx"
                class="grid grid-cols-[minmax(180px,1fr)_80px_120px] gap-4 items-center border-b border-gray-200 dark:border-gray-700/60 py-3 text-sm px-1 hover:bg-gray-50 dark:hover:bg-gray-700/30"
              >
                <!-- Parameter Name -->
                <div class="text-left text-xs text-gray-800 dark:text-gray-200 flex items-center min-w-0">
                  <span class="truncate">{{ item.name }}</span>
                </div>
                
                <!-- Enable Checkbox -->
                <div class="flex justify-center w-[80px]">
                  <input
                    type="checkbox"
                    v-model="item.enabled"
                    class="w-4 h-4 cursor-pointer text-violet-600 focus:ring-violet-500 border-gray-300 rounded"
                  />
                </div>
                
                <!-- Select Level Combobox -->
                <div class="flex justify-center w-[120px]">
                  <select
                    v-model="item.level"
                    :disabled="!item.enabled"
                    class="w-24 px-3 py-1.5 text-xs border rounded-md transition-all"
                    :class="
                      !item.enabled
                        ? 'bg-gray-100 dark:bg-gray-800 text-gray-400 dark:text-gray-500 border-gray-300 dark:border-gray-600 cursor-not-allowed'
                        : 'bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500'
                    "
                  >
                    <template v-if="activeTab === 'diagnostic'">
                      <option :value="2">{{ t("diagnosis.tabContext.st2") }}</option>
                      <option :value="3">{{ t("diagnosis.tabContext.st3") }}</option>
                      <option :value="4">{{ t("diagnosis.tabContext.st4") }}</option>
                    </template>
                    <template v-else>
                      <option :value="2">{{ type === 'diagnostic'? t("diagnosis.tabContext.pqfe2") : t("diagnosis.tabContext.fe2") }}</option>
                      <option :value="3">{{ type === 'diagnostic'? t("diagnosis.tabContext.pqfe3") : t("diagnosis.tabContext.fe3") }}</option>
                      <option :value="4">{{ type === 'diagnostic'? t("diagnosis.tabContext.pqfe4") : t("diagnosis.tabContext.fe4") }}</option>
                    </template>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject, watch, onMounted, nextTick } from "vue";
import { useI18n } from "vue-i18n";
const isDataLoaded = ref(false);
const { t } = useI18n();
const props = defineProps({
  channel: { type: String, default: '' },
  type: { type: String, default: 'diagnostic' } // 'diagnostic' | 'faults'
});

const mainData = inject('channel_main');
const subData = inject('channel_sub');

// 탭 설정 - type에 따라 다른 탭
const tabs = computed(() => {
  if (props.type === 'diagnostic') {
    return [
      { id: 'diagnostic', label: 'Equipment' },
      { id: 'pq', label: 'PQ' }
    ];
  } else {
    return [
      { id: 'faults', label: 'Faults' },
      { id: 'events', label: 'Events' }
    ];
  }
});

const activeTab = ref(props.type === 'diagnostic' ? 'diagnostic' : 'faults');

const stDict = computed(() => {
  if (props.channel === 'Main' && mainData.value) {
    return mainData.value.status_Info;
  } else if (props.channel === 'Sub' && subData.value) {
    return subData.value.status_Info;
  }
  return { diagnosis: [], pq: [], faults: [], events: [] };
});

const AssetType = computed(() => {
  if (props.channel === 'Main' && mainData.value) {
    return mainData.value.assetInfo?.type || '';
  } else if (props.channel === 'Sub' && subData.value) {
    return subData.value.assetInfo?.type || '';
  }
  return '';
});

// 데이터 저장소
const diagnosticData = ref([]);
const pqData = ref([]);
const faultsData = ref([]);
const eventsData = ref([]);

const currentTabData = computed(() => {
  switch (activeTab.value) {
    case 'diagnostic':
      return diagnosticData.value;
    case 'pq':
      return pqData.value;
    case 'faults':
      return faultsData.value;
    case 'events':
      return eventsData.value;
    default:
      return [];
  }
});

const fetchData = async () => {
  if (!AssetType.value || AssetType.value === '') {
    diagnosticData.value = [];
    pqData.value = [];
    faultsData.value = [];
    eventsData.value = [];
    return false;
  }

  try {
    const res = await fetch('/diagnosis_items.json');
    const data = await res.json();
    
    if (props.type === 'diagnostic') {
      // Diagnostic/PQ 데이터만 로드
      const diagnosisItems = data.Diagnosis[AssetType.value] || [];
      const pqItems = data.PQ || [];
      
      diagnosticData.value = diagnosisItems.map(item => ({
        name: item.Title,
        component: item.Component,
        enabled: false,
        level: 2
      }));
      
      pqData.value = pqItems.map(item => ({
        name: item.Title,
        component: item.Component,
        enabled: false,
        level: 2
      }));
    } else {
      // Faults/Events 데이터만 로드
      const faultsItems = data.Faults || [];
      const eventsItems = data.Events || [];
      
      faultsData.value = faultsItems.map(item => ({
        name: item.Title,
        component: item.Component,
        enabled: false,
        level: 2
      }));
      
      eventsData.value = eventsItems.map(item => ({
        name: item.Title,
        component: item.Component,
        enabled: false,
        level: 2
      }));
    }
    
    return true;
  } catch (error) {
    console.error("데이터 가져오기 실패:", error);
    return false;
  }
};

const loadSavedData = () => {
  const statusInfo = stDict.value;
  if (!statusInfo) return;

  if (props.type === 'diagnostic') {
    // Diagnostic 데이터 초기화 및 로드
    diagnosticData.value.forEach(item => {
      item.enabled = false;
      item.level = 2;
    });
    pqData.value.forEach(item => {
      item.enabled = false;
      item.level = 2;
    });

    const diagnosisArray = statusInfo.diagnosis || [];
    diagnosisArray.forEach(savedItem => {
      const item = diagnosticData.value.find(d => d.name === savedItem.name);
      if (item) {
        item.enabled = true;
        item.level = savedItem.level;
      }
    });

    const pqArray = statusInfo.pq || [];
    pqArray.forEach(savedItem => {
      const item = pqData.value.find(p => p.name === savedItem.name);
      if (item) {
        item.enabled = true;
        item.level = savedItem.level;
      }
    });
  } else {
    // Faults/Events 데이터 초기화 및 로드
    faultsData.value.forEach(item => {
      item.enabled = false;
      item.level = 2;
    });
    eventsData.value.forEach(item => {
      item.enabled = false;
      item.level = 2;
    });

    const faultsArray = statusInfo.faults || [];
    faultsArray.forEach(savedItem => {
      const item = faultsData.value.find(f => f.name === savedItem.name);
      if (item) {
        item.enabled = true;
        item.level = savedItem.level;
      }
    });

    const eventsArray = statusInfo.events || [];
    eventsArray.forEach(savedItem => {
      const item = eventsData.value.find(e => e.name === savedItem.name);
      if (item) {
        item.enabled = true;
        item.level = savedItem.level;
      }
    });
  }
};

const updateInputDict = () => {
  const statusInfo = stDict.value;
  if (!statusInfo) return;

  if (props.type === 'diagnostic') {
    statusInfo["diagnosis"] = diagnosticData.value
      .filter(item => item.enabled)
      .map(item => ({ name: item.name, level: item.level }));

    statusInfo["pq"] = pqData.value
      .filter(item => item.enabled)
      .map(item => ({ name: item.name, level: item.level }));
  } else {
    statusInfo["faults"] = faultsData.value
      .filter(item => item.enabled)
      .map(item => ({ name: item.name, level: item.level }));

    statusInfo["events"] = eventsData.value
      .filter(item => item.enabled)
      .map(item => ({ name: item.name, level: item.level }));
  }
};

// Watchers
watch(AssetType, async (newType, oldType) => {
  if (newType !== oldType) {
    await fetchData();
    if (isDataLoaded.value) {
      loadSavedData();
    }
  }
});

watch([diagnosticData, pqData, faultsData, eventsData], () => {
  if (isDataLoaded.value) {
    updateInputDict();
  }
}, { deep: true });

onMounted(async () => {
  await fetchData();
  
  let retries = 0;
  const maxRetries = 10;
  
  while (retries < maxRetries) {
    const data = props.channel === 'Main' ? mainData.value : subData.value;
    
    if (data && data.status_Info) {
      await nextTick();
      loadSavedData();
      isDataLoaded.value = true;
      break;
    }
    
    await new Promise(resolve => setTimeout(resolve, 100));
    retries++;
  }
  
  if (retries === maxRetries) {
    isDataLoaded.value = true;
  }
});
</script>

<style scoped>
.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  @apply bg-gray-100 dark:bg-gray-800;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  @apply bg-gray-300 dark:bg-gray-600 rounded;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  @apply bg-gray-400 dark:bg-gray-500;
}

select:focus {
  outline: none;
}

input[type="checkbox"]:focus {
  outline: none;
}
</style>