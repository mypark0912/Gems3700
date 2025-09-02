<template>
  <div ref="pdfContent" class="grow flex flex-col md:translate-x-0 transition-transform duration-300 ease-in-out max-h-[70vh] overflow-y-auto p-4 bg-gradient-to-br from-gray-50 to-gray-100" :class="profileSidebarOpen ? 'translate-x-1/3' : 'translate-x-0'">

    <!-- Content -->
    <div class="relative px-4 sm:px-6 pb-8">
      
      <!-- Channel Information Section -->
      <section ref="channelInfo" class="mt-6">
        <div v-if="pdf_hidden" class="inline-flex items-start mb-6 py-3">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
              <svg class="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
              </svg>
            </div>
            <h1 class="text-2xl text-gray-800 dark:text-gray-100 font-bold">{{ channelComputed }} Channel Report</h1>
          </div>
        </div>

        <div v-if="channelComputed === 'Main'? setupMenu.MainDiagnosis : setupMenu.SubDiagnosis" class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header with gradient -->
          <div class="bg-gradient-to-r from-slate-500 to-slate-700 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <img :src="equipImageSrc" alt="장비 이미지"
                    class="w-8 h-8 object-cover rounded-lg shadow-md border border-gray-300 dark:border-gray-600" />
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.Channel`) }}</h3>
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="p-6">
            <div class="space-y-6">
            <!-- 첫 번째 줄: 이름/타입 (col-12 꽉 차게) -->
            <div class="w-full">
                <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-600">
                <div class="grid grid-cols-3 gap-6">
                    <div class="flex flex-col">
                    <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{{ t(`diagnosis.info.name`) }}</span>
                    <span class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ channelComputed.toLowerCase() == 'main'? asset.assetNickname_main : asset.assetNickname_sub }}</span>
                    </div>
                    <div class="flex flex-col">
                    <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{{ t(`diagnosis.info.type`) }}</span>
                    <span class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ channelComputed.toLowerCase() == 'main'? asset.assetType_main : asset.assetType_sub }}</span>
                    </div>
                    <div v-if="devLocation != ''" class="flex flex-col">
                    <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{{ t(`diagnosis.info.location`) }}</span>
                    <span class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ devLocation }}</span>
                    </div>
                </div>
                </div>
            </div>

            <!-- 두 번째 줄: 나머지 항목들 (col-12 박스 안에 수평 정렬) -->
            <div class="w-full">
                <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-600">
                <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                    <div v-for="(item, index) in datalist" :key="index" class="group hover:scale-105 transition-transform duration-200">
                    <div class="flex flex-col h-full">
                        <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">{{ t(`dashboard.transDiag.${item.Name}`) }}</span>
                        <span class="text-lg font-bold text-gray-800 dark:text-gray-100 group-hover:text-blue-600 transition-colors">{{ item.Value }} {{ item.Unit }}</span>
                    </div>
                    </div>
                </div>
                </div>
            </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Diagnosis Section -->
      <section v-if="mode" ref="DiagnosisInfo" class="mt-8">
        <div class="space-y-6">
          <!-- Status Information -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
            <div class="bg-gradient-to-r from-sky-500 to-blue-600 p-3 text-white">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14z"/>
                    <path d="M3 10h18"/>
                    <path d="M10 3v18"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-xl font-bold">{{ t(`diagnosis.info.status`) }}</h3>
                </div>
              </div>
            </div>
            <div class="p-6">
               <Diagnosis_Barchart v-if="ischartDataValid" :channel="channelComputed" :data="diagData.chartdata" :mode="'PowerQuality'" class="h-auto" />
            </div>
            <div class="flex flex-col gap-4 p-2 pl-4 pr-4 mb-1">
                <template v-for="item in diagData.items">
                <StatusReport :data="item" />
                </template>
            </div>
          </div>
        </div>
      </section>

      <section v-if="mode && ischartDataValid && diagData.items.length > 0" ref="DiagnosisDetail" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
            <div class="bg-gradient-to-r from-pink-500 to-rose-600 p-3 text-white">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 17a.833.833 0 01-.833-.833 3.333 3.333 0 00-3.334-3.334.833.833 0 110-1.666 3.333 3.333 0 003.334-3.334.833.833 0 111.666 0 3.333 3.333 0 003.334 3.334.833.833 0 110 1.666 3.333 3.333 0 00-3.334 3.334c0 .46-.373.833-.833.833z"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-xl font-bold">진단 세부항목 트렌드</h3>
                </div>
              </div>
            </div>
            <div class="p-6">
                <div class="grid grid-cols-12 gap-4 p-2 pl-4 pr-4">
                    <div v-for="option in diagData.chartOption" class="col-span-12">
                    <ReportTrend :data="option" />
                    </div>
                </div>
            </div>
          </div>
      </section>

      <section v-if="mode" ref="PQStatusisInfo" class="mt-8">
        <div class="space-y-6">
          <!-- Status Information -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
            <div class="bg-gradient-to-r from-sky-500 to-blue-600 p-3 text-white">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 5a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v14a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-14z"/>
                    <path d="M3 10h18"/>
                    <path d="M10 3v18"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-xl font-bold">{{ t(`dashboard.diagnosis.pq`)}} <!--전력품질 상태--></h3>
                </div>
              </div>
            </div>
            <div class="p-6">
               <Diagnosis_Barchart v-if="ispqValid" :channel="channelComputed" :data="pqdata.chartdata" :mode="'PowerQuality'" class="h-auto" />
            </div>
            <div class="flex flex-col gap-4 p-2 pl-4 pr-4 mb-1">
                <template v-for="item in pqdata.items">
                <StatusReport :data="item" />
                </template>
            </div>
          </div>
        </div>
      </section>

      <section v-if="mode && ispqValid && pqdata.items.length > 0" ref="PQDetail" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
            <div class="bg-gradient-to-r from-pink-500 to-rose-600 p-3 text-white">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 17a.833.833 0 01-.833-.833 3.333 3.333 0 00-3.334-3.334.833.833 0 110-1.666 3.333 3.333 0 003.334-3.334.833.833 0 111.666 0 3.333 3.333 0 003.334 3.334.833.833 0 110 1.666 3.333 3.333 0 00-3.334 3.334c0 .46-.373.833-.833.833z"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-xl font-bold">전력품질 세부항목 트렌드</h3>
                </div>
              </div>
            </div>
            <div class="p-6">
                <div class="grid grid-cols-12 gap-4 p-2 pl-4 pr-4">
                    <div v-for="option in pqdata.chartOption" class="col-span-12">
                    <ReportTrend :data="option" />
                    </div>
                </div>
            </div>
          </div>
      </section>

      <section ref="PQTable" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-blue-500 to-blue-700 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M3 20h18M5 16l4-4 4 4 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">
                 {{ t(`report.cardTitle.PowerQuality`) }} <!--Power Quality - EN 50160 Report-->

                </h3>
              </div>
            </div>
          </div>

          <!-- Table Content -->
          <div class="p-6">
            <div class="overflow-x-auto">
                <!-- No Data 표시 또는 테이블 표시 -->
                <div v-if="!tbdata || Object.keys(tbdata).length === 0" 
                     class="flex flex-col items-center justify-center py-12">
                  <svg 
                    class="w-16 h-16 text-gray-400 dark:text-gray-600 mb-4" 
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24" 
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path 
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="1.5" 
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                  <p class="text-lg font-medium text-gray-500 dark:text-gray-400">No Data Available</p>
                  <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">EN50160 data will appear here when available</p>
                </div>
              <table v-else class="w-full">
                <thead>
                  <tr class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-700 dark:to-gray-800">
                    <th class="px-6 py-4 text-left text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider rounded-tl-lg">Parameter</th>
                    <th class="px-6 py-4 text-center text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">L1</th>
                    <th class="px-6 py-4 text-center text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">L2</th>
                    <th class="px-6 py-4 text-center text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">L3</th>
                    <th class="px-6 py-4 text-center text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider">Compliance</th>
                    <th class="px-6 py-4 text-center text-sm font-bold text-gray-700 dark:text-gray-300 uppercase tracking-wider rounded-tr-lg">Required Quality</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                  <tr v-for="(row, index) in pqTableData" :key="index" class="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                    <td class="px-6 py-4 text-sm font-medium text-gray-900 dark:text-gray-100">{{ row.parameter }}</td>
                    <td class="px-6 py-4 text-sm text-center text-gray-700 dark:text-gray-300">
                      {{ tbdata[makeKey(row.parameter, 'L1')] ?? '-' }}
                    </td>
                    <td class="px-6 py-4 text-sm text-center text-gray-700 dark:text-gray-300">
                      {{ tbdata[makeKey(row.parameter, 'L2')] ?? '-' }}
                    </td>
                    <td class="px-6 py-4 text-sm text-center text-gray-700 dark:text-gray-300">
                      {{ tbdata[makeKey(row.parameter, 'L3')] ?? '-' }}
                    </td>
                    <td class="px-6 py-4 text-center">
                      <span :class="getComplianceClass(getComp(row.parameter))" class="px-3 py-1 text-xs font-semibold">
                        {{ getComp(row.parameter) ?? '-' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 text-sm text-center text-gray-700 dark:text-gray-300">{{ row.required }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>

      <!-- Power Quality Chart Section -->
      <section ref="PQChart" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <path d="M3 12c1.5-2 3-2 4-2s2 .5 3 2 3 2 4 2 2-.5 3-2 3-2 4-2 2 .5 3 2" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <circle cx="3" cy="12" r="1.5" fill="currentColor"/>
                  <circle cx="21" cy="12" r="1.5" fill="currentColor"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">ITIC Curve Analysis</h3>
              </div>
            </div>
          </div>

          <!-- Chart Content -->
          <div class="p-6">
            <div class="bg-gradient-to-br from-gray-50 to-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <LineChart :data="linechartData" width="595" height="248" />
            </div>
          </div>
        </div>
      </section>

      <!-- Energy Section - Hourly -->
      <section ref="EnergyHourly" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-amber-500 to-orange-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2L3 14h7l-1 8 8-12h-6l2-8z"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.hourlyEngUsage`) }}<!--Hourly Energy Usage--></h3>
              </div>
            </div>
          </div>

          <!-- Chart Content -->
          <div class="p-6">
            <div class="bg-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <BarChart2 :data="hourlyChartData" width="595" height="248" />
            </div>
          </div>
        </div>
      </section>

      <!-- Energy Section - Daily -->
      <section ref="EnergyDaily" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-amber-500 to-orange-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2L3 14h7l-1 8 8-12h-6l2-8z"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.dailyEngUsage`) }}<!--Daily Energy Usage--></h3>
              </div>
            </div>
          </div>

          <!-- Chart Content -->
          <div class="p-6">
            <div class="bg-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <BarChart2 :data="dailyChartData" width="595" height="248" />
            </div>
          </div>
        </div>
      </section>

      <!-- Energy Section - Monthly -->
      <section ref="EnergyMonthly" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-amber-500 to-orange-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2L3 14h7l-1 8 8-12h-6l2-8z"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.monthlyEngUsage`) }}<!--Monthly Energy Usage--></h3>
              </div>
            </div>
          </div>

          <!-- Chart Content -->
          <div class="p-6">
            <div class="bg-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <BarChart2 :data="monthlyChartData" width="595" height="248" />
            </div>
          </div>
        </div>
      </section>

      <!-- Power Load Trend Section - 전력량 & 부하율 추이 -->
      <section v-if="assettypes == 'Transformer'" ref="PowerLoadTrend" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-cyan-500 to-blue-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M13 2L3 14h7l-1 8 8-12h-6l2-8z"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.Energy2`) }}<!--전력량 & 부하율 추이--></h3>
              </div>
            </div>
          </div>

          <!-- Chart Content -->
          <div class="p-6">
            <div class="bg-gradient-to-br from-gray-50 to-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <div ref="dualAxisChart" class="dual-axis-chart"></div>
            </div>
            
            <!-- 통계 정보 -->
            <div class="grid grid-cols-3 gap-4 mt-4">
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-3 text-center">
                <div class="text-sm text-gray-600 dark:text-gray-400">{{ t(`report.cardContext.averageLoadRate`) }}<!--평균 부하율--></div>
                <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ averageLoadRate.toFixed(1) }}%</div>
              </div>
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-3 text-center">
                <div class="text-sm text-gray-600 dark:text-gray-400">{{ t(`report.cardContext.maxLoadRate`) }}<!--최대 부하율--></div>
                <div class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ maxLoadRate.toFixed(1) }}%</div>
              </div>
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-3 text-center">
                <div class="text-sm text-gray-600 dark:text-gray-400">{{ t(`report.cardContext.overloadCount`) }}<!--과부하 횟수--></div>
                <div class="text-2xl font-bold" :class="overloadCount > 0 ? 'text-red-600' : 'text-gray-900 dark:text-gray-100'">
                  {{ overloadCount }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Load Pattern Analysis Section -->
      <section v-if="assettypes == 'Transformer'" ref="LoadPatternAnalysis" class="mt-8">
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700/60 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-indigo-500 to-purple-600 p-3 text-white">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
                <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 3h18v18H3zM9 9h6v6H9z"/>
                  <path d="M9 3v6M15 3v6M21 9H9M21 15H9"/>
                </svg>
              </div>
              <div>
                <h3 class="text-xl font-bold">{{ t(`report.cardTitle.LoadPatternAnalysis`) }}<!--부하율 패턴 분석 (요일 × 시간)--></h3>
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="p-6">
            <!-- Heatmap Chart -->
            <div class="bg-gradient-to-br from-gray-50 to-white dark:from-gray-700 dark:to-gray-800 rounded-xl p-4">
              <div ref="loadPatternChart" class="heatmap-chart"></div>
            </div>

            <!-- Weekly Load Pattern Summary -->
            <div class="mt-6">
              <h4 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">{{ t(`report.cardTitle.weeklyLoadDistribution`) }}<!--주간 부하율 분포 통계--></h4>
              <div class="grid grid-cols-4 gap-4">
                <div class="border-l-4 border-green-500 pl-4">
                  <div class="text-sm text-gray-600 dark:text-gray-400 mb-1">{{ t(`report.cardContext.lightLoadHours`) }}<!--경부하 시간--></div>
                  <div class="text-2xl font-bold text-green-600">{{ lightLoadPercentage }}%</div>
                  <div class="text-xs text-gray-500">0-50% {{ t(`report.cardContext.loadRate`) }}<!--부하율--></div>
                </div>
                <div class="border-l-4 border-yellow-500 pl-4">
                  <div class="text-sm text-gray-600 dark:text-gray-400 mb-1">{{ t(`report.cardContext.mediumLoadHours`) }}<!--중간부하 시간--></div>
                  <div class="text-2xl font-bold text-yellow-600">{{ mediumLoadPercentage }}%</div>
                  <div class="text-xs text-gray-500">50-80% {{ t(`report.cardContext.loadRate`) }}<!--부하율--></div>
                </div>
                <div class="border-l-4 border-orange-500 pl-4">
                  <div class="text-sm text-gray-600 dark:text-gray-400 mb-1">{{ t(`report.cardContext.heavyLoadHours`) }}<!--고부하 시간--></div>
                  <div class="text-2xl font-bold text-orange-600">{{ highLoadPercentage }}%</div>
                  <div class="text-xs text-gray-500">80-100% {{ t(`report.cardContext.loadRate`) }}<!--부하율--></div>
                </div>
                <div class="border-l-4 border-red-500 pl-4">
                  <div class="text-sm text-gray-600 dark:text-gray-400 mb-1">{{ t(`report.cardContext.overLoadHours`) }}<!--과부하 시간--></div>
                  <div class="text-2xl font-bold text-red-600">{{ overLoadPercentage }}%</div>
                  <div class="text-xs text-gray-500">100% {{ t(`report.cardContext.overload`) }}<!--초과부하율--></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
import LineChart from '../../../charts/connect/LineChart02.vue'
import BarChart2 from '../../../charts/connect/BarChart01_Energy.vue'
import BarChart from '../../../charts/connect/BarChart03.vue'
import Diagnosis_Barchart from '../diagnosis/Diagnosis_BarChart2.vue'
import DashboardCard_table3 from './DashboardCard_table4.vue'
import Diagnosis_Info from '../diagnosis/Diagnosis_Info2.vue'
import axios from 'axios'
import html2canvas from "html2canvas";
import { useSetupStore } from '@/store/setup';
import { jsPDF } from "jspdf";
import { tailwindConfig } from '../../../utils/Utils'
import { useRoute } from 'vue-router'
import motorImg from '@/images/motor_m.png'
import fanImg from '@/images/fan_m.png'
import pumpImg from '@/images/pump_m.png'
import compImg from '@/images/comp_m.png'
import powerImg from '@/images/power_m.png'
import defaultImg from '@/images/cleaned_logo.png'
import transImg from '@/images/trans.png'
import { useI18n } from 'vue-i18n'
import { useReportData } from '@/composables/reportDict'
import StatusReport from './StatusReport.vue'
import ReportTrend from './ReportTrend.vue'
import dayjs from 'dayjs'
import * as echarts from 'echarts';


export default {
  name: 'PdfView',
  components:{
      LineChart,
      DashboardCard_table3,
      Diagnosis_Info,
      BarChart,
      BarChart2,
      Diagnosis_Barchart,
      ReportTrend,
      StatusReport,
  },
  props:{
    channel:String,
    data:{
      type:Object,
    },
  },
  emits: ['render-complete'],
  setup(props, context){
    
    const emit = context.emit;
    const { t, locale } = useI18n();    
    const data = ref(props.data);
    const datalist = ref({});
    const setupStore = useSetupStore();
    const route = useRoute()
    const tbdata = ref([]);
    const assetdata = ref([]);
    const rawdata = ref([]);
    const pqdata = ref({ items: [], chartdata: null, chartOption: [] });
    const asset = computed(() => setupStore.getAssetConfig);
    const channelInfo = ref(null);
    const PQTable = ref(null);
    const PQChart = ref(null);
    const DiagnosisInfo = ref(null);
    const EnergyHourly = ref(null);
    const EnergyDaily = ref(null);
    const EnergyMonthly = ref(null);
    const LoadPatternAnalysis = ref(null);
    const DiagnosisDetail = ref(null);
    const PQStatusisInfo = ref(null);
    const PQDetail = ref(null);
    const PQTrends = ref(null);
    const PowerLoadTrend = ref(null);  // 전력량 & 부하율 추이 섹션 추가
    const diagData = ref({ items: [], chartdata: null, chartOption: [] });
    const pdf_hidden = ref(false);
    const channelStatus = computed(() => setupStore.getChannelSetting);
    const channelComputed = computed(() => props.channel || route.params.channel || 'Default')
    const setupMenu = ref({});
    const selectedX = ref(0.02);
    const devLocation = computed(() => setupStore.getDevLocation);
    const channel = ref(props.channel);
    const assettypes = computed(()=> {
      const assetInfo = setupStore.getAssetConfig;
      if(channel.value == 'Main')
        return assetInfo.assetType_main;
      else
        return assetInfo.assetType_sub;
    })
    // useReportData에서 실제 데이터 로드 함수들 가져오기
    const { 
      loadInfoData, 
      loadDiagnosisData,
      loadEnergyHourlyData,
      loadEnergyDailyData,
      loadEnergyMonthlyData,
      getLoadFactorCalculated,
      getHeatmapLoadFactorData,
      reportData 
    } = useReportData()

    // Power Quality Trends Variables
    const timeseriesChart = ref(null)
    const timeseriesData = ref([])
    let timeseriesChartInstance = null
    
    // Power Load Trend Variables (이중 Y축 차트)
    const dualAxisChart = ref(null)
    let dualAxisChartInstance = null
    
    // Load Pattern Analysis Variables
    const loadPatternChart = ref(null)
    const loadPatternData = ref([])
    const apiHeatmapData = ref([])
    const apiHeatmapStats = ref({})
    let loadPatternChartInstance = null
    
    const equipImageSrc = computed(() => {
      const eqType = channelComputed.value.toLowerCase() == 'main'? asset.value.assetType_main : asset.value.assetType_sub;
      switch (eqType) {
        case 'Motor':
          return motorImg;
        case 'MotorFeed':
          return motorImg;          
        case 'Pump':
          return pumpImg;//'/images/motor_pump.png';
        case 'Fan':
          return fanImg;
        case 'Compressor':
          return compImg;//'/images/fan_m.png';
        case 'PSupply':
          return powerImg;//'/images/power.png';
        case 'PowerSupply':
          return powerImg;//'/images/power.png';  
        case 'PrimaryTransformer':
          return transImg;//'/images/trans.png';       
        case 'Transformer':
          return transImg;//'/images/trans.png';                  
        default:
          return defaultImg;//'@/images/cleaned_logo.png'
      }
    });

    // formatToISOString 함수 추가
    const formatToISOString = (date, addDay) => {
      const d = new Date(date);
      if (addDay) {
        d.setDate(d.getDate() + addDay);
      }
      return d.toISOString();
    };

    const today = new Date();
    const weekAgo = computed(() => {
      const date = new Date()
      date.setDate(date.getDate() - 7)
      return date
    });

    // 데이터 변환 함수들 추가
    const generateHourlyLabels = () => {
      return Array.from({ length: 24 }, (_, i) => `${i.toString().padStart(2, '0')}:00`)
    }

    const generateDailyLabels = () => {
      const today = dayjs()
      const daysInMonth = today.daysInMonth()
      return Array.from({ length: daysInMonth }, (_, i) =>
        today.date(i + 1).format('MM-DD')
      )
    }

    const generateMonthlyLabels = () => {
      return ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    }

    // API 데이터를 차트 형식으로 변환하는 함수들
    const transformHourlyData = (hourlyData) => {
      //console.log('transformHourlyData 호출됨, 입력 데이터:', hourlyData);
      
      if (!hourlyData || !Array.isArray(hourlyData) || hourlyData.length === 0) {
        console.log('hourly 데이터가 없어서 기본 데이터 생성');
        // 데이터가 없으면 더미 데이터 생성
        return {
          labels: generateHourlyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: 24 }, () => 0),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      //console.log('실제 hourly 데이터 사용, 길이:', hourlyData.length);
      
      const labels = hourlyData.map(item => {
        const hour = item?.hour?.toString().padStart(2, '0') || '00'
        return `${hour}:00`
      })

      const data = hourlyData.map(item => item?.value || 0)

      //console.log('변환된 labels:', labels);
      //console.log('변환된 data:', data);

      return {
        labels,
        datasets: [
          {
            label: 'Import (kWh)',
            data,
            backgroundColor: tailwindConfig().theme.colors.sky[500],
            hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
            barPercentage: 0.7,
            categoryPercentage: 0.7,
            borderRadius: 4,
          },
        ],
      }
    }

    const transformDailyData = (dailyData) => {
      console.log('transformDailyData 호출됨, 입력 데이터:', dailyData);
      
      if (!dailyData || !Array.isArray(dailyData) || dailyData.length === 0) {
        console.log('daily 데이터가 없어서 기본 데이터 생성');
        // 데이터가 없으면 더미 데이터 생성
        return {
          labels: generateDailyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: generateDailyLabels().length }, () => 0),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      const labels = dailyData.map(item => {
        const date = new Date(item?.date || new Date())
        return dayjs(date).format('MM-DD')
      })

      const data = dailyData.map(item => item?.value || 0)

      return {
        labels,
        datasets: [
          {
            label: 'Import (kWh)',
            data,
            backgroundColor: tailwindConfig().theme.colors.sky[500],
            hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
            barPercentage: 0.7,
            categoryPercentage: 0.7,
            borderRadius: 4,
          },
        ],
      }
    }

    const transformMonthlyData = (monthlyData) => {
      console.log('transformMonthlyData 호출됨, 입력 데이터:', monthlyData);
      
      if (!monthlyData || !Array.isArray(monthlyData) || monthlyData.length === 0) {
        console.log('monthly 데이터가 없어서 기본 데이터 생성');
        // 데이터가 없으면 더미 데이터 생성
        return {
          labels: generateMonthlyLabels(),
          datasets: [
            {
              label: 'Import (kWh)',
              data: Array.from({ length: 12 }, () => Math.floor(Math.random() * 3000 + 500)),
              backgroundColor: tailwindConfig().theme.colors.sky[500],
              hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
              barPercentage: 0.7,
              categoryPercentage: 0.7,
              borderRadius: 4,
            },
          ],
        }
      }

      const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
      const labels = monthlyData.map(item => monthNames[(item?.month || 1) - 1])
      const data = monthlyData.map(item => item?.value || 0)

      return {
        labels,
        datasets: [
          {
            label: 'Import (kWh)',
            data,
            backgroundColor: tailwindConfig().theme.colors.sky[500],
            hoverBackgroundColor: tailwindConfig().theme.colors.sky[600],
            barPercentage: 0.7,
            categoryPercentage: 0.7,
            borderRadius: 4,
          },
        ],
      }
    }

    // Energy Chart Data - ref로 변경하여 동적 업데이트 가능하도록
    const hourlyChartData = ref({})
    const dailyChartData = ref({})
    const monthlyChartData = ref({})

    // 에너지 데이터 로드 함수 - reportData 객체를 통해 데이터 접근
    const loadEnergyData = async () => {
      try {
        //console.log('Loading energy data for channel:', channelComputed.value);
        
        // Hourly 데이터 로드
        await loadEnergyHourlyData(channelComputed.value);
        //console.log('Hourly data from reportData:', reportData?.energyHourlyData);
        hourlyChartData.value = transformHourlyData(reportData?.energyHourlyData);
        
        // Daily 데이터 로드
        await loadEnergyDailyData(channelComputed.value);
        //console.log('Daily data from reportData:', reportData?.energyDailyData);
        dailyChartData.value = transformDailyData(reportData?.energyDailyData);
        
        // Monthly 데이터 로드
        await loadEnergyMonthlyData(channelComputed.value);
        //console.log('Monthly data from reportData:', reportData?.energyMonthlyData);
        monthlyChartData.value = transformMonthlyData(reportData?.energyMonthlyData);
        
  
      } catch (error) {
        console.error('Failed to load energy data:', error);
        // 실패 시 더미 데이터 사용
        hourlyChartData.value = transformHourlyData(null);
        dailyChartData.value = transformDailyData(null);
        monthlyChartData.value = transformMonthlyData(null);
      }
    };

    // Diagnosis용 setParamIds 함수
    const setDiagnosisParamIds = async(asset, datalist) => {
      const idList = [];
      try{
        const response = await axios.get(`/api/getParameters/${asset}/diagnostic`);
        if (response.data.success) {
          const paramData = response.data.data;
          //console.log('diagnosis datalist', datalist);
          //console.log('diagnosis paramData', paramData);
          for(let i = 0 ; i < paramData.length ; i++){
            for(let j = 0; j < datalist.length; j++){
              if(paramData[i]["Name"] == datalist[j].Name && paramData[i]["AssemblyID"] == datalist[j].AssemblyID ){
                idList.push({
                  idx: paramData[i]["ID"], 
                  Assembly: paramData[i]["AssemblyID"], 
                  title: paramData[i]["Title"]
                });
              }
            }
          }
        }
      }catch (error){
        console.error(error);
      }
      return idList;
    }

    // PowerQuality용 setParamIds 함수
    const setPQParamIds = async(asset, datalist) => {
      const idList = [];
      try{
        const response = await axios.get(`/api/getParameters/${asset}/powerquality`);
        if (response.data.success) {
          const paramData = response.data.data;
          for(let i = 0 ; i < paramData.length ; i++){
            for(let j = 0; j < datalist.length; j++){
              if(paramData[i]["Name"] == datalist[j].Name && paramData[i]["AssemblyID"] == datalist[j].AssemblyID){
                idList.push({
                  idx: paramData[i]["ID"],
                  Assembly: paramData[i]["AssemblyID"], 
                  title: paramData[i]["Title"]
                });
              }
            }
          }
        }
      }catch (error){
        console.error(error);
      }
      return idList;
    }

    // setChartData 함수
    const setChartData = async (effectiveIds, Title) => {
      let option = {};

      const trendDataRequest = {
        ParametersIds: [effectiveIds],
        StartDate: formatToISOString(weekAgo.value, 0),
        EndDate: formatToISOString(today, 1),
      };

      try {
        const url = `/api/getTrendData`;
        const response = await axios.post(url, trendDataRequest, {
          headers: { "Content-Type": "application/json" },
        });

        if (response.data.success) {
          const resData = response.data.data;
          let datasets = [];
          let labels = [];

          // 1. 메인 데이터 처리 (Thresholds 제외)
          Object.keys(resData).forEach((key) => {
            if (key !== "Thresholds") {
              const dataPoints = resData[key].data;
              if (dataPoints && dataPoints.length > 0) {
                if (labels.length === 0) {
                  labels = dataPoints.map((point) => point.XAxis);
                }
                datasets.push({
                  name: resData[key].Title,
                  data: dataPoints.map((point) => point.YAxis),
                  isThreshold: false,
                });
              }
            }
          });

          // 2. Thresholds 처리
          if (
            resData.Thresholds &&
            resData.Thresholds.length == 2 &&
            labels.length > 0
          ) {
            let timeList = []
            for(let i = 0 ; i < resData.Thresholds.length;i++){
              timeList.push(new Date(resData.Thresholds[i].XAxis));
            }
            const t1 = new Date(resData.Thresholds[0].XAxis);
            const t2 = new Date(resData.Thresholds[1].XAxis);

            resData.Thresholds[0].Thresholds.forEach((value, idx) => {
              if (value !== "NaN" && value !== null && value !== undefined) {
                const secondValue = resData.Thresholds[1].Thresholds[idx];
                if (
                  secondValue === "NaN" ||
                  secondValue === null ||
                  secondValue === undefined
                ) {
                  return;
                }

                const thresholdData = labels.map((lbl) => {
                  const dt = new Date(lbl);
                  return dt < t1 ? value : secondValue;
                });
                const ThresholdString = [
                  "Out of Range(Down side)",
                  "Repair",
                  "Inspect",
                  "Warning",
                  "Warning",
                  "Inspect",
                  "Repair",
                  "Out of Range(Upper side)",
                ];
                datasets.push({
                  name: ThresholdString[idx],
                  data: thresholdData,
                  isThreshold: true,
                });
              }
            });
          }

          option = {
            lineLabels: labels,
            lineData: datasets,
            lineTitle : Title,
          };
        } else {
          console.error("서버 오류:", response.data.error);
          option = {
            lineLabels: [],
            lineData: [],
            lineTitle : '',
          };
        }
      } catch (error) {
        console.error("요청 실패:", error);
      }
      return option;
    };

    // Power Quality Trends - Static Data Generation for PDF
    const generateStaticTrendsData = () => {
      const now = new Date()
      const data = []
      
      for (let i = 23; i >= 0; i--) {
        const time = new Date(now.getTime() - i * 60 * 60 * 1000)
        const hour = time.getHours()
        
        // 시간대별 패턴 시뮬레이션
        let baseVUF = 0.8 + Math.random() * 0.4
        let baseTHDV = 2.0 + Math.random() * 1.0
        let baseTHDI = 3.0 + Math.random() * 1.5
        
        // 특정 시간대에 품질 저하
        if ((hour >= 9 && hour <= 11) || (hour >= 14 && hour <= 16)) {
          baseVUF += Math.random() * 0.8
          baseTHDV += Math.random() * 1.5
          baseTHDI += Math.random() * 2.0
        }
        
        // 가끔 스파이크 발생
        if (Math.random() < 0.05) {
          baseVUF *= 1.5
          baseTHDV *= 1.3
          baseTHDI *= 1.4
        }
        
        data.push({
          time: time.toISOString(),
          vuf: Math.min(baseVUF, 3.0),
          thdv: Math.min(baseTHDV, 8.0),
          thdi: Math.min(baseTHDI, 10.0)
        })
      }
      
      return data
    }

    // Generate static load pattern data for PDF (더미 데이터용)
    const generateLoadPatternData = () => {
      const days = [t('report.cardContext.days.sun'), 
      t('report.cardContext.days.mon'),  t('report.cardContext.days.tue'),  
      t('report.cardContext.days.wed'),  t('report.cardContext.days.thu'),  
      t('report.cardContext.days.fri'),  t('report.cardContext.days.sat')]
      const hours = Array.from({length: 24}, (_, i) => i.toString().padStart(2, '0'))
      const data = []
      
      days.forEach((day, dayIdx) => {
        hours.forEach((hour, hourIdx) => {
          let baseLoad = 25
          
          if (hourIdx >= 8 && hourIdx <= 18) {
            baseLoad = 50 + Math.sin((hourIdx - 8) / 10 * Math.PI) * 30
          } else if (hourIdx >= 19 && hourIdx <= 22) {
            baseLoad = 40 + Math.random() * 20
          }
          
          if (dayIdx === 0 || dayIdx === 6) {
            baseLoad *= 0.7
          }
          
          if (hourIdx === 12 || hourIdx === 18) {
            baseLoad *= 1.2
          }
          
          const loadRate = Math.max(10, baseLoad + (Math.random() - 0.5) * 20)
          
          data.push([hourIdx, dayIdx, Math.round(loadRate)])
        })
      })
      
      return { data, days, hours }
    }

    // Load Pattern Analysis 실제 데이터 로드 함수
    const loadHeatmapData = async () => {
      try {
        console.log('히트맵 API 데이터 로딩 시작:', channelComputed.value);
        
        const response = await getHeatmapLoadFactorData(channelComputed.value, 4);
        
        //console.log('히트맵 API 응답:', response);
        
        if (response && response.success) {
          apiHeatmapData.value = response.heatmapData || [];
          apiHeatmapStats.value = response.distribution || {};
          
          //console.log('히트맵 데이터 개수:', apiHeatmapData.value.length);
          //console.log('히트맵 통계:', apiHeatmapStats.value);
          
          // 히트맵 데이터가 있으면 차트 업데이트
          if (loadPatternChartInstance && apiHeatmapData.value.length > 0) {
            createLoadPatternChart();
          }
        } else {
          console.error('히트맵 API 실패, 더미 데이터 사용');
          // API 실패 시 더미 데이터 사용
          loadPatternData.value = generateLoadPatternData();
          if (loadPatternChartInstance) {
            createLoadPatternChart();
          }
        }
      } catch (error) {
        console.error('히트맵 데이터 로딩 실패:', error);
        // 에러 시 더미 데이터 사용
        loadPatternData.value = generateLoadPatternData();
        if (loadPatternChartInstance) {
          createLoadPatternChart();
        }
      }
    }

    // Computed properties for averages
    const averageVUF = computed(() => {
      if (timeseriesData.value.length === 0) return 0
      const sum = timeseriesData.value.reduce((acc, item) => acc + item.vuf, 0)
      return sum / timeseriesData.value.length
    })

    const averageTHDV = computed(() => {
      if (timeseriesData.value.length === 0) return 0
      const sum = timeseriesData.value.reduce((acc, item) => acc + item.thdv, 0)
      return sum / timeseriesData.value.length
    })

    const averageTHDI = computed(() => {
      if (timeseriesData.value.length === 0) return 0
      const sum = timeseriesData.value.reduce((acc, item) => acc + item.thdi, 0)
      return sum / timeseriesData.value.length
    })

    // 부하율 관련 computed 속성들
    const averageLoadRate = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        const sum = reportData.loadrateData.reduce((acc, item) => acc + (item.load_factor_percent || 0), 0);
        return sum / reportData.loadrateData.length;
      }
      // 더미 데이터 사용
      return 0;
    });

    const maxLoadRate = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        return Math.max(...reportData.loadrateData.map(item => item.load_factor_percent || 0));
      }
      // 더미 데이터 사용
      return 0;
    });

    const overloadCount = computed(() => {
      if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
        return reportData.loadrateData.filter(item => (item.load_factor_percent || 0) > 100).length;
      }
      // 더미 데이터 사용
      return 0;
    });

    // 이중 Y축 차트 생성 함수
