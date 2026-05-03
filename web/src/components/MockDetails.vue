<template>
  <v-form ref="form" v-model="valid" class="mock-details-form">
    <v-progress-linear indeterminate color="primary" :active="loading" height="4" rounded/>
    
    <!--MOCK BASIC INFO-->
    <v-card class="section-card mb-6" elevation="2">
      <v-toolbar color="primary-darken-1" density="compact" dark flat>
        <v-icon start class="ml-4">mdi-cog-outline</v-icon>
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">General Configuration</v-toolbar-title>
      </v-toolbar>
      
      <v-container fluid class="pa-6">
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field label="Key"
                          v-model="mock.key"
                          variant="filled" density="comfortable" disabled
                          prepend-inner-icon="mdi-key-variant"/>
            
            <v-text-field label="Name"
                          v-model="mock.name"
                          :rules="[v => !!v || 'Name is required']"
                          required variant="outlined" density="comfortable"
                          prepend-inner-icon="mdi-format-title"/>
            
            <v-text-field label="Group"
                          v-model="mock.group"
                          placeholder="Examples: users, payments, auth"
                          :rules="[v => !!v || 'Group is required']"
                          required variant="outlined" density="comfortable"
                          prepend-inner-icon="mdi-folder-outline"/>
          </v-col>
          
          <v-col cols="12" md="6">
            <v-text-field label="Path"
                          v-model="mock.path"
                          placeholder="Example: /users/{user_id}"
                          :rules="[
                              v => !!v || 'Path is required',
                              v => /^((?!\?).)*$/.test(v) || 'Type path without query',
                              v => !!v.startsWith('/') || 'Path must start with \'/\'']"
                          required variant="outlined" density="comfortable"
                          prepend-inner-icon="mdi-link-variant"/>
            
            <v-select label="HTTP Method"
                      v-model="mock.method"
                      :items="httpMethods"
                      :rules="[v => !!v || 'HTTP Method is required']"
                      required variant="outlined" density="comfortable"
                      prepend-inner-icon="mdi-api"/>
            
            <v-select label="Strategy"
                      v-model="mock.strategy"
                      :items="strategies"
                      :rules="[v => !!v || 'Strategy is required']"
                      @update:model-value="updateResponses()"
                      required variant="outlined" density="comfortable"
                      prepend-inner-icon="mdi-layers-outline"/>
          </v-col>
          
          <v-col cols="12" class="d-flex align-center justify-space-between pt-0">
            <v-switch v-model="mock.status" color="success"
                      true-value="enabled" false-value="disabled"
                      hide-details class="status-switch">
              <template #label>
                <span class="font-weight-bold ml-2">
                  {{ mock.status === 'enabled' ? 'MOCK ENABLED' : 'MOCK DISABLED' }}
                </span>
              </template>
            </v-switch>
            
            <v-btn
                variant="tonal"
                color="primary"
                prepend-icon="mdi-rocket-launch-outline"
                @click="showExecutionURL()"
                class="text-none"
            >Show Execution URL
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!--RESPONSES SECTION-->
    <v-card class="section-card mb-6" elevation="2">
      <v-toolbar color="primary-darken-1" density="compact" dark flat>
        <v-icon start class="ml-4">mdi-undo-variant</v-icon>
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">Server Responses</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn variant="text" icon="mdi-plus" @click="addResponse()" class="mr-2"></v-btn>
      </v-toolbar>
      
      <v-container fluid class="pa-4" v-if="mock.responses && mock.responses.length > 0">
      <v-expansion-panels variant="accordion">
        <v-expansion-panel
          v-for="(response, index) in mock.responses"
          :key="index"
          class="response-panel"
        >
          <v-expansion-panel-title>
            <div class="d-flex align-center w-100 pr-4">
              <v-chip :color="getStatusColor(response.http_status)" size="small" class="mr-3 font-weight-bold">
                {{ response.http_status }}
              </v-chip>
              <span class="font-weight-medium">{{ response.description || 'Response ' + (index + 1) }}</span>
              <v-spacer></v-spacer>
              <v-btn icon="mdi-delete-outline" variant="text" color="error" density="comfortable" @click.stop="removeResponse(index)"></v-btn>
            </div>
          </v-expansion-panel-title>
          
          <v-expansion-panel-text>
            <v-row class="mt-2">
              <v-col cols="12" md="6">
                <v-text-field label="Description" v-model="response.description" variant="outlined" density="comfortable" placeholder="Descriptive name for this response"/>
                <v-text-field label="Content Type" v-model="response.content_type" variant="outlined" density="comfortable"/>
                <v-row>
                  <v-col cols="6">
                    <v-text-field label="HTTP Status" v-model.number="response.http_status" variant="outlined" density="comfortable" type="number"/>
                  </v-col>
                  <v-col cols="6">
                    <v-text-field label="Delay (ms)" v-model.number="response.delay" variant="outlined" density="comfortable" type="number"/>
                  </v-col>
                </v-row>
                <v-text-field label="Scene Name" v-model="response.scene" variant="outlined" density="comfortable" 
                              :disabled="!isResponseSceneRequired(mock)"
                              placeholder="Required for 'Scene' strategy"/>
              </v-col>
              <v-col cols="12" md="6">
                <v-textarea label="Response Body"
                            v-model="response.body"
                            variant="outlined"
                            class="code-editor"
                            auto-grow
                            rows="10"
                            placeholder='{ "status": "ok" }'/>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
      </v-container>
    </v-card>

    <!--VARIABLES SECTION-->
    <v-card class="section-card mb-8" elevation="2">
      <v-toolbar color="primary-darken-1" density="compact" dark flat>
        <v-icon start class="ml-4">mdi-variable</v-icon>
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">Dynamic Variables</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn variant="text" icon="mdi-plus" @click="addVariable()" class="mr-2"></v-btn>
      </v-toolbar>

      <v-container fluid class="pa-4" v-if="mock.variables && mock.variables.length > 0">
      <v-row v-for="(variable, index) in mock.variables" :key="index" class="variable-row mb-4 pa-2 rounded-lg border">
        <v-col cols="12">
          <div class="d-flex align-center mb-2">
            <v-avatar color="secondary" size="24" class="mr-2 text-caption font-weight-bold">{{ index + 1 }}</v-avatar>
            <span class="font-weight-bold">Variable Configuration</span>
            <v-spacer></v-spacer>
            <v-btn icon="mdi-delete-outline" variant="text" color="error" density="comfortable" @click="removeVariable(index)"></v-btn>
          </div>
          
          <v-row>
            <v-col cols="12" md="4">
              <v-select label="Source Type" v-model="variable.type" :items="varTypes" variant="outlined" density="comfortable" @update:model-value="updateVariables()"/>
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field label="Variable Name" v-model="variable.name" variant="outlined" density="comfortable" placeholder="How to call it in the body"/>
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field label="Source Key/Path" v-model="variable.key" variant="outlined" density="comfortable" 
                            :disabled="!isVariableTypeRequired(variable)"
                            placeholder="JSON Path or Header Key"/>
            </v-col>
          </v-row>

          <!--ASSERTIONS-->
          <div class="mt-2 pl-4 border-left">
            <div class="d-flex align-center mb-2">
              <v-icon size="small" color="grey" class="mr-1">mdi-shield-check-outline</v-icon>
              <span class="text-caption font-weight-bold text-grey">Assertions</span>
              <v-spacer></v-spacer>
              <v-btn variant="text" color="primary" size="x-small" prepend-icon="mdi-plus" :disabled="!isAssertionAllowed(variable)" @click="addAssertion(index)">
                Add Assertion
              </v-btn>
            </div>
            
            <v-row v-for="(assertion, assertIndex) in variable.assertions" :key="assertIndex" class="assertion-box mb-2 pa-3 rounded">
              <v-col cols="12" md="3">
                <v-select label="Assert Type" v-model="assertion.type" :items="assertionTypes" variant="underlined" density="compact" @update:model-value="updateAssertions(index)"/>
              </v-col>
              <v-col cols="12" md="3">
                <v-text-field label="Expected Value" v-model="assertion.value" variant="underlined" density="compact" :disabled="!isAssertionFieldRequired(assertion, 'value')"/>
              </v-col>
              <v-col cols="12" md="2">
                <v-text-field label="Min" v-model.number="assertion.min" variant="underlined" density="compact" type="number" :disabled="!isAssertionFieldRequired(assertion, 'min')"/>
              </v-col>
              <v-col cols="12" md="2">
                <v-text-field label="Max" v-model.number="assertion.max" variant="underlined" density="compact" type="number" :disabled="!isAssertionFieldRequired(assertion, 'max')"/>
              </v-col>
              <v-col cols="12" md="2" class="d-flex align-center justify-end">
                <v-btn icon="mdi-close" variant="text" color="grey" size="small" @click="removeAssertion(index, assertIndex)"></v-btn>
              </v-col>
              <v-col cols="12" class="pt-0">
                <v-switch v-model="assertion.fail_on_error" color="error" label="Fail request if assertion fails" density="compact" hide-details class="text-caption"/>
              </v-col>
            </v-row>
          </div>
        </v-col>
      </v-row>
      </v-container>
    </v-card>

    <!--STICKY BOTTOM ACTIONS-->
    <div class="actions-bar d-flex align-center justify-end pa-4 rounded-lg elevation-10 mt-6">
      <v-btn variant="text" color="error" class="mr-2" prepend-icon="mdi-delete" @click="submitDelete" v-if="theKey">Delete Mock</v-btn>
      <v-spacer></v-spacer>
      <v-btn variant="outlined" color="grey-darken-1" class="mr-3" prepend-icon="mdi-history" @click="resetForm">Discard Changes</v-btn>
      <v-btn variant="elevated" color="primary" size="large" class="px-8 font-weight-bold" prepend-icon="mdi-content-save" @click="submit" :loading="saving">Save Mock</v-btn>
    </div>

    <v-snackbar v-model="alert.show" :color="alert.color" :timeout="alert.timeout" rounded="lg">
      <div class="d-flex align-center">
        <v-icon start>{{ alert.color === 'green' ? 'mdi-check-circle' : 'mdi-alert-circle' }}</v-icon>
        {{ alert.text }}
      </div>
    </v-snackbar>

    <v-dialog
        transition="dialog-bottom-transition"
        max-width="700"
        v-model="executionURL.show"
    >
      <v-card class="rounded-xl">
        <v-toolbar color="primary" theme="dark">
          <v-toolbar-title class="font-weight-bold">Ready for Execution</v-toolbar-title>
          <v-btn icon="mdi-close" @click="executionURL.show = false"></v-btn>
        </v-toolbar>
        <v-card-text class="pa-8 text-center">
          <p class="text-subtitle-1 text-grey-darken-1 mb-4">Use this URL to hit your mock server:</p>
          <v-sheet color="grey-lighten-4" class="pa-4 rounded-lg border d-flex align-center">
            <code class="text-h6 text-primary flex-grow-1 text-truncate">{{ executionURL.value }}</code>
            <v-btn icon="mdi-content-copy" variant="text" color="primary" @click="copyURL()"></v-btn>
          </v-sheet>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer></v-spacer>
          <v-btn variant="flat" color="primary" class="px-6" @click="executionURL.show = false">Got it!</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-form>
