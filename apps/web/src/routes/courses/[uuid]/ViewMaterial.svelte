<script lang="ts">
	import type { Material } from '$lib/types';
	import { slide } from 'svelte/transition';

	import PrimaryButton from '$lib/components/PrimaryButton.svelte';
	import SecondaryButton from '$lib/components/SecondaryButton.svelte';
	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true);

	function getFavicon(url: string) {
		try {
			const domain = new URL(url).hostname;
			return `https://www.google.com/s2/favicons?domain=${domain}&sz=64`;
		} catch {
			return `https://www.google.com/s2/favicons?domain=google.com&sz=64`;
		}
	}

	async function incrementAccessedCounter() {
		await fetch(
			`/api/courses/${page.params.uuid}/modules/${material.moduleId}/materials/${material.uuid}/increment`,
			{ method: 'POST' }
		);
	}
</script>

<div
	class="overflow-hidden rounded-xl border-4 border-s-black bg-p-green shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all"
>
	<button
		class="flex w-full cursor-pointer items-center justify-between p-4 text-left transition-colors hover:bg-black/5"
		type="button"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span
				class="flex h-12 w-12 items-center justify-center rounded-xl border-4 border-s-black bg-white text-2xl shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]"
			>
				{#if material.type === 'file'}
					📁
				{:else if material.type === 'url'}
					<img src={getFavicon(material.url)} alt="site icon" class="h-7 w-7 rounded-sm" />
				{/if}
			</span>
			<span class="text-xl font-black tracking-tight uppercase md:text-2xl">{material.name}</span>
			{#if auth.user?.isAdmin}
				<span
					class="rounded-lg border-2 border-s-black bg-p-blue px-2 py-1 text-xs font-bold text-white uppercase"
				>
					Times Accessed: {material.timesAccessed}
				</span>
			{/if}
		</div>

		<span class="text-xl transition-transform duration-300 {collapsed ? '' : 'rotate-180'}">
			▼
		</span>
	</button>

	{#if !collapsed}
		<div transition:slide class="space-y-6 border-t-4 border-s-black bg-white p-6">
			{#if material.description}
				<p class="text-lg leading-relaxed font-bold text-s-black/80 italic">
					{material.description}
				</p>
			{/if}

			<!-- This is used for stats tracking -> no need to do any aria stuff -->

			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div class="flex flex-wrap gap-4" onclick={incrementAccessedCounter}>
				{#if material.type === 'file'}
					<SecondaryButton href={material.fileUrl} target="_blank" class="text-sm! md:text-base!">
						<span>👁️</span> View File
					</SecondaryButton>

					<PrimaryButton
						href={material.fileUrl}
						download={material.name}
						class="!md:text-base! text-sm"
					>
						<span>📥</span> Download
					</PrimaryButton>
				{:else if material.type === 'url'}
					<PrimaryButton href={material.url} target="_blank" class="text-sm! md:text-base!">
						<span>🌐</span> Open Link
					</PrimaryButton>
				{/if}
			</div>
		</div>
	{/if}
</div>
