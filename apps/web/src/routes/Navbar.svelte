<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import NavbarLink from './NavbarLink.svelte';
</script>

<nav class="mx-3 flex gap-4">
	<NavbarLink name="Home" url="/" />
	<NavbarLink name="Courses" url="/courses" />

	{#if auth.user}
		<NavbarLink name="Dashboard" url="/dashboard" />
		<NavbarLink name="Profile" url="/profile" />

		<button class="rounded-md px-3 py-2 hover:bg-stone-200" onclick={() => auth.logout()}
			>Logout</button
		>
	{:else}
		<NavbarLink name="Login" url="/login" />
		<NavbarLink name="Register" url="/register" />
	{/if}
</nav>
