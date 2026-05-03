import { createRouter, createWebHashHistory } from 'vue-router'
const ListMocks = () => import('../components/ListMocks.vue')
const MockDetails = () => import('../components/MockDetails.vue')
const Help = () => import('../components/Help.vue')
const Logs = () => import('../components/Logs.vue')

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'ListMocks',
      component: ListMocks,
    },
    {
      path: '/details/:theKey',
      name: 'MockDetails',
      component: MockDetails,
      props: true
    },
    {
      path: '/new',
      name: 'NewMock',
      component: MockDetails,
    },
    {
      path: '/help',
      name: 'Help',
      component: Help,
    },
    {
      path: '/logs',
      name: 'Logs',
      component: Logs,
    },
  ],
})

export default router
