import { createRouter, createWebHashHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import Home from "@/views/Home.vue";
const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    redirect: "/home",
  },
  {
    path: "/home",
    component: Home,
    meta: { title: "首页" }
  },
];
const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
