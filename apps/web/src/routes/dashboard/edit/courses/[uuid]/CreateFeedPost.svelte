<script lang="ts">
	let { courseId }: { courseId: string } = $props();

	let collapsed = $state(true);

	async function newFeedPost(e: Event) {
		e.preventDefault();

		const form = e.currentTarget as HTMLFormElement;
		const formData = new FormData(form);
		const data = Object.fromEntries(formData);

		let res = await fetch(`/api/courses/${courseId}/feed`, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify(data)
		});
	}
</script>

<button
	type="button"
	class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200"
	onclick={() => (collapsed = !collapsed)}
>
	New Announcement
</button>

{#if !collapsed}
	<br />
	<br />

	<div class="rounded bg-white p-4 shadow">
		<!-- <h3 class="mb-2 font-bold">New Announcement</h3> -->
		<form onsubmit={newFeedPost}>
			<textarea
				name="message"
				class="mb-2 w-full rounded border p-2"
				placeholder="What's new?"
				required
			></textarea>
			<button class="rounded bg-blue-600 px-4 py-2 text-white">Post</button>
		</form>
	</div>
{/if}
