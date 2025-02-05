<script setup lang="ts">
import Toaster from '@/components/ui/toast/Toaster.vue'
</script>

<template>
  <Toaster/>
  <router-view />
</template>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
  /* width: 100%; */
}

::-webkit-scrollbar {
  background-color: #71717a;
  border-radius: 10px;
  width: 5px;
  height:5px;
}

::-webkit-scrollbar-thumb {
  border-radius: 10px;
  background: #71717a;
}

::-webkit-scrollbar-track {
  border-radius: 10px;
  background: #fff;
}
</style>
