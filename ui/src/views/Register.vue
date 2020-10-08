<template>
  <div class="flex justify-center items-center h-full">
    <form class="xl:w-1/2 w-full flex items-center flex-col" @submit.prevent>
      <div class="w-1/2">
        <div class="mb-4">
          <label class="block text-sm" for="email">
            email
          </label>
          <input type="email" v-model="email" id="email" required />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="teamname">
            teamname
          </label>
          <input type="text" v-model="teamname" id="teamname" required />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="password">
            password
          </label>
          <input type="password" v-model="password" id="password" required />
        </div>

        <div class="mv-4">
          <input type="submit" value="Register" @click="register" />
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { message, errorHandle } from "@/message";

export default Vue.extend({
  data() {
    return {
      email: "",
      password: "",
      teamname: ""
    };
  },
  mounted() {
    if (this.$store.teamname != null) {
      this.$router.push("/");
    }
  },
  methods: {
    register() {
      API.post("/register", {
        teamname: this.teamname,
        email: this.email,
        password: this.password
      })
        .then(r => {
          message(this, r.data.message);
          this.$router.push("/login");
        })
        .catch(e => {
          errorHandle(this, e);
        });
    }
  }
});
</script>
