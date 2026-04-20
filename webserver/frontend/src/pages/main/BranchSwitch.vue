<template>
  <Branch v-if="viewMode === 'classic'" :id="id" />
  <BranchDesign v-else :id="id" />
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'
import Branch from './Branch.vue'
import BranchDesign from './BranchDesign.vue'

export default {
  name: 'BranchSwitch',
  components: { Branch, BranchDesign },
  props: {
    id: { type: String, required: true },
  },
  setup() {
    const authStore = useAuthStore()
    const viewMode = computed(() => authStore.getViewMode)
    return { viewMode }
  },
}
</script>
