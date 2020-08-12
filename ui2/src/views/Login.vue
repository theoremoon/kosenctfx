<template>
  <div class="flex justify-center items-center h-full">
    <form class="w-1/2 flex items-center flex-col" @submit.prevent="login">
      <div class="w-1/2">
        <div class="mb-4">
          <label class="block text-sm" for="username">
            username
          </label>
          <input type="text" v-model="username" id="username" required />
        </div>
        <div class="mb-4">
          <label class="block text-sm" for="password">
            password
          </label>
          <input type="password" v-model="password" id="password" required />
        </div>

        <input type="submit" value="Login" class="float-right" />
      </div>
      <div class="w-1/2">
        <p style="text-align: right;">
          forget your password?
          <router-link to="/reset-request">reset from here</router-link>
        </p>
      </div>
    </form>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { errorHandle, message } from "@/message";
export default Vue.extend({
  data() {
    return {
      username: "",
      password: ""
    };
  },
  mounted() {
    if (this.$store.username != null) {
      this.$router.push("/");
    }
  },
  methods: {
    login() {
      API.post("/login", {
        username: this.username,
        password: this.password
      })
        .then(r => {
          message(this, r.data.message);
          this.$eventHub.$emit("login-check");
          this.$router.push("/");
        })
        .catch(e => {
          errorHandle(this, e);
        });
    }
  }
});
</script>
