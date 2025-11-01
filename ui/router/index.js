import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '',
      component: () => import('@/views/HomeView.vue'),
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

export default router
