<script setup lang="ts">
import Menu from '@/views/home/Menu.vue'
import Header from '@/views/home/Header.vue'
</script>

<template>
  <div class="flex flex-col w-full h-full">
    <Header class="w-full h-full"/>
    <div class="flex flex-row flex-grow w-full overflow-hidden">
      <Menu/>
      <div class="w-full overflow-y-scroll">
        <router-view />
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>