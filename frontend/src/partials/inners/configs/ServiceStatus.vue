<template>
    <div class="col-span-full xl:col-span-6 2xl:col-span-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-lg">
              <!-- Card content -->
              <div class="flex flex-col h-full p-5">
                <div v-if="mode == 'System'" class="grow">
                  <header class="flex items-center mb-4">
                    <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">{{ title }}</h3>
                  </header>
                  <div class="flex items-center py-2 ml-2">
                      <template v-for="item in data">
                        <div  class="mr-5">
                          <div class="flex items-center">
                              <div class="text-lg font-bold text-gray-800 mr-2 dark:text-gray-100" > OK </div>
                          </div>
                          <div class="text-sm text-gray-500 dark:text-gray-400">{{ item.name }}</div>
                        </div>
                        <div class="hidden md:block w-px h-8 bg-gray-200 dark:bg-gray-700 mr-5" aria-hidden="true"></div>
                      </template>
                  </div>
                </div>
                <div v-if="mode == 'Disk' && DiagState" class="grow">
                     <header class="flex items-center mb-4">
                      <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">{{ title }}</h3>
                    </header>
                    <div class="flex items-center py-2">
                        <div class="mr-5 ml-2">
                            <div class="flex items-center">
                                <div class="text-lg font-bold text-gray-800 mr-2 dark:text-gray-100" > {{ data["totalGB"] }} </div>
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">Total GB</div>
                        </div>
                        <div class="hidden md:block w-px h-8 bg-gray-200 dark:bg-gray-700 mr-5" aria-hidden="true"></div>
                        <div class="mr-5 ml-2">
                            <div class="flex items-center">
                                <div class="text-lg font-bold text-gray-800 mr-2 dark:text-gray-100" > {{ data["freeGB"] }} </div>
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">Free GB</div>
                        </div>
                        <div class="hidden md:block w-px h-8 bg-gray-200 dark:bg-gray-700 mr-5" aria-hidden="true"></div>
                        <div class="mr-5 ml-2">
                            <div class="flex items-center">
                                <div class="text-lg font-bold text-gray-800 mr-2 dark:text-gray-100" > {{ data["status"] == 'ok'?'OK':'ERROR' }} </div>
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">Status</div>
                        </div>
                    </div>
                </div>
                <!-- Card footer -->
              </div>
            </div>
</template>
<script>
import { ref, watch, computed, onMounted, inject } from 'vue'
import { useSetupStore } from "@/store/setup"; // ✅ Pinia Store 사용
import axios from 'axios';
export default {
  name: 'ServiceCard',
  props: {
    title:String,
    mode:String,
    state:Boolean,
    data:Object
  },
  setup(props){

    const title = ref(props.title);
    const mode = ref(props.mode);
    const DiagState = ref(props.state);
    const data = ref(props.data);
    const setupStore = useSetupStore();
    // const item = ref(props.item);
    // const Status = inject('sysStatus');
    // const apiHealth = inject('health');

    // const mode = ref(props.mode);
    // const modalSelectItem = ref('all');
    // const message = ref('');
    // const health = ref(false);
    // const healthAPI = ref(false);

    // watch(() => props.item, (newChannel) => {
    //     item.value = newChannel;
    // }, { immediate: true });

    // watch([Status, apiHealth], () => {
    //   if (item.value === "DBMS" || item.value === "App") {
    //     health.value = Status?.value?.[itemDict.value[item.value]] !== 'error';
    //   } else if (item.value === "System") {
    //     healthAPI.value = ['ok', 'OK'].includes(apiHealth.value);
    //   }
    // }, { immediate: true });


    // const serviceOp = async (cmd) => {
    //   try {
    //     const response = await axios.get(`/setting/SysService/${cmd}/${item.value}`);

    //     if (response.data.success) {
    //       message.value = `${item.value} ${cmd}!`;
    //     } else {
    //       message.value = `${item.value} ${cmd} Failed: Restful API Service error`;
    //     }
    //   } catch (error) {
    //     message.value =
    //       `${item.value} ${cmd} Failed: `;
    //   }
    //   emit('service-done', message.value);
    // };

    // const start = async () => {
    //   serviceOp('start');
    // };

    // const restart = async () => {
    //   serviceOp('restart');
    // };

    // const stop = async () => {
    //   serviceOp('stop');
    // };

    // const init = async () => {
    //   try {
    //     const response = await axios.get("/setting/initDB");
    //     if (response.data.success){
    //       message.value = "InfluxDB initiated"
    //       authStore.setInstall(3);
    //     }else{
    //       message.value = "InfluxDB initialization Failed"
    //     }
    //   } catch (error) {
    //     message.value = "InfluxDB initialization Failed";
    //   }

    //   emit('service-done', message.value);
    // };

    // const Reset = async () => {
    //   try {
    //     const response = await axios.get("/setting/Reset");
    //     if (response.data.success){
    //       message.value = "Setup initiated"
    //       authStore.setInstall(0);
    //     }else{
    //       message.value = "Setup initialization Failed"
    //     }
    //   } catch (error) {
    //     message.value = "Setup initialization Failed";
    //   }

    //   emit('service-done', message.value);
    // };


    // onMounted(()=>{
    //   healthCheck();
    // });

    // const donwloadBackup = async () => {
    //   try {
    //     const response = await axios.get(`/setting/DownloadBackup/${modalSelectItem.value}`, {
    //       responseType: 'blob' // blob 타입으로 받아야 다운로드 가능
    //     });

    //     const messageHeader = response.headers['x-message'] || response.headers['X-Message'];

    //     if (messageHeader === 'OK') {
    //       const blob = new Blob([response.data], { type: response.headers['content-type'] });
    //       const url = window.URL.createObjectURL(blob);
    //       const link = document.createElement('a');
    //       link.href = url;
    //       link.setAttribute('download', `${modalSelectItem.value}.zip`);
    //       document.body.appendChild(link);
    //       link.click();
    //       link.remove();
    //       window.URL.revokeObjectURL(url);
    //       message.value = `${modalSelectItem.value} file downloaded!`;
    //     } else {
    //       // 서버에서 파일이 아닌 JSON 에러가 왔을 수도 있음
    //       const reader = new FileReader();
    //       reader.onload = () => {
    //         try {
    //           const json = JSON.parse(reader.result);
    //           message.value = `Download Failed: ${json.message || 'Unknown error'}`;
    //         } catch {
    //           message.value = "Download Failed: Invalid response format";
    //         }
    //         emit('service-done', message.value);
    //       };
    //       reader.readAsText(response.data);
    //       return; // emit 여기서 이미 호출됨
    //     }
    //   } catch (error) {
    //     message.value = "Download Failed";
    //   }

    //   emit('service-done', message.value);
    // };

    return {
      title,
      mode,
      data,
    //   modalSelectItem,
    //   start,
    //   stop,
    //   restart,
    //   downloadUrl,
    //   //donwloadBackup,
    //   message,
    //   installed,
    //   init,
    //   health,
    //   healthAPI,
    //   Reset,
      DiagState,
    }
  }
}
</script>