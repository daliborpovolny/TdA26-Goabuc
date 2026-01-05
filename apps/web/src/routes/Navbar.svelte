<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import NavbarLink from './NavbarLink.svelte';
</script>

<nav class="flex gap-4 bg-[#0257A5] p-4">
	<img src="../../static/resources/Think-different-Academy_LOGO_erb.png" alt="Logo TdA" /> <!--fsr to odmítá načíst ten obrázek, asi su kkt nebo idk-->
	<NavbarLink name="Home" url="/" />
	<NavbarLink name="Courses" url="/courses" />

	{#if auth.user}
		<NavbarLink name="Dashboard" url="/dashboard" />
		<NavbarLink name="Profile" url="/profile" />

		<button class="rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#000]" onclick={() => auth.logout()}
			>Logout</button
		>
	{:else}
		<NavbarLink name="Login" url="/login" />
		<NavbarLink name="Register" url="/register" />
	{/if}
</nav>
