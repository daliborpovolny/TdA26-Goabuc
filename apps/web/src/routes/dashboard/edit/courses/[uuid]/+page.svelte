<script lang="ts">
	import { page } from '$app/state';
	import { fade, slide } from 'svelte/transition';
	import type { Course, Module } from '$lib/types';

	// Subcomponents
	import EditCourse from './EditCourse.svelte';

	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';

	import EditQuiz from './EditQuiz.svelte';

	import CreateFeedPost from './CreateFeedPost.svelte';
	import EditFeed from './EditFeed.svelte';

	import CreateModule from './CreateModule.svelte';
	import EditModule from './EditModule.svelte';
	import UniLink from '../../../../UniLink.svelte';
	import UniButton from '../../../../UniButton.svelte';

	let courseId = page.params.uuid!;

	let sections = [
		{ id: 'general', label: 'General Info', icon: '📝' },
		{ id: 'modules', label: 'Modules', icon: '📦' },
		{ id: 'materials', label: 'Materials', icon: '📁' },
		{ id: 'quizzes', label: 'Quizzes', icon: '📝' },
		{ id: 'feed', label: 'Course Feed', icon: '💬' }
	];

	let activeSection = $state('general'); // 'general', 'materials', 'quizzes', 'feed', 'modules'

	let showCreateQuiz = $state(false);
	let showCreateModule = $state(false);

	let course = $state<Course | null>(null);
	// $inspect(course);

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

	function onCreateModuleSubmit() {
		loadCourse();
		showCreateModule = false;
	}

	import { flip } from 'svelte/animate';
	// ... other imports
	import SuccessButton from '$lib/components/SuccessButton.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	// ... existing loadCourse logic

	let localModules = $state<Module[]>([]);
	let orderChanged = $state(false);
	let isSavingOrder = $state(false);

	// Sync localModules whenever course updates
	$effect(() => {
		if (course) {
			// Filter out Unassigned and sort by current order
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

				const res = await fetch(`/api/courses/${courseId}/modules/${mod.uuid}/state`, {
					method: 'PUT',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({
						state: mod.state,
						order: i + 2 // offset by 2 because buffer module has order 1
					})
				});

				if (!res.ok) throw new Error(`Failed to update module ${mod.name}`);
			}

			orderChanged = false;
			loadCourse();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save order';
		} finally {
			isSavingOrder = false;
		}
	}

	import ModuleSelector from './ModuleSelector.svelte';

	let highlightedModuleId = $state(course?.highlightedModuleId || '');
	let highlightedMessage = $state(course?.highlightedModuleMessage || '');
	let originalMessage = course?.highlightedModuleMessage || '';

	let isSavingHighlight = $state(false);

	async function updateHighlight(newUuid?: string) {
		isSavingHighlight = true;

		// Use the newUuid if provided (from the selector), otherwise use the current state
		const targetUuid = newUuid || highlightedModuleId;

		try {
			const res = await fetch(`/api/courses/${courseId}/state`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					highlightedModuleId: targetUuid,
					highlightedModuleMessage: highlightedMessage,
					state: course?.modules.filter((m) => m.uuid !== targetUuid)[0].state
				})
			});

			if (res.ok) {
				originalMessage = highlightedMessage; // Sync original to hide save button
				if (newUuid) highlightedModuleId = newUuid;
				loadCourse();
			}
		} finally {
			isSavingHighlight = false;
		}
	}

	// Derived to find the full module object for the "Preview" card
	let activeHighlight = $derived(course?.modules.find((m) => m.uuid === highlightedModuleId));
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
			<UniButton
				onclick={loadCourse}
				content="Retry Connection"
				bgcolor="bg-[#444]"
				hv_bgcolor=""
				more_style="text-white hover:text-s-black hover:bg-white"
			/>
		</div>
	{:else if course}
		<div class="mx-auto w-full max-w-[1600px]">
			<header
				class="mb-10 flex flex-col justify-between gap-6 border-b-8 border-s-black pb-8 md:flex-row md:items-end"
			>
				<div>
					<a href="/dashboard" class="text-sm font-black text-p-blue uppercase hover:underline"
						>← Dashboard</a
					>
					<h1 class="text-4xl font-black tracking-tighter uppercase md:text-6xl">
						Editor: <span class="text-p-blue">{course.name}</span>
					</h1>
				</div>
				<UniLink href="/courses/{course.uuid}" content="View 👁️" target="_blank" />
			</header>

			<div class="grid grid-cols-1 gap-10 lg:grid-cols-12">
				<nav class="sticky top-8 h-fit space-y-3 lg:col-span-3 xl:col-span-2">
					{#each sections as section}
						<UniButton
							onclick={() => (activeSection = section.id)}
							more_style={activeSection === section.id
								? 'translate-x-1 translate-y-1 text-white shadow-none w-full'
								: 'w-full'}
							bgcolor={activeSection === section.id ? 'bg-p-blue' : undefined}
							hv_bgcolor={activeSection === section.id ? '' : undefined}
						>
							<span>{section.icon}</span>
							{section.label}
						</UniButton>
					{/each}
				</nav>

				<main class="lg:col-span-9 xl:col-span-10">
					<div
						class="rounded-3xl border-4 border-s-black bg-white p-8 shadow-[12px_12px_0px_0px_rgba(26,26,26,1)] md:p-12"
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
									<CreateMaterial
										courseUuid={courseId}
										onchange={loadCourse}
										modules={course.modules}
									/>
								</div>

								{#if course.materials.length > 0}
									<div class="space-y-4 border-t-2 border-dashed border-gray-200 pt-6">
										<h3 class="text-xl font-black uppercase">
											Existing Materials ({course.materials.length})
										</h3>
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
								<div class="flex items-center justify-between space-x-1">
									<h2 class="text-3xl font-black text-p-blue uppercase">Examinations</h2>
									<UniButton
										onclick={() => (showCreateQuiz = !showCreateQuiz)}
										content={showCreateQuiz ? 'Cancel' : '+ New Quiz'}
									/>
								</div>

								{#if showCreateQuiz}
									<div transition:slide class="rounded-xl border-2 border-s-black bg-gray-50 p-4">
										<EditQuiz
											modules={course.modules}
											edit={false}
											courseId={course.uuid}
											onchange={onCreateQuizSubmit}
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
						{:else if activeSection === 'modules'}
							<div in:fade class="space-y-8">
								<section
									class="rounded-2xl border-4 border-s-black bg-p-blue/5 p-6 shadow-[4px_4px_0px_0px_rgba(2,87,165,1)]"
								>
									<div
										class="mb-6 flex flex-col items-start justify-between gap-4 md:flex-row md:items-center"
									>
										<div>
											<h2 class="text-xl font-black tracking-tight text-p-blue uppercase">
												Featured Module
											</h2>
											<p class="text-xs font-bold text-gray-500 uppercase">
												Direct student focus to a specific module
											</p>
										</div>

										<ModuleSelector
											modules={localModules}
											bind:selectedId={highlightedModuleId}
											label="Choose Highlight"
										/>
									</div>

									{#if activeHighlight}
										<div transition:slide class="space-y-4">
											<div
												class="flex items-center gap-4 rounded-xl border-4 border-s-black bg-white p-4 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
											>
												<span class="text-4xl">🌟</span>
												<div>
													<p class="text-xs font-black text-p-blue uppercase">Target:</p>
													<h3 class="text-xl font-black tracking-tighter uppercase">
														{activeHighlight.name}
													</h3>
												</div>
											</div>

											<div class="relative space-y-2">
												<label
													for="highlight_msg"
													class="text-[10px] font-black tracking-widest text-gray-400 uppercase"
												>
													Lecturer's Note (Explain why this is featured)
												</label>
												<div class="flex flex-col gap-3 sm:flex-row">
													<input
														id="highlight_msg"
														type="text"
														bind:value={highlightedMessage}
														placeholder="e.g., Focus on this module for this week's midterm!"
														class="flex-1 rounded-xl border-4 border-s-black bg-white p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
													/>

													<!-- {#if messageChanged} -->
													<div in:fade>
														<SuccessButton
															onclick={() => updateHighlight()}
															isSaving={isSavingHighlight}
															class="!py-3 !text-sm whitespace-nowrap"
														>
															Highlight
														</SuccessButton>
													</div>
													<!-- {/if} -->
												</div>
											</div>
										</div>
									{:else}
										<p
											class="rounded-xl border-2 border-dashed border-p-blue/30 p-4 text-center text-sm font-bold text-p-blue/50 italic"
										>
											No module currently featured. Use the selector above to pin one.
										</p>
									{/if}
								</section>

								<hr class="border-2 border-dashed border-gray-200" />

								<div class="flex items-center justify-between">
									<h2 class="text-3xl font-black text-s-black uppercase">All Modules</h2>
									<div class="flex gap-2">
										{#if orderChanged}
											<div in:fade>
												<SuccessButton
													onclick={saveModuleOrder}
													isSaving={isSavingOrder}
													class="!py-2 !text-sm"
												>
													Save New Order
												</SuccessButton>
											</div>
										{/if}
										<PrimaryButton
											onclick={() => (showCreateModule = !showCreateModule)}
											class="!py-2 !text-sm"
										>
											{showCreateModule ? 'Cancel' : '+ New Module'}
										</PrimaryButton>
									</div>
								</div>

								{#if showCreateModule}
									<div transition:slide class="rounded-xl border-4 border-s-black bg-gray-50 p-4">
										<CreateModule courseId={course.uuid} onchange={onCreateModuleSubmit} />
									</div>
								{/if}

								<div class="space-y-4">
									{#each localModules as module, i (module.uuid)}
										<div animate:flip={{ duration: 300 }} class="group relative">
											<div
												class="absolute top-1/2 -left-12 hidden -translate-y-1/2 flex-col gap-1 opacity-0 transition-opacity group-hover:opacity-100 md:flex"
											>
												<button
													disabled={i === 0}
													onclick={() => moveModule(i, 'up')}
													class="rounded border-2 border-s-black bg-white p-1 hover:bg-p-green disabled:opacity-30"
												>
													▲
												</button>
												<button
													disabled={i === localModules.length - 1}
													onclick={() => moveModule(i, 'down')}
													class="rounded border-2 border-s-black bg-white p-1 hover:bg-p-green disabled:opacity-30"
												>
													▼
												</button>
											</div>

											<div class="relative">
												{#if course.highlightedModuleId === module.uuid}
													<div
														class="absolute -top-2 -right-2 z-10 rounded-full border-2 border-s-black bg-p-blue px-2 py-0.5 text-[10px] font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
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
						{/if}
					</div>
				</main>
			</div>
		</div>
	{/if}
</div>
