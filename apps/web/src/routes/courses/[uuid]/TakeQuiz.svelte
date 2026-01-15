<script lang="ts">
	import type { Quiz, QuizSubmit, QuizMarked, Answer } from '$lib/types';
	import { slide } from 'svelte/transition';

	let { quiz, courseId }: { quiz: Quiz; courseId: string } = $props();

	let collapsed = $state(true);
	let quizSubmit: QuizSubmit = { answers: [] };
	let quizMarked: QuizMarked | null = $state(null);
	let missingAnswers: number[] = $state([]);
	let attempsCount: number = $state(quiz.attemptsCount);

	// initialize answers
	for (let question of quiz.questions) {
		let answer: Answer = { uuid: question.uuid };
		if (question.type === 'multipleChoice') answer.selectedIndices = [];
		quizSubmit.answers.push(answer);
	}

	async function submitQuiz(e: Event) {
		e.preventDefault();
		missingAnswers = [];

		quizSubmit.answers.forEach((ans, i) => {
			const isSingleMissing =
				ans.selectedIndex === undefined && quiz.questions[i].type === 'singleChoice';
			const isMultiMissing =
				(!ans.selectedIndices || ans.selectedIndices.length === 0) &&
				quiz.questions[i].type === 'multipleChoice';
			if (isSingleMissing || isMultiMissing) missingAnswers.push(i);
		});

		if (missingAnswers.length > 0) return;

		let res = await fetch(`/api/courses/${courseId}/quizzes/${quiz.uuid}/submit`, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify(quizSubmit)
		});

		if (res.ok) {
			attempsCount += 1;
			quizMarked = await res.json();
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
	class="overflow-hidden rounded-xl border-2 border-s-black bg-p-green shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
>
	<button
		type="button"
		class="flex w-full items-center justify-between p-4 text-left hover:bg-black/5"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span class="text-2xl">📝</span>
			<span class="text-xl font-bold md:text-2xl">{quiz.title}</span>
			<span class="text-xl">Attempts Taken: {attempsCount}</span>
		</div>
		<span class="text-xl transition-transform {collapsed ? '' : 'rotate-180'}">▼</span>
	</button>

	{#if !collapsed}
		<div transition:slide class="border-t-2 border-s-black bg-white p-4 md:p-6">
			{#if quizMarked}
				<div
					class="mb-8 rounded-xl border-2 border-s-black bg-p-blue p-4 text-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
				>
					<h3 class="text-2xl font-bold">Quiz Results</h3>
					<p class="text-4xl font-black">{quizMarked.score} / {quizMarked.maxScore}</p>
					<p class="text-sm opacity-90">Submitted at: {quizMarked.submittedAt}</p>
				</div>
			{/if}

			<form onsubmit={submitQuiz} class="space-y-6">
				{#each quiz.questions as q, qi}
					<div
						class="space-y-4 rounded-xl border-2 border-s-black p-4 transition-colors
                        {missingAnswers.includes(qi) ? 'border-red-500 bg-red-50' : 'bg-white'} 
                        {quizMarked?.correctPerQuestion[qi] ? 'border-p-green bg-green-50' : ''}"
					>
						<div class="flex items-center justify-between gap-2">
							<h3 class="font-bold tracking-tight text-s-black uppercase">
								Question {qi + 1}
								<span class="text-xs font-normal opacity-60"
									>— {q.type === 'singleChoice' ? 'Single' : 'Multiple'}</span
								>
							</h3>

							{#if quizMarked}
								<span
									class="rounded-lg px-2 py-1 text-xs font-bold uppercase {quizMarked
										.correctPerQuestion[qi]
										? 'bg-p-green'
										: 'bg-red-500 text-white'}"
								>
									{quizMarked.correctPerQuestion[qi] ? 'Correct' : 'Incorrect'}
								</span>
							{/if}
						</div>

						<p class="text-xl leading-tight font-semibold">{q.question}</p>

						<div class="grid gap-2">
							{#each q.options as option, oi}
								<label
									class="flex cursor-pointer items-center justify-between rounded-lg border-2 border-s-black p-3 transition-all
            {quizMarked ? 'cursor-default' : 'hover:bg-gray-50'}
            {quizMarked && isCorrectOption(qi, oi)
										? 'border-p-green bg-p-green/20 ring-2 ring-p-green'
										: ''}
            {quizMarked && isSelectedOption(qi, oi) && !isCorrectOption(qi, oi)
										? 'border-red-500 bg-red-100'
										: ''}
            {!quizMarked ? 'has-checked:border-p-blue has-checked:bg-p-blue/10' : ''}"
								>
									<div class="flex items-center gap-2">
										<span class="text-lg font-medium">{option}</span>

										{#if quizMarked && isCorrectOption(qi, oi)}
											<span
												class="rounded bg-p-green px-1.5 py-0.5 text-[10px] font-black text-s-black uppercase"
											>
												Correct Answer
											</span>
										{/if}

										{#if quizMarked && isSelectedOption(qi, oi) && !isCorrectOption(qi, oi)}
											<span
												class="rounded bg-red-500 px-1.5 py-0.5 text-[10px] font-black text-white uppercase"
											>
												Your Choice
											</span>
										{/if}
									</div>

									{#if q.type === 'singleChoice'}
										<input
											type="radio"
											name={`q-${qi}`}
											disabled={!!quizMarked}
											checked={quizSubmit.answers[qi].selectedIndex === oi}
											onchange={() => (quizSubmit.answers[qi].selectedIndex = oi)}
											class="h-5 w-5 border-2 border-s-black accent-p-blue"
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
											class="h-5 w-5 rounded border-2 border-s-black accent-p-blue"
										/>
									{/if}
								</label>
							{/each}
						</div>
					</div>
				{/each}

				<div class="pt-4">
					{#if !quizMarked}
						<button
							type="submit"
							class="w-full rounded-xl border-2 border-s-black bg-p-green py-4 text-2xl font-black tracking-widest uppercase shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none active:bg-s-1"
						>
							Submit Results
						</button>
					{:else}
						<button
							type="button"
							onclick={() => resetState()}
							class="w-full rounded-xl border-2 border-s-black bg-s-2 py-4 text-xl font-bold text-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] hover:shadow-none"
						>
							Try Again
						</button>
					{/if}
				</div>
			</form>

			{#if missingAnswers.length > 0}
				<div class="mt-4 rounded-lg bg-red-100 p-3 font-bold text-red-700">
					⚠️ Please answer question{missingAnswers.length > 1 ? 's' : ''}:
					{missingAnswers.map((i) => i + 1).join(', ')}
				</div>
			{/if}
		</div>
	{/if}
</div>
