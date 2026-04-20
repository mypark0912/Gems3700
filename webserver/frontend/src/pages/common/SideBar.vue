<template>
  <div class="min-w-fit">
    <!-- Sidebar backdrop (mobile only) -->
    <div
      class="fixed inset-0 bg-gray-900 bg-opacity-30 z-40 lg:hidden lg:z-auto transition-opacity duration-200"
      :class="sidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'"
      aria-hidden="true"
    ></div>

    <!-- Sidebar -->
    <div
      id="sidebar"
      ref="sidebar"
      class="flex lg:!flex flex-col absolute z-40 left-0 top-0 lg:static lg:left-auto lg:top-auto lg:translate-x-0 h-[100dvh] overflow-y-scroll lg:overflow-y-auto no-scrollbar w-64 lg:w-20 lg:sidebar-expanded:!w-64 2xl:!w-64 shrink-0 bg-white dark:bg-gray-800 p-4 transition-all duration-200 ease-in-out"
      :class="[
        variant === 'v2'
          ? 'border-r border-gray-200 dark:border-gray-700/60'
          : 'rounded-r-2xl shadow-sm',
        sidebarOpen ? 'translate-x-0' : '-translate-x-64',
      ]"
    >
      <!-- Sidebar header -->
      <div class="flex justify-between mb-10 pr-3 sm:px-2">
        <!-- Close button -->
        <button
          ref="trigger"
          class="lg:hidden text-gray-500 hover:text-gray-400"
          @click="$emit('close-sidebar')"
          aria-controls="sidebar"
          :aria-expanded="sidebarOpen"
        >
          <span class="sr-only">Close sidebar</span>
          <svg
            class="w-6 h-6 fill-current"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M10.7 18.7l1.4-1.4L7.8 13H20v-2H7.8l4.3-4.3-1.4-1.4L4 12z"
            />
          </svg>
        </button>
        <!-- Logo -->
        <router-link class="block" to="/">
          <img :src="logoSrc" alt="LOGO" width="132" height="132" />
        </router-link>
      </div>

      <!-- Links -->
      <div class="space-y-8">
        <!-- Pages group -->
        <div>
          <h3
            class="text-xs uppercase text-gray-400 dark:text-gray-500 font-semibold pl-3"
          >
            <span
              class="hidden lg:block lg:sidebar-expanded:hidden 2xl:hidden text-center w-6"
              aria-hidden="true"
              >•••</span
            >
          </h3>
          <ul class="mt-3">
            <!-- Dashboard -->
            <router-link
              to="/dashboard"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      class="shrink-0 fill-current"
                      :class="
                        currentRoute.fullPath === '/' ||
                        currentRoute.fullPath.includes('dashboard')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 16 16"
                    >
                      <path
                        d="M5.936.278A7.983 7.983 0 0 1 8 0a8 8 0 1 1-8 8c0-.722.104-1.413.278-2.064a1 1 0 1 1 1.932.516A5.99 5.99 0 0 0 2 8a6 6 0 1 0 6-6c-.53 0-1.045.076-1.548.21A1 1 0 1 1 5.936.278Z"
                      />
                      <path
                        d="M6.068 7.482A2.003 2.003 0 0 0 8 10a2 2 0 1 0-.518-3.932L3.707 2.293a1 1 0 0 0-1.414 1.414l3.775 3.775Z"
                      />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.dashboard") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- Meter -->
            <router-link
              to="/meter"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      class="shrink-0"
                      :class="
                        currentRoute.fullPath.includes('meter')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path d="M12 20V10" />
                      <path d="M18 20V4" />
                      <path d="M6 20v-4" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.meter") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- Power Quality -->
            <router-link
              to="/powerq"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      class="shrink-0"
                      :class="
                        currentRoute.fullPath.includes('powerq')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.pq") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- 브랜치 -->
            <SidebarLinkGroup
              v-slot="parentLink"
              :activeCondition="
                currentRoute.fullPath.includes('branch') ||
                currentRoute.fullPath.includes('ibsm') ||
                currentRoute.fullPath.includes('mcs')
              "
            >
              <a
                class="block text-gray-800 dark:text-gray-100 truncate transition"
                :class="
                  currentRoute.fullPath.includes('branch')
                    ? ''
                    : 'hover:text-gray-900 dark:hover:text-white'
                "
                href="#0"
                @click.prevent="
                  parentLink.handleClick();
                  sidebarExpanded = true;
                "
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <svg
                      class="shrink-0"
                      :class="
                        currentRoute.fullPath.includes('branch')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <rect x="4" y="4" width="6" height="6" rx="1" />
                      <rect x="14" y="4" width="6" height="6" rx="1" />
                      <rect x="4" y="14" width="6" height="6" rx="1" />
                      <rect x="14" y="14" width="6" height="6" rx="1" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t('sidebar.module') }}</span
                    >
                  </div>
                  <!-- Icon -->
                  <div class="flex shrink-0 ml-2">
                    <svg
                      class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500"
                      :class="parentLink.expanded && 'rotate-180'"
                      viewBox="0 0 12 12"
                    >
                      <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                    </svg>
                  </div>
                </div>
              </a>
              <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                  <router-link
                    to="/branch/1"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >IPSM72 #1</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    to="/branch/2"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >IPSM72 #2</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    to="/ibsm"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >IBSM</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    to="/mcs"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >MCS</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <!-- <router-link
  to="/ibsm2"
  custom
  v-slot="{ href, navigate, isExactActive }"
