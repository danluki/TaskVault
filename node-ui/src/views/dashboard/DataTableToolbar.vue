<script setup lang="ts">
import type { Table } from '@tanstack/vue-table'
import type {NodeInfo} from "@/views/dashboard/schema.ts";
import DataTableViewOptions from './DataTableViewOptions.vue'

interface DataTableToolbarProps {
  table: Table<NodeInfo>
}

defineProps<DataTableToolbarProps>()

</script>

<template>
  <div class="flex items-center justify-between">
    <div class="flex flex-1 items-center space-x-2">
    </div>
    <DataTableViewOptions :table="table" />
  </div>
</template>

<style scoped>

</style>