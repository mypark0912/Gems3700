<template>
  <div class="space-y-5">

    <!-- 헤더 -->
    <div class="flex items-center justify-between">
      <span class="inline-block px-3 py-1 text-sm font-semibold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-500/10 rounded">
        Main Channel
      </span>
      <button class="px-5 py-2 text-sm font-semibold rounded-lg border transition-all bg-violet-500 hover:bg-violet-600 text-white border-violet-500 shadow-sm">
        Update Settings
      </button>
    </div>

    <!-- 상단 3카드: CT / PT / Demand -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">

      <!-- CT -->
      <div class="setting-card">
        <div class="card-header border-green-400">
          <div class="w-6 h-6 rounded-full bg-green-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>
          </div>
          <h3 class="card-title">CT</h3>
        </div>
        <div class="card-body">
          <div class="field">
            <label class="field-label">CT Directions</label>
            <div class="space-y-2 mt-1">
              <div v-for="ct in ['CT1','CT2','CT3']" :key="ct" class="flex items-center gap-2">
                <span class="text-sm text-gray-500 dark:text-gray-400 w-8">{{ ct }}</span>
                <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600 flex-1">
                  <button @click="form.ctDir[ct]='Positive'" class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                    :class="form.ctDir[ct]==='Positive' ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">Positive</button>
                  <button @click="form.ctDir[ct]='Negative'" class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                    :class="form.ctDir[ct]==='Negative' ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">Negative</button>
                </div>
              </div>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 mt-1">
            <div class="field"><label class="field-label">Starting Current (mA)</label><input type="number" v-model="form.startingCurrent" class="field-input" /></div>
            <div class="field"><label class="field-label">Rated Current (A)</label><input type="number" v-model="form.ratedCurrent" class="field-input" /></div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Primary</label><input type="number" v-model="form.ctPrimary" class="field-input" /></div>
            <div class="field">
              <label class="field-label">Secondary</label>
              <select v-model="form.ctSecondary" class="field-select"><option>100mA</option><option>1A</option><option>5A</option></select>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">ZCT Scale</label><input type="number" v-model="form.zctScale" class="field-input" /></div>
            <div class="field">
              <label class="field-label">ZCT Type</label>
              <select v-model="form.zctType" class="field-select"><option>None</option><option>Type A</option><option>Type B</option></select>
            </div>
          </div>
        </div>
      </div>

      <!-- PT -->
      <div class="setting-card">
        <div class="card-header border-blue-400">
          <div class="w-6 h-6 rounded-full bg-blue-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>
          </div>
          <h3 class="card-title">PT</h3>
        </div>
        <div class="card-body">
          <div class="field"><label class="field-label">Line Frequency (Hz)</label><input type="number" v-model="form.lineFreq" class="field-input" /></div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field">
              <label class="field-label">Wiring Mode</label>
              <select v-model="form.wiringMode" class="field-select"><option>3P3W</option><option>3P4W</option><option>1P2W</option></select>
            </div>
            <div class="field"><label class="field-label">Rated Voltage (V)</label><input type="number" v-model="form.ratedVoltage" class="field-input" /></div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Primary</label><input type="number" v-model="form.ptPrimary" class="field-input" /></div>
            <div class="field"><label class="field-label">Secondary</label><input type="number" v-model="form.ptSecondary" class="field-input" /></div>
          </div>
          <div class="field">
            <label class="field-label">대시보드 출력설정</label>
            <select v-model="form.dashOutput" class="field-select"><option>Line Voltage</option><option>Phase Voltage</option></select>
          </div>
        </div>
      </div>

      <!-- Demand -->
      <div class="setting-card">
        <div class="card-header border-indigo-400">
          <div class="w-6 h-6 rounded-full bg-indigo-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/></svg>
          </div>
          <h3 class="card-title">Demand</h3>
        </div>
        <div class="card-body">
          <div class="field"><label class="field-label">Target (W)</label><input type="number" v-model="form.demandTarget" class="field-input" /></div>
          <div class="field">
            <label class="field-label">측정 간격</label>
            <select v-model="form.demandInterval" class="field-select"><option>15 분</option><option>30 분</option><option>60 분</option></select>
          </div>
          <div class="field">
            <label class="field-label">트렌드 수집 여부</label>
            <select v-model="form.demandTrend" class="field-select"><option>No</option><option>Yes</option></select>
          </div>
        </div>
      </div>
    </div>

    <!-- 알람 섹션 -->
    <div class="setting-card">
      <div class="card-header border-red-400">
        <div class="w-6 h-6 rounded-full bg-red-400 flex items-center justify-center">
          <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6 6 0 10-12 0v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>
        </div>
        <h3 class="card-title">알람</h3>
      </div>
      <div class="card-body">
        <div class="mb-4 max-w-xs">
          <label class="field-label">비교 시간 지연 (초)</label>
          <input type="number" v-model="form.alarmDelay" class="field-input" />
        </div>
        <div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
          <div v-for="half in [0,1]" :key="half" class="overflow-x-auto">
            <table class="w-full text-xs table-fixed">
              <colgroup>
                <col style="width: 28px" />
                <col style="width: 35%" />
                <col style="width: 20%" />
                <col style="width: 20%" />
                <col style="width: 20%" />
              </colgroup>
              <thead>
                <tr class="bg-gray-50 dark:bg-gray-700/50">
                  <th class="alarm-th text-center">알람</th>
                  <th class="alarm-th">발생 요소</th>
                  <th class="alarm-th">발생 조건</th>
                  <th class="alarm-th">편차 범위(%)</th>
                  <th class="alarm-th">기준 값</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="i in 16" :key="i" class="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-700/30">
                  <td class="alarm-td text-center text-gray-400">{{ i + half * 16 }}</td>
                  <td class="alarm-td">
                    <select v-model="form.alarms[i-1+half*16].element" class="alarm-select">
                      <option>None</option>
                      <option>Phase Voltage L1</option>
                      <option>Phase Voltage L2</option>
                      <option>Phase Voltage L3</option>
                      <option>Line Voltage L12</option>
                      <option>Line Voltage L23</option>
                      <option>Line Voltage L31</option>
                      <option>Current L1</option>
                      <option>Current L2</option>
                      <option>Current L3</option>
                      <option>Frequency</option>
                      <option>Power Factor</option>
                      <option>Active Power</option>
                      <option>THD Voltage</option>
                      <option>THD Current</option>
                    </select>
                  </td>
                  <td class="alarm-td">
                    <select v-model="form.alarms[i-1+half*16].condition" class="alarm-select">
                      <option>Over</option>
                      <option>Under</option>
                    </select>
                  </td>
                  <td class="alarm-td">
                    <input type="number" v-model="form.alarms[i-1+half*16].deviation" class="alarm-input" />
                  </td>
                  <td class="alarm-td">
                    <input type="number" v-model="form.alarms[i-1+half*16].reference" class="alarm-input" />
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- 이벤트 1/2/3 + 트렌드 -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-5">

      <!-- 이벤트 설정 1 -->
      <div class="setting-card">
        <div class="card-header border-orange-400">
          <div class="w-6 h-6 rounded-full bg-orange-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="card-title">이벤트 설정 1</h3>
        </div>
        <div class="card-body">
          <h4 class="section-label">TRANSIENT VOLTAGE</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event1.tvAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event1.tvAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event1.tvLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event1.tvHoldoff" class="field-input" /></div>
          </div>
          <div class="field"><label class="field-label">Fast Change (%)</label><input type="number" v-model="form.event1.tvFastChange" class="field-input" /></div>

          <h4 class="section-label mt-4">TRANSIENT CURRENT</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event1.tcAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event1.tcAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event1.tcLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event1.tcHoldoff" class="field-input" /></div>
          </div>
          <div class="field"><label class="field-label">Fast Change (%)</label><input type="number" v-model="form.event1.tcFastChange" class="field-input" /></div>
        </div>
      </div>

      <!-- 이벤트 설정 2 -->
      <div class="setting-card">
        <div class="card-header border-yellow-400">
          <div class="w-6 h-6 rounded-full bg-yellow-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="card-title">이벤트 설정 2</h3>
        </div>
        <div class="card-body">
          <h4 class="section-label">SAG</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event2.sagAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event2.sagAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event2.sagLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event2.sagHoldoff" class="field-input" /></div>
          </div>

          <h4 class="section-label mt-4">SWELL</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event2.swellAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event2.swellAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event2.swellLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event2.swellHoldoff" class="field-input" /></div>
          </div>
        </div>
      </div>

      <!-- 이벤트 설정 3 -->
      <div class="setting-card">
        <div class="card-header border-pink-400">
          <div class="w-6 h-6 rounded-full bg-pink-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          </div>
          <h3 class="card-title">이벤트 설정 3</h3>
        </div>
        <div class="card-body">
          <h4 class="section-label">OVER CURRENT</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event3.ocAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event3.ocAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event3.ocLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event3.ocHoldoff" class="field-input" /></div>
          </div>

          <h4 class="section-label mt-4">INTERRUPTION</h4>
          <div class="field">
            <label class="field-label">Action</label>
            <div class="flex rounded-lg overflow-hidden border border-gray-200 dark:border-gray-600">
              <button v-for="a in ['None','Event','Capture']" :key="a"
                @click="form.event3.intAction = a"
                class="flex-1 py-1.5 text-xs font-semibold transition-colors"
                :class="form.event3.intAction===a ? 'bg-violet-500 text-white' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400'">
                {{ a }}
              </button>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="field"><label class="field-label">Level (%)</label><input type="number" v-model="form.event3.intLevel" class="field-input" /></div>
            <div class="field"><label class="field-label">Holdoff time (ms)</label><input type="number" v-model="form.event3.intHoldoff" class="field-input" /></div>
          </div>
        </div>
      </div>

      <!-- 트렌드 설정 -->
      <div class="setting-card">
        <div class="card-header border-teal-400">
          <div class="w-6 h-6 rounded-full bg-teal-400 flex items-center justify-center">
            <svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"/></svg>
          </div>
          <h3 class="card-title">트렌드 설정</h3>
        </div>
        <div class="card-body">
          <h4 class="section-label">SAVE OPTION</h4>
          <div class="field">
            <label class="field-label">Period (min)</label>
            <input type="number" v-model="form.trendPeriod" class="field-input" />
          </div>
          <table class="w-full text-xs table-fixed mt-2">
            <colgroup>
              <col style="width: 36px" />
              <col />
            </colgroup>
            <thead>
              <tr class="bg-gray-50 dark:bg-gray-700/50">
                <th class="alarm-th text-center">Index</th>
                <th class="alarm-th">Parameter</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, i) in form.trendParams" :key="i" class="border-b border-gray-100 dark:border-gray-700/50 hover:bg-gray-50 dark:hover:bg-gray-700/30">
                <td class="alarm-td text-center text-gray-400">{{ i + 1 }}</td>
                <td class="alarm-td">
                  <select v-model="form.trendParams[i]" class="alarm-select">
                    <option>TDD</option>
                    <option>Frequency</option>
                    <option>Line Voltage</option>
                    <option>Phase Voltage</option>
                    <option>Current</option>
                    <option>Unbalance</option>
                    <option>PF</option>
                    <option>THD</option>
                    <option>Active Power</option>
                    <option>Reactive Power</option>
                    <option>Apparent Power</option>
                  </select>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const form = reactive({
  ctDir: { CT1: 'Positive', CT2: 'Positive', CT3: 'Positive' },
  startingCurrent: 100,
  ratedCurrent: 1.2,
  ctPrimary: 100,
  ctSecondary: '100mA',
  zctScale: 1000,
  zctType: 'None',
  lineFreq: 60,
  wiringMode: '3P3W',
  ratedVoltage: 72,
  ptPrimary: 72,
  ptSecondary: 72,
  dashOutput: 'Line Voltage',
  demandTarget: 1,
  demandInterval: '15 분',
  demandTrend: 'No',
  alarmDelay: 1,
  alarms: Array.from({ length: 32 }, (_, i) => ({
    element: i === 0 ? 'Phase Voltage L1' : i === 1 ? 'Frequency' : i === 2 ? 'Frequency' : 'None',
    condition: i === 0 ? 'Over' : 'Under',
    deviation: 1,
    reference: i === 0 ? 230 : i === 1 ? 64 : i === 2 ? 58 : 0,
  })),
  event1: { tvAction: 'None', tvLevel: 200, tvHoldoff: 200, tvFastChange: 4, tcAction: 'None', tcLevel: 120, tcHoldoff: 200, tcFastChange: 4 },
  event2: { sagAction: 'None', sagLevel: 80, sagHoldoff: 500, swellAction: 'None', swellLevel: 120, swellHoldoff: 500 },
  event3: { ocAction: 'None', ocLevel: 120, ocHoldoff: 500, intAction: 'None', intLevel: 5, intHoldoff: 100 },
  trendPeriod: 5,
  trendParams: ['TDD', 'Frequency', 'Line Voltage', 'Phase Voltage', 'Current', 'Unbalance', 'PF', 'THD'],
})
</script>

