import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { useSettingStore } from '@/stores/setting'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/install',
      name: 'install',
      component: () => import('@/views/InstallView.vue'),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/AboutView.vue'),
    },
    {
      path: '',
      component: HomeView,
      children: [
        {
          path: '',
          name: 'root',
          component: () => import('@/views/RootView.vue'),
        },
        {
          path: 'issues',
          name: 'issue-list',
          component: () => import('@/views/IssueListView.vue'),
        },
        {
          path: 'issues/:issueId',
          name: 'issue-details',
          component: () => import('@/views/IssueDetailsView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, from, next) => {
  const settingStore = useSettingStore()
  if (!settingStore.isInstalled && to.name !== 'install') {
    return next({ name: 'install' })
  }

  return next()
})

export default router
