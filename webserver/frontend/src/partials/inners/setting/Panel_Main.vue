<template>
  <div class="grow dark:text-white">
    <div class="p-6 space-y-6">
      <section v-if="getInputDict().Enable">

        <!-- 1행: CT / PT / Demand / Opt -->
        <div class="grid grid-cols-12 gap-5">
          <CTCard class="col-span-12 md:col-span-6 xl:col-span-3" />
          <PTCard class="col-span-12 md:col-span-6 xl:col-span-3" />
          <DemandCard class="col-span-12 md:col-span-6 xl:col-span-3" />
          <OptCard class="col-span-12 md:col-span-6 xl:col-span-3" />
        </div>

        <!-- 2행: Event / Trend -->
        <div class="grid grid-cols-12 gap-5 mt-5">
          <EventCard2 class="col-span-12 md:col-span-8 xl:col-span-8" />
          <TrendPanelCard class="col-span-12 md:col-span-4 xl:col-span-4" />
        </div>

        <!-- 3행: Alarm (full width) -->
        <div class="mt-5">
          <AlarmCard :parameterOptions="parameterOptions" />
        </div>

      </section>
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted, inject, computed, provide, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useInputDict } from '@/composables/useInputDict'

import CTCard from './CTCard.vue'
import PTCard from './PTCard.vue'
import DemandCard from './DemandCard.vue'
import AlarmCard from './AlarmCard.vue'
import EventCard2 from './EventCard2.vue'
import TrendPanelCard from './TrendPanelCard.vue'
import OptCard from './OptCard.vue'

export default {
  name: 'PanelMain',
  props: ['channel'],
  components: {
    CTCard,
    PTCard,
    DemandCard,
    AlarmCard,
    EventCard2,
    TrendPanelCard,
    OptCard,
  },
  setup(props) {
    const { t } = useI18n()
    const General_inputDict = inject('inputDict', ref({}))
    const inputDict_main = inject('channel_main', ref({ Enable: false }))
    const inputDict_sub = inject('channel_sub', ref({ Enable: false }))
    const saveAllSettings = inject('saveAllSettings', null)

    const componentKey = ref(0)
    const channel = ref(props.channel)

    const { parameterOptions, selectedTrendSetup } = useInputDict()

    // Current channel inputDict
    const getInputDict = () => {
      return props.channel === 'Main'
        ? inputDict_main.value
        : inputDict_sub.value
    }

    // Field update functions
    const updateField = (field, value) => {
      const currentDict = getInputDict()
      currentDict[field] = value
      componentKey.value++
    }

    const updateNestedField = (parent, field, value) => {
      const currentDict = getInputDict()
      if (!currentDict[parent]) {
        currentDict[parent] = {}
      }
      currentDict[parent][field] = value
      componentKey.value++
    }

    const updateArrayField = (parent, field, index, value) => {
      const currentDict = getInputDict()
      if (!currentDict[parent]) {
        currentDict[parent] = {}
      }
      if (!currentDict[parent][field]) {
        currentDict[parent][field] = []
      }
      currentDict[parent][field][index] = value
      componentKey.value++
    }

    watch(
      () => props.channel,
      (newChannel) => {
        channel.value = newChannel
        componentKey.value++
      },
    )

    const syncTrendSetup = () => {
      const currentDict = getInputDict()
      if (currentDict?.trendInfo) {
        const trendInfo = currentDict.trendInfo
        const periodValue = typeof trendInfo.period === 'number' ? trendInfo.period : 5
        const TREND_PARAM_COUNT = 10
        let params = Array.isArray(trendInfo.params) ? [...trendInfo.params] : [...selectedTrendSetup.value.params]
        while (params.length < TREND_PARAM_COUNT) {
          params.push("None")
        }
        params = params.slice(0, TREND_PARAM_COUNT)
        selectedTrendSetup.value = {
          period: periodValue,
          params,
        }
      }
    }

    watch(
      () => getInputDict()?.trendInfo,
      () => {
        syncTrendSetup()
      },
      { deep: true }
    )

    // Provide for child components
    provide('channel_inputDict', computed(() => getInputDict()))
    provide('updateNestedField', updateNestedField)
    provide('updateField', updateField)
    provide('updateArrayField', updateArrayField)
    provide('selectedTrendSetup', selectedTrendSetup)
    provide('parameterOptions', parameterOptions)
    provide('getInputDict', getInputDict)

    return {
      channel,
      getInputDict,
      updateField,
      updateNestedField,
      updateArrayField,
      parameterOptions,
      selectedTrendSetup,
      componentKey,
      t,
    }
  },
}
</script>
