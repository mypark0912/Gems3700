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
      <header class="flex items-center mb-2">
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
          Diagnostic asset management : {{ selectedAsset.name }}
        </h3>
      </header>
    </div>
    <div class="px-4 py-3">
      <div class="grid grid-cols-[4fr_2fr_5fr_2fr] gap-6 items-start">
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
            Add Asset
          </button>
          <button
            v-if="assetMode.name !== ''"
            @click="unregisterAsset"
            class="bg-sky-500 text-white px-4 py-2 rounded"
          >
            Unregister Asset
          </button>
          <button
            @click="registerAsset"
            class="bg-green-500 text-white px-4 py-2 rounded"
          >
            Register Asset
          </button>
          <button
            @click="showCheckDelete = true"
            class="bg-red-500 text-white px-4 py-2 rounded"
          >
            Delete Asset
          </button>
        </div>

        <!-- 3. 상세 설정 영역 -->
        <div v-if="selectedbtn > 0" class="space-y-4">
          <!-- 1행: Asset Type + Asset Name 수평 정렬 -->
          <div class="grid grid-cols-2 gap-4">
            <div class="flex items-center gap-2">
              <label class="w-24 text-sm font-medium">Asset Type</label>
              <select
                v-model="assetMode.type"
                class="form-select w-full bg-gray-100"
                :disabled="selectedbtn == 2"
              >
                <option v-for="item in asseTypeList" :key="item" :value="item">
                  {{ item }}
                </option>
              </select>
            </div>

            <div class="flex items-center gap-2">
              <label class="w-24 text-sm font-medium">Asset Name</label>
              <input
                v-model="assetMode.name"
                type="text"
                class="form-input w-full"
              />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div v-if="showEquip" class="flex items-center gap-2">
              <label class="w-24 text-sm font-medium">Equipment Name</label>
              <input
                v-model="inputDict.assetInfo.nickname"
                type="text"
                class="form-input w-full"
              />
            </div>
            <div
              v-if="assetMode.type == 'Transformer'"
              class="flex items-center gap-2"
            >
              <label class="w-48 text-sm font-medium"
                >Transformer Capacity</label
              >
              <input
                v-model="inputDict.n_kva"
                type="text"
                class="form-input w-full"
              />
            </div>
          </div>

          <!-- 버튼 영역 -->
          <div class="flex gap-2">
            <button
              v-if="selectedbtn == 2"
              @click="updateAsset"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              Change
            </button>
            <button
              v-if="selectedbtn == 1"
              @click="createAsset"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              Create
            </button>
            <button
              v-if="selectedbtn >= 1"
              @click="cancel"
              class="bg-gray-500 text-white px-3 py-2 rounded"
            >
              Cancel
            </button>
          </div>
        </div>

        <!--div v-if="selectedbtn > 0" class="space-y-4">
                    <div class="flex items-center gap-2">
                      <label class="w-24 text-sm font-medium">Asset Type</label>
                      <select v-model="assetMode.type" class="form-select w-full bg-gray-100" :disabled="selectedbtn == 2">
                        <option value="Transformer">Transformer</option>
                        <option value="PowerSupply">PowerSupply</option>
                        <option value="Motor">Motor</option>
                        <option value="Pump">Pump</option>
                        <option value="Fan">Fan</option>
                        <option value="Compressor">Compressor</option>
                      </select>
                    </div>

                    <div class="flex items-center gap-2">
                      <label class="w-24 text-sm font-medium">Asset Name</label>
                      <input v-model="assetMode.name" type="text" class="form-input w-full" />
                    </div>
                    <div v-if="assetMode.type == 'Transformer'" class="flex items-center gap-2">
                      <label class="w-24 text-sm font-medium">Transformer Capacity</label>
                      <input v-model="inputDict.n_kva" type="text" class="form-input w-full" />
                    </div>

                    <div class="flex gap-2">
                      <button v-if="selectedbtn == 2" @click="updateAsset" class="bg-gray-500 text-white px-3 py-2 rounded">Change</button>
                      <button v-if="selectedbtn == 1" @click="createAsset" class="bg-gray-500 text-white px-3 py-2 rounded">Create</button>
                      <button v-if="selectedbtn >= 1" @click="cancel" class="bg-gray-500 text-white px-3 py-2 rounded">Cancel</button>
                    </div>
                  </div-->

        <!-- 4. Edit Configuration 영역 -->
        <div v-if="selectedbtn == 2" class="space-y-3">
          <label class="text-sm font-medium block">Edit Configuration</label>
          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              v-model="isEditNameplates"
              class="form-checkbox text-violet-500"
            />
            <span>Nameplates</span>
          </label>
          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              v-model="isEditParameters"
              class="form-checkbox text-violet-500"
            />
            <span>Threshold</span>
          </label>
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
      <!-- 헤더 -->
      <div class="text-sm">
        Do you want to delete this asset? All associated data will be
        permanently removed.
      </div>

      <!-- Footer -->
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
import { useI18n } from "vue-i18n"; //
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
const asseTypeList = ref([]);
const assetMode = ref({ name: "", type: "Motor" });
const diagnosis_detail = inject("diagnosis_detail");

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
const existAsset = computed(() => setupStore.getAssetConfig);

