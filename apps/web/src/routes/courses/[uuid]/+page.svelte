<script lang="ts">
	import { page } from '$app/state';
	import ViewMaterial from './ViewMaterial.svelte';
	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course } from '$lib/types';
	import TakeQuiz from './TakeQuiz.svelte';

	let coursesPromise: Promise<any> = loadCourseDetail();

	async function loadCourseDetail() {
		return fetch('/api/courses/' + page.params.uuid, {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		})
			.then(async (res) => {
				if (res.status == 404) {
					throw new Error('Unknown course');
				}

				if (!res.ok) {
					throw new Error('Failed to get course info');
				}
				return res.json();
			})
			.then((data) => {
				return data;
			});
	}
</script>

<br />
<DataLoader promise={coursesPromise}>
	{#snippet children(course: Course)}
		<div class="mx-auto max-w-2xl space-y-8 p-6">
		<title>{course.name}</title>
		<div class="ml-5">
			<div>
				<h1 class="text-5xl underline">{course.name}</h1>
				<br />
				<p>{course.description}</p>
			</div>

			<div>
				{#if course.materials && course.materials.length > 0}
					<br />
					<h1 class="text-3xl">Materials</h1>
					<br />

					{#each course.materials as material}
						<ViewMaterial {material} />
						<br />
						<br />
					{/each}
				{/if}
			</div>

			<div>
				{#if course.quizzes && course.quizzes.length > 0}
					<br>
					<h1 class="text-3xl">Quizzes</h1>
					<br>

					{#each course.quizzes as quiz}
						<TakeQuiz {quiz} courseId={course.uuid}/>
						<br>
						<br>
					{/each}

				{/if}
			</div>
		</div>
		</div>
	{/snippet}
</DataLoader>
