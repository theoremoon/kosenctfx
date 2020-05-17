
<template>
  <div class="m-1">
    <div>
      <Button class="mx-1">Challenges</Button>
      <Button class="mx-1">Teams</Button>
      <Button class="mx-1">Users</Button>
    </div>

    <router-view></router-view>
  </div>
</template>

<script>
import Button from "@/components/Button";

export default {
  components: {
    Button
  }
};
</script>