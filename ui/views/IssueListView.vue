<script setup>
import { onMounted, ref, computed } from 'vue'
import { projectListAPI } from '@/api/projects'
import { eventListAPI } from '@/api/events'
import { Search, ChevronDown } from 'lucide-vue-next'
import {
  CheckCircleOutlined,
  InboxOutlined,
  FilterOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

const projects = ref([])
const events = ref([])
const selectedRowKeys = ref([])
const searchText = ref('')
const statusFilter = ref(null)
const projectFilter = ref(null)

const columns = [
  {
    title: 'Issue',
    dataIndex: 'summary',
    key: 'summary',
    width: '40%',
  },
  {
    title: 'Project',
    dataIndex: ['project', 'name'],
    key: 'project',
    width: '15%',
  },
  {
    title: 'Events',
    dataIndex: 'eventCount',
    key: 'eventCount',
    width: '10%',
    align: 'center',
  },
  {
    title: 'Assignee',
    dataIndex: ['assignee', 'firstName'],
    key: 'assignee',
    width: '15%',
  },
  {
    title: 'Last Seen',
    dataIndex: 'lastSeenAt',
    key: 'lastSeenAt',
    width: '20%',
  },
]

const filteredIssues = computed(() => {
  let result = events.value.map((event) => ({
    ...event,
    summary: event.message,
    status: event.group?.status || 'unresolved',
    assignee: event.group?.assignee,
    eventCount: event.group?.eventCount || 0,
    lastSeenAt: event.timestamp,
  }))

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (issue) =>
        issue.summary?.toLowerCase().includes(search) ||
        issue.assignee?.email?.toLowerCase().includes(search),
    )
  }

  if (statusFilter.value) {
    result = result.filter((issue) => issue.status === statusFilter.value)
  }

  if (projectFilter.value) {
    result = result.filter((issue) => issue.project?.id === projectFilter.value)
  }

  return result
})

const unresolvedCount = computed(() => {
  return events.value.filter((e) => e.group?.status === 'unresolved').length
})

const projectOptions = computed(() => {
  return projects.value.map((p) => ({
    label: p.name,
    value: p.id,
  }))
})

const rowSelection = {
  selectedRowKeys: selectedRowKeys,
  onChange: (keys) => {
    selectedRowKeys.value = keys
  },
}

const loadProjects = async () => {
  projects.value = await projectListAPI()
}

const loadEvents = async () => {
  events.value = await eventListAPI()
}

const handleResolve = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('Please select at least one issue')
    return
  }
  message.success(`Resolved ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}

const handleArchive = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('Please select at least one issue')
    return
  }
  message.success(`Archived ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}

const clearFilters = () => {
  searchText.value = ''
  statusFilter.value = null
  projectFilter.value = null
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  if (days < 7) return `${days}d ago`

  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  })
}

const getTypeIcon = (type) => {
  return type === 'exception' ? 'error' : 'info'
}

onMounted(async () => {
  await loadProjects()
  await loadEvents()
})
</script>

