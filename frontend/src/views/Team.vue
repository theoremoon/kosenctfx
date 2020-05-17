<template>
  <div>
    <h1 class="text-5xl">
      {{team.name}}
      <span class="text-lg">[{{team.score}}pts / {{team.pos}}{{pospostfix}}]</span>
    </h1>

    <template v-if="team.token">
      <!-- change team name will be here -->

      <div class="flex">
        <div class="text-lg flex items-center">
          <div>team token:</div>
        </div>
        <input
          type="text"
          class="inline-block flex-grow bg-transparent text-lg ml-1 font-mono"
          v-model="team.token"
          readonly
        />
        <Button class="flex-grow-0 ml-1">Regenerate</Button>
      </div>
    </template>

    <div class="text-lg ml-4">
      <ul>
        <li v-for="m in team.members" :key="m.name">
          <img :src="m.icon" class="w-5 h-5 inline-block mr-1 align-text-bottom rounded-full" />
          {{m.name}}
        </li>
      </ul>
    </div>

    <ul class="my-4 text-2xl">
      <li v-for="chal in solved_timeline" :key="chal.name" class="timeline">
        <img :src="user_icon(chal.solved_by)" class="icon mr-1 rounded-full" />
        solved {{chal.name}} and get {{chal.score}} points
        <Time class="text-lg" :time="chal.time" />.
      </li>
    </ul>
  </div>
</template>

<script>
import _ from "lodash";
import Button from "@/components/Button";
import Time from "@/components/Time";

export default {
    components: {
        Button,
        Time,
    },
  data() {
    return {
      team: {
        name: "zer0pts",
        token: "0000xxxxyyyyyfffffABCDABCD",
        id: 3,
        pos: 1,
        score: 1900,
        members: [
          {
            name: "ptr-yudai",
            icon: "https://github.com/ptr-yudai.png"
          },
          {
            name: "theoremoon",
            icon: "https://github.com/theoremoon.png"
          }
        ],
        solved_tasks: [
          {
            name: "padrsa",
            score: 1000,
            time: 1589091200,
            solved_by: "ptr-yudai"
          },
          {
            name: "sanity check",
            score: 900,
            time: 1589091230,
            solved_by: "theoremoon"
          }
        ]
      }
    };
  },
  methods: {
    user_icon(username) {
      return _.findLast(this.team.members, { name: username }).icon;
    },
  },
  computed: {
    pospostfix() {
      if (this.team.pos % 10 == 1) {
        return "st";
      } else if (this.team.pos % 10 == 2) {
        return "nd";
      } else if (this.team.pos % 10 == 3) {
        return "rd";
      } else {
        return "th";
      }
    },
    solved_timeline() {
      return _.sortBy(this.team.solved_tasks, "time");
    }
  }
};
</script>

<style lang="scss" scoped>
@import "@/vars.scss";

.timeline:not(:last-child) {
  &:after {
    display: block;
    position: relative;
    left: 0.4em;
    content: "";
    height: 1em;
    border: 1px solid $light-color;
    width: 0;
  }
}
</style>