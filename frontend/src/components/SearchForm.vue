<template>
  <form class="relative w-full flex items-center h-[28px]"> <!-- 높이 조정 -->
    <label for="action-search" class="sr-only">Search</label>

    <!-- ✅ Input Field -->
    <input 
      id="action-search"
      class="form-input w-full pl-8 pr-8 py-1 text-xs h-full bg-white dark:bg-gray-800 border rounded-md focus:ring-1 focus:ring-indigo-400"
      type="search"
      :placeholder="placeholder"
    />

    <!-- ✅ Search Icon Button -->
    <button 
      class="absolute right-1 top-1/2 -translate-y-1/2 p-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-400"
      aria-label="Search"
    >
      <svg class="w-4 h-4" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <path d="M7 14c-3.86 0-7-3.14-7-7s3.14-7 7-7 7.14 3.14 7 7-3.14 7-7 7zM7 2C4.243 2 2 4.243 2 7s2.243 5 5 5 5-2.243 5-5-2.243-5-5-5z" />
        <path d="M15.707 14.293L13.314 11.9a8.019 8.019 0 01-1.414 1.414l2.393 2.393a.997.997 0 001.414 0 .999.999 0 000-1.414z" />
      </svg>
    </button>
  </form>
</template>

<script>
export default {
  name: 'SearchForm',
  props: {
    placeholder: {
      type: String,
      default: 'Search…',
    },
  },
}
</script>
