<script lang="ts" setup>
import 'bootstrap/js/dist/dropdown'
import 'bootstrap/js/dist/collapse'
import { nanoid } from 'nanoid'
import { RouterLink, useRoute } from 'vue-router'

const navbarCollapseId = nanoid()
const expandSize = 'sm'
const navItems: Array<{ meAuto?: boolean, items: Array<{ text: string, name: string }> }> = [
  {
    meAuto: true,
    items: [
      {
        name: 'ParseView',
        text: '视频解析'
      },
      {
        name: 'DownloadView',
        text: '下载管理'
      },
      {
        name: 'SettingView',
        text: '系统设置'
      }
    ]
  },
  {
    items: [
      {
        name: 'LoginView',
        text: '登录 / 注册'
      }
    ]
  }
]
const route = useRoute()
</script>

<template>
  <div :class="`navbar navbar-expand-${expandSize} bg-success-subtle`" data-bs-theme="dark">
    <div class="container">
      <RouterLink class="navbar-brand text-success-emphasis" :to="{ name: 'HomeView' }">Bilidown
        <span class="text-body-secondary">v3</span>
      </RouterLink>
      <button class="navbar-toggler" data-bs-toggle="collapse" :data-bs-target="`#${navbarCollapseId}`">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" :id="navbarCollapseId">
        <div v-for="group, index in navItems" :key="index"
          :class="`navbar-nav ${group.meAuto && 'me-auto'} mb-2 mb-${expandSize}-0`">
          <div v-for="navItem, index2 in group.items" :key="index2" class="nav-item">
            <RouterLink :class="{ 'nav-link': true, active: route.name === navItem.name }" :to="{ name: navItem.name }">
              {{ navItem.text }}
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
