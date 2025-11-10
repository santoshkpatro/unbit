<script setup>
import { onMounted, ref, computed, h } from 'vue'
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
  Bug as LBug,
  FileCode2 as LFileCode,
  Clock as LClock,
  Folder as LFolder,
  User as LUser,
  AlertTriangle as LAlert,
} from 'lucide-vue-next'

dayjs.extend(relativeTime)

/** ---- State ---- **/
const route = useRoute()
const loading = ref(true)
const error = ref(null)
const issue = ref(null)
const resolving = ref(false)

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

const assignToMe = () => {
  antdMessage.success('Assigned to you')
}

const copyText = async (text, label = 'Copied') => {
  try {
    await navigator.clipboard.writeText(text)
    antdMessage.success(label)
  } catch {
    antdMessage.error('Copy failed')
  }
}

/** ---- Snippet rendering ----
 * If backend returns a preformatted/HTML-highlighted snippet, prefer it.
 * Expected fields (examples):
 *  - frame.snippetHtml  (string: already-highlighted HTML)
 *  - frame.snippetText  (string: plain text with newlines)
 *  - frame.code         (fallback single-line or multiline string)
 * We will render EXACTLY what backend sends without re-formatting.
 */
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
</script>

<template>
  <main class="issue-details">
    <a-card v-if="issue" class="header-card" :bordered="true">
      <div class="header">
        <div class="title-wrap">
          <div class="title-area">
            <div class="title">{{ titleText }}</div>
            <div class="subtle">
              <span v-if="issue.project?.name" class="mono">
                <LFolder class="icon" /> {{ issue.project?.name }}
              </span>
              <a-divider type="vertical" />
              <span class="mono">Event {{ issue.eventCount }}</span>
              <a-divider type="vertical" />
              <span class="mono">Last seen: {{ occurredAt.abs }} ({{ occurredAt.rel }})</span>
            </div>
          </div>
        </div>
        <a-space align="center" wrap>
          <a-badge :status="statusBadge.status" :text="statusBadge.text" />
          <a-tag :color="levelTag.color">{{ levelTag.text }}</a-tag>
          <a-tag :color="typeTag.color">{{ typeTag.text }}</a-tag>
          <a-divider type="vertical" />
          <a-button :icon="h(CopyOutlined)" @click="() => copyText(issue.id, 'Issue ID copied')">
            Copy ID
          </a-button>
          <a-button :icon="h(UserOutlined)" @click="assignToMe">Assign to me</a-button>
          <a-button
            type="primary"
            :loading="resolving"
            :icon="h(CheckCircleOutlined)"
            @click="resolveIssue"
          >
            Resolve
          </a-button>
          <a-button :icon="h(ReloadOutlined)" @click="() => location.reload()">Refresh</a-button>
          <a-button
            :icon="h(SendOutlined)"
            @click="() => copyText(shareableSnippet, 'Shareable snippet copied')"
          >
            Share
          </a-button>
        </a-space>
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
            <!-- Stacktrace -->
            <a-card :bordered="true" :title="'Stack Trace'">
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

                  <!-- Prefer exact backend formatting -->
                  <div
                    v-if="hasHtmlSnippet(frame)"
                    class="code-block"
                    v-html="frame.snippetHtml"
                  ></div>
                  <pre v-else class="code-block"><code>{{ plainSnippetFrom(frame) }}</code></pre>

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
            </a-card>

            <!-- Details -->
            <a-card class="mt" :bordered="true" :title="'Details'">
              <a-descriptions bordered :column="1" size="small">
                <a-descriptions-item label="Issue ID">{{ issue.id }}</a-descriptions-item>
                <a-descriptions-item label="Event ID">{{ issue.eventId }}</a-descriptions-item>
                <a-descriptions-item label="Type">{{ issue.type }}</a-descriptions-item>
                <a-descriptions-item label="Level">{{ issue.level }}</a-descriptions-item>
                <a-descriptions-item label="Status">{{ issue.status }}</a-descriptions-item>
                <a-descriptions-item label="Timestamp"
                  >{{ occurredAt.abs }} ({{ occurredAt.rel }})</a-descriptions-item
                >
                <a-descriptions-item
                  v-if="issue.project?.id || issue.project?.name"
                  label="Project"
                >
                  {{ issue.project?.name
                  }}<template v-if="issue.project?.id">
                    <a-divider type="vertical" /> <span class="weak">ID:</span>
                    {{ issue.project?.id }}</template
                  >
                </a-descriptions-item>
              </a-descriptions>

              <a-collapse class="mt">
                <a-collapse-panel key="raw" header="Raw payload (JSON)">
                  <div class="code-block dark">
                    <pre><code>{{ JSON.stringify(issue, null, 2) }}</code></pre>
                  </div>
                  <a-button
                    size="small"
                    class="mt-8"
                    :icon="h(CopyOutlined)"
                    @click="() => copyText(JSON.stringify(issue, null, 2), 'Payload copied')"
                  >
                    Copy payload
                  </a-button>
                </a-collapse-panel>
              </a-collapse>
            </a-card>
          </a-col>

          <a-col :xs="24" :md="8">
            <!-- Assignee -->
            <a-card :bordered="true" title="Assignee">
              <div class="meta">
                <LUser class="icon" />
                <div>
                  <div class="strong">{{ issue.assignee?.name || 'Unassigned' }}</div>
                  <div class="weak">{{ issue.assignee?.email || '—' }}</div>
                </div>
              </div>
              <a-divider />
              <a-space wrap>
                <a-button :icon="h(UserOutlined)" @click="assignToMe">Assign to me</a-button>
                <a-button
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(issue.assignee?.email || '', 'Email copied')"
                >
                  Copy email
                </a-button>
              </a-space>
            </a-card>

            <!-- Quick Actions -->
            <a-card class="mt" title="Quick actions">
              <a-space direction="vertical" style="width: 100%">
                <a-button
                  block
                  type="primary"
                  :loading="resolving"
                  :icon="h(CheckCircleOutlined)"
                  @click="resolveIssue"
                >
                  Mark as resolved
                </a-button>
                <a-button block :icon="h(ReloadOutlined)" @click="() => location.reload()"
                  >Refresh data</a-button
                >
                <a-button
                  block
                  :icon="h(CopyOutlined)"
                  @click="() => copyText(issue.id, 'Issue ID copied')"
                >
                  Copy Issue ID
                </a-button>
                <a-button
                  block
                  :icon="h(SendOutlined)"
                  @click="() => copyText(shareableSnippet, 'Shareable snippet copied')"
                >
                  Share summary snippet
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
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 10px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.title-area {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.title {
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
  word-break: break-word;
}
.subtle {
  color: rgba(0, 0, 0, 0.45);
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
  width: 22px;
  height: 22px;
}
.meta {
  display: flex;
  align-items: center;
  gap: 8px;
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
</style>
