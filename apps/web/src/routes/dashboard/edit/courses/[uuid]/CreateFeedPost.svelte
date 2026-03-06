<script lang="ts">
	import { fade, slide } from 'svelte/transition';
	import UniButton from '../../../../UniButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';

	let { courseId }: { courseId: string } = $props();

	let collapsed = $state(true);
	let isSaving = $state(false);
	let showSuccess = $state(false);

	async function newFeedPost(e: Event) {
		e.preventDefault();
		isSaving = true;

		const form = e.currentTarget as HTMLFormElement;
		const formData = new FormData(form);
		const data = Object.fromEntries(formData);

		let res = await fetch(`/api/courses/${courseId}/feed`, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify(data)
		});

		if (res.ok) {
			showSuccess = true;
			setTimeout(() => {
				collapsed = true;
				showSuccess = false;
				form.reset();
			}, 1000);
		}
		isSaving = false;
	}
</script>

<div class="mb-6">
	<div class="flex items-center gap-4">
		<UniButton
			type="button"
			more_style="flex items-center gap-2 tracking-widest {collapsed
				? 'shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]'
				: 'translate-x-1 translate-y-1 bg-p-blue text-white shadow-none'}"
			shadow={collapsed}
			uppercase
			text="text-l"
			bgcolor={collapsed ? undefined : 'bg-p-blue'}
			hv_bgcolor={collapsed ? undefined : ''}
			onclick={() => (collapsed = !collapsed)}
		>
			<span>📢</span>
			{collapsed ? 'New Announcement' : 'Cancel Post'}
		</UniButton>

		{#if showSuccess}
			<span transition:fade class="font-bold tracking-tighter text-p-green uppercase">
				✓ Post Broadcasted
			</span>
		{/if}
	</div>

	{#if !collapsed}
		<div transition:slide class="relative mt-6 max-w-2xl">
			<div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>

			<div class="relative rounded-2xl border-4 border-s-black bg-white p-6">
				<form onsubmit={newFeedPost} class="space-y-4">
					<div class="space-y-1">
						<label for="message" class="text-xs font-black tracking-widest text-gray-500 uppercase">
							Message to Students
						</label>
						<textarea
							id="message"
							name="message"
							class="min-h-[120px] w-full resize-none rounded-xl border-4 border-s-black bg-white p-4 text-xl font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
							placeholder="Type your update here..."
							required
						></textarea>
					</div>

					<div class="flex justify-end pt-2">
						<SuccessButton type="submit" disabled={isSaving}>
							{isSaving ? 'Sending...' : 'Post Update →'}
						</SuccessButton>
					</div>
				</form>
			</div>
		</div>
	{/if}
</div>
