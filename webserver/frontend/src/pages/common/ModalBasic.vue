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
    <div v-if="modalOpen" :id="id" class="fixed inset-0 z-50 overflow-hidden flex items-center justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
      <div ref="modalContent" class="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-auto min-w-[300px] max-w-screen-lg max-h-full">
        <!-- Modal header -->
        <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700/60">
          <div class="flex justify-between items-center">
            <div class="font-semibold text-gray-800 dark:text-gray-100">{{ title }}</div>
            <button class="text-gray-400 dark:text-gray-500 hover:text-gray-500 dark:hover:text-gray-400" @click.stop="$emit('close-modal')">
              <div class="sr-only">Close</div>
              <svg class="fill-current" width="16" height="16" viewBox="0 0 16 16">
                <path d="M7.95 6.536l4.242-4.243a1 1 0 111.415 1.414L9.364 7.95l4.243 4.242a1 1 0 11-1.415 1.415L7.95 9.364l-4.243 4.243a1 1 0 01-1.414-1.415L6.536 7.95 2.293 3.707a1 1 0 011.414-1.414L7.95 6.536z" />
              </svg>
            </button>
          </div>
        </div>
        <slot />       
      </div>
    </div>
  </transition>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'

export default {
  name: 'ModalBasic',
  props: ['id', 'modalOpen', 'title'],
  emits: ['close-modal'],
  setup(props, { emit }) {

    const modalContent = ref(null)

    // close on click outside
    // const clickHandler = ({ target }) => {
    //   if (!props.modalOpen || modalContent.value.contains(target)) return
    //   emit('close-modal')
    // }

    const clickHandler = (event) => {
      if (!props.modalOpen) return;

      // 클릭한 대상이 모달 내부인지 확인
      if (modalContent.value && modalContent.value.contains(event.target)) {
        event.stopPropagation(); // 이벤트 전파 방지
        return;
      }

      emit('close-modal');
    };


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

    return {
      modalContent,
    }
  }  
}
</script>

<style scoped>
table {
  table-layout: fixed;
  width: 100%;
  min-height: calc(40px * 20); /* 40px * 20줄 높이 고정 */
  border-collapse: collapse;
}

th,
td {
  text-align: center;
  padding: 8px;
  border: 1px solid #ddd;
  height: 40px; /* 각 행 높이 고정 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
th:nth-child(1) {
  width: 300px;
} /* 첫 번째 컬럼 */
th:nth-child(2) {
  width: 150px;
}
th:nth-child(3) {
  width: 150px;
}
th:nth-child(4) {
  width: 200px;
}
th:nth-child(5) {
  width: 100px;
}
.empty-row {
  visibility: hidden; /* 빈 행 안보이게 처리 */
}


</style>
