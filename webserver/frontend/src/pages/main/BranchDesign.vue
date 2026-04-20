<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">
          <!-- 타이틀 + 탭 -->
          <div class="mb-4 flex items-center gap-3 flex-wrap">
            <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">
              {{ t('sidebar.module') }} > IPSM72 #{{ id }}
            </h2>
            <span class="text-xl md:text-2xl text-gray-400 dark:text-gray-500 font-bold">&gt;</span>
            <div class="flex items-center gap-1">
              <button
                class="px-3 py-1 text-lg md:text-xl font-bold rounded-md transition-colors"
                :class="activeTab === 'sv'
                  ? 'text-violet-600 dark:text-violet-400 bg-violet-50 dark:bg-violet-900/20'
                  : 'text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
                @click="activeTab = 'sv'"
              >System / Voltage</button>
              <span class="text-gray-300 dark:text-gray-600">·</span>
              <button
                class="px-3 py-1 text-lg md:text-xl font-bold rounded-md transition-colors"
                :class="activeTab === 'feeder'
                  ? 'text-violet-600 dark:text-violet-400 bg-violet-50 dark:bg-violet-900/20'
                  : 'text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
                @click="activeTab = 'feeder'"
              >Feeder</button>
            </div>
          </div>

          <!-- 탭 컨텐츠 -->
          <IpsmSystem
            v-if="activeTab === 'sv'"
            :systemData="systemData"
            :voltageData="voltageData"
          />
          <IpsmFeeder
            v-else-if="activeTab === 'feeder'"
            :feederData="feederData"
            :energyData="energyData"
            :voltageData="voltageData"
            :selectedDevice="selectedDevice"
            @select="onFeederSelect"
          />
        </div>
      </main>
      <Footer />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'
import IpsmSystem from './IpsmSystem.vue'
import IpsmFeeder from './IpsmFeeder.vue'

const { t } = useI18n()

const props = defineProps({
  id: { type: String, required: true },
})

const sidebarOpen = ref(false)
const activeTab = ref('feeder') // 'sv' | 'feeder'

const systemData = ref(null)
const voltageData = ref(null)
const feederData = ref(null)
const energyData = ref(null)
const selectedDevice = ref(1)

let ipsmPollInterval = null
let feederPollInterval = null

async function fetchIpsmData(device) {
  if (!device) return
  try {
    const res = await axios.get(`/api/getIpsmData/${device}`)
    if (res.data.success) {
      systemData.value = res.data.data.system
      voltageData.value = res.data.data.voltage
    } else {
      console.error('[IpsmData] API error:', res.data.error)
    }
  } catch (error) {
    console.error('[IpsmData] fetch error:', error)
  }
}

async function fetchFeeder(device, feederId) {
  if (!device || !feederId) return
  try {
    const res = await axios.get(`/api/getFeeder/${device}/${feederId}`)
    if (res.data.success) {
      feederData.value = res.data.data.feeder
      energyData.value = res.data.data.energy
    } else {
      console.error('[Feeder] API error:', res.data.error)
    }
  } catch (error) {
    console.error('[Feeder] fetch error:', error)
  }
}

function stopPolling() {
  if (ipsmPollInterval) { clearInterval(ipsmPollInterval); ipsmPollInterval = null }
  if (feederPollInterval) { clearInterval(feederPollInterval); feederPollInterval = null }
}

function startPollingForTab(tab, deviceId) {
  stopPolling()
  if (tab === 'sv') {
    ipsmPollInterval = setInterval(() => fetchIpsmData(deviceId), 5000)
  } else if (tab === 'feeder') {
    // Feeder 탭은 feeder만 폴링. voltage/system 은 탭 진입 시 1회 fetch 로 충분.
    feederPollInterval = setInterval(() => fetchFeeder(deviceId, selectedDevice.value), 5000)
  }
}

async function loadTabData(tab, deviceId) {
  if (tab === 'sv') {
    await fetchIpsmData(deviceId)
  } else if (tab === 'feeder') {
    await Promise.all([
      fetchIpsmData(deviceId),
      fetchFeeder(deviceId, selectedDevice.value),
    ])
  }
}

function onFeederSelect(idx) {
  selectedDevice.value = selectedDevice.value === idx ? null : idx
  if (selectedDevice.value) {
    fetchFeeder(props.id, selectedDevice.value)
  }
}

onMounted(async () => {
  await loadTabData(activeTab.value, props.id)
  startPollingForTab(activeTab.value, props.id)
})

onUnmounted(() => {
  stopPolling()
})

watch(activeTab, async (tab) => {
  await loadTabData(tab, props.id)
  startPollingForTab(tab, props.id)
})

watch(() => props.id, async (newId) => {
  selectedDevice.value = 1
  await loadTabData(activeTab.value, newId)
  startPollingForTab(activeTab.value, newId)
})
</script>
