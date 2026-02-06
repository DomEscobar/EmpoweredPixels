import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "home",
    component: () => import("../../pages/Home.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("../../pages/Login.vue"),
  },
  {
    path: "/register",
    name: "register",
    component: () => import("../../pages/Register.vue"),
  },
  {
    path: "/roster",
    name: "roster",
    component: () => import("../../pages/Roster.vue"),
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: () => import("../../pages/Dashboard.vue"),
  },
  {
    path: "/matches",
    name: "matches",
    component: () => import("../../pages/Matches.vue"),
  },
  {
    path: "/inventory",
    name: "inventory",
    component: () => import("../../pages/Inventory.vue"),
  },
  {
    path: "/leagues",
    name: "leagues",
    component: () => import("../../pages/Leagues.vue"),
  },
  {
    path: "/matches/:id",
    name: "match-viewer",
    component: () => import("../../pages/MatchViewer.vue"),
  },
  {
    path: "/shop",
    name: "shop",
    component: () => import("../../pages/Shop.vue"),
  },
  {
    path: "/attunement",
    name: "attunement",
    component: () => import("../../pages/Attunement.vue"),
  },
  {
    path: "/leaderboard",
    name: "leaderboard",
    component: () => import("../../pages/Leaderboard.vue"),
  },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
