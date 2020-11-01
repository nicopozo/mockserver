<template>
  <div>
    <br />
    <!-- CONTAINER FORM -->
    <b-container>
      <b-form @submit="submit">
        <!-- INPUT KEY -->
        <b-form-group
          id="input-group-key"
          label-cols="4"
          label-cols-lg="2"
          label="Key:"
          label-for="input-key"
          v-if="mock.key"
        >
          <b-form-input
            id="input-key"
            v-model="mock.key"
            disabled
          ></b-form-input>
        </b-form-group>

        <!-- INPUT NAME -->
        <b-form-group
          id="input-group-name"
          label-cols="4"
          label-cols-lg="2"
          label="Name:"
          label-for="input-name"
        >
          <b-form-input
            id="input-name"
            v-model="mock.name"
            required
            placeholder="Enter a name for your mock"
          ></b-form-input>
        </b-form-group>

        <!-- INPUT APPLICATION -->
        <b-form-group
          id="input-group-application"
          label-cols="4"
          label-cols-lg="2"
          label="Application:"
          label-for="input-application"
        >
          <b-form-input
            id="input-application"
            v-model="mock.application"
            required
            placeholder="Examples: core, payments, simetrik, movistar, etc"
          ></b-form-input>
        </b-form-group>

        <!-- INPUT PATH -->
        <b-form-group
          id="input-group-path"
          label="Path:"
          label-cols="4"
          label-cols-lg="2"
          label-for="input-path"
        >
          <b-form-input
            id="input-path"
            v-model="mock.path"
            required
            placeholder="Example: /users/{user_id}"
          ></b-form-input>
        </b-form-group>

        <!-- INPUT PATH -->
        <b-form-group
          id="input-group-method"
          required
          label-cols="4"
          label-cols-lg="2"
          label="Method:"
        >
          <b-form-select v-model="mock.method">
            <option
              v-for="httpMethod in httpMethods"
              :key="httpMethod.text"
              :value="httpMethod.value"
            >
              {{ httpMethod.text }}
            </option>
          </b-form-select>
        </b-form-group>

        <!-- INPUT STRATEGY -->
        <b-form-group
          id="input-group-strategy"
          required
          label-cols="4"
          label-cols-lg="2"
          label="Strategy:"
        >
          <b-form-select v-model="mock.strategy">
            <option
              v-for="strategy in strategies"
              :key="strategy.text"
              :value="strategy.value"
            >
              {{ strategy.text }}
            </option>
          </b-form-select>
        </b-form-group>

        <!-- RESPONSES -->
        <b-card class="mt-3" header="Responses:">
          <!-- FOR EACH RESPONSE -->
          <b-card
            class="mt-3"
            v-for="(response, index) in mock.responses"
            v-bind:key="response"
          >
            <!-- RESPONSE CONTENT TYPE-->
            <b-form-group
              id="input-group-content-type"
              label-cols="4"
              label-cols-lg="2"
              label="Content Type:"
              label-for="input-content-type"
            >
              <b-form-input
                id="input-content-type"
                v-model="response.content_type"
                required
                placeholder="Example: application/json"
              ></b-form-input>
            </b-form-group>

            <!-- RESPONSE HTTP STATUS-->
            <b-form-group
              id="input-group-status"
              label-cols="4"
              label-cols-lg="2"
              label="HTTP Status:"
              label-for="input-status"
            >
              <b-form-input
                id="input-status"
                type="number"
                number
                v-model="response.http_status"
                required
                placeholder="Examples: 200, 201, 400, 404, 500"
              ></b-form-input>
            </b-form-group>

            <!-- RESPONSE DELAY-->
            <b-form-group
              id="input-group-delay"
              label-cols="4"
              label-cols-lg="2"
              label="Delay Time:"
              label-for="input-delay"
            >
              <b-form-input
                id="input-delay"
                type="number"
                number
                v-model="response.delay"
                required
                value="0"
                placeholder="Time to delay the response from server in Miliseconds."
              ></b-form-input>
            </b-form-group>

            <!-- RESPONSE BODY-->
            <b-form-group
              id="input-group-body"
              label-cols="4"
              label-cols-lg="2"
              label="body:"
              label-for="input-body"
            >
              <b-form-textarea
                id="input-group-body"
                v-model="response.body"
                placeholder="Add a response body..."
                rows="3"
                max-rows="6"
              ></b-form-textarea>
            </b-form-group>

            <!-- RESPONSE SCENE-->
            <b-form-group
              id="input-group-scene"
              label-cols="4"
              label-cols-lg="2"
              label="Scene:"
              label-for="input-scene"
            >
              <b-form-input
                id="input-scene"
                v-model="response.scene"
                required
                placeholder="Example: name"
              ></b-form-input>
            </b-form-group>

            <b-col cols="100">
              <b-button
                pill
                variant="outline-danger"
                v-on:click="removeResponse(index)"
                >Remove</b-button
              >
            </b-col>
          </b-card>

          <br />

          <!-- ADD NEW VARIABLE BUTTON-->
          <b-container>
            <b-row align-h="end">
              <b-col cols="100">
                <b-button
                  pill
                  variant="outline-primary"
                  v-on:click="addResponse()"
                  >New Response</b-button
                >
              </b-col>
            </b-row>
          </b-container>
        </b-card>

        <!-- V-CARD VARIABLES -->
        <b-card class="mt-3" header="Variables:">
          <!-- FOR EACH VARIABLE -->
          <b-card
            class="mt-3"
            v-for="(variable, index) in mock.variables"
            v-bind:key="variable"
          >
            <b-row align-h="between">
              <b-col cols="200">
                <div class="form-inline">
                  <label class="mr-sm-2" for="inline-form-custom-select-pref"
                    >Type:</label
                  >
                  <b-form-select
                    v-model="variable.type"
                    @change="variable.key = null"
                    class="mb-2 mr-sm-2 mb-sm-0"
                  >
                    <option
                      v-for="varType in varTypes"
                      :key="varType.text"
                      :value="varType.value"
                    >
                      {{ varType.text }}
                    </option>
                  </b-form-select>
                  <label class="mr-sm-2" for="inline-form-custom-select-pref"
                    >Name:</label
                  >
                  <b-form-input
                    class="mb-2 mr-sm-2 mb-sm-0"
                    v-model="variable.name"
                    required
                  ></b-form-input>
                  <label class="mr-sm-2" for="inline-form-custom-select-pref"
                    >Value:</label
                  >
                  <b-form-input
                    :disabled="
                      variable.type != 'body' &&
                      variable.type != 'query' &&
                      variable.type != 'header'
                    "
                    :required="
                      variable.type === 'body' ||
                      variable.type === 'query' ||
                      variable.type === 'header'
                    "
                    class="mb-2 mr-sm-2 mb-sm-0"
                    v-model="variable.key"
                  ></b-form-input>
                </div>
              </b-col>
              <b-col cols="100">
                <b-button
                  pill
                  variant="outline-danger"
                  v-on:click="removeVariable(index)"
                  >Remove</b-button
                >
              </b-col>
            </b-row>
          </b-card>

          <br />

          <!-- ADD NEW VARIABLE BUTTON-->
          <b-container>
            <b-row align-h="end">
              <b-col cols="100">
                <b-button
                  pill
                  variant="outline-primary"
                  v-on:click="addVariable()"
                  >New Variable</b-button
                >
              </b-col>
            </b-row>
          </b-container>
        </b-card>

        <br />

        <!-- SUBMIT BUTTON-->
        <b-container>
          <b-row align-h="end">
            <b-col cols="1000">
              <div class="form-inline">
                <b-button variant="info" v-on:click="resetForm()"
                  >Reset</b-button
                >
                <label
                  class="mr-sm-2"
                  for="inline-form-custom-select-pref"
                ></label>
                <b-button
                  variant="danger"
                  v-on:click="submitDelete()"
                  v-if="theKey"
                  >Delete</b-button
                >
                <label
                  class="mr-sm-2"
                  for="inline-form-custom-select-pref"
                  v-if="theKey"
                ></label>
                <b-button type="submit" variant="primary">Submit</b-button>
              </div>
            </b-col>
          </b-row>
        </b-container>
      </b-form>
    </b-container>
    <br />
  </div>
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
  },
  data() {
    return {
      mock: {},
      httpMethods: [
        { text: "Select One", value: null },
        { text: "GET", value: "GET" },
        { text: "POST", value: "POST" },
        { text: "PUT", value: "PUT" },
        { text: "PATCH", value: "PATCH" },
        { text: "DELETE", value: "DELETE" },
        { text: "OPTIONS", value: "OPTIONS" },
        { text: "HEAD", value: "HEAD" },
      ],

      varTypes: [
        { text: "Select One", value: null },
        { text: "BODY", value: "body" },
        { text: "HEADER", value: "header" },
        { text: "QUERY", value: "query" },
        { text: "RANDOM", value: "random" },
        { text: "HASH", value: "hash" },
        { text: "PATH", value: "path" },
      ],

      strategies: [
        { text: "Select One", value: null },
        { text: "NORMAL", value: "normal" },
        { text: "SCENE", value: "scene" },
        { text: "RAMDOM", value: "random" },
      ],
    };
  },
  methods: {
    submit(ev) {
      var confirmMsg = "";
      var confirmTitle = "";
      if (this.theKey) {
        confirmMsg = "Please confirm you want to update this mock";
        confirmTitle = "Updating Mock: " + this.theKey;
      } else {
        confirmMsg = "Please confirm you want to  create this mock";
        confirmTitle = "Creating New Mock";
      }

      this.$bvModal
        .msgBoxConfirm(confirmMsg, {
          title: confirmTitle,
          okTitle: "OK",
          cancelTitle: "Cancel",
          centered: true,
        })
        .then((ok) => {
          if (ok) {
            if (this.theKey) {
              this.updateMock();
            } else {
              this.createMock();
            }
          }
        })
        .catch((err) => {
          console.log(err);
        });
      ev.preventDefault();
    },
    submitDelete() {
      var confirmMsg = "Please confirm you want to delete this mock";
      var confirmTitle = "Deleting Mock: " + this.theKey;

      this.$bvModal
        .msgBoxConfirm(confirmMsg, {
          title: confirmTitle,
          okTitle: "OK",
          cancelTitle: "Cancel",
          centered: true,
        })
        .then((ok) => {
          if (ok) {
            this.deleteMock();
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    resetForm() {
      var confirmMsg = "All changes will be lost, are you sure?";
      var confirmTitle = "Reset Form";

      this.$bvModal
        .msgBoxConfirm(confirmMsg, {
          title: confirmTitle,
          okTitle: "OK",
          cancelTitle: "Cancel",
          centered: true,
        })
        .then((ok) => {
          if (ok) {
            this.initialize();
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    updateMock() {
      axios
        .put(
          "http://localhost:8081/mock-server/rules/" + this.theKey,
          this.mock,
          {
            headers: {
              "Content-Type": "application/json",
            },
          }
        )
        .then((response) => {
          var title = "Success!!";
          var msg = "Mock successfully updated!";
          this.showSuccessModal(title, msg, false);
          console.log(response);
        })
        .catch((err) => {
          var msg = err.message;
          if (
            typeof err.response !== "undefined" &&
            err.response.data !== "undefined" &&
            err.response.data.message !== "undefined" &&
            err.response.data.cause[0] !== "undefined" &&
            err.response.data.cause[0].description !== "undefined"
          ) {
            msg =
              err.response.data.cause[0].description +
              " - " +
              err.response.data.message;
          }
          this.showErrorModal("Error updating mock", msg);
        });
    },
    createMock() {
      axios
        .post("http://localhost:8081/mock-server/rules", this.mock, {
          headers: {
            "Content-Type": "application/json",
          },
        })
        .then((resp) => {
          var msg = "Mock successfully created!";
          var title = "Success!!";
          this.showSuccessModal(title, msg, true);
          console.log(resp);
        })
        .catch((err) => {
          var msg = err.message;
          if (
            typeof err.response !== "undefined" &&
            err.response.data !== "undefined" &&
            err.response.data.message !== "undefined" &&
            err.response.data.cause[0] !== "undefined" &&
            err.response.data.cause[0].description !== "undefined"
          ) {
            msg =
              err.response.data.cause[0].description +
              " - " +
              err.response.data.message;
          }
          this.showErrorModal("Error creating mock", msg);
        });
    },
    deleteMock() {
      axios
        .delete("http://localhost:8081/mock-server/rules/" + this.theKey)
        .then((res) => {
          var msg = "Mock successfully deleted! ";
          var title = "Success!!";
          this.showSuccessModal(title, msg, true);
          console.log(res);
        })
        .catch((err) => {
          var msg = err.message;
          if (
            typeof err.response !== "undefined" &&
            err.response.data !== "undefined" &&
            err.response.data.message !== "undefined" &&
            err.response.data.cause[0] !== "undefined" &&
            err.response.data.cause[0].description !== "undefined"
          ) {
            msg =
              err.response.data.cause[0].description +
              " - " +
              err.response.data.message;
          }
          var title = "ERROR!";

          this.showErrorModal(title, msg);
        });
    },
    showSuccessModal(title, msg, goHome) {
      const router = this.$router;
      this.$bvModal
        .msgBoxOk(msg, {
          title: title,
          okVariant: "success",
          centered: true,
        })
        .then((value) => {
          if (goHome) {
            router.push({ name: "ListMocks" });
          }
          console.log(value);
        })
        .catch((err) => {
          console.log(err);
        });
    },
    showErrorModal(title, msg) {
      this.$bvModal.msgBoxOk(msg, {
        title: title,
        okVariant: "danger",
        centered: true,
      });
    },
    addVariable() {
      var newVar = {
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
      var newResponse = {
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
    newMock() {
      return {
        key: "",
        application: "",
        name: "",
        path: "",
        strategy: "",
        method: "",
        status: "",
        responses: [
          {
            body: "",
            content_type: "",
            http_status: "",
            delay: 0,
            scene: "",
          },
        ],
        variables: [],
      };
    },
    initialize() {
      if (this.theKey) {
        axios
          .get("http://localhost:8081/mock-server/rules/" + this.theKey)
          .then((res) => {
            this.mock = res.data;
          })
          .catch((err) => {
            var msg = err.message;
            if (
              typeof err.response !== "undefined" &&
              err.response.data !== "undefined" &&
              err.response.data.message !== "undefined" &&
              err.response.data.cause[0] !== "undefined" &&
              err.response.data.cause[0].description !== "undefined"
            ) {
              msg =
                err.response.data.cause[0].description +
                " - " +
                err.response.data.message;
            }

            this.showErrorModal("ERROR!", msg);
          });
      } else {
        this.mock = this.newMock();
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