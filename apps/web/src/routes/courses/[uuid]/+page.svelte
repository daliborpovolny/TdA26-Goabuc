<script lang="ts">
	import { page } from '$app/state';

	let coursesPromise: Promise<any> = loadCourseDetail();

	async function loadCourseDetail() {
		return fetch('/api/courses/' + page.params.uuid, {
			method: 'GET',
			headers: { 'Content-type': 'application/json' }
		})
			.then(async (res) => {
				if (res.status == 404) {
					throw new Error('Unknown course');
				}

				if (!res.ok) {
					try {
						const err = await res.json();
					} catch {
						throw new Error('Failed to get course info');
					}
				}
				return res.json();
			})
			.then((data) => {
				return data;
			});
	}
</script>

<br />
{#await coursesPromise}
	<p>Loading course detail...</p>
{:then data}
	<div>
		<h1 class="text-2xl">{data.name}</h1>
		<br />
		<p>{data.description}</p>
	</div>

	{#if data.materials && data.materials.length > 0}
		<br />
		<h1 class="text-xl">Materials</h1>
		<br />

		{console.log(data)}

		{#each data.materials as material}
			{#if material.Type === 'file'}
				<div>
					<p>Name: {material.Name}</p>
					<p>Description: {material.Description}</p>
					<a class="border border-black bg-stone-200 p-1" href={material.FileUrl}>View</a>
					<a
						class="border border-black bg-stone-200 p-1"
						download={material.Name + material.FileUrl.split('.').pop()}
						href={material.FileUrl}>Download</a
					>
				</div>
			{:else if material.Type === 'url'}
				<div>
					<p>Name: {material.Name}</p>
					<p>Description: {material.Description}</p>
					<a class="border border-black bg-stone-200 p-1" href={material.Url}>View</a>
				</div>
			{:else}
				<p>Invalid material type</p>
			{/if}
			<br />
			<br />
		{/each}
	{/if}
{:catch error}
	<p class="text-red-500 capitalize">{error.message}</p>
{/await}
