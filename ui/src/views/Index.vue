<template>
  <div class="my-4 mx-8">
    <h1 class="text-4xl">{{ $store.ctfName }}</h1>
    <div class="ml-4">
      <p>{{ startTime }} 〜 {{ endTime }}</p>
      <p v-if="willHold">CTF will start in {{ countDown }}</p>
      <p v-else-if="nowRunning">CTF now running!</p>
      <p v-else-if="hasEnd">CTF is over. Thanks for playing!</p>
    </div>

    <div class="mt-4 justify-center flex">
      <img src="../assets/neko.png" />
    </div>

    <div class="mt-4">
      <h2 class="text-2xl">[ About ]</h2>
      <p class="ml-4">TBD</p>
    </div>
    <div class="mt-4">
      <h2 class="text-2xl">[ Contact ]</h2>
      <p class="ml-4">TBD</p>
    </div>
    <div class="mt-4">
      <h2 class="text-2xl">[ Rules ]</h2>
      <div class="ml-4">
        <ul class="list-disc">
          <li>
            There is no restriction on the team size, your age nationality. It
            is unrelated whether you are in the Kosen or not.
          </li>
          <li>
            The team who earns more points will be placed higher. If two teams
            have the same amount of points, the team who reached the score
            earlier will win (except for the survey*).
          </li>
          <li>
            Sharing solutions or hints with other teams during the competition
            is forbidden.
          </li>
          <li>
            Attacking the score server is forbidden. We may disqualify and ban
            the team which attacks the score server. Attacking other teams is
            forbidden as well.
          </li>
          <li>
            You are not allowed to brute-force the flag. The form will be locked
            for a while if you submit wrong flags 5 times successively.
          </li>
          <li>You may not play the CTF in multiple teams.</li>
          <li>
            You may not have multiple accounts. In case you can't log in to your
            account, please contact us in Discord.
          </li>
        </ul>
        <p class="text-sm">
          *In the late of the competition we will open a survey as a challenge,
          which has points as well as other challenges. However, even if the
          lower-ranked team among those who have the same amount of points
          solves the survey earlier, the originally higher-ranked teams will be
          placed higher as long as they solve the survey during the competition.
          Be noticed this is an exception to give the participants enough time
          to answer the survey while encouraging them to submit it.
        </p>
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