const createDualAxisChart = () => {
  if (!dualAxisChartInstance) return;
  
  console.log('=== createDualAxisChart 실행 ===');
  //console.log('reportData.energyHourlyData:', reportData.energyHourlyData);
  //console.log('reportData.loadrateData:', reportData.loadrateData);

  let times = [];
  let powerData = [];
  let loadData = [];
  
  // 실제 전력량 데이터 처리
  if (reportData.energyHourlyData && Array.isArray(reportData.energyHourlyData) && reportData.energyHourlyData.length > 0) {
    //'실제 전력량 데이터 사용');
    
    const now = dayjs();
    const currentHour = now.hour();
    
    const powerDataMap = {};
    reportData.energyHourlyData.forEach(item => {
      if (item?.hour !== undefined && item?.value !== undefined) {
        const hourLabel = `${item.hour.toString().padStart(2, '0')}:00`;
        powerDataMap[hourLabel] = parseFloat(item.value) || 0;
      }
    });
    
    for (let i = 0; i <= Math.min(currentHour, 23); i++) {
      const hourLabel = `${i.toString().padStart(2, '0')}:00`;
      times.push(hourLabel);
      powerData.push(powerDataMap[hourLabel] || 0);
    }
  } else {
    // 더미 데이터 생성 - 모든 값을 0으로
    //console.log('전력량 더미 데이터 사용');
    for (let i = 0; i < 24; i++) {
      times.push(`${i.toString().padStart(2, '0')}:00`);
      powerData.push(0); // ✅ 모든 값을 0으로 변경
    }
  }

  // 실제 부하율 데이터 처리
  if (reportData.loadrateData && Array.isArray(reportData.loadrateData) && reportData.loadrateData.length > 0) {
    //console.log('실제 부하율 데이터 사용');
    
    const loadDataMap = {};
    reportData.loadrateData.forEach(item => {
      if (item?.hour !== undefined && item?.load_factor_percent !== undefined) {
        const hourLabel = `${item.hour.toString().padStart(2, '0')}:00`;
        loadDataMap[hourLabel] = parseFloat(item.load_factor_percent) || 0;
      }
    });
    
    loadData = times.map(time => loadDataMap[time] || 0);
  } else {
    // 더미 데이터 생성 - 모든 값을 0으로
    console.log('부하율 더미 데이터 사용');
    loadData = times.map(() => 0); // ✅ 모든 값을 0으로 변경
  }

  const maxPowerValue = Math.max(...powerData.filter(val => !isNaN(val)), 100);
  const maxLoadValue = Math.max(...loadData.filter(val => !isNaN(val)), 100);
  
  const powerAxisMax = Math.ceil(maxPowerValue * 1.2);
  const loadAxisMax = Math.ceil(Math.max(maxLoadValue * 1.2, 120));

  const option = {
    grid: {
      left: '10%',
      right: '10%',
      top: '12%',
      bottom: '10%',
      containLabel: true
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross', crossStyle: { color: '#999' } },
      formatter: function (params) {
        let result = params[0].axisValueLabel + '<br/>';
        params.forEach(function (item) {
          if (item.seriesName === '전력량') {
            result += item.marker + item.seriesName + ': ' + item.value.toFixed(2) + ' kWh<br/>';
          } else {
            result += item.marker + item.seriesName + ': ' + item.value.toFixed(1) + '%<br/>';
          }
        });
        return result;
      }
    },
    legend: { 
      data: ['전력량', '부하율'],
      top: '2%',
      textStyle: { fontSize: 12 }
    },
    xAxis: [{ 
      type: 'category', 
      data: times,
      boundaryGap: false,
      axisLabel: {
        formatter: function (value) {
          if (times.length > 12) {
            const hourNum = parseInt(value.split(':')[0]);
            return hourNum % 2 === 0 ? value : '';
          }
          return value;
        },
        rotate: times.length > 20 ? 45 : 0
      }
    }],
    yAxis: [
      {
        type: 'value',
        name: t(`report.cardTitle.Energy`)+'(kWh)',
        position: 'left',
        axisLabel: { 
          formatter: '{value}'
        },
        splitLine: { show: true, lineStyle: { type: 'dashed', color: '#e0e0e0' } },
        min: 0,
        max: powerAxisMax
      },
      {
        type: 'value',
        name: t(`report.cardContext.loadRate`)+'(%)',
        position: 'right',
        axisLabel: { 
          formatter: '{value}%'
        },
        min: 0,
        max: loadAxisMax,
        splitLine: { show: false }
      }
    ],
    series: [
      {
        name: t(`report.cardTitle.Energy`),
        type: 'bar',
        yAxisIndex: 0,
        data: powerData,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        }
      },
      {
        name: t(`report.cardContext.loadRate`),
        type: 'line',
        yAxisIndex: 1,
        data: loadData,
        lineStyle: { width: 3, color: '#ff6b6b' },
        itemStyle: { color: '#ff6b6b' },
        areaStyle: { opacity: 0.1, color: '#ff6b6b' },
        smooth: true
      }
    ]
  };
  
  dualAxisChartInstance.setOption(option);
};

    // Load distribution statistics - API 데이터 기반
