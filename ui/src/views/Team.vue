<template>
  <div class="mt-4">
    <h1 class="text-lg">
      {{ teamname }} - <span class="text-2xl">{{ score }}</span
      >pts
    </h1>
    <div class="px-4 py-2">
      <div class="member" v-for="m in members" :key="m.id">
        <router-link :to="'/user/' + m.id"> {{ m.username }}</router-link>
      </div>
    </div>

    <div v-if="token" class="inline-form px-4">
      <label>teamtoken:</label> <input type="text" readonly v-model="token" />
      <input type="submit" value="Regenerate" @click="regenerate" />
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
  },
  filters: {
    dateFormat(t) {
      return dateFormat(t);
    }
  },
  computed: {
    score() {
      if (!this.$store.ranking || !this.$store.ranking.standings) {
        return 0;
      }
      let score = null;
      for (let i = 0; i < this.$store.ranking.standings.length; i++) {
        if (this.$store.ranking.standings[i].team == this.teamname) {
          score = this.$store.ranking.standings[i];
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
      if (!this.$store.ranking || !this.$store.ranking.standings) {
        return [];
      }

      let score = null;
      for (let i = 0; i < this.$store.ranking.standings.length; i++) {
        if (this.$store.ranking.standings[i].team == this.teamname) {
          score = this.$store.ranking.standings[i];
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
.member {
  display: inline-block;
  &:not(:last-child)::after {
    padding-left: 1rem;
    padding-right: 1rem;
    content: "/";
  }
}
</style>
