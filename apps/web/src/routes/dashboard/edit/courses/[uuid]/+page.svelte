<script lang="ts">
	import { page } from '$app/state';

	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course } from '$lib/types';

	import CreateMaterial from './CreateMaterial.svelte';
	import EditMaterial from './EditMaterial.svelte';
	import EditCourse from './EditCourse.svelte';

	let courseId: string = page.params.uuid!;

	let refreshCount = $state(0);
	const reload = () => refreshCount++;

	let coursePromise = $derived.by(async () => {
		refreshCount;
		const res = await fetch(`/api/courses/${courseId}`);
		if (!res.ok) throw new Error('Failed to load course');
		return res.json();
	});
</script>

<DataLoader promise={coursePromise}>
	{#snippet children(course: Course)}
	<div class="mx-auto max-w-2xl space-y-8 p-6">
	<section>
				<EditCourse {course} onchange={reload} />
			</section>

			<hr />

			<section>
				<h1 class="text-3xl">Materials</h1>
				<br>

				<CreateMaterial courseUuid={courseId} onchange={reload} />
				<br>


				{#if course.materials.length > 0}
					
					<h1 class="mb-3 text-2xl font-semibold text-gray-800">Edit Materials</h1>

					{#each course.materials as material}
						<EditMaterial {material} courseUuid={courseId} onchange={reload} />
						<br>
					{/each}
				{/if}

				</section>
		</div>
	{/snippet}
</DataLoader>
