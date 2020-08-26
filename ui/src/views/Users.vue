<template>
  <div class="my-4 mx-8">
    <h1 class="text-2xl">User Ranking</h1>
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
        <tr v-for="u in orderdUsers" :key="u.team_id">
          <td class="text-center">{{ u.pos }}</td>
          <td>
            <router-link :to="'/user/' + u.team_id">{{ u.team }}</router-link>
          </td>
          <td class="text-right">{{ u.points }}</td>
          <td v-for="c in orderedChallenges" :key="c">
            <font-awesome-icon icon="flag" v-if="u.taskStats[c]" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import Vue from "vue";

export default Vue.extend({
  computed: {
    orderedChallenges() {
      if (!this.$store.userRanking || !this.$store.userRanking.tasks) {
        return [];
      }

      return this.$store.userRanking.tasks.slice().sort((a, b) => {
        if (a.name < b.name) {
          return -1;
        } else if (a.name > b.name) {
          return 1;
        } else {
          return 0;
        }
      });
    },
    orderdUsers() {
      if (!this.$store.userRanking) {
        return [];
      }
      return this.$store.userRanking.standings;
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
