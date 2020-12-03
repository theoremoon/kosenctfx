<template>
  <tr>
    <td>{{ c.name }}</td>
    <td>{{ c.score }}</td>
    <td>{{ c.solved_by.length }}</td>
    <td>
      <span v-for="tag in c.tags" :key="tag" class="challenge-tag"
        >{{ tag }}
      </span>
    </td>

    <td>{{ c.author }}</td>
    <td>
      <input type="checkbox" v-model="c.is_open" />
    </td>
    <td><span v-if="c.is_survey">🙆</span><span v-else>🙅</span></td>
    <td @click="onFocus" class="cursor-pointer">👀</td>
    <td>
      <pre>{{ c.flag }}</pre>
    </td>
  </tr>
</template>

<script>
import Vue from "vue";
export default Vue.extend({
  props: {
    c: {
      name: String,
      flag: String,
      author: String,
      is_open: Boolean,
      is_survey: Boolean
    }
  },
  methods: {
    onFocus() {
      this.$emit("focus", this.c);
    },
    onCheck(ev) {
      this.$emit("check", this.c, ev.target.checked);
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../../../assets/vars.scss";
@import "../../../assets/tailwind.css";
td {
  border: solid 1px $fg-color;
}

.challenge-tag {
  display: inline-block;
  padding: 0 0.25rem;
  background: rgba(#ffffff, 0.2);
  border-radius: 0.25rem;
  margin-right: 0.5rem;
}
</style>
