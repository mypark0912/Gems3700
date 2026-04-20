<template>
  <div
    v-if="isUseAIEnabled"
    class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center mb-2">
        <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
          <svg
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
        </div>
        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
          AI Module MODBUS Configuration
        </h3>
      </header>
    </div>

    <div class="px-4 py-3 space-y-4">
      <!-- 버튼 영역 -->
      <div class="flex space-x-2">
        <button
          @click="addRow"
          class="px-3 py-1.5 text-sm font-medium text-white bg-green-600 hover:bg-green-700 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
        >
          + Add Row
        </button>
        <button
          @click="handleReset"
          class="px-3 py-1.5 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-violet-500"
        >
          Reset
        </button>
      </div>

      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700/60"
      >
        <!-- 테이블 -->
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b-2 border-gray-300 dark:border-gray-600">
                <th
                  class="px-2 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50 w-14"
                >
                  Index
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Enable
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Dev ID
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Start<br />Address
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Count
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Data<br />Type
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Destination
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Period<br />(ms)
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  FC
                </th>
                <th
                  class="px-3 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50"
                >
                  Scale
                </th>
                <th
                  class="px-2 py-3 text-center text-xs font-semibold text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700/50 w-20"
                >
                  Action
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(row, idx) in modbusConfig"
                :key="row.id"
                class="border-b border-gray-200 dark:border-gray-700/60 hover:bg-gray-50 dark:hover:bg-gray-700/30"
              >
                <!-- Index -->
                <td
                  class="px-2 py-2 text-center text-xs text-gray-800 dark:text-gray-200 font-medium"
                >
                  {{ idx + 1 }}
                </td>
                <!-- Enable -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="checkbox"
                    :checked="row.enable === 1"
                    @change="row.enable = $event.target.checked ? 1 : 0"
                    class="w-4 h-4 text-violet-600 border-gray-300 rounded focus:ring-violet-500"
                  />
                </td>
                <!-- Dev ID -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="number"
                    v-model.number="row.devId"
                    class="w-16 px-1.5 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                    min="0"
                    max="255"
                  />
                </td>
                <!-- Start Address -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="number"
                    v-model.number="row.startAddr"
                    class="w-20 px-1.5 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                    min="0"
                  />
                </td>
                <!-- Count -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="number"
                    v-model.number="row.count"
                    class="w-16 px-1.5 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                    min="0"
                  />
                </td>
                <!-- Data Type -->
                <td class="px-2 py-2 text-center">
                  <select
                    v-model.number="row.dataType"
                    class="w-20 px-1 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="0">UINT16</option>
                    <option :value="1">UINT32</option>
                    <option :value="2">Float</option>
                    <option :value="3">F002</option>
                  </select>
                </td>
                <!-- Destination -->
                <td class="px-2 py-2 text-center">
                  <select
                    v-model.number="row.destination"
                    class="w-24 px-1 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="0">{{t("config.channelPanel.tempR")}}</option>
                    <option :value="1">{{t("config.channelPanel.tempS")}}</option>
                    <option :value="2">{{t("config.channelPanel.tempT")}}</option>
                  </select>
                </td>
                <!-- Period (ms) -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="number"
                    v-model.number="row.period"
                    class="w-20 px-1.5 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                    min="0"
                    step="100"
                  />
                </td>
                <!-- FC -->
                <td class="px-2 py-2 text-center">
                  <select
                    v-model.number="row.fc"
                    class="w-16 px-1 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                  >
                    <option :value="0">-</option>
                    <option :value="1">01</option>
                    <option :value="2">02</option>
                    <option :value="3">03</option>
                    <option :value="4">04</option>
                    <option :value="5">05</option>
                    <option :value="6">06</option>
                    <option :value="15">15</option>
                    <option :value="16">16</option>
                  </select>
                </td>
                <!-- Scale -->
                <td class="px-2 py-2 text-center">
                  <input
                    type="number"
                    v-model.number="row.scale"
                    class="w-20 px-1.5 py-1 text-xs text-center border rounded-md bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 border-gray-300 dark:border-gray-500 focus:ring-violet-500 focus:border-violet-500"
                    step="0.01"
                  />
                </td>
                <!-- Action -->
                <td class="px-2 py-2 text-center">
                  <button
                    @click="removeRow(idx)"
                    :disabled="modbusConfig.length <= 1"
                    class="px-2 py-1 text-xs font-medium text-white bg-red-600 hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed rounded-md focus:outline-none focus:ring-2 focus:ring-red-500"
                    title="Delete row"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 설명 -->
        <div
          class="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md"
        >
          <p class="text-xs text-blue-800 dark:text-blue-200">
            <strong>Note:</strong>
            • FC: Function Code (03=Read Holding Registers, 04=Read Input Registers, 16=Write Multiple Registers)
            <br />
            • Data Type: UINT16, UINT32, Float, F002 (F002=Float value/100)
            <br />
            • Destination: 목적지 주소 (온도 R상, 온도 S상, 온도 T상)
            <br />
            • Scale: 스케일 값 (예: 0.01, 0.1, 1.0, 100.0)
          </p>
        </div>

        <!-- 하단 정보 -->
        <div class="flex justify-between items-center mt-6">
          <div class="text-sm text-gray-600 dark:text-gray-400">
            Total Configurations:
            <span class="font-semibold">{{ modbusConfig.length }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, computed, watch, onMounted } from "vue";
