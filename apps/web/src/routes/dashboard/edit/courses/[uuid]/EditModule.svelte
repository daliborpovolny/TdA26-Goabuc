<script lang="ts">
	import type { Module } from '$lib/types';
	import { fade, slide } from 'svelte/transition';
	import { modal } from '$lib/modal.svelte';
	import UniButton from '../../../../UniButton.svelte';
	import UniLink from '../../../../UniLink.svelte';

	let { module, courseId, onchange }: { module: Module; courseId: string; onchange: () => void } =
		$props();

	const STAGES = ['preparation', 'open', 'closed'] as const;

	let collapsed = $state(true);
	let isSaving = $state(false);
	let showSuccess = $state(false);
	let showStateDropdown = $state(false);

	// Initial state for editing
	let name = $state(module.name);
	let description = $state(module.description);
	let currentState = $state(module.state);

	async function handleUpdate(e: Event) {
		e.preventDefault();
		isSaving = true;

		try {
			const res = await fetch(`/api/courses/${courseId}/modules/${module.uuid}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ name, description, state: currentState })
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
		const confirmed = await modal.confirm(
			`Delete module "${module.name}"? This action cannot be undone.`
		);
		if (!confirmed) return;

		await fetch(`/api/courses/${courseId}/modules/${module.uuid}`, { method: 'DELETE' });
		onchange();
	}

	// Helper for state colors
	const stateColors = {
		preparation: 'text-p-blue bg-p-blue/10 border-p-blue/20',
		open: 'text-p-green bg-p-green/10 border-p-green/20',
		closed: 'text-red-500 bg-red-50 border-red-200'
	};
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
>
	<button
		type="button"
		class="flex w-full items-center justify-between p-4 text-left hover:bg-p-green/5"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span
				class="flex h-10 w-10 items-center justify-center rounded-lg border-2 border-s-black bg-white text-xl"
			>
				📦
			</span>
			<div>
				<span class="block text-xl font-black tracking-tight text-s-black uppercase">
					{module.name}
				</span>
				<span
					class="inline-block rounded border px-1.5 py-0.5 text-[10px] font-black tracking-widest uppercase {stateColors[
						module.state
					]}"
				>
					{module.state}
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
			<form onsubmit={handleUpdate} class="space-y-6">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<div class="space-y-1 md:col-span-2">
						<label
							class="text-xs font-black tracking-widest text-gray-500 uppercase"
							for="edit_name">Module Name</label
						>
						<input
							id="edit_name"
							type="text"
							bind:value={name}
							required
							class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>

					<div class="relative space-y-1 md:col-span-1">
						<label class="text-xs font-black tracking-widest text-gray-500 uppercase"
							>Module State</label
						>

						<div class="relative">
							<UniButton
								type="button"
								onclick={() => (showStateDropdown = !showStateDropdown)}
								more_style="tracking-tight w-full flex items-center"
								text="text-l"
								uppercase
								px="px-4"
							>
								<span class="flex items-center gap-2">
									<span
										class="h-3 w-3 rounded-full {currentState === 'open'
											? 'bg-p-green'
											: currentState === 'preparation'
												? 'bg-p-blue'
												: 'bg-red-500'}"
									></span>
									{currentState}
								</span>
								<span class="w-full text-right">{showStateDropdown ? '▲' : '▼'}</span>
							</UniButton>

							{#if showStateDropdown}
								<div
									transition:slide={{ duration: 150 }}
									class="absolute top-[calc(100%+8px)] right-0 left-0 z-50 overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
								>
									{#each STAGES as stage}
										<UniButton
											content={stage}
											border={false}
											shadow={false}
											more_style="w-full"
											px="px-0"
											text="text-l"
											uppercase
											onclick={() => {
												currentState = stage;
												showStateDropdown = false;
											}}
											translate={false}
										/>
									{/each}
								</div>
							{/if}
						</div>
					</div>
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

				<div
					class="border-t-2 border-gray-200 pt-4 max-md:space-y-5 md:flex md:items-center md:justify-between"
				>
					<UniButton
						type="button"
						onclick={deleteModule}
						content="Delete Module"
						bgcolor="bg-red-200"
						hv_bgcolor="bg-red-400"
					/>
					<UniButton
						type="submit"
						disabled={isSaving}
						content={isSaving ? 'Saving...' : 'Save Changes'}
						bgcolor="bg-p-green"
						hv_bgcolor="bg-green-500"
					/>
				</div>
			</form>
		</div>
	{/if}
</div>

{#if showStateDropdown}
	<button
		tabindex="-1"
		class="fixed inset-0 z-0 h-full w-full cursor-default bg-transparent outline-none"
		onclick={() => (showStateDropdown = false)}
	></button>
{/if}
