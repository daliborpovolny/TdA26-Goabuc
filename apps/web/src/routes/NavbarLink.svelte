<script lang="ts">
	let { name, url }: { name: string; url: string } = $props();
</script>

<a class="rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#000] h-[66px] flex items-center" href={url}>{name}</a>
