<template>
  <div>Submission {{ submissionCount }}</div>
</template>

<script>
import Vue from "vue";
import gql from "graphql-tag";
export default Vue.extend({
  data() {
    return {
      submissionCount: -1
    };
  },
  apollo: {
    submissionCount: {
      query: gql`
        query submissionCount {
          submissionCount: GetNumberOfSubmissions
        }
      `
    }
  }
});
</script>

<style></style>
