<script lang="ts">
	import type { Material } from '$lib/types';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true);
	$inspect(material, collapsed);
</script>

<div class="space-y-2 rounded-xl bg-[#91F5AD] p-4">
	<button
		class="w-[100%] cursor-pointer text-left text-2xl font-bold"
		type="button"
		onclick={() => (collapsed = !collapsed)}
	>
		{material.name}
	</button>
	{#if !collapsed}
		<div>
			{#if material.type === 'file'}
				<a class="border border-[#1a1a1a] p-1 text-xl" href={material.fileUrl}>View</a>
				<a
					class="border border-[#1a1a1a] p-1 text-xl"
					download={material.name + material.fileUrl.split('.').pop()}
					href={material.fileUrl}>Download</a
				>
			{:else if material.type === 'url'}
				<a
					target="_blank"
					href={material.url}
					class="w-fit rounded-xl border border-[#1a1a1a] bg-[#49B3B4] p-1 text-xl text-white hover:bg-[#6DD4B1]"
				>
					View
				</a>
			{/if}
		</div>
		<div class="text-xl">
			<p>{material.description}</p>
		</div>
	{/if}
</div>
