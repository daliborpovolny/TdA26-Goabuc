<script lang="ts">
	import type { Course } from '$lib/types';
	import { fade, fly } from 'svelte/transition';

	let searchTerm = $state('');
	let allCourses = $state<Course[]>([]);
	let loading = $state(true);

	// Initial load
	const coursesPromise = loadCourses().then((data) => {
		allCourses = data;
		loading = false;
		return data;
	});

	async function loadCourses() {
		const res = await fetch('/api/courses');
		if (!res.ok) throw new Error('Failed to get list of courses');
		return res.json();
	}

	// Derived state for filtering
	let filteredCourses = $derived(
		allCourses.filter((course) => course.name.toLowerCase().includes(searchTerm.toLowerCase()))
	);

	const stateStyles = {
		open: {
			bg: 'bg-white group-hover:bg-p-green',
			badge: 'bg-p-blue text-white',
			label: 'OPEN',
			clickable: true
		},
		preparation: {
			bg: 'bg-yellow-300',
			badge: 'bg-s-black text-white',
			label: 'COMING SOON',
			clickable: false
		},
		closed: {
			bg: 'bg-red-400',
			badge: 'bg-s-black text-white',
			label: 'CLOSED',
			clickable: false
		}
	};
</script>

<svelte:head>
	<title>Courses | TdA</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-6 md:p-12">
	<header class="mb-12 text-center">
		<h1 class="inline-block text-5xl font-black tracking-tighter uppercase md:text-7xl">
			Our <span class="text-p-blue">Courses</span>
		</h1>
		<div class="mx-auto mt-2 h-2 w-24 bg-p-green"></div>

		<div class="mx-auto mt-12 max-w-xl">
			<div class="relative">
				<div class="absolute inset-0 translate-x-1.5 translate-y-1.5 rounded-xl bg-s-black"></div>
				<input
					type="text"
					bind:value={searchTerm}
					placeholder="Search by course name..."
					class="relative w-full rounded-xl border-4 border-s-black bg-white p-4 text-xl font-bold placeholder:text-gray-400 focus:ring-4 focus:ring-p-green focus:outline-none"
				/>
				<span class="absolute top-1/2 right-4 -translate-y-1/2 text-2xl">🔍</span>
			</div>
			{#if searchTerm}
				<p class="mt-4 font-bold tracking-widest text-p-blue uppercase" in:fade>
					Found {filteredCourses.length} matching courses
				</p>
			{/if}
		</div>
	</header>

	<div class="mx-auto max-w-6xl">
		{#await coursesPromise}
			<div class="flex h-64 items-center justify-center">
				<p class="animate-bounce text-2xl font-black text-p-blue uppercase">
					Fetching Curriculum...
				</p>
			</div>
		{:then _}
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
				{#each filteredCourses as course (course.uuid)}
					{@const style = stateStyles[course.state] || stateStyles.closed}

					<div in:fly={{ y: 20, duration: 400 }}>
						<svelte:element
							this={style.clickable ? 'a' : 'div'}
							href={style.clickable ? `/courses/${course.uuid}` : undefined}
							class="group relative block transition-transform {style.clickable
								? 'hover:-translate-x-1 hover:-translate-y-1'
								: 'cursor-not-allowed opacity-90'}"
						>
							<div
								class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"
							></div>

							<div
								class="relative flex h-full flex-col border-2 border-s-black p-6 transition-colors {style.bg}"
							>
								<div class="mb-4 flex items-start justify-between">
									<span
										class="rounded-lg border-2 border-s-black px-3 py-1 text-xs font-black uppercase {style.badge}"
									>
										{style.label}
									</span>
									<span class="text-3xl">
										{course.state === 'open' ? '🎓' : course.state === 'preparation' ? '🛠️' : '🔒'}
									</span>
								</div>

								<h2 class="mb-4 text-3xl leading-tight font-black text-s-black uppercase">
									{course.name}
								</h2>

								<div class="mt-auto border-t-2 border-s-black/10 pt-6">
									<div
										class="flex items-center gap-2 font-black tracking-widest text-s-black uppercase"
									>
										{#if style.clickable}
											Enter Course
											<span class="transition-transform group-hover:translate-x-2">→</span>
										{:else}
											Access Locked
										{/if}
									</div>
								</div>
							</div>
						</svelte:element>
					</div>
				{/each}
			</div>
		{/await}
	</div>
</div>
