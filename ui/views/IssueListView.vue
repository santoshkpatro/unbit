<script setup>
import { onMounted, ref, watch, nextTick, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { recentIssuesAPI } from '@/api/issues'
import { projectListAPI } from '@/api/projects'
import { Search, TrendingUp, MoreVertical } from 'lucide-vue-next'
import { message } from 'ant-design-vue'
import { Chart, registerables } from 'chart.js'
Chart.register(...registerables)

const router = useRouter()

const issues = ref([])
const projects = ref([])
const selectedRowKeys = ref([])
const searchText = ref('')
const projectFilter = ref(null)
const chartInstances = ref({})

const columns = [
  { title: 'Issue', dataIndex: 'message', key: 'message', width: '38%' },
  { title: 'Events', dataIndex: 'eventCount', key: 'count', width: '10%', align: 'center' },
  { title: 'Age', dataIndex: 'age', key: 'age', width: '8%', align: 'center' },
  { title: 'Assignee', dataIndex: 'assignee', key: 'assignee', width: '14%' },
  { title: 'Trend', dataIndex: 'issueCountReport', key: 'trend', width: '20%', align: 'center' },
  { title: '', key: 'actions', width: '10%', align: 'center' },
]

const totalIssues = computed(() => issues.value.length)

const rowSelection = {
  selectedRowKeys,
  onChange: (keys) => (selectedRowKeys.value = keys),
}

const loadProjects = async () => {
  projects.value = await projectListAPI()
}

const loadIssues = async () => {
  const params = {}

  if (projectFilter.value) {
    params.project_id = projectFilter.value
  }

  issues.value = await recentIssuesAPI(params)
}

const copyIssueLink = (id) => {
  const link = `${window.location.origin}/issues/${id}`
  navigator.clipboard
    .writeText(link)
    .then(() => message.success('Issue link copied'))
    .catch(() => message.error('Copy failed'))
}

const clearFilters = () => {
  searchText.value = ''
  projectFilter.value = null
  loadIssues()
}

watch(projectFilter, () => {
  loadIssues()
})

const formatTimestamp = (ts) => {
  if (!ts || ts === '0001-01-01T00:00:00Z') return 'Just now'
  const d = new Date(ts)
  const n = new Date()
  const diff = n - d
  const m = Math.floor(diff / 60000)
  const h = Math.floor(diff / 3600000)
  const day = Math.floor(diff / 86400000)
  if (m < 1) return 'Just now'
  if (m < 60) return `${m}m ago`
  if (h < 24) return `${h}h ago`
  if (day < 7) return `${day}d ago`
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: d.getFullYear() !== n.getFullYear() ? 'numeric' : undefined,
  })
}

const formatAge = (sec) => {
  const s = Number(sec)
  if (!Number.isFinite(s) || s < 0) return '-'
  const days = Math.floor(s / 86400)
  const hours = Math.floor((s % 86400) / 3600)
  const minutes = Math.floor((s % 3600) / 60)
  const seconds = Math.floor(s % 60)
  const parts = []
  if (days) parts.push(`${days}d`)
  if (hours) parts.push(`${hours}h`)
  if (minutes && parts.length < 2) parts.push(`${minutes}m`)
  if (!parts.length) parts.push(`${seconds}s`)
  return parts.slice(0, 2).join(' ')
}

