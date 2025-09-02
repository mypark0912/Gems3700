<template>
  <!-- Modal backdrop -->
  <transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition ease-out duration-100"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div v-show="modalOpen" class="fixed inset-0 bg-gray-900 bg-opacity-30 z-50 transition-opacity" aria-hidden="true"></div>
  </transition>
  
  <!-- Modal dialog -->
  <transition
    enter-active-class="transition ease-in-out duration-200"
    enter-from-class="opacity-0 translate-y-4"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition ease-in-out duration-200"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-4"
  >
    <div v-show="modalOpen" class="fixed inset-0 z-50 overflow-hidden flex items-center my-4 justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
      <div ref="modalContent" class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-auto max-w-2xl w-full max-h-full">
        
        <!-- Close button -->
        <div class="absolute top-4 right-4 z-10">
          <button class="text-gray-400 dark:text-gray-500 hover:text-gray-500 dark:hover:text-gray-400" @click.stop="closeModal">
            <div class="sr-only">Close</div>
            <svg class="fill-current" width="16" height="16" viewBox="0 0 16 16">
              <path d="M7.95 6.536l4.242-4.243a1 1 0 111.415 1.414L9.364 7.95l4.243 4.242a1 1 0 11-1.415 1.415L7.95 9.364l-4.243 4.243a1 1 0 01-1.414-1.415L6.536 7.95 2.293 3.707a1 1 0 011.414-1.414L7.95 6.536z" />
            </svg>
          </button>
        </div>

        <div class="p-6">
          <!-- Progress bar (not shown for device0 mode) -->
          <div class="mb-8" v-if="OpMode !== 'device0'">
            <div class="relative">
              <div class="absolute left-0 top-1/2 -mt-px w-full h-0.5 bg-gray-200 dark:bg-gray-700/60" aria-hidden="true"></div>
              <ul class="relative flex justify-between w-full">
                <li v-for="(step, index) in availableSteps" :key="step.id">
                  <div 
                    class="flex items-center justify-center w-6 h-6 rounded-full text-xs font-semibold transition-colors"
                    :class="currentStepIndex >= index ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-900 text-gray-500 dark:text-gray-400'"
                  >
                    {{ index + 1 }}
                  </div>
                </li>
              </ul>
            </div>
          </div>

          <!-- Step 1: Settings Validation (Only for non-device0 mode) -->
          <div v-show="currentStep === 1 && OpMode !== 'device0'" class="step-content">
            <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">Settings Validation</h2>
            
            <div v-if="isLoadingSettingsValidation" class="text-center py-10">
              <div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-violet-100 dark:bg-violet-900/20 rounded-full">
                <svg class="animate-spin w-8 h-8 text-violet-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <div class="text-sm text-gray-500">Validating configuration settings...</div>
            </div>
            
            <div v-else class="space-y-4">
              <!-- Validation Summary -->
              <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                <h3 class="font-semibold mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  Validation Summary
                </h3>
                <div class="space-y-2 text-sm">
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Overall Status:</span>
                    <span :class="validationResult.isValid ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'">
                      {{ validationResult.isValid ? 'VALID' : 'INVALID' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Errors Found:</span>
                    <span :class="validationResult.errors.length > 0 ? 'text-red-600 font-semibold' : 'text-green-600'">
                      {{ validationResult.errors.length }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Warnings Found:</span>
                    <span :class="validationResult.warnings.length > 0 ? 'text-yellow-600 font-semibold' : 'text-green-600'">
                      {{ validationResult.warnings.length }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Validation Time:</span>
                    <span class="text-gray-700 dark:text-gray-300">{{ validationTimestamp }}</span>
                  </div>
                </div>
              </div>
              
              <!-- Validation Errors (if any) -->
              <div v-if="validationResult.errors.length > 0" class="bg-red-50 dark:bg-red-900/20 rounded-lg p-4">
                <h3 class="font-semibold mb-3 text-red-700 dark:text-red-400 flex items-center">
                  <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                  </svg>
                  Critical Errors (Must be fixed to proceed)
                </h3>
                <div class="max-h-40 overflow-y-auto">
                  <ul class="space-y-1 text-sm">
                    <li v-for="(error, index) in validationResult.errors" :key="index" class="flex items-start">
                      <span class="text-red-600 mr-2 mt-0.5">•</span>
                      <span class="text-red-700 dark:text-red-400">{{ error }}</span>
                    </li>
                  </ul>
                </div>
              </div>
              
              <!-- Validation Warnings (if any) -->
              <div v-if="validationResult.warnings.length > 0" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
                <h3 class="font-semibold mb-3 text-yellow-700 dark:text-yellow-400 flex items-center">
                  <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                  </svg>
                  Warnings (Can proceed with caution)
                </h3>
                <div class="max-h-40 overflow-y-auto">
                  <ul class="space-y-1 text-sm">
                    <li v-for="(warning, index) in validationResult.warnings" :key="index" class="flex items-start">
                      <span class="text-yellow-600 mr-2 mt-0.5">•</span>
                      <span class="text-yellow-700 dark:text-yellow-400">{{ warning }}</span>
                    </li>
                  </ul>
                </div>
              </div>
              
              <!-- Success Message -->
              <div v-if="validationResult.isValid && validationResult.warnings.length === 0" class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span class="text-green-700 dark:text-green-400 font-medium">All settings validation passed successfully</span>
                </div>
              </div>
              
              <!-- Success with Warnings Message -->
              <div v-else-if="validationResult.isValid && validationResult.warnings.length > 0" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-yellow-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                  </svg>
                  <span class="text-yellow-700 dark:text-yellow-400 font-medium">Settings validation completed with warnings. You can proceed but please review the warnings above.</span>
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between mt-6 space-x-2">
              <button 
                  type="button" 
                  class="btn bg-gray-500 text-white hover:bg-gray-600"
                  @click="closeModal"
              >
                  Cancel
              </button>
              <button 
                  type="button" 
                  class="btn"
                  :class="validationResult.isValid ? 
                    'bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white' :
                    'bg-gray-400 text-gray-600 cursor-not-allowed'"
                  @click="nextStep"
                  :disabled="!validationResult.isValid"
              >
                  {{ validationResult.isValid ? 'Next Step →' : 'Fix Errors First' }}
              </button>
            </div>
          </div>

          <!-- Step 2: Main Test Result (When Diagnosis_main is true) -->
          <div v-show="currentStep === 2 && diagnosis_main" class="step-content">
            <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">Main Test Result</h2>
            
            <div v-if="isLoadingMain" class="text-center py-10">
              <div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-violet-100 dark:bg-violet-900/20 rounded-full">
                <svg class="animate-spin w-8 h-8 text-violet-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <div class="text-sm text-gray-500">Loading main test result...</div>
            </div>
            
            <div v-else class="space-y-4 text-xs text-gray-800 dark:text-gray-100">
              <div class="text-sm font-semibold uppercase">Main Device Report</div>
              <div class="flex flex-wrap justify-between gap-4">
                <div>Asset Type: {{ mainTestData.AssetName || 'N/A' }}</div>
                <div>Serial Number: {{ mainTestData.SerialNumber || 'N/A' }}</div>
                <div>Channel: {{ mainTestData.Channel || 'N/A' }}</div>
              </div>
              <div class="font-semibold mb-4 mt-2">
                Result: {{ mainTestResult.err || 0 }} Errors, {{ mainTestResult.warn || 0 }} Warnings
              </div>
              <div
        class="grid grid-cols-[10%_1fr_10%] gap-2 font-semibold text-gray-500 border-b border-gray-200 pb-1"
      >
        <div class="text-left">AssemblyId</div>
        <div class="text-center">Detail</div>
        <div class="text-center">Status</div>
      </div>
      <div class="space-y-1">
        <div
          v-for="(item, index) in mainTestData.Commissions"
          :key="index"
          class="grid grid-cols-[10%_1fr_10%] gap-2 px-1 py-1 rounded hover:bg-gray-50"
        >
          <div class="text-left">{{ item.AssemblyID }}</div>
          <div class="text-left text-gray-600">
            <span v-for="(msg, i) in item.Messages" :key="i">
              • {{ msg }}
            </span>
          </div>
          <div
            class="text-center font-bold"
            :class="{
              'text-green-600': item.Status < 2,
              'text-red-600': item.Status === 3,
              'text-yellow-500': item.Status === 2,
            }"
          >
            {{ stList[item.Status] }}
          </div>
        </div>
      </div>
              
              <!-- Chart View Toggle Button -->
              <div class="mt-4">
                <button
                  @click="getWaveShow('main')"
                  class="btn bg-violet-500 text-white hover:bg-violet-600"
                >
                  <svg class="w-4 h-4 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                  View Charts
                </button>
              </div>
              
              <!-- Chart View -->
              <div v-if="showMainDiagChart" class="space-y-4">
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="option in ChartTypes"
                    :key="option"
                    @click.prevent="selectedMainChart = option; updateMainChart()"
                    :class="[
                      'btn border px-3 py-1 text-xs transition-colors duration-200 rounded-lg',
                      selectedMainChart === option
                        ? 'bg-violet-500 text-white border-violet-500'
                        : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                    ]"
                  >
                    {{ option }}
                  </button>
                </div>
                
                <div class="text-sm font-semibold text-gray-600 dark:text-gray-300 mt-4">
                  {{ selectedMainChart }} Chart
                </div>
                
                <!-- Chart Components -->
                <div class="max-h-[30vh] overflow-y-auto">
                  <LineChart
                    v-if="selectedMainChart == 'Time Domain(Voltage)'"
                    :label="mainWaveformLabelT"
                    :data="mainWaveformDataT"
                    :title="'3-Phase ' + selectedMainChart"
                    :legend="['V1', 'V2', 'V3']"
                    :mode="'3phase'"
                  />
                  <LineChart
                    v-if="selectedMainChart == 'Time Domain(Current)'"
                    :label="mainWaveformLabelT"
                    :data="mainWaveformDataT"
                    :title="'3-Phase ' + selectedMainChart"
                    :legend="['I1', 'I2', 'I3']"
                    :mode="'3phase'"
                  />
                  <LineChart
                    v-if="selectedMainChart == 'Frequency Domain(Voltage)'"
                    :label="mainWaveformLabelT"
                    :data="mainWaveformDataT"
                    :title="selectedMainChart"
                    :legend="['VFFT']"
                    :mode="'1phase'"
                  />
                  <LineChart
                    v-if="selectedMainChart == 'Frequency Domain(Current)'"
                    :label="mainWaveformLabelT"
                    :data="mainWaveformDataT"
                    :title="selectedMainChart"
                    :legend="['IFFT']"
                    :mode="'1phase'"
                  />
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between mt-6">
              <button 
                  type="button" 
                  class="btn bg-gray-500 text-white hover:bg-gray-600"
                  @click="closeModal"
              >
                  Cancel
              </button>
              <button 
                type="button" 
                class="btn"
                :class="mainTestResult.err > 0 ? 
                  'bg-gray-400 text-gray-600 cursor-not-allowed' :
                  'bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white'"
                @click="handleMainTestNext"
                :disabled="mainTestResult.err > 0"
              >
                {{ mainTestResult.err > 0 ? 'Fix Errors First' : 'Next Step →' }}
              </button>
            </div>
          </div>

          <!-- Step 3: Sub Test Result (When Diagnosis_sub is true) -->
          <div v-show="currentStep === 3 && diagnosis_sub" class="step-content">
            <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">Sub Test Result</h2>
            
            <div v-if="isLoadingSub" class="text-center py-10">
              <div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-violet-100 dark:bg-violet-900/20 rounded-full">
                <svg class="animate-spin w-8 h-8 text-violet-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <div class="text-sm text-gray-500">Loading sub test result...</div>
            </div>
            
            <div v-else class="space-y-4 text-xs text-gray-800 dark:text-gray-100">
              <div class="text-sm font-semibold uppercase">Sub Device Report</div>
              <div class="flex flex-wrap justify-between gap-4">
                <div>Asset Type: {{ subTestData.AssetName || 'N/A' }}</div>
                <div>Serial Number: {{ subTestData.SerialNumber || 'N/A' }}</div>
                <div>Channel: {{ subTestData.Channel || 'N/A' }}</div>
              </div>
              <div class="font-semibold mb-4 mt-2">
                Result: {{ subTestResult.err || 0 }} Errors, {{ subTestResult.warn || 0 }} Warnings
              </div>
              <div
        class="grid grid-cols-[10%_1fr_10%] gap-2 font-semibold text-gray-500 border-b border-gray-200 pb-1"
      >
        <div class="text-left">AssemblyId</div>
        <div class="text-center">Detail</div>
        <div class="text-center">Status</div>
      </div>
      <div class="space-y-1">
        <div
          v-for="(item, index) in subTestData.Commissions"
          :key="index"
          class="grid grid-cols-[10%_1fr_10%] gap-2 px-1 py-1 rounded hover:bg-gray-50"
        >
          <div class="text-left">{{ item.AssemblyID }}</div>
          <div class="text-left text-gray-600">
            <span v-for="(msg, i) in item.Messages" :key="i">
              • {{ msg }}
            </span>
          </div>
          <div
            class="text-center font-bold"
            :class="{
              'text-green-600': item.Status < 2,
              'text-red-600': item.Status === 3,
              'text-yellow-500': item.Status === 2,
            }"
          >
            {{ stList[item.Status] }}
          </div>
        </div>
      </div>
              
              <!-- Chart View Toggle Button -->
              <div class="mt-4">
                <button
                  @click="getWaveShow('sub')"
                  class="btn bg-violet-500 text-white hover:bg-violet-600"
                >
                  <svg class="w-4 h-4 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                  </svg>
                  View Charts
                </button>
              </div>
              
              <!-- Sub Chart View -->
              <div v-if="showSubDiagChart" class="space-y-4">
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="option in ChartTypes"
                    :key="option"
                    @click.prevent="selectedSubChart = option; updateSubChart()"
                    :class="[
                      'btn border px-3 py-1 text-xs transition-colors duration-200 rounded-lg',
                      selectedSubChart === option
                        ? 'bg-violet-500 text-white border-violet-500'
                        : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                    ]"
                  >
                    {{ option }}
                  </button>
                </div>
                
                <div class="text-sm font-semibold text-gray-600 dark:text-gray-300 mt-4">
                  {{ selectedSubChart }} Chart
                </div>
                
                <!-- Sub Chart Components -->
                <div class="max-h-[30vh] overflow-y-auto">
                  <LineChart
                    v-if="selectedSubChart == 'Time Domain(Voltage)'"
                    :label="subWaveformLabelT"
                    :data="subWaveformDataT"
                    :title="'3-Phase ' + selectedSubChart"
                    :legend="['V1', 'V2', 'V3']"
                    :mode="'3phase'"
                  />
                  <LineChart
                    v-if="selectedSubChart == 'Time Domain(Current)'"
                    :label="subWaveformLabelT"
                    :data="subWaveformDataT"
                    :title="'3-Phase ' + selectedSubChart"
                    :legend="['I1', 'I2', 'I3']"
                    :mode="'3phase'"
                  />
                  <LineChart
                    v-if="selectedSubChart == 'Frequency Domain(Voltage)'"
                    :label="subWaveformLabelT"
                    :data="subWaveformDataT"
                    :title="selectedSubChart"
                    :legend="['VFFT']"
                    :mode="'1phase'"
                  />
                  <LineChart
                    v-if="selectedSubChart == 'Frequency Domain(Current)'"
                    :label="subWaveformLabelT"
                    :data="subWaveformDataT"
                    :title="selectedSubChart"
                    :legend="['IFFT']"
                    :mode="'1phase'"
                  />
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between mt-6">
              <button 
                  type="button" 
                  class="btn bg-gray-500 text-white hover:bg-gray-600"
                  @click="closeModal"
              >
                  Cancel
              </button>
              <button 
                type="button" 
                class="btn"
                :class="subTestResult.err > 0 ? 
                  'bg-gray-400 text-gray-600 cursor-not-allowed' :
                  'bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white'"
                @click="nextStep"
                :disabled="subTestResult.err > 0"
              >
                {{ subTestResult.err > 0 ? 'Fix Errors First' : 'Next Step →' }}
              </button>
            </div>
          </div>

          <!-- Step 4: Validation Result (Final for all modes) -->
          <div v-show="currentStep === 4" class="step-content">
            <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">
              {{ OpMode === 'device0' ? 'Validation Result' : 'Complete Setup' }}
            </h2>
            
            <!-- Validation Result for device0 mode -->
            <div v-if="OpMode === 'device0'" class="space-y-4">
              <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                <h3 class="font-semibold mb-3">Validation Summary</h3>
                <div class="space-y-2 text-sm">
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Device Status:</span>
                    <span :class="validationResult.isValid ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'">
                      {{ validationResult.isValid ? 'PASSED' : 'FAILED' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Validation Time:</span>
                    <span>{{ validationTimestamp || new Date().toLocaleString() }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Error Count:</span>
                    <span :class="validationResult.errors.length > 0 ? 'text-red-600' : 'text-green-600'">
                      {{ validationResult.errors.length || 0 }}
                    </span>
                  </div>
                </div>
              </div>
              
              <!-- Validation Details -->
              <div v-if="validationResult.errors.length > 0 || validationResult.warnings.length > 0" class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                <h4 class="font-semibold mb-2 text-sm">Validation Details</h4>
                <ul class="space-y-1 text-sm">
                  <li v-for="(error, index) in validationResult.errors" :key="'error-' + index" class="flex items-start">
                    <svg class="w-4 h-4 shrink-0 fill-current text-red-500 mr-2 mt-0.5" viewBox="0 0 12 12">
                      <path d="M11.707 1.293a1 1 0 0 0-1.414 0L6 5.586 1.707 1.293a1 1 0 0 0-1.414 1.414L4.586 7 .293 11.293a1 1 0 0 0 1.414 1.414L6 8.414l4.293 4.293a1 1 0 0 0 1.414-1.414L7.414 7l4.293-4.293a1 1 0 0 0 0-1.414z"/>
                    </svg>
                    <span>{{ error }}</span>
                  </li>
                  <li v-for="(warning, index) in validationResult.warnings" :key="'warning-' + index" class="flex items-start">
                    <svg class="w-4 h-4 shrink-0 fill-current text-yellow-500 mr-2 mt-0.5" viewBox="0 0 12 12">
                      <path d="M6 0C2.686 0 0 2.686 0 6s2.686 6 6 6 6-2.686 6-6-2.686-6-6-6zm0 10c-.55 0-1-.45-1-1s.45-1 1-1 1 .45 1 1-.45 1-1 1zm1-3H5V3h2v4z"/>
                    </svg>
                    <span>{{ warning }}</span>
                  </li>
                </ul>
              </div>
            </div>
            
            <!-- Setup Complete for other modes -->
            <div v-else class="space-y-4">
              <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                <h3 class="font-semibold mb-3">Configuration Summary</h3>
                <div class="space-y-2 text-sm">
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Settings Validation:</span>
                    <span class="text-green-600 font-semibold">PASSED</span>
                  </div>
                  <div v-if="Diagnosis_main" class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Main Test:</span>
                    <span :class="mainTestResult.err > 0 ? 'text-red-600' : 'text-green-600'" class="font-semibold">
                      {{ mainTestResult.err > 0 ? `${mainTestResult.err} Errors` : 'PASSED' }}
                    </span>
                  </div>
                  <div v-if="Diagnosis_sub" class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Sub Test:</span>
                    <span :class="subTestResult.err > 0 ? 'text-red-600' : 'text-green-600'" class="font-semibold">
                      {{ subTestResult.err > 0 ? `${subTestResult.err} Errors` : 'PASSED' }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">Overall Status:</span>
                    <span class="text-green-600 font-semibold">READY</span>
                  </div>
                </div>
              </div>
              
              <div class="text-xs text-gray-500">
                By completing this setup, you confirm that all diagnosis tests have been reviewed and the system is ready for operation.
              </div>
            </div>
            
            <!-- Modal footer -->
            <div class="flex flex-wrap justify-between items-center mt-6">
              <button 
                  type="button" 
                  class="btn bg-gray-500 text-white hover:bg-gray-600"
                  @click="closeModal"
              >
                  Cancel
              </button>
              <div :class="OpMode === 'device0' ? 'ml-auto flex space-x-3' : 'ml-auto space-x-2'" class="flex">
                <button                  
                  class="btn bg-violet-500 text-white hover:bg-violet-600"
                  @click="restartValidation"
                >
                  Restart
                </button>
                <!--button
                  class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
                  @click="closeModal"
                >
                  {{ OpMode === 'device0' ? 'Close' : 'Complete' }}
                </button-->
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import LineChart from '../charts/connect/LineChart01_Echart2.vue'
import { settingValidator } from "@/utils/validation.js";
import axios from "axios";

export default {
  name: 'OnboardingModal',
  components: {
    LineChart
  },
  props: {
    modalOpen: {
      type: Boolean,
      default: false
    },
    OpMode: {
      type: String,
      default: 'device0' // 'device0' or other modes
    }
  },
  emits: ['close-modal', 'complete-setup', 'restart-validation'],
  setup(props, { emit }) {
    const modalContent = ref(null)
    const currentStep = ref(1)
    const setupDict = ref({})
    
    // Loading states
    const isLoadingSettingsValidation = ref(false)
    const isLoadingMain = ref(false)
    const isLoadingSub = ref(false)
    const stList = ref(["Info", "Pass", "Warning", "Error"]);
    

    
    // Validation data
    const validationTimestamp = ref(new Date().toLocaleString())
    const validationResult = ref({
      isValid: true,
      errors: [],
      warnings: [],
      hasWarnings: false,
      hasErrors: false
    })
    
    // Main test data
    const mainTestData = ref({
      AssetName: 'N/A',
      SerialNumber: 'N/A',
      Channel: 'N/A',
      Commissions: []
    })

    const diagnosis_main = computed(() => {
    // props가 명시적으로 전달되지 않으면 undefined일 수 있음

    if(setupDict.value.main.Enable){
      if (setupDict.value.main.assetInfo.name != "" && setupDict.value.General.useFuction.diagnosis_main)
        return true;
      else
        return false;
    }else{
      return false;
    }
    // return setupDict.value.main.Enable 
  })
  
  const diagnosis_sub = computed(() => {
    if(setupDict.value.sub.Enable){
      if (setupDict.value.sub.assetInfo.name != "" && setupDict.value.General.useFuction.diagnosis_sub)
        return true;
      else
        return false;
    }else{
      return false;
    }
    // return setupDict.value.sub.Enable
  })

    const mainTestResult = ref({ err: 0, warn: 0 })
    const showMainDiagChart = ref(false)
    const mainWaveformData = ref({})
    const mainWaveformDataT = ref([])
    const mainWaveformLabelT = ref([])
    const selectedMainChart = ref('Time Domain(Voltage)')
    
    // Sub test data
    const subTestData = ref({
      AssetName: 'N/A',
      SerialNumber: 'N/A',
      Channel: 'N/A',
      Commissions: []
    })
    const subTestResult = ref({ err: 0, warn: 0 })
    const showSubDiagChart = ref(false)
    const subWaveformData = ref({})
    const subWaveformDataT = ref([])
    const subWaveformLabelT = ref([])
    const selectedSubChart = ref('Time Domain(Voltage)')
    
    // Chart types
    const ChartTypes = ref([
      "Time Domain(Voltage)",
      "Time Domain(Current)",
      "Frequency Domain(Voltage)",
      "Frequency Domain(Current)",
    ])
    
    // Settings validation function
    const validateSettings = () => {
      if (!setupDict.value || Object.keys(setupDict.value).length === 0) {
        validationResult.value = {
          isValid: false,
          errors: ['No settings data loaded'],
          warnings: [],
          hasWarnings: false,
          hasErrors: true
        }
        return validationResult.value
      }
      
      try {
        const result = settingValidator.validateAllSettings(
          setupDict.value["General"],
          setupDict.value["main"],
          setupDict.value["sub"]
        )
        validationResult.value = result
        validationTimestamp.value = new Date().toLocaleString()
        return result
      } catch (error) {
        console.error('Validation error:', error)
        validationResult.value = {
          isValid: false,
          errors: ['Validation process failed: ' + error.message],
          warnings: [],
          hasWarnings: false,
          hasErrors: true
        }
        return validationResult.value
      }
    }
    
    // Load settings data
    const GetSettingData = async () => {
      try {
        console.log('Loading settings data...')
        isLoadingSettingsValidation.value = true
        
        const response = await axios.get(`/setting/getSetting`);

        if (response.data.passOK == 1) {
          setupDict.value = response.data.data;
          console.log('Settings loaded:', setupDict.value)
          
          // Validate settings after loading
          validateSettings()
        }
      } catch (error) {
        console.error("Settings data loading failed:", error);
        validationResult.value = {
          isValid: false,
          errors: ['Failed to load settings data'],
          warnings: [],
          hasWarnings: false,
          hasErrors: true
        }
      } finally {
        isLoadingSettingsValidation.value = false
      }
    }
    
    // Get commission test data

    const getWaveShow = async(channel)=>{
      if (channel == 'main')
        showMainDiagChart.value = true;
      else
        showSubDiagChart.value = true;
      await getWaveform(channel)
    }

    const getWaveform = async (channel) =>{
      const isMain = channel === 'main'?true:false;
      const waveformDataRef = isMain ? mainWaveformData : subWaveformData
      const assetName = isMain ? setupDict.value["main"]["assetInfo"]["name"] : setupDict.value["sub"]["assetInfo"]["name"];
      const showChartRef = isMain ? showMainDiagChart : showSubDiagChart
      const response = await axios.get(`/setting/testwave/${assetName}`);

      if (response.data.success) {
        console.log('getWave',response.data);
        if (response.data.waveT && Object.keys(response.data.waveT).length > 0) {
            showChartRef.value = true
            waveformDataRef.value = response.data.waveT
            
            // Update chart display
            if (isMain) {
              updateMainChart()
            } else {
              updateSubChart()
            }
          }
      }
    }

    const getCommision = async (channel) => {
      const isMain = channel === 'main'?true:false;
      const loadingRef = isMain ? isLoadingMain : isLoadingSub
      const testDataRef = isMain ? mainTestData : subTestData
      const testResultRef = isMain ? mainTestResult : subTestResult
      
      try {
        loadingRef.value = true
        
        // Get asset name from setupStore
        const assetName = isMain ? 
          setupDict.value["main"]["assetInfo"]["name"] : setupDict.value["sub"]["assetInfo"]["name"]
        console.log('Setup:',isMain?setupDict.value.main:setupDict.value.sub);
        console.log(`Loading ${channel} test data for asset: ${assetName}`)
        
        const response = await axios.get(`/setting/test/${assetName}`)

        if (response.data.success) {
          testDataRef.value = response.data.data
          //console.log(subTestData.value);
          // Count errors and warnings
          let errCount = 0, warnCount = 0
          for (let i = 0; i < testDataRef.value["Commissions"].length; i++) {
            if (testDataRef.value["Commissions"][i]["Status"] == 3) errCount += 1
            else if (testDataRef.value["Commissions"][i]["Status"] == 2) warnCount += 1
          }
          testResultRef.value = { err: errCount, warn: warnCount }
          
          // Check if waveform data exists
          
        }
      } catch (e) {
        console.error(`Error loading ${channel} test:`, e)
        testResultRef.value = { err: 1, warn: 0 }
      } finally {
        loadingRef.value = false
      }
    }
    
    // Update main chart data
    const updateMainChart = () => {
      if (!mainWaveformData.value || Object.keys(mainWaveformData.value).length === 0) return
      
      const SR = mainWaveformData.value["SR"]
      const Duration = mainWaveformData.value["Duration"]
      const deltaT = 1 / SR
      const deltaF = 1 / Duration
      
      const chartResult = {
        "Time Domain(Voltage)": mainWaveformData.value["Vwave"],
        "Time Domain(Current)": mainWaveformData.value["Iwave"],
        "Frequency Domain(Voltage)": mainWaveformData.value["Vfft"],
        "Frequency Domain(Current)": mainWaveformData.value["Ifft"],
      }
      
      if (selectedMainChart.value.includes("Time")) {
        mainWaveformLabelT.value = Array.from(
          { length: chartResult[selectedMainChart.value][0].length },
          (_, i) => i * deltaT
        )
        mainWaveformDataT.value = chartResult[selectedMainChart.value]
      } else {
        mainWaveformLabelT.value = Array.from(
          { length: chartResult[selectedMainChart.value].length },
          (_, i) => i * deltaF
        )
        mainWaveformDataT.value = [chartResult[selectedMainChart.value]]
      }
    }
    
    // Update sub chart data
    const updateSubChart = () => {
      if (!subWaveformData.value || Object.keys(subWaveformData.value).length === 0) return
      
      const SR = subWaveformData.value["SR"]
      const Duration = subWaveformData.value["Duration"]
      const deltaT = 1 / SR
      const deltaF = 1 / Duration
      
      const chartResult = {
        "Time Domain(Voltage)": subWaveformData.value["Vwave"],
        "Time Domain(Current)": subWaveformData.value["Iwave"],
        "Frequency Domain(Voltage)": subWaveformData.value["Vfft"],
        "Frequency Domain(Current)": subWaveformData.value["Ifft"],
      }
      
      if (selectedSubChart.value.includes("Time")) {
        subWaveformLabelT.value = Array.from(
          { length: chartResult[selectedSubChart.value][0].length },
          (_, i) => i * deltaT
        )
        subWaveformDataT.value = chartResult[selectedSubChart.value]
      } else {
        subWaveformLabelT.value = Array.from(
          { length: chartResult[selectedSubChart.value].length },
          (_, i) => i * deltaF
        )
        subWaveformDataT.value = [chartResult[selectedSubChart.value]]
      }
    }
    
    // Compute available steps based on OpMode and diagnosis settings  
    const availableSteps = computed(() => {
      if (props.OpMode === 'device0') {
        return [{ id: 4, name: 'Validation' }]
      }
      
      const steps = [
        { id: 1, name: 'Settings Check' }
      ]
      
      if (diagnosis_main.value) {
        steps.push({ id: 2, name: 'Main Test' })
      }
      
      if (diagnosis_sub.value) {
        steps.push({ id: 3, name: 'Sub Test' })
      }
      
      steps.push({ id: 4, name: 'Complete' })
      
      return steps
    })
    
    // Get current step index in available steps
    const currentStepIndex = computed(() => {
      return availableSteps.value.findIndex(step => step.id === currentStep.value)
    })

    // Navigation functions
    const nextStep = async () => {
      // 현재 단계가 Main Test이고 에러가 있으면 진행 불가
      if (currentStep.value === 2 && mainTestResult.value.err > 0) {
        return
      }
      
      // 현재 단계가 Sub Test이고 에러가 있으면 진행 불가
      if (currentStep.value === 3 && subTestResult.value.err > 0) {
        return
      }
      
      const currentIndex = currentStepIndex.value
      if (currentIndex < availableSteps.value.length - 1) {
        const nextStepId = availableSteps.value[currentIndex + 1].id
        currentStep.value = nextStepId
        
        // Load data when entering main or sub test steps
        if (nextStepId === 2 && diagnosis_main) {
          await getCommision('main')
        } else if (nextStepId === 3 && diagnosis_sub) {
          await getCommision('sub')
        }
      }
    }
    
    // Handle main test next button - skip sub if not needed
    const handleMainTestNext = async () => {
      // Main test에 에러가 있으면 진행 불가
      if (mainTestResult.value.err > 0) {
        return
      }
      
      if (diagnosis_sub.value) {
        currentStep.value = 3
        await getCommision('sub')
      } else {
        currentStep.value = 4
      }
    }

    const prevStep = () => {
      const currentIndex = currentStepIndex.value
      if (currentIndex > 0) {
        currentStep.value = availableSteps.value[currentIndex - 1].id
      }
    }

    const closeModal = () => {
      emit('close-modal')
    }

    const completeSetup = () => {
      emit('complete-setup')
      closeModal()
    }
    
    const restartValidation = async () => {
      // Reset all data
      validationResult.value = {
        isValid: true,
        errors: [],
        warnings: [],
        hasWarnings: false,
        hasErrors: false
      }
      
      // Reload settings and validate
      //await GetSettingData()
      
      emit('restart-validation')
    }

    // Initialize step based on OpMode
    watch(() => props.modalOpen, async (newValue) => {
      if (newValue) {
        // Load settings data
        await GetSettingData()
        
        if (props.OpMode === 'device0') {
          currentStep.value = 4
        } else {
          currentStep.value = 1
          
          // If validation passes and only has main test, load it immediately
          if (validationResult.value.isValid && 
              diagnosis_main.value && 
              !diagnosis_sub.value &&
              availableSteps.value.length === 3) {
            // Only settings, main, and complete steps
            setTimeout(() => {
              getCommision('main')
            }, 100)
          }
        }
      }
    })
    
    // Update current step if diagnosis settings change
    watch([() => diagnosis_main, () => diagnosis_sub], () => {
      const validStepIds = availableSteps.value.map(s => s.id)
      if (!validStepIds.includes(currentStep.value)) {
        currentStep.value = validStepIds[0]
      }
    })

    // Close on click outside
    const clickHandler = ({ target }) => {
      if (!props.modalOpen || modalContent.value.contains(target)) return
      closeModal()
    }

    // Close if the esc key is pressed
    const keyHandler = ({ keyCode }) => {
      if (!props.modalOpen || keyCode !== 27) return
      closeModal()
    }

    onMounted(() => {
      document.addEventListener('click', clickHandler)
      document.addEventListener('keydown', keyHandler)
      
      // Load data if modal is already open
      if (props.modalOpen) {
        GetSettingData()
      }
    })

    onUnmounted(() => {
      document.removeEventListener('click', clickHandler)
      document.removeEventListener('keydown', keyHandler)
    })

    return {
      modalContent,
      currentStep,
      currentStepIndex,
      availableSteps,
      validationResult,
      isLoadingSettingsValidation,
      isLoadingMain,
      isLoadingSub,
      validationTimestamp,
      mainTestData,
      mainTestResult,
      showMainDiagChart,
      mainWaveformDataT,
      mainWaveformLabelT,
      selectedMainChart,
      subTestData,
      subTestResult,
      showSubDiagChart,
      subWaveformDataT,
      subWaveformLabelT,
      selectedSubChart,
      ChartTypes,
      nextStep,
      prevStep,
      handleMainTestNext,
      closeModal,
      completeSetup,
      restartValidation,
      updateMainChart,
      updateSubChart,
      stList,
      diagnosis_main,
      diagnosis_sub,
      getWaveShow,
    }
  }
}
</script>

<style scoped>
/* Animation for step transitions */
.step-content {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Button styles */
.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors duration-200;
}

.btn-sm {
  @apply px-3 py-1.5 rounded-lg font-medium transition-colors duration-200 text-sm;
}
</style>