<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { ICurrentPage } from "../../common/model";
import store from "../../common/store";
import { ElDatePicker } from "element-plus";

const currentPage = computed(() => store.state.currentPage);
const isDisableQuickTime = ref(false);
const isShowDatePicker = ref(false);
const keySearch = ref<string>("");
const timeRangePicker = ref("");
const timeRangeArr = [
  {
    value: 30,
    label: "30 min",
  },
  {
    value: 60,
    label: "60 min",
  },
  {
    value: 240,
    label: "4 hours",
  },
  {
    value: 720,
    label: "12 hours",
  },
  {
    value: 1440,
    label: "1 day",
  },
];
const timeRangeSelected = ref(30);

const showDatePicker = (val) => {
  isShowDatePicker.value = val;
};

const clickTimePicker = () => {
  showDatePicker(true);
  isDisableQuickTime.value = true;
};

const submitTimepicker = () => {
  showDatePicker(false);
  if (!timeRangePicker.value) {
    isDisableQuickTime.value = false;
  }
};
const cancelTimpicker = () => {
  showDatePicker(false);
  timeRangePicker.value = "";
  isDisableQuickTime.value = false;
};

watch(
  () => currentPage.value,
  () => {
    isDisableQuickTime.value = false;
    timeRangePicker.value = "";
    keySearch.value = "";
  }
);
</script>

<template>
  <div class="searchbox">
    <div class="searchbox-input">
      <el-input v-model="keySearch" placeholder="Enter key search" />

      <el-select
        size="medium"
        v-model="timeRangeSelected"
        class="m-2"
        placeholder="Select"
        :disabled="isDisableQuickTime"
      >
        <el-option
          v-for="item in timeRangeArr"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
    </div>
    <el-popover
      placement="bottom"
      trigger="click"
      :visible="isShowDatePicker"
      :width="425"
    >
      <template #reference>
        <div @click="clickTimePicker" class="btn search-btn">
          <i class="ms-Icon ms-font-xl ms-Icon--DateTime cursor-pointer"></i>
        </div>
      </template>
      <template #default>
        <div class="time-table">
          <el-date-picker
            v-model="timeRangePicker"
            type="datetimerange"
            start-placeholder="Start date"
            end-placeholder="End date"
            format="YYYY-MM-DD HH:mm:ss"
            date-format="YYYY/MM/DD"
            time-format="HH:mm"
          />
          <div class="time-table-btns d-flex justify-content-end mt-2">
            <div
              class="time-table-btn mr-2"
              style="color: rgb(255, 88, 88)"
              @click="cancelTimpicker"
            >
              Clear/Cancel
            </div>
            <div
              class="time-table-btn"
              style="background-color: rgb(255, 88, 88); color: white"
              @click="submitTimepicker"
            >
              OK
            </div>
          </div>
        </div>
      </template>
    </el-popover>

    <div class="btn search-btn">
      <i class="ms-Icon ms-font-xl ms-Icon--Search cursor-pointer"></i>
    </div>
    <div class="btn download-btn">
      <i class="ms-Icon ms-font-xl ms-Icon--Download cursor-pointer"></i>
    </div>
    <div v-if="currentPage === ICurrentPage.Alert" class="btn flag-btn">
      <i class="ms-Icon ms-font-xl ms-Icon--Warning cursor-pointer"></i>
    </div>
  </div>
</template>

<style scoped>
.searchbox {
  width: 100%;
  display: flex;
  justify-content: center;
  /* background-color: #e2e2e2; */
  padding-block: 8px;
}
.searchbox-input {
  display: flex;
  gap: 5px;
  align-items: center;
  color: black !important;
  height: 42px;
  max-width: 800px;
  width: 50%;
  border: solid rgb(152, 152, 152) 1.5px;
  padding-inline: 10px 10px;
  border-radius: 20px;
  background-color: #fff;
}

.btn {
  height: 42px;
  width: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: solid #d4d2d2 1.5px;
  background-color: white;
}
.btn:hover {
  border: solid #929292 1.5px;
}
.search-btn {
  margin-inline: 10px;
}

.flag-btn {
  margin-left: 10px;
}

.ms-Icon:hover {
  font-size: 25px;
  font-weight: 600;
}
.ms-Icon--Search {
  transform: scaleX(-1);
}
.time-table-btn {
  height: 30px;
  width: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  border: 1px solid black;
  cursor: pointer;
}
</style>
