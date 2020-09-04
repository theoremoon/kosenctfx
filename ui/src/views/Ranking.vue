<template>
  <div class="my-4 mx-8">
    <h1 class="text-2xl">Team Ranking</h1>

    <graph :chartdata="chartData" v-if="dataReady"></graph>

    <table class="ranking">
      <thead>
        <tr>
          <th>Rank</th>
          <th>Team</th>
          <th>Score</th>
          <th v-for="c in orderedChallenges" :key="c" class="challenge-name">
            <span>
              {{ c }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in orderedTeams" :key="t.team_id">
          <td class="text-center">{{ t.pos }}</td>
          <td>
            <router-link :to="'/team/' + t.team_id">{{ t.team }}</router-link>
          </td>
          <td class="text-right">{{ t.points }}</td>
          <td v-for="c in orderedChallenges" :key="c">
            <font-awesome-icon icon="flag" v-if="t.taskStats[c]" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import Vue from "vue";
import Graph from "../components/Graph";
import colorhash from "../colorhash";

export default Vue.extend({
  components: {
    graph: Graph
  },
  computed: {
    orderedChallenges() {
      if (!this.$store.ranking || !this.$store.ranking.tasks) {
        return [];
      }

      return this.$store.ranking.tasks.slice().sort((a, b) => {
        if (a.name < b.name) {
          return -1;
        } else if (a.name > b.name) {
          return 1;
        } else {
          return 0;
        }
      });
    },
    orderedTeams() {
      if (!this.$store.ranking) {
        return [];
      }
      return this.$store.ranking.standings;
    },
    dataReady() {
      if (!this.$store.ranking) {
        return false;
      }
      if (!this.$store.ranking.standings) {
        return false;
      }
      return true;
    },
    chartData() {
      if (!this.$store.ranking.standings) {
        return [];
      }
      return this.$store.ranking.standings.slice(0, 10).map(t => {
        let current_score = 0;
        let data = [];
        let tasks = Object.entries(t.taskStats).map(([name, info]) => {
          return { name, info };
        });
        tasks.sort((a, b) => a.info.time - b.info.time);
        tasks.forEach(e => {
          let name = e.name;
          let info = e.info;
          current_score += info.points;
          data.push({
            t: new Date(info.time * 1000),
            y: current_score,
            name: name,
            score: info.points,
            team: t.team
          });
        });
        return {
          label: t.team,
          lineTension: 0,
          borderColor: colorhash(t.team),
          backgroundColor: colorhash(t.team),
          fill: false,
          data: data
        };
      });
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../assets/vars.scss";

.ranking {
  padding-top: 10rem;

  display: block;
  overflow-x: auto;

  tr {
    border-bottom: 1px solid $fg-color;
  }
}
.challenge-name {
  span {
    display: inline-block;
    transform-origin: left;
    transform: rotate(-45deg);
    width: 2em;
  }

  white-space: pre;
}
</style>
