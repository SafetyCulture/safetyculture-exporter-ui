<script lang="ts">
    import { Badge } from "$lib/components/ui/badge";

    let {
        name = '',
        clazz = '',
        type = 'neutral'
    }: {
        name?: string;
        clazz?: string;
        type?: 'success' | 'info' | 'neutral' | 'error' | 'cancelled';
    } = $props();

    const styleMap = {
        success: 'bg-[#E8FDF5] text-[#005438] hover:bg-[#E8FDF5]',
        info: 'bg-[#E5FAFF] text-[#0D75B5] hover:bg-[#E5FAFF]',
        neutral: 'bg-[#EEF1F7] text-[#3F495A] hover:bg-[#EEF1F7]',
        error: 'bg-[#A9282E] text-[#FFF0F1] hover:bg-[#A9282E]',
        cancelled: 'bg-[#FFF0F1] text-[#A9282E] hover:bg-[#FFF0F1]',
    } as const;
</script>

<Badge variant="secondary" class="rounded-2xl px-4 py-1 font-normal {styleMap[type]} {clazz}">{name}</Badge>
