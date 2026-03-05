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
        group /* Neobrutalist Shadow & Movement */
        /* Slightly darker/richer green on hover */ /*
        Tactical Click
        
        Depth */ /* Logic for Disabled/Saving
        states
        */ relative cursor-pointer
        overflow-hidden rounded-xl border-4 border-s-black bg-p-green px-8 py-3 text-xl
        
        font-black tracking-widest text-s-black uppercase shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]
        transition-all duration-100
        
        hover:translate-x-1 hover:translate-y-1 hover:bg-[#7ce499] hover:shadow-none active:translate-x-1.5 active:translate-y-1.5
        disabled:translate-x-0 disabled:translate-y-0 disabled:cursor-not-allowed disabled:opacity-50
    "
>
	{#if isSaving}
		<div class="flex items-center gap-2">
			<svg
				class="h-5 w-5 animate-spin text-s-black"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
			>
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
				></path>
			</svg>
			<span>Saving...</span>
		</div>
	{:else}
		{@render children()}
	{/if}
</button>
