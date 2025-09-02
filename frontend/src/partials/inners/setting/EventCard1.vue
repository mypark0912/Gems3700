<template>
  <div
    class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div class="absolute top-0 left-0 right-0 h-0.5 bg-violet-500"></div>
    <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
      <header class="flex items-center mb-2">
        <div class="w-6 h-6 rounded-full shrink-0 bg-violet-500 mr-3">
          <svg
            class="w-6 h-6 fill-current text-white"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M12 2C9.5 2 7.5 4 7.5 6.5V10c0 1.2-.4 2.4-1.2 3.3L5 14.5c-1 1 0 2.5 1.3 2.5h11.4c1.3 0 2.3-1.5 1.3-2.5l-1.3-1.2c-.8-.9-1.2-2.1-1.2-3.3V6.5C16.5 4 14.5 2 12 2z"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
            />
            <path
              d="M10 19a2 2 0 0 0 4 0"
              stroke="currentColor"
              stroke-width="2"
              fill="none"
            />
          </svg>
        </div>
        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
          Event Setup 1
        </h3>
      </header>
    </div>

    <div class="px-4 py-3 space-y-4">
      <div class="space-y-6">
        <!-- Transient Voltage -->
        <div>
          <div class="text-xs font-semibold uppercase mb-4">Transient Voltage</div>
            <div>
              <label class="block text-sm font-medium mb-2">Action</label>
              <div class="flex w-full">
                <button
                  v-for="option in tv_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="settv_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.tv_action === option.value
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
              <input v-model="inputDict.eventInfo.tv_level" class="form-input w-full" type="text" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Holdoff time (ms)</label>
              <input v-model="inputDict.eventInfo.tv_holdofftime" class="form-input w-full" type="text" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Fast Change(%)</label>
              <input v-model="inputDict.eventInfo.tv_fastchange" class="form-input w-full" type="text" />
            </div>
          </div>
        </div>

        <!-- Transient Current -->
        <div>
          <div class="text-xs font-semibold uppercase mb-4">Transient Current</div>
            <div>
              <label class="block text-sm font-medium mb-2">Action</label>
              <div class="flex w-full">
                <button
                  v-for="option in tc_action"
                  :key="option.value"
                  :value="option.value"
                  @click.prevent="settc_action(option.value)"
                  :class="[
                    'btn border px-4 py-2 transition-colors duration-200 rounded-none first:rounded-l-lg last:rounded-r-lg w-1/2',
                    inputDict.eventInfo.tc_action === option.value
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
              <input v-model="inputDict.eventInfo.tc_level" class="form-input w-full" type="text" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Holdoff time (ms)</label>
              <input v-model="inputDict.eventInfo.tc_holdofftime" class="form-input w-full" type="text" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1.5">Fast Change(%)</label>
              <input v-model="inputDict.eventInfo.tc_fastchange" class="form-input w-full" type="text" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, computed } from "vue";

const inputDict = inject("channel_inputDict");

const tc_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];
const tv_action = [
  { value: 0, label: "None" },
  { value: 1, label: "Event" },
  { value: 2, label: "Capture" },
];

const settc_action = (value) => {
  inputDict.value.eventInfo.tc_action = parseInt(value);
};
const settv_action = (value) => {
  inputDict.value.eventInfo.tv_action = parseInt(value);
};

// ✅ 숫자 0/1로 저장하며 Boolean처럼 동작하는 computed
const tv_event_bool = computed({
  get: () => parseInt(inputDict.value.eventInfo.tv_event) === 1,
  set: (val) => (inputDict.value.eventInfo.tv_event = val ? 1 : 0),
});

const tc_event_bool = computed({
  get: () => parseInt(inputDict.value.eventInfo.tc_event) === 1,
  set: (val) => (inputDict.value.eventInfo.tc_event = val ? 1 : 0),
});
</script>
