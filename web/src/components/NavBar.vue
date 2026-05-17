<template>
  <v-app-bar 
    app 
    :color="color"
    theme="dark"
    elevation="2" 
    height="80"
    class="main-navbar px-4"
  >
    <template #prepend>
      <img :src="mockLogo" alt="Mocks Server" class="ml-2 flex-shrink-0" style="width: 40px; height: 40px; object-fit: contain;" />
    </template>

    <div class="font-weight-black ml-2" style="letter-spacing: 0.5px; font-size: 1.4rem; white-space: nowrap;">
      Mocks Server
    </div>

    <v-divider vertical class="mx-6 my-4" color="white"></v-divider>

    <div class="nav-links d-flex ga-2">
      <v-btn :to="{name: 'ListMocks'}" variant="text" class="nav-btn" exact rounded="lg">
        <v-icon start size="20">mdi-home-variant-outline</v-icon>
        <span class="font-weight-medium">Home</span>
      </v-btn>

      <v-btn :to="{name: 'NewMock'}" variant="text" class="nav-btn" exact rounded="lg">
        <v-icon start size="20">mdi-plus-circle-outline</v-icon>
        <span class="font-weight-medium">New Mock</span>
      </v-btn>

      <v-btn :to="{name: 'Logs'}" variant="text" class="nav-btn" exact rounded="lg">
        <v-icon start size="20">mdi-chart-timeline-variant</v-icon>
        <span class="font-weight-medium">Logs</span>
      </v-btn>
    </div>

    <v-spacer></v-spacer>

    <div class="d-flex align-center ga-2">
      <v-btn icon variant="text" size="small" @click="$emit('open-settings')" title="Settings">
        <v-icon size="20">mdi-cog</v-icon>
      </v-btn>

      <v-btn icon variant="text" size="small" @click="toggleTheme" title="Toggle Theme">
        <v-icon size="20">{{ isDark ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
      </v-btn>

      <v-btn :to="{name: 'Help'}" icon variant="text" size="small" exact title="Help">
        <v-icon size="20">mdi-help-circle-outline</v-icon>
      </v-btn>
    </div>
  </v-app-bar>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from 'vuetify'
import mockLogo from '@/assets/mock.png'

defineProps<{
  color: string
}>()

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'dark')

defineEmits(['open-settings'])

function toggleTheme() {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('mockserver-theme', newTheme)
}
</script>

<style scoped>
.nav-btn {
  text-transform: none;
  letter-spacing: 0;
  transition: all 0.2s ease;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.v-btn--active.nav-btn {
  background-color: rgba(255, 255, 255, 0.2);
}
</style>
