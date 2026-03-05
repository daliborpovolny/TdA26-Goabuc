<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		onclick,
		href,
		type = 'button',
		disabled = false,
		isSaving = false,
		class: className = ''
	}: {
		children: Snippet;
		onclick?: (e: MouseEvent | KeyboardEvent) => void;
		href?: string;
		type?: 'button' | 'submit';
		disabled?: boolean;
		isSaving?: boolean;
		class?: string;
	} = $props();

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onclick?.(e);
		}
	}

	const baseClass = `
        group relative inline-flex items-center justify-center cursor-pointer 
        overflow-hidden rounded-xl border-4 border-s-black 
        bg-p-blue px-8 py-3 text-xl font-black uppercase tracking-widest text-white
        transition-all duration-100
        
        /* Neobrutalist Shadow & Movement */
        shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]
        active:translate-x-1 active:translate-y-1 active:shadow-none
        hover:bg-s-2 /* Your lighter blue for hover */
        
        /* State Handling */
        disabled:cursor-not-allowed disabled:opacity-50 disabled:translate-x-0 disabled:translate-y-0
    `;
</script>

<svelte:element
	this={href ? 'a' : 'button'}
	{href}
	{type}
	{onclick}
	onkeydown={handleKeyDown}
	{disabled}
	role="button"
	tabindex="0"
	class="{baseClass} {className}"
>
	{#if isSaving}
		<div class="flex items-center gap-2">
			<svg class="h-5 w-5 animate-spin text-white" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
				></path>
			</svg>
			<span>Wait...</span>
		</div>
	{:else}
		{@render children()}
	{/if}
</svelte:element>
