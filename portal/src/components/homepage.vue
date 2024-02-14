<template>
  <div class="homepage">
    <div class="homepage-header">
      <div
        class="homepage-header-top d-flex justify-content-between align-items-center"
      >
        <div class="homepage-header-left d-flex align-items-center">
          <img src="../assets/img/icon-white.png" alt="" />
          <span class="project-name">Project Name</span>
        </div>
        <el-popover width="100">
          <template #reference>
            <div class="homepage-header-right d-flex align-items-center">
              <i class="ms-Icon ms-Icon--Contact"></i>
              <span>Login User</span>
            </div>
          </template>
          <template #default>
            <div class="user-items">
              <div class="user-item">Profile</div>
              <div class="user-item">Log out</div>
            </div>
          </template>
        </el-popover>
      </div>
      <div class="homepage-header-options d-flex">
        <div
          v-for="option in navOptions"
          class="option"
          :class="isSelected(option) && 'option-checked'"
          @click="routeChange(option.id)"
        >
          <i :class="`ms-Icon ms-font-xl ms-Icon--${option.icon}`"></i>
          <span>{{ option.label }}</span>
        </div>
      </div>
    </div>

    <div class="homepage-content">
      <router-view></router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import router from "../common/router";
import store from "../common/store";
// console.log(router);

const currentPage = computed(() => store.state.currentPage);
const navOptions = [
  { id: "event", icon: "Event", label: "Event" },
  { id: "alert", icon: "Warning", label: "Alert" },
  { id: "config", icon: "Settings", label: "Config" },
  { id: "usermanagement", icon: "People", label: "User Management" },
  { id: "personalpage", icon: "", label: "" },
];

const isSelected = (item: any) => {
  return item.id === currentPage.value;
};

const routeChange = (item: string) => {
  router.push({ name: item });
  store.dispatch("getCurrentPage", item);
};
</script>

<style scoped lang="scss">
.user-item {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.user-item:hover {
  background-color: #3d444e;
  color: white;
}
.homepage-header {
  background-color: #181d21;

  color: white;
  gap: 3px;
  position: relative;
  z-index: 2;
  box-shadow: 0 3px 5px rgba(0, 0, 0, 0.29);
  .homepage-header-top {
    padding: 10px 30px 0;
    background-color: #181d21;
  }
  .homepage-header-left,
  .homepage-header-right {
    gap: 10px;
  }
  .homepage-header-left img {
    height: 60px;
  }
  .project-name {
    font-size: 25px;
    font-weight: 600;
  }
  .homepage-header-options {
    padding-top: 10px;
    padding-left: 30px;
    background-color: #3d444e;
    gap: 20px;
    border-radius: 4px;
  }
  .option {
    display: flex;
    align-items: center;
    // padding: 10px 20px;
    gap: 5px;
    cursor: pointer;
    border-bottom: 5px solid #3d444e;
  }
  .option:hover {
    border-bottom: 5px solid #499451;
  }
  .option-checked {
    border-bottom: 5px solid #499451;
  }
}
.homepage-content {
  background-color: #181d21;
  // min-height: calc((100% - 117px));
  padding: 10px 30px 0;
}
</style>
