<script lang="ts">
	import SuccessButton from '$lib/components/SuccessButton.svelte';
	import { fade } from 'svelte/transition';

	let { courseId, onchange }: { courseId: string; onchange: () => void } = $props();

	let name = $state('');
	let description = $state('');
	let isSaving = $state(false);
	let errorMsg = $state('');

	let forbidden = 'Unassigned';

	async function handleCreate(e: Event) {
		if (name === forbidden) {
			return;
		}

		e.preventDefault();
		isSaving = true;
		errorMsg = '';

		try {
			const res = await fetch(`/api/courses/${courseId}/modules`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name, description })
			});

			if (res.ok) {
				onchange();
				name = '';
				description = '';
			} else {
				const data = await res.json();
				errorMsg = data.message || 'Failed to create module';
			}
		} catch (e) {
			errorMsg = 'Connection error';
		} finally {
			isSaving = false;
		}
	}
</script>

<form onsubmit={handleCreate} class="space-y-4">
	<div class="grid grid-cols-1 gap-4">
		<div class="space-y-1">
			<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="mod_name"
				>Module Name</label
			>
			<input
				id="mod_name"
				type="text"
				bind:value={name}
				required
				placeholder="e.g., Week 1: Introduction"
				class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
			/>

			{#if name === forbidden}
				<p class="text-red-500">The name {forbidden} is not allowed</p>
			{/if}
		</div>

		<div class="space-y-1">
			<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="mod_desc"
				>Short Description</label
			>
			<textarea
				id="mod_desc"
				bind:value={description}
				rows="2"
				placeholder="What is the objective of this module?"
				class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
			></textarea>
		</div>
	</div>

	<div class="flex items-center justify-between">
		{#if errorMsg}
			<span transition:fade class="text-xs font-bold text-red-500 uppercase">⚠️ {errorMsg}</span>
		{:else}
			<div></div>
		{/if}
		<SuccessButton {isSaving} type="submit">Create Module</SuccessButton>
	</div>
</form>
