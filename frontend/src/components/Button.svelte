<script lang="ts">
    import { Button as ShadcnButton } from "$lib/components/ui/button";

    let {
        label = '',
        onClick,
        clazz = '',
        error = false,
        type = 'active-purple',
        active = true,
    }: {
        label?: string;
        onClick: () => void;
        clazz?: string;
        error?: boolean;
        type?: 'active-purple' | 'active-white' | 'active-red';
        active?: boolean;
    } = $props();

    function actionFn() {
        if (active === false) {
            return
        }
        onClick();
    }

    const variantMap = {
        'active-purple': 'default',
        'active-white': 'outline',
        'active-red': 'destructive',
    } as const;
</script>

<ShadcnButton
    variant={error ? 'default' : variantMap[type]}
    class="rounded-lg px-4 py-2.5 text-base font-medium leading-6 {active ? 'cursor-pointer' : 'cursor-wait bg-bg-subtle text-text-muted'} {clazz}"
    disabled={!active}
    onclick={actionFn}
>{label}</ShadcnButton>
