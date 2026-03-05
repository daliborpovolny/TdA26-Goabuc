<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import UniButton from '../UniButton.svelte';
	import UniLink from '../UniLink.svelte';

	if (!auth.user?.isAdmin) {
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
			<UniLink
				href="/dashboard/edit/courses"
				more_style="md:w-auto text-white hover:text-black"
				content="+Create New Course"
				bgcolor="bg-s-2"
			/>
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
						<div class="group relative transition hover:-translate-x-1 hover:-translate-y-1">
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
									<UniLink
										href="/dashboard/edit/courses/{course.uuid}"
										content="Edit Course"
										bgcolor="bg-[#444]"
										hv_bgcolor=""
										more_style="text-white hover:bg-p-blue"
									/>
									<UniLink href="/courses/{course.uuid}" content="👁️" hv_bgcolor='' more_style="hover:bg-p-blue"/>
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
					<UniButton onclick={() => location.reload()} content="Try Again" />
				</div>
			{/await}
		</main>
	</div>
</div>
