<script lang="ts">
	import type { Material } from '$lib/types';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true);
	$inspect(material, collapsed);
</script>

<div class="space-y-2 rounded-xl border border-black bg-[#91F5AD] p-4">
	<button class="cursor-pointer w-[100%] text-left text-2xl" type="button" onclick={() => (collapsed = !collapsed)}>
		{material.name}
	</button>
	{#if !collapsed}
		<div class="text-xl">
			<p>{material.description}</p>
		</div>
		<div>
			{#if material.type === 'file'}
				<a class="border border-[#1a1a1a] p-1 text-xl" href={material.fileUrl}>View</a>
				<a
					class="border border-[#1a1a1a] p-1 text-xl"
					download={material.name + material.fileUrl.split('.').pop()}
					href={material.fileUrl}>Download</a
				>
			{:else if material.type === 'url'}
				<a target="_blank" href={material.url} class="border border-[#1a1a1a] p-1 text-xl w-fit bg-[#49B3B4] hover:bg-[#6DD4B1]">
					View
				</a>
			{/if}
		</div>
	{/if}
</div>
