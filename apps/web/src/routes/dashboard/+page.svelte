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

<title>Dashboard</title>

<br />

<div class="text-3xl flex justify-center w-[100%]">
	<dir class="w-fit rounded-xl bg-[#0070BB] px-5 pt-5 pb-3">
		<h1 class="text-5xl font-bold">Manage Courses</h1>
		<br />

		<div class="text-semibold">
			<a
				href="/dashboard/edit/courses"
				class="rounded-md bg-[#2592B8] px-3 py-2 text-[#FFF] hover:bg-[#6DD4B1]"
				>Create new course</a
			>
		</div>

		<br />

		{#await coursesPromise}
			<p>Loading...</p>
		{:then data}
			{#each data as course}
				<a href="/dashboard/edit/courses/{course.uuid}">
					<div class="mb-2 w-[100%] rounded-xl bg-[#91F5AD] p-5">
						<h2 class="text-medium mb-2 font-semibold underline decoration-3">{course.name}</h2>
					</div>
				</a>
			{/each}
		{:catch error}
			<p></p>
		{/await}
	</dir>
</div>
