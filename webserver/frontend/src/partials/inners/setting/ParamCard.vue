<template>
  <div
    class="relative col-span-full xl:col-span-8 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
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
          {{ t("config.parameters") }}
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700/60"
      >
        <div
          class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-2"
        >
          {{ t("config.parametersConfiguration") }}
        </div>
        <div class="mt-6 pt-2">
          <!-- 컨테이너를 flex로 변경하고 스크롤바 공간 확보 -->
          <div class="flex flex-col">
            <!-- 헤더 - 스크롤바 영역만큼 패딩 추가 -->
            <div
              class="grid grid-cols-[25%_6%_8%_1fr_8%_10%] gap-2 text-xs font-semibold text-gray-500 dark:text-gray-400 px-1 mb-4"
              style="padding-right: 16px"
            >
              <div class="text-left">Parameter</div>
              <div class="text-left">Unit</div>
              <div class="text-left">Module</div>
              <div class="text-left">Default Thresholds(Min/Max)</div>
              <div class="text-left">Ack.</div>
              <div class="text-left">Test Type</div>
            </div>

            <!-- 스크롤 가능한 데이터 목록 -->
            <div class="overflow-y-auto max-h-[300px] space-y-2">
              <div
                v-for="({ row, originalIndex }, idx) in filteredParamData"
                :key="originalIndex"
                class="grid grid-cols-[25%_6%_8%_1fr_8%_10%] gap-2 items-center border-b border-gray-200 dark:border-gray-700/60 py-2 text-sm px-1"
              >
                <div class="text-left" hidden>{{ row.originalIndex }}</div>
                <!-- Name -->
                <div class="text-left text-xs text-gray-800 dark:text-gray-200">
                  {{ row.Titles[locale] }}
                </div>
                <!-- Unit -->
                <div class="text-left text-xs text-gray-600 dark:text-gray-400">
                  {{ row.Unit || '-' }}
                </div>
                <!-- Module -->
                <div class="text-left text-xs text-gray-800 dark:text-gray-200">
                  {{ row.AssemblyID }}
                </div>
                <!-- Value 입력 -->
                <div v-if="row.DefaultThresholds">
                  <template
                    v-for="(item, index) in row.DefaultThresholds"
                    :key="index"
                  >
                    <input
                      v-if="!isNaNValue(item)"
                      v-model.number="row.DefaultThresholds[index]"
                      type="text"
                      class="text-center border rounded-md p-1 w-12 text-xs mr-1 transition-all text-gray-800 dark:text-gray-800"
                      :class="[
                        !isEditParameters
                          ? 'disabled-input cursor-not-allowed'
                          : 'enabled-input focus:ring-violet-500 focus:border-violet-500',
                        getBackgroundColorClass(index),
                      ]"
                      :disabled="!isEditParameters"
                      :maxlength="20"
                    />
                  </template>
                </div>
                <div
                  v-else
                  class="text-gray-400 dark:text-gray-500 text-xs italic"
                >
                  No thresholds
                </div>
                <div>
                  <input
                    type="checkbox"
                    v-model="row.Acknowledged"
                    :disabled="!isEditParameters"
                    class="transition-all"
                    :class="
                      !isEditParameters
                        ? 'opacity-50 cursor-not-allowed'
                        : 'cursor-pointer'
                    "
                    
                  />
                </div>
                <div class="text-left text-xs text-gray-800 dark:text-gray-200">
                  {{ checkTestType(row.TestType) }}
                </div>
                <!-- Unit -->
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, defineProps, computed, onMounted, watch } from "vue";
import axios from "axios";
import { useAuthStore } from "@/store/auth"; // ✅ Pinia store 사용
import { useI18n } from "vue-i18n";
defineProps({
  channel: String,
});
// const authStore = useAuthStore(); // ✅ Pinia store 사용
const { locale } = useI18n();
const { t } = useI18n();
const inputDict = inject("channel_inputDict");
const paramData = inject("paramData");
const isEditParameters = inject("isEditParameters");
const filteredParamData = computed(() =>
  paramData.value
    .map((row, originalIndex) => ({ row, originalIndex }))
    .filter(
      (item) =>
        item.row &&
        Array.isArray(item.row.DefaultThresholds) &&
        item.row.DefaultThresholds.length > 0
    )
);
const checkTestType = (val) => {
  let msg = "";
  switch (val) {
    case 0:
      msg = "None";
      break;
    case 1:
      msg = "Max";
      break;
    case 2:
      msg = "Min";
      break;
    case 3:
      msg = "Both";
      break;
  }
  return msg;
};

