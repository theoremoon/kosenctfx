<template>
  <div class="flex">
    <div>
      <h2 class="text-lg">TABLES</h2>
      <ul class="list-none font-mono">
        <li
          v-for="t in tables"
          :key="t"
          @click="getTableData(t)"
          class="hover:underline cursor-pointer"
        >
          {{ t }}
        </li>
      </ul>
    </div>
    <div class="flex-grow px-4">
      <h2 class="text-lg">QUERY</h2>
      <form @submit.prevent="userQuery(query)">
        <div>
          <input
            type="text"
            v-model="query"
            class="font-mono"
            placeholder="SELECT 1"
          />
        </div>
        <div>
          <input type="submit" value="DO QUERY" />
        </div>
        <div>
          <label for="autolimit"
            >Automatially append `LIMIT 100` to query
            <input id="autolimit" type="checkbox" v-model="autolimit" />
          </label>
        </div>
      </form>
      <div style="w-full overflow-auto">
        <table class="border-collapse">
          <thead>
            <tr>
              <th>INDEX</th>
              <th v-for="c in columns" :key="c">{{ c }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(r, i) in rows" :key="'row' + i">
              <td>{{ i + 1 }}</td>
              <td
                v-for="(d, j) in r"
                :key="'row' + i + 'col' + j"
                @click="autoQuery(d, columns[j], current_table)"
                :title="
                  `SELECT * FROM \`${current_table}\` WHERE \`${columns[j]}\` = ${d}`
                "
                class="hover:underline cursor-pointer"
              >
                {{ d }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import API from "@/api";
import { errorHandle } from "@/message";
export default Vue.extend({
  data() {
    return {
      query: "",
      tables: [],
      current_table: null,
      columns: [],
      rows: [],
      autolimit: true
    };
  },
  async mounted() {
    this.getTables();
  },
  methods: {
    doQuery(query) {
      return API.post("/admin/sql", { query: query })
        .then(r => {
          return r.data;
        })
        .catch(e => errorHandle(this, e));
    },
    async getTables() {
      const { rows } = await this.doQuery("");
      rows.forEach(x => {
        this.tables.push(x["TABLE_NAME"]);
      });
    },
    async getTableData(table) {
      this.query = `SELECT * FROM \`${table}\``;
      this.current_table = table;
      this.userQuery();
    },
    async setData(cols, rows) {
      this.columns = cols;
      this.rows.length = [];
      rows.forEach(row => {
        this.rows.push(this.columns.map(c => row[c]));
      });
    },
    async userQuery() {
      const { columns, rows } = await this.doQuery(
        this.query.replace(/;\s*$/, "") + (this.autolimit ? " LIMIT 100" : "")
      );
      this.setData(columns, rows);
    },
    async autoQuery(data, column, table) {
      this.query = `SELECT * FROM \`${table}\` WHERE \`${column}\` = "${data}"`;
      this.userQuery();
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../../assets/vars.scss";

table {
  width: 100%;
  overflow: auto;
  border: 2px solid $fg-color;
  border-collapse: collapse;
}
table::-webkit-scrollbar {
  display: none;
}
th::-webkit-scrollbar {
  display: none;
}
td::-webkit-scrollbar {
  display: none;
}

thead th {
  border-bottom: 2px solid $fg-color;
}

td,
th {
  word-break: keep-all;
  white-space: nowrap;
  overflow: auto;

  max-width: 5rem;
  border-right: 1px solid $fg-color;
}
</style>
