<script setup lang="ts">
import {onMounted, ref} from "vue";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {RefreshCcw} from "lucide-vue-next";

const metrics = ref("");
const loading = ref(false);

const fetchMetrics = async () => {
  loading.value = true;
  try {
    const response = await fetch("http://localhost:8080/metrics");
    metrics.value = await response.text();
  } catch (error) {
    console.error("Failed to fetch metrics:", error);
  } finally {
    loading.value = false;
  }
};

// Auto-refresh every 10 seconds
onMounted(() => {
  fetchMetrics();
  setInterval(fetchMetrics, 10000);
});
</script>

<template>
  <div class="p-6">
    <Card>
      <CardHeader class="flex justify-between items-center">
        <CardTitle>System Metrics</CardTitle>
        <Button @click="fetchMetrics" :disabled="loading">
          <RefreshCcw class="h-4 w-4 mr-2" /> Refresh
        </Button>
      </CardHeader>
      <CardContent>
        <pre class="p-4 rounded-md text-sm overflow-auto">{{ metrics }}</pre>
      </CardContent>
    </Card>
  </div>
</template>
