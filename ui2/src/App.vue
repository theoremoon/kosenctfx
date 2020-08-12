<template>
  <div class="h-full">
    <header>
      <nav class="flex w-full p-4 navbar-border">
        <div class="flex text-lg mr-6">
          <img src="./assets/neko.png" class="brand" />
          <router-link to="/">
            InterKosenCTF2020
          </router-link>
        </div>

        <div class="flex flex-grow">
          <div class="mr-4"><router-link to="/">CHALLENGES</router-link></div>
          <div class="mr-4"><router-link to="/">RANKING</router-link></div>
        </div>

        <div class="flex flex-0">
          <template v-if="username != null">
            <div class="mr-4">
              <router-link :to="'/user/' + $store.userid">{{
                $store.username
              }}</router-link>
            </div>
            <div class="mr-4">
              <router-link :to="'/team/' + $store.teamid">{{
                $store.teamname
              }}</router-link>
            </div>
            <div class="mr-4">
              <button @click="logout">LOGOUT</button>
            </div>
          </template>
          <template v-else>
            <div class="mr-4">
              <router-link to="/login">LOGIN</router-link>
            </div>
            <div class="mr-4">
              <router-link to="/register">REGISTER</router-link>
            </div>
            <div class="mr-4">
              <button @click="logout">LOGOUT</button>
            </div>
          </template>
        </div>
      </nav>
    </header>
    <main class="container mx-auto h-full">
      <router-view />
    </main>
    <div class="messages">
      <div
        v-for="m in messages"
        :key="m.id"
        class="message"
        :class="{ error: m.type == 'error' }"
        @click="deleteMessage(m.id)"
      >
        {{ m.text }}
      </div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "./api";

export default Vue.extend({
  data() {
    return {
      messages: []
    };
  },
  mounted() {
    this.checkLogin();
    this.$eventHub.$on("login-check", () => {
      this.checkLogin();
    });

    this.$eventHub.$on("message", msg => {
      const id = this.newMessage(msg);
      setTimeout(() => {
        this.deleteMessage(id);
      }, 3000);
    });
  },
  methods: {
    logout() {
      API.post("/logout");
      this.$eventHub.$emit("login-check");
    },
    checkLogin() {
      API.get("/info")
        .then(r => {
          if ("username" in r.data) {
            Vue.set(this.$store, "username", r.data.username);
            Vue.set(this.$store, "teamname", r.data.teamname);
            Vue.set(this.$store, "userid", r.data.userid);
            Vue.set(this.$store, "teamid", r.data.teamid);
          } else {
            Vue.set(this.$store, "username", null);
            Vue.set(this.$store, "teamname", null);
            Vue.set(this.$store, "userid", null);
            Vue.set(this.$store, "teamid", null);
          }
          Vue.set(this.$store, "ctfStart", r.data.ctf_start);
          Vue.set(this.$store, "ctfEnd", r.data.ctf_end);
          Vue.set(this.$store, "ctfName", r.data.ctf_name);
        })
        .catch(() => {});
    },
    deleteMessage(id) {
      this.messages = this.messages.filter(m => m.id != id);
    },
    newMessage(msg) {
      const id = new Date().getMilliseconds();
      this.messages.push({
        text: msg.text,
        type: msg.type,
        id: id
      });
      return id;
    }
  },
  computed: {
    username() {
      return this.$store.username;
    }
  }
});
</script>

<style lang="scss">
@import "./assets/vars.scss";
@import "./assets/tailwind.css";

html,
body {
  @extend .bg-bg-color;
  color: $fg-color;
  height: 100%;
}
a {
  &:hover {
    color: $accent-color;
  }
}

input[type="text"],
input[type="email"],
input[type="password"] {
  text-align: center;
  background-color: transparent;
  border-bottom: 1px solid $accent-color;
  display: inline-block;
  width: 100%;
}
input[type="submit"] {
  margin: 0.25rem;
  padding: 0.25rem 0.5rem;
  background-color: transparent;
  border: 1px solid $accent-color;
  border-radius: 0.25rem;
}
input[type="submit"]:hover {
  background-color: $accent-color;
  cursor: pointer;
}
</style>

<style lang="scss" scoped>
@import "./assets/vars.scss";
.brand {
  height: 1.5em;
  width: auto;
}

.navbar-border {
  border-bottom: 1px solid $accent-color;
}

.messages {
  position: fixed;
  right: 20px;
  bottom: 20px;
}

.message {
  width: 15rem;
  margin-top: 0.25rem;
  padding: 0.5rem 1rem;
  border: 1px solid $accent-color;
  background-color: rgba($accent-color, 0.6);

  border-radius: 0.25rem;
}
.message.error {
  border: 1px solid $warn-color;
  background-color: rgba($warn-color, 0.6);
}
.message:hover {
  cursor: pointer;
}
</style>
