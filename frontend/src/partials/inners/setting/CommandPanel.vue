<template>
    <div class="grow">
  
      <!-- Panel body -->
      <div class="p-6">
        <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-5"> {{ t('sidebar.system') }} </h2>
  
        <!-- General -->
        <div class="mb-6">
          <!-- Filters -->
          <div class="mt-4 flex gap-2 border-b border-gray-200 dark:border-gray-700/60 pb-4">
            <!--div class="flex items-center">
              <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold mr-4"> SV-500 Status</h3>
                <span v-if="health" class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-green-500/20 text-green-700 font-semibold">
                  Active
                </span>
                <span v-else class="text-sm rounded-full px-3 py-1 min-w-[100px] w-fit text-center whitespace-nowrap transition-all bg-red-500/20 text-red-700 font-semibold">
                  Disabled
                </span>
            </div-->
          </div>
          <div v-if="message"
            class="text-sm text-gray-800 dark:text-gray-100 mt-2"
          >
            {{ message }}🙌
          </div>
        </div>
  
        <!-- Connected Apps cards -->
        <section class="pb-6 border-b border-gray-200 dark:border-gray-700/60">
          <div class="grid grid-cols-12 gap-6">
            <CommandItem :item="'Clear'" :channel="'main'" @service-done="showMessage"/>
            <CommandItem :item="'Clear'" :channel="'sub'" @service-done="showMessage"/>
          </div>
        </section>

        <!-- Settings Section -->
        <section class="mt-6 pb-6 border-b border-gray-200 dark:border-gray-700/60">
          <div class="grid grid-cols-12 gap-6">
            <CommandItem :item="'Settings'" @service-done="showMessage"/>
            <!-- <CommandItem :item="'System'" @service-done="showMessage"/> -->
          </div>
        </section>

        <!-- Disk Status Section -->
        <section class="mt-6">

        </section>
      </div>
  
    </div>
</template>
  
<script>
import CommandItem from './CommandItem.vue'
import { onMounted, ref , watch, provide, computed} from 'vue'
import { useSetupStore } from '@/store/setup'
import { useAuthStore } from '@/store/auth'
import { useI18n } from "vue-i18n"; // 
import axios from 'axios';
export default {
  name: 'ServicePannel',
  components:{
      CommandItem,
  },
  setup(){
    const { t } = useI18n();
    const setupStore = useSetupStore();
    const authStore = useAuthStore();
    const message = ref('')
    const health = ref('');
    const sysStatus = ref({});
    const diskStatus = ref([]);
    const showMessage = (text) => {
      message.value = text
    }
    const ChannelState = computed(()=> {
      const re = setupStore.getChannelSetting;
      return re.MainDiagnosis || re.SubDiagnosis
    });

    const devMode = computed(()=> authStore.getOpMode);
    
    // Helper functions for disk status display
    const formatSize = (sizeInGB) => {
      if (sizeInGB >= 1024) {
        return (sizeInGB / 1024).toFixed(1) + ' TB';
      }
      return sizeInGB.toFixed(1) + ' GB';
    };

    const getUsagePercentage = (item) => {
      const used = item.totalGB - item.freeGB;
      return Math.round((used / item.totalGB) * 100);
    };

    const getProgressBarColor = (item) => {
      const usage = getUsagePercentage(item);
      if (usage >= 90) return 'bg-red-500';
      if (usage >= 75) return 'bg-yellow-500';
      return 'bg-green-500';
    };

  const CheckAPI = async () => {
    try {
      const response = await axios.get("/setting/checkAPI");
      if (response.data.success){
        health.value = response.data.data;
      }else{
        health.value = '';
      }
    } catch (error) {
      console.log(error);
      message.value = "Restful API Service is not running";
    }
  };

  const SysCheck = async () => {
    try {
      const response = await axios.get("/setting/SysCheck");
      if (response.data.success){
         diskStatus.value = response.data.disk;
         sysStatus.value = response.data.data;
      }else{
        message.value = "System Check API is not respond"
      }
    } catch (error) {
      message.value = "System Check Failed";
      console.log(error);
    }
  };

  onMounted(()=>{
    SysCheck();
  })

  provide('sysStatus',sysStatus);
  provide('health', health);

  function parseFullSystemStatus(rawStr) {
    const lines = rawStr.split('\n');
    const result = {};
    let currentKey = null;
    let buffer = '';

    for (let line of lines) {
      line = line.trim();
      if (line.endsWith(':')) {
        currentKey = line.slice(0, -1);
        buffer = '';
      } else if (line.includes(':') && !line.includes('{')) {
        const [key, value] = line.split(':').map(s => s.trim());
        result[key] = value;
        currentKey = null;
      } else if (currentKey) {
        buffer += line;
        if (currentKey === 'Disk Status') {
          result[currentKey] = parseCustomBlock(buffer);
        }
        else {
          try {
            result[currentKey] = JSON.parse(buffer);
          } catch (e) {
            result[currentKey] = buffer;
          }
        }
        currentKey = null;
      }
    }

    return result;
  }

  function parseCustomBlock(str) {
    const pattern = /(\w+)\s*=\s*([^,}]+)/g;
    const result = {};
    let match;
    while ((match = pattern.exec(str)) !== null) {
      const key = match[1];
      const val = match[2].trim();
      result[key] = isNaN(val) ? val : parseFloat(val);
    }
    return result;
  }

  watch(
    () => ChannelState.value,
    (newVal, oldVal) => {
      if (newVal) {
        CheckAPI();
      } else {
        console.log('둘 다 비활성화됨');
      }
    },
    { immediate: true }
  );

    return {
      message,
      showMessage,
      health,
      CheckAPI,
      sysStatus,
      diskStatus,
      ChannelState,
      devMode,
      formatSize,
      getUsagePercentage,
      getProgressBarColor,
      t,
    }

  }
}
</script>