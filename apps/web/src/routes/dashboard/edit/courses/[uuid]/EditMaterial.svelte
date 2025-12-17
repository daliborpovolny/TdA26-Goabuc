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

	async function remove() {
		if (!confirm('Are you sure?')) return;
		await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, { method: 'DELETE' });
		onchange(); // This triggers the reload in the parent!
	}

	async function save(e: Event) {
		// ... fetch logic ...
		onchange();
	}
</script>

<div>
	{#if material.type === 'file'}
		<form
			class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
			method="POST"
			enctype="multipart/form-data"
			onsubmit={save}
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
				<label class="text-sm font-medium text-gray-700" for="file">File</label>
				<input type="file" name="file" class="mt-1 text-gray-900" />
			</div>

			<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
				Save
			</button>

			<button class="rounded-md bg-red-800 px-4 py-2 text-white hover:bg-red-700"> Remove </button>
		</form>
	{:else if material.type === 'url'}
		<form
			class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
			method="POST"
			onsubmit={save}
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
			<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
				Save
			</button>

			<button class="rounded-md bg-red-800 px-4 py-2 text-white hover:bg-red-700"> Remove </button>
		</form>
	{:else}
		<p>Invalid material type</p>
	{/if}
</div>
