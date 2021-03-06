<template>
  <div class="mt-4">
    <h1 class="text-lg">
      {{ teamname }} {{ country | countryFlag }} -
      <span class="text-2xl">{{ score }}</span
      >pts
    </h1>

    <div v-if="$store.teamname == teamname">
      <p v-if="hasEnd()">
        <router-link :to="'/team/' + $store.teamid + '/diploma'">
          diploma
        </router-link>
      </p>
      <button @click="editmode = !editmode">edit profile</button>
      <div class="w-1/4 ml-4" v-if="editmode">
        <div class="mb-4">
          <label class="block text-sm" for="teamname"> teamname </label>
          <input type="text" v-model="edit_teamname" id="teamname" />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="password"> password </label>
          <input type="password" v-model="edit_password" id="password" />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="country">
            Country Code {{ edit_country | countryFlag }}
          </label>
          <input type="text" v-model="edit_country" id="country" />
        </div>

        <input
          type="submit"
          value="update"
          class="float-right"
          @click="updateProfile"
        />
      </div>
    </div>

    <div class="lg:w-3/4 m-auto">
      <graph :chartdata="chartData" />
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
import Graph from "../components/Graph";
import colorhash from "../colorhash";

export default Vue.extend({
  components: {
    graph: Graph,
  },
  data() {
    return {
      token: "",
      teamname: "",
      country: "",

      editmode: false,
      edit_teamname: "",
      edit_password: "",
      edit_country: "",
      now: new Date().valueOf() / 1000,
    };
  },
  mounted() {
    this.getInfo();
  },
  methods: {
    getInfo() {
      API.get("team/" + this.$route.params.id).then((r) => {
        this.teamname = r.data.teamname;
        this.country = r.data.country;
      });
    },
    updateProfile() {
      API.post("/update-profile", {
        teamname: this.edit_teamname,
        country: this.edit_country,
      })
        .then((r) => {
          message(this, r.data.message);
          this.$eventHub.$emit("login-check");
          this.getInfo();

          this.edit_teamname = "";
          this.edit_password = "";
          this.edit_country = "";
          this.editmode = false;
        })
        .catch((e) => {
          errorHandle(this, e);
        });
    },
    hasEnd() {
      return this.now > this.$store.ctfEnd;
    },
  },
  filters: {
    dateFormat(t) {
      return dateFormat(t);
    },
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
        .map((v) => v["points"])
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
          ...v,
        };
      });
      challenges.sort((a, b) => a.time - b.time);
      return challenges;
    },
    chartData() {
      let current_score = 0;
      let data = [];

      this.solvedChallenges.forEach((c) => {
        current_score += c.points;
        data.push({
          t: new Date(c.time * 1000),
          y: current_score,
          name: c.name,
          score: c.points,
          team: "",
        });
      });
      return [
        {
          lineTension: 0,
          borderColor: colorhash(this.teamname),
          backgroundColor: colorhash(this.teamname),
          fill: false,
          data: data,
          label: this.teamname,
        },
      ];
    },
  },
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