</template>


<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import type { Mock, Variable, Assertion, Response } from '@/types';

const props = defineProps<{
  theKey?: string;
  theName?: string;
}>();

const route = useRoute();
const router = useRouter();
const form = ref<any>(null);

const mock = ref<Mock>(newMock());
const valid = ref(false);
const loading = ref(false);
const saving = ref(false);

const httpMethods = [
  {title: "GET", value: "GET"},
  {title: "POST", value: "POST"},
  {title: "PUT", value: "PUT"},
  {title: "PATCH", value: "PATCH"},
  {title: "DELETE", value: "DELETE"},
  {title: "OPTIONS", value: "OPTIONS"},
  {title: "HEAD", value: "HEAD"},
];
const strategies = [
  {title: "Normal", value: "normal"},
  {title: "Scene", value: "scene"},
  {title: "Random", value: "random"},
  {title: "Sequential", value: "sequential"},
];
const varTypes = [
  {title: "Body (JSON Path)", value: "body"},
  {title: "Header", value: "header"},
  {title: "Query Param", value: "query"},
  {title: "Random Number", value: "random"},
  {title: "SHA256 Hash", value: "hash"},
  {title: "Path Variable", value: "path"},
];
const assertionTypes = [
  {title: "Equals", value: "equals"},
  {title: "Is string", value: "string"},
  {title: "Is number", value: "number"},
  {title: "Is present", value: "present"},
  {title: "Numeric Range", value: "range"},
];

