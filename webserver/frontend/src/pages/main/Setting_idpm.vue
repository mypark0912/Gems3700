<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <!-- Content area -->
    <div
      class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden"
    >
      <!-- Site header -->
      <Header
        :sidebarOpen="sidebarOpen"
        @toggle-sidebar="sidebarOpen = !sidebarOpen"
      />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <div class="mb-4 flex items-center gap-3">
            <h2
              class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold"
            >
              {{ t("config.sitemap.title") }}
            </h2>
            <button
              class="btn h-7 bg-gray-500 text-white hover:bg-gray-600 dark:bg-gray-600 dark:hover:bg-gray-500 px-4 py-1 text-xs rounded-full shadow-sm flex items-center gap-1.5"
              @click="systemModalOpen = true"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
                <line x1="8" y1="21" x2="16" y2="21" />
                <line x1="12" y1="17" x2="12" y2="21" />
              </svg>
              System
            </button>
          </div>

          <!-- Tab Bar (sticky) -->
          <div
            class="sticky top-16 z-20 bg-white dark:bg-gray-800 rounded-t-xl border border-b-0 border-gray-200 dark:border-gray-700/60 shadow-sm"
          >
            <div class="flex items-center justify-between px-4 py-2">
              <!-- Tabs -->
              <div class="flex items-center">
                <button
                  v-for="tab in tabs"
                  :key="tab.key"
                  class="px-5 py-3 text-sm font-medium border-b-2 transition-colors"
                  :class="mode === tab.key
                    ? 'border-violet-500 text-violet-600 dark:text-violet-400'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'"
                  @click="mode = tab.key"
                >
                  {{ tab.label }}
                </button>
              </div>
              <!-- Action Buttons -->
              <div class="flex items-center space-x-3">
                <!-- Setting Mode Switch -->
                <div class="flex items-center gap-2">
                  <span class="text-sm font-medium" :class="!isSetupMode ? 'text-green-600 dark:text-green-400' : 'text-gray-400 dark:text-gray-500'">
                    Local
                  </span>
                  <label class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" class="sr-only peer" :checked="!isSetupMode" @click.prevent="toggleSetupMode" />
                    <div class="w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-600 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-green-500"></div>
                  </label>
                </div>
                <button
                  class="btn h-9 dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-gray-300 px-4 py-1.5 text-sm rounded-lg"
                  :class="{ 'opacity-50 cursor-not-allowed': !isSetupMode }"
                  :disabled="!isSetupMode"
                  @click="save"
                >
                  Save
                </button>
                <button
                  class="btn h-9 bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white px-4 py-1.5 text-sm rounded-lg shadow-sm"
                  :class="{ 'opacity-50 cursor-not-allowed': !isSetupMode }"
                  :disabled="!isSetupMode"
                  @click="apply"
                >
                  Apply
                </button>
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="bg-white dark:bg-gray-800 shadow-sm rounded-b-xl border border-t-0 border-gray-200 dark:border-gray-700/60">
            <GeneralSetting v-if="mode === 'general'" />
            <MainSetting v-else-if="mode === 'main'" channel="Main" />
            <BranchSetting v-else-if="mode === 'ibsm'" :isSetupMode="isSetupMode" />
            <IpsmModuleSetting v-else-if="mode === 'ipsm'" :setupDict="setupDict" />
            <McsModuleSetting v-else-if="mode === 'mcs'" :isSetupMode="isSetupMode" />
          </div>
        </div>
      </main>

      <Footer />
    </div>

    <!-- System Modal -->
    <SystemModal :open="systemModalOpen" :isSetupMode="isSetupMode" @close="systemModalOpen = false" />

    <!-- Validation Result Modal -->
    <ValidationResultModal
      :modalOpen="validationModalOpen"
      :inputDict="inputDict"
      :channelMain="channel_main"
      @close="validationModalOpen = false"
      @save="handleSave"
    />
  </div>
</template>

