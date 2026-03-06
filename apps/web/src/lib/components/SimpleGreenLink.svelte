<script lang="ts">
	import type { Snippet } from 'svelte';

	let {
		children,
		link,
		class: className = ''
	}: {
		children: Snippet;
		link: string;
		class?: string;
	} = $props();

	const baseClass = `
		rounded-xl border-4 border-s-black bg-p-green px-8 py-4 text-2xl font-black text-s-black uppercase shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-1.5 active:translate-y-1.5"
    `;
</script>

<a href={link} class="{baseClass} {className}">
	{@render children()}
</a>