// Load distribution statistics - API 데이터 기반
const lightLoadPercentage = computed(() => {
  // API 데이터가 있으면 사용
  if (apiHeatmapStats.value && Object.keys(apiHeatmapStats.value).length > 0) {
    const total = (apiHeatmapStats.value.light || 0) + 
                 (apiHeatmapStats.value.medium || 0) + 
                 (apiHeatmapStats.value.high || 0) + 
                 (apiHeatmapStats.value.overload || 0);
    
    if (total === 0) return 0;
    return Math.round(((apiHeatmapStats.value.light || 0) / total) * 100);
  }
  
  // 더미 데이터 fallback - 모든 값을 0으로
  return 0; // ✅ 더미 데이터 계산 대신 0 반환
})

const mediumLoadPercentage = computed(() => {
  // API 데이터가 있으면 사용
  if (apiHeatmapStats.value && Object.keys(apiHeatmapStats.value).length > 0) {
    const total = (apiHeatmapStats.value.light || 0) + 
                 (apiHeatmapStats.value.medium || 0) + 
                 (apiHeatmapStats.value.high || 0) + 
                 (apiHeatmapStats.value.overload || 0);
    
    if (total === 0) return 0;
    return Math.round(((apiHeatmapStats.value.medium || 0) / total) * 100);
  }
  
  // 더미 데이터 fallback - 모든 값을 0으로
  return 0; // ✅ 더미 데이터 계산 대신 0 반환
})

