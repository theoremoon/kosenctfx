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
      </div>

      <div class="flex w-full">
        <div class="w-1/2 p-4">
          <label for="teamname" class="block text-sm">
            teamname
          </label>
          <input type="text" v-model="teamname" id="teamname" />
          <input
            type="submit"
            value="Create New Team"
            class="float-right"
            @click="createTeam"
          />
        </div>
        <div class="w-1/2 p-4">
          <label for="token" class="block text-sm">team token</label>
          <input type="text" v-model="token" id="token" />
          <input
            type="submit"
            value="Join to the Team"
            class="float-right"
            @click="joinTeam"
          />
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
      username: "",
      email: "",
      password: "",
      teamname: "",
      token: ""
    };
  },
  mounted() {
    if (this.$store.username != null) {
      this.$router.push("/");
    }
  },
  methods: {
    createTeam() {
      API.post("/register-with-team", {
        teamname: this.teamname,
        username: this.username,
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
    },
    joinTeam() {
      API.post("/register-and-join-team", {
        token: this.token,
        username: this.username,
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
