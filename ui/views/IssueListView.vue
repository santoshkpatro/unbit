<script setup>
import { onMounted, ref, computed, nextTick, watch } from 'vue'
import { projectListAPI } from '@/api/projects'
import { eventIssuesAPI } from '@/api/events'
import { Search, TrendingUp } from 'lucide-vue-next'
import {
  CheckCircleOutlined,
  InboxOutlined,
  FilterOutlined,
  DeleteOutlined,
  ExportOutlined,
  MoreOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

/**
 * Raw data, same shape as API:
 * [
 *  {
 *    event: {...},
 *    project: {...},
 *    group: {...},
 *    last7Days: [{day,count}, ...]
 *  }
 * ]
 */
const projects = ref([]) // optional: used to show project filter labels
const events = ref([]) // raw payload array
const selectedRowKeys = ref([])
const searchText = ref('')
const statusFilter = ref(null) // 'unresolved' | 'resolved' | 'archived'
const projectFilter = ref(null) // project.id
const chartInstances = ref({})

const columns = [
  { title: 'Issue', dataIndex: ['event', 'message'], key: 'message', width: '40%' },
  {
    title: 'Events',
    dataIndex: ['group', 'eventCount'],
    key: 'count',
    width: '10%',
    align: 'center',
  },
  { title: 'Assignee', dataIndex: ['group', 'assignee'], key: 'assignee', width: '12%' },
  { title: 'Trend', dataIndex: 'last7Days', key: 'trend', width: '13%', align: 'center' },
  { title: 'Last Seen', dataIndex: ['event', 'timestamp'], key: 'timestamp', width: '15%' },
  { title: '', key: 'actions', width: '10%', align: 'center' },
]

// filters/search kept simple, operate on RAW events array
const filteredEvents = computed(() => {
  let list = events.value

  if (searchText.value) {
    const s = searchText.value.toLowerCase()
    list = list.filter((row) => {
      const msg = row.event?.message?.toLowerCase() || ''
      const proj = row.project?.name?.toLowerCase() || ''
      const email = row.group?.assignee?.email?.toLowerCase() || ''
      return msg.includes(s) || proj.includes(s) || email.includes(s)
    })
  }

  if (statusFilter.value) {
    list = list.filter((row) => row.group?.status === statusFilter.value)
  }

  if (projectFilter.value) {
    list = list.filter((row) => row.project?.id === projectFilter.value)
  }

  return list
})

const unresolvedCount = computed(
  () => events.value.filter((row) => row.group?.status === 'unresolved').length,
)

const projectOptions = computed(() => projects.value.map((p) => ({ label: p.name, value: p.id })))

const hasActiveFilters = computed(
  () => !!(searchText.value || statusFilter.value || projectFilter.value),
)

const rowSelection = {
  selectedRowKeys: selectedRowKeys,
  onChange: (keys) => (selectedRowKeys.value = keys),
}

const loadProjects = async () => {
  try {
    projects.value = await projectListAPI()
  } catch {
    // non-fatal
  }
}

const loadEvents = async () => {
  try {
    events.value = await eventIssuesAPI()
  } catch {
    message.error('Failed to load events')
  }
}

// tiny helpers that use payload fields directly
const copyIssueLink = (eventId) => {
  const link = `${window.location.origin}/issues/${eventId}`
  navigator.clipboard
    .writeText(link)
    .then(() => message.success('Issue link copied to clipboard'))
    .catch(() => message.error('Failed to copy link'))
}

const clearFilters = () => {
  searchText.value = ''
  statusFilter.value = null
  projectFilter.value = null
}

const formatTimestamp = (ts) => {
  if (!ts || ts === '0001-01-01T00:00:00Z') return 'Just now'
  const date = new Date(ts)
  const now = new Date()
  const diff = now - date
  const m = Math.floor(diff / 60000)
  const h = Math.floor(diff / 3600000)
  const d = Math.floor(diff / 86400000)
  if (m < 60) return `${m}m ago`
  if (h < 24) return `${h}h ago`
  if (d < 7) return `${d}d ago`
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
  })
}

