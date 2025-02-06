<script setup lang="ts">
import {Card} from "@/components/ui/card";
import { IconGitBranch } from "@tabler/icons-vue";
import { defineProps } from "vue";

defineProps<{
    leaderName?: string,
}>()

</script>

<template>
  <Card class="w-[250px] rounded-xl p-4 flex items-center">
    <div class="flex items-center">
      <div class="p-2 rounded-full">
        <IconGitBranch class="w-6 h-6" />
      </div>
      <div class="ml-3">
        <p class="text-gray-400 text-sm">Leader</p>
        <p class="text-lg font-bold">{{ leaderName }}</p>
      </div>
    </div>
  </Card>
</template>

<style scoped>

</style>