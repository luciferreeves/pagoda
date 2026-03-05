import { lazy } from "solid-js";
import type { RouteDefinition } from "@solidjs/router";

import Home from "./pages/home";

export const routes: RouteDefinition[] = [
  { path: "/", component: Home },
  { path: "/login", component: lazy(() => import("./pages/login")) },
  { path: "/register", component: lazy(() => import("./pages/register")) },
  { path: "/account/verify", component: lazy(() => import("./pages/account/verify")) },
  { path: "/account/reactivate", component: lazy(() => import("./pages/account/reactivate")) },
  { path: "**", component: lazy(() => import("./errors/404")) },
];
