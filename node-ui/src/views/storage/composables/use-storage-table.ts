import {
    type ColumnDef,
    getCoreRowModel,
    getFacetedRowModel,
    getFacetedUniqueValues,
    useVueTable,
} from '@tanstack/vue-table'
import {createGlobalState} from '@vueuse/core'
import {computed, h, ref, watchEffect} from 'vue'
import {usePagination} from "@/composables/use-pagination.ts";
import type {PairInfo} from "@/views/storage/schema.ts";
import {getManyReference} from "@/api/data-provider.ts";
import {valueUpdater} from "@/lib/utils.ts";
import {
    useStorageTableActions
} from "@/views/storage/composables/use-storage-table-actions.ts";

export const TABLE_ACCESSOR_KEYS = {
    key: 'Key',
    value: 'Value',
}

export const useStorageTable = createGlobalState(() => {
    const {
        columnVisibility
    } = useStorageTableActions()

    const { pagination, setPagination} = usePagination()

    const data = ref<PairInfo[]>([]);
    const loading = ref<boolean>(false);
    const error = ref<string | null>(null);
    const pageCount = ref<number>(1); 
    const totalPairs = ref<number>(0);

    const params = computed(() => ({
        page: pagination.value.pageIndex,
        perPage: pagination.value.pageSize,
    }));

    const fetchData = async () =>  {
        loading.value = true;
        error.value = null;
        try {
            const resp = await getManyReference<PairInfo>("storage", params.value);

            data.value = resp.data;
            pageCount.value = Math.ceil(resp.total / params.value.perPage);
            totalPairs.value = resp.total;
        } catch (err) {
            error.value = (err as Error).message;
            console.log(err)
        } finally {
            loading.value = false;
        }
    }

    watchEffect(fetchData);


    const tableColumns = computed<ColumnDef<PairInfo>[]>(() => [
        {
            accessorKey: TABLE_ACCESSOR_KEYS.key,
            size: 20,
            header: () => h('div', {}, 'Key'),
            cell: ({ row }) => {
                return h('div', (row.original as any).Key)
            },
        },
        {
            accessorKey: TABLE_ACCESSOR_KEYS.value,
            size: 20,
            header: () => h('div', {}, 'Value'),
            cell: ({ row }) => {
                return h('div', (row.original as any).Value)
            },
        },
    ])

    const table = useVueTable({
        get pageCount() {
            return pageCount.value
        },
        get data() {
            return data.value
        },
        get columns() {
            return tableColumns.value
        },
        state: {
            get columnVisibility() {
                return columnVisibility.value
            },
            get pagination() {
                return pagination.value
            },
        },
        manualPagination: true,
        enableRowSelection: true,
        onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
        onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
        getCoreRowModel: getCoreRowModel(),
        getFacetedRowModel: getFacetedRowModel(),
        getFacetedUniqueValues: getFacetedUniqueValues(),
    })
    return {
        isLoading: loading,
        table,
        totalPairs,
        pagination,
        setPagination,
        fetchData,
    }
})