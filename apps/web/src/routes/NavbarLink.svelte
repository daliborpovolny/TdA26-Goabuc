<script lang="ts">
	let { name, url, onclick }: { name: string; url: string; onclick?: () => void } = $props();
</script>

<a
	{onclick}
	class="flex h-[61px] items-center rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#1A1A1A]"
	href={url}>{name}</a
>
