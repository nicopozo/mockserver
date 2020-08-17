<template>
  <div id="app">
    <nav-bar></nav-bar>
    <span>
      <router-view></router-view>
    </span>
  </div>
</template>

<script>
import NavBar from "./components/NavBar.vue";

export default {
  name: "App",
  components: {
    NavBar
  }
};
</script>