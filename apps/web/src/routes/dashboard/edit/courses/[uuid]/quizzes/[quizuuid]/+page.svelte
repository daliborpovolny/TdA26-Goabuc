<script lang="ts">
	import { page } from '$app/state';
	import { fade } from 'svelte/transition';
	import type { QuizOutcome } from '$lib/types';

	let outcomes = $state<QuizOutcome[]>([]);
	let loading = $state(true);

	$inspect(outcomes);

	async function loadOutcomes() {
		loading = true;
		let res = await fetch(
			`/api/courses/${page.params.uuid}/quizzes/${page.params.quizuuid}/answers`
		);
		if (res.ok) {
			outcomes = await res.json();
		}
		loading = false;
		``;
	}

	// Initial load
	loadOutcomes();

	function formatTime(timestamp: number) {
		return new Date(timestamp).toLocaleString();
	}

	function getScoreColor(score: number, max: number) {
		const percentage = (score / max) * 100;
		if (percentage >= 80) return 'text-p-green';
		if (percentage >= 50) return 'text-p-blue';
		return 'text-red-500';
	}
</script>

<div class="p-6 md:p-12">
	<header class="mb-10 flex items-center justify-between border-b-4 border-s-black pb-6">
		<div>
			<h1 class="text-5xl font-black tracking-tighter uppercase">
				Quiz <span class="text-p-blue">Outcomes</span>
			</h1>
			<p class="font-bold tracking-widest text-gray-500 uppercase">Reviewing student performance</p>
		</div>
		<button
			onclick={loadOutcomes}
			class="cursor-pointer rounded-xl border-2 border-s-black bg-white px-4 py-2 font-bold shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] hover:bg-p-green active:translate-y-0.5 active:shadow-none"
		>
			🔄 Refresh
		</button>
	</header>

	{#if loading}
		<div class="flex h-64 items-center justify-center">
			<p class="animate-pulse text-2xl font-black text-p-blue uppercase">Fetching Grades...</p>
		</div>
	{:else if outcomes.length === 0}
		<div class="rounded-2xl border-4 border-dashed border-gray-300 p-20 text-center">
			<p class="text-2xl font-bold text-gray-400">No submissions found yet.</p>
		</div>
	{:else}
		<div class="relative">
			<div class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"></div>

			<div class="relative overflow-scroll rounded-2xl border-4 border-s-black bg-white">
				<table class="w-full border-collapse text-left">
					<thead class="bg-s-black text-sm tracking-widest text-white uppercase">
						<tr>
							<th class="p-4">Student ID</th>
							<th class="p-4">Attempt</th>
							<th class="p-4">Score</th>
							<th class="p-4">Submitted At</th>
							<th class="p-4">Comment</th>
						</tr>
					</thead>
					<tbody class="font-bold text-s-black">
						{#each outcomes as outcome}
							<tr class="border-b-2 border-gray-100 transition-colors hover:bg-p-green/5">
								<td class="p-4">
									<span class="rounded bg-p-blue/10 px-2 py-1 text-p-blue">
										{outcome.user_id != 0 ? `User #${outcome.user_id}` : 'Anonymous'}
									</span>
								</td>
								<td class="p-4">
									<span class="rounded-full border-2 border-s-black bg-white px-3 py-1 text-xs">
										#{outcome.attempt_number}
									</span>
								</td>
								<td class="p-4 text-2xl font-black">
									<span class={getScoreColor(outcome.score, outcome.max_score)}>
										{outcome.score}
									</span>
									<span class="text-sm text-gray-400">/ {outcome.max_score}</span>
								</td>
								<td class="p-4 text-sm text-gray-500">
									{outcome.submitted_at}
								</td>
								<td class="p-4 text-gray-400 italic">
									{outcome.comment || '—'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>
