<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import type { Quiz, QuizSubmit, QuizMarked, Answer } from '$lib/types';
	import { slide } from 'svelte/transition';

	// Import your new button components
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import SuccessButton from '$lib/components/SuccessButton.svelte';

	let { quiz, courseId }: { quiz: Quiz; courseId: string } = $props();

	let collapsed = $state(true);
	let quizSubmit: QuizSubmit = { answers: [], id: null };
	let quizMarked: QuizMarked | null = $state(null);
	let missingAnswers: number[] = $state([]);
	let attempsCount: number = $state(quiz.attemptsCount);
	let isSubmitting = $state(false);

	let showResultsButton = auth.user?.isAdmin;

	// initialize answers
	for (let question of quiz.questions) {
		let answer: Answer = { uuid: question.uuid };
		if (question.type === 'multipleChoice') answer.selectedIndices = [];
		quizSubmit.answers.push(answer);
	}

	async function submitQuiz(e: Event) {
		e.preventDefault();
		missingAnswers = [];
		isSubmitting = true;

		quizSubmit.answers.forEach((ans, i) => {
			const isSingleMissing =
				ans.selectedIndex === undefined && quiz.questions[i].type === 'singleChoice';
			const isMultiMissing =
				(!ans.selectedIndices || ans.selectedIndices.length === 0) &&
				quiz.questions[i].type === 'multipleChoice';
			if (isSingleMissing || isMultiMissing) missingAnswers.push(i);
		});

		if (missingAnswers.length > 0) {
			isSubmitting = false;
			return;
		}

		if (auth.user) {
			quizSubmit.id = auth.user.id;
		}

		try {
			let res = await fetch(
				`/api/courses/${courseId}/modules/${quiz.moduleId}/quizzes/${quiz.uuid}/submit`,
				{
					method: 'POST',
					headers: { 'Content-type': 'application/json' },
					body: JSON.stringify(quizSubmit)
				}
			);

			if (res.ok) {
				attempsCount += 1;
				quizMarked = await res.json();
			}
		} finally {
			isSubmitting = false;
		}
	}

	function isCorrectOption(qi: number, oi: number) {
		const q = quiz.questions[qi];
		if (q.type === 'singleChoice') {
			return q.correctIndex === oi;
		} else {
			return q.correctIndices.includes(oi);
		}
	}

	function isSelectedOption(qi: number, oi: number) {
		const ans = quizSubmit.answers[qi];
		if (quiz.questions[qi].type === 'singleChoice') {
			return ans.selectedIndex === oi;
		} else {
			return ans.selectedIndices?.includes(oi);
		}
	}

	function resetState() {
		quizMarked = null;
		missingAnswers = [];
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-p-green shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all"
>
	<div class="flex w-full items-center bg-p-green hover:bg-black/5 max-md:overflow-scroll">
		<button
			type="button"
			class="flex flex-1 cursor-pointer items-center justify-between p-4 text-left"
			onclick={() => (collapsed = !collapsed)}
		>
			<div class="flex items-center gap-3">
				<span class="text-2xl">📝</span>
				<span class="text-xl font-black tracking-tight uppercase md:text-2xl">{quiz.title}</span>
				<span
					class="rounded-lg border-2 border-s-black bg-p-blue px-2 py-1 text-xs font-bold text-white uppercase"
				>
					Attempts: {attempsCount}
				</span>
			</div>
		</button>

		<div class="flex items-center gap-2 pr-4">
			{#if showResultsButton}
				<SecondaryButton
					href={`/dashboard/edit/courses/${courseId}/modules/${quiz.moduleId}/quizzes/${quiz.uuid}`}
					class="border-2 !px-3 !py-1.5 !text-xs shadow-[2px_2px_0px_0px_rgba(26,26,26,1)]"
					onclick={(e: MouseEvent | KeyboardEvent) => e.stopPropagation()}
				>
					Results 📊
				</SecondaryButton>
			{/if}

			<span
				class="pointer-events-none text-xl transition-transform {collapsed ? '' : 'rotate-180'}"
			>
				▼
			</span>
		</div>
	</div>

	{#if !collapsed}
		<div transition:slide class="border-t-4 border-s-black bg-white p-4 md:p-6">
			{#if quizMarked}
				<div
					class="mb-8 rounded-xl border-4 border-s-black bg-p-blue p-6 text-white shadow-[6px_6px_0px_0px_rgba(26,26,26,1)]"
				>
					<div class="flex items-center justify-between">
						<div>
							<h3 class="text-xs font-black tracking-widest uppercase opacity-80">Quiz Results</h3>
							<p class="text-5xl leading-none font-black">
								{quizMarked.score} / {quizMarked.maxScore}
							</p>
						</div>

						<!-- {#if showResultsButton}
							<SecondaryButton
								href={`/courses/${courseId}/quizzes/${quiz.uuid}/results`}
								class="border-2 !py-2 !text-sm"
							>
								Detailed View →
							</SecondaryButton>
						{/if} -->
					</div>
				</div>
			{/if}

			<form onsubmit={submitQuiz} class="space-y-8">
				{#each quiz.questions as q, qi}
					<div
						class="space-y-4 rounded-xl border-4 border-s-black p-5 transition-colors
                        {missingAnswers.includes(qi) ? 'border-red-500 bg-red-50' : 'bg-white'} 
                        {quizMarked?.correctPerQuestion[qi] ? 'border-p-green bg-green-50/30' : ''}"
					>
						<div class="flex items-center justify-between gap-2">
							<h3 class="text-xs font-black tracking-widest text-gray-400 uppercase">
								Question {qi + 1}
								<span class="ml-2 font-normal lowercase opacity-60"
									>({q.type === 'singleChoice' ? 'Single Choice' : 'Multiple Choice'})</span
								>
							</h3>

							{#if quizMarked}
								<span
									class="rounded-lg border-2 border-s-black px-2 py-1 text-[10px] font-black uppercase {quizMarked
										.correctPerQuestion[qi]
										? 'bg-p-green text-s-black'
										: 'bg-red-500 text-white'}"
								>
									{quizMarked.correctPerQuestion[qi] ? '✓ Correct' : '✗ Incorrect'}
								</span>
							{/if}
						</div>

						<p class="text-2xl leading-tight font-black text-s-black">{q.question}</p>

						<div class="grid gap-3">
							{#each q.options as option, oi}
								<label
									class="flex cursor-pointer items-center justify-between rounded-xl border-2 border-s-black p-4 transition-all
                                    {quizMarked
										? 'cursor-default'
										: 'hover:bg-gray-50 active:translate-y-0.5'}
                                    {quizMarked && isCorrectOption(qi, oi)
										? 'border-p-green bg-p-green/10 ring-2 ring-p-green'
										: ''}
                                    {quizMarked &&
									isSelectedOption(qi, oi) &&
									!isCorrectOption(qi, oi)
										? 'border-red-500 bg-red-50'
										: ''}
                                    {!quizMarked
										? 'has-checked:border-p-blue has-checked:bg-p-blue/5'
										: ''}"
								>
									<div class="flex items-center gap-3">
										<span class="text-lg font-bold text-s-black">{option}</span>

										{#if quizMarked && isCorrectOption(qi, oi)}
											<span
												class="rounded border border-s-black/10 bg-p-green px-2 py-0.5 text-[10px] font-black text-s-black uppercase"
											>
												Correct
											</span>
										{/if}
										{#if quizMarked && isSelectedOption(qi, oi) && !isCorrectOption(qi, oi)}
											<span
												class="rounded bg-red-500 px-2 py-0.5 text-[10px] font-black text-white uppercase"
											>
												Your Choice
											</span>
										{/if}
									</div>

									<div class="flex items-center">
										{#if q.type === 'singleChoice'}
											<input
												type="radio"
												name={`q-${qi}`}
												disabled={!!quizMarked}
												checked={quizSubmit.answers[qi].selectedIndex === oi}
												onchange={() => (quizSubmit.answers[qi].selectedIndex = oi)}
												class="h-6 w-6 cursor-pointer border-4 border-s-black accent-p-blue"
											/>
										{:else}
											<input
												type="checkbox"
												disabled={!!quizMarked}
												checked={quizSubmit.answers[qi].selectedIndices?.includes(oi)}
												onchange={(e) => {
													if (e.currentTarget.checked)
														quizSubmit.answers[qi].selectedIndices?.push(oi);
													else
														quizSubmit.answers[qi].selectedIndices = quizSubmit.answers[
															qi
														].selectedIndices?.filter((i) => i !== oi);
												}}
												class="h-6 w-6 cursor-pointer rounded border-4 border-s-black accent-p-blue"
											/>
										{/if}
									</div>
								</label>
							{/each}
						</div>
					</div>
				{/each}

				<div class="grid gap-4 pt-4">
					{#if !quizMarked}
						<SuccessButton type="submit" isSaving={isSubmitting}>Submit Results</SuccessButton>
					{:else}
						<SecondaryButton onclick={resetState} class="w-full text-xl">
							🔄 Try Again
						</SecondaryButton>
					{/if}
				</div>
			</form>

			{#if missingAnswers.length > 0}
				<div
					class="mt-6 rounded-xl border-2 border-red-500 bg-red-50 p-4 text-sm font-black tracking-tight text-red-600 uppercase"
				>
					⚠️ Missing Answers for Questions: {missingAnswers.map((i) => i + 1).join(', ')}
				</div>
			{/if}
		</div>
	{/if}
</div>
