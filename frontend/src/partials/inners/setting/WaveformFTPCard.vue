<template>
  <!-- Waveform FTP -->
  <div
    v-if="showFTP"
    class="relative col-span-full xl:col-span-8 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg"
  >
    <div
      class="absolute top-0 left-0 right-0 h-0.5 bg-orange-500"
      aria-hidden="true"
    ></div>
    <div
      class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60"
    >
      <header class="flex items-center justify-between mb-2">
        <div class="flex items-center">
          <div class="w-6 h-6 rounded-full shrink-0 bg-orange-500 mr-3">
            <svg
              class="w-6 h-6 fill-current text-white"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M19 18H6a4 4 0 01-1-7.9A5 5 0 0112 5a5 5 0 014.9 3.8A4 4 0 0119 18z"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M12 12v-4m0 0l-3 3m3-3l3 3"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M12 14v4m0 0l-3-3m3 3l3-3"
                stroke="currentColor"
                stroke-width="2"
                fill="none"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            Waveform FTP
          </h3>
        </div>
        
        <!-- Auto/Manual Toggle -->
        <div class="flex items-center space-x-3 h-10">
          <label class="relative inline-flex items-center cursor-pointer">
            <input
              type="checkbox"
              v-model="inputDict.ftpInfo.isManual"
              class="sr-only peer"
            />
            <div
              class="w-11 h-6 bg-gray-300 rounded-full relative transition-all peer-checked:bg-violet-500"
              :class="{
                'bg-violet-500': inputDict.ftpInfo.isManual,
              }"
            >
              <div
                class="absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition-all transform"
                :class="inputDict.ftpInfo.isManual ? 'translate-x-5' : ''"
              ></div>
            </div>
          </label>
          <span class="text-sm font-medium">
            {{ inputDict.ftpInfo.isManual ? "Manual" : "Auto" }}
          </span>
        </div>
      </header>
    </div>
    
    <!-- FTP Settings - Always show credentials, conditionally show upload directories -->
    <div class="grid grid-cols-2 gap-6 px-5 py-6">
      <!-- Left: Host / Account / Password -->
      <div class="space-y-4">
        <div
          class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-1 mt-2"
        >
          FTP Credentials
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5">Host</label>
          <input
            v-model="inputDict.ftpInfo.host"
            class="form-input w-full"
            type="text"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5">Port</label>
          <input
            v-model="inputDict.ftpInfo.port"
            class="form-input w-full"
            type="text"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5">Account</label>
          <input
            v-model="inputDict.ftpInfo.id"
            class="form-input w-full"
            type="text"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5">Password</label>
          <input
            v-model="inputDict.ftpInfo.pass"
            class="form-input w-full"
            type="password"
          />
        </div>
      </div>

      <!-- Right: Upload Directories - Only show when Manual mode is selected -->
      <div class="space-y-4">
        <!-- Manual mode: Show upload directories -->
        <div v-if="inputDict.ftpInfo.isManual">
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-1 mt-2"
          >
            Main Channel
          </div>
          <div>
            <label class="block text-sm font-medium mb-1.5"
              >Upload Directory</label
            >
            <input
              v-model="inputDict.ftpInfo.upload_main"
              class="form-input w-full"
              type="text"
            />
          </div>
          <div
            class="text-xs text-gray-800 dark:text-gray-100 font-semibold uppercase mb-1 mt-2"
          >
            Sub Channel
          </div>
          <div>
            <label class="block text-sm font-medium mb-1.5"
              >Upload Directory</label
            >
            <input
              v-model="inputDict.ftpInfo.upload_sub"
              class="form-input w-full"
              type="text"
            />
          </div>
        </div>
        
        <!-- Auto mode: Show message -->
        <div v-else class="flex items-center justify-center h-full min-h-[200px]">
          <div class="text-center text-gray-500 dark:text-gray-400">
            <svg class="w-12 h-12 mx-auto mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path>
            </svg>
            <p class="text-sm">Auto mode enabled<br>Upload directories will be configured automatically</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from "vue";

const inputDict = inject("inputDict");

// Props
defineProps({
  showFTP: {
    type: Boolean,
    required: true,
  },
});
</script>