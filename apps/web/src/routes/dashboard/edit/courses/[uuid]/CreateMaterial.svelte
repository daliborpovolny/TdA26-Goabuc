<script lang="ts">
	import { fade, slide } from 'svelte/transition';

	let { courseUuid, onchange }: { courseUuid: string; onchange: () => void } = $props();

	let materialType: 'file' | '' | 'url' = $state('');
	let isSaving = $state(false);
	let showSuccess = $state(false);

	async function handleUpload(e: Event, type: 'file' | 'url') {
		e.preventDefault();
		isSaving = true;

		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const options: RequestInit = {
			method: 'POST',
			body: type === 'url' ? JSON.stringify(Object.fromEntries(formData)) : formData
		};

		if (type === 'url') {
			options.headers = { 'Content-type': 'application/json' };
		}

		try {
			const res = await fetch(`/api/courses/${courseUuid}/materials`, options);
			if (res.ok) {
				materialType = '';
				showSuccess = true;
				onchange();
				setTimeout(() => (showSuccess = false), 2000);
			}
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="space-y-4">
	<div class="flex flex-wrap gap-4">
		<button
			onclick={() => (materialType = materialType === 'file' ? '' : 'file')}
			class="group relative flex items-center gap-2 rounded-xl border-4 border-s-black px-6 py-3 font-black tracking-widest uppercase transition-all
            {materialType === 'file'
				? 'translate-x-1 translate-y-1 bg-p-blue text-white shadow-none'
				: 'bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] hover:bg-p-green'}"
		>
			<span>📁</span> File
		</button>

		<button
			onclick={() => (materialType = materialType === 'url' ? '' : 'url')}
			class="group relative flex items-center gap-2 rounded-xl border-4 border-s-black px-6 py-3 font-black tracking-widest uppercase transition-all
            {materialType === 'url'
				? 'translate-x-1 translate-y-1 bg-p-blue text-white shadow-none'
				: 'bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] hover:bg-p-green'}"
		>
			<span>🔗</span> Link
		</button>

		{#if showSuccess}
			<div
				transition:fade
				class="flex items-center font-bold tracking-tighter text-p-green uppercase"
			>
				✓ Material Added to Course
			</div>
		{/if}
	</div>

	{#if materialType !== ''}
		<div
			transition:slide
			class="rounded-2xl border-4 border-s-black bg-white p-6 shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
		>
			<form
				onsubmit={(e) => handleUpload(e, materialType as 'file' | 'url')}
				class="space-y-4"
				enctype={materialType === 'file' ? 'multipart/form-data' : undefined}
			>
				<input type="hidden" name="type" value={materialType} />

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="space-y-1">
						<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="name"
							>Display Name</label
						>
						<input
							type="text"
							name="name"
							required
							placeholder={materialType === 'file' ? 'Lecture_Notes.pdf' : 'Useful Resource Name'}
							class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>

					{#if materialType === 'url'}
						<div class="space-y-1">
							<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="url"
								>URL Address</label
							>
							<input
								type="url"
								name="url"
								required
								placeholder="https://..."
								class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
							/>
						</div>
					{:else}
						<div class="space-y-1">
							<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="file"
								>Select File</label
							>
							<input
								type="file"
								name="file"
								required
								class="w-full cursor-pointer rounded-xl border-2 border-dashed border-s-black p-2 font-bold file:mr-4 file:rounded-lg file:border-0 file:bg-s-black file:px-4 file:py-1 file:text-sm file:font-semibold file:text-white"
							/>
						</div>
					{/if}
				</div>

				<div class="space-y-1">
					<label
						class="text-xs font-black tracking-widest text-gray-500 uppercase"
						for="description">Brief Description (Optional)</label
					>
					<textarea
						name="description"
						rows="2"
						class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
					></textarea>
				</div>

				<div class="flex justify-end pt-2">
					<button
						type="submit"
						disabled={isSaving}
						class="rounded-xl border-4 border-s-black bg-p-green px-8 py-2 text-lg font-black tracking-widest uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none disabled:opacity-50"
					>
						{isSaving ? 'Uploading...' : 'Confirm Creation'}
					</button>
				</div>
			</form>
		</div>
	{/if}
</div>
