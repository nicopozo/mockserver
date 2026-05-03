<template>
  <v-form ref="form" v-model="valid">
    <v-progress-linear indeterminate color="primary" :active="loading"/>
    <!--MOCK-->
    <v-card class="elevation-2 pb-0">
      <v-container fluid>
        <v-card-title class="px-0 py-0 pb-2">
          Mock
        </v-card-title>
        <v-row>
          <v-col cols="6">
            <!--MOCK KEY-->
            <v-text-field label="Key"
                          v-model="mock.key"
                          variant="outlined" density="compact" disabled/>
            <!--MOCK NAME-->
            <v-text-field label="Name"
                          v-model="mock.name"
                          :rules="[v => !!v || 'Name is required']"
                          required variant="outlined" density="compact"/>
            <!--MOCK GROUP-->
            <v-text-field label="Group"
                          v-model="mock.group"
                          placeholder="Examples: users, payments, auth, etc"
                          :rules="[v => !!v || 'Group is required']"
                          required variant="outlined" density="compact"/>
          </v-col>
          <v-col cols="6">
            <!--MOCK PATH-->
            <v-text-field label="Path"
                          v-model="mock.path"
                          placeholder="Example: /users/{user_id}"
                          :rules="[
                              v => !!v || 'Path is required',
                              v => /^((?!\?).)*$/.test(v) || 'Type path without query',
                              v => !!v.startsWith('/') || 'Path must start with \'/\'']"
                          required variant="outlined" density="compact"/>
            <!--MOCK METHOD-->
            <v-select label="HTTP Method"
                      v-model="mock.method"
                      :items="httpMethods"
                      :rules="[v => !!v || 'HTTP Method is required']"
                      required variant="outlined" density="compact"/>
            <!--MOCK STRATEGY-->
            <v-select label="Strategy"
                      v-model="mock.strategy"
                      :items="strategies"
                      :rules="[v => !!v || 'Strategy is required']"
                      @update:model-value="updateResponses()"
                      required variant="outlined" density="compact"/>
          </v-col>
          <v-col cols="6">
            <v-switch v-model="mock.status" color="info"
                      true-value="enabled" false-value="disabled"
                      hide-details class="v-input--reverse mx-0 my-0">
              <template #label>
                Enable mock?
              </template>
            </v-switch>
          </v-col>
          <v-col cols="6">
            <div class="text-center">
              <v-btn
                  color="primary"
                  @click="showExecutionURL()"
              >Show Execution URL
              </v-btn>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!--RESPONSE-->
    <br>
    <v-card>
      <v-container fluid>
        <v-card class="elevation-3 py-0 my-2" style="border-color: #aaa" v-for="(response, index) in mock.responses"
                v-bind:key="index">
          <v-container fluid class="pt-0">
            <v-card-title class="px-0 py-0 pb-2 d-flex align-center">
              <span class="text-subtitle-1 font-weight-bold mr-4">{{ getResponseDescriptionPrefix(index) }}</span>
              <!--RESPONSE DESCRIPTION-->
              <v-text-field v-model="response.description"
                            placeholder="Type Response Description (Optional)"
                            hide-details
                            variant="outlined"
                            density="compact"
                            class="description flex-grow-1"/>
              <!--REMOVE RESPONSE BUTTON-->
              <v-btn icon color="red" variant="text" class="ml-2" @click="removeResponse(index)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </v-card-title>
            <v-row>
              <v-col cols="6">
                <!--RESPONSE TYPE-->
                <v-text-field label="Content Type"
                              v-model="response.content_type"
                              placeholder="Example: application/json"
                              :rules="[v => !!v || 'Content Type is required']"
                              required variant="outlined" density="compact"/>
                <!--RESPONSE STATUS CODE-->
                <v-text-field label="HTTP Status"
                              v-model.number="response.http_status"
                              placeholder="Examples: 200, 201, 400, 404, 500"
                              :rules="[v => (!isNaN(parseFloat(v)) && v >= 0) || 'HTTP Status is required and greater than or equal to 0']"
                              required variant="outlined" density="compact" type="number"/>
                <!--RESPONSE DELAY-->
                <v-text-field label="Delay"
                              v-model.number="response.delay"
                              placeholder="Time to delay the response from server in milliseconds"
                              :rules="[v => (!isNaN(parseFloat(v)) && v >= 0) || 'Delay is required and greater than or equal to 0']"
                              required variant="outlined" density="compact" type="number"/>
                <!--RESPONSE SCENE-->
                <v-text-field label="Scene"
                              v-model="response.scene"
                              placeholder="Value of 'scene' variable when SCENE strategy is selected."
                              :rules="isResponseSceneRequired(mock) ? [v => !!v || 'Scene is required'] : []"
                              variant="outlined" density="compact"
                              :disabled="!isResponseSceneRequired(mock)"
                              :required="isResponseSceneRequired(mock)"/>
              </v-col>
              <v-col cols="6">
                <v-textarea label="Body"
                            v-model="response.body"
                            :rules="[v => !!v || 'Body is required']"
                            required variant="outlined" density="compact"
                            rows="8" height="238"/>
              </v-col>
            </v-row>
          </v-container>
        </v-card>
        <v-col cols="12" class="text-right">
          <v-btn variant="flat" color="primary" @click="addResponse()">New Response</v-btn>
        </v-col>
      </v-container>
    </v-card>

    <!--VARIABLES-->
    <br>
    <v-card>
      <v-container fluid>
        <v-card class="elevation-3 py-0 my-2" style="border-color: #aaa" v-for="(variable, index) in mock.variables"
                v-bind:key="index">
          <v-container fluid>
            <v-card-title class="px-0 py-0 pb-2 d-flex align-center">
              <span>Variable {{ Number(index) + 1 }}</span>
              <v-spacer/>
              <!--REMOVE VARIABLE BUTTON-->
              <v-btn icon color="red" variant="text" density="compact" @click="removeVariable(index)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </v-card-title>
            <v-row>
              <v-col cols="4">
                <!--VARIABLE TYPE-->
                <v-select label="Type"
                          v-model="variable.type"
                          :items="varTypes"
                          :rules="[v => !!v || 'Type is required']"
                          @update:model-value="updateVariables()"
                          required variant="outlined" density="compact"/>
              </v-col>
              <v-col cols="4">
                <!--VARIABLE NAME-->
                <v-text-field label="Name"
                              v-model="variable.name"
                              :rules="[v => !!v || 'Name is required']"
                              required variant="outlined" density="compact"/>
              </v-col>
              <v-col cols="4">
                <!--VARIABLE KEY-->
                <v-text-field label="Key"
                              v-model="variable.key"
                              :rules="isVariableTypeRequired(variable) ? [v => !!v || 'Key is required'] : []"
                              variant="outlined" density="compact"
                              :disabled="!isVariableTypeRequired(variable)"
                              :required="isVariableTypeRequired(variable)"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <!--ASSERTIONS-->
                <v-card class="elevation-3 py-0 my-2" style="border-color: #aaa"
                        v-for="(assertion, assertIndex) in variable.assertions"
                        v-bind:key="assertIndex">
                  <v-container fluid>
                    <v-card-title class="px-0 py-0 pb-0 d-flex align-center">
                      <span>Assertion {{ Number(assertIndex) + 1 }}</span>
                      <v-spacer/>
                      <!--REMOVE ASSERTION BUTTON-->
                      <v-btn icon color="red" variant="text" density="compact" @click="removeAssertion(index, assertIndex)">
                        <v-icon>mdi-delete</v-icon>
                      </v-btn>
                    </v-card-title>
                    <v-row>
                      <v-col cols="4">
                        <!--ASSERTION TYPE-->
                        <v-select label="Type"
                                  v-model="assertion.type"
                                  :items="assertionTypes"
                                  :rules="[v => !!v || 'Type is required']"
                                  @update:model-value="updateAssertions(index)"
                                  required variant="outlined" density="compact"/>
                      </v-col>
                      <v-col cols="4">
                        <!--VARIABLE VALUE-->
                        <v-text-field label="Value"
                                      v-model="assertion.value"
                                      :rules="isAssertionFieldRequired(assertion, 'value') ? [v => !!v || 'Value is required when Type is Equals'] : []"
                                      variant="outlined" density="compact"
                                      :disabled="!isAssertionFieldRequired(assertion, 'value')"
                                      :required="isAssertionFieldRequired(assertion, 'value')"
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="4">
                        <!--Fail on error-->
                        <v-switch v-model="assertion.fail_on_error" color="info"
                                  hide-details class="v-input--reverse mx-0 my-0">
                          <template #label>
                            Return error on failure
                          </template>
                        </v-switch>
                      </v-col>
                      <v-col cols="4">
                        <!--RANGE MIN-->
                        <v-text-field label="Range Min"
                                      v-model.number="assertion.min"
                                      placeholder="Range Min"
                                      :rules="isAssertionFieldRequired(assertion, 'min') ? [v => v < (assertion.max ?? 0)  || 'Range Min is required and lower than Max'] : []"
                                      :disabled="!isAssertionFieldRequired(assertion, 'min')"
                                      :required="isAssertionFieldRequired(assertion, 'min')"
                                      variant="outlined" density="compact" type="number"/>
                      </v-col>
                      <v-col cols="4">
                        <!--RANGE MAX-->
                        <v-text-field label="Range Max"
                                      v-model.number="assertion.max"
                                      placeholder="Range Max"
                                      :rules="isAssertionFieldRequired(assertion, 'max') ? [v => v > (assertion.min ?? 0) || 'Range Max is required and greater than Min'] : []"
                                      :disabled="!isAssertionFieldRequired(assertion, 'max')"
                                      :required="isAssertionFieldRequired(assertion, 'max')"
                                      variant="outlined" density="compact" type="number"/>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card>


                <v-col cols="12" class="text-right">
                  <v-btn variant="flat" color="primary" :disabled="!isAssertionAllowed(variable)" @click="addAssertion(index)">New Assertion</v-btn>
                </v-col>
              </v-col>
            </v-row>
          </v-container>
        </v-card>
        <v-col cols="12" class="text-right">
          <v-btn variant="flat" color="primary" @click="addVariable()">New Variable</v-btn>
        </v-col>
      </v-container>
    </v-card>

    <v-row class="mt-4">
      <v-col cols="12" class="text-right">
        <v-btn variant="flat" color="error" class="mx-1" @click="submitDelete" v-if="theKey">Delete</v-btn>
        <v-btn variant="flat" color="warning" class="mx-1" @click="resetForm">Reset</v-btn>
        <v-btn variant="flat" color="primary" class="mx-1" @click="submit" :loading="saving">Save</v-btn>
      </v-col>
    </v-row>

    <v-snackbar v-model="alert.show" :color="alert.color" :timeout="alert.timeout">{{ alert.text }}</v-snackbar>

    <v-dialog
        transition="dialog-top-transition"
        max-width="600"
        v-model="executionURL.show"
    >
      <v-card>
        <v-toolbar
            color="primary"
            theme="dark"
        >
          <v-toolbar-title>Mock execution URL:</v-toolbar-title>
        </v-toolbar>
        <v-card-text>
          <div class="text-h6 pa-12">{{ executionURL.value }}</div>
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn
              variant="text"
              @click="copyExecutionURL()"
          >Done
          </v-btn>
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
  {title: "Body", value: "body"},
  {title: "Header", value: "header"},
  {title: "Query", value: "query"},
  {title: "Random", value: "random"},
  {title: "Hash", value: "hash"},
  {title: "Path", value: "path"},
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

function getResponseDescriptionPrefix(index: number) {
  return "Response " + (index + 1).toString() + ": "
}

function newMock(): Mock {
  return {
    key: "",
    group: "",
    name: "",
    path: "",
    strategy: "",
    method: "",
    status: "enabled",
    responses: [{
      description: "",
      body: "",
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

function copyExecutionURL() {
  executionURL.show = false
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

