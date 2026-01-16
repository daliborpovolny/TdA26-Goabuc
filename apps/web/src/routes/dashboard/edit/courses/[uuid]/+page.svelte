<script lang="ts">
	import { page } from '$app/state';
	import { fade, slide } from 'svelte/transition';
	import type { Course } from '$lib/types';

	// Subcomponents
	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';
	import EditCourse from './EditCourse.svelte';
	import EditQuiz from './EditQuiz.svelte';
	import CreateFeedPost from './CreateFeedPost.svelte';
	import EditFeed from './EditFeed.svelte';

	let courseId = page.params.uuid!;
	let showCreateQuiz = $state(false);
	let activeSection = $state('general'); // 'general', 'materials', 'quizzes', 'feed'

	let course = $state<Course | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	async function loadCourse() {
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

	loadCourse();

	function onCreateQuizSubmit() {
		loadCourse();
		showCreateQuiz = false;
	}
</script>

<svelte:head>
	<title>Editing: {course?.name || 'Course'} | TdA</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-4 md:p-8">
	{#if loading}
		<div class="flex h-64 flex-col items-center justify-center space-y-4">
			<div
				class="h-12 w-12 animate-spin rounded-full border-4 border-p-blue border-t-p-green"
			></div>
			<p class="font-black tracking-widest text-s-black uppercase">Syncing Academy Data...</p>
		</div>
	{:else if error}
		<div
			class="mx-auto max-w-md rounded-xl border-4 border-s-black bg-red-500 p-6 text-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
		>
			<p class="text-2xl font-black">SYSTEM ERROR</p>
			<p class="mt-2 font-bold">{error}</p>
			<button
				onclick={loadCourse}
				class="mt-4 rounded-lg bg-s-black px-4 py-2 font-bold hover:bg-white hover:text-s-black"
				>Retry Connection</button
			>
		</div>
	{:else if course}
		<div class="mx-auto max-w-5xl">
			<header
				class="mb-8 flex flex-col justify-between gap-4 border-b-4 border-s-black pb-6 md:flex-row md:items-end"
			>
				<div>
					<a href="/dashboard" class="text-sm font-black text-p-blue uppercase hover:underline"
						>← Dashboard</a
					>
					<h1 class="text-4xl font-black tracking-tighter uppercase md:text-6xl">
						Editor: <span class="text-p-blue">{course.name}</span>
					</h1>
				</div>
				<a
					href="/courses/{course.uuid}"
					target="_blank"
					class="rounded-lg border-2 border-s-black bg-p-green px-4 py-2 font-bold shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] transition hover:shadow-none"
				>
					View 👁️
				</a>
			</header>

			<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
				<nav class="space-y-2 lg:col-span-3">
					{#each [{ id: 'general', label: 'General Info', icon: '📝' }, { id: 'materials', label: 'Materials', icon: '📁' }, { id: 'quizzes', label: 'Quizzes', icon: '📝' }, { id: 'feed', label: 'Course Feed', icon: '💬' }] as section}
						<button
							onclick={() => (activeSection = section.id)}
							class="flex w-full cursor-pointer items-center gap-3 rounded-xl border-2 border-s-black p-4 text-left font-black tracking-tight uppercase transition-all
                            {activeSection === section.id
								? 'translate-x-1 bg-p-blue text-white shadow-none'
								: 'bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] hover:bg-p-green'}"
						>
							<span>{section.icon}</span>
							{section.label}
						</button>
					{/each}
				</nav>

				<main class="lg:col-span-9">
					<div
						class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[8px_8px_0px_0px_rgba(26,26,26,1)]"
					>
						{#if activeSection === 'general'}
							<div in:fade>
								<h2 class="mb-6 text-3xl font-black text-p-blue uppercase italic">
									Course Settings
								</h2>
								<EditCourse {course} onchange={loadCourse} />
							</div>
						{:else if activeSection === 'materials'}
							<div in:fade class="space-y-8">
								<div>
									<h2 class="mb-4 text-3xl font-black text-s-2 uppercase">Add Content</h2>
									<CreateMaterial courseUuid={courseId} onchange={loadCourse} />
								</div>

								{#if course.materials.length > 0}
									<div class="space-y-4 border-t-2 border-dashed border-gray-200 pt-6">
										<h3 class="text-xl font-black uppercase">
											Existing Materials ({course.materials.length})
										</h3>
										{#each course.materials as material (material.uuid)}
											<EditMaterial {material} courseUuid={courseId} onchange={loadCourse} />
										{/each}
									</div>
								{/if}
							</div>
						{:else if activeSection === 'quizzes'}
							<div in:fade class="space-y-6">
								<div class="flex items-center justify-between">
									<h2 class="text-3xl font-black text-p-blue uppercase">Examinations</h2>
									<button
										onclick={() => (showCreateQuiz = !showCreateQuiz)}
										class="rounded-lg border-2 border-s-black bg-p-green px-4 py-2 font-bold uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] hover:shadow-none"
									>
										{showCreateQuiz ? 'Cancel' : '+ New Quiz'}
									</button>
								</div>

								{#if showCreateQuiz}
									<div transition:slide class="rounded-xl border-2 border-s-black bg-gray-50 p-4">
										<EditQuiz edit={false} courseId={course.uuid} onchange={onCreateQuizSubmit} />
									</div>
								{/if}

								<div class="space-y-4">
									{#each course.quizzes as quiz (quiz.uuid)}
										<EditQuiz edit={true} {quiz} courseId={course.uuid} onchange={loadCourse} />
									{/each}
								</div>
							</div>
						{:else if activeSection === 'feed'}
							<div in:fade class="space-y-8">
								<div>
									<h2 class="mb-4 text-3xl font-black text-s-3 uppercase">Post Update</h2>
									<CreateFeedPost courseId={course.uuid} />
								</div>
								<div class="border-t-2 border-dashed border-gray-200 pt-6">
									<h3 class="mb-4 text-xl font-black text-gray-500 uppercase">History</h3>
									<EditFeed courseId={course.uuid} />
								</div>
							</div>
						{/if}
					</div>
				</main>
			</div>
		</div>
	{/if}
</div>
