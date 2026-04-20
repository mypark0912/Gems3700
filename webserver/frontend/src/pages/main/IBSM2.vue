<template>
  <div class="flex h-[100dvh] overflow-hidden">

    <!-- Sidebar -->
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <!-- Content area -->
    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden">

      <!-- Site header -->
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-2 sm:px-4 lg:px-6 py-4 w-full max-w-full">

          <!-- Title -->
          <div class="sm:flex sm:justify-between sm:items-center mb-4">
            <div class="mb-4 sm:mb-0">
              <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold">IBSM2</h2>
            </div>
          </div>

          <!-- Cards -->
          <div class="md:col-span-12 bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 w-full">
            <GemsDashCard v-if="Object.keys(mainData).length > 0" :data="mainData" :channel="'A'" />
            <div v-else class="py-8 text-center text-gray-400 dark:text-gray-500 text-sm">
              데이터 없음
            </div>
          </div>

        </div>
      </main>

      <Footer />
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'
import GemsDashCard from '../../partials/inners/dashboard/GemsDashCard.vue'
import axios from 'axios'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'IBSM2',
  components: {
    Sidebar,
    Header,
    Footer,
    GemsDashCard,
  },
  setup() {
    const sidebarOpen = ref(false)
    const authStore = useAuthStore()
    const langset = computed(() => authStore.getLang)

    // ── 더미 데이터 (GemsDashCard는 JSON.parse(data[key])를 하므로 JSON 문자열로 넣어야 함) ──
    const mainData = ref({
      1: JSON.stringify({
        status: 1,
        updateTime: Date.now() / 1000,
        temp: 35.2,
        cblist: [{
          cbtype: 5,
          irms: [12.50, 11.80, 13.20],
          power: [2750],
          kwh: 15.60,
          ig: 0.12,
          pthd: 3.45,
        }]
      }),
      2: JSON.stringify({
        status: 1,
        updateTime: Date.now() / 1000,
        temp: 32.1,
        cblist: [{
          cbtype: 1,
          irms: [8.30],
          power: [1820],
          kwh: 8.90,
          ig: 0.05,
          pthd: 2.10,
        }]
      }),
      3: JSON.stringify({
        status: 1,
        updateTime: Date.now() / 1000,
        temp: 31.5,
        cblist: [{
          cbtype: 2,
          irms: [7.60],
          power: [1670],
          kwh: 7.20,
          ig: 0.03,
          pthd: 1.80,
        }]
      }),
      4: JSON.stringify({
        status: 0,
        updateTime: Date.now() / 1000,
        temp: 28.0,
        cblist: [{
          cbtype: 6,
          irms: [0],
          power: [0],
          kwh: 3.20,
          ig: 15.20,
          pthd: 0,
        }]
      }),
      5: JSON.stringify({
        status: 3,
        updateTime: Date.now() / 1000,
        temp: 38.5,
        cblist: [{
          cbtype: 4,
          irms: [25.10, 24.80, 26.00],
          power: [5520],
          kwh: 42.00,
          ig: 0.08,
          pthd: 5.20,
        }]
      }),
      6: JSON.stringify({
        status: 1,
        updateTime: Date.now() / 1000,
        temp: 30.0,
        cblist: [{
          cbtype: 3,
          irms: [5.20],
          power: [1140],
          kwh: 4.50,
          ig: 0.01,
          pthd: 2.80,
        }]
      }),
    })

    // TODO: 더미 데이터 제거 후 아래 API 호출 활성화
    // let updateInterval = null
    // const fetchData = async () => {
    //   try {
    //     const response = await axios.get('/api/getSubData/A')
    //     if (response.data.success) {
    //       mainData.value = { ...mainData.value, ...response.data.data }
    //     }
    //   } catch (error) {
    //     console.log('IBSM2 데이터 가져오기 실패:', error)
    //   }
    // }
    //
    // onMounted(() => {
    //   fetchData()
    //   updateInterval = setInterval(fetchData, 5000)
    // })
    //
    // onUnmounted(() => {
    //   if (updateInterval) clearInterval(updateInterval)
    // })

    return {
      sidebarOpen,
      langset,
      mainData,
    }
  }
}
</script>