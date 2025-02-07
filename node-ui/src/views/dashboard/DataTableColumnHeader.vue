<script setup lang="ts">
import type {Column} from '@tanstack/vue-table'
import type { NodeInfo } from './schema.ts'

interface DataTableColumnHeaderProps {
  column: Column<NodeInfo, any>
  title: string
}

defineProps<DataTableColumnHeaderProps>()

</script>

<template>
  <div :class="$attrs.class">
    {{ title }}
  </div>
</template>

<style scoped>

</style>