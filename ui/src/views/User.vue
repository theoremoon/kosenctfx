<template>
  <div class="mt-4">
    <h1 class="text-lg">
      {{ username }} @
      <router-link :to="'/team/' + teamid">{{ teamname }}</router-link>
    </h1>

    <div v-if="$store.userid == userid" class="inline-form px-4">
      <input
        type="password"
        v-model="new_password"
        id="new_password"
        placeholder="your new password"
        required
      />

      <input type="submit" value="Update Password" @click="passwordUpdate" />
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
      username: "",
      teamname: "",
      userid: 0,
      teamid: 0,
      new_password: ""
    };
  },
  mounted() {
    API.get("user/" + this.$route.params.id).then(r => {
      this.username = r.data.username;
      this.teamname = r.data.teamname;
      this.userid = r.data.userid;
      this.teamid = r.data.teamid;
    });
  },
  methods: {
    passwordUpdate() {
      API.post("/password-update", {
        new_password: this.new_password
      })
        .then(r => {
          message(this, r.data.message);
        })
        .catch(e => {
          errorHandle(this, e);
        });
    }
  }
});
</script>