const isNaNValue = (value) => {
  // 문자열 "NaN" 또는 숫자 NaN 모두 감지
  return value === "NaN" || Number.isNaN(Number(value));
};

const getBackgroundColorClass = (index) => {
  
  // TestType: 0=None, 1=Max, 2=Min, 3=Both
  let result = '';
    const colorClasses = [
    '',
    'bg-red-200 dark:bg-red-300',      // 0: 빨강
    'bg-orange-200 dark:bg-orange-300', // 1: 오렌지
    'bg-yellow-200 dark:bg-yellow-300', // 2: 노랑
    'bg-yellow-200 dark:bg-yellow-300', // 3: 노랑
    'bg-orange-200 dark:bg-orange-300', // 4: 오렌지
    'bg-red-200 dark:bg-red-300'       // 5: 빨강
  ];
  result = colorClasses[index] || '';
  
  return result;
};

// const setAssetParams = async () => {
//   try {
//     // Proxy 제거: 순수 JS 객체로 변환
//     const plainTableData = paramData.value.map((item) => ({ ...item }));
//     // 또는 const plainTableData = toRaw(tableData.value);

//     const response = await axios.post(
//       `/setting/setAssetParams/${inputDict.value.assetInfo.name}`,
//       plainTableData, // 평범한 배열 넘기기
//       {
//         headers: { "Content-Type": "application/json" },
//       }
//     );

//     if (response.data?.success) {
//       alert("✅ Asset settings saved successfully!");
//     } else {
//       alert(
//         "❌ Failed to save asset settings: " +
//           (response.data.error || "unknown error")
//       );
//     }
//   } catch (error) {
//     console.error("Error occurred while saving asset:", error);
//     alert(
//       "❌ Error occurred while saving asset: " +
//         (error.response?.data?.error || error.message)
//     );
//   }
// };

// const saveChanges = async () => {
//   try {
//     await setAssetParams();
//   } catch (error) {
//     console.error("Error occurred while saving:", error);
//   }
// };
</script>

<style scoped>
/* 스크롤바 스타일링 - 필요시 사용 */
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

/* ✅ Disabled Input Styles - 조건부 비활성화 */
.disabled-input {
  @apply border-gray-300 dark:border-gray-600 
         bg-gray-100 dark:bg-gray-800 
         text-gray-500 dark:text-gray-400
         opacity-75;
}

/* ✅ Enabled Input Styles - 정상 상태 */
.enabled-input {
  @apply border-gray-300 dark:border-gray-500 
         bg-white dark:bg-gray-900 
         text-gray-800 dark:text-gray-100;
}

/* 배경색 input 스타일 - 텍스트는 어두운 색상으로 고정 */
input.text-gray-800 {
  color: #1f2937 !important; /* 라이트 모드와 다크 모드 모두 어두운 텍스트 */
}

/* 배경색 강제 적용 */
.bg-red-200 {
  background-color: #fecaca !important;
}

.bg-orange-200 {
  background-color: #fed7aa !important;
}

.bg-yellow-200 {
  background-color: #fef3c7 !important;
}

.dark .bg-red-300 {
  background-color: #fca5a5 !important;
}

.dark .bg-orange-300 {
  background-color: #fdba74 !important;
}

.dark .bg-yellow-300 {
  background-color: #fde047 !important;
}
/* 포커스 상태에서도 텍스트 색상 유지 */
input:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}

/* Disabled 상태에서는 포커스 효과 제거 */
input:disabled:focus {
  box-shadow: none;
  outline: none;
}

button:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}

/* Transition 효과 추가 */
input {
  transition: all 0.2s ease-in-out;
}
</style>