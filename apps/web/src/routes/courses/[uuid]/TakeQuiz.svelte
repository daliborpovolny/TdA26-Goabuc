<script lang="ts">
	import type { Quiz, QuizSubmit, QuizMarked, Answer } from '$lib/types';

	let { quiz, courseId }: { quiz: Quiz; courseId: string } = $props();

	let collapsed = $state(true);

	let quizSubmit: QuizSubmit = { answers: [] };

	$inspect(quizSubmit.answers);

	let quizMarked: QuizMarked | null = $state(null);

	let missingAnswers: number[] = $state([]);

	for (let question of quiz.questions) {
		let answer: Answer = { uuid: question.uuid };

		if (question.type === 'multipleChoice') {
			answer.selectedIndices = [];
		}

		quizSubmit.answers.push(answer);
	}

	async function submitQuiz(e: Event) {
		e.preventDefault();

		missingAnswers = [];
		quizMarked = null;

		for (let [id, answer] of quizSubmit.answers.entries()) {
			if (
				answer.selectedIndex === undefined &&
				(answer.selectedIndices === undefined || answer.selectedIndices.length === 0)
			) {
				missingAnswers.push(id);
			}
		}

		if (missingAnswers.length > 0) {
			return;
		}

		let stringifiedQuizSubmit = JSON.stringify(quizSubmit);

		let route = `/api/courses/${courseId}/quizzes/${quiz.uuid}/submit`;

		let res = await fetch(route, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: stringifiedQuizSubmit
		});

		if (res.status == 200) {
			quizMarked = await res.json();
		}
	}
</script>

<div class="space-y-2 rounded-xl bg-[#91F5AD] p-4">
	<button
		type="button"
		class="w-[100%] cursor-pointer text-left text-2xl font-bold"
		onclick={() => (collapsed = !collapsed)}>{quiz.title}</button
	>
	{#if !collapsed}
		<form onsubmit={submitQuiz} class="space-y-2">
			{#each quiz.questions as q, qi}
				<div class="space-y-4 rounded-xl border border-[#1A1A1A] bg-white p-4 shadow-sm">
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-semibold">
							Question {qi + 1} - {q.type === 'singleChoice' ? 'Single choice' : 'Multiple choice'}
							{#if missingAnswers.includes(qi)}
								<span class="text-red-500">Missing answer</span>
							{/if}

							{#if quizMarked}
								{#if quizMarked.correctPerQuestion[qi]}
									<span class="text-green-500"> Correct</span>
								{:else}
									<span class="text-red-500"> Incorrect</span>
								{/if}
							{/if}
						</h3>
					</div>

					<p
						class="w-full rounded-md border border-gray-300 px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
					>
						{q.question}
					</p>

					<div class="space-y-2">
						{#each q.options as _, oi}
							<div class="flex items-center gap-2">
								<p
									class="flex-1 rounded-md border border-gray-300 px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none"
								>
									{q.options[oi]}
								</p>

								{#if q.type === 'singleChoice'}
									<input
										type="radio"
										name={`${qi}-single-choice`}
										onchange={() => {
											// quizSubmit.answers[qi] = {
											// 	uuid: q.uuid,
											// 	comment: '',
											// 	selectedIndex: oi
											// };
											quizSubmit.answers[qi].selectedIndex = oi;
										}}
										class="h-4 w-4"
									/>
								{:else}
									<input
										type="checkbox"
										name={`${qi}-multiple-choice`}
										onchange={(e) => {
											// if (quizSubmit.answers[qi] == null) {
											// 	quizSubmit.answers[qi] = {
											// 		uuid: q.uuid,
											// 		comment: '',
											// 		selectedIndices: []
											// 	};
											// }

											if (e.currentTarget.checked) {
												quizSubmit.answers[qi].selectedIndices?.push(oi);
											} else {
												console.log('filtering');
												quizSubmit.answers[qi].selectedIndices = quizSubmit.answers[
													qi
												].selectedIndices?.filter((i) => i !== oi);
											}

											console.log(quizSubmit.answers[qi]);
										}}
										class="h-4 w-4"
									/>
								{/if}
							</div>
						{/each}
					</div>
				</div>
			{/each}

			<div class="flex gap-3 pt-4">
				<button
					type="submit"
					class="rounded-md bg-green-600 px-4 py-2 text-white hover:bg-green-700"
				>
					Submit Result
				</button>
			</div>
		</form>

		{#if missingAnswers}
			<p>
				Missing an answer at question{missingAnswers.length > 1 ? 's' : ''}: {missingAnswers.join(
					', '
				)}
			</p>
		{/if}

		{#if quizMarked}
			<p>Attempt at {quizMarked.submittedAt}</p>
			<p>Score: {quizMarked.score}/{quizMarked.maxScore}</p>
		{/if}
	{/if}
</div>
