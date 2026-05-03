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
        <v-switch v-model="raw(item).status" color="info" true-value="enabled" false-value="disabled" hide-details density="compact" style="margin: 0"
                  @update:model-value="callStatus(raw(item))"></v-switch>
      </template>
      <template v-slot:item.name="{ item }">
        <router-link :to="{name: 'MockDetails', params:{theKey:raw(item).key, theName:raw(item).name}}">
          <span class="mr-2">{{ raw(item).name }}</span>
        </router-link>
      </template>
      <template v-slot:item.method="{ item }">
        <span :class="getHTTPMethodColor(raw(item).method)">{{ raw(item).method }}</span>
      </template>
      <template v-slot:item.delete="{ item }">
        <v-btn icon color="red" variant="text" density="compact" @click="callDelete(raw(item))">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table-server>

    <v-snackbar v-model="alert.show" :color="alert.color">{{ alert.text }}</v-snackbar>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import axios from "axios";
import type { Mock, PaginatedMocks } from '@/types';

const httpMethods = [
  {text: "GET", value: "GET", color: "text-blue"},
  {text: "POST", value: "POST", color: "text-green"},
  {text: "PUT", value: "PUT", color: "text-orange"},
  {text: "PATCH", value: "PATCH", color: "text-purple"},
  {text: "DELETE", value: "DELETE", color: "text-red"},
  {text: "OPTIONS", value: "OPTIONS", color: "text-grey"},
  {text: "HEAD", value: "HEAD", color: "text-black"}
];

const strategies = [
  {text: "Normal", value: "normal"},
  {text: "Scene", value: "scene"},
  {text: "Random", value: "random"}
];

const filters = reactive({
  group: null as string | null,
  path: null as string | null,
  strategy: null as string | null,
  method: null as string | null,
});

const table = reactive({
  columns: [
    {title: "Enabled", key: "status"},
    {title: "Name", key: "name"},
    {title: "Group", key: "group"},
    {title: "Path", key: "path", width: "35%"},
    {title: "Strategy", key: "strategy"},
    {title: "Method", key: "method"},
    {title: "", key: "delete", width: "1%"},
  ] as any[],
  rows: [] as Mock[],
  footer: {
    showFirstLastPage: false,
    "items-per-page-options": [10, 30, 50],
  },
  total: 0,
  loading: true,
});

const options = ref({
  page: 1,
  itemsPerPage: 10
});

const alert = reactive({
  show: false,
  color: "green",
  text: ""
});

function handleOptionsUpdate(newOptions: any) {
  options.value = newOptions;
}

function raw(item: any): Mock {
  return item.raw || item;
}

function baseURL() {
  if (import.meta.env.PROD) {
    return "/mock-service/rules"
  }
  return "http://localhost:8080/mock-service/rules"
}

function getHTTPMethodColor(httpMethod: string) {
  return httpMethods.find(x => x.value == httpMethod)?.color || 'text-black';
}

function queryParams() {
  const { group, path, strategy, method } = filters;
  const { page, itemsPerPage } = options.value;

  let params: any = {
    limit: itemsPerPage,
    offset: (page - 1) * itemsPerPage,
  };

  if (group) params.group = group;
  if (path) params.path = path;
  if (strategy) params.strategy = strategy;
  if (method) params.method = method;

  return params;
}

function reset() {
  filters.group = null;
  filters.path = null;
  filters.strategy = null;
  filters.method = null;
  search();
}

function search() {
  options.value = {
    page: 1,
    itemsPerPage: 10,
  };
}

function showAlert(text: string, err?: any) {
  alert.text = text;
  alert.color = err == null ? "green" : "red";
  alert.show = true;
  if (err) console.error(err);
}

function restSearch() {
  table.loading = true;
  axios
      .get<PaginatedMocks>(baseURL(), {
        params: queryParams(),
      })
      .then((res) => {
        table.rows = res.data.results;
        table.total = res.data.paging.total;
      })
      .catch((err) => {
        table.rows = [];
        table.total = 0;
        showAlert("Something went wrong searching mocks!", err)
      })
      .finally(() => {
        table.loading = false
      });
}

function callStatus(item: Mock) {
  axios
    .put(
        baseURL() + "/" + item.key + "/status",
        {status: item.status},
        {
          headers: {
            "Content-Type": "application/json"
          }
        }
    ).then(() => {
      showAlert("Mock successfully " + item.status + "!");
    }).catch((err) => {
      item.status = item.status === "enabled" ? "disabled" : "enabled"; //rollback
      showAlert("Something went wrong updated mock status!", err)
    })
}

async function callDelete(item: Mock) {
  const confirmTitle = "Deleting Mock: " + item.key;
  const confirmMsg = confirmTitle + "\n\nPlease confirm you want to delete this mock";
  const confirmation = window.confirm(confirmMsg);
  if (confirmation) {
    axios
        .delete(baseURL() + "/" + item.key)
        .then(() => {
          showAlert("Mock successfully deleted!");
        })
        .catch((err) => {
          showAlert("Error deleting mock!", err);
        }).finally(() => {
          search();
        });
  }
}

watch(options, () => {
  restSearch();
}, { deep: true });

</script>