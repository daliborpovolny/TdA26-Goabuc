<title>Dashboard</title>

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

<div class="ml-4 text-3xl">

	<br />
	<h1 class="text-5xl font-bold">Manage Courses</h1>
	<br />

	<div class="text-semibold">
		<a href="/dashboard/edit/courses" class="px-3 py-2 rounded-md bg-[#2592B8] hover:bg-[#6DD4B1] text-[#FFF]">Create new course</a>
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
