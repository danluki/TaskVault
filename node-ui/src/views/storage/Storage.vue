<script setup lang="ts">
import {h, ref} from "vue";
import Table from "@/components/table.vue";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {Card, CardContent, CardHeader} from "@/components/ui/card";
import { Plus, PackageOpen } from "lucide-vue-next";
import Pagination from "@/components/pagination.vue";
import {
  useStorageTable
} from "@/views/storage/composables/use-storage-table.ts";
import {ToastAction, useToast} from "@/components/ui/toast";


const keyInput = ref("");
const valueInput = ref("");
const {toast} = useToast()

const storageTable = useStorageTable()

const addPair = async () => {
  if (!keyInput.value || !valueInput.value) {
    toast({
      variant: 'destructive',
      description: 'Key and Value must not be empty.',
    });
    return;
  }

  const requestBody = JSON.stringify({
    key: keyInput.value,
    value: valueInput.value,
  });

  const sendRequest = async (url: any) => {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: requestBody,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
  };

  try {
    await sendRequest("http://localhost:8080/v1/storage");
    
    toast({
      description: 'Pair added successfully.',
    });

    keyInput.value = "";
    valueInput.value = "";
    storageTable.pagination.value.pageIndex = storageTable.pagination.value.pageIndex; 
  } catch (error) {
    console.error('First request failed, trying fallback...', error);

    try {
      await sendRequest("http://localhost:8081/v1/storage");

      toast({
        description: 'Pair added successfully.',
      });

      keyInput.value = "";
      valueInput.value = "";
      storageTable.pagination.value.pageIndex = storageTable.pagination.value.pageIndex;
    } catch (fallbackError) {
      console.error('Both requests failed.', fallbackError);

      toast({
        title: 'Uh oh! Something went wrong.',
        description: 'Could not add the pair after retrying.',
        variant: 'destructive',
        action: h(ToastAction, {
          altText: 'Try again',
        }, {
          default: () => 'Try again',
        }),
      });
    }
  }
};

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