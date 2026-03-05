<script lang="ts">
	import { slide, fade } from 'svelte/transition';
	import SuccessButton from './SuccessButton.svelte';
	import type { Course } from '$lib/types';

	let {
		currentState = $bindable(),
		onchange,
		enableTimer = false
	}: {
		currentState: Course['state'];
		onchange: (state: Course['state'], date?: string) => void;
		enableTimer?: boolean;
	} = $props();

	const STAGES: Course['state'][] = ['preparation', 'open', 'closed'];

	let selectedState = $state(currentState);
	let releaseDate = $state('');
	let isSaving = $state(false);
	let showTimer = $state(false);

	async function handleUpdate() {
		isSaving = true;
		// If timer is hidden or empty, we don't send the date
		await onchange(selectedState, showTimer ? releaseDate : undefined);
		isSaving = false;
	}

	const stageMeta = {
		preparation: { label: 'In Preparation', color: 'bg-p-blue', icon: '🛠️' },
		open: { label: 'Open / Live', color: 'bg-p-green', icon: '🎓' },
		closed: { label: 'Temporarily Closed', color: 'bg-red-500', icon: '🔒' }
	};
</script>

<div
	class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
>
	<div class="mb-6 flex items-center justify-between border-b-2 border-gray-100 pb-4">
		<h3 class="text-xl font-black tracking-tight text-s-black uppercase">Course Visibility</h3>
		<span
			class="rounded-lg border-2 border-s-black bg-gray-100 px-2 py-1 text-[10px] font-black tracking-widest uppercase"
		>
			Current: {currentState}
		</span>
	</div>

	<div class="space-y-6">
		<div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
			{#each STAGES as stage}
				<button
					type="button"
					onclick={() => (selectedState = stage)}
					class="flex flex-col items-center justify-center rounded-xl border-4 border-s-black p-4 transition-all
                    {selectedState === stage
						? '-translate-y-1 bg-s-black text-white shadow-[4px_4px_0px_0px_rgba(145,245,173,1)] cursor-not-allowed'
						: 'bg-white text-s-black hover:bg-gray-50 cursor-pointer'}"
				>
					<span class="text-2xl">{stageMeta[stage].icon}</span>
					<span class="mt-1 text-xs font-black tracking-tighter uppercase"
						>{stageMeta[stage].label}</span
					>
				</button>
			{/each}
		</div>

		{#if enableTimer}
			<div class="rounded-xl border-2 border-dashed border-gray-300 p-4">
				<label class="mb-2 flex cursor-pointer items-center gap-2">
					<input
						type="checkbox"
						bind:checked={showTimer}
						class="h-5 w-5 rounded border-2 border-s-black text-p-blue focus:ring-p-green"
					/>
					<span class="text-sm font-bold text-s-black uppercase"
						>Optionally Schedule State Change</span
					>
				</label>

				{#if showTimer}
					<div transition:slide class="mt-4 space-y-2">
						<label
							class="text-[10px] font-black tracking-widest text-gray-400 uppercase"
							for="release_date">Apply change at:</label
						>
						<input
							id="release_date"
							type="datetime-local"
							bind:value={releaseDate}
							class="w-full rounded-lg border-4 border-s-black bg-white p-2 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>
				{/if}
			</div>
		{/if}

		<div class="pt-2">
			<SuccessButton
				{isSaving}
				disabled={selectedState === currentState && !showTimer}
				onclick={handleUpdate}
			>
				{showTimer ? 'Schedule Change' : 'Apply Now'}
			</SuccessButton>
		</div>
	</div>
</div>
