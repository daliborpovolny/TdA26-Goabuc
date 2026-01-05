<script lang="ts">
	let { name, url }: { name: string; url: string } = $props();
</script>

<a class="rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#000]" href={url}>{name}</a>
