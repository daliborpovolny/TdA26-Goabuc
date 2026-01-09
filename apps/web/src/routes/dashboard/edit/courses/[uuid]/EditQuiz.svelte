<script lang="ts">
	
	import type {
		Quiz,
		Question
	} from "$lib/types";

	function uuid() {
		return crypto.randomUUID();
	}

	const props = $props<{
		edit: boolean;
		quiz?: Quiz;
		courseId: string;
		onchange?: (quiz: Quiz) => void;
	}>();

	$inspect(props.quiz)

	const quiz = $state<Quiz>(
		props.quiz ?? {
			uuid: '',
			title: '',
			attemptsCount: 0,
			questions: []
		}
	);

	let collapsed = $state(true)

	let savedTitle = $state(quiz.title)

	async function updateQuiz(e: Event) {
		e.preventDefault();

		for (const q of quiz.questions) {
			if (q.type === 'multipleChoice' && q.correctIndices.length === 0) {
				alert('Each multiple-choice question must have at least one correct answer.');
				return;
			}
		}

		let stringifiedQuiz = JSON.stringify(quiz);

		let putRoute = `/api/courses/${props.courseId}/quizzes/${quiz.uuid}`
		let postRoute = `/api/courses/${props.courseId}/quizzes`

		let res = await fetch(props.edit ? putRoute : postRoute, {
			method: props.edit ? 'PUT' : 'POST',
			headers: { 'Content-type': 'application/json' },
			body: stringifiedQuiz
		});

		if (res.status == 201) {
			savedTitle = quiz.title
		}

		props.onchange?.(quiz);
	}

	function addQuestion(type: 'singleChoice' | 'multipleChoice') {
		let q: Question;

		if (type === 'singleChoice') {
			q = {
				uuid: "",
				type: 'singleChoice',
				question: '',
				options: [''],
				correctIndex: 0
			};
		} else {
			q = {
				uuid: "",
				type: 'multipleChoice',
				question: '',
				options: [''],
				correctIndices: []
			};
		}

		quiz.questions.push(q);
	}


	function removeQuestion(index: number) {
		quiz.questions.splice(index, 1);
	}

	function addOption(q: Question) {
		q.options.push('');
	}

	function removeOption(q: Question, index: number) {
		q.options.splice(index, 1);

		if (q.type === 'singleChoice' && q.correctIndex >= q.options.length) {
			q.correctIndex = 0;
		}

		if (q.type === 'multipleChoice') {
			q.correctIndices = q.correctIndices
				.filter(i => i !== index)
				.map(i => (i > index ? i - 1 : i));
		}

	}

	async function  deleteQuiz(e: Event) {
		e.preventDefault()

		if (!confirm('Are you sure?')) return;

		await fetch(`/api/courses/${props.courseId}/quizzes/${quiz.uuid}`, {
			method: 'DELETE'
		});

		props.onchange?.()
	}

</script>

<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">
{#if props.edit}
	
	<button type="button" class="text-2xl" onclick={() => (collapsed = !collapsed)} >{savedTitle}</button>
	<br />
{/if}

{#if !collapsed || !props.edit}
<form  onsubmit={updateQuiz}>
	<div>
		<label class="block text-sm font-medium text-gray-700 mb-1" for="title">
			Quiz title
		</label>
		<input
			type="text"
			required
			bind:value={quiz.title}
			name="title"
			class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
		/>
	</div>

	{#each quiz.questions as q, qi}
		<div class="rounded-lg border border-gray-200 p-4 space-y-4 bg-white shadow-sm">
			<div class="flex justify-between items-center">
				<h3 class="font-semibold text-lg">
					Question {qi + 1} - {q.type === 'singleChoice' ? 'Single choice' : 'Multiple choice'}
				</h3>
				<button
					type="button"
					onclick={() => removeQuestion(qi)}
					class="text-sm text-red-600 hover:text-red-800"
				>
					Remove
				</button>
			</div>

			<input
				type="text"
				placeholder="Question text"
				required
				bind:value={q.question}
				class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
			/>

			<div class="space-y-2">
				{#each q.options as _, oi}
					<div class="flex items-center gap-2">
						<input
							type="text"
							required
							placeholder={`Option ${oi + 1}`}
							bind:value={q.options[oi]}
							class="flex-1 rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>

						{#if q.type === 'singleChoice'}
							<input
								type="radio"
								name={`${qi}-single-choice-${oi}`}
								checked={q.correctIndex === oi}
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
										q.correctIndices = q.correctIndices.filter(i => i !== oi);
									}
								}}
								class="h-4 w-4"
							/>
						{/if}

						<button
							type="button"
							onclick={() => removeOption(q, oi)}
							class="text-gray-400 hover:text-red-600"
						>
							✕
						</button>
					</div>
				{/each}

				<button
					type="button"
					onclick={() => addOption(q)}
					class="text-sm text-blue-600 hover:text-blue-800"
				>
					+ Add option
				</button>
			</div>
		</div>
	{/each}

	<div class="flex gap-3 pt-4">
		<button
			type="button"
			onclick={() => addQuestion('singleChoice')}
			class="rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
		>
			+ Single choice
		</button>

		<button
			type="button"
			onclick={() => addQuestion('multipleChoice')}
			class="rounded-md bg-yellow-600 px-4 py-2 text-white hover:bg-yellow-700"
		>
			+ Multiple choice
		</button>
		
		<button 
			type="submit"
			class="rounded-md bg-green-600 px-4 py-2 text-white hover:bg-green-700"
		>
		
			Save
		</button>

		<button
			type="button"
			onclick={deleteQuiz}
			class="rounded-md bg-red-600 px-4 py-2 text-white hover:bg-red-700"
		>

		Delete
	</button>

	</div>

</form>
{/if}
</div>