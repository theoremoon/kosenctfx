<template>
  <div class="mt-4">
    <h1 class="text-lg">{{ teamname }}</h1>
    <div class="px-4 py-2">
      <div class="member" v-for="m in members" :key="m.id">
        <router-link :to="'/user/' + m.id"> {{ m.username }}</router-link>
      </div>
    </div>

    <div v-if="token" class="inline-form px-4">
      <label>teamtoken:</label> <input type="text" readonly v-model="token" />
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
      teamname: "",
      members: []
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
        this.members = r.data.members;
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

<style lang="scss" scoped>
.member {
  display: inline-block;
  &:not(:last-child)::after {
    padding-left: 1rem;
    padding-right: 1rem;
    content: "/";
  }
}
</style>
