<template>
  <ModalBasic
    id="system-modal"
    :modalOpen="open"
    @close-modal="handleClose"
    title="System"
  >
    <!-- Password Verification -->
    <div v-if="!authenticated" class="px-6 py-5 dark:text-white">
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Admin Password</label>
          <input
            v-model="password"
            type="password"
            class="form-input w-full px-3 py-2 text-sm rounded-lg"
            placeholder="Enter admin password"
            @keyup.enter="verifyPassword"
          />
        </div>
        <div v-if="authError" class="text-sm text-red-500">{{ authError }}</div>
        <div class="flex justify-end">
          <button
            class="btn h-9 px-4 bg-gray-900 text-sm text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white rounded-lg"
            @click="verifyPassword"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>

    <!-- System Content -->
    <div v-else class="px-6 py-5 space-y-6 dark:text-white">

      <!-- Message -->
      <div v-if="message"
        class="flex items-center gap-2 px-4 py-3 rounded-lg text-sm font-medium bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-400 border border-green-200 dark:border-green-800">
        {{ message }}
      </div>

      <!-- Clear Commands -->
      <section v-if="isSetupMode">
        <h4 class="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-3">
          {{ t("config.system.title_1") }}
        </h4>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="cmd in clearCommands"
            :key="cmd.key"
            class="btn h-9 px-4 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white rounded-lg"
            @click="command(cmd.key, 0)"
          >
            {{ t(cmd.label) }}
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-1">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
              <path d="M3 3v5h5" />
            </svg>
          </button>
        </div>
      </section>

      <hr v-if="isSetupMode" class="border-gray-200 dark:border-gray-700/60" />

      <!-- Settings -->
      <section>
        <h4 class="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-3">
          {{ t("config.system.title_3") }}
        </h4>
        <div class="flex flex-wrap gap-2">
          <button
            class="btn h-9 px-4 bg-sky-900 text-xs text-white hover:bg-sky-800 dark:bg-sky-100 dark:text-sky-800 dark:hover:bg-white rounded-lg"
            @click="uploadModalOpen = true"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
              <polyline points="14,2 14,8 20,8" />
              <line x1="16" y1="13" x2="8" y2="13" />
              <line x1="16" y1="17" x2="8" y2="17" />
            </svg>
            &nbsp;{{ t("config.system.import") }}
          </button>
          <button
            class="btn h-9 px-4 bg-yellow-900 text-xs text-white hover:bg-yellow-800 dark:bg-yellow-100 dark:text-yellow-800 dark:hover:bg-white rounded-lg"
            @click="download"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="7,10 12,15 17,10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
            &nbsp;{{ t("config.system.export") }}
          </button>
          <button
            v-if="isSetupMode"
            class="btn h-9 px-4 bg-violet-900 text-xs text-white hover:bg-violet-800 dark:bg-violet-100 dark:text-violet-800 dark:hover:bg-white rounded-lg"
            @click="restore"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
              <rect x="7" y="7" width="3" height="9" />
              <rect x="14" y="7" width="3" height="5" />
            </svg>
            &nbsp;{{ t("config.system.restore") }}
          </button>
          <button
            v-if="isSetupMode"
            class="btn h-9 px-4 bg-red-900 text-xs text-white hover:bg-red-800 dark:bg-red-100 dark:text-red-800 dark:hover:bg-white rounded-lg"
            @click="confirmReset"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 6h18" />
              <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
              <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              <line x1="10" y1="11" x2="10" y2="17" />
              <line x1="14" y1="11" x2="14" y2="17" />
            </svg>
            &nbsp;{{ t("config.system.reset") }}
          </button>
          <button
            v-if="isSetupMode"
            class="btn h-9 px-4 bg-pink-900 text-xs text-white hover:bg-pink-800 dark:bg-pink-100 dark:text-pink-800 dark:hover:bg-white rounded-lg"
            @click="command('reboot', 1)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="12" y1="6" x2="12" y2="12" />
              <path d="M17 7.5a6.5 6.5 0 1 1-10 0" />
            </svg>
            &nbsp;{{ t("config.system.reboot") }}
          </button>
        </div>
      </section>

      <hr class="border-gray-200 dark:border-gray-700/60" />

      <!-- Backup -->
      <section>
        <h4 class="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-3">Backup</h4>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="bk in backupOptions"
            :key="bk.type"
            class="btn h-9 px-4 text-xs text-white rounded-lg"
            :class="bk.btnClass"
            :disabled="isBackingUp"
            @click="downloadBackup(bk.type)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="7,10 12,15 17,10" />
              <line x1="12" y1="15" x2="12" y2="3" />
            </svg>
            &nbsp;{{ isBackingUp && backingUpType === bk.type ? 'Downloading...' : bk.label }}
          </button>
        </div>
      </section>
    </div>

    <!-- File Upload Sub-Modal -->
    <div v-if="uploadModalOpen" class="px-6 pb-5">
      <div class="border border-gray-200 dark:border-gray-700/60 rounded-lg p-4 space-y-3">
        <h4 class="text-sm font-semibold text-gray-800 dark:text-gray-100">Setting file upload</h4>
        <div>
          <label class="block text-sm font-medium mb-1" for="system-filename">
            file path <span class="text-red-500">*</span>
          </label>
          <input
            id="system-filename"
            class="form-input w-full px-2 py-1 text-sm"
            @change="handleFileUpload"
            type="file"
            required
          />
        </div>
        <div v-if="uploadMessage" class="text-sm text-gray-800 dark:text-gray-100">
          {{ uploadMessage }}
        </div>
        <div class="flex justify-end space-x-2">
          <button
            class="btn-sm bg-gray-900 text-gray-100 hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-800 dark:hover:bg-white"
            @click="upload"
          >
            Import
          </button>
          <button
            class="btn-sm border-gray-200 dark:border-gray-700/60 hover:border-gray-300 dark:hover:border-gray-600 text-gray-800 dark:text-white"
            @click="uploadModalOpen = false"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </ModalBasic>