<script setup>
import { ref, computed, provide, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useAuthStore } from '@/store/auth'
import { useSetupStore } from '@/store/setup'
import { useInputDict } from '@/composables/useInputDict'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'

import GeneralSetting from '../../partials/inners/setting/GeneralSetting.vue'
import MainSetting from '../../partials/inners/setting/Panel_Main.vue'
import BranchSetting from '../../partials/inners/setting/BranchSetting.vue'
import IpsmModuleSetting from '../../partials/inners/setting/IpsmModuleSetting.vue'
import McsModuleSetting from '../../partials/inners/setting/mcsd_setting.vue'
import SystemModal from '../../partials/inners/setting/SystemModal.vue'
import ValidationResultModal from '../../components/ValidationResultModal.vue'

const { t } = useI18n()

const props = defineProps({
  mode: {
    type: String,
    default: 'general'
  }
})

const sidebarOpen = ref(false)
const mode = ref(props.mode)
const systemModalOpen = ref(false)
const validationModalOpen = ref(false)

const authStore = useAuthStore()
const setupStore = useSetupStore()
const { inputDict, channel_main, channel_sub, setupDict } = useInputDict()
const devMode = computed(() => authStore.getOpMode)
const isSetupMode = ref(false)
const devLang = ref('')

const tabs = computed(() => {
  const list = [
    { key: 'general', label: 'General' },
    { key: 'main', label: 'Main' },
  ]
  if (devMode.value === 'device1') {
    // TODO: opMode에 따라 탭 하나만 렌더링하도록 변경
    // if (opMode === 'ibsm') {
    //   list.push({ key: 'ibsm', label: 'IBSM Module' })
    // } else if (opMode === 'ipsm') {
    //   list.push({ key: 'ipsm', label: 'IPSM Module' })
    // } else if (opMode === 'mcs') {
    //   list.push({ key: 'mcs', label: 'MCS Module' })
    // }
    list.push({ key: 'ibsm', label: 'IBSM Module' })
    list.push({ key: 'ipsm', label: 'IPSM Module' })
    list.push({ key: 'mcs', label: 'MCS Module' })
  }
  return list
})

      const GetSettingData = async () => {
        try {
          const response = await axios.get(`/setting/getSetting`);
          if (response.data.passOK == 1) {
            setupDict.value = response.data.data;
            Object.assign(inputDict.value, setupDict.value["General"]);
            Object.assign(channel_main.value, setupDict.value["main"]);
            // Object.assign(channel_sub.value, setupDict.value["sub"]);
            devLang.value = setupDict.value["lang"];
            inputDict.value.MQTT.device_id = inputDict.value.deviceInfo.mac_address;
            setupStore.setsetupFromFile(true);
            isSetupMode.value = setupDict.value.setupMode === 1
            if (setupDict.value.mode === 'device1') {
              await setupStore.fetchIbsm()
              if (response.data.data.ibsm) {
                setupDict.value.ibsm = response.data.data.ibsm
              }
            }
            if (!setupDict.value.ipsm72) {
              setupDict.value.ipsm72 = {
                channel: "ipsm72", Enable: 1, enable1: 0, enable2: 0,
                di60_1: 0, di60_2: 0, di60Target1: 0, di60Target2: 0,
                ipsm72_1: Array.from({ length: 72 }, () => ({ mod: null, pt: null })),
                ipsm72_2: Array.from({ length: 72 }, () => ({ mod: null, pt: null })),
              }
            }
            if (!setupDict.value.mcs) {
              setupDict.value.mcs = {
                channel: "mcs", Enable: 1,
                feeders: Array.from({ length: 32 }, () => ({
                  serialNumber: null, touName: '', systemType: 0,
                  modbusId: 0, ct1: 1, ct2: 1, ct3: 1,
                })),
              }
            }
          }
        } catch (error) {
          console.error("데이터 가져오기 실패:", error);
        }
      };

provide('inputDict', inputDict)
provide('channel_main', channel_main)
provide('channel_sub', channel_sub)
provide('setupDict', setupDict)
provide('devMode', devMode)
provide('GetSettingData', GetSettingData)

