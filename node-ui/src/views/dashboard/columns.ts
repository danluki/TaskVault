import type {ColumnDef} from '@tanstack/vue-table'
import type {NodeInfo} from '@/views/dashboard/schema.ts'
import DataTableColumnHeader from "@/views/dashboard/DataTableColumnHeader.vue";
import { h } from 'vue'
import {Badge} from '@/components/ui/badge'
import { labels } from './data.ts'
export const columns: ColumnDef<NodeInfo>[] = [
    {
        accessorKey: 'Name',
        header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Name'}),
        cell: ({ row }) => {
            const label = labels.find((label: any) => label.value === (row.original as any).label)

            return h('div', { class: 'flex space-x-2' }, [
                label ? h(Badge, { variant: 'outline' }, () => label.label) : null,
                h('span', { class: 'max-w-[400px] truncate font-medium' }, row.getValue('Name')),
            ])
        }
    },
    {
        accessorKey: 'Addr',
        header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Addr'}),
        cell: ({ row }) => {
            const label = labels.find((label: any) => label.value === (row.original as any).label)

            return h('div', { class: 'flex space-x-2' }, [
                label ? h(Badge, { variant: 'outline' }, () => label.label) : null,
                h('span', { class: 'max-w-[100px] truncate font-medium' }, row.getValue('Addr')),
            ])
        }
    },
    {
        accessorKey: 'Port',
        header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Port'}),
        cell: ({ row }) => {
            const label = labels.find((label: any) => label.value === (row.original as any).label)

            return h('div', { class: 'flex space-x-2' }, [
                label ? h(Badge, { variant: 'outline' }, () => label.label) : null,
                h('span', { class: 'max-w-[500px] truncate font-medium' }, row.getValue('Port')),
            ])
        }
    },
    {
        accessorKey: 'statusText',
        header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Status'}),
        cell: ({ row }) => {
            const label = labels.find((label: any) => label.value === (row.original as any).label)

            return h('div', { class: 'flex space-x-2' }, [
                label ? h(Badge, { variant: 'outline' }, () => label.label) : null,
                h('span', { class: 'max-w-[500px] truncate font-medium' }, row.getValue('statusText')),
            ])
        }
    },
    {
        accessorKey: 'Tags',
        header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Tags'}),
        cell: ({ row }) => {
            const tags = Object.entries(row.original.Tags).map(([key, value]) =>
                h('span', { class: 'px-2 py-1 rounded-lg border border-gray-400 text-sm' }, `${key}: ${value}`)
            );

            return h('div', { class: 'flex flex-wrap gap-2' }, tags);
        }
    },
]