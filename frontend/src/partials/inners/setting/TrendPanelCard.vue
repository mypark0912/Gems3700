<template>
  <!-- device0 모드일 때는 컴포넌트를 숨김 -->
  <div
    v-if="opMode !== 'device0' "
    class="relative col-span-full xl:col-span-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
    <div class="absolute top-0 left-0 right-0 h-0.5 bg-pink-500" aria-hidden="true"></div>
    <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
      <header class="flex items-center mb-2">
        <div class="w-6 h-6 rounded-full shrink-0 bg-pink-500 mr-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 3v18h18" />
            <polyline points="5 15 9 10 13 13 17 7 21 11" />
          </svg>
        </div>
        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
          Trend Setup
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <div class="space-y-6">
        <div>
          <div class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-4">
            Save Option
          </div>

          <div class="space-y-2">
            <div>
              <label class="block text-sm font-medium mb-2">Period (min)</label>

              <select 
                :value="selectedTrendSetup?.period || 5" 
                @change="onPeriodChange"
                class="form-select w-full">
                <option value="1">1</option>
                <option value="5">5</option>
                <option value="10">10</option>
                <option value="15">15</option>
                <option value="30">30</option>
              </select>
            </div>

            <div class="flex justify-between items-center text-xs font-semibold text-gray-500 px-1">
              <div class="w-6 block text-sm font-medium mb-2 mt-2">Index</div>
              <div class="w-1/2 max-w-[240px] text-left block text-sm font-medium mb-2 mt-2">
                Parameter
              </div>
            </div>

            <div v-for="(param, index) in (selectedTrendSetup?.params || [])" :key="`param-${index}`"
              class="flex justify-between items-center border-b py-2 text-sm">
              <div class="w-16 text-left ml-5">{{ index + 1 }}</div>
              
              <select 
                :value="selectedTrendSetup?.params?.[index] || 'None'" 
                @change="onParameterChange(index, $event)"
                class="form-select flex-1 text-sm">
                <option v-for="(option, i) in parameterOptions" :key="i" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, ref, computed, onMounted, nextTick } from "vue";
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()
const opMode = computed(() => authStore.getOpMode)

// inject된 값들
const selectedTrendSetup = inject("selectedTrendSetup");
const updateNestedField = inject("updateNestedField");

const parameterOptions = ref([
  { label: "None", value: "None" },
  { label: "Temperature", value: "Temperature" },
  { label: "Frequency", value: "Frequency" },
  { label: "Line Voltage", value: "Line Voltage" },
  { label: "Phase Voltage", value: "Phase Voltage" },
  { label: "Current", value: "Current" },
  { label: "Power", value: "Power" },
  { label: "PF", value: "PF" },
  { label: "Unbalance", value: "Unbalance" },
  { label: "THD", value: "THD" },
  { label: "TDD", value: "TDD" },
]);

onMounted(async () => {
  await nextTick();
});

function onPeriodChange(event) {
  const newPeriod = Number(event.target.value);
  
  if (selectedTrendSetup?.value) {
    selectedTrendSetup.value.period = newPeriod;
  }
  
  if (updateNestedField) {
    updateNestedField('trendInfo', 'period', newPeriod);
  }
}

function onParameterChange(index, event) {
  const newValue = event.target.value;
  
  if (selectedTrendSetup?.value?.params) {
    selectedTrendSetup.value.params[index] = newValue;
    
    if (updateNestedField) {
      updateNestedField('trendInfo', 'params', [...selectedTrendSetup.value.params]);
    }
  }
}

</script>