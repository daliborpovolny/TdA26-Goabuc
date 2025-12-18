<script lang="ts">
	import type { Course } from '$lib/types';

	let { course, onchange }: { course: Course; onchange: () => void } = $props();

	async function updateCourse(e: Event) {
		e.preventDefault();

		let formData = new FormData(e.target as HTMLFormElement);
		let formJson = JSON.stringify(Object.fromEntries(formData));

		await fetch(`/api/courses/${course.uuid}`, {
			method: 'PUT',
			headers: { 'Content-type': 'application/json' },
			body: formJson
		});
		onchange();
	}
</script>

<div>
	<h2 class="mb-4 text-4xl font-semibold text-gray-800">Edit Course</h2>

	<div class="rounded-lg border border-stone-300 bg-stone-50 p-6">
		<form onsubmit={updateCourse} class="space-y-4">
			<div>
				<label class="text-sm font-medium text-gray-700" for="course_name">Name</label><br />
				<input
					type="text"
					name="name"
					value={course.name}
					class="mt-1 w-full rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				/>
			</div>

			<div>
				<label class="text-sm font-medium text-gray-700" for="course_description">Description</label
				><br />
				<textarea
					name="description"
					value={course.description}
					cols="60"
					rows="3"
					class="mt-1 w-full rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				></textarea>
			</div>

			<div class="pt-2">
				<button
					type="submit"
					class="rounded-md border border-stone-400 bg-stone-200 px-4 py-2 text-gray-800 hover:bg-stone-300"
				>
					Update
				</button>
			</div>
		</form>
	</div>
</div>
