<script lang="ts">
	import { page } from '$app/state';
	import ViewMaterial from './ViewMaterial.svelte';
	import type { Course } from '$lib/types';
	import TakeQuiz from './TakeQuiz.svelte';
	import ViewFeed from './ViewFeed.svelte';

	let activeTab = $state('materials');

	let message: string = $state('loading');
	let course: null | Course = $state(null);
	loadCourseDetail();

	async function loadCourseDetail() {
		const res = await fetch('/api/courses/' + page.params.uuid);
		if (!res.ok) {
			if (res.status === 404) {
				message = 'Unknown course';
			} else {
				message = 'Failed to load course detail';
			}
			return;
		}
		course = await res.json();
	}
</script>

<svelte:head>
	<title>{course?.name || 'Tda course'}</title>
</svelte:head>

{#if course !== null}
	<div class="min-h-screen bg-s-white p-4 md:p-10">
		<div class="mx-auto flex w-full flex-col gap-6 md:max-w-8/9">
			<div class="space-y-4">
				<h1 class="text-4xl font-bold text-s-black md:text-6xl">{course.name}</h1>
				<div class="rounded-xl bg-p-blue p-1.5 shadow-lg">
					<div class="rounded-lg bg-p-green p-4">
						<p class="text-xl text-s-black md:text-2xl">{course.description}</p>
					</div>
				</div>
			</div>

			<div class="sticky top-2 z-10 flex gap-2 rounded-xl bg-s-black p-1 shadow-md lg:hidden">
				{#each ['materials', 'quizzes', 'feed'] as tab}
					<button
						onclick={() => (activeTab = tab)}
						class="flex-1 cursor-pointer rounded-lg py-2 text-lg font-bold capitalize transition-all
                            {activeTab === tab
							? 'bg-p-green text-s-black'
							: 'text-s-white hover:bg-p-blue'}"
					>
						{tab}
					</button>
				{/each}
			</div>

			<div class="mt-2 flex min-h-[400px] flex-col lg:flex-row lg:space-x-4">
				<div
					class="w-full space-y-3 lg:w-1/3 {activeTab === 'materials'
						? 'block'
						: 'hidden lg:block'}"
				>
					<h2 class="text-center text-2xl font-bold">Course Materials</h2>
					{#if course.materials?.length}
						<div class="space-y-2 rounded-xl bg-p-blue p-2">
							{#each course.materials as material}
								<ViewMaterial {material} />
							{/each}
						</div>
					{:else}
						<p class="text-center text-gray-500 italic">No materials available yet.</p>
					{/if}
				</div>

				<div
					class="w-full space-y-3 lg:w-1/3 {activeTab === 'quizzes' ? 'block' : 'hidden lg:block'}"
				>
					<h2 class="text-center text-2xl font-bold">Available Quizzes</h2>
					{#if course.quizzes?.length}
						<div class="space-y-2 rounded-xl bg-p-blue p-2">
							{#each course.quizzes as quiz}
								<TakeQuiz {quiz} courseId={course.uuid} />
							{/each}
						</div>
					{:else}
						<p class="text-center text-gray-500 italic">No quizzes assigned.</p>
					{/if}
				</div>

				<div class="w-full space-y-3 lg:w-1/3 {activeTab === 'feed' ? 'block' : 'hidden lg:block'}">
					<h2 class="text-center text-2xl font-bold">News Feed</h2>
					<ViewFeed courseId={course.uuid} />
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="text-2xl">
		{message}
	</div>
{/if}