const alert = reactive({
  show: false,
  color: "green",
  text: "",
  timeout: "5000"
});

const executionURL = reactive({
  value: "",
  show: false
});

function getStatusColor(status: number) {
  if (status < 300) return 'success'
  if (status < 400) return 'warning'
  return 'error'
}

function baseURL() {
  if (import.meta.env.PROD) {
    return "/mock-service/rules"
  }
  return "http://localhost:8080/mock-service/rules"
}

function submit() {
  if (!form.value?.validate()) {
    showAlert("Some fields are not valid!", "validation error");
    return;
  }

  if (props.theKey) {
    submitUpdate();
  } else {
    submitCreate();
  }
}

async function submitCreate() {
  const confirmTitle = "Creating New Mock";
  const confirmMsg = confirmTitle + "\n\nPlease confirm you want to create this mock";
  if (window.confirm(confirmMsg)) {
    createMock();
  }
}

async function submitUpdate() {
  const confirmTitle = "Updating Mock: " + props.theKey;
  const confirmMsg = confirmTitle + "\n\nPlease confirm you want to update this mock";
  if (window.confirm(confirmMsg)) {
    updateMock();
  }
}

async function submitDelete() {
  const confirmTitle = "Deleting Mock: " + props.theKey;
  const confirmMsg = confirmTitle + "\n\nPlease confirm you want to delete this mock";
  if (window.confirm(confirmMsg)) {
    deleteMock();
  }
}

