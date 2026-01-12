<script lang="ts">
	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course } from '$lib/types';

	let coursesPromise: Promise<Course[]> = loadCourses();

	async function loadCourses() {
		return fetch('/api/courses', {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		}).then(async (res) => {
			if (!res.ok) {
				try {
					const err = await res.json();
				} catch {
					throw new Error('Failed to get list of courses');
				}
			}
			return res.json();
		});
	}
</script>

<title>Courses</title>

<br />

<div class="ml-5 text-3xl">
	<div class="w-fit rounded-xl bg-[#0070BB] px-5 pt-5 pb-3">
		<h1 class="text-5xl font-bold">Courses</h1>
		<br />

		<DataLoader promise={coursesPromise}>
			{#snippet children(courses: Course[])}
				{#each courses as course}
					<a href="/courses/{course.uuid}">
						<div class="mb-2 w-[100%] rounded-xl bg-[#91F5AD] p-5">
							<h2 class="text-medium mb-2 font-semibold underline decoration-3">{course.name}</h2>
							<p>{course.description}</p>
						</div>
					</a>
				{/each}
			{/snippet}
		</DataLoader>
	</div>
</div>
