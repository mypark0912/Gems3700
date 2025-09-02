<template>
  <div class="flex items-center justify-between w-full">
    <!-- 가운데 정렬된 페이지 정보 -->
    <div class="flex-1 text-center text-sm text-gray-500">
      Showing 
      <span class="font-medium text-gray-600 dark:text-white">{{ curPage }}</span> 
      to 
      <span class="font-medium text-gray-600 dark:text-white">{{ pages }}</span> 
      of 
      <span class="font-medium text-gray-600 dark:text-white">{{ record }}</span> 
      results
    </div>

    <!-- 오른쪽 정렬된 Prev/Next 버튼 -->
    <nav role="navigation" aria-label="Navigation">
      <ul class="flex">
        <li class="ml-3 first:ml-0">
          <span v-if="curPage == 1" class="btn bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 text-gray-300 dark:text-gray-600">
            &lt;- Previous
          </span>
          <a v-else @click="movePrev" class="btn bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white" href="#0">
            &lt;- Previous
          </a>
        </li>
        <li class="ml-3 first:ml-0">
          <span v-if="curPage == pages" class="btn bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 text-gray-300 dark:text-gray-600">
            Next -&gt;
          </span>
          <a v-else @click="moveNext" class="btn bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white" href="#0">
            Next -&gt;
          </a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script>
import { ref, inject } from 'vue';
export default {
  name: 'PaginationClassic',
  props:['channel'],
  setup(props){
    const record = inject('totalRecord');
    const pages = inject('totalPages');
    const channel = ref(props.channel);
    const curPage = inject('curPage');
    const moveNext = ()=>{
      curPage.value = curPage.value + 1;
    };

    const movePrev = ()=>{
      curPage.value = curPage.value - 1;
    };

    return {
      record,
      pages,
      channel,
      curPage,
      movePrev,
      moveNext,
    }
  }
}
</script>