async function resetForm() {
  const confirmTitle = "Reset Form";
  const confirmMsg = confirmTitle + "\n\nAll changes will be lost, are you sure?";
  if (window.confirm(confirmMsg)) {
    form.value?.reset();
    initialize();
  }
}

function createMock() {
  saving.value = true;
  axios
      .post<Mock>(baseURL(), mock.value, {
        headers: { "Content-Type": "application/json" },
      })
      .then((res) => {
        router.push({name: 'MockDetails', params: {theKey: res.data.key, theName: res.data.name}});
      })
      .catch((err) => showAlert("Error creating mock", err))
      .finally(() => saving.value = false);
}

function updateMock() {
  saving.value = true;
  axios
      .put(baseURL() + "/" + props.theKey, mock.value, {
        headers: { "Content-Type": "application/json" },
      })
      .then(() => showAlert("Mock successfully updated!"))
      .catch((err) => showAlert("Error updating mock", err))
      .finally(() => saving.value = false);
}

function deleteMock() {
  axios
      .delete(baseURL() + "/" + props.theKey)
      .then(() => router.push({name: 'ListMocks'}))
      .catch((err) => showAlert("Error deleting mock!", err));
}

function showAlert(text: string, err?: any) {
  alert.text = text;
  alert.color = err == null ? "green" : "red";
  alert.show = true;
  if (err) console.error(err);
}

function addVariable() {
  const newVar: Variable = {
    type: "body",
    name: "",
    key: "",
    assertions:[]
  };
  if (!mock.value.variables) mock.value.variables = [];
  mock.value.variables.push(newVar);
}

function removeVariable(i: number) {
  mock.value.variables.splice(i, 1);
}

function addAssertion(variableIndex: number) {
  const newAssertion: Assertion = {
    fail_on_error: true,
    type: "equals",
    variable_name: "",
    value: "",
    min: 0,
    max: 0
  };
  if (!mock.value.variables[variableIndex].assertions) {
    mock.value.variables[variableIndex].assertions = [];
  }
  mock.value.variables[variableIndex].assertions.push(newAssertion);
}

function removeAssertion(variableIndex: number, assertionIndex: number) {
  mock.value.variables[variableIndex].assertions.splice(assertionIndex, 1);
}

function addResponse() {
  const newResponse: Response = {
    description: "",
    body: "",
    content_type: "application/json",
    http_status: 200,
    delay: 0,
    scene: "",
  };
  if (!mock.value.responses) {
    mock.value.responses = [newResponse];
  } else {
    mock.value.responses.push(newResponse);
  }
}

function removeResponse(i: number) {
  mock.value.responses.splice(i, 1);
}

function updateResponses() {
  if (mock.value.strategy !== "scene") {
    mock.value.responses?.forEach(r => {
      r.scene = ""
    });
  }
}

