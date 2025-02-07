import { createGlobalState, useLocalStorage } from '@vueuse/core'

import type { VisibilityState } from '@tanstack/vue-table'

const COLUMN_VISIBLE_STORAGE_KEY = 'syncraStorageColumnVisibility'

export const useStorageTableActions = createGlobalState(() => {
    const columnVisibility = useLocalStorage<VisibilityState>(COLUMN_VISIBLE_STORAGE_KEY, {})

    return {
        columnVisibility
    }
})