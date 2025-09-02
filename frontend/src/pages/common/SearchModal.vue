<template>
  <!-- Modal backdrop -->
  <transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition ease-out duration-100"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div v-show="modalOpen" class="fixed inset-0 bg-gray-900 bg-opacity-30 z-50 transition-opacity" aria-hidden="true"></div>
  </transition>
  <!-- Modal dialog -->
  <transition
    enter-active-class="transition ease-in-out duration-200"
    enter-from-class="opacity-0 translate-y-4"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition ease-in-out duration-200"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-4"
  >
    <div v-show="modalOpen" :id="id" class="fixed inset-0 z-50 overflow-hidden flex items-start top-20 mb-4 justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
      <div ref="modalContent" class="bg-white dark:bg-gray-800 border border-transparent dark:border-gray-700/60 overflow-auto max-w-md w-full max-h-full rounded-lg shadow-lg">
        <!-- Search form -->
        <form class="border-b border-gray-200 dark:border-gray-700/60">
          <div class="relative">
            <label :for="searchId" class="sr-only">Search</label>
            <input :id="searchId" v-model="searchText" class="w-full dark:text-white bg-white dark:bg-gray-800 border-0 focus:ring-transparent placeholder-gray-400 dark:placeholder-gray-500 appearance-none py-3 pl-10 pr-4" type="search" placeholder="Search Anything…" ref="searchInput" />
            <button class="absolute inset-0 right-auto group" aria-label="Search">
              <svg class="shrink-0 fill-current text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400 ml-4 mr-2" width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
                <path d="M7 14c-3.86 0-7-3.14-7-7s3.14-7 7-7 7 3.14 7 7-3.14 7-7 7zM7 2C4.243 2 2 4.243 2 7s2.243 5 5 5 5-2.243 5-5-2.243-5-5-5z" />
                <path d="M15.707 14.293L13.314 11.9a8.019 8.019 0 01-1.414 1.414l2.393 2.393a.997.997 0 001.414 0 .999.999 0 000-1.414z" />
              </svg>
            </button>
              <button
                type="button"
                @click="$emit('close-modal')"
                class="absolute right-0 top-1/2 -translate-y-1/2 p-2 text-gray-400 hover:text-gray-600"
                aria-label="Close"
            >
                <svg class="w-4 h-4" viewBox="0 0 16 16" fill="currentColor">
                <path d="M2 2L14 14M2 14L14 2" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                </svg>
            </button>
          </div>
        </form>
        <div class="py-4 px-2">
          <!-- Recent searches -->
          <div class="mb-3 last:mb-0">
            <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase px-2 mb-2">Recent searches</div>
            <ul class="text-sm divide-y divide-gray-200 dark:divide-gray-700">
                <li
                    v-for="(item, index) in filteredValues"
                    :key="index"
                    class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700/20 cursor-pointer flex items-center"
                    @click="$emit('select', item.Name)"
                >
                    <svg class="w-4 h-4 text-gray-400 dark:text-gray-500 mr-2" viewBox="0 0 16 16">
                    <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1" fill="none" />
                    </svg>
                    <span class="text-gray-800 dark:text-gray-100">{{ item.Name }}</span>
                </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { ref, nextTick, onMounted, onUnmounted, watch, inject, computed } from 'vue'

export default {
  name: 'SearchModal',
  props: ['id', 'searchId', 'modalOpen'],
  emits: ['open-modal', 'close-modal'],
  setup(props, { emit }) {

    const modalContent = ref(null)
    const searchInput = ref(null)
    const BearingOptions = inject('BearingOptions')
    const BearingValues = inject('BearingValues')

    const searchText = ref('')  // 입력값

    const filteredValues = computed(() => {
    //if (!searchText.value) return []

    return BearingValues.value
        .filter(b => b.Name.toLowerCase().includes(searchText.value.toLowerCase()))
        .slice(0, 5) // ✅ 최대 50개만 출력
    })
    
    // close on click outside
    const clickHandler = ({ target }) => {
      if (!props.modalOpen || modalContent.value.contains(target)) return
      emit('close-modal')
    }

    // close if the esc key is pressed
    const keyHandler = ({ keyCode }) => {
      if (!props.modalOpen || keyCode !== 27) return
      emit('close-modal')
    }

    onMounted(() => {
      //document.addEventListener('click', clickHandler)
      document.addEventListener('keydown', keyHandler)
    })

    onUnmounted(() => {
      //document.removeEventListener('click', clickHandler)
      document.removeEventListener('keydown', keyHandler)
    })

    // watch(
    //     () => props.modalOpen,
    //     (val) => {
    //         console.log('[WATCH TRIGGERED] modalOpen =', val)
    //     },
    //     { immediate: true }
    // );

    // watch(() => props.modalOpen, (open) => {
    //     console.log('SearchModal modalOpen:', open)
    //   open && nextTick(() => searchInput.value.focus())
    // })

    
    // onMounted(async () => {
    //     try {
    //         const response = await axios.get("/setting/checkBearing");
    //         if (response.data.passOK == "1") {
    //         for (let i = 0; i < response.data.data.length; i++) {
    //             BearingOptions.value.push(response.data.data[i]["Name"]);
    //             BearingValues.value.push(response.data.data[i]);
    //         }
    //         }
    //     } catch (error) {
    //         message.value = "업로드 실패: " + error.response.data.error;
    //     }
    //     });

    return {
      modalContent,
      searchInput,
      BearingOptions,
      BearingValues,
      filteredValues,
      searchText,
    }
  }
}
</script>