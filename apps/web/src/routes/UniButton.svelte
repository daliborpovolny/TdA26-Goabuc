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
		disabled = false,
		text = 'text-2xl',
		px = 'px-8',
		py = 'py-4',
		shadow = true,
		border = true,
		translate = true,
		text_color = 'text-s-black'
	}: {
		content?: string;
		onclick?: (() => void) | ((e: Event) => Promise<void>);
		bgcolor?: string;
		hv_bgcolor?: string;
		uppercase?: boolean;
		children?: () => any;
		more_style?: string;
		type?: any;
		disabled?: boolean;
		text?: string;
		px?: string;
		py?: string;
		shadow?: boolean;
		border?: boolean;
		translate?: boolean;
		text_color?: string;
	} = $props();
</script>

<button
	{type}
	{onclick}
	{disabled}
	class="cursor-pointer {border
		? 'rounded-xl border-4 border-s-black'
		: ''} {bgcolor} {px} {py} {text} font-black {text_color} {uppercase ? 'uppercase' : ''} {shadow
		? 'shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]'
		: ''} transition-all hover:{hv_bgcolor} {translate
		? 'active:translate-x-1 active:translate-y-1'
		: ''} active:shadow-none {more_style}"
>
	{content}
	{@render children?.()}
</button>
