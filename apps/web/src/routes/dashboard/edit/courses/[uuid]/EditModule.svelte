<script lang="ts">
	import type { Module } from '$lib/types';
	import { fade, slide } from 'svelte/transition';

	let { module, courseId, onchange }: { module: Module; courseId: string; onchange: () => void } =
		$props();

	let collapsed = $state(true);
	let isSaving = $state(false);
	let showSuccess = $state(false);

	// Initial state for editing
	let name = $state(module.name);
	let description = $state(module.description);

	async function handleUpdate(e: Event) {
		e.preventDefault();
		isSaving = true;

		try {
			const res = await fetch(`/api/courses/${courseId}/modules/${module.uuid}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name, description })
			});

			if (res.ok) {
				showSuccess = true;
				onchange();
				setTimeout(() => (showSuccess = false), 2000);
			}
		} finally {
			isSaving = false;
		}
	}

	async function deleteModule() {
		if (!confirm(`Delete module "${module.name}"? This will move items to unassigned!`)) return;
		await fetch(`/api/courses/${courseId}/modules/${module.uuid}`, { method: 'DELETE' });
		onchange();
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
>
	<button
		type="button"
		class="flex w-full items-center justify-between p-4 text-left hover:bg-p-green/10"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span
				class="flex h-10 w-10 items-center justify-center rounded-lg border-2 border-s-black bg-white text-xl"
			>
				📦
			</span>
			<div>
				<span class="block text-xl font-black tracking-tight text-s-black uppercase"
					>{module.name}</span
				>
				<span class="text-[10px] font-bold tracking-widest text-gray-400 uppercase">
					State: <span class="text-p-blue">{module.state}</span>
				</span>
			</div>
		</div>

		<div class="flex items-center gap-4">
			{#if showSuccess}
				<span transition:fade class="text-xs font-bold text-p-green uppercase">✓ Saved</span>
			{/if}
			<span class="text-xl transition-transform duration-300 {collapsed ? '' : 'rotate-180'}"
				>▼</span
			>
		</div>
	</button>

	{#if !collapsed}
		<div transition:slide class="border-t-4 border-s-black bg-gray-50 p-6">
			<form onsubmit={handleUpdate} class="space-y-4">
				<div class="space-y-1">
					<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="edit_name"
						>Module Name</label
					>
					<input
						id="edit_name"
						type="text"
						bind:value={name}
						required
						class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div class="space-y-1">
					<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="edit_desc"
						>Description</label
					>
					<textarea
						id="edit_desc"
						bind:value={description}
						rows="2"
						class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
					></textarea>
				</div>

				<div class="flex items-center justify-between border-t-2 border-gray-200 pt-4">
					<button
						type="button"
						onclick={deleteModule}
						class="rounded-lg border-2 border-s-black bg-red-500 px-4 py-2 text-xs font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
					>
						Delete Module
					</button>

					<button
						type="submit"
						disabled={isSaving}
						class="rounded-lg border-2 border-s-black bg-p-green px-6 py-2 text-xs font-black text-s-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none disabled:opacity-50"
					>
						{isSaving ? 'Saving...' : 'Save Changes'}
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>
