<template>
  <div class="css-bar-chart">
    <div class="chart-bars">
      <div 
        v-for="(item, index) in chartItems" 
        :key="index"
        class="bar-item"
      >
        <!-- 바 컨테이너 -->
        <div class="bar-container">
          <div 
            class="bar-fill"
            :class="item.colorClass"
            :style="{ height: item.height + '%' }"
          >
            <!-- 값 라벨 -->
            <div class="bar-label">{{ item.value }}%</div>
          </div>
        </div>
        
        <!-- 하단 라벨 -->
        <div class="bar-name">{{ item.label }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'CSSBarChart',
  props: {
    data: {
      type: Object,
      default: () => ({})
    },
    height: {
      type: [String, Number],
      default: 120
    },
    dataKeys: {
      type: Array,
      default: () => ['thdu total', 'thdi total', 'tddi total']
    },
    labels: {
      type: Array,
      default: () => ['THD-U', 'THD-I', 'TDD-I']
    },
    colors: {
      type: Array,
      default: () => ['purple', 'blue', 'green']
    }
  },
  emits: ['data-change'],
  setup(props, { emit }) {
    
    const chartItems = computed(() => {
      const values = props.dataKeys.map(key => parseFloat(props.data[key] || 0))
      const maxValue = Math.max(...values, 10) // 최소 10으로 설정
      
      const items = props.labels.map((label, index) => {
        const value = values[index].toFixed(1)
        const height = (values[index] / maxValue) * 100
        
        return {
          label,
          value,
          height: Math.max(height, 5), // 최소 5% 높이
          colorClass: `bar-${props.colors[index]}`
        }
      })
      
      // 데이터 변경 이벤트 발생
      emit('data-change', { data: items })
      
      return items
    })

    return {
      chartItems
    }
  }
}
</script>

<style scoped>
.css-bar-chart {
  @apply w-full p-4 pt-8;
}

.chart-bars {
  @apply flex justify-around items-end gap-4 h-24;
  @apply relative;
}

.bar-item {
  @apply flex-1 flex flex-col items-center;
  @apply max-w-[60px];
}

.bar-container {
  @apply relative w-full h-16 flex items-end justify-center;
}

.bar-fill {
  @apply w-full rounded-t transition-all duration-700 ease-out;
  @apply relative flex items-start justify-center;
  @apply min-h-[8px];
}

.bar-label {
  @apply text-xs font-bold;
  @apply absolute bottom-full left-1/2 transform -translate-x-1/2;
  @apply text-center whitespace-nowrap;
  @apply text-gray-800 dark:text-gray-200;
  @apply mb-1;
}

.bar-name {
  @apply text-xs font-semibold text-gray-700 dark:text-gray-300;
  @apply mt-2 text-center;
}

/* 바 색상 */
.bar-purple {
  @apply bg-gradient-to-t from-purple-500 to-purple-400;
}

.bar-blue {
  @apply bg-gradient-to-t from-blue-500 to-blue-400;
}

.bar-green {
  @apply bg-gradient-to-t from-green-500 to-green-400;
}

/* 반응형 */
@media (max-width: 640px) {
  .chart-bars {
    @apply gap-2 h-20;
  }
  
  .bar-container {
    @apply h-12;
  }
  
  .bar-label {
    @apply text-xs;
  }
  
  .bar-name {
    @apply text-xs;
  }
}
</style>