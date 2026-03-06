<script lang="ts">
	import { page } from '$app/state';
	import { fade, slide } from 'svelte/transition';
	import type { Quiz, Question, Module } from '$lib/types';
	import { modal } from '$lib/modal.svelte';

	// Components
	import ModuleSelector from './ModuleSelector.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import DangerButton from '$lib/components/DangerButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';

	let {
		edit,
		quiz: initialQuiz,
		courseId,
		onchange,
		modules
	} = $props<{
		edit: boolean;
		quiz?: Quiz;
		courseId: string;
		onchange?: (quiz: Quiz) => void;
		modules: Module[];
	}>();

	const quiz = $state<Quiz>(
		// svelte-ignore state_referenced_locally
		JSON.parse(
			JSON.stringify(
				initialQuiz ?? {
					uuid: '',
					title: '',
					attemptsCount: 0,
					questions: [],
					moduleId: '',
					moduleOrder: 0
				}
			)
		)
	);

	// svelte-ignore state_referenced_locally
	let collapsed = $state(edit);
	let isSaving = $state(false);
	let showSuccess = $state(false);
	let savedTitle = $state(quiz.title);
	let selectedModuleUuid = $state(quiz.moduleId || '');

	let activeModule = $derived(
		modules.find((m: Module) => m.uuid === (edit ? quiz.moduleId : selectedModuleUuid))
	);

	async function updateQuiz(e: Event) {
		e.preventDefault();

		// Validation logic
		for (const q of quiz.questions) {
			if (q.type === 'multipleChoice' && q.correctIndices.length === 0) {
				alert(`Question "${q.question}" needs at least one correct answer.`);
				return;
			}
			if (q.options.length < 2) {
				alert(`Question "${q.question}" must have at least two options.`);
				return;
			}
		}

		if (!edit && !selectedModuleUuid) {
			alert('Please select a module first.');
			return;
		}

		isSaving = true;
		const targetModuleUuid = edit ? quiz.moduleId : selectedModuleUuid;
		const url = edit
			? `/api/courses/${courseId}/modules/${targetModuleUuid}/quizzes/${quiz.uuid}`
			: `/api/courses/${courseId}/modules/${targetModuleUuid}/quizzes`;

		const res = await fetch(url, {
			method: edit ? 'PUT' : 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify(quiz)
		});

		if (res.ok) {
			// savedTitle = quiz.title;
			showSuccess = true;
			setTimeout(() => (showSuccess = false), 2000);
			onchange?.(quiz);
		}
		isSaving = false;
	}

	function addQuestion(type: 'singleChoice' | 'multipleChoice') {
		let q: Question;

		if (type === 'singleChoice') {
			q = {
				uuid: crypto.randomUUID(),
				type: 'singleChoice',
				question: '',
				options: [''],
				correctIndex: 0
			};
		} else {
			q = {
				uuid: crypto.randomUUID(),
				type: 'multipleChoice',
				question: '',
				options: [''],
				correctIndices: []
			};
		}

		quiz.questions.push(q);
	}

	function removeQuestion(index: number) {
		if (confirm('Delete this question?')) quiz.questions.splice(index, 1);
	}

	function addOption(q: Question) {
		q.options.push('');
	}

	function removeOption(q: Question, index: number) {
		q.options.splice(index, 1);
		if (q.type === 'singleChoice' && q.correctIndex >= q.options.length) q.correctIndex = 0;
		if (q.type === 'multipleChoice') {
			q.correctIndices = q.correctIndices
				.filter((i) => i !== index)
				.map((i) => (i > index ? i - 1 : i));
		}
	}

	async function deleteQuiz() {
		const confirmed = await modal.confirm(`Delete quiz "${quiz.title}"?`);
		if (!confirmed) return;

		await fetch(`/api/courses/${courseId}/modules/${quiz.moduleId}/quizzes/${quiz.uuid}`, {
			method: 'DELETE'
		});
		onchange?.(quiz);
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
>
	{#if edit}
		<button
			type="button"
			class="flex w-full cursor-pointer items-center justify-between bg-white p-4 text-left hover:bg-p-green/10"
			onclick={() => (collapsed = !collapsed)}
		>
			<div class="flex items-center gap-3">
				<span class="text-2xl">📝</span>
				<span class="text-xl font-black tracking-tight uppercase"
					>{savedTitle || 'Untitled Quiz'}</span
				>
			</div>
			<div class="flex items-center gap-4">
				{#if showSuccess}
					<span transition:fade class="text-xs font-black tracking-widest text-p-green uppercase"
						>✓ Saved</span
					>
				{/if}
				<span class="text-[10px] font-black tracking-widest text-gray-400 uppercase">
					📦 {activeModule?.name || 'Unassigned'}
				</span>
				<span class="transition-transform duration-300 {collapsed ? '' : 'rotate-180'}">▼</span>
			</div>
		</button>
	{/if}

	{#if !collapsed || !edit}
		<div transition:slide class="border-t-4 border-s-black bg-gray-50 p-4 md:p-8">
			<form onsubmit={updateQuiz} class="space-y-8">
				<div class="space-y-4">
					<div>
						<label
							class="mb-2 block text-xs font-black tracking-widest text-gray-500 uppercase"
							for="title">Quiz Title</label
						>
						<input
							type="text"
							required
							bind:value={quiz.title}
							class="w-full rounded-xl border-4 border-s-black p-3 font-bold outline-none focus:ring-4 focus:ring-p-green"
						/>
					</div>

					{#if !edit}
						<div>
							<label
								for=""
								class="mb-2 block text-xs font-black tracking-widest text-gray-500 uppercase"
								>Target Module</label
							>
							<ModuleSelector {modules} bind:selectedId={selectedModuleUuid} />
						</div>
					{/if}
				</div>

				<div class="space-y-6">
					{#each quiz.questions as q, qi}
						<div
							class="relative space-y-4 rounded-xl border-4 border-s-black bg-white p-5 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
						>
							<div class="flex items-center justify-between border-b-2 border-gray-100 pb-3">
								<span class="text-[10px] font-black tracking-widest text-p-blue uppercase">
									Q{qi + 1} — {q.type === 'singleChoice' ? 'Single Choice' : 'Multi Choice'}
								</span>
								<DangerButton
									type="button"
									onclick={() => removeQuestion(qi)}
									class="px-3! py-1! text-[10px]!"
								>
									Remove Question
								</DangerButton>
							</div>

							<input
								type="text"
								placeholder="Write your question here..."
								required
								bind:value={q.question}
								class="w-full border-b-4 border-s-black bg-transparent py-2 text-xl font-black outline-none placeholder:text-gray-300 focus:border-p-green"
							/>

							<div class="space-y-3">
								{#each q.options as _, oi}
									<div class="group flex items-center gap-3">
										<div
											class="flex flex-1 items-center gap-3 rounded-xl border-2 border-s-black bg-gray-50 p-3 transition-colors focus-within:border-p-blue focus-within:bg-white"
										>
											<input
												type="text"
												required
												placeholder={`Option ${oi + 1}`}
												bind:value={q.options[oi]}
												class="w-full flex-1 bg-transparent font-bold outline-none"
											/>

											{#if q.type === 'singleChoice'}
												<input
													type="radio"
													name={`correct-${qi}`}
													checked={q.correctIndex === oi}
													onchange={() => (q.correctIndex = oi)}
													class="h-6 w-6 cursor-pointer border-4 border-s-black accent-p-green"
												/>
											{:else}
												<input
													type="checkbox"
													checked={q.correctIndices.includes(oi)}
													onchange={(e) => {
														if (e.currentTarget.checked) q.correctIndices.push(oi);
														else q.correctIndices = q.correctIndices.filter((i) => i !== oi);
													}}
													class="h-6 w-6 cursor-pointer border-4 border-s-black accent-p-green"
												/>
											{/if}
										</div>
										<button
											type="button"
											onclick={() => removeOption(q, oi)}
											class="text-xl font-black text-red-500 transition-transform hover:scale-125"
											>✕</button
										>
									</div>
								{/each}

								<SecondaryButton
									type="button"
									onclick={() => addOption(q)}
									class="py-1! text-[10px]!"
								>
									+ Add Option
								</SecondaryButton>
							</div>
						</div>
					{/each}
				</div>

				<div class="flex flex-wrap items-center gap-4 border-t-4 border-s-black pt-8">
					<PrimaryButton onclick={() => addQuestion('singleChoice')} class="text-xs!">
						+ Single Choice
					</PrimaryButton>

					<PrimaryButton onclick={() => addQuestion('multipleChoice')} class="text-xs!">
						+ Multiple Choice
					</PrimaryButton>

					<div class="ml-auto flex gap-3">
						{#if edit}
							<SecondaryButton
								href={`/dashboard/edit/courses/${page.params.uuid}/modules/${quiz.moduleId}/quizzes/${quiz.uuid}`}
								class="text-xs!"
							>
								Results 📊
							</SecondaryButton>

							<DangerButton type="button" onclick={deleteQuiz} class="text-xs!">
								Delete
							</DangerButton>
						{/if}

						<SuccessButton type="submit" {isSaving} class="px-8! text-xs!">
							{edit ? 'Update Quiz' : 'Create Quiz'}
						</SuccessButton>
					</div>
				</div>
			</form>
		</div>
	{/if}
</div>
