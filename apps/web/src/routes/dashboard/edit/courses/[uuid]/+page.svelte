<script lang="ts">
	import { page } from '$app/state';
	import { fade, slide } from 'svelte/transition';
	import { flip } from 'svelte/animate';
	import type { Course, Module } from '$lib/types';

	import EditCourse from './EditCourse.svelte';
	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';
	import EditQuiz from './EditQuiz.svelte';
	import CreateFeedPost from './CreateFeedPost.svelte';
	import EditFeed from './EditFeed.svelte';
	import CreateModule from './CreateModule.svelte';
	import EditModule from './EditModule.svelte';
	import ModuleSelector from './ModuleSelector.svelte';

	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';
	import WarningButton from '$lib/components/WarningButton.svelte';

	let courseId = page.params.uuid!;
	let activeSection = $state('general');
	let showCreateQuiz = $state(false);
	let showCreateModule = $state(false);
	let course = $state<Course | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let sections = [
		{ id: 'general', label: 'General', icon: '📝' },
		{ id: 'modules', label: 'Modules', icon: '📦' },
		{ id: 'materials', label: 'Materials', icon: '📁' },
		{ id: 'quizzes', label: 'Quizzes', icon: '📝' },
		{ id: 'feed', label: 'Feed', icon: '💬' }
	];

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

	//* Module reordering
	let localModules = $state<Module[]>([]);
	let orderChanged = $state(false);
	let isSavingOrder = $state(false);

	$effect(() => {
		if (course) {
			localModules = course.modules
				.filter((m) => m.name !== 'Unassigned')
				.sort((a, b) => (a.order ?? 0) - (b.order ?? 0));
		}
	});

	function moveModule(index: number, direction: 'up' | 'down') {
		const newIndex = direction === 'up' ? index - 1 : index + 1;
		if (newIndex < 0 || newIndex >= localModules.length) return;
		const temp = localModules[index];
		localModules[index] = localModules[newIndex];
		localModules[newIndex] = temp;
		orderChanged = true;
	}

	async function saveModuleOrder() {
		isSavingOrder = true;
		try {
			for (let i = 0; i < localModules.length; i++) {
				const mod = localModules[i];
				await fetch(`/api/courses/${courseId}/modules/${mod.uuid}/state`, {
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ state: mod.state, order: i + 2 })
				});
			}
			orderChanged = false;
			loadCourse();
		} finally {
			isSavingOrder = false;
		}
	}

	//* Highlight stuff
	let highlightedModuleId = $derived(course?.highlightedModuleId || '');
	let highlightedMessage = $derived(course?.highlightedModuleMessage || '');
	let isSavingHighlight = $state(false);

	async function updateHighlight(newUuid?: string) {
		isSavingHighlight = true;
		const targetUuid = newUuid || highlightedModuleId;
		try {
			await fetch(`/api/courses/${courseId}/state`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					highlightedModuleId: targetUuid,
					highlightedModuleMessage: highlightedMessage,
					state: course?.state
				})
			});
			loadCourse();
		} finally {
			isSavingHighlight = false;
		}
	}

	let activeHighlight = $derived(course?.modules.find((m) => m.uuid === highlightedModuleId));
</script>

