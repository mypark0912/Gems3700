// store/realtime.js
import { defineStore } from 'pinia'
import { ref } from 'vue'  // 🔥 ref 추가
import axios from 'axios'

export const useRealtimeStore = defineStore('realtime', () => {
  // 🔥 state를 ref로 직접 정의
  const meterDictMain = ref({})
  const meterDictSub = ref({})
  const isMainDataLoaded = ref(false)
  const isSubDataLoaded = ref(false)
  const isLoading = ref(false)
  const updateInterval = ref(null)
  const lastUpdateTime = ref(null)
  const error = ref(null)

  // getters
  const isDataReady = (channelState, opMode) => {
    if (!opMode || opMode === '') {
      return false
    }

    if (opMode === 'device2') {
      return isMainDataLoaded.value
    } else {
      const needsSubData = channelState?.SubEnable
      return isMainDataLoaded.value && (needsSubData ? isSubDataLoaded.value : true)
    }
  }

  const getMainData = () => meterDictMain.value || {}
  const getSubData = () => meterDictSub.value || {}
  const getChannelData = (channel) => {
    const data = channel === 'Main' ? meterDictMain.value : meterDictSub.value
    return data || {}
  }

  // actions
  const fetchChannelData = async (channel, opMode, channelState) => {
    if (opMode === '') return false

    const startTime = Date.now()
    //console.log(`[${channel}] 요청 시작:`, new Date().toISOString())

    try {
      const response = await axios.get(`/api/getMeterRedisNew/${channel}/${opMode}`)
      
      //console.log(`[${channel}] 응답 완료: ${Date.now() - startTime}ms`)
      
      if (response.data.success) {
        const newData = response.data.data
        const DashPT = response.data.DashPT;
        // console.log(`[${channel}] 받은 데이터:`, {
        //   U4: newData.U4,
        //   Itot: newData.Itot
        // })
        
        if (channel === 'Main') {
          // 🔥 ref.value에 새 객체 할당
          meterDictMain.value = {
            ...newData,
            P4: newData.P4 ? parseFloat((newData.P4 / 1000).toFixed(2)) : newData.P4
          }
          meterDictMain.value.DashPT = DashPT;
          isMainDataLoaded.value = true
          //console.log('[Main] 스토어 업데이트 완료:', meterDictMain.value.U4)
        } else {
          // 🔥 ref.value에 새 객체 할당
          meterDictSub.value = {
            ...newData,
            P4: newData.P4 ? parseFloat((newData.P4 / 1000).toFixed(2)) : newData.P4
          }
          meterDictSub.value.DashPT = DashPT;
          isSubDataLoaded.value = true
          //console.log('[Sub] 스토어 업데이트 완료:', meterDictSub.value.U4)
        }
        
        lastUpdateTime.value = new Date()
        return true
      } else {
        console.log(`❌ ${channel} 채널: 서버에서 success=false 응답`)
        return false
      }
    } catch (err) {
      console.error(`❌ ${channel} 채널 실패:`, err)
      error.value = err
      return false
    }
  }

  const fetchTempData = async (channel) => {
    try {
      const response = await axios.get(`/api/getAIdata/${channel}`)
      
      if (response.data.success) {
        const newData = response.data.data  // 리스트
        
        if (channel === 'Main') {
          meterDictMain.value = {
            ...meterDictMain.value,
            Temp2: newData
          }
        } else {
          meterDictSub.value = {
            ...meterDictSub.value,
            Temp2: newData
          }
        }
        
        return true
      } else {
        // 실패 시 빈 리스트
        if (channel === 'Main') {
          meterDictMain.value = { ...meterDictMain.value, Temp2: [] }
        } else {
          meterDictSub.value = { ...meterDictSub.value, Temp2: [] }
        }
        return false
      }
    } catch (err) {
      console.error(`❌ ${channel} 온도 데이터 실패:`, err)
      if (channel === 'Main') {
        meterDictMain.value = { ...meterDictMain.value, Temp2: [] }
      } else {
        meterDictSub.value = { ...meterDictSub.value, Temp2: [] }
      }
      return false
    }
  }

  const loadInitialData = async (opMode, channelState) => {
    if (!opMode || opMode === '') {
      console.log('❌ 초기 조건 불만족으로 데이터 로딩 중단')
      return
    }

    //console.log('🔄 초기 데이터 로딩 시작')
    // console.log('🔄 초기 데이터 로딩 시작', {
    //   opMode,
    //   SubEnable: channelState?.SubEnable,
    //   MainEnable: channelState?.MainEnable
    // })
    isLoading.value = true
    error.value = null
    
    try {
      if (channelState?.SubEnable) {
        //console.log('✅ Sub 활성화됨 - 병렬 로딩 시작')
        await Promise.all([
          fetchChannelData('Main', opMode, channelState),
          fetchChannelData('Sub', opMode, channelState),
          //fetchTempData('Main'),
          //fetchTempData('Sub')
        ])
      } else {
        await fetchChannelData('Main', opMode, channelState)
        //await fetchTempData('Main');
        isSubDataLoaded.value = true
      }
      // if (opMode === 'device2') {
      //   await fetchChannelData('Main', opMode, channelState)
      // } else {
        
      // }
    } catch (err) {
      console.error('초기 데이터 로딩 실패:', err)
      error.value = err
    } finally {
      isLoading.value = false
    }
  }

  const startPolling = (opMode, channelState, interval = 10000) => {
    stopPolling()
    loadInitialData(opMode, channelState)

    updateInterval.value = setInterval(async () => {
      if (!opMode) return

      if (channelState?.SubEnable) {
        await Promise.all([
          fetchChannelData('Main', opMode, channelState),
          fetchChannelData('Sub', opMode, channelState),
          //fetchTempData('Main'),
          //fetchTempData('Sub')
        ])
      } else {
        await fetchChannelData('Main', opMode, channelState);
        //await fetchTempData('Main');
      }
    }, interval)
  }

  const stopPolling = () => {
    if (updateInterval.value) {
      clearInterval(updateInterval.value)
      updateInterval.value = null
    }
  }

  const resetData = () => {
    meterDictMain.value = {}
    meterDictSub.value = {}
    isMainDataLoaded.value = false
    isSubDataLoaded.value = false
    error.value = null
  }

  const onChannelStateChange = async (opMode, channelState) => {
    console.log('🔄 채널 상태 변경 감지')
    stopPolling()
    isMainDataLoaded.value = false
    isSubDataLoaded.value = false
    startPolling(opMode, channelState)
  }

  return {
    // state
    meterDictMain,
    meterDictSub,
    isMainDataLoaded,
    isSubDataLoaded,
    isLoading,
    updateInterval,
    lastUpdateTime,
    error,
    // getters
    isDataReady,
    getMainData,
    getSubData,
    getChannelData,
    // actions
    fetchChannelData,
    loadInitialData,
    startPolling,
    stopPolling,
    resetData,
    onChannelStateChange
  }
})