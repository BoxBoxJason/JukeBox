import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw, RouteMeta } from 'vue-router';
import HomeView from '../views/HomeView.vue';

// Define a type for the meta fields you use
interface CustomRouteMeta extends RouteMeta {
  title?: string;
}

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: {
      title: 'JukeBox'
    } as CustomRouteMeta
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('../views/AboutView.vue'),
    meta: {
      title: 'About | JukeBox'
    } as CustomRouteMeta
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
});

// Navigation guard to dynamically update the document title (has to be done like that, single page does not reload the document)
router.beforeEach((to, from, next) => {
  const metaTitle = (to.meta as CustomRouteMeta).title;
  if (metaTitle) {
    document.title = metaTitle;
  } else {
    document.title = 'JukeBox'; // Fallback title
  }
  next();
});

export default router;
