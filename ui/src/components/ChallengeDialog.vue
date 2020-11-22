<template>
  <div class="dialog">
    <p class="challenge-name">{{ c.name }}</p>
    <div class="flex justify-between">
      <div class="w-3/4 break-words">
        <p>
          <span v-for="tag in c.tags" class="challenge-tag" :key="tag">{{
            tag
          }}</span>
        </p>
        <div class="flex">
          <p class="mr-4">
            <span class="challenge-value">{{ c.score }}</span
            >pts
          </p>
          <p>
            <span class="challenge-value">{{ c.solved_by.length }}</span
            >solves
          </p>
        </div>

        <div class="p-4">
          <div v-html="c.description"></div>
          <div class="text-right">author:{{ c.author }}</div>

          <div v-if="c.attachments">
            <a
              v-for="a in c.attachments"
              :key="a.name"
              :href="a.url"
              download
              class="attachment"
              @click.stop
            >
              {{ a.name }}
            </a>
          </div>
        </div>
      </div>

      <div class="w-1/4">
        <h2 class="text-xl">solved by ({{ c.solved_by.length }})</h2>
        <div v-for="t in topteams(c.solved_by)" :key="t.team_id">
          <router-link :to="'/team/' + t.team_id">{{
            t.team_name
          }}</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
export default Vue.extend({
  props: {
    c: Object
  },
  methods: {
    topteams(xs) {
      return xs
        .slice()
        .sort((a, b) => a.solved_at - b.solved_at)
        .slice(0, 10);
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../assets/vars.scss";
.challenge-name {
  font-size: 2rem;
  font-weight: bold;
}
.challenge-value {
  font-size: 1.5rem;
}
.attachment {
  display: inline-block;
  margin: 0.25rem 0;
  padding: 0.25rem 0.5rem;
  background-color: rgba($accent-color, 0.5);
  border: 1px solid $accent-color;
  border-radius: 0.25rem;
}
.dialog {
  z-index: 2;
  background: $fg-color;
  color: $bg-color;

  width: 100%;

  position: relative;
  border-radius: 0.5rem;
  padding: 1rem;
}

.challenge-tag {
  display: inline-block;
  padding: 0 0.25rem;
  background: rgba(#000000, 0.1);
  border-radius: 0.25rem;
  margin-right: 0.5rem;
}
</style>
