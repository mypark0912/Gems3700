<template>
  <form  :class="['relative flex items-center', widthClass, 'self-start']">
    <label for="action-search" class="sr-only">Search</label>

    <!-- ✅ Input Field -->
    <input 
      id="action-search"
      class="form-input w-full pl-10 pr-10 py-2 bg-white dark:bg-gray-800 border rounded-md focus:ring-2 focus:ring-indigo-500"
      type="search"
      :placeholder="placeholder"
    />

    <!-- ✅ Search Button (아이콘) -->
    <button 
      class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-400"
      type="submit"
      aria-label="Search"
    >
      <svg class="w-5 h-5" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <path d="M7 14c-3.86 0-7-3.14-7-7s3.14-7 7-7 7.14 3.14 7 7-3.14 7-7 7zM7 2C4.243 2 2 4.243 2 7s2.243 5 5 5 5-2.243 5-5-2.243-5-5-5z" />
        <path d="M15.707 14.293L13.314 11.9a8.019 8.019 0 01-1.414 1.414l2.393 2.393a.997.997 0 001.414 0 .999.999 0 000-1.414z" />
      </svg>
    </button>
  </form>
</template>
<script>
export default {
  name: 'SearchNewForm',
  props: {
    placeholder: {
      type: String,
      default: 'Search…'
    },
    widthClass: {
      type: String,
      default: 'w-full' // 기본값은 전체 너비
    }
  }
}
</script>