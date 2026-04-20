<template>
  <div class="overflow-x-auto">
    <table class="table-auto w-full dark:text-white">
      <!-- Table header -->
      <thead class="text-xs uppercase text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-700 dark:bg-opacity-50 rounded-sm">
        <tr>
          <th class="p-2">
            <div class="font-bold text-left">Harmonics</div>
          </th>
          <th class="p-2">
            <div class="font-bold text-center">L1</div>
          </th>
          <th class="p-2">
            <div class="font-bold text-center">L2</div>
          </th>
          <th class="p-2">
            <div class="font-bold text-center">L3</div>
          </th>
        </tr>
      </thead>
      <!-- Table body -->
      <tbody class="text-sm font-medium divide-y divide-gray-100 dark:divide-gray-700/60">
        <template v-for="(row, index) in convertedData" :key="index">
          <tr v-if="index > 1">
            <td class="p-2 text-left">{{ index }}</td>
            <td class="p-2 text-center">{{ row.L1 }}</td>
            <td class="p-2 text-center">{{ row.L2 }}</td>
            <td class="p-2 text-center">{{ row.L3 }}</td>
          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template>

<script>
import { computed } from 'vue'

export default {
  name: 'PowerQ_Table',
  props: {
    data: {
      type: Object,
      default: () => ({})
    },
    option: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const harmonicsMap = {
      "phaseVoltage": ["U1", "U2", "U3"],
      "current": ["I1", "I2", "I3"]
    };

    const convertedData = computed(() => {
      const [k1, k2, k3] = harmonicsMap[props.option] || harmonicsMap.phaseVoltage;
      const L1 = props.data[k1] || [];
      const L2 = props.data[k2] || [];
      const L3 = props.data[k3] || [];
      const rowCount = Math.max(L1.length, L2.length, L3.length);

      const rows = [];
      for (let i = 0; i < rowCount; i++) {
        rows.push({
          L1: L1[i] ?? '-',
          L2: L2[i] ?? '-',
          L3: L3[i] ?? '-'
        });
      }
      return rows;
    });

    return {
      convertedData,
    }
  }
}
</script>
