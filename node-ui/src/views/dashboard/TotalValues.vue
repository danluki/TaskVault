<script setup lang="ts">
import {Card} from "@/components/ui/card";
import { IconActivity } from "@tabler/icons-vue";
import { defineProps } from "vue";

defineProps<{
  totalValues?: string,
}>()

</script>

<template>
  <Card class="w-[250px] rounded-xl p-4 flex items-center">
    <div class="flex items-center">
      <div class="p-2 rounded-full">
        <IconActivity class="w-6 h-6" />
      </div>
      <div class="ml-3">
        <p class="text-gray-400 text-sm">Total values</p>
        <p class="text-lg font-bold">{{ totalValues }}</p>
      </div>
    </div>
  </Card>
</template>

<style scoped>

</style>