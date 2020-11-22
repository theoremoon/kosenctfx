<template>
  <div class="my-4 mx-8">
    <h1 class="text-2xl">Administration</h1>

    <div>
      <router-link to="/admin/config" class="inline-block mr-2"
        >Config</router-link
      >
      <router-link to="/admin/challenges" class="inline-block mr-2"
        >Challenges</router-link
      >
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { errorHandle } from "@/message";

export default Vue.extend({
  mounted() {
    this.getConfig();
  },
  methods: {
    getConfig() {
      API.get("admin/get-config").catch(() => {
        errorHandle(this, "forbidden");
        this.$router.push("/");
      });
    }
  }
});
</script>