onMounted(() => {
  GetSettingData()
})

const toggleSetupMode = async () => {
  if (isSetupMode.value) {
    // ON → OFF: setupRemote=0
    try {
      await axios.get('/setting/releaseSetupMode')
      isSetupMode.value = false
    } catch (error) {
      console.error('releaseSetupMode 실패:', error)
    }
  } else {
    // OFF → ON: setupLocal 체크
    try {
      const resp = await axios.get('/setting/checkSetupMode')
      if (resp.data.status) {
        isSetupMode.value = true
      } else {
        alert('Local setup mode is active. Cannot switch to Remote.')
      }
    } catch (error) {
      console.error('checkSetupMode 실패:', error)
    }
  }
}

const transNumber = (channelData) => {
  const data = { ...channelData }
  ;["Enable"].forEach((field) => {
    if (data.hasOwnProperty(field)) data[field] = data[field] === true || data[field] === 1 ? 1 : 0
  })
  if (data.ctInfo) {
    for (const key in data.ctInfo) {
      if (key === "direction") data.ctInfo[key] = data.ctInfo.direction.map((d) => parseInt(d))
      else if (key === "inorminal") data.ctInfo[key] = parseFloat(data.ctInfo[key])
      else data.ctInfo[key] = parseInt(data.ctInfo[key])
    }
  }
  if (data.trendInfo?.params) {
    data.trendInfo.params = data.trendInfo.params.filter(p => p !== "None")
  }
  if (data.ptInfo) for (const key in data.ptInfo) data.ptInfo[key] = parseInt(data.ptInfo[key])
  if (data.eventInfo) for (const key in data.eventInfo) data.eventInfo[key] = parseInt(data.eventInfo[key])
  return data
}

const transFormat = () => {
  const generalData = { ...inputDict.value }
  if (generalData.useFunction) {
    generalData.useFunction.sntp = generalData.useFunction.sntp === true || generalData.useFunction.sntp === 1 ? 1 : 0
    generalData.useFunction.modbus_serial = Number(generalData.useFunction.modbus_serial) || 0
    generalData.useFunction.mqtt = generalData.useFunction.mqtt === true || generalData.useFunction.mqtt === 1 ? 1 : 0
  }
  if (generalData.tcpip) {
    generalData.tcpip.dhcp = generalData.tcpip.dhcp === true || generalData.tcpip.dhcp === 1 ? 1 : 0
  }

  const mainData = transNumber(channel_main.value)
  mainData.channel = "Main"

  const channels = [mainData]
  if (setupDict.value.ibsm) {
    channels.push(setupDict.value.ibsm)
  }
  if (setupDict.value.ipsm72) {
    channels.push(setupDict.value.ipsm72)
  }
  if (setupDict.value.mcs) {
    channels.push(setupDict.value.mcs)
  }
  return { mode: devMode.value, lang: setupDict.value?.lang || '', General: generalData, channel: channels }
}

const save = () => {
  validationModalOpen.value = true
}

const handleSave = async ({ done }) => {
  try {
    const formattedData = transFormat()
    const response = await axios.post("/setting/savefileNew", formattedData, {
      headers: { "Content-Type": "application/json;charset=utf-8" },
      withCredentials: true,
    })

    if (response.data?.status === "1") {
      alert("Configuration saved successfully!")
      validationModalOpen.value = false
    } else {
      alert("Save failed: " + (response.data?.error || "Unknown error"))
    }
  } catch (error) {
    alert("Save failed: " + error.message)
  } finally {
    done()
  }
}

const apply = async () => {
  try {
    const response = await axios.get('/setting/apply');
    if (response.data?.status === '1') {
      alert('Settings applied successfully!')
    } else {
      alert('Apply failed: ' + (response.data?.error || 'Unknown error'))
    }
  } catch (error) {
    alert('Apply failed: ' + error.message)
  }
}
</script>
