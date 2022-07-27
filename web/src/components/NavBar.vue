<template>
  <div>
    <v-app-bar app color="primary" dark>
      <v-img
          alt="Mocks Server"
          class="shrink mr-2"
          contain
          src="../assets/mock.png"
          transition="scale-transition"
          width="40"
      />

      <v-toolbar-title class="mr-3">Mocks Server</v-toolbar-title>

      <v-btn :to="{name: 'ListMocks'}" text class="mx-1" exact>
        <v-icon>mdi-home</v-icon>
        <span class="mx-1">Home</span>
      </v-btn>

      <v-btn :to="{name: 'NewMock'}" text class="mx-1" exact>
        <v-icon>mdi-plus-circle-outline</v-icon>
        <span class="mx-1">New Mock</span>
      </v-btn>

      <v-spacer></v-spacer>

      <v-btn :to="{name: 'Help'}" text class="mx-1" exact>
        <v-icon>mdi-help-circle</v-icon>
      </v-btn>

    </v-app-bar>
  </div>
</template>

<script>
export default {
  name: "NavBar"
};
</script>