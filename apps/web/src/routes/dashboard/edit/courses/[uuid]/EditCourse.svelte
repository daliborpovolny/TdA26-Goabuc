<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Course } from '$lib/types';
	import { fade, slide } from 'svelte/transition';
	import { modal } from '$lib/modal.svelte';
	import UniButton from '../../../../UniButton.svelte';

	import DangerButton from '$lib/components/DangerButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';

	let { course, onchange }: { course: Course; onchange: () => void } = $props();

	const STAGES = ['preparation', 'open', 'closed'] as const;

	let showSuccess = $state(false);
	let showStateDropdown = $state(false);
	// svelte-ignore state_referenced_locally
	let currentState = $state(course.state);

	let isScheduled = $state(false);
	let scheduledTime = $state('');

	function getISOTime() {
		if (!isScheduled || !scheduledTime) return null;
		try {
			return new Date(scheduledTime).toISOString();
		} catch (e) {
			return null;
		}
	}

	let isUpdating = $state(false);
	async function updateCourse(e: Event) {
		e.preventDefault();
		isUpdating = true;

		let formData = new FormData(e.target as HTMLFormElement);
		let formEntries = Object.fromEntries(formData);

		let formJson
		if (isScheduled) {

			formJson = JSON.stringify({
			...formEntries,
			state: currentState,
			openTime: getISOTime() // This maps to your Go struct
		})
		} else {
			formJson = JSON.stringify({
			...formEntries,
			state: currentState,
		})	
		}

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
				// Reset scheduling after success
				isScheduled = false;
				scheduledTime = '';
			}
		} finally {
			isUpdating = false;
		}
	}

	let isDeleting = $state(false);
	async function deleteCourse(e: Event) {
		e.preventDefault();

		isDeleting = true;

		const confirmed = await modal.confirm(`Delete entire course "${course.name}"?`);
		if (!confirmed) return;
		isUpdating = true;
		try {
			const res = await fetch(`/api/courses/${course.uuid}`, { method: 'DELETE' });
			if (res.ok) goto('/dashboard');
		} finally {
			isDeleting = false;
		}
	}


	import StateController from '$lib/components/StateController.svelte';
    
    async function updateCourseStatus(newState: Course["state"], date?: string) {

		let body = {state: newState, openTime: date}

		await fetch(`/api/courses/${course.uuid}/state`, {
			method: 'PUT',
			body: JSON.stringify(body),

			headers: { 'Content-type': 'application/json' },
		})
    
		onchange()
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
	<StateController 
		bind:currentState={course.state} 
		onchange={updateCourseStatus} 
		enableTimer={true} 
	/>

	<br>

		<form onsubmit={updateCourse} class="space-y-6">
			<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
				<div class="space-y-2 md:col-span-2">
					<label class="block text-lg font-black tracking-wide text-s-black uppercase" for="name"
						>Course Name</label
					>
					<input
						type="text"
						id="name"
						name="name"
						value={course.name}
						required
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<!-- <div class="relative space-y-2 md:col-span-1">
					<label class="block text-lg font-black tracking-wide text-s-black uppercase"
						>Course State</label
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
								{currentState == 'closed' ? 'paused' : currentState}
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
										content={stage == 'closed' ? 'paused' : stage}
										border={false}
										shadow={false}
										more_style="w-full"
										px="px-0"
										text="text-l"
										uppercase
										translate={false}
										onclick={() => {
											currentState = stage;
											showStateDropdown = false;
										}}
									/>
								{/each}
							</div>
						{/if}
					</div>

					<div class="mt-4 space-y-3">
						<label class="flex cursor-pointer items-center gap-3">
							<input
								type="checkbox"
								bind:checked={isScheduled}
								class="h-6 w-6 cursor-pointer appearance-none border-4 border-s-black bg-white checked:bg-p-green"
							/>
							<span class="text-sm font-black tracking-wide uppercase">Schedule Change?</span>
						</label>

						{#if isScheduled}
							<div transition:slide={{ duration: 150 }} class="space-y-1">
								<label class="block text-xs font-black text-s-black/60 uppercase"
									>Execution Time</label
								>
								<input
									type="datetime-local"
									bind:value={scheduledTime}
									required={isScheduled}
									class="w-full rounded-xl border-4 border-s-black bg-white p-2 font-bold shadow-[3px_3px_0px_0px_rgba(26,26,26,1)] focus:outline-none"
								/>
							</div>
						{/if}
					</div>
				</div> -->
			</div>

			<div class="space-y-2">
				<label
					class="block text-lg font-black tracking-wide text-s-black uppercase"
					for="description">Detailed Description</label
				>
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
				<!-- <UniButton
					type="submit"
					disabled={isSaving}
					content={isSaving ? 'Syncing...' : 'Save Settings'}
					bgcolor="bg-p-green"
					hv_bgcolor="bg-green-400"
				/> -->
				<!-- <UniButton
					type="button"
					onclick={deleteCourse}
					disabled={isSaving}
					content={isSaving ? 'Deleting...' : 'Delete'}
					bgcolor="bg-red-400"
					hv_bgcolor="bg-red-500"
				/> -->
				<SuccessButton isSaving={isUpdating} type="submit">Save Changes</SuccessButton>

				<DangerButton isSaving={isUpdating} onclick={deleteCourse}>Delete Course</DangerButton>
			</div>
		</form>
	</div>
</div>

<!-- {#if showStateDropdown}
	<button
		title="show dropdown"
		tabindex="-1"
		class="fixed inset-0 z-40 h-full w-full cursor-default bg-transparent outline-none"
		onclick={() => (showStateDropdown = false)}
	></button>
{/if} -->
