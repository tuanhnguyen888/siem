<script setup lang="ts">
import { computed } from "vue";
import store from "./common/store";
import OpeningPage from "./components/opening-page.vue";
import Homepage from "./components/homepage.vue";

const isLogin = computed(() => store.state.isLogin);
</script>

<template>
  <template v-if="isLogin">
    <homepage />
  </template>
  <template v-else>
    <opening-page />
  </template>
</template>
