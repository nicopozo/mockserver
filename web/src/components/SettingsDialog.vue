<template>
  <v-dialog v-model="internalValue" max-width="500px">
    <v-card>
      <v-card-title class="d-flex align-center pa-4">
        <v-icon start>mdi-cog</v-icon>
        <span>Application Settings</span>
        <v-spacer></v-spacer>
        <v-btn icon variant="text" @click="internalValue = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-text class="pa-4">
        <div class="text-subtitle-1 font-weight-bold mb-4">Navigation Appearance</div>
        
        <div class="mb-6">
          <div class="text-body-2 mb-2">Theme Color</div>
          <div class="d-flex flex-wrap ga-3">
            <v-btn
              v-for="color in colors"
              :key="color.value"
              :color="color.value"
              icon
              size="large"
              :class="{ 'selected-color': modelColor === color.value }"
              @click="modelColor = color.value"
            >
              <v-icon v-if="modelColor === color.value">mdi-check</v-icon>
            </v-btn>
          </div>
        </div>

        <v-divider class="mb-4"></v-divider>

        <div>
          <div class="text-body-2 mb-2">Layout Preference</div>
          <v-btn-toggle
            v-model="modelLayout"
            mandatory
            color="primary"
            variant="outlined"
            divided
            class="w-100"
          >
            <v-btn value="sidebar" class="flex-grow-1">
              <v-icon start>mdi-page-layout-sidebar-left</v-icon>
              Sidebar
            </v-btn>
            <v-btn value="top" class="flex-grow-1">
              <v-icon start>mdi-page-layout-header</v-icon>
              Top Bar
            </v-btn>
          </v-btn-toggle>
        </div>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions class="pa-4 d-flex align-center">
        <span class="text-caption text-medium-emphasis ml-2">
          Version {{ version }}
        </span>
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="flat" @click="internalValue = false">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const version = __APP_VERSION__

const props = defineProps<{
  modelValue: boolean
  color: string
  layout: string
}>()

const emit = defineEmits(['update:modelValue', 'update:color', 'update:layout'])

const internalValue = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const modelColor = computed({
  get: () => props.color,
  set: (val) => emit('update:color', val)
})

const modelLayout = computed({
  get: () => props.layout,
  set: (val) => emit('update:layout', val)
})

const colors = [
  { name: 'Pizarra', value: '#475569' },
  { name: 'Noche', value: '#1e293b' },
  { name: 'Acero', value: '#64748b' },
  { name: 'Bosque', value: '#14532d' },
  { name: 'Vino', value: '#7f1d1d' },
  { name: 'Grafito', value: '#27272a' },
  { name: 'Petróleo', value: '#0c4a6e' }
]
</script>

<style scoped>
.selected-color {
  outline: 3px solid #fff;
  outline-offset: 2px;
  box-shadow: 0 0 10px rgba(0,0,0,0.5);
}
</style>
