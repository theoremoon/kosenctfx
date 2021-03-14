<template>
  <div class="my-4 mx-8">

    <div class="mt-4 justify-center flex">
      <img src="../assets/piyo.png">
    </div>

    <h1 class="text-4xl text-center">{{ $store.ctfName }}</h1>
    <div class="text-center">
      <p>{{ startTime }} &ndash; {{ endTime }}</p>
      <p v-if="willHold">CTF will start in {{ countDown }}.</p>
      <p v-else-if="nowRunning">CTF is now running!</p>
      <p v-else-if="hasEnd">CTF is over. Thanks for playing!</p>
    </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ About ]</h2>
      <p class="ml-4">We provide many fun challenges of varying difficulty and categories, none of them requiring any guessing skill.</p>
    </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ Contact ]</h2>
      <p class="ml-4">
        Discord:
      </p>
    </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ Rules ]</h2>
      <div class="ml-4">
        <ul class="list-disc list-inside">
          <li>Your team can be of any size.</li>
          <li>Anyone is allowed to participate: no restriction on age or nationality.</li>
          <li>Your position on the scoreboard depends on 2 things: 1) your total number of points (higher is better); 2) the timestamp of your last solved challenge (lower is better).</li>
          <li>The survey challenge is special: it does award you some points, but it doesn't update your "last solved challenge" timestamp. You can't get ahead simply by solving the survey faster.</li>
          <li>You can't brute-force flags. If you submit 5 incorrect flags in a short succession, the flag submission form will get locked for 5 minutes.</li>
          <li>One person can participate in only one team.</li>
          <li>Sharing solutions, hints or flags with other teams during the competition is strictly forbidden.</li>
          <li>You are not allowed to attack the scoreserver.</li>
          <li>You are not allowed to attack other teams.</li>
          <li>You are not allowed to have multiple accounts. If you can't log in to your account, please contact us on Discord.</li>
          <li>We reserve the right to ban and disqualify any team that chooses to break any of these rules.</li>
          <li>The flag format is <code>flag\{[\x20-\x7e]+\}</code>, unless specified otherwise.</li>
          <li>Most importantly: good luck and have fun!</li>
        </ul>
      </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ Prizes ]</h2>
      <p class="ml-4">TBD</p>
    </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ Sponsors ]</h2>
      <div class="sponsors ml-4">
      <p class="ml-4">TBD</p>
      </div>
    </div>

    </div>
  </div>
</template>

<script>
import Vue from "vue";
import { dateFormat } from "../dateformat";
export default Vue.extend({
  data() {
    return {
      now: new Date().valueOf() / 1000,
    };
  },
  mounted() {
    setInterval(() => {
      this.now = Math.floor(new Date().valueOf() / 1000);
    }, 1000);
  },
  methods: {
    dateFormat(t) {
      return dateFormat(t);
    },
  },
  computed: {
    startTime() {
      return this.dateFormat(this.$store.ctfStart);
    },
    endTime() {
      return this.dateFormat(this.$store.ctfEnd);
    },
    willHold() {
      return this.now < this.$store.ctfStart;
    },
    nowRunning() {
      return this.$store.ctfStart <= this.now && this.now <= this.$store.ctfEnd;
    },
    hasEnd() {
      return this.now > this.$store.ctfEnd;
    },
    countDown() {
      const d = this.$store.ctfStart - this.now;

      const days = ("" + Math.floor(d / (60 * 60 * 24))).padStart(2, "0");
      const hours = (
        "" + Math.floor((d % (60 * 60 * 24)) / (60 * 60))
      ).padStart(2, "0");
      const minutes = ("" + Math.floor((d % (60 * 60)) / 60)).padStart(2, 0);
      const seconds = ("" + Math.floor(d % 60)).padStart(2, 0);
      return days + "d " + hours + ":" + minutes + ":" + seconds;
    },
  },
});
</script>


<style lang="scss" scoped>
.sponsors {
  padding-bottom: 20px;
  text-align: center;
}
.sponsors a {
  display: inline-block;
  margin-right: 1.5rem;
  margin-top: 0.5rem;
  text-align: center;
  font-weight: bold;

}
.sponsors img {
  width: 12rem;
  max-width: 100%;
  border-radius: 5px;
}
ul, ol {
  list-style-position: inside;
}
ol {
  list-style-type: decimal;
  margin-left: 1.5rem;
}
.ml-4 {
  margin-left: 2rem;
}
p, ul li {
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
}
</style>
