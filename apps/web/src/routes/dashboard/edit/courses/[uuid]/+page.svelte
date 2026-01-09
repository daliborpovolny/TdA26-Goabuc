<!-- <script lang="ts">
	import { page } from '$app/state';

	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course } from '$lib/types';

	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';
	import EditCourse from './EditCourse.svelte';
	import EditQuiz from './EditQuiz.svelte';

	let courseId: string = page.params.uuid!;

	let showCreateQuiz: boolean = $state(false)

	let refreshCount = $state(0);
	const reload = () => refreshCount++;

	let coursePromise = $derived.by(async () => {
		refreshCount;
		const res = await fetch(`/api/courses/${courseId}`);
		if (!res.ok) throw new Error('Failed to load course');
		return res.json();
	});
	$inspect(coursePromise)


</script>

<DataLoader promise={coursePromise}>
	{#snippet children(course: Course)}
		<div class="mx-auto max-w-2xl space-y-8 p-6">
			<section>
				<EditCourse {course} onchange={reload} />
			</section>

			<hr />

			<section>
				<h1 class="text-3xl">Materials</h1>
				<br />

				<CreateMaterial courseUuid={courseId} onchange={reload} />
				<br />

				{#if course.materials.length > 0}
					<h1 class="mb-3 text-2xl font-semibold text-gray-800">Edit Materials</h1>

					{#each course.materials as material}
						<EditMaterial {material} courseUuid={courseId} onchange={reload} />
						<br />
					{/each}
				{/if}

				<h1 class="text-3xl">Quizzes</h1>
				<br />

				{#if course.quizzes.length > 0} 

					{#each course.quizzes as quiz}
						
						<br>
						<h2 class="text-2xl">{quiz.title}</h2>
						<br>

						<EditQuiz edit={true} {quiz} courseId={course.uuid} onchange={reload} />
						<br>
						<hr>

						{/each}
				{/if}

				<button type="button" class="text-3xl" onclick={() => (showCreateQuiz = !showCreateQuiz)} >Create New Quiz</button>
				{#if showCreateQuiz}
					<br />
					
					<EditQuiz edit={false} courseId={course.uuid} onchange={() => (reload()) } />
			
				{/if}
					</section>
		</div>
	{/snippet}
</DataLoader> -->

<script lang="ts">

	import { page } from '$app/state';

	// import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course } from '$lib/types';

	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';
	import EditCourse from './EditCourse.svelte';
	import EditQuiz from './EditQuiz.svelte';

	let courseId: string = page.params.uuid!;

	let showCreateQuiz: boolean = $state(false)

let course = $state<Course | null>(null);
let loading = $state(true);
let error = $state<string | null>(null);

async function loadCourse() {
	if (course === null) {
		loading = true;
	}
	error = null;

	try {
		const res = await fetch(`/api/courses/${courseId}`);
		if (!res.ok) throw new Error('Failed to load course');
		course = await res.json();
	} catch (e) {
		error = e instanceof Error ? e.message : 'Unknown error';
	} finally {
		loading = false;
	}
}

// initial load
loadCourse();

function onCreateQuizSubmit () {
	loadCourse()
	showCreateQuiz = false
}

</script>


{#if loading}
	<div class="loading">Loading data...</div>
{:else if error}
	<div class="error">
		<p>Something went wrong: {error}</p>
		<button onclick={loadCourse}>Retry</button>
	</div>
{:else if course}
	<div class="mx-auto max-w-2xl space-y-8 p-6">
		<section>
			<EditCourse {course} onchange={loadCourse} />
		</section>

		<hr />

		<section>
			<h1 class="text-3xl">Materials</h1>
			<br />

			<CreateMaterial courseUuid={courseId} onchange={loadCourse} />
			<br />

			{#if course.materials.length > 0}
				<h1 class="mb-3 text-2xl font-semibold text-gray-800">Edit Materials</h1>

				{#each course.materials as material (material.uuid)}
					<EditMaterial {material} courseUuid={courseId} onchange={loadCourse} />
					<br />
				{/each}
			{/if}

			<h1 class="text-3xl">Quizzes</h1>
			<br />

			<button
				type="button"
				class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200"				
			onclick={() => (showCreateQuiz = !showCreateQuiz)}
			>
				New Quiz
			</button>
			{#if showCreateQuiz}
				<br />
				<EditQuiz edit={false} courseId={course.uuid} onchange={onCreateQuizSubmit} />
			{/if}
			<br><br>


			{#each course.quizzes as quiz (quiz.uuid)}

				<EditQuiz edit={true} {quiz} courseId={course.uuid} onchange={loadCourse} />
				<br />
			{/each}

			<br>

		</section>
	</div>
{/if}
