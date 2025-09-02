<template>
    <div class="min-w-fit">
      <!-- Sidebar backdrop (mobile only) -->
      <div class="fixed inset-0 bg-gray-900 bg-opacity-30 z-40 lg:hidden lg:z-auto transition-opacity duration-200" :class="sidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'" aria-hidden="true"></div>
  
      <!-- Sidebar -->
      <div
        id="sidebar"
        ref="sidebar"
        class="flex lg:!flex flex-col absolute z-40 left-0 top-0 lg:static lg:left-auto lg:top-auto lg:translate-x-0 h-[100dvh] overflow-y-scroll lg:overflow-y-auto no-scrollbar w-64 lg:w-20 lg:sidebar-expanded:!w-64 2xl:!w-64 shrink-0 bg-white dark:bg-gray-800 p-4 transition-all duration-200 ease-in-out"
        :class="[variant === 'v2' ?  'border-r border-gray-200 dark:border-gray-700/60' : 'rounded-r-2xl shadow-sm', sidebarOpen ? 'translate-x-0' : '-translate-x-64']"
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
            <svg class="w-6 h-6 fill-current" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path d="M10.7 18.7l1.4-1.4L7.8 13H20v-2H7.8l4.3-4.3-1.4-1.4L4 12z" />
            </svg>
          </button>
          <!-- Logo -->
          <router-link class="block" to="/">
            <!--svg class="fill-violet-500" xmlns="http://www.w3.org/2000/svg" width="32" height="32">
              <path d="M31.956 14.8C31.372 6.92 25.08.628 17.2.044V5.76a9.04 9.04 0 0 0 9.04 9.04h5.716ZM14.8 26.24v5.716C6.92 31.372.63 25.08.044 17.2H5.76a9.04 9.04 0 0 1 9.04 9.04Zm11.44-9.04h5.716c-.584 7.88-6.876 14.172-14.756 14.756V26.24a9.04 9.04 0 0 1 9.04-9.04ZM.044 14.8C.63 6.92 6.92.628 14.8.044V5.76a9.04 9.04 0 0 1-9.04 9.04H.044Z" />
            </svg-->
            <!--img src="../images/logo2.png" alt="LOGO" width="64" height="64" /-->
            <img :src="logoSrc" alt="LOGO" width="132" height="132" />
          </router-link>
        </div>
  
        <!-- Links -->
        <div class="space-y-8">
          <!-- Pages group -->
          <div>
            <h3 class="text-xs uppercase text-gray-400 dark:text-gray-500 font-semibold pl-3">
              <span class="hidden lg:block lg:sidebar-expanded:hidden 2xl:hidden text-center w-6" aria-hidden="true">•••</span>
            </h3>
            <ul class="mt-3">
              <!-- Dashboard -->
              <router-link to="/dashboard" custom v-slot="{ href, navigate, isExactActive }">
              <li class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]" :class="isExactActive && 'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="isExactActive ? '' : 'hover:text-gray-900 dark:hover:text-white'" :href="href" @click="navigate">
                  <div class="flex items-center">
                    <svg class="shrink-0 fill-current" :class="currentRoute.fullPath === '/' || currentRoute.fullPath.includes('dashboard') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                      <path d="M5.936.278A7.983 7.983 0 0 1 8 0a8 8 0 1 1-8 8c0-.722.104-1.413.278-2.064a1 1 0 1 1 1.932.516A5.99 5.99 0 0 0 2 8a6 6 0 1 0 6-6c-.53 0-1.045.076-1.548.21A1 1 0 1 1 5.936.278Z" />
                      <path d="M6.068 7.482A2.003 2.003 0 0 0 8 10a2 2 0 1 0-.518-3.932L3.707 2.293a1 1 0 0 0-1.414 1.414l3.775 3.775Z" />
                    </svg> 
                    <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.dashboard') }}</span>
                  </div>
                </a>
              </li>
            </router-link>
              <SidebarLinkGroup v-if="setupMenu.MainEnable" v-slot="parentLink" :activeCondition="isMainActive">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="isMainActive ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                      <svg class="shrink-0 fill-current" :class="isMainActive ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                      <path d="M9 6.855A3.502 3.502 0 0 0 8 0a3.5 3.5 0 0 0-1 6.855v1.656L5.534 9.65a3.5 3.5 0 1 0 1.229 1.578L8 10.267l1.238.962a3.5 3.5 0 1 0 1.229-1.578L9 8.511V6.855ZM6.5 3.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0Zm4.803 8.095c.005-.005.01-.01.013-.016l.012-.016a1.5 1.5 0 1 1-.025.032ZM3.5 11c.474 0 .897.22 1.171.563l.013.016.013.017A1.5 1.5 0 1 1 3.5 11Z" />
                    </svg> 
                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.main') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link to="/meter/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.meter') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link to="/powerq/Main"  v-if="setupMenu.MainPowerQuality" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.pq') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link to="/diagnosis/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">Power Quality Diagnosis</span>
                        </a>
                      </li>
                    </router-link>                                
                  </ul>
                </div>
              </SidebarLinkGroup>            
              <!-- E-Commerce  -->
              <router-link to="/Sub" custom v-slot="{ href, navigate, isExactActive }">
              <li class="pl-4 pr-3 py-2 rounded-lg mb-0.5 last:mb-0 bg-[linear-gradient(135deg,var(--tw-gradient-stops))]" :class="isExactActive && 'from-violet-500/[0.12] dark:from-violet-500/[0.24] to-violet-500/[0.04]'">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="isExactActive ? '' : 'hover:text-gray-900 dark:hover:text-white'" :href="href" @click="navigate">
                  <div class="flex items-center">
                    <svg class="shrink-0 fill-current" :class="isSubActive ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                        <path d="M12 1a1 1 0 1 0-2 0v2a3 3 0 0 0 3 3h2a1 1 0 1 0 0-2h-2a1 1 0 0 1-1-1V1ZM1 10a1 1 0 1 0 0 2h2a1 1 0 0 1 1 1v2a1 1 0 1 0 2 0v-2a3 3 0 0 0-3-3H1ZM5 0a1 1 0 0 1 1 1v2a3 3 0 0 1-3 3H1a1 1 0 0 1 0-2h2a1 1 0 0 0 1-1V1a1 1 0 0 1 1-1ZM12 13a1 1 0 0 1 1-1h2a1 1 0 1 0 0-2h-2a3 3 0 0 0-3 3v2a1 1 0 1 0 2 0v-2Z" />
                      </svg>
                    <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">Sub Channel</span>
                  </div>
                </a>
              </li>
            </router-link>
              <!-- 알람/이벤트 -->
              <SidebarLinkGroup v-slot="parentLink" :activeCondition="currentRoute.fullPath.includes('event')">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="currentRoute.fullPath.includes('event') ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                      <!--svg class="shrink-0" 
                          :class="currentRoute.fullPath.includes('event') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'"
                          xmlns="http://www.w3.org/2000/svg" 
                          width="16" 
                          height="16" 
                          viewBox="0 0 24 24" 
                          fill="none" 
                          stroke="currentColor" 
                          stroke-linecap="round" 
                          stroke-linejoin="round" 
                          stroke-width="3">
                        <path d="M10 5a2 2 0 1 1 4 0a7 7 0 0 1 4 6v3a4 4 0 0 0 2 3H4a4 4 0 0 0 2-3v-3a7 7 0 0 1 4-6"/>
                        <path d="M9 17v1a3 3 0 0 0 6 0v-1"/>
                      </svg-->
                      <svg
                          :class="currentRoute.fullPath.includes('event') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'"
                          width="16" 
                          height="16"
                          fill="none"  
                      viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <!-- 🔔 벨 아이콘 -->
                      <path
                        d="M12 2C9.5 2 7.5 4 7.5 6.5V10c0 1.2-.4 2.4-1.2 3.3L5 14.5c-1 1 0 2.5 1.3 2.5h11.4c1.3 0 2.3-1.5 1.3-2.5l-1.3-1.2c-.8-.9-1.2-2.1-1.2-3.3V6.5C16.5 4 14.5 2 12 2z"
                        stroke="currentColor"
                        stroke-width="2.8"
                        fill="none"
                      />

                      <!-- 🔔 벨의 하단 둥근 부분 -->
                      <path
                        d="M10 19a2 2 0 0 0 4 0"
                        stroke="currentColor"
                        stroke-width="2.8"
                        fill="none"
                      />
                    </svg>
                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.event') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link v-if="setupMenu.MainEnable" to="/event/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.main') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link v-if="setupMenu.SubEnable" to="/event/Sub" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.sub') }}</span>
                        </a>
                      </li>
                    </router-link>             
                  </ul>
                </div>
              </SidebarLinkGroup>    
              <!-- 리포트 -->
              <SidebarLinkGroup v-slot="parentLink" :activeCondition="currentRoute.fullPath.includes('report')">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="currentRoute.fullPath.includes('report') ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                      <!--svg class="shrink-0 fill-current" :class="currentRoute.fullPath.includes('report') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                        <path d="M6.649 1.018a1 1 0 0 1 .793 1.171L6.997 4.5h3.464l.517-2.689a1 1 0 1 1 1.964.378L12.498 4.5h2.422a1 1 0 0 1 0 2h-2.807l-.77 4h2.117a1 1 0 1 1 0 2h-2.501l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H5.46l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H1a1 1 0 1 1 0-2h2.807l.77-4H2.46a1 1 0 0 1 0-2h2.5l.518-2.689a1 1 0 0 1 1.17-.793ZM9.307 10.5l.77-4H6.612l-.77 4h3.464Z" />
                      </svg-->
                      <svg class="shrink-0 fill-current" :class="currentRoute.fullPath.includes('report') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" width="16" height="16" viewBox="0 0 16 16">
                          <path d="M14 0H2c-.6 0-1 .4-1 1v14c0 .6.4 1 1 1h8l5-5V1c0-.6-.4-1-1-1zM3 2h10v8H9v4H3V2z" />
                      </svg>
                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.report') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link v-if="setupMenu.MainEnable" to="/report/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.main') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link v-if="setupMenu.SubEnable" to="/report/Sub" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.sub') }}</span>
                        </a>
                      </li>
                    </router-link>             
                  </ul>
                </div>
              </SidebarLinkGroup>
              <!-- 트렌드 -->
              <SidebarLinkGroup v-slot="parentLink" :activeCondition="currentRoute.fullPath.includes('trend')">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="currentRoute.fullPath.includes('trend') ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                      <!--svg class="shrink-0 fill-current" :class="currentRoute.fullPath.includes('report') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                        <path d="M6.649 1.018a1 1 0 0 1 .793 1.171L6.997 4.5h3.464l.517-2.689a1 1 0 1 1 1.964.378L12.498 4.5h2.422a1 1 0 0 1 0 2h-2.807l-.77 4h2.117a1 1 0 1 1 0 2h-2.501l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H5.46l-.517 2.689a1 1 0 1 1-1.964-.378l.444-2.311H1a1 1 0 1 1 0-2h2.807l.77-4H2.46a1 1 0 0 1 0-2h2.5l.518-2.689a1 1 0 0 1 1.17-.793ZM9.307 10.5l.77-4H6.612l-.77 4h3.464Z" />
                      </svg-->
                      <!--svg class="shrink-0 fill-current" :class="currentRoute.fullPath.includes('trend') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" width="16" height="16" viewBox="0 0 16 16">
                          <path d="M14 0H2c-.6 0-1 .4-1 1v14c0 .6.4 1 1 1h8l5-5V1c0-.6-.4-1-1-1zM3 2h10v8H9v4H3V2z" />
                      </svg-->
                      <svg
                        xmlns="http://www.w3.org/2000/svg" :class="currentRoute.fullPath.includes('trend') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'"
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

                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.trend') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link v-if="setupMenu.MainEnable" to="/trend/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.main') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link v-if="setupMenu.SubEnable" to="/trend/Sub" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.sub') }}</span>
                        </a>
                      </li>
                    </router-link>             
                  </ul>
                </div>
              </SidebarLinkGroup>
              <!-- 세팅 -->
              <SidebarLinkGroup v-if="isAdmin" v-slot="parentLink" :activeCondition="currentRoute.fullPath.includes('settings')">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="currentRoute.fullPath.includes('settings') ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                      <svg class="shrink-0 fill-current" :class="currentRoute.fullPath.includes('settings') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                        <path d="M10.5 1a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2h-1.145a3.502 3.502 0 0 1-6.71 0H1a1 1 0 0 1 0-2h6.145A3.502 3.502 0 0 1 10.5 1ZM9 4.5a1.5 1.5 0 1 1 3 0 1.5 1.5 0 0 1-3 0ZM5.5 9a3.502 3.502 0 0 1 3.355 2.5H15a1 1 0 1 1 0 2H8.855a3.502 3.502 0 0 1-6.71 0H1a1 1 0 1 1 0-2h1.145A3.502 3.502 0 0 1 5.5 9ZM4 12.5a1.5 1.5 0 1 0 3 0 1.5 1.5 0 0 0-3 0Z" fill-rule="evenodd" />
                      </svg>                    
                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.setup') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link to="/settings/Device/General" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.general') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link to="/settings/Device/Main" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.main') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link to="/settings/Device/Sub" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.sub') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link to="/settings/Device/Command" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">System</span>
                        </a>
                      </li>
                    </router-link>              
                  </ul>
                </div>
              </SidebarLinkGroup>
              <SidebarLinkGroup v-if="isAdmin && isNtek" v-slot="parentLink" :activeCondition="currentRoute.fullPath.includes('config')">
                <a class="block text-gray-800 dark:text-gray-100 truncate transition" :class="currentRoute.fullPath.includes('config') ? '' : 'hover:text-gray-900 dark:hover:text-white'" href="#0" @click.prevent="parentLink.handleClick(); sidebarExpanded = true">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2.2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="shrink-0" :class="currentRoute.fullPath.includes('config') ? 'text-violet-500' : 'text-gray-400 dark:text-gray-500'"
                        >
                          <circle cx="12" cy="12" r="3" />
                          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
                        </svg>

                                          
                      <span class="text-sm font-medium ml-4 lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.config') }}</span>
                    </div>
                    <!-- Icon -->
                    <div class="flex shrink-0 ml-2">
                      <svg class="w-3 h-3 shrink-0 ml-1 fill-current text-gray-400 dark:text-gray-500" :class="parentLink.expanded && 'rotate-180'" viewBox="0 0 12 12">
                        <path d="M5.9 11.4L.5 6l1.4-1.4 4 4 4-4L11.3 6z" />
                      </svg>
                    </div>
                  </div>
                </a>
                <div class="lg:hidden lg:sidebar-expanded:block 2xl:block">
                  <ul class="pl-8 mt-1" :class="!parentLink.expanded && 'hidden'">
                    <router-link to="/config/Service" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.service') }}</span>
                        </a>
                      </li>
                    </router-link>
                    <router-link v-if="devMode != 'device0'" to="/config/Calibrate" custom v-slot="{ href, navigate, isExactActive }">
                      <li class="mb-1 last:mb-0">
                        <a class="block transition truncate" :class="isExactActive ? 'text-violet-500' : 'text-gray-500/90 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'" :href="href" @click="navigate">
                          <span class="text-sm font-medium lg:opacity-0 lg:sidebar-expanded:opacity-100 2xl:opacity-100 duration-200">{{ t('sidebar.calibration') }}</span>
                        </a>
                      </li>
                    </router-link>            
                  </ul>
                </div>
              </SidebarLinkGroup>                      
              <!-- Messages -->
              <!-- Inbox -->       
            </ul>
          </div>
        </div>
  
        <!-- Expand / collapse button -->
        <div class="pt-3 hidden lg:inline-flex 2xl:hidden justify-end mt-auto">
          <div class="w-12 pl-4 pr-3 py-2">
            <button class="text-gray-400 hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400" @click.prevent="sidebarExpanded = !sidebarExpanded">
              <span class="sr-only">Expand / collapse sidebar</span>
              <svg class="shrink-0 fill-current text-gray-400 dark:text-gray-500 sidebar-expanded:rotate-180" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16">
                  <path d="M15 16a1 1 0 0 1-1-1V1a1 1 0 1 1 2 0v14a1 1 0 0 1-1 1ZM8.586 7H1a1 1 0 1 0 0 2h7.586l-2.793 2.793a1 1 0 1 0 1.414 1.414l4.5-4.5A.997.997 0 0 0 12 8.01M11.924 7.617a.997.997 0 0 0-.217-.324l-4.5-4.5a1 1 0 0 0-1.414 1.414L8.586 7M12 7.99a.996.996 0 0 0-.076-.373Z" />
              </svg>
            </button>
          </div>
        </div>
  
      </div>
    </div>
  </template>
  
  <script>
  import { useRouter, useRoute } from 'vue-router'
  import SidebarLinkGroup from './SidebarLinkGroup.vue'
  import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
  import { useAuthStore } from '@/store/auth'; // ✅ Pinia Store 사용
  import { useSetupStore } from '@/store/setup'; // ✅ Pinia Store 사용
  import { useDark } from '@vueuse/core'
  import { useI18n } from 'vue-i18n'  // ✅ 추가
  import logoLight from '@/images/CI_logo.png'
 import logoDark from '@/images/ci_logo_light.png'
  
  export default {
    name: 'Sidebar',
    props: [
      'sidebarOpen',
      'variant',
      'channel'
    ],
    components: {
      SidebarLinkGroup,
    },  
    setup(props, { emit }) {
      const isDark = useDark({
        selector: 'html',
        attribute: 'class',
      });
      const { t, locale } = useI18n();      
      const trigger = ref(null)
      const sidebar = ref(null)
      const authStore = useAuthStore(); // ✅ Pinia Store 사용
      const setupStore = useSetupStore(); // ✅ Pinia Store 사용
      const setupMenu = ref({});
      const logoSrc = computed(() => isDark.value ? logoDark : logoLight)

  
      const storedSidebarExpanded = localStorage.getItem('sidebar-expanded')
      const sidebarExpanded = ref(storedSidebarExpanded === null ? false : storedSidebarExpanded === 'true')
  
      const currentRoute = useRouter().currentRoute.value

      const setupEnabled = computed(() => setupStore.getSetup); // ✅ Vuex getters 변경
      const channelStatus = computed(() => setupStore.getChannelSetting); // ✅ Vuex getters 변경
      //const hidesetting = computed(() => authStore.getUser =='admin'?true:false); // ✅ Vuex getters 변경
      const isAdmin = computed(() => (authStore.getUserRole =='2' || authStore.getUserRole =='3')?true:false); // ✅ Vuex getters 변경
      //const installed  = computed(() => authStore.getInstalled > 1?true:false);
      const isNtek = computed(()=>{
        const userName = authStore.getUser;
        if (userName == 'ntek' && isAdmin.value)
          return true;
        else
          return false;
      })

      const devMode = computed(()=> authStore.getOpMode);

      watch(() => currentRoute.fullPath, (newPath, oldPath) => {
        //console.log(`Route changed: ${oldPath} -> ${newPath}`);
        }, { immediate: true }); 

      const isMainActive = computed(() => {
          const propdata = props.channel;
          //console.log(currentRoute.fullPath + ',' + propdata);
          if(currentRoute.fullPath.includes('/meter')){
            if(currentRoute.params === 'Main' || propdata === 'Main'){
              return true;
            }else{
              return false;
            }
          }else if(currentRoute.fullPath.includes('/powerq')){
            if(currentRoute.params === 'Main' || propdata === 'Main'){
              return true;
            }else{
              return false;
            }
          }else{
            return false;
          }
       }); // ✅ Vuex getters 변경

       const isSubActive = computed(() => {
          const propdata = props.channel;
          //(currentRoute.fullPath + ',' + propdata);
          if(currentRoute.fullPath.includes('/meter')){
            if(currentRoute.params === 'Sub' || propdata === 'Sub'){
              return true;
            }else{
              return false;
            }
          }else if(currentRoute.fullPath.includes('/powerq')){
            if(currentRoute.params === 'Sub' || propdata === 'Sub'){
              return true;
            }else{
              return false;
            }
          }else{
            return false;
          }
       }); // ✅ Vuex getters 변경
    
    onMounted(async () => {
      if (!setupStore.applysetup) {
         await setupStore.checkSetting(); // ✅ Vuex dispatch → Pinia action 사용
         //console.log("✅ CheckSetting 호출됨");
      }
      if (setupEnabled.value) {
        Object.assign(setupMenu.value, channelStatus.value); // ✅ 기존 객체에 병합하여 반응형 유지
        //console.log("✅ setupMenu 업데이트됨:", setupMenu.value);
      } else {
        console.warn("⚠️ Setup이 활성화되지 않음.");
      }
    });

    // ✅ 데이터 변경 감지하여 값이 바뀌면 업데이트

    watch(channelStatus, (newVal) => {
      //console.log("🚀 Channel Status 변경 감지됨:", newVal);
      Object.assign(setupMenu.value, newVal);
    });
  
      // close on click outside
      const clickHandler = ({ target }) => {
        if (!sidebar.value || !trigger.value) return
        if (
          !props.sidebarOpen ||
          sidebar.value.contains(target) ||
          trigger.value.contains(target)
        ) return
        emit('close-sidebar')
      }
  
      // close if the esc key is pressed
      const keyHandler = ({ keyCode }) => {
        if (!props.sidebarOpen || keyCode !== 27) return
        emit('close-sidebar')
      } 
  
      onMounted(() => {
        //console.log("authStore.getUserRole",authStore.getUserRole,"authStore.getUser",authStore.getUser)
        document.addEventListener('click', clickHandler)
        document.addEventListener('keydown', keyHandler)
      })
  
      onUnmounted(() => {
        document.removeEventListener('click', clickHandler)
        document.removeEventListener('keydown', keyHandler)
      })
  
      watch(sidebarExpanded, () => {
        localStorage.setItem('sidebar-expanded', sidebarExpanded.value)
        if (sidebarExpanded.value) {
          document.querySelector('body').classList.add('sidebar-expanded')
        } else {
          document.querySelector('body').classList.remove('sidebar-expanded')
        }
      })
  
      return {
        trigger,
        sidebar,
        sidebarExpanded,
        currentRoute,
        setupMenu,
        channelStatus,
        isAdmin,
        isMainActive,
        isSubActive,
        logoSrc,
        isDark,
        devMode,
        //installed,
        t,
        isNtek,
      }
    },
  }
  </script>