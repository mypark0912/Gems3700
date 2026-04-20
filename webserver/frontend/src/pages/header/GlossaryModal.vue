// GlossaryModal.vue
<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/50" @click="close"></div>
        
        <div class="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-5xl h-[80vh] flex flex-col">
          <!-- 헤더 -->
          <div class="flex items-center justify-between px-5 py-4 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">
              {{ isKorean ? '용어설명집' : isJapanese ? '用語集' : 'Glossary' }}
            </h3>
            <button @click="close" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- 검색창 -->
          <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
            <div class="flex items-center gap-2">
              <div class="relative flex-1">
                <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                </svg>
                <input 
                  v-model="searchQuery"
                  type="text" 
                  :placeholder="isKorean ? '용어 검색...' : isJapanese ? '用語検索...' : 'Search terms...'"
                  class="w-full pl-10 pr-10 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-violet-500"
                  @keydown.enter.exact="goToNext"
                  @keydown.shift.enter.exact="goToPrev"
                />
                <button v-if="searchQuery" @click="clearSearch" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
              
              <template v-if="searchQuery && searchResults.length > 0">
                <span class="text-sm text-gray-500 whitespace-nowrap">{{ currentSearchIndex + 1 }} / {{ searchResults.length }}</span>
                <button @click="goToPrev" class="p-2 text-gray-500 hover:text-violet-600 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"/>
                  </svg>
                </button>
                <button @click="goToNext" class="p-2 text-gray-500 hover:text-violet-600 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                  </svg>
                </button>
              </template>
              <span v-else-if="searchQuery" class="text-sm text-red-500 whitespace-nowrap">
                {{ isKorean ? '결과 없음' : isJapanese ? '結果なし' : 'No results' }}
              </span>
            </div>
          </div>
          
          <!-- 메인 콘텐츠 -->
          <div class="flex flex-1 overflow-hidden p-4 min-h-0">
            <!-- 왼쪽: 탭 + 트리뷰 -->
            <div class="w-80 flex flex-col h-full overflow-hidden">
              <div class="flex flex-col flex-1 border-r border-gray-200 dark:border-gray-700 pr-4 min-h-0">
                <div class="flex border-b border-gray-200 dark:border-gray-700 pb-1 mb-1 flex-shrink-0">
                  <button
                    v-for="tab in tabs"
                    :key="tab.key"
                    @click="activeTab = tab.key"
                    class="flex-1 px-3 py-2 text-xs font-medium transition-colors"
                    :class="activeTab === tab.key 
                      ? 'text-violet-600 border-b-2 border-violet-600 bg-violet-50 dark:bg-violet-900/20' 
                      : 'text-gray-500 hover:text-gray-700'"
                  >
                    {{ tab.label }}
                  </button>
                </div>
                
                <div class="flex-1 overflow-y-auto p-3 min-h-0">
                  <template v-if="activeTab === 'equipment'">
                    <div v-for="category in equipmentData" :key="category.key" class="mb-2">
                      <button
                        @click="toggleCategory(category.key)"
                        class="flex items-center gap-1 w-full text-left text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-violet-600 py-1"
                      >
                        <svg class="w-4 h-4 transition-transform" :class="expandedCategories.includes(category.key) ? 'rotate-90' : ''" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
                        </svg>
                        {{ category.label }}
                      </button>
                      <ul v-show="expandedCategories.includes(category.key)" class="ml-5">
                        <li v-for="item in category.items" :key="item.key">
                          <button
                            @click="scrollToItem(`equipment-${category.key}-${item.key}`)"
                            class="text-sm text-gray-600 dark:text-gray-400 hover:text-violet-600 py-0.5 w-full text-left"
                            :class="{ 'text-violet-600 font-medium': activeItem === `equipment-${category.key}-${item.key}` }"
                          >
                            {{ getLabel(item) }}
                          </button>
                        </li>
                      </ul>
                    </div>
                  </template>

                  <template v-else>
                    <ul class="space-y-1">
                      <li v-for="item in currentTabItems" :key="item.key">
                        <button
                          @click="scrollToItem(`${activeTab}-${item.key}`)"
                          class="text-sm text-gray-600 dark:text-gray-400 hover:text-violet-600 py-1 w-full text-left"
                          :class="{ 'text-violet-600 font-medium': activeItem === `${activeTab}-${item.key}` }"
                        >
                          {{ getLabel(item) }}
                        </button>
                      </li>
                    </ul>
                  </template>
                </div>
              </div>
            </div>

            <!-- 오른쪽: 내용 -->
            <div ref="contentArea" class="flex-1 overflow-y-auto p-5">
              <!-- Equipment -->
              <section class="mb-8">
                <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100 mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
                  {{ isKorean ? '장비' : isJapanese ? '設備' : 'Equipment' }}
                </h2>
                <div v-for="category in equipmentData" :key="category.key" class="mb-6">
                  <h3 class="text-lg font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ category.label }}</h3>
                  <dl class="space-y-4 ml-4">
                    <div 
                      v-for="item in category.items" 
                      :key="item.key"
                      :id="`equipment-${category.key}-${item.key}`"
                      class="pb-4 border-b border-gray-100 dark:border-gray-700 last:border-0 scroll-mt-4"
                    >
                      <dt class="font-medium text-violet-700 dark:text-violet-500 mb-2" v-html="highlightText(getLabel(item))"></dt>
                      <dd class="mt-2 text-sm text-gray-600 dark:text-gray-400 leading-relaxed" 
                          style="line-height: 1.8; word-break: keep-all; overflow-wrap: break-word; white-space: pre-wrap;"
                          v-html="highlightText(formatDescription(getDescription(item)))"></dd>
                      <img v-if="item.image" :src="item.image" :alt="getLabel(item)" class="mt-3 max-w-full h-auto rounded-lg border border-gray-200 dark:border-gray-600" />
                    </div>
                  </dl>
                </div>
              </section>

              <!-- PQ -->
              <section class="mb-8">
                <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100 mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
                  {{ isKorean ? '전력품질 (PQ)' : isJapanese ? '電力品質 (PQ)' : 'Power Quality (PQ)' }}
                </h2>
                <dl class="space-y-4">
                  <div v-for="item in pqData" :key="item.key" :id="`pq-${item.key}`" class="pb-4 border-b border-gray-100 dark:border-gray-700 last:border-0 scroll-mt-4">
                    <dt class="font-medium text-violet-700 dark:text-violet-500 mb-2" v-html="highlightText(getLabel(item))"></dt>
                    <dd class="mt-2 text-sm text-gray-600 dark:text-gray-400 leading-relaxed" 
                        style="line-height: 1.8; word-break: keep-all; overflow-wrap: break-word; white-space: pre-wrap;"
                        v-html="highlightText(formatDescription(getDescription(item)))"></dd>
                    <img v-if="item.image" :src="item.image" :alt="getLabel(item)" class="mt-3 max-w-full h-auto rounded-lg border border-gray-200 dark:border-gray-600" />
                  </div>
                </dl>
              </section>

              <!-- Fault -->
              <section class="mb-8">
                <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100 mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
                  {{ isKorean ? '결함' : isJapanese ? '故障' : 'Fault' }}
                </h2>
                <dl class="space-y-4">
                  <div v-for="item in faultData" :key="item.key" :id="`fault-${item.key}`" class="pb-4 border-b border-gray-100 dark:border-gray-700 last:border-0 scroll-mt-4">
                    <dt class="font-medium text-violet-700 dark:text-violet-500 mb-2" v-html="highlightText(getLabel(item))"></dt>
                    <dd class="mt-2 text-sm text-gray-600 dark:text-gray-400 leading-relaxed" 
                        style="line-height: 1.8; word-break: keep-all; overflow-wrap: break-word; white-space: pre-wrap;"
                        v-html="highlightText(formatDescription(getDescription(item)))"></dd>
                    <img v-if="item.image" :src="item.image" :alt="getLabel(item)" class="mt-3 max-w-full h-auto rounded-lg border border-gray-200 dark:border-gray-600" />
                  </div>
                </dl>
              </section>

              <!-- Event -->
              <section class="mb-8">
                <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100 mb-4 pb-2 border-b border-gray-200 dark:border-gray-700">
                  {{ isKorean ? '이벤트' : isJapanese ? 'イベント' : 'Event' }}
                </h2>
                <dl class="space-y-4">
                  <div v-for="item in eventData" :key="item.key" :id="`event-${item.key}`" class="pb-4 border-b border-gray-100 dark:border-gray-700 last:border-0 scroll-mt-4">
                    <dt class="font-medium text-violet-700 dark:text-violet-500 mb-2" v-html="highlightText(getLabel(item))"></dt>
                    <dd class="mt-2 text-sm text-gray-600 dark:text-gray-400 leading-relaxed" 
                        style="line-height: 1.8; word-break: keep-all; overflow-wrap: break-word; white-space: pre-wrap;"
                        v-html="highlightText(formatDescription(getDescription(item)))"></dd>
                    <img v-if="item.image" :src="item.image" :alt="getLabel(item)" class="mt-3 max-w-full h-auto rounded-lg border border-gray-200 dark:border-gray-600" />
                  </div>
                </dl>
              </section>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { equipmentData, pqData, faultData, eventData } from './glossaryData'

