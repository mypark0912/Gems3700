<template>
    <footer
      class="sticky bottom-0 before:absolute before:inset-0 before:backdrop-blur-md max-lg:before:bg-white/90 dark:max-lg:before:bg-gray-800/90 before:-z-10 z-30"
      :class="[
        variant === 'v2' || variant === 'v3' ? 'before:bg-white after:absolute after:h-px after:inset-x-0 after:top-full after:bg-gray-200 dark:after:bg-gray-700/60 after:-z-10' : 'max-lg:shadow-sm lg:before:bg-gray-100/90 dark:lg:before:bg-gray-900/90',
        variant === 'v2' ? 'dark:before:bg-gray-800' : '',
        variant === 'v3' ? 'dark:before:bg-gray-900' : '',
      ]"
    >
      <div class="px-4 sm:px-6 lg:px-8">
        <div
          class="flex items-center justify-between h-12"
          :class="variant === 'v2' || variant === 'v3' ? '' : 'lg:border-t border-gray-200 dark:border-gray-700/60'"
        >
          <!-- Footer: Left side -->
          <div class="flex">
          </div>
  
          <!-- Footer: Right side -->
          <div class="flex items-center space-x-3">
            <p class="text-sm text-gray-300 dark:text-gray-100">
              &copy; 2025 Signals Sense. All rights reserved.
            </p>

          </div>
        </div>
      </div>
    </footer>
  </template>
  
  <script>
  import { reactive, watch } from 'vue'
  

  
  export default {
    name: 'Footer',
    props: [
      'variant',
      'langset'
    ],
    components: {

    },
    setup(props) {
      // ✅ langset이 null일 경우 기본 객체 설정
      const langs = reactive(props.langset || {});
      //console.log("langs.value",langs.value);
      // ✅ props.langset이 변경될 때마다 langs 업데이트 (반응형 유지)
      watch(() => props.langset, (newLangset) => {
        Object.assign(langs, newLangset); // ✅ 기존 객체 유지하면서 업데이트
        console.log("✅ langset 업데이트됨:", langs);
      });
  
      return {
        langs,
      }
    }
  }
  </script>
  
  <style scoped>
  /* 원하는 스타일 추가 */
  </style>
  