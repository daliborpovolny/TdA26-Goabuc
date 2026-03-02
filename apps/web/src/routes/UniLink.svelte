<script lang="ts">
	import { page } from '$app/state';
	import { on } from 'svelte/events';

	let {
		content = 'fakebutton',
		href,
		onclick,
		bgcolor = 'bg-white',
		hv_bgcolor = 'bg-p-green',
		uppercase = true
	}: {
		content?: string;
		href: string;
		onclick?: () => void;
		bgcolor?: string;
		hv_bgcolor?: string;
		uppercase?: boolean;
	} = $props();
</script>

<a
	{href}
	{onclick}
	class="cursor-pointer rounded-xl border-4 border-s-black {bgcolor} px-8 py-4 text-2xl font-black text-s-black {uppercase
		? 'uppercase'
		: ''} shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:{hv_bgcolor} active:translate-x-1 active:translate-y-1 active:shadow-none"
>
	{content}
</a>
