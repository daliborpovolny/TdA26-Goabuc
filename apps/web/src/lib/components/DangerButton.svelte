<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		onclick,
		type = 'button',
		disabled = false,
		isSaving = false
	}: {
		children: Snippet;
		onclick?: (e: MouseEvent) => void;
		type?: 'button' | 'submit';
		disabled?: boolean;
		isSaving?: boolean;
	} = $props();
</script>

<button
	{type}
	{onclick}
	disabled={disabled || isSaving}
	class="
        group /* The Shadow and Movement */
        /* Extra depth for the actual click */
        /* Disabled
        
        State */ relative cursor-pointer overflow-hidden rounded-xl
        border-4
        border-s-black bg-red-500 px-8
        py-3
        
        text-xl font-black tracking-widest text-white uppercase shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all duration-100
        hover:translate-x-1 hover:translate-y-1
        
        hover:bg-red-600 hover:shadow-none active:translate-x-1.5 active:translate-y-1.5
        disabled:translate-x-0 disabled:translate-y-0 disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]
    "
>
	{#if isSaving}
		<span class="flex items-center gap-2">
			<span class="animate-pulse">⏳</span> Processing...
		</span>
	{:else}
		{@render children()}
	{/if}
</button>
