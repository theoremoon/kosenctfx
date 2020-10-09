<template>
  <div class="flex justify-center items-center h-full">
    <form
      class="xl:w-1/2 w-full flex items-center flex-col"
      @submit.prevent="reset"
    >
      <div class="w-1/2">
        <div class="mb-4">
          <label class="block text-sm" for="token">
            token
          </label>
          <input type="text" v-model="token" id="token" required />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="new_password">
            newpassword
          </label>
          <input
            type="password"
            v-model="new_password"
            id="new_password"
            required
          />
        </div>

        <input type="submit" value="Reset Password" class="float-right" />
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
      token: "",
      new_password: ""
    };
  },
  mounted() {
    if (this.$store.teamname != null) {
      this.$router.push("/");
    }
  },
  methods: {
    reset() {
      API.post("/passwordreset", {
        token: this.token,
        new_password: this.new_password
      })
        .then(r => {
          message(this, r.data.message);
          this.$router.push("/login");
        })
        .catch(e => {
          console.log("A", e.response.data.message);
          errorHandle(this, e);
        });
    }
  }
});
</script>
