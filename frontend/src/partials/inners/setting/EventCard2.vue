<template>
  <div
    class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-purple-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center mb-2">
        <div class="w-6 h-6 rounded-full shrink-0 bg-purple-500 mr-3">
          <svg
            class="w-6 h-6 fill-current text-white"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <!-- 🔔 벨 아이콘 -->
            <path
              d="M12 2C9.5 2 7.5 4 7.5 6.5V10c0 1.2-.4 2.4-1.2 3.3L5 14.5c-1 1 0 2.5 1.3 2.5h11.4c1.3 0 2.3-1.5 1.3-2.5l-1.3-1.2c-.8-.9-1.2-2.1-1.2-3.3V6.5C16.5 4 14.5 2 12 2z"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
            />

            <!-- 🔔 벨의 하단 둥근 부분 -->
            <path
              d="M10 19a2 2 0 0 0 4 0"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
            />
          </svg>
        </div>
        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
          {{ title }}
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <!-- ✅ 전체 여백 조정 -->
      <div class="space-y-6">
        <!-- Sag -->
        <div v-if="option1 == 'Sag'">
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-4"
          >
            {{ option1 }}
          </div>
            <div>
              <label class="block text-sm font-medium mb-2" for="fw-version"
                >Action</label
              >
              <div class="flex w-full">
                <button
                  v-for="option in sag_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="setsag_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.sag_action === option.value
                      ? 'bg-violet-500 text-white border-violet-500' // ✅ 선택된 버튼: 보라색 배경 + 흰색 텍스트
                      : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900', // ✅ 선택되지 않은 버튼: 흰 배경 + 보라 텍스트
                  ]"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          <div
            class="sm:flex sm:items-center space-y-3 sm:space-y-0 sm:space-x-3 mt-4"
          >
            <!-- ✅ 그룹 간격 조정 -->
            <div class="sm:w-1/2">
              <label class="block text-sm font-medium mb-2" for="fw-version"
                >Level(%)</label
              >
              <input
                v-model="inputDict.eventInfo.sag_level"
                id="fw-version"
                class="form-input w-full"
                type="text"
              />
            </div>

            <div class="sm:w-1/2">
              <label class="block text-sm font-medium mb-2" for="fw-version"
                >HoldOff time(ms)</label
              >
              <input
                v-model="inputDict.eventInfo.sag_holdofftime"
                id="fw-version"
                class="form-input w-full"
                type="text"
              />
            </div>
          </div>
        </div>
        <div v-else>
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-4"
          >
            {{ option1 }}
          </div>
                      <div>
              <label class="block text-sm font-medium mb-2">Action</label>
              <div class="flex w-full">
                <button
                  v-for="option in oc_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="setoc_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.oc_action === option.value
                      ? 'bg-violet-500 text-white border-violet-500' // ✅ 선택된 버튼: 보라색 배경 + 흰색 텍스트
                      : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900', // ✅ 선택되지 않은 버튼: 흰 배경 + 보라 텍스트
                  ]"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          <div class="grid grid-cols-2 gap-4 mt-4">
            <div>
              <label class="block text-sm font-medium mb-2">Level (%)</label>
              <input
                v-model="inputDict.eventInfo.oc_level"
                class="form-input w-full"
                type="text"
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2"
                >Holdoff time (ms)</label
              >
              <input
                v-model="inputDict.eventInfo.oc_holdofftime"
                class="form-input w-full"
                type="text"
              />
            </div>
          </div>
        </div>

        <!-- Swell -->
        <div v-if="option2 == 'Swell'">
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-4 mt-4"
          >
            {{ option2 }}
          </div>
            <div>
              <label class="block text-sm font-medium mb-2">Action</label>
              <div class="flex w-full">
                <button
                  v-for="option in swell_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="setswell_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.swell_action === option.value
                      ? 'bg-violet-500 text-white border-violet-500'
                      : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                  ]"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          <div class="grid grid-cols-2 gap-4 mt-4">
            <div>
              <label class="block text-sm font-medium mb-2">Level (%)</label>
              <input
                v-model="inputDict.eventInfo.swell_level"
                class="form-input w-full"
                type="text"
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2"
                >Holdoff time (ms)</label
              >
              <input
                v-model="inputDict.eventInfo.swell_holdofftime"
                class="form-input w-full"
                type="text"
              />
            </div>
          </div>
        </div>
        <div v-else>
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-4"
          >
            {{ option2 }}
          </div>
            <div>
              <label class="block text-sm font-medium mb-2">Action</label>
              <div class="flex w-full">
                <button
                  v-for="option in inter_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="setinter_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.inter_action === option.value
                      ? 'bg-violet-500 text-white border-violet-500' // ✅ 선택된 버튼: 보라색 배경 + 흰색 텍스트
                      : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900', // ✅ 선택되지 않은 버튼: 흰 배경 + 보라 텍스트
                  ]"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          <div class="grid grid-cols-2 gap-4 mt-4">
            <div>
              <label class="block text-sm font-medium mb-2">Level (%)</label>
              <input
                v-model="inputDict.eventInfo.inter_level"
                class="form-input w-full"
                type="text"
              />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2"
                >Holdoff time (ms)</label
              >
              <input
                v-model="inputDict.eventInfo.inter_holdofftime"
                class="form-input w-full"
                type="text"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, defineProps, ref, computed } from "vue";

const props = defineProps({
  title: String,
  option1: String,
  option2: String,
});
const title = ref(props.title);
const option1 = ref(props.option1);
const option2 = ref(props.option2);

const inputDict = inject("channel_inputDict");

// ✅ Action 버튼 정의
const sag_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];
const swell_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];
const inter_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];
const oc_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];

// ✅ Action setter
const setsag_action = (value) => {
  inputDict.value.eventInfo.sag_action = parseInt(value);
};
const setswell_action = (value) => {
  inputDict.value.eventInfo.swell_action = parseInt(value);
};
const setinter_action = (value) => {
  inputDict.value.eventInfo.inter_action = parseInt(value);
};
const setoc_action = (value) => {
  inputDict.value.eventInfo.oc_action = parseInt(value);
};

// ✅ Event 체크박스용 computed: 숫자 저장 + 불린 UI
const sag_event_bool = computed({
  get: () => Number(inputDict.value?.eventInfo?.sag_event) === 1,
  set: (val) => (inputDict.value.eventInfo.sag_event = val ? 1 : 0),
});
const swell_event_bool = computed({
  get: () => Number(inputDict.value?.eventInfo?.swell_event) === 1,
  set: (val) => (inputDict.value.eventInfo.swell_event = val ? 1 : 0),
});
const inter_event_bool = computed({
  get: () => Number(inputDict.value?.eventInfo?.inter_event) === 1,
  set: (val) => (inputDict.value.eventInfo.inter_event = val ? 1 : 0),
});
const oc_event_bool = computed({
  get: () => Number(inputDict.value?.eventInfo?.oc_event) === 1,
  set: (val) => (inputDict.value.eventInfo.oc_event = val ? 1 : 0),
});
</script>
