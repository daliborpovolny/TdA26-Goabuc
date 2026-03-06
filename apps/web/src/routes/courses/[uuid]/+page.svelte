<script lang="ts">
	import { page } from '$app/state';
	import { slide, fade } from 'svelte/transition';
	import type { Course, FullModule, Quiz, Material, Heading, Module } from '$lib/types';
	import CourseClosed from './CourseClosed.svelte';

	import ViewMaterial from './ViewMaterial.svelte';
	import TakeQuiz from './TakeQuiz.svelte';
	import ViewFeed from './ViewFeed.svelte';
	import { goto } from '$app/navigation';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import { auth } from '$lib/auth.svelte';
	import CourseInPreparation from './CourseInPreparation.svelte';

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
	let regularModules = $derived(
		(course?.modules ?? [])
			.filter((m) => m.name !== 'Unassigned' && m.state != 'preparation')
			.sort((a: Module, b: Module) => a.order - b.order)
	);

	let highlightedModule = $derived(
		course?.modules.find((m) => m.uuid === course?.highlightedModuleId)
	);
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
				<div class="flex flex-col justify-between gap-4 md:flex-row md:items-center">
					<div class="flex flex-wrap items-center gap-4">
						<h1 class="text-5xl font-black tracking-tighter text-s-black uppercase md:text-7xl">
							{course.name}
						</h1>

						{#if auth.user?.isAdmin}
							<div
								in:fade
								class="flex items-center gap-2 rounded-full border-4 border-s-black px-4 py-1 shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]
                    {course.state === 'open' ? 'bg-p-green' : ''}
                    {course.state === 'preparation' ? 'bg-amber-400' : ''}
                    {course.state === 'closed' ? 'bg-red-500 text-white' : ''}"
							>
								<span class="text-xs font-black tracking-widest uppercase">
									{#if course.state === 'open'}
										OPENED
									{:else if course.state === 'preparation'}
										IN PREPARATION
									{:else}
										CLOSED
									{/if}
								</span>
							</div>
						{/if}
					</div>

					{#if auth.user?.isAdmin}
						<div in:fade>
							<SecondaryButton
								href={`/dashboard/edit/courses/${course.uuid}`}
								class="!px-4 !py-2 !text-sm md:!text-base"
							>
								Edit Course ✏️
							</SecondaryButton>
						</div>
					{/if}
				</div>

				<div class="relative">
					<div class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"></div>
					<div class="relative rounded-2xl border-4 border-s-black bg-p-green p-6">
						<p class="text-xl font-bold text-s-black italic md:text-2xl">
							{course.description}
						</p>
					</div>
				</div>
			</header>

			{#if course.state === 'closed' && !auth.user?.isAdmin}
				<CourseClosed />
			{:else if course.state === 'preparation' && !auth.user?.isAdmin}
				<CourseInPreparation />
			{:else}
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
						class="w-full space-y-8 lg:w-2/3 {activeTab === 'modules'
							? 'block'
							: 'hidden lg:block'}"
					>
						{#if highlightedModule}
							<section in:fade class="relative mb-12">
								<div
									class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"
								></div>
								<div class="relative rounded-2xl border-4 border-s-black bg-p-blue p-6 text-white">
									<div class="mb-4 flex items-center gap-3">
										<span class="animate-bounce text-4xl">🌟</span>
										<div>
											<h2 class="text-xs font-black tracking-[0.2em] text-p-green uppercase">
												Current Focus
											</h2>
											<h3 class="text-2xl font-black tracking-tight uppercase">
												{highlightedModule.name}
											</h3>
										</div>
									</div>

									<div
										class="rounded-xl border-2 border-white/20 bg-s-black/20 p-4 font-bold italic"
									>
										"{course.highlightedModuleMessage || 'Focus on this module.'}"
									</div>
								</div>
							</section>
						{/if}

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
									{#if regularModules.length > 0}
										<div class="space-y-12">
											{#each regularModules as module, i}
												{#if module.state !== 'preparation'}
													{@const isHighlighted = module.uuid === course?.highlightedModuleId}
													{@const isClosed = module.state === 'closed'}

													<section class="relative">
														<div
															class="absolute -top-3 -left-3 z-10 flex h-10 w-10 items-center justify-center rounded-full border-4 border-s-black
            {isClosed
																? 'bg-red-500 text-white'
																: isHighlighted
																	? 'bg-p-blue text-white'
																	: 'bg-p-green text-s-black'} 
            font-black shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]"
														>
															{isClosed ? '✕' : isHighlighted ? '⭐' : i + 1}
														</div>

														<div
															class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]
            {isClosed ? 'border-red-600 grayscale-[0.5]' : ''} 
            {isHighlighted && !isClosed ? 'ring-4 ring-p-blue/30' : ''}"
														>
															<div
																class="mb-6 flex items-start justify-between border-b-2 border-gray-100 pb-4"
															>
																<div>
																	<h3
																		class="text-2xl font-black tracking-tight uppercase {isClosed
																			? 'text-red-600'
																			: ''}"
																	>
																		{module.name}
																	</h3>
																	<p class="font-bold text-gray-500 italic">
																		{isClosed
																			? 'This module is currently closed.'
																			: module.description}
																	</p>
																</div>

																{#if isClosed}
																	<span
																		class="rounded-lg border-2 border-s-black bg-red-600 px-3 py-1 text-[10px] font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
																	>
																		Closed
																	</span>
																{:else if isHighlighted}
																	<span
																		class="rounded-lg border-2 border-s-black bg-p-blue px-3 py-1 text-[10px] font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
																	>
																		Featured
																	</span>
																{/if}
															</div>

															{#if !isClosed}
																<div class="space-y-3">
																	{#each module.items as item}
																		{#if isQuiz(item)}
																			<TakeQuiz quiz={item} courseId={course.uuid} />
																		{:else if isMaterial(item)}
																			<ViewMaterial material={item} />
																		{/if}
																	{/each}
																</div>
															{:else}
																<div class="py-4 text-center">
																	<p
																		class="text-sm font-bold tracking-widest text-red-400 uppercase"
																	>
																		Content Locked
																	</p>
																</div>
															{/if}
														</div>
													</section>
												{/if}
											{/each}
										</div>
									{/if}
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

					<div
						class="w-full space-y-8 lg:w-1/3 {activeTab === 'feed' ? 'block' : 'hidden lg:block'}"
					>
						<div class="sticky top-10">
							<h2 class="mb-8 text-3xl font-black tracking-tight text-s-black uppercase">
								News Feed
							</h2>
							<div
								class="rounded-2xl border-4 border-s-black bg-p-blue/5 p-2 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
							>
								<ViewFeed courseId={course.uuid} onUpdate={loadCourse} />
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>
