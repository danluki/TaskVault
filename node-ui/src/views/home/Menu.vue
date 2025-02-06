<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import {SidebarProvider, SidebarContent, Sidebar} from "@/components/ui/sidebar";

useLocalStorage('taskvaultSidebarIsOpen', true)

</script>

<template>
  <SidebarProvider>
    <Sidebar collapsible="icon">
      <SidebarContent>
        <div>123</div>
      </SidebarContent>

      <SidebarSeparator />

      <SidebarContent>
        <div>1234</div>
      </SidebarContent>
    </Sidebar>
  </SidebarProvider>

</template>

<style scoped>

</style>