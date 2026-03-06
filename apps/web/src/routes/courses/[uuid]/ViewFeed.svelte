<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	import type { FeedPost } from '$lib/types';

	import ViewFeedItem from './ViewFeedItem.svelte';

	let { courseId, onUpdate }: { courseId: string; onUpdate?: () => void } = $props();

	let collapsed = $state(false);

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

				if (newPost.type != 'manual' && onUpdate) {
					console.log('system update detected!');
					onUpdate();
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

<div class="space-y-4 px-4">
	{#if posts.length > 0}
		{#if !collapsed}
			<div class="flex flex-col gap-4">
				{#each posts as post (post.uuid)}
					{#if post.type === 'manual' || post.type === 'system'}
						<ViewFeedItem {post} />
					{/if}
				{/each}
			</div>
		{/if}
	{/if}
</div>
