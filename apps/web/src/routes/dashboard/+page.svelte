<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';

	if (!auth.isLoggedIn) {
		goto('/login');
	}

	let coursesPromise = $state(loadCourses());

	async function loadCourses() {
		const res = await fetch('/api/courses');
		if (!res.ok) throw new Error('Failed to load courses');
		return res.json();
	}
</script>

<svelte:head>
	<title>Dashboard | TdA</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-4 md:p-8">
	<div class="mx-auto max-w-6xl">
		<header
			class="mb-10 flex flex-col items-center justify-between gap-6 border-b-4 border-s-black pb-8 md:flex-row"
		>
			<div>
				<h1 class="text-5xl font-black tracking-tighter text-s-black uppercase md:text-6xl">
					Dashboard
				</h1>
				<p class="font-bold tracking-widest text-p-blue uppercase">Manage your curriculum</p>
			</div>

			<a href="/dashboard/edit/courses" class="group relative inline-block w-full md:w-auto">
				<div class="absolute inset-0 translate-x-1 translate-y-1 rounded-xl bg-s-black"></div>
				<div
					class="relative rounded-xl border-2 border-s-black bg-s-2 px-6 py-3 text-center text-xl font-black text-white transition-all group-hover:bg-s-1 group-hover:text-s-black group-active:translate-x-1 group-active:translate-y-1"
				>
					+ Create New Course
				</div>
			</a>
		</header>

		<main>
			{#await coursesPromise}
				<div class="flex h-64 items-center justify-center">
					<p class="animate-bounce text-2xl font-black text-p-blue uppercase">
						Loading Academy Data...
					</p>
				</div>
			{:then courses}
				<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
					{#each courses as course}
						<div class="group relative">
							<div
								class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"
							></div>

							<div
								class="relative flex flex-col border-4 border-s-black bg-white p-6 transition-all group-hover:bg-p-green"
							>
								<div class="mb-4 flex justify-between">
									<span
										class="rounded-md border-2 border-s-black bg-p-blue px-2 py-1 text-xs font-bold text-white uppercase"
									>
										ID: {course.uuid.slice(0, 8)}
									</span>
									<span class="text-2xl">⚙️</span>
								</div>

								<h2 class="text-2xl leading-tight font-black text-s-black">
									{course.name}
								</h2>

								<div class="mt-8 flex gap-2">
									<a
										href="/dashboard/edit/courses/{course.uuid}"
										class="flex-1 rounded-lg border-2 border-s-black bg-s-black py-2 text-center font-bold text-white transition-colors hover:bg-p-blue"
									>
										Edit Course
									</a>
									<a
										href="/courses/{course.uuid}"
										class="rounded-lg border-2 border-s-black bg-white px-3 py-2 transition-colors hover:bg-p-blue"
										title="View Public Page"
									>
										👁️
									</a>
								</div>
							</div>
						</div>
					{/each}
				</div>

				{#if courses.length === 0}
					<div class="rounded-2xl border-4 border-dashed border-gray-300 p-12 text-center">
						<p class="text-xl font-bold text-gray-400">
							No courses found. Time to create your first one!
						</p>
					</div>
				{/if}
			{:catch error}
				<div class="rounded-xl border-4 border-s-black bg-red-500 p-6 text-white">
					<p class="text-xl font-bold tracking-tight uppercase">Error Loading Dashboard</p>
					<p>{error.message}</p>
					<button onclick={() => location.reload()} class="mt-4 underline">Try Again</button>
				</div>
			{/await}
		</main>
	</div>
</div>
