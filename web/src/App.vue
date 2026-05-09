<template>
  <v-app>
    <!--SIDEBAR-->
    <side-bar 
      v-if="layout === 'sidebar'" 
      :color="navColor"
      @open-settings="showSettings = true"
    ></side-bar>
    <!--NAVBAR-->
    <nav-bar 
      v-else 
      :color="navColor"
      @open-settings="showSettings = true"
    ></nav-bar>

    <!--CONTAINER-->
    <v-main>
      <v-container fluid :style="{ 'max-width': layout === 'sidebar' ? '95%' : '90%' }">
        <router-view/>
      </v-container>
    </v-main>

    <!--SETTINGS-->
    <settings-dialog
      v-model="showSettings"
      v-model:color="navColor"
      v-model:layout="layout"
      @update:color="saveColor"
      @update:layout="saveLayout"
    ></settings-dialog>
  </v-app>
</template>

<script lang="ts">
import { ref, onMounted } from 'vue'
import SideBar from "./components/SideBar.vue";
import NavBar from "./components/NavBar.vue";
import SettingsDialog from "./components/SettingsDialog.vue";

export default {
  name: 'App',
  components: {
    SideBar,
    NavBar,
    SettingsDialog
  },
  setup() {
    const layout = ref('sidebar')
    const navColor = ref('primary-darken-1')
    const showSettings = ref(false)

    onMounted(() => {
      const savedLayout = localStorage.getItem('mockserver-layout')
      if (savedLayout) layout.value = savedLayout

      const savedColor = localStorage.getItem('mockserver-nav-color')
      if (savedColor) navColor.value = savedColor
    })

    function toggleLayout() {
      layout.value = layout.value === 'sidebar' ? 'top' : 'sidebar'
      saveLayout(layout.value)
    }

    function saveLayout(val: string) {
      localStorage.setItem('mockserver-layout', val)
    }

    function saveColor(val: string) {
      localStorage.setItem('mockserver-nav-color', val)
    }

    return {
      layout,
      navColor,
      showSettings,
      toggleLayout,
      saveLayout,
      saveColor
    }
  }
};
</script>
<style>
/* Estilos Globales para Tablas (Vuetify 3) */
.v-data-table .v-data-table__th,
.v-data-table thead th {
  background-color: #1565C0 !important;
  color: white !important;
  text-transform: uppercase;
  font-weight: 800 !important;
  font-size: 0.75rem !important;
  letter-spacing: 0.5px;
}

/* Iconos de ordenamiento blancos */
.v-data-table .v-data-table__th .v-icon,
.v-data-table thead th .v-icon {
  color: white !important;
}

/* Ajuste para que el hover no manche el azul del header */
.v-data-table .v-data-table__th:hover,
.v-data-table thead th:hover {
  background-color: #1565C0 !important;
}
</style>
