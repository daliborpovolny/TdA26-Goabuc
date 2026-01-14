<script lang="ts">
	import type { FeedPost } from '$lib/types';

	let { post }: { post: FeedPost } = $props();

	$inspect(post);
</script>

<div class="border p-1">
	<p>{post.message}</p>
	<br />
	<p>{post.type == 'auto' ? 'Course change alert' : 'Message from the teacher'}</p>

	{#if post.edited}
		<p>Edited at: {post.updatedAt}</p>
	{/if}
	<p>{post.createdAt}</p>
</div>
