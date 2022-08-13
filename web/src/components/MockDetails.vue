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
                          outlined dense disabled/>
            <!--MOCK NAME-->
            <v-text-field label="Name"
                          v-model="mock.name"
                          :rules="[v => !!v || 'Name is required']"
                          required outlined dense/>
            <!--MOCK GROUP-->
            <v-text-field label="Group"
                          v-model="mock.group"
                          placeholder="Examples: users, payments, auth, etc"
                          :rules="[v => !!v || 'Group is required']"
                          required outlined dense/>
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
                          required dense outlined/>
            <!--MOCK METHOD-->
            <v-select label="HTTP Method"
                      v-model="mock.method"
                      :items="httpMethods"
                      :rules="[v => !!v || 'HTTP Method is required']"
                      required outlined dense/>
            <!--MOCK STRATEGY-->
            <v-select label="Strategy"
                      v-model="mock.strategy"
                      :items="strategies"
                      :rules="[v => !!v || 'Strategy is required']"
                      v-on:change="updateResponses()"
                      required outlined dense/>
          </v-col>
          <v-col cols="6">
            <v-switch v-model="mock.status"
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
            <v-card-title class="px-0 py-0 pb-2">
              <!--RESPONSE DESCRIPTION-->
              <v-text-field :prefix="getResponseDescriptionPrefix(index)"
                            v-model="response.description"
                            placeholder="Type Response Description (Optional)"
                            class="description"/>
              <!--REMOVE RESPONSE BUTTON-->
              <v-btn icon color="red" @click="removeResponse(index)">
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
                              required outlined dense/>
                <!--RESPONSE STATUS CODE-->
                <v-text-field label="HTTP Status"
                              v-model.number="response.http_status"
                              placeholder="Examples: 200, 201, 400, 404, 500"
                              :rules="[v => (!isNaN(parseFloat(v)) && v >= 0) || 'HTTP Status is required and greater than or equal to 0']"
                              required outlined dense type="number"/>
                <!--RESPONSE DELAY-->
                <v-text-field label="Delay"
                              v-model.number="response.delay"
                              placeholder="Time to delay the response from server in milliseconds"
                              :rules="[v => (!isNaN(parseFloat(v)) && v >= 0) || 'Delay is required and greater than or equal to 0']"
                              required outlined dense type="number"/>
                <!--RESPONSE SCENE-->
                <v-text-field label="Scene"
                              v-model="response.scene"
                              placeholder="Value of 'scene' variable when SCENE strategy is selected."
                              :rules="isResponseSceneRequired(mock) ? [v => !!v || 'Scene is required'] : []"
                              outlined dense
                              :disabled="!isResponseSceneRequired(mock)"
                              :required="isResponseSceneRequired(mock)"/>
              </v-col>
              <v-col cols="6">
                <v-textarea label="Body"
                            v-model="response.body"
                            :rules="[v => !!v || 'Body is required']"
                            required outlined dense
                            rows="8" height="238"/>
              </v-col>
            </v-row>
          </v-container>
        </v-card>
        <v-col cols="12" class="text-right">
          <v-btn depressed color="primary" @click="addResponse()">New Response</v-btn>
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
            <v-card-title class="px-0 py-0 pb-2">
              Variable {{ index + 1 }}
              <v-spacer/>
              <!--REMOVE VARIABLE BUTTON-->
              <v-btn icon color="red" @click="removeVariable(index)">
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
                          v-on:change="updateVariables()"
                          required outlined dense/>
              </v-col>
              <v-col cols="4">
                <!--VARIABLE NAME-->
                <v-text-field label="Name"
                              v-model="variable.name"
                              :rules="[v => !!v || 'Name is required']"
                              required outlined dense/>
              </v-col>
              <v-col cols="4">
                <!--VARIABLE KEY-->
                <v-text-field label="Key"
                              v-model="variable.key"
                              :rules="isVariableTypeRequired(variable) ? [v => !!v || 'Key is required'] : []"
                              outlined dense
                              :disabled="!isVariableTypeRequired(variable)"
                              :required="isVariableTypeRequired(variable)"
                />
              </v-col>
            </v-row>
          </v-container>
        </v-card>
        <v-col cols="12" class="text-right">
          <v-btn depressed color="primary" @click="addVariable()">New Variable</v-btn>
        </v-col>
      </v-container>
    </v-card>

    <v-col cols="12" class="text-right">
      <v-btn depressed color="error" class="mx-1" @click="submitDelete" v-if="theKey">Delete</v-btn>
      <v-btn depressed color="warning" class="mx-1" @click="resetForm">Reset</v-btn>
      <v-btn depressed color="primary" class="mx-1" @click="submit" :loading="saving">Save</v-btn>
    </v-col>

    <v-snackbar v-model="alert.show" :color="alert.color" :timeout="alert.timeout">{{ alert.text }}</v-snackbar>

    <div class="text-center">
      <v-dialog

          width="500"
      >
        <v-card>
          <v-card-title class="text-h5 grey lighten-2">
            Execution URL
          </v-card-title>

          <v-card-text class="text-h6" ref="executionURL">
            {{ executionURL.value }}
          </v-card-text>
        </v-card>
      </v-dialog>
    </div>

    <v-dialog
        transition="dialog-top-transition"
        max-width="600"
        v-model="executionURL.show"
    >
      <template>
        <v-card>
          <v-toolbar

              color="primary"
              dark
          >Mock execution URL:
          </v-toolbar>
          <v-card-text>
            <div class="text-h6 pa-12">{{ executionURL.value }}</div>
          </v-card-text>
          <v-card-actions class="justify-end">
            <v-btn
                text
                @click="copyExecutionURL()"
            >Done
            </v-btn>
          </v-card-actions>
        </v-card>
      </template>
    </v-dialog>


  </v-form>
