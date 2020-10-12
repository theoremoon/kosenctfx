<template>
  <div class="my-4 mx-8">
    <h2 class="text-2xl">Configuration</h2>
    <div class="flex">
      <form class="w-1/2 ml-4" @submit.prevent="update">
        <div class="mb-4">
          <label class="block text-sm" for="ctfname">
            CTF Name
          </label>
          <input type="text" v-model="ctfname" id="ctfname" />
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="start_at">
            CTF Start at
          </label>
          <input type="text" v-model="start_at" id="start_at" />
          <div class="text-sm" v-if="!checkTimeString(start_at)">
            invalid format
          </div>
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="end_at">
            CTF End at
          </label>
          <input type="text" v-model="end_at" id="end_at" />
          <div class="text-sm" v-if="!checkTimeString(end_at)">
            invalid format
          </div>
        </div>

        <div class="mb-4">
          <label class="block text" for="ctf_open"
            ><input type="checkbox" v-model="ctf_open" id="ctf_open" /> CTF is
            <span v-if="ctf_open">Open</span><span v-else>Closed</span>
          </label>
        </div>
        <div class="mb-4">
          <label class="block text" for="register_open"
            ><input
              type="checkbox"
              v-model="register_open"
              id="register_open"
            />
            Registration is <span v-if="register_open">Open</span
            ><span v-else>Closed</span>
          </label>
        </div>

        <div class="mb-4">
          <label class="block text-sm" for="score_expr">
            Score Expr
          </label>
          <textarea
            v-model="score_expr"
            id="score_expr"
            style="font-family: monospace"
          ></textarea>
          <button @click.prevent="drawGraph">draw graph</button>
        </div>

        <input type="submit" value="apply changes" class="float-right" />
      </form>
      <div class="w-1/2">
        <canvas ref="chart"></canvas>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { errorHandle } from "../../message";
import { dateFormat, parseDateString } from "@/dateformat";
import Chart from "chart.js/dist/Chart";

export default Vue.extend({
  data() {
    return {
      ctfname: "",
      start_at: "",
      end_at: "",
      score_expr: "",
      ctf_open: false,
      register_open: false,
      chart: null
    };
  },
  mounted() {
    this.getConfig();
    let ctx = this.$refs.chart.getContext("2d");
    this.chart = new Chart(ctx, {
      type: "line",
      data: [],
      options: {
        responsive: true,
        maintainAspectRatio: false
      }
    });
  },
  methods: {
    getConfig() {
      API.get("admin/get-config")
        .then(r => {
          this.ctfname = r.data.ctf_name;
          this.start_at = dateFormat(r.data.start_at);
          this.end_at = dateFormat(r.data.end_at);
          this.score_expr = r.data.score_expr;
          this.ctf_open = r.data.ctf_open;
          this.register_open = r.data.register_open;
        })
        .catch(() => {
          errorHandle(this, "forbidden");
          this.$router.push("/");
        });
    },
    update() {
      API.post("admin/set-config", {
        name: this.ctfname,
        start_at: parseDateString(this.start_at),
        end_at: parseDateString(this.end_at),
        score_expr: this.score_expr,
        ctf_open: this.ctf_open,
        register_open: this.register_open
      })
        .then(() => {
          this.getConfig();
          this.$eventHub.$emit("login-check");
        })
        .catch(e => {
          errorHandle(this, e);
        });
    },
    checkTimeString(s) {
      if (parseDateString(s)) {
        return true;
      }
      return false;
    },
    drawGraph() {
      API.get("admin/score-emulate", {
        params: {
          maxCount: 100,
          expr: this.score_expr
        }
      })
        .then(r => {
          this.chart.data.datasets = [
            {
              label: "score",
              borderColor: "#4491cf",
              backgroundColor: "#4491cf",
              pointRadius: 0,
              fill: false,
              data: r.data
            }
          ];
          this.chart.data.labels = r.data.map((_, i) => i);
          this.chart.update();
        })
        .catch(e => {
          console.log(e);
          errorHandle(this, e);
        });
    }
  },
  computed: {}
});
</script>
