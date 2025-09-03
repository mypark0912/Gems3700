<template>
    <div class="flex flex-nowrap overflow-x-scroll no-scrollbar md:block md:overflow-auto px-3 py-6 border-b md:border-b-0 md:border-l border-gray-200 dark:border-gray-700/60 min-w-60 md:space-y-3">
      <!-- Group 1 -->
  <div class="flex flex-col gap-2 border-b border-gray-200 dark:border-gray-700/60 p-2">
    <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase">Setting</div>
    <!-- 첫 번째 줄 -->
    <div class="flex items-center gap-1">
      <label for="reference" class="text-xs font-bold flex items-center text-gray-700 dark:text-gray-300">
        Serial Number
      </label>
      <input class="h-9 w-48 p-2 text-xs border border-gray-300 rounded-md" type="text" disabled />
    </div>
    <div class="flex flex-col gap-1">
        <label for="filename" class="text-xs font-bold flex items-center text-gray-700 dark:text-gray-300">
        Setting file path <span class="text-red-500">*</span>
      </label>
      <div class="flex items-center gap-1">
        <input
        id="filename" name="file"
        class="h-9 w-48 p-2 text-xs border border-gray-300 rounded-md"
        @change="handleFileUpload"
        type="file"
        required
      />
      <button
        class="h-9 w-24 px-4 text-sm bg-gray-900 text-white rounded-md hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
        @click.prevent="upload"
      >
        Upload
      </button>
      </div>
      <div
        class="font-medium text-xs text-gray-800 dark:text-gray-100 mb-3"
      >
        {{ message }}🙌
      </div>
    </div>
    <!-- 두 번째 줄 -->
  </div>
    <div class="flex flex-col gap-2 border-b border-gray-200 dark:border-gray-700/60 p-2">
    <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase">Time & ErrorLimit</div>
    <div class="flex items-center gap-1">
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        DateTime
        </label>
        <input
        id="reference" 
        type="text"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-32"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        Error Limit
        </label>
        <input
        id="reference" 
        type="text"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-32"
        />
      </div>
    </div>
  </div>
  <div class="flex flex-col gap-2 border-b border-gray-200 dark:border-gray-700/60 p-2">
    <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase">Reference</div>
    <div class="flex items-center gap-1">
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        U ref
        </label>
        <input
        id="reference" v-model.number="refDict.U"
        type="text"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-24"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        I Ref
        </label>
        <input
        id="reference" v-model.number="refDict.I"
        type="text"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-24"
        />
      </div>

    </div>
    <div class="flex items-center gap-1">
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        In Ref
        </label>
        <input
        id="reference"
        type="text" v-model.number="refDict.In"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-24"
        />
      </div>
      <div class="flex flex-col gap-1">
        <label for="reference" class="text-xs font-bold text-gray-700 dark:text-gray-300">
        Phase Angle
        </label>
        <input
        id="reference"
        type="text" v-model.number="refDict.P"
        class="form-input h-8 px-3 rounded-md border border-gray-300 w-24"
        />
      </div>
    </div>
  </div>
  <div class="flex flex-col gap-2 border-b border-gray-200 dark:border-gray-700/60 p-2">
    <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase">Command</div>
    <div class="grid grid-cols-2 gap-2 w-full max-w-md">
      
      <button v-for="(row, index) in commands"
        class="self-end h-8 w-32 px-4 text-sm dark:hover:bg-white" :class="btnClass(index)" @click = sendCmd(index)
      > {{  row.Label }}
      </button>

    </div>
  </div>
    </div>
  </template>
  
  <script>
  import { ref , computed, onMounted} from 'vue'
  import axios from 'axios';
  import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
  export default {
    name: 'CaliSidebar',
    props:['commands'],
    setup(props){
      const setupStore = useSetupStore();
      const selectedFile = ref(null);
      const message = ref('');
      const setupList = ref([]);
      const commands = ref([]);
      const refDict = ref({
        U:0,
        I:0,
        In:0,
        P:0
      });
      onMounted(()=>{
        commands.value = props.commands;
      });

      const btnClass =(index)=>{
        if (index%2 == 0)
          return 'bg-green-900 text-white rounded-md hover:bg-green-800 dark:bg-green-100 dark:text-green-800';
        else
          return 'bg-pink-900 text-white rounded-md hover:bg-pink-800 dark:bg-pink-100 dark:text-pink-800';
      };

      const sendCmd = async(index) =>{
        //console.log(commands.value[index]["name"]);
        if (index == 0){
          try {
            const response = await axios.get("/config/calibrate/start");
            if (response.data.passOK == "1") {
              alert(commands.value[index]["name"] +' Success');
            } else {
              alert(response.data.error);
            }
          } catch (error) {
            alert(error);
          }
        }else if(index == 1){
          try {
            const response = await axios.get("/config/calibrate/end");
            if (response.data.passOK == "1") {
              alert(commands.value[index]["name"] +' Success');
            } else {
              alert(response.data.error);
            }
          } catch (error) {
            alert(error);
          }
        }else{
          let data = {};
          if("Param" in commands.value[index]){
            data["ref"] = refDict.value[commands.value[index]["Param"]].toString();
          }else{
            data["ref"] = 'None'
          }
          data["cmd"] = commands.value[index]["name"]
          //console.log(formData);
          try {
            const response = await axios.post(`/config/calibrate/cmd/${commands.value[index]["Type"]}`, data, {
              headers: { "Content-Type": 'application/json' },
            });
            if (response.data.passOK == "1") {
              alert(commands.value[index]["name"] +' Success');
            } else {
              alert(response.data.error);
            }
          } catch (error) {
            alert(error);
          }
        }
      }
      //const cali = computed(() => setupStore.getCalib);

      const handleFileUpload = (event) => {
          selectedFile.value = event.target.files[0]; // 파일 객체 저장
        };

      const upload = async() =>{
        if (!selectedFile.value) {
          message.value = "파일을 선택하세요!";
          return;
        }

        const formData = new FormData();
        formData.append("file", selectedFile.value);
        //console.log(formData);
        try {
          const response = await axios.post("/config/upload", formData, {
            headers: { "Content-Type": "multipart/form-data" },
          });
          if (response.data.passOK == "1") {
            message.value = "Upload Success : " + response.data.file_path;
            setupList.value = response.data.data;
            setupStore.setCalib(true);
          } else {
            message.value = response.data.error;
          }
        } catch (error) {
          message.value = "업로드 실패: " + error.response.data.error;
        }
      };

      return {
        handleFileUpload,
        upload,
        setupList,
        message,
        commands,
        btnClass,
        sendCmd,
        refDict,
      }
    }
  }
  </script>