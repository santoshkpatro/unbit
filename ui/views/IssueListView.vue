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
  SoundOutlined,
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
const projects = ref([])
const events = ref([])
const selectedRowKeys = ref([])
const searchText = ref('')
const statusFilter = ref(null)
const projectFilter = ref(null)
const chartInstances = ref({})

// ---- Theme helpers (use Ant Design Vue CSS variable or fallback to your new primary) ----
const getPrimaryColor = () => {
  const cssVar = getComputedStyle(document.documentElement)
    .getPropertyValue('--ant-color-primary')
    ?.trim()
  return cssVar || '#E1306C'
}

const toRgba = (hexOrRgb, alpha = 1) => {
  const clamp = (n) => Math.max(0, Math.min(255, n))
  if (!hexOrRgb) return `rgba(225, 48, 108, ${alpha})` // fallback #E1306C
  const s = hexOrRgb.replace(/\s+/g, '')
  if (s.startsWith('rgb')) {
    const nums = s
      .replace(/rgba?\(/, '')
      .replace(/\)/, '')
      .split(',')
      .slice(0, 3)
      .map((v) => clamp(parseInt(v, 10) || 0))
    return `rgba(${nums[0]}, ${nums[1]}, ${nums[2]}, ${alpha})`
  }
  // hex (#rgb or #rrggbb)
  let r = 225,
    g = 48,
    b = 108 // default #E1306C
  const m = s.match(/^#?([a-fA-F0-9]{3}|[a-fA-F0-9]{6})$/)
  if (m) {
    const h = m[1]
    if (h.length === 3) {
      r = parseInt(h[0] + h[0], 16)
      g = parseInt(h[1] + h[1], 16)
      b = parseInt(h[2] + h[2], 16)
    } else {
      r = parseInt(h.slice(0, 2), 16)
      g = parseInt(h.slice(2, 4), 16)
      b = parseInt(h.slice(4, 6), 16)
    }
  }
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

// Columns: "Last seen" inline in Issue column; Assignee is a dropdown; Trend is wider
const columns = [
  { title: 'Issue', dataIndex: ['event', 'message'], key: 'message', width: '45%' },
  {
    title: 'Events',
    dataIndex: ['group', 'eventCount'],
    key: 'count',
    width: '10%',
    align: 'center',
  },
  { title: 'Assignee', dataIndex: ['group', 'assignee'], key: 'assignee', width: '14%' },
  { title: 'Trend', dataIndex: 'last7Days', key: 'trend', width: '21%', align: 'center' },
  { title: '', key: 'actions', width: '10%', align: 'center' },
]

// Filters/search on raw events
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

  if (statusFilter.value) list = list.filter((row) => row.group?.status === statusFilter.value)
  if (projectFilter.value) list = list.filter((row) => row.project?.id === projectFilter.value)

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
    /* non-fatal */
  }
}

const loadEvents = async () => {
  try {
    events.value = await eventIssuesAPI()
  } catch {
    message.error('Failed to load events')
  }
}

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

/** Trend charts — now use your primary color */
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

    const primary = getPrimaryColor()

    chartInstances.value[canvasId] = new Chart(ctx, {
      type: 'line',
      data: {
        labels,
        datasets: [
          {
            label: 'Events',
            data: counts,
            borderColor: primary,
            backgroundColor: toRgba(primary, 0.12),
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
        elements: { line: { capBezierPoints: true } },
      },
    })
  })
}

const renderAllCharts = () => {
  nextTick(() => {
    filteredEvents.value.forEach((row) => {
      if (row.last7Days?.length) renderTrendChart(`trend-${row.event?.id}`, row.last7Days)
    })
  })
}

watch(filteredEvents, () => renderAllCharts())

onMounted(async () => {
  await loadProjects()
  await loadEvents()
  renderAllCharts()
})

/** bulk actions (demo toasts only) */
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
const handleMute = () => {
  if (!selectedRowKeys.value.length) return message.warning('Please select at least one issue')
  message.success(`Muted ${selectedRowKeys.value.length} issue(s)`)
  selectedRowKeys.value = []
}
</script>

