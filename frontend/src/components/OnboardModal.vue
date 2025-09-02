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
            <div v-show="currentStep === 2 && Diagnosis_main" class="step-content">
              <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">Main Test Result</h2>
              
              <div v-if="isLoading" class="text-center py-10 text-sm text-gray-500">
                Loading test result...
              </div>
              
              <div v-else class="space-y-4 text-xs text-gray-800 dark:text-gray-100">
                <div class="text-sm font-semibold uppercase">Main Device Report</div>
                <div class="flex flex-wrap justify-between gap-4">
                  <div>Asset Type: {{ testData.AssetName?.AssemblyType || 'N/A' }}</div>
                  <div>Serial Number: {{ testData.SerialNumber || 'N/A' }}</div>
                  <div>Channel: {{ testData.Channel || 'N/A' }}</div>
                </div>
                <div class="font-semibold mb-4 mt-2">
                  Result: {{ testResult.err || 0 }} Errors, {{ testResult.warn || 0 }} Warnings
                </div>
                
                <!-- Chart View -->
                <div v-if="showDiagChart" class="space-y-4">
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="option in ChartTypes"
                      :key="option"
                      @click.prevent="selectedChart = option"
                      :class="[
                        'btn border px-3 py-1 text-xs transition-colors duration-200 rounded-lg',
                        selectedChart === option
                          ? 'bg-violet-500 text-white border-violet-500'
                          : 'bg-white text-violet-500 border-gray-200 hover:bg-gray-100 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-900',
                      ]"
                    >
                      {{ option }}
                    </button>
                  </div>
                  
                  <div class="text-sm font-semibold text-gray-600 dark:text-gray-300 mt-4">
                    {{ selectedChart }} Chart
                  </div>
                  
                  <!-- Chart Components -->
                  <div class="max-h-[30vh] overflow-y-auto">
                    <LineChart
                      v-if="selectedChart == 'Time Domain(Voltage)'"
                      :label="waveformLabelT"
                      :data="waveformDataT"
                      :title="'3-Phase ' + selectedChart"
                      :legend="['V1', 'V2', 'V3']"
                      :mode="'3phase'"
                    />
                    <LineChart
                      v-if="selectedChart == 'Time Domain(Current)'"
                      :label="waveformLabelT"
                      :data="waveformDataT"
                      :title="'3-Phase ' + selectedChart"
                      :legend="['I1', 'I2', 'I3']"
                      :mode="'3phase'"
                    />
                    <LineChart
                      v-if="selectedChart == 'Frequency Domain(Voltage)'"
                      :label="waveformLabelT"
                      :data="waveformDataT"
                      :title="selectedChart"
                      :legend="['VFFT']"
                      :mode="'1phase'"
                    />
                    <LineChart
                      v-if="selectedChart == 'Frequency Domain(Current)'"
                      :label="waveformLabelT"
                      :data="waveformDataT"
                      :title="selectedChart"
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
                  class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white ml-auto"
                  @click="nextStep"
                >
                  Next Step →
                </button>
              </div>
            </div>
  
            <!-- Step 3: Sub Test Result (When Diagnosis_sub is true) -->
            <div v-show="currentStep === 3 && Diagnosis_sub" class="step-content">
              <h2 class="text-2xl text-gray-800 dark:text-gray-100 font-bold mb-6">Sub Test Result</h2>
              
              <div v-if="isLoadingSub" class="text-center py-10 text-sm text-gray-500">
                Loading sub test result...
              </div>
              
              <div v-else class="space-y-4 text-xs text-gray-800 dark:text-gray-100">
                <div class="text-sm font-semibold uppercase">Sub Device Report</div>
                <div class="flex flex-wrap justify-between gap-4">
                  <div>Asset Type: {{ subTestData.AssetName?.AssemblyType || 'N/A' }}</div>
                  <div>Serial Number: {{ subTestData.SerialNumber || 'N/A' }}</div>
                  <div>Channel: {{ subTestData.Channel || 'N/A' }}</div>
                </div>
                <div class="font-semibold mb-4 mt-2">
                  Result: {{ subTestResult.err || 0 }} Errors, {{ subTestResult.warn || 0 }} Warnings
                </div>
                
                <!-- Sub Chart View -->
                <div v-if="showSubDiagChart" class="space-y-4">
                  <div class="flex flex-wrap gap-2">
                    <button
                      v-for="option in ChartTypes"
                      :key="option"
                      @click.prevent="selectedSubChart = option"
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
                  class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white ml-auto"
                  @click="nextStep"
                >
                  Next Step →
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
                      <span :class="validationResult.status === 'passed' ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'">
                        {{ validationResult.status?.toUpperCase() || 'PENDING' }}
                      </span>
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-gray-600 dark:text-gray-400">Validation Time:</span>
                      <span>{{ validationResult.timestamp || new Date().toLocaleString() }}</span>
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-gray-600 dark:text-gray-400">Error Count:</span>
                      <span :class="validationResult.errors > 0 ? 'text-red-600' : 'text-green-600'">
                        {{ validationResult.errors || 0 }}
                      </span>
                    </div>
                  </div>
                </div>
                
                <!-- Validation Details -->
                <div v-if="validationResult.details" class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                  <h4 class="font-semibold mb-2 text-sm">Validation Details</h4>
                  <ul class="space-y-1 text-sm">
                    <li v-for="(detail, index) in validationResult.details" :key="index" class="flex items-start">
                      <svg v-if="detail.passed" class="w-4 h-4 shrink-0 fill-current text-green-500 mr-2 mt-0.5" viewBox="0 0 12 12">
                        <path d="M10.28 1.28L3.989 7.575 1.695 5.28A1 1 0 00.28 6.695l3 3a1 1 0 001.414 0l7-7A1 1 0 0010.28 1.28z" />
                      </svg>
                      <svg v-else class="w-4 h-4 shrink-0 fill-current text-red-500 mr-2 mt-0.5" viewBox="0 0 12 12">
                        <path d="M11.707 1.293a1 1 0 0 0-1.414 0L6 5.586 1.707 1.293a1 1 0 0 0-1.414 1.414L4.586 7 .293 11.293a1 1 0 0 0 1.414 1.414L6 8.414l4.293 4.293a1 1 0 0 0 1.414-1.414L7.414 7l4.293-4.293a1 1 0 0 0 0-1.414z"/>
                      </svg>
                      <span>{{ detail.message }}</span>
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
                      <span class="text-gray-600 dark:text-gray-400">Check validation</span>
                      <span class="text-green-600 font-semibold"> OK </span>
                    </div>
                    <div class="flex items-center justify-between">
                      <span class="text-gray-600 dark:text-gray-400">Tests Completed:</span>
                      <span class="font-medium">
                        {{ (Diagnosis_main ? 1 : 0) + (Diagnosis_sub ? 1 : 0) }} of 2
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
                <div :class="OpMode === 'device0' ? 'w-full flex justify-center space-x-3' : 'ml-auto space-x-2'" class="flex">
                  <button
                    v-if="OpMode === 'device0'"
                    class="btn bg-violet-500 text-white hover:bg-violet-600"
                    @click="restartValidation"
                  >
                    Restart
                  </button>
                  <button
                    class="btn bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
                    @click="closeModal"
                  >
                    {{ OpMode === 'device0' ? 'Close' : 'Complete' }}
                  </button>
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
  // LineChart 컴포넌트가 없으면 주석 처리
  // import LineChart from '../charts/connect/LineChart01_Echart2.vue'
  import { settingValidator } from "@/utils/validation.js";
  import axios from "axios";
  
  export default {
    name: 'OnboardModal',
    components: {
      // LineChart  // 컴포넌트가 없으면 주석 처리
    },
    props: {
      modalOpen: {
        type: Boolean,
        default: false
      },
      OpMode: {
        type: String,
        default: 'standard' // 'device0' or other modes
      },
      Diagnosis_main: {
        type: Boolean,
        default: true
      },
      Diagnosis_sub: {
        type: Boolean,
        default: true
      },
      // Main test data
      testData: {
        type: Object,
        default: () => ({
          AssetName: { AssemblyType: 'N/A' },
          SerialNumber: 'N/A',
          Channel: 'N/A',
          Commissions: []
        })
      },
      testResult: {
        type: Object,
        default: () => ({ err: 0, warn: 0 })
      },
      isLoading: {
        type: Boolean,
        default: false
      },
      showDiagChart: {
        type: Boolean,
        default: true
      },
      // Sub test data
      subTestData: {
        type: Object,
        default: () => ({
          AssetName: { AssemblyType: 'N/A' },
          SerialNumber: 'N/A',
          Channel: 'N/A',
          Commissions: []
        })
      },
      subTestResult: {
        type: Object,
        default: () => ({ err: 0, warn: 0 })
      },
      isLoadingSub: {
        type: Boolean,
        default: false
      },
      showSubDiagChart: {
        type: Boolean,
        default: true
      },
      // Chart data
      ChartTypes: {
        type: Array,
        default: () => [
          'Time Domain(Voltage)',
          'Time Domain(Current)',
          'Frequency Domain(Voltage)',
          'Frequency Domain(Current)'
        ]
      },
      waveformLabelT: {
        type: Array,
        default: () => []
      },
      waveformDataT: {
        type: Array,
        default: () => []
      },
      subWaveformLabelT: {
        type: Array,
        default: () => []
      },
      subWaveformDataT: {
        type: Array,
        default: () => []
      },
      // Validation result for device0 mode
      validationResult: {
        type: Object,
        default: () => ({
          status: 'passed',
          timestamp: new Date().toLocaleString(),
          errors: 0,
          details: [
            { passed: true, message: 'Device connectivity check passed' },
            { passed: true, message: 'Configuration validation successful' },
            { passed: false, message: 'Firmware version needs update' }
          ]
        })
      }
    },
    emits: ['close-modal', 'complete-setup', 'restart-validation', 'retry-settings-validation'],
    setup(props, { emit }) {
      const modalContent = ref(null)
      const currentStep = ref(1) // Settings validation부터 시작 (ID를 1로 변경)
      const setupDict = ref({})
      const isLoadingSettingsValidation = ref(false)
      const validationTimestamp = ref(new Date().toLocaleString())
      
      // 초기 검증 결과
      const internalValidationResult = ref({
        isValid: true,
        errors: [],
        warnings: [],
        hasWarnings: false,
        hasErrors: false
      })
      
      // Chart selections
      const selectedChart = ref('Time Domain(Voltage)')
      const selectedSubChart = ref('Time Domain(Voltage)')
      const stList = ['OK', 'Warning', 'Alert', 'Error']
      
      // Settings validation function
      const validateSettings = () => {
        if (!setupDict.value || Object.keys(setupDict.value).length === 0) {
          internalValidationResult.value = {
            isValid: false,
            errors: ['No settings data loaded'],
            warnings: [],
            hasWarnings: false,
            hasErrors: true
          }
          return internalValidationResult.value
        }
        
        try {
          const result = settingValidator.validateAllSettings(
            setupDict.value["General"], // General 설정
            setupDict.value["main"], // Main Channel 설정
            setupDict.value["sub"] // Sub Channel 설정
          )
          internalValidationResult.value = result
          validationTimestamp.value = new Date().toLocaleString()
          return result
        } catch (error) {
          console.error('Validation error:', error)
          internalValidationResult.value = {
            isValid: false,
            errors: ['Validation process failed: ' + error.message],
            warnings: [],
            hasWarnings: false,
            hasErrors: true
          }
          return internalValidationResult.value
        }
      }
      
      // Compute available steps based on OpMode and diagnosis settings  
      const availableSteps = computed(() => {
        console.log('Computing available steps, OpMode:', props.OpMode)
        
        if (props.OpMode === 'device0') {
          return [{ id: 4, name: 'Validation' }]
        }
        
        const steps = [
          { id: 1, name: 'Settings Check' }  // 설정값 검증 스텝
        ]
        
        if (props.Diagnosis_main) {
          steps.push({ id: 2, name: 'Main Test' })
        }
        
        if (props.Diagnosis_sub) {
          steps.push({ id: 3, name: 'Sub Test' })
        }
        
        steps.push({ id: 4, name: 'Complete' })
        
        console.log('Available steps computed:', steps)
        return steps
      })
      
      // Get current step index in available steps
      const currentStepIndex = computed(() => {
        return availableSteps.value.findIndex(step => step.id === currentStep.value)
      })

      // Navigation functions
      const nextStep = () => {
        const currentIndex = currentStepIndex.value
        if (currentIndex < availableSteps.value.length - 1) {
          currentStep.value = availableSteps.value[currentIndex + 1].id
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
      
      const restartValidation = () => {
        emit('restart-validation')
      }
      
      const retrySettingsValidation = async () => {
        isLoadingSettingsValidation.value = true
        emit('retry-settings-validation')
        
        try {
          // 설정 데이터 다시 로드
          await GetSettingData()
          // 검증 실행
          validateSettings()
        } catch (error) {
          console.error('Settings validation retry failed:', error)
        } finally {
          isLoadingSettingsValidation.value = false
        }
      }

      // Initialize step based on OpMode
      watch(() => props.modalOpen, (newValue) => {
        console.log('Modal open changed:', newValue)
        if (newValue) {
          // 설정 데이터 로드
          GetSettingData();
          
          if (props.OpMode === 'device0') {
            currentStep.value = 4
          } else {
            currentStep.value = 1  // 설정값 검증부터 시작 (ID 1로 변경)
          }
          console.log('Modal opened, current step:', currentStep.value)
          console.log('Available steps:', availableSteps.value)
        }
      })
      
      // Update current step if diagnosis settings change
      watch([() => props.Diagnosis_main, () => props.Diagnosis_sub], () => {
        // Ensure current step is still valid
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
        
        // 모달이 이미 열려있다면 데이터 로드
        if (props.modalOpen) {
          GetSettingData();
        }
      })

      onUnmounted(() => {
        document.removeEventListener('click', clickHandler)
        document.removeEventListener('keydown', keyHandler)
      })
      
      const GetSettingData = async () => {
        try {
          console.log('Loading settings data...')
          isLoadingSettingsValidation.value = true
          
          const response = await axios.get(`/setting/getSetting`);

          if (response.data.passOK == 1) {
            setupDict.value = response.data.data;
            console.log('Settings loaded:', setupDict.value)
            
            // 설정 로드 후 검증 실행
            validateSettings()
            
            // 콘솔에 검증 결과 출력
            if (internalValidationResult.value.hasErrors) {
              console.log('Validation errors:', internalValidationResult.value.errors);
            }
            if (internalValidationResult.value.hasWarnings) {
              console.log('Validation warnings:', internalValidationResult.value.warnings);
            }
          }
        } catch (error) {
          console.error("Settings data loading failed:", error);
          internalValidationResult.value = {
            isValid: false,
            errors: ['Failed to load settings data'],
            warnings: [],
            hasWarnings: false,
            hasErrors: true
          }
        } finally {
          isLoadingSettingsValidation.value = false
        }
      };

      return {
        modalContent,
        currentStep,
        currentStepIndex,
        availableSteps,
        selectedChart,
        selectedSubChart,
        stList,
        validationResult: internalValidationResult,
        isLoadingSettingsValidation,
        validationTimestamp,
        nextStep,
        prevStep,
        closeModal,
        completeSetup,
        restartValidation,
        retrySettingsValidation,
        GetSettingData,
        validateSettings
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