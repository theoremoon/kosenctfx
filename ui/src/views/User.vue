<template>
  <div class="mt-4">
    <h1 class="text-lg">
      {{ username }} @
      <router-link :to="'/team/' + teamid">{{ teamname }}</router-link>
      - <span class="text-2xl">{{ score }}</span> pts
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

    <div class="mt-4 ml-4 text-xl">
      <table class="table-auto w-full xl:w-3/4 mx-auto">
        <thead>
          <tr class="bottomline">
            <th>Challenge</th>
            <th>Score</th>
            <th>Solved at</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in solvedChallenges" :key="c.name">
            <td class="text-center">{{ c.name }}</td>
            <td class="text-right">{{ c.points }}</td>
            <td class="text-center">{{ c.time | dateFormat }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { message, errorHandle } from "@/message";
import { dateFormat } from "../dateformat";

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
  },
  filters: {
    dateFormat(t) {
      return dateFormat(t);
    }
  },
  computed: {
    score() {
      if (!this.$store.userRanking || !this.$store.userRanking.standings) {
        return 0;
      }
      let score = null;
      for (let i = 0; i < this.$store.userRanking.standings.length; i++) {
        if (this.$store.userRanking.standings[i].team == this.username) {
          score = this.$store.userRanking.standings[i];
          break;
        }
      }
      if (!score) {
        return 0;
      }
      return Object.values(score["taskStats"])
        .map(v => v["points"])
        .reduce((a, b) => a + b, 0);
    },
    solvedChallenges() {
      if (!this.$store.userRanking || !this.$store.userRanking.standings) {
        return [];
      }

      let score = null;
      for (let i = 0; i < this.$store.userRanking.standings.length; i++) {
        if (this.$store.userRanking.standings[i].team == this.username) {
          score = this.$store.userRanking.standings[i];
          break;
        }
      }
      if (!score) {
        return [];
      }
      let challenges = Object.entries(score["taskStats"]).map(([k, v]) => {
        return {
          name: k,
          ...v
        };
      });
      challenges.sort((a, b) => a.time - b.time);
      console.log(challenges);
      return challenges;
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../assets/vars.scss";

.bottomline {
  border-bottom: 1px solid $fg-color;
}
</style>
