<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	import type { FeedPost } from '$lib/types';
	import EditFeedPost from './EditFeedPost.svelte';

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

			try {
				let newPost: FeedPost = JSON.parse(event.data);

				const index = posts.findIndex((p) => p.uuid === newPost.uuid);

				if (index !== -1) {
					posts[index] = newPost;
					posts = [...posts];
				} else {
					posts = [newPost, ...posts];
				}
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
</script>

<div class="space-y-4 p-4">
	{#if posts.length > 0}
		<button onclick={() => (collapsed = !collapsed)} type="button" class="text-3xl"
			>News Feed</button
		>

		{#if !collapsed}
			<div class="flex flex-col gap-4">
				{#each posts as post (post.uuid)}
					<EditFeedPost {courseId} {post} />
				{/each}
			</div>
		{/if}
	{/if}
</div>
