<script lang="ts">
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import { on } from 'svelte/events';

	let {
		content = 'fakebutton',
		href,
		onclick,
		bgcolor = 'bg-white',
		hv_bgcolor = 'bg-p-green',
		uppercase = false,
		children,
		more_style = '',
		title = '',
		target = '',
		px = 'px-8',
		py = 'py-4'
	}: {
		content?: string;
		href: string;
		onclick?: () => void;
		bgcolor?: string;
		hv_bgcolor?: string;
		uppercase?: boolean;
		children?: Snippet;
		more_style?: string;
		title?: string;
		target?: string;
		px?: string;
		py?: string;
	} = $props();
</script>

<a
	{href}
	{onclick}
	{target}
	{title}
	class="cursor-pointer rounded-xl border-4 border-s-black {bgcolor} {px} {py} text-2xl font-black text-s-black {uppercase
		? 'uppercase'
		: ''} shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:{hv_bgcolor} hover:translate-x-1 hover:translate-y-1 active:translate-x-1.5 active:translate-y-1.5 hover:shadow-none {more_style}"
>
	{content}
	{@render children?.()}
</a>
