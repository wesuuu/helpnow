<script lang="ts">
    interface Column {
        key?: string;
        label: string;
        class?: string;
        format?: (value: any, item: any) => string | number;
    }

    let {
        columns = [],
        data = [],
        onRowClick = undefined,
        rowClass = "",
        children,
    } = $props<{
        columns: Column[];
        data: any[];
        onRowClick?: (item: any) => void;
        rowClass?: string | ((item: any) => string);
        children?: import("svelte").Snippet<[any, string]>; // (item, columnKey)
    }>();
</script>

<div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
    <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700">
            <tr>
                {#each columns as col}
                    <th
                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider {col.class ||
                            ''}"
                    >
                        {col.label}
                    </th>
                {/each}
            </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            {#each data as item}
                <tr
                    class="hover:bg-gray-50 dark:hover:bg-gray-700 {onRowClick
                        ? 'cursor-pointer'
                        : ''} {typeof rowClass === 'function'
                        ? rowClass(item)
                        : rowClass}"
                    onclick={() => onRowClick && onRowClick(item)}
                >
                    {#each columns as col}
                        <td
                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white"
                        >
                            {#if col.key}
                                {@const val = item[col.key]}
                                {#if children}
                                    {@render children(item, col.key)}
                                {:else if col.format}
                                    {col.format(val, item)}
                                {:else}
                                    {val}
                                {/if}
                            {:else}
                                <!-- Column without key, expects custom render or manually handled if children supports it -->
                                {#if children}
                                    {@render children(item, col.key || "")}
                                {/if}
                            {/if}
                        </td>
                    {/each}
                </tr>
            {/each}
        </tbody>
    </table>
</div>
