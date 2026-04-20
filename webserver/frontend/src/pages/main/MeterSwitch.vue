<template>
  <Meter v-if="viewMode === 'classic'" :channel="channel" />
  <MeterDesign v-else :channel="channel" />
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'
import Meter from './Meter.vue'
import MeterDesign from './MeterDesign.vue'

export default {
  name: 'MeterSwitch',
  components: { Meter, MeterDesign },
  props: {
    channel: {
      type: String,
      default: 'Main',
    },
  },
  setup() {
    const authStore = useAuthStore()
    const viewMode = computed(() => authStore.getViewMode)
    return { viewMode }
  },
}
</script>
