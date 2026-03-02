<script lang="ts">
	import type { disableScrollHandling } from '$app/navigation';
	import { page } from '$app/state';

	let {
		content,
		onclick,
		bgcolor = 'bg-white',
		hv_bgcolor = 'bg-p-green',
		uppercase = false,
		children,
		more_style = '',
		type = 'button',
		disabled = false
	}: {
		content?: string;
		onclick?: () => void;
		bgcolor?: string;
		hv_bgcolor?: string;
		uppercase?: boolean;
		children?: () => any;
		more_style?: string;
		type?: any;
		disabled?: boolean;
	} = $props();
</script>

<button
	{type}
	{onclick}
	{disabled}
	class="cursor-pointer rounded-xl border-4 border-s-black {bgcolor} px-8 py-4 text-2xl font-black text-s-black {uppercase
		? 'uppercase'
		: ''} shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:{hv_bgcolor} active:translate-x-1 active:translate-y-1 active:shadow-none {more_style}"
>
	{content}
	{@render children?.()}
</button>
