<template>
  <div class="my-4 mx-8">
    <h1 class="text-4xl">InterKosenCTF 2020</h1>
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
      <p class="ml-4">
        Welcome to InterKosenCTF 2020!<br />
        InterKosenCTF2020 is a Jeopardy-style Capture The Flag competition
        hosted by insecure(<a href="https://twitter.com/ptrYudai">ptr-yudai</a>,
        <a href="https://twitter.com/theoremoon">theoldmoon0602</a>, and
        <a href="https://twitter.com/y05h1k1ng">yoshiking</a>), a Japanese CTF
        team. There will be challenges of mainly 4 categories (pwn, web, rev,
        crypto) for begginer - midium level players. The flag format is
        <code>KosenCTF{[\x21-\x7a]+}</code> unless otherwise specified.<br />
      </p>
    </div>
    <div class="mt-4">
      <h2 class="text-2xl">[ Contact ]</h2>
      <p class="ml-4">
        Discord:
        <a href="https://discord.gg/" target="_blank">https://discord.gg/</a>
      </p>
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
import dayjs from "dayjs";

export default Vue.extend({
  data() {
    return {
      now: 0
    };
  },
  methods: {
    dateFormat(ts) {
      return dayjs(ts * 1000).format("YYYY-MM-DD HH:mm:ss Z");
    }
  },
  mounted() {
    setInterval(() => {
      this.now = Math.floor(new Date().valueOf() / 1000);
    }, 1000);
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
      return this.now > this.ctfEnd;
    },
    countDown() {
      const d = this.ctfStart - this.now;
      const days = ("" + Math.floor(d / (60 * 60 * 24))).padStart(2, "0");
      const hours = (
        "" + Math.floor((d % (60 * 60 * 24)) / (60 * 60))
      ).padStart(2, "0");
      const minutes = ("" + Math.floor((d % (60 * 60)) / 60)).padStart(2, 0);
      const seconds = ("" + Math.floor(d % 60)).padStart(2, 0);
      return days + "d " + hours + ":" + minutes + ":" + seconds;
    }
  }
});
</script>
