<template>
    <div class="relative inline-flex">
      <button
        ref="trigger"
        class="w-8 h-8 flex items-center justify-center hover:bg-gray-100 lg:hover:bg-gray-200 dark:hover:bg-gray-700/50 dark:lg:hover:bg-gray-800 rounded-full"
        :class="{ 'bg-gray-200 dark:bg-gray-800': dropdownOpen }"
        aria-haspopup="true"
        @click.prevent="dropdownOpen = !dropdownOpen"
        :aria-expanded="dropdownOpen"
      >
        <span class="sr-only">Languages</span>
        <component :is="langIconMap[lang]" />
      </button>
      <transition
        enter-active-class="transition ease-out duration-200 transform"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-out duration-200"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-show="dropdownOpen" class="origin-top-right z-10 absolute top-full min-w-40 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 py-1.5 rounded-lg shadow-lg overflow-hidden mt-1"  :class="align === 'right' ? 'right-0' : 'left-0'">
          <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase pt-1.5 pb-2 px-3">Languages</div>
          <ul ref="dropdown" @focusin="dropdownOpen = true" @focusout="dropdownOpen = false">
            <li v-for="item in languages" :key="item.code">
                <button
                @click="setLanguage(item.code)"
                class="w-full flex items-center py-1 px-3 text-sm font-medium"
                :class="lang === item.code
                    ? 'text-violet-500 dark:text-violet-400'
                    : 'text-gray-600 dark:text-gray-300 hover:text-violet-500 dark:hover:text-violet-400'"
                >
                <component :is="item.icon" class="shrink-0 mr-2" />
                <span>{{ item.label }}</span>
                </button>
            </li>
            </ul>
        </div> 
      </transition>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, onUnmounted, defineComponent, h } from 'vue'
  import { useI18n } from 'vue-i18n'

  const IconKo = defineComponent({
  name: 'IconKo',
  render() {
    return h('svg', {
      class: 'stroke-current text-gray-500/80 dark:text-gray-400/80 shrink-0 mr-2',
      xmlns: 'http://www.w3.org/2000/svg',
      viewBox: '0 0 16 16',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      width: '16',
      height: '16',
      'stroke-width': '1.5'
    }, [
      h('path', { d: 'M4.67 4.67h4c0 1.67 -1.06 5.65 -4 6.67' }),
      h('path', { d: 'M10.67 3.33v9.33z' }),
      h('path', { d: 'M10.67 8h1.33' })
    ])
  }
});

const IconEn = defineComponent({
  name: 'IconEn',
  render() {
    return h('svg', {
      class: 'stroke-current text-gray-500/80 dark:text-gray-400/80 shrink-0 mr-2',
      xmlns: 'http://www.w3.org/2000/svg',
      viewBox: '0 0 16 16',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      width: '16',
      height: '16',
      'stroke-width': '1.5'
    }, [
      h('path', { d: 'M4 6.67h1.33a1.33 1.33 0 0 1 1.33 1.33v3.33H5.33a1.33 1.33 0 1 1 0 -2.67H6.67' }),
      h('path', { d: 'M9.33 4.67v6.67' }),
      h('path', { d: 'M9.33 6.67m0 1.33a1.33 1.33 0 0 1 1.33 -1.33h0.67a1.33 1.33 0 0 1 1.33 1.33v2a1.33 1.33 0 0 1 -1.33 1.33h-0.67a1.33 1.33 0 0 1 -1.33 -1.33z' })
    ])
  }
});

const IconJa = defineComponent({
  name: 'IconJa',
  render() {
    return h('svg', {
      class: 'stroke-current text-gray-500/80 dark:text-gray-400/80 shrink-0 mr-2',
      xmlns: 'http://www.w3.org/2000/svg',
      viewBox: '0 0 16 16',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      width: '16',
      height: '16',
      'stroke-width': '1.5'
    }, [
      h('path', { d: 'M2.67 3.33h4.67' }),
      h('path', { d: 'M4.67 2.67c0 3.23 0 4.67 0.33 5.33' }),
      h('path', { d: 'M6.67 5.67c0 1.52 -1.33 3 -2.33 3s-1.67 -0.76 -1.67 -1.33c0 -1.33 0.67 -2 2 -2s3.33 0.38 3.33 1.9c0 1.02 -0.44 1.71 -1.33 2.1' }),
      h('path', { d: 'M8 13.33l2.67 -6l2.67 6' }),
      h('path', { d: 'M12.73 12h-4.13' })
    ])
  }
});

  export default {
    name: 'DropdownLanguage',
    props: ['align'],
    setup() {
  
      const dropdownOpen = ref(false)
      const trigger = ref(null)
      const dropdown = ref(null)
      const { t, locale } = useI18n()
      const lang = ref(locale.value) ;

    const languages = [
      { code: 'ko', label: 'KOREAN', icon: IconKo },
      { code: 'en', label: 'ENGLISH', icon: IconEn },
      { code: 'ja', label: 'JAPANESE', icon: IconJa }
    ]

    const langIconMap = {
      ko: IconKo,
      en: IconEn,
      ja: IconJa
    }

    const setLanguage = (code) => {
      lang.value = code
      locale.value = code
      dropdownOpen.value = false
    }
  
      // close on click outside
      const clickHandler = ({ target }) => {
        if (!dropdownOpen.value || dropdown.value.contains(target) || trigger.value.contains(target)) return
        dropdownOpen.value = false
      }
  
      // close if the esc key is pressed
      const keyHandler = ({ keyCode }) => {
        if (!dropdownOpen.value || keyCode !== 27) return
        dropdownOpen.value = false
      }
  
      onMounted(() => {
        document.addEventListener('click', clickHandler)
        document.addEventListener('keydown', keyHandler)
      })
  
      onUnmounted(() => {
        document.removeEventListener('click', clickHandler)
        document.removeEventListener('keydown', keyHandler)
      })
  
      return {
        dropdownOpen,
        trigger,
        dropdown,
        t,
        lang,
        languages,
        langIconMap,
        setLanguage,
      }
    }
  }
  </script>