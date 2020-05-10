<template>
  <div>
    <header class="w-full bg-gray-200 p-4">
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

        <div class="flex-0">
          <div>user</div>
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
</style>
