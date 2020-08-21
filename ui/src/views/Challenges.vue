<template>
  <div class="mt-4 mx-8">
    <div class="flex justify-between">
      <form class="inline-form" @submit.prevent="submit">
        <label>flag:</label>
        <input type="text" placeholder="KosenCTF{.+}" v-model="flag" />
        <input type="submit" value="Submit" />
      </form>

      <div class="inline-form">
        <input type="text" placeholder="filter" v-model="filter" />
      </div>
    </div>

    <div v-if="focus">
      <div class="focus-background" @click="focus = ''"></div>
      <div class="focus-challenge">
        <p class="challenge-name">{{ focusedchallenge.name }}</p>
        <div class="flex justify-between">
          <div>
            <p>
              <span
                v-for="tag in focusedchallenge.tags"
                class="challenge-tag"
                :key="tag"
                >{{ tag }}</span
              >
            </p>
            <div class="flex">
              <p class="mr-4">
                <span class="challenge-value">{{ focusedchallenge.score }}</span
                >pts
              </p>
              <p>
                <span class="challenge-value">{{
                  focusedchallenge.solved_by.length
                }}</span
                >solves
              </p>
            </div>

            <div class="p-4">
              <div v-html="focusedchallenge.description"></div>
              <div class="text-right">author:{{ focusedchallenge.author }}</div>

              <div v-if="focusedchallenge.attachments">
                <a
                  v-for="a in focusedchallenge.attachments"
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

          <div class="flex-grow-0">
            <h2 class="text-xl">solved by</h2>
            <div v-for="t in focusedchallenge.solved_by" :key="t">
              {{ t }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div
      class="challenges mt-4 grid gap-4 sm:grid-cols-1 md:grid-cols-3 xl:grid-cols-5"
      v-else
    >
      <div
        class="challenge"
        v-for="c in list_challenges"
        :class="{ 'challenge-solved': c.solved_by.includes($store.teamname) }"
        :key="c.name"
        @click="focus = c.name"
      >
        <p class="challenge-name">{{ c.name }}</p>
        <div class="flex justify-around">
          <p>
            <span class="challenge-value">{{ c.score }}</span
            >pts
          </p>
          <p>
            <span class="challenge-value">{{ c.solved_by.length }}</span
            >solves
          </p>
        </div>
        <div>
          <span class="challenge-tag" v-for="tag in c.tags" :key="tag">{{
            tag
          }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { message, errorHandle } from "@/message";

export default Vue.extend({
  data() {
    return {
      flag: "",
      filter: "",
      focus: null
    };
  },
  methods: {
    submit() {
      API.post("/submit", {
        flag: this.flag
      })
        .then(r => {
          message(this, r.data.message);
        })
        .catch(e => {
          errorHandle(this, e);
        });
      this.flag = "";
    }
  },
  computed: {
    focusedchallenge() {
      for (const c of this.$store.challenges) {
        if (c.name == this.focus) {
          return c;
        }
      }
      return null;
    },
    list_challenges() {
      return this.$store.challenges
        .slice()
        .sort((a, b) => {
          if (a.score == b.score) {
            if (a.name < b.name) {
              return -1;
            } else if (a.name > b.name) {
              return 1;
            } else {
              return 0;
            }
          }
          return a.score < b.score ? -1 : 1;
        })
        .filter(e => {
          if (this.filter === "") {
            return true;
          }
          if (e.name.includes(this.filter)) {
            return true;
          }
          for (const t of e.tags) {
            if (t.includes(this.filter)) {
              return true;
            }
          }
          return false;
        });
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../assets/vars.scss";

.focus-background {
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1;
  &:hover {
    cursor: pointer;
  }
}
.focus-challenge {
  z-index: 2;
  background: $fg-color;
  color: $bg-color;

  width: 100%;

  position: relative;
  border-radius: 0.5rem;
  padding: 1rem;

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
}

.challenges {
  .challenge {
    background: $fg-color;
    &:hover {
      cursor: pointer;
      background: rgba($fg-color, 0.7);
    }
    & + .challenge-solved {
      border: $accent-color 4px solid;
    }
    border-radius: 0.25rem;

    padding: 0.25rem 0.5rem;
    color: $bg-color;

    .challenge-name {
      font-weight: bold;
      font-size: 1.2rem;

      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .challenge-value {
      font-size: 1.5rem;
    }
  }
}

.challenge-tag {
  display: inline-block;
  padding: 0 0.25rem;
  background: rgba(#000000, 0.1);
  border-radius: 0.25rem;
  margin-right: 0.5rem;
}
</style>
