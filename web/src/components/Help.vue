<template>
  <div class="help-container">

    <!-- HERO HEADER -->
    <v-card class="section-card mb-6" elevation="2">
      <v-toolbar color="primary-darken-1" density="compact" dark flat>
        <v-icon start class="ml-4">mdi-book-open-page-variant</v-icon>
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">Documentation</v-toolbar-title>
      </v-toolbar>
      <v-container fluid class="pa-6">
        <div class="text-center">
          <v-icon size="48" color="primary" class="mb-3">mdi-information-outline</v-icon>
          <h1 class="text-h5 font-weight-bold mb-2">Welcome to Mocks Server</h1>
          <p class="text-body-1 text-medium-emphasis" style="max-width: 700px; margin: 0 auto;">
            Mocks Server lets you define HTTP endpoints that return preconfigured responses.
            Use it for integration testing, frontend development, or API prototyping.
          </p>
        </div>
      </v-container>
    </v-card>

    <!-- TABLE OF CONTENTS -->
    <v-card class="section-card mb-6" elevation="2">
      <v-toolbar color="primary-darken-1" density="compact" dark flat>
        <v-icon start class="ml-4">mdi-table-of-contents</v-icon>
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">Quick Navigation</v-toolbar-title>
      </v-toolbar>
      <v-container fluid class="pa-4">
        <v-row>
          <v-col v-for="section in tocSections" :key="section.id" cols="12" sm="6" md="4">
            <v-btn
              variant="tonal"
              color="primary"
              block
              class="text-none justify-start"
              :prepend-icon="section.icon"
              @click="scrollTo(section.id)"
            >
              {{ section.title }}
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- SECTIONS -->
    <v-expansion-panels v-model="openPanels" multiple variant="accordion">

      <!-- 1. GETTING STARTED -->
      <v-expansion-panel id="getting-started" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-rocket-launch-outline</v-icon>
          <span class="font-weight-bold">Getting Started</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <h3 class="mb-3">How it works</h3>
            <p>
              Mocks Server intercepts HTTP requests that match a configured <strong>path</strong> and <strong>method</strong>,
              and returns a preconfigured response. Each configuration is called a <strong>Mock</strong> (or Rule).
            </p>

            <v-alert type="info" variant="tonal" density="compact" class="my-4">
              All mock endpoints are served under <code>/mock-service/mock/*</code>.
              For example, a mock with path <code>/users/123</code> is reachable at
              <code>http://your-host/mock-service/mock/users/123</code>.
            </v-alert>

            <h3 class="mb-3">Basic workflow</h3>
            <v-timeline density="compact" side="end" class="ml-n4">
              <v-timeline-item v-for="step in workflowSteps" :key="step.num" :dot-color="step.color" size="small">
                <div>
                  <span class="font-weight-bold">{{ step.title }}</span>
                  <p class="text-body-2 text-medium-emphasis mb-0">{{ step.desc }}</p>
                </div>
              </v-timeline-item>
            </v-timeline>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 2. CREATING A MOCK -->
      <v-expansion-panel id="creating-mock" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-plus-circle-outline</v-icon>
          <span class="font-weight-bold">Creating a Mock</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>Click the <v-icon size="small">mdi-plus</v-icon> <strong>New Mock</strong> button from the List Mocks page. You need to fill in:</p>

            <v-table density="compact" class="my-4 rounded border">
              <thead><tr><th>Field</th><th>Description</th><th>Example</th></tr></thead>
              <tbody>
                <tr><td class="font-weight-medium">Name</td><td>A descriptive label for the mock</td><td><code>Get User by ID</code></td></tr>
                <tr><td class="font-weight-medium">Group</td><td>Logical group for organizing mocks</td><td><code>users</code></td></tr>
                <tr><td class="font-weight-medium">Path</td><td>The URL path (must start with <code>/</code>). Supports path variables with <code>{curly_braces}</code></td><td><code>/v1/users/{user_id}</code></td></tr>
                <tr><td class="font-weight-medium">HTTP Method</td><td>GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD</td><td><code>GET</code></td></tr>
                <tr><td class="font-weight-medium">Strategy</td><td>How the server selects a response (see below)</td><td><code>normal</code></td></tr>
                <tr><td class="font-weight-medium">Status</td><td>Enable or disable the mock without deleting it</td><td><code>enabled</code></td></tr>
              </tbody>
            </v-table>

            <v-alert type="warning" variant="tonal" density="compact" class="my-4">
              The <strong>Path</strong> must not include query parameters (<code>?key=value</code>). Only the path portion is used for matching.
            </v-alert>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 3. RESPONSE STRATEGIES -->
      <v-expansion-panel id="strategies" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-layers-outline</v-icon>
          <span class="font-weight-bold">Response Strategies</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>Each mock can have <strong>multiple responses</strong>. The strategy determines which response is returned:</p>

            <v-card v-for="s in strategyDetails" :key="s.name" variant="outlined" class="mb-4">
              <v-card-title class="text-subtitle-1 d-flex align-center">
                <v-chip :color="s.color" size="small" class="mr-3" label>{{ s.name }}</v-chip>
                {{ s.title }}
              </v-card-title>
              <v-card-text>
                <p>{{ s.description }}</p>
                <pre v-if="s.example" class="code-block mt-2">{{ s.example }}</pre>
              </v-card-text>
            </v-card>

            <v-alert type="info" variant="tonal" density="compact" class="mt-4">
              <strong>Scene Strategy – Step by Step:</strong>
              <ol class="mt-2 ml-4">
                <li>Set the mock strategy to <strong>Scene</strong>.</li>
                <li>Add a <strong>Dynamic Variable</strong> named exactly <code>scene</code> (e.g. type <em>Body</em>, key <code>$.status</code>).</li>
                <li>In each response, set the <strong>Scene Name</strong> to match possible values of that variable.</li>
                <li>Optionally, set one response's scene name to <code>default</code> as a fallback.</li>
              </ol>
            </v-alert>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 4. DYNAMIC VARIABLES -->
      <v-expansion-panel id="variables" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-variable</v-icon>
          <span class="font-weight-bold">Dynamic Variables</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>
              Variables let you extract values from incoming requests and inject them into the response body.
              Use the syntax <code>{variable_name}</code> in your response body to reference them.
            </p>

            <v-table density="compact" class="my-4 rounded border">
              <thead><tr><th>Type</th><th>Source</th><th>Key / Path</th><th>Example</th></tr></thead>
              <tbody>
                <tr>
                  <td><v-chip size="x-small" color="blue" label>body</v-chip></td>
                  <td>Request JSON body</td>
                  <td>JSONPath expression</td>
                  <td><code>$.user.name</code></td>
                </tr>
                <tr>
                  <td><v-chip size="x-small" color="green" label>header</v-chip></td>
                  <td>Request header</td>
                  <td>Header name</td>
                  <td><code>Authorization</code></td>
                </tr>
                <tr>
                  <td><v-chip size="x-small" color="orange" label>query</v-chip></td>
                  <td>Query parameter</td>
                  <td>Parameter name</td>
                  <td><code>page</code></td>
                </tr>
                <tr>
                  <td><v-chip size="x-small" color="purple" label>path</v-chip></td>
                  <td>URL path variable</td>
                  <td>Variable name from path template</td>
                  <td><code>user_id</code></td>
                </tr>
                <tr>
                  <td><v-chip size="x-small" color="teal" label>random</v-chip></td>
                  <td>Generated random number</td>
                  <td>—</td>
                  <td>Returns a random integer</td>
                </tr>
                <tr>
                  <td><v-chip size="x-small" color="grey" label>hash</v-chip></td>
                  <td>Generated SHA-256 hash</td>
                  <td>—</td>
                  <td>Returns a random hash string</td>
                </tr>
              </tbody>
            </v-table>

            <h3 class="mb-3 mt-4">Example: Echo back a path variable</h3>
            <p>Given a mock with path <code>/v1/users/{user_id}</code>:</p>
            <ol class="ml-4 mb-3">
              <li>Add a variable: type = <code>path</code>, name = <code>userId</code>, key = <code>user_id</code>.</li>
              <li>In the response body, use: <code>{"id": "{userId}"}</code></li>
            </ol>
            <p>A request to <code>/mock-service/mock/v1/users/42</code> will return <code>{"id": "42"}</code>.</p>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 5. ASSERTIONS -->
      <v-expansion-panel id="assertions" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-shield-check-outline</v-icon>
          <span class="font-weight-bold">Assertions</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>
              Assertions validate the value of a dynamic variable against expected conditions.
              They are configured per variable and checked before returning the response.
            </p>

            <v-table density="compact" class="my-4 rounded border">
              <thead><tr><th>Type</th><th>Description</th><th>Config</th></tr></thead>
              <tbody>
                <tr><td class="font-weight-medium">Equals</td><td>Variable must exactly match the expected value</td><td>Value field required</td></tr>
                <tr><td class="font-weight-medium">Is Not Equal</td><td>Variable must not match the expected value</td><td>Value field required</td></tr>
                <tr><td class="font-weight-medium">Regex Match</td><td>Variable must match the specified regular expression</td><td>Value (regex) required</td></tr>
                <tr><td class="font-weight-medium">Contains</td><td>Variable must contain the specified substring</td><td>Value field required</td></tr>
                <tr><td class="font-weight-medium">Starts With</td><td>Variable must start with the specified prefix</td><td>Value field required</td></tr>
                <tr><td class="font-weight-medium">Ends With</td><td>Variable must end with the specified suffix</td><td>Value field required</td></tr>
                <tr><td class="font-weight-medium">Length</td><td>Variable length must be within min–max range</td><td>Min and Max fields required</td></tr>
                <tr><td class="font-weight-medium">Is One Of</td><td>Variable must be one of the comma-separated values</td><td>Value (comma-separated) required</td></tr>
                <tr><td class="font-weight-medium">Is Boolean</td><td>Variable must be "true" or "false"</td><td>No extra config</td></tr>
                <tr><td class="font-weight-medium">JSON Schema</td><td>Variable must validate against the provided JSON schema</td><td>Value (JSON Schema) required</td></tr>
                <tr><td class="font-weight-medium">Is String</td><td>Variable must be a non-numeric string</td><td>No extra config</td></tr>
                <tr><td class="font-weight-medium">Is Number</td><td>Variable must be a valid number</td><td>No extra config</td></tr>
                <tr><td class="font-weight-medium">Is Present</td><td>Variable must exist and not be empty</td><td>No extra config</td></tr>
                <tr><td class="font-weight-medium">Numeric Range</td><td>Variable must be a number within min–max range</td><td>Min and Max fields required</td></tr>
              </tbody>
            </v-table>

            <v-alert type="warning" variant="tonal" density="compact" class="my-4">
              When <strong>"Fail request if assertion fails"</strong> is enabled, the server returns
              a <code>400 Bad Request</code> if the assertion fails. Otherwise, the failure is only logged.
            </v-alert>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 6. IMPORT / EXPORT -->
      <v-expansion-panel id="import-export" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-swap-horizontal</v-icon>
          <span class="font-weight-bold">Import &amp; Export</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <h3 class="mb-3">Export</h3>
            <p>
              From the List Mocks page, click <strong>Export</strong> to download all mocks as a JSON file.
              This is useful for backups or sharing configurations across environments.
            </p>

            <h3 class="mb-3 mt-4">Import</h3>
            <p>
              Click <strong>Import</strong> and select a previously exported JSON file.
              The importer will create new mocks or update existing ones (matched by key).
            </p>

            <v-alert type="info" variant="tonal" density="compact" class="my-4">
              Import is non-destructive: existing mocks not present in the file are left unchanged.
            </v-alert>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 7. REQUEST LOGS -->
      <v-expansion-panel id="logs" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-format-list-bulleted</v-icon>
          <span class="font-weight-bold">Request Logs</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>
              The <strong>Logs</strong> page displays all requests made to your mock endpoints, including:
            </p>
            <ul class="ml-4 mb-3">
              <li>Timestamp, HTTP method, URL</li>
              <li>Request body, headers, and query parameters</li>
              <li>Response status code and body</li>
            </ul>
            <p>
              Use the <strong>Auto-refresh</strong> toggle for real-time monitoring,
              or <strong>Clear</strong> to reset the log history.
            </p>
          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

      <!-- 8. API REFERENCE -->
      <v-expansion-panel id="api" class="section-card mb-4" elevation="2">
        <v-expansion-panel-title class="panel-title">
          <v-icon color="primary" class="mr-3">mdi-api</v-icon>
          <span class="font-weight-bold">REST API Reference</span>
        </v-expansion-panel-title>
        <v-expansion-panel-text>
          <div class="help-content pa-4">
            <p>Mocks Server exposes a REST API for programmatic management:</p>

            <v-table density="compact" class="my-4 rounded border">
              <thead><tr><th>Method</th><th>Endpoint</th><th>Description</th></tr></thead>
              <tbody>
                <tr><td><v-chip size="x-small" color="green" label>POST</v-chip></td><td><code>/mock-service/rules</code></td><td>Create a new mock</td></tr>
                <tr><td><v-chip size="x-small" color="blue" label>GET</v-chip></td><td><code>/mock-service/rules</code></td><td>Search / list mocks</td></tr>
                <tr><td><v-chip size="x-small" color="blue" label>GET</v-chip></td><td><code>/mock-service/rules/:key</code></td><td>Get a specific mock</td></tr>
                <tr><td><v-chip size="x-small" color="orange" label>PUT</v-chip></td><td><code>/mock-service/rules/:key</code></td><td>Update a mock</td></tr>
                <tr><td><v-chip size="x-small" color="red" label>DELETE</v-chip></td><td><code>/mock-service/rules/:key</code></td><td>Delete a mock</td></tr>
                <tr><td><v-chip size="x-small" color="blue" label>GET</v-chip></td><td><code>/mock-service/rules/export</code></td><td>Export all mocks as JSON</td></tr>
                <tr><td><v-chip size="x-small" color="green" label>POST</v-chip></td><td><code>/mock-service/rules/import</code></td><td>Import mocks from JSON</td></tr>
                <tr><td><v-chip size="x-small" color="blue" label>GET</v-chip></td><td><code>/mock-service/logs</code></td><td>Get request logs</td></tr>
                <tr><td><v-chip size="x-small" color="red" label>DELETE</v-chip></td><td><code>/mock-service/logs</code></td><td>Clear request logs</td></tr>
              </tbody>
            </v-table>

          </div>
        </v-expansion-panel-text>
      </v-expansion-panel>

    </v-expansion-panels>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const openPanels = ref<number[]>([])

const tocSections = [
  { id: 'getting-started', title: 'Getting Started', icon: 'mdi-rocket-launch-outline' },
  { id: 'creating-mock', title: 'Creating a Mock', icon: 'mdi-plus-circle-outline' },
  { id: 'strategies', title: 'Response Strategies', icon: 'mdi-layers-outline' },
  { id: 'variables', title: 'Dynamic Variables', icon: 'mdi-variable' },
  { id: 'assertions', title: 'Assertions', icon: 'mdi-shield-check-outline' },
  { id: 'import-export', title: 'Import & Export', icon: 'mdi-swap-horizontal' },
  { id: 'logs', title: 'Request Logs', icon: 'mdi-format-list-bulleted' },
  { id: 'api', title: 'REST API Reference', icon: 'mdi-api' },
]

const workflowSteps = [
  { num: 1, color: 'primary', title: 'Create a mock', desc: 'Define the path, HTTP method, and strategy.' },
  { num: 2, color: 'primary', title: 'Configure responses', desc: 'Set the status code, body, and content type.' },
  { num: 3, color: 'primary', title: 'Add variables (optional)', desc: 'Extract values from the request to use in responses.' },
  { num: 4, color: 'primary', title: 'Add assertions (optional)', desc: 'Validate incoming request data before responding.' },
  { num: 5, color: 'success', title: 'Hit the endpoint', desc: 'Use the Execution URL to call your mock and see it in action.' },
]

const strategyDetails = [
  {
    name: 'Normal',
    title: 'Always returns the first response',
    color: 'blue',
    description: 'The simplest strategy. No matter how many responses are configured, it always returns the first one. Best for simple, predictable mocks.',
    example: null,
  },
  {
    name: 'Sequential',
    title: 'Cycles through responses in order',
    color: 'green',
    description: 'Returns responses in order: 1st call → Response 1, 2nd call → Response 2, etc. After the last response, it wraps back to the first. Useful for simulating state changes across multiple calls.',
    example: null,
  },
  {
    name: 'Random',
    title: 'Returns a random response each time',
    color: 'orange',
    description: 'Picks a random response from the configured list on each call. Useful for simulating unpredictable behavior or load testing.',
    example: null,
  },
  {
    name: 'Scene',
    title: 'Selects response based on request data',
    color: 'purple',
    description: 'The most powerful strategy. It reads a value from the request (via a variable named "scene") and matches it to a response\'s Scene Name. If no match is found, a response with scene "default" is used as fallback.',
    example: `// Example: Route by request body field "status"
// 1. Variable: type=body, name=scene, key=$.status
// 2. Responses:
//    - Scene "approved"  → { "result": "Payment OK" }
//    - Scene "rejected"  → { "result": "Payment failed" }
//    - Scene "default"   → { "result": "Unknown status" }
//
// POST with {"status": "approved"} → returns "Payment OK"
// POST with {"status": "rejected"} → returns "Payment failed"
// POST with {"status": "xyz"}      → returns "Unknown status"`,
  },
]

function scrollTo(id: string) {
  const index = tocSections.findIndex(s => s.id === id)
  if (index >= 0 && !openPanels.value.includes(index)) {
    openPanels.value.push(index)
  }
  setTimeout(() => {
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }, 300)
}
</script>

<style scoped>
.help-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 16px;
}

.section-card {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.panel-title {
  font-size: 1rem;
}

.help-content h3 {
  font-size: 1.1rem;
  font-weight: 700;
  color: rgb(var(--v-theme-primary));
}

.help-content p,
.help-content li {
  line-height: 1.7;
}

.help-content code {
  font-family: 'Fira Code', 'Roboto Mono', monospace;
  font-size: 0.85em;
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(25, 118, 210, 0.08);
  color: rgb(var(--v-theme-primary));
}

.code-block {
  font-family: 'Fira Code', 'Roboto Mono', monospace;
  font-size: 0.8rem;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  white-space: pre-wrap;
  line-height: 1.6;
}
</style>

<style>
/* Theme-aware help styles */
.v-theme--light .code-block {
  background: #f5f5f5 !important;
  color: #37474f !important;
}

.v-theme--dark .code-block {
  background: #1a1a1a !important;
  color: #a5d6a7 !important;
}

.v-theme--dark .help-content code {
  background: rgba(100, 255, 218, 0.1);
}
</style>