>
  <li class="mb-1 last:mb-0">
    <a
      class="block transition truncate"
      :class="
        isExactActive
          ? 'text-violet-500'
          : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
      "
      :href="href"
      @click="navigate"
    >
      <span
        class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
        >IBSM2</span
      >
    </a>
  </li>
</router-link> -->
                </ul>
              </div>
            </SidebarLinkGroup>
            <!-- 알람/이벤트 -->
            <router-link
              to="/event"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      :class="
                        currentRoute.fullPath.includes('event')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      width="16"
                      height="16"
                      fill="none"
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M12 2C9.5 2 7.5 4 7.5 6.5V10c0 1.2-.4 2.4-1.2 3.3L5 14.5c-1 1 0 2.5 1.3 2.5h11.4c1.3 0 2.3-1.5 1.3-2.5l-1.3-1.2c-.8-.9-1.2-2.1-1.2-3.3V6.5C16.5 4 14.5 2 12 2z"
                        stroke="currentColor"
                        stroke-width="2.8"
                        fill="none"
                      />
                      <path
                        d="M10 19a2 2 0 0 0 4 0"
                        stroke="currentColor"
                        stroke-width="2.8"
                        fill="none"
                      />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.event") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- 리포트 -->
            <router-link
              to="/report"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      class="shrink-0 fill-current"
                      :class="
                        currentRoute.fullPath.includes('report')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      width="16"
                      height="16"
                      viewBox="0 0 16 16"
                    >
                      <path
                        d="M14 0H2c-.6 0-1 .4-1 1v14c0 .6.4 1 1 1h8l5-5V1c0-.6-.4-1-1-1zM3 2h10v8H9v4H3V2z"
                      />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.report") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- 트렌드 -->
            <router-link
              to="/trend"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  isExactActive &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    isExactActive
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      :class="
                        currentRoute.fullPath.includes('trend')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      width="16"
                      height="16"
                      viewBox="0 0 16 16"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="1.5"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path d="M2 2v12h12" />
                      <path d="M13 12v2" />
                      <path d="M10.5 10v4" />
                      <path d="M8 8v6" />
                      <path d="M5.5 10v4" />
                      <path d="M2 7c4 0 3 -3.5 6 -3.5s2 3.5 6 3.5" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.trend") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- 설정 -->
            <router-link
              v-if="isAdmin"
              to="/settings/general"
              custom
              v-slot="{ href, navigate, isExactActive }"
            >
              <li
                class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]"
                :class="
                  currentRoute.fullPath.includes('settings') &&
                  'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'
                "
              >
                <a
                  class="block text-gray-800 dark:text-gray-100 truncate transition"
                  :class="
                    currentRoute.fullPath.includes('settings')
                      ? ''
                      : 'hover:text-gray-900 dark:hover:text-white'
                  "
                  :href="href"
                  @click="navigate"
                >
                  <div class="flex items-center">
                    <svg
                      class="shrink-0"
                      :class="currentRoute.fullPath.includes('settings') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'"
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2.2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <circle cx="12" cy="12" r="3" />
                      <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.setup") }}</span
                    >
                  </div>
                </a>
              </li>
            </router-link>
            <!-- 관리자 -->
            <SidebarLinkGroup
              v-if="isAdmin"
              v-slot="parentLink"
              :activeCondition="currentRoute.fullPath.includes('config')"
            >
              <a
                class="block text-gray-800 dark:text-gray-100 truncate transition"
                :class="
                  currentRoute.fullPath.includes('config')
                    ? ''
                    : 'hover:text-gray-900 dark:hover:text-white'
                "
                href="#0"
                @click.prevent="
                  parentLink.handleClick();
                  sidebarExpanded = true;
                "
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <svg
                      class="shrink-0 fill-current"
                      :class="
                        currentRoute.fullPath.includes('config')
                          ? 'text-violet-500'
                          : 'text-gray-400 dark:text-gray-500'
                      "
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 16 16"
                    >
                      <path d="M10.5 1a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2h-1.145a3.502 3.502 0 0 1-6.71 0H1a1 1 0 0 1 0-2h6.145A3.502 3.502 0 0 1 10.5 1ZM9 4.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM5.5 9a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2H8.855a3.502 3.502 0 0 1-6.71 0H1a1 1 0 1 1 0-2h1.145A3.502 3.502 0 0 1 5.5 9ZM4 12.5a1.5 1.5 0 1 0 3 0 1.5 1.5 0 0 0-3 0Z" />
                    </svg>
                    <span
                      class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                      >{{ t("sidebar.config") }}</span
                    >
                  </div>
                  <!-- Icon -->
                  <div class="flex shrink-0 ml-2">
                    <svg
                      class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500"
                      :class="parentLink.expanded && 'rotate-180'"
                      viewBox="0 0 12 12"
                    >
                      <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                    </svg>
                  </div>
                </div>
              </a>
              <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                  <router-link
                    to="/config/Calibrate"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >{{ t("sidebar.calibration") }}</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    to="/config/Service"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >{{ t("sidebar.service") }}</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    v-if="isNtek"
                    to="/config/Maintenance"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >{{ t("header.maintenance") }}</span
                        >
                      </a>
                    </li>
                  </router-link>
                  <router-link
                    v-if="isNtek"
                    to="/config/Log"
                    custom
                    v-slot="{ href, navigate, isExactActive }"
                  >
                    <li class="mb-1 last:mb-0">
                      <a
                        class="block transition truncate"
                        :class="
                          isExactActive
                            ? 'text-violet-500'
                            : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                        "
                        :href="href"
                        @click="navigate"
                      >
                        <span
                          class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200"
                          >User Action Log</span
                        >
                      </a>
                    </li>
                  </router-link>
                </ul>
              </div>
            </SidebarLinkGroup>
          </ul>
        </div>
      </div>

      <!-- Expand / collapse button -->
      <div class="pt-3 hidden lg:inline-flex 2xl:hidden justify-end mt-auto">
        <div class="w-12 pl-4 pr-3 py-2">
          <button
            class="text-gray-400 hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400"
            @click.prevent="sidebarExpanded = !sidebarExpanded"
          >
            <span class="sr-only">Expand / collapse sidebar</span>
            <svg
              class="shrink-0 fill-current text-gray-400 dark:text-gray-500 sidebar-expanded:rotate-180"
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 16 16"
            >
              <path
                d="M15 16a1 1 0 0 1-1-1V1a1 1 0 1 1 2 0v14a1 1 0 0 1-1 1ZM8.586 7H1a1 1 0 1 0 0 2h7.586l-2.793 2.793a1 1 0 1 0 1.414 1.414l4.5-4.5A.997.997 0 0 0 12 8.01M11.924 7.617a.997.997 0 0 0-.217-.324l-4.5-4.5a1 1 0 0 0-1.414 1.414L8.586 7M12 7.99a.996.996 0 0 0-.076-.373Z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import SidebarLinkGroup from "./SidebarLinkGroup.vue";
