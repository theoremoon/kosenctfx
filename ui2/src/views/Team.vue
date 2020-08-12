<template>
  <div class="mt-4">
    <h1 class="text-lg">{{ teamname }}</h1>
    <div v-if="token" class="inline-form">
      teamtoken: <input type="text" readonly v-model="token" />
      <input type="submit" value="Regenerate" @click="regenerate" />
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { message, errorHandle } from "@/message";

export default Vue.extend({
  data() {
    return {
      token: "",
      teamname: ""
    };
  },
  mounted() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      API.get("team/" + this.$route.params.id).then(r => {
        if ("token" in r.data) {
          this.token = r.data.token;
        }
        this.teamname = r.data.teamname;
      });
    },
    regenerate() {
      API.post("/renew-teamtoken", {})
        .then(r => {
          message(this, r.data.message);
          this.getInfo();
        })
        .catch(e => {
          errorHandle(this, e);
        });
    }
  }
});
</script>
