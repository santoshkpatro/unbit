<script setup>
import { ref, watchEffect, computed } from 'vue'
import { useRoute, useRouter, RouterView } from 'vue-router'
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
  Star,
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useSettingStore } from '@/stores/setting'
import { storeToRefs } from 'pinia'

const router = useRouter()
const route = useRoute()

const { loggedInUser } = storeToRefs(useAuthStore())
const { setting } = storeToRefs(useSettingStore())

// no collapsed state at all now
const selectedKeys = ref([route.path])
const openKeys = ref([])

const GITHUB_STAR_URL = 'https://github.com/santoshkpatro/unbit'
const DEMO_URL = 'https://demo.unbit.com'

watchEffect(() => {
  selectedKeys.value = [route.path]
})

const onMenuClick = ({ key }) => {
  if (key === '__github_star') {
    window.open(GITHUB_STAR_URL, '_blank', 'noopener,noreferrer')
    return
  }
  if (key === '__demo') {
    window.open(DEMO_URL, '_blank', 'noopener,noreferrer')
    return
  }
  if (key && key !== route.path) router.push(key)
}

const initials = computed(() => {
  const parts = (loggedInUser.value.firstName || '').trim().split(/\s+/)
  const chars = (parts[0]?.[0] || '') + (parts[1]?.[0] || '')
  return chars.toUpperCase() || 'U'
})
</script>

<template>
  <a-layout class="app-layout">
    <!-- Sidebar only -->
    <a-layout-sider class="app-sider" width="216" theme="light" :collapsed="false" :trigger="null">
      <!-- Sider header (brand NO collapse button) -->
      <div class="sider-header">
        <div class="brand" @click="router.push('/')" role="button" tabindex="0">
          <img class="logo" alt="Logo" />
          <span class="title">{{ setting.org.siteName }}</span>
        </div>
      </div>

      <!-- Logged-in user -->
      <div class="user-card">
        <a-avatar :size="36" :src="loggedInUser.avatar">
          {{ initials }}
        </a-avatar>
        <div class="user-meta">
          <div class="name" :title="loggedInUser.firstName">{{ loggedInUser.firstName }}</div>
          <div class="role">{{ loggedInUser.email }}</div>
        </div>
      </div>

      <!-- Navigation -->
      <a-menu
        class="nav-menu"
        theme="light"
        mode="inline"
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        @click="onMenuClick"
      >
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

        <a-menu-divider class="menu-divider" />

        <a-menu-item key="/docs">
          <template #icon><BookOpen class="menu-icon" /></template>
          Docs & SDKs
        </a-menu-item>

        <a-menu-item key="/support">
          <template #icon><HelpCircle class="menu-icon" /></template>
          Help & Support
        </a-menu-item>

        <a-menu-divider class="menu-divider" />

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

      <div class="sidebar-bottom">
        <div class="external-links">
          <a
            :href="GITHUB_STAR_URL"
            target="_blank"
            rel="noopener noreferrer"
            class="external-link"
          >
            <Star class="link-icon" />
            <span>Star Us on GitHub</span>
          </a>
          <a :href="DEMO_URL" target="_blank" rel="noopener noreferrer" class="external-link">
            <Activity class="link-icon" />
            <span>See Demo</span>
          </a>
        </div>
        <div class="sidebar-footer">
          <span>Powered by UnBiT</span>
        </div>
      </div>
    </a-layout-sider>

    <a-layout class="main-layout">
      <a-layout-content class="app-content">
        <div class="content-inner">
          <RouterView />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<style scoped>
/* Layout */
.app-layout {
  height: 100vh;
  background: #f5f7fa;
}