<template>
  <main class="issue-tracker">
    <!-- Fixed header with actions on the right -->
    <div class="header">
      <div class="header-left">
        <a-typography-title :level="3" style="margin: 0">Issues</a-typography-title>
        <a-badge :count="unresolvedCount" :number-style="{ backgroundColor: '#f5222d' }" />
      </div>

      <div class="header-actions">
        <a-space>
          <a-button type="primary" :disabled="!selectedRowKeys.length" @click="handleResolve">
            <template #icon><CheckCircleOutlined /></template>Resolve
          </a-button>
          <a-button :disabled="!selectedRowKeys.length" @click="handleArchive">
            <template #icon><InboxOutlined /></template>Archive
          </a-button>
          <a-button :disabled="!selectedRowKeys.length" @click="handleMute">
            <template #icon><SoundOutlined /></template>Mute
          </a-button>
          <a-dropdown>
            <a-button :disabled="!selectedRowKeys.length">
              <template #icon><MoreOutlined /></template>More actions
            </a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item key="delete" @click="handleDelete" :disabled="!selectedRowKeys.length">
                  <DeleteOutlined />Delete
                </a-menu-item>
                <a-menu-item key="export" @click="handleExport" :disabled="!selectedRowKeys.length">
                  <ExportOutlined />Export
                </a-menu-item>
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
        :data-source="filteredEvents"
        :row-selection="rowSelection"
        :row-key="(record) => record.event?.id"
      >
        <template #title>
          <!-- Filters only -->
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
          <!-- Issue cell: title + meta; includes Last Seen inline -->
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
                  <span class="meta-separator">|</span>
                  <span class="last-seen"
                    >Last seen {{ formatTimestamp(record.event?.timestamp) }}</span
                  >
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

          <!-- assignee: dropdown (empty options for now) -->
          <template v-if="column.key === 'assignee'">
            <a-select
              :options="[]"
              placeholder="Select assignee"
              style="width: 160px"
              allow-clear
            />
          </template>

          <!-- trend (LINE spark) -->
          <template v-if="column.key === 'trend'">
            <div v-if="record.last7Days?.length" class="trend-chart-container">
              <canvas :id="`trend-${record.event?.id}`" class="trend-chart"></canvas>
            </div>
            <div v-else class="trend-placeholder">
              <TrendingUp :size="16" style="color: #d9d9d9" />
            </div>
          </template>

          <!-- row actions: dropdown only; copy moved inside menu -->
          <template v-if="column.key === 'actions'">
            <a-dropdown placement="bottomRight" :trigger="['click']">
              <a-button type="text" size="small" aria-label="More actions">
                <template #icon><MoreOutlined style="font-size: 16px" /></template>
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="copy" @click="copyIssueLink(record.event?.id)">
                    <ExportOutlined /> Copy Issue Link
                  </a-menu-item>
                  <a-menu-item key="jira"> <MoreOutlined /> Create Jira Ticket </a-menu-item>
                  <a-menu-item key="email"> <MoreOutlined /> Send Email </a-menu-item>
                  <a-menu-item key="custom"> <MoreOutlined /> Custom Action </a-menu-item>
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
/* page fills the app-content area */
.issue-tracker {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
  background: #fff;
}

/* Fixed header */
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

.header-actions :deep(.ant-btn + .ant-btn) {
  margin-left: 8px;
}

/* inner scroller (everything below header) */
.table-container {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  padding: 0 24px 24px;
  border-radius: 0;
}

/* sticky toolbar inside scroller (filters) */
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
  color: var(--ant-color-primary);
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
  color: var(--ant-color-primary);
}

.issue-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  flex-wrap: wrap;
}

.last-seen {
  color: #8c8c8c;
}

.project-name {
  color: var(--ant-color-primary);
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
  color: var(--ant-color-primary);
  font-weight: 500;
}

.assignee-cell {
  display: flex;
  align-items: center;
}

/* trend chart (line) */
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
:deep(.ant-table-tbody > tr:hover) {
  cursor: pointer;
}
:deep(.ant-table-title) {
  padding: 0 !important;
}
</style>
