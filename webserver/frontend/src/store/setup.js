import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { useAuthStore } from './auth'

export const useSetupStore = defineStore('setup', () => {
  const authStore = useAuthStore();
  const calib = ref(false)
  const setup = ref(false)
  const applysetup = ref(false)
  const setupFromFile = ref(false)
  const mkva = ref(-1);
  const skva = ref(-1);
  const pf_sign = ref(-1);
  const va_type = ref(-1);
  const unbalance = ref(-1);
  const channelSetting = ref(null);

  // IBSM (device1 only)
  const ibsmTapboxCount = ref(0);
  const ibsmTapboxs = ref([]);

  const devLocation = ref(localStorage.getItem('devLocation') || '')

  const initialState = {
    calib: false,
    setup: false,
    applysetup: false,
    setupFromFile: false,
    devLocation: '',
  }

  function setPf_sign(value) {
    pf_sign.value = value
  }

  function setVa_type(value) {
    va_type.value = value
  }

  function setUnbalance(value) {
    unbalance.value = value
  }

  function setCalib(value) {
    calib.value = value
  }

  function setSetup(status) {
    setup.value = status
  }

  function setApplySetup(status) {
    applysetup.value = status
  }

  function setsetupFromFile(status) {
    setupFromFile.value = status
  }

  function setDevLocation(value){
    devLocation.value = value
  }

  async function checkSetting(forceUpdate = false) {
    if (!forceUpdate && applysetup.value) return
    try {
      const response = await axios.get('/setting/checkSettingFile')
      if (response.data?.result === '1') {
        localStorage.setItem('opMode', response.data.mode);
        authStore.setOpMode(response.data.mode)
        localStorage.setItem('devLocation', response.data.location);
        setDevLocation(response.data.location);
        localStorage.setItem('langDefault', response.data.lang);
        authStore.setLangDefault(response.data.lang);
        setPf_sign(response.data.pf_sign);
        setVa_type(response.data.va_type);
        setUnbalance(response.data.unbalance);
        channelSetting.value = {
          MainEnable: response.data.enable_main,
          SubEnable: response.data.enable_sub,
        }
        if (response.data.mode === 'device1') {
          await fetchIbsm()
        }
        setSetup(true)
        setApplySetup(true)
      } else {
        setSetup(false)
        setApplySetup(true)
      }
    } catch (err) {
      console.error('checkSetting Error:', err)
    }
  }

  async function fetchIbsm() {
    try {
      const res = await axios.get('/setting/checkIbsm')
      if (res.data?.result === '1') {
        ibsmTapboxCount.value = res.data.tapboxCount
        ibsmTapboxs.value = res.data.tapboxs
      }
    } catch (err) {
      console.error('fetchIbsm Error:', err)
    }
  }

  function $reset() {
    calib.value = initialState.calib
    setup.value = initialState.setup
    applysetup.value = initialState.applysetup
    setupFromFile.value = initialState.setupFromFile
    devLocation.value = initialState.devLocation
    pf_sign.value = -1
    va_type.value = -1
    unbalance.value = -1
    channelSetting.value = null
    ibsmTapboxCount.value = 0
    ibsmTapboxs.value = []

    localStorage.removeItem('devLocation')
    localStorage.removeItem('opMode')
    localStorage.removeItem('langDefault')
  }

  const getCalib = computed(() => calib.value)
  const getSetup = computed(() => setup.value)
  const getsetupFromFile = computed(() => setupFromFile.value)
  const getDevLocation = computed(() => devLocation.value)
  const getMkva = computed(() => mkva.value)
  const getSkva = computed(() => skva.value)
  const getPf_sign = computed(() => pf_sign.value)
  const getVa_type = computed(() => va_type.value)
  const getUnbalance = computed(() => unbalance.value)
  const getChannelSetting = computed(() => channelSetting.value)
  const getIbsmTapboxCount = computed(() => ibsmTapboxCount.value)
  const getIbsmTapboxs = computed(() => ibsmTapboxs.value)

  return {
    calib,
    setup,
    applysetup,
    setupFromFile,
    devLocation,
    mkva,
    skva,
    pf_sign,
    va_type,
    unbalance,
    channelSetting,
    ibsmTapboxCount,
    ibsmTapboxs,

    setCalib,
    setSetup,
    setApplySetup,
    setsetupFromFile,
    setDevLocation,
    setPf_sign,
    setVa_type,
    setUnbalance,

    checkSetting,
    fetchIbsm,

    getCalib,
    getSetup,
    getsetupFromFile,
    getDevLocation,
    getMkva,
    getSkva,
    getPf_sign,
    getVa_type,
    getUnbalance,
    getChannelSetting,
    getIbsmTapboxCount,
    getIbsmTapboxs,
    $reset,
  }
})