<svelte:head>
	<title>Editor: {course?.name || 'Course'} | TdA</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-4 md:p-8">
	{#if loading}
		<div class="flex h-64 flex-col items-center justify-center space-y-4">
			<div
				class="h-12 w-12 animate-spin rounded-full border-4 border-p-blue border-t-p-green"
			></div>
			<p class="font-black tracking-widest uppercase">Syncing Data...</p>
		</div>
	{:else if error}
		<div
			class="mx-auto max-w-md rounded-xl border-4 border-s-black bg-red-500 p-6 text-white shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]"
		>
			<p class="text-2xl font-black tracking-tighter uppercase">System Error</p>
			<p class="mt-2 font-bold">{error}</p>
			<div class="mt-6 flex gap-4">
				<PrimaryButton onclick={loadCourse} class="border-2">Retry Connection</PrimaryButton>
				<SecondaryButton onclick={() => history.back()} class="border-2">Go Back</SecondaryButton>
			</div>
		</div>
	{:else if course}
		<div class="mx-auto w-full max-w-[1600px]">
			<header
				class="mb-10 flex flex-col justify-between gap-6 border-b-8 border-s-black pb-8 md:flex-row md:items-end"
			>
				<div class="space-y-2">
					<a href="/dashboard" class="text-xs font-black text-p-blue uppercase hover:underline"
						>← Dashboard</a
					>
					<h1 class="text-4xl font-black tracking-tighter uppercase md:text-7xl">
						Editor: <span class="text-p-blue">{course.name}</span>
					</h1>
				</div>
				<SecondaryButton href="/courses/{course.uuid}" target="_blank" class="py-2!">
					View Live 👁️
				</SecondaryButton>
			</header>

			<div class="grid grid-cols-1 gap-6 lg:grid-cols-12 lg:gap-10">
				<div class="nav-mask lg:col-span-3 xl:col-span-2">
					<div class="mb-2 flex items-center justify-between px-1 lg:hidden">
						<span class="text-[10px] font-black tracking-widest text-gray-400 uppercase"
							>Sections</span
						>
						<span class="animate-pulse text-[10px] font-black tracking-widest text-p-blue uppercase"
							>Swipe for more →</span
						>
					</div>

					<nav
						class="brutal-scroll sticky top-4 z-40 flex flex-nowrap gap-2 overflow-x-auto bg-s-white/80 pb-4 backdrop-blur-sm lg:static lg:flex-col lg:overflow-visible lg:pb-0"
					>
						{#each sections as section}
							{#if activeSection === section.id}
								<PrimaryButton
									class="min-w-fit shrink-0 translate-x-1 translate-y-1 shadow-none! lg:w-full lg:justify-start lg:gap-4"
								>
									<span class="text-xl">{section.icon}</span>
									<span class="text-xs font-black uppercase">{section.label}</span>
								</PrimaryButton>
							{:else}
								<SecondaryButton
									onclick={() => (activeSection = section.id)}
									class="min-w-fit shrink-0 lg:w-full lg:justify-start lg:gap-4"
								>
									<span class="text-xl">{section.icon}</span>
									<span class="text-xs font-black uppercase">{section.label}</span>
								</SecondaryButton>
							{/if}
						{/each}
					</nav>
				</div>

				<main class="lg:col-span-9 xl:col-span-10">
					<div
						class="rounded-2xl border-4 border-s-black bg-white p-4 shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] md:rounded-3xl md:p-12 md:shadow-[12px_12px_0px_0px_rgba(0,0,0,1)]"
					>
						{#if activeSection === 'general'}
							<section in:fade>
								<h2 class="mb-8 text-3xl font-black text-p-blue uppercase italic">
									Course Settings
								</h2>
								<EditCourse {course} onchange={loadCourse} />
							</section>
						{:else if activeSection === 'modules'}
							<div in:fade class="space-y-10">
								{#if course.modules.length > 1}
									<section
										class="rounded-2xl border-4 border-s-black bg-p-blue/5 p-4 shadow-[6px_6px_0px_0px_rgba(2,87,165,1)] md:p-8"
									>
										<div
											class="mb-6 flex flex-col justify-between gap-4 md:flex-row md:items-center"
										>
											<div>
												<h2 class="text-xl font-black text-p-blue uppercase">Featured Focus</h2>
												<p class="text-[10px] font-bold text-gray-500 uppercase">
													Pinned to the top for students
												</p>
											</div>
											<ModuleSelector
												modules={localModules}
												bind:selectedId={highlightedModuleId}
												label="Select Highlight"
											/>
										</div>

										{#if activeHighlight}
											<div transition:slide class="space-y-6">
												<div
													class="flex items-center gap-4 rounded-xl border-4 border-s-black bg-white p-4 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
												>
													<span class="text-4xl">🌟</span>
													<h3 class="text-xl font-black tracking-tighter uppercase">
														{activeHighlight.name}
													</h3>
												</div>
												<div class="flex flex-col gap-3 md:flex-row">
													<input
														type="text"
														bind:value={highlightedMessage}
														placeholder="Add a note for students..."
														class="flex-1 rounded-xl border-4 border-s-black bg-white p-3 font-bold outline-none focus:ring-4 focus:ring-p-green"
													/>
													<SuccessButton
														onclick={() => updateHighlight()}
														isSaving={isSavingHighlight}
														class="py-3!">Update Highlight</SuccessButton
													>
												</div>
											</div>
										{/if}
									</section>
								{/if}

								<div class="space-y-6">
									<div
										class="flex flex-col justify-between gap-4 border-b-4 border-gray-100 pb-4 md:flex-row md:items-center"
									>
										<h2 class="text-3xl font-black uppercase">Curriculum</h2>
										<div class="flex gap-2">
											{#if orderChanged}
												<SuccessButton
													onclick={saveModuleOrder}
													isSaving={isSavingOrder}
													class="py-2! text-xs!">Save Order</SuccessButton
												>
											{/if}

											{#if showCreateModule}
												<WarningButton
													onclick={() => (showCreateModule = false)}
													class="py-2! text-xs!">Cancel</WarningButton
												>
											{:else}
												<PrimaryButton
													onclick={() => (showCreateModule = true)}
													class="py-2! text-xs!">+ New Module</PrimaryButton
												>
											{/if}
										</div>
									</div>

									{#if showCreateModule}
										<div transition:slide class="rounded-xl border-4 border-s-black bg-gray-50 p-6">
											<CreateModule
												courseId={course.uuid}
												onchange={() => {
													loadCourse();
													showCreateModule = false;
												}}
											/>
										</div>
									{/if}

									<div class="space-y-6">
										{#each localModules as module, i (module.uuid)}
											<div
												animate:flip={{ duration: 300 }}
												class="group relative flex items-start gap-4"
											>
												<div
													class="flex flex-col gap-1 transition-opacity lg:absolute lg:top-1/2 lg:-left-12 lg:-translate-y-1/2 lg:opacity-0 lg:group-hover:opacity-100"
												>
													<button
														disabled={i === 0}
														onclick={() => moveModule(i, 'up')}
														class="rounded border-2 border-s-black bg-white p-2 hover:bg-p-green disabled:opacity-30"
														>▲</button
													>
													<button
														disabled={i === localModules.length - 1}
														onclick={() => moveModule(i, 'down')}
														class="rounded border-2 border-s-black bg-white p-2 hover:bg-p-green disabled:opacity-30"
														>▼</button
													>
												</div>
												<div class="relative flex-1">
													{#if course.highlightedModuleId === module.uuid}
														<div
															class="absolute -top-3 -right-3 z-10 rounded-full border-2 border-s-black bg-p-blue px-3 py-1 text-[10px] font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
														>
															Featured
														</div>
													{/if}
													<EditModule {module} courseId={course.uuid} onchange={loadCourse} />
												</div>
											</div>
										{/each}
									</div>
								</div>
							</div>
						{:else if activeSection === 'materials'}
							<div in:fade class="space-y-10">
								<h2 class="text-3xl font-black uppercase">Materials</h2>
								<CreateMaterial
									courseUuid={courseId}
									onchange={loadCourse}
									modules={course.modules}
								/>
								{#if course.materials.length > 0}
									<div class="space-y-4 border-t-4 border-dashed border-gray-100 pt-10">
										{#each course.materials as material (material.uuid)}
											<EditMaterial
												{material}
												modules={course.modules}
												courseUuid={courseId}
												onchange={loadCourse}
											/>
										{/each}
									</div>
								{/if}
							</div>
						{:else if activeSection === 'quizzes'}
							<div in:fade class="space-y-6">
								<div class="flex items-center justify-between">
									<h2 class="text-3xl font-black uppercase">Exams</h2>
									{#if showCreateQuiz}
										<WarningButton onclick={() => (showCreateQuiz = false)} class="py-2!"
											>Cancel</WarningButton
										>
									{:else}
										<PrimaryButton onclick={() => (showCreateQuiz = true)} class="py-2!"
											>+ New Quiz</PrimaryButton
										>
									{/if}
								</div>
								{#if showCreateQuiz}
									<div transition:slide class="rounded-xl border-4 border-s-black bg-gray-50 p-6">
										<EditQuiz
											modules={course.modules}
											edit={false}
											courseId={course.uuid}
											onchange={() => {
												loadCourse();
												showCreateQuiz = false;
											}}
										/>
									</div>
								{/if}
								<div class="space-y-4">
									{#each course.quizzes as quiz (quiz.uuid)}
										<EditQuiz
											modules={course.modules}
											edit={true}
											{quiz}
											courseId={course.uuid}
											onchange={loadCourse}
										/>
									{/each}
								</div>
							</div>
						{:else if activeSection === 'feed'}
							<div in:fade class="space-y-10">
								<h2 class="text-3xl font-black uppercase">Course Feed</h2>
								<CreateFeedPost courseId={course.uuid} />
								<EditFeed courseId={course.uuid} />
							</div>
						{/if}
					</div>
				</main>
			</div>
		</div>
	{/if}
</div>

<style>
	.brutal-scroll::-webkit-scrollbar {
		height: 8px;
	}
	.brutal-scroll::-webkit-scrollbar-track {
		background: #f3f4f6;
		border-top: 2px solid #1a1a1a;
		border-bottom: 2px solid #1a1a1a;
	}
	.brutal-scroll::-webkit-scrollbar-thumb {
		background: #1a1a1a;
		border: 1px solid #fff;
	}

	.nav-mask {
		position: relative;
	}

	@media (max-width: 1024px) {
		.nav-mask::after {
			content: '';
			position: absolute;
			top: 0;
			right: 0;
			height: calc(100% - 8px); /* Account for scrollbar height */
			width: 50px;
			background: linear-gradient(to right, transparent, rgba(255, 255, 255, 0.9));
			pointer-events: none;
			z-index: 50;
		}
	}
</style>
