<script lang="ts">
	import type { FeedPost } from '$lib/types';

	let { post, courseId }: { post: FeedPost; courseId: string } = $props();

	let message = $state(post.message);
	let initialMessage = $state(post.message);

	async function editFeedPost(e: Event) {
		e.preventDefault();

		const form = e.currentTarget as HTMLFormElement;
		const formData = new FormData(form);
		const data = Object.fromEntries(formData);

		const res = await fetch(`/api/courses/${courseId}/feed/${post.uuid}`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(data)
		});

		if (res.ok) {
			initialMessage = message;
		}
	}


	async function deleteFeedPost(e: Event) {
		e.preventDefault();

		const res = await fetch(`/api/courses/${courseId}/feed/${post.uuid}`, {
			method: 'DELETE',
			headers: { 'Content-Type': 'application/json' },
		});

	}
</script>

<form onsubmit={editFeedPost} class="rounded border">

	{#if post.type === 'manual'}
		<textarea
			name="message"
			class="mb-2 w-full p-2"
			bind:value={message}
			required>
		</textarea>

		{#if message !== initialMessage}
			<button class="rounded bg-blue-600 px-4 py-2 text-white">
				Edit
			</button>
		{/if}

		<button  type="button" onclick={deleteFeedPost} class="rounded bg-red-600 px-4 py-2 text-white">
			Delete
		</button>
	
	{:else}
		<p class="mb-2 w-full p-2"
>
		{message}

		</p>
	
	{/if}

	{#if post.createdAt != post.updatedAt}
		<p>Updated at: {post.updatedAt}</p>
	{/if}
	<p>Created at: {post.createdAt}</p>

</form>
