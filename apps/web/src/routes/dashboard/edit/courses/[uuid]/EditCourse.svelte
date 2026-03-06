<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Course } from '$lib/types';
	import { fade } from 'svelte/transition';
	import { modal } from '$lib/modal.svelte';

	import DangerButton from '$lib/components/DangerButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';

	let { course, onchange }: { course: Course; onchange: () => void } = $props();

	let showSuccess = $state(false);

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

		let formJson;
		if (isScheduled) {
			formJson = JSON.stringify({
				...formEntries,
				state: currentState,
				openTime: getISOTime()
			});
		} else {
			formJson = JSON.stringify({
				...formEntries,
				state: currentState
			});
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

	let isArchiving = $state(false);
	async function archiveCourse(e: Event) {
		e.preventDefault();

		isArchiving = true;

		const confirmed = await modal.confirm(`Archive entire course "${course.name}"?`);
		if (!confirmed) return;
		isUpdating = true;
		try {
			const res = await fetch(`/api/courses/${course.uuid}/archive`, { method: 'POST' });
			if (res.ok) goto('/dashboard');
		} finally {
			isArchiving = false;
		}
	}

	import StateController from '$lib/components/StateController.svelte';
	import WarningButton from '$lib/components/WarningButton.svelte';

	async function updateCourseStatus(newState: Course['state'], date?: string) {
		let body = { state: newState, openTime: date };

		await fetch(`/api/courses/${course.uuid}/state`, {
			method: 'PUT',
			body: JSON.stringify(body),

			headers: { 'Content-type': 'application/json' }
		});

		onchange();
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

		<br />

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
				<SuccessButton isSaving={isUpdating} type="submit">Save Changes</SuccessButton>

				<DangerButton isSaving={isDeleting} onclick={deleteCourse}>Delete Course</DangerButton>

				<WarningButton isSaving={isArchiving} onclick={archiveCourse}>Archive Course</WarningButton>
			</div>
		</form>
	</div>
</div>
