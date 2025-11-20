<script lang="ts">
	let coursesPromise: Promise<any[]> = loadCourses();

	async function loadCourses() {
		return fetch('/api/courses', {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		})
			.then(async (res) => {
				if (!res.ok) {
					const err = await res.json();
					throw new Error(err.message || 'Login failed');
				}
				return res.json();
			})
			.then((data) => {
				return data;
			});
	}
</script>

<h1>Courses</h1>
<br />

{#await coursesPromise}
	<p>Loading</p>
{:then data}
	<ul>
		{#each data as course}
			<a href="/courses/{course.uuid}"> {course.name} </a>
			<p>{course.description}</p>
			<br />
		{/each}
	</ul>
{:catch error}
	<p></p>
{/await}
