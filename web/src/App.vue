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
    const navColor = ref('primary')
    const showSettings = ref(false)

    onMounted(() => {
      const savedLayout = localStorage.getItem('mockserver-layout')
      if (savedLayout) layout.value = savedLayout

      const savedColor = localStorage.getItem('mockserver-nav-color')
      if (savedColor) {
        // Migrate old default color to new primary
        if (savedColor === 'primary-darken-1') {
          navColor.value = 'primary'
          localStorage.setItem('mockserver-nav-color', 'primary')
        } else {
          navColor.value = savedColor
        }
      }
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

