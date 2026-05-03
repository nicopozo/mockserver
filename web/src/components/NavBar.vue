<template>
  <div>
    <v-app-bar app color="primary" theme="dark" class="px-4">
      <v-img
          alt="Mocks Server"
          class="shrink mr-2"
          :src="mockLogo"
          max-height="40"
          max-width="40"
      />

      <v-app-bar-title class="mr-3">Mocks Server</v-app-bar-title>

      <v-btn :to="{name: 'ListMocks'}" variant="text" class="mx-1" exact>
        <v-icon>mdi-home</v-icon>
        <span class="mx-1">Home</span>
      </v-btn>

      <v-btn :to="{name: 'NewMock'}" variant="text" class="mx-1" exact>
        <v-icon>mdi-plus-circle-outline</v-icon>
        <span class="mx-1">New Mock</span>
      </v-btn>

      <v-spacer></v-spacer>

      <v-btn :to="{name: 'Help'}" variant="text" class="mx-1" exact>
        <v-icon>mdi-help-circle</v-icon>
      </v-btn>

    </v-app-bar>
  </div>
</template>

<script>
import mockLogo from '@/assets/mock.png'

export default {
  name: "NavBar",
  data() {
    return {
      mockLogo
    }
  }
};
</script>