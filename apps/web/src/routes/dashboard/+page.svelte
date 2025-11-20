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

<h1>Manage Courses</h1>
<br />

<div>
	<a href="/edit/courses">Create new course</a>
	<p>Click on a course to edit it</p>
</div>

<br />

{#await coursesPromise}
	<p>Loading</p>
{:then data}
	<ul>
		{#each data as course}
			<a href="/edit/courses/{course.uuid}"> {course.name} </a>
			<br />
			<br />
		{/each}
	</ul>
{:catch error}
	<p></p>
{/await}
