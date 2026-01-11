<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import NavbarLink from './NavbarLink.svelte';
</script>

<nav class="flex items-center gap-4 bg-[#0257A5] p-4 text-3xl">
	<img
		src="../resources/Think-different-Academy_LOGO_bily.svg"
		alt="Logo TdA"
		width="100"
		class="m-[calc(100px/3)]"
	/>

	<NavbarLink name="Home" url="/" />
	<NavbarLink name="Courses" url="/courses" />

	{#if auth.user}
		<NavbarLink name="Dashboard" url="/dashboard" />
		<NavbarLink name="Profile" url="/profile" />

		<button
			class="flex h-[61px] cursor-pointer items-center rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#000]"
			onclick={() => auth.logout()}
			>Logout
		</button>
	{:else}
		<NavbarLink name="Login" url="/login" />
		<NavbarLink name="Register" url="/register" />
	{/if}
</nav>
