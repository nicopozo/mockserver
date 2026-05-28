<template>
  <div class="list-mocks-container">
    <!--SEARCH-->
    <v-card class="search-card mb-6" elevation="0">
      <v-container fluid class="pa-6">
        <v-row align="center">
          <!--GROUP FILTER-->
          <v-col cols="12" md="3">
            <v-text-field label="Group"
                          v-model="filters.group"
                          @keyup.enter="search()"
                          variant="solo-filled" density="comfortable" flat clearable hide-details
                          prepend-inner-icon="mdi-folder-outline"
                          class="custom-field"></v-text-field>
          </v-col>
          <!--PATH FILTER-->
          <v-col cols="12" md="3">
            <v-text-field label="Path"
                          v-model="filters.path"
                          @keyup.enter="search()"
                          variant="solo-filled" density="comfortable" flat clearable hide-details
                          prepend-inner-icon="mdi-link-variant"
                          class="custom-field"></v-text-field>
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
                variant="solo-filled" density="comfortable" flat clearable hide-details
                prepend-inner-icon="mdi-layers-outline"
                class="custom-field"
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
                variant="solo-filled" density="comfortable" flat clearable hide-details
                prepend-inner-icon="mdi-api"
                class="custom-field"
            ></v-select>
          </v-col>
          <!--SEARCH BUTTONS-->
          <v-col cols="12" class="d-flex align-center justify-end ga-3 pt-4">
            <v-btn variant="text" color="grey-darken-1" @click="reset()" class="text-none">
              <v-icon start>mdi-refresh</v-icon> Reset
            </v-btn>
            <v-btn color="primary" @click="search()" class="text-none px-6" elevation="2">
              <v-icon start>mdi-magnify</v-icon> Search
            </v-btn>
            <v-divider vertical class="mx-2"></v-divider>
            <v-btn variant="tonal" color="secondary" prepend-icon="mdi-download" @click="exportMocks()" class="text-none">Export</v-btn>
            <v-btn variant="tonal" color="secondary" prepend-icon="mdi-upload" @click="triggerFileInput()" class="text-none">Import</v-btn>
            <input type="file" ref="fileInput" hidden accept=".json" @change="importMocks" />
          </v-col>
        </v-row>
      </v-container>
    </v-card>

    <!-- RESULT -->
    <v-card class="table-card" elevation="4">
      <v-data-table-server
          density="comfortable"
          :headers="table.columns"
          :items="table.rows"
          :items-length="table.total"
          :loading="table.loading"
          @update:options="handleOptionsUpdate"
          hover
          class="custom-table"
      >
        <template v-slot:item.status="{ item }">
          <div class="d-flex align-center">
            <v-switch v-model="raw(item).status"
                      color="success"
                      true-value="enabled"
                      false-value="disabled"
                      hide-details
                      density="compact"
                      @update:model-value="callStatus(raw(item))"
                      class="status-switch"
            ></v-switch>
            <span class="text-caption ml-2 text-uppercase font-weight-bold" :class="raw(item).status === 'enabled' ? 'text-success' : 'text-grey'">
              {{ raw(item).status }}
            </span>
          </div>
        </template>

        <template v-slot:item.name="{ item }">
          <div class="d-flex align-center">
            <v-avatar color="primary" variant="tonal" size="32" class="mr-3">
              <v-icon size="18">mdi-file-code-outline</v-icon>
            </v-avatar>
            <router-link :to="{name: 'MockDetails', params:{theKey:raw(item).key, theName:raw(item).name}}" class="mock-link">
              <span class="font-weight-bold">{{ raw(item).name }}</span>
            </router-link>
          </div>
        </template>

        <template v-slot:item.path="{ item }">
          <code class="path-code">{{ raw(item).path }}</code>
        </template>

        <template v-slot:item.strategy="{ item }">
          <v-chip size="small" variant="outlined" color="primary" class="text-uppercase font-weight-bold">
            {{ raw(item).strategy }}
          </v-chip>
        </template>

        <template v-slot:item.method="{ item }">
          <v-chip :color="getHTTPMethodColor(raw(item).method)" size="small" variant="flat" class="font-weight-black px-3">
            {{ raw(item).method }}
          </v-chip>
        </template>

        <template v-slot:item.delete="{ item }">
          <v-btn icon color="error" variant="text" density="comfortable" @click="callDelete(raw(item))">
            <v-icon size="20">mdi-delete-outline</v-icon>
            <v-tooltip activator="parent" location="top">Delete Mock</v-tooltip>
          </v-btn>
        </template>

        <template v-slot:no-data>
          <div class="py-10 text-center">
            <v-icon size="64" color="grey-lighten-1">mdi-database-off-outline</v-icon>
            <p class="text-h6 text-grey-darken-1 mt-4">No mocks found matching your criteria</p>
            <v-btn variant="text" color="primary" @click="reset()">Clear Filters</v-btn>
          </div>
        </template>
      </v-data-table-server>
    </v-card>

    <v-snackbar v-model="alert.show" :color="alert.color" elevation="10" rounded="lg">
      <div class="d-flex align-center">
        <v-icon start>{{ alert.color === 'green' ? 'mdi-check-circle' : 'mdi-alert-circle' }}</v-icon>
        {{ alert.text }}
      </div>
    </v-snackbar>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import axios from "axios";
