<template>
  <div class="flex justify-center items-center h-full">
    <form
      class="xl:w-1/2 w-full flex items-center flex-col"
      @submit.prevent="reset"
    >
      <div class="w-1/2">
        <div class="mb-4">
          <label class="block text-sm" for="email">
            email
          </label>
          <input type="email" v-model="email" id="email" />
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
      email: "",
      password: ""
    };
  },
  mounted() {
    if (this.$store.username != null) {
      this.$router.push("/");
    }
  },
  methods: {
    reset() {
      API.post("/passwordreset-request", {
        email: this.email,
        password: this.password
      })
        .then(r => {
          message(this, r.data.message);
          this.$router.push("/reset");
        })
        .catch(e => {
          console.log("A", e.response.data.message);
          errorHandle(this, e);
        });
    }
  }
});
</script>
