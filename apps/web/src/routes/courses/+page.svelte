<title>Courses</title>

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

<br />

<div class="ml-5 text-3xl">
	<h1 class="text-5xl font-bold">Courses</h1>
	<br />

	<DataLoader promise={coursesPromise}>
		{#snippet children(courses: Course[])}
			<ul>
				{#each courses as course}
					<a class="font-medium" href="/courses/{course.uuid}"> {course.name} </a>
					<p>{course.description}</p>
					<br />
				{/each}
			</ul>
		{/snippet}
	</DataLoader>
</div>
