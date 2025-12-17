<script lang="ts">
	let coursesPromise: Promise<any[]> = loadCourses();

	async function loadCourses() {
		return fetch('/api/courses', {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		})
			.then(async (res) => {
				if (!res.ok) {
					try {
						const err = await res.json();
					} catch {
						throw new Error('Failed to get list of courses');
					}
				}
				return res.json();
			})
			.then((data) => {
				return data;
			});
	}
</script>

<h1 class="ml-2">Courses</h1>
<br />

{#await coursesPromise}
	<p>Loading</p>
{:then data}
	<ul class="ml-4">
		{#each data as course}
			<a href="/courses/{course.uuid}"> {course.name} </a>
			<p>{course.description}</p>
			<br />
		{/each}
	</ul>
{:catch error}
	<p class="text-red-500 capitalize">{error.message}</p>
{/await}
