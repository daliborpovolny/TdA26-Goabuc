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
		<div class="min-w-fit p-6 flex justify-center">
			<title>{course.name}</title>
			<div class="space-y-8">
				<div>
					<h1 class="text-5xl font-bold underline decoration-5">{course.name}</h1>
					<br />
					<p class="text-2xl">{course.description}</p>
				</div>

				<div>
					{#if course.materials && course.materials.length > 0}
						<h1 class="text-3xl font-semibold">Materials</h1>
						<br>
						{#each course.materials as material}
							<ViewMaterial {material} />
						{/each}
					{/if}
				</div>

				<div>
					{#if course.quizzes && course.quizzes.length > 0}
						<h1 class="text-3xl font-semibold">Quizzes</h1>
						<br>
						{#each course.quizzes as quiz}
							<TakeQuiz {quiz} courseId={course.uuid} />
						{/each}
					{/if}
				</div>
			</div>
		</div>
	{/snippet}
</DataLoader>
