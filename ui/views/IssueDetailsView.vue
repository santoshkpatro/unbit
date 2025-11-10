<script setup>
import { onMounted, ref, computed, h, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { issueDetailsAPI } from '@/api/issues.js'

// Ant icons — used ONLY on buttons
import {
  CopyOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  SendOutlined,
  UserOutlined,
} from '@ant-design/icons-vue'
import { message as antdMessage } from 'ant-design-vue'

// Lucide icons — section/visual icons
import {
  FileCode2 as LFileCode,
  Clock as LClock,
  Folder as LFolder,
  User as LUser,
  AlertTriangle as LAlert,
  Server as LServer,
  Cpu as LCpu,
  TerminalSquare as LTerminal,
  Network as LNetwork,
} from 'lucide-vue-next'

dayjs.extend(relativeTime)

/** ---- State ---- **/
const route = useRoute()
const loading = ref(true)
const error = ref(null)
const issue = ref(null)
const resolving = ref(false)

// NEW: Hardcoded previous events (placeholder)
const previousEvents = ref([
  {
    id: 'evt_1001',
    message: 'TypeError: Cannot read properties of undefined (reading "foo")',
    timestamp: dayjs().subtract(2, 'day').toISOString(),
    level: 'error',
    status: 'unresolved',
  },
  {
    id: 'evt_1000',
    message: 'ReferenceError: bar is not defined',
    timestamp: dayjs().subtract(5, 'day').toISOString(),
    level: 'error',
    status: 'resolved',
  },
  {
    id: 'evt_0999',
    message: 'Warning: setState on unmounted component',
    timestamp: dayjs().subtract(8, 'day').toISOString(),
    level: 'warning',
    status: 'unresolved',
  },
])

onMounted(async () => {
  try {
    loading.value = true
    const id = route.params.issueId
    issue.value = await issueDetailsAPI(id)
  } catch (e) {
    error.value = e?.message || 'Failed to load issue'
  } finally {
    loading.value = false
  }
})

/** ---- UI computed ---- **/
const titleText = computed(() => issue.value?.message || 'Issue Detail')

const statusBadge = computed(() => {
  const s = (issue.value?.status || '').toLowerCase()
  if (s === 'unresolved') return { status: 'error', text: 'Unresolved' }
  if (s === 'resolved') return { status: 'success', text: 'Resolved' }
  return { status: 'warning', text: issue.value?.status || 'Unknown' }
})

const levelTag = computed(() => {
  const lvl = (issue.value?.level || '').toLowerCase()
  if (lvl === 'error') return { color: 'red', text: 'Error' }
  if (lvl === 'warning') return { color: 'orange', text: 'Warning' }
  if (lvl === 'info') return { color: 'blue', text: 'Info' }
  return { color: 'default', text: lvl || 'Unknown' }
})

const typeTag = computed(() => ({ color: 'purple', text: issue.value?.type || 'exception' }))

const occurredAt = computed(() => {
  const ts = issue.value?.timestamp
  if (!ts) return { abs: '-', rel: '-' }
  const d = dayjs(ts)
  return { abs: d.format('YYYY-MM-DD HH:mm:ss'), rel: d.fromNow() }
})

const ageText = computed(() => {
  const sec = issue.value?.age
  if (typeof sec !== 'number') return '-'
  return dayjs().subtract(sec, 'second').fromNow(true)
})

const stack = computed(() => issue.value?.stacktrace || [])

/** ---- Environment & Context ---- **/
const runtimeInfo = computed(() => issue.value?.runtime || {})
const osInfo = computed(() => issue.value?.os || {})
const processInfo = computed(() => issue.value?.process || {})
const threadInfo = computed(() => issue.value?.thread || {})
const hostInfo = computed(() => issue.value?.host || {})

const runtimeSummary = computed(() => {
  const r = runtimeInfo.value
  if (!r || Object.keys(r).length === 0) return '-'
  const parts = [r.name, r.implementation].filter(Boolean)
  const ver = r.version ? ` ${r.version}` : ''
  return parts.join(' · ') + ver
})

const osSummary = computed(() => {
  const o = osInfo.value
  if (!o || Object.keys(o).length === 0) return '-'
  return [o.name, o.platform, o.machine].filter(Boolean).join(' · ')
})

const processSummary = computed(() => {
  const p = processInfo.value
  if (!p || Object.keys(p).length === 0) return '-'
  const pid = p.pid != null ? `pid ${p.pid}` : null
  const ppid = p.ppid != null ? `ppid ${p.ppid}` : null
  return [pid, ppid].filter(Boolean).join(' · ')
})

const threadSummary = computed(() => {
  const t = threadInfo.value
  if (!t || Object.keys(t).length === 0) return '-'
  const name = t.name ? `“${t.name}”` : null
  const ident = t.ident != null ? `id ${t.ident}` : null
  return [name, ident].filter(Boolean).join(' · ')
})

const hostSummary = computed(() => hostInfo.value?.hostname || '-')

/** ---- Actions ---- **/
const resolveIssue = () => {
  if (!issue.value) return
  resolving.value = true
  setTimeout(() => {
    issue.value.status = 'resolved'
    resolving.value = false
    antdMessage.success('Issue marked as resolved')
  }, 200)
}

// Header assignee dropdown
const assignee = ref('unassigned')
watchEffect(() => {
  if (!issue.value?.assignee) {
    assignee.value = 'unassigned'
  } else if (issue.value.assignee?.name === 'You') {
    assignee.value = 'me'
  } else {
    assignee.value = 'unassigned'
  }
})
const onAssigneeChange = (val) => {
  if (!issue.value) return
  if (val === 'me') {
    issue.value.assignee = { name: 'You', email: '' }
    antdMessage.success('Assigned to you')
  } else {
    issue.value.assignee = null
    antdMessage.success('Unassigned')
  }
}

const assignToMe = () => {
  onAssigneeChange('me')
}

const copyText = async (text, label = 'Copied') => {
  try {
    await navigator.clipboard.writeText(String(text ?? ''))
    antdMessage.success(label)
  } catch {
    antdMessage.error('Copy failed')
  }
}

/** ---- Snippet rendering ---- */
const hasHtmlSnippet = (frame) => typeof frame?.snippetHtml === 'string' && frame.snippetHtml.trim()
const hasTextSnippet = (frame) => typeof frame?.snippetText === 'string' && frame.snippetText.trim()

// Build a copy/share friendly text block (no HTML) from best-available source
const plainSnippetFrom = (frame) => {
  if (hasTextSnippet(frame)) return frame.snippetText
  if (typeof frame?.code === 'string' && frame.code.trim()) return frame.code
  // fallback minimal
  return '/* no code available */'
}

const shareableSnippet = computed(() => {
  if (!stack.value.length) return ''
  const first = stack.value[0]
  const header = `${first.file || '<unknown>'}:${first.line ?? '?'} in ${first.function || '<module>'}`
  const body = plainSnippetFrom(first)
  return `# ${issue.value?.message || 'Exception'}\n${header}\n\n\`\`\`\n${body}\n\`\`\``
})

const shareableEnv = computed(() => {
  const payload = {
    id: issue.value?.id,
    eventId: issue.value?.eventId,
    timestamp: issue.value?.timestamp,
    runtime: runtimeInfo.value,
    os: osInfo.value,
    process: processInfo.value,
    thread: threadInfo.value,
    host: hostInfo.value,
  }
  return JSON.stringify(payload, null, 2)
})
</script>

<template>
  <main class="issue-details">
    <a-card v-if="issue" class="header-card" :bordered="false">
      <div class="header">
        <div class="header-left">
          <div class="title-line">
            <div class="title">{{ titleText }}</div>
          </div>

          <div class="meta-line">
            <span v-if="issue.project?.name" class="mono">
              <LFolder class="icon" /> {{ issue.project?.name }}
            </span>
            <a-divider type="vertical" />
            <a-badge :status="statusBadge.status" :text="statusBadge.text" />
            <a-tag :color="levelTag.color">{{ levelTag.text }}</a-tag>
            <a-tag :color="typeTag.color">{{ typeTag.text }}</a-tag>
            <a-divider type="vertical" />
            <span class="mono">Event {{ issue.eventCount }}</span>
            <a-divider type="vertical" />
            <span class="mono">Last seen: {{ occurredAt.abs }} ({{ occurredAt.rel }})</span>
            <a-divider type="vertical" />
            <span class="mono"><LClock class="icon" /> Age: {{ ageText }}</span>
          </div>
        </div>

        <div class="header-right">
          <div class="actions-inline">
            <a-space :size="8" align="center" wrap>
              <a-button
                :icon="h(CopyOutlined)"
                @click="() => copyText(issue.id, 'Issue ID copied')"
              >
                Copy ID
              </a-button>
              <a-select
                size="small"
                style="min-width: 160px"
                :value="assignee"
                @change="onAssigneeChange"
              >
                <a-select-option value="unassigned">Unassigned</a-select-option>
                <a-select-option value="me">Assign to me</a-select-option>
              </a-select>
              <a-button :icon="h(UserOutlined)" @click="assignToMe">Assign to me</a-button>
              <a-button
                type="primary"
                :loading="resolving"
                :icon="h(CheckCircleOutlined)"
                @click="resolveIssue"
              >
                Resolve
              </a-button>
              <a-button :icon="h(ReloadOutlined)" @click="() => location.reload()"
                >Refresh</a-button
              >
              <a-button
                :icon="h(SendOutlined)"
                @click="() => copyText(shareableSnippet, 'Shareable snippet copied')"
              >
                Share
              </a-button>
            </a-space>
          </div>
        </div>
      </div>
    </a-card>

    <a-card v-if="error" class="mt">
      <a-alert type="error" :message="error" show-icon />
    </a-card>

    <a-skeleton v-if="loading" :active="true" :title="true" :paragraph="{ rows: 6 }" />

    <template v-else>
      <a-empty v-if="!issue" description="No issue found" />

      <div v-else class="mt">
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :md="16">
            <a-card :bordered="true">
              <a-tabs default-active-key="stack">
                <a-tab-pane key="stack" tab="Stactrace">
                  <a-collapse accordion>
                    <a-collapse-panel
                      v-for="(frame, i) in stack"
                      :key="i"
                      :header="`${frame.function || '<module>'} — ${frame.file || ''}:${frame.line ?? ''}`"
                    >
                      <div class="frame-head">
                        <LFileCode class="icon" />
                        <code class="mono">{{ frame.file }}</code>
                        <span class="weak">:</span>
                        <span class="mono">{{ frame.line }}</span>
                      </div>

                      <div
                        v-if="hasHtmlSnippet(frame)"
                        class="code-block"
                        v-html="frame.snippetHtml"
                      ></div>
                      <pre
                        v-else
                        class="code-block"
                      ><code>{{ plainSnippetFrom(frame) }}</code></pre>

                      <a-space class="mt-8" wrap>
                        <a-button
                          size="small"
                          :icon="h(CopyOutlined)"
                          @click="() => copyText(plainSnippetFrom(frame), 'Frame code copied')"
                        >
                          Copy frame code
                        </a-button>
                        <a-button
                          size="small"
                          :icon="h(CopyOutlined)"
                          @click="() => copyText(`${frame.file}:${frame.line}`, 'Location copied')"
                        >
                          Copy location
                        </a-button>
                        <a-button
                          size="small"
                          :icon="h(SendOutlined)"
                          @click="
                            () =>
                              copyText(
                                `\`\`\`\n${plainSnippetFrom(frame).trim()}\n\`\`\``,
                                'Snippet copied',
                              )
                          "
                        >
                          Share snippet
                        </a-button>
                      </a-space>
                    </a-collapse-panel>
                  </a-collapse>
                </a-tab-pane>

                <a-tab-pane key="events" tab="Previous Events">
                  <a-list :data-source="previousEvents" item-layout="horizontal" :split="true">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta>
                          <template #title>
                            <div class="event-title">
                              <span class="mono">{{ item.id }}</span>
                              <a-tag
                                :color="
                                  (item.level === 'error' && 'red') ||
                                  (item.level === 'warning' && 'orange') ||
                                  'blue'
                                "
                              >
                                {{ item.level }}
                              </a-tag>
                              <a-badge
                                :status="item.status === 'resolved' ? 'success' : 'error'"
                                :text="item.status"
                              />
                            </div>
                          </template>
                          <template #description>
                            <div class="event-desc">
                              <div class="truncate">{{ item.message }}</div>
                              <div class="weak mono">
                                {{ dayjs(item.timestamp).format('YYYY-MM-DD HH:mm:ss') }} ({{
                                  dayjs(item.timestamp).fromNow()
                                }})
                              </div>
                            </div>
                          </template>
                        </a-list-item-meta>
                        <template #actions>
                          <a @click="() => copyText(item.id, 'Event ID copied')">Copy ID</a>
                          <a-divider type="vertical" />
                          <a @click="() => copyText(item.message, 'Message copied')"
                            >Copy message</a
                          >
                        </template>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </a-card>

            <a-card class="mt" :bordered="true" :title="'Details'">
              <a-descriptions bordered :column="1" size="small">
                <a-descriptions-item label="Issue ID">{{ issue.id }}</a-descriptions-item>
                <a-descriptions-item label="Event ID">{{ issue.eventId }}</a-descriptions-item>
                <a-descriptions-item label="Type">{{ issue.type }}</a-descriptions-item>
                <a-descriptions-item label="Level">{{ issue.level }}</a-descriptions-item>
                <a-descriptions-item label="Status">{{ issue.status }}</a-descriptions-item>
                <a-descriptions-item label="Timestamp">
                  {{ occurredAt.abs }} ({{ occurredAt.rel }})
                </a-descriptions-item>
                <a-descriptions-item
                  v-if="issue.project?.id || issue.project?.name"
                  label="Project"
                >
                  {{ issue.project?.name }}
                  <template v-if="issue.project?.id">
                    <a-divider type="vertical" /> <span class="weak">ID:</span>
                    {{ issue.project?.id }}
                  </template>
                </a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>

          <a-col :xs="24" :md="8">
            <a-card title="Runtime & Environment">
              <a-descriptions :column="1" size="small" bordered>
                <a-descriptions-item label="Runtime">
                  <div class="kv">
                    <LCpu class="icon" />
                    <span>{{ runtimeSummary }}</span>
                  </div>
                </a-descriptions-item>
                <a-descriptions-item label="OS">
                  <div class="kv">
                    <LServer class="icon" />
                    <span>{{ osSummary }}</span>
                  </div>
                </a-descriptions-item>
                <a-descriptions-item label="Process">
                  <div class="kv">
                    <LTerminal class="icon" />
                    <span>{{ processSummary }}</span>
                  </div>
                </a-descriptions-item>
                <a-descriptions-item label="Thread">
                  <div class="kv">
                    <LTerminal class="icon" />
                    <span>{{ threadSummary }}</span>
                  </div>
                </a-descriptions-item>
                <a-descriptions-item label="Host">
                  <div class="kv">
                    <LNetwork class="icon" />
                    <span>{{ hostSummary }}</span>
                  </div>
                </a-descriptions-item>
              </a-descriptions>

              <a-space class="mt-8" wrap>
                <a-button
                  size="small"
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(runtimeSummary, 'Runtime copied')"
                >
                  Copy runtime
                </a-button>
                <a-button
                  size="small"
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(osSummary, 'OS copied')"
                >
                  Copy OS
                </a-button>
                <a-button
                  size="small"
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(processSummary, 'Process copied')"
                >
                  Copy process
                </a-button>
                <a-button
                  size="small"
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(threadSummary, 'Thread copied')"
                >
                  Copy thread
                </a-button>
                <a-button
                  size="small"
                  :icon="h(SendOutlined)"
                  @click="() => copyText(shareableEnv, 'Environment JSON copied')"
                >
                  Copy env JSON
                </a-button>
              </a-space>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </template>
  </main>
</template>

<style scoped>
/* Layout */
.header-card {
  background: linear-gradient(180deg, #ffffff, #fafafa);
  border-radius: 12px;
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.02);
}

/* This is now the main flex container */
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  flex-wrap: wrap; /* Actions will stack below on mobile */
}
.header-left {
  min-width: 0;
}

