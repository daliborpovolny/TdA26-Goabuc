<script lang="ts">
	import { slide } from 'svelte/transition';

	import type { Module } from '$lib/types';

	let {
		modules,
		selectedId = $bindable(),
		label = 'Select Module'
	}: {
		modules: Module[];
		selectedId: string;
		label?: string;
	} = $props();

	let isOpen = $state(false);

	let selectedModuleName = $derived(
		modules.find((m) => m.uuid === selectedId)?.name || 'Choose a module...'
	);

	function selectModule(uuid: string) {
		selectedId = uuid;
		isOpen = false;
	}

	selectedId = modules[0].uuid; // this could be included only if we can always be certain that there's allways a module present
</script>

<div class="relative w-full max-w-xs">
	<span class="mb-1 block text-xs font-black tracking-widest text-gray-500 uppercase">{label}</span>

	<button
		type="button"
		onclick={() => (isOpen = !isOpen)}
		class="relative flex w-full items-center justify-between rounded-xl border-4 border-s-black bg-white p-3 font-bold transition-all hover:bg-gray-50 active:translate-y-1 active:shadow-none"
		style="box-shadow: {isOpen ? 'none' : '4px 4px 0px 0px rgba(26,26,26,1)'}"
	>
		<span class="truncate">{selectedModuleName}</span>
		<span class="transition-transform duration-200 {isOpen ? 'rotate-180' : ''}">▼</span>
	</button>

	{#if isOpen}
		<div
			transition:slide={{ duration: 200 }}
			class="absolute z-50 mt-2 w-full overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
		>
			<div class="max-h-60 overflow-y-auto">
				{#each modules as module}
					<button
						type="button"
						onclick={() => selectModule(module.uuid)}
						class="w-full border-b-2 border-gray-100 p-3 text-left font-bold transition-colors last:border-0 hover:bg-p-green"
						class:bg-p-blue={selectedId === module.uuid}
						class:text-white={selectedId === module.uuid}
					>
						{module.name}
					</button>
				{:else}
					<p class="p-3 text-sm italic text-gray-400">No modules available</p>
				{/each}
			</div>
		</div>
	{/if}
</div>

{#if isOpen}
	<button
		title="selected module"
		tabindex="-1"
		class="fixed inset-0 z-40 h-full w-full cursor-default bg-transparent outline-none"
		onclick={() => (isOpen = false)}
	></button>
{/if}
