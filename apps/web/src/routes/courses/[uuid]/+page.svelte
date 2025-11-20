<script lang="ts">
	import { page } from '$app/state';

	let coursesPromise: Promise<any[]> = loadCourses();

	async function loadCourses() {
		return fetch('api/courses', {
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

<div>
	<h1>UUID {page.params.uuid}</h1>
</div>
