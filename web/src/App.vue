<template>
  <v-app>
    <!--APP BAR-->
    <nav-bar></nav-bar>
    <!--CONTAINER-->
    <v-main>
      <v-container fluid style="max-width: 90%">
        <router-view/>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import NavBar from "./components/NavBar.vue";

export default {
  name: 'App',
  components: {
    NavBar
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
