<script lang="ts">
	import { page } from '$app/state';

	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { Course, Material } from '$lib/types';

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
		<section>
			<EditCourse {course} onchange={reload} />
		</section>

		<hr />

		<section>
			<h3>Materials</h3>

			<CreateMaterial courseUuid={courseId} onchange={reload} />

			{#each course.materials as material}
				<EditMaterial {material} courseUuid={courseId} onchange={reload} />
			{/each}
		</section>
	{/snippet}
</DataLoader>