/** charts — simple and direct from last7Days (no remap) */
const renderTrendChart = (canvasId, last7Days) => {
  if (!last7Days?.length) return

  if (chartInstances.value[canvasId]) {
    chartInstances.value[canvasId].destroy()
  }

  nextTick(() => {
    const canvas = document.getElementById(canvasId)
    if (!canvas) return
    const ctx = canvas.getContext('2d')

    const labels = last7Days.map((d) =>
      new Date(d.day).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    )
    const counts = last7Days.map((d) => d.count)

    chartInstances.value[canvasId] = new Chart(ctx, {
      type: 'line',
      data: {
        labels,
        datasets: [
          {
            label: 'Events',
            data: counts,
            borderColor: '#51bc8f',
            backgroundColor: 'rgba(81,188,143,0.1)',
            borderWidth: 2,
            fill: true,
            tension: 0.4,
            pointRadius: 0,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: { x: { display: false }, y: { display: false, beginAtZero: true } },
      },
    })
  })
}

const renderAllCharts = () => {
  nextTick(() => {
    filteredEvents.value.forEach((row) => {
      if (row.last7Days?.length) {
        renderTrendChart(`trend-${row.event?.id}`, row.last7Days)
      }
    })
  })
}

watch(filteredEvents, () => renderAllCharts())

onMounted(async () => {
  await loadProjects()
  await loadEvents()
  renderAllCharts()
})

/** simple bulk actions (demo toasts only) */
const handleResolve = () => {
  if (!selectedRowKeys.value.length) return message.warning('Please select at least one issue')
  message.success(`Resolved ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}
const handleArchive = () => {
  if (!selectedRowKeys.value.length) return message.warning('Please select at least one issue')
  message.success(`Archived ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}
const handleDelete = () => {
  if (!selectedRowKeys.value.length) return message.warning('Please select at least one issue')
  message.success(`Deleted ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}
const handleExport = () => {
  if (!selectedRowKeys.value.length) return message.warning('Please select at least one issue')
  message.success(`Exported ${selectedRowKeys.value.length} issue(s)`)
}
</script>

<template>
  <main class="issue-tracker">
    <div class="header">
      <a-typography-title :level="3" style="margin: 0">Issues</a-typography-title>
      <a-badge :count="unresolvedCount" :number-style="{ backgroundColor: '#f5222d' }" />
    </div>

    <div class="table-container">
      <a-table
        size="small"
        :pagination="false"
        :columns="columns"
        :data-source="filteredEvents"
        :row-selection="{
          selectedRowKeys,
          onChange: (keys) => (selectedRowKeys = keys),
        }"
        :row-key="(record) => record.event?.id"
      >
        <template #title>
          <div v-if="selectedRowKeys.length" class="table-header-actions">
            <span class="selected-count">{{ selectedRowKeys.length }} selected</span>
            <a-space>
              <a-button type="primary" @click="handleResolve">
                <template #icon><CheckCircleOutlined /></template>Resolve
              </a-button>
              <a-button @click="handleArchive">
                <template #icon><InboxOutlined /></template>Archive
              </a-button>
              <a-dropdown>
                <a-button
                  ><template #icon><MoreOutlined /></template>More actions</a-button
                >
                <template #overlay>
                  <a-menu>
                    <a-menu-item key="delete" @click="handleDelete">
                      <DeleteOutlined />Delete
                    </a-menu-item>
                    <a-menu-item key="export" @click="handleExport">
                      <ExportOutlined />Export
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </div>

          <div v-else class="table-header-filters">
            <div class="filters-left">
              <a-input
                v-model:value="searchText"
                placeholder="Search by message, project, or assignee..."
                style="width: 350px"
                allow-clear
              >
                <template #prefix><Search :size="16" /></template>
              </a-input>

              <a-select
                v-model:value="statusFilter"
                placeholder="All Statuses"
                style="width: 150px"
                allow-clear
              >
                <a-select-option value="unresolved">Unresolved</a-select-option>
                <a-select-option value="resolved">Resolved</a-select-option>
                <a-select-option value="archived">Archived</a-select-option>
              </a-select>

              <a-select
                v-model:value="projectFilter"
                placeholder="All Projects"
                style="width: 180px"
                allow-clear
                :options="projectOptions"
              />

              <a-button @click="clearFilters" v-if="hasActiveFilters">
                <template #icon><FilterOutlined /></template>Clear
              </a-button>
            </div>
          </div>
        </template>

        <template #bodyCell="{ column, record }">
          <!-- message / context -->
          <template v-if="column.key === 'message'">
            <div class="issue-cell">
              <div class="issue-content">
                <a-typography-link class="issue-title" href="#">
                  {{ record.event?.message }}
                </a-typography-link>

                <div class="issue-meta">
                  <span class="project-name">{{ record.project?.name || 'Unknown' }}</span>
                  <span class="meta-separator">|</span>
                  <span class="issue-type-text">{{ record.event?.type }}</span>
                </div>

                <div v-if="record.event?.stacktraceFirst" class="stacktrace-preview">
                  <span class="stacktrace-function">{{
                    record.event.stacktraceFirst.function
                  }}</span>
                  <span class="stacktrace-separator">in</span>
                  <span class="stacktrace-file">{{ record.event.stacktraceFirst.file }}</span>
                  <span class="stacktrace-separator">at line</span>
                  <span class="stacktrace-line">{{ record.event.stacktraceFirst.line }}</span>
                </div>
              </div>
            </div>
          </template>

          <!-- counts -->
          <template v-if="column.key === 'count'">
            <a-badge
              :count="record.group?.eventCount || 0"
              :overflow-count="999"
              :number-style="{ backgroundColor: '#f0f0f0', color: '#595959', fontWeight: 600 }"
            />
          </template>

          <!-- assignee -->
          <template v-if="column.key === 'assignee'">
            <div v-if="record.group?.assignee" class="assignee-cell">
              <a-avatar size="small" :style="{ backgroundColor: '#722ed1', marginRight: '8px' }">
                {{ record.group.assignee.firstName?.[0] }}
              </a-avatar>
              <span>{{ record.group.assignee.firstName }}</span>
            </div>
            <a-typography-text type="secondary" v-else>Unassigned</a-typography-text>
          </template>

          <!-- trend -->
          <template v-if="column.key === 'trend'">
            <div v-if="record.last7Days?.length" class="trend-chart-container">
              <canvas :id="`trend-${record.event?.id}`" class="trend-chart"></canvas>
            </div>
            <div v-else class="trend-placeholder">
              <TrendingUp :size="16" style="color: #d9d9d9" />
            </div>
          </template>

          <!-- last seen -->
          <template v-if="column.key === 'timestamp'">
            <a-typography-text type="secondary">
              {{ formatTimestamp(record.event?.timestamp) }}
            </a-typography-text>
          </template>

          <!-- actions -->
          <template v-if="column.key === 'actions'">
            <a-tooltip title="Copy issue link">
              <a-button
                type="text"
                size="small"
                @click.stop="copyIssueLink(record.event?.id)"
                aria-label="Copy issue link"
              >
                <template #icon><ExportOutlined style="font-size: 14px" /></template>
              </a-button>
            </a-tooltip>
          </template>
        </template>
      </a-table>
    </div>
  </main>
</template>

<style scoped>
/* page fills the app-content area */
.issue-tracker {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 24px;
  background: #fff;
}

.header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

/* inner scroller */
.table-container {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  border-radius: 8px;
}

/* sticky toolbar inside scroller */
.table-header-actions,
.table-header-filters {
  position: sticky;
  top: 0;
  z-index: 2;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
  padding: 12px 16px;
}

.filters-left {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
}

.selected-count {
  font-weight: 600;
  color: #1890ff;
}

.issue-cell {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.issue-content {
  flex: 1;
  min-width: 0;
}

.issue-title {
  font-weight: 500;
  margin-bottom: 4px;
  word-break: break-word;
  color: #262626;
}
.issue-title:hover {
  color: #51bc8f;
}

.issue-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  flex-wrap: wrap;
}

.project-name {
  color: #1890ff;
  font-weight: 500;
  font-size: 12px;
}
.meta-separator {
  color: #d9d9d9;
  margin: 0 2px;
}
.issue-type-text {
  color: #8c8c8c;
  font-size: 12px;
  text-transform: capitalize;
}

.stacktrace-preview {
  margin-top: 6px;
  font-size: 11px;
  color: #8c8c8c;
  font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', 'Courier New', monospace;
  line-height: 1.4;
  padding: 4px 8px;
  background: #fafafa;
  border-radius: 3px;
  border-left: 2px solid #e8e8e8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.stacktrace-function {
  color: #722ed1;
  font-weight: 500;
}
.stacktrace-separator {
  color: #bfbfbf;
  margin: 0 4px;
}
.stacktrace-file {
  color: #595959;
}
.stacktrace-line {
  color: #1890ff;
  font-weight: 500;
}

.assignee-cell {
  display: flex;
  align-items: center;
}

.trend-chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 0;
}
.trend-chart {
  width: 120px;
  height: 40px;
}
.trend-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
}
:deep(.ant-table-tbody > tr:hover) {
  cursor: pointer;
}
:deep(.ant-table-title) {
  padding: 0 !important;
}
</style>
