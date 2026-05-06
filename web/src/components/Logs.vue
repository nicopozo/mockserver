<template>
  <div>
    <!-- TOOLBAR -->
    <v-card class="elevation-2 mb-4">
      <v-container fluid>
        <v-row align="center">
          <v-col cols="12" md="6">
            <div class="d-flex align-center ga-2">
              <v-icon size="24" color="primary">mdi-format-list-bulleted</v-icon>
              <span class="text-h6 font-weight-bold">Request Logs</span>
              <v-chip color="primary" size="small" class="ml-2">{{ totalItems }}</v-chip>
            </div>
          </v-col>
          <v-col cols="12" md="6" class="d-flex justify-end ga-2">
            <!-- Auto-refresh toggle -->
            <v-btn
              :color="autoRefresh ? 'success' : 'grey'"
              :variant="autoRefresh ? 'flat' : 'outlined'"
              size="small"
              prepend-icon="mdi-refresh-auto"
              @click="toggleAutoRefresh"
            >
              {{ autoRefresh ? 'Auto-refresh ON' : 'Auto-refresh OFF' }}
            </v-btn>
            <!-- Manual refresh -->
            <v-btn
              color="primary"
              variant="outlined"
              size="small"
              :loading="loading"
              prepend-icon="mdi-refresh"
              @click="fetchLogs"
            >
              Refresh
            </v-btn>
            <!-- Clear logs -->
            <v-btn
              color="red"
              variant="outlined"
              size="small"
              prepend-icon="mdi-delete-sweep"
              @click="clearLogs"
            >
              Clear
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- LOG TABLE -->
    <v-card class="table-card elevation-2">
      <v-data-table-server
        density="compact"
        :headers="columns"
        :items="logs"
        :loading="loading"
        :items-per-page="itemsPerPage"
        :items-length="totalItems"
        v-model:options="tableOptions"
        @update:options="onTableUpdate"
        hover
        show-expand
      >
      <!-- Method chip -->
      <template v-slot:item.method="{ item }">
        <v-chip :color="methodColor(item.method)" size="small" label>
          {{ item.method }}
        </v-chip>
      </template>

      <!-- URL truncated -->
      <template v-slot:item.url="{ item }">
        <span class="text-mono text-body-2" :title="item.url">{{ item.url }}</span>
      </template>

      <!-- Response status chip -->
      <template v-slot:item.response_status="{ item }">
        <v-chip :color="statusColor(item.response_status)" size="small" label>
          {{ item.response_status }}
        </v-chip>
      </template>

      <!-- Timestamp formatted -->
      <template v-slot:item.timestamp="{ item }">
        <span class="text-caption text-mono">{{ formatTimestamp(item.timestamp) }}</span>
      </template>

      <!-- Request body preview -->
      <template v-slot:item.request_body="{ item }">
        <span class="text-caption text-mono text-truncate-cell" :title="item.request_body">
          {{ truncate(item.request_body, 60) }}
        </span>
      </template>

      <!-- Assertion result chip -->
      <template v-slot:item.assertion_errors="{ item }">
        <v-chip
          v-if="item.assertion_errors && item.assertion_errors.length > 0"
          color="red" size="small" label
          prepend-icon="mdi-close-circle-outline"
        >
          Fail ({{ item.assertion_errors.length }})
        </v-chip>
        <v-chip v-else color="green" size="small" label prepend-icon="mdi-check-circle-outline">
          Pass
        </v-chip>
      </template>

      <!-- Expanded row: full details -->
      <template v-slot:expanded-row="{ columns: cols, item }">
        <tr>
          <td :colspan="cols.length" class="pa-0">
            <v-card flat class="ma-2" border>
              <v-container fluid>
                <v-row>
                  <!-- REQUEST BODY -->
                  <v-col cols="12" md="6">
                    <div class="text-overline text-primary mb-1">Request Body</div>
                    <pre class="code-block">{{ formatJson(item.request_body) }}</pre>
                  </v-col>
                  <!-- RESPONSE BODY -->
                  <v-col cols="12" md="6">
                    <div class="text-overline text-primary mb-1">Response Body</div>
                    <pre class="code-block">{{ formatJson(item.response_body) }}</pre>
                  </v-col>
                  <!-- HEADERS -->
                  <v-col cols="12" md="6">
                    <div class="text-overline text-primary mb-1">Request Headers</div>
                    <v-table density="compact" class="rounded border">
                      <tbody>
                        <tr v-for="(value, key) in item.request_headers" :key="key">
                          <td class="font-weight-medium text-caption" style="width:40%">{{ key }}</td>
                          <td class="text-caption text-mono">{{ value }}</td>
                        </tr>
                        <tr v-if="!Object.keys(item.request_headers).length">
                          <td colspan="2" class="text-caption text-disabled">No headers</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-col>
                  <!-- QUERY PARAMS -->
                  <v-col cols="12" md="6">
                    <div class="text-overline text-primary mb-1">Query Params</div>
                    <v-table density="compact" class="rounded border">
                      <tbody>
                        <tr v-for="(value, key) in item.query_params" :key="key">
                          <td class="font-weight-medium text-caption" style="width:40%">{{ key }}</td>
                          <td class="text-caption text-mono">{{ value }}</td>
                        </tr>
                        <tr v-if="!Object.keys(item.query_params).length">
                          <td colspan="2" class="text-caption text-disabled">No query params</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-col>
                  <!-- ASSERTION ERRORS -->
                  <v-col cols="12" v-if="item.assertion_errors && item.assertion_errors.length > 0">
                    <div class="text-overline text-error mb-1">
                      <v-icon size="small" class="mr-1">mdi-shield-alert-outline</v-icon>
                      Assertion Errors ({{ item.assertion_errors.length }})
                    </div>
                    <v-alert
                      v-for="(err, errIdx) in item.assertion_errors"
                      :key="errIdx"
                      type="error"
                      variant="tonal"
                      density="compact"
                      class="mb-1 text-caption"
                    >
                      {{ err }}
                    </v-alert>
                  </v-col>
                </v-row>
              </v-container>
            </v-card>
          </td>
        </tr>
      </template>

      <!-- Empty state -->
      <template v-slot:no-data>
        <div class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1">mdi-format-list-bulleted</v-icon>
          <p class="text-h6 text-disabled mt-4">No logs yet</p>
          <p class="text-body-2 text-disabled">Make a request to <code>/mock-service/mock/*</code> to see it here.</p>
        </div>
      </template>
      </v-data-table-server>
    </v-card>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
      {{ snackbar.text }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import type { LogEntry, LogList } from '@/types'

const logs = ref<LogEntry[]>([])
const totalItems = ref(0)
const loading = ref(false)
const itemsPerPage = ref(25)
const tableOptions = ref({
  page: 1,
  itemsPerPage: 25,
  sortBy: [{ key: 'id', order: 'desc' }]
})

// Store the last ID seen for each page to support keyset pagination
// pageCursors[page] = the ID of the last item of page (page - 1)
const pageCursors = ref<Record<number, string>>({})

const autoRefresh = ref(true)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const snackbar = ref({ show: false, color: 'green', text: '' })

const columns = [
  { title: 'Time', key: 'timestamp', width: '160px', sortable: false },
  { title: 'Method', key: 'method', width: '90px', sortable: false },
  { title: 'URL', key: 'url', sortable: false },
  { title: 'Request Body', key: 'request_body', sortable: false },
  { title: 'Status', key: 'response_status', width: '80px', sortable: false },
  { title: 'Assertions', key: 'assertion_errors', width: '110px', sortable: false },
  { title: '', key: 'data-table-expand', width: '40px' },
]

const HTTP_METHOD_COLORS: Record<string, string> = {
  GET: 'blue',
  POST: 'green',
  PUT: 'orange',
  PATCH: 'purple',
  DELETE: 'red',
  OPTIONS: 'grey',
  HEAD: 'blue-grey',
}

function methodColor(method: string): string {
  return HTTP_METHOD_COLORS[method] ?? 'grey'
}

function statusColor(status: number): string {
  if (status >= 500) return 'red'
  if (status >= 400) return 'orange'
  if (status >= 300) return 'blue'
  if (status >= 200) return 'green'
  return 'grey'
}

function formatTimestamp(ts: string): string {
  const d = new Date(ts)
  const time = d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  const ms = String(d.getMilliseconds()).padStart(3, '0')
  return `${time}.${ms}`
}

function truncate(text: string, maxLen: number): string {
  if (!text || text === '<nil>') return '—'
  return text.length > maxLen ? text.slice(0, maxLen) + '…' : text
}

function formatJson(text: string): string {
  if (!text || text === '<nil>' || text === '<error>') return text || '—'
  try {
    return JSON.stringify(JSON.parse(text), null, 2)
  } catch {
    return text
  }
}

function baseURL(): string {
  return import.meta.env.PROD
    ? '/mock-service/logs'
    : 'http://localhost:8080/mock-service/logs'
}

async function fetchLogs(isAuto = false) {
  if (!isAuto) loading.value = true
  
  const { page, itemsPerPage } = tableOptions.value
  const offset = (page - 1) * itemsPerPage
  
  // Try to use the cursor for the current page
  const lastId = pageCursors.value[page]
  
  try {
    const res = await axios.get<LogList>(baseURL(), {
      params: {
        limit: itemsPerPage,
        last_id: lastId || undefined,
        // Only send offset if we don't have a cursor (for initial page jumps)
        offset: lastId ? undefined : offset
      }
    })
    
    const newLogs = res.data.results ?? []
    totalItems.value = res.data.paging?.total ?? 0

    // Store the cursor for the NEXT page
    if (newLogs.length > 0) {
      pageCursors.value[page + 1] = newLogs[newLogs.length - 1].id
    }

    // Only update if data actually changed
    if (newLogs.length !== logs.value.length || (newLogs.length > 0 && newLogs[0].id !== logs.value[0]?.id)) {
      logs.value = newLogs
    }
  } catch (err) {
    if (!isAuto) {
      showSnackbar('Error fetching logs', true)
      console.error(err)
    }
  } finally {
    if (!isAuto) loading.value = false
  }
}

function onTableUpdate(options: any) {
  // If page size changed, clear cursors because they are no longer valid
  if (options.itemsPerPage !== itemsPerPage.value) {
    pageCursors.value = {}
    itemsPerPage.value = options.itemsPerPage
  }
  fetchLogs()
}

async function clearLogs() {
  try {
    await axios.delete(baseURL())
    logs.value = []
    showSnackbar('Logs cleared')
  } catch (err) {
    showSnackbar('Error clearing logs', true)
    console.error(err)
  }
}

function toggleAutoRefresh() {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startTimer()
  } else {
    stopTimer()
  }
}

function startTimer() {
  stopTimer()
  refreshTimer = setInterval(() => fetchLogs(true), 2000)
}

function stopTimer() {
  if (refreshTimer !== null) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

function showSnackbar(text: string, isError = false) {
  snackbar.value = { show: true, color: isError ? 'red' : 'green', text }
}

onMounted(() => {
  // Initial fetch will be triggered by @update:options on v-data-table-server
  startTimer()
})

onUnmounted(() => {
  stopTimer()
})
</script>

<style scoped>
.code-block {
  font-family: 'Fira Code', 'Courier New', monospace;
  font-size: 0.75rem;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 6px;
  padding: 10px;
  overflow-x: auto;
  max-height: 200px;
  white-space: pre-wrap;
  word-break: break-all;
}

.text-mono {
  font-family: 'Fira Code', 'Courier New', monospace;
}

.text-truncate-cell {
  display: block;
  max-width: 300px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.table-card {
  border-radius: 16px;
  overflow: hidden;
}

</style>
