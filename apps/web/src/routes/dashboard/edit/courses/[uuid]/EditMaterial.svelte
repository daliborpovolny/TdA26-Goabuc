<script lang="ts">
	import type { Material, Module } from '$lib/types';
	import { fade, slide } from 'svelte/transition';

	let {
		material,
		courseUuid,
		onchange,
		modules
	}: {
		material: Material;
		courseUuid: string;
		onchange: () => void;
		modules: Module[];
	} = $props();

	let collapsed = $state(true);
	let isSaving = $state(false);
	let showSuccess = $state(false);

	async function remove(e: Event) {
		e.preventDefault();
		if (!confirm(`Are you sure you want to delete "${material.name}"?`)) return;

		await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, { method: 'DELETE' });
		onchange();
	}

	async function handleUpdate(e: Event, type: 'file' | 'url') {
		e.preventDefault();
		isSaving = true;

		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const options: RequestInit = {
			method: 'PUT',
			body: type === 'url' ? JSON.stringify(Object.fromEntries(formData)) : formData
		};

		if (type === 'url') {
			options.headers = { 'Content-type': 'application/json' };
		}

		try {
			const res = await fetch(`/api/courses/${courseUuid}/materials/${material.uuid}`, options);
			if (res.ok) {
				showSuccess = true;
				onchange();
				setTimeout(() => (showSuccess = false), 2000);
			}
		} finally {
			isSaving = false;
		}
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all"
>
	<button
		type="button"
		class="flex w-full cursor-pointer items-center justify-between p-4 text-left hover:bg-p-green/10"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span class="text-xl">{material.type === 'file' ? '📁' : '🔗'}</span>
			<span class="text-xl font-black tracking-tight text-s-black uppercase">{material.name}</span>
			<span>
				📦: {modules.find((x: Module) => x.uuid === material.moduleId)?.name}
			</span>
		</div>

		<div class="flex items-center gap-4">
			{#if showSuccess}
				<span transition:fade class="text-xs font-bold text-p-green uppercase">✓ Saved</span>
			{/if}
			<span class="text-xl transition-transform duration-300 {collapsed ? '' : 'rotate-180'}"
				>▼</span
			>
		</div>
	</button>

	{#if !collapsed}
		<div transition:slide class="border-t-4 border-s-black bg-gray-50 p-6">
			<form
				onsubmit={(e) => handleUpdate(e, material.type as 'file' | 'url')}
				enctype={material.type === 'file' ? 'multipart/form-data' : undefined}
				class="space-y-5"
			>
				<input type="hidden" name="type" value={material.type} />

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="space-y-1">
						<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="name"
							>Display Name</label
						>
						<input
							type="text"
							name="name"
							value={material.name}
							required
							class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>

					{#if material.type === 'url'}
						<div class="space-y-1">
							<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="url"
								>Resource URL</label
							>
							<input
								type="url"
								name="url"
								value={material.url}
								required
								class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
							/>
						</div>
					{:else}
						<div class="space-y-1">
							<label class="text-xs font-black tracking-widest text-gray-500 uppercase" for="file"
								>Replace File (Optional)</label
							>
							<input
								type="file"
								name="file"
								class="w-full cursor-pointer rounded-xl border-2 border-dashed border-s-black p-2 font-bold file:mr-4 file:rounded-lg file:border-0 file:bg-s-black file:px-4 file:py-1 file:text-sm file:font-semibold file:text-white"
							/>
						</div>
					{/if}
				</div>

				<div class="space-y-1">
					<label
						class="text-xs font-black tracking-widest text-gray-500 uppercase"
						for="description">Description</label
					>
					<textarea
						name="description"
						rows="2"
						value={material.description}
						class="w-full rounded-xl border-2 border-s-black p-3 font-bold focus:ring-4 focus:ring-p-green focus:outline-none"
					></textarea>
				</div>

				<div class="flex items-center justify-between gap-3 border-t-2 border-gray-200 pt-4">
					<div>
						{#if material.type === 'file'}
							<a
								class="text-md font-bold text-p-blue hover:text-s-2"
								href={material.fileUrl}
								target="_blank"
							>
								<div
									class="rounded-lg border-2 border-s-black bg-p-blue p-2 shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
								>
									👁️
								</div>
							</a>
						{/if}
					</div>

					<div class="flex gap-3">
						<button
							type="button"
							onclick={remove}
							class="cursor-pointer rounded-lg border-2 border-s-black bg-red-500 px-4 py-2 text-xs font-black text-white uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none"
						>
							Delete
						</button>

						<button
							type="submit"
							disabled={isSaving}
							class="cursor-pointer rounded-lg border-2 border-s-black bg-p-green px-6 py-2 text-xs font-black text-s-black uppercase shadow-[2px_2px_0px_0px_rgba(26,26,26,1)] active:translate-y-0.5 active:shadow-none disabled:opacity-50"
						>
							{isSaving ? 'Saving...' : 'Save Changes'}
						</button>
					</div>
				</div>
			</form>
		</div>
	{/if}
</div>
