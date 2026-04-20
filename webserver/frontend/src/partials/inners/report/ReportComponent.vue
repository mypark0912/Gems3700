<template>
  <div class="col-span-full xl:col-span-12 bg-white dark:bg-gray-800 shadow-sm rounded-xl">
    <div
      class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg mb-6 mt-6"
    >
      <div
        class="absolute top-0 left-0 right-0 h-0.5 bg-blue-500"
        aria-hidden="true"
      ></div>
      <div
        class="px-5 pt-6 pb-6 border-b border-gray-200 dark:border-gray-700/60"
      >
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-blue-500 mr-3">
            <svg
              class="w-6 h-6 fill-current text-white"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M3 20h18M5 16l4-4 4 4 6-6"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">
            EN 50160
          </h3>
        </header>
      </div>
      <div class="px-4 py-4 space-y-3">
        <!-- No Data 표시 -->
        <div v-if="!tbdata || Object.keys(tbdata).length === 0" 
             class="flex flex-col items-center justify-center py-12">
          <svg class="w-16 h-16 text-gray-100 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          <p class="text-lg font-medium text-gray-500 dark:text-gray-100">No Summary data available</p>
        </div>
        
        <!-- 테이블 -->
        <table v-else class="table-auto w-full dark:text-white">
          <thead class="text-xs uppercase text-gray-100 bg-gray-50 dark:bg-gray-400 dark:text-gray-200 dark:bg-opacity-50 rounded-sm">
            <tr>
              <th class="p-2"><div class="font-bold text-left">{{ t('report.cardContext.pqInfo.parameter') }}</div></th>
              <th class="p-2"><div class="font-bold text-center">L1</div></th>
              <th class="p-2"><div class="font-bold text-center">L2</div></th>
              <th class="p-2"><div class="font-bold text-center">L3</div></th>
              <th class="p-2"><div class="font-bold text-center">{{ t('report.cardContext.pqInfo.compliance') }}</div></th>
              <th class="p-2"><div class="font-bold text-center">{{ t('report.cardContext.pqInfo.RequiredQuality') }}</div></th>
            </tr>
          </thead>
          <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
            <tr v-for="(row, index) in data" :key="index">
              <td class="p-2 text-left">{{ row.parameter }}</td>
              <td class="p-2 text-center">{{ tbdata[makeKey(row.parameter, 'L1')] ?? '-' }}</td>
              <td class="p-2 text-center">{{ tbdata[makeKey(row.parameter, 'L2')] ?? '-' }}</td>
              <td class="p-2 text-center">{{ tbdata[makeKey(row.parameter, 'L3')] ?? '-' }}</td>
              <td class="p-2 text-center">{{ getComp(row.parameter) ?? '-' }}</td>
              <td class="p-2 text-center">{{ row.required }}</td>
            </tr>
          </tbody>
        </table>
        <div v-if="!enReport_data || Object.keys(enReport_data).length === 0" 
             class="flex flex-col items-center justify-center py-12">
          <svg class="w-16 h-16 text-gray-100 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          <p class="text-lg font-medium text-gray-500 dark:text-gray-100">No Detail data available</p>
        </div>
        <!-- Frequency 차트 -->
         <div v-else>
          <EN50160Chart title="Frequency" :timeseries="frequencyTimeseries" :yAxis="frequencyYAxis"
          :limitLines="frequencyLimits" :markAreas="frequencyAreas" :tableColumns="frequencyTableColumns"
          :tableData="frequencyTableData" :legendItems="frequencyLegend" chartHeight="300px"
          :showHistogram="true" :histogramData="frequencyHistogram" histogramUnit="Hz"
          histogramColor="#008080" :histogramLimitMin="59.4" :histogramLimitMax="60.6"/>

        <!-- 전압 차트 -->
        <EN50160Chart title="Supply Voltage Variations" :timeseries="voltageTimeseries" :yAxis="voltageYAxis"
          :limitLines="voltageLimits" :markAreas="voltageAreas" :tableColumns="voltageTableColumns"
          :tableData="voltageTableData" :legendItems="voltageLegend" chartHeight="320px"
          :showPhaseHistograms="true" :phaseHistograms="voltagePhaseHistograms" histogramUnit="V"
          :histogramLimitMin="voltageLimitMin" :histogramLimitMax="voltageLimitMax"/>

        <!-- THD 차트 -->
        <EN50160Chart title="Voltage THD" :timeseries="thdTimeseries" :yAxis="thdYAxis"
          :limitLines="thdLimits" :tableColumns="thdTableColumns" :tableData="thdTableData"
          :legendItems="thdLegend" chartHeight="280px" :showPhaseHistograms="true"
          :phaseHistograms="thdPhaseHistograms" histogramUnit="%" :histogramLimitValue="8"/>

        <!-- 불평형 차트 -->
        <EN50160Chart title="Voltage Unbalance" :timeseries="unbalanceTimeseries" :yAxis="unbalanceYAxis"
          :limitLines="unbalanceLimits" :tableColumns="unbalanceTableColumns" :tableData="unbalanceTableData"
          chartHeight="250px" :showHistogram="true" :histogramData="unbalanceHistogram"
          histogramUnit="%" histogramColor="#9c27b0" :histogramLimitValue="2"/>
        
        <!-- Harmonics 테이블 -->
        <EN50160Harmonics :harmonicsData="harmonicsData"/>
         </div>
      </div>
    </div>

    <!-- ITIC Section -->
    <div class="relative col-span-full xl:col-span-12 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 shadow-sm rounded-b-lg">
      <div class="absolute top-0 left-0 right-0 h-0.5 bg-green-500" aria-hidden="true"></div>
      <div class="px-5 pt-5 pb-6 border-b border-gray-200 dark:border-gray-700/60">
        <header class="flex items-center mb-2">
          <div class="w-6 h-6 rounded-full shrink-0 bg-green-500 mr-3">
            <svg class="w-6 h-6 fill-current text-white" viewBox="0 0 24 24">
              <path d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="3" cy="12" r="1.5" fill="currentColor"/><circle cx="21" cy="12" r="1.5" fill="currentColor"/>
              <path d="M5 12v7m14-7v7" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
            </svg>
          </div>
          <h3 class="text-lg text-gray-800 dark:text-gray-100 font-semibold">ITIC</h3>
        </header>
      </div>
      <div class="px-4 py-4 space-y-3">
        <div class="flex items-center gap-x-6 mb-4">
          <div class="flex items-center gap-x-2">
            <span class="w-3 h-3 rounded-full" style="background-color: orange"></span>
            <span class="text-sm text-gray-600 dark:text-gray-100">sag/interruption</span>
          </div>
          <div class="flex items-center gap-x-2">
            <span class="w-3 h-3 rounded-full" style="background-color: green"></span>
            <span class="text-sm text-gray-600 dark:text-gray-100">swell</span>
          </div>
        </div>
        <LineChart :data="linechartData" width="595" height="248" />
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue'
import LineChart from '../../../charts/connect/LineChart_ITIC.vue'
import Report_PQ from './Report_PQ.vue'
import Report_PQ_detail from './Report_PQ_detail.vue'
import EN50160Chart from './EN_Frequency.vue'
import EN50160Harmonics from './EN50160Harmonics.vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { useReportData } from '@/composables/reportDict'

