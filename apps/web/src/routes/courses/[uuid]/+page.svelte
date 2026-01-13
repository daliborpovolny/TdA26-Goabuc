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

<DataLoader promise={coursesPromise}>
	{#snippet children(course: Course)}
		<div class="p-6 flex justify-center">
			<title>{course.name}</title>
			<div class="space-y-8 w-[50%]">
				<h1 class="text-5xl font-bold mb-0">{course.name}</h1>
				<div class="bg-[#0070BB] rounded-xl p-2">
					<div class="bg-[#91F5AD] rounded-xl p-4">
						<p class="text-2xl">{course.description}</p>
					</div>
				</div>

				<div>
					{#if course.materials && course.materials.length > 0}
						<h1 class="text-3xl font-bold">Materials</h1>
						<div class="bg-[#0070BB] p-2 space-y-2 rounded-xl">
							{#each course.materials as material}
								<ViewMaterial {material} />
							{/each}
						</div>
					{/if}
				</div>

				<div>
					{#if course.quizzes && course.quizzes.length > 0}
						<h1 class="text-3xl font-bold">Quizzes</h1>
						<div class="bg-[#0070BB] p-2 space-y-2 rounded-xl">
							{#each course.quizzes as quiz}
								<TakeQuiz {quiz} courseId={course.uuid} />
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/snippet}
</DataLoader>
