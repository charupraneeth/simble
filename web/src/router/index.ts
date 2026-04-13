import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SitesListView from '../views/SitesListView.vue'
import SiteDashboardView from '../views/SiteDashboardView.vue'
import OnboardingView from '../views/OnboardingView.vue'
import DemoView from '../views/DemoView.vue'
import { useAuth } from '../composables/useAuth'

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
      component: SitesListView,
      meta: { requiresAuth: true }
    },
    {
      path: '/onboarding',
      name: 'onboarding',
      component: OnboardingView,
      meta: { requiresAuth: true }
    },
    {
      path: '/sites/:id',
      name: 'siteDashboard',
      component: SiteDashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/demo',
      name: 'demo',
      component: DemoView
    }
  ]
})

router.beforeEach(async (to) => {
  const { login, isLoggedIn } = useAuth()

  // Await the login check to guarantee we know their identity before routing
  await login()

  // Guard 1: Unauthenticated users visiting a secure route
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    // Optionally: trigger a toast here if they were kicked out mid-session 
    // by tracking a previouslyLoggedIn variable.
    return { name: 'home' } // redirect back to landing
  }

  // Guard 2: Authenticated users trying to view the public landing page
  if (to.name === 'home' && isLoggedIn.value) {
    return { name: 'sitesList' } // kick them cleanly to their dashboard
  }
})

export default router