export default {
  name: 'ReportComponent',
  components: { LineChart, Report_PQ, Report_PQ_detail, EN50160Chart, EN50160Harmonics },
  props: {
    data: { type: Array, default: () => [] },
    channel: { type: String, default: '' },
    mode: { type: Boolean, default: false },
    reportData: { type: Object, default: null }
  },
  setup(props) {
    const { t } = useI18n()
    const channel = ref(props.channel)
    const mode = computed(() => props.mode)
    const iticDataList = ref([])

    // props에서 받은 데이터 사용
    const enReport_data = computed(() => props.reportData)

    // === Frequency ===
    const frequencyTimeseries = computed(() => {
      const freq = enReport_data.value?.frequency
      if (!freq?.timeseries) return { labels: [], datasets: [] }
      return { labels: freq.timeseries.labels, datasets: [{ name: 'Frequency', data: freq.timeseries.data, color: '#008080' }] }
    })
    const frequencyYAxis = { name: '[Hz]', min: 56, max: 63 }
    const frequencyLimits = computed(() => {
      const limits = enReport_data.value?.frequency?.limits || {}
      return [
        { value: limits.limit_100?.max || 62.4, color: '#f44336' },
        { value: limits.limit_100?.min || 56.4, color: '#f44336' },
        { value: limits.limit_99_5?.max || 60.6, color: '#ffc107' },
        { value: limits.limit_99_5?.min || 59.4, color: '#ffc107' }
      ]
    })
    const frequencyAreas = computed(() => {
      const limits = enReport_data.value?.frequency?.limits || {}
      return [{ from: limits.limit_99_5?.min || 59.4, to: limits.limit_99_5?.max || 60.6, color: 'rgba(255, 255, 200, 0.3)' }]
    })
    const frequencyTableColumns = [
      { key: 'requirement', label: 'EN50160 Requirement' },
      { key: 'measured', label: 'Measured Frequency' },
      { key: 'result', label: 'Result' }
    ]
    const frequencyTableData = computed(() => {
      const freq = enReport_data.value?.frequency
      if (!freq) return []
      const stats = freq.statistics || {}, limits = freq.limits || {}
      return [
        { requirement: `99.5% of the week: ${limits.limit_99_5?.min}Hz - ${limits.limit_99_5?.max}Hz`, measured: `${stats.min}Hz ~ ${stats.max}Hz`, result: stats.result_99_5 },
        { requirement: `100% of the week: ${limits.limit_100?.min}Hz - ${limits.limit_100?.max}Hz`, measured: `${stats.min}Hz ~ ${stats.max}Hz`, result: stats.result_100 }
      ]
    })
    const frequencyLegend = [{ name: 'Flagged Data', color: 'rgba(255, 255, 200, 0.8)' }, { name: 'Frequency', color: '#008080' }]
    const frequencyHistogram = computed(() => enReport_data.value?.frequency?.histogram || null)

    // === Voltage ===
    const voltageTimeseries = computed(() => {
      const volt = enReport_data.value?.voltage
      if (!volt?.phases) return { labels: [], datasets: [] }
      const phases = volt.phases, firstPhase = Object.values(phases)[0]
      return {
        labels: firstPhase?.timeseries?.labels || [],
        datasets: [
          { name: 'L1', data: phases.L1?.timeseries?.data || [], color: '#e53935', areaStyle: false },
          { name: 'L2', data: phases.L2?.timeseries?.data || [], color: '#43a047', areaStyle: false },
          { name: 'L3', data: phases.L3?.timeseries?.data || [], color: '#1e88e5', areaStyle: false }
        ]
      }
    })
    const voltageYAxis = computed(() => {
      const nominal = enReport_data.value?.voltage?.limits?.nominal || 22900
      return { name: '[V]', min: nominal * 0.8, max: nominal * 1.2, formatter: (v) => (v / 1000).toFixed(1) + 'k' }
    })
    const voltageLimits = computed(() => {
      const limits = enReport_data.value?.voltage?.limits || {}
      return [
        { value: limits.limit_100?.max, color: '#f44336' }, { value: limits.limit_100?.min, color: '#f44336' },
        { value: limits.limit_95?.max, color: '#ffc107' }, { value: limits.limit_95?.min, color: '#ffc107' }
      ]
    })
    const voltageAreas = computed(() => {
      const limits = enReport_data.value?.voltage?.limits || {}
      return [{ from: limits.limit_95?.min, to: limits.limit_95?.max, color: 'rgba(255, 255, 200, 0.3)' }]
    })
    const voltageTableColumns = [
      { key: 'requirement', label: 'EN50160 Requirement' }, { key: 'phase', label: 'Phase' },
      { key: 'measured', label: 'Measured Voltage' }, { key: 'inRange', label: 'In Range %' }, { key: 'result', label: 'Result' }
    ]
    const voltageTableData = computed(() => {
      const volt = enReport_data.value?.voltage
      if (!volt?.phases) return []
      const phases = volt.phases, limits = volt.limits || {}, rows = []
      Object.entries(phases).forEach(([key, phase], i) => {
        rows.push({ requirement: i === 0 ? `95% of the week: ${limits.limit_95?.min?.toFixed(0)}V - ${limits.limit_95?.max?.toFixed(0)}V (±10%)` : '',
          phase: key, measured: `${phase.statistics?.min?.toFixed(0)}V ~ ${phase.statistics?.max?.toFixed(0)}V`,
          inRange: `${phase.statistics?.in_range_95_percent}%`, result: phase.statistics?.result_95 })
      })
      Object.entries(phases).forEach(([key, phase], i) => {
        rows.push({ requirement: i === 0 ? `100% of the week: ${limits.limit_100?.min?.toFixed(0)}V - ${limits.limit_100?.max?.toFixed(0)}V (-15%/+10%)` : '',
          phase: key, measured: `${phase.statistics?.min?.toFixed(0)}V ~ ${phase.statistics?.max?.toFixed(0)}V`,
          inRange: `${phase.statistics?.in_range_100_percent}%`, result: phase.statistics?.result_100 })
      })
      return rows
    })
    const voltageLegend = [{ name: 'L1', color: '#e53935' }, { name: 'L2', color: '#43a047' }, { name: 'L3', color: '#1e88e5' }]
    const voltagePhaseHistograms = computed(() => {
      const volt = enReport_data.value?.voltage
      if (!volt?.phases) return null
      return { L1: volt.phases.L1?.histogram, L2: volt.phases.L2?.histogram, L3: volt.phases.L3?.histogram }
    })
    const voltageLimitMin = computed(() => enReport_data.value?.voltage?.limits?.limit_95?.min || null)
    const voltageLimitMax = computed(() => enReport_data.value?.voltage?.limits?.limit_95?.max || null)

    // === THD ===
    const thdTimeseries = computed(() => {
      const thd = enReport_data.value?.thd
      if (!thd?.phases) return { labels: [], datasets: [] }
      const phases = thd.phases, firstPhase = Object.values(phases)[0]
      return {
        labels: firstPhase?.timeseries?.labels || [],
        datasets: [
          { name: 'L1', data: phases.L1?.timeseries?.data || [], color: '#e53935', areaStyle: false },
          { name: 'L2', data: phases.L2?.timeseries?.data || [], color: '#43a047', areaStyle: false },
          { name: 'L3', data: phases.L3?.timeseries?.data || [], color: '#1e88e5', areaStyle: false }
        ]
      }
    })
    const thdYAxis = { name: '[%]', min: 0, max: 15 }
    const thdLimits = computed(() => [{ value: enReport_data.value?.thd?.limits?.limit_95 || 8, color: '#f44336' }])
    const thdTableColumns = [
      { key: 'phase', label: 'Phase' }, { key: 'max', label: 'Max' }, { key: 'avg', label: 'Avg' },
      { key: 'inRange', label: 'In Range %' }, { key: 'result', label: 'Result' }
    ]
    const thdTableData = computed(() => {
      const thd = enReport_data.value?.thd
      if (!thd?.phases) return []
      return Object.entries(thd.phases).map(([key, phase]) => ({
        phase: key, max: `${phase.statistics?.max?.toFixed(2)}%`, avg: `${phase.statistics?.avg?.toFixed(2)}%`,
        inRange: `${phase.statistics?.in_range_95_percent}%`, result: phase.statistics?.result
      }))
    })
    const thdLegend = [{ name: 'L1', color: '#e53935' }, { name: 'L2', color: '#43a047' }, { name: 'L3', color: '#1e88e5' }]
    const thdPhaseHistograms = computed(() => {
      const thd = enReport_data.value?.thd
      if (!thd?.phases) return null
      return { L1: thd.phases.L1?.histogram, L2: thd.phases.L2?.histogram, L3: thd.phases.L3?.histogram }
    })

    // === Harmonics ===
    const harmonicsData = computed(() => enReport_data.value?.harmonics || null)

    // === Unbalance ===
    const unbalanceTimeseries = computed(() => {
      const unbal = enReport_data.value?.unbalance
      if (!unbal?.timeseries) return { labels: [], datasets: [] }
      return { labels: unbal.timeseries.labels, datasets: [{ name: 'Unbalance', data: unbal.timeseries.data, color: '#9c27b0' }] }
    })
    const unbalanceYAxis = { name: '[%]', min: 0, max: 5 }
    const unbalanceLimits = computed(() => [{ value: enReport_data.value?.unbalance?.limits?.limit_95 || 2, color: '#f44336' }])
    const unbalanceTableColumns = [
      { key: 'label', label: 'Parameter' }, { key: 'max', label: 'Max' }, { key: 'avg', label: 'Avg' },
      { key: 'inRange', label: 'In Range %' }, { key: 'result', label: 'Result' }
    ]
    const unbalanceTableData = computed(() => {
      const unbal = enReport_data.value?.unbalance
      if (!unbal?.statistics) return []
      const stats = unbal.statistics
      return [{ label: 'Voltage Unbalance', max: `${stats.max?.toFixed(3)}%`, avg: `${stats.avg?.toFixed(3)}%`,
        inRange: `${stats.in_range_95_percent}%`, result: stats.result }]
    })
    const unbalanceHistogram = computed(() => enReport_data.value?.unbalance?.histogram || null)

    // === 기타 ===
    const tbdata = ref(null)
    const COMPLIANCE_BITS = {
      "Frequency Variation 1": [0], "Frequency Variation 2": [1], "Voltage Variation 1": [2, 3, 4],
      "Voltage Variation 2": [5, 6, 7], "Voltage Unbalance": [8], "THD": [9, 10, 11],
      "Harmonics": [12, 13, 14], "Pst": [15, 16, 17], "Plt": [18, 19, 20], "Signal Vol.": [21, 22, 23]
    }
    const { baseChart, parseMask, getfinValue, makeKey } = useReportData()

    onMounted(() => { fetchData(); fetchITICData() })

    const fetchData = async () => {
      try {
        const response = await axios.get(`/api/getEn50160/${props.channel}`)
        tbdata.value = response.data.success ? response.data.data : {}
      } catch (error) { tbdata.value = {} }
    }

const fetchITICData = async () => {
  try {
    const response = await axios.get(`/api/getITIC/${props.channel}`)
    if (response.data.success) {
      const data = response.data.data, ratedV = parseInt(response.data.ratedV)
      iticDataList.value = []
      for (let i = 0; i < data.length; i++) {
        let phaselist = parseMask(data[i]["mask"]), levellist = [data[i]["level_l1"], data[i]["level_l2"], data[i]["level_l3"]]
        let Yvalue = (getfinValue(phaselist, levellist, data[i].event_type) * 100) / ratedV
        let Xvalue = data[i]["duration"] / 1000.0
        
        // ✅ 숫자로 변환, 로그스케일에서 0 방지
        let yNum = parseFloat(Yvalue.toFixed(2))
        if (yNum <= 0) yNum = 0.5  // 로그스케일에서 0 대신 0.5%
        if (Xvalue <= 0) Xvalue = 0.0001  // X도 0 방지
        
        iticDataList.value.push({
          label: data[i].event_type == 'SWELL' ? 'Swell Event' : 'Sag Event',
          type: 'scatter',
          data: [{ x: Xvalue, y: yNum }],  // ✅ 숫자 타입 보장
          backgroundColor: data[i].event_type == 'SWELL' ? 'green' : 'orange',
          pointRadius: 5,
          pointHoverRadius: 8,
          showLine: false
        })
      }
    }
  } catch (error) { console.log("ITIC 데이터 가져오기 실패:", error) }
}

    const linechartData = computed(() => {
      const base = JSON.parse(JSON.stringify(baseChart))
      iticDataList.value.forEach(item => base.datasets.push(item))     
      return base
    })

    const getComp = (param) => {
      if (!tbdata.value || !tbdata.value["status&compliance"]) return "-"
      const comp = tbdata.value["status&compliance"], bits = COMPLIANCE_BITS[param]
      if (!bits) return "-"
      return bits.some(bit => (comp & (1 << bit)) !== 0) ? "Failed" : "OK"
    }

    watch(() => props.channel, (newChannel) => { channel.value = newChannel; fetchData(); fetchITICData() })

    return {
      channel, linechartData, mode, tbdata, makeKey, getComp, t, iticDataList, enReport_data,
      frequencyTimeseries, frequencyYAxis, frequencyLimits, frequencyAreas, frequencyTableColumns, frequencyTableData, frequencyLegend, frequencyHistogram,
      voltageTimeseries, voltageYAxis, voltageLimits, voltageAreas, voltageTableColumns, voltageTableData, voltageLegend, voltagePhaseHistograms, voltageLimitMin, voltageLimitMax,
      thdTimeseries, thdYAxis, thdLimits, thdTableColumns, thdTableData, thdLegend, thdPhaseHistograms,
      harmonicsData,
      unbalanceTimeseries, unbalanceYAxis, unbalanceLimits, unbalanceTableColumns, unbalanceTableData, unbalanceHistogram
    }
  }
}
</script>
