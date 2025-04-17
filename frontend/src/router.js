import { createRouter, createWebHashHistory } from 'vue-router'
import LabelPrint from './views/LabelPrint.vue'
import PrintHistory from './views/PrintHistory.vue'
import TemplateManage from './views/TemplateManage.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      redirect: '/label-print'
    },
    {
      path: '/label-print',
      name: 'LabelPrint',
      component: LabelPrint
    },
    {
      path: '/print-history',
      name: 'PrintHistory',
      component: PrintHistory
    },
    {
      path: '/template-manage',
      name: 'TemplateManage',
      component: TemplateManage
    }
  ]
})

export default router