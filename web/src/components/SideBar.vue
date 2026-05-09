<template>
  <v-navigation-drawer
    v-model="drawer"
    :rail="rail"
    permanent
    :color="color"
    theme="dark"
    :width="width"
    elevation="2"
    @click="rail = false"
    class="sidebar-drawer"
  >
    <div class="pa-2 d-flex align-center header-container" :class="{ 'flex-column': rail, 'px-4': !rail }">
      <v-img :src="mockLogo" alt="Mocks Server" :width="rail ? 40 : 50" :height="rail ? 40 : 50" class="flex-shrink-0" />
      <div v-if="!rail" class="font-weight-black text-h6 ml-2 title-text">
        Mocks Server
      </div>
      <v-spacer v-if="!rail"></v-spacer>
      <v-btn
        variant="text"
        :icon="rail ? 'mdi-chevron-right' : 'mdi-chevron-left'"
        size="small"
        @click.stop="rail = !rail"
        class="toggle-btn"
      ></v-btn>
    </div>

    <v-divider class="mx-2 mb-4" color="white"></v-divider>

    <v-list nav density="compact">
      <v-list-item
        :to="{name: 'ListMocks'}"
        prepend-icon="mdi-home-variant-outline"
        title="Home"
        exact
        rounded="lg"
        class="mb-2 nav-item"
      >
        <v-tooltip v-if="rail" activator="parent" location="right">Home</v-tooltip>
      </v-list-item>

      <v-list-item
        :to="{name: 'NewMock'}"
        prepend-icon="mdi-plus-circle-outline"
        title="New Mock"
        exact
        rounded="lg"
        class="mb-2 nav-item"
      >
        <v-tooltip v-if="rail" activator="parent" location="right">New Mock</v-tooltip>
      </v-list-item>

      <v-list-item
        :to="{name: 'Logs'}"
        prepend-icon="mdi-chart-timeline-variant"
        title="Logs"
        exact
        rounded="lg"
        class="nav-item"
      >
        <v-tooltip v-if="rail" activator="parent" location="right">Logs</v-tooltip>
      </v-list-item>
    </v-list>

    <template v-slot:append>
      <v-divider class="mx-2 mb-2" color="white"></v-divider>
      <v-list nav density="compact">
        <v-list-item
          @click="$emit('open-settings')"
          prepend-icon="mdi-cog"
          title="Settings"
          rounded="lg"
          class="mb-2 nav-item"
        >
          <v-tooltip v-if="rail" activator="parent" location="right">Settings</v-tooltip>
        </v-list-item>

        <v-list-item
          @click="toggleTheme"
          :prepend-icon="isDark ? 'mdi-weather-sunny' : 'mdi-weather-night'"
          title="Toggle Theme"
          rounded="lg"
          class="mb-2 nav-item"
        >
          <v-tooltip v-if="rail" activator="parent" location="right">Toggle Theme</v-tooltip>
        </v-list-item>

        <v-list-item
          :to="{name: 'Help'}"
          prepend-icon="mdi-help-circle-outline"
          title="Help"
          exact
          rounded="lg"
          class="nav-item"
        >
          <v-tooltip v-if="rail" activator="parent" location="right">Help</v-tooltip>
        </v-list-item>
      </v-list>
    </template>

    <!-- Resize Handle -->
    <div v-if="!rail" class="resize-handle" @mousedown="startResize"></div>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import mockLogo from '@/assets/mock.png'

const theme = useTheme()
const isDark = computed(() => theme.global.name.value === 'dark')

const drawer = ref(true)
const rail = ref(false)
const width = ref(280)
const isResizing = ref(false)

defineProps<{
  color: string
}>()

defineEmits(['open-settings'])

function toggleTheme() {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('mockserver-theme', newTheme)
}

// Resizing Logic
function startResize() {
  isResizing.value = true
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

function handleResize(e: MouseEvent) {
  if (!isResizing.value) return
  // Constrain width between 200 and 600
  const newWidth = Math.max(200, Math.min(600, e.clientX))
  width.value = newWidth
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
  localStorage.setItem('sidebar-width', width.value.toString())
}

onMounted(() => {
  const savedWidth = localStorage.getItem('sidebar-width')
  if (savedWidth) width.value = parseInt(savedWidth)
})
</script>

<style scoped>
.sidebar-drawer {
  transition: width 0.1s ease;
}

.header-container {
  min-height: 80px;
}

.title-text {
  letter-spacing: 0.5px;
  white-space: nowrap;
}

.nav-item {
  transition: all 0.2s ease;
}

.nav-item:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

.v-list-item--active {
  background-color: rgba(255, 255, 255, 0.2) !important;
  font-weight: bold;
}

.resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 6px;
  cursor: col-resize;
  background: transparent;
  transition: background 0.2s;
  z-index: 10;
}

.resize-handle:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>
