<script setup>
import { onMounted, ref, computed } from 'vue'
import { projectListAPI } from '@/api/projects'
import { issueListAPI } from '@/api/issues'
import { Search } from 'lucide-vue-next'
import { CheckCircleOutlined, InboxOutlined, FilterOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

const projects = ref([])
const issues = ref([])
const selectedRowKeys = ref([])
const searchText = ref('')
const statusFilter = ref(null)
const projectFilter = ref(null)

const columns = [
  {
    title: 'Summary',
    dataIndex: 'summary',
    key: 'summary',
    ellipsis: true,
    width: '35%',
  },
  {
    title: 'Project',
    dataIndex: ['project', 'name'],
    key: 'project',
    width: '15%',
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    width: '12%',
  },
  {
    title: 'Assignee',
    dataIndex: ['assignee', 'firstName'],
    key: 'assignee',
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
    title: 'Last Seen',
    dataIndex: 'lastSeenAt',
    key: 'lastSeenAt',
    width: '13%',
  },
]

const filteredIssues = computed(() => {
  let result = issues.value

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

const loadIssues = async () => {
  issues.value = await issueListAPI()
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
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(async () => {
  await loadProjects()
  await loadIssues()
})
</script>

<template>
  <main class="issue-tracker">
    <a-typography-title :level="3">Issues</a-typography-title>

    <div class="toolbar">
      <div class="filters">
        <a-input
          v-model:value="searchText"
          placeholder="Search issues..."
          style="width: 300px"
          allow-clear
        >
          <template #prefix>
            <Search :size="16" />
          </template>
        </a-input>

        <a-select
          v-model:value="statusFilter"
          placeholder="Status"
          style="width: 150px"
          allow-clear
        >
          <a-select-option value="unresolved">Unresolved</a-select-option>
          <a-select-option value="resolved">Resolved</a-select-option>
          <a-select-option value="archived">Archived</a-select-option>
        </a-select>

        <a-select
          v-model:value="projectFilter"
          placeholder="Project"
          style="width: 200px"
          allow-clear
          :options="projectOptions"
        />

        <a-button @click="clearFilters">
          <template #icon><FilterOutlined /></template>
          Clear Filters
        </a-button>
      </div>

      <div class="actions" v-if="selectedRowKeys.length > 0">
        <a-badge :count="selectedRowKeys.length" :offset="[-5, 5]">
          <a-space>
            <a-button type="primary" @click="handleResolve">
              <template #icon><CheckCircleOutlined /></template>
              Resolve
            </a-button>
            <a-button @click="handleArchive">
              <template #icon><InboxOutlined /></template>
              Archive
            </a-button>
          </a-space>
        </a-badge>
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
        showTotal: (total) => `Total ${total} issues`,
      }"
      :scroll="{ x: 1200 }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'summary'">
          <div class="summary-cell">
            <a-typography-text :ellipsis="{ tooltip: record.summary }" strong>
              {{ record.summary }}
            </a-typography-text>
          </div>
        </template>

        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 'unresolved' ? 'red' : 'green'">
            {{ record.status?.toUpperCase() }}
          </a-tag>
        </template>

        <template v-if="column.key === 'assignee'">
          <div v-if="record.assignee" class="assignee-cell">
            <a-avatar size="small" :style="{ backgroundColor: '#1890ff', marginRight: '8px' }">
              {{ record.assignee.firstName?.[0] }}
            </a-avatar>
            <span>{{ record.assignee.firstName }}</span>
          </div>
          <span v-else>-</span>
        </template>

        <template v-if="column.key === 'eventCount'">
          <a-typography-text type="secondary">
            {{ record.eventCount }}
          </a-typography-text>
        </template>

        <template v-if="column.key === 'lastSeenAt'">
          <a-typography-text type="secondary">
            {{ formatDate(record.lastSeenAt) }}
          </a-typography-text>
        </template>
      </template>
    </a-table>
  </main>
</template>

<style scoped>
.issue-tracker {
  padding: 24px;
  background: #fff;
  height: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  flex: 1;
}

.actions {
  display: flex;
  align-items: center;
}

.summary-cell {
  max-width: 100%;
}

.assignee-cell {
  display: flex;
  align-items: center;
}

:deep(.ant-table) {
  font-size: 14px;
}

:deep(.ant-table-thead > tr > th) {
  background: #fafafa;
  font-weight: 600;
}

:deep(.ant-table-tbody > tr:hover) {
  cursor: pointer;
}

:deep(.ant-table-cell) {
  padding: 12px 16px;
}
</style>
