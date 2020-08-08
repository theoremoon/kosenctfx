<template>
  <div>
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
    if (this.$store.username != null) {
      this.$route.push("/");
    }

    API.get("team/" + this.$route.params.id).then(r => {
      this.token = r.data.token;
    });
  },
  methods: {}
});
</script>
