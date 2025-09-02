import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import { useAuthStore } from './auth'
import { useStorage } from '@vueuse/core'
export const useSetupStore = defineStore('setup', () => {
  const authStore = useAuthStore();
  const assetConfig = useStorage('asset-config', {
    assetType_main: '',
    assetName_main: '',
    assetType_sub: '',
    assetName_sub: '',
    assetNickname_main: '',
    assetNickname_sub:'',
  });
  const calib = ref(false)
  const setup = ref(false)
  const applysetup = ref(false)
  const setupFromFile = ref(false)
  const mkva = ref(-1);
  const skva = ref(-1);
  const MainDiagnosis = ref(false)
  const SubDiagnosis = ref(false)
  const MainEnable = ref(false)
  const SubEnable = ref(false)
  const MainPowerQuality = ref(false)
  const SubPowerQuality = ref(false)
  const pf_sign = ref(-1);
  const va_type = ref(-1);
  const unbalance = ref(-1);

  const devLocation = ref(localStorage.getItem('devLocation') || '')

  const initialState = {
    calib: false,
    setup: false,
    applysetup: false,
    setupFromFile: false,
    mkva: -1,
    skva: -1,
    MainDiagnosis: false,
    SubDiagnosis: false,
    MainEnable: false,
    SubEnable: false,
    MainPowerQuality: false,
    SubPowerQuality: false,
    devLocation: '',
    assetConfig: {
      assetType_main: '',
      assetName_main: '',
      assetType_sub: '',
      assetName_sub: '',
      assetNickname_main: '',
      assetNickname_sub:'',
    }
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

  function setMkva(value){
    mkva.value = value
  }
  function setSkva(value){
    skva.value = value
  }

  function setDevLocation(value){
    devLocation.value = value
  }

  function setChannelSetting({ channel, Diagnosis, PowerQuality, Enable }) {
    if (channel === 'Main') {
      MainDiagnosis.value = Diagnosis
      MainEnable.value = Enable
      MainPowerQuality.value = PowerQuality
    } else if (channel === 'Sub') {
      SubDiagnosis.value = Diagnosis
      SubEnable.value = Enable
      SubPowerQuality.value = PowerQuality
    }
  }

  function setDiagSetting({ channel, Diagnosis }) {
    if (channel === 'Main') {
      MainDiagnosis.value = Diagnosis
    } else if (channel === 'Sub') {
      SubDiagnosis.value = Diagnosis
    }
  }

  function setAssetConfig({ channel, name, type, nick }) {
    if (channel === 'Main') {
      assetConfig.value.assetName_main = name;
      assetConfig.value.assetType_main = type;
      assetConfig.value.assetNickname_main = nick;
    } else if (channel === 'Sub') {
      assetConfig.value.assetName_sub = name;
      assetConfig.value.assetType_sub = type;
      assetConfig.value.assetNickname_sub = nick;
    }
  }

  async function checkSetting(forceUpdate = false) {
    if (!forceUpdate && applysetup.value) return
    try {
      const response = await axios.get('/setting/checkSettingFile')
      if (response.data?.result === '1') {
        setChannelSetting({
          channel: 'Main',
          Diagnosis: response.data.Diag_main,
          Enable: response.data.enable_main,
          PowerQuality: response.data.pq_main
        });
        setMkva(response.data.main_kva);
        setChannelSetting({
          channel: 'Sub',
          Diagnosis: response.data.Diag_sub,
          Enable: response.data.enable_sub,
          PowerQuality: response.data.pq_sub
        })
        setSkva(response.data.sub_kva);
        setAssetConfig({
          channel: 'Main',
          name: response.data.assetName_main,
          type: response.data.assetType_main,
          nick : response.data.assetNickname_main
        })
        setAssetConfig({
          channel: 'Sub',
          name: response.data.assetName_sub,
          type: response.data.assetType_sub,
          nick : response.data.assetNickname_sub
        })
        localStorage.setItem('opMode',response.data.mode);
        authStore.setOpMode(response.data.mode)
        localStorage.setItem('devLocation', response.data.location);
        setDevLocation(response.data.location);
        localStorage.setItem('langDefault', response.data.lang);
        authStore.setLangDefault(response.data.lang);
        setPf_sign(response.data.pf_sign);
        setVa_type(response.data.va_type);
        setUnbalance(response.data.unbalance);
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

  function $reset() {
    calib.value = initialState.calib
    setup.value = initialState.setup
    applysetup.value = initialState.applysetup
    setupFromFile.value = initialState.setupFromFile
    mkva.value = initialState.mkva
    skva.value = initialState.skva
    MainDiagnosis.value = initialState.MainDiagnosis
    SubDiagnosis.value = initialState.SubDiagnosis
    MainEnable.value = initialState.MainEnable
    SubEnable.value = initialState.SubEnable
    MainPowerQuality.value = initialState.MainPowerQuality
    SubPowerQuality.value = initialState.SubPowerQuality
    devLocation.value = initialState.devLocation
    
    // ✅ useStorage로 관리되는 assetConfig도 초기화
    Object.assign(assetConfig.value, initialState.assetConfig)
  
    // ✅ localStorage에서도 완전히 제거
    localStorage.removeItem('asset-config') // 이 줄이 빠져있었음
    localStorage.removeItem('devLocation')
    localStorage.removeItem('opMode')
    localStorage.removeItem('langDefault')
    
    console.log('✅ SetupStore reset complete')
  }


  const getCalib = computed(() => calib.value)
  const getSetup = computed(() => setup.value)
  const getsetupFromFile = computed(() => setupFromFile.value)
  const getMainDiagnosis = computed(() => MainDiagnosis.value)
  const getSubDiagnosis = computed(() => SubDiagnosis.value)
  const getSubEnable = computed(() => SubEnable.value)  
  const getAssetConfig = computed(() => assetConfig.value);
  const getMkva = computed(() => mkva.value);
  const getSkva = computed(() => skva.value);
  const getDevLocation = computed(() => devLocation.value);

  const getChannelSetting = computed(() => ({
    MainDiagnosis: MainDiagnosis.value,
    SubDiagnosis: SubDiagnosis.value,
    MainPowerQuality: MainPowerQuality.value,
    SubPowerQuality: SubPowerQuality.value,
    MainEnable: MainEnable.value,
    SubEnable: SubEnable.value
  }))

  const getPf_sign = computed(() => pf_sign.value)
  const getVa_type = computed(() => va_type.value)
  const getUnbalance = computed(() => unbalance.value)


  return {
    calib,
    setup,
    applysetup,
    setupFromFile,
    MainDiagnosis,
    SubDiagnosis,
    MainEnable,
    SubEnable,
    MainPowerQuality,
    SubPowerQuality,
    assetConfig,
    mkva,
    skva,
    devLocation,
    pf_sign,
    va_type,
    unbalance,

    setCalib,
    setSetup,
    setApplySetup,
    setsetupFromFile,

    setChannelSetting,
    setAssetConfig,
    checkSetting,
    setDiagSetting,
    setMkva,
    setSkva,
    setDevLocation,
    setPf_sign,
    setVa_type,
    setUnbalance,

    getCalib,
    getSetup,
    getsetupFromFile,
    getMainDiagnosis,
    getSubDiagnosis,
    getSubEnable,
    getMkva,
    getSkva,
    getAssetConfig,
    getChannelSetting,
    getDevLocation,
    getPf_sign,
    getVa_type,
    getUnbalance,
    $reset,
  }
})