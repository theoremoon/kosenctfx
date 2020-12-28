<template>
  <div>
    filter: <input type="text" v-model="filterExpr" />
    <table>
      <thead>
        <tr>
          <th>Team</th>
          <th>Flag</th>
          <th>IsCorrect</th>
          <th>IsValid</th>
          <th>Challenge</th>
          <th>SubmittedAt</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in applyFilter(submissions, filterExpr)" :key="s.id">
          <td>{{ s.team.name }}</td>
          <td>
            <pre>{{ s.flag }}</pre>
          </td>
          <td>{{ s.isCorrect ? "⭕" : "❌" }}</td>
          <td>{{ s.isValid ? "✅" : "" }}</td>
          <td>{{ s.challenge ? s.challenge.name : "-" }}</td>
          <td>{{ s.submittedAt }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import Vue from "vue";
import gql from "graphql-tag";
export default Vue.extend({
  data() {
    return {
      submissionCount: -1,
      submissions: [],
      filterExpr: ""
    };
  },
  apollo: {
    submissions: {
      query: gql`
        query submissions($page: PaginationInput!) {
          submissionCount: GetNumberOfSubmissions
          submissions: ListSubmissions(page: $page) {
            id
            flag
            challenge {
              name
            }
            isCorrect
            isValid
            submittedAt
            team {
              id
              name
            }
          }
        }
      `,
      variables: {
        page: {
          offset: 0,
          limit: 100
        }
      }
    }
  },
  methods: {
    applyFilter(list, filterexpr) {
      if (!filterexpr) {
        return list;
      }

      return list.filter(x =>
        function(e) {
          return new Function(e);
        }.call(x, filterexpr)
      );
    }
  }
});
</script>

<style lang="scss" scoped>
@import "../../assets/vars.scss";
@import "../../assets/tailwind.css";

td,
th {
  max-width: 20rem;
  border-bottom: 1px solid $fg-color;

  pre {
    white-space: pre-line;
  }
}
tbody {
  tr {
    &:nth-child(odd) {
      background-color: rgba($fg-color, 0.1);
    }
  }
}
</style>
