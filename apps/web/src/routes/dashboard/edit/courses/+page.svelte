<script lang="ts">
	let course_name = '';
	let course_description = '';

	let createCoursePromise: Promise<any> | null = null;

	async function createCourse() {
		createCoursePromise = fetch('/api/courses', {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: JSON.stringify({
				name: course_name,
				description: course_description
			})
		}).then(async (res) => {
			if (!res.ok) {
				let err = await res.json();
				console.log(err);
				throw new Error(err.message || 'Failed to create the course');
			}
			return res.json();
		});
	}
</script>

<div class="ml-4">
	<h1 class="text-2xl">Create new course</h1>

	<form on:submit={createCourse}>
		<label for="course_name">Name</label><br />
		<input type="text" bind:value={course_name} name="course_name" /><br />

		<label for="course_description">Description</label><br />
		<input type="text" bind:value={course_description} name="course_description" /><br />

		<button class="" type="submit">Create</button>
		{#if createCoursePromise}
			{#await createCoursePromise}
				<p>Creating the course...</p>
			{:then data}
				<p class="text-green-500">Created the course</p>
				<a class="bg-stone-400 hover:bg-stone-600" href={'/courses/' + data.uuid}
					>See the course page</a
				>
			{:catch error}
				<p class="text-red-500">{error.message}</p>
			{/await}
		{/if}
	</form>
</div>
