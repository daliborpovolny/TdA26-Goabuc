<script lang="ts">
	import { slide } from 'svelte/transition';

	import type { Module } from '$lib/types';
	import UniButton from '../../../../UniButton.svelte';

	let {
		modules,
		selectedId = $bindable(),
		label = 'Select Module',
		onChange // Define it in the props
	}: {
		modules: Module[];
		selectedId: string;
		label?: string;
		onChange?: (id: string) => void; // Type it to accept the ID
	} = $props();

	let isOpen = $state(false);

	let selectedModuleName = $derived(
		modules.find((m) => m.uuid === selectedId)?.name || 'Choose a module...'
	);

	function selectModule(uuid: string) {
		selectedId = uuid;
		isOpen = false;
		// Call it immediately after the state change
		if (onChange) onChange(uuid);
	}

	if (selectedId == '') {
		selectedId = modules[0].uuid; // this could be included only if we can always be certain that there's allways a module present
	}
</script>

<div class="relative w-full max-w-xs">
	<span class="mb-1 block text-xs font-black tracking-widest text-gray-500 uppercase">{label}</span>

	<UniButton
		type="button"
		onclick={() => (isOpen = !isOpen)}
		more_style="flex justify-between tracking-tight w-full flex items-center"
		text="text-l"
		uppercase
		px="px-4"
	>
		<span class="truncate">{selectedModuleName}</span>
		<span class="transition-transform duration-200 {isOpen ? 'rotate-180' : ''}">▼</span>
	</UniButton>

	{#if isOpen}
		<div
			transition:slide={{ duration: 200 }}
			class="absolute z-50 mt-2 w-full overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
		>
			<div class="max-h-60 overflow-y-auto">
				{#each modules as module}
					<UniButton
						type="button"
						onclick={() => selectModule(module.uuid)}
						more_style={'w-full'}
						px="px-3"
						py="py-2"
						bgcolor={selectedId === module.uuid ? 'bg-p-blue' : undefined}
						text_color={selectedId === module.uuid ? 'text-white' : undefined}
						uppercase
						text="text-l"
						border={false}
						shadow={false}
						translate={false}
					>
						{module.name}
					</UniButton>
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
