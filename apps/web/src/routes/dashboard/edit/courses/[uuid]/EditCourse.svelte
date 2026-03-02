<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Course } from '$lib/types';
	import { fade, slide } from 'svelte/transition';
	import { modal } from '$lib/modal.svelte';

	let { course, onchange }: { course: Course; onchange: () => void } = $props();

	const STAGES = ['preparation', 'open', 'closed'] as const;

	let isSaving = $state(false);
	let showSuccess = $state(false);
	let showStateDropdown = $state(false);
	let currentState = $state(course.state);

	async function updateCourse(e: Event) {
		e.preventDefault();
		isSaving = true;

		let formData = new FormData(e.target as HTMLFormElement);
		let formEntries = Object.fromEntries(formData);

		// Ensure the current state is included in the payload
		let formJson = JSON.stringify({
			...formEntries,
			state: currentState
		});

		try {
			const res = await fetch(`/api/courses/${course.uuid}`, {
				method: 'PUT',
				headers: { 'Content-type': 'application/json' },
				body: formJson
			});

			if (res.ok) {
				onchange();
				showSuccess = true;
				setTimeout(() => (showSuccess = false), 2000);
			}
		} finally {
			isSaving = false;
		}
	}

	async function deleteCourse(e: Event) {
		e.preventDefault();
		const confirmed = await modal.confirm(
			`Delete entire course "${course.name}"? This action cannot be undone.`
		);
		if (!confirmed) return;

		isSaving = true;
		try {
			const res = await fetch(`/api/courses/${course.uuid}`, {
				method: 'DELETE'
			});
			if (res.ok) goto('/dashboard');
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="space-y-6">
	<header class="flex items-center justify-between">
		<h2 class="text-3xl font-black tracking-tight text-s-black uppercase italic">General Info</h2>
		{#if showSuccess}
			<span
				transition:fade
				class="rounded-lg border-2 border-s-black bg-p-green px-3 py-1 text-sm font-bold uppercase"
			>
				✓ Changes Saved
			</span>
		{/if}
	</header>

	<div
		class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] md:p-8"
	>
		<form onsubmit={updateCourse} class="space-y-6">
			<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
				<div class="space-y-2 md:col-span-2">
					<label class="block text-lg font-black tracking-wide text-s-black uppercase" for="name">
						Course Name
					</label>
					<input
						type="text"
						id="name"
						name="name"
						value={course.name}
						required
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div class="relative space-y-2 md:col-span-1">
					<label class="block text-lg font-black tracking-wide text-s-black uppercase">
						Course State
					</label>
					<div class="relative">
						<button
							type="button"
							onclick={() => (showStateDropdown = !showStateDropdown)}
							class="flex w-full items-center justify-between rounded-xl border-4 border-s-black bg-white p-3 font-bold tracking-tight uppercase transition-all active:translate-y-1"
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
							<span>{showStateDropdown ? '▲' : '▼'}</span>
						</button>

						{#if showStateDropdown}
							<div
								transition:slide={{ duration: 150 }}
								class="absolute top-[calc(100%+8px)] right-0 left-0 z-50 overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
							>
								{#each STAGES as stage}
									<button
										type="button"
										class="w-full p-4 text-left font-black tracking-widest uppercase transition-colors hover:bg-p-green"
										class:bg-p-blue={currentState === stage}
										class:text-white={currentState === stage}
										onclick={() => {
											currentState = stage;
											showStateDropdown = false;
										}}
									>
										{stage}
									</button>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>

			<div class="space-y-2">
				<label
					class="block text-lg font-black tracking-wide text-s-black uppercase"
					for="description"
				>
					Detailed Description
				</label>
				<textarea
					id="description"
					name="description"
					value={course.description}
					required
					rows="4"
					class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
				></textarea>
			</div>

			<div class="flex flex-wrap gap-4 pt-4">
				<button
					type="submit"
					disabled={isSaving}
					class="group relative flex-1 cursor-pointer overflow-hidden rounded-xl border-4 border-s-black bg-p-blue px-8 py-3 text-xl font-black tracking-widest text-white uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
				>
					{isSaving ? 'Syncing...' : 'Save Settings'}
				</button>

				<button
					type="button"
					onclick={deleteCourse}
					disabled={isSaving}
					class="group relative cursor-pointer overflow-hidden rounded-xl border-4 border-s-black bg-red-400 px-8 py-3 text-xl font-black tracking-widest text-white uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
				>
					{isSaving ? 'Deleting...' : 'Delete'}
				</button>
			</div>
		</form>
	</div>
</div>

{#if showStateDropdown}
	<button
		tabindex="-1"
		class="fixed inset-0 z-40 h-full w-full cursor-default bg-transparent outline-none"
		onclick={() => (showStateDropdown = false)}
	></button>
{/if}
