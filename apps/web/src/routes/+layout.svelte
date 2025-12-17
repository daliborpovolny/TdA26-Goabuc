<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	import { auth } from '$lib/auth.svelte';
	import { onMount } from 'svelte';

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
	<div class="mx-3">
		<nav
			class="mx-3 flex gap-4
            [&>a]:rounded-md [&>a]:px-3
            [&>a]:py-2
            [&>a:hover]:bg-stone-200"
		>
			<a href="/">Home</a>
			<a href="/courses">Courses</a>

			{#if auth.user}
				<button class="rounded-md px-3 py-2 hover:bg-stone-200" onclick={() => auth.logout()}
					>Logout</button
				>
				<a href="/dashboard">Dashboard</a>
			{:else}
				<a href="/login">Login</a>
				<a href="/register">Register</a>
			{/if}
		</nav>
	</div>

	{@render children()}
{/if}
