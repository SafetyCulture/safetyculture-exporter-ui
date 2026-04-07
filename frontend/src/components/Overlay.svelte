<script lang="ts">
    import type { Snippet } from "svelte";

    let { children }: { children: Snippet } = $props();
</script>

<div class="fixed inset-0 z-[9999] flex cursor-wait items-center justify-center bg-black/75">
    <div class="absolute rounded-lg border border-primary bg-white px-5 py-2.5 text-2xl">
        {@render children()}
    </div>
</div>
