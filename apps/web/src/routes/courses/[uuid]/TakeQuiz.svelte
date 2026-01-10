<script lang="ts">

    import type { Quiz, QuizSubmit } from "$lib/types";

    let {quiz, courseId}: {quiz: Quiz, courseId: string} = $props()

    let collapsed = $state(true)

    async function submitQuiz(e: Event) {
        e.preventDefault()
    
    
    	let stringifiedQuiz = JSON.stringify(quiz);

		let route = `/api/courses/${courseId}/quizzes/${quiz.uuid}`;

		let res = await fetch(route, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: stringifiedQuiz
		});
        console.log(res)
    }

</script>

<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">
		<button type="button" class="text-2xl" onclick={() => (collapsed = !collapsed)}
			>{quiz.title}</button
		>
		<br />

	{#if !collapsed}
		<form onsubmit={submitQuiz}>


			{#each quiz.questions as q, qi}
				<div class="space-y-4 rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
					<div class="flex items-center justify-between">
						<h3 class="text-lg font-semibold">
							Question {qi + 1} - {q.type === 'singleChoice' ? 'Single choice' : 'Multiple choice'}
						</h3>

					</div>

					<p class="w-full rounded-md border border-gray-300 px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none">
                        {q.question}
                    </p>
						
					<div class="space-y-2">
						{#each q.options as _, oi}
							<div class="flex items-center gap-2">
								<p class="flex-1 rounded-md border border-gray-300 px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none">
                                {q.options[oi]}
                                </p>

								{#if q.type === 'singleChoice'}
									<input
										type="radio"
										name={`${qi}-single-choice-${oi}`}
										onchange={() => {
											q.correctIndex = oi;
										}}
										class="h-4 w-4"
									/>
								{:else}
									<input
										type="checkbox"
										checked={q.correctIndices.includes(oi)}
										name={`${qi}-multiple-choice-${oi}`}
										onchange={(e) => {
											if (e.currentTarget.checked) {
												q.correctIndices.push(oi);
											} else {
												q.correctIndices = q.correctIndices.filter((i) => i !== oi);
											}
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
	{/if}
</div>
