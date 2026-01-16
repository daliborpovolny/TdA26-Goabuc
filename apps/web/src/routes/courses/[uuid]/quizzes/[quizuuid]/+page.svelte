<script lang="ts">
	import { page } from '$app/state';
	import type { QuizOutcome } from '$lib/types';

	let outcomes: QuizOutcome[] = getOutcomes();

	// pasted from lib for reference

	async function getOutcomes() {
		let res = await fetch(
			`/api/courses/${page.params.uuid}/quizzes/${page.params.quizuuid}/answers`
		);

		if (!res.ok) {
			return [];
		}

		let data: QuizOutcome[] = await res.json();
		return data;
	}
</script>

<div>
	{#each outcomes as outcome}
		<div>
			<p>{outcome.user_id ? outcome.user_id : 'anonymous'}</p>
			<p></p>
		</div>
	{/each}
</div>
