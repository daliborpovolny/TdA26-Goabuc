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

<div class="p-6">
	<h1 class="text-center text-5xl font-bold underline">Courses</h1>
	<div class="flex w-[100%] justify-center text-center text-3xl">
		<div class="max-w-[50%] space-y-2 rounded-xl bg-[#0070BB] px-3 pt-3 pb-1">
			<DataLoader promise={coursesPromise}>
				{#snippet children(courses: Course[])}
					{#each courses as course}
						<a href="/courses/{course.uuid}">
							<div class="mb-2 w-[100%] rounded-xl bg-[#91F5AD] p-3 text-left">
								<h2 class="text-medium mb-2 font-semibold underline decoration-3">{course.name}</h2>
								<p class="text-ellipsis">{course.description}</p>
							</div>
						</a>
					{/each}
				{/snippet}
			</DataLoader>
		</div>
	</div>
</div>
