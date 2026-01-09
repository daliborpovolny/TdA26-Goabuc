<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	import { auth } from '$lib/auth.svelte';
	import { onMount } from 'svelte';
	import Navbar from './Navbar.svelte';

	let { children } = $props();

	onMount(() => {
		auth.check();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#if !auth.initialized}
	<div class="splash-screen">Loading App...</div>
{:else}
	<div class="text-xl font-semibold">
		<Navbar />
	</div>

	{@render children()}
{/if}
