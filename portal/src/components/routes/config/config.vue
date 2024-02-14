<script setup lang="ts">
import { ref } from "vue";
import { demoLogSource } from "../../demo-data";
import {
  ElMessageBox,
  ElNotification,
  ElScrollbar,
  ElTable,
} from "element-plus";
import moment from "moment";
import { ILogSource } from "../../../common/model";
import DetailData from "../../common/detail-data.vue";
import ConfigForm from "./config-form.vue";

const tableRef = ref();
const detailRef = ref();
const loading = ref(false);
const configFormRef = ref();

const tableData = ref<ILogSource[]>(demoLogSource.slice(0, 7));

const loadMore = () => {
  loading.value = true;
  if (tableData.value.length > demoLogSource.length) {
    return;
  }
  const num = tableData.value.length;
  const arr = demoLogSource.slice(num, num + 5);

  tableData.value = tableData.value.concat(arr);
  setTimeout(() => {
    loading.value = false;
  }, 1000);
};

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

const openFormDialog = (val?: object) => {
  if (configFormRef?.value.showFormDialog) {
    configFormRef?.value.showFormDialog(val);
  }
};

const deleteLogSource = async () => {
  ElNotification({
    type: "success",
    title: "Successful",
    message: "Deleted",
    duration: 3000,
  });
};

const handleDelete = () => {
  ElMessageBox.confirm("Do you want to delete it?", "Warning", {
    confirmButtonText: "OK",
    cancelButtonText: "Cancel",
    type: "warning",
  })
    .then(async () => {
      await deleteLogSource();
    })
    .catch(() => {
      console.log("cancel");
    });
};
</script>

<template>
  <div class="config-wrap">
    <div class="config-addnew mb-3" @click="openFormDialog({})">
      <i class="ms-Icon ms-font ms-Icon--Add cursor-pointer mr-2"></i>
      <span>Add New</span>
    </div>
    <div class="config-table">
      <div class="table-titles d-flex">
        <div class="title-item config-title-inputname">Input Name</div>
        <div class="title-item config-title-protocol">Protocol</div>
        <div class="title-item config-title-logtype">Log Type</div>
        <div class="title-item config-title-created">Created At</div>
        <div class="title-item config-title-updated">Updated At</div>
        <div class="title-item config-title-properties">Properties</div>
        <div class="title-item config-title-actions"></div>
      </div>
      <el-scrollbar @scroll="handleScroll" max-height="260px">
        <el-table
          :data="tableData"
          style="width: 100%"
          v-loading="loading"
          ref="tableRef"
        >
          <!-- INPUT NAME -->
          <el-table-column width="180">
            <template #default="scope">
              <span>{{ scope.row.inputName }}</span>
            </template>
          </el-table-column>
          <!-- PROTOCOL -->
          <el-table-column width="150" show-overflow-tooltip>
            <template #default="scope">
              <span>{{ scope.row.protocol }}</span>
            </template>
          </el-table-column>
          <!-- LOG TYPE -->
          <el-table-column width="150">
            <template #default="scope">
              <span>{{ scope.row.properties.log_type }}</span>
            </template>
          </el-table-column>
          <!-- CREATED -->
          <el-table-column width="160" show-overflow-tooltip>
            <template #default="scope">
              <span>{{
                moment(scope.row.createddAt).format("DD/MM/YYYY HH:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <!-- UPDATED -->
          <el-table-column width="160" show-overflow-tooltip>
            <template #default="scope">
              <span>{{
                moment(scope.row.updatedAt).format("DD/MM/YYYY HH:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <!-- PROPERTIES  -->
          <el-table-column>
            <template #default="scope">
              <span
                class="data"
                @click="openDetailDialog(scope.row.properties, scope.row.id)"
                >Click to see detail</span
              >
            </template>
          </el-table-column>

          <!-- DELETE/EDIT  -->
          <el-table-column width="100">
            <template #default="scope">
              <i
                class="ms-Icon ms-font-xl mr-3 ms-Icon--Delete cursor-pointer"
                @click="handleDelete"
              >
              </i>
              <i
                @click="openFormDialog(scope.row)"
                class="ms-Icon ms-font-xl ms-Icon--Edit cursor-pointer"
              ></i>
            </template>
          </el-table-column>
        </el-table>
      </el-scrollbar>
    </div>
    <config-form ref="configFormRef" />
    <detail-data ref="detailRef" />
  </div>
</template>

<style scoped lang="scss">
@import url("../../css/el-table.css");
.config-wrap {
  height: 100%;
  width: 100%;
  background-color: #363c41;
  border-radius: 4px;
  padding: 10px 10px 20px 10px;
  color: white;
}

.config-addnew {
  display: flex;
  align-items: center;
  cursor: pointer;
}
</style>
<style lang="scss">
.config-wrap .el-dialog {
  max-height: 510px;
  .el-dialog__header {
    font-weight: bold;
    margin-right: 0;
    border-bottom: 1px solid black;
  }
}
</style>