const m_kva = computed(() => setupStore.getMkva);
const s_kva = computed(() => setupStore.getSkva);

const showEquip = computed(() => {
  console.log(locale.value);
  if(locale.value !== 'en')
    return true;
});

// locale 변경 시 nickname 업데이트
watch(locale, (newLocale) => {
  console.log('Locale changed to:', newLocale);
  if(newLocale === 'en') {
    inputDict.value.assetInfo.nickname = inputDict.value.assetInfo.name;
  }
}, { immediate: true });

function onCheckChange({ id, name, type, connected, connectedCh }) {
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
  }
  //console.log(selectedAsset.value);
}

watch(
  () => inputDict.value.assetInfo.name, // inputDict.assetName을 감시
  (newValue, oldValue) => {
    if (newValue) {
      tempAssetName.value = newValue;
      if(newLocale === 'en')
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

watch(
  () => checkList.value, // inputDict.assetName을 감시
  (newValue, oldValue) => {
    //console.log("watchCheckList", newValue);
    if (newValue) {
      if (channel.value == "Main") {
        if (existAsset.value.assetName_main != "") {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            existAsset.value.assetName_main,
            existAsset.value.assetName_sub,
          ]);
          const result = findAssetByName(
            assetList.value,
            existAsset.value.assetName_main
          );
          //console.log('ChekcList, compare result:',result,'assetName:',existAsset.value.assetName_main);
          if (result) {
            selectedAsset.value.name = existAsset.value.assetName_main;
            selectedAsset.value.type = existAsset.value.assetType_main;
            selectedAsset.value.connected = true;
            selectedAsset.value.connectedCh = result.connectedCh;
            selectedAsset.value.id = result.id;
            selectedbtn.value = 2;
            if (selectedAsset.value.connected) {
              selectedbtn.value = 2;
              assetMode.value.name = selectedAsset.value.name;
              assetMode.value.type = selectedAsset.value.type;
              if (selectedAsset.value.type == "Transformer" && m_kva.value > 0)
                inputDict.value.n_kva = m_kva.value;
            } else {
              selectedbtn.value = 0;
            }
          }
        } else {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            existAsset.value.assetName_main,
            existAsset.value.assetName_sub,
          ]);
        }
      } else {
        if (existAsset.value.assetName_sub != "") {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            existAsset.value.assetName_main,
            existAsset.value.assetName_sub,
          ]);
          const result = findAssetByName(
            assetList.value,
            existAsset.value.assetName_sub
          );
          //console.log('compare result:',result,'assetName:',existAsset.value.assetName_sub);
          if (result) {
            selectedAsset.value.name = existAsset.value.assetName_sub;
            selectedAsset.value.type = existAsset.value.assetType_sub;
            selectedAsset.value.connected = true;
            selectedAsset.value.connectedCh = result.connectedCh;
            selectedAsset.value.id = result.id;
            selectedbtn.value = 2;
            if (selectedAsset.value.connected) {
              selectedbtn.value = 2;
              assetMode.value.name = selectedAsset.value.name;
              assetMode.value.type = selectedAsset.value.type;
              if (selectedAsset.value.type == "Transformer" && s_kva.value > 0)
                inputDict.value.n_kva = s_kva.value;
            } else {
              selectedbtn.value = 0;
            }
          }
        } else {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            existAsset.value.assetName_main,
            existAsset.value.assetName_sub,
          ]);
        }
      }
    }
  },
  { immediate: true }
);

watch(
  () => existAsset.value, // inputDict.assetName을 감시
  (newValue, oldValue) => {
    console.log(newValue);
    if (newValue) {
      if (channel.value == "Main") {
        if (newValue.assetName_main != "") {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            newValue.assetName_main,
            newValue.assetName_sub,
          ]);
          const result = findAssetByName(
            assetList.value,
            newValue.assetName_main
          );
          //console.log('existAsset compare result:',result,'assetName:',newValue.assetName_main);
          if (result) {
            selectedAsset.value.name = newValue.assetName_main;
            selectedAsset.value.type = newValue.assetType_main;
            selectedAsset.value.connected = true;
            selectedAsset.value.connectedCh = result.connectedCh;
            selectedAsset.value.id = result.id;
            selectedbtn.value = 2;
            if (selectedAsset.value.connected) {
              selectedbtn.value = 2;
              assetMode.value.name = selectedAsset.value.name;
              assetMode.value.type = selectedAsset.value.type;
            } else {
              selectedbtn.value = 0;
            }
          }
        } else {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            newValue.assetName_main,
            newValue.assetName_sub,
          ]);
        }
      } else {
        if (newValue.assetName_sub != "") {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            newValue.assetName_main,
            newValue.assetName_sub,
          ]);
          const result = findAssetByName(
            assetList.value,
            newValue.assetName_sub
          );
          //console.log('compare result:',result,'assetName:',newValue.assetName_sub);
          if (result) {
            selectedAsset.value.name = newValue.assetName_sub;
            selectedAsset.value.type = newValue.assetType_sub;
            selectedAsset.value.connected = true;
            selectedAsset.value.connectedCh = result.connectedCh;
            selectedAsset.value.id = result.id;
            selectedbtn.value = 2;
            if (selectedAsset.value.connected) {
              selectedbtn.value = 2;
              assetMode.value.name = selectedAsset.value.name;
              assetMode.value.type = selectedAsset.value.type;
            } else {
              selectedbtn.value = 0;
            }
          }
        } else {
          assetList.value = transformGroupedAssetObject(checkList.value, [
            newValue.assetName_main,
            newValue.assetName_sub,
          ]);
        }
      }
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

