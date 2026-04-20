<template>
  <transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition ease-out duration-100" leave-from-class="opacity-100" leave-to-class="opacity-0">
    <div v-show="modalOpen" class="fixed inset-0 bg-gray-900 bg-opacity-30 z-50 transition-opacity" aria-hidden="true"></div>
  </transition>

  <transition enter-active-class="transition ease-in-out duration-200" enter-from-class="opacity-0 translate-y-4" enter-to-class="opacity-100 translate-y-0" leave-active-class="transition ease-in-out duration-200" leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-4">
    <div v-show="modalOpen" class="fixed inset-0 z-50 overflow-hidden flex items-center my-4 justify-center px-4 sm:px-6" role="dialog" aria-modal="true">
      <div ref="modalContent" class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-auto max-w-2xl w-full max-h-full">
        <div class="p-6">
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ t('saveModal.title') }}</h2>
          </div>

          <!-- Progress Steps -->
          <div class="mb-6">
            <div class="flex items-center justify-between">
              <template v-for="(step, index) in visibleSteps" :key="step.id">
                <div class="flex flex-col items-center">
                  <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-semibold transition-all duration-300" :class="getStepClass(step)">
                    <svg v-if="step.status === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                    <svg v-else-if="step.status === 'error'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" /></svg>
                    <svg v-else-if="step.status === 'processing'" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                    <span v-else>{{ index + 1 }}</span>
                  </div>
                  <span class="text-xs mt-1 text-center max-w-[80px]" :class="getStepLabelClass(step)">{{ step.name }}</span>
                </div>
                <div v-if="index < visibleSteps.length - 1" class="flex-1 h-0.5 mx-2" :class="getConnectorClass(index)"></div>
              </template>
            </div>
          </div>

          <!-- Step 1: Validation -->
          <div v-if="currentStep === 1" class="step-content">
            <div v-if="isProcessing" class="text-center py-8">
              <div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-violet-100 dark:bg-violet-900/20 rounded-full">
                <svg class="animate-spin w-8 h-8 text-violet-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-300">
                {{ isSaving ? t('saveModal.messages.saving') : t('saveModal.validation.validating') }}
              </div>
            </div>

            <div v-else class="space-y-4">
              <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                <h3 class="font-semibold mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                  {{ t('saveModal.validation.summary') }}
                </h3>
                <div class="space-y-2 text-sm">
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">{{ t('saveModal.validation.overallStatus') }}:</span>
                    <span :class="validationResult.isValid ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'">
                      {{ validationResult.isValid ? t('saveModal.validation.valid') : t('saveModal.validation.invalid') }}
                    </span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">{{ t('saveModal.validation.errors') }}:</span>
                    <span :class="validationResult.errors.length > 0 ? 'text-red-600 font-semibold' : 'text-green-600'">{{ validationResult.errors.length }}</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-gray-600 dark:text-gray-400">{{ t('saveModal.validation.warnings') }}:</span>
                    <span :class="validationResult.warnings.length > 0 ? 'text-yellow-600 font-semibold' : 'text-green-600'">{{ validationResult.warnings.length }}</span>
                  </div>
                </div>
              </div>

              <div v-if="validationResult.errors.length > 0" class="bg-red-50 dark:bg-red-900/20 rounded-lg p-4">
                <h3 class="font-semibold mb-3 text-red-700 dark:text-red-400 flex items-center">
                  <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>
                  {{ t('saveModal.validation.criticalErrors') }}
                </h3>
                <ul class="space-y-1 text-sm max-h-40 overflow-y-auto">
                  <li v-for="(error, idx) in validationResult.errors" :key="idx" class="text-red-700 dark:text-red-400">• {{ error }}</li>
                </ul>
              </div>

              <div v-if="validationResult.warnings.length > 0" class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
                <h3 class="font-semibold mb-3 text-yellow-700 dark:text-yellow-400 flex items-center">
                  <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>
                  {{ t('saveModal.validation.warnings') }}
                </h3>
                <ul class="space-y-1 text-sm max-h-40 overflow-y-auto">
                  <li v-for="(warning, idx) in validationResult.warnings" :key="idx" class="text-yellow-700 dark:text-yellow-400">• {{ warning }}</li>
                </ul>
              </div>

              <div v-if="validationResult.isValid && validationResult.errors.length === 0" class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                  <span class="text-green-700 dark:text-green-400 font-medium">{{ t('saveModal.validation.allPassed') }}</span>
                </div>
              </div>
            </div>

            <!-- Nameplate 변경 경고 -->
            <div v-if="showNameplateWarning" class="mt-4 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-300 dark:border-yellow-700 rounded-lg">
              <div class="flex items-start">
                <svg class="w-5 h-5 text-yellow-600 dark:text-yellow-400 mr-3 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
                <div class="flex-1">
                  <h4 class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2">{{ t('saveModal.nameplateWarning.title') }}</h4>
                  <p class="text-sm text-yellow-700 dark:text-yellow-300 mb-3">
                    {{ nameplateWarningMessage }}
                  </p>
                  <label class="flex items-center cursor-pointer">
                    <input type="checkbox" v-model="nameplateConfirmed" class="w-4 h-4 text-yellow-600 border-yellow-300 rounded focus:ring-yellow-500" />
                    <span class="ml-2 text-sm text-yellow-800 dark:text-yellow-200 font-medium">{{ t('saveModal.nameplateConfirm') }}</span>
                  </label>
                </div>
              </div>
            </div>

            <div class="flex justify-between mt-6">
              <button class="btn px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 rounded-lg" :disabled="isProcessing" @click="cancelAndReload">
                {{ t('saveModal.buttons.cancel') }}
              </button>
              <button v-if="!isProcessing" class="btn px-4 py-2 rounded-lg" :class="canProceedToNext ? 'bg-blue-600 text-white hover:bg-blue-700' : 'bg-gray-400 text-gray-600 cursor-not-allowed'" :disabled="!canProceedToNext" @click="handleStep1Next">
                {{ useDiagnosis ? t('saveModal.buttons.next') : t('saveModal.buttons.save') }}
              </button>
            </div>
          </div>

          <!-- Step 2: Diagnosis Config -->
          <div v-if="currentStep === 2" class="step-content">
            <div v-if="isProcessing" class="text-center py-8">
              <div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-violet-100 dark:bg-violet-900/20 rounded-full">
                <svg class="animate-spin w-8 h-8 text-violet-500" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-300">{{ processingMessage }}</div>
            </div>

            <div v-else class="space-y-4">
              <h3 class="font-semibold text-gray-800 dark:text-gray-100 mb-4">{{ t('saveModal.diagnosis.results') }}</h3>
              
              <div v-for="result in diagnosisResults" :key="result.id" class="rounded-lg p-4" :class="result.status === 'error' ? 'bg-red-50 dark:bg-red-900/20' : result.status === 'success' ? 'bg-green-50 dark:bg-green-900/20' : 'bg-gray-50 dark:bg-gray-700/50'">
                <div class="flex items-start">
                  <div class="flex-shrink-0 mr-3">
                    <svg v-if="result.status === 'success'" class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                    <svg v-else-if="result.status === 'error'" class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>
                    <svg v-else class="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z" clip-rule="evenodd" /></svg>
                  </div>
                  <div class="flex-1">
                    <div class="font-medium text-sm" :class="result.status === 'error' ? 'text-red-700 dark:text-red-400' : result.status === 'success' ? 'text-green-700 dark:text-green-400' : 'text-gray-700 dark:text-gray-300'">{{ result.name }}</div>
                    <div v-if="result.message" class="text-xs mt-1" :class="result.status === 'error' ? 'text-red-600' : 'text-gray-500'">{{ result.message }}</div>
                    <ul v-if="result.errors && result.errors.length > 0" class="mt-2 text-xs space-y-1">
                      <li v-for="(err, idx) in result.errors" :key="idx" class="text-red-600">• {{ err }}</li>
                    </ul>
                  </div>
                </div>
              </div>

              <div class="mt-4 p-4 rounded-lg" :class="hasDiagnosisErrors ? 'bg-red-100 dark:bg-red-900/30 border border-red-300' : 'bg-green-100 dark:bg-green-900/30 border border-green-300'">
                <div class="flex items-center">
                  <svg v-if="hasDiagnosisErrors" class="w-6 h-6 text-red-500 mr-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>
                  <svg v-else class="w-6 h-6 text-green-500 mr-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                  <div>
                    <div class="font-semibold" :class="hasDiagnosisErrors ? 'text-red-700' : 'text-green-700'">
                      {{ hasDiagnosisErrors ? t('saveModal.diagnosis.configFailed') : t('saveModal.diagnosis.configSuccess') }}
                    </div>
                    <div class="text-sm" :class="hasDiagnosisErrors ? 'text-red-600' : 'text-green-600'">
                      {{ hasDiagnosisErrors ? t('saveModal.diagnosis.fixErrors') : t('saveModal.diagnosis.readyToSave') }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex justify-between mt-6">
              <button class="btn px-4 py-2 bg-gray-500 text-white hover:bg-gray-600 rounded-lg" :disabled="isProcessing" @click="cancelAndReload">
                {{ t('saveModal.buttons.cancel') }}
              </button>
              <button v-if="!isProcessing && !hasDiagnosisErrors" class="btn px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-lg" @click="saveAndClose">
                {{ t('saveModal.buttons.save') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { ref, computed, watch, inject } from "vue";
import axios from "axios";
import { settingValidator } from "@/utils/validation.js";
import { useI18n } from "vue-i18n";

export default {
  name: "SaveProgressModal",
  props: { modalOpen: { type: Boolean, default: false } },
  emits: ["close-modal", "save-complete"],
  setup(props, { emit }) {
    const { t } = useI18n();
    const inputDict = inject("inputDict");
    const channel_main = inject("channel_main");
    const channel_sub = inject("channel_sub");
    const useDiagnosis = inject("useDiagnosis");
    const diagnosisData = inject("diagnosisData");
    const advancedData = inject("advancedData");
    const diagnosis_detail = inject("diagnosis_detail");
    const devMode = inject("devMode");
    const devLang = inject("devLang");
    const currentStep = ref(1);
    const isProcessing = ref(false);
    const isSaving = ref(false);
    const processingMessage = ref("");
    const checkNameplateResult = ref({ main: false, sub: false });

    const steps = ref([
      { id: 1, name: t('saveModal.validation.title') || "Validation", status: "current" },
      { id: 2, name: t('saveModal.diagnosis.title') || "Diagnosis Config", status: "pending" },
    ]);

    const validationResult = ref({ isValid: false, errors: [], warnings: [] });
    const diagnosisResults = ref([]);

    const showNameplateWarning = ref(false);
    const nameplateWarningMessage = ref("");
    const nameplateConfirmed = ref(false);

    const canProceedToNext = computed(() => {
      if (!validationResult.value.isValid) return false;
      if (showNameplateWarning.value && !nameplateConfirmed.value) return false;
      return true;
    });

    const visibleSteps = computed(() => {
      if (useDiagnosis.value) {
        return steps.value;
      }
      return [steps.value[0]];
    });

    const hasDiagnosisErrors = computed(() => diagnosisResults.value.some((r) => r.status === "error"));

    const getStepClass = (step) => {
      if (step.status === "success") return "bg-green-500 text-white";
      if (step.status === "error") return "bg-red-500 text-white";
      if (step.status === "current" || step.status === "processing") return "bg-violet-500 text-white";
      return "bg-gray-200 dark:bg-gray-700 text-gray-500";
    };

    const getStepLabelClass = (step) => {
      if (step.status === "error") return "text-red-600";
      if (step.status === "success") return "text-green-600";
      if (step.status === "current") return "text-violet-600 font-semibold";
      return "text-gray-500";
    };

    const getConnectorClass = (index) => {
      const allSteps = visibleSteps.value;
      const nextStep = allSteps[index + 1];
      if (nextStep && (nextStep.status === "success" || nextStep.status === "current")) return "bg-violet-500";
      return "bg-gray-200 dark:bg-gray-700";
    };

    const getChannelLabel = (channelName) => {
      return channelName === "main" 
        ? t('common.mainChannel') || 'Main channel'
        : t('common.subChannel') || 'Sub channel';
    };

    const resetState = () => {
      currentStep.value = 1;
      isProcessing.value = false;
      isSaving.value = false;
      processingMessage.value = "";
      checkNameplateResult.value = { main: false, sub: false };
      showNameplateWarning.value = false;
      nameplateWarningMessage.value = "";
      nameplateConfirmed.value = false;
      steps.value = [
        { id: 1, name: t('saveModal.validation.title') || "Validation", status: "current" },
        { id: 2, name: t('saveModal.diagnosis.title') || "Diagnosis Config", status: "pending" },
      ];
      validationResult.value = { isValid: false, errors: [], warnings: [] };
      diagnosisResults.value = [];
    };

    const runValidation = async () => {
      isProcessing.value = true;
      try {
        const result = settingValidator.validateAllSettings(inputDict.value, channel_main.value, channel_sub.value);

        if (useDiagnosis.value) {
          const diagErrors = validateDiagnosisAssets();
          if (diagErrors.length > 0) {
            result.errors.push(...diagErrors);
            result.isValid = false;
          }

          await checkNameplateChanges();
        }

        validationResult.value = result;
        steps.value[0].status = result.isValid ? "success" : "error";
      } catch (error) {
        validationResult.value = { isValid: false, errors: ["Validation failed: " + error.message], warnings: [] };
        steps.value[0].status = "error";
      } finally {
        isProcessing.value = false;
      }
    };

    const checkNameplateChanges = async () => {
      let changedAssets = [];
      checkNameplateResult.value = { main: false, sub: false };

      for (const channelName in diagnosis_detail.value) {
        const channelData = diagnosis_detail.value[channelName];
        const isDiagEnabled = channelName === "main" 
          ? inputDict.value.useFunction?.diagnosis_main 
          : inputDict.value.useFunction?.diagnosis_sub;

        if (channelData?.use && isDiagEnabled && 
            channelData.assetName && channelData.assetName !== "" && 
            channelData.tableData?.length > 0) {
          try {
            const combinedData = [...channelData.tableData, ...(channelData.modalData || [])];
            const changed = await checkNameplateConfig(combinedData, channelData.assetName);
            if (changed) {
              checkNameplateResult.value[channelName] = true;
              
              const assetType = channelName === "main" 
                ? channel_main.value.assetInfo?.type 
                : channel_sub.value.assetInfo?.type;
              
              if (assetType) {
                const channelLabel = getChannelLabel(channelName);
                changedAssets.push({ channel: channelLabel, type: assetType });
              }
            }
          } catch (error) {
            console.error(`Error checking nameplate for ${channelName}:`, error);
          }
        }
      }

      if (changedAssets.length > 0) {
        showNameplateWarning.value = true;
        
        if (changedAssets.length === 1) {
          const asset = changedAssets[0];
          nameplateWarningMessage.value = t('saveModal.nameplateWarning.single', {
            channel: asset.channel,
            type: asset.type
          }) || `Changes to the ${asset.channel} ${asset.type}'s nameplate parameters invalidate previous calculations. All stored data for ${asset.type} will be deleted and rebuilt from scratch.`;
        } else {
          const assetDescriptions = changedAssets.map(asset => 
            `"${asset.channel}'s ${asset.type}"`
          ).join(` ${t('common.and') || 'and'} `);
          
          nameplateWarningMessage.value = t('saveModal.nameplateWarning.multiple', {
            assets: assetDescriptions
          }) || `Changes to ${assetDescriptions} nameplate parameters invalidate previous calculations. All stored data for these assets will be deleted and rebuilt from scratch.`;
        }
      } else {
        showNameplateWarning.value = false;
        nameplateWarningMessage.value = "";
      }
    };

    const validateDiagnosisAssets = () => {
      const errors = [];
      if (!diagnosis_detail.value) return errors;
      for (const channelName in diagnosis_detail.value) {
        const channelData = diagnosis_detail.value[channelName];
        const isDiagEnabled = channelName === "main" ? inputDict.value.useFunction?.diagnosis_main : inputDict.value.useFunction?.diagnosis_sub;
        if (!isDiagEnabled && channelData?.assetName && channelData.assetName !== "") {
          errors.push(`Diagnosis disabled for ${channelName} channel, but asset (${channelData.assetName}) is registered`);
        }
      }
      return errors;
    };

    const handleStep1Next = async () => {
      if (useDiagnosis.value) {
        currentStep.value = 2;
        steps.value[1].status = "current";
        await runDiagnosisConfig();
      } else {
        await saveAndClose();
      }
    };

    const runDiagnosisConfig = async () => {
      isProcessing.value = true;
      diagnosisResults.value = [];

      try {
        processingMessage.value = t('saveModal.diagnosis.configuring') || "Configuring nameplate settings...";
        const nameplateResult = await runNameplateConfig();
        diagnosisResults.value.push(nameplateResult);
        if (nameplateResult.status === "error") { steps.value[1].status = "error"; isProcessing.value = false; return; }

        processingMessage.value = t('saveModal.diagnosis.settingParams') || "Setting asset parameters...";
        const paramsResult = await runAssetParams();
        diagnosisResults.value.push(paramsResult);
        if (paramsResult.status === "error") { steps.value[1].status = "error"; isProcessing.value = false; return; }

        processingMessage.value = t('saveModal.diagnosis.savingSettings') || "Saving diagnosis settings...";
        const diagResult = await runDiagnosisSettings();
        diagnosisResults.value.push(diagResult);
        if (diagResult.status === "error") { steps.value[1].status = "error"; isProcessing.value = false; return; }

        processingMessage.value = t('saveModal.diagnosis.savingProfile') || "Saving advanced profile...";
        const advResult = await runAdvancedProfile();
        diagnosisResults.value.push(advResult);
        steps.value[1].status = advResult.status === "error" ? "error" : "success";
      } catch (error) {
        diagnosisResults.value.push({ id: "error", name: "Diagnosis Configuration", status: "error", message: error.message, errors: [error.message] });
        steps.value[1].status = "error";
      } finally {
        isProcessing.value = false;
      }
    };

    const runNameplateConfig = async () => {
      const result = { id: "nameplate", name: t('saveModal.diagnosis.nameplate') || "Nameplate Configuration", status: "success", message: "", errors: [] };

      try {
        for (const channelName in diagnosis_detail.value) {
          const channelData = diagnosis_detail.value[channelName];
          const isDiagEnabled = channelName === "main" ? inputDict.value.useFunction?.diagnosis_main : inputDict.value.useFunction?.diagnosis_sub;

          if (channelData?.use && isDiagEnabled && channelData.assetName && channelData.assetName !== "" && channelData.tableData?.length > 0) {
            const combinedData = [...channelData.tableData, ...(channelData.modalData || [])];

            const plainData = combinedData.map((item) => ({ ...item }));
            const response = await axios.post(`/setting/setAssetConfig/${channelData.assetName}`, plainData, { headers: { "Content-Type": "application/json" } });

            const { status, success, error } = response.data || {};
            if (status === "0" || status === 0 || ((status === "1" || status === 1) && !success)) {
              result.errors.push(`[${channelName}] ${Array.isArray(error) ? error.join(", ") : error || "Unknown error"}`);
              result.status = "error";
            }
          }
        }
        const hasAnyChange = checkNameplateResult.value.main || checkNameplateResult.value.sub;
        if (result.status !== "error") result.message = hasAnyChange ? "Nameplate updated (requires recommissioning)" : "Nameplate saved";
      } catch (error) {
        result.status = "error";
        result.errors.push(error.message);
      }
      return result;
    };

    const checkNameplateConfig = async (tableData, assetName) => {
      try {
        if (!tableData?.length || !assetName) return false;
        const plainData = tableData.map((item) => ({ ...item }));
        const response = await axios.post(`/setting/checkAssetConfig/${assetName}`, plainData, { headers: { "Content-Type": "application/json" } });
        return response.data.status == 1 && response.data.success ? response.data.result : false;
      } catch {
        return false;
      }
    };

    const runAssetParams = async () => {
      const result = { id: "params", name: t('saveModal.diagnosis.parameters') || "Asset Parameters", status: "success", message: "Parameters configured", errors: [] };
      try {
        for (const channelName in diagnosis_detail.value) {
          const channelData = diagnosis_detail.value[channelName];
          const isDiagEnabled = channelName === "main" ? inputDict.value.useFunction?.diagnosis_main : inputDict.value.useFunction?.diagnosis_sub;

          if (channelData?.use && isDiagEnabled && channelData.assetName && channelData.paramData?.length > 0) {
            const response = await axios.post(`/setting/setAssetParams/${channelData.assetName}`, channelData.paramData.map((item) => ({ ...item })), { headers: { "Content-Type": "application/json" } });
            const { status, success, error } = response.data || {};
            if (status === "0" || status === 0 || ((status === "1" || status === 1) && !success)) {
              result.errors.push(`[${channelName}] ${Array.isArray(error) ? error.join(", ") : error || "Unknown error"}`);
              result.status = "error";
            }
          }
        }
      } catch (error) { result.status = "error"; result.errors.push(error.message); }
      return result;
    };

    const runDiagnosisSettings = async () => {
      const result = { id: "diagSettings", name: t('saveModal.diagnosis.settings') || "Diagnosis Settings", status: "success", message: "Settings saved", errors: [] };
      try {
        if (diagnosisData.value && Object.keys(diagnosisData.value).length > 0) {
          const response = await axios.post("/setting/setDiagnosisSetting", diagnosisData.value, { headers: { "Content-Type": "application/json" }, withCredentials: true });
          const { status, success, error } = response.data;
          if (status === "0" || status === 0 || ((status === "1" || status === 1) && !success)) {
            result.errors.push(Array.isArray(error) ? error.join(", ") : error || "Save failed");
            result.status = "error";
          }
        } else { result.message = "No settings to save"; result.status = "skipped"; }
      } catch (error) { result.status = "error"; result.errors.push(error.message); }
      return result;
    };

    const runAdvancedProfile = async () => {
      const result = { id: "advanced", name: t('saveModal.diagnosis.profile') || "Advanced Profile", status: "success", message: "Profile saved", errors: [] };
      try {
        if (advancedData.value && Object.keys(advancedData.value).length > 0) {
          const response = await axios.post("/setting/setDiagnosisProfile", advancedData.value, { headers: { "Content-Type": "application/json" }, withCredentials: true });
          const { status, success, error } = response.data;
          if (status === "0" || status === 0 || ((status === "1" || status === 1) && !success)) {
            result.errors.push(Array.isArray(error) ? error.join(", ") : error || "Save failed");
            result.status = "error";
          }
        } else { result.message = "No profile to save"; result.status = "skipped"; }
      } catch (error) { result.status = "error"; result.errors.push(error.message); }
      return result;
    };

    const saveAndClose = async () => {
      isProcessing.value = true;
      isSaving.value = true;
      processingMessage.value = t('saveModal.messages.saving') || "Saving configuration...";

      try {
        const formattedData = transFormat();

        const response = await axios.post("/setting/savefileNew", formattedData, { headers: { "Content-Type": "application/json;charset=utf-8" }, withCredentials: true });

        if (response.data?.status === "1") {
          alert(t('saveModal.messages.saveSuccess') || "Configuration saved successfully!");
          emit("save-complete", { success: true });
          emit("close-modal");
        } else {
          alert((t('saveModal.messages.saveFailed') || "Save failed") + ": " + (response.data?.error || "Unknown error"));
        }
      } catch (error) {
        alert((t('saveModal.messages.saveFailed') || "Save failed") + ": " + error.message);
      } finally {
        isProcessing.value = false;
        isSaving.value = false;
      }
    };

    const transFormat = () => {
      const generalData = { ...inputDict.value };
      if (generalData.useFunction) {
        generalData.useFunction.ftp = generalData.useFunction.ftp === true || generalData.useFunction.ftp === 1 ? 1 : 0;
        generalData.useFunction.sntp = generalData.useFunction.sntp === true || generalData.useFunction.sntp === 1 ? 1 : 0;
      }
      if (generalData.modbus?.rtu_use !== undefined) generalData.modbus.rtu_use = generalData.modbus.rtu_use === true || generalData.modbus.rtu_use === 1 ? 1 : 0;

      const mainData = transNumber(channel_main.value);
      mainData.channel = "Main";
      const subData = transNumber(channel_sub.value);
      subData.channel = "Sub";

      if (devMode.value === "device0") {
        const params = ["Temperature", "Frequency", "Line Voltage", "Phase Voltage", "Current", "Unbalance", "PF", "THD", "TDD", "Power"];
        if (!mainData.trendInfo) mainData.trendInfo = {};
        mainData.trendInfo.params = [...params];
        if (!subData.trendInfo) subData.trendInfo = {};
        subData.trendInfo.params = [...params];
      }

      return { mode: devMode.value, lang: devLang.value, General: generalData, channel: [mainData, subData] };
    };

    const transNumber = (channelData) => {
      const data = { ...channelData };
      ["Enable", "PowerQuality", "useDO", "confStatus", "useAI"].forEach((field) => {
        if (data.hasOwnProperty(field)) data[field] = data[field] === true || data[field] === 1 ? 1 : 0;
      });
      if (data.confStatus === 0 && data.status_Info) {
        data.status_Info.diagnosis = []; data.status_Info.pq = []; data.status_Info.faults = []; data.status_Info.events = [];
      }
      if (data.ctInfo) {
        for (const key in data.ctInfo) {
          if (key === "direction") data.ctInfo[key] = data.ctInfo.direction.map((d) => parseInt(d));
          else if (key === "inorminal") data.ctInfo[key] = parseFloat(data.ctInfo[key]);
          else data.ctInfo[key] = parseInt(data.ctInfo[key]);
        }
      }
      if (data.ptInfo) for (const key in data.ptInfo) data.ptInfo[key] = parseInt(data.ptInfo[key]);
      if (data.eventInfo) for (const key in data.eventInfo) data.eventInfo[key] = parseInt(data.eventInfo[key]);
      return data;
    };

    const cancelAndReload = async () => {
      try {
        const response = await axios.get(`/setting/getDiagnosisSetting`);
        if (response.data.success && diagnosisData.value) {
          Object.assign(diagnosisData.value, response.data.data);
        }
      } catch (e) { console.error(e); }
      finally { emit("close-modal"); }
    };

    watch(() => props.modalOpen, async (newVal) => { if (newVal) { resetState(); await runValidation(); } });

    return { 
      currentStep, steps, visibleSteps, isProcessing, isSaving, processingMessage, 
      validationResult, diagnosisResults, hasDiagnosisErrors, useDiagnosis,
      showNameplateWarning, nameplateWarningMessage, nameplateConfirmed, canProceedToNext,
      getStepClass, getStepLabelClass, getConnectorClass, 
      handleStep1Next, saveAndClose, cancelAndReload,devLang,
      t
    };
  },
};
</script>

<style scoped>
.step-content { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.btn { @apply font-medium transition-colors duration-200; }
</style>