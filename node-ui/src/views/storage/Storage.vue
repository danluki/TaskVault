<script setup lang="ts">
import { ref, watchEffect, computed } from "vue";
import axios from "axios";
import {apiUrl} from "@/api/data-provider.ts"
import Table from "@/components/table.vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import { Plus, PackageOpen } from "lucide-vue-next";
import type {PairInfo} from "@/views/storage/schema.ts";
import Pagination from "@/components/pagination.vue";
import {
  useStorageTable
} from "@/views/storage/composables/use-storage-table.ts";


// Reactive state for pagination and sorting
const page = ref(1);
const pageSize = ref(10);
const sortBy = ref("dc"); // Default sorting column
const sortOrder = ref("asc"); // "asc" or "desc"
const totalItems = ref(0);
const tagsData = ref<Array<Record<string, string>>>([]); // Store fetched tags

// Fetch data function
// const fetchData = async () => {
//   console.log(`${apiUrl}/storage`)
//   try {
//     const response = await axios.get(`${apiUrl}/storage`, {
//       params: {
//         _page: page.value,
//         _limit: pageSize.value,
//         _sort: sortBy.value,
//         _order: sortOrder.value,
//       },
//     });
//
//     tagsData.value = response.data;
//     totalItems.value = parseInt(response.headers["x-total-count"], 10) || 0;
//   } catch (error) {
//     console.error("Error fetching data:", error);
//   }
// };
//
// // Fetch data when pagination, sorting, or page size changes
// watchEffect(fetchData);

// Computed values for pagination
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value));

// Pagination handlers
const nextPage = () => {
  if (page.value < totalPages.value) page.value++;
};

const prevPage = () => {
  if (page.value > 1) page.value--;
};

const keyInput = ref("");
const valueInput = ref("");
const storageValues = ref<PairInfo[]>([]);

const addPair = async () => {

}

const storageTable = useStorageTable()
</script>

<template>
  <div class="p-4">
    <Card class="mb-3">
      <CardHeader>Add value</CardHeader>
      <CardContent>
        <div class="flex space-x-2 mb-4">
          <div class="flex mr-[0.5em]">
            <Input v-model="keyInput" placeholder="Key" class="w-1/3 mr-2" />
            <Input v-model="valueInput" placeholder="Value" class="w-1/3 mr-2" />
            <Button @click="addPair" class="h-10 w-20">
              <Plus class="w-4 h-4" /> Add
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
    <Card>
      <CardHeader>Storage</CardHeader>
      <CardContent>
        <Pagination :total="storageTable.totalPairs.value"
                    :table="storageTable.table"
                    :pagination="storageTable.pagination.value"
                    @update:page="(page) => storageTable.pagination.value.pageIndex = page"
                    @update:page-size="(pageSize) => storageTable.pagination.value.pageSize = pageSize"
        />
        <Table :table="storageTable.table" :is-loading="storageTable.isLoading.value">
          <template #empty-message>
            <div class="flex items-center justify-center h-[400px]">
              <PackageOpen class="w-6 h-6 mr-2" />
              Storage is empty
            </div>
          </template>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>

</style>