const getAssetList = async () => {
  try {
    const response = await axios.get(`/setting/getAssetList`);

    if (response.data.success) {
      //const groupedAssets = response.data.data;
      checkList.value = response.data.data;
      //console.log('getAssetList',checkList.value);
    }
  } catch (error) {
    console.log(error);
  }
};

const getAssetTypes = async () => {
  try {
    const response = await axios.get(`/setting/getAssetTypes`);

    if (response.data.success) {
      //const groupedAssets = response.data.data;
      asseTypeList.value = response.data.data;
      console.log(asseTypeList.value);
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
    }
  } catch (error) {
    console.log(error);
  }
  if (flag) {
    inputDict.value.assetInfo.name = "";
    inputDict.value.assetInfo.type = "";
    inputDict.value.assetInfo.nickname = "";
    isAssetExist.value = false;
    assetMode.value.name = "";
    assetMode.value.type = "";
    selectedAsset.value.connected = false;
    selectedAsset.value.connectedCh = "";
    selectedbtn.value = 0;

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
      });

     
    } else {
      console.warn(
        `⚠️ diagnosis_detail key '${diagnosisChannelKey}' not found`
      );
    }

    setupStore.setAssetConfig({
      channel: channel.value,
      name: inputDict.value.assetInfo.name,
      type: inputDict.value.assetInfo.type,
      nick: inputDict.value.assetInfo.nickname,
    });

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
  await props.savefile("fromChild");
  try {
    const response = await axios.get(
      `/setting/registerAsset/${channel.value}/${selectedAsset.value.name}/${selectedAsset.value.type}`
    );

    if (response.data.success) {
      flag = true;
    }
  } catch (error) {
    console.log(error);
  }
  if (flag) {
    inputDict.value.assetInfo.name = selectedAsset.value.name; //tempAssetName.value; // 수정된 이름 반영
    inputDict.value.assetInfo.type = selectedAsset.value.type;

    selectedAsset.value.connected = true;
    isAssetExist.value = true;
    // setupStore.setAssetConfig(channel.value, inputDict.value.assetInfo.name,inputDict.value.assetInfo.type,inputDict.value.assetInfo.nickname );
    setupStore.setAssetConfig({
      channel: channel.value,
      name: inputDict.value.assetInfo.name,
      type: inputDict.value.assetInfo.type,
      nick: inputDict.value.assetInfo.nickname,
    });
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

  try {
    // if(assetMode.value.type == 'Transformer'){
    //   const kva = parseInt(inputDict.value.n_kva);
    //   inputDict.value.n_kva = kva;
    // }
    const payload = {
      assetType: assetMode.value.type, //inputDict.value.assetInfo.type,
      assetName: assetMode.value.name, //tempAssetName.value,
      assetNickname: assetMode.value.nickname,
    };

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
    console.log("payload", payload);
    const response = await axios.post("/setting/modifyAsset", payload, {
      headers: { "Content-Type": "application/json" },
    });

    if (response.data.status === "1") {
      inputDict.value.assetInfo.name = newName; // ✅ 새 이름 적용
      isEditing.value = false; // ✅ 수정모드 종료
      //setupStore.setAssetConfig(channel.value, inputDict.value.assetInfo.name,inputDict.value.assetInfo.type,inputDict.value.assetInfo.nickname );
      setupStore.setAssetConfig({
        channel: channel.value,
        name: inputDict.value.assetInfo.name,
        type: inputDict.value.assetInfo.type,
        nick: inputDict.value.assetInfo.nickname,
      });
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
        inputDict.value.assetInfo.nickname = "";
        isAssetExist.value = false;
        //setupStore.setAssetConfig(props.channel, inputDict.value.assetInfo.name,inputDict.value.assetInfo.type,inputDict.value.assetInfo.nickname );

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
          });
        }

        setupStore.setAssetConfig({
          channel: channel.value,
          name: inputDict.value.assetInfo.name,
          type: inputDict.value.assetInfo.type,
          nick: inputDict.value.assetInfo.nickname,
        });
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
