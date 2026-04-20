<template>
  <div class="border-b border-gray-200 dark:border-gray-700/60 bg-white dark:bg-gray-800">
    <div class="px-6 py-4">
      <!-- Group Title -->
      <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase mb-4">
        {{ t('sidebar.setup') }}
      </div>
      
      <!-- Horizontal Tab Navigation -->
      <nav class="flex space-x-8" aria-label="Device Settings Tabs">
        <!-- General Tab -->
        <router-link to="/settings/Device/General" custom v-slot="{ href, navigate, isExactActive }">
          <a 
            :href="href" 
            @click="navigate"
            :class="[
              'group inline-flex items-center py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
              isExactActive 
                ? 'border-violet-500 text-violet-600 dark:text-violet-400' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300 dark:hover:border-gray-500'
            ]"
          >
            <svg 
              class="shrink-0 fill-current mr-2 transition-colors duration-200" 
              :class="isExactActive ? 'text-violet-500 dark:text-violet-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400'" 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M10.5 1a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2h-1.145a3.502 3.502 0 0 1-6.71 0H1a1 1 0 0 1 0-2h6.145A3.502 3.502 0 0 1 10.5 1ZM9 4.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM5.5 9a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2H8.855a3.502 3.502 0 0 1-6.71 0H1a1 1 0 1 1 0-2h1.145A3.502 3.502 0 0 1 5.5 9ZM4 12.5a1.5 1.5 0 1 0 3 0 1.5 1.5 0 0 0-3 0Z" fill-rule="evenodd" />
            </svg>
            <span class="whitespace-nowrap">{{ t('sidebar.general') }}</span>
          </a>
        </router-link>

        <!-- Main Channel Tab -->
        <router-link to="/settings/Device/Main" custom v-slot="{ href, navigate, isExactActive }">
          <a 
            :href="href" 
            @click="navigate"
            :class="[
              'group inline-flex items-center py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
              isExactActive 
                ? 'border-violet-500 text-violet-600 dark:text-violet-400' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300 dark:hover:border-gray-500'
            ]"
          >
            <svg 
              class="shrink-0 fill-current mr-2 transition-colors duration-200" 
              :class="isExactActive ? 'text-violet-500 dark:text-violet-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400'" 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M9 6.855A3.502 3.502 0 0 0 8 0a3.5 3.5 0 0 0-1 6.855v1.656L5.534 9.65a3.5 3.5 0 1 0 1.229 1.578L8 10.267l1.238.962a3.5 3.5 0 1 0 1.229-1.578L9 8.511V6.855ZM6.5 3.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0Zm4.803 8.095c.005-.005.01-.01.013-.016l.012-.016a1.5 1.5 0 1 1-.025.032ZM3.5 11c.474 0 .897.22 1.171.563l.013.016.013.017A1.5 1.5 0 1 1 3.5 11Z" />
            </svg>
            <span class="whitespace-nowrap">{{ t('sidebar.main') }}</span>
          </a>
        </router-link>

        <!-- Sub Channel Tab -->
        <router-link to="/settings/Device/Sub" custom v-slot="{ href, navigate, isExactActive }">
          <a 
            :href="href" 
            @click="navigate"
            :class="[
              'group inline-flex items-center py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
              isExactActive 
                ? 'border-violet-500 text-violet-600 dark:text-violet-400' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300 dark:hover:border-gray-500'
            ]"
          >
            <svg 
              class="shrink-0 fill-current mr-2 transition-colors duration-200" 
              :class="isExactActive ? 'text-violet-500 dark:text-violet-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400'" 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M12 1a1 1 0 1 0-2 0v2a3 3 0 0 0 3 3h2a1 1 0 1 0 0-2h-2a1 1 0 0 1-1-1V1ZM1 10a1 1 0 1 0 0 2h2a1 1 0 0 1 1 1v2a1 1 0 1 0 2 0v-2a3 3 0 0 0-3-3H1ZM5 0a1 1 0 0 1 1 1v2a3 3 0 0 1-3 3H1a1 1 0 0 1 0-2h2a1 1 0 0 0 1-1V1a1 1 0 0 1 1-1ZM12 13a1 1 0 0 1 1-1h2a1 1 0 1 0 0-2h-2a3 3 0 0 0-3 3v2a1 1 0 1 0 2 0v-2Z" />
            </svg>
            <span class="whitespace-nowrap">{{ t('sidebar.sub') }}</span>
          </a>
        </router-link>

        <!-- System Tab -->
        <router-link to="/settings/Device/Command" custom v-slot="{ href, navigate, isExactActive }">
          <a 
            :href="href" 
            @click="navigate"
            :class="[
              'group inline-flex items-center py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
              isExactActive 
                ? 'border-violet-500 text-violet-600 dark:text-violet-400' 
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300 dark:hover:border-gray-500'
            ]"
          >
            <svg 
              class="shrink-0 fill-current mr-2 transition-colors duration-200" 
              :class="isExactActive ? 'text-violet-500 dark:text-violet-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400'" 
              xmlns="http://www.w3.org/2000/svg" 
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M10.5 1a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2h-1.145a3.502 3.502 0 0 1-6.71 0H1a1 1 0 0 1 0-2h6.145A3.502 3.502 0 0 1 10.5 1ZM9 4.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM5.5 9a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2H8.855a3.502 3.502 0 0 1-6.71 0H1a1 1 0 1 1 0-2h1.145A3.502 3.502 0 0 1 5.5 9ZM4 12.5a1.5 1.5 0 1 0 3 0 1.5 1.5 0 0 0-3 0Z" fill-rule="evenodd" />
            </svg>
            <span class="whitespace-nowrap">System</span>
          </a>
        </router-link>
      </nav>
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n'

export default {
  name: 'DeviceSettingsTabs',
  setup() {
    const { t } = useI18n();

    return {
      t,
    }
  }
}
</script>