/* This is the new container for actions */
.header-right {
  flex-shrink: 0; /* Prevents it from being crushed */
}

.title {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.2;
  word-break: break-word;
}
.actions-inline {
  /* This rule is fine, it contains the a-space */
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

/* Give the meta line a bit more breathing room from the title */
.meta-line {
  color: rgba(0, 0, 0, 0.65);
  margin-top: 12px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

/* Reusable */
.icon {
  width: 16px;
  height: 16px;
  margin-right: 6px;
  vertical-align: -2px;
}
.icon-lg {
  width: 24px;
  height: 24px;
}
.meta {
  display: flex;
  align-items: center;
  gap: 8px;
}
.kv {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.mono {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
}
.mt {
  margin-top: 12px;
}
.mt-8 {
  margin-top: 8px;
}
.strong {
  font-weight: 600;
}
.weak {
  color: rgba(0, 0, 0, 0.45);
}

/* Frame header */
.frame-head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

/* Code blocks — match Sentry-like look */
.code-block {
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  background: #0b1220;
  color: #e2e8f0;
  overflow: auto;
  padding: 10px 12px;
}
.code-block.dark {
  background: #0f172a;
}
.code-block pre {
  margin: 0;
  white-space: pre;
}

/* Previous events list */
.event-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.event-desc {
  display: grid;
  gap: 2px;
}
.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
