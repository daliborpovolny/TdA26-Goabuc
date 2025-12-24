<script lang="ts">
	let { name, url }: { name: string; url: string } = $props();
</script>

<a class="rounded-md px-3 py-2 hover:bg-stone-200" href={url}>{name}</a>
