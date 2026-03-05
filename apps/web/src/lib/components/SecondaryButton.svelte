<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		onclick,
		href,
		type = 'button',
		disabled = false,
		class: className = ''
	}: {
		children: Snippet;
		onclick?: (e: MouseEvent | KeyboardEvent) => void;
		href?: string;
		type?: 'button' | 'submit';
		disabled?: boolean;
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
        bg-white px-8 py-3 text-xl font-black uppercase tracking-widest text-s-black
        transition-all duration-100
        
        shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]
        active:translate-x-1 active:translate-y-1 active:shadow-none
        hover:bg-gray-100
        
        
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
	{@render children()}
</svelte:element>
