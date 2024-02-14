<template>
  <div class="table">
    <div class="table-titles d-flex">
      <div class="title-item title-id">ID</div>
      <div class="title-item title-localtime">Local Time</div>
      <div class="title-item title-source">Source</div>
      <div class="title-item title-utctime">UTC Time</div>
      <div class="title-item title-data">Data</div>
    </div>
    <div class="table-content">
      <el-table
        :data="props.tableData"
        style="width: 100%"
        height="260px"
        v-loading="props.loading"
        v-el-table-infinite-scroll="loadMore"
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
    </div>
  </div>
  <detail-data ref="detailRef" />
</template>

<script setup lang="ts">
import { ref } from "vue";
import moment from "moment";

interface IProps {
  tableData: Array<any>;
  loading: boolean;
}
const detailRef = ref();
const props = defineProps<IProps>();
const emits = defineEmits(["loadMore"]);
const openDetailDialog = (val: object, id: any) => {
  if (detailRef?.value.showDetailData) {
    detailRef?.value.showDetailData(val, id);
  }
};

const loadMore = () => {
  emits("loadMore");
};
</script>

<!-- <style scoped>
.title-item {
  padding-left: 12px;
  height: 40px;
  display: flex;
  align-items: center;
  background-color: #d5d4d4;
  border-left: 1px solid black;
  border-top: 1px solid black;
  border-bottom: 1px solid black;
  font-weight: 600;
}

.title-id {
  width: 209.5px;
}
.title-localtime {
  width: 180px;
}
.title-source {
  width: 120px;
}
.title-utctime {
  width: 180px;
}
.title-data {
  width: calc(100% - 689.5px);
  border-right: 1px solid black;
}

.data {
  overflow: hidden;
  cursor: pointer;
  color: #0022ff;
  text-decoration: underline;
}
.data:hover {
  color: #ff1e00;
}
:deep(.el-table tr) {
  background-color: #eae9e9;
}

:deep(.el-table__cell) {
  /* border-bottom: solid 1px black !important; */
  border-top: solid 1px black !important;
  /* border-left: solid 1px black !important; */
  color: black;
}
:deep(.el-table__header) {
  display: none;
}

:deep(.el-table:not(.el-table--border) .el-table__cell) {
  border-right: solid 1px black;
}
:deep(.el-table__row .el-table__cell:first-child) {
  border-left: solid 1px black !important;
}
</style> -->
