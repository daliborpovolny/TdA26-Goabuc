<script lang="ts">
	import DangerButton from '$lib/components/DangerButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';
	import type { FeedPost } from '$lib/types';
	import { fade } from 'svelte/transition';

	let { post, courseId }: { post: FeedPost; courseId: string } = $props();

	let message = $derived(post.message);
	let initialMessage = $derived(post.message);
	let isSaving = $state(false);
	let showSuccess = $state(false);

	const isManual = $derived(post.type === 'manual');
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
		class="relative flex flex-col rounded-xl border-2 border-s-black p-4 transition-colors
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

			<div class="flex items-center justify-between">
				{#if post.createdAt !== post.updatedAt}
					<p class="text-[9px] font-bold text-gray-300 uppercase">
						Edited: {post.updatedAt}
					</p>
				{:else}
					<div></div>
				{/if}

				<div class="flex items-center gap-2">
					{#if isEdited}
						<SuccessButton type="submit" disabled={isSaving} class="text-xs">
							{isSaving ? '...' : 'Save Edit'}
						</SuccessButton>
					{/if}

					<DangerButton type="button" onclick={deleteFeedPost} class="text-center text-xs">
						Delete
					</DangerButton>
				</div>
			</div>
		{:else}
			<p class="py-2 font-medium text-gray-600 italic">
				{message}
			</p>
		{/if}
	</form>
</div>
