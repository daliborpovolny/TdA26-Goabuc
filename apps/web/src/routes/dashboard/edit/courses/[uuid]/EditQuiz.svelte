<script lang="ts">
	import { page } from '$app/state';
	import type { Quiz, Question } from '$lib/types';
	import { fade, slide } from 'svelte/transition';

	const props = $props<{
		edit: boolean;
		quiz?: Quiz;
		courseId: string;
		onchange?: (quiz: Quiz) => void;
	}>();

	const quiz = $state<Quiz>(
		props.quiz ?? {
			uuid: '',
			title: '',
			attemptsCount: 0,
			questions: []
		}
	);

	let collapsed = $state(props.edit); // Default open for new quizzes, collapsed for existing
	let isSaving = $state(false);
	let showSuccess = $state(false);
	let savedTitle = $state(quiz.title);

	async function updateQuiz(e: Event) {
		e.preventDefault();

		// Validation for multiple choice
		for (const q of quiz.questions) {
			if (q.type === 'multipleChoice' && q.correctIndices.length === 0) {
				alert(`Question "${q.question}" needs at least one correct answer.`);
				return;
			}

			if (q.options.length < 2) {
				alert(`Question "${q.question}" must have at least two options to choose`);
			}
		}

		isSaving = true;
		const putRoute = `/api/courses/${props.courseId}/quizzes/${quiz.uuid}`;
		const postRoute = `/api/courses/${props.courseId}/quizzes`;

		const res = await fetch(props.edit ? putRoute : postRoute, {
			method: props.edit ? 'PUT' : 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify(quiz)
		});

		if (res.ok) {
			savedTitle = quiz.title;
			showSuccess = true;
			setTimeout(() => (showSuccess = false), 2000);
			props.onchange?.(quiz);
		}
		isSaving = false;
	}

	function addQuestion(type: 'singleChoice' | 'multipleChoice') {
		const q: Question =
			type === 'singleChoice'
				? {
						uuid: crypto.randomUUID(),
						type: 'singleChoice',
						question: '',
						options: [''],
						correctIndex: 0
					}
				: {
						uuid: crypto.randomUUID(),
						type: 'multipleChoice',
						question: '',
						options: [''],
						correctIndices: []
					};
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
		if (!confirm('Permanently delete this entire quiz?')) return;
		await fetch(`/api/courses/${props.courseId}/quizzes/${quiz.uuid}`, { method: 'DELETE' });
		props.onchange?.(quiz);
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
>
	{#if props.edit}
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
					<span transition:fade class="text-xs font-bold text-p-green uppercase">✓ Saved</span>
				{/if}
				<span class="transition-transform {collapsed ? '' : 'rotate-180'}">▼</span>
			</div>
		</button>
	{/if}

	{#if !collapsed || !props.edit}
		<div transition:slide class="border-t-4 border-s-black bg-gray-50 p-4 md:p-6">
			<form onsubmit={updateQuiz} class="space-y-8">
				<div class="space-y-2">
					<label class="block text-sm font-black tracking-widest text-s-black uppercase" for="title"
						>Quiz Title</label
					>
					<input
						type="text"
						required
						bind:value={quiz.title}
						class="w-full rounded-xl border-4 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div class="space-y-6">
					{#each quiz.questions as q, qi}
						<div
							class="relative space-y-4 rounded-xl border-4 border-s-black bg-white p-4 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
						>
							<div class="flex items-center justify-between border-b-2 border-gray-100 pb-2">
								<span class="text-xs font-black tracking-widest text-p-blue uppercase">
									Question {qi + 1} — {q.type === 'singleChoice' ? 'Single' : 'Multi'}
								</span>
								<button
									type="button"
									onclick={() => removeQuestion(qi)}
									class="cursor-pointer text-xs font-bold text-red-500 hover:underline"
								>
									Remove Question
								</button>
							</div>

							<input
								type="text"
								placeholder="What is the question?"
								required
								bind:value={q.question}
								class="w-full border-b-2 border-s-black bg-transparent py-2 text-xl font-bold focus:border-p-green focus:outline-none"
							/>

							<div class="space-y-2">
								{#each q.options as _, oi}
									<div class="group flex items-center gap-3">
										<div
											class="flex flex-1 items-center gap-2 rounded-lg border-2 border-s-black bg-white p-2 focus-within:bg-p-green/5"
										>
											<input
												type="text"
												required
												placeholder={`Option ${oi + 1}`}
												bind:value={q.options[oi]}
												class="flex-1 bg-transparent font-medium focus:outline-none"
											/>

											{#if q.type === 'singleChoice'}
												<input
													type="radio"
													name={`correct-${qi}`}
													checked={q.correctIndex === oi}
													onchange={() => (q.correctIndex = oi)}
													class="h-5 w-5 cursor-pointer accent-p-green"
												/>
											{:else}
												<input
													type="checkbox"
													checked={q.correctIndices.includes(oi)}
													onchange={(e) => {
														if (e.currentTarget.checked) q.correctIndices.push(oi);
														else q.correctIndices = q.correctIndices.filter((i) => i !== oi);
													}}
													class="h-5 w-5 cursor-pointer accent-p-green"
												/>
											{/if}
										</div>
										<button
											type="button"
											onclick={() => removeOption(q, oi)}
											class="cursor-pointer px-1 font-bold text-red-500 opacity-0 group-hover:opacity-100"
											>✕</button
										>
									</div>
								{/each}

								<button
									type="button"
									onclick={() => addOption(q)}
									class="cursor-pointer text-xs font-black tracking-widest text-p-blue uppercase hover:text-s-2"
								>
									+ Add Option
								</button>
							</div>
						</div>
					{/each}
				</div>

				<div class="flex flex-wrap gap-3 border-t-4 border-s-black pt-6">
					<button
						type="button"
						onclick={() => addQuestion('singleChoice')}
						class="cursor-pointer rounded-lg border-2 border-s-black bg-s-3 px-4 py-2 text-xs font-black text-white uppercase hover:bg-s-2"
					>
						+ Single Choice
					</button>
					<button
						type="button"
						onclick={() => addQuestion('multipleChoice')}
						class="cursor-pointer rounded-lg border-2 border-s-black bg-s-3 px-4 py-2 text-xs font-black text-white uppercase hover:bg-s-2"
					>
						+ Multiple Choice
					</button>

					<div class="ml-auto flex gap-3">

						<a
							href={"/dashboard/edit/courses/" + page.params.uuid + "/quizzes/" + quiz.uuid}
							class="rounded-lg border-2 border-s-black bg-s-2 px-8 py-2 text-xs font-black text-s-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
						>
							Results
					</a>

						<!-- <button
							type="submit"
							disabled={isSaving}
							class="rounded-lg border-2 border-s-black bg-p-green px-8 py-2 text-xs font-black text-s-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
						>
							{isSaving ? 'Saving...' : 'Save Quiz'}
						</button> -->

						{#if props.edit}
							<button
								type="button"
								onclick={deleteQuiz}
								class="cursor-pointer rounded-lg border-2 border-s-black bg-red-500 px-4 py-2 text-xs font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
							>
								Delete Quiz
							</button>
						{/if}
						<button
							type="submit"
							disabled={isSaving}
							class="cursor-pointer rounded-lg border-2 border-s-black bg-p-green px-8 py-2 text-xs font-black text-s-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
						>
							{isSaving ? 'Saving...' : 'Save Quiz'}
						</button>
					</div>
				</div>
			</form>
		</div>
	{/if}
</div>