<template>
  <main class="issue-tracker">
    <div class="header">
      <div class="header-content">
        <a-typography-title :level="2" style="margin: 0">Issues</a-typography-title>
        <a-badge :count="unresolvedCount" :number-style="{ backgroundColor: '#f5222d' }">
          <a-tag color="default" style="margin: 0; padding: 4px 12px">
            {{ unresolvedCount }} Unresolved
          </a-tag>
        </a-badge>
      </div>
    </div>

    <div class="content-wrapper">
      <div class="toolbar">
        <div class="toolbar-left">
          <a-input
            v-model:value="searchText"
            placeholder="Search by message, project, or assignee..."
            class="search-input"
            allow-clear
            size="large"
          >
            <template #prefix>
              <Search :size="18" style="color: #8c8c8c" />
            </template>
          </a-input>

          <a-select
            v-model:value="statusFilter"
            placeholder="All Statuses"
            class="filter-select"
            allow-clear
            size="large"
          >
            <a-select-option value="unresolved">Unresolved</a-select-option>
            <a-select-option value="resolved">Resolved</a-select-option>
            <a-select-option value="archived">Archived</a-select-option>
          </a-select>

          <a-select
            v-model:value="projectFilter"
            placeholder="All Projects"
            class="filter-select"
            allow-clear
            size="large"
            :options="projectOptions"
          />

          <a-button
            @click="clearFilters"
            size="large"
            v-if="searchText || statusFilter || projectFilter"
          >
            <template #icon><FilterOutlined /></template>
            Clear
          </a-button>
        </div>

        <div class="toolbar-right" v-if="selectedRowKeys.length > 0">
          <a-space :size="12">
            <a-button type="primary" @click="handleResolve" size="large">
              <template #icon><CheckCircleOutlined /></template>
              Resolve ({{ selectedRowKeys.length }})
            </a-button>
            <a-button @click="handleArchive" size="large">
              <template #icon><InboxOutlined /></template>
              Archive
            </a-button>
          </a-space>
        </div>
      </div>

      <a-table
        :columns="columns"
        :data-source="filteredIssues"
        :row-selection="rowSelection"
        :row-key="(record) => record.id"
        :pagination="{
          pageSize: 20,
          showSizeChanger: true,
          pageSizeOptions: ['10', '20', '50', '100'],
          showTotal: (total) => `${total} issue${total !== 1 ? 's' : ''}`,
        }"
        class="issues-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'summary'">
            <div class="issue-cell">
              <div class="issue-icon">
                <ExclamationCircleOutlined style="color: #f5222d; font-size: 16px" />
              </div>
              <div class="issue-content">
                <div class="issue-title">{{ record.summary }}</div>
                <div class="issue-meta">
                  <a-tag :color="record.type === 'exception' ? 'red' : 'blue'" style="margin: 0">
                    {{ record.type }}
                  </a-tag>
                  <span class="issue-group-id">{{ record.group?.id }}</span>
                </div>
              </div>
            </div>
          </template>

          <template v-if="column.key === 'project'">
            <div class="project-cell">
              <a-typography-text strong>{{ record.project?.name || '-' }}</a-typography-text>
            </div>
          </template>

          <template v-if="column.key === 'eventCount'">
            <a-badge
              :count="record.eventCount"
              :number-style="{ backgroundColor: '#f0f0f0', color: '#595959', fontWeight: '600' }"
              :overflow-count="999"
            />
          </template>

          <template v-if="column.key === 'assignee'">
            <div v-if="record.assignee" class="assignee-cell">
              <a-avatar
                size="small"
                :style="{
                  backgroundColor: '#722ed1',
                  marginRight: '8px',
                  fontWeight: '500',
                }"
              >
                {{ record.assignee.firstName?.[0] }}
              </a-avatar>
              <span>{{ record.assignee.firstName }}</span>
            </div>
            <a-typography-text type="secondary" v-else>Unassigned</a-typography-text>
          </template>

          <template v-if="column.key === 'lastSeenAt'">
            <a-typography-text type="secondary" style="font-size: 13px">
              {{ formatDate(record.lastSeenAt) }}
            </a-typography-text>
          </template>
        </template>
      </a-table>
    </div>
  </main>
</template>

<style scoped>
.issue-tracker {
  min-height: 100vh;
  background: #f5f5f7;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  padding: 20px 32px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 16px;
  max-width: 1600px;
  margin: 0 auto;
}

.content-wrapper {
  max-width: 1600px;
  margin: 0 auto;
  padding: 24px 32px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  flex: 1;
}

.toolbar-right {
  display: flex;
  align-items: center;
}

.search-input {
  width: 400px;
  border-radius: 6px;
}

.filter-select {
  width: 180px;
}

.issue-cell {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 4px 0;
}

.issue-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.issue-content {
  flex: 1;
  min-width: 0;
}

.issue-title {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  margin-bottom: 6px;
  word-break: break-word;
  line-height: 1.4;
}

.issue-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.issue-group-id {
  color: #8c8c8c;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 11px;
}

.project-cell {
  color: #1890ff;
  font-size: 13px;
}

.assignee-cell {
  display: flex;
  align-items: center;
  font-size: 13px;
}

.issues-table {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow:
    0 1px 2px 0 rgba(0, 0, 0, 0.03),
    0 1px 6px -1px rgba(0, 0, 0, 0.02),
    0 2px 4px 0 rgba(0, 0, 0, 0.02);
}

:deep(.ant-table) {
  font-size: 13px;
}

:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
  color: #595959;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 2px solid #e8e8e8;
  padding: 14px 16px;
}

:deep(.ant-table-tbody > tr) {
  transition: all 0.2s;
}

:deep(.ant-table-tbody > tr:hover) {
  cursor: pointer;
  background: #fafafa;
}

:deep(.ant-table-tbody > tr > td) {
  border-bottom: 1px solid #f0f0f0;
  padding: 16px;
}

:deep(.ant-table-cell) {
  vertical-align: middle;
}

:deep(.ant-checkbox-wrapper) {
  margin-right: 0;
}

:deep(.ant-pagination) {
  margin: 16px 16px;
}

:deep(.ant-select-selector) {
  border-radius: 6px !important;
}

:deep(.ant-input) {
  border-radius: 6px;
}

:deep(.ant-btn) {
  border-radius: 6px;
  font-weight: 500;
}

:deep(.ant-tag) {
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  padding: 2px 8px;
  border: none;
}
</style>