</template>

<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store/auth'
import axios from 'axios'
import ModalBasic from '../../../pages/common/ModalBasic.vue'

const props = defineProps({
  open: { type: Boolean, default: false },
  isSetupMode: { type: Boolean, default: false },
})
const emit = defineEmits(['close'])

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const GetSettingData = inject('GetSettingData', null)

const authenticated = ref(false)
const password = ref('')
const authError = ref('')
const message = ref('')
const uploadMessage = ref('')
const uploadModalOpen = ref(false)
const selectedFile = ref(null)
const isBackingUp = ref(false)
const backingUpType = ref('')

const clearCommands = [
  { key: 'maxmin', label: 'config.system.minmax' },
  { key: 'alarm', label: 'config.system.alarm' },
  { key: 'event', label: 'config.system.eventCount' },
  { key: 'demand', label: 'config.system.demand' },
  { key: 'energy', label: 'config.system.energy' },
]

const backupOptions = [
  { type: 'all', label: 'All (DB + Logs)', btnClass: 'bg-green-700 hover:bg-green-600 dark:bg-green-100 dark:text-green-800 dark:hover:bg-white' },
  { type: 'dbbackup', label: 'Database Only', btnClass: 'bg-blue-700 hover:bg-blue-600 dark:bg-blue-100 dark:text-blue-800 dark:hover:bg-white' },
  { type: 'log', label: 'Logs Only', btnClass: 'bg-amber-700 hover:bg-amber-600 dark:bg-amber-100 dark:text-amber-800 dark:hover:bg-white' },
]

const verifyPassword = async () => {
  authError.value = ''
  try {
    const response = await axios.post('/auth/verify-admin', {
      password: password.value,
    }, { withCredentials: true })
    if (response.data.success) {
      authenticated.value = true
    } else {
      authError.value = response.data.message || 'Invalid password'
    }
  } catch (error) {
    authError.value = error.response?.data?.message || 'Verification failed'
  }
}