export default {
  name: 'GlossaryModal',
  props: {
    isOpen: Boolean
  },
  emits: ['close'],
  setup(props, { emit }) {
    const { locale } = useI18n()
    const isKorean = computed(() => locale.value === 'ko')
    const isJapanese = computed(() => locale.value === 'ja')
    const activeTab = ref('equipment')
    const activeItem = ref('')
    const expandedCategories = ref(equipmentData.map(c => c.key))
    const contentArea = ref(null)
    const searchQuery = ref('')
    const currentSearchIndex = ref(0)

    const getLabel = (item) => {
      if (locale.value === 'ko' && item.label_ko) return item.label_ko
      if (locale.value === 'ja' && (item.label_jp || item.label_ja)) return item.label_jp || item.label_ja
      return item.label
    }

    const getDescription = (item) => {
      if (locale.value === 'ko' && item.description_ko) return item.description_ko
      if (locale.value === 'ja' && (item.description_jp || item.description_ja)) return item.description_jp || item.description_ja
      return item.description
    }

    const highlightText = (text) => {
      if (!searchQuery.value || !text) return text
      const query = searchQuery.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
      return text.replace(new RegExp(`(${query})`, 'gi'), '<mark class="bg-yellow-300 dark:bg-yellow-600 px-0.5 rounded">$1</mark>')
    }

    const formatDescription = (text) => {
      if (!text) return ''
      
      let formatted = text
      
      // 줄바꿈이 있으면 처리
      if (formatted.includes('\n')) {
        const lines = formatted.split('\n')
        const processedLines = []
        let isFirstParagraph = true
        
        for (let i = 0; i < lines.length; i++) {
          let line = lines[i].trim()
          
          // 빈 줄은 건너뛰기
          if (!line) {
            continue
          }
          
          // "원인" 또는 "영향" 라인 처리 (연한 보라색으로 표시, "-" 없음)
          if (line === '원인' || line === '영향' || line === 'Causes' || line === 'Effects' || line === 'Cause' || line === 'Effect') {
            processedLines.push(`<div class="mb-2 mt-3"><span class="text-violet-400 dark:text-violet-300 font-semibold">${line}</span></div>`)
            isFirstParagraph = false
            continue
          }
          
          // 첫 번째 문단 처리 (정의 부분, 일반 텍스트로 표시)
          if (isFirstParagraph) {
            processedLines.push(`<div class="mb-2" style="word-break: keep-all; overflow-wrap: break-word;">${line}</div>`)
            isFirstParagraph = false
            continue
          }
          
          // 불릿 포인트 라인 처리 (라인 시작이 • 또는 - 또는 *)
          if (line.match(/^[•\-\*]\s+/)) {
            const content = line.replace(/^[•\-\*]\s+/, '')
            processedLines.push(`<div class="ml-4 mb-2 flex items-start"><span class="text-violet-600 dark:text-violet-400 mr-2 flex-shrink-0 font-bold">•</span><span class="flex-1">${content}</span></div>`)
          } else {
            // 두 번째부터는 "-" 추가
            processedLines.push(`<div class="ml-4 mb-2 flex items-start"><span class="text-gray-600 dark:text-gray-400 mr-2 flex-shrink-0">-</span><span class="flex-1" style="word-break: keep-all; overflow-wrap: break-word;">${line}</span></div>`)
          }
        }
        
        return processedLines.join('')
      }
      
      // 줄바꿈이 없으면 자동 포맷팅 (기존 설명이 줄바꿈 없는 경우)
      formatted = formatted
        .replace(/([.!?])\s+([A-Z가-힣][^.!?]{20,})/g, '$1<br><br>$2')
        .replace(/([.!?])\s+([A-Z가-힣])/g, '$1<br>$2')
        .replace(/\s+(However|Additionally|Furthermore|Moreover|Therefore|Thus|그러나|하지만|또한|따라서|그러므로)/gi, '<br><br>$1')
        .replace(/(<br>\s*){3,}/gi, '<br><br>')
        .trim()
      
      return formatted
    }

    const tabs = computed(() => [
      { key: 'equipment', label: isKorean.value ? '장비' : isJapanese.value ? '設備' : 'Equipment' },
      { key: 'pq',        label: isKorean.value ? '전력품질' : isJapanese.value ? '電力品質' : 'Power Quality' },
      { key: 'fault',     label: isKorean.value ? '결함' : isJapanese.value ? '故障' : 'Fault' },
      { key: 'event',     label: isKorean.value ? '이벤트' : isJapanese.value ? 'イベント' : 'Event' },
    ])

    const allItems = computed(() => {
      const items = []
      equipmentData.forEach(category => {
        category.items.forEach(item => {
          items.push({ id: `equipment-${category.key}-${item.key}`, ...item })
        })
      })
      pqData.forEach(item => items.push({ id: `pq-${item.key}`, ...item }))
      faultData.forEach(item => items.push({ id: `fault-${item.key}`, ...item }))
      eventData.forEach(item => items.push({ id: `event-${item.key}`, ...item }))
      return items
    })

    const searchResults = computed(() => {
      if (!searchQuery.value) return []
      const q = searchQuery.value.toLowerCase()
      return allItems.value.filter(item => 
        item.label.toLowerCase().includes(q) ||
        (item.label_ko && item.label_ko.toLowerCase().includes(q)) ||
        (item.label_ja && item.label_ja.toLowerCase().includes(q)) ||
        item.description.toLowerCase().includes(q) ||
        (item.description_ko && item.description_ko.toLowerCase().includes(q)) ||
        (item.description_ja && item.description_ja.toLowerCase().includes(q))
      )
    })

    watch(searchQuery, (newVal) => {
      currentSearchIndex.value = 0
      if (newVal && searchResults.value.length > 0) {
        goToResult(0)
      }
    })

    const goToResult = (index) => {
      if (searchResults.value.length === 0) return
      currentSearchIndex.value = index
      const result = searchResults.value[index]
      document.getElementById(result.id)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }

    const goToNext = () => {
      if (searchResults.value.length === 0) return
      goToResult((currentSearchIndex.value + 1) % searchResults.value.length)
    }

    const goToPrev = () => {
      if (searchResults.value.length === 0) return
      goToResult(currentSearchIndex.value - 1 < 0 ? searchResults.value.length - 1 : currentSearchIndex.value - 1)
    }

    const clearSearch = () => {
      searchQuery.value = ''
      currentSearchIndex.value = 0
    }

    const currentTabItems = computed(() => {
      switch (activeTab.value) {
        case 'pq': return pqData
        case 'fault': return faultData
        case 'event': return eventData
        default: return []
      }
    })

    const toggleCategory = (key) => {
      const idx = expandedCategories.value.indexOf(key)
      if (idx > -1) expandedCategories.value.splice(idx, 1)
      else expandedCategories.value.push(key)
    }

    const scrollToItem = (id) => {
      activeItem.value = id
      document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }

    const close = () => {
      activeItem.value = ''
      searchQuery.value = ''
      currentSearchIndex.value = 0
      emit('close')
    }

    return {
      tabs, activeTab, activeItem, expandedCategories, contentArea, searchQuery,
      currentSearchIndex, searchResults, isKorean, isJapanese, getLabel, getDescription, highlightText, formatDescription,
      equipmentData, pqData, faultData, eventData, currentTabItems,
      goToNext, goToPrev, clearSearch, toggleCategory, scrollToItem, close,
    }
  }
}
</script>