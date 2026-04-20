<template>
  <div
    class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center justify-between mb-2">
        <div class="flex items-center">
          <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
            <svg
              class="w-6 h-6 fill-current text-white"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M12 15.5a3.5 3.5 0 1 1 3.5-3.5 3.5 3.5 0 0 1-3.5 3.5zM19.5 12a7.5 7.5 0 0 1-1 3.7l2.1 2.1-1.8 1.8-2.1-2.1a7.5 7.5 0 0 1-3.7 1V21h-2.5v-2.5a7.5 7.5 0 0 1-3.7-1l-2.1 2.1-1.8-1.8 2.1-2.1a7.5 7.5 0 0 1-1-3.7H3v-2.5h2.5a7.5 7.5 0 0 1 1-3.7L4.4 4.4l1.8-1.8 2.1 2.1a7.5 7.5 0 0 1 3.7-1V3h2.5v2.5a7.5 7.5 0 0 1 3.7 1l2.1-2.1 1.8 1.8-2.1 2.1a7.5 7.5 0 0 1 1 3.7H21V12h-1.5z"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
           {{ t('config.channelPanel.Diagnosis.title') }}  : {{ selectedAsset.name }}
          </h3>
        </div>
        
        <!-- Edit Configuration in Header -->
<div v-if="selectedbtn == 2" class="flex items-center gap-4">
  <label class="flex items-center space-x-2">
    <input
      type="checkbox"
      v-model="inputDict.confStatus"
      :true-value="1"
      :false-value="0"
      class="form-checkbox text-violet-500"
    />
    <span class="text-sm">Configure status</span>
  </label>
  <label class="flex items-center space-x-2">
    <input
      type="checkbox"
      v-model="inputDict.useDO"
      :true-value="1"
      :false-value="0"
      class="form-checkbox text-violet-500"
    />
    <span class="text-sm">Use DO Alarm</span>
  </label>
  <label class="flex items-center space-x-2">
    <input
      type="checkbox"
      v-model="inputDict.useAI"
      :true-value="1"
      :false-value="0"
      class="form-checkbox text-violet-500"
    />
    <span class="text-sm">Use AI Setup</span>
  </label>
  <label class="flex items-center space-x-2">
    <input
      type="checkbox"
      v-model="isEditNameplates"
      :true-value="1"
      :false-value="0"
      class="form-checkbox text-violet-500"
    />
    <span class="text-sm">Edit Nameplates</span>
  </label>
  <label class="flex items-center space-x-2">
    <input
      type="checkbox"
      v-model="isEditParameters"
      :true-value="1"
      :false-value="0"
      class="form-checkbox text-violet-500"
    />
    <span class="text-sm">Edit Threshold</span>
  </label>
