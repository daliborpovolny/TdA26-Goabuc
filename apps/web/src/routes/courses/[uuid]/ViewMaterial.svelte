<script lang="ts">
	import type { Material } from '$lib/types';

	let { material }: { material: Material } = $props();

	let collapsed = $state(true)
</script>

<div>
	{#if !collapsed}
		<button type="button" onclick={() => (collapsed = !collapsed)} >{material.name}</button>

		{#if material.type === 'file'}
			<div>
				<p>{material.description}</p>
				<a class="border border-black bg-stone-200 p-1" href={material.fileUrl}>View</a>
				<a
					class="border border-black bg-stone-200 p-1"
					download={material.name + material.fileUrl.split('.').pop()}
					href={material.fileUrl}>Download</a
				>
			</div>
		{:else if material.type === 'url'}
			<div>
				<p>{material.description}</p>
				<a class="border border-black bg-stone-200 p-1" href={material.url}>View</a>
			</div>
		{/if}
	{/if}
	</div>
