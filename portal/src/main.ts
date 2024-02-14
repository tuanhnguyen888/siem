import { createApp } from "vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import App from "./App.vue";
import "./style.css";
import "office-ui-fabric-vue/dist/index.css";
import ElTableInfiniteScroll from "el-table-infinite-scroll";

import router from "./common/router";

const app = createApp(App);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}
app.use(ElTableInfiniteScroll);
app.use(ElementPlus).use(router);
app.mount("#app");
