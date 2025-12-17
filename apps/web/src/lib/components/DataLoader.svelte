<script lang="ts">
	let { promise, children } = $props();
</script>

{#await promise}
	<div class="loading">Loading data...</div>
{:then value}
	{@render children(value)}
{:catch error}
	<div class="error">
		<p>Something went wrong: {error.message}</p>
		<button onclick={() => location.reload()}>Retry</button>
	</div>
{/await}