import type { Mock, PaginatedMocks } from '@/types';
const fileInput = ref<HTMLInputElement | null>(null);

const httpMethods = [
  {text: "GET", value: "GET", color: "blue-darken-1"},
  {text: "POST", value: "POST", color: "green-darken-1"},
  {text: "PUT", value: "PUT", color: "orange-darken-2"},
  {text: "PATCH", value: "PATCH", color: "purple-darken-1"},
  {text: "DELETE", value: "DELETE", color: "red-darken-1"},
  {text: "OPTIONS", value: "OPTIONS", color: "grey-darken-1"},
  {text: "HEAD", value: "HEAD", color: "blue-grey-darken-3"}
];

const strategies = [
  {text: "Normal", value: "normal"},
  {text: "Scene", value: "scene"},
  {text: "Random", value: "random"},
  {text: "Sequential", value: "sequential"}
];

const filters = reactive({
  group: null as string | null,
  path: null as string | null,
  strategy: null as string | null,
  method: null as string | null,
});

const table = reactive({
  columns: [
    {title: "Status", key: "status", align: 'start', sortable: false},
    {title: "Name", key: "name", width: "30%", align: 'start'},
    {title: "Group", key: "group", width: "15%", align: 'start'},
    {title: "Path", key: "path", width: "35%", align: 'start'},
    {title: "Strategy", key: "strategy", align: 'center'},
    {title: "Method", key: "method", align: 'center'},
    {title: "", key: "delete", align: 'center', sortable: false},
  ] as any[],
  rows: [] as Mock[],
  total: 0,
  loading: true,
});

const options = ref({
  page: 1,
  itemsPerPage: 10
});

// Store the last ID seen for each page to support keyset pagination
const pageCursors = ref<Record<number, string>>({});

const alert = reactive({
  show: false,
  color: "green",
  text: ""
});

function triggerFileInput() {
  fileInput.value?.click();
}

function handleOptionsUpdate(newOptions: any) {
  // If page size changed, clear cursors because they are no longer valid
  if (newOptions.itemsPerPage !== options.value.itemsPerPage) {
    pageCursors.value = {};
  }
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
  return httpMethods.find(x => x.value == httpMethod)?.color || 'grey';
}

function queryParams() {
  const { group, path, strategy, method } = filters;
  const { page, itemsPerPage } = options.value;

  // Try to use the cursor for the current page
  const lastId = pageCursors.value[page];

  let params: any = {
    limit: itemsPerPage,
    last_id: lastId || undefined,
    // Only send offset if we don't have a cursor and it's not the first page
    offset: (lastId || page === 1) ? undefined : (page - 1) * itemsPerPage,
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
  pageCursors.value = {}; // Reset cursors on new search
  options.value = {
    page: 1,
    itemsPerPage: options.value.itemsPerPage,
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
        table.rows = res.data.results || [];
        table.total = res.data.paging.total;

        // Store the cursor for the NEXT page
        const { page } = options.value;
        if (table.rows.length > 0) {
          const lastKey = table.rows[table.rows.length - 1].key;
          if (lastKey) {
            pageCursors.value[page + 1] = lastKey;
          }
        }
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
          restSearch();
        });
  }
}

async function exportMocks() {
  try {
    const res = await axios.get(baseURL() + "/export");
    const dataStr = JSON.stringify(res.data, null, 2);
    const blob = new Blob([dataStr], { type: "application/json" });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "mocks_export.json";
    link.click();
    window.URL.revokeObjectURL(url);
    showAlert("Mocks exported successfully!");
  } catch (err) {
    showAlert("Error exporting mocks!", err);
  }
}

async function importMocks(event: any) {
  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = async (e) => {
    try {
      const content = e.target?.result as string;
      const rules = JSON.parse(content);
      const res = await axios.post(baseURL() + "/import", rules);
      const { created, updated, failed } = res.data;
      showAlert(`Import result: ${created} created, ${updated} updated, ${failed} failed.`);
      restSearch();
    } catch (err) {
      showAlert("Error importing mocks! Check file format.", err);
    } finally {
      event.target.value = ""; // reset input
    }
  };
  reader.readAsText(file);
}

watch(options, () => {
  restSearch();
}, { deep: true });

</script>

<style scoped>
.list-mocks-container {
  padding: 12px;
}

.search-card {
  backdrop-filter: blur(10px);
  border-radius: 16px;
  transition: background 0.3s ease;
}

.table-card {
  border-radius: 16px;
  overflow: hidden;
}

.custom-field :deep(.v-field) {
  border-radius: 12px;
}

</style>
