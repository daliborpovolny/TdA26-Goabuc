<script lang="ts">
	import type { FeedPost } from '$lib/types';

	let { post }: { post: FeedPost } = $props();
</script>

<div
	class="group relative flex flex-col gap-2 rounded-xl border-2 border-s-black bg-white p-4 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none"
>
	<div class="flex items-center justify-between">
		<span
			class="rounded-lg px-2 py-1 text-xs font-bold tracking-wider uppercase
            {post.type === 'system' ? 'bg-s-3 text-white' : 'bg-p-green text-s-black'}"
		>
			{post.type === 'system' ? 'System Alert' : 'Teacher Message'}
		</span>

		<div class="text-right text-xs font-medium text-gray-500">
			{post.createdAt}
		</div>
	</div>

	<p class="text-lg leading-relaxed text-s-black">
		{post.message}
	</p>

	<div class="mt-2 flex items-center gap-3 pt-2 text-[10px] font-bold text-gray-400 uppercase">
		{#if post.edited}
			<span>Edited: {post.updatedAt}</span>
		{/if}
	</div>
</div>