const renderTrendChart = (canvasId, issueCountReport) => {
  if (!issueCountReport?.length) return
  if (chartInstances.value[canvasId]) chartInstances.value[canvasId].destroy()

  nextTick(() => {
    const el = document.getElementById(canvasId)
    if (!el) return
    const ctx = el.getContext('2d')
    const labels = issueCountReport.map((d) =>
      new Date(d.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
    )
    const data = issueCountReport.map((d) => d.eventCount || 0)
    const primary =
      getComputedStyle(document.documentElement).getPropertyValue('--ant-color-primary')?.trim() ||
      '#1677ff'

    chartInstances.value[canvasId] = new Chart(ctx, {
      type: 'line',
      data: {
        labels,
        datasets: [
          {
            label: 'Events',
            data,
            borderColor: primary,
            backgroundColor: 'rgba(22, 119, 255, 0.2)',
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
        plugins: { legend: { display: false }, tooltip: { enabled: true } },
        scales: { x: { display: false }, y: { display: false, beginAtZero: true } },
      },
    })
  })
}

const renderAllCharts = () =>
  nextTick(() =>
    issues.value.forEach(
      (r) => r.issueCountReport?.length && renderTrendChart(`trend-${r.id}`, r.issueCountReport),
    ),
  )

watch(issues, renderAllCharts)

onMounted(async () => {
  await loadProjects()
  await loadIssues()
})

const requireSelection = (fnName) => {
  if (!selectedRowKeys.value.length) {
    message.warning('Please select at least one issue')
    return false
  }
  return true
}

const handleResolve = () =>
  requireSelection() &&
  (message.success(`Resolved ${selectedRowKeys.value.length} issue(s)`),
  (selectedRowKeys.value = []))

const handleArchive = () =>
  requireSelection() &&
  (message.success(`Archived ${selectedRowKeys.value.length} issue(s)`),
  (selectedRowKeys.value = []))

const handleDelete = () =>
  requireSelection() &&
  (message.success(`Deleted ${selectedRowKeys.value.length} issue(s)`),
  (selectedRowKeys.value = []))

const handleExport = () =>
  requireSelection() && message.success(`Exported ${selectedRowKeys.value.length} issue(s)`)

const handleMute = () =>
  requireSelection() &&
  (message.success(`Muted ${selectedRowKeys.value.length} issue(s)`), (selectedRowKeys.value = []))
</script>

<template>
  <main class="issue-tracker">
    <div class="header">
      <div class="header-left">
        <a-typography-title :level="3" style="margin: 0">Issues</a-typography-title>
        <a-badge :count="totalIssues" :number-style="{ backgroundColor: '#2f54eb' }" />
      </div>
      <div class="header-actions">
        <a-space>
          <a-button type="primary" :disabled="!selectedRowKeys.length" @click="handleResolve"
            >Resolve</a-button
          >
          <a-button :disabled="!selectedRowKeys.length" @click="handleArchive">Archive</a-button>
          <a-button :disabled="!selectedRowKeys.length" @click="handleMute">Mute</a-button>
          <a-dropdown>
            <a-button :disabled="!selectedRowKeys.length">More actions</a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item key="delete" :disabled="!selectedRowKeys.length" @click="handleDelete"
                  >Delete</a-menu-item
                >
                <a-menu-item key="export" :disabled="!selectedRowKeys.length" @click="handleExport"
                  >Export</a-menu-item
                >
              </a-menu>
            </template>
          </a-dropdown>
        </a-space>
      </div>
    </div>

    <div class="table-container">
      <a-table
        size="small"
        :pagination="false"
        :columns="columns"
        :data-source="issues"
        :row-selection="rowSelection"
        :row-key="(r) => r.id"
      >
        <template #title>
          <div class="table-header-filters">
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
                v-model:value="projectFilter"
                style="width: 180px"
                placeholder="All projects"
                allow-clear
              >
                <a-select-option :value="null">All projects</a-select-option>
                <a-select-option v-for="p in projects" :key="p.id" :value="p.id">
                  {{ p.name }}
                </a-select-option>
              </a-select>
              <a-button v-if="hasActiveFilters" @click="clearFilters">Clear</a-button>
            </div>
          </div>
        </template>

        <template #bodyCell="{ column, record }">
          <!-- Issue -->
          <template v-if="column.key === 'message'">
            <div class="issue-cell">
              <div class="issue-content">
                <router-link
                  class="issue-title"
                  :to="{ name: 'issue-details', params: { issueId: record.id } }"
                >
                  {{ record.message }}
                </router-link>

                <div class="issue-meta">
                  <a-typography-text strong>
                    {{ record.project?.name || 'Unknown' }}
                  </a-typography-text>
                  <span class="meta-separator">|</span>
                  <span class="issue-type-text">{{ record.type }}</span>
                  <span class="meta-separator">|</span>
                  <span class="last-seen">Last seen {{ formatTimestamp(record.timestamp) }}</span>
                </div>

                <div v-if="record.firstStackTrace" class="stacktrace-preview">
                  <span class="stacktrace-function">{{ record.firstStackTrace.function }}</span>
                  <span class="stacktrace-separator">in</span>
                  <span class="stacktrace-file">{{ record.firstStackTrace.file }}</span>
                  <span class="stacktrace-separator">at line</span>
                  <span class="stacktrace-line">{{ record.firstStackTrace.line }}</span>
                </div>
              </div>
            </div>
          </template>

          <!-- Count -->
          <template v-if="column.key === 'count'">
            <a-badge
              :count="record.eventCount || 0"
              :overflow-count="999"
              :number-style="{ backgroundColor: '#f0f0f0', color: '#595959', fontWeight: 600 }"
            />
          </template>

          <!-- Age (formatted from seconds) -->
          <template v-if="column.key === 'age'">
            <a-typography-text>{{ formatAge(record.age) }}</a-typography-text>
          </template>

          <!-- Assignee -->
          <template v-if="column.key === 'assignee'">
            <a-select
              :options="[]"
              :value="
                record.assignee?.email || record.assignee?.name || record.assignee || undefined
              "
              placeholder="Select assignee"
              style="width: 160px"
              allow-clear
            />
          </template>

          <!-- Trend -->
          <template v-if="column.key === 'trend'">
            <div v-if="record.issueCountReport?.length" class="trend-chart-container">
              <canvas :id="`trend-${record.id}`" class="trend-chart"></canvas>
            </div>
            <div v-else class="trend-placeholder">
              <TrendingUp :size="16" style="color: #d9d9d9" />
            </div>
          </template>

          <!-- Actions -->
          <template v-if="column.key === 'actions'">
            <a-dropdown placement="bottomRight" :trigger="['click']">
              <a-button type="text" size="small" shape="circle" aria-label="More actions">
                <MoreVertical :size="16" />
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="copy" @click="copyIssueLink(record.id)"
                    >Copy Issue Link</a-menu-item
                  >
                  <a-menu-item key="jira">Create Jira Ticket</a-menu-item>
                  <a-menu-item key="email">Send Email</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </template>
        </template>
      </a-table>
    </div>
  </main>
</template>

<style scoped>
.issue-tracker {
  height: 100%;
  min-height: 98vh;
  display: flex;
  flex-direction: column;
  background: #fff;
}
.header {
  position: sticky;
  top: 0;
  z-index: 3;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* make the header count badge use your brand color */
:deep(.header-left .ant-badge .ant-badge-count) {
  background: #e1306c !important;
}

.table-container {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: 0 24px 24px;
  display: flex;
  flex-direction: column;
}

/* Make table always take the available height */
:deep(.ant-table-wrapper),
:deep(.ant-spin-nested-loading),
:deep(.ant-spin-container),
:deep(.ant-table),
:deep(.ant-table-container) {
  height: 100%;
}
:deep(.ant-table) {
  display: flex;
  flex-direction: column;
}
:deep(.ant-table-container) {
  flex: 1 1 auto;
}
:deep(.ant-table-body),
:deep(.ant-table-content) {
  max-height: none !important;
}

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
  font-weight: 600;
  margin-bottom: 4px;
  color: #1f1f1f;
  transition:
    color 0.15s ease,
    transform 0.15s ease,
    text-decoration-color 0.15s ease;
  text-decoration: none;
}
.issue-title:hover {
  color: #e1306c;
  text-decoration: underline;
  text-underline-offset: 2px;
  transform: translateY(-1px);
}

.issue-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  flex-wrap: wrap;
}

.meta-separator {
  color: #d9d9d9;
  margin: 0 2px;
}
.issue-type-text {
  color: #8c8c8c;
  text-transform: capitalize;
}
.last-seen {
  color: #8c8c8c;
}

.stacktrace-preview {
  margin-top: 6px;
  font-size: 11px;
  color: #8c8c8c;
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
  line-height: 1.4;
  padding: 4px 8px;
  background: #fafafa;
  border-radius: 4px;
  border-left: 2px solid #e8e8e8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.stacktrace-function {
  color: #722ed1;
  font-weight: 600;
}
.stacktrace-separator {
  color: #bfbfbf;
  margin: 0 4px;
}
.stacktrace-file {
  color: #595959;
}
.stacktrace-line {
  color: #e1306c;
  font-weight: 600;
}

/* charts */
.trend-chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 0;
}
.trend-chart {
  width: 220px;
  height: 48px;
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
:deep(.ant-table-title) {
  padding: 0 !important;
}
</style>
