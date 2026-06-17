import type { RouteRecordRaw } from 'vue-router'

export const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    component: () => import('@/layouts/AuthLayout.vue'),
    meta: { requiresAuth: false },
    children: [
      {
        path: '',
        name: 'Login',
        component: () => import('@/features/auth/pages/LoginPage.vue'),
      },
    ],
  },
  {
    path: '/setup',
    component: () => import('@/layouts/AuthLayout.vue'),
    meta: { requiresAuth: false },
    children: [
      {
        path: '',
        name: 'Setup',
        component: () => import('@/features/auth/pages/SetupPage.vue'),
      },
    ],
  },
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/features/dashboard/pages/DashboardPage.vue'),
      },
      {
        path: 'users',
        name: 'UsersList',
        component: () => import('@/features/users/pages/UsersListPage.vue'),
        meta: { requiresAuth: true, requiresAdmin: true },
      },
      {
        path: 'questions',
        name: 'QuestionsList',
        component: () => import('@/features/questions/pages/QuestionsListPage.vue'),
        meta: { requiresNotAdmin: true },
      },
      {
        path: 'questions/:id',
        name: 'QuestionDetail',
        component: () => import('@/features/questions/pages/QuestionDetailPage.vue'),
        meta: { requiresNotAdmin: true },
      },
      {
        path: 'topics',
        name: 'TopicsList',
        component: () => import('@/features/topics/pages/TopicsListPage.vue'),
      },
      {
        path: 'sessions',
        name: 'SessionsList',
        component: () => import('@/features/sessions/pages/SessionsListPage.vue'),
        meta: { requiresNotAdmin: true },
      },
      {
        path: 'sessions/:id/study',
        name: 'SessionStudy',
        component: () => import('@/features/sessions/pages/SessionStudyPage.vue'),
        meta: { requiresNotAdmin: true },
      },
      {
        path: 'progress',
        name: 'Progress',
        component: () => import('@/features/progress/pages/ProgressPage.vue'),
        meta: { requiresNotAdmin: true },
      },
    ],
  },
]
