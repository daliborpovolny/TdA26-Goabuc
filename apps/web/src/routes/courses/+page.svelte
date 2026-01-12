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

<div class="ml-5 text-3xl text-center">
	<h1 class="text-5xl font-bold">Courses</h1>
	<br />

	<DataLoader promise={coursesPromise}>
		{#snippet children(courses: Course[])}
			<div>
				{#each courses as course}
					<div class="bg-[#91F5AD] p-5 m-2 w-fit rounded-xl">
						<a class="font-medium underline" href="/courses/{course.uuid}"> {course.name} </a>
						<p>{course.description}</p>
					</div>
				{/each}
			</div>
		{/snippet}
	</DataLoader>
</div>
