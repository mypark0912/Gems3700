<template>
  <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition ease-out duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
    <div v-show="modalOpen" class="fixed inset-0 bg-gray-900 bg-opacity-30 z-50 transition-opacity" aria-hidden="true"></div>
  </transition>

  <transition enter-active-class="transition ease-in-out duration-200" enter-from-class="opacity-0 translate-y-4" enter-to-class="opacity-100 translate-y-0" leave-active-class="transition ease-in-out duration-200" leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-4">
    <div v-show="modalOpen" class="fixed inset-0 z-50 overflow-hidden flex items-center my-4 justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
      <div ref="modalContent" class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-auto max-w-lg w-full max-h-full">
        <div class="p-6">
          <h2 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-4">Validation Result</h2>

          <!-- Summary -->
          <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 mb-4">
            <div class="space-y-2 text-sm">
              <div class="flex items-center justify-between">
                <span class="text-gray-600 dark:text-gray-400">Status:</span>
                <span :class="validationResult.isValid ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'">
                  {{ validationResult.isValid ? 'PASS' : 'FAIL' }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600 dark:text-gray-400">Errors:</span>
                <span :class="validationResult.errors.length > 0 ? 'text-red-600 font-semibold' : 'text-green-600'">{{ validationResult.errors.length }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600 dark:text-gray-400">Warnings:</span>
                <span :class="validationResult.warnings.length > 0 ? 'text-yellow-600 font-semibold' : 'text-green-600'">{{ validationResult.warnings.length }}</span>
              </div>
            </div>
          </div>

          <!-- Errors -->
          <div v-if="validationResult.errors.length > 0" class="bg-red-50 dark:bg-red-900/20 rounded-lg p-4 mb-4">
            <h3 class="font-semibold mb-2 text-red-700 dark:text-red-400 flex items-center text-sm">
              <svg class="w-4 h-4 mr-1.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>
              Errors
            </h3>
            <ul class="space-y-1 text-sm max-h-40 overflow-y-auto">
              <li v-for="(error, idx) in validationResult.errors" :key="idx" class="text-red-700 dark:text-red-400">• {{ error }}</li>
            </ul>
          </div>

          <!-- Warnings -->
          <div v-if="validationResult.warnings.length > 0" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4 mb-4">
            <h3 class="font-semibold mb-2 text-yellow-700 dark:text-yellow-400 flex items-center text-sm">
              <svg class="w-4 h-4 mr-1.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
              Warnings
            </h3>
            <ul class="space-y-1 text-sm max-h-40 overflow-y-auto">
              <li v-for="(warning, idx) in validationResult.warnings" :key="idx" class="text-yellow-700 dark:text-yellow-400">• {{ warning }}</li>
            </ul>
          </div>

          <!-- All Passed -->
          <div v-if="validationResult.isValid && validationResult.errors.length === 0" class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4 mb-4">
            <div class="flex items-center">
              <svg class="w-5 h-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
              <span class="text-green-700 dark:text-green-400 font-medium">All validations passed.</span>
            </div>
          </div>

          <!-- Buttons -->
          <div class="flex justify-end gap-3 mt-2">
            <button
              v-if="validationResult.isValid"
              class="btn px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-lg text-sm"
              :disabled="isSaving"
              @click="doSave"
            >
              {{ isSaving ? 'Saving...' : 'Save' }}
            </button>
            <button class="btn px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 rounded-lg text-sm" @click="$emit('close')">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch } from 'vue'
import { settingValidator } from '@/utils/validation.js'

const props = defineProps({
  modalOpen: { type: Boolean, default: false },
  inputDict: { type: Object, required: true },
  channelMain: { type: Object, required: true },
})

const emit = defineEmits(['close', 'save'])

const validationResult = ref({ isValid: false, errors: [], warnings: [] })
const isSaving = ref(false)

watch(() => props.modalOpen, (val) => {
  if (val) {
    const result = settingValidator.validateAllSettings(props.inputDict, props.channelMain)
    validationResult.value = result
  }
})

const doSave = () => {
  isSaving.value = true
  emit('save', {
    done: () => { isSaving.value = false }
  })
}
</script>

<style scoped>
.btn { @apply font-medium transition-colors duration-200; }
</style>
