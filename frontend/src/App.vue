<template>
  <div>
    <header class="w-full p-4 navbar">
      <nav class="lg:flex">
        <div class="flex flex-shrink-0 text-l mr-6">
          <img src="./assets/neko.png" class="h-6 w-6" />
          InterKosenCTF2020
        </div>

        <div class="flex-grow lg:flex lg:w-auto">
          <div class="lg:mr-4">
            <router-link to="/challenges">challenges</router-link>
          </div>
          <div class="lg:mr-4">
            <router-link to="/ranking">ranking</router-link>
          </div>
        </div>

        <div class="flex-0 flex">
          <div class="mr-4">user</div>
          <div class="mr-4">
            <router-link to="/team">team</router-link>
          </div>
          <div class="mr-4">
            <router-link to="/admin">admin</router-link>
          </div>
        </div>
      </nav>
    </header>
    <div class="container mx-auto">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
@import "../node_modules/tailwindcss/dist/tailwind.css";
@import "../node_modules/tailwindcss/dist/components.css";
@import "../node_modules/tailwindcss/dist/utilities.css";
@import "@/vars.scss";

.navbar {
  background-color: $dark-color;
}
</style>
<style lang="scss">
@import "@/vars.scss";
* {
  color: $fg-color;
}
body {
  background-color: $bg-color;
}
</style>