<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte.ts';

	if (!auth.isLoggedIn) {
		goto('/login');
	}

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

<div class="ml-4">
	<title>Dashboard</title>

	<br />
	<h1>Manage Courses</h1>
	<br />

	<div>
		<a href="/dashboard/edit/courses">Create new course</a>
	</div>

	<br />

	{#await coursesPromise}
		<p>Loading</p>
	{:then data}
		<ul>
			{#each data as course}
				<a href="/dashboard/edit/courses/{course.uuid}"> {course.name} </a>
				<br />
				<br />
			{/each}
		</ul>
	{:catch error}
		<p></p>
	{/await}
</div>
