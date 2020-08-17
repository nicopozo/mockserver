<template>
  <div>
    <br />
    <b-container>
      <!-- V-CARD SEARCH-->
      <b-card class="mt-3" header="Filters:">
        <b-row align-h="between">
          <b-col cols="0">
            <div class="form-inline">
              <label class="mr-sm-2" for="inline-form-custom-select-pref">Application:</label>
              <b-form-input
                class="mb-2 mr-sm-2 mb-sm-0"
                v-model="searchFields.application"
                v-on:keyup.enter="applyFilters"
              ></b-form-input>
              <label class="mr-sm-2" for="inline-form-custom-select-pref">Path:</label>
              <b-form-input
                class="mb-2 mr-sm-2 mb-sm-0"
                v-model="searchFields.path"
                v-on:keyup.enter="applyFilters"
              ></b-form-input>
              <label class="mr-sm-2" for="inline-form-custom-select-pref">Strategy:</label>
              <b-form-select
                v-model="searchFields.strategy"
                class="mb-2 mr-sm-2 mb-sm-0"
                @change="applyFilters"
              >
                <option
                  v-for="strategy in strategies"
                  :key="strategy.text"
                  :value="strategy.value"
                >{{ strategy.text }}</option>
              </b-form-select>
              <label class="mr-sm-2" for="inline-form-custom-select-pref">HTTP Method:</label>
              <b-form-select
                v-model="searchFields.httpMethod"
                class="mb-2 mr-sm-2 mb-sm-0"
                @change="applyFilters"
              >
                <option
                  v-for="httpMethod in httpMethods"
                  :key="httpMethod.text"
                  :value="httpMethod.value"
                >{{ httpMethod.text }}</option>
              </b-form-select>
            </div>
          </b-col>
          <b-col cols="200">
            <div class="form-inline">
              <b-button pill variant="outline-danger" v-on:click="resetFilters()">Reset</b-button>
              <label class="mr-sm-2" for="inline-form-custom-select-pref"></label>
              <b-button pill variant="outline-primary" v-on:click="applyFilters()">Apply</b-button>
            </div>
          </b-col>
        </b-row>
      </b-card>
      <br />

      <!-- MOCKS TABLE -->
      <b-table
        id="mocks-table"
        striped
        small
        :fields="fields"
        :items="mocks.results"
        responsive="sm"
      >
        <template v-slot:cell(name)="data">
          <router-link
            :to="{name: 'MockDetails', params:{theKey:data.item.key}}"
          >{{ data.item.name }}</router-link>
        </template>
        <template v-slot:cell(application)="data">{{ data.item.application }}</template>
        <template v-slot:cell(path)="data">{{ data.item.path }}</template>
        <template v-slot:cell(strategy)="data">{{ data.item.strategy }}</template>

        <template v-slot:cell(method)="data">
          <b v-if="data.item.method === 'GET'" style="color: #436f8a;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'POST'" style="color: #96bb7c;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'PUT'" style="color: #ffbd69;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'PATCH'" style="color: #45046a;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'DELETE'" style="color: #e84a5f;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'OPTIONS'" style="color: #b5076b;">{{ data.item.method }}</b>
          <b v-if="data.item.method === 'HEAD'" style="color: #342b38;">{{ data.item.method }}</b>
        </template>

        <!-- DROPDOWN ENABLED/DISABLED -->
        <template v-slot:cell(status)="data">
          <b-dropdown
            :text="data.item.status"
            variant="outline-success"
            size="sm"
            v-if="data.item.status === 'enabled'"
          >
            <b-dropdown-item
              href="#"
              v-on:click="submitUpdateStatus(data.item.key, 'disabled')"
            >Disable</b-dropdown-item>
          </b-dropdown>

          <b-dropdown
            :text="data.item.status"
            variant="outline-danger"
            size="sm"
            v-if="data.item.status === 'disabled'"
          >
            <b-dropdown-item
              href="#"
              v-on:click="submitUpdateStatus(data.item.key, 'enabled')"
            >Enable</b-dropdown-item>
          </b-dropdown>
        </template>

        <!-- BUTTON DELETE RULE -->
        <template v-slot:cell(delete)="data">
          <b-button
            pill
            variant="outline-danger"
            size="sm"
            v-on:click="submitDelete(data.item.key)"
          >Delete</b-button>
        </template>
      </b-table>

      <!-- PAGINATION -->
      <div class="mt-3">
        <b-pagination
          v-model="paging.page"
          :total-rows="paging.total"
          :per-page="paging.limit"
          align="center"
          aria-controls="mocks-table"
          @change="pageChange"
        >></b-pagination>
      </div>
    </b-container>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ListMocks",
  props: {},
  data() {
    return {
      fields: [
        "name",
        "application",
        "path",
        "strategy",
        "method",
        "status",
        "delete"
      ],
      mocks: [],
      httpMethods: [
        { text: "", value: null },
        { text: "GET", value: "GET" },
        { text: "POST", value: "POST" },
        { text: "PUT", value: "PUT" },
        { text: "PATCH", value: "PATCH" },
        { text: "DELETE", value: "DELETE" },
        { text: "OPTIONS", value: "OPTIONS" },
        { text: "HEAD", value: "HEAD" }
      ],
      strategies: [
        { text: "", value: null },
        { text: "Normal", value: "normal" }
      ],
      searchFields: {
        httpMethod: "",
        application: "",
        path: "",
        strategy: ""
      },
      paging: {
        offset: 0,
        limit: 15,
        total: 0,
        page: 0
      }
    };
  },
  methods: {
    search() {
      axios
        .get("http://localhost:8081/mock-server/rules", {
            params: this.queryParams()
          }
        )
        .then(res => {
          this.mocks = res.data;
          this.paging.total = res.data.paging.total;
          this.paging.page =
            Math.floor(this.paging.offset / this.paging.limit) + 1;
        })
        .catch(err => {
          console.log(err);
        });
    },
    submitDelete(theKey) {
      var confirmMsg = "Please confirm you want to delete this mock";
      var confirmTitle = "Deleting Mock: " + theKey;

      this.$bvModal
        .msgBoxConfirm(confirmMsg, {
          title: confirmTitle,
          okTitle: "OK",
          cancelTitle: "Cancel",
          centered: true
        })
        .then(ok => {
          if (ok) {
            this.deleteMock(theKey);
          }
        })
        .catch(err => {
          console.log(err);
        });
    },
    submitUpdateStatus(theKey, newStatus) {
      var confirmMsg = "Please confirm you want to update this mock´s Status";
      var confirmTitle = "New status: " + newStatus;

      this.$bvModal
        .msgBoxConfirm(confirmMsg, {
          title: confirmTitle,
          okTitle: "OK",
          cancelTitle: "Cancel",
          centered: true
        })
        .then(ok => {
          if (ok) {
            this.updateMockStatus(theKey, newStatus);
          }
        })
        .catch(err => {
          console.log(err);
        });
    },
    deleteMock(theKey) {
      axios
        .delete("http://localhost:8081/mock-server/rules/" + theKey)
        .then(res => {
          var msg = "Mock successfully deleted!";
          var title = "Success!!";
          this.showSuccessModal(title, msg, true);
          console.log(res);
        })
        .catch(err => {
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
          this.$bvModal.msgBoxOk(msg, {
            title: "Error deleting mock",
            okVariant: "danger",
            centered: true
          });
        });
    },
    updateMockStatus(theKey, newStatus) {
      const mocks = this.mocks;
      axios
        .put("http://localhost:8081/mock-server/rules/" + theKey + "/status",
          { status: newStatus },
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        )
        .then(res => {
          mocks.results.forEach(function(item) {
            if (item.key === theKey) {
              item.status = newStatus;
            }
          });
          var msg = "Mock status successfully updated!";
          var title = "Success!!";
          this.showSuccessModal(title, msg, false);
          console.log(res);
        })
        .catch(err => {
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
          this.$bvModal.msgBoxOk(msg, {
            title: "Error updating mock status",
            okVariant: "danger",
            centered: true
          });
        });
    },
    showSuccessModal(title, msg, goHome) {
      this.$bvModal
        .msgBoxOk(msg, {
          title: title,
          okVariant: "success",
          centered: true
        })
        .then(value => {
          if (goHome) {
            this.applyFilters();
          }
          console.log(value);
        })
        .catch(err => {
          console.log(err);
        });
    },
    queryParams() {
      var params = {
        limit: this.paging.limit,
        offset: this.paging.offset
      };

      if (this.searchFields.httpMethod) {
        params.method = this.searchFields.httpMethod;
      }

      if (this.searchFields.application) {
        params.application = this.searchFields.application;
      }

      if (this.searchFields.path) {
        params.path = this.searchFields.path;
      }

      if (this.searchFields.strategy) {
        params.strategy = this.searchFields.strategy;
      }

      return params;
    },
    applyFilters() {
      this.paging.page = 1;
      this.paging.offset = 0;
      this.search();
    },
    resetFilters() {
      this.searchFields.httpMethod = "";
      this.searchFields.application = "";
      this.searchFields.path = "";
      this.searchFields.strategy = "";
      this.applyFilters();
    },
    pageChange(page) {
      this.paging.page = page;
      this.paging.offset = (page - 1) * this.paging.limit;
      this.search();
    }
  },
  created() {
    this.search();
  }
};
</script>