<script lang="ts">
	import type { FeedPost } from '$lib/types';
	import { fade } from 'svelte/transition';

	let { post, courseId }: { post: FeedPost; courseId: string } = $props();

	let message = $state(post.message);
	let initialMessage = $state(post.message);
	let isSaving = $state(false);
	let showSuccess = $state(false);

	const isManual = post.type === 'manual';
	const isEdited = $derived(message !== initialMessage);

	async function editFeedPost(e: Event) {
		e.preventDefault();
		isSaving = true;

		const res = await fetch(`/api/courses/${courseId}/feed/${post.uuid}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ message })
		});

		if (res.ok) {
			initialMessage = message;
			showSuccess = true;
			setTimeout(() => (showSuccess = false), 2000);
		}
		isSaving = false;
	}

	async function deleteFeedPost() {
		if (!confirm('Delete this post forever?')) return;
		await fetch(`/api/courses/${courseId}/feed/${post.uuid}`, { method: 'DELETE' });
	}
</script>

<div class="group relative">
	<div class="absolute inset-0 translate-x-1.5 translate-y-1.5 rounded-xl bg-s-black"></div>

	<form
		onsubmit={editFeedPost}
		class="relative flex flex-col gap-3 rounded-xl border-2 border-s-black p-4 transition-colors
        {isManual ? 'bg-white' : 'bg-gray-100 opacity-80'}"
	>
		<div class="flex items-center justify-between border-b border-gray-100 pb-2">
			<span
				class="text-[10px] font-black tracking-widest uppercase {isManual
					? 'text-p-blue'
					: 'text-gray-500'}"
			>
				{isManual ? 'Teacher Post' : 'System Automation'}
			</span>
			<div class="flex items-center gap-2">
				{#if showSuccess}
					<span transition:fade class="text-[10px] font-bold text-p-green uppercase">✓ Updated</span
					>
				{/if}
				<span class="text-[10px] font-bold text-gray-400">{post.createdAt}</span>
			</div>
		</div>

		{#if isManual}
			<textarea
				name="message"
				class="min-h-20 w-full resize-none bg-transparent font-bold text-s-black focus:outline-none"
				bind:value={message}
				required
			></textarea>

			<div class="flex items-center justify-end gap-2 pt-2">
				{#if isEdited}
					<button
						transition:fade
						disabled={isSaving}
						class="cursor-pointer rounded-lg border-2 border-s-black bg-p-green px-3 py-1 text-xs font-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
					>
						{isSaving ? '...' : 'Save Edit'}
					</button>
				{/if}

				<button
					type="button"
					onclick={deleteFeedPost}
					class="cursor-pointer rounded-lg border-2 border-s-black bg-red-500 px-3 py-1 text-xs font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] transition active:translate-y-0.5 active:shadow-none"
				>
					Delete
				</button>
			</div>
		{:else}
			<p class="py-2 font-medium text-gray-600 italic">
				{message}
			</p>
		{/if}

		{#if post.createdAt !== post.updatedAt}
			<p class="text-[9px] font-bold text-gray-300 uppercase">Edited: {post.updatedAt}</p>
		{/if}
	</form>
</div>