const handleClose = () => {
  authenticated.value = false
  password.value = ''
  authError.value = ''
  message.value = ''
  emit('close')
}

const command = async (cmdItem, cmd) => {
  try {
    const typemap = {
      demand: 0, maxmin: 1, energy: 2, alarm: 3, event: 4, reboot: 5, runhour: 6,
    }
    let ch = typemap[cmdItem] === 5 && cmd === 1 ? 2 : 0
    const data = { type: ch, cmd, item: typemap[cmdItem] }
    const response = await axios.post('/setting/command', data, {
      headers: { 'Content-Type': 'application/json' },
      withCredentials: true,
    })
    if (response.data.success) {
      message.value = cmd === 0
        ? `${cmdItem} reset completed (Main)`
        : `${cmdItem} command completed`
    } else {
      message.value = `${cmdItem} failed: ${response.data.message}`
    }
  } catch (error) {
    message.value = `${cmdItem} failed: ${error.message}`
  }
}

const downloadBackup = async (backupType) => {
  isBackingUp.value = true
  backingUpType.value = backupType
  message.value = ''
  try {
    const response = await axios.get(`/setting/backup/download/${backupType}`, {
      responseType: 'blob',
      timeout: 600000,
    })
    const contentDisposition = response.headers['content-disposition']
    let filename = `backup_${backupType}.tar.gz`
    if (contentDisposition) {
      const match = contentDisposition.match(/filename="?(.+?)"?$/)
      if (match) filename = match[1]
    }
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    message.value = `Backup download completed (${backupType})`
  } catch (error) {
    message.value = `Backup failed: ${error.message}`
  } finally {
    isBackingUp.value = false
    backingUpType.value = ''
  }
}

const handleFileUpload = (event) => {
  selectedFile.value = event.target.files[0]
}

const upload = async () => {
  if (!selectedFile.value) {
    uploadMessage.value = 'Please select a file.'
    return
  }
  const formData = new FormData()
  formData.append('file', selectedFile.value)
  try {
    const response = await axios.post('/setting/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    if (response.data.passOK == '1') {
      uploadModalOpen.value = false
      if (GetSettingData) {
        await GetSettingData()
        message.value = 'Upload and reload completed successfully!'
      } else {
        message.value = 'Upload success. Please refresh the page.'
      }
    } else {
      uploadMessage.value = response.data.error
    }
  } catch (error) {
    uploadMessage.value = 'Upload failed: ' + (error.response?.data?.error || error.message)
  }
}

const download = async () => {
  try {
    const response = await axios.get('/setting/download', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'setting.json')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    message.value = 'Download success!'
  } catch (error) {
    message.value = 'Download failed: ' + (error.response?.data?.error || error.message)
  }
}

const restore = async () => {
  try {
    const response = await axios.post('/setting/restoreSetting')
    if (response.data.passOK === '1') {
      if (GetSettingData) {
        await GetSettingData()
        message.value = 'Restore and reload completed successfully!'
      } else {
        message.value = 'Restore success. Please refresh the page.'
      }
    } else {
      message.value = response.data.error || 'Restore failed'
    }
  } catch (error) {
    message.value = 'Restore failed: ' + (error.response?.data?.error || error.message)
  }
}

const confirmReset = () => {
  const confirmed = confirm('All settings will be deleted and reset to default.\nDo you want to proceed?')
  if (confirmed) resetSettings()
}

const resetSettings = async () => {
  try {
    const response = await axios.get('/setting/Reset')
    if (response.data.success) {
      message.value = 'Setup initiated'
      authStore.setInstall(0)
      if (authStore.logout) {
        await authStore.logout()
        router.push('/signin')
      }
    } else {
      message.value = 'Reset failed'
    }
  } catch (error) {
    message.value = 'Reset failed: ' + error.message
  }
}
</script>
