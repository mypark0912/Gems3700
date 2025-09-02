<template>
<div class="relative col-span-full xl:col-span-8 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
      <div class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500" aria-hidden="true"></div>
      <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
            <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <circle cx="17" cy="17" r="4" stroke="currentColor" stroke-width="2" fill="none"/>
              <path d="M20 20l2 2" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">Parameters/Thresholds</h3>
        </header>
      </div>
      <div class="px-4 py-3 space-y-4">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700/60">
            <div class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-2">Parameter Configuration</div>
          <div class="mt-6 pt-2">
                    <div class="grid grid-cols-[30%_10%_1fr_10%_10%] gap-2 text-xs font-semibold text-gray-500 dark:text-gray-400 px-1 mb-4">
                    <div class="text-left">Parameter</div>
                    <div class="text-left">Module</div>
                    <div class="text-left">Default Thresholds(Min/Max)</div>
                    <div class="text-left">Ack.</div>
                    <div class="text-left">Test Type</div>
                    </div>

                    <!-- 스크롤 가능한 데이터 목록 -->
                    <div class="overflow-y-auto max-h-[300px] space-y-2 pr-1">
                        <div
                           v-for="({ row, originalIndex }, idx) in filteredParamData" :key="originalIndex"
                            class="grid grid-cols-[30%_10%_1fr_10%_10%] gap-2 items-center border-b border-gray-200 dark:border-gray-700/60 py-2 text-sm"
                        >
                            <div class="text-left" hidden>{{ row.originalIndex }}</div>
                            <!-- Name -->
                            <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.Title }}</div>
                            <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.AssemblyID }}</div>
                            <!-- Value 입력 -->
                             <div v-if="row.DefaultThresholds">
                                <template v-for="(item, index) in row.DefaultThresholds" :key="index">
                                    <input
                                        v-if="!isNaNValue(item)"
                                        v-model.number="row.DefaultThresholds[index]"
                                        type="text"
                                        class="text-center border border-gray-300 dark:border-gray-600 rounded-md p-1 w-12 text-xs mr-1 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 focus:ring-violet-500 focus:border-violet-500"
                                        :disabled="!isEditParameters"
                                    />
                                </template>
                            </div>
                            <div v-else class="text-gray-400 dark:text-gray-500 text-xs italic">No thresholds</div>
                            <div>
                              <input type="checkbox" v-model="row.Acknowledged" :disabled="!isEditParameters" />
                            </div>
                            <div class="text-left text-gray-800 dark:text-gray-200">{{ checkTestType(row.TestType) }}</div>
                            <!-- Unit -->
                        </div>
                    </div>


                </div>
        </div>
        <!--div class="mt-4 flex justify-end space-x-4">
          <button class="btn h-6 px-5 bg-blue-900 text-blue-100 hover:bg-blue-800 dark:bg-blue-100 dark:text-blue-800 dark:hover:bg-white" @click="saveChanges">Save</button>
        </div-->
      </div>
    </div>
</template>

<script setup>
import { inject, defineProps, computed, onMounted, watch } from 'vue'
import axios from 'axios';
import { useAuthStore } from "@/store/auth"; // ✅ Pinia store 사용

defineProps({
  channel: String
})
// const authStore = useAuthStore(); // ✅ Pinia store 사용
const inputDict = inject('channel_inputDict')
const paramData = inject('paramData');
const isEditParameters = inject('isEditParameters');
const filteredParamData = computed(() =>
  paramData.value
    .map((row, originalIndex) => ({ row, originalIndex }))
    .filter(item => item.row && Array.isArray(item.row.DefaultThresholds) && item.row.DefaultThresholds.length > 0)
);
    const checkTestType = (val) => {
        let msg = "";
        switch(val){
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
    return value === 'NaN' || Number.isNaN(Number(value));
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
/* 다크모드를 위한 포커스 및 호버 스타일 향상 */
input:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}

button:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}
</style>