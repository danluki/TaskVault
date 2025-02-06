<script setup lang="ts">
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import Leader from "@/views/dashboard/Leader.vue";
import {onMounted, ref} from "vue";

const leaderName = ref("Unknown Leader");

onMounted(() => {
  leaderName.value = (window as any).TASKVAULT_LEADER || "Unknown Leader";
});

</script>

<template>
  <div>
    <Card>
      <CardHeader>Welcome</CardHeader>
      <CardContent>
        <div class="flex">
          <div class="flex mr-[0.5em]">
            <div class="flex">
                <Leader :leader-name="leaderName"/>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>

</style>