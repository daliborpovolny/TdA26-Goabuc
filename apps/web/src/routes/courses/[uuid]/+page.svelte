<script lang="ts">
	import { page } from '$app/state';
	import { slide, fade } from 'svelte/transition';
	import type { Course, FullModule, Quiz, Material, Heading } from '$lib/types';

	import ViewMaterial from './ViewMaterial.svelte';
	import TakeQuiz from './TakeQuiz.svelte';
	import ViewFeed from './ViewFeed.svelte';
	import { goto } from '$app/navigation';

	let activeTab = $state('modules');

	let course = $state<Course | null>(null);

	let loading = $state(true);
	let error = $state<string | null>(null);

	async function loadCourse() {
		error = null;
		try {
			const res = await fetch(`/api/courses/${page.params.uuid}`);
			if (!res.ok) {
				if (res.status == 404) {
					throw new Error('Unknown Course');
				}

				throw new Error('Failed to load course');
			}
			course = await res.json();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			loading = false;
		}
	}
	loadCourse();

	// Helper to check if item is Quiz or Material
	function isQuiz(item: Quiz | Material | Heading): item is Quiz {
		return 'questions' in item;
	}

	function isMaterial(item: Material | Heading): item is Material {
		return !('content' in item);
	}
	let unassignedModule = $derived(course?.modules.find((m) => m.name === 'Unassigned'));
	let regularModules = $derived(course?.modules.filter((m) => m.name !== 'Unassigned') ?? []);
</script>

<svelte:head>
	<title>{course ? course.name + ' | TdA' : 'Loading...'}</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-4 md:p-10">
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
			<p class="text-2xl font-black">ERROR</p>
			<p class="mt-2 font-bold">{error}</p>
			<button
				onclick={loadCourse}
				class="mt-4 rounded-lg bg-s-black px-4 py-2 font-bold hover:bg-white hover:text-s-black"
				>Retry Connection</button
			>
			<button
				onclick={() => goto('/courses')}
				class="mt-4 rounded-lg bg-s-black px-4 py-2 font-bold hover:bg-white hover:text-s-black"
				>Go back</button
			>
		</div>
	{:else if course}
		<div class="mx-auto flex w-full flex-col gap-8 md:max-w-[90%]">
			<header class="space-y-4">
				<h1 class="text-5xl font-black tracking-tighter text-s-black uppercase md:text-7xl">
					{course.name}
				</h1>
				<div class="relative">
					<div class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"></div>
					<div class="relative rounded-2xl border-4 border-s-black bg-p-green p-6">
						<p class="text-xl font-bold text-s-black italic md:text-2xl">
							{course.description}
						</p>
					</div>
				</div>
			</header>

			<div
				class="sticky top-4 z-20 flex rounded-xl border-4 border-s-black bg-s-black p-1 shadow-xl lg:hidden"
			>
				<button
					onclick={() => (activeTab = 'modules')}
					class="flex-1 rounded-lg py-3 text-lg font-black tracking-widest uppercase transition-all
                        {activeTab === 'modules' ? 'bg-p-green text-s-black' : 'text-white'}"
				>
					Modules
				</button>
				<button
					onclick={() => (activeTab = 'feed')}
					class="flex-1 rounded-lg py-3 text-lg font-black tracking-widest uppercase transition-all
                        {activeTab === 'feed' ? 'bg-p-green text-s-black' : 'text-white'}"
				>
					Feed
				</button>
			</div>

			<div class="flex flex-col gap-8 lg:flex-row">
				<div
					class="w-full space-y-8 lg:w-2/3 {activeTab === 'modules' ? 'block' : 'hidden lg:block'}"
				>
					<h2 class="text-3xl font-black tracking-tight text-p-blue uppercase">Learning Path</h2>

					<div class="space-y-12">
						{#if unassignedModule && unassignedModule.items.length > 0}
							<section class="space-y-4">
								<div class="flex items-center gap-2">
									<span class="text-2xl">🎒</span>
									<h3 class="text-xl font-black tracking-widest text-gray-400 uppercase">
										General Resources
									</h3>
								</div>
								<div
									class="space-y-3 rounded-2xl border-4 border-dashed border-s-black/20 bg-gray-50/50 p-6"
								>
									{#each unassignedModule.items as item}
										{#if isQuiz(item)}
											<TakeQuiz quiz={item} courseId={course.uuid} />
										{:else if isMaterial(item)}
											<ViewMaterial material={item} />
										{/if}
									{/each}
								</div>
							</section>
						{/if}

						{#if regularModules.length > 0}
							<div class="space-y-12">
								{#each regularModules as module, i}
									<section class="relative">
										<div
											class="absolute -top-3 -left-3 z-10 flex h-10 w-10 items-center justify-center rounded-full border-4 border-s-black bg-p-green font-black shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]"
										>
											{i + 1}
										</div>

										<div
											class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
										>
											<div class="mb-6 border-b-2 border-gray-100 pb-4">
												<h3 class="text-2xl font-black tracking-tight uppercase">{module.name}</h3>
												<p class="font-bold text-gray-500 italic">{module.description}</p>
											</div>

											<div class="space-y-3">
												{#each module.items as item}
													{#if isQuiz(item)}
														<TakeQuiz quiz={item} courseId={course.uuid} />
													{:else if isMaterial(item)}
														<ViewMaterial material={item} />
													{/if}
												{/each}
											</div>
										</div>
									</section>
								{/each}
							</div>
						{:else if !unassignedModule || unassignedModule.items.length === 0}
							<div
								class="rounded-2xl border-4 border-dashed border-gray-300 p-12 text-center font-bold text-gray-400"
							>
								No modules have been released for this course yet.
							</div>
						{/if}
					</div>
				</div>

				<div class="w-full space-y-8 lg:w-1/3 {activeTab === 'feed' ? 'block' : 'hidden lg:block'}">
					<div class="sticky top-10">
						<h2 class="mb-8 text-3xl font-black tracking-tight text-s-black uppercase">
							News Feed
						</h2>
						<div
							class="rounded-2xl border-4 border-s-black bg-p-blue/5 p-2 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
						>
							<ViewFeed courseId={course.uuid} />
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
