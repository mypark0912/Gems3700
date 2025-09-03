<template>
    <div class="col-span-full xl:col-span-6 2xl:col-span-6 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-lg">
              <!-- Card content -->
              <div class="flex flex-col h-full p-5">
                <div class="grow">
                  <header class="flex items-center mb-4">
                    <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">{{ item.name }}</h3>
                  </header>
                  <div class="flex justify-center">
                    <div class="grid grid-cols-3 gap-4 mt-2">
                      <!-- Main Channel -->
                      <div class="text-left">
                        <div class="text-sm font-bold text-black-400 dark:text-black-500 uppercase mb-1">Value</div>
                        <div class="flex flex-col gap-1">
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 dark:bg-violet-600/40 dark:text-violet-300 rounded-full">L1</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.value[0]).toFixed(2) }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 dark:bg-sky-600/40 dark:text-sky-300 rounded-full">L2</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.value[1]).toFixed(2) }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-lime-700 px-1.5 bg-lime-500/20 dark:bg-lime-600/40 dark:text-lime-300 rounded-full">L3</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.value[2]).toFixed(2) }}</span>
                          </div>
                        </div>
                      </div>

                      <!-- Sub Channel -->
                      <div class="text-left">
                        <div class="text-sm font-bold text-black-400 dark:text-black-500 uppercase mb-1">Error</div>
                        <div class="flex flex-col gap-1">
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-violet-700 px-1.5 bg-violet-500/20 dark:bg-violet-600/40 dark:text-violet-300 rounded-full">L1</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.error[0]).toFixed(2) }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-sky-700 px-1.5 bg-sky-500/20 dark:bg-sky-600/40 dark:text-sky-300 rounded-full">L2</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.error[1]).toFixed(2) }}</span>
                          </div>
                          <div class="flex items-center gap-2">
                            <span class="text-xs font-medium text-lime-700 px-1.5 bg-lime-500/20 dark:bg-lime-600/40 dark:text-lime-300 rounded-full">L3</span>
                            <span class="text-basic font-bold text-gray-800 dark:text-gray-100">{{ Number(item.error[2]).toFixed(2) }}</span>
                          </div>
                        </div>
                      </div>
                    </div>

                    </div>
                </div>
                <!-- Card footer -->
                <!--footer class="mt-4">
                    <div class="flex justify-center items-center gap-3">
                        <label for="reference" class="text-sm font-medium text-gray-700 dark:text-gray-300">
                        Reference
                        </label>

                        <input
                        id="reference"
                        type="text"
                        class="form-input h-6 px-3 rounded-md border border-gray-300 w-24"
                        />

                        <button
                        class="h-6 px-4 bg-sky-900 text-white text-sm rounded-md hover:bg-sky-800"
                        >
                        Send
                        </button>
                    </div>
                </footer-->
              </div>
            </div>
</template>
<script>
import { ref, watch, watchEffect } from 'vue'


export default {
  name: 'DashCard',
  props: ['item'],
  setup(props){

    const item = ref(props.item);

    return {
      item,
    }
  }
}
</script>