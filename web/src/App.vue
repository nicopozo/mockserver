<template>
  <v-app>
    <!--APP BAR-->
    <nav-bar></nav-bar>
    <!--CONTAINER-->
    <v-main style="background-color: #f1f1f1;">
      <v-container fuild style="max-width: 90%">
        <router-view/>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import NavBar from "./components/NavBar.vue";

export default {
  name: 'App',
  components: {
    NavBar
  }
};
</script>
