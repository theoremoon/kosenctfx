<template>
  <div class="my-4 mx-8">
    <h2 class="text-2xl">Current Challenges</h2>
    <div>
      <input
        type="button"
        value="Open / Close Challenges"
        @click="openCloseChallenges"
      />
    </div>
    <div class="mx-4">
      <table class="w-full">
        <thead>
          <tr>
            <th>Name</th>
            <th>Score</th>
            <th>#Solve</th>
            <th>Tags</th>
            <th>Flag</th>
            <th>Author</th>
            <th>Is Open?</th>
            <th>Is Survey?</th>
            <th>Preview</th>
          </tr>
        </thead>
        <tbody>
          <challenge-raw
            v-for="c in challenges"
            :c="c"
            @focus="focusToChallenge(c)"
            :key="c.id"
          ></challenge-raw>
        </tbody>
      </table>
    </div>

    <template v-if="focus">
      <h2 class="text-2xl mt-4">Challenge Preview</h2>
      <div class="flex">
        <div class="mx-8 py-4 flex dialog-background" @click="focus = null">
          <div class="w-3/4 mx-auto">
            <challenge-dialog :c="focus"></challenge-dialog>
          </div>
        </div>
      </div>
    </template>

    <h2 class="text-2xl mt-4">Add New Challenge</h2>
    <challenge-register @update="loadChallenges()" />
  </div>
</template>

<script>
import Vue from "vue";

import API from "@/api";
import { errorHandle, message } from "../../message";

import ChallengeRaw from "./components/ChallengeRaw.vue";
import ChallengeDialog from "../../components/ChallengeDialog";
import ChallengeRegister from "./components/ChallengeRegister";

export default Vue.extend({
  components: {
    ChallengeRaw,
    ChallengeDialog,
    ChallengeRegister
  },
  data() {
    return {
      challenges: [],
      checks: {},
      focus: null
    };
  },
  mounted() {
    this.loadChallenges();
  },
  methods: {
    loadChallenges() {
      Vue.set(this, "challenges", []);
      Vue.set(this, "checks", {});

      API.get("admin/list-challenges")
        .then(r => {
          Vue.set(this, "challenges", r.data.challenges);
        })
        .catch(e => {
          errorHandle(this, e);
          this.$router.push("/");
        });
    },

    focusToChallenge(c) {
      this.focus = Vue.util.extend({}, c);
    },

    async openCloseChallenges() {
      await this.openCloseChallenges_impl();
      this.loadChallenges();
    },

    async openCloseChallenges_impl() {
      const promises = [];
      this.challenges.forEach(c => {
        const endpoint = c.is_open
          ? "admin/open-challenge"
          : "admin/close-challenge";
        promises.push(
          API.post(endpoint, {
            name: c.name
          }).then(r => {
            if (!r.data.includes("already")) {
              message(this, r.data);
            }
          })
        );
      });
      return Promise.all(promises);
    }
  }
});
</script>

<style lang="scss">
@import "../../assets/vars.scss";
@import "../../assets/tailwind.css";
td,
th {
  border: solid 1px $fg-color;
}

.dialog-background {
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1;
  &:hover {
    cursor: pointer;
  }
}

#description {
  width: 100%;
  height: 5rem;
}
</style>