/* Sider */
.app-sider {
  background: #fff;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* Sider header (brand + collapse) */
.sider-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 10px;
  border-bottom: 1px solid #f2f4f7;
  background: #fff;
  position: sticky;
  top: 0;
  z-index: 2;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.logo {
  width: 22px;
  height: 22px;
  background: linear-gradient(135deg, #1677ff 0%, #69b1ff 100%);
  border-radius: 6px;
}
.title {
  color: #111827;
  font-weight: 700;
  letter-spacing: 0.2px;
  font-size: 14px;
}
.collapse-btn {
  color: #111827;
}

/* Logged-in user card */
.user-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  margin: 10px 10px 6px;
  border: 1px solid #f1f1f1;
  border-radius: 12px;
  background: #fafafa;
  transition:
    padding 0.2s ease,
    margin 0.2s ease;
}
.user-card.collapsed {
  justify-content: center;
  margin: 10px 8px 6px;
  padding: 10px 8px;
}
.user-meta {
  min-width: 0;
}
.user-meta .name {
  font-weight: 600;
  line-height: 1.1;
  font-size: 13px;
  color: #111827;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
.user-meta .role {
  font-size: 11px;
  opacity: 0.7;
}

/* Menu (modern + compact + consistent) */
.nav-menu {
  padding: 6px 8px 10px;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}
.app-sider :deep(.ant-menu) {
  background: #fff;
  border-right: 0;
  padding-inline: 0 !important;
}

/* Row sizing - more compact */
.app-sider :deep(.ant-menu-item),
.app-sider :deep(.ant-menu-submenu-title) {
  height: 32px !important;
  line-height: 32px !important;
  display: flex;
  align-items: center;
  border-radius: 8px;
  margin-block: 2px;
}

/* Consistent padding when expanded */
.app-sider :deep(.ant-menu-item),
.app-sider :deep(.ant-menu-submenu-title) {
  padding-inline: 12px !important;
}

/* When collapsed - center everything */
.app-sider.ant-layout-sider-collapsed :deep(.ant-menu-item),
.app-sider.ant-layout-sider-collapsed :deep(.ant-menu-submenu-title) {
  padding-inline: 0 !important;
  justify-content: center;
}

/* Force icon centering when collapsed */
.app-sider.ant-layout-sider-collapsed :deep(.ant-menu-item-icon),
.app-sider.ant-layout-sider-collapsed :deep(.ant-menu-submenu-icon) {
  margin-inline-end: 0 !important;
}

/* Indent ONLY level-2 children (items inside sub menus) */
.app-sider :deep(.ant-menu-sub .ant-menu-item),
.app-sider :deep(.ant-menu-sub .ant-menu-submenu-title) {
  padding-inline-start: 32px !important;
}

/* Typography */
.app-sider :deep(.ant-menu-title-content) {
  font-size: 13px;
  font-weight: 500;
}

/* Icons */
.menu-icon {
  width: 16px;
  height: 16px;
}

.menu-divider {
  margin: 4px 10px !important;
}

/* Sidebar Bottom (Links + Footer) */
.sidebar-bottom {
  margin-top: auto;
  background: #fff;
}

.external-links {
  padding: 6px 8px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-bottom: 1px solid #f0f0f0;
}

.external-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  height: 32px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: inherit;
  text-decoration: none;
  transition: all 0.2s ease;
}

.external-link:hover {
  background: var(--ant-primary-color-hover-bg, #e6f4ff);
}

.link-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.external-links-collapsed {
  padding: 6px 8px 10px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 2px;
  border-bottom: 1px solid #f0f0f0;
}

.external-link-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 32px;
  border-radius: 8px;
  color: inherit;
  text-decoration: none;
  transition: all 0.2s ease;
}

.external-link-icon:hover {
  background: var(--ant-primary-color-hover-bg, #e6f4ff);
}

.sidebar-footer {
  text-align: center;
  padding: 16px;
  font-size: 11px;
  font-weight: 400;
  color: rgba(0, 0, 0, 0.35);
  background: #fff;
}

.sidebar-footer.collapsed {
  padding: 16px 8px;
}

.footer-icon {
  display: inline-block;
  font-size: 10px;
  color: rgba(0, 0, 0, 0.35);
  font-weight: 400;
}

/* Main */
.main-layout {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 100vh;
  overflow: hidden;
}
.app-content {
  flex: 1 1 auto;
  background: #f5f7fa;
  padding: 2px;
  overflow-y: auto;
  overflow-x: hidden;
  height: 100%;
}
.content-inner {
  border-radius: 10px;
  padding: 6px;
}

@media (max-width: 768px) {
  .title {
    display: none;
  }
}
</style>
