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
                          @keyup.enter="search()"
                          variant="outlined" density="compact" clearable hide-details></v-text-field>
          </v-col>
          <!--PATH FILTER-->
          <v-col cols="12" md="3">
            <v-text-field label="Path"
                          v-model="filters.path"
                          @keyup.enter="search()"
                          variant="outlined" density="compact" clearable hide-details></v-text-field>
          </v-col>
          <!--STRATEGY FILTER-->
          <v-col cols="12" md="3">
            <v-select
                label="Strategy"
                v-model="filters.strategy"
                @keyup.enter="search()"
                :items="strategies"
                item-title="text"
                item-value="value"
                variant="outlined" density="compact" clearable hide-details
            ></v-select>
          </v-col>
          <!--HTTP METHOD FILTER-->
          <v-col cols="12" md="3">
            <v-select
                label="HTTP Method"
                v-model="filters.method"
                @keyup.enter="search()"
                :items="httpMethods"
                item-title="text"
                item-value="value"
                variant="outlined" density="compact" clearable hide-details
            ></v-select>
          </v-col>
          <!--SEARCH BUTTONS-->
          <v-col cols="12" class="text-right">
            <v-btn variant="flat" color="grey-lighten-2" class="mr-2" @click="reset()">Reset</v-btn>
            <v-btn variant="flat" color="primary" @click="search()">Search</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- RESULT -->
    <br>
    <v-data-table-server
        class="elevation-2" density="compact"
        :headers="table.columns"
        :items="table.rows"
        :items-length="table.total"
        :loading="table.loading"
        @update:options="handleOptionsUpdate"
        hover
    >
      <template v-slot:item.status="{ item }">
        <v-switch v-model="(item.raw || item).status" color="info" true-value="enabled" false-value="disabled" hide-details density="compact" style="margin: 0"
                  @update:model-value="callStatus(item.raw || item)"></v-switch>
      </template>
      <template v-slot:item.name="{ item }">
        <router-link :to="{name: 'MockDetails', params:{theKey:(item.raw || item).key, theName:(item.raw || item).name}}">
          <span class="mr-2">{{ (item.raw || item).name }}</span>
        </router-link>
      </template>
      <template v-slot:item.method="{ item }">
        <span :class="getHTTPMethodColor((item.raw || item).method)">{{ (item.raw || item).method }}</span>
      </template>
      <template v-slot:item.delete="{ item }">
        <v-btn icon color="red" variant="text" density="compact" @click="callDelete(item.raw || item)">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table-server>

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
        {text: "GET", value: "GET", color: "text-blue"},
        {text: "POST", value: "POST", color: "text-green"},
        {text: "PUT", value: "PUT", color: "text-orange"},
        {text: "PATCH", value: "PATCH", color: "text-purple"},
        {text: "DELETE", value: "DELETE", color: "text-red"},
        {text: "OPTIONS", value: "OPTIONS", color: "text-grey"},
        {text: "HEAD", value: "HEAD", color: "text-black"}
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
          {title: "Enabled", key: "status"},
          {title: "Name", key: "name"},
          {title: "Group", key: "group"},
          {title: "Path", key: "path", width: "35%"},
          {title: "Strategy", key: "strategy"},
          {title: "Method", key: "method"},
          {title: "", key: "delete", width: "1%"},
        ],
        rows: [],
        footer: {
          showFirstLastPage: false,
          "items-per-page-options": [10, 30, 50],
        },
        total: 0,
        loading: true,
      },
      options: {
        page: 1,
        itemsPerPage: 10
      },
      alert: {
        show: false,
        color: "green",
        text: ""
      },
    };
  },
  methods: {
    handleOptionsUpdate(newOptions) {
      this.options = newOptions;
    },
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
      const confirmMsg = confirmTitle + "\n\nPlease confirm you want to delete this mock";
      const confirmation = window.confirm(confirmMsg);
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