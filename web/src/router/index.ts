import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SitesListView from '../views/SitesListView.vue'
import SiteDashboardView from '../views/SiteDashboardView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/sites',
      name: 'sitesList',
      component: SitesListView
    },
    {
      path: '/sites/:id',
      name: 'siteDashboard',
      component: SiteDashboardView
    }
  ]
})


export default router
