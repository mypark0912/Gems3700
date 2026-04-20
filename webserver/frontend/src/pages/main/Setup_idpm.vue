<template>
  <div class="flex h-[100dvh] overflow-hidden">
    <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" />

    <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden bg-gray-50 dark:bg-gray-950">
      <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <main class="grow">
        <div class="px-4 sm:px-6 lg:px-8 py-8 w-full">

          <div class="mb-6">
            <h2 class="text-xl md:text-2xl text-gray-800 dark:text-gray-100 font-bold ml-1">
              Setup > {{ modeName }}
            </h2>
          </div>

          <GeneralSetting v-if="mode === 'general'" />
          <MainSetting v-else-if="mode === 'main'" channel="Main" />
          <BranchSetting v-else-if="mode === 'branch'" />

          <div v-else class="bg-white dark:bg-gray-800 shadow-sm rounded-xl border border-gray-200 dark:border-gray-700/60 p-5">
            <p class="text-center text-gray-500 py-12">설정 항목을 선택해주세요.</p>
          </div>
        </div>
      </main>

      <Footer />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, provide } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useInputDict } from '@/composables/useInputDict'
import Sidebar from '../common/SideBar.vue'
import Header from '../common/Header.vue'
import Footer from '../common/Footer.vue'

import GeneralSetting from '../../partials/inners/setting/GeneralSetting.vue'
import MainSetting from '../../partials/inners/setting/Panel_Main.vue'
import BranchSetting from '../../partials/inners/setting/BranchSetting.vue'

const props = defineProps({
  mode: {
    type: String,
    default: 'general'
  }
})

const sidebarOpen = ref(true)

const authStore = useAuthStore()
const { inputDict, channel_main, channel_sub } = useInputDict()
const devMode = computed(() => authStore.getOpMode)

provide('inputDict', inputDict)
provide('channel_main', channel_main)
provide('channel_sub', channel_sub)
provide('devMode', devMode)

const modeName = computed(() => {
  if (props.mode === 'general') return 'General'
  if (props.mode === 'main') return 'Main'
  if (props.mode === 'branch') return 'Branch'
  return props.mode
})
</script>