const highLoadPercentage = computed(() => {
  // API 데이터가 있으면 사용
  if (apiHeatmapStats.value && Object.keys(apiHeatmapStats.value).length > 0) {
    const total = (apiHeatmapStats.value.light || 0) + 
                 (apiHeatmapStats.value.medium || 0) + 
                 (apiHeatmapStats.value.high || 0) + 
                 (apiHeatmapStats.value.overload || 0);
    
    if (total === 0) return 0;
    return Math.round(((apiHeatmapStats.value.high || 0) / total) * 100);
  }
  
  // 더미 데이터 fallback - 모든 값을 0으로
  return 0; // ✅ 더미 데이터 계산 대신 0 반환
})

const overLoadPercentage = computed(() => {
  // API 데이터가 있으면 사용
  if (apiHeatmapStats.value && Object.keys(apiHeatmapStats.value).length > 0) {
    const total = (apiHeatmapStats.value.light || 0) + 
                 (apiHeatmapStats.value.medium || 0) + 
                 (apiHeatmapStats.value.high || 0) + 
                 (apiHeatmapStats.value.overload || 0);
    
    if (total === 0) return 0;
    return Math.round(((apiHeatmapStats.value.overload || 0) / total) * 100);
  }
  
  // 더미 데이터 fallback - 모든 값을 0으로
  return 0; // ✅ 더미 데이터 계산 대신 0 반환
})

    // Create timeseries chart
    const createTimeseriesChart = () => {
      if (!timeseriesChartInstance || timeseriesData.value.length === 0) return
      
      const times = timeseriesData.value.map(item => new Date(item.time))
      const vufData = timeseriesData.value.map(item => item.vuf.toFixed(2))
      const thdvData = timeseriesData.value.map(item => item.thdv.toFixed(2))
      const thdiData = timeseriesData.value.map(item => item.thdi.toFixed(2))
      
      const option = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          },
          formatter: function (params) {
            let result = params[0].axisValueLabel + '<br/>'
            params.forEach(function (item) {
              result += item.marker + item.seriesName + ': ' + item.value + '%<br/>'
            })
            return result
          }
        },
        legend: {
          data: ['불평형률', '전압 고조파왜곡률', '전류 고조파왜곡률']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            boundaryGap: false,
            data: times.map(time => time.toLocaleTimeString('ko-KR', {
              hour: '2-digit',
              minute: '2-digit'
            }))
          }
        ],
        yAxis: [
          {
            type: 'value',
            name: '왜곡률 (%)',
            min: 0,
            max: 10,
            axisLabel: {
              formatter: '{value}%'
            }
          }
        ],
        series: [
          {
            name: '불평형률',
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: vufData,
            color: '#ef4444',
            markLine: {
              data: [{ yAxis: 2, name: '불평형률 기준선' }]
            }
          },
          {
            name: '전압 고조파왜곡률',
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: thdvData,
            color: '#14b8a6',
            markLine: {
              data: [{ yAxis: 5, name: 'THDV 기준선' }]
            }
          },
          {
            name: '전류 고조파왜곡률',
            type: 'line',
            smooth: true,
            lineStyle: { width: 2 },
            showSymbol: false,
            areaStyle: { opacity: 0.1 },
            data: thdiData,
            color: '#3b82f6',
            markLine: {
              data: [{ yAxis: 5, name: 'THDI 기준선' }]
            }
          }
        ]
      }
      
      timeseriesChartInstance.setOption(option)
    }

    // Create load pattern heatmap chart - API 데이터 우선 사용
    const createLoadPatternChart = () => {
      if (!loadPatternChartInstance) return
      
      // API 데이터가 있으면 우선 사용
      if (apiHeatmapData.value && apiHeatmapData.value.length > 0) {
        console.log('API 히트맵 데이터로 차트 생성');
        
        const days = [t('report.cardContext.days.sun'), 
      t('report.cardContext.days.mon'),  t('report.cardContext.days.tue'),  
      t('report.cardContext.days.wed'),  t('report.cardContext.days.thu'),  
      t('report.cardContext.days.fri'),  t('report.cardContext.days.sat')]
        const hours = Array.from({length: 24}, (_, i) => i.toString().padStart(2, '0'))
        
        // API 데이터 형식 변환
        let chartData = [];
        if (apiHeatmapData.value.length > 0 && Array.isArray(apiHeatmapData.value[0]) && apiHeatmapData.value[0].length === 3) {
          chartData = apiHeatmapData.value;
        } else {
          chartData = apiHeatmapData.value.map(item => {
            if (Array.isArray(item)) {
              return item;
            } else if (typeof item === 'object') {
              return [item.hour || 0, item.day_of_week || 0, item.load_factor_percent || 0];
            } else {
              return [0, 0, 0];
            }
          });
        }
        
        const values = chartData.map(item => item[2] || 0);
        const maxValue = Math.max(...values, 100);
        
        const option = {
          tooltip: {
            position: 'top',
            formatter: function (params) {
              const hour = hours[params.data[0]]
              const day = days[params.data[1]]
              const value = params.data[2]
              return `${day}t('report.cardContext.days.week') ${hour}:00<br/>t('report.cardContext.loadRate') : ${value}%`
            }
          },
          grid: { height: '50%', top: '10%' },
          xAxis: {
            type: 'category',
            data: hours,
            splitArea: { show: true },
            axisLabel: { 
              formatter: function (value) { 
                const hourNum = parseInt(value);
                return hourNum % 2 === 0 ? value + ':00' : '';
              }
            }
          },
          yAxis: { 
            type: 'category', 
            data: days, 
            splitArea: { show: true } 
          },
          visualMap: {
            min: 0,
            max: Math.max(maxValue * 1.2, 100),
            calculable: true,
            orient: 'horizontal',
            left: 'center',
            bottom: '15%',
            inRange: {
              color: ['#313695', '#4575b4', '#74add1', '#abd9e9', '#e0f3f8', '#fee090', '#fdae61', '#f46d43', '#d73027', '#a50026']
            },
            text: [ t('report.cardContext.hight'),t('report.cardContext.low')]
          },
          series: [{
            name: t('report.cardContext.loadRate'),
            type: 'heatmap',
            data: chartData,
            label: { show: false },
            emphasis: { 
              itemStyle: { 
                shadowBlur: 10, 
                shadowColor: 'rgba(0, 0, 0, 0.5)' 
              } 
            }
          }]
        }
        
        loadPatternChartInstance.setOption(option)
        
      } else if (loadPatternData.value?.data) {
        // 더미 데이터로 fallback
        console.log('더미 히트맵 데이터로 차트 생성');
        
        const { data, days, hours } = loadPatternData.value
        
        // ✅ 더미 데이터를 모두 0으로 변경
        const zeroData = data.map(item => [item[0], item[1], 0]);
        
        const option = {
          tooltip: {
            position: 'top',
            formatter: function (params) {
              const hour = hours[params.data[0]]
              const day = days[params.data[1]]
              const value = params.data[2]
              //return `${day}요일 ${hour}:00<br/>부하율: ${value}%`
              return `${day}${t('report.cardContext.days.week') }${hour}:00<br/>${t('report.cardContext.loadRate') }: ${value}%`
          
            }
          },
          grid: { height: '50%', top: '10%' },
          xAxis: {
            type: 'category',
            data: hours,
            splitArea: { show: true },
            axisLabel: { formatter: function (value) { return value + ':00' } }
          },
          yAxis: { type: 'category', data: days, splitArea: { show: true } },
          visualMap: {
            min: 0,
            max: 100,
            calculable: true,
            orient: 'horizontal',
            left: 'center',
            bottom: '15%',
            inRange: {
              color: ['#313695', '#4575b4', '#74add1', '#abd9e9', '#e0f3f8', '#fee090', '#fdae61', '#f46d43', '#d73027', '#a50026']
            },
            text: [t('report.cardContext.hight'),t('report.cardContext.low')]
          },
          series: [{
            name: t('report.cardContext.loadRate'),
            type: 'heatmap',
            data: zeroData, // ✅ 0으로 변경된 데이터 사용
            label: { show: false },
            emphasis: { itemStyle: { shadowBlur: 10, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
          }]
        }
        
        loadPatternChartInstance.setOption(option)
      }
    }

    const baseChart = {
      datasets: [
        {
          label: 'Series 0',
          data: [
            { x: 0.0001, y: 500 },
            { x: 0.001, y: 200 },
            { x: 0.003, y: 140 },
            { x: 0.003, y: 120 },
            { x: 0.02, y: 120 },
            { x: 0.5, y: 120 },
            { x: 0.5, y: 110 },
            { x: 10, y: 110 },
            { x: 100, y: 110 },
          ],
          borderColor: 'blue',
          backgroundColor: 'transparent',
          borderWidth: 2,
          pointRadius: 0,
          tension: 0,
        },
        {
          label: 'Series 1',
          data: [
            { x: 0.02, y: 0 },
            { x: 0.02, y: 70 },
            { x: 0.5, y: 70 },
            { x: 0.5, y: 80 },
            { x: 10, y: 80 },
            { x: 10, y: 90 },
            { x: 100, y: 90 },
          ],
          borderColor: 'red',
          backgroundColor: 'transparent',
          borderWidth: 2,
          pointRadius: 0,
          tension: 0,
        }
      ]
    }

    const linechartData = computed(() => {
      const base = JSON.parse(JSON.stringify(baseChart));
      if (selectedX.value != null) {
        base.datasets.push({
            label: 'Selected Point',
            type: 'scatter',
            data: [{ x: selectedX.value, y: 125 }],
            backgroundColor: 'orange',
            pointRadius: 6,
            pointHoverRadius: 8,
            showLine: false,
          })
      }
      return base
    })

    const ischartDataValid = computed(() => {
      return diagData.value.chartdata && 
     diagData.value.chartdata.Names && 
     Array.isArray(diagData.value.chartdata.Names) &&
     diagData.value.chartdata.Names.length > 0 &&
     diagData.value.chartdata.Values &&
     Array.isArray(diagData.value.chartdata.Values) &&
     diagData.value.chartdata.Values.length > 0
    })

    const ispqValid = computed(() => {
      return pqdata.value.chartdata && 
     pqdata.value.chartdata.Names && 
     Array.isArray(pqdata.value.chartdata.Names) &&
     pqdata.value.chartdata.Names.length > 0 &&
     pqdata.value.chartdata.Values &&
     Array.isArray(pqdata.value.chartdata.Values) &&
     pqdata.value.chartdata.Values.length > 0
    })

    const fetchData = async () => {
        if(!asset.value)
          await setupStore.checkSetting(); 
        const chName = channelComputed.value.toLowerCase() == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
        if(chName != ''){
            const infos = await loadInfoData(chName);
            datalist.value = infos;
        }else{
            console.log('There are no registered Asset.');
        }
    };

    watch(channelStatus, (newVal) => {
      setupMenu.value = newVal;
    }, { immediate: true });

    const mode = computed(() => {
      if (channelComputed.value === 'Main')
          return channelComputed.value === 'Main' && setupMenu.value.MainDiagnosis;
      else
          return channelComputed.value === 'Sub' && setupMenu.value.SubDiagnosis;
    });

    const COMPLIANCE_BITS = {
        "Frequency Variation 1": [0],
        "Frequency Variation 2": [1],
        "Voltage Variation 1": [2, 3, 4],
        "Voltage Variation 2": [5, 6, 7],
        "Voltage Unbalance": [8],
        "THD": [9, 10, 11],
        "Harmonics": [12, 13, 14],
        "Pst": [15, 16, 17],
        "Plt": [18, 19, 20],
        "Signal Vol.": [21, 22, 23]
    }

    const makeKey = (param, phase) => {
      const suffixMap = {
        L1: "L1",
        L2: "L2",
        L3: "L3",
        "Multi Phase": "Multi Phase"
      }

      if (param === "Frequency Variation 1") {
        return phase === 'L1' ? "Frequency Variation 1(%)" : undefined
      }
      if (param === "Frequency Variation 2") {
        return phase === 'L1' ? "Frequency Variation 2(%)" : undefined
      }
      if (param === "Voltage Unbalance") {
        return phase === 'L1' ? "Voltage Unbalance(%)" : undefined
      }
      if (param === "Voltage Variation 1") return `Voltage Variation 1 ${suffixMap[phase]}(%)`
      if (param === "Voltage Variation 2") return `Voltage Variation 2 ${suffixMap[phase]}(%)`
      if (param === "THD") return `THDs Variation ${suffixMap[phase]}(%)`
      if (param === "Harmonics") return `Harmonics Variatiopn ${suffixMap[phase]}(%)`
      if (param === "Pst") return `Flickers Pst ${suffixMap[phase]}(%)`
      if (param === "Plt") return `Flickers Plt ${suffixMap[phase]}(%)`
      if (param === "Signal Vol.") return `Signaling Voltage ${suffixMap[phase]}(%)`
      if (param === "Voltage Sag") return `Voltage Dips ${suffixMap[phase]}`
      if (param === "Voltage Swell") return `Voltage Swells ${suffixMap[phase]}`
      if (param === "Short Interruption") return `Short Interruptions ${suffixMap[phase]}`
      if (param === "Long Interruption") return `Long Interruptions ${suffixMap[phase]}`
      if (param === "Signaling Volt.") return `Signaling Voltage ${suffixMap[phase]}(%)`

      return `${param} ${suffixMap[phase]}`
    }
    
    const getComp = (param) => {
      const comp = tbdata.value["status&compliance"]
      const bits = COMPLIANCE_BITS[param]
      if (!bits) return "-"

      const isAnyOK = bits.some(bit => (comp & (1 << bit)) !== 0)
      return isAnyOK ? "Failed" : "OK"
    }

    const getComplianceClass = (status) => {
      if (status === 'OK') {
        return 'text-green-800'
      } else if (status === 'Failed') {
        return 'text-red-800'
      } else {
        return 'text-gray-800'
      }
    }

    const fetchEn = async() =>{
      try {
          const response = await axios.get(`/api/getEn50160/${channelComputed.value}`);
          if (response.data.success) {
              tbdata.value = response.data.data
          } else {
              console.warn("서버 응답이 success: false 입니다.");
              tbdata.value = {};
          }
        } catch (error) {
          console.log("데이터 가져오기 실패:", error);
        }
    }
     
    // fetchRawData 함수를 더 정확하게 구현 (Diagnosis용)
    const fetchRawData = async() => {
      diagData.value = { items: [], chartdata: null, chartOption: [] };
      const chName = channelComputed.value.toLowerCase() == 'main'? asset.value.assetName_main : asset.value.assetName_sub;
      
      try {
        const response = await axios.get(`/api/getDiagnosisDetail/${chName}`);

        if (response.data.success) {
          let itemlist = [], valuelist=[], datalist=[];
          for(let i = 0; i < response.data.data_status.length;i++){
            itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"])
          }
          diagData.value.chartdata = {"Names" : itemlist, "Values" : valuelist, "Titles": datalist};

          const treeItems = response.data.data_tree;
          const chartList = [];
          
          for(let i = 0 ; i< treeItems.length;i++){
            if(treeItems[i]["Status"] > 1 ){
              const childDict = [];
              if(treeItems[i]["children"].length > 0){
                for(let j = 0 ; j < treeItems[i]["children"].length ; j++ ){
                  if(treeItems[i]["children"][j]["isSub"]){
                    for(let k = 0 ; k < treeItems[i]["children"][j]["children"].length ; k++ ){
                      if( treeItems[i]["children"][j]["children"][k]["Status"] > 1){
                        chartList.push(treeItems[i]["children"][j]["children"][k]);
                        childDict.push({
                          Title:treeItems[i]["children"][j]["children"][k]["Title"], 
                          Assembly:treeItems[i]["children"][j]["AssemblyID"], 
                          Value:treeItems[i]["children"][j]["children"][k]["Value"]
                        })
                      }
                    }
                  }else{
                    if(treeItems[i]["children"][j]["Status"] > 1){
                      chartList.push(treeItems[i]["children"][j]);
                      childDict.push({
                        Title:treeItems[i]["children"][j]["Title"],
                        Assembly:treeItems[i]["children"][j]["AssemblyID"], 
                        Value:treeItems[i]["children"][j]["Value"]
                      })
                    }
                  }
                }
              }else{
                childDict = [{
                  Title:treeItems[i]["Title"],
                  Assembly:treeItems[i]["AssemblyID"], 
                  Value:treeItems[i]["Value"]
                }]
              }
              diagData.value.items.push({Item:treeItems[i], Child:childDict});            
            }
          }

          // 차트 데이터 생성 (수정된 부분 - setDiagnosisParamIds 사용)
          if(chartList.length > 0){
            let chartValue=null;
            const effectiveIds = await setDiagnosisParamIds(chName, chartList)  // Diagnosis 전용 함수 사용
            let idxList=[],idList=[];
            for (let i = 0 ; i < effectiveIds.length;i++){
              if (idxList.includes(effectiveIds[i].idx)){
                continue;
              }else{
                idList.push(effectiveIds[i])
                idxList.push(effectiveIds[i].idx)
              }
            }
            console.log('Diagnosis PDF idxList:', idxList);
            
            for (let i = 0 ; i < idList.length;i++){
              const titleName = '['+idList[i].Assembly + ']' + idList[i].title;
              chartValue = await setChartData(idList[i].idx, titleName);
              diagData.value.chartOption.push(chartValue);
            }
          }

        }else{
          console.log('No Diagnosis Data');
        }
      }catch (error) {
        console.log("Diagnosis 데이터 가져오기 실패:", error);
      }
    }

    // fetchPQData 함수에서 setPQParamIds 사용하도록 수정
    const fetchPQData = async() => {
      pqdata.value = { items: [], chartdata: null, chartOption: [] };
      
      const chName = channelComputed.value.toLowerCase() == 'main' ? asset.value.assetName_main : asset.value.assetName_sub;
      
      try {
        const response = await axios.get(`/api/getDiagPQ/${chName}`);

        if (response.data.success) {
          let itemlist = [], valuelist = [], datalist = [];
          
          for (let i = 0; i < response.data.data_status.length; i++) {
            itemlist.push(response.data.data_status[i]["Titles"][locale.value]);
            valuelist.push(response.data.data_status[i]["Status"]);
            datalist.push(response.data.data_status[i]["Titles"]);
          }
          
          pqdata.value.chartdata = {"Names": itemlist, "Values": valuelist, "Titles": datalist};
          
          const treeItems = response.data.data_tree;
          const chartList = [];
          
          for (let i = 0; i < treeItems.length; i++) {
            if (treeItems[i]["Status"] > 1) {
              const childDict = [];
              
              if (treeItems[i]["children"].length > 0) {
                for (let j = 0; j < treeItems[i]["children"].length; j++) {
                  if (treeItems[i]["children"][j]["isSub"]) {
                    for (let k = 0; k < treeItems[i]["children"][j]["children"].length; k++) {
                      if (treeItems[i]["children"][j]["children"][k]["Status"] > 1) {
                        chartList.push(treeItems[i]["children"][j]["children"][k]);
                        childDict.push({
                          Title: treeItems[i]["children"][j]["children"][k]["Title"],
                          Assembly: treeItems[i]["children"][j]["AssemblyID"],
                          Value: treeItems[i]["children"][j]["children"][k]["Value"]
                        });
                      }
                    }
                  } else {
                    if (treeItems[i]["children"][j]["Status"] > 1) {
                      chartList.push(treeItems[i]["children"][j]);
                      childDict.push({
                        Title: treeItems[i]["children"][j]["Title"],
                        Assembly: treeItems[i]["children"][j]["AssemblyID"],
                        Value: treeItems[i]["children"][j]["Value"]
                      });
                    }
                  }
                }
              } else {
                childDict = [{
                  Title: treeItems[i]["Title"],
                  Assembly: treeItems[i]["AssemblyID"],
                  Value: treeItems[i]["Value"]
                }];
              }
              
              pqdata.value.items.push({Item: treeItems[i], Child: childDict});
            }
          }
          
          // 차트 데이터 생성 (수정된 부분 - setPQParamIds 사용)
          if (chartList.length > 0) {
            let chartValue = null;
            const effectiveIds = await setPQParamIds(chName, chartList);  // PowerQuality 전용 함수 사용
            let idxList = [], idList = [];
            
            for (let i = 0; i < effectiveIds.length; i++) {
              if (idxList.includes(effectiveIds[i].idx)) {
                continue;
              } else {
                idList.push(effectiveIds[i]);
                idxList.push(effectiveIds[i].idx);
              }
            }
            
            for (let i = 0; i < idList.length; i++) {
              const titleName = '[' + idList[i].Assembly + ']' + idList[i].title;
              chartValue = await setChartData(idList[i].idx, titleName);
              pqdata.value.chartOption.push(chartValue);
            }
            
            console.log('PQ PDF idList:', idxList);
          }
        } else {
          console.log('No PQ Data');
        }
      } catch (error) {
        console.log("PQ 데이터 가져오기 실패:", error);
      }
    };

    onMounted(async ()=>{
      console.log('PdfView component mounted');      
      
      // 기본 데이터 로드
      fetchData();
      fetchEn();
      fetchPQData();
      fetchRawData();
      
      // 에너지 데이터 로드 추가
      console.log('Starting energy data load...');
      await loadEnergyData();
      console.log('Energy data load complete');
      
      // 부하율 데이터 로드
      console.log('Starting load rate data load...');
      try {
        await getLoadFactorCalculated(channelComputed.value);
      } catch (error) {
        console.error('Failed to load load rate data:', error);
      }
      
      // 부하율 히트맵 데이터 로드
      console.log('Starting heatmap data load...');
      await loadHeatmapData();
      console.log('Heatmap data load complete');
      
      // Initialize trends data and chart
      timeseriesData.value = generateStaticTrendsData()
      
      // 히트맵 데이터가 없으면 더미 데이터 사용
      if (!apiHeatmapData.value || apiHeatmapData.value.length === 0) {
        loadPatternData.value = generateLoadPatternData()
      }
      
      nextTick(() => {
        if (timeseriesChart.value) {
          timeseriesChartInstance = echarts.init(timeseriesChart.value)
          createTimeseriesChart()
          
          window.addEventListener('resize', () => {
            timeseriesChartInstance?.resize()
          })
        }
        
        if (dualAxisChart.value) {
          dualAxisChartInstance = echarts.init(dualAxisChart.value)
          createDualAxisChart()
          
          window.addEventListener('resize', () => {
            dualAxisChartInstance?.resize()
          })
        }
        
        if (loadPatternChart.value) {
          loadPatternChartInstance = echarts.init(loadPatternChart.value)
          createLoadPatternChart()
          
          window.addEventListener('resize', () => {
            loadPatternChartInstance?.resize()
          })
        }
      })
      // 모든 차트가 렌더링된 후 약간의 지연을 주고 완료 이벤트 발생
      setTimeout(() => {
        console.log('PDF rendering complete, emitting event');
        emit('render-complete');
      }, 1000);

    });

    // Cleanup on unmount
    onUnmounted(() => {
      timeseriesChartInstance?.dispose()
      dualAxisChartInstance?.dispose()
      loadPatternChartInstance?.dispose()
      window.removeEventListener('resize', () => {})
    })

    // PDF 데이터 테이블 더미 데이터
    const pqTableData = [
      { parameter: 'Frequency Variation 1', required: '±1% 99.5%' },
      { parameter: 'Frequency Variation 2', required: '+4%/-6% 100%' },
      { parameter: 'Voltage Variation 1', required: '±10% 95%' },
      { parameter: 'Voltage Variation 2', required: '+10%/-15% 100%' },
      { parameter: 'Voltage Unbalance', required: '2% 95%' },
      { parameter: 'THD', required: '8% 95%' },
      { parameter: 'Harmonics', required: 'IEC 61000-2-2' },
      { parameter: 'Pst', required: '1.2 95%' },
      { parameter: 'Plt', required: '1.0 95%' },
      { parameter: 'Signal Vol.', required: '9% 95%' },
    ]

const generatePDF = async() => {
  pdf_hidden.value = true;
  
  document.documentElement.classList.add('generating-pdf');
  await new Promise(resolve => setTimeout(resolve, 100));
  
  const pdf = new jsPDF({
    orientation: "p",
    unit: "mm",
    format: "a4",
    compress: true
  });

  const canvasOptions = {
    scale: 1.5,
    useCORS: true,
    logging: false,
    imageTimeout: 0,
    backgroundColor: '#ffffff',
    removeContainer: true,
    letterRendering: true,
    allowTaint: true,
    foreignObjectRendering: false,
    windowWidth: document.documentElement.scrollWidth,
    windowHeight: document.documentElement.scrollHeight
  };

  let pageCount = 0;

  try {
    // 1. Channel Info
    if (channelInfo.value) {
      await processSection(pdf, channelInfo.value, 'Channel Info', canvasOptions, pageCount++ > 0);
    }

    // 2. Diagnosis Status
    if (DiagnosisInfo.value && mode.value) {
      await processSection(pdf, DiagnosisInfo.value, 'Diagnosis Status', canvasOptions, pageCount++ > 0);
    }

    // 3. Diagnosis Detail - 전체 섹션을 캡처하되 차트는 숨기고 헤더만
    if (DiagnosisDetail.value && mode.value && ischartDataValid.value && diagData.value.items.length > 0) {
      pdf.addPage();
      
      // 먼저 차트들을 임시로 숨김
      const trendCharts = DiagnosisDetail.value.querySelectorAll('.col-span-12');
      trendCharts.forEach(chart => {
        chart.style.display = 'none';
      });
      
      // 헤더만 있는 상태로 캡처
      await processSection(pdf, DiagnosisDetail.value, 'Diagnosis Detail Header', canvasOptions, false);
      
      // 차트들을 다시 표시
      trendCharts.forEach(chart => {
        chart.style.display = '';
      });
      
      // 각 트렌드 차트를 개별적으로 처리
      let chartY = 60; // 헤더 아래 시작 위치
      
      for (let i = 0; i < trendCharts.length; i++) {
        if (i % 2 === 0 && i > 0) {
          pdf.addPage();
          chartY = 10;
        }
        await processChartElement(pdf, trendCharts[i], `Diagnosis Trend ${i + 1}`, canvasOptions, chartY);
        chartY += 140;
      }
      pageCount++;
    }

    // 4. PQ Status
    if (PQStatusisInfo.value && mode.value) {
      pdf.addPage();
      await processSection(pdf, PQStatusisInfo.value, 'PQ Status', canvasOptions, false);
      pageCount++;
    }

    // 5. PQ Detail - 전체 섹션을 캡처하되 차트는 숨기고 헤더만
    if (PQDetail.value && mode.value && ispqValid.value && pqdata.value.items.length > 0) {
      pdf.addPage();
      
      // 먼저 차트들을 임시로 숨김
      const trendCharts = PQDetail.value.querySelectorAll('.col-span-12');
      trendCharts.forEach(chart => {
        chart.style.display = 'none';
      });
      
      // 헤더만 있는 상태로 캡처
      await processSection(pdf, PQDetail.value, 'PQ Detail Header', canvasOptions, false);
      
      // 차트들을 다시 표시
      trendCharts.forEach(chart => {
        chart.style.display = '';
      });
      
      // 각 트렌드 차트를 개별적으로 처리
      let chartY = 60; // 헤더 아래 시작 위치
      
      for (let i = 0; i < trendCharts.length; i++) {
        if (i % 2 === 0 && i > 0) {
          pdf.addPage();
          chartY = 10;
        }
        await processChartElement(pdf, trendCharts[i], `PQ Trend ${i + 1}`, canvasOptions, chartY);
        chartY += 140;
      }
      pageCount++;
    }

    // 나머지 섹션들은 동일...
    // 6. PQ Trends
    if (PQTrends.value) {
      pdf.addPage();
      await processSection(pdf, PQTrends.value, 'PQ Trends', canvasOptions, false);
      pageCount++;
    }

    // 7. PQ Table
    if (PQTable.value) {
      pdf.addPage();
      await processSection(pdf, PQTable.value, 'PQ Table', canvasOptions, false);
      pageCount++;
    }

    // 8. ITIC Curve
    if (PQChart.value) {
      pdf.addPage();
      await processSection(pdf, PQChart.value, 'ITIC Curve', canvasOptions, false);
      pageCount++;
    }

    // 9. Energy Charts
    if (EnergyHourly.value) {
      pdf.addPage();
      await processSection(pdf, EnergyHourly.value, 'Hourly Energy', canvasOptions, false);
      pageCount++;
    }

    if (EnergyDaily.value) {
      pdf.addPage();
      await processSection(pdf, EnergyDaily.value, 'Daily Energy', canvasOptions, false);
      pageCount++;
    }

    if (EnergyMonthly.value) {
      pdf.addPage();
      await processSection(pdf, EnergyMonthly.value, 'Monthly Energy', canvasOptions, false);
      pageCount++;
    }

    // 10. Power Load Trend (전력량 & 부하율 추이)
    if (PowerLoadTrend.value) {
      pdf.addPage();
      await processSection(pdf, PowerLoadTrend.value, 'Power Load Trend', canvasOptions, false);
      pageCount++;
    }

    // 11. Load Pattern
    if (LoadPatternAnalysis.value) {
      pdf.addPage();
      await processSection(pdf, LoadPatternAnalysis.value, 'Load Pattern', canvasOptions, false);
      pageCount++;
    }

    // PDF 저장
    const pdfBlob = pdf.output('blob');
    const fileSize = (pdfBlob.size / 1024 / 1024).toFixed(2);
    console.log(`PDF generated successfully. Size: ${fileSize} MB`);
    
    const fileName = `power-facility-report-${channelComputed.value}-${new Date().toISOString().split('T')[0]}.pdf`;
    pdf.save(fileName);
    
  } catch (error) {
    console.error('PDF generation failed:', error);
    alert('PDF 생성 중 오류가 발생했습니다. 다시 시도해주세요.');
  } finally {
    document.documentElement.classList.remove('generating-pdf');
    await new Promise(resolve => setTimeout(resolve, 100));
    pdf_hidden.value = false;
  }
};

// processSection 함수 수정 - 헤더만 캡처할 때 높이 제한
async function processSection(pdf, element, name, options, addPage = false) {
  if (!element) return;
  
  console.log(`Processing ${name}...`);
  
  if (addPage) {
    pdf.addPage();
  }
  
  await nextTick();
  element.scrollIntoView({ behavior: 'instant' });
  await new Promise(resolve => setTimeout(resolve, 200));
  
  const canvas = await html2canvas(element, options);
  const imgData = canvas.toDataURL("image/jpeg", 0.85);
  
  const imgWidth = 190;
  const imgHeight = (canvas.height * imgWidth) / canvas.width;
  
  // 헤더만 캡처하는 경우 높이 제한
  let maxHeight = 270;
  if (name.includes('Header')) {
    maxHeight = 50; // 헤더는 작게
  }
  
  if (imgHeight > maxHeight) {
    const ratio = maxHeight / imgHeight;
    pdf.addImage(imgData, "JPEG", 10, 10, imgWidth * ratio, maxHeight);
  } else {
    pdf.addImage(imgData, "JPEG", 10, 10, imgWidth, imgHeight);
  }
  
  canvas.remove();
  await new Promise(resolve => setTimeout(resolve, 50));
}

// processChartElement 함수는 그대로
async function processChartElement(pdf, element, name, options, yPosition) {
  if (!element) return;
  
  console.log(`Processing ${name}...`);
  
  await nextTick();
  await new Promise(resolve => setTimeout(resolve, 100));
  
  const canvas = await html2canvas(element, options);
  const imgData = canvas.toDataURL("image/jpeg", 0.85);
  
  const imgWidth = 190;
  const imgHeight = (canvas.height * imgWidth) / canvas.width;
  const maxHeight = 130;
  
  pdf.addImage(imgData, "JPEG", 10, yPosition, imgWidth, Math.min(imgHeight, maxHeight));
  
  canvas.remove();
  await new Promise(resolve => setTimeout(resolve, 50));
}

    return {
      equipImageSrc,
      datalist,
      tbdata,
      assetdata,
      asset,
      rawdata,
      fetchData,
      fetchEn,
      fetchRawData,
      linechartData,
      hourlyChartData,
      dailyChartData,
      monthlyChartData,
      channelInfo,
      PQTable,
      PQChart,
      DiagnosisInfo,
      DiagnosisDetail,
      EnergyHourly,
      EnergyDaily,
      EnergyMonthly,
      PowerLoadTrend,  // 전력량 & 부하율 추이 추가
      LoadPatternAnalysis,
      generatePDF,
      pdf_hidden,
      channelComputed,
      setupMenu,
      channelStatus,
      mode,
      data,
      pqTableData,  // 변경된 변수명
      makeKey,
      getComp,
      getComplianceClass,
      t,
      locale,
      diagData,
      ischartDataValid,
      PQStatusisInfo,
      pqdata,
      ispqValid,
      fetchPQData,
      PQDetail,
      // Power Quality Trends
      PQTrends,
      timeseriesChart,
      timeseriesData,
      averageVUF,
      averageTHDV,
      averageTHDI,
      devLocation,
      // Power Load Trend
      dualAxisChart,
      averageLoadRate,
      maxLoadRate,
      overloadCount,
      // Load Pattern Analysis
      loadPatternChart,
      loadPatternData,
      apiHeatmapData,
      apiHeatmapStats,
      lightLoadPercentage,
      mediumLoadPercentage,
      highLoadPercentage,
      overLoadPercentage,
      // 추가된 함수들
      setDiagnosisParamIds,
      setPQParamIds,
      setChartData,
      formatToISOString,
      loadEnergyData,  // 에너지 데이터 로드 함수 추가
      loadHeatmapData, // 히트맵 데이터 로드 함수 추가
      t,
      assettypes,
    }
  }
}
</script>

<style scoped>
/* 호버 효과 */
.group:hover .group-hover\:text-blue-600 {
  color: #2563eb;
}

/* 트랜지션 */
.transition-transform {
  transition-property: transform;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 200ms;
}

.transition-colors {
  transition-property: color, background-color, border-color;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

/* 차트 영역 */
.dual-axis-chart {
  height: 320px;
  width: 100%;
}

.heatmap-chart {
  height: 320px;
  width: 100%;
}

/* PDF 인쇄 스타일 */
@media print {
  .bg-gradient-to-br {
    background: white !important;
  }
  
  .shadow-xl {
    box-shadow: none !important;
  }
  
  .rounded-2xl {
    border-radius: 8px !important;
  }
  
  .p-8 {
    padding: 1rem !important;
  }
  
  .text-3xl {
    font-size: 1.5rem !important;
  }
}

/* 그라데이션 텍스트 */
.bg-clip-text {
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 반응형 */
@media (max-width: 1024px) {
  .time-selector {
    @apply text-xs;
  }
}

@media (max-width: 640px) {
  .grid.grid-cols-3 {
    @apply grid-cols-1;
  }
}
</style>