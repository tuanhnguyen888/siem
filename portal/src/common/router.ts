import { RouteRecordRaw } from "vue-router";
import { createRouter, createWebHistory } from "vue-router";

import Event from "../components/routes/event.vue";
import Alert from "../components/routes/alert.vue";
import Config from "../components/routes/config/config.vue";
import UserManagement from "../components/routes/user-management.vue";
import Profile from "../components/routes/profile.vue";
import OpeningPage from "../components/opening-page.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "opening",
    component: OpeningPage,
  },
  {
    path: "/event",
    name: "event",
    component: Event,
  },
  {
    path: "/alert",
    name: "alert",
    component: Alert,
  },
  {
    path: "/config",
    name: "config",
    component: Config,
  },
  {
    path: "/usermanagement",
    name: "usermanagement",
    component: UserManagement,
  },
  {
    path: "/profile",
    name: "profile",
    component: Profile,
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