function updateVariables() {
  mock.value.variables?.forEach(v => {
    if (v.type !== "body" && v.type !== "query" && v.type !== "header" && v.type !== "path") {
      v.key = "";
    }
  });
}

function updateAssertions(variableIndex: number) {
  mock.value.variables[variableIndex].assertions?.forEach(a => {
    if (a.type !== "range") {
      a.min = 0;
      a.max = 0;
      a.value = "";
    } else {
      a.min = 0;
      a.max = 1;
      a.value = "";
    }
  });
}

function isResponseSceneRequired(m: Mock) {
  return m.strategy === "scene";
}

function isVariableTypeRequired(variable: Variable) {
  return variable.type === 'body' || variable.type === 'query' || variable.type === 'header' || variable.type === 'path';
}

function isAssertionFieldRequired(assertion: Assertion, field: string) {
  switch (assertion.type) {
    case "equals":
      return field === "value"
    case "range":
      return field === "min" || field === "max";
    default:
      return false;
  }
}

function isAssertionAllowed(variable: Variable) {
  return variable.type === 'body' || variable.type === 'query' || variable.type === 'header' || variable.type === 'path';
}

function newMock(): Mock {
  return {
    key: "",
    group: "",
    name: "",
    path: "",
    strategy: "normal",
    method: "GET",
    status: "enabled",
    responses: [{
      description: "Default Response",
      body: '{\n  "message": "Hello World"\n}',
      content_type: "application/json",
      http_status: 200,
      delay: 0,
      scene: "",
    }],
    variables: [],
  };
}

function showExecutionURL() {
  let url = getExecutionURL()
  if (url === "") {
    showAlert("Path cannot be empty.", "error");
  } else {
    executionURL.value = url
    executionURL.show = true
  }
}

function getExecutionURL() {
  let path = mock.value.path;
  if (!path) return ""

  if (!path.startsWith("/")) path = "/" + path
  return window.location.protocol + "//" + window.location.host + "/mock-service/mock" + path
}

function copyURL() {
  navigator.clipboard.writeText(executionURL.value);
  showAlert("URL copied to clipboard!");
}

function initialize() {
  if (props.theKey) {
    loading.value = true;
    axios
        .get<Mock>(baseURL() + "/" + props.theKey)
        .then((res) => {
          mock.value = res.data;
        })
        .catch((err) => showAlert("Error getting mock info!", err))
        .finally(() => loading.value = false);
  } else {
    mock.value = newMock();
    form.value?.resetValidation();
  }
}

onMounted(() => {
  initialize();
});

watch(() => route.path, () => {
  initialize();
});

</script>

<style scoped>
.mock-details-form {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px;
}

.section-card {
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(0,0,0,0.05);
}

.response-panel {
  border: 1px solid rgba(0,0,0,0.05);
  margin-bottom: 8px !important;
  border-radius: 8px !important;
}

.variable-row {
  border: 1px solid rgba(0,0,0,0.1);
}

.border-left {
  border-left: 3px solid #1976D2;
}

.code-editor :deep(textarea) {
  font-family: 'Fira Code', 'Roboto Mono', monospace !important;
  font-size: 0.85rem !important;
  line-height: 1.5 !important;
}

.actions-bar {
  position: sticky;
  bottom: 16px;
  z-index: 10;
  border: 1px solid rgba(0,0,0,0.05);
}

.status-switch {
  transform: scale(0.9);
}

.assertion-box {
  border: 1px dashed #ccc;
}
</style>

<style>
/* Estilos Globales para Temas */
.v-theme--light .variable-row {
  background: white !important;
}

.v-theme--dark .variable-row {
  background: #2a2a2a !important;
  border: 1px solid rgba(255, 255, 255, 0.05) !important;
}

.v-theme--light .actions-bar {
  background: white !important;
}

.v-theme--dark .actions-bar {
  background: #1e1e1e !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.v-theme--dark .assertion-box {
  background: #212121 !important;
  border: 1px dashed rgba(255, 255, 255, 0.3) !important;
}

.v-theme--light .code-editor :deep(textarea) {
  background-color: #fafafa !important;
  color: #2c3e50 !important;
}

.v-theme--dark .code-editor :deep(textarea) {
  background-color: #1a1a1a !important;
  color: #64ffda !important; /* Un toque de color para el código en dark mode */
}
</style>
