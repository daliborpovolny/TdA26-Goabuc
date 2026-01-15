<script lang="ts">
	import { goto } from '$app/navigation';
	import { fade } from 'svelte/transition';

	let course_name = $state('');
	let course_description = $state('');
	let errorMsg = $state('');
	let isSubmitting = $state(false);

	async function handleCreateCourse(e: Event) {
		e.preventDefault();
		isSubmitting = true;
		errorMsg = '';

		try {
			const res = await fetch('/api/courses', {
				method: 'POST',
				headers: { 'Content-type': 'application/json' },
				body: JSON.stringify({
					name: course_name,
					description: course_description
				})
			});

			const data = await res.json();

			if (res.ok) {
				// Quick transition to the new course page
				goto('/dashboard/edit/courses/' + data.uuid);
			} else {
				errorMsg = data.message || 'Failed to create the course';
				isSubmitting = false;
			}
		} catch (e) {
			errorMsg = 'Network error. Please try again.';
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>New Course | TdA</title>
</svelte:head>

<div class="flex min-h-screen items-start justify-center bg-s-white p-4 md:p-12">
	<div class="relative w-full max-w-2xl">
		<div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>

		<div class="relative rounded-2xl border-4 border-s-black bg-white p-6 md:p-10">
			<header class="mb-8">
				<a href="/dashboard" class="text-sm font-black text-p-blue uppercase hover:underline"
					>← Back to Dashboard</a
				>
				<h1 class="mt-2 text-4xl font-black tracking-tighter uppercase md:text-5xl">
					Launch New <span class="text-p-blue">Course</span>
				</h1>
				<div class="mt-2 h-2 w-20 bg-p-green"></div>
			</header>

			<form onsubmit={handleCreateCourse} class="space-y-6">
				<div class="space-y-2">
					<label for="course_name" class="block text-xl font-black tracking-wide uppercase">
						Course Title
					</label>
					<input
						id="course_name"
						type="text"
						bind:value={course_name}
						required
						class="w-full rounded-xl border-4 border-s-black bg-white p-4 text-xl font-bold placeholder:text-gray-300 focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div class="space-y-2">
					<label for="course_description" class="block text-xl font-black tracking-wide uppercase">
						Description
					</label>
					<textarea
						id="course_description"
						bind:value={course_description}
						required
						rows="4"
						class="w-full rounded-xl border-4 border-s-black bg-white p-4 text-xl font-bold placeholder:text-gray-300 focus:ring-4 focus:ring-p-green focus:outline-none"
					></textarea>
				</div>

				<button
					type="submit"
					disabled={isSubmitting}
					class="group relative w-full overflow-hidden rounded-xl border-4 border-s-black bg-p-green py-4 text-2xl font-black tracking-widest uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
				>
					<div
						class="absolute inset-0 translate-x-1 translate-y-1 bg-s-black opacity-0 group-hover:opacity-10"
					></div>
					{isSubmitting ? 'Initializing...' : 'Create Course →'}
				</button>
			</form>

			{#if errorMsg}
				<div
					transition:fade
					class="mt-6 rounded-xl border-4 border-s-black bg-red-500 p-4 font-bold text-white"
				>
					⚠️ {errorMsg}
				</div>
			{/if}
		</div>
	</div>
</div>
