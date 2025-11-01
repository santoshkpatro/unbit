<!-- AppLayout.vue -->
<script setup>
import { ref, watchEffect } from 'vue'
import { useRoute, useRouter, RouterView } from 'vue-router'
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons-vue'
import {
  Home,
  ListChecks,
  FolderClosed,
  Settings,
  BookOpen,
  HelpCircle,
  Shield,
  KeyRound,
  ScrollText,
  Plug,
  Activity,
  Users,
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()

const collapsed = ref(false)
const selectedKeys = ref([route.path])
const openKeys = ref([])

watchEffect(() => {
  selectedKeys.value = [route.path]
})

const onMenuClick = ({ key }) => {
  if (key && key !== route.path) router.push(key)
}
</script>

<template>
  <a-layout class="app-layout">
    <a-layout-header class="app-header">
      <div class="brand" @click="router.push('/')">
        <img class="logo" alt="Logo" />
        <span class="title">My App</span>
      </div>

      <a-button
        type="text"
        class="collapse-btn"
        @click="collapsed = !collapsed"
        :aria-label="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
      >
        <template #icon>
          <component :is="collapsed ? MenuUnfoldOutlined : MenuFoldOutlined" />
        </template>
      </a-button>
    </a-layout-header>

    <a-layout class="body-layout">
      <a-layout-sider
        class="app-sider"
        :collapsed="collapsed"
        collapsible
        :trigger="null"
        width="240"
        theme="light"
      >
        <a-menu
          theme="light"
          mode="inline"
          v-model:selectedKeys="selectedKeys"
          v-model:openKeys="openKeys"
          @click="onMenuClick"
        >
          <!-- Profile / Tenant -->
          <a-menu-item :disabled="true" class="profile-item">
            <div class="profile">
              <div class="avatar" aria-hidden="true">S</div>
              <div class="who">
                <div class="name">Santosh Kumar Patro</div>
                <div class="role">Owner</div>
              </div>
            </div>
          </a-menu-item>

          <!-- Primary nav -->
          <a-menu-item key="/">
            <template #icon><Home class="menu-icon" /></template>
            Dashboard
          </a-menu-item>

          <a-menu-item key="/issues">
            <template #icon><ListChecks class="menu-icon" /></template>
            Issues
          </a-menu-item>

          <a-menu-item key="/projects">
            <template #icon><FolderClosed class="menu-icon" /></template>
            Projects
          </a-menu-item>

          <a-menu-item key="/settings">
            <template #icon><Settings class="menu-icon" /></template>
            Settings
          </a-menu-item>

          <a-menu-divider />

          <!-- Docs & Support -->
          <a-menu-item key="/docs">
            <template #icon><BookOpen class="menu-icon" /></template>
            Docs & SDKs
          </a-menu-item>

          <a-menu-item key="/support">
            <template #icon><HelpCircle class="menu-icon" /></template>
            Help & Support
          </a-menu-item>

          <!-- Security suite -->
          <a-sub-menu key="/security">
            <template #icon><Shield class="menu-icon" /></template>
            <template #title>Security</template>

            <a-menu-item key="/access-control">
              <template #icon><KeyRound class="menu-icon" /></template>
              Access Control (RBAC)
            </a-menu-item>

            <a-menu-item key="/audit-logs">
              <template #icon><ScrollText class="menu-icon" /></template>
              Audit Logs
            </a-menu-item>

            <a-menu-item key="/compliance">
              <template #icon><ScrollText class="menu-icon" /></template>
              Compliance
            </a-menu-item>

            <a-menu-item key="/organizations">
              <template #icon><Users class="menu-icon" /></template>
              Organizations & Teams
            </a-menu-item>
          </a-sub-menu>

          <!-- Platform ops -->
          <a-sub-menu key="/system">
            <template #icon><Activity class="menu-icon" /></template>
            <template #title>System</template>

            <a-menu-item key="/system-status">
              <template #icon><Activity class="menu-icon" /></template>
              System Status
            </a-menu-item>

            <a-menu-item key="/integrations">
              <template #icon><Plug class="menu-icon" /></template>
              Integrations
            </a-menu-item>

            <a-menu-item key="/webhooks">
              <template #icon><Plug class="menu-icon" /></template>
              Webhooks
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-sider>

      <a-layout class="main-layout">
        <a-layout-content class="app-content">
          <div class="content-inner">
            <RouterView />
          </div>
        </a-layout-content>

        <a-layout-footer class="app-footer"> Powered by UnBiT </a-layout-footer>
      </a-layout>
    </a-layout>
  </a-layout>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 16px;
  background: #001529;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.logo {
  width: 24px;
  height: 24px;
  background: #fff;
  border-radius: 6px;
}

.title {
  color: #fff;
  font-weight: 600;
  letter-spacing: 0.2px;
}

.collapse-btn {
  color: #fff;
}

.body-layout {
  min-height: calc(100vh - 56px);
}

.app-sider {
  background: #fff;
}

.app-sider :deep(.ant-menu) {
  background: #fff;
  border-right: none;
}

.menu-icon {
  width: 18px;
  height: 18px;
}

.main-layout {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.app-content {
  flex: 1 1 auto;
  background: #f5f7fa;
  padding: 8px;
}

.content-inner {
  border-radius: 8px;
  padding: 10px;
}

.app-footer {
  text-align: center;
  padding: 10px 0 14px;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.55);
  background: #fff;
}

/* Profile item styles */
.profile-item {
  cursor: default !important;
}

.profile {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #1677ff;
  color: #fff;
  font-weight: 700;
  display: grid;
  place-items: center;
}

.who .name {
  font-weight: 600;
  line-height: 1.1;
}

.who .role {
  font-size: 12px;
  opacity: 0.65;
}

@media (max-width: 768px) {
  .title {
    display: none;
  }
}
</style>
