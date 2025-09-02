<template>
  <div class="border-b border-gray-200 dark:border-gray-700/60 bg-white dark:bg-gray-800">
    <div class="px-6 py-4">
      <!-- Group Title -->
      <!-- <div class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase mb-4">
        {{ t('User.sidebar.title') }}
      </div> -->
      
      <!-- Horizontal Tab Navigation -->
      <nav class="flex space-x-8" aria-label="Tabs">
        <!-- Profile Tab -->
        <router-link to="/settings/Account/Profile" custom v-slot="{ href, navigate, isExactActive }">
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
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M8 9a4 4 0 1 1 0-8 4 4 0 0 1 0 8Zm0-2a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm-5.143 7.91a1 1 0 1 1-1.714-1.033A7.996 7.996 0 0 1 8 10a7.996 7.996 0 0 1 6.857 3.877 1 1 0 1 1-1.714 1.032A5.996 5.996 0 0 0 8 12a5.996 5.996 0 0 0-5.143 2.91Z" />
            </svg>
            <span class="whitespace-nowrap">{{ t('User.sidebar.profile') }}</span>
          </a>
        </router-link>

        <!-- User Management Tab -->
        <router-link to="/settings/Account/User" custom v-slot="{ href, navigate, isExactActive }">
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
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M5 9a1 1 0 1 1 0-2h6a1 1 0 0 1 0 2H5ZM1 4a1 1 0 1 1 0-2h14a1 1 0 0 1 0 2H1Zm0 10a1 1 0 0 1 0-2h14a1 1 0 0 1 0 2H1Z" />
            </svg>
            <span class="whitespace-nowrap">{{ t('User.sidebar.user') }}</span>
          </a>
        </router-link>

        <!-- API User Tab (conditional) -->
        <router-link 
          v-if="api && isNtek" 
          to="/settings/Account/APIUser" 
          custom 
          v-slot="{ href, navigate, isExactActive }"
        >
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
              width="16" 
              height="16" 
              viewBox="0 0 16 16"
            >
              <path d="M6.649 1.018a1 1 0 0 1 .793 1.171L6.997 4.5h3.464l.517-2.689a1 1 0 1 1 1.964.378L12.498 4.5h2.422a1 1 0 0 1 0 2h-2.807l-.77 4h2.117a1 1 0 1 1 0 2h-2.501l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H5.46l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H1a1 1 0 1 1 0-2h2.807l.77-4H2.46a1 1 0 0 1 0-2h2.5l.518-2.689a1 1 0 0 1 1.17-.793ZM9.307 10.5l.77-4H6.612l-.77 4h3.464Z" />
            </svg>
            <span class="whitespace-nowrap">{{ t('User.sidebar.api') }}</span>
          </a>
        </router-link>
      </nav>
    </div>
  </div>
</template>

<script>
import { ref, computed } from "vue";
import { useI18n } from 'vue-i18n'
import { useAuthStore } from "@/store/auth";

export default {
  name: 'SettingsSidebar',
  props: ["api"],
  setup(props) {
    const api = ref(props.api);
    const authStore = useAuthStore();
    const { t } = useI18n();
    const isAdmin = computed(() => (authStore.getUserRole == '2' || authStore.getUserRole == '3') ? true : false);
    const isNtek = computed(() => {
      const userName = authStore.getUser;
      if (userName == 'ntek' && isAdmin.value)
        return true;
      else
        return false;
    });

    return {
      api,
      t,
      isNtek,
      isAdmin,
    }
  }
}
</script>