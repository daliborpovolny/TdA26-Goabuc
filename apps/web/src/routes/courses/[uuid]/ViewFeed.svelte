<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	import type { FeedPost } from '$lib/types';

	import ViewFeedItem from './ViewFeedItem.svelte';

	let { courseId }: { courseId: string } = $props();

	let collapsed = $state(true);

	let posts: FeedPost[] = $state([]);
	loadCourseFeed();

	async function loadCourseFeed() {
		let res = await fetch(`/api/courses/${courseId}/feed`, {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		});

		let data: FeedPost[] = await res.json();
		posts = data;
	}

	let eventSource: EventSource;

	onMount(() => {
		eventSource = new EventSource(`/api/courses/${courseId}/feed/stream`);

		// 1. Debugging: Check if connection opens
		eventSource.onopen = () => {
			console.log('SSE Connection Open');
		};

		// 2. ERROR: This only catches unnamed messages
		eventSource.onmessage = (e) => {
			console.log('Received unnamed message:', e.data);
		};

		// 3. FIX: Add a specific listener for "new_post"
		eventSource.addEventListener('new_post', (event) => {
			if (!event.data) return;

			console.log('Received new_post:', event.data); // Debug log

			try {
				let newPost: FeedPost = JSON.parse(event.data);
				console.log(newPost);

				posts.unshift(newPost);

				// postsStore.update((currentMap) => {
				//     // 1. Create a copy of the current map to ensure reactivity triggers
				//     const newMap = new Map(currentMap);

				//     // 2. Add or update the post
				//     newMap.set(newPost.uuid, newPost);

				//     // 3. Return the new map to update the store
				//     return newMap;
				// });
			} catch (err) {
				console.error('Error parsing SSE message:', err);
			}
		});

		eventSource.onerror = (err) => {
			console.error('SSE connection error:', err);
		};

		return () => {
			eventSource?.close();
		};
	});

	onDestroy(() => {
		if (eventSource) {
			eventSource.close();
		}
	});

	// // Convert map to array and sort by date descending for display
	// $: sortedPosts = Array.from($postsStore.values()).sort((a, b) =>
	//     new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
	// );
</script>

<div class="space-y-4 p-4">
	<button onclick={() => (collapsed = !collapsed)} type="button" class="text-3xl">News Feed</button>

	{#if !collapsed}
		<div class="flex flex-col gap-4">
			{#each posts as post (post.uuid)}
				<ViewFeedItem {post} />
			{/each}
		</div>
	{/if}
</div>
