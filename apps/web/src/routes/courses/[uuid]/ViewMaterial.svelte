<script lang="ts">
	import type { Material } from '$lib/types';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true);
	$inspect(material, collapsed);
</script>

<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">
	<button type="button" onclick={() => (collapsed = !collapsed)}>{material.name}</button>

	{#if !collapsed}
		<div>
			<p>{material.description}</p>
			<br />
			{#if material.type === 'file'}
				<a class="border border-black bg-stone-200 p-1" href={material.fileUrl}>View</a>
				<a
					class="border border-black bg-stone-200 p-1"
					download={material.name + material.fileUrl.split('.').pop()}
					href={material.fileUrl}>Download</a
				>
			{:else if material.type === 'url'}
				<a class="border border-black bg-stone-200 p-1" href={material.url}>View</a>
			{/if}
		</div>
	{/if}
</div>
