<script lang="ts">
	let { courseUuid, onchange }: { courseUuid: string; onchange: () => void } = $props();

	let materialType: 'file' | '' | 'url' = $state('');

	async function createUrlMaterial(e: Event) {
		e.preventDefault();

		let formData = new FormData(e.target as HTMLFormElement);
		let formJson = JSON.stringify(Object.fromEntries(formData));

		await fetch(`/api/courses/${courseUuid}/materials`, {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			body: formJson
		});

		materialType = ''
		onchange();
	}

	async function createFileMaterial(e: Event) {
		e.preventDefault();

		let formData = new FormData(e.target as HTMLFormElement);
		await fetch(`/api/courses/${courseUuid}/materials`, {
			method: 'POST',
			body: formData
		});

		materialType = ''
		onchange();
	}
</script>

<div>
	<div class="flex gap-4">
		<button
			onclick={() => (materialType = materialType === 'url' || materialType == '' ? 'file' : '')}
			class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200 {materialType ===
			'file'
				? 'bg-stone-200'
				: ''}"
		>
			New File Material
		</button>
		<button
			onclick={() => (materialType = materialType === 'file' || materialType == '' ? 'url' : '')}
			class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200 {materialType ===
			'url'
				? 'bg-stone-200'
				: ''}"
		>
			New Url Material
		</button>
	</div>

	<br />

	{#if materialType === 'file'}
		<form
			class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
			method="POST"
			enctype="multipart/form-data"
			onsubmit={createFileMaterial}
		>
			<!-- required by API -->
			<input type="hidden" name="type" value="file" />

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="name">Name</label>
				<input
					type="text"
					name="name"
					required
					class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				/>
			</div>

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="description">Description</label>
				<textarea
					name="description"
					class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				></textarea>
			</div>

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="file">File</label>
				<input type="file" name="file" required class="mt-1 text-gray-900" />
			</div>

			<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
				Create
			</button>
		</form>
	{:else if materialType === 'url'}
		<form
			method="POST"
			onsubmit={createUrlMaterial}
			class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
		>
			<input type="hidden" name="type" value="url" />

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="name">Name</label>
				<input
					type="text"
					name="name"
					class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				/>
			</div>

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="description">Description</label>
				<textarea
					name="description"
					class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				></textarea>
			</div>

			<div class="flex flex-col">
				<label class="text-sm font-medium text-gray-700" for="url">URL</label>
				<input
					type="url"
					name="url"
					class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
				/>
			</div>

			<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
				Create
			</button>
		</form>
	{/if}
</div>
