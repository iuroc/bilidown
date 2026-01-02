import { createRouter, createWebHashHistory } from 'vue-router'
import ParseView from '@/views/parse/parse-view.vue'
import DownloadView from '@/views/download/download-view.vue'
import SettingView from '@/views/setting/setting-view.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'HomeView',
      redirect: { name: 'ParseView' }
    },
    {
      path: '/parse',
      name: 'ParseView',
      component: ParseView
    },
    {
      path: '/download',
      name: 'DownloadView',
      component: DownloadView
    },
    {
      path: '/setting',
      name: 'SettingView',
      component: SettingView
    },
    {
      path: '/login',
      name: 'LoginView',
      component: SettingView
    }
  ],
})

export default router
