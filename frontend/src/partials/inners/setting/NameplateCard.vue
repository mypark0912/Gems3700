<template>
  <div
    class="relative col-span-full xl:col-span-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
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
          Nameplates 
        </h3>
      </header>
    </div>
    <div class="px-4 py-3 space-y-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700/60">
        <div class="flex justify-between items-center mb-2">
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase"
          >
            Nameplates Configuration
          </div>
          <button v-if="isAdmin"
            class="btn h-6 px-5 ml-auto mr-5 bg-sky-900 text-sky-100 hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
            @click="feedbackModalOpen = true"
          >
            Add Bearing
          </button>
          <button
            v-if="authStore.getUserRole !== '2'"
            @click="showAdvancedModal = true"
            class="btn h-6 px-5 bg-violet-900 text-violet-100 hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
          >
            Advanced
          </button>
        </div>
        <div class="mt-6 pt-2">
          <!-- 헤더 -->
          <div
            class="grid grid-cols-[30%_10%_1fr_10%] gap-2 text-xs font-semibold text-gray-500 dark:text-gray-400 px-1 mb-4"
          >
            <div class="text-left">Item</div>
            <div class="text-left">Module</div>
            <div class="text-center">Value</div>
            <div class="text-center">Unit</div>
          </div>

          <!-- 스크롤 가능한 데이터 목록 -->
          <div class="overflow-y-auto max-h-[300px] space-y-2 pr-1">
            <div
              v-for="(row, index) in readonlyRows"
              :key="index"
              class="grid grid-cols-[30%_10%_1fr_10%] gap-2 items-center border-b border-gray-200 dark:border-gray-700/60 py-2 text-sm"
            >
              <!-- Name -->
              <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.Title }}</div>
              <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.AssemblyID }}</div>
              <!-- Value 입력 -->
              <input
                v-model.number="row.Value"
                type="number"
                class="text-center border border-gray-300 dark:border-gray-600 rounded-md p-1 w-full text-xs bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200"
                :min="getMinValue(row)"
                :max="getMaxValue(row)"
                disabled
              />
              <div class="text-left" hidden>{{ row.Min }}</div>
              <div class="text-left" hidden>{{ row.Max }}</div>
              <!-- Unit -->
              <div class="text-center text-xs text-gray-600 dark:text-gray-300">{{ row.Unit }}</div>
            </div>

            <div
              v-for="(row, index) in filteredEditableRows"
              :key="index"
              class="grid grid-cols-[30%_10%_1fr_10%] gap-2 items-center border-b border-gray-200 dark:border-gray-700/60 py-2 text-sm"
            >
              <!--div class="text-left text-gray-800 dark:text-gray-200">
                {{ `${row.Assembly}_${row.AssemblyID}: ${row.Title}` }}
              </div-->
              <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.Title }}</div>
              <div class="text-left text-xs text-gray-800 dark:text-gray-200">{{ row.AssemblyID }}</div>
              <div class="text-left" hidden>{{ `${row.Name}` }}</div>
              <template v-if="row.DataType === 3">
                <select
                  v-model.number="row.Value"
                  class="text-xs w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded p-1 text-center focus:ring-violet-500 focus:border-violet-500" 
                  :disabled="!isEditNameplates"
                >
                  <option
                    v-for="(label, i) in row.DataInfo"
                    :key="i"
                    :value="i"
                  >
                    {{ label }}
                  </option>
                </select>
              </template>

              <template v-else-if="row.DataType === 4">
                <div v-if="row.Name == 'BearingName'" class="flex items-center gap-2">
                  <select
                  v-model="row.Value"
                  class="text-xs w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded p-1 text-center focus:ring-violet-500 focus:border-violet-500" 
                  :disabled="!isEditNameplates"
                  @change="onBearingChanged(row)"
                >
                  <option v-for="data in BearingOptions" :value="data">
                    {{ data }}
                  </option>
                </select>
                <button
                class="text-xs px-2 py-1 border border-gray-300 dark:border-gray-600 rounded bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200" 
                :disabled="!isEditNameplates"
                @click="openSearchModal(row)"
              >
                🔍
              </button>
                </div>                
                <input v-else 
                  type="text"
                  v-model="row.Value"
                  class="text-xs w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded p-1 text-center focus:ring-violet-500 focus:border-violet-500" 
                  :disabled="!isEditNameplates"
                />
              </template>

              <template v-else>
                <input
                  type="number"
                  v-model.number="row.Value"
                  :min="row.Min"
                  :max="row.Max"
                  class="text-xs w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded p-1 text-center focus:ring-violet-500 focus:border-violet-500" 
                  :disabled="!isEditNameplates"
                />
              </template>
              <div class="text-left" hidden>{{ row.Min }}</div>
              <div class="text-left" hidden>{{ row.Max }}</div>
              <div class="text-center text-xs text-gray-600 dark:text-gray-300">{{ row.Unit }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ModalBasic
    id="feedback-modal"
    :modalOpen="feedbackModalOpen"
    @close-modal="feedbackModalOpen = false"
    title="File Upload"
  >
    <!-- Modal content -->
    <div class="px-5 py-4">
      <div class="text-sm">
        <div class="font-medium text-gray-800 dark:text-gray-100 mb-3">
          Setting file upload
        </div>
      </div>
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="name"
            >file path <span class="text-red-500">*</span></label
          >
          <input
            id="filename"
            class="form-input w-full px-2 py-1 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md focus:ring-violet-500 focus:border-violet-500"
            @change="handleFileUpload"
            type="file"
            required
          />
        </div>
      </div>
      <div class="text-sm">
        <div class="font-medium text-gray-800 dark:text-gray-100 mb-3">
          {{ message }}🙌
        </div>
      </div>
    </div>
    <!-- Modal footer -->
    <div class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60">
      <div class="flex flex-wrap justify-end space-x-2">
        <button
          class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
          @click.prevent="uploadBearing"
        >
          Import
        </button>
        <button
          class="btn-sm border border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-gray-200"
          @click.stop="feedbackModalOpen = false"
        >
          Cancel
        </button>
      </div>
    </div>
  </ModalBasic>
  <ModalBasic
    :modalOpen="showAdvancedModal"
    @close-modal="showAdvancedModal = false"
    title="Advanced Settings"
  >
    <div class="w-[600px] max-w-full px-6">
      <!-- 헤더 -->
      <div
        class="flex justify-between items-center text-base font-semibold text-gray-500 dark:text-gray-400 mt-2 mb-1"
      >
        <div class="w-1/2 text-left">Title</div>
        <div class="w-1/4 text-center">Value</div>
        <div class="w-1/4 text-right">Unit</div>
      </div>

      <!-- ✅ 리스트 스크롤 영역 -->
      <div class="max-h-[50vh] overflow-y-auto">
        <div
          v-for="(row, index) in modalData"
          :key="index"
          class="flex items-center text-base border-b border-gray-200 dark:border-gray-700/60 py-2 relative"
        >
          <!-- Title -->
          <div class="w-1/2 text-left text-gray-800 dark:text-gray-200">
            {{ row.Title }}
          </div>

          <!-- Value -->
          <div class="w-1/4 text-center relative">
            <input
              v-model.number="row.Value"
              type="number"
              :id="`modal-input-${row.Title}`"
              class="w-full text-center border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md px-2 py-1 text-base focus:ring-violet-500 focus:border-violet-500"
            />
            <div
              v-if="row.invalid"
              class="absolute -top-6 left-1/2 -translate-x-1/2 z-10 text-xs bg-red-500 text-white px-2 py-1 rounded shadow"
            >
              {{ row.tooltip }}
            </div>
          </div>

          <!-- Unit -->
          <div class="w-1/4 text-right pr-2 text-gray-600 dark:text-gray-300">
            {{ row.Unit }}
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60">
        <div class="flex justify-end space-x-2">
          <button
            class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 dark:bg-gray-600 dark:hover:bg-gray-700"
            @click="showAdvancedModal = false"
          >
            Close
          </button>
          <button
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700"
            @click="saveAdvanced"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </ModalBasic>
  <SearchModal id="search-modal" searchId="search" :modalOpen="searchModalOpen" @open-modal="searchModalOpen = true" @close-modal="searchModalOpen = false"  @select="onBearingSelected" />
</template>

<script setup>
import {
  inject,
  defineProps,
  ref,
  watchEffect,
  computed,
  onMounted,
  watch,
  provide,
} from "vue";
import ModalBasic from "../../../pages/common/ModalBasic.vue";
import SearchModal from "../../../pages/common/SearchModal.vue";
import axios from "axios";
import { useAuthStore } from "@/store/auth"; // ✅ Pinia store 사용

defineProps({
  channel: String,
});
const authStore = useAuthStore(); // ✅ Pinia store 사용
const inputDict = inject("channel_inputDict");
const tableData = inject("tableData");
const modalData = inject("modalData");
const isEditNameplates = inject("isEditNameplates");
const showAdvancedModal = ref(false);
const getMinValue = inject("getMinValue");
const getMaxValue = inject("getMaxValue");
const feedbackModalOpen = ref(false);
const searchModalOpen = ref(false);
const readonlyRows = ref([]);
const editableRows = ref([]);
// const searchData = ref(null)
const targetRow = ref(null)
const message = ref("");
const selectedFile = ref("");
const mapping = {
  "Rated Current": "ctInfo.inorminal",
  "Rated Voltage": "ptInfo.vnorminal",
  "Sampling Rate": "sampling.rate",
  Duration: "sampling.duration",
  Interval: "sampling.period",
  Channel : "channel",
  "Connection Type": "ptInfo.wiringmode",
  "Rated Frequency": "ptInfo.linefrequency"
};

// const hz = [8000, 4000, 2000, 1000];
// const duration = [5, 10, 15, 20, 25, 30];
// const intervalList = Array.from({ length: 13 }, (_, i) => i * 5);

const BearingOptions = ref(["My Bearing"]);
const BearingValues = ref([]);
const isAdmin = computed(()=>{
   console.log(authStore.getUserRole);
   return (authStore.getUserRole == 0 && authStore.getUser == 'ntek')
})
watchEffect(() => {
  if (!inputDict.value || !Array.isArray(tableData.value)) return;

  readonlyRows.value = [];
  editableRows.value = [];
  console.log(inputDict.value);
  tableData.value.forEach((row) => {
    //const children = row["children"]
    const path = mapping[row.Title];
    if (!path) {
      // 매핑되지 않은 항목은 사용자 수정 가능
      editableRows.value.push(row);
      return;
    }

    const keys = path.split(".");
    let value = inputDict.value;

    for (const key of keys) {
      if (value && key in value) {
        if (key === "rate") value = parseInt(value[key]);
        else if (key === "duration") value = parseInt(value[key]);
        else if (key === "period") value = parseInt(value[key]);
        else if (key === "channel") value = (value[key] == 'Main'?1:2);
        else if (key === "wiringmode") value = parseInt(value[key]);
        else if (key === "linefrequency") value = parseInt(value[key]);
        else value = value[key];
      } else {
        value = null;
        break;
      }
    }

    if (value !== null && value !== undefined) {
      row.Value = Number(value); // tableData의 실제 값을 갱신
      readonlyRows.value.push(row);
    } else {
      editableRows.value.push(row); // path는 있었지만 값이 없는 경우도 사용자 변경 가능
    }
  });
  // editableRows.value = setStructNameplate(editableRows.value);
});

onMounted(async () => {
  try {
    const response = await axios.get("/setting/checkBearing");
    if (response.data.passOK == "1") {
      for (let i = 0; i < response.data.data.length; i++) {
        BearingValues.value.push(response.data.data[i]);
      }
    }
  } catch (error) {
    //message.value = "업로드 실패: " + error.response.data.error;
  }
});

provide('BearingOptions',BearingOptions);
provide('BearingValues',BearingValues);
// provide('searchData',searchData);

// function onBearingSelected(name) {
//   const targetRow = editableRows.value.find(r => r.Name === 'BearingName');
//   if (targetRow) {
//     // 선택된 항목이 BearingOptions에 없으면 추가
//     if (!BearingOptions.value.includes(name)) {
//       BearingOptions.value.splice(BearingOptions.value.length - 1, 0, name);
//     }
//     targetRow.Value = name;
//   }
//   searchModalOpen.value = false;
// }

function onBearingSelected(selectedName) {
  const row = targetRow.value
  if (!row) return

  // 옵션에 없으면 추가
  if (!BearingOptions.value.includes(selectedName)) {
    BearingOptions.value.push(selectedName)
  }

  row.Value = selectedName
  onBearingChanged(row)
  searchModalOpen.value = false
}

function openSearchModal(row) {
  targetRow.value = row
  searchModalOpen.value = true
}

const onBearingChanged = (bearingRow) => {
  const selectedName = bearingRow.Value;
  const groupKey = `${bearingRow.Assembly}_${bearingRow.AssemblyID}`;

  const selectedBearing = BearingValues.value.find(
    (b) => b.Name === selectedName
  );
  if (!selectedBearing) return;

  editableRows.value.forEach((r) => {
    const rGroupKey = `${r.Assembly}_${r.AssemblyID}`;
    if (rGroupKey !== groupKey) return;

    if (["BPFO", "BPFI", "BSF", "FTF"].includes(r.Name)) {
      const paramVal = selectedBearing[r.Name];
      if (paramVal !== undefined) {
        r.Value = Number(paramVal);
      }
    }
  });
};

const validateValue = (row) => {
  const min = getMinValue(row);
  const max = getMaxValue(row);

  // Min보다 작으면 Min 값으로 설정
  if (row.Value < min) {
    row.Value = min;
  }
  // Max보다 크면 Max 값으로 설정
  if (row.Value > max) {
    row.Value = max;
  }
};

const filteredEditableRows = computed(() => {
  const grouped = {};

  editableRows.value.forEach((row) => {
    const key = `${row.Assembly}_${row.AssemblyID}`;
    if (!grouped[key]) grouped[key] = [];
    grouped[key].push(row);
  });

  const result = [];

  for (const key in grouped) {
    const group = grouped[key];
    const bearingType = group.find((r) => r.Name === "BearingType");
    const value = parseInt(bearingType?.Value); // ✅ 확실한 정수 변환
    const allowOthers = !(bearingType && value < 2);

    for (const row of group) {
      if (
        row.Name === "BearingName" &&
        !row.Value &&
        BearingOptions.value.length > 0
      ) {
        row.Value = BearingOptions.value[0]; // ✅ 기본 선택값 설정
      }
      if (row.Name === "BearingType" || allowOthers) {
        result.push(row);
      }
    }
  }

  return result;
});

const handleFileUpload = (event) => {
  selectedFile.value = event.target.files[0]; // 파일 객체 저장
};
// const OptionChanged = (row) =>{
//   console.log('OptionChanged called with:', row.Value);
//   //searchData = event.target.name;
//   if (row.Value === 'Search..') {
//     searchModalOpen.value = true;
//   }else{
//     onBearingChanged(row);
//   }
// }

const uploadBearing = async () => {
  if (!selectedFile.value) {
    message.value = "파일을 선택하세요!";
    return;
  }

  const formData = new FormData();
  formData.append("file", selectedFile.value);
  //console.log(formData);
  try {
    const response = await axios.post("/setting/uploadBearing", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    if (response.data.passOK == "1") {
      message.value = "Upload Success : " + response.data.file_path;
      feedbackModalOpen.value = false;
      const uploadData = response.data.data;
      for (let i = 0; i < uploadData.length; i++) {
        BearingOptions.value.push(uploadData[i]["Name"]);
        BearingValues.value.push(uploadData[i]);
      }
    } else {
      message.value = response.data.error;
    }
  } catch (error) {
    message.value = "업로드 실패: " + error.response.data.error;
  }
};

// const saveChanges = async () => {
//   try {
//     const validation = validateTableData();
//     if (!validation.isValid) {
//       alert(validation.message);
//       return;
//     }
//     await setAssetConfig();
//   } catch (error) {
//     console.error("Error occurred while saving:", error);
//   }
// };

const saveAdvanced = async () => {
  if (!validateModalData()) return;
  try {
    const filtered = modalData.value.map((item) => ({ ...item }));
    const response = await axios.post(
      `/setting/setAssetConfig/${inputDict.value.assetInfo.name}`,
      filtered,
      {
        headers: { "Content-Type": "application/json" },
      }
    );
    if (response.data?.success) {
      alert("✅ Advanced settings saved!");
      showAdvancedModal.value = false;
    } else {
      alert("❌ Save failed: " + (response.data.error || "Unknown error"));
    }
  } catch (error) {
    console.error("Advanced save error:", error);
    alert("❌ Error: " + error.message);
  }
};
// const validateTableData = () => {
//   for (const row of tableData.value) {
//     const min = getMinValue(row);
//     const max = getMaxValue(row);
//     if (row.Value < min || row.Value > max) {
//       return {
//         isValid: false,
//         message: `${row.Title} 값이 범위를 벗어났습니다. (${min} ~ ${max})`,
//       };
//     }
//   }
//   return { isValid: true };
// };

const validateModalData = () => {
  let firstInvalid = null;

  for (const row of modalData.value) {
    const min = getMinValue(row);
    const max = getMaxValue(row);
    const val = Number(row.Value);

    row.invalid = false;
    row.tooltip = "";

    if (isNaN(val)) {
      row.invalid = true;
      row.tooltip = "숫자를 입력하세요";
      if (!firstInvalid) firstInvalid = row;
      continue;
    }

    if ((min !== null && val < min) || (max !== null && val > max)) {
      row.invalid = true;
      row.tooltip = `범위: ${min} ~ ${max}`;
      if (!firstInvalid) firstInvalid = row;
    }
  }

  if (firstInvalid) {
    setTimeout(() => {
      const el = document.getElementById(`modal-input-${firstInvalid.Title}`);
      if (el) el.focus();
    }, 100);
    return false;
  }

  return true;
};
</script>

<style scoped>
table {
  table-layout: fixed;
  width: 100%;
  min-height: calc(40px * 20); /* 40px * 20줄 높이 고정 */
  border-collapse: collapse;
}

th,
td {
  text-align: center;
  padding: 8px;
  border: 1px solid #ddd;
  height: 40px; /* 각 행 높이 고정 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (prefers-color-scheme: dark) {
  th,
  td {
    border-color: rgba(55, 65, 81, 0.6);
  }
}

th:nth-child(1) {
  width: 300px;
} /* 첫 번째 컬럼 */
th:nth-child(2) {
  width: 150px;
}
th:nth-child(3) {
  width: 150px;
}
th:nth-child(4) {
  width: 200px;
}
th:nth-child(5) {
  width: 100px;
}
.empty-row {
  visibility: hidden; /* 빈 행 안보이게 처리 */
}

/* 다크모드를 위한 포커스 및 호버 스타일 향상 */
input:focus,
select:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}

button:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.5);
}
</style>