</template>


<script>
import axios from "axios";

export default {
  name: "MockDetails",
  props: {
    theKey: {
      type: String,
      required: false,
    },
    theName: {
      type: String,
      required: false,
    },
  },
  title() {
    return this.theName ? this.theName : "New Mock";
  },
  data() {
    return {
      mock: {},
      httpMethods: [
        {text: "GET", value: "GET"},
        {text: "POST", value: "POST"},
        {text: "PUT", value: "PUT"},
        {text: "PATCH", value: "PATCH"},
        {text: "DELETE", value: "DELETE"},
        {text: "OPTIONS", value: "OPTIONS"},
        {text: "HEAD", value: "HEAD"},
      ],
      strategies: [
        {text: "Normal", value: "normal"},
        {text: "Scene", value: "scene"},
        {text: "Random", value: "random"},
      ],
      varTypes: [
        {text: "Body", value: "body"},
        {text: "Header", value: "header"},
        {text: "Query", value: "query"},
        {text: "Random", value: "random"},
        {text: "Hash", value: "hash"},
        {text: "Path", value: "path"},
      ],
      alert: {
        show: false,
        color: "green",
        text: "",
        timeout: "5000"
      },
      valid: false,
      loading: false,
      saving: false,
      executionURL: {
        value: "",
        show: false
      }
    };
  },
  methods: {
    baseURL() {
      if (process.env.NODE_ENV === 'production') {
        return "/mock-service/rules"
      }
      return "http://localhost:8080/mock-service/rules"
    },
    submit() {
      if (!this.$refs.form.validate()) {
        this.showAlert("Some fields are not valid!", "validation error");
        return;
      }

      if (this.theKey) {
        this.submitUpdate();
      } else {
        this.submitCreate();
      }
    },
    async submitCreate() {
      const confirmTitle = "Creating New Mock";
      const confirmMsg = "Please confirm you want to create this mock";
      const confirmation = await this.$confirm(confirmMsg, {title: confirmTitle, color: "primary"});
      if (confirmation) {
        this.createMock();
      }
    },
    async submitUpdate() {
      const confirmTitle = "Updating Mock: " + this.theKey;
      const confirmMsg = "Please confirm you want to update this mock";
      const confirmation = await this.$confirm(confirmMsg, {title: confirmTitle, color: "warning"});
      if (confirmation) {
        this.updateMock();
      }
    },
    async submitDelete() {
      const confirmTitle = "Deleting Mock: " + this.theKey;
      const confirmMsg = "Please confirm you want to delete this mock";
      const confirmation = await this.$confirm(confirmMsg, {title: confirmTitle, color: "error"});
      if (confirmation) {
        this.deleteMock();
      }
    },
    async resetForm() {
      const confirmTitle = "Reset Form";
      const confirmMsg = "All changes will be lost, are you sure?";
      const confirmation = await this.$confirm(confirmMsg, {title: confirmTitle, color: "warning"});
      if (confirmation) {
        this.$refs.form.reset();
        this.initialize();
      }
    },
    createMock() {
      this.saving = true;
      axios
          .post(this.baseURL(), this.mock, {
            headers: {
              "Content-Type": "application/json",
            },
          })
          .then((res) => {
            this.$router.push({name: 'MockDetails', params: {theKey: res.data.key, theName: res.data.name}});
          })
          .catch((err) => {
            this.showAlert("Error creating mock", err);
          }).finally(() => {
        this.saving = false;
      });
    },
    updateMock() {
      this.saving = true;
      axios
          .put(this.baseURL() + "/" + this.theKey,
              this.mock,
              {
                headers: {
                  "Content-Type": "application/json",
                },
              }
          )
          .then(() => {
            this.showAlert("Mock successfully updated!", null);
          })
          .catch((err) => {
            this.showAlert("Error updating mock", err)
          }).finally(() => {
        this.saving = false;
      });
    },
    deleteMock() {
      axios
          .delete(this.baseURL() + "/" + this.theKey)
          .then(() => {
            this.$router.push({name: 'ListMocks'});
          })
          .catch((err) => {
            this.showAlert("Error deleting mock!", err);
          });
    },
    showAlert(text, err) {
      this.alert = {text: text, color: err == null ? "green" : "red", show: true};
      console.log(err);
    },
    addVariable() {
      let newVar = {
        type: "body",
        name: "",
        key: "",
      };
      if (!this.mock.variables) {
        this.mock.variables = [newVar];
      } else {
        this.mock.variables.push(newVar);
      }
    },
    removeVariable(i) {
      this.mock.variables.splice(i, 1);
    },
    addResponse() {
      let newResponse = {
        body: "",
        content_type: "application/json",
        http_status: 200,
        delay: 0,
        scene: "",
      };
      if (!this.mock.variables) {
        this.mock.responses = [newResponse];
      } else {
        this.mock.responses.push(newResponse);
      }
    },
    removeResponse(i) {
      this.mock.responses.splice(i, 1);
    },
    updateResponses() {
      if (this.mock.strategy !== "scene") {
        this.mock.responses.forEach(r => {
          r.scene = ""
        });
      }
    },
    updateVariables() {
      this.mock.variables.forEach(v => {
        if (v.type !== "body" && v.type !== "query" && v.type !== "header" && v.type !== "path") {
          v.key = "";
        }
      });
    },
    isResponseSceneRequired(mock) {
      return mock.strategy === "scene";
    },
    isVariableTypeRequired(variable) {
      return variable.type === 'body' || variable.type === 'query' || variable.type === 'header' || variable.type === 'path';
    },
    getResponseDescriptionPrefix(index) {
      return "Response " + (index + 1).toString() + ": "
    },
    newMock() {
      return {
        key: "",
        group: "",
        name: "",
        path: "",
        strategy: "",
        method: "",
        status: "enabled",
        responses: [
          {
            description: "",
            body: "",
            content_type: "application/json",
            http_status: 200,
            delay: 0,
            scene: "",
          },
        ],
        variables: [],
      };
    },
    showExecutionURL() {
      let executionURL = this.getExecutionURL()
      if (executionURL === "") {
        this.showAlert("Path cannot be empty.", "error");
      } else {
        this.executionURL.value = executionURL
        this.executionURL.show = true
      }
    },
    getExecutionURL() {
      let path = this.mock.path;
      if (path === "") {
        return ""
      }

      if (!path.startsWith("/")) {
        path = "/" + path
      }

      return window.location.protocol + "//" + window.location.host + "/mock-service/mock" + path
    },
    copyExecutionURL() {
      // navigator.clipboard.writeText(this.executionURL.value)
      //     .then(() => {
      //       this.showAlert("Successfully copied Execution URL to clipboard", null);
      //     })
      //     .catch((err) => {
      //       this.showAlert("Error copying Execution URL to clipboard", err);
      //     });
      this.executionURL.show = false
    },
    initialize() {
      if (this.theKey) {
        this.loading = true;
        axios
            .get(this.baseURL() + "/" + this.theKey)
            .then((res) => {
              this.mock = res.data;
            }).catch((err) => {
          this.showAlert("Error getting mock info!", err);
        }).finally(() => {
          this.loading = false;
        });
      } else {
        this.mock = this.newMock();
        this.$refs.form.resetValidation();
      }
    },
  },
  async created() {
    this.initialize();
  },
  watch: {
    // will fire on route changes
    //'$route.params.id': function(val, oldVal){ // Same
    "$route.path": function (val, oldVal) {
      console.log(val + oldVal);
      this.initialize();
    },
  },
};
</script>

<style>
.v-input--reverse .v-input__slot {
  flex-direction: row-reverse;
  justify-content: flex-end;
}

.v-application--is-ltr,
.v-input--selection-controls__input {
  margin-right: 0;
  margin-left: 8px;
}

.v-application--is-rtl,
.v-input--selection-controls__input {
  margin-left: 0;
  margin-right: 8px;
}

.description > .v-input__control > .v-input__slot::before {
  display: none !important;
}

</style>
