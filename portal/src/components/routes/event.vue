<template>
  <div class="event">
    <div class="event-header">
      <search-box />
    </div>
    <div class="event-result">
      <div class="event-result-chart">
        <div class="event-result-chart-wrap">
          <BarChart :chartData="testData2" :height="150" :options="options" />
        </div>
      </div>
      <div class="event-result-table">
        <div class="table-titles d-flex">
          <div class="title-item event-title-id">ID</div>
          <div class="title-item event-title-localtime">Local Time</div>
          <div class="title-item event-title-source">Source</div>
          <div class="title-item event-title-utctime">UTC Time</div>
          <div class="title-item event-title-data">Data</div>
        </div>
        <el-scrollbar @scroll="handleScroll" height="260px">
          <el-table
            :highlight-current-row="false"
            :data="tableData"
            style="width: 100%"
            v-loading="loading"
            ref="tableRef"
          >
            <el-table-column width="210">
              <template #default="scope">
                <span>{{ scope.row.id }}</span>
              </template>
            </el-table-column>

            <el-table-column width="180">
              <template #default="scope">
                <span>{{
                  moment(scope.row.localTime).format("DD/MM/YYYY HH:mm:ss")
                }}</span>
              </template>
            </el-table-column>
            <el-table-column width="120">
              <template #default="scope">
                <span>{{ scope.row.source }}</span>
              </template>
            </el-table-column>
            <el-table-column width="180">
              <template #default="scope">
                <span>{{
                  moment(scope.row.utcTime).format("DD/MM/YYYY HH:mm:ss")
                }}</span>
              </template>
            </el-table-column>
            <el-table-column>
              <template #default="scope">
                <span
                  class="data"
                  @click="openDetailDialog(scope.row.data, scope.row.id)"
                  >Click to see detail</span
                >
              </template>
            </el-table-column>
          </el-table>
        </el-scrollbar>
      </div>
    </div>
  </div>
  <detail-data ref="detailRef" />
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";
import moment from "moment";
import SearchBox from "../common/search-box.vue";
import { IEvent } from "../../common/model";
// import store from "../../common/store";
import { demoEvent } from "../demo-data";
import DetailData from "../common/detail-data.vue";
import { ElScrollbar, ElTable } from "element-plus";
import { DoughnutChart, BarChart } from "vue-chart-3";
import { Chart, registerables } from "chart.js";

Chart.register(...registerables);
Chart.defaults.color = "#fff";
Chart.defaults.borderColor = "#535353";

const detailRef = ref();
const loading = ref(false);
const isMounted = ref(false);
const tableData = ref<IEvent[]>(demoEvent.slice(0, 8));
const tableRef = ref();
const loadMore = () => {
  console.log("load");
  loading.value = true;
  if (tableData.value.length > demoEvent.length) {
    return;
  }
  const num = tableData.value.length;
  const arr = demoEvent.slice(num, num + 5);

  tableData.value = tableData.value.concat(arr);
  setTimeout(() => {
    loading.value = false;
  }, 1000);
};

const testData2 = {
  labels: ["Paris", "Nîmes", "Toulon", "Perpignan", "Autre"],
  datasets: [
    {
      label: "abc",
      data: [30, 40, 60, 70, 5],
      backgroundColor: ["#77CEFF"],
    },
  ],
};

const options = ref({
  responsive: true,
  plugins: {
    legend: {
      position: "top",
    },
  },
});

const handleScroll = ({ scrollTop }) => {
  if (Math.abs(scrollTop - tableRef.value.$el.clientHeight + 260) <= 1) {
    loadMore();
  }
};

const openDetailDialog = (val: object, id: any) => {
  if (detailRef?.value.showDetailData) {
    detailRef?.value.showDetailData(val, id);
  }
};

onMounted(() => {
  nextTick(() => (isMounted.value = true));
});
</script>

<style scoped>
@import url("../css//el-table.css");
.event {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30px;
}

.event-header {
  width: 100%;
  background-color: #363c41;
  border-radius: 4px;
}

/* RESULT  */
.event-result {
  width: 100%;
  height: calc(100% - 102px);
  display: flex;
  flex-direction: column;
  gap: 20px;
}
/* RESUTL CHART  */
.event-result-chart {
  height: 150px;
  background-color: #363c41;
  padding-inline: 30px;
  border-radius: 4px;
  width: 100%;
  display: flex;
  justify-content: center;
}
.event-result-chart-wrap {
  width: 50%;
}
/* RESULT TABLE  */
.event-result-table {
  background-color: #363c41;
  padding: 10px;
  border-radius: 4px;
}

:deep(.el-input:not(.el-input--medium) .el-input__wrapper) {
  box-shadow: unset;
  border-radius: unset;
}
</style>
