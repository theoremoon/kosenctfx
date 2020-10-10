<template>
  <div class="my-4 mx-8">
    <h1 class="text-2xl">Administration</h1>

    <div>
      <router-link to="/admin/config">Config</router-link>
      <router-link to="/admin/challenges">Challenges</router-link>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
import Vue from "vue";

export default Vue.extend({});
</script>
