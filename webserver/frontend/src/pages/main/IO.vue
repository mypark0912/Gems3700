<template>
    <div class="flex h-[100dvh] overflow-hidden">
      <!-- Sidebar -->
      <Sidebar :sidebarOpen="sidebarOpen" @close-sidebar="sidebarOpen = false" :channel="channel" />
  
      <!-- Content area -->
      <div class="relative flex flex-col flex-1 overflow-y-auto overflow-x-hidden bg-gray-50 dark:bg-gray-950">
        <!-- Site header -->
        <Header :sidebarOpen="sidebarOpen" @toggle-sidebar="sidebarOpen = !sidebarOpen" />
  
        <main class="grow">
          <div class="px-4 sm:px-6 lg:px-8 py-6 w-full">

            <!-- DI & DO Row (7:5 비율) -->
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-6">
              
              <!-- DI Section (7 columns) -->
              <div class="premium-card overflow-hidden lg:col-span-7">
                <!-- Header -->
                <div class="premium-header">
                  <div class="flex items-center justify-between">
                    <h2 class="gradient-title">Digital Input (DI)</h2>
                    <div class="flex items-center gap-4 text-xs">
                      <div class="flex items-center gap-1.5">
                        <div class="w-2 h-2 rounded-full bg-green-500"></div>
                        <span class="text-gray-600 dark:text-gray-400 font-medium">DI ON</span>
                      </div>
                      <div class="flex items-center gap-1.5">
                        <div class="w-2 h-2 rounded-full bg-gray-400"></div>
                        <span class="text-gray-600 dark:text-gray-400 font-medium">DI OFF</span>
                      </div>
                      <div class="flex items-center gap-1.5">
                        <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                        <span class="text-gray-600 dark:text-gray-400 font-medium">PI</span>
                      </div>
                    </div>
                  </div>
                </div>
                
                <!-- DI Grid -->
                <div class="p-6">
                  <div class="grid grid-cols-6 sm:grid-cols-8 md:grid-cols-12 gap-2">
                    <div v-for="ch in 24" :key="`di-${ch}`" 
                         class="flex flex-col items-center justify-center p-3 rounded-lg border-2 transition-all"
                         :class="getDICardClass(ch)">
                      <div class="text-lg font-bold mb-1">{{ ch }}</div>
                      <div class="text-xs font-semibold">{{ getDILabel(ch) }}</div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- DO Section (5 columns) -->
              <div class="premium-card overflow-hidden lg:col-span-5">
                <div class="premium-header">
                  <div class="flex items-center justify-between">
                    <h2 class="gradient-title">Digital Output (DO)</h2>
                    <div class="flex gap-2 text-xs">
                      <span class="px-2 py-1 rounded-md bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 font-medium">
                        Alarm
                      </span>
                      <span class="px-2 py-1 rounded-md bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 font-medium">
                        Event
                      </span>
                      <span class="px-2 py-1 rounded-md bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400 font-medium">
                        Output
                      </span>
                    </div>
                  </div>
                </div>
                
                <div class="p-7 flex items-center justify-center h-[calc(100%-4rem)]">
                  <div class="grid grid-cols-3 gap-5 w-full max-w-lg">
                    <!-- Module 1 -->
                    <div class="aspect-square rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 flex flex-col">
                      <div class="flex flex-col items-center justify-center p-3 bg-blue-50 dark:bg-blue-900/10 border-b border-blue-200 dark:border-blue-800 flex-1">
                        <span class="text-xs font-semibold text-blue-700 dark:text-blue-400 mb-2">Module 1</span>
                        <div class="flex items-center gap-2">
                          <div class="w-3 h-3 rounded-full bg-gray-400"></div>
                          <span class="text-xl font-bold text-gray-600 dark:text-gray-400">OFF</span>
                        </div>
                      </div>
                      <button @click="openDOControl(1, 'OFF')" 
                              class="w-full py-2 bg-white dark:bg-gray-900 text-blue-600 dark:text-blue-400 text-xs font-semibold hover:bg-blue-50 dark:hover:bg-blue-900/10 transition-colors">
                        제어
                      </button>
                    </div>

                    <!-- Module 2 -->
                    <div class="aspect-square rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 flex flex-col">
                      <div class="flex flex-col items-center justify-center p-3 bg-green-50 dark:bg-green-900/10 border-b border-green-200 dark:border-green-800 flex-1">
                        <span class="text-xs font-semibold text-green-700 dark:text-green-400 mb-2">Module 2</span>
                        <div class="flex items-center gap-2">
                          <div class="w-3 h-3 rounded-full bg-gray-400"></div>
                          <span class="text-xl font-bold text-gray-600 dark:text-gray-400">OFF</span>
                        </div>
                      </div>
                      <button @click="openDOControl(2, 'OFF')"
                              class="w-full py-2 bg-white dark:bg-gray-900 text-green-600 dark:text-green-400 text-xs font-semibold hover:bg-green-50 dark:hover:bg-green-900/10 transition-colors">
                        제어
                      </button>
                    </div>

                    <!-- Module 3 -->
                    <div class="aspect-square rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 flex flex-col">
                      <div class="flex flex-col items-center justify-center p-3 bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 flex-1">
                        <span class="text-xs font-semibold text-gray-700 dark:text-gray-300 mb-2">Module 3</span>
                        <div class="flex items-center gap-2">
                          <div class="w-3 h-3 rounded-full bg-green-500"></div>
                          <span class="text-xl font-bold text-green-600 dark:text-green-400">ON</span>
                        </div>
                      </div>
                      <button @click="openDOControl(3, 'ON')"
                              class="w-full py-2 bg-white dark:bg-gray-900 text-gray-700 dark:text-gray-300 text-xs font-semibold hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                        제어
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- DC, AI & Temperature Row (4:4:4 비율) -->
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
              
              <!-- DC Section (4 columns) -->
              <div class="premium-card overflow-hidden lg:col-span-4">
                <div class="premium-header">
                  <h2 class="gradient-title">DC Power Monitor</h2>
                </div>
                
                <div class="p-6">
                  <div class="space-y-3">
                    <!-- DC Voltage -->
                    <div class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-3">
                        <div class="w-9 h-9 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                          </svg>
                        </div>
                        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">DC Voltage</span>
                      </div>
                      <div class="text-right">
                        <span class="text-xl font-bold text-gray-900 dark:text-white tabular-nums">0.000</span>
                        <span class="text-sm font-medium text-gray-500 ml-1">V</span>
                      </div>
                    </div>

                    <!-- DC Current -->
                    <div class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-3">
                        <div class="w-9 h-9 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
                          </svg>
                        </div>
                        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">DC Current</span>
                      </div>
                      <div class="text-right">
                        <span class="text-xl font-bold text-gray-900 dark:text-white tabular-nums">0.000</span>
                        <span class="text-sm font-medium text-gray-500 ml-1">A</span>
                      </div>
                    </div>

                    <!-- Battery Current -->
                    <div class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-3">
                        <div class="w-9 h-9 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"></path>
                          </svg>
                        </div>
                        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Battery Current</span>
                      </div>
                      <div class="text-right">
                        <span class="text-xl font-bold text-gray-900 dark:text-white tabular-nums">0.000</span>
                        <span class="text-sm font-medium text-gray-500 ml-1">A</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- AI Section (4 columns) -->
              <div class="premium-card overflow-hidden lg:col-span-4">
                <div class="premium-header">
                  <h2 class="gradient-title">Analog Input (AI)</h2>
                </div>
                
                <div class="p-6">
                  <div class="grid grid-cols-2 gap-4">
                    <!-- AI 1 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Channel 1</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">8.685</div>
                    </div>

                    <!-- AI 2 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Channel 2</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">8.199</div>
                    </div>

                    <!-- AI 3 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Channel 3</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">0.000</div>
                    </div>

                    <!-- AI 4 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Channel 4</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">0.000</div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Temperature Section (4 columns) -->
              <div class="premium-card overflow-hidden lg:col-span-4">
                <div class="premium-header">
                  <h2 class="gradient-title">Temperature Monitor</h2>
                </div>
                
                <div class="p-6">
                  <div class="grid grid-cols-2 gap-4">
                    <!-- Temp 1 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Sensor 1</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">
                        28.3<span class="text-sm font-medium text-gray-500">°C</span>
                      </div>
                    </div>

                    <!-- Temp 2 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Sensor 2</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">
                        25.3<span class="text-sm font-medium text-gray-500">°C</span>
                      </div>
                    </div>

                    <!-- Temp 3 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Sensor 3</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">
                        0.0<span class="text-sm font-medium text-gray-500">°C</span>
                      </div>
                    </div>

                    <!-- Temp 4 -->
                    <div class="p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700">
                      <div class="flex items-center gap-2 mb-3">
                        <div class="w-8 h-8 rounded-lg bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                          <svg class="w-4 h-4 text-gray-600 dark:text-gray-300" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"></path>
                          </svg>
                        </div>
                        <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Sensor 4</span>
                      </div>
                      <div class="text-3xl font-bold text-gray-900 dark:text-white tabular-nums">
                        0.0<span class="text-sm font-medium text-gray-500">°C</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </main>
        <Footer />
      </div>
    </div>

    <!-- Password Modal -->
    <div v-if="showPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl p-6 max-w-xs w-full mx-4 border border-gray-200 dark:border-gray-800">
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">인증 비밀번호</h3>
        
        <!-- Password Display -->
        <div class="mb-4">
          <input 
            type="password" 
            v-model="password" 
            readonly
            class="w-full px-3 py-3 text-center text-xl tracking-widest border-2 border-gray-300 dark:border-gray-700 rounded-xl bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white font-mono focus:outline-none"
            placeholder="••••••"
          />
        </div>

        <!-- Keypad -->
        <div class="grid grid-cols-3 gap-2 mb-4">
          <button v-for="num in [1,2,3,4,5,6,7,8,9]" :key="num"
                  @click="addPasswordDigit(num)"
                  class="p-3 text-lg font-semibold rounded-xl bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 active:scale-95 transition-all">
            {{ num }}
          </button>
          <button @click="clearPassword"
                  class="p-3 text-sm font-semibold rounded-xl bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 active:scale-95 transition-all">
            AC
          </button>
          <button @click="addPasswordDigit(0)"
                  class="p-3 text-lg font-semibold rounded-xl bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 active:scale-95 transition-all">
            0
          </button>
          <button @click="deletePasswordDigit"
                  class="p-3 text-sm font-semibold rounded-xl bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 active:scale-95 transition-all">
            ←
          </button>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button @click="verifyPassword"
                  class="flex-1 px-3 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition-all text-sm">
            확인
          </button>
          <button @click="closePasswordModal"
                  class="flex-1 px-3 py-2.5 bg-gray-500 hover:bg-gray-600 text-white font-semibold rounded-xl transition-all text-sm">
            취소
          </button>
        </div>
      </div>
    </div>

    <!-- DO Control Modal -->
    <div v-if="showDOControlModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl p-8 max-w-md w-full mx-4 border-2 border-blue-500 dark:border-blue-600">
        <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Digital Output 제어</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">채널 {{ selectedDOChannel }}</p>
        
        <div class="bg-blue-50 dark:bg-blue-900/20 rounded-xl p-6 mb-6 border border-blue-200 dark:border-blue-800">
          <p class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-6 text-center">
            DO 출력 값을 변경하시겠습니까?
          </p>

          <!-- Toggle Switch -->
          <div class="flex items-center justify-center gap-6">
            <button @click="toggleDOValue" 
                    class="relative inline-flex items-center h-14 rounded-full w-28 transition-all duration-300"
                    :class="selectedDOValue === 'ON' ? 'bg-green-500' : 'bg-gray-400'">
              <span class="absolute inline-block w-10 h-10 transform rounded-full bg-white shadow-lg transition-transform duration-300"
                    :class="selectedDOValue === 'ON' ? 'translate-x-16' : 'translate-x-2'">
              </span>
            </button>
            <div class="text-center min-w-[80px]">
              <div class="text-4xl font-bold"
                   :class="selectedDOValue === 'ON' ? 'text-green-600 dark:text-green-400' : 'text-gray-500 dark:text-gray-400'">
                {{ selectedDOValue }}
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-3">
          <button @click="applyDOControl"
                  class="flex-1 px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl transition-all">
            적용
          </button>
          <button @click="closeDOControlModal"
                  class="flex-1 px-4 py-3 bg-gray-500 hover:bg-gray-600 text-white font-semibold rounded-xl transition-all">
            취소
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, computed, onUnmounted } from 'vue'
  import Sidebar from '../common/SideBar.vue'
  import Header from '../common/Header.vue'
  import Footer from "../common/Footer.vue"
  import { useAuthStore } from '@/store/auth'
  import { useRealtimeStore } from '@/store/realtime'
  import { useI18n } from 'vue-i18n'
  
  export default {
    name: 'IOModuleDashboard',
    props: ['user'],
    components: {
      Sidebar,
      Header,
      Footer,
    },
    setup(props) {
      const sidebarOpen = ref(false)
      const user = ref(props.user)
      const { t } = useI18n()
      const channel = ref('main')
      
      const authStore = useAuthStore()
      const realtimeStore = useRealtimeStore()

      // DI 상태 정의
      const diStates = ref({
        1: { type: 'DI', value: 'OFF' },
        2: { type: 'DI', value: 'ON' },
        3: { type: 'DI', value: 'ON' },
        4: { type: 'DI', value: 'OFF' },
        5: { type: 'DI', value: 'ON' },
        6: { type: 'DI', value: 'ON' },
        7: { type: 'DI', value: 'OFF' },
        8: { type: 'DI', value: 'ON' },
        9: { type: 'DI', value: 'ON' },
        10: { type: 'DI', value: 'ON' },
        11: { type: 'DI', value: 'OFF' },
        12: { type: 'DI', value: 'OFF' },
        13: { type: 'DI', value: 'OFF' },
        14: { type: 'DI', value: 'OFF' },
        15: { type: 'DI', value: 'ON' },
        16: { type: 'DI', value: 'ON' },
        17: { type: 'PI', value: 21 },
        18: { type: 'DI', value: 'OFF' },
        19: { type: 'DI', value: 'ON' },
        20: { type: 'DI', value: 'ON' },
        21: { type: 'DI', value: 'ON' },
        22: { type: 'PI', value: 12 },
        23: { type: 'PI', value: 8 },
        24: { type: 'DI', value: 'ON' }
      })

      // DO 제어 관련 상태
      const showPasswordModal = ref(false)
      const showDOControlModal = ref(false)
      const password = ref('')
      const selectedDOChannel = ref(null)
      const selectedDOValue = ref('OFF')
      const correctPassword = '1234'

      const getDICardClass = (ch) => {
        const state = diStates.value[ch]
        if (state.type === 'DI') {
          if (state.value === 'ON') {
            return 'bg-green-50 dark:bg-green-900/20 border-green-400 dark:border-green-600 text-green-700 dark:text-green-300'
          } else {
            return 'bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-700 text-gray-600 dark:text-gray-400'
          }
        } else if (state.type === 'PI') {
          return 'bg-blue-50 dark:bg-blue-900/20 border-blue-400 dark:border-blue-600 text-blue-700 dark:text-blue-300'
        }
      }

      const getDILabel = (ch) => {
        return diStates.value[ch].value
      }

      const openDOControl = (channel, currentValue) => {
        selectedDOChannel.value = channel
        selectedDOValue.value = currentValue
        showPasswordModal.value = true
        password.value = ''
      }

      const addPasswordDigit = (digit) => {
        if (password.value.length < 6) {
          password.value += digit.toString()
        }
      }

      const deletePasswordDigit = () => {
        password.value = password.value.slice(0, -1)
      }

      const clearPassword = () => {
        password.value = ''
      }

      const verifyPassword = () => {
        if (password.value === correctPassword) {
          showPasswordModal.value = false
          showDOControlModal.value = true
        } else {
          alert('비밀번호가 일치하지 않습니다')
          password.value = ''
        }
      }

      const closePasswordModal = () => {
        showPasswordModal.value = false
        password.value = ''
      }

      const toggleDOValue = () => {
        selectedDOValue.value = selectedDOValue.value === 'ON' ? 'OFF' : 'ON'
      }

      const applyDOControl = () => {
        console.log(`DO Channel ${selectedDOChannel.value} set to ${selectedDOValue.value}`)
        alert(`DO ${selectedDOChannel.value} → ${selectedDOValue.value} 변경 완료`)
        closeDOControlModal()
      }

      const closeDOControlModal = () => {
        showDOControlModal.value = false
      }
      
      const isDataReady = computed(() => {
        return realtimeStore.isDataReady()
      })
  
      onMounted(async () => {
        console.log('=== IO Module Dashboard onMounted 시작 ===')
        
        if (user.value != null) {
          authStore.setUser(user.value)
          authStore.setLogin(true)
        }
        realtimeStore.startPolling()
        
        console.log('=== IO Module Dashboard onMounted 완료 ===')
      })
  
      onUnmounted(() => {
        realtimeStore.stopPolling()
      })
  
      return {
        sidebarOpen,
        user,
        t,
        realtimeStore,
        isDataReady,
        diStates,
        getDICardClass,
        getDILabel,
        channel,
        showPasswordModal,
        showDOControlModal,
        password,
        selectedDOChannel,
        selectedDOValue,
        openDOControl,
        addPasswordDigit,
        deletePasswordDigit,
        clearPassword,
        verifyPassword,
        closePasswordModal,
        toggleDOValue,
        applyDOControl,
        closeDOControlModal,
      }
    }
  }
  </script>

  <style scoped>
  /* 카드 스타일 */
  .premium-card {
    background: linear-gradient(to bottom right, rgb(255, 255, 255), rgb(249, 250, 251));
    border-radius: 0.75rem;
    box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
    border: 1px solid rgb(229, 231, 235);
  }

  .dark .premium-card {
    background: linear-gradient(to bottom right, rgb(31, 41, 55), rgb(17, 24, 39));
    border-color: rgb(55, 65, 81);
  }

  /* 헤더 스타일 */
  .premium-header {
    padding: 0.75rem 1.5rem;
    border-bottom: 1px solid rgb(229, 231, 235);
    background: linear-gradient(to right, rgba(239, 246, 255, 0.5), rgba(250, 245, 255, 0.5));
    border-top-left-radius: 0.75rem;
    border-top-right-radius: 0.75rem;
  }

  .dark .premium-header {
    border-bottom-color: rgb(55, 65, 81);
    background: linear-gradient(to right, rgba(30, 58, 138, 0.2), rgba(88, 28, 135, 0.2));
  }

  /* 제목 그라데이션 텍스트 - 수정된 버전 */
  .gradient-title {
    font-size: 1.125rem;
    font-weight: 600;
    background: linear-gradient(90deg, rgb(37, 99, 235) 0%, rgb(147, 51, 234) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    display: inline-block;
  }

  /* 반응형 */
  @media (max-width: 640px) {
    .premium-header {
      padding: 0.5rem 1rem;
    }
  }
  </style>