import { ref, onMounted, onUnmounted, watch, computed } from "vue";
import { useAuthStore } from "@/store/auth";
import { useSetupStore } from "@/store/setup";
import { useDark } from "@vueuse/core";
import { useI18n } from "vue-i18n";
import logoLight from "@/images/idpm300_logo_4.png";
import logoDark from "@/images/idpm300_logo_4.png";

export default {
  name: "Sidebar",
  props: ["sidebarOpen", "variant"],
  components: {
    SidebarLinkGroup,
  },
  setup(props, { emit }) {
    const isDark = useDark({
      selector: "html",
      attribute: "class",
    });
    const { t, locale } = useI18n();
    const trigger = ref(null);
    const sidebar = ref(null);
    const authStore = useAuthStore();
    const setupStore = useSetupStore();
    const logoSrc = computed(() => (isDark.value ? logoDark : logoLight));
    const storedSidebarExpanded = localStorage.getItem("sidebar-expanded");
    const sidebarExpanded = ref(
      storedSidebarExpanded === null ? false : storedSidebarExpanded === "true",
    );

    const currentRoute = useRouter().currentRoute.value;

    const isAdmin = computed(() =>
      authStore.getUserRole == "2" || authStore.getUserRole == "3"
        ? true
        : false,
    );
    const isNtek = computed(() => {
      const userName = authStore.getUser;
      if (userName == "ntek" && isAdmin.value) return true;
      else return false;
    });

    const devMode = computed(() => authStore.getOpMode);

    onMounted(async () => {
      if (!setupStore.applysetup) {
        await setupStore.checkSetting();
      }
    });

    // close on click outside
    const clickHandler = ({ target }) => {
      if (!sidebar.value || !trigger.value) return;
      if (
        !props.sidebarOpen ||
        sidebar.value.contains(target) ||
        trigger.value.contains(target)
      )
        return;
      emit("close-sidebar");
    };

    // close if the esc key is pressed
    const keyHandler = ({ keyCode }) => {
      if (!props.sidebarOpen || keyCode !== 27) return;
      emit("close-sidebar");
    };

    onMounted(() => {
      document.addEventListener("click", clickHandler);
      document.addEventListener("keydown", keyHandler);
    });

    onUnmounted(() => {
      document.removeEventListener("click", clickHandler);
      document.removeEventListener("keydown", keyHandler);
    });

    watch(sidebarExpanded, () => {
      localStorage.setItem("sidebar-expanded", sidebarExpanded.value);
      if (sidebarExpanded.value) {
        document.querySelector("body").classList.add("sidebar-expanded");
      } else {
        document.querySelector("body").classList.remove("sidebar-expanded");
      }
    });

    return {
      trigger,
      sidebar,
      sidebarExpanded,
      currentRoute,
      isAdmin,
      logoSrc,
      isDark,
      devMode,
      t,
      isNtek,
    };
  },
};
</script>
