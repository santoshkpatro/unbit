<script setup>
import { onMounted, ref, computed, nextTick, watch } from 'vue'
import { projectListAPI } from '@/api/projects'
import { eventIssuesAPI } from '@/api/events'
import { Search, TrendingUp } from 'lucide-vue-next'
import {
  CheckCircleOutlined,
  InboxOutlined,
  FilterOutlined,
  ExclamationCircleOutlined,
  DeleteOutlined,
  ExportOutlined,
  MoreOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { Chart, registerables } from 'chart.js'

// Register Chart.js components
Chart.register(...registerables)

const projects = ref([])
const events = ref([])
const selectedRowKeys = ref([])
const searchText = ref('')
const statusFilter = ref(null)
const projectFilter = ref(null)
const chartInstances = ref({})

const columns = [
  {
    title: 'Issue',
    dataIndex: 'summary',
    key: 'summary',
    width: '40%',
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
    width: '12%',
  },
  {
    title: 'Trend',
    dataIndex: 'trend',
    key: 'trend',
    width: '13%',
    align: 'center',
  },
  {
    title: 'Last Seen',
    dataIndex: 'lastSeenAt',
    key: 'lastSeenAt',
    width: '15%',
  },
  {
    title: '',
    key: 'actions',
    width: '10%',
    align: 'center',
  },
]

const filteredIssues = computed(() => {
  let result = events.value.map((item) => ({
    id: item.event?.id,
    type: item.event?.type,
    level: item.event?.level,
    summary: item.event?.message,
    timestamp: item.event?.timestamp,
    project: item.project,
    status: item.group?.status || 'unresolved',
    assignee: item.group?.assignee,
    eventCount: item.group?.eventCount || 0,
    lastSeenAt: item.event?.timestamp,
    trend: item.last7Days || [],
    stacktrace: item.event?.stacktraceFirst,
  }))

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(
      (issue) =>
        issue.summary?.toLowerCase().includes(search) ||
        issue.assignee?.email?.toLowerCase().includes(search) ||
        issue.project?.name?.toLowerCase().includes(search),
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

const hasActiveFilters = computed(() => {
  return searchText.value || statusFilter.value || projectFilter.value
})

const rowSelection = {
  selectedRowKeys: selectedRowKeys,
  onChange: (keys) => {
    selectedRowKeys.value = keys
  },
}

const loadProjects = async () => {
  try {
    projects.value = await projectListAPI()
  } catch (error) {
    message.error('Failed to load projects')
    console.error(error)
  }
}

const loadEvents = async () => {
  try {
    events.value = await eventIssuesAPI()
  } catch (error) {
    message.error('Failed to load events')
    console.error(error)
  }
}

const renderTrendChart = (canvasId, trendData) => {
  if (!trendData || trendData.length === 0) return

  // Destroy existing chart if it exists
  if (chartInstances.value[canvasId]) {
    chartInstances.value[canvasId].destroy()
  }

  nextTick(() => {
    const canvas = document.getElementById(canvasId)
    if (!canvas) return

    const ctx = canvas.getContext('2d')

    const labels = trendData.map((d) => {
      const date = new Date(d.day)
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    })
    const counts = trendData.map((d) => d.count)

    try {
      chartInstances.value[canvasId] = new Chart(ctx, {
        type: 'line',
        data: {
          labels: labels,
          datasets: [
            {
              label: 'Events',
              data: counts,
              borderColor: '#51bc8f',
              backgroundColor: 'rgba(81, 188, 143, 0.1)',
              borderWidth: 2,
              fill: true,
              tension: 0.4,
              pointRadius: 0,
              pointHoverRadius: 3,
              pointHoverBackgroundColor: '#51bc8f',
              pointHoverBorderColor: '#fff',
              pointHoverBorderWidth: 1.5,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              enabled: true,
              mode: 'index',
              intersect: false,
              backgroundColor: 'rgba(0, 0, 0, 0.8)',
              padding: 6,
              bodyFont: {
                size: 10,
              },
              titleFont: {
                size: 10,
              },
              boxPadding: 3,
              callbacks: {
                label: function (context) {
                  return `${context.parsed.y} event${context.parsed.y !== 1 ? 's' : ''}`
                },
              },
            },
          },
          scales: {
            x: {
              display: false,
            },
            y: {
              display: false,
              beginAtZero: true,
            },
          },
          interaction: {
            intersect: false,
            mode: 'index',
          },
        },
      })
    } catch (error) {
      console.error('Error rendering chart:', error)
    }
  })
}

const renderAllCharts = () => {
  nextTick(() => {
    filteredIssues.value.forEach((issue) => {
      if (issue.trend && issue.trend.length > 0) {
        renderTrendChart(`trend-${issue.id}`, issue.trend)
      }
    })
  })
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

const handleDelete = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('Please select at least one issue')
    return
  }
  message.success(`Deleted ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}

const handleExport = () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('Please select at least one issue')
    return
  }
  message.success(`Exported ${selectedRowKeys.value.length} issue(s)`)
}

const copyIssueLink = (issueId) => {
  const link = `${window.location.origin}/issues/${issueId}`
  navigator.clipboard
    .writeText(link)
    .then(() => {
      message.success('Issue link copied to clipboard')
    })
    .catch(() => {
      message.error('Failed to copy link')
    })
}

const clearFilters = () => {
  searchText.value = ''
  statusFilter.value = null
  projectFilter.value = null
}

const formatDate = (dateString) => {
  if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'Just now'
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

// Watch for filter changes to re-render charts
watch(filteredIssues, () => {
  renderAllCharts()
})

onMounted(async () => {
  await loadProjects()
  await loadEvents()
  renderAllCharts()
})
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
        :data-source="filteredIssues"
        :row-selection="rowSelection"
        :row-key="(record) => record.id"
      >
        <template #title>
          <div v-if="selectedRowKeys.length > 0" class="table-header-actions">
            <span class="selected-count">{{ selectedRowKeys.length }} selected</span>
            <a-space>
              <a-button type="primary" @click="handleResolve">
                <template #icon><CheckCircleOutlined /></template>
                Resolve
              </a-button>
              <a-button @click="handleArchive">
                <template #icon><InboxOutlined /></template>
                Archive
              </a-button>
              <a-dropdown>
                <a-button>
                  <template #icon><MoreOutlined /></template>
                  More actions
                </a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item key="delete" @click="handleDelete">
                      <DeleteOutlined />
                      Delete
                    </a-menu-item>
                    <a-menu-item key="export" @click="handleExport">
                      <ExportOutlined />
                      Export
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
                <template #prefix>
                  <Search :size="16" />
                </template>
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
                <template #icon><FilterOutlined /></template>
                Clear
              </a-button>
            </div>
          </div>
        </template>

        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'summary'">
            <div class="issue-cell">
              <div class="issue-content">
                <a-typography-link class="issue-title" href="#">
                  {{ record.summary }}
                </a-typography-link>
                <div class="issue-meta">
                  <span class="project-name">{{ record.project?.name || 'Unknown' }}</span>
                  <span class="meta-separator">|</span>
                  <span class="issue-type-text">{{ record.type }}</span>
                </div>
                <div v-if="record.stacktrace" class="stacktrace-preview">
                  <span class="stacktrace-function">{{ record.stacktrace.function }}</span>
                  <span class="stacktrace-separator">in</span>
                  <span class="stacktrace-file">{{ record.stacktrace.file }}</span>
                  <span class="stacktrace-separator">at line</span>
                  <span class="stacktrace-line">{{ record.stacktrace.line }}</span>
                </div>
              </div>
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
                }"
              >
                {{ record.assignee.firstName?.[0] }}
              </a-avatar>
              <span>{{ record.assignee.firstName }}</span>
            </div>
            <a-typography-text type="secondary" v-else>Unassigned</a-typography-text>
          </template>

          <template v-if="column.key === 'trend'">
            <div v-if="record.trend && record.trend.length > 0" class="trend-chart-container">
              <canvas :id="`trend-${record.id}`" class="trend-chart"></canvas>
            </div>
            <div v-else class="trend-placeholder">
              <TrendingUp :size="16" style="color: #d9d9d9" />
            </div>
          </template>

          <template v-if="column.key === 'lastSeenAt'">
            <a-typography-text type="secondary">
              {{ formatDate(record.lastSeenAt) }}
            </a-typography-text>
          </template>

          <template v-if="column.key === 'actions'">
            <a-tooltip title="Copy issue link">
              <a-button type="text" size="small" @click.stop="copyIssueLink(record.id)">
                <template #icon>
                  <ExportOutlined style="font-size: 14px" />
                </template>
              </a-button>
            </a-tooltip>
          </template>
        </template>
      </a-table>
    </div>
  </main>
</template>

<style scoped>
.issue-tracker {
  padding: 24px;
  background: #fff;
}

.header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.table-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #e6f7ff;
  border-bottom: 1px solid #91d5ff;
}

.table-header-filters {
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
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
