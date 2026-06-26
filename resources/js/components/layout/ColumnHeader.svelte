<script lang="ts">
    import Menu from '@lucide/svelte/icons/menu';
    import type { Snippet } from 'svelte';

    import {
        availableColumnOptions,
        COLUMN_CONTENT_LABELS,
    } from '@/lib/studyLayout';
    import { setColumnContent, study } from '@/lib/study.svelte.ts';
    import type { ColumnContentType } from '@/lib/types/study';

    let {
        contentType,
        slotIndex,
        showViewSelector = false,
        children,
    }: {
        contentType: ColumnContentType | 'primary-bible';
        slotIndex?: number;
        showViewSelector?: boolean;
        children?: Snippet;
    } = $props();

    let menuOpen = $state(false);
    let triggerEl = $state<HTMLButtonElement | null>(null);
    let panelEl = $state<HTMLUListElement | null>(null);

    const columnOptions = $derived(
        slotIndex === undefined
            ? []
            : availableColumnOptions(
                slotIndex,
                study.columns,
                study.translationBId,
                study.translationCId,
            ),
    );

    function closeMenu(): void {
        menuOpen = false;
    }

    function selectView(type: ColumnContentType): void {
        if (slotIndex === undefined) {
            return;
        }

        setColumnContent(slotIndex, type);
        closeMenu();
    }

    function toggleMenu(event: MouseEvent): void {
        event.stopPropagation();
        menuOpen = ! menuOpen;
    }

    $effect(() => {
        if (! menuOpen) {
            return;
        }

        function handlePointerDown(event: PointerEvent): void {
            const target = event.target as Node;

            if (triggerEl?.contains(target) || panelEl?.contains(target)) {
                return;
            }

            closeMenu();
        }

        document.addEventListener('pointerdown', handlePointerDown);

        return () => {
            document.removeEventListener('pointerdown', handlePointerDown);
        };
    });
</script>

<div class="border-base-300 flex shrink-0 items-center gap-2 border-b px-3 py-2">
    {#if showViewSelector && slotIndex !== undefined}
        <div
            class="dropdown"
            class:dropdown-open={menuOpen}
            class:dropdown-end={slotIndex === 0}
            class:z-50={menuOpen}
        >
            <button
                bind:this={triggerEl}
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                aria-label="Change column view"
                aria-expanded={menuOpen}
                aria-haspopup="menu"
                onclick={toggleMenu}
            >
                <Menu size={16} aria-hidden="true" />
            </button>
            <ul
                bind:this={panelEl}
                class="dropdown-content menu bg-base-100 rounded-box z-50 w-44 border border-base-300 p-1 shadow-lg"
                role="menu"
            >
                {#each columnOptions as option (option)}
                    <li role="none">
                        <button
                            type="button"
                            role="menuitem"
                            class:menu-active={contentType === option}
                            onclick={() => selectView(option)}
                        >
                            {COLUMN_CONTENT_LABELS[option]}
                        </button>
                    </li>
                {/each}
            </ul>
        </div>
    {/if}

    {#if children}
        <div class="flex min-w-0 flex-1 items-center gap-2">
            {@render children()}
        </div>
    {/if}
</div>
