<script lang="ts">
	import type { Material } from '$lib/types';

	let {
		material,
		courseUuid,
		onchange
	}: {
		material: Material;
		courseUuid: string;
		onchange: () => void;
	} = $props();

	let collapsed = $state(true)

	async function remove(e: Event) {
		e.preventDefault();

		if (!confirm('Are you sure?')) return;
		await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, { method: 'DELETE' });
		onchange();
	}

	async function saveUrlMaterial(e: Event) {
		e.preventDefault();

		let formData = new FormData(e.target as HTMLFormElement);
		let formJson = JSON.stringify(Object.fromEntries(formData));

		await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, {
			method: 'PUT',
			headers: { 'Content-type': 'application/json' },
			body: formJson
		});

		onchange();
	}

	async function saveFileMaterial(e: Event) {
		e.preventDefault();

		let formData = new FormData(e.target as HTMLFormElement);
		await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, {
			method: 'PUT',
			body: formData
		});
		onchange();
	}
</script>

<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">

	<button type="button" class="text-2xl" onclick={() => (collapsed = !collapsed)} >{material.name}</button>

	{#if !collapsed}
		
		{#if material.type === 'file'}
			<form
				method="POST"
				enctype="multipart/form-data"
				onsubmit={saveFileMaterial}
			>
				<input type="hidden" name="type" value="file" />

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="name">Name</label>
					<input
						type="text"
						name="name"
						value={material.name}
						required
						class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					/>
				</div>

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="description">Description</label>
					<textarea
						name="description"
						value={material.description}
						class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					></textarea>
				</div>

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="file">Replace Existing File</label>
					<input type="file" name="file" class="mt-1 text-gray-900" />
				</div>

				<a class="text-sm font-medium text-gray-700 hover:text-gray-500" href={material.fileUrl}
					>View Existing File</a
				>

				<br />

				<button class="mt-5 rounded-md bg-green-800 px-4 py-2 text-white hover:bg-green-700">
					Save
				</button>

				<button onclick={remove} class="rounded-md bg-red-800 px-4 py-2 text-white hover:bg-red-700">
					Remove
				</button>
			</form>
		{:else if material.type === 'url'}
			<form
				class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
				method="POST"
				onsubmit={saveUrlMaterial}
			>
				<input type="hidden" name="type" value="url" />

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="name">Name</label>
					<input
						type="text"
						name="name"
						value={material.name}
						required
						class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					/>
				</div>

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="description">Description</label>
					<textarea
						name="description"
						value={material.description}
						class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					></textarea>
				</div>

				<div class="flex flex-col">
					<label class="text-sm font-medium text-gray-700" for="name">Url</label>
					<input
						type="text"
						name="url"
						value={material.url}
						required
						class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					/>
				</div>
				<button class="rounded-md bg-green-800 px-4 py-2 text-white hover:bg-green-700">
					Save
				</button>

				<button onclick={remove} class="rounded-md bg-red-800 px-4 py-2 text-white hover:bg-red-700">
					Remove
				</button>
			</form>
		{/if}
	{/if}
</div>
