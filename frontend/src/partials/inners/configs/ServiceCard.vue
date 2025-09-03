<template>
    <div class="col-span-full xl:col-span-6 2xl:col-span-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-lg">
              <!-- Card content -->
              <div class="flex flex-col h-full p-5">
                <div class="grow">
                  <header class="flex items-center mb-4">
                    <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4">{{ itemDict[item] }}</h3>
                    
                    <!-- Redis, InfluxDB,'Core', 'WebServer' 상태 -->
                    <template v-if="['Redis', 'InfluxDB','Core', 'WebServer' ].includes(item)">
                      <span v-if="health" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                         Running
                      </span>
                      <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                        Stopped
                      </span>
                    </template>

                    <!--  SmartSystems, SmartAPI상태 (App 카테고리) -->
                    <template v-if="['SmartSystems', 'SmartAPI'].includes(item)">
                      <template v-if="DiagState">
                        <span v-if="health" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                           Running
                        </span>
                        <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                          Stopped
                        </span>
                      </template>
                      <template v-else>
                        <span class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-gray-500/20 text-gray-700 font-semibold">
                          Disabled
                        </span>
                      </template>
                    </template>

                    <!-- 기존 DBMS 상태 -->
                    <template v-if="item == 'DBMS'">
                      <span v-if="health" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                         Running
                      </span>
                      <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                        Stopped
                      </span>
                    </template>

                    <!-- 기존 App 상태 -->
                    <template v-if="item == 'App'">
                      <template v-if="DiagState">
                        <span v-if="health" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                           Running
                        </span>
                        <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                          Stopped
                        </span>
                      </template>
                      <template v-else>
                        <span class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-gray-500/20 text-gray-700 font-semibold">
                          Disabled
                        </span>
                      </template>
                    </template>

                    <!-- 기존 System 상태 -->
                    <template v-if="item == 'System'">
                      <template v-if="DiagState">
                        <span v-if="healthAPI" class="mr-2 text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                          API  Running
                        </span>
                        <span v-else class="mr-2 text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                          API Stopped
                        </span>
                      </template>
                      <span v-if="coreHealth" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                          Core  Running
                      </span>
                        <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                          Core Stopped
                      </span>
                    </template>
                  </header>
                  <div class="flex flex-col">
                    <div v-if="mode =='Service'" class="flex justify-center space-x-4">
                      <button v-if="!dbInit && (item =='DBMS' || item == 'InfluxDB')"
                        class="btn h-9 px-5 bg-violet-900 text-violet-100 hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
                        @click="init"
                      ><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="24" height="24" stroke-width="1.25">
                        <path d="M19 14v-2h2l-9 -9l-9 9h2v7a2 2 0 0 0 2 2h2.5"></path>
                        <path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 1.75 1.032"></path>
                        <path d="M15.536 17.586a2.123 2.123 0 0 0 -2.929 0a1.951 1.951 0 0 0 0 2.828c.809 .781 2.12 .781 2.929 0c.809 -.781 -.805 .778 0 0l1.46 -1.41l1.46 -1.419"></path>
                        <path d="M15.54 17.582l1.46 1.42l1.46 1.41c.809 .78 -.805 -.779 0 0s2.12 .781 2.929 0a1.951 1.951 0 0 0 0 -2.828a2.123 2.123 0 0 0 -2.929 0"></path>
                      </svg>
                      &nbsp;Init
                      </button>
                      <button
                        class="btn h-9 px-5 bg-green-900 text-green-100 hover:bg-green-800 dark:bg-green-100 dark:text-green-800 dark:hover:bg-white"
                        @click="start"
                      ><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="24" height="24" stroke-width="1.25">
                        <path d="M7 4v16l13 -8z"></path>
                      </svg>
                      &nbsp;Start
                      </button>
                      <button
                        class="btn h-9 px-5 bg-yellow-900 text-yellow-100 hover:bg-yellow-800 dark:bg-yellow-100 dark:text-yellow-800 dark:hover:bg-white"
                        @click="restart"
                      >
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="24" height="24" stroke-width="1.25">
                        <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                        <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
                      </svg>
                        &nbsp;Restart
                      </button>
                      <button
                        class="btn h-9 px-5 bg-red-900 text-red-100 hover:bg-red-800 dark:bg-red-100 dark:text-red-800 dark:hover:bg-white"
                        @click="stop"
                      >
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="24" height="24" stroke-width="1.25">
                        <path d="M5 5m0 2a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v10a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2z"></path>
                      </svg>
                      &nbsp;Stop
                      </button>
                      <button v-if="item =='System'"
                        class="btn h-9 px-5 bg-violet-900 text-violet-100 hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white"
                        @click="ResetAll"
                      ><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" width="24" height="24" stroke-width="1.25">
                        <path d="M19 14v-2h2l-9 -9l-9 9h2v7a2 2 0 0 0 2 2h2.5"></path>
                        <path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 1.75 1.032"></path>
                        <path d="M15.536 17.586a2.123 2.123 0 0 0 -2.929 0a1.951 1.951 0 0 0 0 2.828c.809 .781 2.12 .781 2.929 0c.809 -.781 -.805 .778 0 0l1.46 -1.41l1.46 -1.419"></path>
                        <path d="M15.54 17.582l1.46 1.42l1.46 1.41c.809 .78 -.805 -.779 0 0s2.12 .781 2.929 0a1.951 1.951 0 0 0 0 -2.828a2.123 2.123 0 0 0 -2.929 0"></path>
                      </svg>
                      &nbsp; Factory Default
                      </button>
                    </div>
                    <div v-else class="flex justify-center space-x-4">
                      <label
                        for="reference"
                        class="text-sm flex items-center text-gray-700 dark:text-gray-300"
                      >
                        Select Item
                      </label>

                      <!-- Item 선택 -->
                      <select
                        class="h-9 w-32 p-2 text-sm border border-gray-300 rounded-md"
                        v-model="modalSelectItem"
                      >
                        <option value="all">All</option>
                        <option value="log">Log</option>
                        <option value="project">Project</option>
                        <option value="dbbackup">DBbackup</option>
                        <option value="backup">Backup</option>
                      </select>
                      <a
                        :href="downloadUrl"
                        class="btn h-9 px-5 bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
                        download
                      >
                        Download
                      </a>
                    </div>
                  </div>
                </div>
                <!-- Card footer -->
              </div>
            </div>
