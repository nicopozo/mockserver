<template>
  <div>
    <!--SEARCH-->
    <v-card class="elevation-2">
      <v-container fluid>
        <v-row>
          <!--GROUP FILTER-->
          <v-col cols="12" md="3">
            <v-text-field label="Group"
                          v-model="filters.group"
                          @keyup.enter.native="search()"
                          outlined dense clearable hide-details></v-text-field>
          </v-col>
          <!--PATH FILTER-->
          <v-col cols="12" md="3">
            <v-text-field label="Path"
                          v-model="filters.path"
                          @keyup.enter.native="search()"
                          outlined dense clearable hide-details></v-text-field>
          </v-col>
          <!--STRATEGY FILTER-->
          <v-col cols="12" md="3">
            <v-select
                label="Strategy"
                v-model="filters.strategy"
                @keyup.enter.native="search()"
                :items="strategies"
                outlined dense clearable hide-details
            ></v-select>
          </v-col>
          <!--HTTP METHOD FILTER-->
          <v-col cols="12" md="3">
            <v-select
                label="HTTP Method"
                v-model="filters.method"
                @keyup.enter.native="search()"
                :items="httpMethods"
                outlined dense clearable hide-details
            ></v-select>
          </v-col>
          <!--SEARCH BUTTONS-->
          <v-col cols="12" class="text-right">
            <v-btn depressed class="mr-2" @click="reset()">Reset</v-btn>
            <v-btn depressed color="primary" @click="search()">Search</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- RESULT -->
    <br>
    <v-data-table
        class="elevation-2" dense
        :headers="table.columns"
        :items="table.rows"
        :footer-props="table.footer"
        :server-items-length="table.total"
        :loading="table.loading"
        :options.sync="options"
        disable-sort
    >
      <template v-slot:item.status="{ item }">
        <v-switch v-model="item.status" true-value="enabled" false-value="disabled" hide-details style="margin: 0"
                  @change="callStatus(item)"></v-switch>
      </template>
      <template v-slot:item.name="{ item }">
        <router-link :to="{name: 'MockDetails', params:{theKey:item.key, theName:item.name}}">
          <span class="mr-2">{{ item.name }}</span>
        </router-link>
      </template>
      <template v-slot:item.method="{ item }">
        <span :class="getHTTPMethodColor(item.method)">{{ item.method }}</span>
      </template>
      <template v-slot:item.delete="{ item }">
        <v-btn icon color="red" @click="callDelete(item)">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table>

    <v-snackbar v-model="alert.show" :color="alert.color">{{ alert.text }}</v-snackbar>

  </div>
</template>

<script>
import axios from "axios";

export default {
  title() {
    return "Mocks";
  },
  data() {
    return {
      httpMethods: [
        {text: "GET", value: "GET", color: "blue--text"},
        {text: "POST", value: "POST", color: "green--text"},
        {text: "PUT", value: "PUT", color: "orange--text"},
        {text: "PATCH", value: "PATCH", color: "violet--text"},
        {text: "DELETE", value: "DELETE", color: "red--text"},
        {text: "OPTIONS", value: "OPTIONS", color: "gray--text"},
        {text: "HEAD", value: "HEAD", color: "black--text"}
      ],
      strategies: [
        {text: "Normal", value: "normal"},
        {text: "Scene", value: "scene"},
        {text: "Random", value: "random"}
      ],
      filters: {
        group: null,
        path: null,
        strategy: null,
        method: null,
      },
      table: {
        columns: [
          {text: "Enabled", value: "status"},
          {text: "Name", value: "name"},
          {text: "Group", value: "group"},
          {text: "Path", value: "path", width: "35%"},
          {text: "Strategy", value: "strategy"},
          {text: "Method", value: "method"},
          {text: "", value: "delete", width: "1%"},
        ],
        rows: [],
        footer: {
          showFirstLastPage: false,
          "items-per-page-options": [10, 30, 50],
        },
        total: 0,
        loading: true,
      },
      options: {},
      alert: {
        show: false,
        color: "green",
        text: ""
      },
    };
  },
  methods: {
    baseURL() {
      if (process.env.NODE_ENV === 'production') {
        return "/mock-service/rules"
      }
      return "http://localhost:8080/mock-service/rules"
    },
    getHTTPMethodColor(httpMethod) {
      return this.httpMethods.find(x => x.value == httpMethod).color;
    },
    queryParams() {
      const {
        group,
        path,
        strategy,
        method,
      } = this.filters;

      const {
        page,
        itemsPerPage,
      } = this.options;

      let params = {
        limit: itemsPerPage,
        offset: (page - 1) * itemsPerPage,
      };

      if (group) {
        params.group = group;
      }

      if (path) {
        params.path = path;
      }

      if (strategy) {
        params.strategy = strategy;
      }

      if (method) {
        params.method = method;
      }

      return params;
    },
    reset() {
      this.filters = {
        group: null,
        path: null,
        strategy: null,
        method: null,
      };
      this.search();
    },
    search() {
      this.options = {
        page: 1,
        itemsPerPage: 10,
      };
    },
    showAlert(text, err) {
      this.alert = {text: text, color: err == null ? "green" : "red", show: true};
      console.log(err);
    },
    //API calls
    restSearch() {
      this.table.loading = true;
      axios
          .get(this.baseURL(), {
            params: this.queryParams(),
          })
          .then((res) => {
            this.table.rows = res.data.results;
            this.table.total = res.data.paging.total;
          })
          .catch((err) => {
            this.table.rows = [];
            this.table.total = 0;

            const msg = "Something went wrong searching mocks!";
            this.showAlert(msg, err)
          })
          .finally(() => {
            this.table.loading = false
          });
    },
    callStatus(item) {
      axios
        .put(
            this.baseURL() + "/" + item.key + "/status",
            {status: item.status},
            {
              headers: {
                "Content-Type": "application/json"
              }
            }
        ).then(() => {
          this.showAlert("Mock successfully " + item.status + "!", null);
        }).catch((err) => {
          item.status = item.status === "enabled" ? "disabled" : "enabled"; //rollback

          const msg = "Something went wrong updated mock status!";
          this.showAlert(msg, err)
        })
    },
    async callDelete(item) {
      const confirmTitle = "Deleting Mock: " + item.key;
      const confirmMsg = "Please confirm you want to delete this mock";
      const confirmation = await this.$confirm(confirmMsg, {title: confirmTitle, color: "error"});
      if (confirmation) {
        axios
            .delete(this.baseURL() + "/" + item.key)
            .then(() => {
              this.showAlert("Mock successfully deleted!", null);
            })
            .catch((err) => {
              this.showAlert("Error deleting mock!", err);
            }).finally(() => {
              this.search();
            });
      }
    }
  },
  watch: {
    options: {
      handler() {
        this.restSearch();
      },
      deep: true,
    },
  },
};
</script>