<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { fade, slide } from 'svelte/transition';
	import type { FeedPost } from '$lib/types';
	import EditFeedPost from './EditFeedPost.svelte';

	let { courseId }: { courseId: string } = $props();
	let collapsed = $state(false);
	let posts: FeedPost[] = $state([]);
	let isConnected = $state(false);

	loadCourseFeed();

	async function loadCourseFeed() {
		let res = await fetch(`/api/courses/${courseId}/feed`);
		posts = await res.json();
	}

	let eventSource: EventSource;

	onMount(() => {
		eventSource = new EventSource(`/api/courses/${courseId}/feed/stream`);

		eventSource.onopen = () => (isConnected = true);
		eventSource.onerror = () => (isConnected = false);

		eventSource.addEventListener('new_post', (event) => {
			if (!event.data) return;
			try {
				let newPost: FeedPost = JSON.parse(event.data);
				const index = posts.findIndex((p) => p.uuid === newPost.uuid);
				if (index !== -1) {
					posts[index] = newPost;
				} else {
					posts = [newPost, ...posts];
				}
			} catch (err) {
				console.error('Error parsing SSE message:', err);
			}
		});

		return () => eventSource?.close();
	});
</script>

<div class="space-y-4">
	{#if posts.length > 0}
		<button
			onclick={() => (collapsed = !collapsed)}
			type="button"
			class="flex cursor-pointer items-center gap-3 text-3xl font-black tracking-tighter uppercase hover:text-p-blue"
		>
			<span>News Feed</span>
			<span class="text-sm transition-transform {collapsed ? '' : 'rotate-180'}">▼</span>

			{#if isConnected}
				<span class="flex items-center gap-1 text-xs font-bold text-p-green">
					<span class="h-2 w-2 animate-pulse rounded-full bg-p-green"></span>
					Live
				</span>
			{/if}
		</button>

		{#if !collapsed}
			<div transition:slide class="flex flex-col gap-6 pt-4">
				{#each posts as post (post.uuid)}
					<EditFeedPost {courseId} {post} />
				{/each}
			</div>
		{/if}
	{/if}
</div>
