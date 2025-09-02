<template>
  <table class="table-fixed w-full dark:text-white">
    <!-- colgroup으로 컬럼 너비 비율 고정 -->
    <colgroup>
      <col style="width: 70%">
      <col style="width: 20%">
      <col style="width: 10%">
    </colgroup>
    
    <!-- Table header -->
    <thead
      class="text-xs uppercase text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-700 dark:bg-opacity-50 rounded-sm"
    >
      <tr>
        <th class="px-4 py-2">
          <div class="font-bold text-left">
            {{ t("trend.treetable.Item") }}
          </div>
        </th>
        <th class="px-4 py-2">
          <div class="font-bold text-left">
            {{ t("trend.treetable.Assembly") }}
          </div>
        </th>
        <th class="px-4 py-2">
          <div class="font-bold text-center">
            {{ t("trend.treetable.Select") }}
          </div>
        </th>
      </tr>
    </thead>
    <!-- Table body -->
    <tbody>
      <TreeRow
        v-for="item in data"
        :key="item.id"
        :item="item"
        :level="0"
        :mode="mode"
        :checked-ids="checkedIds"
        :checked-names="checkedNames"
        :expanded="false"
        @check-change="onCheckChange"
      />
    </tbody>
  </table>
</template>

<script>
import { ref } from "vue";
import TreeRow from "./TreeRowTrend.vue";
import { useI18n } from "vue-i18n";

export default {
  name: "Trend_treetable",
  components: {
    TreeRow,
  },
  props: {
    checkedIds: Array,
    checkedNames: Array,
    data: {
      type: Array,
      default: () => [], // ✅ required 제거
    },
    channel: {
      type: String,
      default: "",
    },
    mode: {
      type: String,
      default: "",
    },
  },
  setup() {
    const { t } = useI18n();
    return { t };
  },
  methods: {
    onCheckChange(payload) {
      this.$emit("check-change", payload);
    },
  },
  // setup(props){
  //   const items = ref(props.data);
  //   const mode = ref(props.mode);

  //   function onCheckChange(payload) {
  //     emit('check-change', payload) // 그대로 상위로 전달
  //   }
  //   return {
  //     items,
  //     mode,
  //     onCheckChange,
  //   }
  // }
};
</script>