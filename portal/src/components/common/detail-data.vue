<template>
  <el-dialog v-model="dialogVisible" :title="detailID" width="800px">
    <div class="detail-table">
      <div class="detail-entry d-flex" v-for="index in detailKeyArr.length">
        <div class="entry-item entry-key">{{ detailKeyArr[index - 1] }}</div>
        <div class="entry-item entry-value">
          {{ detailValueArr[index - 1] }}
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";

const dialogVisible = ref(false);
const detailKeyArr = ref<any>({});
const detailValueArr = ref<any>({});
const detailID = ref("");

const showDetailData = (val: object, id: any) => {
  dialogVisible.value = true;
  detailID.value = id;
  detailKeyArr.value = Object.keys(val);
  detailValueArr.value = Object.values(val);
};

defineExpose({
  showDetailData,
});
</script>

<style scoped>
.detail-entry {
  border-top: 1px solid black;
  border-left: 1px solid black;
  border-right: 1px solid black;
  color: black;
}
.detail-entry:last-child {
  border-bottom: 1px solid black;
}
.entry-item {
  padding: 5px 10px;
}
.entry-key {
  border-right: 1px solid black;
  width: 40%;
  background-color: #d8d6d6;
}
.entry-value {
  width: 60%;
  overflow: hidden;
}
</style>