<style scoped>
.setting-card {
  @apply bg-white dark:bg-gray-800 shadow-sm rounded-xl border border-gray-200 dark:border-gray-700/60 flex flex-col overflow-hidden;
}

.card-header {
  @apply flex items-center gap-2.5 px-5 py-4 border-t-[3px];
}

.card-title {
  @apply text-base font-bold text-gray-800 dark:text-gray-100;
}

.card-body {
  @apply px-5 pb-5 pt-2 flex-1;
}

.section-label {
  @apply text-[11px] font-bold text-gray-800 dark:text-gray-200 uppercase tracking-wider mb-3;
}

.field {
  @apply mb-3;
}

.field-label {
  @apply block text-sm text-gray-500 dark:text-gray-400 mb-1;
}

.field-input {
  @apply w-full h-9 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-all;
}

.field-select {
  @apply w-full h-9 px-3 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-all;
}

/* 알람/트렌드 테이블 — field-input/select와 동일한 스타일, 높이만 h-8로 축소 */
.alarm-th {
  @apply px-2 py-2 text-left text-[11px] font-semibold text-gray-500 dark:text-gray-400;
  @apply border-b border-gray-200 dark:border-gray-700;
  white-space: normal;
  word-break: keep-all;
  line-height: 1.3;
}

.alarm-td {
  @apply px-1 py-1;
  min-width: 0;
}

.alarm-select {
  @apply w-full h-9 px-2 text-sm rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-all;
  min-width: 0;
  max-width: 100%;
}

.alarm-input {
  @apply w-full h-9 px-2 text-xs rounded-lg border border-gray-200 dark:border-gray-600;
  @apply bg-gray-50 dark:bg-gray-700 text-gray-800 dark:text-gray-200;
  @apply focus:outline-none focus:ring-1 focus:ring-violet-400 focus:border-violet-400 transition-all;
  min-width: 0;
}
</style>