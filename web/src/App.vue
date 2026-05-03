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

<script>
import NavBar from "./components/NavBar.vue";

export default {
  name: 'App',
  components: {
    NavBar
  }
};
</script>
