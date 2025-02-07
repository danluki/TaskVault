<script setup lang="ts">
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import Leader from "@/views/dashboard/Leader.vue";
import TotalValues from "@/views/dashboard/TotalValues.vue";
import {onMounted, ref} from "vue";
import members from './members.json';
import DataTable from '@/views/dashboard/DataTable.vue';
import {columns} from "./columns.ts"
import {getManyReference} from "@/api/data-provider.ts";
const leaderName = ref("Unknown Leader");
const totalValues = ref("0")

// async function fetchMembers(): Promise<NodeList[]> {
//   getManyReference("members", )
// }

onMounted(() => {
  leaderName.value = (window as any).TASKVAULT_LEADER || "Unknown Leader";
  totalValues.value = (window as any).TASKVAULT_TOTAL_PAIRS || "0";
});

</script>

<template>
  <div class="p-4">
    <Card class="mb-3">
      <CardHeader>Welcome</CardHeader>
      <CardContent>
        <div class="flex">
          <div class="flex mr-[0.5em]">
            <div class="flex">
                <Leader :leader-name="leaderName"/>
                <TotalValues :total-values="totalValues" />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
    <Card>
      <CardHeader>
        Nodes
      </CardHeader>
      <CardContent>
        <DataTable :data="members as any"  :columns="columns"/>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>

</style>