<template>
  <div class="mt-4">
    <pre>{{ token }}</pre>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";

export default Vue.extend({
  data() {
    return {
      token: ""
    };
  },
  mounted() {
    API.get("team/" + this.$route.params.id).then(r => {
      this.token = r.data.token;
    });
  },
  methods: {}
});
</script>
