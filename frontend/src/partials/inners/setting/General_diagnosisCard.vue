<template>
  <div
    v-if="showDiagnosis"
    class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500" aria-hidden="true"></div>
    <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
      <header class="flex items-center justify-between mb-2">
        <div class="flex items-center space-x-3">
          <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
          <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 15.5a3.5 3.5 0 1 1 3.5-3.5 3.5 3.5 0 0 1-3.5 3.5zM19.5 12a7.5 7.5 0 0 1-1 3.7l2.1 2.1-1.8 1.8-2.1-2.1a7.5 7.5 0 0 1-3.7 1V21h-2.5v-2.5a7.5 7.5 0 0 1-3.7-1l-2.1 2.1-1.8-1.8 2.1-2.1a7.5 7.5 0 0 1-1-3.7H3v-2.5h2.5a7.5 7.5 0 0 1 1-3.7L4.4 4.4l1.8-1.8 2.1 2.1a7.5 7.5 0 0 1 3.7-1V3h2.5v2.5a7.5 7.5 0 0 1 3.7 1l2.1-2.1 1.8 1.8-2.1 2.1a7.5 7.5 0 0 1 1 3.7H21V12h-1.5z" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" />
            <circle cx="17" cy="17" r="4" stroke="currentColor" stroke-width="2" fill="none" />
            <path d="M20 20l2 2" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </div>
        <h3 class="text-left text-lg text-gray-800 dark:text-gray-100 font-semibold">
          {{t("config.plansPanel.dig.title") }}<!--Diagnostic Setup-->
        </h3>
        </div>
          <button
            v-if="isNtek"
            @click="showProfile"
            class="btn h-6 px-5 bg-violet-900 text-violet-100 hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
          >
            {{t("config.plansPanel.dig.advanced") }}<!--Advanced-->
          </button>
      </header>
    </div>

    <div class="px-5 py-6">
      <div class="grid grid-cols-4 gap-6">
        <div
          v-for="(value, key) in diagnosisData"
          :key="key"
          :class="[
            'col-span-1',
            (key === 'DBToken' || key === 'DataPath') ? 'col-span-2' : '',
            (typeof value === 'string' && value.length > 100) ? 'col-span-4' : ''
          ]"
        >
          <label class="block text-sm font-medium mb-2">{{ key }}</label>

          <!-- Boolean → Select Box -->
          <select
            v-if="typeof value === 'boolean'"
            v-model="diagnosisData[key]"
            class="form-select w-full"
          >
            <option :value="true">True</option>
            <option :value="false">False</option>
          </select>

          <!-- Array → Comma-separated Text -->
          <input
            v-else-if="Array.isArray(value)"
            :value="diagnosisData[key].join(', ')"
            @input="onArrayInput($event, key)"
            class="form-input w-full"
            type="text"
          />

          <!-- Number → Input type number -->
          <input
            v-else-if="typeof value === 'number'"
            v-model.number="diagnosisData[key]"
            class="form-input w-full"
            type="number"
          />

          <!-- Default → Text -->
          <input
            v-else
            v-model="diagnosisData[key]"
            class="form-input w-full"
            type="text"
          />
        </div>
      </div>
    </div>
  </div>
    <ModalBasic
    id="advanced-modal"
    :modalOpen="advancedModalOpen"
    @close-modal="advancedModalOpen = false"
    title="Adanced Setup"
  >
  <!-- Modal content -->
    <div class="px-5 py-6">
      <div class="grid grid-cols-2 gap-6">
          <div
            v-for="(value, key) in advancedData"
            :key="key"
            :class="[
              'col-span-1',
               ((typeof value === 'string' && value.length > 100) || Array.isArray(value)) ? 'col-span-2' : ''
            ]"
          >
            <label class="block text-sm font-medium mb-2">{{ key }}</label>

            <!-- Boolean → Select Box -->
            <select
              v-if="typeof value === 'boolean'"
              v-model="advancedData[key]"
              class="form-select w-full"
            >
              <option :value="true">True</option>
              <option :value="false">False</option>
            </select>

            <!-- Array → Comma-separated Text -->
            <!--input
              v-else-if="Array.isArray(value)"
              :value="advancedData[key].join(', ')"
              @input="onArrayInput($event, key)"
              class="form-input w-full"
              type="text"
            /-->
            <input
              v-else-if="Array.isArray(value)"
              v-model="arrayBufferMap[key]"
              class="form-input w-full"
              type="text"
            />

            <!-- Number → Input type number -->
            <input
              v-else-if="typeof value === 'number'"
              v-model.number="advancedData[key]"
              class="form-input w-full"
              type="number"
            />

            <!-- Default → Text -->
            <input
              v-else
              v-model="advancedData[key]"
              class="form-input w-full"
              type="text"
            />
          </div>
        </div>
    </div>
  <!-- Modal footer -->
  <div
    class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60"
  >
    <div class="flex flex-wrap justify-end space-x-2">
      <button
        class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
        @click.stop="onSaveAdvanced"
      >
        Save
      </button>
      <button
        class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
        @click.stop="advancedModalOpen = false"
      >
        Cancel
      </button>
    </div>
  </div>
</ModalBasic>
</template>

<script setup>
import { inject, ref, computed, watchEffect } from "vue";
import ModalBasic from "../../../pages/common/ModalBasic.vue";
import { useAuthStore } from "@/store/auth"; // ✅ Pinia store 사용
import axios from "axios";
import { useI18n } from "vue-i18n"; // 


// const Diagnosis_inputDict = inject("Diagnosis_inputDict");
const { t } = useI18n();
const diagnosisData = inject('diagnosisData');
const advancedModalOpen = ref(false);
const authStore = useAuthStore();
const isNtek = computed(()=>{
   return (authStore.getUserRole == '3' && authStore.getUser == 'ntek')
})
const advancedData  = inject('advancedData');
const arrayBufferMap = ref({});
const props = defineProps({
  showDiagnosis: {
    type: Boolean,
    required: true,
  }
});

function onArrayInput(event, key) {
  const val = event.target.value;
  advancedData[key] = val
    .split(',')
    .map(v => parseFloat(v.trim()))
    .filter(v => !isNaN(v));
}

watchEffect(() => {
  for (const key in advancedData.value) {
    if (Array.isArray(advancedData.value[key])) {
      //console.log("✔ 배열 초기화:", key, advancedData.value[key]);
      arrayBufferMap.value[key] = advancedData.value[key].join(', ');
    }
  }
});


    const GetAdvanced = async () => {
      try {
        const response = await axios.get("/setting/getDiagnosisProfile");

        if (response.data.success) {
          Object.assign(advancedData.value, response.data.data); // ✅ 기존 데이터 유지하면서 새로운 데이터 병합
          //alert("Data Loaded!");
        }
      } catch (error) {
        console.error("데이터 가져오기 실패:", error);
        errorMessage.value = "데이터를 불러오는 중 오류가 발생했습니다.";
      }
    };

  const showProfile = async() =>{
    await GetAdvanced();
    advancedModalOpen.value = true
  }

  function onSaveAdvanced() {
  for (const key in arrayBufferMap.value) {
    const raw = arrayBufferMap.value[key];
    const parsed = raw
      .split(',')
      .map(v => parseFloat(v.trim()))
      .filter(v => !isNaN(v)); // 문자열 배열 원하면 이 줄 제거

    advancedData.value[key] = parsed;
  }
  advancedModalOpen.value = false;
}
</script>
