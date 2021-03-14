<template>
  <div>
    <div>
      <button class="button" @click="openPrint">PRINT</button>
    </div>

    <div class="print-cover"></div>
    <div class="print">
      <p class="text-4xl text-center tracking-wide font-bold high-shadow">
        DIPLOMA
      </p>
      <p class="text-4xl text-center tracking-wide font-bold high-shadow">
        {{ $store.ctfName }}
      </p>
      <p class="text-center low-shadow">{{ startTime }} 〜 {{ endTime }}</p>
      <img src="../assets/neko.png" class="mx-auto m-4" />
      <p
        class="text-4xl text-center tracking-wide font-bold text-shadow high-shadow"
      >
        <span class="border-b">
          <span v-if="pos == 1"> 🥇 </span>
          <span v-else-if="pos == 2"> 🥈 </span>
          <span v-else-if="pos == 3"> 🥉 </span>
          {{ pos | suffix }} place - {{ teamname }} {{ country | countryFlag }}
        </span>
      </p>
      <p class="text-2xl text-center low-shadow">{{ points }} points</p>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import { message, errorHandle } from "../message";
import API from "@/api";
import { dateFormat } from "@/dateformat";

export default Vue.extend({
  data() {
    return {
      now: new Date().valueOf() / 1000,
      teamname: "",
      country: "",
    };
  },
  mounted() {
    if (!this.hasEnd()) {
      message(this, "Diploma is available after Competition ends");
      this.$router.push("/");
    }
    this.getInfo();
  },
  methods: {
    getInfo() {
      API.get("team/" + this.$route.params.id)
        .then((r) => {
          this.teamname = r.data.teamname;
          this.country = r.data.country;
        })
        .catch((e) => {
          errorHandle(this, e);
          this.$router.push("/");
        });
    },
    hasEnd() {
      return this.now > this.$store.ctfEnd;
    },
    openPrint() {
      window.print();
    },
  },
  computed: {
    startTime() {
      return dateFormat(this.$store.ctfStart);
    },
    endTime() {
      return dateFormat(this.$store.ctfEnd);
    },
    teamScore() {
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
      return score;
    },
    pos() {
      const score = this.teamScore;
      if (!score) {
        return 0;
      }
      return score.pos;
    },
    points() {
      const score = this.teamScore;
      if (!score) {
        return 0;
      }
      return score.points;
    },
  },
  filters: {
    suffix(x) {
      const s = ["th", "st", "nd", "rd"];
      const v = x % 100;
      return x + (s[(v - 20) % 10] || s[v] || s[0]);
    },
  },
});
</script>

<style lang="scss" scoped>
@import "../assets/vars.scss";
@media print {
  .print-cover {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100vh;
    background-color: $bg-color;
    z-index: 1;
  }
}

.print {
  z-index: 3;
  position: relative;
}

.high-shadow {
  text-shadow: 0 0 5px $accent-color;
}
.low-shadow {
  text-shadow: 0 0 2px $bg-color;
}
</style>
