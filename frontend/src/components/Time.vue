<template>
  <span>{{formatTime}}</span>
</template>

<script>
import dayjs from "dayjs";

export default {
  props: ["time"],
  computed: {
    formatTime() {
      return dayjs(this.time * 1000).format("YYYY-MM-DD HH:mm:ss Z");
    }
  }
};
</script>