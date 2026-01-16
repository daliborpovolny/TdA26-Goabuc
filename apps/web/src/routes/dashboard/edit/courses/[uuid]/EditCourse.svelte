<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Course } from '$lib/types';
	import { fade } from 'svelte/transition';

	let { course, onchange }: { course: Course; onchange: () => void } = $props();

	let isSaving = $state(false);
	let showSuccess = $state(false);

	async function updateCourse(e: Event) {
		e.preventDefault();
		isSaving = true;

		let formData = new FormData(e.target as HTMLFormElement);
		let formJson = JSON.stringify(Object.fromEntries(formData));

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
		isSaving = true;

		try {
			const res = await fetch(`/api/courses/${course.uuid}`, {
				method: 'DELETE'
			});

			if (res.ok) {
				goto('/dashboard');
				onchange();
				showSuccess = true;
				setTimeout(() => (showSuccess = false), 2000);
			}
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
			<div class="space-y-2">
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

			<div class="pt-2">
				<button
					type="submit"
					disabled={isSaving}
					class="group relative cursor-pointer overflow-hidden rounded-xl border-4 border-s-black bg-p-blue px-8 py-3 text-xl font-black tracking-widest text-white uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
				>
					{isSaving ? 'Syncing...' : 'Save Settings'}
				</button>

				<button
					type="button"
					onclick={deleteCourse}
					disabled={isSaving}
					class=" group relative ml-3 cursor-pointer overflow-hidden rounded-xl border-4 border-s-black bg-red-400 px-8 py-3 text-xl font-black tracking-widest text-white uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
				>
					{isSaving ? 'Deleting...' : 'Delete'}
				</button>
			</div>
		</form>
	</div>
</div>
