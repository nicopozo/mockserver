<template>
  <div>
    <v-container fuild style="max-width: 70%">
      <v-card>
        <v-card-text>
          <v-responsive :aspect-ratio="16/9">
            <iframe :src="videoSrc" width="100%" height="100%" style="border: none;"></iframe>
          </v-responsive>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
export default {
  name: "Help",
  data() {
    return {
      videoSrc: "https://docs.google.com/file/d/1aKCu18gcgDt7gpwgWnGmWfKWJCbOKkTJ/preview"
    }
  },
}
</script>
