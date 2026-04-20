<template>
  <div
    class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
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
          {{ t("config.trendSetup") }}
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <div class="space-y-6">
        <div>
          <div class="space-y-2">
            <div>
              <label class="block text-sm font-medium mb-2">Collect period (min)</label>

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
                :value="selectedTrendSetup?.params?.[index] ?? 0"
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
import { useI18n } from "vue-i18n";

// inject된 값들
const selectedTrendSetup = inject("selectedTrendSetup");
const updateNestedField = inject("updateNestedField");
const { t } = useI18n();
const parameterOptions = ref([
  { label: "None", value: 0 },
  { label: "Frequency", value: 1 },
  { label: "Line Voltage", value: 2 },
  { label: "Phase Voltage", value: 3 },
  { label: "Current", value: 4 },
  { label: "Power", value: 5 },
  { label: "PF", value: 6 },
  { label: "Unbalance", value: 7 },
  { label: "THD", value: 8 },
  { label: "TDD", value: 9 },
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
  const newValue = Number(event.target.value);

  if (selectedTrendSetup?.value?.params) {
    const newParams = [...selectedTrendSetup.value.params];
    newParams[index] = newValue;
    selectedTrendSetup.value = {
      ...selectedTrendSetup.value,
      params: newParams,
    };

    if (updateNestedField) {
      updateNestedField('trendInfo', 'params', newParams);
    }
  }
}

</script>