import { useI18n } from "vue-i18n";
// 고유 ID 생성용 카운터
let nextId = 1;
const { t } = useI18n();
const props = defineProps({
  channel: { type: String, default: '' }
});

const mainData = inject('channel_main');
const subData = inject('channel_sub');

// 현재 채널 데이터 computed
const currentChannelData = computed(() => {
  return props.channel === 'Main' ? mainData.value : subData.value;
});

// Modbus 설정 데이터
const modbusConfig = ref([]);

// 기본 행 생성 함수
const createDefaultRow = (id) => ({
  id: id,
  enable: 0,
  devId: 0,
  startAddr: 0,
  count: 0,
  dataType: 0,
  destination: 0,
  period: 1000,
  fc: 0,
  scale: 1.0,
  min: 0.0,
  max: 500.0
});

// 초기 기본값
const defaultConfig = [
  {
    id: 1,
    enable: 1,
    devId: 100,
    startAddr: 14,
    count: 28,
    dataType: 2,
    destination: 0,
    period: 1000,
    fc: 3,
    scale: 1.0,
  },
  {
    id: 2,
    enable: 1,
    devId: 100,
    startAddr: 90,
    count: 12,
    dataType: 2,
    destination: 1,
    period: 1000,
    fc: 3,
    scale: 1.0,
  },
];

// ai_modbus에서 데이터 로드
const loadFromChannel = () => {
  const aiModbus = currentChannelData.value?.ai_modbus;
  if (aiModbus && Array.isArray(aiModbus) && aiModbus.length > 0) {
    // 기존 데이터가 있으면 로드
    modbusConfig.value = aiModbus.map((item, idx) => ({
      id: item.id || idx + 1,
      enable: item.enable ?? 0,
      devId: item.devId ?? 0,
      startAddr: item.startAddr ?? 0,
      count: item.count ?? 0,
      dataType: item.dataType ?? 0,
      destination: item.destination ?? 0,
      period: item.period ?? 1000,
      fc: item.fc ?? 0,
      scale: item.scale ?? 1.0,
    }));
    nextId = Math.max(...modbusConfig.value.map(r => r.id)) + 1;
  } else {
    // 없으면 기본값 사용
    modbusConfig.value = JSON.parse(JSON.stringify(defaultConfig));
    nextId = 3;
  }
};

// modbusConfig 변경 시 ai_modbus에 동기화
const syncToChannel = () => {
  if (currentChannelData.value) {
    currentChannelData.value.ai_modbus = JSON.parse(JSON.stringify(modbusConfig.value));
  }
};

// 컴포넌트 마운트 시 데이터 로드
onMounted(() => {
  loadFromChannel();
});

// 채널 변경 시 데이터 다시 로드
watch(() => props.channel, () => {
  loadFromChannel();
});

// modbusConfig 변경 감지하여 동기화
watch(modbusConfig, () => {
  syncToChannel();
}, { deep: true });

// 행 추가
const addRow = () => {
  modbusConfig.value.push(createDefaultRow(nextId++));
};

// 행 삭제
const removeRow = (index) => {
  if (modbusConfig.value.length > 1) {
    modbusConfig.value.splice(index, 1);
  }
};

// 초기값으로 리셋
const handleReset = () => {
  if (confirm("Are you sure you want to reset to default configuration?")) {
    nextId = 3;
    modbusConfig.value = JSON.parse(JSON.stringify(defaultConfig));
  }
};

// useAI 상태 확인
const isUseAIEnabled = computed(() => {
  const channelUseAI =
    currentChannelData.value?.useAI === 1 || currentChannelData.value?.useAI === true;
  return channelUseAI;
});
</script>

<style scoped>
input[type="number"]:focus,
select:focus {
  outline: none;
}

input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type="number"] {
  -moz-appearance: textfield;
}
</style>