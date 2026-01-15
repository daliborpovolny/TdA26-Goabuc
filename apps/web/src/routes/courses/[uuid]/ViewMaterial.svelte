<script lang="ts">
	import type { Material } from '$lib/types';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true); // Default to collapsed for a cleaner list
</script>

<div
	class="overflow-hidden rounded-xl border-2 border-s-black bg-p-green shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] transition-all"
>
	<button
		class="flex w-full items-center justify-between p-4 text-left transition-colors hover:bg-black/5"
		type="button"
		onclick={() => (collapsed = !collapsed)}
	>
		<div class="flex items-center gap-3">
			<span class="text-2xl">
				{material.type === 'file' ? '📁' : '🔗'}
			</span>
			<span class="text-xl font-bold md:text-2xl">{material.name}</span>
		</div>

		<span class="text-xl transition-transform duration-300 {collapsed ? '' : 'rotate-180'}">
			▼
		</span>
	</button>

	{#if !collapsed}
		<div class="space-y-4 border-t-2 border-s-black bg-white p-4">
			{#if material.description}
				<p class="text-lg leading-relaxed text-s-black/80">
					{material.description}
				</p>
			{/if}

			<div class="flex flex-wrap gap-3">
				{#if material.type === 'file'}
					<a
						href={material.fileUrl}
						class="flex items-center gap-2 rounded-lg border-2 border-s-black bg-s-2 px-4 py-2 font-bold text-white transition-all hover:bg-s-1 hover:text-s-black active:translate-y-1"
					>
						<span>👁️</span> View File
					</a>
					<a
						download={material.name}
						href={material.fileUrl}
						class="flex items-center gap-2 rounded-lg border-2 border-s-black bg-p-blue px-4 py-2 font-bold text-white transition-all hover:opacity-90 active:translate-y-1"
					>
						<span>📥</span> Download
					</a>
				{:else if material.type === 'url'}
					<a
						target="_blank"
						href={material.url}
						class="flex items-center gap-2 rounded-lg border-2 border-s-black bg-s-2 px-4 py-2 font-bold text-white transition-all hover:bg-s-1 hover:text-s-black active:translate-y-1"
					>
						<span>🌐</span> Open Link
					</a>
				{/if}
			</div>
		</div>
	{/if}
</div>
