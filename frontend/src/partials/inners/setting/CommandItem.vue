<template>
    <div class="col-span-full xl:col-span-6 2xl:col-span-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-lg">
              <!-- Card content -->
              <div class="flex flex-col h-full p-5">
                <div v-if="item == 'Clear'" class="grow">
                    <header class="flex items-center mb-4">
                        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">
                          {{ channel === 'main' ? t("config.system.title_1")
                          : channel === 'sub' ? t("config.system.title_2"): '초기화 명령어' }}  <!--//{{ t("config.system.title_1") }} Clear Command-->
                        </h3>
                    </header>
                    <div class="flex flex-col space-y-4">
                      <!-- 첫 번째 줄: 3개 버튼 -->
                      <div class="flex justify-start space-x-4">                        
                        <button
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('maxmin',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp;{{ t("config.system.minmax") }} <!--MaxMin-->
                        </button>
                        <button
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('alarm',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp;{{ t("config.system.alarm") }} <!--Alarm-->
                        </button>
                        
                        <button
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('event',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp;{{ t("config.system.eventCount") }} <!--Event Count-->
                        </button>
                      </div>
                      
                      <!-- 두 번째 줄: 3개 버튼 -->
                      <div class="flex justify-start space-x-4">
                        <button 
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('demand',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp; {{ t("config.system.demand") }} <!--Demand-->
                        </button>              
                        <button
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('energy',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp;{{ t("config.system.energy") }} <!--Energy-->
                        </button>
                        <button
                          class="btn w-40 h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                          @click="command('runhour',0)"
                        >
                          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                          </svg>
                          &nbsp;{{ t("config.system.runhour") }} <!--Run Hour-->
                        </button>
                      </div>
                    </div>
                </div>
                <div v-else-if="item == 'System'" class="grow">
                    <header class="flex items-start mb-4">
                        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">
                          {{ t("config.system.title_2") }} <!--Device--></h3>
                    </header>
                    <div class="flex justify-center space-x-4">
                        <button
                        class="btn h-11 px-5 bg-violet-900 text-xs text-white hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
                        @click="ResetAll"
                      ><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="20" height="20" stroke-width="1.25">
                        <path d="M19 14v-2h2l-9 -9l-9 9h2v7a2 2 0 0 0 2 2h2.5"></path>
                        <path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 1.75 1.032"></path>
                        <path d="M15.536 17.586a2.123 2.123 0 0 0 -2.929 0a1.951 1.951 0 0 0 0 2.828c.809 .781 2.12 .781 2.929 0c.809 -.781 -.805 .778 0 0l1.46 -1.41l1.46 -1.419"></path>
                        <path d="M15.54 17.582l1.46 1.42l1.46 1.41c.809 .78 -.805 -.779 0 0s2.12 .781 2.929 0a1.951 1.951 0 0 0 0 -2.828a2.123 2.123 0 0 0 -2.929 0"></path>
                      </svg>
                      &nbsp;{{ t("config.system.init") }} <!--Init setting-->
                      </button>
                      <button
                        class="btn h-11 px-5 bg-pink-900 text-xs text-white hover:bg-pink-800 dark:bg-pink-100 dark:text-pink-800 dark:hover:bg-white"
                        @click="command('reboot',1)"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="12" y1="6" x2="12" y2="12" />
                        <path d="M17 7.5a6.5 6.5 0 1 1-10 0" />
                      </svg>
                      &nbsp;{{ t("config.system.reboot") }} <!--Reboot-->
                      </button>

                    </div>
                </div>
                <div v-else-if="item == 'Settings'" class="grow">
                    <header class="flex items-center mb-4">
                        <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">
                          {{ t("config.system.title_3") }} <!--Settings--></h3>
                    </header>
                    <div class="flex justify-start space-x-4">
                        <button
                        class="btn h-11 px-5 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white"
                        @click.stop="feedbackModalOpen = true"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                        <polyline points="14,2 14,8 20,8"/>
                        <line x1="16" y1="13" x2="8" y2="13"/>
                        <line x1="16" y1="17" x2="8" y2="17"/>
                        <polyline points="10,9 9,9 8,9"/>
                      </svg>
                      &nbsp;{{ t("config.system.import") }} <!--Import-->
                      </button>
                      <button
                        class="btn h-11 px-5 bg-yellow-900 text-xs text-white hover:bg-yellow-800 dark:bg-yellow-100 dark:text-yellow-800 dark:hover:bg-white"
                        @click="download"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                        <polyline points="7,10 12,15 17,10"/>
                        <line x1="12" y1="15" x2="12" y2="3"/>
                      </svg>
                      &nbsp;{{t("config.system.export") }} <!--Export-->
                      </button>
                      <button
                        class="btn h-11 px-5 bg-violet-900 text-xs text-white hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
                        @click="restore"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                        <rect x="7" y="7" width="3" height="9"/>
                        <rect x="14" y="7" width="3" height="5"/>
                      </svg>
                      &nbsp;{{t("config.system.restore") }} <!--Restore-->
                      </button>
                      <button
                        class="btn h-11 px-5 bg-red-900 text-xs text-white hover:bg-red-800 dark:bg-red-100 dark:text-red-800 dark:hover:bg-white"
                        @click="confirmReset"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M3 6h18"/>
                        <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                        <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                        <line x1="10" y1="11" x2="10" y2="17"/>
                        <line x1="14" y1="11" x2="14" y2="17"/>
                      </svg>
                      &nbsp;{{t("config.system.reset") }} <!--Reset-->
                      </button>
                      <button
                        class="btn h-11 px-5 bg-pink-900 text-xs text-white hover:bg-pink-800 dark:bg-pink-100 dark:text-pink-800 dark:hover:bg-white"
                        @click="command('reboot',1)"
                      ><svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="12" y1="6" x2="12" y2="12" />
                        <path d="M17 7.5a6.5 6.5 0 1 1-10 0" />
                      </svg>
                      &nbsp;{{ t("config.system.reboot") }} <!--Reboot-->
                      </button>
                    </div>
                </div>
                <!-- Card footer -->
              </div>



              <!-- File Upload Modal -->
              <ModalBasic v-if="item == 'Settings'" id="feedback-modal" :modalOpen="feedbackModalOpen" @close-modal="feedbackModalOpen = false" title="File Upload">
                <!-- Modal content -->
                <div class="px-5 py-4">
                  <div class="text-sm">
                    <div
                      class="font-medium text-gray-800 dark:text-gray-100 mb-3"
                    >
                      Setting file upload
                    </div>
                  </div>
                  <div class="space-y-3">
                    <div>
                      <label class="block text-sm font-medium mb-1" for="name"
                        >file path <span class="text-red-500">*</span></label
                      >
                      <input
                        id="filename"
                        class="form-input w-full px-2 py-1"
                        @change="handleFileUpload"
                        type="file"
                        required
                      />
                    </div>
                  </div>
                  <div class="text-sm">
                    <div
                      class="font-medium text-gray-800 dark:text-gray-100 mb-3"
                    >
                      {{ message }}🙌
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
                      @click.prevent="upload"
                    >
                      Import
                    </button>
                    <button
                      class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
                      @click.stop="feedbackModalOpen = false"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </ModalBasic>
            </div>
</template>
<script>
import { useRouter } from 'vue-router'; 
import { ref, watch, computed, onMounted, inject } from 'vue'
import axios from 'axios';
import { useAuthStore } from "@/store/auth"; // ✅ Pinia Store 사용
import ModalBasic from "../../../pages/common/ModalBasic.vue";
import { useI18n } from "vue-i18n"; // 

export default {
  name: 'ServiceCard',
  components: {
    ModalBasic,
  },
  props: ['item', 'mode', 'state','channel'],
  emits: ['service-done'],
  setup(props, { emit }){
    const { t } = useI18n();
    const authStore = useAuthStore();
    const router = useRouter();
    const installed = computed(() => authStore.getInstalled);
    const DiagState = ref(props.state);
    const feedbackModalOpen = ref(false);
    // const dangerModalOpen = fef(false);
    const selectedFile = ref(null);

      const downloadUrl = computed(() => {
        return `http://192.168.1.24:5000/api/getFolder?name=${modalSelectItem.value}`
      })

    const item = ref(props.item);
    const Status = inject('sysStatus');
    const apiHealth = inject('health');

    const itemDict = ref({
      "DBMS":"Database Service",
      "App":"Smart Systems Service",
      "System":"System",
      "Backup Download":"Download"
    });
    const mode = ref(props.mode);
    const modalSelectItem = ref('all');
    const message = ref('');
    const health = ref(false);
    const healthAPI = ref(false);
    const coreHealth = ref(true);
    const dbInit = ref(false);
    
    watch(() => props.item, (newChannel) => {
        item.value = newChannel;
    }, { immediate: true });

    watch([Status, apiHealth], () => {
      switch(item.value){
        case "DBMS":
          health.value = Status?.value?.["InfluxDB"];
          break;
        case "App":
          health.value = Status?.value?.["smartsystemsservice"];
          break;
        case "System":
          healthAPI.value = ['ok', 'OK'].includes(apiHealth.value);
          coreHealth.value = Status?.value?.["core"];
          break;
      }
    }, { immediate: true });

    const serviceOp = async (cmd) => {
      try {
        const response = await axios.get(`/setting/SysService/${cmd}/${item.value}`);

        if (response.data.success) {
          message.value = `${item.value} ${cmd}!`;
        } else {
          message.value = `${item.value} ${cmd} Failed: Restful API Service error`;
        }
      } catch (error) {
        message.value =
          `${item.value} ${cmd} Failed: `;
      }
      emit('service-done', message.value);
    };
    const command = async (cmdItem, command) => {
      try {
        const typemap = { 'demand':0 , 'maxmin':1, 'energy':2, 'alarm':3, 'event':4, 'reboot':5, 'runhour':6 };
        let ch = -1;
        if (typemap[cmdItem] == 5 && command == 1){
          ch = 2;
        }else{
          ch = props.channel == 'Main'? 0: 1;
        }
        //console.log(ch);

        const data ={
          type : ch,
          cmd : command,
          item: typemap[cmdItem]
        }
        // const response = await axios.post(`/setting/command`);
        const response = await axios.post("/setting/command", data, {
        headers: { "Content-Type": "application/json" },
        withCredentials: true,
      });

        if (response.data.success) {
          if (command == 0)
            message.value = `ITEM:${cmdItem} COMMAND:${command} CHANNEL:${props.channel}!`;
          else
            message.value = `ITEM:${cmdItem} COMMAND:${command} !`;
        } else {
          message.value = `ITEM:${cmdItem} COMMAND:${command} CHANNEL:${props.channel} Failed: ${response.data.message}`;
        }
      } catch (error) {
        console.error(error);
        message.value =
          `ITEM:${cmdItem} COMMAND:${command} CHANNEL:${props.channel} Failed: ${error}`;
      }
      emit('service-done', message.value);
    };
    const checkDBMS = async () => {
      try {
        const response = await axios.get(`/auth/checkDBMS`);
        if (response.data.result == 1) {
          dbInit.value = true;
        }else{
          dbInit.value = false;
        } 
      } catch (error) {
        console.error(error);
      }
    };

    const init = async () => {
      try {
        const response = await axios.get("/setting/initDB");
        if (response.data.success){
          message.value = "InfluxDB initiated"
          authStore.setInstall(3);
        }else{
          message.value = "InfluxDB initialization Failed"
        }
      } catch (error) {
        message.value = "InfluxDB initialization Failed";
      }

      emit('service-done', message.value);
    };

    const ResetAll = async () => {
      try {
        const response = await axios.get("/setting/ResetAll");
        if (response.data.success){
          message.value = "Setup initiated"
          authStore.setInstall(0);

          if (authStore.logout) {
            await authStore.logout();
            router.push("/signin"); 
          } else {
            console.error("❌ Error: logout 함수가 정의되지 않음!");
          }
          //dangerModalOpen.value = false;
        }else{
          console.error("데이터 가져오기 실패:", error);
        }
      } catch (error) {
        console.error("데이터 가져오기 실패:", error);
      }
    };

    const Reboot = async () => {
      try {
        const response = await axios.get("/setting/Reboot");
        if (response.data.success){
          message.value = "System Reboot initiated"
        }else{
          message.value = "Reboot Failed"
        }
      } catch (error) {
        message.value = "Reboot Failed";
      }
      emit('service-done', message.value);
    };

    // ✅ confirm 확인창을 사용한 리셋 함수
    const confirmReset = () => {
      const confirmed = confirm('All settings will be deleted and reset to default.\nDo you want to proceed?');
      if (confirmed) {
        Reset();
      }
    };

    const Reset = async () => {
      //console.log('Reset function called!'); // 디버깅용
      try {
        const response = await axios.get("/setting/Reset");
        //console.log('Reset API response:', response.data); // 디버깅용
        
        if (response.data.success){
          message.value = "Setup initiated"
          authStore.setInstall(0);

          if (authStore.logout) {
            await authStore.logout();
            router.push("/signin"); 
          } else {
            console.error("❌ Error: logout 함수가 정의되지 않음!");
          }
        }else{
          console.error("데이터 가져오기 실패:", response.data);
        }
      } catch (error) {
        console.error("Reset 실패:", error);
      }
      emit('service-done', message.value);
    };

    // Settings 카드 관련 함수들
    const handleFileUpload = (event) => {
      selectedFile.value = event.target.files[0];
    };

    const upload = async () => {
      if (!selectedFile.value) {
        message.value = "파일을 선택하세요!";
        return;
      }

      const formData = new FormData();
      formData.append("file", selectedFile.value);
      try {
        const response = await axios.post("/setting/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        if (response.data.passOK == "1") {
          message.value = "Upload Success : " + response.data.file_path;
          feedbackModalOpen.value = false;
          // await GetSettingData(); // 필요시 추가
        } else {
          message.value = response.data.error;
        }
      } catch (error) {
        message.value = "업로드 실패: " + error.response.data.error;
      }
      emit('service-done', message.value);
    };

    const download = async () => {
      try {
        const response = await axios.get("/setting/download", {
          responseType: "blob",
        });

        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", "setting.json");
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        message.value = "Download Success!";
      } catch (error) {
        message.value = "다운로드 실패: " + error.response.data.error;
      }
      emit('service-done', message.value);
    };

    const restore = async () => {
      try {
        const response = await axios.post("/setting/restoreSetting");

        if (response.data.passOK === "1") {
          message.value = "Restore Success!";
          // await GetSettingData(); // 필요시 추가
        } else {
          message.value = response.data.error || "Unknown error";
        }
      } catch (error) {
        message.value =
          "복원 실패: " + error.response?.data?.error || error.message;
      }
      emit('service-done', message.value);
    };

    onMounted(()=>{
      if(item.value == 'DBMS')
        checkDBMS();
    });

    return {
      item,
      itemDict,
      mode,
      modalSelectItem,     
      stop,
      downloadUrl,
      message,
      installed,
      init,
      health,
      healthAPI,
      DiagState,
      coreHealth,
      checkDBMS,
      dbInit,
      ResetAll,
      Reboot,
      Reset,
      confirmReset, // ✅ 추가
      feedbackModalOpen,
      selectedFile,
      handleFileUpload,
      upload,
      download,
      restore,
      command,
      // dangerModalOpen,
      t,
    }
  }
}
</script>