</div>
      </header>
    </div>
    <div class="px-4 py-3">
      <div class="grid grid-cols-[4fr_2fr_6fr] gap-6 items-start">
        <!-- 1. Tree 테이블 -->
        <div class="border rounded shadow-sm">
          <table class="w-full">
            <tbody>
              <TreeRowAsset
                v-for="item in assetList"
                :key="item.ID"
                :item="item"
                :level="0"
                :checked-id="selectedAsset.id"
                :checked-name="selectedAsset.name"
                :channel="channel"
                :expanded-id="selectedAsset.id"
                @check-change="onCheckChange"
              />
            </tbody>
          </table>
        </div>

        <!-- 2. 버튼 영역 -->
        <div class="flex flex-col gap-3">
          <button
            @click="addAsset"
            class="bg-blue-500 text-white px-4 py-2 rounded"
          >
           {{ t('config.channelPanel.Diagnosis.AddAsset') }}
          </button>
          <button
            v-if="assetMode.name !== '' && selectedbtn == 2"
            @click="unregisterAsset"
            class="bg-sky-500 text-white px-4 py-2 rounded"
          >
           {{ t('config.channelPanel.Diagnosis.UnregisterAsset') }}
          </button>
          <button
            v-if="selectedbtn != 2"
            @click="registerAsset"
            class="bg-green-500 text-white px-4 py-2 rounded"
          >
            {{ t('config.channelPanel.Diagnosis.registerAsset') }}
          </button>
        </div>

        <!-- 3. 상세 설정 영역 -->
        <div class="space-y-4">
          <!-- 1행: Asset Type + Drive Type -->
          <div class="grid grid-cols-2 gap-4">
            <div class="flex items-center gap-2">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">{{ t('config.channelPanel.Diagnosis.AssetType') }}</label>
              <select
                v-model="assetMode.type"
                class="form-select flex-1 bg-gray-100"
                :disabled="selectedbtn != 1"
              >
                <option v-for="item in assetTypeList" :key="item" :value="item">
                  {{ item }}
                </option>
              </select>
            </div>

            <div v-if="canSelectDriveType" class="flex items-center gap-2">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">Drive Type</label>
              <select
                v-model="inputDict.assetInfo.driveType"
                class="form-select flex-1 bg-gray-100"
                :disabled="selectedbtn != 1"
              >
                <option value="DOL">DOL</option>
                <option value="VFD">VFD</option>
              </select>
            </div>
          </div>

          <!-- 2행: Asset Name + Equipment Name -->
          <div class="grid grid-cols-2 gap-4">
            <div class="flex items-center gap-2" v-if="selectedbtn >= 1">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">{{ t('config.channelPanel.Diagnosis.AssetName') }}</label>
              <input
                v-model="assetMode.name"
                type="text"
                class="form-input flex-1"
                :disabled="selectedbtn == 2"
                :maxlength="20"
              />
            </div>
            <div class="flex items-center gap-2" v-if="selectedbtn == 0">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">{{ t('config.channelPanel.Diagnosis.AssetName') }}</label>
              <input
                v-model="assetMode.newname"
                type="text"
                class="form-input flex-1"
                :maxlength="20"
              />
            </div>

            <div v-if="showEquip && selectedbtn >= 2" class="flex items-center gap-2">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">{{ t('config.channelPanel.Diagnosis.EquipmentName') }}</label>
              <input
                v-model="inputDict.assetInfo.nickname"
                type="text"
                class="form-input flex-1"
                :maxlength="20"
              />
            </div>
          </div>

          <!-- 3행: Transformer Capacity (변압기일 때만) -->
          <div v-if="assetMode.type == 'Transformer' && selectedbtn == 2" class="grid grid-cols-2 gap-4">
            <div class="flex items-center gap-2">
              <label class="w-32 text-sm font-medium shrink-0 whitespace-nowrap">{{ t('config.channelPanel.Diagnosis.transformerCapacity') }}</label>
              <input
                v-model="inputDict.n_kva"
                type="number"
                class="form-input flex-1"
                :maxlength="20"
              />
            </div>
          </div>

          <!-- 버튼 영역 -->
          <div class="flex gap-2">
            <button
              v-if="selectedbtn == 2 && false"
              @click="updateAsset"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              Change
            </button>
            <button
              v-if="selectedbtn == 0"
              @click="renameAsset"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              {{ t('config.channelPanel.Diagnosis.NameChange') }}
            </button>
            <button
              v-if="selectedbtn == 0"
              @click="showCheckDelete = true"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              {{ t('config.channelPanel.Diagnosis.DeleteAsset') }}
            </button>
            <button
              v-if="selectedbtn == 1"
              @click="createAsset"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              {{ t('config.channelPanel.Diagnosis.CreateAsset') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ModalBasic
    :modalOpen="showCheckDelete"
    @close-modal="showCheckDelete = false"
    title="Confirm to delete asset"
  >
    <div class="w-[600px] max-w-full px-6">
      <div class="text-sm">
        Do you want to delete this asset? All associated data will be
        permanently removed.
      </div>

      <div class="px-5 py-4 border-t border-gray-200 dark:border-gray-700/60">
        <div class="flex justify-end space-x-2">
          <button
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            @click="deleteAsset"
          >
            YES
          </button>
          <button
            class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
            @click="showCheckDelete = false"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </ModalBasic>
</template>
<script setup>
import { inject, ref, watch, defineProps, onMounted, computed } from "vue";
import axios from "axios";
import { useSetupStore } from "@/store/setup"; // ✅ Pinia store 사용
import TreeRowAsset from "./TreeRowAsset.vue";
import ModalBasic from "../../../pages/common/ModalBasic.vue";
import { useI18n } from "vue-i18n"; 
const { t } = useI18n();
const { locale } = useI18n();
const inputDict = inject("channel_inputDict");
const changeDiagnosis = inject("changeDiagnosis");
const isEditing = ref(false); // 수정 모드 여부
const tempAssetName = ref(""); // 임시로 수정할 asset name 값
//const savefile = inject('savefile');
const isEditNameplates = inject("isEditNameplates");
const isEditParameters = inject("isEditParameters");
const assetList = ref([]);
const checkList = ref([]);
const selectedbtn = ref(0);
const showCheckDelete = ref(false);
const assetTypeList = ref([
            "PSupply",
            "PrimaryTransformer",
            "Transformer",
            "MotorFeed",
            "Motor",
            "Pump",
            "Fan",
            "Compressor"
        ]);
const assetDriveType = inject('assetDriveType');
//const driveType = ref('DOL');  // 'DOL' 또는 'VFD'


const vfdSupportedTypes = [
    "MotorFeed",
    "Motor", 
    "Pump",
    "Fan",
    "Compressor"
];
const assetMode = ref({ name: "", type: "Motor", newname:"" });
const diagnosis_detail = inject("diagnosis_detail");
const emit = defineEmits(['update-selectedbtn']);
watch(selectedbtn, (newValue) => {
  emit('update-selectedbtn', newValue);
}, { immediate: true });
//const norminal_kva = ref(0);
const props = defineProps({
  channel: String,
  isAsset: Boolean,
  savefile: Function,
});

const setupStore = useSetupStore(); // ✅ Pinia store 사용
const isAssetExist = ref(false);
const channel = ref(props.channel);
const selectedAsset = ref({
  id: "",
  name: "",
  type: "",
  connected: false,
  connectedCh: "",
});
//console.log('savefile',savefile);
const canSelectDriveType = computed(() => {
    return vfdSupportedTypes.includes(assetMode.value.type);
});

// 실제 전송할 AssetType
const finalAssetType = computed(() => {
    if (driveType.value === 'VFD' && canSelectDriveType.value) {
        return `${assetMode.value.type}-VFD`;
    }
    return assetMode.value.type;
});

// AssetType 변경 시 DriveType 초기화
watch(() => assetMode.value.type, (newType) => {
    if (!vfdSupportedTypes.includes(newType)) {
        inputDict.value.assetInfo.driveType = 'DOL';
    }
});

// watch(assetDriveType, (newValue) => {
//   inputDict.value.assetInfo.driveType = newValue;
// },{ immediate: true });

const m_kva = computed(() => setupStore.getMkva);
const s_kva = computed(() => setupStore.getSkva);

const showEquip = computed(() => {
  console.log(locale.value);
  if(locale.value !== 'en')
    return true;
});

// locale 변경 시 nickname 업데이트
watch(locale, (newLocale) => {
  //console.log('Locale changed to:', newLocale);
  if(newLocale === 'en') {
    inputDict.value.assetInfo.nickname = inputDict.value.assetInfo.name;
  }
}, { immediate: true });

function onCheckChange({ id, name, type, connected, connectedCh }) {
  console.log('id, name, type, connected, connectedCh',id, name, type, connected, connectedCh);
  selectedAsset.value.id = id;
  selectedAsset.value.name = name;
  selectedAsset.value.type = type;
  selectedAsset.value.connected = connected;
  if (selectedAsset.value.connected) {
    selectedAsset.value.connectedCh = connectedCh;
    selectedbtn.value = 2;
    assetMode.value.name = selectedAsset.value.name;
    assetMode.value.type = selectedAsset.value.type;
  } else {
    selectedbtn.value = 0;
    assetMode.value.name = selectedAsset.value.name;
    assetMode.value.newname = selectedAsset.value.name;
    assetMode.value.type = selectedAsset.value.type;
  }
  //console.log(selectedAsset.value);
}

watch(
  () => inputDict.value.assetInfo.name, // inputDict.assetName을 감시
  (newValue, oldValue) => {
    if (newValue) {
      tempAssetName.value = newValue;
      if(locale.value === 'en')
        inputDict.value.assetInfo.nickname = tempAssetName;
    }
  },
  { immediate: true }
);

watch(
  () => props.isAsset, // inputDict.assetName을 감시
  (newValue, oldValue) => {
    if (newValue) {
      isAssetExist.value = newValue;
    }
  },
  { immediate: true }
);

function transformGroupedAssetObject(obj, allConnectedAssetNames = []) {
  //console.log('list',allConnectedAssetNames);
  //const validNames = allConnectedAssetNames.filter(name => name && name.trim() !== "");
  const mainName = allConnectedAssetNames[0] || null;
  const subName = allConnectedAssetNames[1] || null;

  return Object.entries(obj).map(([type, items]) => {
    const children = items.map((asset, index) => {
      let connectedCh = "";

      if (asset.Name === mainName) {
        connectedCh = "Main";
      } else if (asset.Name === subName) {
        connectedCh = "Sub";
      }

      return {
        ID: `${type}-${index}`,
        Type: asset.Type,
        Connected: connectedCh !== "",
        connectedCh: connectedCh,
        Name: asset.Name,
        Title: asset.Name,
      };
    });

    const isExpanded = children.some(
      (child) => child.connectedCh === channel.value
    );

    return {
      ID: type,
      Type: type,
      Title: type,
      isExpanded,
      children: children,
    };
  });
}

const addAsset = () => {
  selectedbtn.value = 1;
  assetMode.value.name = "";
  assetMode.value.type = "Motor";
};

const cancel = () => {
  selectedbtn.value = 2;
  assetMode.value.name = selectedAsset.value.name;
  assetMode.value.type = selectedAsset.value.type;
};

const resetAssetInfo = () => {
  inputDict.value.assetInfo.name = "";
  inputDict.value.assetInfo.type = "";
  inputDict.value.assetInfo.nickname = "";
  isAssetExist.value = false;
  assetMode.value.name = "";
  assetMode.value.type = "";
  selectedAsset.value = {
    id: "",
    name: "",
    type: "",
    connected: false,
    connectedCh: "",
  };
  
  selectedbtn.value = -1;
  console.log('Asset reset - ', inputDict.value.assetInfo);
};

const getAssetList = async () => {
  try {
    const response = await axios.get(`/setting/getAssetList`);

    if (response.data.success) {
      //const groupedAssets = response.data.data;
      checkList.value = response.data.data;
      //console.log('getAssetList',checkList.value);
      if (Object.keys(checkList.value).length === 0){
        resetAssetInfo();
      }else{
        const storedAssetName = inputDict.value.assetInfo.name;
        const storedAssetType = inputDict.value.assetInfo.type;
        
        let isStoredAssetFound = false;
        
        if (storedAssetName && storedAssetType) {
          // checkList에서 해당 타입의 리스트를 찾음
          const assetListOfType = checkList.value[storedAssetType];
          //console.log('checkList',checkList.value);
          if (assetListOfType && Array.isArray(assetListOfType)) {
            // 해당 타입 리스트에서 이름이 일치하는지 확인
            isStoredAssetFound = assetListOfType.some(
              asset => asset.Name === storedAssetName
            );
          }
        }
        
        // 스토어에 저장된 asset이 조회 결과에 없으면 초기화
        if (!isStoredAssetFound && storedAssetName) {
          console.log('Stored asset not found in list - resetting');
          resetAssetInfo();
        }
      }
    }
  } catch (error) {
    console.log(error);
  }
};

// const getAssetTypes = async () => {
//   try {
//     const response = await axios.get(`/setting/getAssetTypes`);

//     if (response.data.success) {
//       //const assetTypeList = response.data.data;
//        asseTypeList.value = response.data.data;
//       // console.log(asseTypeList.value);
//     }
//   } catch (error) {
//     console.log(error);
//   }
// };

const getAssetTypes = async () => {
  try {
    const response = await axios.get('/setting/getAssetTypes');
    if (response.data.success) {
      const rawTypes = response.data.data;
      
      // VFD 제거하고 중복 제거
      const uniqueTypes = [...new Set(
        rawTypes.map(type => type.replace('-VFD', ''))
      )];
      
      assetTypeList.value = uniqueTypes;
    }
  } catch (error) {
    console.log(error);
  }
};

function findAssetByName(assetTree, targetName) {
  for (const group of assetTree) {
    const match = group.children.find(
      // Children 대소문자 주의
      (child) => child.Name === targetName
      //child => child.Name.trim() === targetName.trim()
    );
    if (match) {
      return {
        id: match.ID,
        type: match.Type,
        connected: match.Connected,
        connectedCh: match.connectedCh || "", // 없으면 빈 문자열
      };
    }
  }
  return null;
}

onMounted(async () => {
  getAssetList();
  getAssetTypes();
});

// const startEdit = () => {
//       tempAssetName.value = inputDict.value.assetInfo.name; // 현재 자산 이름을 임시 변수에 저장
//       isEditing.value = true;
//     };

const unregisterAsset = async () => {
  if (
    selectedAsset.value.connected &&
    selectedAsset.value.connectedCh &&
    selectedAsset.value.connectedCh !== channel.value
  ) {
    alert(
      `❌ This asset is already registered on ${selectedAsset.value.connectedCh} channel. Please select another asset or select asset connected this channel`
    );
    return; // 등록 중단
  }

  // 🔥 추가 검증: 선택된 자산이 없는 경우
  if (!selectedAsset.value.name || !selectedAsset.value.type) {
    alert("❌ Please select an asset to unregister.");
    return;
  }
  let flag = false;
  try {
    const response = await axios.get(
      `/setting/unregisterAsset/${channel.value}/${selectedAsset.value.name}`
    );

    if (response.data.success) {
      flag = true;
    }else {
    // ✅ 에러 메시지 배열 처리
    const errorMessages = Array.isArray(response.data.error)
      ? response.data.error.join('\n')
      : response.data.error || 'Unknown error';
    alert(`❌ Unregister failed:\n${errorMessages}`);
    return;
  }
  } catch (error) {
    console.log(error);
  }
  if (flag) {

    // 🔥 status_Info 초기화 추가
    if (inputDict.value.status_Info) {
      inputDict.value.status_Info.diagnosis = [];
      inputDict.value.status_Info.pq = [];
      inputDict.value.status_Info.faults = [];
      inputDict.value.status_Info.events = [];
      inputDict.value.confStatus = 0;
      inputDict.value.useDO = 0;
      inputDict.value.useAI = 0;

    }

    inputDict.value.assetInfo.name = "";
    inputDict.value.assetInfo.type = "";
    //inputDict.value.assetInfo.nickname = "";
    isAssetExist.value = false;
    assetMode.value.name = "";
    assetMode.value.type = "";
    selectedAsset.value.connected = false;
    selectedAsset.value.connectedCh = "";
    selectedbtn.value = 0;

    // ✅ TreeRow 체크 해제를 위한 selectedAsset 초기화
    selectedAsset.value = {
      id: "",
      name: "",
      type: "",
      connected: false,
      connectedCh: "",
    };
    // ✨ 채널명 매핑 및 diagnosis_detail 초기화
    const channelMapping = {
      Main: "main",
      Sub: "sub",
    };

    const diagnosisChannelKey =
      channelMapping[channel.value] || channel.value.toLowerCase();

    if (diagnosis_detail.value[diagnosisChannelKey]) {
  

      // diagnosis_detail 완전 초기화
      Object.assign(diagnosis_detail.value[diagnosisChannelKey], {
        use: false,
        assetName: "",
        tableData: [],
        paramData: [],
        modalData: [],
      });

     
    } else {
      console.warn(
        `⚠️ diagnosis_detail key '${diagnosisChannelKey}' not found`
      );
    }

    selectedbtn.value = -1;
    await getAssetList();
    alert("Asset is unregistered on this channel");
  } else {
    alert("Unregistering asset is failed");
  }
};

const registerAsset = async () => {
  //console.log(selectedAsset.value.connectedCh, channel.value);
  if (
    selectedAsset.value.connected &&
    selectedAsset.value.connectedCh &&
    selectedAsset.value.connectedCh !== channel.value
  ) {
    alert(
      `❌ This asset is already registered on ${selectedAsset.value.connectedCh} channel. Please select another asset or unregister it first.`
    );
    return; // 등록 중단
  }

  // 🔥 추가 검증: 선택된 자산이 없는 경우
  if (!selectedAsset.value.name || !selectedAsset.value.type) {
    alert("❌ Please select an asset to register.");
    return;
  }

  let flag = false;
  // try {
  //   if (props.savefile && typeof props.savefile === 'function') {
  //     await props.savefile("fromChild");
  //   }
  // } catch (e) {
  //   console.warn("savefile 호출 실패:", e);
  // }

  try {
    const response = await axios.get(
      `/setting/registerAsset/${channel.value}/${selectedAsset.value.name}/${selectedAsset.value.type}`
    );

    if (response.data.success) {
      flag = true;
    } else {
    // ✅ 에러 메시지 배열 처리
    const errorMessages = Array.isArray(response.data.error)
      ? response.data.error.join('\n')
      : response.data.error || 'Unknown error';
    alert(`❌ Registration failed:\n${errorMessages}`);
    return; // 에러 시 함수 종료
  }
  } catch (error) {
    console.log(error);
  }
  if (flag) {
    inputDict.value.assetInfo.name = selectedAsset.value.name; //tempAssetName.value; // 수정된 이름 반영
    inputDict.value.assetInfo.type = selectedAsset.value.type;

    selectedAsset.value.connected = true;
    isAssetExist.value = true;
    if (selectedAsset.value.connected) {
      selectedbtn.value = 2;
      assetMode.value.name = selectedAsset.value.name;
      assetMode.value.type = selectedAsset.value.type;
    } else {
      selectedbtn.value = 0;
    }
    await getAssetList();
    alert("✅ Asset is registered on this channel");
  } else {
    alert("❌ Registering asset is failed. Save this channel information");
  }
};

const createAsset = async () => {
  //tempAssetName.value = assetMode.value.name;

  if (!assetMode.value.name || !assetMode.value.type) {
    alert("❌ Please enter both Asset Type and Asset Name.");
    return;
  }
  const validNamePattern = /^[a-zA-Z0-9]+$/;
  
  if (!validNamePattern.test(assetMode.value.name)) {
    alert("❌ Asset name can only contain English letters and numbers. No spaces, special characters, or Korean characters allowed.");
    return;
  }
  try {
    // if(assetMode.value.type == 'Transformer'){
    //   const kva = parseInt(inputDict.value.n_kva);
    //   inputDict.value.n_kva = kva;
    // }
    

    //getAssetList();

    // 이름 중복 체크
    const newAssetName = assetMode.value.name;
    let isDuplicate = false;

    // checkList.value가 Proxy 객체이므로 Object.values()로 배열들을 가져옴
    const assetArrays = Object.values(checkList.value);
    
    for (const assetArray of assetArrays) {
      if (Array.isArray(assetArray)) {
        const duplicate = assetArray.find(asset => asset.Name === newAssetName);
        if (duplicate) {
          isDuplicate = true;
          break;
        }
      }
    }

    if (isDuplicate) {
      alert("❌ Asset name already exists. Please choose a different name.");
      return; // 중복이면 함수 종료
    }
    let sendType;
    if(inputDict.value.assetInfo.driveType == 'VFD')
      sendType = assetMode.value.type +"-" + inputDict.value.assetInfo.driveType;
    else
      sendType = assetMode.value.type;
    //console.log(sendType);
    const payload = {
      assetType: sendType, //inputDict.value.assetInfo.type,
      assetName: assetMode.value.name, //tempAssetName.value,
      assetNickname: assetMode.value.nickname,
    };
    console.log(payload);
    const response = await axios.post("/setting/createAsset", payload, {
      headers: { "Content-Type": "application/json" },
    });
    if (response.data.status === "1") {
      //savefile(); // 저장 호출
      getAssetList();
      changeDiagnosis.value.asset = true;
      alert("Asset is created!");
    } else {
      alert("❌ Asset creation failed");
    }
  } catch (error) {
    console.log(error);
    alert("❌ An error occurred while creating the asset");
  }
};


// const cancel = () =>{
//   if (isEditing.value){
//     isEditing.value = false;
//     tempAssetName.value ="";
//   }
// }
const updateAsset = async () => {
  if (!assetMode.value.name || !assetMode.value.type) {
    //if (!inputDict.value.assetInfo.name || !inputDict.value.assetInfo.type) {
    alert("❌ Please enter both Asset Type and Asset Name.");
    return;
  }
  // if(assetMode.value.type == 'Transformer'){
  //   const kva = parseInt(norminal_kva.value);
  //   inputDict.value.n_kva = kva;
  // }
  
  const oldName = inputDict.value.assetInfo.name;
  const newName = assetMode.value.name;

  // 수정된 이름 서버에 전송
  const payload = {
    assetName: oldName,
    newName: newName,
    assetNickname: assetMode.value.nickname,
  };

  try {
    await unregisterAsset();
    const response = await axios.post("/setting/modifyAsset", payload, {
      headers: { "Content-Type": "application/json" },
    });

    if (response.data.status === "1") {
      inputDict.value.assetInfo.name = newName; // ✅ 새 이름 적용
      inputDict.value.assetInfo.nickname = assetMode.value.nickname; // ✅ 새 닉네임 적용
      isEditing.value = false; // ✅ 수정모드 종료

      await getAssetList();
      //savefile(); // ✅ 저장 호출
      changeDiagnosis.value.asset = true;
     
      alert("✅ Asset Updated!");
    } else {
      alert("❌Asset name modification failed: " + (response.data.error || ""));
    }
  } catch (err) {
    alert("❌Error occurred while modifying Asset: " + err.message);
  }
};
const renameAsset = async () => {
  if (!assetMode.value.name || !assetMode.value.type) {
    //if (!inputDict.value.assetInfo.name || !inputDict.value.assetInfo.type) {
    alert("❌ Please enter both Asset Type and Asset Name.");
    return;
  }
  const validNamePattern = /^[a-zA-Z0-9]+$/;
  
  if (!validNamePattern.test(assetMode.value.newname)) {
    alert("❌ Asset name can only contain English letters and numbers. No spaces, special characters, or Korean characters allowed.");
    return;
  }
  const oldName = assetMode.value.name;
  const newName = assetMode.value.newname;

  let isDuplicate = false;

  const assetArrays = Object.values(checkList.value);
  
  for (const assetArray of assetArrays) {
    if (Array.isArray(assetArray)) {
      const duplicate = assetArray.find(asset => asset.Name === newName);
      if (duplicate) {
        isDuplicate = true;
        break;
      }
    }
  }

  if (isDuplicate) {
    alert("❌ Asset name already exists. Please choose a different name.");
    return; // 중복이면 함수 종료
  }
  // 수정된 이름 서버에 전송
  const payload = {
    assetName: oldName,
    newName: newName,
    assetNickname: assetMode.value.nickname,
  };

  console.log('renameAsset payload',payload);

  try {
    
    const response = await axios.post("/setting/modifyAsset", payload, {
      headers: { "Content-Type": "application/json" },
    });

    if (response.data.status === "1") {
      inputDict.value.assetInfo.name = newName; // ✅ 새 이름 적용

      isEditing.value = false; // ✅ 수정모드 종료
      await getAssetList();
      //savefile(); // ✅ 저장 호출
      changeDiagnosis.value.asset = true;
     
      alert("✅ Asset name Updated!");
    } else {
      alert("❌Asset name modification failed: " + (response.data.error || ""));
    }
  } catch (err) {
    alert("❌Error occurred while modifying Asset: " + err.message);
  }
};
const deleteAsset = async () => {
  const assetName = selectedAsset.value.name;

  if (!assetName) {
    alert("❌ There are no assets to delete.");
    return;
  }

  try {
    const response = await axios.get(`/setting/deleteAsset/${assetName}`);

    if (response.data.status === "1") {
      // ✅ 삭제 성공 시 로컬 상태 초기화
      if (selectedAsset.value.connected) {
        inputDict.value.assetInfo.name = "";
        inputDict.value.assetInfo.type = "";
        //inputDict.value.assetInfo.nickname = "";
        isAssetExist.value = false;

        const channelMapping = {
          Main: "main",
          Sub: "sub",
        };

        const diagnosisChannelKey =
          channelMapping[channel.value] || channel.value.toLowerCase();

        if (
          diagnosis_detail.value &&
          diagnosis_detail.value[diagnosisChannelKey]
        ) {
          Object.assign(diagnosis_detail.value[diagnosisChannelKey], {
            use: false,
            assetName: "",
            tableData: [],
            paramData: [],
            modalData: [],
          });
        }

      }
      selectedAsset.value = {
        id: "",
        name: "",
        type: "",
        connected: false,
      };
      await getAssetList();
      changeDiagnosis.value.asset = true;
      //savefile();
      alert("✅ Asset deletion complete!");
      showCheckDelete.value = false;
    } else {
      alert("❌ Asset deletion failed: " + (response.data.error || ""));
    }
  } catch (error) {
    alert("❌ Error occurred while deleting: " + error.message);
  }
};
</script>