</template>
<script>
import { useRouter } from 'vue-router'; 
import { ref, watch, computed, onMounted, inject } from 'vue'
import axios from 'axios';
import { useAuthStore } from "@/store/auth";

export default {
  name: 'ServiceCard',
  props: ['item', 'mode', 'state'],
  emits: ['service-done'],
  setup(props, { emit }){
    const authStore = useAuthStore();
    const router = useRouter();
    const installed = computed(() => authStore.getInstalled);
    const DiagState = ref(props.state);


    const item = ref(props.item);
    const Status = inject('sysStatus');
    const apiHealth = inject('health');

    const itemDict = ref({
      "Redis":"Redis",
      "InfluxDB":"Influxdb", 
      "SmartSystems":"SmartSystem",
      "SmartAPI":"Smart RestAPI",
      "Core":"Core",
      "WebServer":"WebServer",
    });

    const mode = ref(props.mode);

    const message = ref('');
    const health = ref(false);
    const healthAPI = ref(false);
    const coreHealth = ref(true);
    const dbInit = ref(false);

    watch(() => props.item, (newChannel) => {
        item.value = newChannel;
    }, { immediate: true });

    watch([Status, apiHealth], () => {
      console.log('Status:', Status.value, 'Item:', item.value);
      
      switch(item.value){
        case "Redis":
          health.value = Status?.value?.["redis"];
          break;
        case "InfluxDB":
          health.value = Status?.value?.["influxdb"];
          break;
        case "SmartSystems":
          health.value = Status?.value?.["smartsystem"];
          break;
        case "SmartAPI":
          health.value = Status?.value?.["smartapi"];
          break;
        case "Core":
          health.value = Status?.value?.["core"];
          break;
        case "WebServer":
          health.value = Status?.value?.["webserver"];
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
        message.value = `${item.value} ${cmd} Failed: `;
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

    const start = async () => {
      serviceOp('start');
    };

    const restart = async () => {
      serviceOp('restart');
    };

    const stop = async () => {
      serviceOp('stop');
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



    onMounted(()=>{
      if(item.value == 'DBMS' || item.value == 'InfluxDB')
        checkDBMS();
    });

    return {
      item,
      itemDict,
      mode,
      start,
      stop,
      restart,
      message,
      installed,
      init,
      health,
      healthAPI,
      DiagState,
      coreHealth,
      checkDBMS,
      dbInit,
    }
  }
}
</script>