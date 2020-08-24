<template>
  <div class="mt-4 mx-8">
    <table class="ranking">
      <thead>
        <tr>
          <th>Rank</th>
          <th>Team</th>
          <th>Score</th>
          <th v-for="c in orderedChallenges" :key="c.id" class="challenge-name">
            <span>
              {{ c.name }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="t in orderedTeams" :key="t.id">
          <td class="text-center">{{ t.pos }}</td>
          <td>{{ t.team }}</td>
          <td class="text-right">{{ t.points }}</td>
          <td v-for="c in orderedChallenges" :key="c.name">
            <font-awesome-icon icon="flag" v-if="t.taskStats[c.name]" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import Vue from "vue";

export default Vue.extend({
  data() {
    return {
      flag: "",
      filter: "",
      focus: null
    };
  },

  computed: {
    orderedChallenges() {
      return this.$store.challenges.slice().sort((a, b) => {
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
      return this.$store.ranking;
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
    transform-origin: left bottom;
    transform: rotate(-45deg);
    width: 2em;
  }

  white-space: pre;
}
</style>
