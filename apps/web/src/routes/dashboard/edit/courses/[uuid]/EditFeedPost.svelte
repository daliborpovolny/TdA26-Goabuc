<script lang="ts">
	import { enhance } from '$app/forms';
	import type { FeedPost } from '$lib/types';
	import { fade } from 'svelte/transition';

	export let post: FeedPost;
	export let courseId: string;

	let isEditing = false;
	let isDeleting = false;

	// Snapshot for reverting if optimistic update fails
	let originalMessage = '';

	function toggleEdit() {
		if (!isEditing) {
			originalMessage = post.message;
		}
		isEditing = !isEditing;
	}
</script>

<div class="rounded-lg border bg-white p-4 shadow-sm transition-all" class:opacity-50={isDeleting}>
	<div class="mb-2 flex items-start justify-between">
		<div class="flex items-center gap-2">
			<span
				class={`rounded-full px-2 py-1 text-xs font-bold uppercase ${post.type === 'manual' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'}`}
			>
				{post.type}
			</span>
			<span class="text-xs text-gray-500">
				{new Date(post.createdAt).toLocaleString()}
				{#if post.edited}
					<span class="ml-1 text-gray-400 italic">(edited)</span>
				{/if}
			</span>
		</div>

		{#if post.type === 'manual'}
			<div class="flex gap-2">
				<button on:click={toggleEdit} class="text-sm text-gray-500 hover:text-blue-600">
					{isEditing ? 'Cancel' : 'Edit'}
				</button>

				<form
					method="POST"
					action={`/courses/${courseId}/feed?/${post.uuid}/delete`}
					use:enhance={() => {
						isDeleting = true;
						return async ({ result }) => {
							// If it fails, undo the visual deletion state
							if (result.type === 'error') isDeleting = false;
						};
					}}
				>
					<button class="text-sm text-red-400 hover:text-red-600">Delete</button>
				</form>
			</div>
		{/if}
	</div>

	{#if isEditing}
		<form
			method="POST"
			action={`/courses/${courseId}/feed?/${post.uuid}/edit`}
			use:enhance={({ formData }) => {
				// OPTIMISTIC UPDATE: Update UI immediately
				const newMessage = formData.get('message')?.toString();
				if (newMessage) {
					post.message = newMessage;
					post.edited = true;
					isEditing = false;
				}

				return async ({ result, update }) => {
					// If server fails, revert to snapshot
					if (result.type === 'error') {
						post.message = originalMessage;
						alert('Failed to update post');
					}
					// Apply server response (which will confirm the edit)
					await update();
				};
			}}
		>
			<textarea
				name="message"
				class="w-full rounded border p-2 ring-blue-500 outline-none focus:ring-2"
				rows="3">{post.message}</textarea
			>

			<div class="mt-2 flex justify-end gap-2">
				<button type="button" on:click={toggleEdit} class="px-3 py-1 text-gray-600">Cancel</button>
				<button type="submit" class="rounded bg-blue-600 px-3 py-1 text-white shadow">Save</button>
			</div>
		</form>
	{:else}
		<p class="whitespace-pre-wrap text-gray-800">{post.message}</p>
	{/if}
</div>
